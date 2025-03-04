package db

import (
	"database/sql"
	"errors"
	"fmt"
	"marina/types"
	"time"
)

const getLastUpdatedQuery = `
	SELECT Timestamp
	FROM LastFetchedDate
	WHERE Owner = ? AND Repository = ?
`

func GetLastUpdated(repo *marina.Repository) *time.Time {
	row := db.QueryRow(getLastUpdatedQuery, repo.Owner, repo.Repository)

	var timeString string

	err := row.Scan(&timeString)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		panic(fmt.Errorf("Error reading last fetched date from manifest: %w", err))
	}

	time, err := time.Parse(time.Layout, timeString)
	if err != nil {
		panic(fmt.Errorf("Error reading last fetched date from manifest: %w", err))
	}

	return &time
}

const addLastUpdatedQuery = `
	INSERT INTO LastFetchedDate
	(
		Owner,
		Repository,
		Timestamp
	)
	Values (?, ?, ?)
`

func AddLastUpdated(repo *marina.Repository, datetime time.Time) {
	_, err := db.Exec(addLastUpdatedQuery, repo.Owner, repo.Repository, datetime.Format(time.Layout))
	if err != nil {
		panic(fmt.Errorf("Error adding last fetched date from manifest: %w", err))
	}
}

const updateLastUpdatedQuery = `
	UPDATE LastFetchedDate
	SET Timestamp = ?
	WHERE Owner = ? AND Repository = ?
`

func UpdateLastUpdated(repo *marina.Repository, datetime time.Time) {
	_, err := db.Exec(updateLastUpdatedQuery, datetime.Format(time.Layout), repo.Owner, repo.Repository)
	if err != nil {
		panic(fmt.Errorf("Error updating last fetched date from manifest: %w", err))
	}
}
