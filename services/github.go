package services

import (
	"context"
	"fmt"
	"marina/types"
	"strings"

	"github.com/google/go-github/v68/github"
)

func fetchReleaseVersions(definition *marina.Repository) ([]marina.Version, error) {
	client := github.NewClient(nil)
	ctx := context.Background()
	listRequestOptions := github.ListOptions{}
	versions := []marina.Version{}

	for {
		list, resp, err := client.Repositories.ListReleases(ctx, definition.Owner, definition.Repository, &listRequestOptions)
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("Rate Limit for Public api reached, Try again later.")
		}
		if err != nil {
			return nil, fmt.Errorf("Error Accessing GitHub: %w", err)
		}

		for _, i := range list {

			item := marina.Version{
				Repository: definition,
				Name:       (*i).GetName(),
				TagName:    (*i).GetTagName(),
			}
			for _, asset := range (*i).Assets {
				name := asset.GetName()

				if !isValidAssetType(asset.GetContentType()) || strings.Contains(name, "Source Code") {
					continue
				}

				asset.GetCreatedAt()

				containsLinux := strings.Contains(name, "Linux")
				containsCompatibility := strings.Contains(name, "Compatibility")

				switch {
				case containsLinux && containsCompatibility:
					item.DownloadUrls.LinuxCompatibility = asset.GetBrowserDownloadURL()
				case containsLinux && !containsCompatibility:
					item.DownloadUrls.Linux = asset.GetBrowserDownloadURL()
				case strings.HasSuffix(name, "Mac.zip"):
					item.DownloadUrls.Mac = asset.GetBrowserDownloadURL()
				case true:
					item.DownloadUrls.Windows = asset.GetBrowserDownloadURL()
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

	return versions, nil
}

func fetchLatestCommit(definition *marina.Repository) (*marina.UnstableVersion, error) {
	client := github.NewClient(nil)
	ctx := context.Background()

	repo, _, err := client.Repositories.Get(ctx, definition.Owner, definition.Repository)
	if _, ok := err.(*github.RateLimitError); ok {
		return nil, fmt.Errorf("Rate Limit for Public api reached, Try again later.")
	}
	if err != nil {
		return nil, fmt.Errorf("Error Accessing GitHub: %w", err)
	}

	branchName := repo.GetDefaultBranch()

	branch, _, err := client.Repositories.GetBranch(ctx, definition.Owner, definition.Repository, branchName, 5)
	if _, ok := err.(*github.RateLimitError); ok {
		return nil, fmt.Errorf("Rate Limit for Public api reached, Try again later.")
	}
	if err != nil {
		panic(fmt.Errorf("Error Accessing GitHub: %w", err))
	}

	latestCommit := branch.GetCommit()

	return &marina.UnstableVersion{
		Repository:  definition,
		Hash:        *latestCommit.SHA,
		ReleaseDate: latestCommit.Commit.Committer.GetDate().Time,
	}, nil
}

func isValidAssetType(name string) bool {
	return name == "application/zip" || name == "application/x-zip-compressed"
}
