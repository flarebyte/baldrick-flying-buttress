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
