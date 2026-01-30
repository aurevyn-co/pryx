package vault

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/crypto/argon2"
)

const (
	argon2Time    = 3
	argon2Memory  = 64 * 1024
	argon2Threads = 4
	argon2KeyLen  = 32

	minPasswordLength = 12

	defaultKeyCacheTTL = 15 * time.Minute
)

var (
	ErrVaultLocked        = errors.New("vault is locked")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrPasswordTooShort   = errors.New("password must be at least 12 characters")
	ErrVaultAlreadyLocked = errors.New("vault is already locked")
	ErrRateLimited        = errors.New("too many failed attempts, please wait")
)

type DerivedKey struct {
	Key       []byte
	Salt      []byte
	CreatedAt time.Time
	ExpiresAt time.Time
}

type MasterKeyManager struct {
	mu sync.RWMutex

	derivedKey *DerivedKey
	salt       []byte
	ttl        time.Duration

	failedAttempts int
	lastAttempt    time.Time
	lockoutUntil   time.Time

	isLocked bool
}

func NewMasterKeyManager() *MasterKeyManager {
	return &MasterKeyManager{
		ttl:      defaultKeyCacheTTL,
		isLocked: true,
	}
}

func (m *MasterKeyManager) SetTTL(ttl time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ttl = ttl
}

func (m *MasterKeyManager) Unlock(password string) error {
	if err := m.checkRateLimit(); err != nil {
		return err
	}

	if err := m.validatePassword(password); err != nil {
		m.recordFailedAttempt()
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.salt == nil {
		salt, err := generateSalt()
		if err != nil {
			return fmt.Errorf("failed to generate salt: %w", err)
		}
		m.salt = salt
	}

	key := argon2.IDKey(
		[]byte(password),
		m.salt,
		argon2Time,
		argon2Memory,
		argon2Threads,
		argon2KeyLen,
	)

	now := time.Now()
	m.derivedKey = &DerivedKey{
		Key:       key,
		Salt:      m.salt,
		CreatedAt: now,
		ExpiresAt: now.Add(m.ttl),
	}
	m.isLocked = false
	m.failedAttempts = 0

	return nil
}

// Lock immediately locks the vault and clears the derived key from memory
func (m *MasterKeyManager) Lock() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isLocked {
		return ErrVaultAlreadyLocked
	}

	// Securely clear the key from memory
	if m.derivedKey != nil && m.derivedKey.Key != nil {
		clear(m.derivedKey.Key)
	}

	m.derivedKey = nil
	m.isLocked = true

	return nil
}

// IsLocked returns whether the vault is currently locked
func (m *MasterKeyManager) IsLocked() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isLocked
}

// GetKey returns the derived key if the vault is unlocked and key is not expired
func (m *MasterKeyManager) GetKey() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.isLocked {
		return nil, ErrVaultLocked
	}

	if m.derivedKey == nil {
		return nil, ErrVaultLocked
	}

	// Check if key has expired
	if time.Now().After(m.derivedKey.ExpiresAt) {
		return nil, ErrVaultLocked
	}

	// Return a copy of the key
	keyCopy := make([]byte, len(m.derivedKey.Key))
	copy(keyCopy, m.derivedKey.Key)

	return keyCopy, nil
}

// GetSalt returns the salt used for key derivation
func (m *MasterKeyManager) GetSalt() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.salt == nil {
		return nil, errors.New("salt not initialized")
	}

	saltCopy := make([]byte, len(m.salt))
	copy(saltCopy, m.salt)
	return saltCopy, nil
}

// SetSalt sets the salt (used when loading from storage)
func (m *MasterKeyManager) SetSalt(salt []byte) error {
	if len(salt) == 0 {
		return errors.New("salt cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.salt = make([]byte, len(salt))
	copy(m.salt, salt)

	return nil
}

// RotateKey re-derives the key with a new password
func (m *MasterKeyManager) RotateKey(oldPassword, newPassword string) error {
	if err := m.validatePassword(newPassword); err != nil {
		return err
	}

	// Verify old password first
	if err := m.Unlock(oldPassword); err != nil {
		return fmt.Errorf("invalid old password: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate new salt
	newSalt, err := generateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate new salt: %w", err)
	}

	// Derive new key
	newKey := argon2.IDKey(
		[]byte(newPassword),
		newSalt,
		argon2Time,
		argon2Memory,
		argon2Threads,
		argon2KeyLen,
	)

	// Clear old key
	if m.derivedKey != nil && m.derivedKey.Key != nil {
		clear(m.derivedKey.Key)
	}

	// Update with new key and salt
	now := time.Now()
	m.derivedKey = &DerivedKey{
		Key:       newKey,
		Salt:      newSalt,
		CreatedAt: now,
		ExpiresAt: now.Add(m.ttl),
	}
	m.salt = newSalt

	return nil
}

// VerifyPassword checks if the provided password matches without unlocking
func (m *MasterKeyManager) VerifyPassword(password string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.salt == nil {
		return false
	}

	// Derive key with provided password
	testKey := argon2.IDKey(
		[]byte(password),
		m.salt,
		argon2Time,
		argon2Memory,
		argon2Threads,
		argon2KeyLen,
	)

	// Compare with current key if available
	if m.derivedKey != nil && m.derivedKey.Key != nil {
		return subtle.ConstantTimeCompare(testKey, m.derivedKey.Key) == 1
	}

	return false
}

// GetStatus returns the current vault status
func (m *MasterKeyManager) GetStatus() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := map[string]interface{}{
		"locked":          m.isLocked,
		"failed_attempts": m.failedAttempts,
	}

	if !m.isLocked && m.derivedKey != nil {
		status["key_created_at"] = m.derivedKey.CreatedAt
		status["key_expires_at"] = m.derivedKey.ExpiresAt
		status["time_remaining"] = time.Until(m.derivedKey.ExpiresAt).String()
	}

	if !m.lockoutUntil.IsZero() && time.Now().Before(m.lockoutUntil) {
		status["locked_out"] = true
		status["lockout_remaining"] = time.Until(m.lockoutUntil).String()
	}

	return status
}

// validatePassword checks password strength
func (m *MasterKeyManager) validatePassword(password string) error {
	if len(password) < minPasswordLength {
		return ErrPasswordTooShort
	}
	return nil
}

// checkRateLimit enforces rate limiting on unlock attempts
func (m *MasterKeyManager) checkRateLimit() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.lockoutUntil.IsZero() && time.Now().Before(m.lockoutUntil) {
		return fmt.Errorf("%w: %v", ErrRateLimited, time.Until(m.lockoutUntil))
	}

	return nil
}

// recordFailedAttempt tracks failed unlock attempts
func (m *MasterKeyManager) recordFailedAttempt() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.failedAttempts++
	m.lastAttempt = time.Now()

	// Exponential backoff: 1s, 2s, 4s, 8s, 16s, 32s, 64s, 128s, 256s, 512s
	if m.failedAttempts > 0 {
		delay := time.Duration(1<<uint(m.failedAttempts-1)) * time.Second
		if delay > 5*time.Minute {
			delay = 5 * time.Minute
		}
		m.lockoutUntil = time.Now().Add(delay)
	}
}

// generateSalt creates a random salt
func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// EncodeSalt encodes salt to base64 string
func EncodeSalt(salt []byte) string {
	return base64.StdEncoding.EncodeToString(salt)
}

// DecodeSalt decodes salt from base64 string
func DecodeSalt(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}
