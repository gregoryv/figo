package main

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
	"path"
	"path/filepath"
	"strings"

	. "github.com/gregoryv/web"
)

// Generate go documentation for the given directory and its children.
func Generate(dir string) (p *Page, err error) {
	if !isPackage(dir) {
		return nil, fmt.Errorf("%v: not a package", dir)
	}

	docs, err := godoc(dir)
	if err != nil {
		return nil, err
	}

	p = NewPage(Html(
		Head(
			Meta(Charset("utf-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
			Meta(Name("theme-color"), Content("#375EAB")),
			Style(theme()),
		),
		Body(
			Div(Class("top"), "FiGo1"),
			docs,
		)),
	)
	return
}

func godoc(dir string) (*Element, error) {
	w := Article()
	pkgName, err := golist(dir)
	if err != nil {
		return nil, err
	}

	s, err := docPkg(pkgName, dir)
	if err != nil {
		return nil, err
	}
	w.With(s)
	return w, nil
}

func docPkg(pkgName, dir string) (*Element, error) {
	// Parse files
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
	// Build section
	s := Section(
		H1("Package ", path.Base(pkgName)),
		Dl(
			Dt(`import "`, pkgName, `"`),
			Dt("Overview"),
			Dt("Index"),
			Dt("Examples"),
		),
	)

	dl := Dl()
	indexSection := Section(
		H2("Index"),
		dl,
	)
	s.With(indexSection)
	for _, f := range p.Funcs {
		dl.With(
			Dd(printHTML(fset, f.Decl)),
		)
	}

	for _, t := range p.Types {
		dl.With(Dd(t.Name))
		for _, f := range t.Funcs {
			dl.With(
				Dd("&nbsp;&nbsp;", printHTML(fset, f.Decl)),
			)
		}
	}
	return s, nil
}

func isPackage(dir string) bool {
	name, err := golist(dir)
	return err == nil && name != ""
}

func golist(dir string) (string, error) {
	out, err := exec.Command("go", "list", dir).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s", string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

func printHTML(fset *token.FileSet, node interface{}) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, node)
	return buf.String()
}

func funcNames(funcs []*doc.Func) []string {
	names := make([]string, len(funcs))
	for i, f := range funcs {
		names[i] = f.Name
	}
	return names
}

func addFunc(s *Element, fset *token.FileSet, funcs []*doc.Func) {
	for _, f := range funcs {
		fn := printHTML(fset, f.Decl)
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
