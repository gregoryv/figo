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
			Div(Class("top"), "Go Documentation", Span(Class("right"), "generated by FiGo1")),
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
	// Prepare sections
	pkgExamplesSection := Span()

	index := Dl()
	examplesIndex := Dl()
	indexSection := Section(
		A(Name("pkg-index")),
		H2("Index"),
		index,
		Section(
			A(Name("pkg-examples")),
			H3("Examples"),
			examplesIndex,
		),
	)
	docSection := Section(H2("Variables"))

	// Examples index
	for _, ex := range p.Examples {
		addExample(examplesIndex, pkgExamplesSection, ex, fset)
	}

	// Package funcs
	for _, f := range p.Funcs {
		lnk := genFuncLink(fset, f)
		index.With(
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
		index.With(Dd(A(Href("#"+t.Name), "type ", t.Name)))
		docSection.With(
			A(Name(t.Name)),
			H2("type ", t.Name),
			P(template.HTMLEscapeString(t.Doc)),
			Pre(Code(printHTML(fset, t.Decl))),
		)
		// Constructors
		for _, f := range t.Funcs {
			index.With(
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
			index.With(
				Dd("&nbsp;&nbsp;", genFuncLink(fset, f)),
			)
			docSection.With(
				A(Name(f.Name)),
				H3("func ", f.Name),
				Pre(Code(printHTML(fset, f.Decl))),
				P(template.HTMLEscapeString(f.Doc)),
			)
		}
		for _, ex := range t.Examples {
			addExample(examplesIndex, pkgExamplesSection, ex, fset)
		}
	}
	s := Article(
		H1("Package ", path.Base(pkgName)),
		Dl(
			Dd(`import "`, pkgName, `"`),
		),
		Dl(
			Dd(A(Href("#pkg-overview"), "Overview")),
			Dd(A(Href("#pkg-index"), "Index")),
			Dd(A(Href("#pkg-examples"), "Examples")),
		),
		Section(
			A(Name("pkg-overview")),
			H2("Overview"),
			toHTML(p.Doc),
			pkgExamplesSection,
		),
		indexSection,
		docSection,
	)
	return s, nil
}

func addExample(index, section *Element, ex *doc.Example, fset *token.FileSet) {
	name := ex.Name
	id := ex.Name
	if name == "" {
		name = "Package"
		id = "example_"
	}
	index.With(Dd(
		A(Href("#"+id), name),
	))
	var output interface{}
	if ex.Output != "" {
		output = Wrap("Output:", Br(),
			Pre(Code(ex.Output)),
		)
	}

	section.With(
		A(Name(id)),
		A("Example"), Br(),
		"Code:", Br(),
		Pre(Code(printHTML(fset, ex.Code))),
		output,
	)
}

// ----------------------------------------

func genFuncLink(fset *token.FileSet, f *doc.Func) interface{} {
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
