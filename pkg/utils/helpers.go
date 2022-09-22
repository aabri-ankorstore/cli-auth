package utils

import (
	"fmt"
	"github.com/ankorstore/ankorstore-cli-core/pkg/util"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"path/filepath"
)

type AuthStatus struct {
	IsConnected bool `json:"is_connected"`
}

func IsAuthenticated(r *http.Request) bool {
	session, err := SessionStore.Get(r, CookieName)
	if err != nil || session.Values["access_token"] == nil || session.Values["access_token"] == "" {
		log.Debug().Err(err)
		return false
	}
	return true
}

func CreateTmpFile() (string, error) {
	dirs := util.NewDirs()
	f, err := os.CreateTemp(fmt.Sprintf("%s/%s", dirs.GetPluginsDir(), PluginPath), "*-auth.lock")
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}

func IsAuthenticatedOffline() bool {
	dirs := util.NewDirs()
	pattern := "*-auth.lock"
	file := fmt.Sprintf("%s/%s/%s", dirs.GetPluginsDir(), PluginPath, pattern)
	matches, err := filepath.Glob(file)
	if err != nil {
		return false
	}
	if len(matches) > 0 {
		return true
	}
	return false
}

func RemoveAuth() {
	_ = os.Remove(LockFile)
}
