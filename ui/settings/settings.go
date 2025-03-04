package settings

import (
	"marina/settings"
	"marina/ui/dialogs"

	g "github.com/AllenDang/giu"
)

var installDirectoryInput string

func GetSettingsDialog() *g.PopupModalWidget {
	return g.PopupModal("Settings").Layout(
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
					saveChangesButton(),
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
	btn := g.Button("Select")

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

func saveChangesButton() *g.ButtonWidget {
	btn := g.Button("Save Changes")

	btn.OnClick(func() {
		if settings.GetInstallDirName() != installDirectoryInput {
			settings.SetInstallDir(installDirectoryInput)
			g.Context.Backend().SetShouldClose(true)
		} else {
			g.CloseCurrentPopup()
		}
	})

	return btn
}
