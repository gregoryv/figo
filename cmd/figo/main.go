package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/figo"
	"github.com/gregoryv/nexus"
)

// Set during build
var retailInfo string = "none"

func main() {
	var (
		cli           = cmdline.NewBasicParser()
		rinfo         = cli.Option("--retail-info", "Show who bought this software").Bool()
		writeToStdout = cli.Flag("-w, --write-to-stdout")
	)
	u := cli.Usage()
	u.Preface(
		"figo - generates go documetation to HTML",
		"",
		"If BROWSER is set, the generated file is automatically opened.",
	)

	cli.Parse()
	sh := cmdline.DefaultShell

	if rinfo {
		p, _ := nexus.NewPrinter(sh.Stdout())
		p.Println(retailInfo)
		sh.Exit(0)
		return
	}

	// Must be a go package
	dir := "."
	imp, err := golist(dir)
	if err != nil || imp == "" {
		sh.Fatal(err)
		return
	}
	fset, files := goFiles(dir)
	pkg, err := doc.NewFromFiles(fset, files, imp)
	if err != nil {
		sh.Fatal(err)
		return
	}

	fidoc := figo.FiDocs{
		Import:  imp,
		Package: pkg,
		FileSet: fset,
	}

	switch {
	case writeToStdout:
		page := fidoc.NewPage()
		page.WriteTo(sh.Stdout())

	default:
		// Create output file
		tmp := os.TempDir()
		filename := tmp + "/figo_" + pkg.Name + ".html"
		fh, err := os.Create(filename)
		if err != nil {
			sh.Fatal(err)
			return
		}
		defer fh.Close()

		page := fidoc.NewPage()
		_, err = page.WriteTo(fh)
		if err != nil {
			sh.Fatal(err)
			return
		}
		browser := sh.Getenv("BROWSER")
		if browser != "" {
			exec.Command(browser, filename).Run()
		}
		fmt.Println(filename)
	}
}

func goFiles(dir string) (*token.FileSet, []*ast.File) {
	files := make([]*ast.File, 0)
	fset := token.NewFileSet()
	gofiles, _ := filepath.Glob(dir + "/*.go")
	for _, f := range gofiles {
		data, _ := ioutil.ReadFile(f)
		if bytes.Contains(data, []byte("+build ignore")) {
			continue
		}
		files = append(files, mustParse(fset, f, string(data)))
	}
	return fset, files
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
