package main

import "os"
import "testing"

func TestMainRuns(t *testing.T) {
	orig := os.Args
	os.Args = []string{"surreal-orm", "--help"}
	defer func() { os.Args = orig }()
	main()
}
