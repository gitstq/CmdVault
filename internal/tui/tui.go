// Package tui provides the terminal user interface for CmdVault
package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/cmdvault/cmdvault/internal/storage"
	"github.com/cmdvault/cmdvault/internal/search"
)

// Styles for the TUI
var (
	titleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#7C3AED")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Padding(1, 2)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C3AED")).
			Bold(true)

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981"))

	favoriteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B"))

	statsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3B82F6"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7C3AED")).
			Padding(0, 1)
)

// item represents a list item
type item struct {
	command storage.Command
}

func (i item) FilterValue() string {
	return i.command.Command
}

func (i item) Title() string {
	prefix := ""
	if i.command.Favorite {
		prefix = "⭐ "
	}
	return prefix + i.command.Command
}

func (i item) Description() string {
	tags := i.command.Tags
	if tags == "" {
		tags = "no tags"
	}
	return fmt.Sprintf("exec: %d | %s | %s", 
		i.command.ExecCount, 
		tags,
		i.command.LastUsed.Format("2006-01-02 15:04"),
	)
}

// Model represents the TUI model
type Model struct {
	list          list.Model
	textInput     textinput.Model
	db            *storage.Database
	searcher      *search.FuzzySearcher
	commands      []storage.Command
	filtered      []storage.Command
	mode          string // "search", "normal", "detail"
	selectedID    int64
	showFavorites bool
	err           error
	width         int
	height        int
}

// NewModel creates a new TUI model
func NewModel(db *storage.Database) Model {
	// Initialize text input
	ti := textinput.New()
	ti.Placeholder = "Search commands..."
	ti.Focus()

	// Initialize list
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = selectedStyle
	delegate.Styles.NormalTitle = commandStyle

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "🔐 CmdVault"
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return Model{
		list:      l,
		textInput: ti,
		db:        db,
		mode:      "search",
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.loadCommands,
	)
}

// loadCommands loads commands from the database
func (m Model) loadCommands() tea.Msg {
	commands, err := m.db.GetCommands(1000, 0, m.showFavorites)
	if err != nil {
		return errMsg{err}
	}
	return commandsMsg{commands}
}

// messages
type commandsMsg struct {
	commands []storage.Command
}

type errMsg struct{ error }

// Update handles updates to the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.mode == "detail" {
				m.mode = "normal"
				return m, nil
			}
			return m, tea.Quit

		case "enter":
			if m.mode == "search" && m.textInput.Value() != "" {
				// Perform search
				m.performSearch()
				return m, nil
			}
			if m.mode == "normal" {
				if i, ok := m.list.SelectedItem().(item); ok {
					// Output selected command
					fmt.Println(i.command.Command)
					return m, tea.Quit
				}
			}

		case "f":
			// Toggle favorite filter
			m.showFavorites = !m.showFavorites
			return m, m.loadCommands

		case "tab":
			// Switch mode
			if m.mode == "search" {
				m.mode = "normal"
			} else {
				m.mode = "search"
				m.textInput.Focus()
			}

		case "s":
			// Toggle favorite for selected item
			if m.mode == "normal" {
				if i, ok := m.list.SelectedItem().(item); ok {
					m.db.ToggleFavorite(i.command.ID)
					return m, m.loadCommands
				}
			}

		case "d":
			// Delete selected item
			if m.mode == "normal" {
				if i, ok := m.list.SelectedItem().(item); ok {
					m.db.DeleteCommand(i.command.ID)
					return m, m.loadCommands
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-3)

	case commandsMsg:
		m.commands = msg.commands
		m.filtered = msg.commands
		m.searcher = search.NewFuzzySearcher(msg.commands)
		items := make([]list.Item, len(msg.commands))
		for i, cmd := range msg.commands {
			items[i] = item{command: cmd}
		}
		m.list.SetItems(items)
	}

	// Update text input
	if m.mode == "search" {
		var tiCmd tea.Cmd
		m.textInput, tiCmd = m.textInput.Update(msg)
		cmds = append(cmds, tiCmd)
	}

	// Update list
	var lCmd tea.Cmd
	m.list, lCmd = m.list.Update(msg)
	cmds = append(cmds, lCmd)

	return m, tea.Batch(cmds...)
}

// performSearch performs fuzzy search
func (m *Model) performSearch() {
	if m.searcher == nil {
		return
	}

	results := m.searcher.Search(m.textInput.Value(), 100)
	m.filtered = make([]storage.Command, len(results))
	items := make([]list.Item, len(results))

	for i, r := range results {
		m.filtered[i] = r.Command
		items[i] = item{command: r.Command}
	}

	m.list.SetItems(items)
}

// View renders the TUI
func (m Model) View() string {
	var b strings.Builder

	// Search bar
	searchBar := boxStyle.Render(m.textInput.View())
	b.WriteString(searchBar + "\n")

	// Mode indicator
	modeText := "🔍 Search Mode"
	if m.mode == "normal" {
		modeText = "📋 Browse Mode"
	}
	if m.showFavorites {
		modeText += " | ⭐ Favorites Only"
	}
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render(modeText) + "\n\n")

	// List
	b.WriteString(m.list.View())

	// Help
	help := "tab: switch mode | enter: select | f: filter favorites | s: star | d: delete | esc/ctrl+c: quit"
	b.WriteString("\n" + helpStyle.Render(help))

	return b.String()
}

// RunTUI starts the terminal UI
func RunTUI(db *storage.Database) error {
	p := tea.NewProgram(
		NewModel(db),
		tea.WithAltScreen(),
	)

	_, err := p.Run()
	return err
}
