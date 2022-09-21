package utils

import (
	"embed"
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
)
