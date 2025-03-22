package mainwindow

import (
	"fmt"
	"marina/files"
	"marina/launcher"
	"marina/stores"
	"marina/types"
	"marina/ui/dialogs"
	"marina/ui/fonts"
	"time"

	g "github.com/AllenDang/giu"
)

type VersionListItem struct {
	StableVersion   *marina.Version
	UnstableVersion *marina.UnstableVersion
}

func (i *VersionListItem) isStableVersion() bool {
	return i.StableVersion != nil
}

func (i *VersionListItem) getName() string {
	if i.isStableVersion() {
		return i.StableVersion.Name
	}
	return fmt.Sprintf("Unstable - %s", i.UnstableVersion.ReleaseDate.Format(time.DateTime))
}

func (i *VersionListItem) getCaption() string {
	if i.isStableVersion() {
		return i.StableVersion.ReleaseDate.Format(time.DateOnly)
	}
	return fmt.Sprintf("Commit: %s", i.UnstableVersion.Hash)
}

func (i *VersionListItem) getPopoupName() string {
	return fmt.Sprintf("Options##%s", i.getName())
}

func (i *VersionListItem) openPopoup() {
	g.OpenPopup(i.getPopoupName())
	g.Update()
}

func (i *VersionListItem) getWidget() *g.TableRowWidget {
	return g.TableRow(
		g.Column(
			g.Spacing(),
			g.Spacing(),
			g.Spacing(),
			g.Row(
				i.getInfo(),
				g.Align(g.AlignRight).To(
					i.getButtons(),
				),
			),
			g.Spacing(),
			g.Spacing(),
			g.Spacing(),
		),
	)
}

func (i *VersionListItem) getInfo() *g.ColumnWidget {
	return g.Column(
		g.Label(i.getName()),
		g.Style().SetColor(g.StyleColorText, fonts.ColorCaption).SetFontSize(fonts.CaptionSize).To(
			g.Label(i.getCaption()),
		),
	)
}

func (i *VersionListItem) getButtons() *g.RowWidget {
	isInstalled := (i.isStableVersion() && i.StableVersion.Installed) ||
		(!i.isStableVersion() && i.UnstableVersion.Installed)

	if isInstalled {
		return g.Row(
			g.Button(" Play").OnClick(i.play),
			g.Button("").OnClick(i.openPopoup),
			g.Popup(i.getPopoupName()).Layout(
				g.Column(
					g.Button(" Open Dir").OnClick(i.openDir),
					// g.Button("󰋺 Import").OnClick(i.importConfig),
					g.Separator(),
					g.Style().SetColor(g.StyleColorText, fonts.ColorDestructive).To(
						g.Button(" Delete").OnClick(i.delete),
					),
				),
			),
		)
	}

	var canDownload bool
	if i.isStableVersion() {
		canDownload = i.StableVersion.CanDownload()
	} else {
		canDownload = i.UnstableVersion.CanDownload()
	}

	return g.Row(
		g.Button(" Install").Disabled(!canDownload).OnClick(i.install),
	)
}

func (i *VersionListItem) install() {
	var err error
	if i.isStableVersion() {
		err = files.DownloadVersion(i.StableVersion)
	} else {
		err = files.DownloadUnstableVersion(i.UnstableVersion)
	}

	dialogs.ShowDialogIfError(err)
	g.Update()
}

func (i *VersionListItem) play() {
	var err error

	if i.isStableVersion() {
		err = launcher.LaunchGame(i.StableVersion, dialogs.ShowDialogIfError)
	} else {
		err = launcher.LaunchUnstableGame(i.UnstableVersion, dialogs.ShowDialogIfError)
	}

	dialogs.ShowDialogIfError(err)
}

func (i *VersionListItem) delete() {
	if !dialogs.ShowConfirmDialog("Delete?", "Delete version? Configurations, Saves, and Mods will be permanently lost.") {
		return
	}

	var err error
	if i.isStableVersion() {
		err = files.DeleteVersion(i.StableVersion)
	} else {
		err = files.DeleteUnstableVersion(i.UnstableVersion)
		stores.RemoveOldUnstableVersions()
	}

	dialogs.ShowDialogIfError(err)
	g.Update()
}

func (i *VersionListItem) openDir() {
	var dir string
	if i.isStableVersion() {
		dir = files.GetVersionInstallDirPath(i.StableVersion)
	} else {
		dir = files.GetUnstableVersionInstallDirPath(i.UnstableVersion)
	}

	dialogs.OpenDirectory(dir)
}

func (i *VersionListItem) importConfig() {
}

func getVersionListItems() *[]VersionListItem {
	stableVersions := stores.GetVersions(SelectedGame)
	unstableVersions := stores.GetUnstableVersions(SelectedGame)

	unstableLen := len(*unstableVersions)
	stableLen := len(*stableVersions)

	versions := make([]VersionListItem, stableLen+unstableLen)

	for i, v := range *unstableVersions {
		versions[i].UnstableVersion = &v
	}

	for i, v := range *stableVersions {
		versions[unstableLen+i].StableVersion = &v
	}

	return &versions
}

func getVersionRows() []*g.TableRowWidget {
	versions := getVersionListItems()
	rows := make([]*g.TableRowWidget, len(*versions))

	for i, v := range *versions {
		rows[i] = v.getWidget()
	}

	return rows
}

func GetGameList() *g.TableWidget {
	return g.Table().FastMode(true).Rows(getVersionRows()...).NoHeader(true)
}
