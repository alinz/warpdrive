package services

import (
	"io"
	"path/filepath"

	"strings"

	"fmt"

	"github.com/blang/semver"
	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/lib/crypto"
	"github.com/pressly/warpdrive/lib/folder"
	"github.com/pressly/warpdrive/lib/warp"
	"upper.io/db.v2/lib/sqlbuilder"
)

func SearchReleases(userID, appID, cycleID int64, platform, version, note string) ([]*data.Release, error) {
	plat, err := data.ParsePlatform(platform)
	if err != nil {
		return nil, err
	}

	_, err = semver.Make(version)
	if err != nil {
		return nil, err
	}

	_, err = FindCycleByID(userID, appID, cycleID)
	if err != nil {
		return nil, err
	}

	releases, err := data.FindReleases(cycleID, plat, version, note)
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

func CreateRelease(userID, appID, cycleID int64, platform data.Platform, rawVersion, note string) (*data.Release, error) {
	version, err := semver.Make(rawVersion)
	if err != nil {
		return nil, err
	}

	_, err = FindCycleByID(userID, appID, cycleID)
	if err != nil {
		return nil, err
	}

	release := &data.Release{
		CycleID:  cycleID,
		Platform: platform,
		Version:  rawVersion,
		Major:    int64(version.Major),
		Minor:    int64(version.Minor),
		Patch:    int64(version.Patch),
		Build:    strings.Join(version.Build, "."),
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

func UpdateRelease(userID, appID, cycleID, releaseID int64, platfrom *data.Platform, rawVersion *string, note *string) (*data.Release, error) {
	version, err := semver.Make(*rawVersion)
	if err != nil {
		return nil, err
	}

	_, err = FindCycleByID(userID, appID, cycleID)
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

		if rawVersion != nil {
			release.Version = *rawVersion
			release.Major = int64(version.Major)
			release.Minor = int64(version.Minor)
			release.Patch = int64(version.Patch)
			release.Build = strings.Join(version.Build, ".")
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

// LatestRelease retusn soft and hard version which soft means can be updated and hard means it
// requires to download from app store or play store
func LatestRelease(appID, cycleID int64, rawVersion string, platform string) (map[string]*data.Release, error) {
	// check if cycle id belongs to app id
	_, err := FindCycleByAppIdCycleId(appID, cycleID)
	if err != nil {
		return nil, err
	}

	version, err := semver.Make(rawVersion)
	if err != nil {
		return nil, err
	}

	plat, err := data.ParsePlatform(platform)
	if err != nil {
		return nil, err
	}

	releases := make(map[string]*data.Release)

	softRelease, err := data.FindLatestSoftRelease(cycleID, plat, version)
	if err == nil {
		releases["soft"] = softRelease
	}

	hardRelease, err := data.FindLatestHardRelease(cycleID, plat, version)
	if err == nil {
		releases["hard"] = hardRelease
	}

	return releases, nil
}

// DownloadRelease downloads the bundle releated to a specific release
func DownloadRelease(appID, cycleID, releaseID int64, encryptedKey []byte, output io.Writer) error {
	// checks if cycle id belongs to app id and also we need cycle object
	// for PrivateKey
	cycle, err := FindCycleByAppIdCycleId(appID, cycleID)
	if err != nil {
		return err
	}

	// checks if release id belongs to cycle if
	_, err = FindLockedRelease(cycleID, releaseID)
	if err != nil {
		return err
	}

	// gets all bundles for a specific release id
	bundles, err := FindAllBundlesByReleaseID(releaseID)
	if err != nil {
		return err
	}

	// decrypts the key by using private key
	key, err := crypto.DecryptByPrivateRSA(encryptedKey, cycle.PrivateKey, "sha1")
	if err != nil {
		return err
	}

	files := make(map[string]string)
	for _, bundle := range bundles {
		path := filepath.Join(warpdrive.Conf.Server.BundlesFolder, bundle.Hash)
		// we need to check if the file exists
		if exists, _ := folder.PathExists(path); !exists {
			return fmt.Errorf("file '%s' not found", path)
		}
		files[path] = bundle.Name
	}

	return warp.Compress(files, crypto.NewAESEncryptWriter(key, output))
}
