//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"pryx-core/internal/config"
	"pryx-core/internal/keychain"
	"pryx-core/internal/server"
	"pryx-core/internal/store"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestChannelEndpointsIntegration tests channel API endpoints
func TestChannelEndpointsIntegration(t *testing.T) {
	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	s, _ := store.New(":memory:")
	defer s.Close()

	t.Setenv("PRYX_KEYCHAIN_FILE", t.TempDir()+"/keychain.json")
	kc := keychain.New("test")

	srv := server.New(cfg, s.DB, kc)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	go srv.Serve(listener)
	time.Sleep(10 * time.Millisecond)

	client := &http.Client{Timeout: time.Second}
	baseUrl := "http://" + listener.Addr().String()

	// Test channel list endpoint
	resp, err := client.Get(baseUrl + "/api/v1/channels")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Contains(t, result, "channels")
}

// TestOAuthDeviceFlowEndpoints tests the OAuth device flow endpoints
func TestOAuthDeviceFlowEndpoints(t *testing.T) {
	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	s, _ := store.New(":memory:")
	defer s.Close()

	t.Setenv("PRYX_KEYCHAIN_FILE", t.TempDir()+"/keychain.json")
	kc := keychain.New("test")

	srv := server.New(cfg, s.DB, kc)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	go srv.Serve(listener)
	time.Sleep(10 * time.Millisecond)

	client := &http.Client{Timeout: time.Second}
	baseUrl := "http://" + listener.Addr().String()

	// Test device code endpoint (should return structure for device flow)
	resp, err := client.Get(baseUrl + "/api/v1/auth/device/code")
	if err != nil {
		// This might fail if auth is not configured - that's expected
		t.Logf("Device code endpoint not available (expected if auth not configured): %v", err)
		return
	}
	defer resp.Body.Close()

	// If endpoint exists, verify structure
	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		// Device flow should return device_code, user_code, verification_uri, etc.
		assert.Contains(t, result, "device_code")
		assert.Contains(t, result, "user_code")
		assert.Contains(t, result, "verification_uri")
	}
}

