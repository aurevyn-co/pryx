package vault

import (
	"testing"
	"time"
)

func TestNewMasterKeyManager(t *testing.T) {
	m := NewMasterKeyManager()

	if m == nil {
		t.Fatal("NewMasterKeyManager() returned nil")
	}

	if !m.IsLocked() {
		t.Error("New vault should be locked")
	}

	if m.ttl != defaultKeyCacheTTL {
		t.Errorf("Expected default TTL %v, got %v", defaultKeyCacheTTL, m.ttl)
	}
}

func TestMasterKeyManager_Unlock(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid password unlocks vault",
			password: "my-secure-password-123",
			wantErr:  false,
		},
		{
			name:        "short password fails",
			password:    "short",
			wantErr:     true,
			errContains: "at least 12 characters",
		},
		{
			name:        "empty password fails",
			password:    "",
			wantErr:     true,
			errContains: "at least 12 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMasterKeyManager()
			err := m.Unlock(tt.password)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tt.errContains != "" {
					if !contains(err.Error(), tt.errContains) {
						t.Errorf("Error %q does not contain %q", err.Error(), tt.errContains)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if m.IsLocked() {
					t.Error("Vault should be unlocked")
				}
			}
		})
	}
}

func TestMasterKeyManager_Lock(t *testing.T) {
	m := NewMasterKeyManager()

	// Try to lock when already locked
	err := m.Lock()
	if err != ErrVaultAlreadyLocked {
		t.Errorf("Expected ErrVaultAlreadyLocked, got %v", err)
	}

	// Unlock first
	if err := m.Unlock("my-secure-password-123"); err != nil {
		t.Fatalf("Failed to unlock: %v", err)
	}

	// Now lock
	if err := m.Lock(); err != nil {
		t.Errorf("Failed to lock: %v", err)
	}

	if !m.IsLocked() {
		t.Error("Vault should be locked")
	}

	// Key should be cleared
	_, err = m.GetKey()
	if err != ErrVaultLocked {
		t.Errorf("Expected ErrVaultLocked when getting key after lock, got %v", err)
	}
}

func TestMasterKeyManager_GetKey(t *testing.T) {
	m := NewMasterKeyManager()

	// Should fail when locked
	_, err := m.GetKey()
	if err != ErrVaultLocked {
		t.Errorf("Expected ErrVaultLocked, got %v", err)
	}

	// Unlock
	if err := m.Unlock("my-secure-password-123"); err != nil {
		t.Fatalf("Failed to unlock: %v", err)
	}

	// Should succeed now
	key, err := m.GetKey()
	if err != nil {
		t.Errorf("Failed to get key: %v", err)
	}
	if len(key) != argon2KeyLen {
		t.Errorf("Expected key length %d, got %d", argon2KeyLen, len(key))
	}
}

