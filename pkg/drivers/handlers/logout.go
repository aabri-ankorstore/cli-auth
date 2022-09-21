package handlers

import (
	"github.com/aabri-ankorstore/cli-auth/pkg/utils"
	"net/http"
	"os"
	"time"
)

func (h *Auth) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := utils.SessionStore.Get(r, utils.CookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// remove session
	_ = os.Remove(utils.LockFile.Name())

	delete(session.Values, "id_token")
	delete(session.Values, "access_token")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
	time.Sleep(3 * time.Second)
	os.Exit(0)
}
