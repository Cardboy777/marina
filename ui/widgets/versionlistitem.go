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
	primaryButton    *widget.Button
	deleteButton     *widget.Button
	name             *widget.Label
	compatability    *widget.Label
	releaseDate      *widget.Label
	latestTag        *widget.Label
	version          *marina.Version
	deleteCallback   func(*marina.Version, func())
	downloadCallback func(*marina.Version, func())
	playCallback     func(*marina.Version, func())
	isPlaying        bool
}

func NewVersionListItemWidget(version *marina.Version, downloadCallback func(*marina.Version, func()), playCallback func(*marina.Version, func()), deleteCallback func(*marina.Version, func())) *VersionListItemWidget {
	var item *VersionListItemWidget
	item = &VersionListItemWidget{
		primaryButton:    widget.NewButtonWithIcon("Download", theme.DownloadIcon(), func() { item.primaryAction() }),
		deleteButton:     widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() { item.deleteAction() }),
		name:             widget.NewLabel(""),
		compatability:    widget.NewLabel("OS not supported"),
		releaseDate:      widget.NewLabel(""),
		latestTag:        widget.NewLabel("- Latest -"),
		version:          version,
		deleteCallback:   deleteCallback,
		downloadCallback: downloadCallback,
		playCallback:     playCallback,
	}

	leftContainer := container.NewStack(
		container.NewVBox(
			container.NewHBox(
				container.NewStack(
					// canvas.NewRectangle(marina.NewColor(0, 0, 0xffff, 0xffff)),
					item.latestTag,
				),
				layout.NewSpacer(),
			),
		),
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
			item.deleteButton,
			item.primaryButton,
		),
		layout.NewSpacer(),
	)

	item.latestTag.Theme()
	item.latestTag.Hide()

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

func (item *VersionListItemWidget) primaryAction() {
	if item == nil || item.version == nil {
		return
	}

	item.setButtonsEnabled(false)
	if item.version.Installed {
		item.isPlaying = true
		item.playCallback(item.version, func() {
			item.isPlaying = false
			item.setButtonsEnabled(true)
		})
	} else {
		item.downloadCallback(item.version, func() {
			item.Update(item.version)
			item.setButtonsEnabled(true)
		})
	}
}

func (item *VersionListItemWidget) deleteAction() {
	if item == nil || item.version == nil {
		return
	}

	item.setButtonsEnabled(false)
	item.deleteCallback(item.version, func() { item.Update(item.version) })
	item.setButtonsEnabled(true)
}

func (item *VersionListItemWidget) Update(version *marina.Version) {
	item.version = version
	if item.version == nil {
		return
	}

	item.releaseDate.SetText(fmt.Sprintf("Release: %s", version.ReleaseDate.Local().Format(time.DateOnly)))

	isOSCompatible := version.IsOSCompatible()
	if !isOSCompatible {
		item.compatability.Show()
	} else {
		item.compatability.Hide()
	}

	item.name.SetText(version.Name)

	if item.version.Installed {
		item.primaryButton.SetText("Play")
		item.primaryButton.SetIcon(theme.MediaPlayIcon())
		item.deleteButton.Show()
	} else {
		item.primaryButton.SetText("Download")
		item.primaryButton.SetIcon(theme.DownloadIcon())
		item.deleteButton.Hide()
	}

	if isOSCompatible && !item.isPlaying {
		item.primaryButton.Enable()
	} else {
		item.primaryButton.Disable()
	}
	//
	// latestVersion := stores.GetLatestVersion(item.version.Repository)
	// latestDev := stores.GetLatestUnstableVersion(item.version.Repository)
	//
	// if latestVersion != nil && latestVersion.TagName == item.version.TagName {
	// 	item.latestTag.SetText("Latest Release")
	// 	item.latestTag.Show()
	// } else if latestDev != nil && latestDev.Hash == item.version.TagName {
	// 	item.latestTag.SetText("Latest Unstable Release")
	// 	item.latestTag.Show()
	// } else {
	// 	item.latestTag.Hide()
	// }
}

func (item *VersionListItemWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(item.content)
}

func (item *VersionListItemWidget) setButtonsEnabled(enabled bool) {
	if enabled {
		item.primaryButton.Enable()
		item.deleteButton.Enable()
	} else {
		item.primaryButton.Disable()
		item.deleteButton.Disable()
	}
}
