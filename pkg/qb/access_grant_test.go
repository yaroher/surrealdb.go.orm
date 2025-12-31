package qb

import "testing"

func TestAccessGrantBuild(t *testing.T) {
	stmt := AccessGrant("acc").
		OnDatabase().
		Subject(Raw("user")).
		Token(PWith("tok", "secret")).
		Duration(Raw("1d"))

	q := Build(stmt)
	expected := "ACCESS acc ON DATABASE GRANT FOR user TOKEN $tok DURATION 1d"
	if q.Text != expected {
		t.Fatalf("unexpected query:\n%s\nexpected:\n%s", q.Text, expected)
	}
	if got := q.Args["tok"]; got != "secret" {
		t.Fatalf("expected bound param tok=secret, got %v", got)
	}
}
