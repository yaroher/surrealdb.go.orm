package qb

import "testing"

func TestBuilderArgBind(t *testing.T) {
	b := NewBuilder()
	b.Write("X")
	if got := b.String(); got != "X" {
		t.Fatalf("expected buffer to be X, got %q", got)
	}
	if p := b.Arg(1); p != "$p1" {
		t.Fatalf("expected $p1, got %s", p)
	}
	if p := b.Arg("a"); p != "$p2" {
		t.Fatalf("expected $p2, got %s", p)
	}
	b.Bind("p2", "override")
	b.Bind("p3", 3)
	args := b.Args()
	if args["p1"] != 1 {
		t.Fatalf("expected p1=1")
	}
	if args["p2"] != "a" {
		t.Fatalf("expected p2 to remain 'a'")
	}
	if args["p3"] != 3 {
		t.Fatalf("expected p3=3")
	}
}

func TestQueryChainAndRaw(t *testing.T) {
	stmt := QueryChain(DefineNamespace("test"), DefineDatabase("db"))
	assertQuery(t, stmt, "DEFINE NAMESPACE test; DEFINE DATABASE db")

	raw := RawStmt("LET $x = $p1", map[string]any{"p1": 42})
	q := Build(raw)
	if q.Text != "LET $x = $p1" {
		t.Fatalf("unexpected raw text: %s", q.Text)
	}
	if q.Args["p1"] != 42 {
		t.Fatalf("expected raw arg p1")
	}
}

func TestSubqueryAndStatementExpr(t *testing.T) {
	sub := Subquery{Stmt: Select(I("id")).From(T("user"))}
	assertQuery(t, Return(sub), "RETURN (SELECT id FROM user)")

	expr := StatementExpr(Select(I("id")).From(T("user")))
	assertQuery(t, Return(expr), "RETURN SELECT id FROM user")
}
