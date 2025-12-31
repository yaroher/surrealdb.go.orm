package surreal

import (
	"context"
	"testing"

	surrealdb "github.com/surrealdb/surrealdb.go"
)

func TestAdapterQuery(t *testing.T) {
	orig := queryFn
	defer func() { queryFn = orig }()

	queryFn = func(ctx context.Context, db *surrealdb.DB, sql string, vars map[string]any) (*[]surrealdb.QueryResult[[]map[string]any], error) {
		res := []surrealdb.QueryResult[[]map[string]any]{
			{Result: []map[string]any{{"a": 1}}},
			{Result: []map[string]any{{"b": 2}, {"c": 3}}},
		}
		return &res, nil
	}

	out, err := (Adapter{}).Query(context.Background(), "select", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(out))
	}
}

func TestAdapterQueryNil(t *testing.T) {
	orig := queryFn
	defer func() { queryFn = orig }()

	queryFn = func(ctx context.Context, db *surrealdb.DB, sql string, vars map[string]any) (*[]surrealdb.QueryResult[[]map[string]any], error) {
		return nil, nil
	}

	out, err := (Adapter{}).Query(context.Background(), "select", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Fatalf("expected nil output")
	}
}
