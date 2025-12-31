package codegen

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/yaroher/surrealdb.go.orm/pkg/migrator"
	"github.com/yaroher/surrealdb.go.orm/pkg/qb"
)

func TestBuildResourceSetVariants(t *testing.T) {
	model := Model{
		Kind:        "edge",
		Table:       "edge_table",
		EdgeIn:      "user",
		EdgeOut:     "account",
		SchemaFull:  true,
		SchemaLess:  true,
		Drop:        true,
		Permissions: "none",
		Fields: []Field{
			{Name: "Meta", Type: "map[string]any", DBName: "meta"},
		},
	}
	res := BuildResourceSet([]Model{model})
	if _, ok := res.Tables["edge_table"]; !ok {
		t.Fatalf("expected edge table")
	}
	if _, ok := res.Fields["edge_table"]["in"]; !ok {
		t.Fatalf("expected edge in field")
	}
	if _, ok := res.Fields["edge_table"]["out"]; !ok {
		t.Fatalf("expected edge out field")
	}
}

func TestAddAccessResourceVariants(t *testing.T) {
	res := migrator.NewResourceSet()
	addAccessResource(&res, AccessConfig{
		Name:         "jwt",
		Scope:        "namespace",
		Type:         "jwt",
		URL:          "https://example.com/jwks",
		Overwrite:    true,
		Authenticate: "auth::check()",
	})
	addAccessResource(&res, AccessConfig{
		Name:          "rec",
		Scope:         "database",
		Type:          "record",
		Algorithm:     "HS256",
		Key:           "secret",
		RecordSignup:  "create user",
		RecordSignin:  "select user",
		RecordIssuer:  "issuer",
		RecordRefresh: true,
		GrantSubject:  "user",
		GrantToken:    "token",
		GrantDuration: "1h",
		IfNotExists:   true,
	})
	addAccessResource(&res, AccessConfig{
		Name:      "bear",
		Scope:     "database",
		Type:      "bearer",
		Algorithm: "record",
	})
	if len(res.Access) != 3 {
		t.Fatalf("expected access definitions")
	}
	if len(res.AccessGrants) != 1 {
		t.Fatalf("expected access grant")
	}
}

func TestApplyPermissionsHelpers(t *testing.T) {
	stmt := qb.DefineTableName("t")
	applyPermissionsTable(stmt, "full")
	if got := qb.Build(stmt).Text; got != "DEFINE TABLE t PERMISSIONS FULL" {
		t.Fatalf("unexpected full: %s", got)
	}

	stmt2 := qb.DefineFieldName("f", "t")
	applyPermissionsField(stmt2, "none")
	if got := qb.Build(stmt2).Text; got != "DEFINE FIELD f ON TABLE t PERMISSIONS NONE" {
		t.Fatalf("unexpected none: %s", got)
	}

	stmt3 := qb.DefineFieldName("f", "t")
	applyPermissionsField(stmt3, "")
	if got := qb.Build(stmt3).Text; got != "DEFINE FIELD f ON TABLE t" {
		t.Fatalf("unexpected empty permissions: %s", got)
	}
}

func TestRenderAccessAndFieldResources(t *testing.T) {
	var buf bytes.Buffer
	ac := AccessConfig{
		Name:            "acc",
		Scope:           "database",
		Type:            "record",
		URL:             "https://example.com",
		RecordIssuer:    "issuer",
		RecordRefresh:   true,
		Authenticate:    "auth()",
		DurationGrant:   "1h",
		DurationToken:   "2h",
		DurationSession: "3h",
	}
	renderAccessResource(&buf, ac)
	out := buf.String()
	if out == "" {
		t.Fatalf("expected renderAccessResource output")
	}

	buf.Reset()
	model := Model{Table: "user"}
	field := Field{
		Name:        "Meta",
		Type:        "map[string]any",
		DBName:      "meta",
		ValueExpr:   "value()",
		AssertExpr:  "true",
		DefaultExpr: "false",
		Permissions: "full",
	}
	renderFieldResource(&buf, model, field)
	if buf.String() == "" {
		t.Fatalf("expected renderFieldResource output")
	}
}

func TestParseDirErrors(t *testing.T) {
	if _, err := ParseDir(t.TempDir()); err == nil {
		t.Fatalf("expected error for empty dir")
	}
}

func TestCollectAndMarkImports(t *testing.T) {
	src := `package sample

import (
	 t "time"
	 _ "net/http"
)

type User struct { CreatedAt t.Time }
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "sample.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	imports := collectImports(file)
	used := map[string]string{}
	var fieldExpr ast.Expr
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gen.Specs {
			tspec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			st, ok := tspec.Type.(*ast.StructType)
			if !ok || len(st.Fields.List) == 0 {
				continue
			}
			fieldExpr = st.Fields.List[0].Type
		}
	}
	if fieldExpr == nil {
		t.Fatalf("expected field expr")
	}
	markTypeImports(fieldExpr, imports, used)
	if imports["t"] != "time" || used["t"] != "time" {
		t.Fatalf("expected time import")
	}
}

func TestNormalizeRef(t *testing.T) {
	if got := normalizeRef(""); got != "" {
		t.Fatalf("expected empty normalize")
	}
	if got := normalizeRef("user"); got != "user" {
		t.Fatalf("expected user normalize")
	}
	if got := normalizeRef("user:1"); got != "user:1" {
		t.Fatalf("expected record ref")
	}
	if got := normalizeRef("pkg.User"); got != "pkg.User" {
		t.Fatalf("expected pkg ref")
	}
	if got := normalizeRef("path/user"); got != "path/user" {
		t.Fatalf("expected path ref")
	}
}
