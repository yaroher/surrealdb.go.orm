package codegen

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateNoDirs(t *testing.T) {
	if err := Generate(Options{}); err == nil {
		t.Fatalf("expected error for no dirs")
	}
}

func TestGenerateSkipsNoModels(t *testing.T) {
	dir := t.TempDir()
	src := `package sample

type Plain struct{}
`
	if err := os.WriteFile(filepath.Join(dir, "plain.go"), []byte(src), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := Generate(Options{Dirs: []string{dir}}); err != nil {
		t.Fatalf("generate: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "orm_gen.go")); err == nil {
		t.Fatalf("did not expect orm_gen.go to be created")
	}
}

func TestGenerateCreatesFile(t *testing.T) {
	dir := t.TempDir()
	src := `package sample

// orm:node table=users
type User struct{ Name string }
`
	if err := os.WriteFile(filepath.Join(dir, "model.go"), []byte(src), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := Generate(Options{Dirs: []string{dir}}); err != nil {
		t.Fatalf("generate: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "orm_gen.go")); err != nil {
		t.Fatalf("expected orm_gen.go to be created")
	}
}

func TestGenerateParseError(t *testing.T) {
	dir := t.TempDir()
	src := `package sample

func {` // invalid
	if err := os.WriteFile(filepath.Join(dir, "bad.go"), []byte(src), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := Generate(Options{Dirs: []string{dir}}); err == nil {
		t.Fatalf("expected parse error")
	}
}
