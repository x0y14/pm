package pm

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstall(t *testing.T) {
	url := "demo/distribute-server/userA/repoA/repoA-0.0.1.tar.gz"
	repoAuthorDir := filepath.Join(getPMRoot(), "github.com/userA")
	repoNameVersion := "repoA@v0.0.1"
	tarGzReader, err := os.Open(url)
	if err != nil {
		t.Fatalf("failed to open file %s: %v", url, err)
	}

	err = extractFilesFromTarGz(tarGzReader, repoAuthorDir, repoNameVersion)
	if err != nil {
		t.Fatalf("failed to unpack tar.gz: %v", err)
	}
}
