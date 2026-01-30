package vault

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewAuditLogger(t *testing.T) {
	tmpDir := t.TempDir()

	logger, err := NewAuditLogger(WithAuditDir(tmpDir))
	if err != nil {
		t.Fatalf("Failed to create audit logger: %v", err)
	}
	defer logger.Close()

	if logger.auditDir != tmpDir {
		t.Errorf("Expected audit dir %s, got %s", tmpDir, logger.auditDir)
	}

	if logger.retentionDays != defaultRetentionDays {
		t.Errorf("Expected retention %d, got %d", defaultRetentionDays, logger.retentionDays)
	}
}

func TestAuditLogger_Log(t *testing.T) {
	tmpDir := t.TempDir()
	logger, err := NewAuditLogger(WithAuditDir(tmpDir))
	if err != nil {
		t.Fatalf("Failed to create audit logger: %v", err)
	}
	defer logger.Close()

	err = logger.Log(ActionUnlock, "vault", "user1", true, nil, map[string]interface{}{"ip": "127.0.0.1"})
	if err != nil {
		t.Errorf("Failed to log: %v", err)
	}

	err = logger.Log(ActionRead, "api-key", "user1", true, nil, nil)
	if err != nil {
		t.Errorf("Failed to log: %v", err)
	}

	testErr := errors.New("test error")
	err = logger.Log(ActionWrite, "secret", "user2", false, testErr, nil)
	if err != nil {
		t.Errorf("Failed to log: %v", err)
	}

	logger.flushBuffer()

	entries, err := logger.Query(QueryOptions{})
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}

	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}
}

func TestAuditLogger_ConvenienceMethods(t *testing.T) {
	tmpDir := t.TempDir()
	logger, err := NewAuditLogger(WithAuditDir(tmpDir))
	if err != nil {
		t.Fatalf("Failed to create audit logger: %v", err)
	}
	defer logger.Close()

	if err := logger.LogUnlock("user1", true, nil); err != nil {
		t.Errorf("LogUnlock failed: %v", err)
	}

	if err := logger.LogLock("user1"); err != nil {
		t.Errorf("LogLock failed: %v", err)
	}

	if err := logger.LogRead("cred1", "user1", true, nil); err != nil {
		t.Errorf("LogRead failed: %v", err)
	}

	if err := logger.LogWrite("cred1", "user1", true, nil); err != nil {
		t.Errorf("LogWrite failed: %v", err)
	}

	if err := logger.LogDelete("cred1", "user1", true, nil); err != nil {
		t.Errorf("LogDelete failed: %v", err)
	}

	if err := logger.LogRotate("user1", true, nil); err != nil {
		t.Errorf("LogRotate failed: %v", err)
	}

	logger.flushBuffer()

	entries, err := logger.Query(QueryOptions{})
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}

	if len(entries) != 6 {
		t.Errorf("Expected 6 entries, got %d", len(entries))
	}
}

func TestAuditLogger_Query(t *testing.T) {
	tmpDir := t.TempDir()
	logger, err := NewAuditLogger(WithAuditDir(tmpDir))
	if err != nil {
		t.Fatalf("Failed to create audit logger: %v", err)
	}
	defer logger.Close()

	logger.Log(ActionUnlock, "vault", "user1", true, nil, nil)
	logger.Log(ActionRead, "cred1", "user1", true, nil, nil)
	logger.Log(ActionRead, "cred2", "user2", true, nil, nil)
	logger.Log(ActionWrite, "cred1", "user1", false, errors.New("fail"), nil)
	logger.flushBuffer()

	tests := []struct {
		name     string
		opts     QueryOptions
		expected int
	}{
		{
			name:     "query all",
			opts:     QueryOptions{},
			expected: 4,
		},
		{
			name: "query by actor",
			opts: QueryOptions{
				Actor: "user1",
			},
			expected: 3,
		},
		{
			name: "query by action",
			opts: QueryOptions{
				Actions: []AuditAction{ActionRead},
			},
			expected: 2,
		},
		{
			name: "query by target",
			opts: QueryOptions{
				Target: "cred1",
			},
			expected: 2,
		},
		{
			name: "query by success",
			opts: QueryOptions{
				Success: boolPtr(true),
			},
			expected: 3,
		},
		{
			name: "query by failure",
			opts: QueryOptions{
				Success: boolPtr(false),
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries, err := logger.Query(tt.opts)
			if err != nil {
				t.Fatalf("Query failed: %v", err)
			}
			if len(entries) != tt.expected {
				t.Errorf("Expected %d entries, got %d", tt.expected, len(entries))
			}
		})
	}
}

func TestAuditLogger_QueryTimeRange(t *testing.T) {
	tmpDir := t.TempDir()
	logger, err := NewAuditLogger(WithAuditDir(tmpDir))
	if err != nil {
		t.Fatalf("Failed to create audit logger: %v", err)
	}
	defer logger.Close()

	oldTime := time.Now().Add(-24 * time.Hour)
	newTime := time.Now()

	logger.Log(ActionUnlock, "vault", "user1", true, nil, nil)
	logger.flushBuffer()

	entries, err := logger.Query(QueryOptions{
		StartTime: &oldTime,
		EndTime:   &newTime,
	})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
	}

	futureTime := time.Now().Add(24 * time.Hour)
	entries, err = logger.Query(QueryOptions{
		StartTime: &futureTime,
	})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Expected 0 entries for future time, got %d", len(entries))
	}
}

