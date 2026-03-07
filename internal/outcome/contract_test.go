package outcome

import (
	"errors"
	"testing"
)

func TestContractExitCodeMapping(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		err  error
		kind Kind
		code int
	}{
		{name: "success", err: nil, kind: KindSuccess, code: ExitCodeSuccess},
		{name: "validation_blocked", err: ValidationBlockedError(), kind: KindValidationBlocked, code: ExitCodeValidationBlocked},
		{name: "runtime_failure", err: errors.New("boom"), kind: KindRuntimeFailure, code: ExitCodeRuntimeFailure},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			exec := FromError(tc.err)
			if exec.Kind != tc.kind {
				t.Fatalf("expected kind %q, got %q", tc.kind, exec.Kind)
			}
			if exec.ExitCode() != tc.code {
				t.Fatalf("expected exit code %d, got %d", tc.code, exec.ExitCode())
			}
		})
	}
}
