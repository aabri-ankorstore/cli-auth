package handlers

import (
	"github.com/aabri-ankorstore/cli-auth/pkg/utils"
	"github.com/rs/zerolog/log"
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
		IsAuthenticated: utils.IsAuthenticated(r),
	}
	err := utils.View.ExecuteTemplate(w, "home.gohtml", data)
	if err != nil {
		log.Trace().Msg(err.Error())
	}
}
