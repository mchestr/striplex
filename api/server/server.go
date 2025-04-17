package server

import (
	"context"
	"net/http"
	"plefi/api/config"
	"plefi/api/services"

	"github.com/gin-gonic/gin"
)

// Server holds the HTTP server and router instances
type Server struct {
	router *gin.Engine
	server *http.Server
}

// Init initializes the server without starting it
func Init(svcs *services.Services, client *http.Client) (*Server, error) {
	r := NewRouter(svcs, client)
	r.SetTrustedProxies(config.C.Server.TrustedProxies)

	srv := &Server{
		router: r,
		server: &http.Server{
			Addr:    config.C.Server.Address,
			Handler: r,
		},
	}

	return srv, nil
}

// Start begins serving HTTP requests
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
