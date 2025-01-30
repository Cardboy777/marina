package main

import (
	"marina/filemanager"
	"marina/settings"
	"marina/ui"
)

func main() {
	settings.ConfigInit()
	filemanager.Init()
	fyneInit.Init()
}