func TestMasterKeyManager_KeyExpiration(t *testing.T) {
	m := NewMasterKeyManager()
	m.SetTTL(100 * time.Millisecond)

	if err := m.Unlock("my-secure-password-123"); err != nil {
		t.Fatalf("Failed to unlock: %v", err)
	}

	// Should work immediately
	_, err := m.GetKey()
	if err != nil {
		t.Errorf("Failed to get key: %v", err)
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should fail after expiration
	_, err = m.GetKey()
	if err != ErrVaultLocked {
		t.Errorf("Expected ErrVaultLocked after expiration, got %v", err)
	}
}

func TestMasterKeyManager_VerifyPassword(t *testing.T) {
	m := NewMasterKeyManager()
	password := "my-secure-password-123"

	// Should fail when vault never unlocked
	if m.VerifyPassword(password) {
		t.Error("VerifyPassword should fail when vault never unlocked")
	}

	// Unlock
	if err := m.Unlock(password); err != nil {
		t.Fatalf("Failed to unlock: %v", err)
	}

	// Should verify correct password
	if !m.VerifyPassword(password) {
		t.Error("VerifyPassword should succeed with correct password")
	}

	// Should fail with wrong password
	if m.VerifyPassword("wrong-password-123") {
		t.Error("VerifyPassword should fail with wrong password")
	}
}

func TestMasterKeyManager_RotateKey(t *testing.T) {
	m := NewMasterKeyManager()
	oldPassword := "my-secure-password-123"
	newPassword := "my-new-password-456"

	if err := m.Unlock(oldPassword); err != nil {
		t.Fatalf("Failed to unlock: %v", err)
	}

	// Get old key
	oldKey, _ := m.GetKey()

	// Rotate
	if err := m.RotateKey(oldPassword, newPassword); err != nil {
		t.Errorf("Failed to rotate key: %v", err)
	}

	// Get new key
	newKey, _ := m.GetKey()

	// Keys should be different
	if string(oldKey) == string(newKey) {
		t.Error("Key should change after rotation")
	}

	// Old password should not work anymore
	if m.VerifyPassword(oldPassword) {
		t.Error("Old password should not work after rotation")
	}

	// New password should work
	if !m.VerifyPassword(newPassword) {
		t.Error("New password should work after rotation")
	}
}

func TestMasterKeyManager_RateLimiting(t *testing.T) {
	m := NewMasterKeyManager()

	// Try multiple failed attempts
	for i := 0; i < 5; i++ {
		m.Unlock("short") // Will fail due to password length
	}

	err := m.Unlock("my-secure-password-123")
	if err == nil {
		t.Error("Should be rate limited after multiple failures")
	}
	if !contains(err.Error(), "too many failed attempts") {
		t.Errorf("Expected rate limit error, got %v", err)
	}
}

func TestMasterKeyManager_GetStatus(t *testing.T) {
	m := NewMasterKeyManager()

	status := m.GetStatus()
	if status["locked"] != true {
		t.Error("Status should show locked when new")
	}

	if err := m.Unlock("my-secure-password-123"); err != nil {
		t.Fatalf("Failed to unlock: %v", err)
	}

	status = m.GetStatus()
	if status["locked"] != false {
		t.Error("Status should show unlocked")
	}

	if _, ok := status["key_created_at"]; !ok {
		t.Error("Status should include key_created_at when unlocked")
	}

	if _, ok := status["key_expires_at"]; !ok {
		t.Error("Status should include key_expires_at when unlocked")
	}
}

func TestMasterKeyManager_Salt(t *testing.T) {
	m := NewMasterKeyManager()

	// Should fail when salt not set
	_, err := m.GetSalt()
	if err == nil {
		t.Error("GetSalt should fail when salt not initialized")
	}

	// Set salt
	testSalt := []byte("test-salt-123456")
	if err := m.SetSalt(testSalt); err != nil {
		t.Errorf("Failed to set salt: %v", err)
	}

	// Get salt
	salt, err := m.GetSalt()
	if err != nil {
		t.Errorf("Failed to get salt: %v", err)
	}

	if string(salt) != string(testSalt) {
		t.Error("Retrieved salt doesn't match")
	}

	// Should fail with empty salt
	if err := m.SetSalt([]byte{}); err == nil {
		t.Error("SetSalt should fail with empty salt")
	}
}

func TestEncodeDecodeSalt(t *testing.T) {
	salt := []byte("test-salt-data")

	encoded := EncodeSalt(salt)
	if encoded == "" {
		t.Error("EncodeSalt returned empty string")
	}

	decoded, err := DecodeSalt(encoded)
	if err != nil {
		t.Errorf("DecodeSalt failed: %v", err)
	}

	if string(decoded) != string(salt) {
		t.Error("Decoded salt doesn't match original")
	}

	// Should fail with invalid base64
	_, err = DecodeSalt("not-valid-base64!!!")
	if err == nil {
		t.Error("DecodeSalt should fail with invalid input")
	}
}

func TestGenerateSalt(t *testing.T) {
	salt1, err := generateSalt()
	if err != nil {
		t.Errorf("generateSalt failed: %v", err)
	}

	if len(salt1) != 16 {
		t.Errorf("Expected salt length 16, got %d", len(salt1))
	}

	salt2, err := generateSalt()
	if err != nil {
		t.Errorf("generateSalt failed: %v", err)
	}

	// Should generate unique salts
	if string(salt1) == string(salt2) {
		t.Error("generateSalt should generate unique salts")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
