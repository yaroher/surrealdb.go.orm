package codegen

import (
	"testing"

	"github.com/yaroher/surrealdb.go.orm/pkg/qb"
)

func TestBuildResourceSetAccessGrant(t *testing.T) {
	model := Model{
		Kind: "access",
		Access: AccessConfig{
			Name:            "acc",
			Scope:           "database",
			Type:            "jwt",
			Algorithm:       "HS256",
			Key:             "secret",
			Authenticate:    "auth::check()",
			DurationGrant:   "1h",
			DurationToken:   "2h",
			DurationSession: "3h",
			GrantSubject:    "user",
			GrantToken:      "token",
			GrantDuration:   "1d",
		},
	}

	res := BuildResourceSet([]Model{model})
	def, ok := res.Access["acc"]
	if !ok {
		t.Fatalf("expected access acc in resources")
	}
	if def.Scope != "database" {
		t.Fatalf("expected scope database, got %s", def.Scope)
	}

	accessQuery := qb.Build(def.Statement)
	expectedAccess := "DEFINE ACCESS acc ON DATABASE TYPE JWT ALGORITHM HS256 KEY $p1 AUTHENTICATE auth::check() DURATION FOR GRANT 1h, FOR TOKEN 2h, FOR SESSION 3h"
	if accessQuery.Text != expectedAccess {
		t.Fatalf("unexpected access query:\n%s\nexpected:\n%s", accessQuery.Text, expectedAccess)
	}
	if _, ok := accessQuery.Args["p1"]; !ok {
		t.Fatalf("expected access key to be bound")
	}

	if len(res.AccessGrants) != 1 {
		t.Fatalf("expected 1 access grant, got %d", len(res.AccessGrants))
	}
	grantText := qb.Build(res.AccessGrants[0].Statement).Text
	expectedGrant := "ACCESS acc ON DATABASE GRANT FOR user TOKEN token DURATION 1d"
	if grantText != expectedGrant {
		t.Fatalf("unexpected grant query:\n%s\nexpected:\n%s", grantText, expectedGrant)
	}
}

func TestBuildResourceSetPermissions(t *testing.T) {
	model := Model{
		Kind:        "node",
		Table:       "user",
		Permissions: "full",
		Fields: []Field{
			{Name: "Name", Type: "string", DBName: "name", Permissions: "none"},
			{Name: "Role", Type: "string", DBName: "role", Permissions: "for select where true"},
		},
	}

	res := BuildResourceSet([]Model{model})
	table := res.Tables["user"]
	if q := qb.Build(table.Statement).Text; q != "DEFINE TABLE user PERMISSIONS FULL" {
		t.Fatalf("unexpected table permissions: %s", q)
	}
	nameField := res.Fields["user"]["name"]
	if q := qb.Build(nameField.Statement).Text; q != "DEFINE FIELD name ON TABLE user TYPE string PERMISSIONS NONE" {
		t.Fatalf("unexpected field permissions: %s", q)
	}
	roleField := res.Fields["user"]["role"]
	if q := qb.Build(roleField.Statement).Text; q != "DEFINE FIELD role ON TABLE user TYPE string PERMISSIONS\nfor select where true" {
		t.Fatalf("unexpected field permissions: %s", q)
	}
}
