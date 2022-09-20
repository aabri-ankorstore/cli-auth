package handlers

import (
	"context"
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/pkg/entities"
	"github.com/aabri-ankorstore/cli-auth/pkg/repository"
	utils2 "github.com/aabri-ankorstore/cli-auth/utils"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h *Auth) CallBackHandler(w http.ResponseWriter, r *http.Request) {
	e, _ := h.manager.ExchangeCode(w, r)
	if e.Error != "" {
		fmt.Println(e.Error)
		fmt.Println(e.ErrorDescription)
		return
	}
	// save access token
	repo := repository.AccessTokensRepository{
		DB:  utils2.DB.DB,
		Ctx: context.Background(),
	}
	profile, er := h.manager.GetProfile(r)
	if er != nil {
		log.Info().Err(er)
		return
	}
	fmt.Println(profile)
	accessToken := entities.AccessToken{
		AccountID:   profile["email"],
		AccessToken: e.AccessToken,
		IdToken:     e.IdToken,
	}
	err := repo.Insert(&accessToken)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
