package files

import (
	"errors"
	"marina/types"
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
		errs = append(errs, cp.Copy(srcSetting, destSetting))
	}

	return errors.Join(errs...)
}
