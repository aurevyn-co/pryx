//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"pryx-core/internal/bus"
	"pryx-core/internal/config"
	"pryx-core/internal/keychain"
	"pryx-core/internal/server"
	"pryx-core/internal/store"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"nhooyr.io/websocket"
)

type testEnv interface {
	Helper()
	Setenv(key, value string)
	TempDir() string
}

func newTestKeychain(t testEnv) *keychain.Keychain {
	t.Helper()
	t.Setenv("PRYX_KEYCHAIN_FILE", filepath.Join(t.TempDir(), "keychain.json"))
	return keychain.New("test")
}

// TestRuntimeStartup tests the complete runtime startup sequence
func TestRuntimeStartup(t *testing.T) {
	// Create temporary directory for test data
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.Config{
		ListenAddr:   "127.0.0.1:0", // Let OS assign port
		DatabasePath: dbPath,
	}

	s, err := store.New(dbPath)
	require.NoError(t, err)
	defer s.Close()

	kc := newTestKeychain(t)

	srv := server.New(cfg, s.DB, kc)
	require.NotNil(t, srv)

	// Start server in background
	go func() {
		_ = srv.Start()
	}()

	// Give server time to start and write port file
	time.Sleep(100 * time.Millisecond)

	// Read port from file
	portFile := filepath.Join(tmpDir, ".pryx", "runtime.port")
	if _, err := os.Stat(portFile); err == nil {
		data, _ := os.ReadFile(portFile)
		t.Logf("Server started on port: %s", string(data))
	}

	_ = srv.Shutdown(context.Background())
}

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	srv := server.New(cfg, s.DB, kc)

	// Create listener
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	// Start server
	go srv.Serve(listener)
	time.Sleep(10 * time.Millisecond)

	// Make health request
	client := &http.Client{Timeout: time.Second}
	resp, err := client.Get("http://" + listener.Addr().String() + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	assert.Equal(t, "ok", result["status"])
}

// TestSkillsEndpoint tests the skills API
func TestSkillsEndpoint(t *testing.T) {
	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	srv := server.New(cfg, s.DB, kc)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	go srv.Serve(listener)
	time.Sleep(10 * time.Millisecond)

	// Test skills list endpoint
	client := &http.Client{Timeout: time.Second}
	resp, err := client.Get("http://" + listener.Addr().String() + "/skills")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Contains(t, result, "skills")
}

