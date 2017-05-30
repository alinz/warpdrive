package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	bundlePathIOS     = ".bundles/ios"
	bundlePathAndroid = ".bundles/android"
)

// RootCmd is the base command for cli
var RootCmd = &cobra.Command{
	Use:   "warp",
	Short: "In-App upgrade service for React-Native! Supporting iOS and Android apps",
	Long: `
A Fast and Flexible upgrade service for React-Native apps!
Created by Ali Najafizadeh (alinz) at Pressly Inc.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("please run 'warp -h' for usage")
	},
}

func init() {}
