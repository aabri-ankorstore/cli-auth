package utils

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
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
	req, reqerr := http.NewRequest("POST", "http://localhost:8080/is-authenticated", bytes.NewBuffer(nil))
	if reqerr != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	fmt.Println(body)
	return true
}
