// purpose: Bootstraps the flyb CLI process and wires runtime dependencies for command execution.
// responsibilities: initialize build metadata references; create cancellable OS-signal context; construct loader/validator wiring; execute root CLI and return process exit code
// architecture_notes: Dependency wiring is centralized here so internal packages stay decoupled from process lifecycle concerns and signal handling.
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
