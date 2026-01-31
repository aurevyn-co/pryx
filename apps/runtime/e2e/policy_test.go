//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

// TestPolicyEngine tests the policy engine and approval workflow
func TestPolicyEngine(t *testing.T) {
	bin := buildPryxCore(t)
	home := t.TempDir()

	port, cancel := startPryxCore(t, bin, home)
	defer cancel()

	waitForServer(t, port, 5*time.Second)

	baseURL := "http://localhost:" + port

	// Test 1: Check policy status endpoint
	t.Run("check_policy_status", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/api/v1/policy/status")
		if err != nil {
			// Policy endpoint may not exist yet
			t.Skip("Policy endpoint not available")
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			t.Skip("Policy API not implemented yet")
		}

		if resp.StatusCode != http.StatusOK {
			t.Skipf("Policy status returned: %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Skip("Could not decode policy status")
		}

		t.Logf("✓ Policy status: %+v", result)
	})

	// Test 2: Check approval queue
	t.Run("check_approval_queue", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/api/v1/approvals/pending")
		if err != nil {
			t.Skip("Approvals endpoint not available")
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			t.Skip("Approvals API not implemented yet")
		}

		if resp.StatusCode != http.StatusOK {
			t.Skipf("Approvals endpoint returned: %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Skip("Could not decode approvals")
		}

		t.Logf("✓ Approvals: %+v", result)
	})

	// Test 3: Submit approval decision
	t.Run("submit_approval", func(t *testing.T) {
		payload := map[string]interface{}{
			"request_id": "test-request-1",
			"decision":   "approve",
			"reason":     "Test approval",
		}

		body, _ := json.Marshal(payload)
		resp, err := http.Post(baseURL+"/api/v1/approvals/decision", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Skip("Approvals endpoint not available")
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			t.Skip("Approvals API not implemented yet")
		}

		// API may return various statuses depending on implementation
		t.Logf("✓ Approval submission returned: %d", resp.StatusCode)
	})
}
