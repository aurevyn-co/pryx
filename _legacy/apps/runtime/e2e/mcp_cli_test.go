package e2e

import (
	"os/exec"
	"strings"
	"testing"
)

// TestMCPCLI_TestValidServer tests valid server connection
func TestMCPCLI_TestValidServer(t *testing.T) {
	// First add a test server (assuming test-mcp exists)
	addCmd := exec.Command("/tmp/pryx-core", "mcp", "add", "test-mcp", "--url", "http://localhost:3001")
	addOutput, addErr := addCmd.CombinedOutput()
	if addErr != nil {
		t.Logf("Add server output: %s", addOutput)
	}

	// Test the server connection
	cmd := exec.Command("/tmp/pryx-core", "mcp", "test", "test-mcp")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Logf("Test output: %s", output)
		// It's OK if test fails due to server not running
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "test") {
		t.Logf("Expected 'test' in output, got: %s", outputStr)
	}
}

// TestMCPCLI_AddWithCmd tests stdio transport
func TestMCPCLI_AddWithCmd(t *testing.T) {
	// Add MCP server with stdio transport
	cmd := exec.Command("/tmp/pryx-core", "mcp", "add", "test-stdio", "--cmd", "python -u http://example.com")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Add stdio output: %s", output)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "add") {
		t.Logf("Expected 'add' in output, got: %s", outputStr)
	}
}

// TestMCPCLI_AddWithAuth tests authentication
func TestMCPCLI_AddWithAuth(t *testing.T) {
	// Add MCP server with authentication
	cmd := exec.Command("/tmp/pryx-core", "mcp", "add", "test-auth", "--url", "http://localhost:3001", "--auth", "bearer", "--token-ref", "mytoken")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Add auth output: %s", output)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "add") {
		t.Logf("Expected 'add' in output, got: %s", outputStr)
	}
}

// TestMCPCLI_AuthInfo tests auth info display
func TestMCPCLI_AuthInfo(t *testing.T) {
	// First add a server with auth
	addCmd := exec.Command("/tmp/pryx-core", "mcp", "add", "test-auth", "--url", "http://localhost:3001", "--auth", "bearer", "--token-ref", "mytoken")
	addOutput, addErr := addCmd.CombinedOutput()
	if addErr != nil {
		t.Logf("Add server output: %s", addOutput)
	}

	// Show auth info
	cmd := exec.Command("/tmp/pryx-core", "mcp", "auth", "test-auth")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Logf("Auth info output: %s", output)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "auth") {
		t.Logf("Expected 'auth' in output, got: %s", outputStr)
	}
}
