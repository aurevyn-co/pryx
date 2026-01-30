// Package keychain provides secure credential storage using the system keyring.
// It abstracts OS-specific keychain/keyring implementations for storing sensitive data like API keys.
package keychain

import (
	"fmt"
	"strings"

	"github.com/zalando/go-keyring"
)

// Keychain provides secure storage for credentials using the system keyring.
// It uses a service name to namespace all stored credentials.
type Keychain struct {
	service string
}

// New creates a new Keychain instance for the specified service.
// The service name is used as a namespace for all stored credentials.
func New(service string) *Keychain {
	return &Keychain{service: service}
}

// Set stores a password for the specified user in the keychain.
// Returns an error if the operation fails.
func (k *Keychain) Set(user, password string) error {
	return keyring.Set(k.service, user, password)
}

// Get retrieves the password for the specified user from the keychain.
// Returns an error if the credential is not found or the operation fails.
func (k *Keychain) Get(user string) (string, error) {
	return keyring.Get(k.service, user)
}

// Delete removes the credential for the specified user from the keychain.
// Returns an error if the operation fails.
func (k *Keychain) Delete(user string) error {
	return keyring.Delete(k.service, user)
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
