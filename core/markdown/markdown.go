package markdown

import (
	"html"
	"regexp"
	"strings"
)

var (
	reHeading = regexp.MustCompile(`(?m)^(#{1,6})\s+(.+)$`)
	reBold    = regexp.MustCompile(`\*\*(.+?)\*\*|__(.+?)__`)
	reItalic  = regexp.MustCompile(`\*(.+?)\*|_(.+?)_`)
	reCode    = regexp.MustCompile("`([^`]+)`")
	reLink    = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	reHR      = regexp.MustCompile(`(?m)^-{3,}$`)
	reUL      = regexp.MustCompile(`(?m)^(?:[-*]\s+.+\n?)+`)
	reOL      = regexp.MustCompile(`(?m)^(?:\d+\.\s+.+\n?)+`)
	reQuote   = regexp.MustCompile(`(?m)^>\s?(.+)$`)
)

// ToHTML converts a small Markdown subset to HTML.
func ToHTML(input string) string {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}

	blocks := strings.Split(input, "\n\n")
	parts := make([]string, 0, len(blocks))
	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}
		parts = append(parts, renderBlock(block))
	}
	return strings.Join(parts, "\n")
}

// ToText strips Markdown markers for a plain-text fallback.
func ToText(input string) string {
	out := input
	out = reHeading.ReplaceAllString(out, "$2")
	out = reBold.ReplaceAllString(out, "$1$2")
	out = reItalic.ReplaceAllString(out, "$1$2")
	out = reCode.ReplaceAllString(out, "$1")
	out = reLink.ReplaceAllString(out, "$1 ($2)")
	out = reHR.ReplaceAllString(out, "")
	out = reQuote.ReplaceAllString(out, "$1")
	out = regexp.MustCompile(`(?m)^[-*]\s+`).ReplaceAllString(out, "• ")
	out = regexp.MustCompile(`(?m)^\d+\.\s+`).ReplaceAllString(out, "")
	return strings.TrimSpace(out)
}

func renderBlock(block string) string {
	if reHR.MatchString(block) && !strings.Contains(block, "\n") {
		return "<hr>"
	}
	if m := reHeading.FindStringSubmatch(block); len(m) == 3 && !strings.Contains(block, "\n") {
		level := len(m[1])
		return "<h" + itoa(level) + ">" + inline(m[2]) + "</h" + itoa(level) + ">"
	}
	if reUL.MatchString(block+"\n") || strings.HasPrefix(block, "- ") || strings.HasPrefix(block, "* ") {
		lines := strings.Split(block, "\n")
		var b strings.Builder
		b.WriteString("<ul>")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			line = strings.TrimPrefix(line, "- ")
			line = strings.TrimPrefix(line, "* ")
			if line == "" {
				continue
			}
			b.WriteString("<li>")
			b.WriteString(inline(line))
			b.WriteString("</li>")
		}
		b.WriteString("</ul>")
		return b.String()
	}
	if reOL.MatchString(block+"\n") || regexp.MustCompile(`^\d+\.\s+`).MatchString(block) {
		lines := strings.Split(block, "\n")
		var b strings.Builder
		b.WriteString("<ol>")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			line = regexp.MustCompile(`^\d+\.\s+`).ReplaceAllString(line, "")
			if line == "" {
				continue
			}
			b.WriteString("<li>")
			b.WriteString(inline(line))
			b.WriteString("</li>")
		}
		b.WriteString("</ol>")
		return b.String()
	}
	if strings.HasPrefix(block, "> ") || strings.HasPrefix(block, ">") {
		lines := strings.Split(block, "\n")
		var b strings.Builder
		b.WriteString("<blockquote><p>")
		parts := make([]string, 0, len(lines))
		for _, line := range lines {
			line = strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "> "), ">"))
			parts = append(parts, inline(line))
		}
		b.WriteString(strings.Join(parts, "<br>"))
		b.WriteString("</p></blockquote>")
		return b.String()
	}
	escapedLines := make([]string, 0)
	for _, line := range strings.Split(block, "\n") {
		escapedLines = append(escapedLines, inline(line))
	}
	return "<p>" + strings.Join(escapedLines, "<br>") + "</p>"
}

func inline(input string) string {
	escaped := html.EscapeString(input)
	// Re-apply markdown on escaped text carefully by matching original patterns on raw then escaping.
	raw := input
	raw = reCode.ReplaceAllStringFunc(raw, func(m string) string {
		match := reCode.FindStringSubmatch(m)
		if len(match) != 2 {
			return html.EscapeString(m)
		}
		return "<code>" + html.EscapeString(match[1]) + "</code>"
	})
	raw = reBold.ReplaceAllStringFunc(raw, func(m string) string {
		match := reBold.FindStringSubmatch(m)
		if len(match) < 3 {
			return html.EscapeString(m)
		}
		text := match[1]
		if text == "" {
			text = match[2]
		}
		return "<strong>" + html.EscapeString(text) + "</strong>"
	})
	raw = reItalic.ReplaceAllStringFunc(raw, func(m string) string {
		match := reItalic.FindStringSubmatch(m)
		if len(match) < 3 {
			return html.EscapeString(m)
		}
		text := match[1]
		if text == "" {
			text = match[2]
		}
		return "<em>" + html.EscapeString(text) + "</em>"
	})
	raw = reLink.ReplaceAllStringFunc(raw, func(m string) string {
		match := reLink.FindStringSubmatch(m)
		if len(match) != 3 {
			return html.EscapeString(m)
		}
		return `<a href="` + html.EscapeString(match[2]) + `">` + html.EscapeString(match[1]) + `</a>`
	})
	if raw != input {
		return raw
	}
	return escaped
}

func itoa(n int) string {
	if n < 1 {
		return "1"
	}
	if n > 6 {
		return "6"
	}
	return string(rune('0' + n))
}
