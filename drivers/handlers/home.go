package handlers

import (
	utils2 "github.com/aabri-ankorstore/cli-auth/utils"
	"net/http"
)

func (h *Auth) HomeHandler(w http.ResponseWriter, r *http.Request) {
	type customData struct {
		Profile         map[string]string
		IsAuthenticated bool
	}

	profile, _ := h.manager.GetProfile(r)
	data := customData{
		Profile:         profile,
		IsAuthenticated: utils2.IsAuthenticated(r),
	}
	_ = utils2.View.ExecuteTemplate(w, "home.gohtml", data)
}
