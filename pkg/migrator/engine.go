package migrator

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yaroher/surrealdb.go.orm/pkg/qb"
)

const migrationsTable = "_migrations"

type Migrator struct {
	DB       DB
	Config   Config
	Prompter Prompter
}

func New(db DB, cfg Config, prompter Prompter) *Migrator {
	if prompter == nil {
		prompter = NoPrompter{}
	}
	return &Migrator{DB: db, Config: cfg, Prompter: prompter}
}

// Init ensures migration metadata table exists.
func (m *Migrator) Init(ctx context.Context) error {
	if m.DB == nil {
		return errors.New("migrator: DB is required")
	}
	stmt := qb.QueryChain(
		qb.DefineTableName(migrationsTable),
		qb.DefineFieldName("id", migrationsTable).Type("string"),
		qb.DefineFieldName("checksum", migrationsTable).Type("string"),
		qb.DefineFieldName("applied_at", migrationsTable).Type("datetime"),
		qb.DefineFieldName("name", migrationsTable).Type("string"),
		qb.DefineFieldName("up", migrationsTable).Type("string"),
		qb.DefineFieldName("down", migrationsTable).Type("string"),
	)
	_, err := m.DB.Query(ctx, qb.Build(stmt).Text, qb.Build(stmt).Args)
	return err
}

// Up applies pending migrations.
func (m *Migrator) Up(ctx context.Context, steps int) error {
	if m.DB == nil {
		return errors.New("migrator: DB is required")
	}
	if err := m.Init(ctx); err != nil {
		return err
	}
	source := FileSource{Dir: m.Config.Dir}
	migs, err := source.Migrations()
	if err != nil {
		return err
	}
	applied, err := m.appliedMigrations(ctx)
	if err != nil {
		return err
	}
	var pending []Migration
	for _, mig := range migs {
		if rec, ok := applied[mig.ID]; ok {
			if m.Config.Mode == ModeStrict && rec.Checksum != mig.Checksum {
				return fmt.Errorf("checksum mismatch for %s", mig.ID)
			}
			continue
		}
		pending = append(pending, mig)
	}
	if steps > 0 && len(pending) > steps {
		pending = pending[:steps]
	}
	for _, mig := range pending {
		if err := m.applyMigration(ctx, mig); err != nil {
			return err
		}
	}
	return nil
}

// Down rolls back migrations.
func (m *Migrator) Down(ctx context.Context, steps int) error {
	if m.DB == nil {
		return errors.New("migrator: DB is required")
	}
	if err := m.Init(ctx); err != nil {
		return err
	}
	source := FileSource{Dir: m.Config.Dir}
	migs, err := source.Migrations()
	if err != nil {
		return err
	}
	appliedList, err := m.appliedMigrationsList(ctx)
	if err != nil {
		return err
	}
	if steps <= 0 {
		steps = 1
	}
	if len(appliedList) < steps {
		steps = len(appliedList)
	}
	for i := 0; i < steps; i++ {
		migID := appliedList[len(appliedList)-1-i].ID
		mig := findMigration(migs, migID)
		if mig.ID == "" {
			return fmt.Errorf("migration %s not found in dir", migID)
		}
		if strings.TrimSpace(mig.DownSQL) == "" {
			return fmt.Errorf("migration %s has no down script", mig.ID)
		}
		if err := m.rollbackMigration(ctx, mig); err != nil {
			return err
		}
	}
	return nil
}

// List returns migration status.
func (m *Migrator) List(ctx context.Context) ([]MigrationStatus, error) {
	source := FileSource{Dir: m.Config.Dir}
	migs, err := source.Migrations()
	if err != nil {
		return nil, err
	}
	applied, err := m.appliedMigrations(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]MigrationStatus, 0, len(migs))
	for _, mig := range migs {
		status := MigrationStatus{Migration: mig, Applied: false}
		if rec, ok := applied[mig.ID]; ok {
			status.Applied = true
			status.AppliedAt = rec.AppliedAt
			status.Checksum = rec.Checksum
		}
		out = append(out, status)
	}
	return out, nil
}

// Reset rolls back all applied migrations.
func (m *Migrator) Reset(ctx context.Context) error {
	applied, err := m.appliedMigrationsList(ctx)
	if err != nil {
		return err
	}
	for range applied {
		if err := m.Down(ctx, 1); err != nil {
			return err
		}
	}
	return nil
}

