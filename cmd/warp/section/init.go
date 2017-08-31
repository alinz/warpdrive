package section

import (
	"context"
	"fmt"

	pb "github.com/pressly/warpdrive/proto"
	"github.com/spf13/cobra"
)

func init() {
	var server string
	var user string
	var pass string
	var app string

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "create a new app or load the configuration for an exisiting one",
		Long: `It either create or load the configuration for given app name.
Usually call this per project.`,
		Run: func(cmd *cobra.Command, args []string) {
			if !stringHasValue(&server, &user, &pass, &app) {
				fmt.Println("please run 'warp init -h' for usage")
			}

			conn, err := grpcConnection(server, "")
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			defer conn.Close()

			client := pb.NewWarpdriveClient(conn)

			certificate, err := client.SetupApp(context.Background(), &pb.Credential{
				Username: user,
				Password: pass,
				AppName:  app,
			})

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			err = saveAdminCertificate(certificate)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		},
	}

	initCmd.Flags().StringVarP(&server, "server", "s", "", "server address")
	initCmd.Flags().StringVarP(&user, "user", "u", "", "admin username")
	initCmd.Flags().StringVarP(&pass, "pass", "p", "", "admin password")
	initCmd.Flags().StringVarP(&app, "app", "a", "", "app name either new or already exisiting one")

	root.AddCommand(initCmd)
}
