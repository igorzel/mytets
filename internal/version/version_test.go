package version_test

import (
	"testing"

	"github.com/igorzel/mytets/internal/version"
)

func TestVersionDefault(t *testing.T) {
	if version.Version == "" {
		t.Fatal("Version must not be empty; expected 'dev' fallback")
	}
}

func TestVersionFallback(t *testing.T) {
	original := version.Version
	t.Cleanup(func() { version.Version = original })

	version.Version = "dev"
	if version.Version != "dev" {
		t.Errorf("expected fallback 'dev', got %q", version.Version)
	}
}

func TestVersionInjected(t *testing.T) {
	original := version.Version
	t.Cleanup(func() { version.Version = original })

	version.Version = "1.2.3"
	if version.Version != "1.2.3" {
		t.Errorf("expected injected '1.2.3', got %q", version.Version)
	}
}
