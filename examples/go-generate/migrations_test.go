package gogenerate

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yaroher/surrealdb.go.orm/pkg/migrator"
)

type mockMigratorDB struct {
	queries []string
}

func (m *mockMigratorDB) Query(ctx context.Context, sql string, vars map[string]any) ([]map[string]any, error) {
	m.queries = append(m.queries, sql)
	if strings.HasPrefix(sql, "INFO FOR DB") {
		return []map[string]any{}, nil
	}
	if strings.HasPrefix(sql, "SELECT id, checksum, applied_at FROM _migrations") {
		return []map[string]any{}, nil
	}
	return []map[string]any{}, nil
}

var _ migrator.DB = (*mockMigratorDB)(nil)

func TestAutoMigrate(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	db := &mockMigratorDB{}

	id, err := AutoMigrate(ctx, db, dir)
	if err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	if id == "" {
		t.Fatalf("expected migration id")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	if len(entries) == 0 {
		t.Fatalf("expected migration files")
	}

	upFound := false
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".up.surql") {
			upFound = true
			break
		}
	}
	if !upFound {
		t.Fatalf("expected .up.surql migration")
	}

	if len(db.queries) == 0 {
		t.Fatalf("expected queries to be executed")
	}

	if _, err := os.Stat(filepath.Join(dir, id+".up.surql")); err != nil {
		t.Fatalf("expected up migration: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, id+".down.surql")); err != nil {
		t.Fatalf("expected down migration: %v", err)
	}
}
