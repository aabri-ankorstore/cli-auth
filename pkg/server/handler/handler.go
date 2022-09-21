package handler

import (
	"github.com/aabri-ankorstore/cli-auth/pkg/server/util/portforward"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	Mux   *mux.Router
	Ports map[string]*Forward
}

type Forward struct {
	portForwarder     *portforward.PortForwarder
	portForwarderStop chan struct{}
	portForwarderPort int
	podUUID           string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	/*w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == "OPTIONS" {
		return
	}*/

	if r.Method != "GET" && r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	h.Mux.ServeHTTP(w, r)
}
