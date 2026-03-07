package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/olivier/baldrick-flying-buttress/internal/app"
	"github.com/olivier/baldrick-flying-buttress/internal/diagnostics"
	"github.com/olivier/baldrick-flying-buttress/internal/pipeline"
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
				payload, err := json.Marshal(report)
				if err != nil {
					return err
				}
				if _, err := cmd.OutOrStdout().Write(append(payload, '\n')); err != nil {
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
				out := listReportsOutput{Reports: make([]listReport, 0, len(validated.Reports))}
				for _, r := range validated.Reports {
					out.Reports = append(out.Reports, listReport{ID: r.ID, Title: r.Title})
				}
				payload, err := json.Marshal(out)
				if err != nil {
					return err
				}
				if _, err := cmd.OutOrStdout().Write(append(payload, '\n')); err != nil {
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
				out := generateJSONOutput{
					Notes:         make([]generateNote, 0, len(validated.Notes)),
					Relationships: make([]generateRelationship, 0, len(validated.Relationships)),
				}
				for _, n := range validated.Notes {
					out.Notes = append(out.Notes, generateNote{ID: n.ID, Label: n.Label})
				}
				for _, r := range validated.Relationships {
					out.Relationships = append(out.Relationships, generateRelationship{From: r.FromID, To: r.ToID, Label: r.Label})
				}
				payload, err := json.Marshal(out)
				if err != nil {
					return err
				}
				if _, err := cmd.OutOrStdout().Write(append(payload, '\n')); err != nil {
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
