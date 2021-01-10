package main

import (
	"testing"
)

func Test_generate(t *testing.T) {
	_, err := Generate(".")
	if err != nil {
		t.Error(err)
	}
}
