package codegen

import "testing"

func TestSplitTokensQuoted(t *testing.T) {
	in := "access name=\"my access\" key='secret value' overwrite"
	got := splitTokens(in)
	want := []string{"access", "name=my access", "key=secret value", "overwrite"}
	if len(got) != len(want) {
		t.Fatalf("expected %d tokens, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected token[%d]: %q expected %q", i, got[i], want[i])
		}
	}
}

func TestParseAnnotationQuoted(t *testing.T) {
	line := "// orm:access name=\"my access\" key='secret value' overwrite"
	ann, ok := ParseAnnotation(line)
	if !ok {
		t.Fatalf("expected annotation to parse")
	}
	if ann.Kind != "access" {
		t.Fatalf("expected kind access, got %s", ann.Kind)
	}
	if ann.Args["name"] != "my access" {
		t.Fatalf("unexpected name: %q", ann.Args["name"])
	}
	if ann.Args["key"] != "secret value" {
		t.Fatalf("unexpected key: %q", ann.Args["key"])
	}
	if ann.Args["overwrite"] != "true" {
		t.Fatalf("expected overwrite=true, got %q", ann.Args["overwrite"])
	}
}

func TestParseAnnotationInvalid(t *testing.T) {
	if _, ok := ParseAnnotation("// nothing here"); ok {
		t.Fatalf("expected no annotation")
	}
	if _, ok := ParseAnnotation("orm:"); ok {
		t.Fatalf("expected no annotation for empty directive")
	}
}

func TestTrimQuotes(t *testing.T) {
	if got := trimQuotes("\"x\""); got != "x" {
		t.Fatalf("expected x, got %s", got)
	}
	if got := trimQuotes("x"); got != "x" {
		t.Fatalf("expected x, got %s", got)
	}
}
