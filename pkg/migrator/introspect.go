package migrator

import (
	"context"
	"fmt"

	"github.com/yaroher/surrealdb.go.orm/pkg/qb"
)

// Introspect reads schema resources from the database using INFO statements.
func Introspect(ctx context.Context, db DB) (ResourceSet, error) {
	res := NewResourceSet()
	rows, err := db.Query(ctx, "INFO FOR DB", nil)
	if err != nil {
		return res, err
	}
	if len(rows) == 0 {
		return res, nil
	}
	root := rows[0]
	tablesMap := toStringMap(root["tables"])
	for table, val := range tablesMap {
		def := extractDefinition(val)
		if def == "" {
			continue
		}
		res.AddTable(table, qb.RawStmt(def, nil))

		if err := introspectTable(ctx, db, table, &res); err != nil {
			return res, err
		}
	}
	accessMap := toStringMap(root["accesses"])
	for name, val := range accessMap {
		def := extractDefinition(val)
		if def == "" {
			continue
		}
		res.AddAccess(name, "database", qb.RawStmt(def, nil))
	}
	return res, nil
}

func introspectTable(ctx context.Context, db DB, table string, res *ResourceSet) error {
	rows, err := db.Query(ctx, fmt.Sprintf("INFO FOR TABLE %s", table), nil)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}
	root := rows[0]
	fieldsMap := toStringMap(root["fields"])
	for name, val := range fieldsMap {
		def := extractDefinition(val)
		if def == "" {
			continue
		}
		res.AddField(table, name, qb.RawStmt(def, nil))
	}
	indexesMap := toStringMap(root["indexes"])
	for name, val := range indexesMap {
		def := extractDefinition(val)
		if def == "" {
			continue
		}
		res.AddIndex(table, name, qb.RawStmt(def, nil))
	}
	eventsMap := toStringMap(root["events"])
	for name, val := range eventsMap {
		def := extractDefinition(val)
		if def == "" {
			continue
		}
		res.AddEvent(table, name, qb.RawStmt(def, nil))
	}
	return nil
}

func toStringMap(val any) map[string]any {
	switch t := val.(type) {
	case map[string]any:
		return t
	case map[any]any:
		out := map[string]any{}
		for k, v := range t {
			ks := fmt.Sprintf("%v", k)
			out[ks] = v
		}
		return out
	default:
		return map[string]any{}
	}
}

func extractDefinition(val any) string {
	switch t := val.(type) {
	case string:
		return t
	case map[string]any:
		if def, ok := t["definition"].(string); ok {
			return def
		}
		if def, ok := t["def"].(string); ok {
			return def
		}
	case map[any]any:
		return extractDefinition(toStringMap(t))
	}
	return ""
}
