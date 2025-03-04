package mainwindow

import (
	"marina/files"
	"marina/stores"
	"marina/ui/dialogs"

	g "github.com/AllenDang/giu"
)

func getRomsList() []g.Widget {
	list := []g.Widget{}

	roms := stores.GetInstalledRomsList(SelectedGame)

	for _, r := range *roms {
		list = append(list, g.Label(r.Name))
	}

	return list
}

func GetRomDisplay() *g.ColumnWidget {
	return g.Column(
		g.Button("Add Rom").OnClick(addRom),
		g.Label("Installed Roms:"),
		g.Column(
			getRomsList()...,
		),
	)
}

var romFileExtensions = []string{"z64", "n64"}

func addRom() {
	file, err := dialogs.ShowFilePickerDialogFiltered("Choose Rom", "Nintendo 64 ROM", romFileExtensions)
	if err != nil {
		dialogs.ShowErrorDialog(err)
		return
	}

	if file == "" {
		return
	}

	err = files.CopyRomToInstallDir(SelectedGame, file)
	if err != nil {
		dialogs.ShowErrorDialog(err)
	}
	g.Update()
}
