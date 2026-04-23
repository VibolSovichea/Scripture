package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SetupModel struct {
	input     textinput.Model
	VaultPath string
	Cancelled bool
	width     int
	height    int
}

func NewSetupModel() SetupModel {
	ti := textinput.New()
	home, _ := os.UserHomeDir()
	ti.Placeholder = filepath.Join(home, ".etch")
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50
	ti.TextStyle = lipgloss.NewStyle().Foreground(sand)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(gold)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(darkStone)

	return SetupModel{input: ti}
}

func (m SetupModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SetupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.VaultPath = m.input.Value()
			if m.VaultPath == "" {
				home, _ := os.UserHomeDir()
				m.VaultPath = filepath.Join(home, ".etch")
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
	modalW := min(m.width-4, 60)
	if modalW < 30 {
		modalW = 30
	}
	innerW := modalW - 2

	title := createTitleStyle.Render(" First Time Setup")
	div := finderPreviewDivStyle.Render(strings.Repeat("─", innerW))

	label := createLabelStyle.Render(" Where should etch store your notes?")
	inputLine := fmt.Sprintf(" %s", m.input.View())

	help := " " + helpKeyStyle.Render("Enter") + helpDescStyle.Render(" confirm  ") +
		helpKeyStyle.Render("Esc") + helpDescStyle.Render(" cancel")

	modal := lipgloss.JoinVertical(lipgloss.Left,
		title,
		div,
		"",
		label,
		"",
		inputLine,
		"",
		div,
		help,
	)

	framed := createBorderStyle.
		Width(modalW).
		Render(modal)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		framed,
	)
}
