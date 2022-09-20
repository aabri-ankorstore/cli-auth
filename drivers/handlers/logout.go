package handlers

import (
	"context"
	"github.com/aabri-ankorstore/cli-auth/pkg/repository"
	"github.com/aabri-ankorstore/cli-auth/utils"
	"github.com/uptrace/bun"
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
	repo := repository.AccessTokensRepository{
		DB:  utils.DB.DB,
		Ctx: context.Background(),
	}
	err = repo.Delete(session.Values["account_id"].(string))
	if err != nil {
		panic(err)
	}
	defer utils.DB.DB.(*bun.DB).DB.Close()

	delete(session.Values, "id_token")
	delete(session.Values, "access_token")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
	time.Sleep(3 * time.Second)
	os.Exit(0)
}
