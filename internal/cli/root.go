package cli

import (
	"errors"
	"fmt"
	"io"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/pipeline"
	"github.com/spf13/cobra"
)

func NewRootCmd(loader pipeline.AppLoader, validator pipeline.AppValidator) *cobra.Command {
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

func Execute(args []string, out io.Writer, errOut io.Writer, loader pipeline.AppLoader, validator pipeline.AppValidator) int {
	cmd := NewRootCmd(loader, validator)
	cmd.SetOut(out)
	cmd.SetErr(errOut)
	cmd.SetArgs(args)

	err := cmd.Execute()
	if err == nil {
		return 0
	}
	if errors.Is(err, pipeline.ErrValidationFailed) {
		return 1
	}
	_, _ = fmt.Fprintln(errOut, err.Error())
	return 1
}

func newValidateCmd(loader pipeline.AppLoader, validator pipeline.AppValidator) *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pipeline.Run(loader, validator, validateAction{out: cmd.OutOrStdout()})
		},
	}
}

func newListCmd(loader pipeline.AppLoader, validator pipeline.AppValidator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List entities",
	}
	cmd.AddCommand(newListReportsCmd(loader, validator))
	return cmd
}

func newListReportsCmd(loader pipeline.AppLoader, validator pipeline.AppValidator) *cobra.Command {
	return &cobra.Command{
		Use:   "reports",
		Short: "List reports",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pipeline.Run(loader, validator, listReportsAction{out: cmd.OutOrStdout()})
		},
	}
}

func newGenerateCmd(loader pipeline.AppLoader, validator pipeline.AppValidator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate artifacts",
	}
	cmd.AddCommand(newGenerateJSONCmd(loader, validator))
	return cmd
}

func newGenerateJSONCmd(loader pipeline.AppLoader, validator pipeline.AppValidator) *cobra.Command {
	return &cobra.Command{
		Use:   "json",
		Short: "Generate JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pipeline.Run(loader, validator, generateJSONAction{out: cmd.OutOrStdout()})
		},
	}
}

type validateAction struct {
	out io.Writer
}

func (a validateAction) Execute(validated domain.ValidatedApp, report domain.ValidationReport) error {
	_ = validated
	if err := emitDiagnostics(a.out, report); err != nil {
		return err
	}
	if report.HasErrors() {
		return pipeline.ErrValidationFailed
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
	return emitReportList(a.out, validated)
}

func (listReportsAction) AllowOnValidationErrors() bool {
	return false
}

type generateJSONAction struct {
	out io.Writer
}

func (a generateJSONAction) Execute(validated domain.ValidatedApp, report domain.ValidationReport) error {
	_ = report
	return emitGraphJSON(a.out, validated)
}

func (generateJSONAction) AllowOnValidationErrors() bool {
	return false
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
