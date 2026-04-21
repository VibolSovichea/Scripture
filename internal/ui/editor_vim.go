package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type vimMode int

const (
	vimNormal vimMode = iota
	vimInsert
	vimCommand
)

func (m vimMode) String() string {
	switch m {
	case vimInsert:
		return "INSERT"
	case vimCommand:
		return "COMMAND"
	default:
		return "NORMAL"
	}
}

type vimCmd int

const (
	vimCmdNone vimCmd = iota
	vimCmdSave
	vimCmdQuit
	vimCmdSaveQuit
	vimCmdPassthrough
)

type vimState struct {
	mode      vimMode
	cmdBuffer string
	pending   string
	register  string
}

func newVimState() vimState {
	return vimState{mode: vimNormal}
}

func (v *vimState) handleKey(msg tea.KeyMsg, ta *textarea.Model) vimCmd {
	key := msg.String()

	// Always-on shortcuts
	switch key {
	case "ctrl+s":
		return vimCmdSave
	case "ctrl+q":
		return vimCmdQuit
	}

	switch v.mode {
	case vimNormal:
		return v.handleNormal(key, ta)
	case vimInsert:
		return v.handleInsert(key, ta)
	case vimCommand:
		return v.handleCommand(key, ta)
	}
	return vimCmdNone
}

func (v *vimState) handleNormal(key string, ta *textarea.Model) vimCmd {
	// Handle pending multi-key commands
	if v.pending != "" {
		return v.handlePending(key, ta)
	}

	switch key {
	// Mode transitions
	case "i":
		v.mode = vimInsert
		return vimCmdNone
	case "a":
		v.mode = vimInsert
		col := ta.LineInfo().ColumnOffset
		ta.SetCursor(col + 1)
		return vimCmdNone
	case "A":
		v.mode = vimInsert
		ta.CursorEnd()
		return vimCmdNone
	case "I":
		v.mode = vimInsert
		ta.CursorStart()
		return vimCmdNone
	case "o":
		v.mode = vimInsert
		ta.CursorEnd()
		// We'll insert newline via passthrough
		return vimCmdNone
	case "O":
		v.mode = vimInsert
		ta.CursorStart()
		return vimCmdNone

	// Movement
	case "h", "left":
		col := ta.LineInfo().ColumnOffset
		if col > 0 {
			ta.SetCursor(col - 1)
		}
	case "l", "right":
		col := ta.LineInfo().ColumnOffset
		ta.SetCursor(col + 1)
	case "j", "down":
		ta.CursorDown()
	case "k", "up":
		ta.CursorUp()
	case "0", "home":
		ta.CursorStart()
	case "$", "end":
		ta.CursorEnd()
	case "w":
		v.wordForward(ta)
	case "b":
		v.wordBackward(ta)
	case "g":
		v.pending = "g"
	case "G":
		gotoLastLine(ta)

	// Editing
	case "x":
		deleteCharAtCursor(ta)
	case "d":
		v.pending = "d"
	case "y":
		v.pending = "y"
	case "p":
		if v.register != "" {
			ta.CursorEnd()
			ta.InsertString("\n" + v.register)
		}

	// Command mode
	case ":":
		v.mode = vimCommand
		v.cmdBuffer = ""

	// Arrow keys pass through
	}
	return vimCmdNone
}

func (v *vimState) handlePending(key string, ta *textarea.Model) vimCmd {
	combo := v.pending + key
	v.pending = ""

	switch combo {
	case "gg":
		gotoFirstLine(ta)
	case "dd":
		v.register = deleteLine(ta)
	case "yy":
		v.register = yankLine(ta)
	}
	return vimCmdNone
}

func (v *vimState) handleInsert(key string, ta *textarea.Model) vimCmd {
	if key == "esc" {
		v.mode = vimNormal
		return vimCmdNone
	}
	return vimCmdPassthrough
}

func (v *vimState) handleCommand(key string, ta *textarea.Model) vimCmd {
	switch key {
	case "esc":
		v.mode = vimNormal
		v.cmdBuffer = ""
		return vimCmdNone
	case "enter":
		cmd := v.executeCommand()
		v.mode = vimNormal
		v.cmdBuffer = ""
		return cmd
	case "backspace":
		if len(v.cmdBuffer) > 0 {
			v.cmdBuffer = v.cmdBuffer[:len(v.cmdBuffer)-1]
		}
		if v.cmdBuffer == "" {
			v.mode = vimNormal
		}
		return vimCmdNone
	default:
		if len(key) == 1 {
			v.cmdBuffer += key
		}
		return vimCmdNone
	}
}

