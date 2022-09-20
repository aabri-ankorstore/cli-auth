package utils

import (
	"context"
	"github.com/aabri-ankorstore/cli-auth/pkg/repository"
	"github.com/rs/zerolog/log"
	"net/http"
)

func IsAuthenticated(r *http.Request) bool {
	session, err := SessionStore.Get(r, CookieName)
	if err != nil || session.Values["access_token"] == nil || session.Values["access_token"] == "" {
		log.Debug().Err(err)
		return false
	}
	return true
}

func IsLoggedIn() bool {
	repo := repository.AccessTokensRepository{
		DB:  DB.DB,
		Ctx: context.Background(),
	}
	s, err := repo.GetAll()
	if err != nil {
		panic(err)
	}
	if len(s) > 0 {
		return true
	}
	return false
}
