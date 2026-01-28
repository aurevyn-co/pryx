package db

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	// Create a temporary database file
	tmpFile, err := os.CreateTemp("", "pryx_test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	db, err := Init(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Verify the database is accessible
	if err := db.Ping(); err != nil {
		t.Errorf("Database ping failed: %v", err)
	}
}

func TestInitInvalidPath(t *testing.T) {
	// Test with an invalid path (directory that doesn't exist)
	_, err := Init("/nonexistent/path/test.db")
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

func TestInitMultipleConnections(t *testing.T) {
	// Create a temporary database file
	tmpFile, err := os.CreateTemp("", "pryx_test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Open first connection
	db1, err := Init(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to initialize first database connection: %v", err)
	}
	defer db1.Close()

	// Open second connection
	db2, err := Init(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to initialize second database connection: %v", err)
	}
	defer db2.Close()

	// Both connections should be valid
	if err := db1.Ping(); err != nil {
		t.Errorf("First database connection ping failed: %v", err)
	}

	if err := db2.Ping(); err != nil {
		t.Errorf("Second database connection ping failed: %v", err)
	}
}
