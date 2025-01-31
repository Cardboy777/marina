package files

import (
	"context"
	"fmt"
	"marina/constants"
	"marina/settings"
	"marina/types"
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
				RepositoryDefinition: definition,
				Name:                 (*i).GetName(),
				TagName:              (*i).GetTagName(),
			}
			for _, asset := range (*i).Assets {
				name := asset.GetName()

				fmt.Printf("Name: %s\nContentType: %s\nKeep: %t\n\n", name, asset.GetContentType(), isValidAssetType(asset.GetContentType()))
				if !isValidAssetType(asset.GetContentType()) || strings.Contains(name, "Source Code") {
					continue
				}

				asset.GetCreatedAt()

				switch {
				case isUsableLinuxAsset(name):
					item.LinuxDownloadUrl = asset.GetBrowserDownloadURL()
				case strings.Contains(name, "Mac"):
					item.MacDownloadUrl = asset.GetBrowserDownloadURL()
				case true:
					item.WindowsDownloadUrl = asset.GetBrowserDownloadURL()
				}
			}
			item.ReleaseDate = i.GetCreatedAt().Time
			versions = append(versions, item)
		}
		if resp.NextPage == 0 {
			break
		}
		listRequestOptions.Page = resp.NextPage
	}

	return versions
}

func isValidAssetType(name string) bool {
	return name == "application/zip" || name == "application/x-zip-compressed"
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
	for _, def := range constants.RepositoryDefinitions {
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
