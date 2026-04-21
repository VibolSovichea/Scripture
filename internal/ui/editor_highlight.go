package ui

import (
	"strings"
)

// highlightLine applies markdown syntax highlighting to a single line.
// inCodeBlock tracks fenced code block state across lines.
// Returns the styled string and updated code block state.
func highlightLine(line string, inCodeBlock bool) (string, bool) {
	trimmed := strings.TrimSpace(line)

	// Fenced code block toggle
	if strings.HasPrefix(trimmed, "```") {
		return mdCodeBlockStyle.Render(line), !inCodeBlock
	}

	// Inside code block — render entire line as code
	if inCodeBlock {
		return mdCodeBlockStyle.Render(line), true
	}

	// Horizontal rule
	if trimmed == "---" || trimmed == "***" || trimmed == "___" {
		return mdHrStyle.Render(line), false
	}

	// Headings
	if strings.HasPrefix(trimmed, "# ") || strings.HasPrefix(trimmed, "## ") ||
		strings.HasPrefix(trimmed, "### ") || strings.HasPrefix(trimmed, "#### ") ||
		strings.HasPrefix(trimmed, "##### ") || strings.HasPrefix(trimmed, "###### ") {
		return mdHeadingStyle.Render(line), false
	}

	// Blockquote
	if strings.HasPrefix(trimmed, "> ") {
		return mdBlockquoteStyle.Render(line), false
	}

	// List items — style the marker, then highlight rest inline
	if len(trimmed) >= 2 {
		if (trimmed[0] == '-' || trimmed[0] == '*' || trimmed[0] == '+') && trimmed[1] == ' ' {
			idx := strings.Index(line, trimmed[:2])
			prefix := line[:idx]
			marker := mdListMarkerStyle.Render(string(trimmed[0]))
			rest := highlightInline(line[idx+2:])
			return prefix + marker + " " + rest, false
		}
		// Numbered list
		for i, c := range trimmed {
			if c == '.' && i > 0 && i < len(trimmed)-1 && trimmed[i+1] == ' ' {
				allDigits := true
				for _, d := range trimmed[:i] {
					if d < '0' || d > '9' {
						allDigits = false
						break
					}
				}
				if allDigits {
					idx := strings.Index(line, trimmed[:i+2])
					prefix := line[:idx]
					marker := mdListMarkerStyle.Render(trimmed[:i+1])
					rest := highlightInline(line[idx+i+2:])
					return prefix + marker + " " + rest, false
				}
				break
			}
		}
	}

	// Regular line — apply inline highlighting
	return highlightInline(line), false
}

// highlightInline applies inline markdown highlighting (bold, italic, code, links).
func highlightInline(line string) string {
	if line == "" {
		return ""
	}

	var result strings.Builder
	i := 0

	for i < len(line) {
		// Inline code: `code`
		if line[i] == '`' {
			end := strings.Index(line[i+1:], "`")
			if end >= 0 {
				result.WriteString(mdCodeInlineStyle.Render(line[i : i+end+2]))
				i += end + 2
				continue
			}
		}

		// Bold: **text**
		if i+1 < len(line) && line[i] == '*' && line[i+1] == '*' {
			end := strings.Index(line[i+2:], "**")
			if end >= 0 {
				result.WriteString(mdBoldStyle.Render(line[i : i+end+4]))
				i += end + 4
				continue
			}
		}

		// Italic: *text* (single asterisk, not preceded by another *)
		if line[i] == '*' && (i == 0 || line[i-1] != '*') {
			end := strings.Index(line[i+1:], "*")
			if end >= 0 && (i+end+2 >= len(line) || line[i+end+2] != '*') {
				result.WriteString(mdItalicStyle.Render(line[i : i+end+2]))
				i += end + 2
				continue
			}
		}

		// Italic: _text_
		if line[i] == '_' {
			end := strings.Index(line[i+1:], "_")
			if end >= 0 {
				result.WriteString(mdItalicStyle.Render(line[i : i+end+2]))
				i += end + 2
				continue
			}
		}

		// Link: [text](url)
		if line[i] == '[' {
			closeBracket := strings.Index(line[i:], "](")
			if closeBracket >= 0 {
				closeParen := strings.Index(line[i+closeBracket:], ")")
				if closeParen >= 0 {
					linkEnd := i + closeBracket + closeParen + 1
					result.WriteString(mdLinkStyle.Render(line[i:linkEnd]))
					i = linkEnd
					continue
				}
			}
		}

		// Regular character
		result.WriteString(edTextStyle.Render(string(line[i])))
		i++
	}

	return result.String()
}
