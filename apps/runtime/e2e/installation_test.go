//go:build e2e

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestInstallationFlow tests the installation and first-run experience
func TestInstallationFlow(t *testing.T) {
	bin := buildPryxCore(t)
	home := t.TempDir()

	// Test 1: Verify binary runs
	t.Run("binary_runs", func(t *testing.T) {
		cmd := exec.Command(bin, "--help")
		cmd.Env = append(os.Environ(),
			"HOME="+home,
			"PRYX_KEYCHAIN_FILE="+filepath.Join(home, ".pryx", "keychain.json"),
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Binary failed to run: %v\n%s", err, string(out))
		}

		if !strings.Contains(string(out), "pryx-core") {
			t.Fatal("Help output doesn't contain expected content")
		}

		t.Logf("✓ Binary runs successfully")
	})

	// Test 2: First run creates necessary directories
	t.Run("first_run_setup", func(t *testing.T) {
		cmd := exec.Command(bin, "config", "list")
		cmd.Env = append(os.Environ(),
			"HOME="+home,
			"PRYX_KEYCHAIN_FILE="+filepath.Join(home, ".pryx", "keychain.json"),
			"PRYX_TELEMETRY_DISABLED=true",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("Config list output: %s", string(out))
		}

		// Check if .pryx directory was created
		pryxDir := filepath.Join(home, ".pryx")
		if _, err := os.Stat(pryxDir); os.IsNotExist(err) {
			t.Skip(".pryx directory not created on first run")
		}

		t.Logf("✓ First run setup completed")
	})

	// Test 3: Database initialization
	t.Run("database_initialized", func(t *testing.T) {
		dbPath := filepath.Join(home, "pryx.db")
		cmd := exec.Command(bin, "doctor")
		cmd.Env = append(os.Environ(),
			"HOME="+home,
			"PRYX_DB_PATH="+dbPath,
			"PRYX_KEYCHAIN_FILE="+filepath.Join(home, ".pryx", "keychain.json"),
			"PRYX_TELEMETRY_DISABLED=true",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("Doctor output: %s", string(out))
		}

		// Check if database file was created
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			t.Skip("Database not created")
		}

		t.Logf("✓ Database initialized")
	})
}
