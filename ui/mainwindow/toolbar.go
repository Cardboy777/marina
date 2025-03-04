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
	refreshBtn := g.Button("Refresh")
	refreshBtn.OnClick(func() { RefreshVersions(true) })

	settingsBtn := g.Button("Settings")
	settingsBtn.OnClick(settings.ShowDialog)

	return g.Align(g.AlignRight).To(
		g.Row(
			settingsBtn,
			refreshBtn,
			settings.GetSettingsDialog(),
		),
	)
}
