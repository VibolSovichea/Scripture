package ui

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/vibolsovichea/etch/internal/asset"
	"github.com/vibolsovichea/etch/internal/config"
	"github.com/vibolsovichea/etch/internal/note"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const maxRecent = 5

type view int

const (
	viewDashboard view = iota
	viewFinder
	viewCreate
	viewEditor
)

type finderMode int

const (
	finderBrowse finderMode = iota
	finderSearch
	finderDelete
)

type createStep int

const (
	createStepTitle createStep = iota
	createStepTags
)

type dashAction int

const (
	dashActionFind dashAction = iota
	dashActionNew
	dashActionQuit
)

var dashActions = []struct {
	key    string
	label  string
	action dashAction
}{
	{"f", "Find Note", dashActionFind},
	{"n", "New Note", dashActionNew},
	{"q", "Quit", dashActionQuit},
}

type AppModel struct {
	cfg    *config.Config
	ascii  string
	notes  []*note.Note
	recent []*note.Note
	width  int
	height int
	view       view
	dashCursor int
	finderMode   finderMode
	filtered     []*note.Note
	finderCursor int
	finderScroll int
	input        textinput.Model
	createStep  createStep
	newTitle    string
	createInput textinput.Model

	editor EditorModel

	err error
}

func NewAppModel(cfg *config.Config) AppModel {
	ti := textinput.New()
	ti.CharLimit = 256
	ti.Width = 40
	ti.PromptStyle = inputLabelStyle
	ti.TextStyle = lipgloss.NewStyle().Foreground(sand)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(gold)

	ci := textinput.New()
	ci.CharLimit = 256
	ci.Width = 50
	ci.PromptStyle = createLabelStyle
	ci.TextStyle = lipgloss.NewStyle().Foreground(sand)
	ci.Cursor.Style = lipgloss.NewStyle().Foreground(gold)

	ascii := strings.TrimRight(asset.ASCIIArt, "\n")

	m := AppModel{
		cfg:         cfg,
		ascii:       ascii,
		input:       ti,
		createInput: ci,
		view:        viewDashboard,
	}
	m.loadNotes()
	return m
}

func (m *AppModel) loadNotes() {
	notes, err := note.ListAll(m.cfg.VaultPath)
	if err != nil {
		m.err = err
		return
	}
	m.notes = notes
	m.filtered = notes

	sorted := make([]*note.Note, len(notes))
	copy(sorted, notes)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Modified.After(sorted[j].Modified)
	})
	if len(sorted) > maxRecent {
		sorted = sorted[:maxRecent]
	}
	m.recent = sorted

	if m.finderCursor >= len(m.filtered) {
		m.finderCursor = max(0, len(m.filtered)-1)
	}
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case editorCloseMsg:
		m.loadNotes()
		m.view = viewDashboard
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.view == viewEditor {
			m.editor.Resize(m.width, m.height)
		}
		return m, nil

	case tea.KeyMsg:
		switch m.view {
		case viewDashboard:
			return m.updateDashboard(msg)
		case viewFinder:
			return m.updateFinder(msg)
		case viewCreate:
			return m.updateCreate(msg)
		case viewEditor:
			var cmd tea.Cmd
			m.editor, cmd = m.editor.Update(msg)
			return m, cmd
		}
	}

	if m.view == viewEditor {
		var cmd tea.Cmd
		m.editor, cmd = m.editor.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m AppModel) updateDashboard(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	totalItems := len(dashActions) + len(m.recent)

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "j", "down":
		if m.dashCursor < totalItems-1 {
			m.dashCursor++
		}
	case "k", "up":
		if m.dashCursor > 0 {
			m.dashCursor--
		}
	case "f":
		m.openFinder()
		return m, textinput.Blink
	case "n":
		m.openCreate()
		return m, textinput.Blink
	case "q":
		return m, tea.Quit
	case "enter":
		return m.handleDashEnter()
	}
	return m, nil
}

