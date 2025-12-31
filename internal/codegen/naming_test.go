package codegen

import "testing"

func TestNamingHelpers(t *testing.T) {
	if got := toSnake("UserID"); got != "user_id" {
		t.Fatalf("expected user_id, got %s", got)
	}
	if got := toPascal("user_id"); got != "UserId" {
		t.Fatalf("expected UserId, got %s", got)
	}
	if got := toCamel("UserID"); got != "userId" {
		t.Fatalf("expected userId, got %s", got)
	}
	if got := applyRenameAll("FieldName", "snake_case"); got != "field_name" {
		t.Fatalf("unexpected snake_case: %s", got)
	}
	if got := applyRenameAll("FieldName", "SCREAMING_SNAKE_CASE"); got != "FIELD_NAME" {
		t.Fatalf("unexpected screaming: %s", got)
	}
	if got := applyRenameAll("FieldName", "PascalCase"); got != "FieldName" {
		t.Fatalf("unexpected pascal: %s", got)
	}
	if got := applyRenameAll("FieldName", "camelCase"); got != "fieldName" {
		t.Fatalf("unexpected camel: %s", got)
	}
	if got := applyRenameAll("FieldName", "lowercase"); got != "fieldname" {
		t.Fatalf("unexpected lowercase: %s", got)
	}
	if got := applyRenameAll("FieldName", "UPPERCASE"); got != "FIELDNAME" {
		t.Fatalf("unexpected uppercase: %s", got)
	}
	if got := applyRenameAll("FieldName", "unknown"); got != "FieldName" {
		t.Fatalf("unexpected default: %s", got)
	}

	if charClass('9') != 1 || charClass('A') != 2 || charClass('a') != 3 || charClass('-') != 4 {
		t.Fatalf("unexpected charClass results")
	}
	if receiverName("User") != "u" || receiverName("") != "m" {
		t.Fatalf("unexpected receiverName")
	}
	if pathBase("github.com/acme/pkg") != "pkg" || pathBase("") != "" {
		t.Fatalf("unexpected pathBase")
	}
	if got := toCamel(""); got != "" {
		t.Fatalf("expected empty toCamel")
	}
	if got := capitalize(""); got != "" {
		t.Fatalf("expected empty capitalize")
	}
	if got := toSnake("user_id"); got != "user_id" {
		t.Fatalf("expected user_id, got %s", got)
	}
}
