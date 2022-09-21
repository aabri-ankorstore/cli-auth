package handlers

import (
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/pkg/utils"
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
	_, err := utils.CreateTmpFile()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
