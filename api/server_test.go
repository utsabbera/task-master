package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	t.Run("should return configured HTTP server", func(t *testing.T) {
		cfg := ServerConfig{Addr: ":9090"}
		server := NewServer(cfg)

		assert.NotNil(t, server)
		assert.Equal(t, ":9090", server.Addr)
		assert.NotNil(t, server.Handler)
	})

	t.Run("should use default Addr when empty", func(t *testing.T) {
		cfg := ServerConfig{Addr: ""}
		server := NewServer(cfg)

		assert.Equal(t, ":8080", server.Addr)
	})
}
