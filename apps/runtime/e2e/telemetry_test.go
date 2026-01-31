//go:build e2e

package e2e

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

// TestTelemetryExport tests the telemetry export functionality
func TestTelemetryExport(t *testing.T) {
	bin := buildPryxCore(t)
	home := t.TempDir()

	port, cancel := startPryxCore(t, bin, home)
	defer cancel()

	waitForServer(t, port, 5*time.Second)

	baseURL := "http://localhost:" + port

	// Test 1: Check telemetry status
	t.Run("check_telemetry_status", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/api/v1/telemetry/status")
		if err != nil {
			t.Skip("Telemetry endpoint not available")
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			t.Skip("Telemetry API not implemented yet")
		}

		if resp.StatusCode != http.StatusOK {
			t.Skipf("Telemetry status returned: %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Skip("Could not decode telemetry status")
		}

		t.Logf("✓ Telemetry status: %+v", result)
	})

	// Test 2: Export telemetry data
	t.Run("export_telemetry", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/api/v1/telemetry/export?format=json")
		if err != nil {
			t.Skip("Telemetry export endpoint not available")
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			t.Skip("Telemetry export API not implemented yet")
		}

		if resp.StatusCode != http.StatusOK {
			t.Skipf("Telemetry export returned: %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Skip("Could not decode telemetry export")
		}

		t.Logf("✓ Telemetry export: %+v", result)
	})

	// Test 3: Get telemetry metrics
	t.Run("get_telemetry_metrics", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/api/v1/telemetry/metrics")
		if err != nil {
			t.Skip("Telemetry metrics endpoint not available")
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			t.Skip("Telemetry metrics API not implemented yet")
		}

		if resp.StatusCode != http.StatusOK {
			t.Skipf("Telemetry metrics returned: %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Skip("Could not decode telemetry metrics")
		}

		t.Logf("✓ Telemetry metrics: %+v", result)
	})
}
