package migrator

import "testing"

func TestNotImplemented(t *testing.T) {
	err := NotImplemented("x")
	if err == nil || err.Error() == "" {
		t.Fatalf("expected error")
	}
}
