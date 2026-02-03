package server

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

// handleMeshPairRequest represents a pairing request from TUI
type handleMeshPairRequest struct {
	Code string `json:"code"`
}

// handleMeshPairResponse represents the pairing response
type handleMeshPairResponse struct {
	Success    bool   `json:"success"`
	DeviceID   string `json:"device_id,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
	Message    string `json:"message,omitempty"`
}

// handleMeshQRCodeResponse represents QR code generation response
type handleMeshQRCodeResponse struct {
	QRCode    string `json:"qr_code"`    // Base64-encoded QR code
	Code      string `json:"code"`       // 6-digit pairing code
	ExpiresAt string `json:"expires_at"` // Expiration time
}

// PairingSession represents an active pairing session
type PairingSession struct {
	Code       string    `json:"code"`
	DeviceID   string    `json:"device_id"`
	DeviceName string    `json:"device_name"`
	ExpiresAt  time.Time `json:"expires_at"`
	Status     string    `json:"status"` // pending, approved, rejected, expired
}

// In-memory pairing session storage (in production, use Redis or D1)
var pairingSessions = make(map[string]*PairingSession)

// generatePairingCode generates a random 6-digit code
func generatePairingCode() string {
	const charset = "0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// generateDeviceID generates a device ID (same format as mesh.generateDeviceID)
func generateDeviceID() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return fmt.Sprintf("pryx-%s", string(b))
}

// handleMeshPair handles device pairing via 6-digit code
func (s *Server) handleMeshPair(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "method not allowed"})
		return
	}

	var req handleMeshPairRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "invalid request body"})
		return
	}

	code := strings.TrimSpace(req.Code)
	if len(code) != 6 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(handleMeshPairResponse{
			Success: false,
			Message: "pairing code must be exactly 6 digits",
		})
		return
	}

	// Validate pairing code
	session, exists := pairingSessions[code]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(handleMeshPairResponse{
			Success: false,
			Message: "invalid or expired pairing code",
		})
		return
	}

	// Check expiration
	if time.Now().After(session.ExpiresAt) {
		delete(pairingSessions, code)
		w.WriteHeader(http.StatusGone)
		_ = json.NewEncoder(w).Encode(handleMeshPairResponse{
			Success: false,
			Message: "pairing code has expired",
		})
		return
	}

	// Check status
	if session.Status != "pending" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(handleMeshPairResponse{
			Success: false,
			Message: fmt.Sprintf("pairing code status is %s", session.Status),
		})
		return
	}

	// Approve pairing
	session.Status = "approved"

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(handleMeshPairResponse{
		Success:    true,
		DeviceID:   session.DeviceID,
		DeviceName: session.DeviceName,
		Message:    "device paired successfully",
	})
}

// handleMeshQRCode generates a new pairing QR code and 6-digit code
func (s *Server) handleMeshQRCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "method not allowed"})
		return
	}

	// Generate pairing code
	code := generatePairingCode()

	// Get device info
	deviceID, _ := s.keychain.Get("device_id")
	if deviceID == "" {
		deviceID = generateDeviceID()
		s.keychain.Set("device_id", deviceID)
	}

	deviceName := "New Device"
	if name, err := s.keychain.Get("device_name"); err == nil && name != "" {
		deviceName = name
	}

	// Create pairing session (expires in 5 minutes)
	session := &PairingSession{
		Code:       code,
		DeviceID:   deviceID,
		DeviceName: deviceName,
		ExpiresAt:  time.Now().Add(5 * time.Minute),
		Status:     "pending",
	}
	pairingSessions[code] = session

	// In production, generate actual QR code containing:
	// - Device ID
	// - Pairing code
	// - Server URL
	// - Cryptographic nonce for secure key exchange
	//
	// For now, return a placeholder that includes the pairing data
	qrData := map[string]string{
		"device_id":  deviceID,
		"code":       code,
		"server_url": s.cfg.CloudAPIUrl,
		"nonce":      generatePairingCode(), // Placeholder for cryptographic nonce
	}
	qrJSON, _ := json.Marshal(qrData)
	qrBase64 := fmt.Sprintf("data:application/json;base64,%s", qrJSON)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(handleMeshQRCodeResponse{
		QRCode:    qrBase64,
		Code:      code,
		ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
	})
}

// handleMeshDevicesList lists paired devices
func (s *Server) handleMeshDevicesList(w http.ResponseWriter, r *http.Request) {
	// In production, query from store or D1 database
	// For now, return empty list
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"devices": []interface{}{},
	})
}

// handleMeshEventsList lists mesh sync events
func (s *Server) handleMeshEventsList(w http.ResponseWriter, r *http.Request) {
	// In production, query from store or D1 database
	// For now, return empty list
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"events": []interface{}{},
	})
}

// handleMeshDevicesUnpair removes a paired device
func (s *Server) handleMeshDevicesUnpair(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")
	if deviceID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "device ID required"})
		return
	}

	// In production, delete from store or D1 database
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "device unpaired successfully",
	})
}
