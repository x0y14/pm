package pm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getPMRoot() string {
	return "demo/packages"
}

func getHostingServices() ([]string, error) {
	root := getPMRoot()

	var hosts []string

	files, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() && f.Name() != root {
			hosts = append(hosts, f.Name())
		}
	}

	return hosts, nil
}

func getRepositoryAuthors(hostingService string) ([]string, error) {
	hsDir := filepath.Join(getPMRoot(), hostingService)

	var authors []string

	files, err := os.ReadDir(hsDir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() && f.Name() != hsDir {
			authors = append(authors, f.Name())
		}
	}

	return authors, nil
}

func getRepositoryNameVersions(hostingService string, author string) ([]string, error) {
	hosting := filepath.Join(getPMRoot(), hostingService)

	authorDir := filepath.Join(hosting, author)

	files, err := os.ReadDir(authorDir)
	if err != nil {
		return nil, err
	}

	var repos []string

	for _, repo := range files {
		if repo.IsDir() {
			repos = append(repos, repo.Name())
		}
	}

	return repos, nil
}

func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getRepository(hostingService, author, repoNameVersion string) (*Repository, error) {
	if !strings.Contains(repoNameVersion, "@") {
		return nil, fmt.Errorf("invalid repository id format, expect: 'name@version', but reserve %s", repoNameVersion)
	}

	pkgJsonPath := filepath.Join(getPMRoot(), hostingService, author, repoNameVersion, "pkg.json")
	if !IsFileExists(pkgJsonPath) {
		return nil, fmt.Errorf("file not found: %s", pkgJsonPath)
	}
	bytes, err := os.ReadFile(pkgJsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %v", pkgJsonPath, err)
	}

	pkgJson, err := UnmarshalPkgJson(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %v", pkgJsonPath, err)
	}

	nameVersion := strings.Split(repoNameVersion, "@")
	repoName := nameVersion[0]
	repoVersion := nameVersion[1]

	return &Repository{
		Host:    hostingService,
		Author:  author,
		Name:    repoName,
		Version: repoVersion,
		Deps:    pkgJson.Deps,
	}, nil
}

func getRepositories(hostingService, author string) ([]*Repository, error) {
	repoNameVersions, err := getRepositoryNameVersions(hostingService, author)
	if err != nil {
		return nil, fmt.Errorf("failed to get repositoryIds: %v", err)
	}

	var repos []*Repository

	for _, repoNV := range repoNameVersions {
		repo, err := getRepository(hostingService, author, repoNV)
		if err != nil {
			return nil, fmt.Errorf("failed to get repository data: %v", err)
		}
		repos = append(repos, repo)
	}

	return repos, nil
}

func GetInstalledRepositories() ([]*Repository, error) {
	hostingServices, err := getHostingServices()
	if err != nil {
		return nil, fmt.Errorf("failed to get hosting service list: %v", err)
	}

	var allRepo []*Repository

	for _, hostingService := range hostingServices {
		authors, err := getRepositoryAuthors(hostingService)
		if err != nil {
			return nil, fmt.Errorf("failed to get repository authors: %v", err)
		}
		for _, author := range authors {
			repos, err := getRepositories(hostingService, author)
			if err != nil {
				return nil, fmt.Errorf("failed to get installed repositories of %s/%s: %v", hostingService, author, err)
			}
			allRepo = append(allRepo, repos...)
		}
	}

	return allRepo, nil
}
