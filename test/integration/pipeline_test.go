package integration

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yaroher/surrealdb.go.orm/internal/codegen"
	"github.com/yaroher/surrealdb.go.orm/pkg/migrator"
)

type stubDB struct {
	calls []string
	fn    func(sql string) ([]map[string]any, error)
}

func (s *stubDB) Query(ctx context.Context, sql string, vars map[string]any) ([]map[string]any, error) {
	s.calls = append(s.calls, sql)
	if s.fn != nil {
		return s.fn(sql)
	}
	return nil, nil
}

func TestIntegrationGenerateAndApply(t *testing.T) {
	codeDir := t.TempDir()
	src := `package sample

// orm:node table=users
type User struct { Name string }
`
	if err := os.WriteFile(filepath.Join(codeDir, "model.go"), []byte(src), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	pkg, err := codegen.ParseDir(codeDir)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	code := codegen.BuildResourceSet(pkg.Models)

	introspectDB := &stubDB{fn: func(sql string) ([]map[string]any, error) {
		if sql == "INFO FOR DB" {
			return []map[string]any{{"tables": map[string]any{}, "accesses": map[string]any{}}}, nil
		}
		return nil, nil
	}}
	dbResources, err := migrator.Introspect(context.Background(), introspectDB)
	if err != nil {
		t.Fatalf("introspect: %v", err)
	}

	migDir := t.TempDir()
	applyDB := &stubDB{fn: func(sql string) ([]map[string]any, error) {
		if strings.HasPrefix(sql, "SELECT id, checksum") {
			return nil, nil
		}
		return nil, nil
	}}
	m := migrator.New(applyDB, migrator.Config{Dir: migDir, TwoWay: true, Mode: migrator.ModeLax}, migrator.NoPrompter{})
	if _, err := m.Generate(context.Background(), code, dbResources, "init"); err != nil {
		t.Fatalf("generate: %v", err)
	}

	files, err := os.ReadDir(migDir)
	if err != nil {
		t.Fatalf("readdir: %v", err)
	}
	if len(files) == 0 {
		t.Fatalf("expected migration files")
	}

	if err := m.Up(context.Background(), 0); err != nil {
		t.Fatalf("up: %v", err)
	}

	found := false
	for _, call := range applyDB.calls {
		if strings.Contains(call, "DEFINE TABLE users") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected define table users in DB calls")
	}
}
