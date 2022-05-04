package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/leads-su/logger"
)

// healthRouteResponse represents structure of health check response object
type healthRouteResponse struct {
	Status    bool   `json:"status"`
	Timestamp string `json:"timestamp"`
	Timezone  string `json:"timezone"`
}

// registerHealthRoute registers health check route for Consul agent
func (server *Server) registerHealthRoute() {
	logger.Tracef("consul:http", "registered health check route")
	http.HandleFunc("/health", server.handleHealthRequest)
}

// handleHealthRequest handles request to '/health' endpoint
func (server *Server) handleHealthRequest(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(healthRouteResponse{
		Status:    true,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Timezone:  "UTC",
	})
}
