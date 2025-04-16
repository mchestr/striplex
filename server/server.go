package server

import (
	"context"
	"net/http"
	"plefi/config"
	"plefi/services"

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
	r.SetTrustedProxies(config.Config.GetStringSlice("server.trusted_proxies"))

	srv := &Server{
		router: r,
		server: &http.Server{
			Addr:    config.Config.GetString("server.address"),
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
