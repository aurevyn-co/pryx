package audit

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pryx-core/internal/store"

	"github.com/go-chi/chi/v5"
)

// Handler provides HTTP handlers for audit log operations
type Handler struct {
	repo  *AuditRepository
	store *store.Store
}

// NewHandler creates a new audit handler
func NewHandler(store *store.Store) *Handler {
	return &Handler{
		repo:  NewAuditRepository(store.DB),
		store: store,
	}
}

// RegisterRoutes registers audit log routes on the router
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/audit", h.handleQuery)
	r.Get("/audit/export", h.handleExport)
	r.Get("/audit/count", h.handleCount)
	r.Delete("/audit", h.handleDelete)
}

// handleQuery handles GET /audit - query audit logs
func (h *Handler) handleQuery(w http.ResponseWriter, r *http.Request) {
	opts, err := parseQueryOptions(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entries, err := h.repo.Query(opts)
	if err != nil {
		http.Error(w, "failed to query audit log: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"entries": entries,
		"count":   len(entries),
		"limit":   opts.Limit,
		"offset":  opts.Offset,
	})
}

// handleExport handles GET /audit/export - export audit logs
func (h *Handler) handleExport(w http.ResponseWriter, r *http.Request) {
	opts, err := parseExportOptions(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ValidateExportOptions(&opts); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exportSvc := NewExportService(h.repo)
	data, err := exportSvc.Export(opts)
	if err != nil {
		http.Error(w, "failed to export audit log: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set appropriate content type
	contentType := "application/json"
	extension := ".json"
	switch opts.Format {
	case "csv":
		contentType = "text/csv"
		extension = ".csv"
	case "jsonl":
		contentType = "application/x-ndjson"
		extension = ".jsonl"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename=audit-log"+extension)
	w.Write(data)
}

// handleCount handles GET /audit/count - get count of matching entries
func (h *Handler) handleCount(w http.ResponseWriter, r *http.Request) {
	opts, err := parseQueryOptions(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	count, err := h.repo.Count(opts)
	if err != nil {
		http.Error(w, "failed to count audit log: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": count,
	})
}

// handleDelete handles DELETE /audit - delete old entries
func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	// Get duration from query parameter (default: 30 days)
	durationStr := r.URL.Query().Get("older_than")
	var duration time.Duration

	if durationStr == "" {
		duration = 30 * 24 * time.Hour // 30 days
	} else {
		d, err := parseDuration(durationStr)
		if err != nil {
			http.Error(w, "invalid duration: "+err.Error(), http.StatusBadRequest)
			return
		}
		duration = d
	}

	count, err := h.repo.DeleteOlderThan(duration)
	if err != nil {
		http.Error(w, "failed to delete audit entries: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"deleted_count": count,
		"older_than":    duration.String(),
	})
}

// parseQueryOptions parses query parameters into QueryOptions
func parseQueryOptions(r *http.Request) (QueryOptions, error) {
	opts := QueryOptions{
		Limit:  100,
		Offset: 0,
	}

	// Parse time range
	if startStr := r.URL.Query().Get("start_time"); startStr != "" {
		t, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			return opts, err
		}
		opts.StartTime = &t
	}

	if endStr := r.URL.Query().Get("end_time"); endStr != "" {
		t, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			return opts, err
		}
		opts.EndTime = &t
	}

	// Parse filters
	opts.SessionID = r.URL.Query().Get("session_id")
	opts.Surface = r.URL.Query().Get("surface")
	opts.Tool = r.URL.Query().Get("tool")
	opts.UserID = r.URL.Query().Get("user_id")

	if actionStr := r.URL.Query().Get("action"); actionStr != "" {
		opts.Action = AuditAction(actionStr)
	}

	if successStr := r.URL.Query().Get("success"); successStr != "" {
		success := successStr == "true" || successStr == "1"
		opts.Success = &success
	}

	// Parse pagination
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			opts.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			opts.Offset = offset
		}
	}

	// Parse ordering
	opts.OrderBy = r.URL.Query().Get("order_by")
	opts.OrderDir = r.URL.Query().Get("order_dir")

	return opts, nil
}

// parseExportOptions parses query parameters into ExportOptions
func parseExportOptions(r *http.Request) (ExportOptions, error) {
	queryOpts, err := parseQueryOptions(r)
	if err != nil {
		return ExportOptions{}, err
	}

	opts := ExportOptions{
		Format:       r.URL.Query().Get("format"),
		QueryOptions: queryOpts,
		Sanitize:     r.URL.Query().Get("sanitize") == "true",
		IncludeCost:  r.URL.Query().Get("include_cost") == "true",
	}

	return opts, nil
}

// parseDuration parses a duration string (e.g., "24h", "7d", "30d")
func parseDuration(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)

	var duration time.Duration
	multiplier := time.Hour

	if strings.HasSuffix(s, "d") {
		multiplier = 24 * time.Hour
		s = strings.TrimSuffix(s, "d")
	} else if strings.HasSuffix(s, "h") {
		multiplier = time.Hour
		s = strings.TrimSuffix(s, "h")
	} else if strings.HasSuffix(s, "m") {
		multiplier = time.Minute
		s = strings.TrimSuffix(s, "m")
	}

	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	duration = time.Duration(value) * multiplier
	return duration, nil
}

// LogMiddleware creates middleware that logs requests
func LogMiddleware(handler *Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Create audit entry for the request
			entry := &AuditEntry{
				Action:    ActionUserAction,
				Surface:   "api",
				Success:   true,
				Timestamp: time.Now().UTC(),
			}

			// Store entry in context for later logging
			ctx = context.WithValue(ctx, "audit_entry", entry)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
