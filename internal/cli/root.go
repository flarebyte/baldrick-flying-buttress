package cli

import (
	"errors"
	"fmt"
	"io"

	"github.com/flarebyte/baldrick-flying-buttress/internal/app"
	"github.com/flarebyte/baldrick-flying-buttress/internal/diagnostics"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
	"github.com/spf13/cobra"
)

var errValidationFailed = errors.New("validation failed")

func NewRootCmd(loader pipeline.Loader, validator pipeline.Validator) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "flyb",
		Short:         "Baldrick Flying Buttress CLI",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(newValidateCmd(loader, validator))
	cmd.AddCommand(newListCmd(loader, validator))
	cmd.AddCommand(newGenerateCmd(loader, validator))
	return cmd
}

func Execute(args []string, out io.Writer, errOut io.Writer, loader pipeline.Loader, validator pipeline.Validator) int {
	cmd := NewRootCmd(loader, validator)
	cmd.SetOut(out)
	cmd.SetErr(errOut)
	cmd.SetArgs(args)

	err := cmd.Execute()
	if err == nil {
		return 0
	}
	if errors.Is(err, errValidationFailed) {
		return 1
	}
	_, _ = fmt.Fprintln(errOut, err.Error())
	return 1
}

func newValidateCmd(loader pipeline.Loader, validator pipeline.Validator) *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pipeline.Run(loader, validator, func(validated app.ValidatedApp, report diagnostics.Report) error {
				_ = validated
				if err := emitDiagnostics(cmd.OutOrStdout(), report); err != nil {
					return err
				}
				if report.HasErrors() {
					return errValidationFailed
				}
				return nil
			})
		},
	}
}

func newListCmd(loader pipeline.Loader, validator pipeline.Validator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List entities",
	}
	cmd.AddCommand(newListReportsCmd(loader, validator))
	return cmd
}

func newListReportsCmd(loader pipeline.Loader, validator pipeline.Validator) *cobra.Command {
	return &cobra.Command{
		Use:   "reports",
		Short: "List reports",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pipeline.Run(loader, validator, func(validated app.ValidatedApp, report diagnostics.Report) error {
				if report.HasErrors() {
					return errValidationFailed
				}
				if err := emitReportList(cmd.OutOrStdout(), validated); err != nil {
					return err
				}
				return nil
			})
		},
	}
}

func newGenerateCmd(loader pipeline.Loader, validator pipeline.Validator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate artifacts",
	}
	cmd.AddCommand(newGenerateJSONCmd(loader, validator))
	return cmd
}

func newGenerateJSONCmd(loader pipeline.Loader, validator pipeline.Validator) *cobra.Command {
	return &cobra.Command{
		Use:   "json",
		Short: "Generate JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pipeline.Run(loader, validator, func(validated app.ValidatedApp, report diagnostics.Report) error {
				if report.HasErrors() {
					return errValidationFailed
				}
				if err := emitGraphJSON(cmd.OutOrStdout(), validated); err != nil {
					return err
				}
				return nil
			})
		},
	}
}

type listReportsOutput struct {
	Reports []listReport `json:"reports"`
}

type listReport struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type generateJSONOutput struct {
	Notes         []generateNote         `json:"notes"`
	Relationships []generateRelationship `json:"relationships"`
}

type generateNote struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type generateRelationship struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Label string `json:"label"`
}
