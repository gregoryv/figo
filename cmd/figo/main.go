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
	)

	switch {
	case !cli.Ok():
		fmt.Fprintln(cmd.Stderr(), cli.Error())
		return cmd.Stop(1)

	case help:
		p, _ := nexus.NewPrinter(cmd.Stderr())
		p.Println(
			cmd.Args()[0], "- generates go documentation of the current working directory to stdout.",
		)
		p.Println("Written by Gregory Vincic <g@7de.se>")
		return cmd.Stop(0)
	}

	page, err := Generate(".")
	if err != nil {
		fmt.Fprintln(cmd.Stderr(), err)
		return cmd.Stop(1)
	}

	page.WriteTo(cmd.Stdout())
	return cmd.Stop(0)
}
