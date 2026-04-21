package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/vibolsovichea/scripture/internal/note"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const gutterWidth = 6 // "  42 │ "

type editorCloseMsg struct{}

type EditorModel struct {
	ta       textarea.Model
	vim      vimState
	note     *note.Note
	original string // original body to track modifications
	width    int
	height   int
	err      error
}

func NewEditorModel(n *note.Note, width, height int) EditorModel {
	ta := textarea.New()
	ta.SetValue(n.Body)
	ta.ShowLineNumbers = false
	ta.SetWidth(width - gutterWidth - 2)
	ta.SetHeight(height - 2) // status bar + command line
	ta.CharLimit = 0         // no limit
	ta.Blur()                // start in normal mode

	// Style the textarea to match theme
	ta.FocusedStyle.Base = lipgloss.NewStyle().Foreground(ivory)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().Foreground(ivory)
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(darkStone)
	ta.BlurredStyle.Base = lipgloss.NewStyle().Foreground(ivory)
	ta.BlurredStyle.CursorLine = lipgloss.NewStyle().Foreground(ivory)
	ta.BlurredStyle.Placeholder = lipgloss.NewStyle().Foreground(darkStone)
	ta.Prompt = ""

	return EditorModel{
		ta:       ta,
		vim:      newVimState(),
		note:     n,
		original: n.Body,
		width:    width,
		height:   height,
	}
}

func (m EditorModel) Modified() bool {
	return m.ta.Value() != m.original
}

func (m *EditorModel) Save() error {
	m.note.SetBody(m.ta.Value())
	err := m.note.Save()
	if err == nil {
		m.original = m.ta.Value()
	}
	return err
}

func (m *EditorModel) Resize(w, h int) {
	m.width = w
	m.height = h
	m.ta.SetWidth(w - gutterWidth - 2)
	m.ta.SetHeight(h - 2)
}

func (m EditorModel) Update(msg tea.Msg) (EditorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmd := m.vim.handleKey(msg, &m.ta)

		switch cmd {
		case vimCmdSave:
			m.Save()
			return m, nil
		case vimCmdQuit:
			return m, func() tea.Msg { return editorCloseMsg{} }
		case vimCmdSaveQuit:
			m.Save()
			return m, func() tea.Msg { return editorCloseMsg{} }
		case vimCmdPassthrough:
			// Insert mode — let textarea handle it
			var taCmd tea.Cmd
			m.ta, taCmd = m.ta.Update(msg)
			return m, taCmd
		}

		// For 'o' key in normal mode — we switched to insert and need
		// to insert a newline
		if m.vim.mode == vimInsert {
			m.ta.Focus()
		} else {
			m.ta.Blur()
		}

		return m, nil
	}
	return m, nil
}

func (m EditorModel) View() string {
	// Get content lines
	content := m.ta.Value()
	lines := strings.Split(content, "\n")

	cursorRow := m.ta.Line()
	cursorCol := m.ta.LineInfo().ColumnOffset

	// Calculate viewport
	editorH := m.height - 2 // status bar + help/command line
	if editorH < 1 {
		editorH = 1
	}

	// Determine scroll offset
	scrollOff := 0
	if cursorRow >= editorH {
		scrollOff = cursorRow - editorH + 1
	}

	// Render visible lines with gutter and highlighting
	var viewLines []string
	inCodeBlock := false

	// Pre-process code block state for lines before viewport
	for i := 0; i < scrollOff && i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
		}
	}

	textW := m.width - gutterWidth
	if textW < 10 {
		textW = 10
	}

	for vi := 0; vi < editorH; vi++ {
		lineIdx := scrollOff + vi

		if lineIdx >= len(lines) {
			// Tilde for empty lines past content (like vim)
			gutter := edGutterStyle.Render("~")
			sep := edGutterSepStyle.Render("│")
			viewLines = append(viewLines, gutter+sep)
			continue
		}

		// Gutter
		lineNum := fmt.Sprintf("%d", lineIdx+1)
		var gutter string
		if lineIdx == cursorRow {
			gutter = edGutterActiveStyle.Render(lineNum)
		} else {
			gutter = edGutterStyle.Render(lineNum)
		}
		sep := edGutterSepStyle.Render("│")

		// Content with highlighting
		line := lines[lineIdx]
		var styledLine string
		styledLine, inCodeBlock = highlightLine(line, inCodeBlock)

		// Cursor rendering in normal mode
		if lineIdx == cursorRow && m.vim.mode == vimNormal {
			styledLine = renderCursorLine(line, cursorCol, inCodeBlock)
		}

		// Truncate if too wide
		viewLines = append(viewLines, gutter+sep+" "+styledLine)
	}

	// Pad to fill
	for len(viewLines) < editorH {
		gutter := edGutterStyle.Render("~")
		sep := edGutterSepStyle.Render("│")
		viewLines = append(viewLines, gutter+sep)
	}

	editor := strings.Join(viewLines, "\n")

	// Status bar
	statusBar := m.renderStatusBar(cursorRow, cursorCol)

	// Command / help line
	cmdLine := m.renderCmdLine()

	return lipgloss.JoinVertical(lipgloss.Left, editor, statusBar, cmdLine)
}

func renderCursorLine(line string, col int, inCodeBlock bool) string {
	if len(line) == 0 {
		return edCursorBlockStyle.Render(" ")
	}

	// Render the line with a block cursor at the cursor position
	var result strings.Builder

	runes := []rune(line)
	for i, r := range runes {
		s := string(r)
		if i == col {
			result.WriteString(edCursorBlockStyle.Render(s))
		} else {
			// Apply basic highlighting
			result.WriteString(edTextStyle.Render(s))
		}
	}

	// If cursor is at end of line
	if col >= len(runes) {
		result.WriteString(edCursorBlockStyle.Render(" "))
	}

	return result.String()
}

func (m EditorModel) renderStatusBar(row, col int) string {
	// Mode badge
	var modeBadge string
	switch m.vim.mode {
	case vimInsert:
		modeBadge = edModeInsertStyle.Render("INSERT")
	case vimCommand:
		modeBadge = edModeCommandStyle.Render("COMMAND")
	default:
		modeBadge = edModeNormalStyle.Render("NORMAL")
	}

	// Filename
	filename := edFileNameStyle.Render(filepath.Base(m.note.Path))

	// Modified indicator
	mod := ""
	if m.Modified() {
		mod = edModifiedStyle.Render("[+]")
	}

	// Position
	pos := edPosStyle.Render(fmt.Sprintf("Ln %d, Col %d", row+1, col+1))

	// Spacer
	leftParts := modeBadge + filename + mod
	leftW := lipgloss.Width(leftParts)
	rightW := lipgloss.Width(pos)
	spacerW := m.width - leftW - rightW
	if spacerW < 0 {
		spacerW = 0
	}
	spacer := edStatusBarStyle.Render(strings.Repeat(" ", spacerW))

	return leftParts + spacer + pos
}

func (m EditorModel) renderCmdLine() string {
	if m.vim.mode == vimCommand {
		return edCmdLineStyle.Render(":" + m.vim.cmdBuffer)
	}

	// Help hints
	if m.vim.mode == vimInsert {
		return helpDescStyle.Render(" Esc normal  Ctrl+S save  Ctrl+Q quit")
	}
	return helpDescStyle.Render(" i insert  :w save  :q quit  Ctrl+S save  Ctrl+Q quit")
}
