package utils

import (
	"encoding/json"
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

func IsAuthenticatedViaWeb() bool {

	resp, err := http.Get("http://localhost:8080/is-authenticated")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	type status struct {
		IsAuthenticated bool `json:"IsAuthenticated"`
	}
	var data status
	_ = decoder.Decode(&data)
	return data.IsAuthenticated
}
