package codegen

import "testing"

func TestParseAccessConfig(t *testing.T) {
	ann := Annotation{Kind: "access", Args: map[string]string{
		"access":           "acc",
		"on":               "database",
		"type":             "jwt",
		"alg":              "HS256",
		"key":              "secret",
		"authenticate":     "auth::check()",
		"duration_grant":   "1h",
		"duration_token":   "2h",
		"duration_session": "3h",
		"grant_subject":    "user",
		"grant_token":      "token",
		"grant_duration":   "1d",
		"overwrite":        "true",
	}}

	ac := parseAccessConfig(ann)
	if ac.Name != "acc" {
		t.Fatalf("expected name acc, got %s", ac.Name)
	}
	if !ac.Overwrite {
		t.Fatalf("expected overwrite=true")
	}
	if ac.Scope != "database" || ac.Type != "jwt" || ac.Algorithm != "HS256" || ac.Key != "secret" {
		t.Fatalf("unexpected access config: %+v", ac)
	}
	if ac.Authenticate != "auth::check()" || ac.DurationGrant != "1h" || ac.DurationToken != "2h" || ac.DurationSession != "3h" {
		t.Fatalf("unexpected auth/duration config: %+v", ac)
	}
	if ac.GrantSubject != "user" || ac.GrantToken != "token" || ac.GrantDuration != "1d" {
		t.Fatalf("unexpected grant config: %+v", ac)
	}
}
