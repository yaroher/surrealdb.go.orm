package migrator

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestFileSourceMigrations(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "20240101010101_init.surql"), []byte("DEFINE TABLE test"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "20240101010102_add.up.surql"), []byte("DEFINE FIELD name ON TABLE test"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "20240101010102_add.down.surql"), []byte("REMOVE FIELD name ON TABLE test"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	src := FileSource{Dir: dir}
	migs, err := src.Migrations()
	if err != nil {
		t.Fatalf("migrations: %v", err)
	}
	if len(migs) != 2 {
		t.Fatalf("expected 2 migrations, got %d", len(migs))
	}
	if migs[0].ID != "20240101010101_init" {
		t.Fatalf("unexpected first id: %s", migs[0].ID)
	}
	if migs[1].DownSQL == "" {
		t.Fatalf("expected down SQL for second migration")
	}
	if migs[1].Checksum == "" {
		t.Fatalf("expected checksum")
	}
}

func TestEmbeddedSourceMigrations(t *testing.T) {
	fsys := fstest.MapFS{
		"m/001_init.surql":     &fstest.MapFile{Data: []byte("DEFINE TABLE test")},
		"m/002_add.up.surql":   &fstest.MapFile{Data: []byte("DEFINE FIELD name ON TABLE test")},
		"m/002_add.down.surql": &fstest.MapFile{Data: []byte("REMOVE FIELD name ON TABLE test")},
	}
	var _ fs.FS = fsys
	src := EmbeddedSource{FS: fsys, Dir: "m"}
	migs, err := src.Migrations()
	if err != nil {
		t.Fatalf("migrations: %v", err)
	}
	if len(migs) != 2 {
		t.Fatalf("expected 2 migrations, got %d", len(migs))
	}
}
