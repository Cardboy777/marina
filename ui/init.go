package fyneInit

import (
	"fmt"
	"marina/constants"
	"marina/filemanager"
	"marina/types"
	"marina/ui/components"
	"marina/versionmanager"
	"marina/versionmanager/gamedefinitions"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var selectedGame = &gamedefinitions.SohDefinition

var versionList = widget.NewList(
	func() int {
		return len(*getCurrentGameVersions())
	},
	func() fyne.CanvasObject {
		return marinacomponents.NewVersionListItemWidget(nil)
	},
	func(i widget.ListItemID, o fyne.CanvasObject) {
		vDef := (*getCurrentGameVersions())[i]
		item := o.(*marinacomponents.VersionListItemWidget)

		item.Update(&vDef)
	})

var (
	a fyne.App
	w fyne.Window
)

func Init() {
	a = app.New()
	w = a.NewWindow("Marina - Ship Launcher")
	w.Resize(fyne.Size{Width: 640, Height: 400})

	initWindow(w)
	go func() {
		versionmanager.SyncReleases()

		selectGame(selectedGame)

		versionList.Refresh()
	}()

	w.ShowAndRun()
}

var (
	gameSelectorBox        = container.NewVBox()
	selectedGameTitleLabel = widget.NewLabel(selectedGame.Name)
	installedRomsLabel     = widget.NewLabel("None")
)

func initWindow(window fyne.Window) {
	for _, def := range gamedefinitions.RepositoryDefinitions {
		gameButton := widget.NewButton(def.Name, func() {
			selectGame(def)
		})
		gameSelectorBox.Add(gameButton)
	}

	gameSelector := container.NewVScroll(
		gameSelectorBox,
	)
	gameSelector.SetMinSize(fyne.NewSize(0, float32(len(gamedefinitions.RepositoryDefinitions)*40)))

	versionSelector := container.NewVScroll(versionList)

	toolbar := container.NewBorder(
		nil,
		nil,
		selectedGameTitleLabel,
		widget.NewToolbar(
			widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
				versionmanager.SyncReleases()
			}),
		),
		container.NewCenter(widget.NewLabel(constants.AppName)),
	)

	addRomsButton := widget.NewButton("Add ROMs", func() {
		fmt.Println("Clicked Button")
		dialog := dialog.NewFileOpen(onFilesSelected, window)
		dialog.Show()
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

func onFilesSelected(reader fyne.URIReadCloser, err error) {
	if err != nil {
		dialog.ShowError(err, w)
	}
	uri := reader.URI()
	path := uri.Path()

	err = filemanager.CopyRomToInstallDir(selectedGame, path)
	if err != nil {
		dialog.ShowError(err, w)
	}

	updateRomText()
}

func selectGame(def *marina.RepositoryDefinition) {
	selectedGame = def
	versionList.Refresh()
	selectedGameTitleLabel.SetText(def.Name)
	updateRomText()
}

func getCurrentGameVersions() *[]marina.VersionDefinition {
	versionList := *versionmanager.GetVersionsList()
	if len(versionList) <= selectedGame.Id {
		return &([]marina.VersionDefinition{})
	}
	return &(versionList[selectedGame.Id])
}

func updateRomText() {
	roms := filemanager.GetInstalledRoms(selectedGame)
	if roms == nil || len(*roms) == 0 {
		installedRomsLabel.SetText("None")
		return
	}
	names := []string{}
	for _, r := range *roms {
		names = append(names, r.Name)
	}
	installedRomsLabel.SetText(strings.Join(names, ", "))
}
