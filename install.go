package pm

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func GetGitHubReleases(author, repoName string) ([]*Release, error) {
	var githubReleases []*GithubRelease
	var releases []*Release
	for i := 1; i < 100; i++ {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases?page=%d", author, repoName, i)

		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(bodyBytes, &githubReleases)
		if err != nil {
			return nil, err
		}

		if len(githubReleases) == 0 {
			break
		}

		for _, gRelease := range githubReleases {
			releases = append(releases, &Release{
				TagName:  gRelease.TagName,
				TarGzUrl: gRelease.TarballUrl,
			})
		}
		resp.Body.Close()
	}

	return releases, nil
}

func GetGitHubLatestRelease(author, repoName string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", author, repoName)

	var githubLatestRelease GithubRelease

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &githubLatestRelease)
	if err != nil {
		return nil, err
	}

	return &Release{
		TagName:  githubLatestRelease.TagName,
		TarGzUrl: githubLatestRelease.TarballUrl,
	}, nil
}

// extractFilesFromTarGz repoAuthorDirPathは取り出されたファイルが設置される場所なので
// PMRoot/host/userで指定すること
// なお、一時的なディレクトリネームを使用して設置したのち、正規の名前に変更する
// tag_name -> "vx.y.z
// regex ? "tarball_url": "https://api.github.com/repos/dillonkearns/mobster/tarball/v0.0.48",
func extractFilesFromTarGz(tarGzReader io.Reader, repoAuthorDirPath string, repoNameVersion string) error {
	gzReader, err := gzip.NewReader(tarGzReader)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	var header *tar.Header

	tempInstallDir := repoAuthorDirPath

	var installedDir string

	for {
		header, err = tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to Next(): %v", err)
		}

		fileName := filepath.Join(tempInstallDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(fileName, 0755); err != nil {
				return fmt.Errorf("extractTarGz: Mkdir() failed: %v", err)
			}
			installedDir = header.Name
		case tar.TypeReg:
			outFile, err := os.Create(fileName)
			if err != nil {
				return fmt.Errorf("extractTarGz: Create() failed: %v", err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				_ = outFile.Close()
				return fmt.Errorf("extractTarGz: Copy() failed: %v", err)
			}

			if err := outFile.Close(); err != nil {
				return fmt.Errorf("extractTarGz: Close() failed: %v", err)
			}

			//if strings.Contains(header.Name, "pkg.json") {
			//	j, err := os.ReadFile(tempInstallDir)
			//	if err != nil {
			//		return fmt.Errorf("failed to read file: %v", err)
			//	}
			//}
		case tar.TypeXGlobalHeader:
			continue

		default:
			return fmt.Errorf("extractTarGz: uknown type: %b in %s", header.Typeflag, header.Name)
		}
	}

	realName := filepath.Join(repoAuthorDirPath, repoNameVersion)
	err = os.Rename(filepath.Join(tempInstallDir, installedDir), realName)
	if err != nil {
		return err
	}
	return nil
}
