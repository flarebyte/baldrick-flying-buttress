package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"runtime"

	"github.com/flarebyte/baldrick-flying-buttress/internal/buildinfo"
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
	var showVersion bool
	cmd := &cobra.Command{
		Use:   "flyb",
		Short: "Build, inspect, lint, and generate structured graph-driven reports",
		Long: "flyb validates application graph configuration and provides deterministic\n" +
			"commands to list entities, run lint checks, and generate JSON or markdown outputs.",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				return emitVersionJSON(cmd.OutOrStdout())
			}
			return cmd.Help()
		},
	}
	cmd.PersistentFlags().StringVar(&configPath, "config", defaultConfigPath, "Path to app config file or directory")
	cmd.PersistentFlags().BoolVar(&showVersion, "version", false, "Print detailed version metadata as JSON")

	cmd.AddCommand(newValidateCmd(loaderFactory, validator, &configPath))
	cmd.AddCommand(newListCmd(loaderFactory, validator, &configPath))
	cmd.AddCommand(newLintCmd(loaderFactory, validator, &configPath))
	cmd.AddCommand(newGenerateCmd(loaderFactory, validator, &configPath))
	return cmd
}

type versionDTO struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
	GoVersion string `json:"goVersion"`
	GOOS      string `json:"goos"`
	GOARCH    string `json:"goarch"`
}

func emitVersionJSON(w io.Writer) error {
	data, err := json.Marshal(versionDTO{
		Version:   buildinfo.Version,
		Commit:    buildinfo.Commit,
		Date:      buildinfo.Date,
		GoVersion: runtime.Version(),
		GOOS:      runtime.GOOS,
		GOARCH:    runtime.GOARCH,
	})
	if err != nil {
		return err
	}
	_, err = w.Write(append(data, '\n'))
	return err
}

func Execute(args []string, out io.Writer, errOut io.Writer, loader pipeline.AppLoader, validator pipeline.AppValidator) int {
	return ExecuteContextWithFactory(context.Background(), args, out, errOut, func(string) pipeline.AppLoader {
		return loader
	}, validator)
}

func ExecuteWithFactory(args []string, out io.Writer, errOut io.Writer, loaderFactory LoaderFactory, validator pipeline.AppValidator) int {
	return ExecuteContextWithFactory(context.Background(), args, out, errOut, loaderFactory, validator)
}

func ExecuteContext(ctx context.Context, args []string, out io.Writer, errOut io.Writer, loader pipeline.AppLoader, validator pipeline.AppValidator) int {
	return ExecuteContextWithFactory(ctx, args, out, errOut, func(string) pipeline.AppLoader {
		return loader
	}, validator)
}

func ExecuteContextWithFactory(ctx context.Context, args []string, out io.Writer, errOut io.Writer, loaderFactory LoaderFactory, validator pipeline.AppValidator) int {
	cmd := NewRootCmdWithFactory(loaderFactory, validator)
	cmd.SetOut(out)
	cmd.SetErr(errOut)
	cmd.SetArgs(args)

	err := cmd.ExecuteContext(ctx)
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
			return runWithConfig(cmd.Context(), loaderFactory, validator, configPath, validateAction{out: cmd.OutOrStdout()})
		},
	}
}

func newListCmd(loaderFactory LoaderFactory, validator pipeline.AppValidator, configPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List entities",
	}
	cmd.AddCommand(newListReportsCmd(loaderFactory, validator, configPath))
	cmd.AddCommand(newListNamesCmd(loaderFactory, validator, configPath))
	return cmd
}

func newListReportsCmd(loaderFactory LoaderFactory, validator pipeline.AppValidator, configPath *string) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "reports",
		Short: "List reports",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateReportsFormat(format); err != nil {
				return err
			}
			return runWithConfig(cmd.Context(), loaderFactory, validator, configPath, listReportsAction{
				out:    cmd.OutOrStdout(),
				format: format,
			})
		},
	}
	cmd.Flags().StringVar(&format, "format", reportsFormatJSON, "Output format: json|table")
	return cmd
}

func newListNamesCmd(loaderFactory LoaderFactory, validator pipeline.AppValidator, configPath *string) *cobra.Command {
	var prefix string
	var kind string
	var format string

	cmd := &cobra.Command{
		Use:   "names",
		Short: "List names",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateNamesKind(kind); err != nil {
				return err
			}
			if err := validateNamesFormat(format); err != nil {
				return err
			}
			return runWithConfig(cmd.Context(), loaderFactory, validator, configPath, namesAction{
				out:    cmd.OutOrStdout(),
				prefix: prefix,
				kind:   kind,
				format: format,
			})
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "Required name prefix filter")
	cmd.Flags().StringVar(&kind, "kind", namesKindAll, "Filter kind: all|notes|relationships")
	cmd.Flags().StringVar(&format, "format", namesFormatTable, "Output format: table|json")
	_ = cmd.MarkFlagRequired("prefix")
	return cmd
}

