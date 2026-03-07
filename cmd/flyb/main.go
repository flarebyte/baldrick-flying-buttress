package main

import (
	"os"

	"github.com/flarebyte/baldrick-flying-buttress/internal/cli"
	"github.com/flarebyte/baldrick-flying-buttress/internal/load"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
	"github.com/flarebyte/baldrick-flying-buttress/internal/validate"
)

func main() {
	os.Exit(cli.ExecuteWithFactory(
		os.Args[1:],
		os.Stdout,
		os.Stderr,
		func(configPath string) pipeline.AppLoader {
			return load.FSAppLoader{ConfigPath: configPath}
		},
		validate.AppDataValidator{},
	))
}
