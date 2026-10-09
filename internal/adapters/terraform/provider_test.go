package terraform

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/providers"
)

func TestAvailableVersions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		writer.Header().Set("Content-Type", "text/html")
		_, _ = writer.Write([]byte(`
<a href="terraform_1.8.5/">terraform_1.8.5</a>
<a href="terraform_1.9.0-rc1/">terraform_1.9.0-rc1</a>
<a href="terraform_1.9.0/">terraform_1.9.0</a>
`))
	}))
	defer server.Close()

	provider := NewProvider(server.URL, server.Client(), nil)
	result, err := provider.AvailableVersions(
		context.Background(),
		modelDependency(),
		providers.QueryOptions{Mode: providers.NetworkRefresh},
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Versions) != 3 {
		t.Fatalf("se esperaban 3 versiones, se encontraron %d", len(result.Versions))
	}

	if !strings.Contains(result.Versions[1].String(), "rc1") {
		t.Fatalf("no se encontró la versión rc esperada: %v", result.Versions)
	}
}

func modelDependency() model.Dependency {
	return model.Dependency{
		Name:   "terraform",
		Source: ".terraform-version",
	}
}
