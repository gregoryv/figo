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
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gregoryv/web"
	. "github.com/gregoryv/web"
)

// Generate go documentation for the given package.
func Generate(imp string, pkg *doc.Package, fset *token.FileSet) *Page {
	page := NewPage(Html(
		Head(
			Meta(Charset("utf-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
			Title(imp, " - figo"),
			Style(theme()),
		),
		body(imp, pkg, fset),
	))
	return page
}

func body(imp string, pkg *doc.Package, fset *token.FileSet) *Element {
	body := Body(
		Div(
			Class("top"), Span(Class("fi"), "Fi"), " - Go Documentation",
		),
		Span(Class("timestamp"), time.Now().Format("2006-01-02 15:04:05")),
		Article(
			H1("Package ", path.Base(imp)),
			Dl(
				Dd(`import "`, imp, `"`),
			),
			Dl(
				Dd(A(Href("#pkg-overview"), "Overview")),
				Dd(A(Href("#pkg-index"), "Index")),
				Dd(A(Href("#pkg-examples"), "Examples")),
			),
			Section(
				A(Name("pkg-overview")),
				H2("Overview"),
				overview(pkg, fset),
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
				A(Name("pkg-constants")),
				H2("Constants"),
				constants(pkg, fset),
			),
			Section(
				A(Name("pkg-variables")),
				H2("Variables"),
				variables(pkg, fset),
			),
			Section(
				docs(pkg, fset),
			),
		),
	)
	return body
}

func overview(pkg *doc.Package, fset *token.FileSet) *Element {
	return Wrap(
		toHTML(pkg.Doc),
		docExamples(fset, pkg.Examples...),
	)
}

func index(p *doc.Package, fset *token.FileSet) *Element {
	index := Dl(
		Dd(A(Href("#pkg-constants"), "Constants")),
		Dd(A(Href("#pkg-variables"), "Variables")),
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

func constants(pkg *doc.Package, fset *token.FileSet) *Element {
	w := Wrap()
	for _, t := range pkg.Consts {
		w.With(
			toHTML(t.Doc),
			Pre(Code(printHTML(fset, t.Decl))),
		)
	}
	return w
}

func variables(pkg *doc.Package, fset *token.FileSet) *Element {
	w := Wrap()
	for _, t := range pkg.Vars {
		w.With(
			toHTML(t.Doc),
			Pre(Code(colorComments(printHTML(fset, t.Decl)))),
		)
	}
	return w
}

func docs(p *doc.Package, fset *token.FileSet) *Element {
	section := Wrap()
	for _, f := range p.Funcs {
		docFunc(section, p, fset, f)
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
			docFunc(section, p, fset, f)
			section.With(docExamples(fset, f.Examples...))
		}
		for _, f := range t.Methods {
			docFunc(section, p, fset, f)
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

func docFunc(section *Element, pkg *doc.Package, fset *token.FileSet, f *doc.Func) {
	w := Wrap(
		A(Name(f.Name)),
		H3("func ", f.Name),
		Pre(Code(printHTML(fset, f.Decl))),
		P(toHTML(f.Doc)),
	)
	refs := refsTypes(pkg)
	web.WalkElements(w, func(e *Element) {
		if e.Name == "code" {
			web.LinkAll(e, refs)
		}
	})
	section.With(w)
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
			A(Class("title"), title), Br(),
			Span(Class("title"), "Code:"), Br(),
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

func refsTypes(pkg *doc.Package) map[string]string {
	r := make(map[string]string)
	for _, v := range pkg.Types {
		r[v.Name] = "#" + v.Name
	}
	return r
}

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

func colorComments(str string) string {
	re := regexp.MustCompile("([^:]//.*\n)")
	span := []byte(`<span class="comment">$1</span>`)
	b := re.ReplaceAll([]byte(str), span)
	return string(b)
}

func printHTML(fset *token.FileSet, node interface{}, comments ...*ast.CommentGroup) string {
	switch n := node.(type) {
	case *ast.BlockStmt:
		var buf bytes.Buffer
		cnode := &printer.CommentedNode{
			Node:     n,
			Comments: comments, // include comments
		}
		conf := &printer.Config{
			Mode:     printer.UseSpaces,
			Tabwidth: 4,
		}
		conf.Fprint(&buf, fset, cnode)
		re := regexp.MustCompile("\n    ") // starting 4 spaces
		block := re.ReplaceAll(buf.Bytes(), []byte("\n"))
		return colorComments(string(block[2 : len(block)-2]))

	default:
		cnode := &printer.CommentedNode{
			Node:     node,
			Comments: comments,
		}
		conf := &printer.Config{
			Mode:     printer.UseSpaces,
			Tabwidth: 4,
		}
		var buf bytes.Buffer
		conf.Fprint(&buf, fset, cnode)
		return colorComments(buf.String())
	}
}

func toHTML(v string) string {
	var buf bytes.Buffer
	doc.ToHTML(&buf, v, nil)
	return colorComments(buf.String())
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
