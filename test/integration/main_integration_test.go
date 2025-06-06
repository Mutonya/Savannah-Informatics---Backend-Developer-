package integration

import (
	"context"
	"net/http"
	_ "net/http/httptest"
	"testing"
	"time"

	"github.com/Mutonya/Savanah/internal/config"
	_ "github.com/Mutonya/Savanah/pkg/oauth2"
	"github.com/stretchr/testify/require"
)

func TestMainIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg := config.LoadTestConfig()

	t.Run("complete application startup and shutdown", func(t *testing.T) {
		// Setup context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Channel to signal when server is ready
		ready := make(chan bool)

		// Start the application in a goroutine
		go func() {
			// Modify main function to use test config
			runTestApplication(cfg, ready)
		}()

		// Wait for server to be ready or timeout
		select {
		case <-ready:
			// Server is ready, run tests
			testServerEndpoints(t, cfg)
		case <-ctx.Done():
			t.Fatal("server did not start within timeout period")
		}
	})
}

func runTestApplication(cfg *config.TestConfig, ready chan<- bool) {

	// Initialize components using test config...

	// Signal that server is ready
	ready <- true

	// Wait for shutdown signal...
}

func testServerEndpoints(t *testing.T, cfg *config.TestConfig) {
	// Test health check endpoint
	resp, err := http.Get("http://localhost:" + cfg.Server.Port + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Add more endpoint tests as needed
}
