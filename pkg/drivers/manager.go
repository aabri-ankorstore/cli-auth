package drivers

import (
	"fmt"
	verifier "github.com/okta/okta-jwt-verifier-golang"
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
