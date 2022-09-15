package utils

import (
	"embed"
	"github.com/gorilla/sessions"
)

var (
	Tpl          embed.FS
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
