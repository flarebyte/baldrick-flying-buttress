package main

import (
	"os"

	"github.com/olivier/baldrick-flying-buttress/internal/cli"
	"github.com/olivier/baldrick-flying-buttress/internal/validate"
)

func main() {
	os.Exit(cli.Execute(os.Args[1:], os.Stdout, os.Stderr, validate.RunStub))
}
