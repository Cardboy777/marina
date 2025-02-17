package ui

import (
	"errors"

	"github.com/skratchdot/open-golang/open"
	"github.com/sqweek/dialog"
)

func ShowErrorDialog(err error) {
	dialog.Message("%s", err).Title("Encounterd an Error").Error()
}

func ShowConfirmDialog(title string, message string, callbackFn func(bool)) {
	ok := dialog.Message("%s", message).Title(title).YesNo()
	callbackFn(ok)
}

func ShowFilePickerDialogFiltered(title string, fileTypesDescription string, filetypes []string, callbackFn func(string, error)) {
	file, err := dialog.File().Filter(fileTypesDescription, filetypes...).Title(title).Load()

	if err != nil && !errors.Is(err, dialog.ErrCancelled) {
		callbackFn("", err)
	}
	callbackFn(file, nil)
}

func OpenDirectory(path string) {
	_ = open.Start(path)
}
