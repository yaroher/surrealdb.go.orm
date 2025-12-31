package codegen

import "testing"

func TestInferSurrealType(t *testing.T) {
	model := Model{Table: "user"}
	cases := []struct {
		name   string
		goType string
		field  Field
		want   string
	}{
		{
			name:   "type hint",
			goType: "int",
			field:  Field{TypeHint: "uuid"},
			want:   "uuid",
		},
		{
			name:   "link one annotation",
			goType: "string",
			field:  Field{LinkOne: "Account"},
			want:   "record<account>",
		},
		{
			name:   "link self annotation",
			goType: "string",
			field:  Field{LinkSelf: "User"},
			want:   "record<user>",
		},
		{
			name:   "link many annotation",
			goType: "string",
			field:  Field{LinkMany: "Account"},
			want:   "array<record<account>>",
		},
		{
			name:   "link one generic",
			goType: "LinkOne[Account]",
			field:  Field{},
			want:   "record<account>",
		},
		{
			name:   "link self generic",
			goType: "LinkSelf[User]",
			field:  Field{},
			want:   "record<user>",
		},
		{
			name:   "link many generic",
			goType: "LinkMany[Account]",
			field:  Field{},
			want:   "array<record<account>>",
		},
		{
			name:   "id type",
			goType: "ID[User]",
			field:  Field{},
			want:   "record<user>",
		},
		{
			name:   "simple id",
			goType: "SimpleID[User]",
			field:  Field{},
			want:   "record<user>",
		},
		{
			name:   "slice",
			goType: "[]string",
			field:  Field{},
			want:   "array",
		},
		{
			name:   "map",
			goType: "map[string]any",
			field:  Field{},
			want:   "object",
		},
		{
			name:   "bytes",
			goType: "[]byte",
			field:  Field{},
			want:   "array",
		},
		{
			name:   "time",
			goType: "time.Time",
			field:  Field{},
			want:   "datetime",
		},
		{
			name:   "duration",
			goType: "time.Duration",
			field:  Field{},
			want:   "duration",
		},
	}

	for _, tc := range cases {
		if got := inferSurrealType(tc.goType, tc.field, model); got != tc.want {
			t.Fatalf("%s: expected %s, got %s", tc.name, tc.want, got)
		}
	}
}

func TestExtractGenericType(t *testing.T) {
	if got := extractGenericType("LinkOne[Account]"); got != "Account" {
		t.Fatalf("unexpected generic: %s", got)
	}
	if got := extractGenericType("LinkOne[*pkg.Account]"); got != "Account" {
		t.Fatalf("unexpected generic pointer: %s", got)
	}
	if got := extractGenericType("LinkOne"); got != "" {
		t.Fatalf("expected empty generic")
	}
}
