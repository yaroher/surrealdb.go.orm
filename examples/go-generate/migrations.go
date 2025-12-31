package gogenerate

import (
	"context"
	"strings"

	"github.com/yaroher/surrealdb.go.orm/pkg/migrator"
)

func AutoMigrate(ctx context.Context, db migrator.DB, dir string) (string, error) {
	m := migrator.New(db, migrator.Config{Dir: dir, TwoWay: true}, nil)
	if err := m.Init(ctx); err != nil {
		return "", err
	}
	current, err := migrator.Introspect(ctx, db)
	if err != nil {
		return "", err
	}
	id, err := m.Generate(ctx, Resources(), current, "auto")
	if err != nil {
		if isNoChanges(err) {
			return "", nil
		}
		return "", err
	}
	if err := m.Up(ctx, 0); err != nil {
		return id, err
	}
	return id, nil
}

func isNoChanges(err error) bool {
	return err != nil && strings.Contains(err.Error(), "no changes detected")
}
