package ui

import (
	"marina/constants"
	"marina/ui/fonts"
	"marina/ui/mainwindow"

	g "github.com/AllenDang/giu"
)

func Start() {
	wnd := g.NewMasterWindow(constants.AppName, 800, 600, 0)

	fonts.InitializeFonts()

	go mainwindow.RefreshVersions(false)

	wnd.Run(mainwindow.Loop)
}
