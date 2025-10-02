/*
Copyright © 2025 Peter Shaan <petershaan12@gmail.com>
*/
package cmd

import (
	"github.com/labstack/gommon/log"
	"github.com/petershaan12/go-auth-clean-arch/package/library"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations using goose here",
	Long: `Run database migrations using goose here.
You can use this command to apply or rollback database migrations as needed.`,
	Run: migrate,
}

var direction string

func init() {
	migrateCmd.Flags().StringVar(&direction, "direction", "", "Migration direction: up or down")
	rootCmd.AddCommand(migrateCmd)
}

func migrate(cmd *cobra.Command, args []string) {
	library.ModuleConfig()
	if direction == "up" || direction == "down" {
		goose.SetDialect("mysql")
		db, err := library.GetSqlDB()
		if err != nil {
			log.Fatal("Failed to connect database: ", err)
		}
		defer db.Close()
		migrationsDir := "migrations"
		if direction == "up" {
			if err := goose.Up(db, migrationsDir); err != nil {
				log.Fatal("Goose up migration failed: ", err)
			}
		} else {
			if err := goose.Down(db, migrationsDir); err != nil {
				log.Fatal("Goose down migration failed: ", err)
			}
		}
		return
	}
}
