package services

import (
	"os"

	"sync"

	"path/filepath"

	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/lib/crypto"
	"upper.io/db.v2/lib/sqlbuilder"
)

var lockCreateBundle sync.Mutex

func isFileExist(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// moveFile calculate the hash value and move the file from temp folder to
// bundles folder. because we hash the content of the file, duplicate files will be
// eliminated and save some space.
func moveFile(srcPath string) (string, error) {
	hash, err := crypto.HashFile(srcPath)
	if err != nil {
		return "", err
	}

	targetPath := filepath.Join(warpdrive.Conf.Server.BundlesFolder, hash)
	if isFileExist(targetPath) {
		os.Remove(srcPath)
		return hash, nil
	}

	err = os.Rename(srcPath, targetPath)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func CreateBundles(userID, appID, cycleID, releaseID int64, mapFiles map[string]string) ([]*data.Bundle, error) {
	_, err := FindReleaseByID(userID, appID, cycleID, releaseID)
	if err != nil {
		return nil, err
	}

	lockCreateBundle.Lock()
	defer lockCreateBundle.Unlock()

	var bundles []*data.Bundle

	for fileName, tempFile := range mapFiles {
		hash, err := moveFile(tempFile)
		if err != nil {
			return nil, err
		}

		bundles = append(bundles, &data.Bundle{
			ReleaseID: releaseID,
			Hash:      hash,
			Name:      fileName,
		})
	}

	err = data.Transaction(func(session sqlbuilder.Tx) error {
		var err error
		for _, bundle := range bundles {
			err = bundle.Save(session)
			if err != nil {
				return nil
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return bundles, nil
}
