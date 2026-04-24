package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_ViewFilesContainZatranoSyntax(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := Run(Options{
		Dir:     dir,
		AppName: "myapp",
		Module:  "github.com/x/myapp",
	}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "views", "layouts", "app.html"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if !strings.Contains(s, `{{assetLink "css/app.css"}}`) {
		t.Fatalf("expected Zatrano assetLink in app layout, got:\n%s", s)
	}
}
