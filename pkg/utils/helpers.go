package utils

import (
	"fmt"
	"github.com/ankorstore/ankorstore-cli-core/pkg/util"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"path/filepath"
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
	pattern := "*-auth.lock"
	file := fmt.Sprintf("%s/%s/%s", dirs.GetPluginsDir(), PluginPath, pattern)
	matches, err := filepath.Glob(file)
	if err != nil {
		return false
	}
	for _, match := range matches {
		fmt.Println(fmt.Sprintf("Match :%s", match))
		fmt.Println(fmt.Sprintf("Lock File :%s", LockFile))
		if match == LockFile {
			return true
		}
	}
	return false
}

func RemoveAuth() {
	_ = os.Remove(LockFile)
}
