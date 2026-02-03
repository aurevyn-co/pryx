package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pryx-core/internal/config"
	"pryx-core/internal/store"

	"github.com/stretchr/testify/assert"
)

// TestCloudLoginWithRealAPI tests cloud login endpoints structure
func TestCloudLoginWithRealAPI(t *testing.T) {
	cfg := &config.Config{
		ListenAddr:  ":0",
		CloudAPIUrl: "https://pryx.dev/api",
	}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	server := New(cfg, s.DB, kc)

	// Test login start endpoint structure
	t.Run("login_start_endpoint", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/cloud/login/start", nil)
		rec := httptest.NewRecorder()
		server.router.ServeHTTP(rec, req)

		// Should return OK or fail gracefully (without pryx.dev API)
		assert.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusBadGateway)
	})

	// Test login poll endpoint structure
	t.Run("login_poll_endpoint", func(t *testing.T) {
		reqBody := `{"device_code":"test-code","interval":1,"expires_in":60}`
		req := httptest.NewRequest("POST", "/api/v1/cloud/login/poll", strings.NewReader(reqBody))
		rec := httptest.NewRecorder()
		server.router.ServeHTTP(rec, req)

		// Should return OK or fail gracefully
		assert.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusBadGateway || rec.Code == http.StatusRequestTimeout)
	})

	// Test cloud status endpoint
	t.Run("cloud_status_endpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/cloud/status", nil)
		rec := httptest.NewRecorder()
		server.router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

// TestCloudLoginStartEndpoint tests login start endpoint
func TestCloudLoginStartEndpoint(t *testing.T) {
	cfg := &config.Config{
		ListenAddr:  ":0",
		CloudAPIUrl: "https://pryx.dev/api",
	}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	server := New(cfg, s.DB, kc)

	req := httptest.NewRequest("POST", "/api/v1/cloud/login/start", nil)
	rec := httptest.NewRecorder()
	server.router.ServeHTTP(rec, req)

	// Should fail gracefully without pryx.dev API
	assert.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusBadGateway)
}

// TestCloudLoginPollEndpoint tests login poll endpoint
func TestCloudLoginPollEndpoint(t *testing.T) {
	cfg := &config.Config{
		ListenAddr:  ":0",
		CloudAPIUrl: "https://pryx.dev/api",
	}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	server := New(cfg, s.DB, kc)

	reqBody := `{"device_code":"test","interval":1,"expires_in":60}`
	req := httptest.NewRequest("POST", "/api/v1/cloud/login/poll", strings.NewReader(reqBody))
	rec := httptest.NewRecorder()
	server.router.ServeHTTP(rec, req)

	// Should fail gracefully without pryx.dev API
	assert.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusBadGateway || rec.Code == http.StatusRequestTimeout)
}

// TestCloudStatusEndpoint tests cloud status endpoint
func TestCloudStatusEndpoint(t *testing.T) {
	cfg := &config.Config{
		ListenAddr:  ":0",
		CloudAPIUrl: "https://pryx.dev/api",
	}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	server := New(cfg, s.DB, kc)

	req := httptest.NewRequest("GET", "/api/v1/cloud/status", nil)
	rec := httptest.NewRecorder()
	server.router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
