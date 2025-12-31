package migrator

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

// Prune removes unapplied local migration files.
func (m *Migrator) Prune(ctx context.Context) (int, error) {
	applied, err := m.appliedMigrations(ctx)
	if err != nil {
		return 0, err
	}
	dir := m.Config.Dir
	if dir == "" {
		dir = "migrations"
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	removed := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".surql") {
			continue
		}
		id := strings.TrimSuffix(name, ".surql")
		id = strings.TrimSuffix(id, ".up")
		id = strings.TrimSuffix(id, ".down")
		if _, ok := applied[id]; ok {
			continue
		}
		if err := os.Remove(filepath.Join(dir, name)); err != nil {
			return removed, err
		}
		removed++
	}
	return removed, nil
}
