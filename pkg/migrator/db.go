package migrator

import "context"

// DB abstracts query execution for migrations.
type DB interface {
	Query(ctx context.Context, sql string, vars map[string]any) ([]map[string]any, error)
}
