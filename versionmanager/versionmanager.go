package versionmanager

import (
	"context"
	"fmt"
	"marina/settings"
	"marina/types"
	"marina/versionmanager/gamedefinitions"
	"strings"

	"github.com/google/go-github/v68/github"
)

var versionLists = [][]marina.VersionDefinition{}

func fetchVersions(definition *marina.RepositoryDefinition) []marina.VersionDefinition {
	versions := []marina.VersionDefinition{}

	client := github.NewClient(nil)
	ctx := context.Background()
	listRequestOptions := github.ListOptions{}

	for {
		list, resp, err := client.Repositories.ListReleases(ctx, definition.Owner, definition.Repository, &listRequestOptions)
		if err != nil {
			// probably rate limit reached
			// or offline
			panic(fmt.Errorf("Error Accessing GitHub: %w", err))
		}

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

				// fmt.Printf("Name: %s\nContentType: %s\nIsCompatible: %t\nURL: %s\n\n", name, asset.GetContentType(), isUsableLinuxAsset(name), asset.GetBrowserDownloadURL())

				switch {
				case isUsableLinuxAsset(name):
					item.LinuxDownloadUrl = asset.GetBrowserDownloadURL()
				case strings.Contains(name, "Mac"):
					item.MacDownloadUrl = asset.GetBrowserDownloadURL()
				case true:
					item.MacDownloadUrl = asset.GetBrowserDownloadURL()
				}
			}
			versions = append(versions, item)
		}
		if resp.NextPage == 0 {
			break
		}
		listRequestOptions.Page = resp.NextPage
	}

	return versions
}

func isUsableLinuxAsset(name string) bool {
	if !strings.Contains(name, "Linux") {
		return false
	}
	if settings.ShouldUseLinuxCompatibilityVersion() && strings.Contains(name, "Compatibility") {
		return true
	}
	if !settings.ShouldUseLinuxCompatibilityVersion() && strings.Contains(name, "Performance") {
		return true
	}
	return !strings.Contains(name, "Compatibility") && !strings.Contains(name, "Performance")
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
