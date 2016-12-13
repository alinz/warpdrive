package cli

import (
	"fmt"
	"os"
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
		"./ios",
		"./package.json",
	}

	for _, path := range paths {
		if exists, _ := pathExists(path); !exists {
			fmt.Println("this is not react-native project")
			return false
		}
	}

	return true
}
