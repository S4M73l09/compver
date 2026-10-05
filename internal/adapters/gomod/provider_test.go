package gomod

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/S4M73l09/compver/internal/model"
	"github.com/S4M73l09/compver/internal/providers"
)

func TestProviderAvailableVersions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			if request.URL.Path != "/example.com/library/@v/list" {
				t.Fatalf("ruta inesperada: %s", request.URL.Path)
			}

			_, _ = writer.Write([]byte(
				"v1.0.0\nv1.1.0-beta.1\nv1.2.0\n",
			))
		},
	))
	defer server.Close()

	provider := NewProvider(server.URL, server.Client())
	result, err := provider.AvailableVersions(
		context.Background(),
		model.Dependency{Name: "example.com/library"},
		providers.QueryOptions{Mode: providers.NetworkAuto},
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Versions) != 3 {
		t.Fatalf(
			"se esperaban 3 versiones, se obtuvieron %d",
			len(result.Versions),
		)
	}
}
