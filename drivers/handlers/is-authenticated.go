package handlers

import (
	"encoding/json"
	utils2 "github.com/aabri-ankorstore/cli-auth/utils"
	"github.com/rs/zerolog/log"
	"net/http"
)

type status struct {
	IsAuthenticated bool
}

func (h *Auth) IsAuthenticated(w http.ResponseWriter, r *http.Request) {
	data := status{
		IsAuthenticated: utils2.IsAuthenticated(r),
	}
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		log.Debug().Msg("Unable to encode JSON")
	}

	r.Header.Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
