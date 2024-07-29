package server

import (
	"testing"
	"tezos-delegation-service/config"
	"tezos-delegation-service/internal/api/handlers"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	// Mock configuration and handler
	cfg := &config.Config{
		Port: "8080",
	}
	handler := &handlers.Handler{} // Assume this is your handler struct
	cfg.RateLimit = 10

	// Create the server
	srv := NewServer(cfg, handler)

	// Validate server configuration
	assert.Equal(t, ":8080", srv.Addr)
	assert.Equal(t, 10*time.Second, srv.ReadTimeout)
	assert.Equal(t, 10*time.Second, srv.WriteTimeout)
	assert.Equal(t, 1<<20, srv.MaxHeaderBytes)
}
