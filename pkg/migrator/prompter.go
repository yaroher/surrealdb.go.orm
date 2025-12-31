package migrator

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Prompter decides how to handle ambiguous changes.
type Prompter interface {
	ConfirmRename(table, from, to string) bool
}

// RealPrompter prompts via stdin.
type RealPrompter struct{}

func (RealPrompter) ConfirmRename(table, from, to string) bool {
	fmt.Printf("Detected field rename in %s: %s -> %s. Rename? [y/N]: ", table, from, to)
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.ToLower(strings.TrimSpace(line))
	return line == "y" || line == "yes"
}

// NoPrompter always returns false.
type NoPrompter struct{}

func (NoPrompter) ConfirmRename(table, from, to string) bool { return false }
