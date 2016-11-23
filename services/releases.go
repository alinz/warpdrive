package services

import (
	"bytes"
	"io"
	"path/filepath"

	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/lib/crypto"
	"github.com/pressly/warpdrive/lib/warp"
	"upper.io/db.v2/lib/sqlbuilder"
)

func SearchReleases(userID, appID, cycleID int64, platform, version, note string) ([]*data.Release, error) {
	plat, err := data.ParsePlatform(platform)
	if err != nil {
		return nil, err
	}

	vers, err := data.ParseVersion(version)
	if err != nil {
		return nil, err
	}

	_, err = FindCycleByID(userID, appID, cycleID)
	if err != nil {
		return nil, err
	}

	releases, err := data.FindReleases(cycleID, plat, vers, note)
	if err != nil {
		return nil, err
	}

	if releases == nil {
		releases = make([]*data.Release, 0)
	}

	return releases, nil
}

func FindReleaseByID(userID, appID, cycleID, releaseID int64) (*data.Release, error) {
	_, err := FindCycleByID(userID, appID, cycleID)
	if err != nil {
		return nil, err
	}

	return data.FindReleaseByID(cycleID, releaseID)
}

func FindLockedRelease(cycleID, releaseID int64) (*data.Release, error) {
	return data.FindLockedReleaseByID(cycleID, releaseID)
}

func CreateRelease(userID, appID, cycleID int64, platform data.Platform, version data.Version, note string) (*data.Release, error) {
	_, err := FindCycleByID(userID, appID, cycleID)
	if err != nil {
		return nil, err
	}

	release := &data.Release{
		CycleID:  cycleID,
		Platform: platform,
		Version:  version,
		Note:     note,
	}

	err = data.Transaction(func(session sqlbuilder.Tx) error {
		return release.Save(session)
	})

	if err != nil {
		return nil, err
	}

	return release, nil
}

func UpdateRelease(userID, appID, cycleID, releaseID int64, platfrom *data.Platform, version *data.Version, note *string) (*data.Release, error) {
	_, err := FindCycleByID(userID, appID, cycleID)
	if err != nil {
		return nil, err
	}

	release := &data.Release{ID: releaseID}

	err = data.Transaction(func(session sqlbuilder.Tx) error {
		err := release.Load(session)
		if err != nil {
			return err
		}

		// we need to make sure that if released is locked, you can not edit it
		if release.Locked {
			return ErrReleaseLocked
		}

		if platfrom != nil {
			release.Platform = *platfrom
		}

		if version != nil {
			release.Version = *version
		}

		if note != nil {
			release.Note = *note
		}

		return release.Save(session)
	})

	if err != nil {
		return nil, err
	}

	return release, nil
}

func RemoveRelease(userID, appID, cycleID, releaseID int64) error {
	_, err := FindCycleByID(userID, appID, cycleID)
	if err != nil {
		return err
	}

	err = data.Transaction(func(session sqlbuilder.Tx) error {
		release := &data.Release{
			ID: releaseID,
		}

		err := release.Load(session)
		if err != nil {
			return err
		}

		// we need to make sure that if released is locked, you can not delete it
		if release.Locked {
			return ErrReleaseLocked
		}

		return release.Remove(session)
	})

	return nil
}

func LockRelease(userID, appID, cycleID, releaseID int64) error {
	_, err := FindCycleByID(userID, appID, cycleID)
	if err != nil {
		return err
	}

	release := &data.Release{ID: releaseID}

	err = data.Transaction(func(session sqlbuilder.Tx) error {
		err := release.Load(session)
		if err != nil {
			return err
		}

		if release.Locked {
			return ErrReleaseAlreadyLocked
		}

		release.Locked = true

		return release.Save(session)
	})

	return nil
}

func UnlockRelease(userID, appID, cycleID, releaseID int64) error {
	_, err := FindCycleByID(userID, appID, cycleID)
	if err != nil {
		return err
	}

	release := &data.Release{ID: releaseID}

	err = data.Transaction(func(session sqlbuilder.Tx) error {
		err := release.Load(session)
		if err != nil {
			return err
		}

		if !release.Locked {
			return ErrReleaseAlreadyUnlocked
		}

		release.Locked = false

		return release.Save(session)
	})

	return nil
}

func LatestRelease(appID, cycleID int64, version string, platform string) (map[string]*data.Release, error) {
	// check if cycle id belongs to app id
	_, err := FindCycleByAppIdCycleId(appID, cycleID)
	if err != nil {
		return nil, err
	}

	vers, err := data.ParseVersion(version)
	if err != nil {
		return nil, err
	}

	plat, err := data.ParsePlatform(platform)
	if err != nil {
		return nil, err
	}

	releases := make(map[string]*data.Release)

	softRelease, err := data.FindLatestSoftRelease(cycleID, plat, vers)
	if err != nil {
		releases["soft"] = softRelease
	}

	hardRelease, err := data.FindLatestHardRelease(cycleID, plat, vers)
	if err != nil {
		releases["hard"] = hardRelease
	}

	return releases, nil
}

func DownloadRelease(appID, cycleID, releaseID int64, encryptedKey []byte) (io.Reader, error) {
	// checks if cycle id belongs to app id and also we need cycle object
	// for PrivateKey
	cycle, err := FindCycleByAppIdCycleId(appID, cycleID)
	if err != nil {
		return nil, err
	}

	// checks if release id belongs to cycle if
	_, err = FindLockedRelease(cycleID, releaseID)
	if err != nil {
		return nil, err
	}

	// gets all bundles for a specific release id
	bundles, err := FindAllBundlesByReleaseID(releaseID)
	if err != nil {
		return nil, err
	}

	// decrypts the key by using private key
	key, err := crypto.DecryptByPrivateRSA(encryptedKey, cycle.PrivateKey, "sha1")
	if err != nil {
		return nil, err
	}

	// creates buffer for warpFile.
	var buffer bytes.Buffer
	warpFile := warp.NewWriter(&buffer)

	// goes over all bundles and add them to warpFile
	for _, bundle := range bundles {
		path := filepath.Join(warpdrive.Conf.Server.BundlesFolder, bundle.Hash)
		warpFile.AddFile(bundle.Name, path)
	}

	// we are encrypting the warpFile with key that we got from client
	encrypted, err := crypto.AESEncrypt(buffer.Bytes(), key)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(encrypted), nil
}