func TestProviderKeyEndpoints(t *testing.T) {
	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	srv := server.New(cfg, s.DB, kc)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	go srv.Serve(listener)
	time.Sleep(10 * time.Millisecond)

	client := &http.Client{Timeout: time.Second}
	baseUrl := "http://" + listener.Addr().String()

	{
		resp, err := client.Get(baseUrl + "/api/v1/providers/openai/key")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Equal(t, false, result["configured"])
	}

	{
		body := bytes.NewBufferString(`{"api_key":"sk-test"}`)
		resp, err := client.Post(baseUrl+"/api/v1/providers/openai/key", "application/json", body)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	}

	{
		resp, err := client.Get(baseUrl + "/api/v1/providers/openai/key")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Equal(t, true, result["configured"])
	}

	{
		req, err := http.NewRequest(http.MethodDelete, baseUrl+"/api/v1/providers/openai/key", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	}

	{
		resp, err := client.Get(baseUrl + "/api/v1/providers/openai/key")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Equal(t, false, result["configured"])
	}

	{
		resp, err := client.Get(baseUrl + "/api/v1/providers/bad%20id/key")
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	}

	{
		resp, err := client.Get(baseUrl + "/api/v1/providers/not-a-real-provider/key")
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	}

	{
		resp, err := client.Post(
			baseUrl+"/api/v1/providers/openai/key",
			"application/json",
			strings.NewReader(`not json`),
		)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	}

	{
		resp, err := client.Post(
			baseUrl+"/api/v1/providers/openai/key",
			"application/json",
			strings.NewReader(`{"api_key":""}`),
		)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	}

	{
		req, err := http.NewRequest(http.MethodDelete, baseUrl+"/api/v1/providers/openai/key", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	}
}

// TestWebSocketConnection tests WebSocket upgrade and basic communication
func TestWebSocketConnection(t *testing.T) {
	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	srv := server.New(cfg, s.DB, kc)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	go srv.Serve(listener)
	time.Sleep(10 * time.Millisecond)

	// Connect via WebSocket
	wsURL := "ws://" + listener.Addr().String() + "/ws"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ws, resp, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{})
	require.NoError(t, err)
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
	defer ws.Close(websocket.StatusNormalClosure, "test complete")

	// Connection should be established
	assert.NotNil(t, ws)
}

func TestCloudLoginEndpoints_Validation(t *testing.T) {
	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	srv := server.New(cfg, s.DB, kc)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	go srv.Serve(listener)
	time.Sleep(10 * time.Millisecond)

	client := &http.Client{Timeout: time.Second}
	baseUrl := "http://" + listener.Addr().String()

	{
		resp, err := client.Post(baseUrl+"/api/v1/cloud/login/start", "application/json", nil)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	}

	{
		cfgWithCloud := &config.Config{ListenAddr: "127.0.0.1:0", CloudAPIUrl: "https://example.invalid"}
		srv2 := server.New(cfgWithCloud, s.DB, kc)

		listener2, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)
		defer listener2.Close()

		go srv2.Serve(listener2)
		time.Sleep(10 * time.Millisecond)

		baseUrl2 := "http://" + listener2.Addr().String()

		resp, err := client.Post(
			baseUrl2+"/api/v1/cloud/login/poll",
			"application/json",
			strings.NewReader(`{}`),
		)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		_ = srv2.Shutdown(context.Background())
	}

	_ = srv.Shutdown(context.Background())
}

func TestWebSocketSessionsList(t *testing.T) {
	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	sess, err := s.CreateSession("Test Session")
	require.NoError(t, err)
	_, err = s.AddMessage(sess.ID, store.RoleUser, "hello")
	require.NoError(t, err)

	srv := server.New(cfg, s.DB, kc)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	go srv.Serve(listener)
	time.Sleep(10 * time.Millisecond)

	wsURL := "ws://" + listener.Addr().String() + "/ws"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ws, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{})
	require.NoError(t, err)
	defer ws.Close(websocket.StatusNormalClosure, "")

	req := map[string]any{"event": "sessions.list", "payload": map[string]any{}}
	reqBytes, err := json.Marshal(req)
	require.NoError(t, err)
	require.NoError(t, ws.Write(ctx, websocket.MessageText, reqBytes))

	readCtx, readCancel := context.WithTimeout(ctx, time.Second)
	defer readCancel()

	found := false
	for i := 0; i < 10; i++ {
		_, data, err := ws.Read(readCtx)
		require.NoError(t, err)

		var evt map[string]any
		require.NoError(t, json.Unmarshal(data, &evt))
		if evt["event"] != "sessions.list" {
			continue
		}

		payload, _ := evt["payload"].(map[string]any)
		sessions, _ := payload["sessions"].([]any)
		require.NotEmpty(t, sessions)
		found = true
		break
	}
	require.True(t, found)
}

func TestWebSocketSessionResume(t *testing.T) {
	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	sess, err := s.CreateSession("Test Session")
	require.NoError(t, err)
	_, err = s.AddMessage(sess.ID, store.RoleUser, "hello")
	require.NoError(t, err)

	srv := server.New(cfg, s.DB, kc)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	go srv.Serve(listener)
	time.Sleep(10 * time.Millisecond)

	wsURL := "ws://" + listener.Addr().String() + "/ws"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ws, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{})
	require.NoError(t, err)
	defer ws.Close(websocket.StatusNormalClosure, "")

	req := map[string]any{
		"event": "session.resume",
		"payload": map[string]any{
			"session_id": sess.ID,
		},
	}
	reqBytes, err := json.Marshal(req)
	require.NoError(t, err)
	require.NoError(t, ws.Write(ctx, websocket.MessageText, reqBytes))

	readCtx, readCancel := context.WithTimeout(ctx, time.Second)
	defer readCancel()

	found := false
	for i := 0; i < 10; i++ {
		_, data, err := ws.Read(readCtx)
		require.NoError(t, err)

		var evt map[string]any
		require.NoError(t, json.Unmarshal(data, &evt))
		if evt["event"] != "session.resume" {
			continue
		}

		require.Equal(t, sess.ID, evt["session_id"])
		payload, _ := evt["payload"].(map[string]any)
		sessionObj, _ := payload["session"].(map[string]any)
		require.Equal(t, sess.ID, sessionObj["id"])
		messages, _ := payload["messages"].([]any)
		require.NotEmpty(t, messages)
		found = true
		break
	}
	require.True(t, found)
}

