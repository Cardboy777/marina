package widgets

import marina "marina/types"

type ListItem struct {
	Release         *marina.Version
	UnstableRelease *marina.UnstableVersion
	IsStableRelease bool
}

func (i ListItem) Installed() bool {
	if i.IsStableRelease {
		return i.Release.Installed
	}
	return i.UnstableRelease.Installed
}
