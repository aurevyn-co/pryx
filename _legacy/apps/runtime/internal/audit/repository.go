package audit

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AuditAction represents the type of audit action
type AuditAction string

const (
	ActionSessionCreate   AuditAction = "session.create"
	ActionSessionUpdate   AuditAction = "session.update"
	ActionSessionDelete   AuditAction = "session.delete"
	ActionMessageSend     AuditAction = "message.send"
	ActionToolRequest     AuditAction = "tool.request"
	ActionToolExecute     AuditAction = "tool.execute"
	ActionToolComplete    AuditAction = "tool.complete"
	ActionToolError       AuditAction = "tool.error"
	ActionApprovalRequest AuditAction = "approval.request"
	ActionApprovalGrant   AuditAction = "approval.grant"
	ActionApprovalDeny    AuditAction = "approval.deny"
	ActionChannelMessage  AuditAction = "channel.message"
	ActionChannelStatus   AuditAction = "channel.status"
	ActionErrorOccurred   AuditAction = "error.occurred"
	ActionUserAction      AuditAction = "user.action"
)

// AuditEntry represents a single audit log entry
type AuditEntry struct {
	ID          string      `json:"id"`
	Timestamp   time.Time   `json:"timestamp"`
	SessionID   string      `json:"session_id,omitempty"`
	Surface     string      `json:"surface,omitempty"`
	Tool        string      `json:"tool,omitempty"`
	Action      AuditAction `json:"action"`
	Description string      `json:"description,omitempty"`
	Payload     interface{} `json:"payload,omitempty"`
	Cost        *CostInfo   `json:"cost,omitempty"`
	Duration    *int64      `json:"duration_ms,omitempty"`
	UserID      string      `json:"user_id,omitempty"`
	Success     bool        `json:"success"`
	ErrorMsg    string      `json:"error,omitempty"`
	Metadata    interface{} `json:"metadata,omitempty"`
}

// CostInfo represents cost tracking information
type CostInfo struct {
	InputTokens  int64   `json:"input_tokens"`
	OutputTokens int64   `json:"output_tokens"`
	TotalTokens  int64   `json:"total_tokens"`
	InputCost    float64 `json:"input_cost"`
	OutputCost   float64 `json:"output_cost"`
	TotalCost    float64 `json:"total_cost"`
	Model        string  `json:"model,omitempty"`
}

// QueryOptions defines filtering options for audit log queries
type QueryOptions struct {
	StartTime *time.Time
	EndTime   *time.Time
	SessionID string
	Surface   string
	Tool      string
	Action    AuditAction
	UserID    string
	Success   *bool
	Limit     int
	Offset    int
	OrderBy   string // "timestamp" or "created_at"
	OrderDir  string // "ASC" or "DESC"
}

// ExportOptions defines options for exporting audit logs
type ExportOptions struct {
	Format       string // "json", "csv", "jsonl"
	QueryOptions QueryOptions
	Sanitize     bool // Remove sensitive data
	IncludeCost  bool // Include cost information
}

// AuditRepository handles audit log storage operations
type AuditRepository struct {
	db *sql.DB
}

// NewAuditRepository creates a new audit repository
func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

