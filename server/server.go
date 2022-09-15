package server

import (
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/server/handler"
	"github.com/aabri-ankorstore/cli-auth/server/util/port"
	"github.com/go-errors/errors"
	"net/http"
	"strconv"
)

// Server is listens on a given port for the ui functionality
type Server struct {
	Server *http.Server
}

// DefaultPort is the default port the ui server will listen to
const DefaultPort = 8090

func NewServer(host string, forcePort *int, h *handler.Handler) (*Server, error) {
	// Find an open port
	usePort := DefaultPort
	if forcePort != nil {
		usePort = *forcePort
		if host == "localhost" {
			available, err := port.IsAvailable(fmt.Sprintf(":%d", usePort))
			if !available {
				return nil, errors.Errorf("Port %d already in use: %v", usePort, err)
			}
		}
	} else {
		if host == "localhost" {
			for i := 0; i < 20; i++ {
				available, _ := port.IsAvailable(fmt.Sprintf(":%d", usePort))
				if available {
					break
				}
				usePort++
			}
		}
	}
	return &Server{
		Server: &http.Server{
			Addr:    host + ":" + strconv.Itoa(usePort),
			Handler: h,
			//ReadTimeout:  5 * time.Second,
			//WriteTimeout: 10 * time.Second,
			//IdleTimeout:  60 * time.Second,
		},
	}, nil
}

// ListenAndServe implements interface
func (s *Server) ListenAndServe() error {
	return s.Server.ListenAndServe()
}
