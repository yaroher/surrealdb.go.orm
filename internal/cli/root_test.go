package cli

import (
	"os"
	"testing"
)

func TestExecute(t *testing.T) {
	orig := os.Args
	os.Args = []string{"surreal-orm", "--help"}
	defer func() { os.Args = orig }()
	Execute()
}
