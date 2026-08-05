package docs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zatrano/framework/core/docs"
)

func TestDocsListAndGet(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "index.md"), []byte("# Hello\n\nWelcome"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "routing.md"), []byte("# Routing\n\nGroups"), 0o644)
	repo := docs.New(dir)
	pages, err := repo.List()
	if err != nil || len(pages) != 2 {
		t.Fatalf("pages=%v err=%v", pages, err)
	}
	html, page, err := repo.HTML("index")
	if err != nil || page.Title != "Hello" || html == "" {
		t.Fatalf("html=%q page=%#v err=%v", html, page, err)
	}
}
