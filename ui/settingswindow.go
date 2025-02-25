package ui

import (
	"marina/settings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func showSettingsDialog(window fyne.Window) {
	dialog := dialog.NewCustom("Settings", "Close", getSettingsDialog(), window)

	dialog.Show()

	dialog.Resize(fyne.NewSize(540, 140))
}

func getSettingsDialog() *fyne.Container {
	dialogContainer := container.NewVBox(
		getInstallLocationRow(),
	)

	return dialogContainer
}

var currentInstallDisplay *widget.Entry

func getInstallLocationRow() *fyne.Container {
	pathBinding := binding.NewString()
	_ = pathBinding.Set(settings.GetInstallDirName())
	pathBinding.AddListener(binding.NewDataListener(func() {
		val, err := pathBinding.Get()
		if err != nil {
			settings.SetInstallDir(val)
		}
	}))

	currentInstallDisplay = widget.NewEntryWithData(pathBinding)
	chooseFolderButton := widget.NewButtonWithIcon("", theme.FolderIcon(), changeInstallDir)

	return container.New(layout.NewFormLayout(),
		widget.NewLabel("Install Location"),
		container.NewBorder(nil, nil, nil, chooseFolderButton, currentInstallDisplay),
	)
}

func changeInstallDir() {
	dir, err := ShowDirectoryPickerDialog("Choose Install Directory")
	if err != nil {
		ShowErrorDialog(err)
		return
	}

	settings.SetInstallDir(dir)
	currentInstallDisplay.SetText(settings.GetInstallDirName())
}
