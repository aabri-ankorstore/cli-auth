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
	req, reqerr := http.NewRequest("GET", "http://localhost:8080/is-authenticated", nil)
	if reqerr != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	type status struct {
		IsAuthenticated bool `json:"IsAuthenticated,omitempty"`
	}
	var data status
	_ = decoder.Decode(&data)
	return data.IsAuthenticated
}
