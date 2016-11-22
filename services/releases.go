package services

import (
	"github.com/pressly/warpdrive/data"
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

	release, err := data.FindReleaseByID(cycleID, releaseID)

	return release, err
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
