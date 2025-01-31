package main

import (
	"marina/files"
	"marina/settings"
	"marina/ui"
)

func main() {
	settings.ConfigInit()
	files.Init()
	ui.Init()
}
