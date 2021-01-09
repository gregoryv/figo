package main

import (
	"fmt"
	"os"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/figo"
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
		dir  = "."
	)

	switch {
	case !cli.Ok():
		fmt.Fprintln(cmd.Stderr(), cli.Error())
		return cmd.Stop(1)

	case help:
		cli.WriteUsageTo(cmd.Stderr())
		return cmd.Stop(0)
	}

	fmt.Println(dir)
	page, err := figo.Generate(dir)
	if err != nil {
		fmt.Fprintln(cmd.Stderr(), err)
		return cmd.Stop(1)
	}

	page.WriteTo(cmd.Stdout())
	return cmd.Stop(0)
}
