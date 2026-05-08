// Package utils provides utility functions for CmdVault
package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// PlatformInfo returns information about the current platform
func PlatformInfo() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}

// ClearScreen clears the terminal screen
func ClearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// TruncateString truncates a string to a maximum length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// FormatDuration formats a duration in milliseconds
func FormatDuration(ms int64) string {
	if ms < 1000 {
		return fmt.Sprintf("%dms", ms)
	}
	if ms < 60000 {
		return fmt.Sprintf("%.1fs", float64(ms)/1000)
	}
	return fmt.Sprintf("%.1fm", float64(ms)/60000)
}

// FormatTime formats a time for display
func FormatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	}
	if diff < time.Hour {
		return fmt.Sprintf("%d min ago", int(diff.Minutes()))
	}
	if diff < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	}
	if diff < 7*24*time.Hour {
		return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
	}
	return t.Format("2006-01-02")
}

// ParseTags parses a comma-separated tag string
func ParseTags(tags string) []string {
	if tags == "" {
		return nil
	}
	
	result := make([]string, 0)
	for _, tag := range strings.Split(tags, ",") {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			result = append(result, tag)
		}
	}
	return result
}

// JoinTags joins tags into a comma-separated string
func JoinTags(tags []string) string {
	return strings.Join(tags, ", ")
}

// IsCommandAvailable checks if a command is available in PATH
func IsCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// GetEditor returns the user's preferred editor
func GetEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		if IsCommandAvailable("nano") {
			editor = "nano"
		} else if IsCommandAvailable("vim") {
			editor = "vim"
		} else {
			editor = "vi"
		}
	}
	return editor
}

// OpenInEditor opens a file in the user's editor
func OpenInEditor(filename string) error {
	editor := GetEditor()
	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CopyToClipboard copies text to the system clipboard
func CopyToClipboard(text string) error {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		if IsCommandAvailable("xclip") {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if IsCommandAvailable("xsel") {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else if IsCommandAvailable("wl-copy") {
			cmd = exec.Command("wl-copy")
		} else {
			return fmt.Errorf("no clipboard utility found (install xclip, xsel, or wl-copy)")
		}
	case "windows":
		cmd = exec.Command("cmd", "/c", "clip")
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// Color codes for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

// Colorize applies a color to a string
func Colorize(s, color string) string {
	return fmt.Sprintf("%s%s%s", color, s, ColorReset)
}

// PrintSuccess prints a success message
func PrintSuccess(format string, args ...interface{}) {
	fmt.Printf("%s✅ %s%s\n", ColorGreen, fmt.Sprintf(format, args...), ColorReset)
}

// PrintError prints an error message
func PrintError(format string, args ...interface{}) {
	fmt.Printf("%s❌ %s%s\n", ColorRed, fmt.Sprintf(format, args...), ColorReset)
}

// PrintWarning prints a warning message
func PrintWarning(format string, args ...interface{}) {
	fmt.Printf("%s⚠️  %s%s\n", ColorYellow, fmt.Sprintf(format, args...), ColorReset)
}

// PrintInfo prints an info message
func PrintInfo(format string, args ...interface{}) {
	fmt.Printf("%sℹ️  %s%s\n", ColorBlue, fmt.Sprintf(format, args...), ColorReset)
}
