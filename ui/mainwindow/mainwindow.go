package mainwindow

import (
	"marina/constants"
	"marina/types"

	g "github.com/AllenDang/giu"
)

var SelectedGame *marina.Repository = &constants.SohDefinition

var mainSplit float32 = 200

func Loop() {
	g.SingleWindow().Layout(
		g.Style().SetFontSize(19).To(
			g.SplitLayout(g.DirectionVertical, &mainSplit,
				g.Column(
					GetGameSelector(),
					g.Spacing(),
					g.Spacing(),
					g.Spacing(),
					g.Spacing(),
					g.Separator(),
					g.Spacing(),
					GetRomDisplay(),
				),
				g.Column(
					GetToolbar(),
					GetGameList(),
				),
			),
		),
	)
}
