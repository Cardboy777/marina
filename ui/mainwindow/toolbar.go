package mainwindow

import (
	"marina/services"
	"marina/ui/dialogs"
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
			g.Button(" Settings").OnClick(settings.ShowDialog),
			g.Button(" Refresh").OnClick(func() { RefreshVersions(true) }),
			settings.GetSettingsDialog(),
		),
	)
}
