package sqlite

import (
	"context"
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/db/migrations"
	"github.com/uptrace/bun"
)

type SqliteDB struct {
	DB bun.IDB
}

func (db *SqliteDB) RunMigrations(ctx context.Context) error {
	// Initialize migrator.
	if err := migrations.Init(ctx, db.DB.(*bun.DB)); err != nil {
		return fmt.Errorf("init migrations: %w", err)
	}
	// Run migrations.
	if err := migrations.Migrate(ctx, db.DB.(*bun.DB)); err != nil {
		return fmt.Errorf("migrations: %w", err)
	}
	return nil
}
