package figo

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/gregoryv/web"
)

func Generate(dir, out string) error {
	page := NewPage(Html(Body(
		fileTree(dir),
	)))
	return page.SaveAs(out)
}

func fileTree(dir string) *Element {
	ul := Ul()
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if skipPath(info) {
			return filepath.SkipDir
		}
		ul.With(Li(path))
		return nil
	})
	return ul
}

func skipPath(info os.FileInfo) bool {
	switch {
	case info.IsDir() && info.Name() == ".git":
	case strings.Contains(info.Name(), "~"):
	default:
		return false
	}
	return true
}