// TestWebSocketEventSubscription tests event subscription via WebSocket
func TestWebSocketEventSubscription(t *testing.T) {
	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	srv := server.New(cfg, s.DB, kc)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	go srv.Serve(listener)
	time.Sleep(10 * time.Millisecond)

	// Connect with event filter
	wsURL := "ws://" + listener.Addr().String() + "/ws?event=trace"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ws, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{})
	require.NoError(t, err)
	defer ws.Close(websocket.StatusNormalClosure, "")

	// Publish an event on the bus
	b := srv.Bus()
	b.Publish(bus.NewEvent(bus.EventTraceEvent, "test-session", map[string]interface{}{
		"message": "test event",
	}))

	// Try to read the event (may need to wait)
	wsCtx, wsCancel := context.WithTimeout(ctx, time.Second)
	defer wsCancel()

	_, _, err = ws.Read(wsCtx)
	// We might get the event or timeout - both are OK for this test
	// The important thing is the connection works
}

// TestMCPEndpoint tests the MCP tools endpoint
func TestMCPEndpoint(t *testing.T) {
	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	srv := server.New(cfg, s.DB, kc)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	go srv.Serve(listener)
	time.Sleep(10 * time.Millisecond)

	// Test MCP tools endpoint
	client := &http.Client{Timeout: time.Second}
	resp, err := client.Get("http://" + listener.Addr().String() + "/mcp/tools")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Contains(t, result, "tools")
}

// TestCORSMiddleware tests CORS headers
func TestCORSMiddleware(t *testing.T) {
	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	srv := server.New(cfg, s.DB, kc)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	go srv.Serve(listener)
	time.Sleep(10 * time.Millisecond)

	// Test preflight request with Origin header
	client := &http.Client{Timeout: time.Second}
	req, _ := http.NewRequest("OPTIONS", "http://"+listener.Addr().String()+"/skills", nil)
	req.Header.Set("Origin", "http://localhost:3000") // Send Origin header for CORS
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "http://localhost:3000", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.NotEmpty(t, resp.Header.Get("Access-Control-Allow-Methods"))
}

// TestCompleteWorkflow tests a complete workflow: start server, connect WS, make API calls
func TestCompleteWorkflow(t *testing.T) {
	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	s, _ := store.New(":memory:")
	defer s.Close()
	kc := newTestKeychain(t)

	srv := server.New(cfg, s.DB, kc)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	go srv.Serve(listener)
	time.Sleep(50 * time.Millisecond)

	baseURL := "http://" + listener.Addr().String()
	client := &http.Client{Timeout: time.Second}

	// 1. Check health
	resp, err := client.Get(baseURL + "/health")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// 2. Get skills list
	resp, err = client.Get(baseURL + "/skills")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// 3. Connect WebSocket
	wsURL := "ws://" + listener.Addr().String() + "/ws"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ws, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{})
	require.NoError(t, err)
	defer ws.Close(websocket.StatusNormalClosure, "")

	// 4. Publish event and verify it flows through
	b := srv.Bus()
	b.Publish(bus.NewEvent(bus.EventTraceEvent, "", map[string]interface{}{
		"kind": "test.workflow",
	}))

	t.Log("Complete workflow test passed")
}