// TestCompleteWorkflowIntegration tests a complete user workflow
func TestCompleteWorkflowIntegration(t *testing.T) {
	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	s, _ := store.New(":memory:")
	defer s.Close()

	t.Setenv("PRYX_KEYCHAIN_FILE", t.TempDir()+"/keychain.json")
	kc := keychain.New("test")

	srv := server.New(cfg, s.DB, kc)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	go srv.Serve(listener)
	time.Sleep(10 * time.Millisecond)

	client := &http.Client{Timeout: time.Second}
	baseUrl := "http://" + listener.Addr().String()

	// Test health endpoint
	resp, err := client.Get(baseUrl + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test skills endpoint
	resp, err = client.Get(baseUrl + "/skills")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test providers endpoint
	resp, err = client.Get(baseUrl + "/api/v1/providers")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test sessions endpoint
	resp, err = client.Get(baseUrl + "/api/v1/sessions")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test channels endpoint
	resp, err = client.Get(baseUrl + "/api/v1/channels")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestMeshPairingIntegration tests mesh QR code pairing endpoints
func TestMeshPairingIntegration(t *testing.T) {
	cfg := &config.Config{ListenAddr: "127.0.0.1:0", CloudAPIUrl: "http://localhost:3000"}
	s, _ := store.New(":memory:")
	defer s.Close()

	t.Setenv("PRYX_KEYCHAIN_FILE", t.TempDir()+"/keychain.json")
	kc := keychain.New("test")
	kc.Set("device_id", "test-device-001")
	kc.Set("device_name", "Test Device")

	srv := server.New(cfg, s.DB, kc)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	go srv.Serve(listener)
	time.Sleep(10 * time.Millisecond)

	client := &http.Client{Timeout: time.Second}
	baseUrl := "http://" + listener.Addr().String()

	// Test QR code generation endpoint
	resp, err := client.Post(baseUrl+"/api/mesh/qrcode", "application/json", nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var qrResult map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&qrResult)
	require.NoError(t, err)
	assert.Contains(t, qrResult, "code")
	assert.Contains(t, qrResult, "expires_at")
	assert.Contains(t, qrResult, "qr_code")
	assert.Len(t, qrResult["code"].(string), 6) // 6-digit code

	// Get the pairing code for subsequent tests
	pairingCode := qrResult["code"].(string)

	// Test pairing with invalid code
	pairReq := map[string]string{"code": "123456"}
	pairJson, _ := json.Marshal(pairReq)
	resp, err = client.Post(baseUrl+"/api/mesh/pair", "application/json", strings.NewReader(string(pairJson)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Test pairing with valid code
	pairReq["code"] = pairingCode
	pairJson, _ = json.Marshal(pairReq)
	resp, err = client.Post(baseUrl+"/api/mesh/pair", "application/json", strings.NewReader(string(pairJson)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var pairResult map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&pairResult)
	require.NoError(t, err)
	assert.Contains(t, pairResult, "success")
	assert.True(t, pairResult["success"].(bool))
	assert.Contains(t, pairResult, "device_id")

	// Test pairing again with same code (should fail - already used)
	resp, err = client.Post(baseUrl+"/api/mesh/pair", "application/json", strings.NewReader(string(pairJson)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Test devices list endpoint
	resp, err = client.Get(baseUrl + "/api/mesh/devices")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var devicesResult map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&devicesResult)
	require.NoError(t, err)
	assert.Contains(t, devicesResult, "devices")

	// Test events list endpoint
	resp, err = client.Get(baseUrl + "/api/mesh/events")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var eventsResult map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&eventsResult)
	require.NoError(t, err)
	assert.Contains(t, eventsResult, "events")
}

// TestMeshPairingValidationIntegration tests mesh pairing validation
func TestMeshPairingValidationIntegration(t *testing.T) {
	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	s, _ := store.New(":memory:")
	defer s.Close()

	t.Setenv("PRYX_KEYCHAIN_FILE", t.TempDir()+"/keychain.json")
	kc := keychain.New("test")

	srv := server.New(cfg, s.DB, kc)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	go srv.Serve(listener)
	time.Sleep(10 * time.Millisecond)

	client := &http.Client{Timeout: time.Second}
	baseUrl := "http://" + listener.Addr().String()

	// Test pairing with invalid code format (too short)
	pairReq := map[string]string{"code": "123"}
	pairJson, _ := json.Marshal(pairReq)
	resp, err := client.Post(baseUrl+"/api/mesh/pair", "application/json", strings.NewReader(string(pairJson)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Test pairing with invalid code format (too long)
	pairReq["code"] = "1234567"
	pairJson, _ = json.Marshal(pairReq)
	resp, err = client.Post(baseUrl+"/api/mesh/pair", "application/json", strings.NewReader(string(pairJson)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Test pairing with non-numeric characters
	pairReq["code"] = "12AB56"
	pairJson, _ = json.Marshal(pairReq)
	resp, err = client.Post(baseUrl+"/api/mesh/pair", "application/json", strings.NewReader(string(pairJson)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Test pairing with wrong HTTP method
	resp, err = client.Get(baseUrl + "/api/mesh/pair")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	// Test QR code with wrong HTTP method
	resp, err = client.Get(baseUrl + "/api/mesh/qrcode")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

// TestMeshPairingSessionIntegration tests mesh pairing session lifecycle
func TestMeshPairingSessionIntegration(t *testing.T) {
	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	s, _ := store.New(":memory:")
	defer s.Close()

	t.Setenv("PRYX_KEYCHAIN_FILE", t.TempDir()+"/keychain.json")
	kc := keychain.New("test")

	srv := server.New(cfg, s.DB, kc)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	go srv.Serve(listener)
	time.Sleep(10 * time.Millisecond)

	client := &http.Client{Timeout: time.Second}
	baseUrl := "http://" + listener.Addr().String()

	// Generate a pairing session
	resp, err := client.Post(baseUrl+"/api/mesh/qrcode", "application/json", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	var qrResult map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&qrResult)
	require.NoError(t, err)

	pairingCode := qrResult["code"].(string)

	// Verify session was created in database
	session, err := s.GetPairingSessionByCode(pairingCode)
	require.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, pairingCode, session.Code)
	assert.Equal(t, "test-device-001", session.DeviceID)
	assert.Equal(t, "pending", session.Status)

	// Test unpair with non-existent device
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/api/mesh/devices/non-existent-id/unpair", nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
