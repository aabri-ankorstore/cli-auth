package handlers

import (
	"github.com/aabri-ankorstore/cli-auth/drivers"
	"github.com/aabri-ankorstore/cli-auth/server/handler"
	"github.com/aabri-ankorstore/cli-auth/utils"
	"io/fs"
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
	fSys, err := fs.Sub(utils.Tpl, "ui/static")
	if err != nil {
		panic(err)
	}
	staticFile := http.FileServer(http.FS(fSys))
	handler.Mux.PathPrefix("/css/").Handler(staticFile)
	handler.Mux.PathPrefix("/fonts/").Handler(staticFile)
	handler.Mux.PathPrefix("/images/").Handler(staticFile)
	handler.Mux.PathPrefix("/js/").Handler(staticFile)
	handler.Mux.PathPrefix("/plugins/").Handler(staticFile)
	//#################

	handler.Mux.HandleFunc("/", a.HomeHandler).Methods("GET")
	handler.Mux.HandleFunc("/is-authenticated", a.IsAuthenticated).Methods("GET")
	handler.Mux.HandleFunc("/login", a.LoginHandler).Methods("GET")
	handler.Mux.HandleFunc("/authorization-code/callback", a.CallBackHandler).Methods("GET")
	handler.Mux.HandleFunc("/logout", a.LogoutHandler).Methods("POST")
	return handler, nil
}
