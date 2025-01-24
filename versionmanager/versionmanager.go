package versionmanager

import (
	"context"
	"strings"
	"marina/settings"
	"marina/types"
	"marina/versionmanager/gamedefinitions"

	"github.com/google/go-github/v68/github"
)

var versionLists = [][]marina.VersionDefinition{}

func fetchVersions(definition *marina.RepositoryDefinition) []marina.VersionDefinition {
	client := github.NewClient(nil)
	ctx := context.Background()
	listRequestOptions := github.ListOptions{
		PerPage: 100,
	}
	list, _, err := client.Repositories.ListReleases(ctx, definition.Owner, definition.Repository, &listRequestOptions)
	if err != nil {
		// probably rate limit reached
		// or offline
	}

	versions := []marina.VersionDefinition{}

	for _, i := range list {

		item := marina.VersionDefinition{
			Name:    (*i).GetName(),
			TagName: (*i).GetTagName(),
		}
		for _, asset := range (*i).Assets {
			name := asset.GetName()
			if asset.GetContentType() != "application/zip" || strings.Contains(name, "Source Code") {
				continue
			}

			switch {
			case strings.Contains(name, "Linux") && !(strings.Contains(name, "Compatibility") || strings.Contains(name, "Performance")), strings.Contains(name, "Compatibility") && settings.ShouldUseLinuxCompatibilityVersion(), strings.Contains(name, "Performance") && !settings.ShouldUseLinuxCompatibilityVersion():
				item.LinuxDownloadUrl = asset.GetBrowserDownloadURL()
			case strings.Contains(name, "Mac"):
				item.MacDownloadUrl = asset.GetBrowserDownloadURL()
			case true:
				item.MacDownloadUrl = asset.GetBrowserDownloadURL()
			}
		}
		versions = append(versions, item)
	}

	return versions
}

func SyncReleases() {
	for _, def := range gamedefinitions.RepositoryDefinitions {
		versions := fetchVersions(def)
		if len(versionLists) <= def.Id {
			versionLists = append(versionLists, versions)
		} else {
			versionLists[def.Id] = versions
		}

	}
	go writeVersionsToManifest()
}

func GetVersionsList() *[][]marina.VersionDefinition {
	return &versionLists
}

func writeVersionsToManifest() {}
