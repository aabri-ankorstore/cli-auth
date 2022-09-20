package drivers

import (
	"context"
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/pkg/database/adapters/sqlite"
	verifier "github.com/okta/okta-jwt-verifier-golang"
	"github.com/uptrace/bun"
	"log"
	"net/http"
)

const (
	host string = "http://localhost"
)

var (
	db  *sqlite.SqliteDB
	con *sqlite.SqliteDB
)

type Manager interface {
	InformUserAndOpenBrowser() error
	ExchangeCode(w http.ResponseWriter, r *http.Request) (Exchange, error)
	VerifyToken(t string) (*verifier.Jwt, error)
	GetProfile(r *http.Request) (map[string]string, error)
}

func GetAuth(authType string) (Manager, error) {
	switch authType {
	case "okta":
		return NewOktaClient(), nil
	case "github":
		return NewGithubClient(), nil
	default:
		return nil, fmt.Errorf("wrong auth type passed")
	}
}
func init() {
	// Initialize db
	con, err := sqlite.InitDB(false)
	if err != nil {
		log.Printf("failed to initialize db", "err", err)
	}
	db = con
	// Run migration
	runMigration()
}
func runMigration() {
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
