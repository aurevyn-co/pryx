package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSSETransport_Basic(t *testing.T) {
	server := NewMockServer()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/sse" {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			flusher, ok := w.(http.Flusher)
			if !ok {
				t.Fatal("ResponseWriter doesn't support flushing")
			}

			for {
				select {
				case <-r.Context().Done():
					return
				default:
					flusher.Flush()
					time.Sleep(100 * time.Millisecond)
				}
			}
		} else if r.URL.Path == "/message" {
			var req RPCRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			resp := server.HandleRequest(r.Context(), req)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer ts.Close()

	transport := NewSSETransport(ts.URL, nil)
	defer transport.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := transport.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	if !transport.IsConnected() {
		t.Error("Should be connected")
	}
}

func TestSSETransport_IsConnected(t *testing.T) {
	transport := NewSSETransport("http://localhost:99999", nil)

	if transport.IsConnected() {
		t.Error("Should not be connected before Connect()")
	}
}

func TestSSETransport_Close(t *testing.T) {
	transport := NewSSETransport("http://localhost:99999", nil)

	if err := transport.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	if transport.IsConnected() {
		t.Error("Should not be connected after Close()")
	}
}

func TestSSETransport_SetEndpoints(t *testing.T) {
	transport := NewSSETransport("http://example.com", nil)

	transport.SetEndpoints("/custom-sse", "/custom-message")

	if transport.sseEndpoint != "/custom-sse" {
		t.Errorf("Expected sse endpoint '/custom-sse', got '%s'", transport.sseEndpoint)
	}

	if transport.postEndpoint != "/custom-message" {
		t.Errorf("Expected post endpoint '/custom-message', got '%s'", transport.postEndpoint)
	}
}
