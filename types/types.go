package marina

import "runtime"

type RepositoryDefinition struct {
	Id                int
	Name              string
	Owner             string
	Repository        string
	AcceptedRomHashes []string
}

type VersionDefinition struct {
	Name               string
	TagName            string
	WindowsDownloadUrl string
	LinuxDownloadUrl   string
	MacDownloadUrl     string
}

func (*VersionDefinition) IsDownloaded() bool {
	return false
}

func (d *VersionDefinition) GetDownloadUrl() string {
	switch {
	case runtime.GOOS == "linux":
		return d.LinuxDownloadUrl
	case runtime.GOOS == "mac":
		return d.MacDownloadUrl
	case runtime.GOOS == "windows":
		return d.WindowsDownloadUrl
	}
	return ""
}

func (d *VersionDefinition) IsOSCompatible() bool {
	return len(d.GetDownloadUrl()) > 0
}

func (rd *RepositoryDefinition) HasRoms() bool {
	return false
}
