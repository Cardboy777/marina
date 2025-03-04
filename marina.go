package main

import (
	"marina/db"
	"marina/files"
	"marina/settings"
	"marina/ui"
)

func main() {
	settings.Init()
	files.Init()
	db.Init()
	for {
		ui.Start()
		if !settings.ShouldRestart() {
			break
		}
		db.ResetDbConnection()
	}
}
