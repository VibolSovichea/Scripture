package ui

import "github.com/charmbracelet/lipgloss"

var (
	gold      = lipgloss.Color("#D4A843")
	amber     = lipgloss.Color("#C4873B")
	sand      = lipgloss.Color("#E8D5B7")
	stone     = lipgloss.Color("#8B8680")
	darkStone = lipgloss.Color("#5C5955")
	ivory     = lipgloss.Color("#FAF3E8")
	rust      = lipgloss.Color("#A0522D")
	sage      = lipgloss.Color("#7D9B76")
	bg        = lipgloss.Color("#1A1816")
)

var (
	borderStyle = lipgloss.RoundedBorder()

	asciiStyle = lipgloss.NewStyle().
			Foreground(gold).
			Bold(true)

	dashActionStyle = lipgloss.NewStyle().
			Foreground(sand)

	dashActionKeyStyle = lipgloss.NewStyle().
				Foreground(amber).
				Bold(true)

	dashActionSelectedStyle = lipgloss.NewStyle().
				Foreground(gold).
				Bold(true)

	dashRecentTitleStyle = lipgloss.NewStyle().
				Foreground(stone).
				Bold(true)

	dashRecentItemStyle = lipgloss.NewStyle().
				Foreground(sand)

	dashRecentSelectedStyle = lipgloss.NewStyle().
				Foreground(gold).
				Bold(true)

	dashRecentDateStyle = lipgloss.NewStyle().
				Foreground(darkStone)

	dashFooterStyle = lipgloss.NewStyle().
			Foreground(darkStone).
			Italic(true)

	finderBorderStyle = lipgloss.NewStyle().
				BorderStyle(borderStyle).
				BorderForeground(darkStone)

	finderTitleStyle = lipgloss.NewStyle().
			Foreground(gold).
			Bold(true)

	finderInputStyle = lipgloss.NewStyle().
			Foreground(sand)

	finderItemStyle = lipgloss.NewStyle().
			Foreground(sand)

	finderItemSelectedStyle = lipgloss.NewStyle().
				Foreground(gold).
				Bold(true)

	finderCursorStyle = lipgloss.NewStyle().
				Foreground(amber).
				Bold(true)

	finderTagStyle = lipgloss.NewStyle().
			Foreground(sage)

	finderDateStyle = lipgloss.NewStyle().
			Foreground(darkStone)

	finderPreviewTitleStyle = lipgloss.NewStyle().
				Foreground(gold).
				Bold(true)

	finderPreviewMetaStyle = lipgloss.NewStyle().
				Foreground(stone)

	finderPreviewTagStyle = lipgloss.NewStyle().
				Foreground(sage).
				Italic(true)

	finderPreviewBodyStyle = lipgloss.NewStyle().
				Foreground(ivory)

	finderPreviewDivStyle = lipgloss.NewStyle().
				Foreground(darkStone)

	finderPreviewEmptyStyle = lipgloss.NewStyle().
				Foreground(darkStone).
				Italic(true)

	finderCountStyle = lipgloss.NewStyle().
			Foreground(stone)

	inputLabelStyle = lipgloss.NewStyle().
			Foreground(amber).
			Bold(true)

	deleteWarnStyle = lipgloss.NewStyle().
			Foreground(rust).
			Bold(true)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(amber).
			Bold(true)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(stone)

	createBorderStyle = lipgloss.NewStyle().
				BorderStyle(borderStyle).
				BorderForeground(gold)

	createTitleStyle = lipgloss.NewStyle().
				Foreground(gold).
				Bold(true)

	createLabelStyle = lipgloss.NewStyle().
				Foreground(stone).
				Bold(true)

	createStepActiveStyle = lipgloss.NewStyle().
				Foreground(gold).
				Bold(true)

	createStepDoneStyle = lipgloss.NewStyle().
				Foreground(sage)

	createStepPendingStyle = lipgloss.NewStyle().
				Foreground(darkStone)

	createValueStyle = lipgloss.NewStyle().
			Foreground(sand)
)
