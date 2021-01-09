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
	"regexp"
	"strings"

	"github.com/gregoryv/web"
	. "github.com/gregoryv/web"
	"github.com/gregoryv/web/toc"
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
	//MakeTOC(nav, docs, "h1", "h2", "h3")
	return page, nil
}

func MakeTOC(dest, root *web.Element, names ...string) *web.Element {
	toc.GenerateIDs(root, names...)
	toc.GenerateAnchors(root, names...)
	ul := ParseTOC(root, names...)
	dest.With(ul)
	return ul
}

func ParseTOC(root *web.Element, names ...string) *web.Element {
	ul := web.Ul()
	web.WalkElements(root, func(e *web.Element) {
		for _, name := range names {
			if e.Name == name {
				if hasClass(e, "empty") {
					ul.With(web.Li(web.Class(name), e.Text()))
					continue
				}

				a := web.A(web.Href("#"+idOf(e)), e.Text())
				ul.With(web.Li(web.Class(name), a))
			}
		}
	})
	return ul
}

func hasClass(e *Element, class string) bool {
	for _, attr := range e.Attributes {
		if attr.Name == "class" {
			return strings.Contains(attr.Val, class)
		}
	}
	return false
}

func idOf(e *Element) string {
	for _, attr := range e.Attributes {
		if attr.Name == "id" {
			return attr.Val
		}
	}
	txt := idChars.ReplaceAllString(e.Text(), "")
	return strings.ToLower(txt)
}

var idChars = regexp.MustCompile(`\W`)

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

func fileTree(dir string) *Element {
	ul := Ul()
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if skipPath(info) {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			ul.With(Li(path))
		}
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