func (m AppModel) handleDashEnter() (tea.Model, tea.Cmd) {
	if m.dashCursor < len(dashActions) {
		switch dashActions[m.dashCursor].action {
		case dashActionFind:
			m.openFinder()
			return m, textinput.Blink
		case dashActionNew:
			m.openCreate()
			return m, textinput.Blink
		case dashActionQuit:
			return m, tea.Quit
		}
	}

	recentIdx := m.dashCursor - len(dashActions)
	if recentIdx >= 0 && recentIdx < len(m.recent) {
		m.openEditor(m.recent[recentIdx])
		return m, nil
	}
	return m, nil
}

func (m *AppModel) openFinder() {
	m.view = viewFinder
	m.finderMode = finderSearch
	m.finderCursor = 0
	m.finderScroll = 0
	m.filtered = m.notes
	m.input.SetValue("")
	m.input.Placeholder = "Search notes..."
	m.input.Focus()
}

func (m *AppModel) openEditor(n *note.Note) {
	m.editor = NewEditorModel(n, m.width, m.height)
	m.view = viewEditor
}

func (m *AppModel) openCreate() {
	m.view = viewCreate
	m.createStep = createStepTitle
	m.newTitle = ""
	m.createInput.SetValue("")
	m.createInput.Placeholder = "Enter a title for your note"
	m.createInput.Focus()
}

func (m AppModel) updateFinder(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.finderMode == finderDelete {
		return m.handleFinderDelete(msg)
	}
	return m.handleFinderNormal(msg)
}

func (m AppModel) handleFinderNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.view = viewDashboard
		m.input.Blur()
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "down", "ctrl+n":
		if m.finderCursor < len(m.filtered)-1 {
			m.finderCursor++
			m.clampFinderScroll()
		}
	case "up", "ctrl+p":
		if m.finderCursor > 0 {
			m.finderCursor--
			m.clampFinderScroll()
		}
	case "enter":
		if len(m.filtered) > 0 {
			m.openEditor(m.filtered[m.finderCursor])
			return m, nil
		}
	case "ctrl+d":
		if len(m.filtered) > 0 {
			m.finderMode = finderDelete
		}
		return m, nil
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)

		query := m.input.Value()
		if query == "" {
			m.filtered = m.notes
		} else {
			m.filtered = filterNotes(m.notes, query)
		}
		m.finderCursor = 0
		m.finderScroll = 0
		return m, cmd
	}
	return m, nil
}

func (m AppModel) handleFinderDelete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		if len(m.filtered) > 0 {
			n := m.filtered[m.finderCursor]
			trashDir := filepath.Join(m.cfg.VaultPath, ".etch", "trash")
			n.Delete(trashDir)
			m.loadNotes()
			m.filtered = m.notes
		}
		m.finderMode = finderSearch
	default:
		m.finderMode = finderSearch
	}
	return m, nil
}

func (m AppModel) updateCreate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		if m.createStep == createStepTags {
			m.createStep = createStepTitle
			m.createInput.SetValue(m.newTitle)
			m.createInput.Placeholder = "Enter a title for your note"
			return m, nil
		}
		m.view = viewDashboard
		m.createInput.Blur()
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "enter":
		if m.createStep == createStepTitle {
			m.newTitle = m.createInput.Value()
			if m.newTitle == "" {
				m.view = viewDashboard
				m.createInput.Blur()
				return m, nil
			}
			m.createStep = createStepTags
			m.createInput.SetValue("")
			m.createInput.Placeholder = "e.g. project, idea, draft (optional)"
			return m, nil
		}
		tags := parseTags(m.createInput.Value())
		dir := filepath.Join(m.cfg.VaultPath, "notes")
		n, err := note.Create(dir, m.newTitle, tags)
		if err != nil {
			m.err = err
			m.view = viewDashboard
			m.createInput.Blur()
			return m, nil
		}
		m.createInput.Blur()
		m.loadNotes()
		m.openEditor(n)
		return m, nil
	}

	var cmd tea.Cmd
	m.createInput, cmd = m.createInput.Update(msg)
	return m, cmd
}

