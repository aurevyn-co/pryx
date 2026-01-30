package server

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

// GetAvailablePort finds an available TCP port on the system.
// It binds to port 0 to let the OS assign an available port, then returns that port number.
func GetAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

// WritePortFile writes the port number to ~/.pryx/runtime.port.
// This allows clients (like the TUI) to discover the runtime's port.
// Creates the .pryx directory if it doesn't exist.
func WritePortFile(port int) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	pryxDir := filepath.Join(homeDir, ".pryx")
	if err := os.MkdirAll(pryxDir, 0755); err != nil {
		return fmt.Errorf("failed to create .pryx directory: %w", err)
	}

	portFile := filepath.Join(pryxDir, "runtime.port")
	content := fmt.Sprintf("%d", port)

	if err := os.WriteFile(portFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write port file: %w", err)
	}

	return nil
}

// ReadPortFile reads the port number from ~/.pryx/runtime.port.
// Returns the port as an integer. Returns an error if the file doesn't exist or contains invalid data.
func ReadPortFile() (int, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return 0, fmt.Errorf("failed to get home directory: %w", err)
	}

	portFile := filepath.Join(homeDir, ".pryx", "runtime.port")
	content, err := os.ReadFile(portFile)
	if err != nil {
		return 0, fmt.Errorf("failed to read port file: %w", err)
	}

	port, err := strconv.Atoi(string(content))
	if err != nil {
		return 0, fmt.Errorf("invalid port in file: %w", err)
	}

	return port, nil
}

// CleanupPortFile removes the port file from ~/.pryx/runtime.port.
// Should be called on shutdown to clean up the port file.
func CleanupPortFile() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	portFile := filepath.Join(homeDir, ".pryx", "runtime.port")
	return os.Remove(portFile)
}
