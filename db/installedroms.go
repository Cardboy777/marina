package db

import (
	"fmt"
	"marina/types"
)

const getInstalledRomsScript = `
SELECT
		Hash,
		Name
FROM InstalledRoms 
WHERE Owner = ? AND Repository = ?
`

func GetInstalledRomList(repo *marina.Repository) []marina.Rom {
	rows, err := db.Query(getInstalledRomsScript, repo.Owner, repo.Repository)
	if err != nil {
		panic(fmt.Errorf("Error reading installed roms from manifest: %w", err))
	}

	roms := []marina.Rom{}

	for rows.Next() {
		var r marina.Rom

		if err := rows.Scan(
			&r.Sha1,
			&r.Name,
		); err != nil {
			panic(fmt.Errorf("Error reading installed roms from manifest: %w", err))
		}
		roms = append(roms, r)
	}

	rows.Close()
	return roms
}

const addRomsScript = `
INSERT INTO InstalledRoms (
	Owner,
	Repository,
	Hash,
	Name
)
VALUES (?, ?, ?, ?)
`

func AddInstalledRom(rom marina.Rom, repo *marina.Repository) {
	_, err := db.Exec(addRomsScript,
		repo.Owner,
		repo.Repository,
		rom.Sha1,
		rom.Name,
	)
	if err != nil {
		panic(fmt.Errorf("Error adding rom to manifest: %w", err))
	}
}
