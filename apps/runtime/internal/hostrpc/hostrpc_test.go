package hostrpc

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestNewClient(t *testing.T) {
	in := bytes.NewReader([]byte{})
	out := bytes.NewBuffer([]byte{})

	client := NewClient(in, out)

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.in == nil {
		t.Error("Expected non-nil input reader")
	}

	if client.out == nil {
		t.Error("Expected non-nil output writer")
	}
}

func TestRequestPermissionApproved(t *testing.T) {
	// Mock response from host
	response := PermissionResult{Approved: true}
	respBody, _ := json.Marshal(rpcResponse{
		JSONRPC: "2.0",
		Result:  json.RawMessage(mustMarshal(response)),
		ID:      2,
	})

	// Add newline to simulate proper protocol
	respBody = append(respBody, '\n')

	in := pipeReader{buf: respBody}
	out := bytes.NewBuffer([]byte{})

	client := NewClient(&in, out)

	req := PermissionRequest{
		Description: "Test permission",
		Intent:      "test-intent",
	}

	approved, err := client.RequestPermission(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !approved {
		t.Error("Expected permission to be approved")
	}
}

func TestRequestPermissionDenied(t *testing.T) {
	// Mock response from host
	response := PermissionResult{Approved: false}
	respBody, _ := json.Marshal(rpcResponse{
		JSONRPC: "2.0",
		Result:  json.RawMessage(mustMarshal(response)),
		ID:      2,
	})

	// Add newline to simulate proper protocol
	respBody = append(respBody, '\n')

	in := pipeReader{buf: respBody}
	out := bytes.NewBuffer([]byte{})

	client := NewClient(&in, out)

	req := PermissionRequest{
		Description: "Test permission",
	}

	approved, err := client.RequestPermission(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if approved {
		t.Error("Expected permission to be denied")
	}
}

// pipeReader simulates a pipe with buffered data
type pipeReader struct {
	buf []byte
	pos int
}

func (r *pipeReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.buf) {
		return 0, nil // Return EOF gracefully
	}
	n = copy(p, r.buf[r.pos:])
	r.pos += n
	return n, nil
}

func TestRequestPermissionError(t *testing.T) {
	// Mock error response from host
	respBody, _ := json.Marshal(rpcResponse{
		JSONRPC: "2.0",
		Error: &rpcError{
			Code:    -32600,
			Message: "Invalid Request",
		},
		ID: 2,
	})

	in := bytes.NewReader(respBody)
	out := bytes.NewBuffer([]byte{})

	client := NewClient(in, out)

	req := PermissionRequest{
		Description: "Test permission",
	}

	_, err := client.RequestPermission(req)
	if err == nil {
		t.Error("Expected error from host RPC")
	}
}

func TestRequestPermissionMismatchedID(t *testing.T) {
	// Mock response with wrong ID
	response := PermissionResult{Approved: true}
	respBody, _ := json.Marshal(rpcResponse{
		JSONRPC: "2.0",
		Result:  json.RawMessage(mustMarshal(response)),
		ID:      999, // Wrong ID
	})

	in := bytes.NewReader(respBody)
	out := bytes.NewBuffer([]byte{})

	client := NewClient(in, out)

	req := PermissionRequest{
		Description: "Test permission",
	}

	_, err := client.RequestPermission(req)
	if err == nil {
		t.Error("Expected error for mismatched response ID")
	}
}

// Helper function to marshal objects to JSON
func mustMarshal(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}
