package audit

import (
	"os"
	"testing"
	"time"

	"pryx-core/internal/store"

	_ "github.com/mattn/go-sqlite3"
)

// TestRepositoryCreate tests creating audit entries
func TestRepositoryCreate(t *testing.T) {
	// Create temporary database
	tmpFile, err := os.CreateTemp("", "audit_test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	s, err := store.New(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	repo := NewAuditRepository(s.DB)

	// Test creating an entry
	entry := &AuditEntry{
		SessionID:   "test-session-123",
		Surface:     "tui",
		Tool:        "filesystem.read",
		Action:      ActionToolExecute,
		Description: "Read file operation",
		Success:     true,
		Duration:    int64Ptr(150),
		Cost: &CostInfo{
			InputTokens:  100,
			OutputTokens: 200,
			TotalTokens:  300,
			InputCost:    0.001,
			OutputCost:   0.002,
			TotalCost:    0.003,
			Model:        "gpt-4",
		},
	}

	err = repo.Create(entry)
	if err != nil {
		t.Fatalf("Failed to create audit entry: %v", err)
	}

	// Verify the entry was created
	if entry.ID == "" {
		t.Error("Expected entry ID to be set")
	}
}

// TestRepositoryQuery tests querying audit entries
func TestRepositoryQuery(t *testing.T) {
	// Create temporary database
	tmpFile, err := os.CreateTemp("", "audit_test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	s, err := store.New(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	repo := NewAuditRepository(s.DB)

	// Create multiple entries
	now := time.Now()
	for i := 0; i < 5; i++ {
		entry := &AuditEntry{
			SessionID:   "session-" + string(rune('A'+i)),
			Surface:     "tui",
			Tool:        "tool-" + string(rune('1'+i)),
			Action:      ActionToolExecute,
			Description: "Test entry " + string(rune('1'+i)),
			Success:     true,
			Timestamp:   now.Add(time.Duration(i) * time.Hour),
		}
		if err := repo.Create(entry); err != nil {
			t.Fatalf("Failed to create audit entry: %v", err)
		}
	}

	// Query all entries
	entries, err := repo.Query(QueryOptions{Limit: 10})
	if err != nil {
		t.Fatalf("Failed to query audit entries: %v", err)
	}

	if len(entries) != 5 {
		t.Errorf("Expected 5 entries, got %d", len(entries))
	}

	// Query with filter
	filteredEntries, err := repo.Query(QueryOptions{
		SessionID: "session-B",
		Limit:     10,
	})
	if err != nil {
		t.Fatalf("Failed to query filtered entries: %v", err)
	}

	if len(filteredEntries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(filteredEntries))
	}

	if filteredEntries[0].SessionID != "session-B" {
		t.Errorf("Expected session-B, got %s", filteredEntries[0].SessionID)
	}
}

// TestRepositoryCount tests counting audit entries
func TestRepositoryCount(t *testing.T) {
	// Create temporary database
	tmpFile, err := os.CreateTemp("", "audit_test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	s, err := store.New(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	repo := NewAuditRepository(s.DB)

	// Create entries with different actions
	actions := []AuditAction{ActionToolExecute, ActionToolComplete, ActionApprovalRequest}
	for i, action := range actions {
		entry := &AuditEntry{
			SessionID: "session-count-test",
			Action:    action,
			Success:   true,
		}
		if err := repo.Create(entry); err != nil {
			t.Fatalf("Failed to create audit entry: %v", err)
		}
		_ = i
	}

	// Count all
	count, err := repo.Count(QueryOptions{})
	if err != nil {
		t.Fatalf("Failed to count entries: %v", err)
	}
	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}

	// Count filtered
	count, err = repo.Count(QueryOptions{Action: ActionToolExecute})
	if err != nil {
		t.Fatalf("Failed to count filtered entries: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}

// TestRepositoryDelete tests deleting old entries
func TestRepositoryDelete(t *testing.T) {
	// Create temporary database
	tmpFile, err := os.CreateTemp("", "audit_test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	s, err := store.New(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	repo := NewAuditRepository(s.DB)

	// Create entries
	for i := 0; i < 3; i++ {
		entry := &AuditEntry{
			SessionID: "session-delete-test",
			Action:    ActionToolExecute,
			Success:   true,
		}
		if err := repo.Create(entry); err != nil {
			t.Fatalf("Failed to create audit entry: %v", err)
		}
	}

	// Delete entries older than 0 hours (should delete all)
	count, err := repo.DeleteOlderThan(0 * time.Hour)
	if err != nil {
		t.Fatalf("Failed to delete entries: %v", err)
	}

	// Should have deleted all entries
	if count != 3 {
		t.Errorf("Expected to delete 3 entries, got %d", count)
	}

	// Verify no entries remain
	entries, err := repo.Query(QueryOptions{})
	if err != nil {
		t.Fatalf("Failed to query entries: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(entries))
	}
}

// TestExportService tests the export service
func TestExportService(t *testing.T) {
	// Create temporary database
	tmpFile, err := os.CreateTemp("", "audit_test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	s, err := store.New(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	repo := NewAuditRepository(s.DB)
	exportSvc := NewExportService(repo)

	// Create entries
	for i := 0; i < 3; i++ {
		entry := &AuditEntry{
			SessionID:   "session-export-test",
			Surface:     "tui",
			Tool:        "test-tool",
			Action:      ActionToolExecute,
			Description: "Export test entry " + string(rune('1'+i)),
			Success:     true,
			Cost: &CostInfo{
				InputTokens:  int64(100 + i*10),
				OutputTokens: int64(200 + i*20),
				TotalTokens:  int64(300 + i*30),
				TotalCost:    float64(0.003 + float64(i)*0.001),
				Model:        "gpt-4",
			},
		}
		if err := repo.Create(entry); err != nil {
			t.Fatalf("Failed to create audit entry: %v", err)
		}
	}

	// Test JSON export
	jsonData, err := exportSvc.Export(ExportOptions{
		Format:       "json",
		QueryOptions: QueryOptions{Limit: 10},
		IncludeCost:  true,
	})
	if err != nil {
		t.Fatalf("Failed to export JSON: %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("Expected JSON data to be non-empty")
	}

	// Test JSONL export
	jsonlData, err := exportSvc.Export(ExportOptions{
		Format:       "jsonl",
		QueryOptions: QueryOptions{Limit: 10},
		IncludeCost:  false,
	})
	if err != nil {
		t.Fatalf("Failed to export JSONL: %v", err)
	}

	if len(jsonlData) == 0 {
		t.Error("Expected JSONL data to be non-empty")
	}

	// Test CSV export
	csvData, err := exportSvc.Export(ExportOptions{
		Format:       "csv",
		QueryOptions: QueryOptions{Limit: 10},
		IncludeCost:  true,
	})
	if err != nil {
		t.Fatalf("Failed to export CSV: %v", err)
	}

	if len(csvData) == 0 {
		t.Error("Expected CSV data to be non-empty")
	}

	// Test sanitized export
	sanitizedData, err := exportSvc.Export(ExportOptions{
		Format:       "json",
		QueryOptions: QueryOptions{Limit: 10},
		Sanitize:     true,
	})
	if err != nil {
		t.Fatalf("Failed to export sanitized JSON: %v", err)
	}

	if len(sanitizedData) == 0 {
		t.Error("Expected sanitized data to be non-empty")
	}
}

// TestValidateExportOptions validates export options
func TestValidateExportOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    ExportOptions
		wantErr bool
	}{
		{
			name: "valid json export",
			opts: ExportOptions{
				Format: "json",
				QueryOptions: QueryOptions{
					Limit: 100,
				},
			},
			wantErr: false,
		},
		{
			name: "valid csv export",
			opts: ExportOptions{
				Format: "csv",
				QueryOptions: QueryOptions{
					Limit: 500,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid format defaults to json",
			opts: ExportOptions{
				Format: "invalid",
				QueryOptions: QueryOptions{
					Limit: 100,
				},
			},
			wantErr: true,
		},
		{
			name: "limit exceeds max",
			opts: ExportOptions{
				Format: "json",
				QueryOptions: QueryOptions{
					Limit: 20000,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateExportOptions(&tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateExportOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function to create int64 pointer
func int64Ptr(v int64) *int64 {
	return &v
}
