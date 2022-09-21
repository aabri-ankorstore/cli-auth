package utils

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

func IsAuthenticated(r *http.Request) bool {
	session, err := SessionStore.Get(r, CookieName)
	if err != nil || session.Values["access_token"] == nil || session.Values["access_token"] == "" {
		log.Debug().Err(err)
		return false
	}
	return true
}

func CreateTmpFile() (*os.File, error) {
	f, err := os.CreateTemp("", "auth")
	if err != nil {
		return nil, err
	}
	//defer os.Remove(f.Name())
	fmt.Println(f.Name())
	return f, nil
}

func IsAuthenticatedOffline() bool {
	fmt.Println(LockFile.Name())
	if _, err := os.Stat(LockFile.Name()); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
