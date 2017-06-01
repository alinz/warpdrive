package cli

import (
	"context"
	"io"
	"log"

	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pressly/warpdrive/helper"
	warpdrive "github.com/pressly/warpdrive/proto"
	"github.com/spf13/cobra"
)

var publishFlag = struct {
	app      string
	platform string
	rollout  string
	version  string
	notes    string
	upgrades []string
}{}

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish react-native to server",
	Long:  `publish react-native project for both ios and android`,
	Run: func(cmd *cobra.Command, args []string) {
		commandEnv := &struct {
			CA   string `require:"true"`
			Crt  string `require:"true"`
			Key  string `require:"true"`
			Addr string `require:"true"`
		}{}

		err := envconfig.Process("command", commandEnv)
		if err != nil {
			log.Fatal(err.Error())
		}

		grpcConfig, err := helper.NewGrpcConfig(commandEnv.CA, commandEnv.Crt, commandEnv.Key)
		if err != nil {
			log.Fatal(err.Error())
		}

		clientConn, err := grpcConfig.CreateClient("command", commandEnv.Addr)
		if err != nil {
			log.Fatal(err.Error())
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
			log.Fatal(fmt.Errorf("unknown platform '%s'", publishFlag.platform))
		}

		createdAt := time.Now().Format("Mon Jan 2 15:04:05 MST 2006")

		command := warpdrive.NewCommandClient(clientConn)

		ctx := context.Background()

		// First a release entity will be created
		release, err := command.CreateRelease(ctx, &warpdrive.Release{
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
		})

		if err != nil {
			log.Fatal(err.Error())
		}

		ctx = context.Background()
		upload, err := command.UploadRelease(ctx)
		if err != nil {
			log.Fatal(err.Error())
		}

		reader, writer := io.Pipe()
		go func() {
			// sending header
			err := upload.Send(&warpdrive.Chunck{
				Value: &warpdrive.Chunck_Header_{
					Header: &warpdrive.Chunck_Header{
						ReleaseId: release.Id,
						Upgrades:  publishFlag.upgrades,
					},
				},
			})

			if err != nil {
				log.Println("error sending header", err.Error())
				reader.CloseWithError(err)
				return
			}

			buffer := make([]byte, 1000)
			for {
				n, err := reader.Read(buffer)

				if err == io.EOF {
					log.Println("io closed")
					reader.Close()
					break
				}

				if err != nil {
					log.Println(err.Error())
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
					log.Println("close during sending data", err.Error())
					reader.CloseWithError(err)
					return
				}
			}
		}()

		err = helper.BundleCompress(bundlePath, writer)
		if err != nil {
			log.Fatal(err.Error())
		}

		release, err = upload.CloseAndRecv()
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println(release)
	},
}

func initPublishFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&publishFlag.app, "app", "a", "", "project app name")
	cmd.Flags().StringVarP(&publishFlag.platform, "platform", "p", "all", "target platform, `ios` or `android`")
	cmd.Flags().StringVarP(&publishFlag.rollout, "rollout", "r", "", "rollout cycle, could be beta, alpha, etc.")
	cmd.Flags().StringVarP(&publishFlag.version, "version", "v", "", "version of this bundle")
	cmd.Flags().StringVarP(&publishFlag.notes, "notes", "n", "", "release notes")
	cmd.Flags().StringArrayVarP(&publishFlag.upgrades, "upgrades", "u", nil, "comma seperate versions which can be upgrade to this version")
}

func init() {
	initPublishFlags(publishCmd)
	RootCmd.AddCommand(publishCmd)
}
