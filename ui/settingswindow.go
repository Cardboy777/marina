package ui

import (
	"image/color"
	"marina/settings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	fyneTheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func showSettingsDialog(window fyne.Window) {
	dialog := dialog.NewCustomConfirm("Settings", "Save", "Cancel", getSettingsDialog(), dialogResultCallback, window)

	dialog.Show()

	dialog.Resize(fyne.NewSize(540, 140))
}

func dialogResultCallback(save bool) {
	if save {
		saveChanges()
	}
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

	currentInstallDisplay = widget.NewEntryWithData(pathBinding)
	chooseFolderButton := widget.NewButtonWithIcon("", fyneTheme.FolderIcon(), func() {
		dir, err := ShowDirectoryPickerDialog("Choose Install Directory")
		if err != nil {
			ShowErrorDialog(err)
			return
		}

		currentInstallDisplay.SetText(dir)
	})

	warning := canvas.NewText("Changing the install location will cause Marina to Restart", color.White)
	warning.TextStyle = fyne.TextStyle{Italic: true}
	warning.TextSize = 14

	return container.NewBorder(
		nil,
		warning,
		nil,
		layout.NewSpacer(),
		container.New(
			layout.NewFormLayout(),
			widget.NewLabel("Install Location"),
			container.NewBorder(
				nil,
				nil,
				nil,
				chooseFolderButton,
				currentInstallDisplay,
			),
		),
	)
}

func saveChanges() {
	path := currentInstallDisplay.Text
	settings.SetInstallDir(path)
}
