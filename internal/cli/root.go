package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/olivier/baldrick-flying-buttress/internal/validate"
	"github.com/spf13/cobra"
)

var errValidationFailed = errors.New("validation failed")

func NewRootCmd(runner validate.Runner) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "flyb",
		Short:         "Baldrick Flying Buttress CLI",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(newValidateCmd(runner))
	return cmd
}

func Execute(args []string, out io.Writer, errOut io.Writer, runner validate.Runner) int {
	cmd := NewRootCmd(runner)
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

func newValidateCmd(runner validate.Runner) *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			report, err := runner(context.Background())
			if err != nil {
				return err
			}

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
		},
	}
}
