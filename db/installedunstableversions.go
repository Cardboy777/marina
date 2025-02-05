package db

import (
	"fmt"
	"marina/types"
	"time"
)

const getInstalledUnstableVersionsScript = `
SELECT
		CommitHash,
		ReleaseDate
FROM InstalledUnstableVersions
WHERE Owner = ? AND Repository = ?
`

func GetInstalledUnstableVersionList(repo *marina.Repository) []marina.UnstableVersion {
	rows, err := db.Query(getInstalledUnstableVersionsScript, repo.Owner, repo.Name)
	if err != nil {
		panic(fmt.Errorf("Error reading unstable versions from manifest: %w", err))
	}
	defer rows.Close()

	versions := []marina.UnstableVersion{}

	for rows.Next() {
		var v marina.UnstableVersion
		var timeString string

		if err := rows.Scan(
			&v.Hash,
			&timeString,
		); err != nil {
			panic(fmt.Errorf("Error reading unstable versions from manifest: %w", err))
		}
		time, err := time.Parse(time.Layout, timeString)
		if err != nil {
			panic(fmt.Errorf("Error reading unstable versions from manifest: %w", err))
		}

		v.ReleaseDate = time
		v.Repository = repo
		versions = append(versions, v)
	}
	return versions
}

const addInstalledUnstableVersionScript = `
INSERT INTO InstalledUnstableVersions (
	Owner,
	Repository,
	CommitHash,
	ReleaseDate
) VALUES (?,?,?,?)
`

func AddInstalledUnstableVersion(devVersion *marina.UnstableVersion) {
	_, err := db.Exec(addInstalledUnstableVersionScript,
		devVersion.Repository.Owner,
		devVersion.Repository.Repository,
		devVersion.Hash,
		devVersion.ReleaseDate.Format(time.Layout),
	)
	if err != nil {
		panic(fmt.Errorf("Error writing unstable versions to manifest: %w", err))
	}
}
