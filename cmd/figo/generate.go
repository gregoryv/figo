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
	"sort"
	"strings"
	"time"

	. "github.com/gregoryv/web"
)

// Generate go documentation for the given directory and its children.
func Generate(name string, pkg *doc.Package, fset *token.FileSet) (page *Page, err error) {
	stamp := time.Now().Format("2006-01-02 15:04:05")
	page = NewPage(Html(
		Head(
			Meta(Charset("utf-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
			Style(theme()),
		),
		Body(
			Div(
				Class("top"), "Go Documentation",
				Span(Class("generated"), "generated by figo, ", stamp),
			),
			Article(
				H1("Package ", path.Base(name)),
				Dl(
					Dd(`import "`, name, `"`),
				),
				Dl(
					Dd(A(Href("#pkg-overview"), "Overview")),
					Dd(A(Href("#pkg-index"), "Index")),
					Dd(A(Href("#pkg-examples"), "Examples")),
				),
				Section(
					A(Name("pkg-overview")),
					H2("Overview"),
					toHTML(pkg.Doc),
					docExamples(fset, pkg.Examples...),
				),
				Section(
					A(Name("pkg-index")),
					H2("Index"),
					index(pkg, fset),

					A(Name("pkg-examples")),
					H3("Examples"),
					examples(pkg, fset),

					H3("Package files"),
					packageFiles(fset),
				),
				Section(
					docs(pkg, fset),
				),
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
	return index
}

func examples(pkg *doc.Package, fset *token.FileSet) *Element {
	dl := Dl()
	for _, ex := range allExamples(pkg, fset) {
		dl.With(Dd(
			A(
				Href("#"+exampleId(ex)),
				ex.Name,
			),
		))
	}
	return dl
}

func packageFiles(fset *token.FileSet) *Element {
	names := make([]string, 0)
	fset.Iterate(func(f *token.File) bool {
		names = append(names, f.Name())
		return true
	})

	v := P()
	for _, name := range names {
		if strings.Contains(name, "_test.go") {
			continue
		}
		v.With(name, " ")
	}
	return v
}

func docs(p *doc.Package, fset *token.FileSet) *Element {
	section := Wrap()
	for _, f := range p.Funcs {
		docFunc(section, fset, f)
		section.With(docExamples(fset, f.Examples...))
	}
	for _, t := range p.Types {
		section.With(
			A(Name(t.Name)),
			H2("type ", t.Name),
			toHTML(t.Doc),
			Pre(Code(printHTML(fset, t.Decl))),
			docExamples(fset, t.Examples...),
		)
		// Constructors
		for _, f := range t.Funcs {
			docFunc(section, fset, f)
			section.With(docExamples(fset, f.Examples...))
		}
		for _, f := range t.Methods {
			docFunc(section, fset, f)
			section.With(docExamples(fset, f.Examples...))
		}
	}
	return section
}

func allExamples(pkg *doc.Package, fset *token.FileSet) []*doc.Example {
	all := make([]*doc.Example, 0, len(pkg.Examples))
	for _, ex := range pkg.Examples {
		all = append(all, ex)
	}
	for _, f := range pkg.Funcs {
		for _, ex := range f.Examples {
			all = append(all, ex)
		}
	}
	for _, t := range pkg.Types {
		for _, ex := range t.Examples {
			all = append(all, ex)
		}
		for _, f := range t.Funcs {
			for _, ex := range f.Examples {
				all = append(all, ex)
			}
		}
		for _, f := range t.Methods {
			for _, ex := range f.Examples {
				all = append(all, ex)
			}
		}
	}
	sort.Sort(exampleByName(all))
	return all
}

type exampleByName []*doc.Example

func (a exampleByName) Len() int           { return len(a) }
func (a exampleByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a exampleByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func docFunc(section *Element, fset *token.FileSet, f *doc.Func) {
	section.With(
		A(Name(f.Name)),
		H3("func ", f.Name),
		Pre(Code(printHTML(fset, f.Decl))),
		P(toHTML(f.Doc)),
	)
}

func docExamples(fset *token.FileSet, examples ...*doc.Example) *Element {
	el := Wrap()
	for _, ex := range examples {
		var output interface{}
		if ex.Output != "" {
			output = Wrap("Output:", Br(),
				Pre(Code(ex.Output)),
			)
		}
		title := "Example"
		if ex.Suffix != "" {
			title = fmt.Sprintf("Example (%s)", strings.Title(ex.Suffix))
		}
		el.With(
			A(Name(exampleId(ex))),
			A(title), Br(),
			"Code:", Br(),
			Pre(Code(printHTML(fset, ex.Code))),
			output,
		)
	}
	return el
}

func exampleId(ex *doc.Example) string {
	return "example_" + ex.Name + "_" + ex.Suffix
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