func (v *vimState) executeCommand() vimCmd {
	cmd := strings.TrimSpace(v.cmdBuffer)
	switch cmd {
	case "w":
		return vimCmdSave
	case "q", "q!":
		return vimCmdQuit
	case "wq", "x":
		return vimCmdSaveQuit
	}
	return vimCmdNone
}

// --- Word movement ---

func (v *vimState) wordForward(ta *textarea.Model) {
	lines := strings.Split(ta.Value(), "\n")
	row := ta.Line()
	col := ta.LineInfo().ColumnOffset
	if row >= len(lines) {
		return
	}
	line := lines[row]

	// Skip current word
	i := col
	for i < len(line) && !isWordSep(line[i]) {
		i++
	}
	// Skip whitespace
	for i < len(line) && isWordSep(line[i]) {
		i++
	}

	if i < len(line) {
		ta.SetCursor(i)
	} else if row < len(lines)-1 {
		ta.CursorDown()
		ta.CursorStart()
	} else {
		ta.CursorEnd()
	}
}

func (v *vimState) wordBackward(ta *textarea.Model) {
	lines := strings.Split(ta.Value(), "\n")
	row := ta.Line()
	col := ta.LineInfo().ColumnOffset
	if row >= len(lines) {
		return
	}
	line := lines[row]

	if col == 0 {
		if row > 0 {
			ta.CursorUp()
			ta.CursorEnd()
		}
		return
	}

	i := col - 1
	// Skip whitespace backwards
	for i > 0 && isWordSep(line[i]) {
		i--
	}
	// Skip word backwards
	for i > 0 && !isWordSep(line[i-1]) {
		i--
	}
	ta.SetCursor(i)
}

// --- Text manipulation ---

func deleteCharAtCursor(ta *textarea.Model) {
	lines := strings.Split(ta.Value(), "\n")
	row := ta.Line()
	col := ta.LineInfo().ColumnOffset
	if row >= len(lines) || col >= len(lines[row]) {
		return
	}
	line := lines[row]
	lines[row] = line[:col] + line[col+1:]
	ta.SetValue(strings.Join(lines, "\n"))
	restoreCursor(ta, row, col)
}

func deleteLine(ta *textarea.Model) string {
	lines := strings.Split(ta.Value(), "\n")
	row := ta.Line()
	if row >= len(lines) {
		return ""
	}
	yanked := lines[row]

	newLines := make([]string, 0, len(lines)-1)
	newLines = append(newLines, lines[:row]...)
	if row+1 < len(lines) {
		newLines = append(newLines, lines[row+1:]...)
	}
	if len(newLines) == 0 {
		newLines = []string{""}
	}

	ta.SetValue(strings.Join(newLines, "\n"))
	if row >= len(newLines) {
		row = len(newLines) - 1
	}
	restoreCursor(ta, row, 0)
	return yanked
}

func yankLine(ta *textarea.Model) string {
	lines := strings.Split(ta.Value(), "\n")
	row := ta.Line()
	if row >= len(lines) {
		return ""
	}
	return lines[row]
}

func gotoFirstLine(ta *textarea.Model) {
	for i := 0; i < ta.LineCount(); i++ {
		ta.CursorUp()
	}
	ta.CursorStart()
}

func gotoLastLine(ta *textarea.Model) {
	for i := 0; i < ta.LineCount(); i++ {
		ta.CursorDown()
	}
	ta.CursorStart()
}

func restoreCursor(ta *textarea.Model, row, col int) {
	gotoFirstLine(ta)
	for i := 0; i < row; i++ {
		ta.CursorDown()
	}
	ta.SetCursor(col)
}

func isWordSep(b byte) bool {
	return b == ' ' || b == '\t' || b == '.' || b == ',' || b == ';' ||
		b == ':' || b == '(' || b == ')' || b == '[' || b == ']' ||
		b == '{' || b == '}'
}
