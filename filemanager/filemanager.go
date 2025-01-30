package filemanager

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"marina/constants"
	"marina/rommanager"
	"marina/settings"
	"marina/types"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

var dirName string

var manifest = []marina.ManifestItem{}

const (
	permission           = 0o700
	executablePermission = 0o700
)

func Init() {
	dirName = getInstallDirName()

	err := os.MkdirAll(dirName, permission)
	if err != nil {
		panic(fmt.Errorf("Cannot create install dir: %w", err))
	}

	_ = readManifestFile()
}

func getInstallDirName() string {
	return filepath.Join(settings.GetInstallDir(), marina.DirName(constants.AppName))
}

func getRomFileName(rom marina.RomDefinition) string {
	return fmt.Sprintf("%s.z64", marina.DirName(rom.Name))
}

func getRomInstallDir() string {
	return filepath.Join(getInstallDirName(), "roms")
}

func getRomPath(rom marina.RomDefinition) string {
	return filepath.Join(getRomInstallDir(), getRomFileName(rom))
}

func GetVersionInstallDirPath(version *marina.VersionDefinition) string {
	return filepath.Join(getInstallDirName(), "versions", marina.DirName(version.RepositoryDefinition.Repository), version.GetVersionDirName())
}

func IsValidRomInstalled(repo *marina.RepositoryDefinition) (bool, *[]marina.RomDefinition) {
	for _, item := range manifest {
		if item.Owner == repo.Owner && item.Repository == repo.Repository && len(*item.InstalledRoms) > 0 {
			return true, item.InstalledRoms
		}
	}
	return false, nil
}

func CopyRomsToVersionInstall(version *marina.VersionDefinition) {
	var installedRoms *[]marina.RomDefinition

	for _, m := range manifest {
		if version.RepositoryDefinition.Repository == m.Repository && version.RepositoryDefinition.Owner == m.Owner {
			installedRoms = m.InstalledRoms
			break
		}
	}
	if installedRoms == nil {
		return
	}

	dirName := GetVersionInstallDirPath(version)
	for _, r := range *installedRoms {
		romPath := getRomPath(r)
		err := os.Link(romPath, filepath.Join(dirName, getRomFileName(r)))
		if err != nil && !errors.Is(err, os.ErrExist) {
			panic(fmt.Errorf("Error linking rom to install directory: %w", err))
		}
	}
}

func IsVersionInstalled(version *marina.VersionDefinition) bool {
	for _, item := range manifest {
		if item.Owner == version.RepositoryDefinition.Owner && item.Repository == version.RepositoryDefinition.Repository && item.InstalledTagNames != nil {
			return slices.Contains(*item.InstalledTagNames, version.TagName)
		}
	}
	return false
}

func DeleteVersion(version *marina.VersionDefinition) error {
	path := GetVersionInstallDirPath(version)
	err := os.RemoveAll(path)
	if err != nil {
		log.Printf("Cannot delete version: %s", err)
	}

	removeVersionFromManifest(version)

	return err
}

func DownloadVersion(version *marina.VersionDefinition) error {
	path := GetVersionInstallDirPath(version)

	zipPath := filepath.Join(path, "download.zip")

	downloadUrl, err := version.GetDownloadUrl()
	if err != nil {
		return err
	}

	err = os.MkdirAll(path, permission)
	if err != nil {
		log.Printf("Cannot create version directory: %s", err)
		return err
	}

	out, err := os.Create(zipPath)
	if err != nil {
		log.Printf("Error downloading file: %s\n", err)
		return err
	}
	defer out.Close()

	resp, err := http.Get(downloadUrl)
	if err != nil {
		log.Printf("Error downloading file: %s\n", err)
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Printf("Error downloading file: %s\n", err)
		return err
	}

	err = unzip(zipPath, path)
	if err != nil {
		log.Printf("Error unzipping version: %s\n", err)
		return err
	}

	addVersionToManifest(version)

	return nil
}

func addRomToInstalledVersions(rom *marina.RomDefinition, repo *marina.RepositoryDefinition) {
	var item *marina.ManifestItem

	for _, m := range manifest {
		if m.Repository == repo.Repository && m.Owner == repo.Owner {
			item = &m
			break
		}
	}
	if item == nil {
		return
	}

	romFileName := getRomFileName(*rom)
	romPath := getRomPath(*rom)
	installPath := getInstallDirName()
	relativePaths := item.GetVersionInstallRelativePaths()

	for _, path := range relativePaths {
		fullPath := filepath.Join(installPath, "versions", path, romFileName)
		err := os.Link(romPath, fullPath)
		if err != nil {
			panic(fmt.Errorf("Error copying rom file: %w", err))
		}
	}
}

