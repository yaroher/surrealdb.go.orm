package qb

import "testing"

func assertQuery(t *testing.T, stmt Statement, expected string) Query {
	q := Build(stmt)
	if q.Text != expected {
		t.Fatalf("unexpected query:\n%s\nexpected:\n%s", q.Text, expected)
	}
	return q
}

func assertArgsLen(t *testing.T, q Query, n int) {
	if len(q.Args) != n {
		t.Fatalf("expected %d args, got %d", n, len(q.Args))
	}
}
