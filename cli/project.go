package cli

import "fmt"
import "github.com/pressly/warpdrive/lib/folder"

const (
	bundlePathIOS     = ".bundles/ios"
	bundlePathAndroid = ".bundles/android"
)

func isReactNativeProject() bool {
	paths := []string{
		"./android/app/src/main",
		"./ios",
		"./package.json",
	}

	for _, path := range paths {
		if exists, _ := folder.PathExists(path); !exists {
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

	exists, _ := folder.PathExists(path)
	if !exists {
		return "", fmt.Errorf("%s bundle not found\n", platform)
	}

	return path, nil
}
