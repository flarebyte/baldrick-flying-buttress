package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/flarebyte/baldrick-flying-buttress/internal/buildinfo"
	"github.com/flarebyte/baldrick-flying-buttress/internal/cli"
	"github.com/flarebyte/baldrick-flying-buttress/internal/load"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
	"github.com/flarebyte/baldrick-flying-buttress/internal/validate"
)

func main() {
	_ = buildinfo.Version
	_ = buildinfo.Commit
	_ = buildinfo.Date

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	os.Exit(cli.ExecuteContextWithFactory(
		ctx,
		os.Args[1:],
		os.Stdout,
		os.Stderr,
		func(configPath string) pipeline.AppLoader {
			return load.FSAppLoader{ConfigPath: configPath}
		},
		validate.AppDataValidator{},
	))
}
