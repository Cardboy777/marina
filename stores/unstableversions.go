package stores

import (
	"marina/constants"
	"marina/db"
	"marina/types"
)

func GetUnstableVersions(repo *marina.Repository) *[]marina.UnstableVersion {
	list := db.GetUnstableVersionList(repo)
	return &list
}

func AddUnstableVersion(version *marina.UnstableVersion) {
	unstableVersions := GetUnstableVersions(version.Repository)

	found := false
	for _, v := range *unstableVersions {
		if v.Hash == version.Hash {
			found = true
			break
		}
	}

	if !found {
		db.AddUnstableVersion(version)
	}
}

func SetUnstableVersionInstalled(version *marina.UnstableVersion, installed bool) {
	version.Installed = installed
	db.SetUnstableVersionInstalled(version)
}

func RemoveOldUnstableVersions() {
	var latest *marina.UnstableVersion

	for _, repo := range constants.Repositories {
		versions := GetUnstableVersions(repo)

		latest = findLatestVersion(versions)
		if latest != nil {
			db.RemoveOldNotInstalledUnstableVersions(repo, latest)
		}

	}
}

func findLatestVersion(list *[]marina.UnstableVersion) (latest *marina.UnstableVersion) {
	if list == nil {
		return nil
	}
	for _, item := range *list {
		if latest == nil || latest.ReleaseDate.Before(item.ReleaseDate) {
			latest = &item
		}
	}

	return latest
}