func CopyRomToInstallDir(repo *marina.RepositoryDefinition, sourcePath string) error {
	bytes, rom, isValid := rommanager.IsValidRom(repo.AcceptedRomHashes, sourcePath)
	if !isValid {
		return errors.New("Invalid ROM")
	}

	err := os.MkdirAll(getRomInstallDir(), permission)
	if err != nil {
		panic(fmt.Errorf("Cannot create version directory: %w", err))
	}

	filename := getRomPath(rom)

	err = os.WriteFile(filename, *bytes, permission)
	if err != nil {
		panic(fmt.Errorf("Error copying rom file: %w", err))
	}

	addRomToManifest(&rom, repo)

	go addRomToInstalledVersions(&rom, repo)

	return nil
}

func GetInstalledRoms(repo *marina.RepositoryDefinition) *[]marina.RomDefinition {
	for _, m := range manifest {
		if m.Owner == repo.Owner && m.Repository == repo.Repository {
			return m.InstalledRoms
		}
	}

	return nil
}

func IsExecutable(fileName string) bool {
	switch {
	case runtime.GOOS == "linux" && strings.HasSuffix(fileName, ".appimage"):
		return true
	case runtime.GOOS == "mac" && strings.HasSuffix(fileName, ".dmg"):
		return true
	case runtime.GOOS == "windows" && strings.HasSuffix(fileName, ".exe"):
		return true
	}
	return false
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(path, f.Mode())
			if err != nil {
				return err
			}
		} else {
			err := os.MkdirAll(filepath.Dir(path), f.Mode())
			if err != nil {
				return err
			}

			isExe := IsExecutable(f.Name)
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())

			if err == nil && isExe {
				err = os.Chmod(path, permission)
			}

			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func addRomToManifest(rom *marina.RomDefinition, repo *marina.RepositoryDefinition) {
	hasUpdated := false
	for _, item := range manifest {
		if item.Owner == repo.Owner && item.Repository == repo.Repository {
			if item.InstalledRoms == nil {
				item.InstalledRoms = &[]marina.RomDefinition{*rom}
			} else {
				*item.InstalledRoms = append(*item.InstalledRoms, *rom)
			}
			hasUpdated = true
			break
		}
	}
	if !hasUpdated {
		manifest = append(manifest, marina.ManifestItem{
			Owner:             repo.Owner,
			Repository:        repo.Repository,
			InstalledRoms:     &[]marina.RomDefinition{*rom},
			InstalledTagNames: &[]string{},
		})
	}

	go updateManifestFile()
}

func removeVersionFromManifest(version *marina.VersionDefinition) {
	for _, item := range manifest {
		if item.Owner == version.RepositoryDefinition.Owner && item.Repository == version.RepositoryDefinition.Repository && item.InstalledTagNames != nil {
			*item.InstalledTagNames = slices.DeleteFunc(*item.InstalledTagNames, func(s string) bool {
				return s == version.TagName
			})
			break
		}
	}

	go updateManifestFile()
}

func addVersionToManifest(version *marina.VersionDefinition) {
	hasUpdated := false
	for _, item := range manifest {
		if item.Owner == version.RepositoryDefinition.Owner && item.Repository == version.RepositoryDefinition.Repository {
			if item.InstalledTagNames == nil {
				item.InstalledTagNames = &[]string{version.TagName}
			} else {
				*item.InstalledTagNames = append(*item.InstalledTagNames, version.TagName)
			}
			hasUpdated = true
			break
		}
	}
	if !hasUpdated {
		manifest = append(manifest, marina.ManifestItem{
			Owner:             version.RepositoryDefinition.Owner,
			Repository:        version.RepositoryDefinition.Repository,
			InstalledRoms:     &[]marina.RomDefinition{},
			InstalledTagNames: &[]string{version.TagName},
		})
	}

	go updateManifestFile()
}

func getManifestFilePath() string {
	return fmt.Sprintf("%s/manifest.json", getInstallDirName())
}

func readManifestFile() error {
	file, err := os.Open(getManifestFilePath())
	if err != nil {
		fmt.Printf("Error reading manifest.json: %s", err)
		return err
	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)

	err = json.Unmarshal(bytes, &manifest)
	if err != nil {
		fmt.Printf("Error reading manifest.json: %s", err)
		return err
	}

	return nil
}

func updateManifestFile() {
	bytes, err := json.Marshal(manifest)
	if err != nil {
		panic(fmt.Errorf("Unable to save manifest file: %w", err))
	}

	// fmt.Print(string(bytes))

	err = os.WriteFile(getManifestFilePath(), bytes, permission)
	if err != nil {
		panic(fmt.Errorf("Unable to save manifest file: %w", err))
	}
}
