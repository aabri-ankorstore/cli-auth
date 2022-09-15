package utils

import (
	"net/http"
)

func IsAuthenticated(r *http.Request) bool {
	session, err := SessionStore.Get(r, CookieName)
	if err != nil || session.Values["access_token"] == nil || session.Values["access_token"] == "" {
		return false
	}
	return true
}
