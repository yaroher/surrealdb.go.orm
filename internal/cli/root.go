package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "surreal-orm",
	Short: "Surreal ORM toolchain",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(generateCmd)
}
