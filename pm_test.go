package pm

import (
	"fmt"
	"testing"
)

func TestGetHostingServices(t *testing.T) {
	hosts, err := getHostingServices()
	if err != nil {
		t.Fatalf("failed to get hosting services: %v", err)
	}

	for _, h := range hosts {
		fmt.Println(h)
	}
}

func TestGetRepositoryAuthors(t *testing.T) {
	hosts, err := getHostingServices()
	if err != nil {
		t.Fatalf("failed to get hosting services: %v", err)
	}

	for _, h := range hosts {
		fmt.Printf("%s\n", h)
		authors, err := getRepositoryAuthors(h)
		if err != nil {
			t.Fatalf("failed to get autors: %v", err)
		}
		for _, a := range authors {
			fmt.Printf("- %s\n", a)
		}
	}
}

func TestGetRepositories(t *testing.T) {
	hosts, err := getHostingServices()
	if err != nil {
		t.Fatalf("failed to get hosting services: %v", err)
	}

	for _, host := range hosts {
		fmt.Printf("%s\n", host)
		authors, err := getRepositoryAuthors(host)
		if err != nil {
			t.Fatalf("failed to get autors: %v", err)
		}

		for _, author := range authors {
			fmt.Printf("  %s\n", author)
			repos, err := getRepositories(host, author)
			if err != nil {
				t.Fatalf("failed to get repos: %v", err)
			}

			for _, repo := range repos {
				fmt.Printf("    %s %s\n", repo.Name, repo.Version)
				fmt.Printf("    deps:\n")
				for _, dep := range repo.Deps {
					fmt.Printf("      - %s %s\n", dep.Version, dep.Url)
				}
			}
		}
	}

}
