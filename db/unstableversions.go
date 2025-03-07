package db

import (
	"fmt"
	"marina/types"
	"time"
)

const getUnstableVersionsScript = `
SELECT
		CommitHash,
		ReleaseDate,
		Installed
FROM UnstableVersions
WHERE Owner = ? AND Repository = ?
ORDER BY ReleaseDate DESC
`

func GetUnstableVersionList(repo *marina.Repository) []marina.UnstableVersion {
	rows, err := db.Query(getUnstableVersionsScript, repo.Owner, repo.Repository)
	if err != nil {
		panic(fmt.Errorf("Error reading unstable versions from manifest: %w", err))
	}

	versions := []marina.UnstableVersion{}

	for rows.Next() {
		var v marina.UnstableVersion
		var timeString string

		if err := rows.Scan(
			&v.Hash,
			&timeString,
			&v.Installed,
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
	rows.Close()
	return versions
}

const addUnstableVersionscript = `
INSERT INTO UnstableVersions (
	Owner,
	Repository,
	CommitHash,
	ReleaseDate,
	Installed
) VALUES (?,?,?,?,?)
`

func AddUnstableVersion(devVersion *marina.UnstableVersion) {
	_, err := db.Exec(addUnstableVersionscript,
		devVersion.Repository.Owner,
		devVersion.Repository.Repository,
		devVersion.Hash,
		devVersion.ReleaseDate.Format(time.Layout),
		devVersion.Installed,
	)
	if err != nil {
		panic(fmt.Errorf("Error writing unstable versions to manifest: %w", err))
	}
}

const setUnstableVersionInstalledScript = `
UPDATE UnstableVersions
	SET Installed = ?
	WHERE Owner = ?
	AND Repository = ?
	AND CommitHash = ?
`

func SetUnstableVersionInstalled(devVersion *marina.UnstableVersion) {
	_, err := db.Exec(setUnstableVersionInstalledScript,
		devVersion.Installed,
		devVersion.Repository.Owner,
		devVersion.Repository.Repository,
		devVersion.Hash,
	)
	if err != nil {
		panic(fmt.Errorf("Error updating unstable versions to manifest: %w", err))
	}
}

const removeUnstableVersionscript = `
DELETE FROM UnstableVersions
	WHERE Owner = ?
	AND Repository = ?
	AND CommitHash = ?
`

func RemoveUnstableVersion(devVersion *marina.UnstableVersion) {
	_, err := db.Exec(removeUnstableVersionscript,
		devVersion.Repository.Owner,
		devVersion.Repository.Repository,
		devVersion.Hash,
	)
	if err != nil {
		panic(fmt.Errorf("Error writing unstable versions to manifest: %w", err))
	}
}

const deleteNotInstalledNotLatestUnstableRelease = `
DELETE FROM UnstableVersions
WHERE Owner = ?
	AND Repository = ?
	AND CommitHash != ?
	AND INSTALLED = 0
`

func RemoveOldNotInstalledUnstableVersions(repo *marina.Repository, latestVersion *marina.UnstableVersion) {
	_, err := db.Exec(deleteNotInstalledNotLatestUnstableRelease, repo.Owner, repo.Repository, latestVersion.Hash)
	if err != nil {
		panic(fmt.Errorf("Error removing old unstable versions: %w", err))
	}
}
