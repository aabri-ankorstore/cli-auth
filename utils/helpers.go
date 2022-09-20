package utils

import (
	"github.com/rs/zerolog/log"
)

func IsAuthenticated() bool {
	session, err := SessionStore.Get(HttpRequest, CookieName)
	if err != nil || session.Values["access_token"] == nil || session.Values["access_token"] == "" {
		log.Debug().Err(err)
		return false
	}
	return true
}
