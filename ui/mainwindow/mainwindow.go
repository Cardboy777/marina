package mainwindow

import (
	"marina/constants"
	"marina/types"
	"marina/ui/fonts"

	"github.com/AllenDang/cimgui-go/imgui"
	g "github.com/AllenDang/giu"
)

var SelectedGame *marina.Repository = &constants.SohDefinition

var mainSplit float32 = 200

func Loop() {
	if false {
		imgui.ShowDemoWindow()
		return
	}

	g.SingleWindow().Layout(
		g.Style().SetFont(fonts.GlyphFont).SetFontSize(fonts.HeaderSize).To(
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
