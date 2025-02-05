package ui

import (
	"fmt"
	"marina/constants"
	"marina/db"
	"marina/files"
	"marina/launcher"
	"marina/stores"
	"marina/types"
	"marina/ui/widgets"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var PrimaryWindow *fyne.Window

func GetPrimaryWindow() *fyne.Window {
	return PrimaryWindow
}

func CreateMainWindow(app *fyne.App) {
	window := (*app).NewWindow("Marina - Ship Launcher")
	PrimaryWindow = &window
	window.Resize(fyne.Size{Width: 1000, Height: 600})

	initWindow(window)
	selectGame(selectedGame)
	window.ShowAndRun()
}

var (
	gameSelectorBox        = container.NewVBox()
	selectedGameTitleLabel = widget.NewLabel(selectedGame.Name)
	installedRomsLabel     = widget.NewLabel("None")
)

func initWindow(window fyne.Window) {
	for _, def := range constants.Repositories {
		gameButton := widget.NewButton(def.Name, func() {
			selectGame(def)
		})
		gameSelectorBox.Add(gameButton)
	}

	gameSelector := container.NewVScroll(
		gameSelectorBox,
	)
	gameSelector.SetMinSize(fyne.NewSize(0, float32(len(constants.Repositories)*40)))

	versionSelector := container.NewVScroll(versionList)

	toolbar := container.NewBorder(
		nil,
		nil,
		selectedGameTitleLabel,
		widget.NewToolbar(
			widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
				syncReleases(true)
			}),
		),
		container.NewCenter(widget.NewLabel(constants.AppName)),
	)

	addRomsButton := widget.NewButton("Add ROMs", func() {
		ShowFilePickerDialogFiltered("Choose ROM", "Nintendo 64 ROMs", []string{"z64"}, onFilesSelected)
	})
	romBox := container.NewVBox(installedRomsLabel, addRomsButton)

	window.SetContent(
		container.NewBorder(
			toolbar,
			nil,
			container.NewBorder(gameSelector, romBox, nil, nil, nil),
			nil,
			versionSelector,
		),
	)
}

func onFilesSelected(path string, err error) {
	if err != nil {
		ShowErrorDialog(err)
		return
	}
	if path == "" {
		// canceled
		return
	}
	err = files.CopyRomToInstallDir(selectedGame, path)
	if err != nil {
		ShowErrorDialog(err)
	}

	updateRomText()
}

func selectGame(def *marina.Repository) {
	selectedGame = def
	versionList.Refresh()
	selectedGameTitleLabel.SetText(def.Name)
	updateRomText()

	versions := getCurrentGameVersions()

	if len(*versions) == 0 {
		go syncReleases(false)
	}
}

func getCurrentGameVersions() *[]marina.Version {
	versions := stores.GetVersions(selectedGame)
	return versions
}

func updateRomText() {
	roms := stores.GetInstalledRomsList(selectedGame)
	if roms == nil || len(*roms) == 0 {
		installedRomsLabel.SetText("None")
		return
	}
	names := []string{}
	for _, r := range *roms {
		names = append(names, r.Name)
	}
	installedRomsLabel.SetText(strings.Join(names, "\n"))
	installedRomsLabel.Refresh()
}

var selectedGame = &constants.SohDefinition

var versionList = widget.NewList(
	func() int {
		return len(*getCurrentGameVersions())
	},
	func() fyne.CanvasObject {
		return widgets.NewVersionListItemWidget(nil, downloadVersion, playVersion, deleteVersion)
	},
	func(i widget.ListItemID, o fyne.CanvasObject) {
		vDef := (*getCurrentGameVersions())[i]
		item := o.(*widgets.VersionListItemWidget)

		item.Update(&vDef)
	})

func playVersion(version *marina.Version, onClose func()) {
	files.CopyRomsToVersionInstall(version)
	err := launcher.LaunchGame(version, func(e error) {
		if e != nil {
			ShowErrorDialog(e)
		}
		onClose()
	})
	if err != nil {
		ShowErrorDialog(err)
	}
}

func downloadVersion(version *marina.Version, update func()) {
	err := files.DownloadVersion(version)
	if err != nil {
		ShowErrorDialog(err)
		return
	}
	update()
}

func deleteVersion(version *marina.Version, update func()) {
	ShowConfirmDialog("Delete", fmt.Sprintf("Delete %s?", version.Name), func(shouldDelete bool) {
		if !shouldDelete {
			return
		}
		err := files.DeleteVersion(version)
		if err != nil {
			ShowErrorDialog(err)
			return
		}
		version.Installed = false
		db.SetInstalled(version, false)
		update()
	})
}

func syncReleases(force bool) {
	err := files.SyncReleases(selectedGame, force)
	if err != nil {
		ShowErrorDialog(err)
	}
	versionList.Refresh()
}
