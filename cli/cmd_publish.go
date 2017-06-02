package cli

import (
	"context"
	"io"

	"fmt"
	"time"

	"github.com/pressly/warpdrive/helper"
	warpdrive "github.com/pressly/warpdrive/proto"
	"github.com/spf13/cobra"
)

var publishFlag = struct {
	caPath   string
	certPath string
	keyPath  string
	addr     string
	app      string
	platform string
	rollout  string
	version  string
	notes    string
	upgrade  string
	root     string
}{}

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish react-native to server",
	Long:  `publish react-native project for both ios and android`,
	Run: func(cmd *cobra.Command, args []string) {

		grpcConfig, err := helper.NewGrpcConfig(publishFlag.caPath, publishFlag.certPath, publishFlag.keyPath)
		if err != nil {
			logError(err)
		}

		clientConn, err := grpcConfig.CreateClient("command", publishFlag.addr)
		if err != nil {
			logError(err)
		}

		var platform warpdrive.Platform
		var bundlePath string
		switch publishFlag.platform {
		case "ios":
			platform = warpdrive.Platform_IOS
			bundlePath = bundlePathIOS
		case "android":
			platform = warpdrive.Platform_ANDROID
			bundlePath = bundlePathAndroid
		default:
			logError(fmt.Errorf("unknown platform '%s'", publishFlag.platform))
		}

		createdAt := time.Now().Format("Mon Jan 2 15:04:05 MST 2006")

		command := warpdrive.NewCommandClient(clientConn)

		ctx := context.Background()

		// First a release entity will be created
		newRelease := &warpdrive.Release{
			Id:            0,
			App:           publishFlag.app,
			Version:       publishFlag.version,
			Notes:         publishFlag.notes,
			Platform:      platform,
			NextReleaseId: 0,
			RolloutAt:     publishFlag.rollout,
			Bundle:        "",
			Lock:          false,
			CreatedAt:     createdAt,
			UpdatedAt:     createdAt,
		}

		ctx = context.Background()
		upload, err := command.UploadRelease(ctx)
		if err != nil {
			logError(err)
		}

		reader, writer := io.Pipe()
		go func() {
			// sending header
			err := upload.Send(&warpdrive.Chunck{
				Value: &warpdrive.Chunck_Header_{
					Header: &warpdrive.Chunck_Header{
						Release: newRelease,
						Root:    publishFlag.root,
						Upgrade: publishFlag.upgrade,
					},
				},
			})

			if err != nil {
				if err == io.EOF {
					_, err = upload.CloseAndRecv()
				}
				reader.CloseWithError(err)
				return
			}

			buffer := make([]byte, 1000)
			for {
				n, err := reader.Read(buffer)

				if err == io.EOF {
					reader.Close()
					break
				}

				if err != nil {
					reader.CloseWithError(err)
					return
				}

				err = upload.Send(&warpdrive.Chunck{
					Value: &warpdrive.Chunck_Body_{
						Body: &warpdrive.Chunck_Body{
							Data: buffer[:n],
						},
					},
				})

				if err != nil {
					_, err = upload.CloseAndRecv()
					reader.CloseWithError(err)
					return
				}
			}
		}()

		err = helper.BundleCompress(bundlePath, writer)
		if err != nil {
			logError(err)
		}

		release, err := upload.CloseAndRecv()
		if err != nil {
			logError(err)
		}

		fmt.Printf("new Release to %s\n\n\tPlatform: %s\n\tVersion: %s\n\tRollout: %s\n\n", release.App, release.Platform, release.Version, release.RolloutAt)
	},
}

func initPublishFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&publishFlag.caPath, "ca-path", "", "", "path to CA certificate file")
	cmd.Flags().StringVarP(&publishFlag.certPath, "cert-path", "", "", "path to certificate file")
	cmd.Flags().StringVarP(&publishFlag.keyPath, "key-path", "", "", "path to certificate's key file")
	cmd.Flags().StringVarP(&publishFlag.addr, "addr", "", "", "warpdrive server address. e.g. 127.0.0.1:1000")
	cmd.Flags().StringVarP(&publishFlag.app, "app", "a", "", "project app name")
	cmd.Flags().StringVarP(&publishFlag.platform, "platform", "p", "all", "target platform, `ios` or `android`")
	cmd.Flags().StringVarP(&publishFlag.rollout, "rollout", "r", "", "rollout cycle, could be beta, alpha, etc.")
	cmd.Flags().StringVarP(&publishFlag.version, "version", "v", "", "version of this bundle")
	cmd.Flags().StringVarP(&publishFlag.notes, "notes", "n", "", "release notes")
	cmd.Flags().StringVarP(&publishFlag.upgrade, "upgrade", "", "", "previous version which can upgrade to this version")
	cmd.Flags().StringVarP(&publishFlag.root, "root", "", "", "the very first version of this chain")
}

func init() {
	initPublishFlags(publishCmd)
	RootCmd.AddCommand(publishCmd)
}
