package importconfigs

import (
	"errors"
	"fmt"
	"marina/files"
	"marina/stores"
	marina "marina/types"
	"marina/ui/dialogs"
	"marina/ui/fonts"

	g "github.com/AllenDang/giu"
)

const importModalId = "Import"

var (
	initialPreviewText        = "Select a Version"
	stableVersion             *marina.Version
	unstableVersion           *marina.UnstableVersion
	installedStableVersions   []marina.Version
	installedUnstableVersions []marina.UnstableVersion
	repository                *marina.Repository
	selectedVersion           int32
	versionList               []string
	previewText               = &initialPreviewText
	importEnabled             = false

	transferSettings   = true
	transferMods       = false
	transferRandomizer = false
	transferSaves      = false
)

func GetImportDialog() *g.PopupModalWidget {
	return g.PopupModal(importModalId).Layout(
		// Selector here
		g.Column(
			g.Label("Import from:"),
			g.Combo("", *previewText, versionList, &selectedVersion).OnChange(updatePreviewText),
			g.Spacing(),
			g.Spacing(),
			g.Label("Select the types of files you want to import:"),
			g.Style().SetDisabled(repository == nil || !repository.Imports.Configuration.Capable).To(
				g.Checkbox("Settings", &transferSettings),
			),
			g.Style().SetDisabled(repository == nil || !repository.Imports.Mods.Capable).To(
				g.Row(
					g.Checkbox("Mods", &transferMods),
					g.Column(
						g.Spacing(),
						g.Spacing(),
						g.Style().SetFontSize(fonts.CaptionSize).To(g.Label("Copying may take time. When installing mods, consider using shortcuts instead.")),
					),
				),
			),
			g.Style().SetDisabled(repository == nil || !repository.Imports.Randomizer.Capable).To(
				g.Checkbox("Randomizer Seeds", &transferRandomizer),
			),
			g.Style().SetDisabled(repository == nil || !repository.Imports.Saves.Capable).To(
				g.Checkbox("Save Files", &transferSaves),
			),
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

	previewText = &(versionList[selectedVersion])

	g.Update()
}

func closeDialog() {
	versionList = []string{}
	g.CloseCurrentPopup()
	g.Update()
}

func ShowDialogStable(version *marina.Version) {
	stableVersion = version
	repository = stableVersion.Repository
	unstableVersion = nil
	show()
}

func ShowDialogUnstable(version *marina.UnstableVersion) {
	stableVersion = nil
	unstableVersion = version
	repository = unstableVersion.Repository
	show()
}

func show() {
	importEnabled = false
	previewText = &initialPreviewText

	transferSettings = true
	transferMods = false
	transferRandomizer = false
	transferSaves = false

	g.OpenPopup(importModalId)
	g.Update()
	generateVersionOptions()
}

func generateVersionOptions() {
	stable := stores.GetVersions(repository)
	unstable := stores.GetUnstableVersions(repository)

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
	errs := []error{}

	if transferSettings {
		err := files.ImportCategory(src, dest, repository.Imports.Configuration)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if transferRandomizer {
		err := files.ImportCategory(src, dest, repository.Imports.Randomizer)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if transferSaves {
		err := files.ImportCategory(src, dest, repository.Imports.Saves)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if transferMods {
		err := files.ImportCategory(src, dest, repository.Imports.Mods)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		dialogs.ShowErrorDialog(errors.Join(errs...))
		return
	}

	dialogs.ShowInformationDialog("Import Successful", "The import was successful.")

	closeDialog()
}
