package audit

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// ExportService handles audit log export operations
type ExportService struct {
	repo *AuditRepository
}

// NewExportService creates a new export service
func NewExportService(repo *AuditRepository) *ExportService {
	return &ExportService{repo: repo}
}

// Export exports audit logs in the specified format
func (s *ExportService) Export(opts ExportOptions) ([]byte, error) {
	entries, err := s.repo.Query(opts.QueryOptions)
	if err != nil {
		return nil, err
	}

	if opts.Sanitize {
		entries = sanitizeEntries(entries)
	}

	switch opts.Format {
	case "json":
		return exportJSON(entries, opts.IncludeCost)
	case "csv":
		return exportCSV(entries, opts.IncludeCost)
	case "jsonl":
		return exportJSONL(entries, opts.IncludeCost)
	default:
		return nil, fmt.Errorf("unsupported format: %s", opts.Format)
	}
}

// sanitizeEntries removes sensitive information from entries
func sanitizeEntries(entries []*AuditEntry) []*AuditEntry {
	for _, entry := range entries {
		// Remove user ID for privacy
		entry.UserID = ""

		// Sanitize payload to remove sensitive keys
		if entry.Payload != nil {
			entry.Payload = sanitizeData(entry.Payload)
		}

		// Clear error messages that might contain sensitive info
		if entry.ErrorMsg != "" && len(entry.ErrorMsg) > 100 {
			entry.ErrorMsg = "[error details redacted]"
		}
	}
	return entries
}

// sanitizeData recursively removes sensitive keys from data
func sanitizeData(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			if isSensitiveKey(key) {
				result[key] = "[REDACTED]"
			} else {
				result[key] = sanitizeData(value)
			}
		}
		return result
	case []interface{}:
		for i, item := range v {
			v[i] = sanitizeData(item)
		}
		return v
	default:
		return v
	}
}

// isSensitiveKey checks if a key might contain sensitive information
func isSensitiveKey(key string) bool {
	sensitiveKeys := []string{
		"password", "secret", "token", "key", "auth", "credential",
		"private", "apikey", "api_key", "bearer", "authorization",
	}

	keyLower := fmt.Sprintf("%s", key)
	for _, sensitive := range sensitiveKeys {
		if containsString(keyLower, sensitive) {
			return true
		}
	}
	return false
}

// containsString checks if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// exportJSON exports entries as formatted JSON
func exportJSON(entries []*AuditEntry, includeCost bool) ([]byte, error) {
	type exportEntry struct {
		ID          string    `json:"id"`
		Timestamp   time.Time `json:"timestamp"`
		SessionID   string    `json:"session_id,omitempty"`
		Surface     string    `json:"surface,omitempty"`
		Tool        string    `json:"tool,omitempty"`
		Action      string    `json:"action"`
		Description string    `json:"description,omitempty"`
		Success     bool      `json:"success"`
		Duration    *int64    `json:"duration_ms,omitempty"`
		Cost        *CostInfo `json:"cost,omitempty"`
	}

	exportData := make([]exportEntry, len(entries))
	for i, entry := range entries {
		e := exportEntry{
			ID:          entry.ID,
			Timestamp:   entry.Timestamp,
			SessionID:   entry.SessionID,
			Surface:     entry.Surface,
			Tool:        entry.Tool,
			Action:      string(entry.Action),
			Description: entry.Description,
			Success:     entry.Success,
			Duration:    entry.Duration,
		}

		if includeCost && entry.Cost != nil {
			e.Cost = entry.Cost
		}

		exportData[i] = e
	}

	return json.MarshalIndent(exportData, "", "  ")
}

// exportJSONL exports entries as JSONL (one JSON object per line)
func exportJSONL(entries []*AuditEntry, includeCost bool) ([]byte, error) {
	var buf bytes.Buffer

	for _, entry := range entries {
		// Create export-friendly entry
		exportEntry := map[string]interface{}{
			"id":        entry.ID,
			"timestamp": entry.Timestamp.Format(time.RFC3339),
			"action":    string(entry.Action),
			"success":   entry.Success,
		}

		if entry.SessionID != "" {
			exportEntry["session_id"] = entry.SessionID
		}
		if entry.Surface != "" {
			exportEntry["surface"] = entry.Surface
		}
		if entry.Tool != "" {
			exportEntry["tool"] = entry.Tool
		}
		if entry.Description != "" {
			exportEntry["description"] = entry.Description
		}
		if entry.Duration != nil {
			exportEntry["duration_ms"] = *entry.Duration
		}
		if includeCost && entry.Cost != nil {
			exportEntry["cost"] = entry.Cost
		}

		line, err := json.Marshal(exportEntry)
		if err != nil {
			return nil, err
		}

		buf.Write(line)
		buf.WriteByte('\n')
	}

	return buf.Bytes(), nil
}

