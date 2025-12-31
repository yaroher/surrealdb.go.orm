package surreal

import (
	"context"

	surrealdb "github.com/surrealdb/surrealdb.go"
)

var queryFn = surrealdb.Query[[]map[string]any]
var fromEndpointFn = surrealdb.FromEndpointURLString
var signInFn = func(ctx context.Context, db *surrealdb.DB, auth surrealdb.Auth) (any, error) {
	return db.SignIn(ctx, auth)
}
var useFn = func(ctx context.Context, db *surrealdb.DB, ns, dbName string) error {
	return db.Use(ctx, ns, dbName)
}

// Adapter wraps surrealdb.DB to satisfy migrator.DB.
type Adapter struct {
	DB *surrealdb.DB
}

func (a Adapter) Query(ctx context.Context, sql string, vars map[string]any) ([]map[string]any, error) {
	res, err := queryFn(ctx, a.DB, sql, vars)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	out := make([]map[string]any, 0)
	for _, r := range *res {
		out = append(out, r.Result...)
	}
	return out, nil
}

// Connect establishes a SurrealDB connection and signs in if credentials provided.
func Connect(ctx context.Context, dsn, ns, db, user, pass string) (*surrealdb.DB, error) {
	client, err := fromEndpointFn(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if user != "" || pass != "" {
		_, err = signInFn(ctx, client, surrealdb.Auth{
			Namespace: ns,
			Database:  db,
			Username:  user,
			Password:  pass,
		})
		if err != nil {
			return nil, err
		}
	}
	if ns != "" || db != "" {
		if err := useFn(ctx, client, ns, db); err != nil {
			return nil, err
		}
	}
	return client, nil
}
