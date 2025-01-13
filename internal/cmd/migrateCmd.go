package cmd

import (
	"fmt"
	"strings"

	"app/app/modules"
	"app/database/migrations"
	"app/internal/modules/log"

	"github.com/spf13/cobra"
	"github.com/uptrace/bun/migrate"
)

// Migrate Command
func Migrate() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "db",
		Args: NotReqArgs,
	}
	cmd.AddCommand(initCMD())
	cmd.AddCommand(createSQL())
	cmd.AddCommand(createGO())
	cmd.AddCommand(migrateCMD())
	cmd.AddCommand(rollbackCMD())
	cmd.AddCommand(statusCMD())
	cmd.AddCommand(markAppliedCMD())
	return cmd
}
func initCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "init",
		Long: "create migration tables",
		Run: func(cmd *cobra.Command, args []string) {
			migrator := migrate.NewMigrator(modules.Get().DB.Svc.DB(), migrations.Migrations)
			migrator.Init(cmd.Context())
		},
	}
	return cmd
}

func createSQL() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create_sql",
		Long: "create up and down SQL migrations",
		Run: func(cmd *cobra.Command, args []string) {
			migrator := migrate.NewMigrator(modules.Get().DB.Svc.DB(), migrations.Migrations)
			name := strings.Join(args, "_")
			files, err := migrator.CreateSQLMigrations(cmd.Context(), name)
			if err != nil {
				log.Error(err.Error())
				return
			}
			for _, mf := range files {
				fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
			}
		},
	}
	return cmd
}
func createGO() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create_go",
		Long: "create Go migration",
		Run: func(cmd *cobra.Command, args []string) {
			migrator := migrate.NewMigrator(modules.Get().DB.Svc.DB(), migrations.Migrations)
			name := strings.Join(args, "_")
			mf, err := migrator.CreateGoMigration(cmd.Context(), name)
			if err != nil {
				log.Error(err.Error())
				return
			}
			fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
		},
	}
	return cmd
}
func migrateCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use: "migrate",
		Run: func(cmd *cobra.Command, args []string) {
			migrator := migrate.NewMigrator(modules.Get().DB.Svc.DB(), migrations.Migrations)
			if err := migrator.Lock(cmd.Context()); err != nil {
				log.Error(err.Error())
				return
			}
			defer migrator.Unlock(cmd.Context()) //nolint:errcheck

			group, err := migrator.Migrate(cmd.Context())
			if err != nil {
				log.Error(err.Error())
				return
			}
			if group.IsZero() {
				fmt.Printf("there are no new migrations to run (database is up to date)\n")
				return
			}
			fmt.Printf("migrated to %s\n", group)
			return
		},
	}
	return cmd
}
func rollbackCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use: "rollback",
		Run: func(cmd *cobra.Command, args []string) {
			migrator := migrate.NewMigrator(modules.Get().DB.Svc.DB(), migrations.Migrations)
			group, err := migrator.Rollback(cmd.Context())
			if err != nil {
				log.Error(err.Error())
				return
			}

			if group.ID == 0 {
				fmt.Printf("there are no groups to roll back\n")
				return
			}

			fmt.Printf("rolled back %s\n", group)
			return
		},
	}
	return cmd
}

func statusCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "status",
		Long: "print migrations status",
		Run: func(cmd *cobra.Command, args []string) {
			migrator := migrate.NewMigrator(modules.Get().DB.Svc.DB(), migrations.Migrations)
			ms, err := migrator.MigrationsWithStatus(cmd.Context())
			if err != nil {
				log.Error(err.Error())
				return
			}
			fmt.Printf("migrations: %s\n", ms)
			fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
			fmt.Printf("last migration group: %s\n", ms.LastGroup())
			return
		},
	}
	return cmd
}

func markAppliedCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "mark_applied",
		Long: "mark migrations as applied without actually running them",
		Run: func(cmd *cobra.Command, args []string) {
			migrator := migrate.NewMigrator(modules.Get().DB.Svc.DB(), migrations.Migrations)
			group, err := migrator.Migrate(cmd.Context(), migrate.WithNopMigration())
			if err != nil {
				log.Error(err.Error())
				return
			}
			if group.IsZero() {
				fmt.Printf("there are no new migrations to mark as applied\n")
				return
			}
			fmt.Printf("marked as applied %s\n", group)
			return
		},
	}
	return cmd
}
