package main

import (
	"os"

	"github.com/flarebyte/baldrick-flying-buttress/internal/cli"
	"github.com/flarebyte/baldrick-flying-buttress/internal/validate"
)

func main() {
	os.Exit(cli.Execute(os.Args[1:], os.Stdout, os.Stderr, validate.LoadStub, validate.ValidateStub))
}
