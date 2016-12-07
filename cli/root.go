package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "warp",
	Short: "In-App upgrade service for React-Native! Supporting iOS and Android apps",
	Long: `A Fast and Flexible upgrade service for React-Native apps!
                love by alinz and Pressly Inc.
                Complete documentation is available at https://pressly.github.io/warpdrive`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hahaha")
		// Do Stuff Here
	},
}

func init() {

}
