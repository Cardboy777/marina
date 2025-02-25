package ui

import (
	"fmt"
	"marina/constants"
	"marina/files"
	"marina/launcher"
	"marina/services"
	"marina/stores"
	"marina/types"
	"marina/ui/widgets"
	"sort"
	"strings"
	"time"

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
				go syncReleases(true)
			}),
			widget.NewToolbarAction(theme.SettingsIcon(), func() {
				go showSettingsDialog(window)
			}),
		),
		container.NewCenter(widget.NewLabel(constants.AppName)),
	)

	addRomsButton := widget.NewButton("Add ROMs", func() {
		file, err := ShowFilePickerDialogFiltered("Choose ROM", "Nintendo 64 ROMs", []string{"z64"})
		onFilesSelected(file, err)
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

	go syncReleases(false)
}

func getCurrentGameVersions() *[]widgets.ListItem {
	var list []widgets.ListItem

	versions := stores.GetVersions(selectedGame)

	for _, v := range *versions {
		list = append(list, widgets.ListItem{
			IsStableRelease: true,
			Release:         &v,
		})
	}

	unstable := stores.GetUnstableVersions(selectedGame)

	for _, v := range *unstable {
		list = append(list, widgets.ListItem{
			IsStableRelease: false,
			UnstableRelease: &v,
		})
	}

	sort.Slice(list, func(i, j int) bool {
		l := list[i]
		r := list[j]

		var ldate time.Time
		if l.IsStableRelease {
			ldate = l.Release.ReleaseDate
		} else {
			ldate = l.UnstableRelease.ReleaseDate
		}

		var rdate time.Time
		if r.IsStableRelease {
			rdate = r.Release.ReleaseDate
		} else {
			rdate = r.UnstableRelease.ReleaseDate
		}

		return ldate.After(rdate)
	})

	return &list
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
		return widgets.NewVersionListItemWidget(nil, downloadVersion, playVersion, deleteVersion, openFolder)
	},
	func(i widget.ListItemID, o fyne.CanvasObject) {
		vDef := (*getCurrentGameVersions())[i]
		item := o.(*widgets.VersionListItemWidget)

		item.Update(&vDef)
	})

func playVersion(item *widgets.ListItem, onClose func()) {
	var err error
	if item.IsStableRelease {
		files.CopyRomsToVersionInstall(item.Release)
		err = launcher.LaunchGame(item.Release, func(e error) {
			if e != nil {
				ShowErrorDialog(e)
			}
			onClose()
		})
	} else {
		files.CopyRomsToUnstableVersionInstall(item.UnstableRelease)
		err = launcher.LaunchUnstableGame(item.UnstableRelease, func(e error) {
			if e != nil {
				ShowErrorDialog(e)
			}
			onClose()
		})

	}

	if err != nil {
		ShowErrorDialog(err)
	}
}

func downloadVersion(item *widgets.ListItem, update func()) {
	var err error
	if item.IsStableRelease {
		err = files.DownloadVersion(item.Release)
	} else {
		err = files.DownloadUnstableVersion(item.UnstableRelease)
	}

	if err != nil {
		ShowErrorDialog(err)
		return
	}
	update()
}

func openFolder(item *widgets.ListItem) {
	var dir string
	if item.IsStableRelease {
		dir = files.GetVersionInstallDirPath(item.Release)
	} else {
		dir = files.GetUnstableVersionInstallDirPath(item.UnstableRelease)
	}
	OpenDirectory(dir)
}

const deletetionWarning string = "Saves, configurations, and mods will be permanently deleted!"

func deleteVersion(item *widgets.ListItem, update func()) {
	var name string

	if item.IsStableRelease {
		name = item.Release.Name
	} else {
		name = fmt.Sprintf("Unstable Version - %s", item.UnstableRelease.ReleaseDate.Format(time.DateTime))
	}

	shouldDelete := ShowConfirmDialog("Delete", fmt.Sprintf("Delete %s?\n\n%s", name, deletetionWarning))
	if !shouldDelete {
		return
	}
	var err error
	if item.IsStableRelease {
		err = files.DeleteVersion(item.Release)
	} else {
		err = files.DeleteUnstableVersion(item.UnstableRelease)
	}

	if err != nil {
		ShowErrorDialog(err)
		return
	}

	update()
}

func syncReleases(force bool) {
	err := services.SyncReleases(selectedGame, force)
	if err != nil {
		ShowErrorDialog(err)
	}
	versionList.Refresh()
}
