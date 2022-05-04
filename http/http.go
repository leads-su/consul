package http

import (
	"fmt"
	"net/http"

	"github.com/leads-su/logger"
)

// Server represents structure of HTTP server
type Server struct {
	Enabled bool
	Port    uint
}

// NewServer creates new instance of HTTP server
func NewServer(port uint, enabled bool) *Server {
	server := &Server{
		Enabled: enabled,
		Port:    port,
	}
	server.registerRoutes()
	return server
}

// registerRoutes registers all routes needed for the package
func (server *Server) registerRoutes() {
	server.registerHealthRoute()

	server.enableBundledServer()
}

// enableBundledServer enables bundled server if it is enabled in the configuration
func (server *Server) enableBundledServer() {
	if server.Enabled {
		endpoint := fmt.Sprintf(":%d", server.Port)
		logger.Infof("consul:http", "starting http server at %s", endpoint)
		go func() {
			err := http.ListenAndServe(endpoint, nil)
			if err != nil {
				logger.Fatalf("consul:http", "failed to create http server - %s", err.Error())
			}
		}()
	}
}
