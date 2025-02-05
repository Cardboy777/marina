package stores

import (
	"marina/constants"
	"marina/db"
	"marina/types"
)

var installedRoms = []*[]marina.Rom{}

func initializeRomStoreIfNeccessary() {
	for len(constants.Repositories) > len(installedRoms) {
		installedRoms = append(installedRoms, nil)
	}
}

func GetInstalledRomsList(repo *marina.Repository) *[]marina.Rom {
	initializeRomStoreIfNeccessary()

	if installedRoms[repo.Id] == nil {
		list := db.GetInstalledRomList(repo)
		installedRoms[repo.Id] = &list
	}

	return installedRoms[repo.Id]
}

func AddInstalledRom(rom marina.Rom, repo *marina.Repository) {
	initializeRomStoreIfNeccessary()

	GetInstalledRomsList(repo)

	isFound := false
	for _, r := range *installedRoms[repo.Id] {
		if r.Sha1 == rom.Sha1 {
			isFound = true
			break
		}
	}

	if !isFound {
		*installedRoms[repo.Id] = append(*installedRoms[repo.Id], rom)
		db.AddInstalledRom(rom, repo)
	}
}
