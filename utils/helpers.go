package utils

import (
	"context"
	"github.com/aabri-ankorstore/cli-auth/pkg/database/adapters/sqlite"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"net/http"
)

func IsAuthenticated(r *http.Request) bool {
	session, err := SessionStore.Get(r, CookieName)
	if err != nil || session.Values["access_token"] == nil || session.Values["access_token"] == "" {
		log.Debug().Err(err)
		return false
	}
	return true
}

func RunMigration() {
	// Run migration
	ctx := context.Background()
	// Initialize db for migrations
	db, err := sqlite.InitDB(false)
	if err != nil {
		log.Printf("failed to initialize db", "err", err)
	}
	// Run migrations.
	if err = db.RunMigrations(ctx); err != nil {
		log.Printf("failed to migrate DB: %w", err)
	}
	if err = db.DB.(*bun.DB).DB.Close(); err != nil {
		log.Printf("failed to close migrations DB: %w", err)
	}
}
