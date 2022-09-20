package migrations

import (
	"context"
	"embed"
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/db/migrator"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"strings"
)

// Migrations are all migrations.
var Migrations *migrator.Migrations

//go:embed *.sql
var migrations embed.FS

// DropTables deletes all database tables in the `public` schema of the configured database.
//
// Note: this method assumes that PostgresSQL is used as the underlying db.
func DropTables(ctx context.Context, db *bun.DB) error {
	rows, err := db.QueryContext(ctx, "SELECT * FROM pg_tables")
	if err != nil {
		return err
	}
	if err = rows.Err(); err != nil {
		return err
	}

	var results []map[string]interface{}
	if err = db.ScanRows(ctx, rows, &results); err != nil {
		return err
	}

	for _, result := range results {
		tableName := string(result["tablename"].([]byte))
		schemaName := string(result["schemaname"].([]byte))

		// Skip internal pg tables.
		if strings.Compare(schemaName, "public") != 0 {
			continue
		}

		log.Info().Msgf("dropping table", "table_name", tableName)

		_, err := db.ExecContext(ctx, fmt.Sprintf("DROP TABLE \"%s\" CASCADE", tableName))
		if err != nil {
			log.Printf("failed dropping table", "table_name", tableName, "err", err)
			continue
		}
	}

	return nil
}

// Init initializes the migrator tables.
func Init(ctx context.Context, db *bun.DB) error {
	// Initialize the migrator.
	m := migrator.NewMigrator(db, Migrations)
	return m.Init(ctx)
}

// Migrate migrates the DB to the latest version.
func Migrate(ctx context.Context, db *bun.DB) error {
	// Initialize the migrator.
	m := migrator.NewMigrator(db, Migrations)

	// Run migrations.
	if err := m.Migrate(ctx); err != nil {
		log.Printf("failed to migrate db", "err", err)
		return err
	}

	status, err := m.MigrationsWithStatus(ctx)
	if err != nil {
		return err
	}
	log.Printf("migration done", "applied", status.Applied(), "status", status.String(), "unaplied", status.Unapplied())

	return nil
}

func init() {
	Migrations = migrator.NewMigrations()
	for _, m := range []migrator.Migration{
		// Initial.
		{
			Name: "20211213143752",
			Up:   migrator.NewSQLMigrationFunc(migrations, "20211213143752_initial.up.sql"),
		},
	} {
		Migrations.Add(m)
	}
}
