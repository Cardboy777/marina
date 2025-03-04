package marina

import (
	"errors"
	"runtime"
	"time"
)

type Version struct {
	Repository   *Repository
	Name         string
	TagName      string
	DownloadUrls DownloadUrls
	ReleaseDate  time.Time
	Installed    bool
}

func (d *Version) GetDownloadUrl() (string, error) {
	switch {
	// case runtime.GOOS == "linux" && settings.ShouldUseLinuxCompatibilityVersion() && len(d.DownloadUrls.LinuxCompatibility) > 0:
	// 	return d.DownloadUrls.LinuxCompatibility, nil
	case runtime.GOOS == "linux" && len(d.DownloadUrls.Linux) > 0:
		return d.DownloadUrls.Linux, nil
	case runtime.GOOS == "mac" && len(d.DownloadUrls.Mac) > 0:
		return d.DownloadUrls.Mac, nil
	case runtime.GOOS == "windows" && len(d.DownloadUrls.Windows) > 0:
		return d.DownloadUrls.Windows, nil
	}
	return "", errors.New("No compatible Version found")
}

func (d *Version) IsOSCompatible() bool {
	_, err := d.GetDownloadUrl()
	return err == nil
}
