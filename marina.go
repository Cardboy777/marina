//go:generate goversioninfo build/windows/versioninfo.json -64
package main

import (
	_ "embed"
	"marina/assets"
	"marina/db"
	"marina/files"
	"marina/settings"
	"marina/ui"
)

func main() {
	settings.Init()
	files.Init()
	assets.LoadAssets()
	db.Init()
	for {
		ui.Start()
		if !settings.ShouldRestart() {
			break
		}
		db.ResetDbConnection()
	}
}
