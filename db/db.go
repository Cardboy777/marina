package db

import (
	"database/sql"
	"fmt"
	"marina/db/scripts"
	"marina/settings"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const dbFileName = "manifest.db"

var db *sql.DB

func Init() {
	installDir := settings.GetInstallDirName()

	dbFilePath := fmt.Sprintf("%s?_busy_timeout=5000", filepath.Join(installDir, dbFileName))
	database, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		panic(fmt.Errorf("Error opening manifest db: %w", err))
	}
	db = database

	_, err = db.Exec(scripts.DbSchemaInit)
	if err != nil {
		panic(fmt.Errorf("Error initializing manifest db: %w", err))
	}
}
