package ui

import "github.com/charmbracelet/lipgloss"

var (
	edGutterStyle = lipgloss.NewStyle().
			Foreground(darkStone).
			Width(5).
			Align(lipgloss.Right).
			PaddingRight(1)

	edGutterActiveStyle = lipgloss.NewStyle().
				Foreground(amber).
				Bold(true).
				Width(5).
				Align(lipgloss.Right).
				PaddingRight(1)

	edGutterSepStyle = lipgloss.NewStyle().
				Foreground(darkStone)

	edTextStyle = lipgloss.NewStyle().
			Foreground(ivory)

	edCursorBlockStyle = lipgloss.NewStyle().
				Reverse(true)

	edStatusBarStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#2A2825")).
				Foreground(stone)

	edModeNormalStyle = lipgloss.NewStyle().
				Background(gold).
				Foreground(lipgloss.Color("#1A1816")).
				Bold(true).
				Padding(0, 1)

	edModeInsertStyle = lipgloss.NewStyle().
				Background(sage).
				Foreground(lipgloss.Color("#1A1816")).
				Bold(true).
				Padding(0, 1)

	edModeCommandStyle = lipgloss.NewStyle().
				Background(amber).
				Foreground(lipgloss.Color("#1A1816")).
				Bold(true).
				Padding(0, 1)

	edFileNameStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#2A2825")).
			Foreground(sand).
			Padding(0, 1)

	edModifiedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#2A2825")).
			Foreground(rust).
			Bold(true).
			Padding(0, 1)

	edPosStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#2A2825")).
			Foreground(stone).
			Padding(0, 1)

	edCmdLineStyle = lipgloss.NewStyle().
			Foreground(sand)

	mdHeadingStyle = lipgloss.NewStyle().
			Foreground(gold).
			Bold(true)

	mdBoldStyle = lipgloss.NewStyle().
			Foreground(sand).
			Bold(true)

	mdItalicStyle = lipgloss.NewStyle().
			Foreground(sand).
			Italic(true)

	mdCodeInlineStyle = lipgloss.NewStyle().
				Foreground(amber)

	mdCodeBlockStyle = lipgloss.NewStyle().
			Foreground(amber)

	mdLinkStyle = lipgloss.NewStyle().
			Foreground(sage)

	mdBlockquoteStyle = lipgloss.NewStyle().
				Foreground(stone).
				Italic(true)

	mdListMarkerStyle = lipgloss.NewStyle().
				Foreground(amber)

	mdHrStyle = lipgloss.NewStyle().
			Foreground(darkStone)
)
