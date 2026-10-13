package terraform

import (
	"testing"
	"time"

	"github.com/S4M73l09/compver/internal/model"
)

func TestProviderLifecycleCanHandle(t *testing.T) {
	provider := NewProviderLifecycleProvider("https://api.github.com", nil, nil)

	if !provider.CanHandle(model.Dependency{
		Name: "registry.terraform.io/hashicorp/aws",
	}) {
		t.Fatal("se esperaba que el provider gestionara un provider del Registry")
	}

	if provider.CanHandle(model.Dependency{Name: "terraform"}) {
		t.Fatal("no se esperaba gestionar Terraform CLI")
	}
}

func TestParseGitHubReleases(t *testing.T) {
	publishedAt := time.Date(2026, time.January, 10, 12, 0, 0, 0, time.UTC)

	releases, err := parseGitHubReleases([]githubRelease{
		{
			TagName:     "v5.31.0",
			PublishedAt: &publishedAt,
		},
		{
			TagName:    "v6.0.0-rc.1",
			Prerelease: true,
		},
		{
			TagName: "draft",
			Draft:   true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(releases) != 2 {
		t.Fatalf("se esperaban dos releases, se obtuvieron %d", len(releases))
	}

	if releases[0].Support != model.SupportUnknown {
		t.Fatalf("soporte inesperado: %s", releases[0].Support)
	}

	if releases[0].PublishedAt == nil {
		t.Fatal("se esperaba la fecha de publicación")
	}
}
