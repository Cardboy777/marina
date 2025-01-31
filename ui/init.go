package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var appInstance *fyne.App

func Init() {
	app := app.New()
	appInstance = &app
	CreateMainWindow(appInstance)
}

func GetApp() *fyne.App {
	return appInstance
}
