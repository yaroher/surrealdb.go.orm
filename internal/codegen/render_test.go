package codegen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderOutput(t *testing.T) {
	dir := t.TempDir()
	pkg := Package{
		Name: "sample",
		Imports: map[string]string{
			"time":      "time",
			"acmealias": "github.com/acme/alias",
			"_":         "ignored",
			".":         "dot",
		},
		Models: []Model{
			{
				Name:        "User",
				Kind:        "node",
				Table:       "user",
				RenameAll:   "snake_case",
				Permissions: "full",
				Fields: []Field{
					{Name: "ID", Type: "string", DBName: "id"},
					{Name: "CreatedAt", Type: "time.Time", DBName: "created_at"},
				},
			},
			{
				Name:    "UserEdge",
				Kind:    "edge",
				Table:   "user_edge",
				EdgeIn:  "user",
				EdgeOut: "account",
			},
			{
				Name: "Access",
				Kind: "access",
				Access: AccessConfig{
					Name:         "acc",
					Scope:        "database",
					Type:         "jwt",
					Algorithm:    "HS256",
					Key:          "secret",
					Authenticate: "auth::check()",
					GrantSubject: "user",
				},
			},
		},
	}

	if err := Render(dir, pkg); err != nil {
		t.Fatalf("render: %v", err)
	}
	outPath := filepath.Join(dir, "orm_gen.go")
	content, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	out := string(content)
	checks := []string{
		"package sample",
		"acmealias \"github.com/acme/alias\"",
		"type UserSchema struct",
		"ID qb.Field[string]",
		"CreatedAt qb.Field[time.Time]",
		"return qb.T(\"user\")",
		"qb.F[string](\"id\")",
		"qb.F[time.Time](\"created_at\")",
		"res.AddTable(\"user\"",
		".PermissionsFull()",
		"res.AddField(\"user\", \"created_at\"",
		"res.AddField(\"user_edge\", \"in\"",
		"res.AddAccess(\"acc\", \"database\"",
		"qb.DefineAccess(\"acc\").OnDatabase().TypeJWT().JWTAlgorithmKey(\"HS256\", \"secret\")",
		".Authenticate(qb.Raw(\"auth::check()\"))",
		"res.AddAccessGrant(\"acc\"",
	}
	for _, c := range checks {
		if !strings.Contains(out, c) {
			t.Fatalf("expected output to contain %q", c)
		}
	}
	if strings.Contains(out, "ignored") || strings.Contains(out, "dot") {
		t.Fatalf("did not expect underscore/dot imports to be rendered")
	}
}
