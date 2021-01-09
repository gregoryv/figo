package main

import (
	"os"
	"testing"

	"github.com/gregoryv/wolf"
)

func Test_default_behaviour(t *testing.T) {
	orig, _ := os.Getwd()
	cmd := wolf.NewTCmd("figo")
	defer cmd.Cleanup()
	defer os.RemoveAll("docs.html")

	os.Chdir(orig)
	code := run(cmd)
	if code != 0 {
		t.Fail()
	}
}

func Test_figo_help(t *testing.T) {
	cmd := wolf.NewTCmd("figo", "-h")
	defer cmd.Cleanup()

	code := run(cmd)
	if code != 0 {
		t.Error(cmd.Dump())
	}
}

func Test_bad_flag(t *testing.T) {
	cmd := wolf.NewTCmd("figo", "-no-such")
	defer cmd.Cleanup()

	code := run(cmd)
	if code != 1 {
		t.Error(cmd.Dump())
	}
}
