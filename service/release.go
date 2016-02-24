package service

import (
	"errors"
	"os"

	"upper.io/db.v2"

	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/lib/crypto"
)

func CreateRelease(
	appID,
	cycleID,
	userID int64,
	platformStr,
	versionStr,
	note string,
	filenames []string,
	filepaths []string) (*data.Release, error) {
	//check the permission
	if !HasPermissionToCreateAppCycle(appID, userID) {
		return nil, errors.New("No access to this app")
	}

	release := data.Release{
		CycleID: cycleID,
		Note:    note,
	}

	if platform, err := data.ParsePlatform(platformStr); err == nil {
		release.Platform = platform
	} else {
		return nil, err
	}

	if version, err := data.ParseVersion(versionStr); err == nil {
		release.Version = version
	} else {
		return nil, err
	}

	release.Locked = false

	//this variable is used to track all the bundles
	//so in case of error, we can easily remove them
	var bundlepaths []string

	//this is a transaction operation to add release record and
	//all bundle files.
	fn := func(session db.Database) error {
		var err error

		//create release
		if err = release.Save(session); err != nil {
			return err
		}

		//loop through all uploaded file and create bundle record.
		for index, filepath := range filepaths {

			//hash the temprary file
			hash, err := crypto.HashFile(filepath)
			if err != nil {
				return err
			}

			//create a bundle record
			bundle := data.Bundle{
				ReleaseID: release.ID,
				Hash:      hash,
				Name:      filenames[index],
				Type:      data.JS,
			}

			err = bundle.Save(session)
			if err != nil {
				return err
			}

			bundlepath := warpdrive.Config.Bundle.BundlesFolder + hash

			//move the bundle file from temp to bundle folder
			err = os.Rename(filepath, bundlepath)

			if err != nil {
				return err
			}

			bundlepaths = append(bundlepaths, bundlepath)
		}

		return nil
	}

	if err := data.Transaction(fn); err != nil {
		//we need to remove all temporary files
		for _, filepath := range filepaths {
			os.Remove(filepath)
		}

		//since there is an error we need to clean up bundle folder
		for _, bundlepath := range bundlepaths {
			os.Remove(bundlepath)
		}

		return nil, err
	}

	return &release, nil
}

func AllCycleReleases(appID, cycleID, userID int64) ([]*data.Release, error) {
	//check the permission
	if !HashPermissionToAccessCycle(appID, cycleID, userID) {
		return nil, errors.New("No access to this app's cycle")
	}

	return data.FindAllReleases(warpdrive.DB, db.Cond{
		"cycle_id": cycleID,
	})
}

func LockRelease(appID, cycleID, userID, releaseID int64) error {
	if !HashPermissionToAccessCycle(appID, cycleID, userID) {
		return errors.New("No access to this app's cycle")
	}

	if HasReleaseLocked(releaseID) {
		return errors.New("release is locked already")
	}

	fn := func(session db.Database) error {
		var err error
		release := data.Release{}
		err = release.Find(session, db.Cond{"id": releaseID})

		if err != nil {
			return err
		}

		release.Locked = true

		err = release.Save(session)
		return err
	}

	err := data.Transaction(fn)

	return err
}

func CheckDownloadableVersion(
	appID,
	cycleID int64,
	platformStr string,
	versionStr string,
) (map[string]data.Version, error) {
	//convert version to int represantation
	//then increment the major side and put it into the sql
	currentVersion, err := data.ParseVersion(versionStr)

	if err != nil {
		return nil, err
	}

	platform, err := data.ParsePlatform(platformStr)

	if err != nil {
		return nil, err
	}

	versions := make(map[string]data.Version)

	//finding soft version

	nextMajorVersion := data.VersionAdd(currentVersion, 1, 0, 0)
	builder := warpdrive.DB.Builder()
	q := builder.
		Select(
			"releases.id",
			"releases.version",
		).
		From("releases").
		Join("cycles").
		On("cycles.id=releases.cycle_id").
		Join("apps").
		On("apps.id=cycles.app_id").
		Where(`
			releases.locked=TRUE AND
			apps.id=? AND releases.cycle_id=? AND releases.platform=? AND
			releases.version > ? AND releases.version < ?
		`, appID, cycleID, platform, currentVersion, nextMajorVersion).
		OrderBy("version")

	release := data.Release{}

	err = q.Iterator().One(&release)
	if err == nil {
		versions["soft"] = release.Version
	}

	//finding hard version

	nextMajorVersion = data.MaskVersion(nextMajorVersion, 1, 0, 0)

	q = builder.
		Select(
			"releases.id",
			"releases.version",
		).
		From("releases").
		Join("cycles").
		On("cycles.id=releases.cycle_id").
		Join("apps").
		On("apps.id=cycles.app_id").
		Where(`
			releases.locked=TRUE AND
			apps.id=? AND releases.cycle_id=? AND releases.platform=? AND
			releases.version >= ?
		`, appID, cycleID, platform, nextMajorVersion).
		OrderBy("version")

	release = data.Release{}

	err = q.Iterator().One(&release)
	if err == nil {
		versions["hard"] = release.Version
	}

	return versions, nil
}
