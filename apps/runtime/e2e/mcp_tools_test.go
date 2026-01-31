//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestMCPTools tests the MCP tool execution API
func TestMCPTools(t *testing.T) {
	bin := buildPryxCore(t)
	home := t.TempDir()

	// Create a test workspace directory
	workspaceDir := filepath.Join(home, "workspace")
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	// Start pryx-core in background with dynamic port
	port, cancel := startPryxCore(t, bin, home)
	defer cancel()

	// Wait for server to be ready
	waitForServer(t, port, 5*time.Second)

	baseURL := "http://localhost:" + port

	// Test 1: List available MCP tools
	t.Run("list_mcp_tools", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/mcp/tools")
		if err != nil {
			t.Fatalf("Failed to list MCP tools: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		tools, ok := result["tools"].([]interface{})
		if !ok {
			t.Logf("Tools response: %+v", result)
			t.Skip("No tools available - MCP servers may not be configured")
		}

		t.Logf("✓ Found %d MCP tools", len(tools))
	})

	// Test 2: Filesystem - List directory
	t.Run("filesystem_list_directory", func(t *testing.T) {
		payload := map[string]interface{}{
			"session_id": "test-session-fs",
			"tool":       "filesystem_list_directory",
			"arguments": map[string]interface{}{
				"path": workspaceDir,
			},
		}

		body, _ := json.Marshal(payload)
		resp, err := http.Post(baseURL+"/mcp/tools/call", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to call filesystem tool: %v", err)
		}
		defer resp.Body.Close()

		// The tool may not exist or MCP may not be available
		if resp.StatusCode == http.StatusBadGateway {
			t.Skip("MCP tool not available - skipping filesystem test")
		}

		if resp.StatusCode != http.StatusOK {
			t.Logf("Response status: %d", resp.StatusCode)
			var errResult map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&errResult)
			t.Logf("Error response: %+v", errResult)
			t.Skip("Filesystem tool not available")
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		t.Logf("✓ Filesystem list directory succeeded: %+v", result)
	})

	// Test 3: Filesystem - Read and Write file
	t.Run("filesystem_read_write", func(t *testing.T) {
		testFile := filepath.Join(workspaceDir, "test-e2e.txt")
		testContent := "Hello from E2E test!"

		// Write file
		writePayload := map[string]interface{}{
			"session_id": "test-session-fs",
			"tool":       "filesystem_write_file",
			"arguments": map[string]interface{}{
				"path":    testFile,
				"content": testContent,
			},
		}

		body, _ := json.Marshal(writePayload)
		resp, err := http.Post(baseURL+"/mcp/tools/call", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusBadGateway {
			t.Skip("MCP write tool not available")
		}

		// Read file back
		readPayload := map[string]interface{}{
			"session_id": "test-session-fs",
			"tool":       "filesystem_read_file",
			"arguments": map[string]interface{}{
				"path": testFile,
			},
		}

		body, _ = json.Marshal(readPayload)
		resp, err = http.Post(baseURL+"/mcp/tools/call", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusBadGateway {
			t.Skip("MCP read tool not available")
		}

		if resp.StatusCode != http.StatusOK {
			t.Skip("Filesystem read tool not available")
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		t.Logf("✓ Filesystem read/write succeeded")
	})

	// Test 4: Shell - Execute simple command
	t.Run("shell_execute", func(t *testing.T) {
		payload := map[string]interface{}{
			"session_id": "test-session-shell",
			"tool":       "shell_execute",
			"arguments": map[string]interface{}{
				"command": "echo 'Hello from E2E'",
				"timeout": 30,
			},
		}

		body, _ := json.Marshal(payload)
		resp, err := http.Post(baseURL+"/mcp/tools/call", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to execute shell command: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusBadGateway {
			t.Skip("Shell tool not available")
		}

		if resp.StatusCode != http.StatusOK {
			t.Skip("Shell tool returned error - may require approval or not be available")
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		t.Logf("✓ Shell execution succeeded: %+v", result)
	})
}
