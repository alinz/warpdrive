package cli

import (
	"fmt"
	"os"
)

const (
	bundlePathIOS     = "./.bundles/ios"
	bundlePathAndroid = "./.bundles/android"
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

func bundleReadyPath(platform string) (string, error) {
	var path string

	switch platform {
	case "ios":
		path = bundlePathIOS
	case "android":
		path = bundlePathAndroid
	default:
		return "", fmt.Errorf("platform '%s' unknown\n", platform)
	}

	exists, _ := pathExists(path)
	if !exists {
		return "", fmt.Errorf("%s bundle not found\n", platform)
	}

	return path, nil
}
