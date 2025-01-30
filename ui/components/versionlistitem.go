package marinacomponents

import (
	"fmt"
	"marina/filemanager"
	"marina/launcher"
	"marina/types"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type VersionListItemWidget struct {
	widget.BaseWidget
	leftContainer  *fyne.Container
	rightContainer *fyne.Container
	primaryButton  *widget.Button
	deleteButton   *widget.Button
	name           *widget.Label
	label          *widget.Label
	error          *widget.Label
	version        *marina.VersionDefinition
	takingAction   bool
}

func NewVersionListItemWidget(version *marina.VersionDefinition) *VersionListItemWidget {
	var item *VersionListItemWidget
	item = &VersionListItemWidget{
		leftContainer:  container.NewVBox(),
		rightContainer: container.NewHBox(),
		primaryButton: widget.NewButtonWithIcon("Download", theme.DownloadIcon(), func() {
			item.setButtonsEnabled(false)
			item.PrimaryAction()
			item.setButtonsEnabled(true)
		}),
		deleteButton: widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
			item.setButtonsEnabled(false)
			item.Delete()
			item.setButtonsEnabled(true)
		}),
		name:    widget.NewLabel(""),
		label:   widget.NewLabel("OS not supported"),
		error:   widget.NewLabel(""),
		version: version,
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

func (item *VersionListItemWidget) PrimaryAction() {
	if filemanager.IsVersionInstalled(item.version) {
		item.Play(item.version)
	} else {
		item.Download()
	}
}

func (item *VersionListItemWidget) Play(version *marina.VersionDefinition) {
	filemanager.CopyRomsToVersionInstall(item.version)
	err := launcher.LaunchGame(version)
	item.error.SetText(fmt.Sprintf("%s", err))
}

func (item *VersionListItemWidget) Download() {
	item.error.Hide()
	if item.version == nil {
		return
	}
	err := filemanager.DownloadVersion(item.version)
	if err != nil {
		item.error.SetText("Error downloading version.")
		item.error.Show()
	}
	item.Update(item.version)
}

func (item *VersionListItemWidget) Delete() {
	item.error.Hide()
	if item.version == nil {
		return
	}
	err := filemanager.DeleteVersion(item.version)
	if err != nil {
		item.error.SetText("Error deleting version.")
		item.error.Show()
	}
	item.Update(item.version)
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

	if filemanager.IsVersionInstalled(item.version) {
		item.primaryButton.SetText("Play")
		item.primaryButton.SetIcon(theme.MediaPlayIcon())
		item.deleteButton.Show()
	} else {
		item.primaryButton.SetText("Download")
		item.primaryButton.SetIcon(theme.DownloadIcon())
		item.deleteButton.Hide()
	}

	if isOSCompatible {
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
