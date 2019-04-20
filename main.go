package main

import (
	"fmt"
	"os"

	"github.com/mui87/gogo/app"
)

const (
	exitCodeOK = iota
	exitCodeError
)

func main() {
	a, err := app.New(os.Args, os.Stdout, os.Stderr)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		os.Exit(exitCodeError)
	}

	if err := a.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		os.Exit(exitCodeError)
	}

	os.Exit(exitCodeOK)
}
