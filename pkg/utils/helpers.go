package utils

import (
	"errors"
	"fmt"
	"github.com/ankorstore/ankorstore-cli-core/pkg/util"
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
	dirs := util.NewDirs()
	f, err := os.CreateTemp(fmt.Sprintf("%s/%s", dirs.GetPluginsDir(), PluginPath), "*-auth.lock")
	if err != nil {
		return nil, err
	}
	return f, nil
}

func IsAuthenticatedOffline() bool {
	dirs := util.NewDirs()
	file := fmt.Sprintf("%s/%s/%s", dirs.GetPluginsDir(), PluginPath, LockFile)
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func RemoveAuth() {
	defer os.Remove(LockFile)
}
