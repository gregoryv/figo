package main

import (
	"os"
	"testing"

	"github.com/gregoryv/wolf"
)

func Test_default_behaviour(t *testing.T) {
	cmd := wolf.NewTCmd("figo")
	defer cmd.Cleanup()

	os.Chdir("/home/gregory/dl/go1/go/src/net/http")
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
