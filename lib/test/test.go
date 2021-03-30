package test

import (
	"flag"
	"os"
	"testing"
)

var long = flag.Bool("long", false, "enable long running tests")

func Long() bool {
	return *long
}

func Quiet() func() {
	stdout := os.Stdout
	os.Stdout = nil
	return func() { os.Stdout = stdout }
}

func RequireLong(t *testing.T) {
	if !Long() {
		t.Skip("long test: use -long to enable")
	}
}

func RequireLongBenchmark(b *testing.B) {
	if !Long() {
		b.Skip("long test: use -long to enable")
	}
}
