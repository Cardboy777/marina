package marina

import (
	"errors"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type RomDefinition struct {
	Name string
	Sha1 string
}

type RepositoryDefinition struct {
	Id                int
	Name              string
	Owner             string
	Repository        string
	PathVariableName  string
	AcceptedRomHashes *[]RomDefinition
}

type VersionDefinition struct {
	RepositoryDefinition *RepositoryDefinition
	Name                 string
	TagName              string
	ReleaseDate          time.Time
	WindowsDownloadUrl   string
	LinuxDownloadUrl     string
	MacDownloadUrl       string
}

type ManifestItem struct {
	Owner             string
	Repository        string
	InstalledRoms     *[]RomDefinition
	InstalledTagNames *[]string
}

func (d *VersionDefinition) IsDownloaded() bool {
	return false
}

func (d *VersionDefinition) GetDownloadUrl() (string, error) {
	switch {
	case runtime.GOOS == "linux" && len(d.LinuxDownloadUrl) > 0:
		return d.LinuxDownloadUrl, nil
	case runtime.GOOS == "mac" && len(d.MacDownloadUrl) > 0:
		return d.MacDownloadUrl, nil
	case runtime.GOOS == "windows" && len(d.WindowsDownloadUrl) > 0:
		return d.WindowsDownloadUrl, nil
	}
	return "", errors.New("No compatible Version found")
}

func (d *VersionDefinition) GetVersionDirName() string {
	return DirName(d.TagName)
}

func (d *VersionDefinition) IsOSCompatible() bool {
	_, err := d.GetDownloadUrl()
	return err == nil
}

func (d *VersionDefinition) GetVersionInstallDirPath(baseDir string) string {
	return filepath.Join(baseDir, d.GetVersionDirName())
}

func (rd *RepositoryDefinition) HasRoms() bool {
	return false
}

func (m *ManifestItem) GetVersionInstallRelativePaths() []string {
	paths := []string{}

	for _, tagName := range *m.InstalledTagNames {
		path := filepath.Join(DirName(m.Repository), DirName(tagName))
		paths = append(paths, path)
	}

	return paths
}

func DirName(name string) string {
	replacer := strings.NewReplacer(" ", "-", ".", "_", "(", "", ")", "")
	return replacer.Replace(strings.ToLower(name))
}
