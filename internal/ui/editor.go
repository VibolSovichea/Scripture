package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/vibolsovichea/etch/internal/note"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const gutterWidth = 6

type editorCloseMsg struct{}

type EditorModel struct {
	ta       textarea.Model
	vim      vimState
	note     *note.Note
	original string
	width    int
	height   int
	err      error
}

func NewEditorModel(n *note.Note, width, height int) EditorModel {
	ta := textarea.New()
	ta.SetValue(n.Body)
	ta.ShowLineNumbers = false
	ta.SetWidth(width - gutterWidth - 2)
	ta.SetHeight(height - 2)
	ta.CharLimit = 0
	ta.Blur()

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
			var taCmd tea.Cmd
			m.ta, taCmd = m.ta.Update(msg)
			return m, taCmd
		}

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
	content := m.ta.Value()
	lines := strings.Split(content, "\n")

	cursorRow := m.ta.Line()
	cursorCol := m.ta.LineInfo().ColumnOffset

	editorH := m.height - 2
	if editorH < 1 {
		editorH = 1
	}

	scrollOff := 0
	if cursorRow >= editorH {
		scrollOff = cursorRow - editorH + 1
	}

	var viewLines []string
	inCodeBlock := false

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
			gutter := edGutterStyle.Render("~")
			sep := edGutterSepStyle.Render("│")
			viewLines = append(viewLines, gutter+sep)
			continue
		}

		lineNum := fmt.Sprintf("%d", lineIdx+1)
		var gutter string
		if lineIdx == cursorRow {
			gutter = edGutterActiveStyle.Render(lineNum)
		} else {
			gutter = edGutterStyle.Render(lineNum)
		}
		sep := edGutterSepStyle.Render("│")

		line := lines[lineIdx]
		var styledLine string
		styledLine, inCodeBlock = highlightLine(line, inCodeBlock)

		if lineIdx == cursorRow && m.vim.mode == vimNormal {
			styledLine = renderCursorLine(line, cursorCol, inCodeBlock)
		}

		viewLines = append(viewLines, gutter+sep+" "+styledLine)
	}

	for len(viewLines) < editorH {
		gutter := edGutterStyle.Render("~")
		sep := edGutterSepStyle.Render("│")
		viewLines = append(viewLines, gutter+sep)
	}

	editor := strings.Join(viewLines, "\n")

	statusBar := m.renderStatusBar(cursorRow, cursorCol)
	cmdLine := m.renderCmdLine()

	return lipgloss.JoinVertical(lipgloss.Left, editor, statusBar, cmdLine)
}

func renderCursorLine(line string, col int, inCodeBlock bool) string {
	if len(line) == 0 {
		return edCursorBlockStyle.Render(" ")
	}

	var result strings.Builder

	runes := []rune(line)
	for i, r := range runes {
		s := string(r)
		if i == col {
			result.WriteString(edCursorBlockStyle.Render(s))
		} else {
			result.WriteString(edTextStyle.Render(s))
		}
	}

	if col >= len(runes) {
		result.WriteString(edCursorBlockStyle.Render(" "))
	}

	return result.String()
}

func (m EditorModel) renderStatusBar(row, col int) string {
	var modeBadge string
	switch m.vim.mode {
	case vimInsert:
		modeBadge = edModeInsertStyle.Render("INSERT")
	case vimCommand:
		modeBadge = edModeCommandStyle.Render("COMMAND")
	default:
		modeBadge = edModeNormalStyle.Render("NORMAL")
	}

	filename := edFileNameStyle.Render(filepath.Base(m.note.Path))
	mod := ""
	if m.Modified() {
		mod = edModifiedStyle.Render("[+]")
	}

	pos := edPosStyle.Render(fmt.Sprintf("Ln %d, Col %d", row+1, col+1))

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

	if m.vim.mode == vimInsert {
		return helpDescStyle.Render(" Esc normal  Ctrl+S save  Ctrl+Q quit")
	}
	return helpDescStyle.Render(" i insert  :w save  :q quit  Ctrl+S save  Ctrl+Q quit")
}
