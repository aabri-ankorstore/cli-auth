package handlers

import (
	"github.com/aabri-ankorstore/cli-auth/utils"
	"net/http"
)

func (h *Auth) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	type customData struct {
		Profile         map[string]string
		IsAuthenticated bool
	}
	profile, _ := h.manager.GetProfile(r)
	data := customData{
		Profile:         profile,
		IsAuthenticated: utils.IsAuthenticated(r),
	}
	_ = utils.Tpl.ExecuteTemplate(w, "profile.gohtml", data)
}
