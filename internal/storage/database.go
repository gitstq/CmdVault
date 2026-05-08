// Package storage provides database operations for CmdVault
package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Command represents a stored command entry
type Command struct {
	ID          int64     `json:"id"`
	Command     string    `json:"command"`
	Shell       string    `json:"shell"`
	SessionID   string    `json:"session_id"`
	WorkingDir  string    `json:"working_dir"`
	ExitCode    int       `json:"exit_code"`
	Duration    int64     `json:"duration"` // milliseconds
	Tags        string    `json:"tags"`
	Notes       string    `json:"notes"`
	Favorite    bool      `json:"favorite"`
	ExecCount   int       `json:"exec_count"`
	LastUsed    time.Time `json:"last_used"`
	CreatedAt   time.Time `json:"created_at"`
}

// Database wraps the SQLite connection
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dbPath string) (*Database, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrency
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &Database{db: db}, nil
}

// createTables creates the necessary database tables
func createTables(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS commands (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		command TEXT NOT NULL,
		shell TEXT DEFAULT 'bash',
		session_id TEXT,
		working_dir TEXT,
		exit_code INTEGER DEFAULT 0,
		duration INTEGER DEFAULT 0,
		tags TEXT,
		notes TEXT,
		favorite INTEGER DEFAULT 0,
		exec_count INTEGER DEFAULT 1,
		last_used TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_commands_command ON commands(command);
	CREATE INDEX IF NOT EXISTS idx_commands_last_used ON commands(last_used DESC);
	CREATE INDEX IF NOT EXISTS idx_commands_exec_count ON commands(exec_count DESC);
	CREATE INDEX IF NOT EXISTS idx_commands_favorite ON commands(favorite);

	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		color TEXT DEFAULT '#3498db',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		shell TEXT,
		started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		ended_at TIMESTAMP
	);
	`

	_, err := db.Exec(schema)
	return err
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// AddCommand adds a new command to the database
func (d *Database) AddCommand(cmd *Command) error {
	// Check if command already exists
	var existingID int64
	var existingCount int
	err := d.db.QueryRow(
		"SELECT id, exec_count FROM commands WHERE command = ?",
		cmd.Command,
	).Scan(&existingID, &existingCount)

	if err == sql.ErrNoRows {
		// Insert new command
		_, err = d.db.Exec(`
			INSERT INTO commands (command, shell, session_id, working_dir, exit_code, duration, tags, notes)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, cmd.Command, cmd.Shell, cmd.SessionID, cmd.WorkingDir, cmd.ExitCode, cmd.Duration, cmd.Tags, cmd.Notes)
		return err
	} else if err != nil {
		return err
	}

	// Update existing command
	_, err = d.db.Exec(`
		UPDATE commands 
		SET exec_count = ?, last_used = CURRENT_TIMESTAMP, working_dir = ?
		WHERE id = ?
	`, existingCount+1, cmd.WorkingDir, existingID)
	return err
}

// GetCommands retrieves commands with optional filtering
func (d *Database) GetCommands(limit int, offset int, favoriteOnly bool) ([]Command, error) {
	query := "SELECT id, command, shell, session_id, working_dir, exit_code, duration, tags, notes, favorite, exec_count, last_used, created_at FROM commands"
	
	if favoriteOnly {
		query += " WHERE favorite = 1"
	}
	query += " ORDER BY last_used DESC LIMIT ? OFFSET ?"

	rows, err := d.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var cmd Command
		err := rows.Scan(
			&cmd.ID, &cmd.Command, &cmd.Shell, &cmd.SessionID, &cmd.WorkingDir,
			&cmd.ExitCode, &cmd.Duration, &cmd.Tags, &cmd.Notes, &cmd.Favorite,
			&cmd.ExecCount, &cmd.LastUsed, &cmd.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}

	return commands, nil
}

// SearchCommands searches commands by keyword
func (d *Database) SearchCommands(keyword string, limit int) ([]Command, error) {
	query := `
		SELECT id, command, shell, session_id, working_dir, exit_code, duration, tags, notes, favorite, exec_count, last_used, created_at
		FROM commands 
		WHERE command LIKE ? OR tags LIKE ? OR notes LIKE ?
		ORDER BY exec_count DESC, last_used DESC
		LIMIT ?
	`

	searchPattern := "%" + keyword + "%"
	rows, err := d.db.Query(query, searchPattern, searchPattern, searchPattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var cmd Command
		err := rows.Scan(
			&cmd.ID, &cmd.Command, &cmd.Shell, &cmd.SessionID, &cmd.WorkingDir,
			&cmd.ExitCode, &cmd.Duration, &cmd.Tags, &cmd.Notes, &cmd.Favorite,
			&cmd.ExecCount, &cmd.LastUsed, &cmd.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}

	return commands, nil
}

// ToggleFavorite toggles the favorite status of a command
func (d *Database) ToggleFavorite(id int64) error {
	_, err := d.db.Exec(`
		UPDATE commands SET favorite = NOT favorite WHERE id = ?
	`, id)
	return err
}

// DeleteCommand deletes a command by ID
func (d *Database) DeleteCommand(id int64) error {
	_, err := d.db.Exec("DELETE FROM commands WHERE id = ?", id)
	return err
}

// UpdateNotes updates the notes for a command
func (d *Database) UpdateNotes(id int64, notes string) error {
	_, err := d.db.Exec("UPDATE commands SET notes = ? WHERE id = ?", notes, id)
	return err
}

// UpdateTags updates the tags for a command
func (d *Database) UpdateTags(id int64, tags string) error {
	_, err := d.db.Exec("UPDATE commands SET tags = ? WHERE id = ?", tags, id)
	return err
}

// GetStats returns statistics about stored commands
func (d *Database) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total commands
	var total int
	err := d.db.QueryRow("SELECT COUNT(*) FROM commands").Scan(&total)
	if err != nil {
		return nil, err
	}
	stats["total_commands"] = total

	// Unique commands
	var unique int
	err = d.db.QueryRow("SELECT COUNT(DISTINCT command) FROM commands").Scan(&unique)
	if err != nil {
		return nil, err
	}
	stats["unique_commands"] = unique

	// Favorites
	var favorites int
	err = d.db.QueryRow("SELECT COUNT(*) FROM commands WHERE favorite = 1").Scan(&favorites)
	if err != nil {
		return nil, err
	}
	stats["favorites"] = favorites

	// Top commands
	rows, err := d.db.Query(`
		SELECT command, exec_count FROM commands 
		ORDER BY exec_count DESC LIMIT 10
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topCommands []map[string]interface{}
	for rows.Next() {
		var cmd string
		var count int
		if err := rows.Scan(&cmd, &count); err != nil {
			return nil, err
		}
		topCommands = append(topCommands, map[string]interface{}{
			"command":    cmd,
			"exec_count": count,
		})
	}
	stats["top_commands"] = topCommands

	return stats, nil
}

// ImportHistory imports commands from a shell history file
func (d *Database) ImportHistory(historyPath, shell string) (int, error) {
	data, err := os.ReadFile(historyPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read history file: %w", err)
	}

	// Parse history based on shell type
	var commands []string
	lines := strings.Split(string(data), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Skip timestamp lines for bash
		if strings.HasPrefix(line, "#") {
			continue
		}
		commands = append(commands, line)
	}

	// Insert commands
	count := 0
	for _, cmd := range commands {
		if err := d.AddCommand(&Command{
			Command: cmd,
			Shell:   shell,
		}); err == nil {
			count++
		}
	}

	return count, nil
}
