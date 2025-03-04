package db

import (
	"fmt"
	"marina/types"
	"sort"
	"time"
)

const getVersionsScript = `
SELECT
	TagName,
	Name,
	WindowsDownloadUrl,
	MacDownloadUrl,
	LinuxPerformanceDownloadUrl,
	LinuxCompatabilityDownloadUrl,
	ReleaseDate,
	Installed
FROM Releases
WHERE Owner = ? AND Repository = ?
ORDER BY ReleaseDate DESC
`

func GetVersionList(repo *marina.Repository) []marina.Version {
	rows, err := db.Query(getVersionsScript, repo.Owner, repo.Repository)
	if err != nil {
		panic(fmt.Errorf("Error reading releases from manifest: %w", err))
	}

	versions := []marina.Version{}

	for rows.Next() {
		var v marina.Version
		var timeString string

		if err := rows.Scan(
			&v.TagName,
			&v.Name,
			&v.DownloadUrls.Windows,
			&v.DownloadUrls.Mac,
			&v.DownloadUrls.Linux,
			&v.DownloadUrls.LinuxCompatibility,
			&timeString,
			&v.Installed,
		); err != nil {
			panic(fmt.Errorf("Error reading releases from manifest: %w", err))
		}
		v.Repository = repo
		time, err := time.Parse(time.Layout, timeString)
		if err != nil {
			panic(fmt.Errorf("Error reading releases from manifest: %w", err))
		}
		v.ReleaseDate = time
		versions = append(versions, v)
	}
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].ReleaseDate.After(versions[j].ReleaseDate)
	})

	rows.Close()

	return versions
}

const addRelease = `
INSERT
INTO Releases
(
	TagName,
	Owner,
	Repository,
	Name,
	WindowsDownloadUrl,
	MacDownloadUrl,
	LinuxPerformanceDownloadUrl,
	LinuxCompatabilityDownloadUrl,
	ReleaseDate,
	Installed
)
Values (?,?,?,?,?,?,?,?,?,?)
`

func AddRelease(v marina.Version) {
	_, err := db.Exec(addRelease,
		v.TagName,
		v.Repository.Owner,
		v.Repository.Repository,
		v.Name,
		v.DownloadUrls.Windows,
		v.DownloadUrls.Mac,
		v.DownloadUrls.Linux,
		v.DownloadUrls.LinuxCompatibility,
		v.ReleaseDate.Format(time.Layout),
		v.Installed,
	)
	if err != nil {
		panic(fmt.Errorf("Error inserting release into manifest: %w", err))
	}
}

const setAsDownloaded = `
UPDATE Releases
SET Installed = ?
WHERE 
	Owner = ? AND
	Repository = ? AND
	TagName = ?
`

func SetInstalled(v *marina.Version, isInstalled bool) {
	_, err := db.Exec(setAsDownloaded,
		isInstalled,
		v.Repository.Owner,
		v.Repository.Repository,
		v.TagName,
	)
	if err != nil {
		panic(fmt.Errorf("Error updating version in manifest: %w", err))
	}
}
