package surreal

import "context"

// Client is a minimal interface to execute SurrealQL.
type Client interface {
	Query(ctx context.Context, sql string, vars map[string]any) (any, error)
}
