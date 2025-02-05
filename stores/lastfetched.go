package stores

import (
	"marina/constants"
	"marina/db"
	"marina/types"
	"time"
)

var lastFetchedStore = []*time.Time{}

func initializeLastFetchedStoreIfNecessary() {
	for len(lastFetchedStore) < len(constants.Repositories) {
		lastFetchedStore = append(lastFetchedStore, nil)
	}
}

func GetLastFetched(repo *marina.Repository) *time.Time {
	initializeLastFetchedStoreIfNecessary()

	if lastFetchedStore[repo.Id] == nil {
		lastFetchedStore[repo.Id] = db.GetLastUpdated(repo)
	}

	return lastFetchedStore[repo.Id]
}

func UpdateLastFetched(repo *marina.Repository, datetime time.Time) {
	initializeLastFetchedStoreIfNecessary()

	GetLastFetched(repo)

	if lastFetchedStore[repo.Id] == nil {
		db.AddLastUpdated(repo, datetime)
	} else {
		db.UpdateLastUpdated(repo, datetime)
	}

	lastFetchedStore[repo.Id] = &datetime
}
