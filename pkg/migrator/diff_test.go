package migrator

import (
	"testing"

	"github.com/yaroher/surrealdb.go.orm/pkg/qb"
)

type testPrompter struct{ confirm bool }

func (t testPrompter) ConfirmRename(table, from, to string) bool { return t.confirm }

func TestDiffTableResourcesRenameWithCopy(t *testing.T) {
	cur := map[string]map[string]Definition{
		"user": {
			"new": {Name: "new", Statement: qb.DefineFieldName("new", "user").Type("string")},
		},
	}
	prev := map[string]map[string]Definition{
		"user": {
			"old": {Name: "old", Statement: qb.DefineFieldName("old", "user").Type("string")},
		},
	}
	opts := DiffOptions{RenameStrategy: "rename", RenameExpr: "coalesce({old}, '')"}

	got := diffTableResources(cur, prev, opts, true, true, func(table, name string) qb.Statement {
		return qb.RemoveField(name).OnTableName(table)
	})

	if len(got) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(got))
	}
	texts := statementsText(got)
	want := []string{
		"UPDATE user SET new = coalesce(old, '')",
		"REMOVE FIELD old ON TABLE user",
		"DEFINE FIELD new ON TABLE user TYPE string",
	}
	assertTexts(t, texts, want)
}

func TestDiffTableResourcesPromptRename(t *testing.T) {
	cur := map[string]map[string]Definition{
		"user": {
			"new": {Name: "new", Statement: qb.DefineFieldName("new", "user").Type("string")},
		},
	}
	prev := map[string]map[string]Definition{
		"user": {
			"old": {Name: "old", Statement: qb.DefineFieldName("old", "user").Type("string")},
		},
	}
	opts := DiffOptions{Prompter: testPrompter{confirm: true}}

	got := diffTableResources(cur, prev, opts, true, true, func(table, name string) qb.Statement {
		return qb.RemoveField(name).OnTableName(table)
	})

	if len(got) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(got))
	}
	texts := statementsText(got)
	want := []string{
		"UPDATE user SET new = old",
		"REMOVE FIELD old ON TABLE user",
		"DEFINE FIELD new ON TABLE user TYPE string",
	}
	assertTexts(t, texts, want)
}

func TestDiffResourcesAccessScope(t *testing.T) {
	code := NewResourceSet()
	code.AddAccess("acc", "database", qb.DefineAccess("acc").OnDatabase().TypeJWT().JWTAlgorithmKey("HS256", "secret"))

	db := NewResourceSet()
	up, down := DiffResources(code, db, DiffOptions{})

	if len(up) != 1 || len(down) != 1 {
		t.Fatalf("expected 1 up and 1 down statement, got %d/%d", len(up), len(down))
	}
	if got := qb.Build(down[0]).Text; got != "REMOVE ACCESS acc ON DATABASE" {
		t.Fatalf("unexpected down statement: %s", got)
	}
}

func TestPickRenameStrategyForce(t *testing.T) {
	if got := pickRenameStrategy(DiffOptions{Force: true}); got != renameStrategyDelete {
		t.Fatalf("expected force to use delete, got %s", got)
	}
	if got := pickRenameStrategy(DiffOptions{RenameStrategy: "keep"}); got != renameStrategyKeep {
		t.Fatalf("expected keep strategy, got %s", got)
	}
}

func statementsText(stmts []qb.Statement) []string {
	out := make([]string, 0, len(stmts))
	for _, stmt := range stmts {
		out = append(out, qb.Build(stmt).Text)
	}
	return out
}

func assertTexts(t *testing.T, got, want []string) {
	if len(got) != len(want) {
		t.Fatalf("expected %d statements, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected statement[%d]:\n%s\nexpected:\n%s", i, got[i], want[i])
		}
	}
}
