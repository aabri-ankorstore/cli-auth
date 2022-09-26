package utils

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

type AuthStatus struct {
	IsConnected bool `json:"is_connected"`
}

func IsAuthenticated(r *http.Request) bool {
	session, err := SessionStore.Get(r, CookieName)
	if err != nil || session.Values["access_token"] == nil || session.Values["access_token"] == "" {
		log.Trace().Msg(err.Error())
		return false
	}
	return true
}
