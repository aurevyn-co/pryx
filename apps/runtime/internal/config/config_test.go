package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("PRYX_LISTEN_ADDR")
	os.Unsetenv("PRYX_DB_PATH")

	config := Load()

	if config.ListenAddr != ":3000" {
		t.Errorf("Expected ListenAddr ':3000', got '%s'", config.ListenAddr)
	}

	if config.DatabasePath != "pryx.db" {
		t.Errorf("Expected DatabasePath 'pryx.db', got '%s'", config.DatabasePath)
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	// Set environment variables
	os.Setenv("PRYX_LISTEN_ADDR", "0.0.0.0:8080")
	os.Setenv("PRYX_DB_PATH", "/tmp/test.db")
	defer func() {
		os.Unsetenv("PRYX_LISTEN_ADDR")
		os.Unsetenv("PRYX_DB_PATH")
	}()

	config := Load()

	if config.ListenAddr != "0.0.0.0:8080" {
		t.Errorf("Expected ListenAddr '0.0.0.0:8080', got '%s'", config.ListenAddr)
	}

	if config.DatabasePath != "/tmp/test.db" {
		t.Errorf("Expected DatabasePath '/tmp/test.db', got '%s'", config.DatabasePath)
	}
}

func TestPartialEnvironment(t *testing.T) {
	// Set only one environment variable
	os.Setenv("PRYX_LISTEN_ADDR", "localhost:9000")
	defer os.Unsetenv("PRYX_LISTEN_ADDR")

	// Clear the other
	os.Unsetenv("PRYX_DB_PATH")

	config := Load()

	if config.ListenAddr != "localhost:9000" {
		t.Errorf("Expected ListenAddr 'localhost:9000', got '%s'", config.ListenAddr)
	}

	// DatabasePath should still have default
	if config.DatabasePath != "pryx.db" {
		t.Errorf("Expected DatabasePath 'pryx.db', got '%s'", config.DatabasePath)
	}
}
