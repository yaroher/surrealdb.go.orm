package qb

import "testing"

func TestDefineAccessJWT(t *testing.T) {
	stmt := DefineAccess("app").
		Overwrite().
		OnDatabase().
		TypeJWT().
		JWTAlgorithmKey("HS256", "secret").
		Authenticate(Raw("auth::check()"))
	stmt.DurationGrant(Raw("1h")).DurationToken(Raw("2h")).DurationSession(Raw("3h"))

	q := Build(stmt)
	expected := "DEFINE ACCESS OVERWRITE app ON DATABASE TYPE JWT ALGORITHM HS256 KEY $p1 AUTHENTICATE auth::check() DURATION FOR GRANT 1h, FOR TOKEN 2h, FOR SESSION 3h"
	if q.Text != expected {
		t.Fatalf("unexpected query:\n%s\nexpected:\n%s", q.Text, expected)
	}
	if _, ok := q.Args["p1"]; !ok {
		t.Fatalf("expected bound param p1")
	}
}

func TestDefineAccessRecord(t *testing.T) {
	stmt := DefineAccess("acc").
		IfNotExists().
		OnNamespace().
		TypeRecord().
		Signup(Raw("create user")).
		Signin(Raw("select user")).
		RecordJWTURL("https://example.com/jwks").
		RecordIssuerKey("issuer").
		WithRefresh()

	q := Build(stmt)
	expected := "DEFINE ACCESS IF NOT EXISTS acc ON NAMESPACE TYPE RECORD SIGNUP { create user } SIGNIN { select user } WITH JWT URL $p1 WITH ISSUER KEY $p2 WITH REFRESH"
	if q.Text != expected {
		t.Fatalf("unexpected query:\n%s\nexpected:\n%s", q.Text, expected)
	}
	if _, ok := q.Args["p1"]; !ok {
		t.Fatalf("expected bound param p1")
	}
	if _, ok := q.Args["p2"]; !ok {
		t.Fatalf("expected bound param p2")
	}
}

func TestRemoveAccess(t *testing.T) {
	stmt := RemoveAccess("acc").IfExists().OnDatabase()
	q := Build(stmt)
	expected := "REMOVE ACCESS IF EXISTS acc ON DATABASE"
	if q.Text != expected {
		t.Fatalf("unexpected query:\n%s\nexpected:\n%s", q.Text, expected)
	}
}
