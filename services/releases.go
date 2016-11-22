package services

import "github.com/pressly/warpdrive/data"

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