// Create inserts a new audit entry into the database
func (r *AuditRepository) Create(entry *AuditEntry) error {
	id := uuid.New().String()
	if entry.ID == "" {
		entry.ID = id
	}

	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}

	// Serialize payload and metadata to JSON
	payloadJSON, _ := json.Marshal(entry.Payload)
	metadataJSON, _ := json.Marshal(entry.Metadata)
	costJSON, _ := json.Marshal(entry.Cost)

	query := `
		INSERT INTO audit_log (
			id, timestamp, session_id, surface, tool, action,
			description, payload, cost, duration, user_id,
			success, error_msg, metadata, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		entry.ID,
		entry.Timestamp,
		entry.SessionID,
		entry.Surface,
		entry.Tool,
		entry.Action,
		entry.Description,
		string(payloadJSON),
		string(costJSON),
		entry.Duration,
		entry.UserID,
		entry.Success,
		entry.ErrorMsg,
		string(metadataJSON),
		time.Now().UTC(),
	)

	return err
}

// Query retrieves audit entries based on the provided options
func (r *AuditRepository) Query(opts QueryOptions) ([]*AuditEntry, error) {
	query := buildQuery(opts)
	rows, err := r.db.Query(query.Query, query.Args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*AuditEntry
	for rows.Next() {
		entry, err := scanEntryFromRows(rows)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// Count returns the total number of entries matching the query
func (r *AuditRepository) Count(opts QueryOptions) (int, error) {
	query := buildCountQuery(opts)
	var count int
	err := r.db.QueryRow(query.Query, query.Args...).Scan(&count)
	return count, err
}

// Delete removes entries older than the specified duration
func (r *AuditRepository) DeleteOlderThan(duration time.Duration) (int64, error) {
	cutoff := time.Now().UTC().Add(-duration)
	result, err := r.db.Exec(
		"DELETE FROM audit_log WHERE created_at < ?",
		cutoff,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// scanEntryFromRows converts database rows to an AuditEntry
func scanEntryFromRows(rows *sql.Rows) (*AuditEntry, error) {
	entry := &AuditEntry{}

	var payloadJSON, metadataJSON, costJSON []byte
	var durationPtr *int64

	err := rows.Scan(
		&entry.ID,
		&entry.Timestamp,
		&entry.SessionID,
		&entry.Surface,
		&entry.Tool,
		&entry.Action,
		&entry.Description,
		&payloadJSON,
		&costJSON,
		&durationPtr,
		&entry.UserID,
		&entry.Success,
		&entry.ErrorMsg,
		&metadataJSON,
	)

	if err != nil {
		return nil, err
	}

	// Parse JSON fields
	if len(payloadJSON) > 0 {
		json.Unmarshal(payloadJSON, &entry.Payload)
	}
	if len(metadataJSON) > 0 {
		json.Unmarshal(metadataJSON, &entry.Metadata)
	}
	if len(costJSON) > 0 {
		json.Unmarshal(costJSON, &entry.Cost)
	}
	if durationPtr != nil {
		entry.Duration = durationPtr
	}

	return entry, nil
}

// QueryResult holds the query string and arguments
type QueryResult struct {
	Query string
	Args  []interface{}
}

// buildQuery constructs the SQL query and arguments from QueryOptions
func buildQuery(opts QueryOptions) QueryResult {
	var conditions []string
	var args []interface{}

	if opts.StartTime != nil {
		conditions = append(conditions, "timestamp >= ?")
		args = append(args, *opts.StartTime)
	}
	if opts.EndTime != nil {
		conditions = append(conditions, "timestamp <= ?")
		args = append(args, *opts.EndTime)
	}
	if opts.SessionID != "" {
		conditions = append(conditions, "session_id = ?")
		args = append(args, opts.SessionID)
	}
	if opts.Surface != "" {
		conditions = append(conditions, "surface = ?")
		args = append(args, opts.Surface)
	}
	if opts.Tool != "" {
		conditions = append(conditions, "tool = ?")
		args = append(args, opts.Tool)
	}
	if opts.Action != "" {
		conditions = append(conditions, "action = ?")
		args = append(args, opts.Action)
	}
	if opts.UserID != "" {
		conditions = append(conditions, "user_id = ?")
		args = append(args, opts.UserID)
	}
	if opts.Success != nil {
		conditions = append(conditions, "success = ?")
		args = append(args, *opts.Success)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + joinStrings(conditions, " AND ")
	}

	orderBy := "timestamp"
	if opts.OrderBy != "" {
		orderBy = opts.OrderBy
	}
	orderDir := "DESC"
	if opts.OrderDir != "" {
		orderDir = opts.OrderDir
	}

	limit := 100
	if opts.Limit > 0 {
		limit = opts.Limit
	}

	query := fmt.Sprintf(`
		SELECT id, timestamp, session_id, surface, tool, action,
		       description, payload, cost, duration, user_id,
		       success, error_msg, metadata
		FROM audit_log
		%s
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, whereClause, orderBy, orderDir)

	args = append(args, limit, opts.Offset)

	return QueryResult{Query: query, Args: args}
}

// buildCountQuery constructs a count query from QueryOptions
func buildCountQuery(opts QueryOptions) QueryResult {
	var conditions []string
	var args []interface{}

	if opts.StartTime != nil {
		conditions = append(conditions, "timestamp >= ?")
		args = append(args, *opts.StartTime)
	}
	if opts.EndTime != nil {
		conditions = append(conditions, "timestamp <= ?")
		args = append(args, *opts.EndTime)
	}
	if opts.SessionID != "" {
		conditions = append(conditions, "session_id = ?")
		args = append(args, opts.SessionID)
	}
	if opts.Surface != "" {
		conditions = append(conditions, "surface = ?")
		args = append(args, opts.Surface)
	}
	if opts.Tool != "" {
		conditions = append(conditions, "tool = ?")
		args = append(args, opts.Tool)
	}
	if opts.Action != "" {
		conditions = append(conditions, "action = ?")
		args = append(args, opts.Action)
	}
	if opts.UserID != "" {
		conditions = append(conditions, "user_id = ?")
		args = append(args, opts.UserID)
	}
	if opts.Success != nil {
		conditions = append(conditions, "success = ?")
		args = append(args, *opts.Success)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + joinStrings(conditions, " AND ")
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM audit_log %s", whereClause)
	return QueryResult{Query: query, Args: args}
}

// joinStrings joins strings with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
