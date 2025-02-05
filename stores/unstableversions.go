package stores

import (
	"marina/constants"
	"marina/db"
	"marina/types"
)

var (
	unstableversions       = []*[]marina.UnstableVersion{}
	latestUnstableVersions = []*marina.UnstableVersion{}
)

func initializeUnstableVersionStoreIfNecessary() {
	for len(constants.Repositories) > len(unstableversions) {
		unstableversions = append(unstableversions, nil)
	}
	for len(constants.Repositories) > len(latestUnstableVersions) {
		latestUnstableVersions = append(latestUnstableVersions, nil)
	}
}

func GetUnstableVersions(repo *marina.Repository) *[]marina.UnstableVersion {
	initializeUnstableVersionStoreIfNecessary()

	if unstableversions[repo.Id] == nil {
		list := db.GetInstalledUnstableVersionList(repo)
		unstableversions[repo.Id] = &list

		var latest *marina.UnstableVersion
		for _, item := range list {
			if latest == nil || item.ReleaseDate.After(latest.ReleaseDate) {
				latest = &item
			}
		}

		latestUnstableVersions[repo.Id] = latest
	}
	return unstableversions[repo.Id]
}

func GetLatestUnstableVersion(repo *marina.Repository) *marina.UnstableVersion {
	initializeUnstableVersionStoreIfNecessary()

	return latestUnstableVersions[repo.Id]
}

func AddUnstableVersion(version *marina.UnstableVersion, installed bool) {
	initializeVersionStoreIfNecessary()

	GetUnstableVersions(version.Repository)

	isFound := false
	if *unstableversions[version.Repository.Id] == nil {
		*unstableversions[version.Repository.Id] = []marina.UnstableVersion{}
	}
	for _, v := range *unstableversions[version.Repository.Id] {
		if v.Hash == version.Hash {
			isFound = true
			break
		}
	}

	if !isFound {
		*unstableversions[version.Repository.Id] = append(*unstableversions[version.Repository.Id], *version)
	}
	if installed {
		db.AddInstalledUnstableVersion(version)
	}
}
