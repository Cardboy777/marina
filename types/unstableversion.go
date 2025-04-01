package marina

import (
	"errors"
	"fmt"
	"runtime"
	"time"
)

type UnstableVersion struct {
	Hash        string
	ReleaseDate time.Time
	Installed   bool
	Repository  *Repository
}

func (d *UnstableVersion) GetName() string {
	return fmt.Sprintf("Unstable - %s", d.ReleaseDate.Format(time.DateTime))
}

func (d *UnstableVersion) GetDownloadUrl() (string, error) {
	switch {
	// case runtime.GOOS == "linux" && settings.ShouldUseLinuxCompatibilityVersion() && len(d.DownloadUrls.LinuxCompatibility) > 0:
	// 	return d.DownloadUrls.LinuxCompatibility, nil
	case runtime.GOOS == "linux":
		return d.Repository.LatestBuildUrls.Linux, nil
	case runtime.GOOS == "mac":
		return d.Repository.LatestBuildUrls.Mac, nil
	case runtime.GOOS == "windows":
		return d.Repository.LatestBuildUrls.Windows, nil
	}
	return "", errors.New("No compatible Version found")
}

func (d *UnstableVersion) CanDownload() bool {
	_, err := d.GetDownloadUrl()

	return err == nil
}
