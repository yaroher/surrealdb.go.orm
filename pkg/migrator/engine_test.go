package migrator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestShouldIncludeAccessGrants(t *testing.T) {
	grants := []Definition{{Name: "acc"}}

	if shouldIncludeAccessGrants("missing-dir", nil) {
		t.Fatalf("expected no grants to return false")
	}

	missing := filepath.Join(t.TempDir(), "missing")
	if !shouldIncludeAccessGrants(missing, grants) {
		t.Fatalf("expected missing dir to return true")
	}

	emptyDir := t.TempDir()
	if !shouldIncludeAccessGrants(emptyDir, grants) {
		t.Fatalf("expected empty dir to return true")
	}

	if err := os.WriteFile(filepath.Join(emptyDir, "001_init.surql"), []byte(""), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if shouldIncludeAccessGrants(emptyDir, grants) {
		t.Fatalf("expected existing .surql to return false")
	}
}
