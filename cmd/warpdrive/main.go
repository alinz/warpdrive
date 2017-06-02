package warpdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pressly/warpdrive/helper"
	pb "github.com/pressly/warpdrive/proto"
)

var config = struct {
	bundlePath     string
	documentPath   string
	warpdrivePath  string
	platform       pb.Platform
	app            string
	rollout        string
	bundleVersion  string
	currentVersion string
	addr           string
	grpcClient     *helper.GrpcConfig
}{}

func update() (*pb.Release, error) {
	clientConn, err := config.grpcClient.CreateClient("query", config.addr)
	if err != nil {
		return nil, err
	}

	query := pb.NewQueryClient(clientConn)
	release, err := query.GetUpgrade(context.Background(), &pb.Release{
		App:       config.app,
		Version:   config.currentVersion,
		RolloutAt: config.rollout,
		Platform:  config.platform,
	})
	if err != nil {
		return nil, err
	}

	return release, nil
}

func upgrade(release *pb.Release) error {
	log.Println("upgrading to", release.Version)

	err := download(release)
	if err != nil {
		return err
	}

	err = saveCurrentVersion(release)
	if err != nil {
		return err
	}

	return nil
}

func download(release *pb.Release) error {
	clientConn, err := config.grpcClient.CreateClient("query", config.addr)
	if err != nil {
		return err
	}

	query := pb.NewQueryClient(clientConn)

	stream, err := query.DownloadRelease(context.Background(), release)
	if err != nil {
		return err
	}

	// need to read the header first
	chunck, err := stream.Recv()
	if err != nil {
		return err
	}

	header := chunck.GetHeader()
	if header == nil {
		return fmt.Errorf("header is not sent")
	}

	if header.Release.Id != release.Id {
		return fmt.Errorf("release id mismatched")
	}

	reader, writer := io.Pipe()

	// this go-routine reads the bytes from grpc services and pump it to
	// io.Pipe. In this way, we don't need to save the tar file and we can
	// extract the data from stream of bytes.
	go func() {
		for {
			chunck, err := stream.Recv()
			if err == io.EOF {
				writer.Close()
				break
			}

			if err != nil {
				writer.CloseWithError(err)
				return
			}

			body := chunck.GetBody()
			if body == nil {
				writer.CloseWithError(fmt.Errorf("body is empty"))
				return
			}

			// ###1
			_, err = writer.Write(body.Data)
			if err != nil {
				return
			}
		}
	}()

	dstPath := releasePath(release)

	return helper.BundleUncompress(reader, dstPath)
}

func saveCurrentVersion(release *pb.Release) error {
	file, err := os.Create(currentVersionPath())
	if err != nil {
		return err
	}

	defer file.Close()

	return json.NewEncoder(file).Encode(release)
}

func loadCurrentVersion() (*pb.Release, error) {
	file, err := os.Open(currentVersionPath())
	if err != nil {
		return nil, err
	}

	defer file.Close()

	release := pb.Release{}
	err = json.NewDecoder(file).Decode(&release)
	if err != nil {
		return nil, err
	}

	return &release, nil
}

func releasePath(release *pb.Release) string {
	return releasePathByVersion(release.Version)
}

func releasePathByVersion(version string) string {
	version = strings.Replace(version, "/", "-", -1)
	version = strings.Replace(version, " ", "-", -1)
	return filepath.Join(config.warpdrivePath, fmt.Sprintf("releases/%s", version))
}

func currentVersionPath() string {
	return filepath.Join(config.warpdrivePath, "current.warpdrive")
}

func pathNormalize(path string) string {
	if strings.HasPrefix(path, "file://") {
		return strings.Replace(path, "file://", "", 1)
	}
	return path
}

// Init initialize warpdrive
func Init(bundlePath, documentPath, platform, app, rollout, bundleVersion, addr, deviceCert, deviceKey, caCert string) error {
	var err error

	config.bundlePath = pathNormalize(bundlePath)
	config.documentPath = pathNormalize(documentPath)
	switch platform {
	case "ios":
		config.platform = pb.Platform_IOS
	case "android":
		config.platform = pb.Platform_ANDROID
	default:
		config.platform = pb.Platform_UNKNOWN
	}

	config.app = app
	config.rollout = rollout
	config.bundleVersion = bundleVersion
	config.addr = addr
	config.grpcClient, err = helper.NewGrpcConfig(pathNormalize(caCert), pathNormalize(deviceCert), pathNormalize(deviceKey))
	if err != nil {
		return err
	}

	// NOTE: we are using bundleVersion for active warpdrive folder
	// now if a new version of app publishes to app store, the app starts fresh if the user updates
	// the app and read from the new fresh folder.
	config.warpdrivePath = filepath.Join(config.documentPath, "/warpdrive", bundleVersion)
	err = os.MkdirAll(config.warpdrivePath, os.ModePerm)
	if err != nil {
		return err
	}

	// set the current version based on previously saved one or
	// bundleVersion
	release, err := loadCurrentVersion()
	if err != nil {
		config.currentVersion = bundleVersion
	} else {
		config.currentVersion = release.Version
	}

	log.Println("current version:", config.currentVersion)

	// when user launches app, the app will pauses until
	// a new update downloads completely.
	release, err = update()
	if err == nil {
		err = upgrade(release)
		if err == nil {
			config.currentVersion = release.Version
			log.Println("upgraded to", release.Version)
		} else {
			log.Println(err.Error())
		}
	}

	return nil
}

// BundlePath returns the path to bundles which needs to be loaded
// it returns empty string if there is no update version available
func BundlePath() string {
	if config.bundleVersion == config.currentVersion {
		return ""
	}

	return filepath.Join(releasePathByVersion(config.currentVersion), "main.jsbundle")
}
