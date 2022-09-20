package handlers

import (
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/utils"
	"github.com/ankorstore/ankorstore-cli-core/pkg/plugin"
	"github.com/ankorstore/ankorstore-cli-core/pkg/util"
	"net/http"
	"os"
	"time"
)

func (h *Auth) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := utils.SessionStore.Get(r, utils.CookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// remove session from db
	dirs := util.NewDirs()
	p := dirs.GetPluginsDir()
	PluginRepo := "https://github.com/ankorstore/ankor-auth-plugin"
	PluginPath := plugin.Encode(PluginRepo)
	_ = os.Remove(fmt.Sprintf("%s/%s/sessions.db", p, PluginPath))
	delete(session.Values, "id_token")
	delete(session.Values, "access_token")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
	time.Sleep(3 * time.Second)
	os.Exit(0)
}
