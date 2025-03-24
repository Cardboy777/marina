package files

import (
	"errors"
	"marina/types"
	"os"
	"path/filepath"

	cp "github.com/otiai10/copy"
)

func ImportCategory(src, dest string, category marina.ImportCategory) error {
	if !category.Capable {
		return nil
	}

	errs := []error{}
	for _, f := range category.Files {
		srcSetting := filepath.Join(src, f)
		destSetting := filepath.Join(dest, f)

		err := cp.Copy(srcSetting, destSetting)
		if !os.IsNotExist(err) {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