func (m AppModel) View() string {
	if m.width == 0 {
		return ""
	}

	switch m.view {
	case viewFinder:
		return m.viewFinder()
	case viewCreate:
		return m.viewCreate()
	case viewEditor:
		return m.editor.View()
	default:
		return m.viewDashboard()
	}
}

func (m AppModel) viewDashboard() string {
	var sections []string

	art := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, asciiStyle.Render(m.ascii))
	sections = append(sections, art, "", "")

	for i, a := range dashActions {
		isSelected := i == m.dashCursor

		indicator := " "
		labelStyle := dashActionStyle

		if isSelected {
			indicator = ">"
			labelStyle = dashActionSelectedStyle
		}

		row := fmt.Sprintf("%s  %s  %s",
			dashActionSelectedStyle.Render(indicator),
			dashActionKeyStyle.Render(a.key),
			labelStyle.Render(a.label),
		)

		sections = append(sections,
			lipgloss.PlaceHorizontal(m.width, lipgloss.Center, row),
			"",
		)
	}

	if len(m.recent) > 0 {
		sections = append(sections, "")

		title := dashRecentTitleStyle.Render("Recent Notes")
		sections = append(sections,
			lipgloss.PlaceHorizontal(m.width, lipgloss.Center, title),
			"",
		)

		for i, n := range m.recent {
			globalIdx := len(dashActions) + i

			noteTitle := truncate(n.Title, 40) // fixed width for alignment
			date := dashRecentDateStyle.Render(relativeTime(n.Modified))

			var row string
			if globalIdx == m.dashCursor {
				row = fmt.Sprintf("%s  %-40s  %s",
					dashRecentSelectedStyle.Render(">"),
					dashRecentSelectedStyle.Render(noteTitle),
					date,
				)
			} else {
				row = fmt.Sprintf("   %-40s  %s",
					dashRecentItemStyle.Render(noteTitle),
					date,
				)
			}

			sections = append(sections,
				lipgloss.PlaceHorizontal(m.width, lipgloss.Center, row),
			)
		}
	}

	sections = append(sections, "", "")

	footer := dashFooterStyle.Render(
		fmt.Sprintf("etch — %d notes in vault", len(m.notes)),
	)

	sections = append(sections,
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, footer),
	)

	content := strings.Join(sections, "\n")

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

func (m AppModel) viewFinder() string {
	modalW := min(m.width-4, 100)
	modalH := min(m.height-4, 30)
	if modalW < 40 {
		modalW = 40
	}
	if modalH < 10 {
		modalH = 10
	}

	listW := modalW * 2 / 5
	previewW := modalW - listW
	innerW := modalW - 2 // inside border

	contentH := modalH - 5
	if contentH < 3 {
		contentH = 3
	}

	titleBar := finderTitleStyle.Render(" Find Notes")

	topDiv := finderPreviewDivStyle.Render(strings.Repeat("─", innerW))

	left := m.finderList(listW, contentH)

	var sepLines []string
	for i := 0; i < contentH; i++ {
		sepLines = append(sepLines, finderPreviewDivStyle.Render("│"))
	}
	sep := strings.Join(sepLines, "\n")
	right := m.finderPreview(previewW, contentH)
	content := lipgloss.JoinHorizontal(lipgloss.Top, left, sep, right)
	bottomDiv := finderPreviewDivStyle.Render(strings.Repeat("─", innerW))
	var searchBar string
	if m.finderMode == finderDelete {
		n := m.filtered[m.finderCursor]
		searchBar = deleteWarnStyle.Render(
			fmt.Sprintf(" Delete \"%s\"? ", truncate(n.Title, innerW-25)),
		) + helpDescStyle.Render("y to confirm, any key to cancel")
	} else {
		prompt := inputLabelStyle.Render(" > ")
		searchBar = prompt + m.input.View()
	}
	help := m.finderHelp()

	modal := lipgloss.JoinVertical(lipgloss.Left,
		titleBar,
		topDiv,
		content,
		bottomDiv,
		searchBar,
		help,
	)

	framed := finderBorderStyle.
		Width(modalW).
		Render(modal)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		framed,
	)
}

