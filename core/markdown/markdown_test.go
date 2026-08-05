package markdown_test

import (
	"strings"
	"testing"

	"github.com/zatrano/framework/core/markdown"
)

func TestMarkdownToHTML(t *testing.T) {
	html := markdown.ToHTML("# Hello\n\nThis is **bold** and a [link](https://zatrano.test).")
	if !strings.Contains(html, "<h1>Hello</h1>") {
		t.Fatalf("missing h1: %s", html)
	}
	if !strings.Contains(html, "<strong>bold</strong>") {
		t.Fatalf("missing bold: %s", html)
	}
	if !strings.Contains(html, `<a href="https://zatrano.test">link</a>`) {
		t.Fatalf("missing link: %s", html)
	}
}
