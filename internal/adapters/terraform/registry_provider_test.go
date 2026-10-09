package terraform

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/providers"
)

func TestRegistryProvider(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		if request.URL.Path != "/v1/providers/hashicorp/aws/versions" {
			t.Fatalf("ruta inesperada: %s", request.URL.Path)
		}
		_ = json.NewEncoder(writer).Encode(map[string]any{
			"versions": []map[string]string{
				{"version": "5.42.0"},
				{"version": "6.0.0-rc1"},
			},
		})
	}))
	defer server.Close()

	provider := NewRegistryProvider(server.URL, server.Client(), nil)
	result, err := provider.AvailableVersions(
		context.Background(),
		model.Dependency{
			Name:   "registry.terraform.io/hashicorp/aws",
			Source: ".terraform.lock.hcl",
		},
		providers.QueryOptions{Mode: providers.NetworkRefresh},
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Versions) != 2 {
		t.Fatalf("se esperaban 2 versiones, se encontraron %d", len(result.Versions))
	}
}
