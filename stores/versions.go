package stores

import (
	"marina/constants"
	"marina/db"
	"marina/types"
)

var (
	versions       = []*[]marina.Version{}
	latestVersions = []*marina.Version{}
)

func initializeVersionStoreIfNecessary() {
	for len(constants.Repositories) > len(versions) {
		versions = append(versions, nil)
	}
	for len(constants.Repositories) > len(latestVersions) {
		latestVersions = append(latestVersions, nil)
	}
}

func GetVersions(repo *marina.Repository) *[]marina.Version {
	initializeVersionStoreIfNecessary()

	if versions[repo.Id] == nil {
		list := db.GetVersionList(repo)
		versions[repo.Id] = &list

		var latest *marina.Version
		for _, item := range list {
			if latest == nil || item.ReleaseDate.After(latest.ReleaseDate) {
				latest = &item
			}
		}

		latestVersions[repo.Id] = latest
	}

	return versions[repo.Id]
}

func SetVersionInstalled(version *marina.Version) {
	initializeVersionStoreIfNecessary()

	version.Installed = true

	db.SetInstalled(version, version.Installed)
}

func GetLatestVersion(repo *marina.Repository) *marina.Version {
	initializeVersionStoreIfNecessary()

	return latestVersions[repo.Id]
}

func AddVersion(version *marina.Version) {
	initializeVersionStoreIfNecessary()

	_ = GetVersions(version.Repository)

	isFound := false
	for _, v := range *versions[version.Repository.Id] {
		if v.TagName == version.TagName {
			isFound = true
			break
		}
	}

	if !isFound {
		*versions[version.Repository.Id] = append(*versions[version.Repository.Id], *version)
		db.AddRelease(*version)
	}
}
