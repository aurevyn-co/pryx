package server

import (
	"net/http"
	"strings"

	"pryx-core/internal/config"
)

// corsMiddleware creates a CORS middleware with configurable allowed origins.
// It validates the Origin header against a whitelist and only allows credentials
// for explicitly allowed origins (never with wildcard *).
func corsMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			if !isOriginAllowed(origin, cfg.AllowedOrigins) {
				// For preflight requests, reject with 403
				if r.Method == "OPTIONS" {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				// For actual requests without origin, proceed without CORS headers
				// For disallowed origins, still proceed but don't set CORS headers
				next.ServeHTTP(w, r)
				return
			}

			// Origin is allowed - set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// isOriginAllowed checks if the given origin is in the allowed list.
// For development, localhost origins are always allowed.
// Empty origin (e.g., curl requests) is allowed for API compatibility.
func isOriginAllowed(origin string, allowed []string) bool {
	// Empty origin is allowed (e.g., curl, server-to-server)
	if origin == "" {
		return true
	}

	// Always allow localhost origins for development
	if strings.HasPrefix(origin, "http://localhost:") ||
		strings.HasPrefix(origin, "https://localhost:") ||
		origin == "http://localhost" ||
		origin == "https://localhost" {
		return true
	}

	// Check against configured allowed origins
	for _, allowedOrigin := range allowed {
		if strings.EqualFold(origin, allowedOrigin) {
			return true
		}
	}

	return false
}
