package cli

import (
	"github.com/spf13/cobra"
	"github.com/yaroher/surrealdb.go.orm/internal/codegen"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate ORM resources from Go source",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("dir")
		return codegen.Generate(codegen.Options{Dirs: []string{dir}})
	},
}

func init() {
	generateCmd.Flags().String("dir", ".", "package directory to scan")
}
