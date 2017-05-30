package warpdrive

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

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
	clientConn, err := config.grpcClient.CreateClient("warpdrive", config.addr)
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
	clientConn, err := config.grpcClient.CreateClient("warpdrive", config.addr)
	if err != nil {
		return err
	}

	query := pb.NewQueryClient(clientConn)

	stream, err := query.DownloadRelease(context.Background(), release)
	if err != nil {
		return err
	}

	return nil
}

func saveCurrentVersion(version string) error {
	path := filepath.Join(config.warpdrivePath, "current.warpdrive")
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	file.WriteString(version)
	file.Sync()

	return nil
}

func loadCurrentVersion() (string, error) {
	path := filepath.Join(config.warpdrivePath, "current.warpdrive")
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// Init initialize warpdrive
func Init(bundlePath, documentPath, platform, app, rollout, bundleVersion, addr, deviceCert, deviceKey, caCert string) error {
	var err error

	config.bundlePath = bundlePath
	config.documentPath = documentPath
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
	config.grpcClient, err = helper.NewGrpcConfig(caCert, deviceCert, deviceKey)
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
	config.currentVersion, err = loadCurrentVersion()
	if err != nil {
		config.currentVersion = bundleVersion
	}

	return nil
}
