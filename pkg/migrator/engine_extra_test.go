package migrator

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yaroher/surrealdb.go.orm/pkg/qb"
)

func writeMigration(t *testing.T, dir, id, up, down string, twoWay bool) {
	t.Helper()
	if twoWay {
		if err := os.WriteFile(filepath.Join(dir, id+".up.surql"), []byte(up), 0o644); err != nil {
			t.Fatalf("write up: %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, id+".down.surql"), []byte(down), 0o644); err != nil {
			t.Fatalf("write down: %v", err)
		}
		return
	}
	if err := os.WriteFile(filepath.Join(dir, id+".surql"), []byte(up), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestMigratorUpDownListReset(t *testing.T) {
	dir := t.TempDir()
	writeMigration(t, dir, "001_init", "DEFINE TABLE test", "REMOVE TABLE test", true)

	db := &fakeDB{}
	db.fn = func(sql string) ([]map[string]any, error) {
		if strings.HasPrefix(sql, "SELECT id, checksum") {
			return nil, nil
		}
		return nil, nil
	}

	m := New(db, Config{Dir: dir, Mode: ModeLax, TwoWay: true}, NoPrompter{})
	if err := m.Up(context.Background(), 0); err != nil {
		t.Fatalf("up: %v", err)
	}

	// simulate applied for down/list/reset
	db.fn = func(sql string) ([]map[string]any, error) {
		if strings.HasPrefix(sql, "SELECT id, checksum") {
			return []map[string]any{{"id": "001_init", "checksum": "x", "applied_at": "now"}}, nil
		}
		return nil, nil
	}
	if err := m.Down(context.Background(), 1); err != nil {
		t.Fatalf("down: %v", err)
	}
	if _, err := m.List(context.Background()); err != nil {
		t.Fatalf("list: %v", err)
	}
	if err := m.Reset(context.Background()); err != nil {
		t.Fatalf("reset: %v", err)
	}
}

func TestMigratorGenerate(t *testing.T) {
	dir := t.TempDir()
	m := New(&fakeDB{}, Config{Dir: dir, TwoWay: false}, NoPrompter{})
	code := NewResourceSet()
	code.AddTable("user", qb.DefineTableName("user"))
	db := NewResourceSet()

	id, err := m.Generate(context.Background(), code, db, "My Migration")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if id == "" {
		t.Fatalf("expected migration id")
	}
	if _, err := os.Stat(filepath.Join(dir, id+".surql")); err != nil {
		t.Fatalf("expected migration file")
	}
}

func TestGenerateNoChanges(t *testing.T) {
	m := New(&fakeDB{}, Config{Dir: t.TempDir()}, NoPrompter{})
	code := NewResourceSet()
	db := NewResourceSet()
	if _, err := m.Generate(context.Background(), code, db, ""); err == nil {
		t.Fatalf("expected no changes error")
	}
}

func TestChecksumMismatchStrict(t *testing.T) {
	dir := t.TempDir()
	writeMigration(t, dir, "001_init", "DEFINE TABLE test", "REMOVE TABLE test", true)
	migs, _ := FileSource{Dir: dir}.Migrations()
	if len(migs) == 0 {
		t.Fatalf("expected migrations")
	}
	badChecksum := migs[0].Checksum + "x"

	db := &fakeDB{}
	db.fn = func(sql string) ([]map[string]any, error) {
		if strings.HasPrefix(sql, "SELECT id, checksum") {
			return []map[string]any{{"id": "001_init", "checksum": badChecksum, "applied_at": "now"}}, nil
		}
		return nil, nil
	}

	m := New(db, Config{Dir: dir, Mode: ModeStrict}, NoPrompter{})
	if err := m.Up(context.Background(), 0); err == nil {
		t.Fatalf("expected checksum mismatch error")
	}
}

func TestDownErrors(t *testing.T) {
	dir := t.TempDir()
	writeMigration(t, dir, "001_init", "DEFINE TABLE test", "", true)

	db := &fakeDB{}
	db.fn = func(sql string) ([]map[string]any, error) {
		if strings.HasPrefix(sql, "SELECT id, checksum") {
			return []map[string]any{{"id": "001_init", "checksum": "x", "applied_at": "now"}}, nil
		}
		return nil, nil
	}

	m := New(db, Config{Dir: dir}, NoPrompter{})
	if err := m.Down(context.Background(), 1); err == nil {
		t.Fatalf("expected error for empty down")
	}
}

func TestHelperFunctions(t *testing.T) {
	if got := sanitizeName("My Name-Here"); got != "my_name_here" {
		t.Fatalf("unexpected sanitize: %s", got)
	}
	if got := toString([]byte("x")); got != "x" {
		t.Fatalf("unexpected toString: %s", got)
	}

	stmts := []qb.Statement{qb.DefineNamespace("ns"), qb.DefineDatabase("db")}
	text := renderStatements(stmts)
	if !strings.Contains(text, "DEFINE NAMESPACE ns") || !strings.Contains(text, "DEFINE DATABASE db") {
		t.Fatalf("unexpected renderStatements: %s", text)
	}

	migs := []Migration{{ID: "a"}, {ID: "b"}}
	if got := findMigration(migs, "b"); got.ID != "b" {
		t.Fatalf("unexpected findMigration")
	}
}

func TestNewMigratorNilDB(t *testing.T) {
	m := New(nil, Config{}, NoPrompter{})
	if err := m.Init(context.Background()); err == nil {
		t.Fatalf("expected error for nil DB")
	}
	if err := m.Up(context.Background(), 0); err == nil {
		t.Fatalf("expected error for nil DB")
	}
	if err := m.Down(context.Background(), 1); err == nil {
		t.Fatalf("expected error for nil DB")
	}
}
