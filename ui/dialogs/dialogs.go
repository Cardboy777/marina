package dialogs

import (
	"errors"

	"github.com/skratchdot/open-golang/open"
	"github.com/sqweek/dialog"
)

func ShowDialogIfError(err error) {
	if err != nil {
		ShowErrorDialog(err)
	}
}

func ShowErrorDialog(err error) {
	dialog.Message("%s", err).Title("Encounterd an Error").Error()
}

func ShowConfirmDialog(title string, message string) bool {
	ok := dialog.Message("%s", message).Title(title).YesNo()
	return ok
}

func ShowFilePickerDialogFiltered(title string, fileTypesDescription string, filetypes []string) (string, error) {
	file, err := dialog.File().Filter(fileTypesDescription, filetypes...).Title(title).Load()

	if err != nil && !errors.Is(err, dialog.ErrCancelled) {
		return "", err
	}

	return file, nil
}

func ShowDirectoryPickerDialog(title string) (string, error) {
	folder, err := dialog.Directory().Title(title).Browse()

	if err != nil && !errors.Is(err, dialog.ErrCancelled) {
		return "", err
	}

	return folder, nil
}

func OpenDirectory(path string) {
	_ = open.Start(path)
}