func TestAuditLogger_VerifyIntegrity(t *testing.T) {
	tmpDir := t.TempDir()
	logger, err := NewAuditLogger(WithAuditDir(tmpDir))
	if err != nil {
		t.Fatalf("Failed to create audit logger: %v", err)
	}
	defer logger.Close()

	logger.Log(ActionUnlock, "vault", "user1", true, nil, nil)
	logger.Log(ActionRead, "cred1", "user1", true, nil, nil)
	logger.flushBuffer()

	start := time.Now().Add(-time.Hour)
	end := time.Now().Add(time.Hour)

	if err := logger.VerifyIntegrity(start, end); err != nil {
		t.Errorf("Integrity check failed: %v", err)
	}
}

func TestAuditLogger_Cleanup(t *testing.T) {
	tmpDir := t.TempDir()

	oldFile := filepath.Join(tmpDir, "2024-01-01.log")
	newFile := filepath.Join(tmpDir, time.Now().Format("2006-01-02")+".log")

	os.WriteFile(oldFile, []byte("old log entry\n"), 0600)
	os.WriteFile(newFile, []byte("new log entry\n"), 0600)

	logger, err := NewAuditLogger(
		WithAuditDir(tmpDir),
		WithRetentionDays(7),
	)
	if err != nil {
		t.Fatalf("Failed to create audit logger: %v", err)
	}
	defer logger.Close()

	if err := logger.Cleanup(); err != nil {
		t.Errorf("Cleanup failed: %v", err)
	}

	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Error("Old file should have been deleted")
	}

	if _, err := os.Stat(newFile); err != nil {
		t.Error("New file should still exist")
	}
}

func TestAuditLogger_Export(t *testing.T) {
	tmpDir := t.TempDir()
	logger, err := NewAuditLogger(WithAuditDir(tmpDir))
	if err != nil {
		t.Fatalf("Failed to create audit logger: %v", err)
	}
	defer logger.Close()

	logger.Log(ActionUnlock, "vault", "user1", true, nil, nil)
	logger.Log(ActionRead, "cred1", "user1", true, nil, nil)
	logger.flushBuffer()

	tmpFile := filepath.Join(tmpDir, "export.json")
	f, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create export file: %v", err)
	}
	defer f.Close()

	start := time.Now().Add(-time.Hour)
	end := time.Now().Add(time.Hour)

	if err := logger.Export(start, end, "json", f); err != nil {
		t.Errorf("Export failed: %v", err)
	}

	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read export: %v", err)
	}

	if !strings.Contains(string(content), "unlock") {
		t.Error("Export should contain unlock action")
	}

	if !strings.Contains(string(content), "read") {
		t.Error("Export should contain read action")
	}
}

func TestAuditLogger_ExportCSV(t *testing.T) {
	tmpDir := t.TempDir()
	logger, err := NewAuditLogger(WithAuditDir(tmpDir))
	if err != nil {
		t.Fatalf("Failed to create audit logger: %v", err)
	}
	defer logger.Close()

	logger.Log(ActionUnlock, "vault", "user1", true, nil, nil)
	logger.flushBuffer()

	tmpFile := filepath.Join(tmpDir, "export.csv")
	f, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create export file: %v", err)
	}
	defer f.Close()

	start := time.Now().Add(-time.Hour)
	end := time.Now().Add(time.Hour)

	if err := logger.Export(start, end, "csv", f); err != nil {
		t.Errorf("Export failed: %v", err)
	}

	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read export: %v", err)
	}

	if !strings.Contains(string(content), "timestamp,action,target") {
		t.Error("CSV should contain header")
	}

	if !strings.Contains(string(content), "unlock") {
		t.Error("CSV should contain unlock action")
	}
}

func TestAuditLogger_BackgroundFlush(t *testing.T) {
	tmpDir := t.TempDir()
	logger, err := NewAuditLogger(WithAuditDir(tmpDir))
	if err != nil {
		t.Fatalf("Failed to create audit logger: %v", err)
	}
	defer logger.Close()

	logger.Log(ActionUnlock, "vault", "user1", true, nil, nil)

	time.Sleep(6 * time.Second)

	entries, err := logger.Query(QueryOptions{})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 entry after background flush, got %d", len(entries))
	}
}

func TestAuditLogger_HashChain(t *testing.T) {
	tmpDir := t.TempDir()
	logger, err := NewAuditLogger(WithAuditDir(tmpDir))
	if err != nil {
		t.Fatalf("Failed to create audit logger: %v", err)
	}
	defer logger.Close()

	logger.Log(ActionUnlock, "vault", "user1", true, nil, nil)
	time.Sleep(10 * time.Millisecond)
	logger.Log(ActionRead, "cred1", "user1", true, nil, nil)
	time.Sleep(10 * time.Millisecond)
	logger.Log(ActionWrite, "cred1", "user1", true, nil, nil)
	logger.flushBuffer()

	start := time.Now().Add(-time.Hour)
	end := time.Now().Add(time.Hour)

	if err := logger.VerifyIntegrity(start, end); err != nil {
		t.Errorf("Hash chain verification failed: %v", err)
	}
}

func TestAuditLogger_EntryHash(t *testing.T) {
	tmpDir := t.TempDir()
	logger, err := NewAuditLogger(WithAuditDir(tmpDir))
	if err != nil {
		t.Fatalf("Failed to create audit logger: %v", err)
	}
	defer logger.Close()

	logger.Log(ActionUnlock, "vault", "user1", true, nil, map[string]interface{}{"key": "value"})
	logger.flushBuffer()

	entries, err := logger.Query(QueryOptions{})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}

	entry := entries[0]
	if entry.Hash == "" {
		t.Error("Entry should have a hash")
	}

	expectedHash := logger.calculateHash(entry)
	if entry.Hash != expectedHash {
		t.Errorf("Hash mismatch: expected %s, got %s", expectedHash, entry.Hash)
	}
}

func boolPtr(b bool) *bool {
	return &b
}
