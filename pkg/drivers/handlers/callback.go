package handlers

import (
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/pkg/filesystem"
	"github.com/aabri-ankorstore/cli-auth/pkg/utils"
	"github.com/ankorstore/ankorstore-cli-core/pkg/util"
	"net/http"
)

func (h *Auth) CallBackHandler(w http.ResponseWriter, r *http.Request) {
	e, _ := h.manager.ExchangeCode(w, r)
	if e.Error != "" {
		fmt.Println(e.Error)
		fmt.Println(e.ErrorDescription)
		return
	}
	// save access token
	dir := util.NewDirs()
	f := filesystem.NewFilesystem(dir.GetPluginsDir())

	file, err := f.CreateTmpFile()
	f.CheckError(err)
	utils.LockFile = file
	http.Redirect(w, r, "/", http.StatusFound)
}
