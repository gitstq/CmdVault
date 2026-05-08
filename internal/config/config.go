// Package config provides configuration management for CmdVault
package config

import (
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	// DatabasePath is the path to the SQLite database file
	DatabasePath string
	// HistoryFile is the path to the shell history file
	HistoryFile string
	// MaxHistoryItems is the maximum number of history items to store
	MaxHistoryItems int
	// EnableAI enables AI-powered semantic search
	EnableAI bool
	// SyncEnabled enables cross-terminal sync
	SyncEnabled bool
	// SyncPort is the port for sync service
	SyncPort int
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".cmdvault")
	
	return &Config{
		DatabasePath:    filepath.Join(configDir, "history.db"),
		HistoryFile:     filepath.Join(homeDir, ".bash_history"),
		MaxHistoryItems: 100000,
		EnableAI:        false,
		SyncEnabled:     false,
		SyncPort:        9876,
	}
}

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() error {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".cmdvault")
	return os.MkdirAll(configDir, 0755)
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".cmdvault", "config.yaml")
}
