package widgets

import (
	"fmt"
	"marina/types"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type VersionListItemWidget struct {
	widget.BaseWidget
	content          *fyne.Container
	playButton       *widget.Button
	downloadButton   *widget.Button
	deleteButton     *widget.Button
	openDirButton    *widget.Button
	name             *widget.Label
	compatability    *widget.Label
	releaseDate      *widget.Label
	listItem         *ListItem
	deleteCallback   func(*ListItem, func())
	downloadCallback func(*ListItem, func())
	playCallback     func(*ListItem, func())
	openDirCallback  func(*ListItem)
}

func NewVersionListItemWidget(version *ListItem, downloadCallback func(*ListItem, func()), playCallback func(*ListItem, func()), deleteCallback func(*ListItem, func()), openDirCallback func(*ListItem)) *VersionListItemWidget {
	var item *VersionListItemWidget
	item = &VersionListItemWidget{
		playButton:       widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), func() { item.playAction() }),
		downloadButton:   widget.NewButtonWithIcon("Download", theme.DownloadIcon(), func() { item.downloadAction() }),
		deleteButton:     widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() { item.deleteAction() }),
		openDirButton:    widget.NewButtonWithIcon("", theme.FolderIcon(), func() { item.openDirAction() }),
		name:             widget.NewLabel(""),
		compatability:    widget.NewLabel("OS-compatible version not found"),
		releaseDate:      widget.NewLabel(""),
		listItem:         version,
		deleteCallback:   deleteCallback,
		downloadCallback: downloadCallback,
		playCallback:     playCallback,
		openDirCallback:  openDirCallback,
	}

	leftContainer := container.NewStack(
		container.NewVBox(
			layout.NewSpacer(),
			container.NewHBox(
				item.name,
				item.compatability,
			),
			item.releaseDate,
			layout.NewSpacer(),
		),
	)

	rightContainer := container.NewVBox(
		layout.NewSpacer(),
		container.NewHBox(
			item.openDirButton,
			item.deleteButton,
			item.playButton,
			item.downloadButton,
		),
		layout.NewSpacer(),
	)

	item.Update(version)

	item.content = container.NewBorder(
		nil,
		nil,
		leftContainer,
		rightContainer,
		nil,
	)

	item.ExtendBaseWidget(item)

	return item
}

func (item *VersionListItemWidget) playAction() {
	if item == nil || item.listItem == nil {
		return
	}
	item.playCallback(item.listItem, func() {})
}

func (item *VersionListItemWidget) downloadAction() {
	if item == nil || item.listItem == nil {
		return
	}

	item.downloadCallback(item.listItem, func() {
		item.Update(item.listItem)
	})
}

func (item *VersionListItemWidget) deleteAction() {
	if item == nil || item.listItem == nil {
		return
	}

	item.deleteCallback(item.listItem, func() { item.Update(item.listItem) })
}

func (item *VersionListItemWidget) openDirAction() {
	if item == nil || item.listItem == nil {
		return
	}

	item.openDirCallback(item.listItem)
}

func (item *VersionListItemWidget) Update(listItem *ListItem) {
	item.listItem = listItem

	if item.listItem == nil {
		return
	} else if item.listItem.IsStableRelease && item.listItem.Release != nil {
		item.UpdateVersion(item.listItem.Release)
	} else if !item.listItem.IsStableRelease && item.listItem.UnstableRelease != nil {
		item.UpdateUnstableVersion(item.listItem.UnstableRelease)
	}
}

func (item *VersionListItemWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(item.content)
}

func (item *VersionListItemWidget) UpdateVersion(version *marina.Version) {
	item.releaseDate.SetText(fmt.Sprintf("Release: %s", version.ReleaseDate.Local().Format(time.DateOnly)))

	isOSCompatible := version.IsOSCompatible()
	if !isOSCompatible {
		item.compatability.Show()
	} else {
		item.compatability.Hide()
	}

	item.name.SetText(version.Name)

	if version.Installed {
		item.downloadButton.Hide()
		item.playButton.Show()
		item.deleteButton.Show()
		item.openDirButton.Show()
	} else {
		item.downloadButton.Show()
		item.playButton.Hide()
		item.deleteButton.Hide()
		item.openDirButton.Hide()
	}

	if isOSCompatible {
		item.downloadButton.Enable()
		item.playButton.Enable()
	} else {
		item.downloadButton.Disable()
		item.playButton.Disable()
	}
}

func (item *VersionListItemWidget) UpdateUnstableVersion(version *marina.UnstableVersion) {
	item.releaseDate.SetText(fmt.Sprintf("Release: %s", version.ReleaseDate.Local().Format(time.DateOnly)))

	item.compatability.Hide()

	name := fmt.Sprintf("Unstable Version - Commit: %s", version.Hash)

	item.name.SetText(name)

	if version.Installed {
		item.downloadButton.Hide()
		item.playButton.Show()
		item.deleteButton.Show()
		item.openDirButton.Show()
	} else {
		item.downloadButton.Show()
		item.playButton.Hide()
		item.deleteButton.Hide()
		item.openDirButton.Hide()
	}

	item.downloadButton.Enable()
	item.playButton.Enable()
}
