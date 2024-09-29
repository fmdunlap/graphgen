package cmd

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"graphgen/internal/database"
)

var (
	migrationDirPath string
)

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.PersistentFlags().StringVarP(&migrationDirPath, "migrationpath", "m", "", "path to the migration directory")
	err := migrateCmd.MarkPersistentFlagRequired("migrationpath")
	if err != nil {
		panic(err)
	}

	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run db migrations",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run up migrations",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getMigrate(migrationDirPath)
		if err != nil {
			panic(err)
		}

		err = m.Up()
		if err != nil {
			panic(err)
		}
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Run down migrations",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := getMigrate(migrationDirPath)
		if err != nil {
			panic(err)
		}

		err = m.Down()
		if err != nil {
			panic(err)
		}
	},
}

func getMigrate(migrationDirPath string) (*migrate.Migrate, error) {
	db := database.New(&EnvConfig.Database)
	driver, err := postgres.WithInstance(db.GetInstance(), &postgres.Config{})
	if err != nil {
		panic(err)
	}
	return migrate.NewWithDatabaseInstance("file://"+migrationDirPath, "postgres", driver)
}
