package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/printer"
	"go/token"
	"os/exec"
	"path"
	"strings"
	"text/template"
	"time"

	. "github.com/gregoryv/web"
)

// Generate go documentation for the given directory and its children.
func Generate(pkg string, p *doc.Package, fset *token.FileSet) (page *Page, err error) {
	stamp := time.Now().Format("2006-01-02 15:04:05")
	page = NewPage(Html(
		Head(
			Meta(Charset("utf-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
			Meta(Name("theme-color"), Content("#375EAB")),
			Style(theme()),
		),
		Body(
			Div(Class("top"), "Go Documentation", Span(Class("generated"), "generated by FiGo1, ", stamp)),
			Article(
				H1("Package ", path.Base(pkg)),
				Dl(
					Dd(`import "`, pkg, `"`),
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
				),
				index(p, fset),
				docs(p, fset),
			),
		)),
	)
	return
}

func index(p *doc.Package, fset *token.FileSet) *Element {
	index := Dl(
		funcLinks(fset, p.Funcs...),
	)

	for _, t := range p.Types {
		index.With(Dd(A(Href("#"+t.Name), "type ", t.Name)))
		index.With(funcLinks(fset, t.Funcs...)) // Constructors
		for _, f := range t.Methods {
			index.With(Dd(Class("method"), genFuncLink(fset, f)))
		}
	}

	return Section(
		A(Name("pkg-index")),
		H2("Index"),
		index,
	)
}

func docs(p *doc.Package, fset *token.FileSet) *Element {
	section := Section()

	// Package funcs
	for _, f := range p.Funcs {
		docFunc(section, fset, f)
	}

	// Types
	for _, t := range p.Types {
		section.With(
			A(Name(t.Name)),
			H2("type ", t.Name),
			toHTML(t.Doc),
			Pre(Code(printHTML(fset, t.Decl))),
		)
		// Constructors
		for _, f := range t.Funcs {
			docFunc(section, fset, f)
		}
		for _, f := range t.Methods {
			docFunc(section, fset, f)
		}
	}
	return section
}

func docFunc(section *Element, fset *token.FileSet, f *doc.Func) {
	section.With(
		A(Name(f.Name)),
		H3("func ", f.Name),
		Pre(Code(printHTML(fset, f.Decl))),
		P(template.HTMLEscapeString(f.Doc)),
	)
}

func docExample(index, section *Element, fset *token.FileSet, examples ...*doc.Example) {
	for _, ex := range examples {
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
}

// ----------------------------------------

func funcLinks(fset *token.FileSet, funcs ...*doc.Func) *Element {
	el := Wrap()
	for _, f := range funcs {
		el.With(Dd(genFuncLink(fset, f)))
	}
	return el
}

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
