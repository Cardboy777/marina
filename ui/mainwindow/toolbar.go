package mainwindow

import (
	"marina/services"
	"marina/ui/dialogs"
	"marina/ui/fonts"
	"marina/ui/settings"

	g "github.com/AllenDang/giu"
)

func RefreshVersions(force bool) {
	err := services.SyncReleases(SelectedGame, force)
	if err != nil {
		dialogs.ShowErrorDialog(err)
	}
}

func GetToolbar() *g.AlignmentSetter {
	return g.Align(g.AlignRight).To(
		g.Row(
			g.Style().SetFont(fonts.GlyphFont).SetFontSize(fonts.IconSize).To(g.Button("").OnClick(settings.ShowDialog)),
			g.Style().SetFont(fonts.GlyphFont).SetFontSize(fonts.IconSize).To(g.Button("").OnClick(func() { RefreshVersions(true) })),
			settings.GetSettingsDialog(),
		),
	)
}
