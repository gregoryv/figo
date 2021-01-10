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
	"text/template"

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
	// Generate index
	dl := Dl()
	indexSection := Section(
		H2("Index"),
		dl,
	)
	docSection := Section(
		H2("Variables"),
	)
	s.With(indexSection, docSection)
	for _, f := range p.Funcs {
		lnk := genFunc(fset, f)
		dl.With(
			Dd(lnk),
		)
		docSection.With(
			H3("func ", f.Name),
			Pre(printHTML(fset, f.Decl)),
			P(f.Doc),
		)
	}

	for _, t := range p.Types {
		dl.With(Dd("type ", t.Name))
		docSection.With(
			H2("type ", t.Name),
			P(template.HTMLEscapeString(t.Doc)),
			Pre(Code(printHTML(fset, t))),
		)

		for _, f := range t.Funcs {
			dl.With(
				Dd("&nbsp;&nbsp;", genFunc(fset, f)),
			)
			docSection.With(
				A(Name(f.Name)),
				H3("func ", f.Name),
				Pre(Code(printHTML(fset, f.Decl))),
				P(template.HTMLEscapeString(f.Doc)),
			)
		}
	}
	return s, nil
}

func genFunc(fset *token.FileSet, f *doc.Func) interface{} {
	if f.Doc == "" {
		return printHTML(fset, f.Decl)
	}
	return A(Href("#"+f.Name), printHTML(fset, f.Decl))
}

func printHTML(fset *token.FileSet, node interface{}) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, node)
	return buf.String()
}

func toHTML(v string) string {
	var buf bytes.Buffer
	doc.ToHTML(&buf, v, nil)
	return buf.String()
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
