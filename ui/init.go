package fyneInit

import (
	"marina/types"
	"marina/ui/components"
	"marina/versionmanager"
	"marina/versionmanager/gamedefinitions"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var selectedGame = &gamedefinitions.SohDefinition

var versionList = widget.NewList(
	func() int {
		return len(*getCurrentGameVersions())
	},
	func() fyne.CanvasObject {
		return marinacomponents.NewVersionListItemWidget("", false, "", false)
	},
	func(i widget.ListItemID, o fyne.CanvasObject) {
		vDef := (*getCurrentGameVersions())[i]
		item := o.(*marinacomponents.VersionListItemWidget)

		item.Update(vDef.Name, vDef.IsDownloaded(), vDef.GetDownloadUrl(), vDef.IsOSCompatible())
	})

func Init() {
	a := app.New()
	w := a.NewWindow("Marina - Ship Launcher")
	w.Resize(fyne.Size{Width: 640, Height: 400})

	initWindow(w)
	go func() {
		versionmanager.SyncReleases()
		selectGame(selectedGame)
		versionList.Refresh()
	}()

	w.ShowAndRun()
}

func initWindow(window fyne.Window) {
	gameSelectorBox := container.NewVBox()

	for _, def := range gamedefinitions.RepositoryDefinitions {
		gameSelectorBox.Add(widget.NewButton(def.Name, func() { selectGame(def) }))
	}

	gameSelector := container.NewVScroll(
		gameSelectorBox,
	)

	versionSelector := container.NewVScroll(versionList)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			versionmanager.SyncReleases()
		}),
	)

	window.SetContent(
		container.NewBorder(
			toolbar,
			nil,
			gameSelector,
			nil,
			versionSelector,
		),
	)
}

func selectGame(def *marina.RepositoryDefinition) {
	selectedGame = def
	versionList.Refresh()
}

func getCurrentGameVersions() *[]marina.VersionDefinition {
	versionList := *versionmanager.GetVersionsList()
	if len(versionList) <= selectedGame.Id {
		return &([]marina.VersionDefinition{})
	}
	return &(versionList[selectedGame.Id])
}
