// Package keychain provides secure credential storage using the system keyring.
// It abstracts OS-specific keychain/keyring implementations for storing sensitive data like API keys.
// For testing, set PRYX_KEYCHAIN_FILE environment variable to use a file-based keychain instead.
package keychain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/zalando/go-keyring"
)

// Keychain provides secure storage for credentials using the system keyring.
// It uses a service name to namespace all stored credentials.
type Keychain struct {
	service  string
	filePath string
	fileData map[string]string
	fileMu   sync.RWMutex
	useFile  bool
}

// New creates a new Keychain instance for the specified service.
// The service name is used as a namespace for all stored credentials.
// If PRYX_KEYCHAIN_FILE is set, uses file-based storage for testing.
func New(service string) *Keychain {
	k := &Keychain{service: service}

	// Check if we should use file-based keychain (for testing)
	if keychainFile := os.Getenv("PRYX_KEYCHAIN_FILE"); keychainFile != "" {
		k.useFile = true
		k.filePath = keychainFile
		k.fileData = make(map[string]string)
		// Load existing data if file exists
		if data, err := os.ReadFile(keychainFile); err == nil {
			json.Unmarshal(data, &k.fileData)
		}
	}

	return k
}

// Set stores a password for the specified user in the keychain.
// Returns an error if the operation fails.
func (k *Keychain) Set(user, password string) error {
	if k.useFile {
		return k.setFile(user, password)
	}
	return keyring.Set(k.service, user, password)
}

// Get retrieves the password for the specified user from the keychain.
// Returns an error if the credential is not found or the operation fails.
func (k *Keychain) Get(user string) (string, error) {
	if k.useFile {
		return k.getFile(user)
	}
	return keyring.Get(k.service, user)
}

// Delete removes the credential for the specified user from the keychain.
// Returns an error if the operation fails.
func (k *Keychain) Delete(user string) error {
	if k.useFile {
		return k.deleteFile(user)
	}
	return keyring.Delete(k.service, user)
}

// File-based keychain implementation for testing
func (k *Keychain) setFile(user, password string) error {
	k.fileMu.Lock()
	defer k.fileMu.Unlock()

	key := k.service + ":" + user
	k.fileData[key] = password

	// Ensure directory exists
	dir := filepath.Dir(k.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write to file
	data, err := json.Marshal(k.fileData)
	if err != nil {
		return err
	}
	return os.WriteFile(k.filePath, data, 0600)
}

func (k *Keychain) getFile(user string) (string, error) {
	k.fileMu.RLock()
	defer k.fileMu.RUnlock()

	key := k.service + ":" + user
	if password, ok := k.fileData[key]; ok {
		return password, nil
	}
	return "", fmt.Errorf("secret not found in file keychain")
}

func (k *Keychain) deleteFile(user string) error {
	k.fileMu.Lock()
	defer k.fileMu.Unlock()

	key := k.service + ":" + user
	delete(k.fileData, key)

	data, err := json.Marshal(k.fileData)
	if err != nil {
		return err
	}
	return os.WriteFile(k.filePath, data, 0600)
}

// SetProviderKey stores an API key for the specified LLM provider.
// The provider ID is used to construct the key name (e.g., "provider:openai").
func (k *Keychain) SetProviderKey(provider, key string) error {
	keyName := fmt.Sprintf("provider:%s", provider)
	return k.Set(keyName, key)
}

// GetProviderKey retrieves the API key for the specified LLM provider.
// Returns an error if the key is not found.
func (k *Keychain) GetProviderKey(provider string) (string, error) {
	keyName := fmt.Sprintf("provider:%s", provider)
	return k.Get(keyName)
}

// DeleteProviderKey removes the API key for the specified LLM provider.
func (k *Keychain) DeleteProviderKey(provider string) error {
	keyName := fmt.Sprintf("provider:%s", provider)
	return k.Delete(keyName)
}

// ListProviderKeys returns a list of all provider keys stored in the keychain.
// Currently returns an empty list (not fully implemented).
func (k *Keychain) ListProviderKeys() ([]string, error) {
	return []string{}, nil
}

// MigrateConfigKey migrates a provider key from configuration to the keychain.
// If the key is empty, no action is taken.
func (k *Keychain) MigrateConfigKey(provider, key string) error {
	if key == "" {
		return nil
	}
	return k.SetProviderKey(provider, key)
}

// GetKeyForProvider returns the keychain key name for a given provider.
// The format is "provider:<provider_id>".
func GetKeyForProvider(provider string) string {
	return fmt.Sprintf("provider:%s", provider)
}

// ExtractProviderFromKey extracts the provider ID from a keychain key name.
// Returns the provider ID and true if the key has the "provider:" prefix.
func ExtractProviderFromKey(key string) (string, bool) {
	if strings.HasPrefix(key, "provider:") {
		return strings.TrimPrefix(key, "provider:"), true
	}
	return "", false
}
