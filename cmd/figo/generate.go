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
	pkg, err := golist(dir)
	if err != nil || pkg == "" {
		return nil, fmt.Errorf("%v: not a package", dir)
	}

	docs, err := godoc(pkg, dir)
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

func godoc(pkgName, dir string) (*Element, error) {
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
	pkgExamplesSection := Span()
	s := Article(
		H1("Package ", path.Base(pkgName)),
		Dl(
			Dd(`import "`, pkgName, `"`),
		),
		Dl(
			Dt(A(Href("#pkg-overview"), "Overview")),
			Dt(A(Href("#pkg-index"), "Index")),
			Dt(A(Href("#pkg-examples"), "Examples")),
		),
		Section(
			A(Name("pkg-overview")),
			H2("Overview"),
			toHTML(p.Doc),
			// todo add package examples here
			pkgExamplesSection,
		),
	)

	// Generate index
	dl := Dl()
	examplesIndex := Dl()
	indexSection := Section(
		A(Name("pkg-index")),
		H2("Index"),
		dl,
		Section(
			A(Name("pkg-examples")),
			H3("Examples"),
			examplesIndex,
		),
	)
	docSection := Section(H2("Variables"))
	s.With(indexSection, docSection)

	// Examples index
	for _, ex := range p.Examples {
		name := ex.Name
		id := ex.Name
		if name == "" {
			name = "Package"
			id = "example_"
		}
		examplesIndex.With(Dd(
			A(Href("#"+id), name),
		))
		pkgExamplesSection.With(
			A(Name(id)),
			A("Example"), Br(),
			"Code:", Br(),
			Pre(Code(printHTML(fset, ex.Code))),
			"Output:", Br(),
			Pre(Code(ex.Output)),
		)
	}

	// Package funcs
	for _, f := range p.Funcs {
		lnk := genFuncLink(fset, f)
		dl.With(
			Dd(lnk),
		)
		docSection.With(
			A(Name(f.Name)),
			H3("func ", f.Name),
			Pre(printHTML(fset, f.Decl)),
			P(template.HTMLEscapeString(f.Doc)),
		)
	}

	// Types
	for _, t := range p.Types {
		dl.With(Dd(A(Href("#"+t.Name), "type ", t.Name)))
		docSection.With(
			A(Name(t.Name)),
			H2("type ", t.Name),
			P(template.HTMLEscapeString(t.Doc)),
			Pre(Code(printHTML(fset, t.Decl))),
		)

		for _, f := range t.Funcs {
			dl.With(
				Dd("&nbsp;&nbsp;", genFuncLink(fset, f)),
			)
			docSection.With(
				A(Name(f.Name)),
				H3("func ", f.Name),
				Pre(Code(printHTML(fset, f.Decl))),
				P(template.HTMLEscapeString(f.Doc)),
			)
		}
		for _, f := range t.Methods {
			dl.With(
				Dd("&nbsp;&nbsp;", genFuncLink(fset, f)),
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

// ----------------------------------------

func genFuncLink(fset *token.FileSet, f *doc.Func) interface{} {
	if f.Doc == "" {
		return printHTML(fset, f.Decl)
	}
	return A(Href("#"+f.Name), printHTML(fset, f.Decl))
}

func genTypeLink(fset *token.FileSet, t *doc.Type) interface{} {
	return A(Href("#"+t.Name), t.Name)
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
