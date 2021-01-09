package figo

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/gregoryv/web"
)

// Generate go documentation for the given directory and its children.
func Generate(dir string) (*Page, error) {
	_, err := golist(dir)
	if err != nil {
		return nil, err
	}

	docs := godoc(dir)
	nav := Nav()
	page := NewPage(Html(
		Head(
			Style(theme()),
		),
		Body(
			nav,
			docs,
		)),
	)
	return page, nil
}

func golist(dir string) (string, error) {
	out, err := exec.Command("go", "list", dir).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s", string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

func godoc(dir string) *Element {
	w := Article()
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if skipPath(info) {
			return filepath.SkipDir
		}
		if info.IsDir() {
			dp := dir + "/" + path
			pkgName, err := golist(dp)
			if err != nil || pkgName == "" {
				return filepath.SkipDir
			}
			s, err := docPkg(pkgName, dp)
			if err != nil {
				return filepath.SkipDir
			}
			w.With(s)
		}
		return nil
	})
	return w
}

func docPkg(pkgName, dir string) (*Element, error) {
	files := make([]*ast.File, 0)
	fset := token.NewFileSet()
	gofiles, _ := filepath.Glob(dir + "/*.go")
	for _, f := range gofiles {
		data, _ := ioutil.ReadFile(f)
		files = append(files, mustParse(fset, f, string(data)))
	}
	p, err := doc.NewFromFiles(fset, files, pkgName)
	if err != nil {
		return nil, err
	}

	s := Section(
		H1(pkgName),
	)
	for _, f := range p.Funcs {
		var buf bytes.Buffer
		printer.Fprint(&buf, fset, f.Decl)
		fn := buf.String()[5:] // remove func
		var class interface{}
		var p interface{}
		if f.Doc == "" {
			//class = Class("empty")
		} else {
			p = P(toHTML(f.Doc))
		}
		s.With(
			Section(Class("func"),
				H2(fn, class),
				p,
			),
		)
	}

	types := make([]string, len(p.Types))
	for i, t := range p.Types {
		types[i] = t.Name
	}
	s.With("Types: ", strings.Join(types, ", "))
	return s, nil
}

func toHTML(v string) string {
	var buf bytes.Buffer
	doc.ToHTML(&buf, v, nil)
	return buf.String()
}

func mustParse(fset *token.FileSet, filename, src string) *ast.File {
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return f
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
