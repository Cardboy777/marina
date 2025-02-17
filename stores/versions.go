package stores

import (
	"marina/db"
	"marina/types"
)

func GetVersions(repo *marina.Repository) *[]marina.Version {
	list := db.GetVersionList(repo)
	return &list
}

func SetVersionInstalled(version *marina.Version, installed bool) {
	version.Installed = installed
	db.SetInstalled(version, version.Installed)
}

func AddVersion(version *marina.Version) {
	versions := GetVersions(version.Repository)

	isFound := false
	for _, v := range *versions {
		if v.TagName == version.TagName {
			isFound = true
			break
		}
	}

	if !isFound {
		db.AddRelease(*version)
	}
}
