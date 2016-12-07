package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func isReactNativeProject() bool {
	paths := []string{
		"./android/app/src/main",
		"./android/app/src/main/assets",
		"./ios",
		"./package.json",
	}

	for _, path := range paths {
		if exists, _ := pathExists(path); !exists {
			return false
		}
	}

	return true
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "creates a temporary WarpFile",
	Long:  `check if the project is react-native and creates proper WarpFile`,
	Run: func(cmd *cobra.Command, args []string) {
		// check if the current path is a react-native project
		// can be done very quickly by checking if the following path exists in current
		// directory
		if isReactNativeProject() {
			fmt.Println("The current path is not a react-native project")
		}

		// need to check if WarpFile already exists in this project
		// if so, then terminate the init process
		paths := []string{
			"./android/app/src/main/assets/WarpFile",
			"./ios/WarpFile",
		}

		for _, path := range paths {
			if exists, _ := pathExists(path); exists {
				fmt.Println("project was initialized before, nothing needs to be done")
				return
			}
		}

		// creates a default WarpFile
		config := localConfig{}
		err := config.Save()
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
