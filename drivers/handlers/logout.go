package handlers

import (
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/utils"
	"net/http"
)

func (h *Auth) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := utils.SessionStore.Get(r, utils.CookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	delete(session.Values, "id_token")
	delete(session.Values, "access_token")
	session.Save(r, w)

	// remove session from db
	fmt.Println("Good Bye")

	http.Redirect(w, r, "/", http.StatusFound)
}
