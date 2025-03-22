package importconfigs

import (
	"marina/settings"
	"marina/ui/dialogs"

	g "github.com/AllenDang/giu"
)

var (
	sourceVersion         *any
	installDirectoryInput string
)

func GetSettingsDialog() *g.PopupModalWidget {
	return g.PopupModal("Import").Layout(
		g.Column(
			g.Row(
				g.Label("Install Directory:"),
				g.InputText(&installDirectoryInput).Size(300).Hint(settings.GetDefaultInstallDir()),
				chooseDirButton(),
			),
			g.Spacing(),
			g.Spacing(),
			g.Spacing(),
			g.Spacing(),
			g.Align(g.AlignCenter).To(
				g.Row(
					cancelButton(),
					importButton(),
				),
			),
		),
	).Flags(g.WindowFlagsNoDocking).Flags(g.WindowFlagsNoResize).Flags(g.WindowFlagsAlwaysAutoResize)
}

func ShowDialog() {
	installDirectoryInput = settings.GetInstallDirName()
	g.OpenPopup("Settings")
	g.Update()
}

func chooseDirButton() *g.ButtonWidget {
	btn := g.Button("")

	btn.OnClick(func() {
		val, err := dialogs.ShowDirectoryPickerDialog("Choose Install Directory")
		if err != nil {
			dialogs.ShowErrorDialog(err)
		}

		if val != "" {
			installDirectoryInput = val
			g.Update()
		}
	})

	return btn
}

func cancelButton() *g.ButtonWidget {
	btn := g.Button("Cancel")
	btn.OnClick(func() {
		g.CloseCurrentPopup()
	})

	return btn
}

func importButton() *g.ButtonWidget {
	btn := g.Button("Import")

	btn.OnClick(func() {
		if settings.GetInstallDirName() != installDirectoryInput {
			// copy files
		} else {
			g.CloseCurrentPopup()
		}
	})

	return btn
}
