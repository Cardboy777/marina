package importconfigs

import (
	"fmt"
	"marina/files"
	"marina/stores"
	marina "marina/types"
	"marina/ui/dialogs"

	g "github.com/AllenDang/giu"
)

const importModalId = "Import"

var (
	stableVersion             *marina.Version
	unstableVersion           *marina.UnstableVersion
	installedStableVersions   []marina.Version
	installedUnstableVersions []marina.UnstableVersion
	selectedVersion           int32
	versionList               []string
	previewText               = "Select a Version"
	importEnabled             = false

	transferSettings   = true
	transferMods       = true
	transferRandomizer = false
	transferSaves      = false
)

func GetImportDialog() *g.PopupModalWidget {
	return g.PopupModal(importModalId).Layout(
		// Selector here
		g.Column(

			g.Label("Import from:"),
			g.Combo("", previewText, versionList, &selectedVersion).OnChange(updatePreviewText),
			g.Spacing(),
			g.Label("Select the types of files you want to import:"),
			g.Checkbox("Settings", &transferSettings),
			g.Checkbox("Mods", &transferMods),
			g.Checkbox("Randomizer Seeds", &transferRandomizer),
			g.Checkbox("Save Files", &transferSaves),
			g.Spacing(),
			g.Spacing(),
			g.Spacing(),
			g.Spacing(),
			g.Align(g.AlignCenter).To(
				g.Row(
					g.Button("Cancel").OnClick(closeDialog),
					g.Button("Import").OnClick(importFiles).Disabled(!importEnabled),
				),
			),
		),
	).Flags(g.WindowFlagsNoDocking).Flags(g.WindowFlagsNoResize).Flags(g.WindowFlagsAlwaysAutoResize)
}

func updatePreviewText() {
	importEnabled = true

	if int(selectedVersion) < len(installedUnstableVersions) {
		previewText = installedUnstableVersions[selectedVersion].GetName()
	} else {
		previewText = installedStableVersions[selectedVersion].GetName()
	}

	g.Update()
}

func closeDialog() {
	versionList = []string{}
	g.CloseCurrentPopup()
}

func ShowDialogStable(version *marina.Version) {
	stableVersion = version
	unstableVersion = nil
	show()
}

func ShowDialogUnstable(version *marina.UnstableVersion) {
	stableVersion = nil
	unstableVersion = version
	show()
}

func show() {
	importEnabled = false
	g.OpenPopup(importModalId)
	g.Update()
	generateVersionOptions()
}

func generateVersionOptions() {
	var repo *marina.Repository

	if stableVersion != nil {
		repo = stableVersion.Repository
	} else {
		repo = unstableVersion.Repository
	}

	stable := stores.GetVersions(repo)
	unstable := stores.GetUnstableVersions(repo)

	for _, v := range *unstable {
		if v.Installed && (unstableVersion == nil || (unstableVersion.Hash != v.Hash)) {
			installedUnstableVersions = append(installedUnstableVersions, v)
			versionList = append(versionList, v.GetName())
		}
	}

	for _, v := range *stable {
		if v.Installed && (stableVersion == nil || (stableVersion.TagName != v.TagName)) {
			installedStableVersions = append(installedStableVersions, v)
			versionList = append(versionList, v.GetName())
		}
	}
}

func getSourceInstallDir() string {
	if int(selectedVersion) < len(installedUnstableVersions) {
		return files.GetUnstableVersionInstallDirPath(&installedUnstableVersions[selectedVersion])
	}
	return files.GetVersionInstallDirPath(&installedStableVersions[selectedVersion])
}

func getDestinationInstallDir() (string, error) {
	if stableVersion != nil {
		return files.GetVersionInstallDirPath(stableVersion), nil
	}
	if unstableVersion != nil {
		return files.GetUnstableVersionInstallDirPath(unstableVersion), nil
	}

	return "", fmt.Errorf("No Version Selected")
}

func importFiles() {
	dest, err := getDestinationInstallDir()
	if err != nil {
		panic(err)
	}
	src := getSourceInstallDir()

	errorText := ""

	if transferSettings {
		err := files.ImportSettings(src, dest)
		if err != nil {
			errorText = fmt.Sprintf("%s\n%s", errorText, err)
		}
	}
	if transferRandomizer {
		err := files.ImportRandomizer(src, dest)
		if err != nil {
			errorText = fmt.Sprintf("%s\n%s", errorText, err)
		}
	}
	if transferSaves {
		err := files.ImportSaves(src, dest)
		if err != nil {
			errorText = fmt.Sprintf("%s\n%s", errorText, err)
		}
	}
	if transferMods {
		err := files.ImportMods(src, dest)
		if err != nil {
			errorText = fmt.Sprintf("%s\n%s", errorText, err)
		}
	}

	if len(errorText) > 0 {
		dialogs.ShowErrorDialog(fmt.Errorf("%s", errorText))
		return
	}

	closeDialog()
}
