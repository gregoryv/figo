package main

import (
	"fmt"
	"os"

	"github.com/gregoryv/cmdline"
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
		p.Println(
			cmd.Args()[0], "- generates go documentation to HTML",
		)
		cli.WriteUsageTo(p)

	case writeToStdout:
		page, err := Generate(".")
		if err != nil {
			return fail(cmd, err, 1)
		}
		page.WriteTo(cmd.Stdout())

	default:
		// Create output file
		tmp := os.TempDir()
		filename := tmp + "/figo.html"
		fh, err := os.Create(filename)
		if err != nil {
			return fail(cmd, err, 1)
		}
		defer fh.Close()

		page, err := Generate(".")
		if err != nil {
			return fail(cmd, err, 1)
		}

		_, err = page.WriteTo(fh)
		if err != nil {
			return fail(cmd, err, 1)
		}
		fmt.Println(filename)
	}
	return cmd.Stop(0)
}

func fail(cmd wolf.Command, err error, exitCode int) int {
	fmt.Fprintln(cmd.Stderr(), err)
	return cmd.Stop(exitCode)
}
