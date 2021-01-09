package figo

import (
	"testing"
)

func Test_generate(t *testing.T) {
	err := Generate(".", "./docs/index.html")
	if err != nil {
		t.Error(err)
	}
}
