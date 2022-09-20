package utils

import (
	"embed"
	"github.com/aabri-ankorstore/cli-auth/pkg/database/adapters/sqlite"
	"github.com/gorilla/sessions"
	"html/template"
)

var (
	Tpl          embed.FS
	View         *template.Template
	SessionStore *sessions.CookieStore
	AuthDriver   string
	ClientSECRET string
	ISSUER       string
	RedirectUri  string
	ClientID     string
	CookieName   = "login-session-store"
	State        = "none"
	Nonce        = "NonceNotSetYet"
	DB           *sqlite.SqliteDB
)
