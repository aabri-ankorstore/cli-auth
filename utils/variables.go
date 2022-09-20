package utils

import (
	"embed"
	"github.com/michaeljs1990/sqlitestore"
	"html/template"
)

var (
	Tpl          embed.FS
	View         *template.Template
	SessionStore *sqlitestore.SqliteStore
	AuthDriver   string
	ClientSECRET string
	ISSUER       string
	RedirectUri  string
	ClientID     string
	CookieName   = "login-session-store"
	State        = "none"
	Nonce        = "NonceNotSetYet"
)
