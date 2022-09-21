package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

func (h *Auth) CallBackHandler(w http.ResponseWriter, r *http.Request) {
	e, _ := h.manager.ExchangeCode(w, r)
	if e.Error != "" {
		fmt.Println(e.Error)
		fmt.Println(e.ErrorDescription)
		return
	}
	// save access token

	http.Redirect(w, r, "/", http.StatusFound)
}

func empty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
