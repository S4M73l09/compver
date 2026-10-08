package cache

import (
	"testing"
	"time"

	"github.com/S4M73l09/compver/internal/version"
)

func TestFileCacheSetGetDelete(t *testing.T) {
	fileCache, err := NewFileCache(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	parsedVersion, err := version.Parse("1.2.3")
	if err != nil {
		t.Fatal(err)
	}

	entry := Entry{
		Versions:    []version.Version{parsedVersion},
		Source:      "test-provider",
		RetrievedAt: time.Now(),
	}

	if err := fileCache.Set("test:dependency", entry); err != nil {
		t.Fatal(err)
	}

	stored, found, err := fileCache.Get("test:dependency")
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatal("se esperaba encontrar la entrada en caché")
	}
	if len(stored.Versions) != 1 {
		t.Fatalf("se esperaba 1 versión, se obtuvieron %d", len(stored.Versions))
	}
	if stored.Source != entry.Source {
		t.Fatalf("fuente inesperada: %s", stored.Source)
	}

	if err := fileCache.Delete("test:dependency"); err != nil {
		t.Fatal(err)
	}

	_, found, err = fileCache.Get("test:dependency")
	if err != nil {
		t.Fatal(err)
	}
	if found {
		t.Fatal("no se esperaba encontrar la entrada eliminada")
	}
}