func newGenerateCmd(loaderFactory LoaderFactory, validator pipeline.AppValidator, configPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate artifacts",
	}
	cmd.AddCommand(newGenerateJSONCmd(loaderFactory, validator, configPath))
	cmd.AddCommand(newGenerateMarkdownCmd(loaderFactory, validator, configPath))
	return cmd
}

func newLintCmd(loaderFactory LoaderFactory, validator pipeline.AppValidator, configPath *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Lint entities",
	}
	cmd.AddCommand(newLintNamesCmd(loaderFactory, validator, configPath))
	cmd.AddCommand(newLintOrphansCmd(loaderFactory, validator, configPath))
	return cmd
}

func newLintNamesCmd(loaderFactory LoaderFactory, validator pipeline.AppValidator, configPath *string) *cobra.Command {
	var prefix string
	var style string
	var pattern string
	var severity string

	cmd := &cobra.Command{
		Use:   "names",
		Short: "Lint names",
		RunE: func(cmd *cobra.Command, args []string) error {
			policy, err := resolveLintNamesPolicy(style, pattern, severity)
			if err != nil {
				return err
			}
			return runWithConfig(cmd.Context(), loaderFactory, validator, configPath, lintNamesAction{
				out:    cmd.OutOrStdout(),
				prefix: prefix,
				policy: policy,
			})
		},
	}
	cmd.Flags().StringVar(&prefix, "prefix", "", "Optional prefix filter")
	cmd.Flags().StringVar(&style, "style", lintStyleDot, "Style matcher: dot|snake|regex")
	cmd.Flags().StringVar(&pattern, "pattern", "", "Regex pattern when style=regex")
	cmd.Flags().StringVar(&severity, "severity", "warning", "Diagnostic severity: warning|error")
	return cmd
}

func newLintOrphansCmd(loaderFactory LoaderFactory, validator pipeline.AppValidator, configPath *string) *cobra.Command {
	var subjectLabel string
	var edgeLabel string
	var counterpartLabel string
	var direction string
	var severity string

	cmd := &cobra.Command{
		Use:   "orphans",
		Short: "Lint orphans with a label-driven query",
		RunE: func(cmd *cobra.Command, args []string) error {
			query, diagSeverity, err := resolveLintOrphansQuery(subjectLabel, edgeLabel, counterpartLabel, direction, severity)
			if err != nil {
				return err
			}
			return runWithConfig(cmd.Context(), loaderFactory, validator, configPath, lintOrphansAction{
				out:      cmd.OutOrStdout(),
				query:    query,
				severity: diagSeverity,
			})
		},
	}
	cmd.Flags().StringVar(&subjectLabel, "subject-label", "", "Required subject note label")
	cmd.Flags().StringVar(&edgeLabel, "edge-label", "", "Optional relationship label filter")
	cmd.Flags().StringVar(&counterpartLabel, "counterpart-label", "", "Optional counterpart note label filter")
	cmd.Flags().StringVar(&direction, "direction", string(defaultOrphanDirection), "Relationship direction: in|out|either")
	cmd.Flags().StringVar(&severity, "severity", "warning", "Diagnostic severity: warning|error")
	_ = cmd.MarkFlagRequired("subject-label")
	return cmd
}

func newGenerateJSONCmd(loaderFactory LoaderFactory, validator pipeline.AppValidator, configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "json",
		Short: "Generate JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWithConfig(cmd.Context(), loaderFactory, validator, configPath, generateJSONAction{out: cmd.OutOrStdout()})
		},
	}
}

func newGenerateMarkdownCmd(loaderFactory LoaderFactory, validator pipeline.AppValidator, configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "markdown",
		Short: "Generate markdown reports",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWithConfig(cmd.Context(), loaderFactory, validator, configPath, generateMarkdownAction{out: cmd.OutOrStdout()})
		},
	}
}

func runWithConfig(ctx context.Context, loaderFactory LoaderFactory, validator pipeline.AppValidator, configPath *string, action pipeline.CommandAction) error {
	if loaderFactory == nil {
		return errors.New("loader factory is required")
	}
	if configPath == nil {
		return errors.New("config path is required")
	}
	loader := loaderFactory(*configPath)
	return pipeline.Run(ctx, loader, validator, action)
}

type validateAction struct {
	out io.Writer
}

func (a validateAction) Execute(_ context.Context, _ domain.ValidatedApp, report domain.ValidationReport) error {
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

type generateJSONAction struct {
	out io.Writer
}

func (a generateJSONAction) Execute(_ context.Context, validated domain.ValidatedApp, report domain.ValidationReport) error {
	_ = report
	return clioutput.EmitGraphJSON(a.out, validated)
}

func (generateJSONAction) AllowOnValidationErrors() bool {
	return false
}
