package widgets

import (
	"marina/files"
	"marina/types"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type VersionListItemWidget struct {
	widget.BaseWidget
	leftContainer    *fyne.Container
	rightContainer   *fyne.Container
	primaryButton    *widget.Button
	deleteButton     *widget.Button
	name             *widget.Label
	label            *widget.Label
	error            *widget.Label
	version          *marina.VersionDefinition
	deleteCallback   func(*marina.VersionDefinition, func())
	downloadCallback func(*marina.VersionDefinition, func())
	playCallback     func(*marina.VersionDefinition, func())
	isPlaying        bool
}

func NewVersionListItemWidget(version *marina.VersionDefinition, downloadCallback func(*marina.VersionDefinition, func()), playCallback func(*marina.VersionDefinition, func()), deleteCallback func(*marina.VersionDefinition, func())) *VersionListItemWidget {
	var item *VersionListItemWidget
	item = &VersionListItemWidget{
		leftContainer:    container.NewVBox(),
		rightContainer:   container.NewHBox(),
		primaryButton:    widget.NewButtonWithIcon("Download", theme.DownloadIcon(), func() { item.primaryAction() }),
		deleteButton:     widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() { item.deleteAction() }),
		name:             widget.NewLabel(""),
		label:            widget.NewLabel("OS not supported"),
		error:            widget.NewLabel(""),
		version:          version,
		deleteCallback:   deleteCallback,
		downloadCallback: downloadCallback,
		playCallback:     playCallback,
	}

	item.error.Hide()

	item.leftContainer.Add(item.name)
	item.leftContainer.Add(item.label)
	item.leftContainer.Add(item.error)

	item.rightContainer.Add(item.deleteButton)
	item.rightContainer.Add(item.primaryButton)

	item.Update(version)

	item.ExtendBaseWidget(item)

	return item
}

func (item *VersionListItemWidget) primaryAction() {
	if item == nil || item.version == nil {
		return
	}
	item.error.Hide()

	item.setButtonsEnabled(false)
	if files.IsVersionInstalled(item.version) {
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
	item.error.Hide()

	item.setButtonsEnabled(false)
	item.deleteCallback(item.version, func() { item.Update(item.version) })
	item.setButtonsEnabled(true)
}

func (item *VersionListItemWidget) Update(version *marina.VersionDefinition) {
	item.version = version
	if item.version == nil {
		return
	}

	isOSCompatible := version.IsOSCompatible()
	if !isOSCompatible {
		item.label.Show()
	} else {
		item.label.Hide()
	}

	item.name.SetText(version.Name)

	if files.IsVersionInstalled(item.version) {
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
}

func (item *VersionListItemWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewBorder(nil, nil, container.NewVBox(layout.NewSpacer(), item.leftContainer, layout.NewSpacer()), container.NewVBox(layout.NewSpacer(), item.rightContainer, layout.NewSpacer()), nil)
	return widget.NewSimpleRenderer(c)
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
