package marinacomponents

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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
	downloadUrl    string
}

func NewVersionListItemWidget(name string, isDownloaded bool, downloadUrl string, isOSCompatible bool) *VersionListItemWidget {
	item := &VersionListItemWidget{
		downloadUrl:    downloadUrl,
		leftContainer:  container.NewVBox(),
		rightContainer: container.NewHBox(),
		primaryButton: widget.NewButton("Download", func() {
			// doDelete
		}),
		deleteButton: widget.NewButton("Delete", func() {
			// doDownload with downloadUrl
		}),
		name:  widget.NewLabel(name),
		label: widget.NewLabel("OS not supported"),
	}

	item.leftContainer.Add(item.name)
	item.leftContainer.Add(item.label)

	item.rightContainer.Add(item.deleteButton)
	item.rightContainer.Add(item.primaryButton)

	item.Update(name, isDownloaded, downloadUrl, isOSCompatible)

	item.ExtendBaseWidget(item)

	return item
}

func (item *VersionListItemWidget) Update(name string, isDownloaded bool, downloadUrl string, isOSCompatible bool) {
	item.downloadUrl = downloadUrl

	if !isOSCompatible {
		item.label.Show()
	} else {
		item.label.Hide()
	}

	item.name.SetText(name)

	if isDownloaded {
		item.primaryButton.SetText("Play")
		item.deleteButton.Show()
	} else {
		item.primaryButton.SetText("Download")
		item.deleteButton.Hide()
	}

	if isOSCompatible {
		item.primaryButton.Enable()
	} else {
		item.primaryButton.Disable()
	}
}

func (item *VersionListItemWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewBorder(nil, nil, item.leftContainer, item.rightContainer, nil)
	return widget.NewSimpleRenderer(c)
}
