package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"pryx-core/internal/store"
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

	// Validate pairing code from database
	session, err := s.store.GetPairingSessionByCode(code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "database error"})
		return
	}
	if session == nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(handleMeshPairResponse{
			Success: false,
			Message: "invalid or expired pairing code",
		})
		return
	}

	// Check expiration
	if time.Now().After(session.ExpiresAt) {
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
	if err := s.store.UpdatePairingSessionStatus(code, "approved"); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "failed to update pairing status"})
		return
	}

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

	// Generate cryptographic nonce for secure key exchange
	nonce := generatePairingCode() // In production, use crypto rand

	// Create pairing session in database (expires in 5 minutes)
	session := &store.MeshPairingSession{
		ID:         uuid.New().String(),
		Code:       code,
		DeviceID:   deviceID,
		DeviceName: deviceName,
		ServerURL:  s.cfg.CloudAPIUrl,
		Nonce:      nonce,
		Status:     "pending",
		ExpiresAt:  time.Now().Add(5 * time.Minute),
		CreatedAt:  time.Now(),
	}
	if err := s.store.CreatePairingSession(session); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "failed to create pairing session"})
		return
	}

	// Create QR code data with pairing information
	qrData := map[string]string{
		"device_id":  deviceID,
		"code":       code,
		"server_url": s.cfg.CloudAPIUrl,
		"nonce":      nonce,
	}
	qrJSON, _ := json.Marshal(qrData)

	// Generate QR code image with pairing data
	qrImage, err := qrcode.New(string(qrJSON), qrcode.Medium)
	if err != nil {
		// Fallback to JSON data if QR code generation fails
		qrBase64 := fmt.Sprintf("data:application/json;base64,%s", qrJSON)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(handleMeshQRCodeResponse{
			QRCode:    qrBase64,
			Code:      code,
			ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
		})
		return
	}

	// Convert QR image to PNG bytes
	qrPngBytes, err := qrImage.PNG(256)
	if err != nil {
		// Fallback to JSON data if PNG conversion fails
		qrBase64 := fmt.Sprintf("data:application/json;base64,%s", qrJSON)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(handleMeshQRCodeResponse{
			QRCode:    qrBase64,
			Code:      code,
			ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
		})
		return
	}

	// Encode as base64 data URL for display
	qrBase64 := fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(qrPngBytes))

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(handleMeshQRCodeResponse{
		QRCode:    qrBase64,
		Code:      code,
		ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
	})
}

// handleMeshDevicesList lists paired devices
func (s *Server) handleMeshDevicesList(w http.ResponseWriter, r *http.Request) {
	devices, err := s.store.ListMeshDevices()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "failed to list devices"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"devices": devices,
	})
}

// handleMeshEventsList lists mesh sync events
func (s *Server) handleMeshEventsList(w http.ResponseWriter, r *http.Request) {
	events, err := s.store.ListMeshSyncEvents(100)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "failed to list events"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"events": events,
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

	if err := s.store.DeactivateMeshDevice(deviceID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "failed to unpair device"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "device unpaired successfully",
	})
}
