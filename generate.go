package figo

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/gregoryv/web"
)

func Generate(dir, out string) error {
	pkg, err := golist(dir)
	if err != nil {
		return err
	}

	page := NewPage(Html(Body(
		H1(pkg),
		fileTree(dir),
	)))
	return page.SaveAs(out)
}

func golist(dir string) (string, error) {
	out, err := exec.Command("go", "list", dir).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
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
