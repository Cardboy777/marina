package mainwindow

import (
	"marina/constants/games"
	"marina/types"

	g "github.com/AllenDang/giu"
)

func selectGame(repository *marina.Repository) {
	SelectedGame = repository
	g.Update()
}

func getGameOptions() []string {
	list := make([]string, len(games.Repositories))

	for i, r := range games.Repositories {
		list[i] = r.Name
	}

	return list
}

var selectedIndex int32 = 0

func GetGameSelector() *g.ColumnWidget {
	return g.Column(
		g.Label("Select Game:"),
		g.ListBox(getGameOptions()).SelectedIndex(&selectedIndex).OnChange(func(index int) {
			selectGame(games.Repositories[index])
		}).Size(0, 200),
	)
}
