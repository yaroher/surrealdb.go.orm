package codegen

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseDirModels(t *testing.T) {
	dir := t.TempDir()
	src := `package sample

import "time"

// orm:node table=users rename_all=snake_case permissions=full
type User struct {
	// orm:field name=first_name permissions=none type=uuid default=now()
	FirstName string
	LastName string
	Account LinkOne[Account]
	CreatedAt time.Time
}

// orm:edge table=user_account in=User out=Account
type UserAccount struct {
	Note string
}

// orm:access name=acc on=database type=jwt alg=HS256 key=secret authenticate=auth::check() duration_grant=1h grant_subject=user
type Access struct{}

// no orm annotation
type Ignore struct{}
`

	path := filepath.Join(dir, "models.go")
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	pkg, err := ParseDir(dir)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if pkg.Name != "sample" {
		t.Fatalf("expected package sample, got %s", pkg.Name)
	}
	if len(pkg.Models) != 3 {
		t.Fatalf("expected 3 models, got %d", len(pkg.Models))
	}
	if pkg.Imports["time"] != "time" {
		t.Fatalf("expected time import to be used")
	}

	user := pkg.Models[0]
	if user.Kind != "node" || user.Table != "users" {
		t.Fatalf("unexpected user model: %+v", user)
	}
	if user.Permissions != "full" {
		t.Fatalf("expected permissions full")
	}
	if len(user.Fields) < 2 {
		t.Fatalf("expected user fields")
	}
	if user.Fields[0].DBName != "first_name" {
		t.Fatalf("expected overridden field name first_name, got %s", user.Fields[0].DBName)
	}
	if user.Fields[1].DBName != "last_name" {
		t.Fatalf("expected renamed field last_name, got %s", user.Fields[1].DBName)
	}
	if user.Fields[0].Permissions != "none" {
		t.Fatalf("expected field permissions none")
	}
	if user.Fields[0].DefaultExpr != "now()" {
		t.Fatalf("expected default expr now(), got %s", user.Fields[0].DefaultExpr)
	}
	if user.Fields[0].TypeHint != "uuid" {
		t.Fatalf("expected type hint uuid, got %s", user.Fields[0].TypeHint)
	}

	edge := pkg.Models[1]
	if edge.Kind != "edge" || edge.EdgeIn != "user" || edge.EdgeOut != "account" {
		t.Fatalf("unexpected edge model: %+v", edge)
	}

	access := pkg.Models[2]
	if access.Kind != "access" || access.Access.Name != "acc" || access.Access.Scope != "database" {
		t.Fatalf("unexpected access model: %+v", access)
	}
}

func TestMergeAnnotationArgs(t *testing.T) {
	anns := []Annotation{
		{Kind: "field", Args: map[string]string{"name": "a"}},
		{Kind: "field", Args: map[string]string{"name": "b", "type": "string"}},
		{Kind: "node", Args: map[string]string{"table": "t"}},
	}
	out := mergeAnnotationArgs(anns, "field")
	if out["name"] != "b" || out["type"] != "string" {
		t.Fatalf("unexpected merge: %+v", out)
	}
}

func TestParseDirSkipsTestAndGeneratedFiles(t *testing.T) {
	dir := t.TempDir()
	src := `package sample

// orm:node table=users
type User struct{}
`
	if err := os.WriteFile(filepath.Join(dir, "model_test.go"), []byte(src), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "model_orm_gen.go"), []byte(src), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if _, err := ParseDir(dir); err == nil {
		t.Fatalf("expected error when only skipped files exist")
	}
}
