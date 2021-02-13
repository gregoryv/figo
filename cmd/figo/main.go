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
	"github.com/gregoryv/wolf"
)

func main() {
	cmd := wolf.NewOSCmd()
	code := run(cmd)
	os.Exit(code)
}

func run(cmd wolf.Command) int {
	var (
		cli  = cmdline.NewParser(cmd.Args()...)
		help = cli.Flag("-h, --help")

		writeToStdout = cli.Flag("-w, --write-to-stdout")
	)

	switch {
	case !cli.Ok():
		return fail(cmd, cli.Error(), 1)

	case help:
		p, _ := nexus.NewPrinter(cmd.Stderr())
		p.Println(cmd.Args()[0], "- generates go documentation to HTML")
		p.Println()
		p.Println("If BROWSER is set, the generated file is automatically opened.")
		cli.WriteUsageTo(p)
		return cmd.Stop(0)
	}

	// Must be a go package
	dir := "."
	imp, err := golist(dir)
	if err != nil || imp == "" {
		return fail(cmd, err, 1)
	}
	fset, files := goFiles(dir)
	pkg, err := doc.NewFromFiles(fset, files, imp)
	if err != nil {
		return fail(cmd, err, 1)
	}

	fidoc := figo.FiDocs{
		Import:  imp,
		Package: pkg,
		FileSet: fset,
	}

	switch {
	case writeToStdout:
		page := fidoc.NewPage()
		page.WriteTo(cmd.Stdout())

	default:
		// Create output file
		tmp := os.TempDir()
		filename := tmp + "/figo_" + pkg.Name + ".html"
		fh, err := os.Create(filename)
		if err != nil {
			return fail(cmd, err, 1)
		}
		defer fh.Close()

		page := fidoc.NewPage()
		_, err = page.WriteTo(fh)
		if err != nil {
			return fail(cmd, err, 1)
		}
		browser := cmd.Getenv("BROWSER")
		if browser != "" {
			exec.Command(browser, filename).Run()
		}
		fmt.Println(filename)
	}
	return cmd.Stop(0)
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

func fail(cmd wolf.Command, err error, exitCode int) int {
	fmt.Fprintln(cmd.Stderr(), err)
	return cmd.Stop(exitCode)
}
