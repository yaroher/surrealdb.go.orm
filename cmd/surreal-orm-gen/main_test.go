package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

func TestMainRuns(t *testing.T) {
	dir := t.TempDir()
	src := `package sample

// orm:node table=users
type User struct{ Name string }
`
	if err := os.WriteFile(filepath.Join(dir, "model.go"), []byte(src), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	origArgs := os.Args
	origCmd := flag.CommandLine
	os.Args = []string{"surreal-orm-gen", "-dir", dir}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	defer func() {
		os.Args = origArgs
		flag.CommandLine = origCmd
	}()

	main()
}