func (m AppModel) finderList(w, h int) string {
	var lines []string

	if len(m.filtered) == 0 {
		empty := finderPreviewEmptyStyle.Render("  No matches")
		lines = append(lines, empty)
	} else {
		for i, n := range m.filtered {
			if i < m.finderScroll {
				continue
			}
			if len(lines) >= h {
				break
			}

			cursor := "  "
			if i == m.finderCursor {
				cursor = finderCursorStyle.Render("> ")
			}

			maxW := w - 6
			title := truncate(n.Title, maxW)
			if i == m.finderCursor {
				title = finderItemSelectedStyle.Render(title)
			} else {
				title = finderItemStyle.Render(title)
			}

			lines = append(lines, cursor+title)
		}
	}

	for len(lines) < h {
		lines = append(lines, "")
	}

	countStr := finderCountStyle.Render(
		fmt.Sprintf("  %d/%d", len(m.filtered), len(m.notes)),
	)
	if len(lines) > 0 {
		lines[len(lines)-1] = countStr
	}

	return lipgloss.NewStyle().
		Width(w).
		Height(h).
		Render(strings.Join(lines, "\n"))
}

func (m AppModel) finderPreview(w, h int) string {
	if len(m.filtered) == 0 {
		empty := finderPreviewEmptyStyle.Render("No note selected")
		return lipgloss.NewStyle().
			Width(w).
			Height(h).
			Padding(0, 1).
			Render(lipgloss.Place(w-2, h, lipgloss.Center, lipgloss.Center, empty))
	}

	n := m.filtered[m.finderCursor]
	var b strings.Builder

	b.WriteString(finderPreviewTitleStyle.Render(n.Title))
	b.WriteString("\n")

	if len(n.Tags) > 0 {
		tags := make([]string, len(n.Tags))
		for i, t := range n.Tags {
			tags[i] = "#" + t
		}
		b.WriteString(finderPreviewTagStyle.Render(strings.Join(tags, "  ")))
		b.WriteString("\n")
	}

	b.WriteString(finderPreviewMetaStyle.Render(
		fmt.Sprintf("Created %s  Modified %s",
			n.Created.Format("Jan 02, 2006"),
			n.Modified.Format("Jan 02, 2006"))))
	b.WriteString("\n")

	divW := w - 4
	if divW < 5 {
		divW = 5
	}
	b.WriteString(finderPreviewDivStyle.Render(strings.Repeat("─", divW)))
	b.WriteString("\n")

	if n.Body == "" {
		b.WriteString(finderPreviewEmptyStyle.Render("Empty note"))
	} else {
		bodyLines := strings.Split(n.Body, "\n")
		maxLines := h - 6
		if maxLines < 1 {
			maxLines = 1
		}
		if len(bodyLines) > maxLines {
			bodyLines = bodyLines[:maxLines]
			bodyLines = append(bodyLines, finderPreviewEmptyStyle.Render("..."))
		}
		b.WriteString(finderPreviewBodyStyle.Render(strings.Join(bodyLines, "\n")))
	}

	return lipgloss.NewStyle().
		Width(w).
		Height(h).
		Padding(0, 1).
		Render(b.String())
}

func (m AppModel) finderHelp() string {
	if m.finderMode == finderDelete {
		return ""
	}

	parts := []string{
		helpKeyStyle.Render("Enter") + helpDescStyle.Render(" open"),
		helpKeyStyle.Render("Ctrl+d") + helpDescStyle.Render(" delete"),
		helpKeyStyle.Render("Esc") + helpDescStyle.Render(" back"),
	}
	return " " + strings.Join(parts, "  ")
}

