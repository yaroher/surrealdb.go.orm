package migrator

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestPrune(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "001_init.surql"), []byte("init"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "002_add.surql"), []byte("add"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	db := &fakeDB{}
	db.fn = func(sql string) ([]map[string]any, error) {
		return []map[string]any{{"id": "001_init", "checksum": "x", "applied_at": "now"}}, nil
	}

	m := New(db, Config{Dir: dir}, NoPrompter{})
	removed, err := m.Prune(context.Background())
	if err != nil {
		t.Fatalf("prune: %v", err)
	}
	if removed != 1 {
		t.Fatalf("expected 1 removed, got %d", removed)
	}
	if _, err := os.Stat(filepath.Join(dir, "002_add.surql")); err == nil {
		t.Fatalf("expected 002_add.surql to be removed")
	}
}
