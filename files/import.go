package files

import (
	"io"
	"os"
	"path/filepath"
)

func ImportSettings(src, dest string) error {
	srcSaveDirSettings := filepath.Join(src, "shipofharkinian.json")
	destSaveDirSettings := filepath.Join(dest, "shipofharkinian.json")

	settingsResult := copyFile(srcSaveDirSettings, destSaveDirSettings)

	srcSaveDirImgui := filepath.Join(src, "imgui.ini")
	destSaveDirImgui := filepath.Join(dest, "imgui.ini")

	imguiResult := copyFile(srcSaveDirImgui, destSaveDirImgui)

	if settingsResult != nil {
		return settingsResult
	}

	return imguiResult
}

func ImportRandomizer(src, dest string) error {
	srcSaveDir := filepath.Join(src, "Randomizer")
	destSaveDir := filepath.Join(dest, "Randomizer")

	return copyDir(srcSaveDir, destSaveDir)
}

func ImportSaves(src, dest string) error {
	srcSaveDir := filepath.Join(src, "Save")
	destSaveDir := filepath.Join(dest, "Save")

	return copyDir(srcSaveDir, destSaveDir)
}

func ImportMods(src, dest string) error {
	srcSaveDir := filepath.Join(src, "mods")
	destSaveDir := filepath.Join(dest, "mods")

	return copyDir(srcSaveDir, destSaveDir)
}

func copyFile(src, dest string) error {
	if !exists(src) {
		return nil
	}

	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	return err
}

func copyDir(src, dest string) error {
	if !exists(src) {
		return nil
	}

	srcinfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(dest, srcinfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dest, entry.Name())
		if entry.IsDir() {
			if err = copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err = copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
