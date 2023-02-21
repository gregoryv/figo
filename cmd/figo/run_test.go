package main

import (
	"os"
	"runtime"
	"testing"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/cmdline/clitest"
)

func Test_default_behaviour(t *testing.T) {
	sh := newShellT(t, "figo")

	os.Chdir(runtime.GOROOT() + "/src/net/http")
	main()
	if sh.ExitCode != 0 {
		t.Error(sh.Dump())
	}
}

func Test_figo_help(t *testing.T) {
	sh := newShellT(t, "figo", "-h")

	main()
	if sh.ExitCode != 0 {
		t.Error(sh.Dump())
	}
}

func Test_bad_flag(t *testing.T) {
	sh := newShellT(t, "figo", "-no-such")
	main()
	if sh.ExitCode != 1 {
		t.Error(sh.Dump())
	}
}

func Test_write_to_stdout(t *testing.T) {
	sh := newShellT(t, "figo", "-w")

	os.Chdir("/home/gregory/dl/go1/go/src/net/http")
	main()
	if sh.ExitCode != 0 {
		t.Error(sh.Dump())
	}
}

func newShellT(t *testing.T, args ...string) *clitest.ShellT {
	sh := clitest.NewShellT(args...)
	cmdline.DefaultShell = sh
	t.Cleanup(sh.Cleanup)
	return sh
}
