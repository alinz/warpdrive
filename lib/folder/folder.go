package folder

import (
	"os"
	"path/filepath"
)

// ListFilePaths returns the list of all path under the
// given path. It also finds all files in nested path and return the
// full path for each file.
func ListFilePaths(path string) ([]string, error) {
	var files []string

	// we need to define this loop to make it available inside loop itself
	var loop func(string, *[]string) error

	loop = func(path string, files *[]string) error {
		dir, err := os.Open(path)
		if err != nil {
			return err
		}

		defer dir.Close()

		dirStat, err := dir.Stat()
		if err != nil {
			return err
		}

		if !dirStat.IsDir() {
			*files = append(*files, path)
			return nil
		}

		fileInfos, err := dir.Readdir(-1)
		if err != nil {
			return err
		}

		for _, fileInfo := range fileInfos {
			err = loop(filepath.Join(path, fileInfo.Name()), files)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := loop(path, &files)

	return files, err
}
