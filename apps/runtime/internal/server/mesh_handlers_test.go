package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pryx-core/internal/config"
	"pryx-core/internal/store"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMeshQRCodeGeneration tests QR code generation endpoint
func TestMeshQRCodeGeneration(t *testing.T) {
	cfg := &config.Config{ListenAddr: ":0"}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	server := New(cfg, s.DB, kc)

	req := httptest.NewRequest("POST", "/api/mesh/qrcode", nil)
	rec := httptest.NewRecorder()

	server.handleMeshQRCode(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response handleMeshQRCodeResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// QR code should be base64 encoded PNG or JSON fallback
	assert.NotEmpty(t, response.QRCode)
	// Should have 6-digit code
	assert.Len(t, response.Code, 6)
	for _, c := range response.Code {
		assert.True(t, c >= '0' && c <= '9')
	}
	assert.NotEmpty(t, response.ExpiresAt)
}

// TestMeshPairWithInvalidCode tests pairing with invalid code format
func TestMeshPairWithInvalidCode(t *testing.T) {
	cfg := &config.Config{ListenAddr: ":0"}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	server := New(cfg, s.DB, kc)

	tests := []struct {
		name           string
		code           string
		expectedStatus int
		expectedMsg    string
	}{
		{"empty code", "", http.StatusBadRequest, "pairing code must be exactly 6 digits"},
		{"too short", "123", http.StatusBadRequest, "pairing code must be exactly 6 digits"},
		{"too long", "1234567", http.StatusBadRequest, "pairing code must be exactly 6 digits"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := `{"code":"` + tt.code + `"}`
			req := httptest.NewRequest("POST", "/api/mesh/pair", strings.NewReader(reqBody))
			rec := httptest.NewRecorder()

			server.handleMeshPair(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			var response handleMeshPairResponse
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.False(t, response.Success)
			assert.Contains(t, response.Message, "pairing code must be exactly 6 digits")
		})
	}
}

// TestMeshPairNotFound tests pairing with non-existent code
func TestMeshPairNotFound(t *testing.T) {
	cfg := &config.Config{ListenAddr: ":0"}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	server := New(cfg, s.DB, kc)

	reqBody := `{"code":"999999"}`
	req := httptest.NewRequest("POST", "/api/mesh/pair", strings.NewReader(reqBody))
	rec := httptest.NewRecorder()

	server.handleMeshPair(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var response handleMeshPairResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Message, "invalid or expired pairing code")
}

// TestMeshDevicesList tests listing paired devices
func TestMeshDevicesList(t *testing.T) {
	cfg := &config.Config{ListenAddr: ":0"}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	server := New(cfg, s.DB, kc)

	req := httptest.NewRequest("GET", "/api/mesh/devices", nil)
	rec := httptest.NewRecorder()

	server.handleMeshDevicesList(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "devices")
}

// TestMeshEventsList tests listing mesh events
func TestMeshEventsList(t *testing.T) {
	cfg := &config.Config{ListenAddr: ":0"}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	server := New(cfg, s.DB, kc)

	req := httptest.NewRequest("GET", "/api/mesh/events", nil)
	rec := httptest.NewRecorder()

	server.handleMeshEventsList(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "events")
}

// TestMeshPairInvalidMethod tests that only POST is allowed
func TestMeshPairInvalidMethod(t *testing.T) {
	cfg := &config.Config{ListenAddr: ":0"}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	server := New(cfg, s.DB, kc)

	req := httptest.NewRequest("GET", "/api/mesh/pair", nil)
	rec := httptest.NewRecorder()

	server.handleMeshPair(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

// TestMeshPairInvalidJSON tests invalid JSON body
func TestMeshPairInvalidJSON(t *testing.T) {
	cfg := &config.Config{ListenAddr: ":0"}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	server := New(cfg, s.DB, kc)

	req := httptest.NewRequest("POST", "/api/mesh/pair", strings.NewReader("not json"))
	rec := httptest.NewRecorder()

	server.handleMeshPair(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestGeneratePairingCode tests pairing code generation
func TestGeneratePairingCode(t *testing.T) {
	// Test multiple generations for uniqueness
	codes := make(map[string]bool)
	for i := 0; i < 100; i++ {
		code := generatePairingCode()
		assert.Len(t, code, 6)
		// All should be digits
		for _, c := range code {
			assert.True(t, c >= '0' && c <= '9')
		}
		assert.False(t, codes[code], "duplicate code generated")
		codes[code] = true
	}
}

// TestGenerateDeviceID tests device ID generation
func TestGenerateDeviceID(t *testing.T) {
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := generateDeviceID()
		assert.True(t, strings.HasPrefix(id, "pryx-"))
		assert.Len(t, id, 13) // "pryx-" (5) + 8 chars
		assert.False(t, ids[id], "duplicate device ID generated")
		ids[id] = true
	}
}
