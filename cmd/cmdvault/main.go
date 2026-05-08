// Package main is the entry point for CmdVault
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cmdvault/cmdvault/internal/config"
	"github.com/cmdvault/cmdvault/internal/storage"
	"github.com/cmdvault/cmdvault/internal/tui"
)

var (
	version = "1.0.0"
	commit  = "none"
	date    = "unknown"
)

// Global database instance
var db *storage.Database

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "cmdvault",
	Short: "🔐 Intelligent command history manager",
	Long: `CmdVault is a smart command history manager that helps you
search, organize, and recall your terminal commands efficiently.

Features:
  🔍 Fuzzy search with intelligent ranking
  ⭐ Favorite and tag commands
  📊 Usage statistics and insights
  🔒 Local-first with privacy protection
  🎨 Beautiful terminal UI`,
	Version: version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize database
		cfg := config.DefaultConfig()
		var err error
		db, err = storage.NewDatabase(cfg.DatabasePath)
		if err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if db != nil {
			return db.Close()
		}
		return nil
	},
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search commands with fuzzy matching",
	Long: `Search through your command history using fuzzy matching.
If no query is provided, opens the interactive TUI.`,
	Aliases: []string{"s", "find"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// Launch TUI
			if err := tui.RunTUI(db); err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}
			return
		}

		query := strings.Join(args, " ")
		limit, _ := cmd.Flags().GetInt("limit")
		
		commands, err := db.SearchCommands(query, limit)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		if len(commands) == 0 {
			fmt.Println("No commands found matching:", query)
			return
		}

		for _, c := range commands {
			fav := ""
			if c.Favorite {
				fav = "⭐ "
			}
			fmt.Printf("%s%s\n", fav, c.Command)
		}
	},
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent commands",
	Long:  `List your most recently used commands.`,
	Aliases: []string{"ls", "recent"},
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		favorites, _ := cmd.Flags().GetBool("favorites")

		commands, err := db.GetCommands(limit, 0, favorites)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		for i, c := range commands {
			fav := ""
			if c.Favorite {
				fav = "⭐ "
			}
			fmt.Printf("%3d. %s%s\n", i+1, fav, c.Command)
		}
	},
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <command>",
	Short: "Add a command to the vault",
	Long:  `Add a new command to your command vault.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		command := strings.Join(args, " ")
		tags, _ := cmd.Flags().GetString("tags")
		notes, _ := cmd.Flags().GetString("notes")

		wd, _ := os.Getwd()
		err := db.AddCommand(&storage.Command{
			Command:    command,
			Shell:      detectShell(),
			WorkingDir: wd,
			Tags:       tags,
			Notes:      notes,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Println("✅ Command added:", command)
	},
}

// favoriteCmd represents the favorite command
var favoriteCmd = &cobra.Command{
	Use:   "favorite <id>",
	Short: "Toggle favorite status",
	Long:  `Toggle the favorite status of a command.`,
	Aliases: []string{"fav", "star"},
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var id int64
		fmt.Sscanf(args[0], "%d", &id)

		if err := db.ToggleFavorite(id); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Println("✅ Favorite status toggled")
	},
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a command",
	Long:  `Delete a command from the vault.`,
	Aliases: []string{"rm", "remove"},
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var id int64
		fmt.Sscanf(args[0], "%d", &id)

		if err := db.DeleteCommand(id); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Println("✅ Command deleted")
	},
}

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show usage statistics",
	Long:  `Display statistics about your command usage.`,
	Run: func(cmd *cobra.Command, args []string) {
		stats, err := db.GetStats()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		fmt.Println("📊 CmdVault Statistics")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("📝 Total Commands: %v\n", stats["total_commands"])
		fmt.Printf("🎯 Unique Commands: %v\n", stats["unique_commands"])
		fmt.Printf("⭐ Favorites: %v\n", stats["favorites"])
		fmt.Println()
		
		if top, ok := stats["top_commands"].([]map[string]interface{}); ok && len(top) > 0 {
			fmt.Println("🏆 Top Commands:")
			for i, t := range top {
				fmt.Printf("  %d. %s (%d times)\n", i+1, t["command"], t["exec_count"])
			}
		}
	},
}

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import [history-file]",
	Short: "Import shell history",
	Long: `Import commands from your shell history file.
If no file is specified, it will try to detect your default history file.`,
	Run: func(cmd *cobra.Command, args []string) {
		shell := detectShell()
		historyPath := ""

		if len(args) > 0 {
			historyPath = args[0]
		} else {
			homeDir, _ := os.UserHomeDir()
			switch shell {
			case "zsh":
				historyPath = filepath.Join(homeDir, ".zsh_history")
			case "fish":
				historyPath = filepath.Join(homeDir, ".local/share/fish/fish_history")
			default:
				historyPath = filepath.Join(homeDir, ".bash_history")
			}
		}

		count, err := db.ImportHistory(historyPath, shell)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Imported %d commands from %s\n", count, historyPath)
	},
}

// initCmd represents the init command for shell integration
var initCmd = &cobra.Command{
	Use:   "init [shell]",
	Short: "Generate shell integration script",
	Long: `Generate shell integration script for automatic command tracking.
Supported shells: bash, zsh, fish`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shell := args[0]
		script := generateShellIntegration(shell)
		fmt.Println(script)
	},
}

// detectShell detects the current shell
func detectShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "bash"
	}
	
	if strings.Contains(shell, "zsh") {
		return "zsh"
	}
	if strings.Contains(shell, "fish") {
		return "fish"
	}
	return "bash"
}

// generateShellIntegration generates shell integration script
func generateShellIntegration(shell string) string {
	binaryPath, _ := exec.LookPath("cmdvault")
	if binaryPath == "" {
		binaryPath = "cmdvault"
	}

	switch shell {
	case "bash":
		return fmt.Sprintf(`# CmdVault Bash Integration
export PROMPT_COMMAND='cmdvault add "$(history 1 | sed "s/^[ ]*[0-9]*[ ]*//")"'
bind '"\C-r": "\C-a cmdvault search \C-m"'
`)
	case "zsh":
		return fmt.Sprintf(`# CmdVault Zsh Integration
precmd() {
    cmdvault add "$(history -1 | sed 's/^[ ]*[0-9]*[ ]*//')"
}
bindkey '^R' 'cmdvault search\n'
`)
	case "fish":
		return `# CmdVault Fish Integration
function cmdvault_precmd --on-event fish_preexec
    cmdvault add $argv
end
bind \cr 'cmdvault search'
`
	default:
		return "# Unsupported shell: " + shell
	}
}

func init() {
	// Add flags
	searchCmd.Flags().IntP("limit", "l", 20, "Maximum number of results")
	listCmd.Flags().IntP("limit", "l", 20, "Maximum number of commands to list")
	listCmd.Flags().BoolP("favorites", "f", false, "Show only favorites")
	addCmd.Flags().StringP("tags", "t", "", "Tags for the command (comma-separated)")
	addCmd.Flags().StringP("notes", "n", "", "Notes for the command")

	// Add commands
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(favoriteCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(initCmd)
}

// getOS returns the current operating system
func getOS() string {
	return runtime.GOOS
}
