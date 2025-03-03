package mainwindow

import (
	"marina/constants"
	"marina/types"

	g "github.com/AllenDang/giu"
)

func selectGame(repository *marina.Repository) {
	SelectedGame = repository
	g.Update()
}

func getGameOptions() []string {
	list := make([]string, len(constants.Repositories))

	for i, r := range constants.Repositories {
		list[i] = r.Name
	}

	return list
}

var selectedIndex int32 = 0

func GetGameSelector() *g.ColumnWidget {
	return g.Column(
		g.Label("Select Game:"),
		g.ListBox(getGameOptions()).SelectedIndex(&selectedIndex).OnChange(func(index int) {
			selectGame(constants.Repositories[index])
		}).Size(0, 200),
	)
}
