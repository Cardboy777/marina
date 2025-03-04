package stores

import (
	"marina/db"
	"marina/types"
)

func GetInstalledRomsList(repo *marina.Repository) *[]marina.Rom {
	list := db.GetInstalledRomList(repo)
	return &list
}

func AddInstalledRom(rom marina.Rom, repo *marina.Repository) {
	list := GetInstalledRomsList(repo)

	isFound := false
	for _, r := range *list {
		if r.Sha1 == rom.Sha1 {
			isFound = true
			break
		}
	}

	if !isFound {
		db.AddInstalledRom(rom, repo)
	}
}
