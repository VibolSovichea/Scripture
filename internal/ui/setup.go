package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	setupTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#D4A843")).
			MarginBottom(1)

	setupSubtleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8B8680"))
)

type SetupModel struct {
	input     textinput.Model
	VaultPath string
	Cancelled bool
}

func NewSetupModel() SetupModel {
	ti := textinput.New()
	home, _ := os.UserHomeDir()
	ti.Placeholder = filepath.Join(home, ".scripture")
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	return SetupModel{input: ti}
}

func (m SetupModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SetupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.VaultPath = m.input.Value()
			if m.VaultPath == "" {
				home, _ := os.UserHomeDir()
				m.VaultPath = filepath.Join(home, ".scripture")
			}
			return m, tea.Quit
		case "ctrl+c", "esc":
			m.Cancelled = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m SetupModel) View() string {
	return fmt.Sprintf(
		"\n%s\n\n%s\n\n%s\n\n%s\n",
		setupTitleStyle.Render("Scripture — First Time Setup"),
		"Where should Scripture store your notes?",
		m.input.View(),
		setupSubtleStyle.Render("Press Enter to confirm • Esc to cancel"),
	)
}
