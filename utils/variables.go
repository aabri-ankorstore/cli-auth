package utils

import (
	"embed"
	"github.com/ankorstore/ankorstore-cli-core/pkg/plugin"
	"github.com/gernest/qlstore"
	"html/template"
)

var (
	Tpl          embed.FS
	View         *template.Template
	SessionStore *qlstore.QLStore
	AuthDriver   string
	ClientSECRET string
	ISSUER       string
	RedirectUri  string
	ClientID     string
	CookieName   = "session-name"
	State        = "none"
	Nonce        = "NonceNotSetYet"
	PluginRepo   = "https://github.com/ankorstore/ankor-auth-plugin"
	PluginPath   = plugin.Encode(PluginRepo)
)

var KeyPair = [][]byte{
	[]byte("ePAPW9vJv7gHoftvQTyNj5VkWB52mlza"),
	[]byte("N8SmpJ00aSpepNrKoyYxmAJhwVuKEWZD"),
}
