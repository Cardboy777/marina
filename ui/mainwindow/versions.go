package mainwindow

import (
	"fmt"
	"image/color"
	"marina/files"
	"marina/launcher"
	"marina/stores"
	"marina/types"
	"marina/ui/dialogs"
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

var subTextColor = color.RGBA{217, 217, 217, 255}

func (i *VersionListItem) getInfo() *g.ColumnWidget {
	if i.isStableVersion() {
		return g.Column(
			g.Label(i.StableVersion.Name),
			g.Style().SetColor(g.StyleColorText, subTextColor).SetFontSize(11).To(
				g.Label(i.StableVersion.ReleaseDate.Format(time.DateOnly)),
			),
		)
	}

	return g.Column(
		g.Label(fmt.Sprintf("Develop - %s", i.UnstableVersion.ReleaseDate.Format(time.DateTime))),
		g.Style().SetColor(g.StyleColorText, subTextColor).SetFontSize(11).To(
			g.Label(fmt.Sprintf("Commit: %s", i.UnstableVersion.Hash)),
		),
	)
}

func (i *VersionListItem) getButtons() *g.RowWidget {
	isInstalled := (i.isStableVersion() && i.StableVersion.Installed) ||
		(!i.isStableVersion() && i.UnstableVersion.Installed)

	if isInstalled {
		return g.Row(
			g.Button("Open").OnClick(i.openDir),
			g.Button("Delete").OnClick(i.delete),
			g.Button("Play").OnClick(i.play),
		)
	}

	return g.Row(g.Button("Install").OnClick(i.install))
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
