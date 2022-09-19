package pm

import (
	"fmt"
	"testing"
)

func TestGetGitHubReleases(t *testing.T) {
	releases, err := GetGitHubReleases("dillonkearns", "mobster")
	if err != nil {
		t.Fatalf("failed to get release from gtihub.com: %v", err)
	}

	for _, rel := range releases {
		fmt.Printf("%s: %s\n", rel.TagName, rel.TarGzUrl)
	}
}

func TestGetGitHubLatestRelease(t *testing.T) {
	release, err := GetGitHubLatestRelease("dillonkearns", "mobster")
	if err != nil {
		t.Fatalf("failed to get release from gtihub.com: %v", err)
	}

	fmt.Printf("%s: %s\n", release.TagName, release.TarGzUrl)
}
