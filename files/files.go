package files

import (
	"errors"
	"fmt"
	"io"
	"marina/constants"
	"marina/settings"
	"marina/stores"
	"marina/types"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var dirName string

func DirName(name string) string {
	replacer := strings.NewReplacer(" ", "-", ".", "_", "(", "", ")", "", ":", "-")
	return replacer.Replace(strings.ToLower(name))
}

func Init() {
	dirName = settings.GetInstallDirName()

	err := os.MkdirAll(dirName, constants.FilePermission)
	if err != nil {
		panic(fmt.Errorf("Cannot create install dir: %w", err))
	}
}

func getRomFileName(rom marina.Rom) string {
	return fmt.Sprintf("%s.z64", DirName(rom.Name))
}

func getRomInstallDir() string {
	return filepath.Join(settings.GetInstallDirName(), "roms")
}

func getRomPath(rom marina.Rom) string {
	return filepath.Join(getRomInstallDir(), getRomFileName(rom))
}

func GetVersionInstallDirPath(version *marina.Version) string {
	return filepath.Join(settings.GetInstallDirName(), "versions", DirName(version.Repository.Repository), DirName(version.TagName))
}

func GetUnstableVersionInstallDirPath(version *marina.UnstableVersion) string {
	unstableTimeName := fmt.Sprintf("unstable-%s", version.ReleaseDate.Format(time.DateTime))

	return filepath.Join(settings.GetInstallDirName(), "versions", DirName(version.Repository.Repository), DirName(unstableTimeName))
}

func IsValidRomInstalled(repo *marina.Repository) (bool, *[]marina.Rom) {
	roms := stores.GetInstalledRomsList(repo)
	return len(*roms) > 0, roms
}

func CopyRomsToVersionInstall(version *marina.Version) error {
	hasRoms, installedRoms := IsValidRomInstalled(version.Repository)

	if !hasRoms {
		return fmt.Errorf("No valid ROMs are installed")
	}

	dirName := GetVersionInstallDirPath(version)
	for _, r := range *installedRoms {
		romPath := getRomPath(r)
		dest := filepath.Join(dirName, getRomFileName(r))
		err := os.Link(romPath, dest)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return (fmt.Errorf("Error linking rom to install directory: %w", err))
		}
	}

	return nil
}

func CopyRomsToUnstableVersionInstall(version *marina.UnstableVersion) error {
	hasRoms, installedRoms := IsValidRomInstalled(version.Repository)

	if !hasRoms {
		return fmt.Errorf("No valid ROMs are installed")
	}

	dirName := GetUnstableVersionInstallDirPath(version)
	for _, r := range *installedRoms {
		romPath := getRomPath(r)
		dest := filepath.Join(dirName, getRomFileName(r))
		err := os.Link(romPath, dest)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return (fmt.Errorf("Error linking rom to install directory: %w", err))
		}
	}

	return nil
}

func DeleteVersion(version *marina.Version) error {
	path := GetVersionInstallDirPath(version)
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Printf("Cannot delete version: %s", err)
	}

	version.Installed = false
	stores.SetVersionInstalled(version, false)

	return err
}

func DeleteUnstableVersion(version *marina.UnstableVersion) error {
	path := GetUnstableVersionInstallDirPath(version)
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Printf("Cannot delete version: %s", err)
	}

	version.Installed = false
	stores.SetUnstableVersionInstalled(version, false)

	return err
}

func DownloadVersion(version *marina.Version) error {
	path := GetVersionInstallDirPath(version)

	zipPath := filepath.Join(path, "download.zip")

	downloadUrl, err := version.GetDownloadUrl()
	if err != nil {
		return err
	}

	err = os.MkdirAll(path, constants.FilePermission)
	if err != nil {
		fmt.Printf("Cannot create version directory: %s", err)
		return err
	}

	out, err := os.Create(zipPath)
	if err != nil {
		fmt.Printf("Error downloading file: %s\n", err)
		return err
	}
	defer out.Close()

	resp, err := http.Get(downloadUrl)
	if err != nil {
		fmt.Printf("Error downloading file: %s\n", err)
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Error downloading file: %s\n", err)
		return err
	}

	err = Unzip(zipPath, path)
	if err != nil {
		fmt.Printf("Error unzipping version: %s\n", err)
		return err
	}

	version.Installed = true
	stores.SetVersionInstalled(version, true)

	return nil
}

func DownloadUnstableVersion(version *marina.UnstableVersion) error {
	path := GetUnstableVersionInstallDirPath(version)

	zipPath := filepath.Join(path, "download.zip")

	downloadUrl, err := version.GetDownloadUrl()
	if err != nil {
		return err
	}

	err = os.MkdirAll(path, constants.FilePermission)
	if err != nil {
		fmt.Printf("Cannot create version directory: %s", err)
		return err
	}

	out, err := os.Create(zipPath)
	if err != nil {
		fmt.Printf("Error downloading file: %s\n", err)
		return err
	}
	defer out.Close()

	resp, err := http.Get(downloadUrl)
	if err != nil {
		fmt.Printf("Error downloading file: %s\n", err)
		return err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Error downloading file: %s\n", err)
		return err
	}
	defer resp.Body.Close()

	err = Unzip(zipPath, path)
	if err != nil {
		fmt.Printf("Error unzipping unstable version: %s\n", err)
		return err
	}

	version.Installed = true
	stores.SetUnstableVersionInstalled(version, true)

	return nil
}

func CopyRomToInstallDir(repo *marina.Repository, sourcePath string) error {
	bytes, rom, isValid := IsValidRom(repo.AcceptedRomHashes, sourcePath)
	if !isValid {
		return errors.New("Invalid ROM")
	}

	err := os.MkdirAll(getRomInstallDir(), constants.FilePermission)
	if err != nil {
		panic(fmt.Errorf("Cannot create version directory: %w", err))
	}

	filename := getRomPath(rom)

	err = os.WriteFile(filename, *bytes, constants.FilePermission)
	if err != nil {
		panic(fmt.Errorf("Error copying rom file: %w", err))
	}

	stores.AddInstalledRom(rom, repo)

	return nil
}

func IsExecutable(file os.FileInfo) bool {
	if file.Mode().IsDir() {
		return false
	}

	switch {
	case runtime.GOOS == "linux" && strings.HasSuffix(file.Name(), ".appimage"):
		return true
	case runtime.GOOS == "mac" && strings.HasSuffix(file.Name(), ".dmg"):
		return true
	case runtime.GOOS == "windows" && strings.HasSuffix(file.Name(), ".exe"):
		return true
	}

	return false
}
