package ui

import (
	"strings"
)

func highlightLine(line string, inCodeBlock bool) (string, bool) {
	trimmed := strings.TrimSpace(line)

	if strings.HasPrefix(trimmed, "```") {
		return mdCodeBlockStyle.Render(line), !inCodeBlock
	}

	if inCodeBlock {
		return mdCodeBlockStyle.Render(line), true
	}

	if trimmed == "---" || trimmed == "***" || trimmed == "___" {
		return mdHrStyle.Render(line), false
	}

	if strings.HasPrefix(trimmed, "# ") || strings.HasPrefix(trimmed, "## ") ||
		strings.HasPrefix(trimmed, "### ") || strings.HasPrefix(trimmed, "#### ") ||
		strings.HasPrefix(trimmed, "##### ") || strings.HasPrefix(trimmed, "###### ") {
		return mdHeadingStyle.Render(line), false
	}

	if strings.HasPrefix(trimmed, "> ") {
		return mdBlockquoteStyle.Render(line), false
	}

	if len(trimmed) >= 2 {
		if (trimmed[0] == '-' || trimmed[0] == '*' || trimmed[0] == '+') && trimmed[1] == ' ' {
			idx := strings.Index(line, trimmed[:2])
			prefix := line[:idx]
			marker := mdListMarkerStyle.Render(string(trimmed[0]))
			rest := highlightInline(line[idx+2:])
			return prefix + marker + " " + rest, false
		}
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

	return highlightInline(line), false
}

func highlightInline(line string) string {
	if line == "" {
		return ""
	}

	var result strings.Builder
	i := 0

	for i < len(line) {
		if line[i] == '`' {
			end := strings.Index(line[i+1:], "`")
			if end >= 0 {
				result.WriteString(mdCodeInlineStyle.Render(line[i : i+end+2]))
				i += end + 2
				continue
			}
		}

		if i+1 < len(line) && line[i] == '*' && line[i+1] == '*' {
			end := strings.Index(line[i+2:], "**")
			if end >= 0 {
				result.WriteString(mdBoldStyle.Render(line[i : i+end+4]))
				i += end + 4
				continue
			}
		}

		if line[i] == '*' && (i == 0 || line[i-1] != '*') {
			end := strings.Index(line[i+1:], "*")
			if end >= 0 && (i+end+2 >= len(line) || line[i+end+2] != '*') {
				result.WriteString(mdItalicStyle.Render(line[i : i+end+2]))
				i += end + 2
				continue
			}
		}

		if line[i] == '_' {
			end := strings.Index(line[i+1:], "_")
			if end >= 0 {
				result.WriteString(mdItalicStyle.Render(line[i : i+end+2]))
				i += end + 2
				continue
			}
		}

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

		result.WriteString(edTextStyle.Render(string(line[i])))
		i++
	}

	return result.String()
}