func (m AppModel) viewCreate() string {
	modalW := min(m.width-4, 60)
	if modalW < 30 {
		modalW = 30
	}
	innerW := modalW - 2

	title := createTitleStyle.Render(" New Note")
	div := finderPreviewDivStyle.Render(strings.Repeat("─", innerW))

	step1Label := "Title"
	step2Label := "Tags"

	var step1, step2 string
	switch m.createStep {
	case createStepTitle:
		step1 = createStepActiveStyle.Render("● " + step1Label)
		step2 = createStepPendingStyle.Render("○ " + step2Label)
	case createStepTags:
		step1 = createStepDoneStyle.Render("✓ " + step1Label)
		step2 = createStepActiveStyle.Render("● " + step2Label)
	}

	steps := fmt.Sprintf(" %s    %s", step1, step2)

	var titlePreview string
	if m.createStep == createStepTags {
		titlePreview = fmt.Sprintf(" %s %s",
			createLabelStyle.Render("Title:"),
			createValueStyle.Render(m.newTitle),
		)
	}

	var inputLabel string
	if m.createStep == createStepTitle {
		inputLabel = createLabelStyle.Render(" Title")
	} else {
		inputLabel = createLabelStyle.Render(" Tags")
	}
	inputLine := fmt.Sprintf(" %s\n %s", inputLabel, m.createInput.View())

	help := " " + helpKeyStyle.Render("Enter") + helpDescStyle.Render(" next  ") +
		helpKeyStyle.Render("Esc") + helpDescStyle.Render(" back")

	var parts []string
	parts = append(parts, title)
	parts = append(parts, div)
	parts = append(parts, "")
	parts = append(parts, steps)
	parts = append(parts, "")
	if titlePreview != "" {
		parts = append(parts, titlePreview)
		parts = append(parts, "")
	}
	parts = append(parts, inputLine)
	parts = append(parts, "")
	parts = append(parts, div)
	parts = append(parts, help)

	modal := lipgloss.JoinVertical(lipgloss.Left, parts...)

	framed := createBorderStyle.
		Width(modalW).
		Render(modal)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		framed,
	)
}

func (m *AppModel) clampFinderScroll() {
	visibleH := m.finderVisibleCount()
	if m.finderCursor < m.finderScroll {
		m.finderScroll = m.finderCursor
	}
	if m.finderCursor >= m.finderScroll+visibleH {
		m.finderScroll = m.finderCursor - visibleH + 1
	}
	if m.finderScroll < 0 {
		m.finderScroll = 0
	}
}

func (m AppModel) finderVisibleCount() int {
	modalH := min(m.height-4, 30)
	h := modalH - 7
	if h < 3 {
		h = 3
	}
	return h
}

func filterNotes(notes []*note.Note, query string) []*note.Note {
	q := strings.ToLower(query)
	var result []*note.Note
	for _, n := range notes {
		if strings.Contains(strings.ToLower(n.Title), q) ||
			strings.Contains(strings.ToLower(n.Body), q) ||
			containsTag(n.Tags, q) {
			result = append(result, n)
		}
	}
	return result
}

func containsTag(tags []string, query string) bool {
	for _, t := range tags {
		if strings.Contains(strings.ToLower(t), query) {
			return true
		}
	}
	return false
}

func parseTags(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	var tags []string
	for _, t := range strings.Split(s, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

func truncate(s string, maxW int) string {
	if maxW <= 0 {
		return ""
	}
	if len(s) <= maxW {
		return s
	}
	if maxW <= 3 {
		return s[:maxW]
	}
	return s[:maxW-3] + "..."
}

func relativeTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		m := int(diff.Minutes())
		if m == 1 {
			return "1 min ago"
		}
		return fmt.Sprintf("%d mins ago", m)
	case diff < 24*time.Hour:
		h := int(diff.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	case diff < 7*24*time.Hour:
		d := int(diff.Hours() / 24)
		if d == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", d)
	default:
		return t.Format("Jan 02, 2006")
	}
}