// exportCSV exports entries as CSV
func exportCSV(entries []*AuditEntry, includeCost bool) ([]byte, error) {
	buf := bytes.Buffer{}
	writer := csv.NewWriter(&buf)

	// Write header
	headers := []string{"id", "timestamp", "session_id", "surface", "tool", "action", "description", "success", "duration_ms"}
	if includeCost {
		headers = append(headers, "input_tokens", "output_tokens", "total_tokens", "total_cost", "model")
	}
	writer.Write(headers)

	// Write data rows
	for _, entry := range entries {
		row := []string{
			entry.ID,
			entry.Timestamp.Format(time.RFC3339),
			entry.SessionID,
			entry.Surface,
			entry.Tool,
			string(entry.Action),
			entry.Description,
			fmt.Sprintf("%t", entry.Success),
		}

		if entry.Duration != nil {
			row = append(row, fmt.Sprintf("%d", *entry.Duration))
		} else {
			row = append(row, "")
		}

		if includeCost && entry.Cost != nil {
			row = append(row,
				fmt.Sprintf("%d", entry.Cost.InputTokens),
				fmt.Sprintf("%d", entry.Cost.OutputTokens),
				fmt.Sprintf("%d", entry.Cost.TotalTokens),
				fmt.Sprintf("%.6f", entry.Cost.TotalCost),
				entry.Cost.Model,
			)
		}

		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), nil
}

// GetExportFormats returns the list of supported export formats
func GetExportFormats() []string {
	return []string{"json", "csv", "jsonl"}
}

// ValidateExportOptions validates export options
func ValidateExportOptions(opts *ExportOptions) error {
	if opts.Format == "" {
		opts.Format = "json"
	}

	validFormats := GetExportFormats()
	valid := false
	for _, f := range validFormats {
		if opts.Format == f {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("invalid format: %s (must be one of: %s)", opts.Format, joinStrings(validFormats, ", "))
	}

	if opts.QueryOptions.Limit <= 0 {
		opts.QueryOptions.Limit = 1000
	}

	if opts.QueryOptions.Limit > 10000 {
		opts.QueryOptions.Limit = 10000
	}

	return nil
}

// ExportWriter is an interface for streaming exports
type ExportWriter interface {
	WriteEntry(entry *AuditEntry) error
	Close() error
}

// JSONLExportWriter writes entries in JSONL format to an io.Writer
type JSONLExportWriter struct {
	writer io.Writer
}

// NewJSONLExportWriter creates a new JSONL export writer
func NewJSONLExportWriter(w io.Writer) *JSONLExportWriter {
	return &JSONLExportWriter{writer: w}
}

// WriteEntry writes a single entry in JSONL format
func (w *JSONLExportWriter) WriteEntry(entry *AuditEntry) error {
	// Create export-friendly entry
	exportEntry := map[string]interface{}{
		"id":        entry.ID,
		"timestamp": entry.Timestamp.Format(time.RFC3339),
		"action":    string(entry.Action),
		"success":   entry.Success,
	}

	if entry.SessionID != "" {
		exportEntry["session_id"] = entry.SessionID
	}
	if entry.Surface != "" {
		exportEntry["surface"] = entry.Surface
	}
	if entry.Tool != "" {
		exportEntry["tool"] = entry.Tool
	}

	line, err := json.Marshal(exportEntry)
	if err != nil {
		return err
	}

	_, err = w.writer.Write(append(line, '\n'))
	return err
}

// Close flushes the writer
func (w *JSONLExportWriter) Close() error {
	if f, ok := w.writer.(flusher); ok {
		return f.Flush()
	}
	return nil
}

type flusher interface {
	Flush() error
}
