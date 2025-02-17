package services

import (
	"marina/stores"
	"marina/types"
	"time"
)

func SyncReleases(repository *marina.Repository, force bool) error {
	if !force {
		timeLastFetched := stores.GetLastFetched(repository)
		timeCutoff := time.Now().Add(time.Duration(-1) * time.Hour)
		if timeLastFetched != nil && timeLastFetched.After(timeCutoff) {
			return nil
		}
	}

	versions, err := fetchReleaseVersions(repository)
	if err != nil {
		return err
	}
	for _, v := range versions {
		stores.AddVersion(&v)
	}

	latestCommit, err := fetchLatestCommit(repository)
	if err != nil {
		return err
	}

	stores.AddUnstableVersion(latestCommit)
	stores.RemoveOldUnstableVersions()

	stores.UpdateLastFetched(repository, time.Now())

	return nil
}
