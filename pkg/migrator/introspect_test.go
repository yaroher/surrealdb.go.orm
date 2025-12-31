package migrator

import (
	"context"
	"strings"
	"testing"
)

func TestIntrospect(t *testing.T) {
	db := &fakeDB{}
	db.fn = func(sql string) ([]map[string]any, error) {
		switch {
		case sql == "INFO FOR DB":
			return []map[string]any{{
				"tables": map[any]any{
					"user": map[string]any{"definition": "DEFINE TABLE user"},
				},
				"accesses": map[string]any{
					"acc": map[string]any{"def": "DEFINE ACCESS acc ON DATABASE"},
				},
			}}, nil
		case strings.HasPrefix(sql, "INFO FOR TABLE user"):
			return []map[string]any{{
				"fields": map[string]any{
					"name": "DEFINE FIELD name ON TABLE user",
				},
				"indexes": map[string]any{
					"idx": map[string]any{"definition": "DEFINE INDEX idx ON TABLE user"},
				},
				"events": map[any]any{
					"ev": map[any]any{"def": "DEFINE EVENT ev ON TABLE user"},
				},
			}}, nil
		default:
			return nil, nil
		}
	}

	res, err := Introspect(context.Background(), db)
	if err != nil {
		t.Fatalf("introspect: %v", err)
	}
	if _, ok := res.Tables["user"]; !ok {
		t.Fatalf("expected user table")
	}
	if _, ok := res.Fields["user"]["name"]; !ok {
		t.Fatalf("expected user.name field")
	}
	if _, ok := res.Indexes["user"]["idx"]; !ok {
		t.Fatalf("expected user.idx index")
	}
	if _, ok := res.Events["user"]["ev"]; !ok {
		t.Fatalf("expected user.ev event")
	}
	if _, ok := res.Access["acc"]; !ok {
		t.Fatalf("expected access acc")
	}
}

func TestExtractDefinitionVariants(t *testing.T) {
	if got := extractDefinition("DEF"); got != "DEF" {
		t.Fatalf("unexpected string def")
	}
	if got := extractDefinition(map[string]any{"definition": "A"}); got != "A" {
		t.Fatalf("unexpected definition: %s", got)
	}
	if got := extractDefinition(map[string]any{"def": "B"}); got != "B" {
		t.Fatalf("unexpected def: %s", got)
	}
	if got := extractDefinition(map[any]any{"definition": "C"}); got != "C" {
		t.Fatalf("unexpected def: %s", got)
	}
}
