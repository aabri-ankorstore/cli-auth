package handlers

import (
	"github.com/aabri-ankorstore/cli-auth/drivers"
	"github.com/aabri-ankorstore/cli-auth/server/handler"
	"net/http"
)

type Auth struct {
	*handler.Handler
	manager drivers.Manager
}

func NewHandler(handler *handler.Handler, manager drivers.Manager) (*handler.Handler, error) { // Get kube config
	a := Auth{
		manager: manager,
	}
	// Serving assets
	fs := http.FileServer(http.Dir("./ui/static"))
	handler.Mux.PathPrefix("/css/").Handler(fs)
	handler.Mux.PathPrefix("/fonts/").Handler(fs)
	handler.Mux.PathPrefix("/images/").Handler(fs)
	handler.Mux.PathPrefix("/js/").Handler(fs)
	handler.Mux.PathPrefix("/plugins/").Handler(fs)
	//#################

	handler.Mux.HandleFunc("/", a.HomeHandler).Methods("GET")
	handler.Mux.HandleFunc("/login", a.LoginHandler).Methods("GET")
	handler.Mux.HandleFunc("/authorization-code/callback", a.CallBackHandler).Methods("GET")
	handler.Mux.HandleFunc("/profile", a.ProfileHandler).Methods("GET")
	handler.Mux.HandleFunc("/logout", a.LogoutHandler).Methods("POST")
	return handler, nil
}
