package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/yaroher/surrealdb.go.orm/internal/codegen"
	"github.com/yaroher/surrealdb.go.orm/internal/surreal"
	"github.com/yaroher/surrealdb.go.orm/pkg/migrator"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migration management",
}

var migrateInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, ctx, cleanup, err := buildMigrator(cmd)
		if err != nil {
			return err
		}
		defer cleanup()
		return m.Init(ctx)
	},
}

var migrateGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate migrations based on code vs DB state",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, ctx, cleanup, err := buildMigrator(cmd)
		if err != nil {
			return err
		}
		defer cleanup()

		codeDir, _ := cmd.Flags().GetString("code")
		name, _ := cmd.Flags().GetString("name")

		pkg, err := codegen.ParseDir(codeDir)
		if err != nil {
			return err
		}
		if len(pkg.Models) == 0 {
			return fmt.Errorf("no models found in %s", codeDir)
		}
		code := codegen.BuildResourceSet(pkg.Models)
		db, err := migrator.Introspect(ctx, m.DB)
		if err != nil {
			return err
		}
		_, err = m.Generate(ctx, code, db, name)
		return err
	},
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, ctx, cleanup, err := buildMigrator(cmd)
		if err != nil {
			return err
		}
		defer cleanup()
		steps, _ := cmd.Flags().GetInt("steps")
		return m.Up(ctx, steps)
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, ctx, cleanup, err := buildMigrator(cmd)
		if err != nil {
			return err
		}
		defer cleanup()
		steps, _ := cmd.Flags().GetInt("steps")
		return m.Down(ctx, steps)
	},
}

var migrateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List migration status",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, ctx, cleanup, err := buildMigrator(cmd)
		if err != nil {
			return err
		}
		defer cleanup()
		list, err := m.List(ctx)
		if err != nil {
			return err
		}
		for _, item := range list {
			status := "pending"
			if item.Applied {
				status = "applied"
			}
			fmt.Printf("%s\t%s\t%s\n", item.Migration.ID, status, item.AppliedAt)
		}
		return nil
	},
}

var migrateResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, ctx, cleanup, err := buildMigrator(cmd)
		if err != nil {
			return err
		}
		defer cleanup()
		return m.Reset(ctx)
	},
}

var migratePruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prune migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, ctx, cleanup, err := buildMigrator(cmd)
		if err != nil {
			return err
		}
		defer cleanup()
		removed, err := m.Prune(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("pruned %d migration files\n", removed)
		return nil
	},
}

func init() {
	migrateCmd.AddCommand(migrateInitCmd)
	migrateCmd.AddCommand(migrateGenerateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateListCmd)
	migrateCmd.AddCommand(migrateResetCmd)
	migrateCmd.AddCommand(migratePruneCmd)

	migrateCmd.PersistentFlags().String("dir", "migrations", "migrations directory")
	migrateCmd.PersistentFlags().String("mode", "strict", "strict or lax")
	migrateCmd.PersistentFlags().String("dsn", "", "SurrealDB connection string")
	migrateCmd.PersistentFlags().String("ns", "", "namespace")
	migrateCmd.PersistentFlags().String("db", "", "database")
	migrateCmd.PersistentFlags().String("username", "", "username")
	migrateCmd.PersistentFlags().String("password", "", "password")
	migrateCmd.PersistentFlags().Bool("force", false, "force non-interactive behavior")
	migrateCmd.PersistentFlags().Bool("two-way", true, "generate two-way migrations")
	migrateCmd.PersistentFlags().String("rename-strategy", "prompt", "rename strategy: prompt|rename|delete|keep")
	migrateCmd.PersistentFlags().String("rename-expr", "", "rename copy expression (use {old}, {new}, {table})")
	migrateCmd.PersistentFlags().Bool("grants-always", false, "always include ACCESS GRANT statements in generate")

	migrateGenerateCmd.Flags().String("code", ".", "package directory to scan for models")
	migrateGenerateCmd.Flags().String("name", "migration", "migration name")

	migrateUpCmd.Flags().Int("steps", 0, "number of migrations to apply")
	migrateDownCmd.Flags().Int("steps", 1, "number of migrations to rollback")

	migrateCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		fmt.Println("\nNote: migrator engine is in active development.")
	})
}

var connectFn = func(ctx context.Context, cfg migrator.Config) (migrator.DB, func(), error) {
	client, err := surreal.Connect(ctx, cfg.DSN, cfg.NS, cfg.DB, cfg.Username, cfg.Password)
	if err != nil {
		return nil, func() {}, err
	}
	cleanup := func() {
		_ = client.Close(context.Background())
	}
	return surreal.Adapter{DB: client}, cleanup, nil
}

func buildMigrator(cmd *cobra.Command) (*migrator.Migrator, context.Context, func(), error) {
	dir, _ := cmd.Flags().GetString("dir")
	mode, _ := cmd.Flags().GetString("mode")
	dsn, _ := cmd.Flags().GetString("dsn")
	ns, _ := cmd.Flags().GetString("ns")
	dbName, _ := cmd.Flags().GetString("db")
	user, _ := cmd.Flags().GetString("username")
	pass, _ := cmd.Flags().GetString("password")
	force, _ := cmd.Flags().GetBool("force")
	twoWay, _ := cmd.Flags().GetBool("two-way")
	renameStrategy, _ := cmd.Flags().GetString("rename-strategy")
	renameExpr, _ := cmd.Flags().GetString("rename-expr")
	grantsAlways, _ := cmd.Flags().GetBool("grants-always")

	if dsn == "" {
		return nil, nil, func() {}, fmt.Errorf("dsn is required")
	}

	cfg := migrator.Config{
		Dir:            dir,
		Mode:           migrator.Mode(mode),
		DSN:            dsn,
		NS:             ns,
		DB:             dbName,
		Username:       user,
		Password:       pass,
		Force:          force,
		TwoWay:         twoWay,
		RenameStrategy: renameStrategy,
		RenameExpr:     renameExpr,
		GrantsAlways:   grantsAlways,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	db, closeFn, err := connectFn(ctx, cfg)
	if err != nil {
		cancel()
		return nil, nil, func() {}, err
	}
	cleanup := func() {
		closeFn()
		cancel()
	}

	var prompter migrator.Prompter
	if cfg.Force {
		prompter = migrator.NoPrompter{}
	} else {
		prompter = migrator.RealPrompter{}
	}
	m := migrator.New(db, cfg, prompter)
	return m, ctx, cleanup, nil
}

func initEnv() {
	if _, ok := os.LookupEnv("SURREAL_ORM_DEBUG"); ok {
		fmt.Println("debug enabled")
	}
}