// Generate creates a new migration file based on resource diff.
func (m *Migrator) Generate(ctx context.Context, code ResourceSet, db ResourceSet, name string) (string, error) {
	if name == "" {
		name = "migration"
	}
	up, down := DiffResources(code, db, DiffOptions{
		Prompter:       m.Prompter,
		Force:          m.Config.Force,
		RenameStrategy: m.Config.RenameStrategy,
		RenameExpr:     m.Config.RenameExpr,
	})
	if m.Config.GrantsAlways || shouldIncludeAccessGrants(m.Config.Dir, code.AccessGrants) {
		for _, grant := range code.AccessGrants {
			up = append(up, grant.Statement)
		}
	}
	if len(up) == 0 && len(down) == 0 {
		return "", fmt.Errorf("no changes detected")
	}
	stamp := time.Now().UTC().Format("20060102150405")
	base := fmt.Sprintf("%s_%s", stamp, sanitizeName(name))
	dir := m.Config.Dir
	if dir == "" {
		dir = "migrations"
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	upText := renderStatements(up)
	downText := renderStatements(down)
	if m.Config.TwoWay {
		upPath := filepath.Join(dir, base+".up.surql")
		downPath := filepath.Join(dir, base+".down.surql")
		if err := os.WriteFile(upPath, []byte(upText), 0o644); err != nil {
			return "", err
		}
		if err := os.WriteFile(downPath, []byte(downText), 0o644); err != nil {
			return "", err
		}
		return base, nil
	}
	path := filepath.Join(dir, base+".surql")
	if err := os.WriteFile(path, []byte(upText), 0o644); err != nil {
		return "", err
	}
	return base, nil
}

func shouldIncludeAccessGrants(dir string, grants []Definition) bool {
	if len(grants) == 0 {
		return false
	}
	if dir == "" {
		dir = "migrations"
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return true
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".surql") {
			return false
		}
	}
	return true
}

func (m *Migrator) applyMigration(ctx context.Context, mig Migration) error {
	if strings.TrimSpace(mig.UpSQL) == "" {
		return fmt.Errorf("migration %s has empty up", mig.ID)
	}
	if _, err := m.DB.Query(ctx, mig.UpSQL, nil); err != nil {
		return err
	}
	insert := qb.Create(qb.T(migrationsTable)).Content(map[string]any{
		"id":         mig.ID,
		"checksum":   mig.Checksum,
		"applied_at": time.Now().UTC(),
		"name":       mig.ID,
		"up":         mig.UpSQL,
		"down":       mig.DownSQL,
	})
	_, err := m.DB.Query(ctx, qb.Build(insert).Text, qb.Build(insert).Args)
	return err
}

func (m *Migrator) rollbackMigration(ctx context.Context, mig Migration) error {
	if strings.TrimSpace(mig.DownSQL) == "" {
		return fmt.Errorf("migration %s has empty down", mig.ID)
	}
	if _, err := m.DB.Query(ctx, mig.DownSQL, nil); err != nil {
		return err
	}
	remove := qb.Delete(qb.T(migrationsTable)).Where(qb.I("id").Eq(mig.ID))
	_, err := m.DB.Query(ctx, qb.Build(remove).Text, qb.Build(remove).Args)
	return err
}

func (m *Migrator) appliedMigrations(ctx context.Context) (map[string]MigrationRecord, error) {
	rows, err := m.DB.Query(ctx, fmt.Sprintf("SELECT id, checksum, applied_at FROM %s ORDER BY id", migrationsTable), nil)
	if err != nil {
		return nil, err
	}
	out := map[string]MigrationRecord{}
	for _, row := range rows {
		id := toString(row["id"])
		if id == "" {
			continue
		}
		rec := MigrationRecord{
			ID:        id,
			Checksum:  toString(row["checksum"]),
			AppliedAt: toString(row["applied_at"]),
		}
		out[id] = rec
	}
	return out, nil
}

func (m *Migrator) appliedMigrationsList(ctx context.Context) ([]MigrationRecord, error) {
	rows, err := m.DB.Query(ctx, fmt.Sprintf("SELECT id, checksum, applied_at FROM %s ORDER BY id", migrationsTable), nil)
	if err != nil {
		return nil, err
	}
	out := make([]MigrationRecord, 0, len(rows))
	for _, row := range rows {
		id := toString(row["id"])
		if id == "" {
			continue
		}
		out = append(out, MigrationRecord{
			ID:        id,
			Checksum:  toString(row["checksum"]),
			AppliedAt: toString(row["applied_at"]),
		})
	}
	return out, nil
}

func findMigration(migs []Migration, id string) Migration {
	for _, m := range migs {
		if m.ID == id {
			return m
		}
	}
	return Migration{}
}

func renderStatements(stmts []qb.Statement) string {
	parts := make([]string, 0, len(stmts))
	for _, stmt := range stmts {
		text := qb.Build(stmt).Text
		if text != "" {
			parts = append(parts, text+";")
		}
	}
	return strings.Join(parts, "\n") + "\n"
}

func sanitizeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	return name
}

func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		return fmt.Sprintf("%v", v)
	}
}
