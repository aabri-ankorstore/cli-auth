package drivers

import (
	"context"
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/pkg/database/adapters/sqlite"
	utils2 "github.com/aabri-ankorstore/cli-auth/utils"
	verifier "github.com/okta/okta-jwt-verifier-golang"
	"log"
	"net/http"
)

const (
	host string = "http://localhost"
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
	// Initialize db for migrations
	var err error
	utils2.DB, err = sqlite.InitDB(false)
	if err != nil {
		log.Printf("failed to initialize db", "err", err)
	}
}
func init() {
	// Run migration
	ctx := context.Background()
	// Run migrations.
	if err := utils2.DB.RunMigrations(ctx); err != nil {
		log.Printf("failed to migrate DB: %w", err)
	}
}
