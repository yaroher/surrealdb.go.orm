package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/yaroher/surrealdb.go.orm/pkg/migrator"
)

type stubDB struct {
	fn func(sql string) ([]map[string]any, error)
}

func (s stubDB) Query(ctx context.Context, sql string, vars map[string]any) ([]map[string]any, error) {
	if s.fn == nil {
		return nil, nil
	}
	return s.fn(sql)
}

func withConnectStub(t *testing.T, db migrator.DB) {
	orig := connectFn
	connectFn = func(ctx context.Context, cfg migrator.Config) (migrator.DB, func(), error) {
		return db, func() {}, nil
	}
	t.Cleanup(func() { connectFn = orig })
}

func setCmdFlags(t *testing.T, cmd *cobra.Command, dir string) {
	t.Helper()
	if cmd.Flags().Lookup("dsn") == nil {
		cmd.Flags().AddFlagSet(migrateCmd.PersistentFlags())
	}
	setFlag := func(name, value string) {
		_ = cmd.Flags().Set(name, value)
	}
	setFlag("dir", dir)
	setFlag("mode", "lax")
	setFlag("dsn", "http://example")
	setFlag("ns", "test")
	setFlag("db", "db")
	setFlag("username", "user")
	setFlag("password", "pass")
	setFlag("force", "true")
	setFlag("two-way", "true")
	setFlag("rename-strategy", "rename")
	setFlag("rename-expr", "{old}")
	setFlag("grants-always", "true")
}

func TestBuildMigratorRequiresDSN(t *testing.T) {
	if migrateCmd.Flags().Lookup("dsn") == nil {
		migrateCmd.Flags().AddFlagSet(migrateCmd.PersistentFlags())
	}
	_ = migrateCmd.PersistentFlags().Set("dsn", "")
	_ = migrateCmd.Flags().Set("dsn", "")
	_, _, _, err := buildMigrator(migrateCmd)
	if err == nil {
		t.Fatalf("expected error for missing dsn")
	}
}

func TestBuildMigratorSuccess(t *testing.T) {
	withConnectStub(t, stubDB{})
	setCmdFlags(t, migrateCmd, t.TempDir())

	m, ctx, cleanup, err := buildMigrator(migrateCmd)
	if err != nil {
		t.Fatalf("buildMigrator: %v", err)
	}
	cleanup()
	if ctx == nil || m == nil {
		t.Fatalf("expected migrator and ctx")
	}
	if _, ok := m.Prompter.(migrator.NoPrompter); !ok {
		t.Fatalf("expected NoPrompter when force=true")
	}
}

func TestGenerateCommand(t *testing.T) {
	dir := t.TempDir()
	src := `package sample

// orm:node table=users
type User struct{ Name string }
`
	if err := os.WriteFile(filepath.Join(dir, "model.go"), []byte(src), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	_ = generateCmd.Flags().Set("dir", dir)
	if err := generateCmd.RunE(generateCmd, nil); err != nil {
		t.Fatalf("generate: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "orm_gen.go")); err != nil {
		t.Fatalf("expected orm_gen.go")
	}
}

func TestMigrateGenerateCommand(t *testing.T) {
	dir := t.TempDir()
	codeDir := t.TempDir()
	src := `package sample

// orm:node table=users
type User struct{ Name string }
`
	if err := os.WriteFile(filepath.Join(codeDir, "model.go"), []byte(src), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	withConnectStub(t, stubDB{fn: func(sql string) ([]map[string]any, error) {
		if strings.HasPrefix(sql, "INFO FOR DB") {
			return []map[string]any{{}}, nil
		}
		return nil, nil
	}})
	setCmdFlags(t, migrateGenerateCmd, dir)
	_ = migrateGenerateCmd.Flags().Set("code", codeDir)
	_ = migrateGenerateCmd.Flags().Set("name", "init")
	if err := migrateGenerateCmd.RunE(migrateGenerateCmd, nil); err != nil {
		t.Fatalf("migrate generate: %v", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("readdir: %v", err)
	}
	if len(entries) == 0 {
		t.Fatalf("expected migration files")
	}
}

func TestMigrateInitListPruneCommands(t *testing.T) {
	dir := t.TempDir()
	withConnectStub(t, stubDB{fn: func(sql string) ([]map[string]any, error) {
		if strings.HasPrefix(sql, "SELECT id, checksum") {
			return []map[string]any{{"id": "001", "checksum": "x", "applied_at": "now"}}, nil
		}
		return nil, nil
	}})
	setCmdFlags(t, migrateInitCmd, dir)
	setCmdFlags(t, migrateListCmd, dir)
	setCmdFlags(t, migratePruneCmd, dir)

	if err := migrateInitCmd.RunE(migrateInitCmd, nil); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := migrateListCmd.RunE(migrateListCmd, nil); err != nil {
		t.Fatalf("list: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "001.surql"), []byte(""), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := migratePruneCmd.RunE(migratePruneCmd, nil); err != nil {
		t.Fatalf("prune: %v", err)
	}
}

func TestInitEnv(t *testing.T) {
	_ = os.Setenv("SURREAL_ORM_DEBUG", "1")
	initEnv()
	_ = os.Unsetenv("SURREAL_ORM_DEBUG")
}

func TestMigrateUpDownResetCommands(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "001_init.up.surql"), []byte("DEFINE TABLE test"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "001_init.down.surql"), []byte("REMOVE TABLE test"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	withConnectStub(t, stubDB{fn: func(sql string) ([]map[string]any, error) {
		if strings.HasPrefix(sql, "SELECT id, checksum") {
			return []map[string]any{{"id": "001_init", "checksum": "x", "applied_at": "now"}}, nil
		}
		return nil, nil
	}})
	setCmdFlags(t, migrateUpCmd, dir)
	setCmdFlags(t, migrateDownCmd, dir)
	setCmdFlags(t, migrateResetCmd, dir)
	_ = migrateUpCmd.Flags().Set("steps", "0")
	_ = migrateDownCmd.Flags().Set("steps", "1")

	if err := migrateUpCmd.RunE(migrateUpCmd, nil); err != nil {
		t.Fatalf("up: %v", err)
	}
	if err := migrateDownCmd.RunE(migrateDownCmd, nil); err != nil {
		t.Fatalf("down: %v", err)
	}
	if err := migrateResetCmd.RunE(migrateResetCmd, nil); err != nil {
		t.Fatalf("reset: %v", err)
	}
}
