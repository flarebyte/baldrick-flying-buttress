package cli

import (
	"errors"
	"fmt"
	"io"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
	clioutput "github.com/flarebyte/baldrick-flying-buttress/internal/output"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
	"github.com/spf13/cobra"
)

const defaultConfigPath = "testdata/app.raw.json"

type LoaderFactory func(configPath string) pipeline.AppLoader

func NewRootCmd(loader pipeline.AppLoader, validator pipeline.AppValidator) *cobra.Command {
	return NewRootCmdWithFactory(func(string) pipeline.AppLoader {
		return loader
	}, validator)
}

func NewRootCmdWithFactory(loaderFactory LoaderFactory, validator pipeline.AppValidator) *cobra.Command {
	var configPath string
	cmd := &cobra.Command{
		Use:           "flyb",
		Short:         "Baldrick Flying Buttress CLI",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.PersistentFlags().StringVar(&configPath, "config", defaultConfigPath, "Path to raw app config file")

	cmd.AddCommand(newValidateCmd(loaderFactory, validator, &configPath))
	cmd.AddCommand(newListCmd(loaderFactory, validator, &configPath))
	cmd.AddCommand(newGenerateCmd(loaderFactory, validator, &configPath))
	return cmd
}

func Execute(args []string, out io.Writer, errOut io.Writer, loader pipeline.AppLoader, validator pipeline.AppValidator) int {
	return ExecuteWithFactory(args, out, errOut, func(string) pipeline.AppLoader {
		return loader
	}, validator)
}

func ExecuteWithFactory(args []string, out io.Writer, errOut io.Writer, loaderFactory LoaderFactory, validator pipeline.AppValidator) int {
	cmd := NewRootCmdWithFactory(loaderFactory, validator)
	cmd.SetOut(out)
	cmd.SetErr(errOut)
	cmd.SetArgs(args)

	err := cmd.Execute()
	exec := outcome.FromError(err)
	if exec.Kind == outcome.KindRuntimeFailure {
		_, _ = fmt.Fprintln(errOut, exec.Err.Error())
	}
	return exec.ExitCode()
}

func newValidateCmd(loaderFactory LoaderFactory, validator pipeline.AppValidator, configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWithConfig(loaderFactory, validator, configPath, validateAction{out: cmd.OutOrStdout()})
		},
	}
}

func newListCmd(loaderFactory LoaderFactory, validator pipeline.AppValidator, configPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List entities",
	}
	cmd.AddCommand(newListReportsCmd(loaderFactory, validator, configPath))
	return cmd
}

func newListReportsCmd(loaderFactory LoaderFactory, validator pipeline.AppValidator, configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "reports",
		Short: "List reports",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWithConfig(loaderFactory, validator, configPath, listReportsAction{out: cmd.OutOrStdout()})
		},
	}
}

func newGenerateCmd(loaderFactory LoaderFactory, validator pipeline.AppValidator, configPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate artifacts",
	}
	cmd.AddCommand(newGenerateJSONCmd(loaderFactory, validator, configPath))
	return cmd
}

func newGenerateJSONCmd(loaderFactory LoaderFactory, validator pipeline.AppValidator, configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "json",
		Short: "Generate JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWithConfig(loaderFactory, validator, configPath, generateJSONAction{out: cmd.OutOrStdout()})
		},
	}
}

func runWithConfig(loaderFactory LoaderFactory, validator pipeline.AppValidator, configPath *string, action pipeline.CommandAction) error {
	if loaderFactory == nil {
		return errors.New("loader factory is required")
	}
	if configPath == nil {
		return errors.New("config path is required")
	}
	loader := loaderFactory(*configPath)
	return pipeline.Run(loader, validator, action)
}

type validateAction struct {
	out io.Writer
}

func (a validateAction) Execute(validated domain.ValidatedApp, report domain.ValidationReport) error {
	_ = validated
	if err := clioutput.EmitDiagnostics(a.out, report); err != nil {
		return err
	}
	if report.HasErrors() {
		return outcome.ValidationBlockedError()
	}
	return nil
}

func (validateAction) AllowOnValidationErrors() bool {
	return true
}

type listReportsAction struct {
	out io.Writer
}

func (a listReportsAction) Execute(validated domain.ValidatedApp, report domain.ValidationReport) error {
	_ = report
	return clioutput.EmitReportList(a.out, validated)
}

func (listReportsAction) AllowOnValidationErrors() bool {
	return false
}

type generateJSONAction struct {
	out io.Writer
}

func (a generateJSONAction) Execute(validated domain.ValidatedApp, report domain.ValidationReport) error {
	_ = report
	return clioutput.EmitGraphJSON(a.out, validated)
}

func (generateJSONAction) AllowOnValidationErrors() bool {
	return false
}
