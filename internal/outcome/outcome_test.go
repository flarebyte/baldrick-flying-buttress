package outcome

import (
	"errors"
	"testing"
)

func TestFromErrorSuccess(t *testing.T) {
	t.Parallel()

	exec := FromError(nil)
	if exec.Kind != KindSuccess {
		t.Fatalf("expected %q, got %q", KindSuccess, exec.Kind)
	}
	if exec.ExitCode() != ExitCodeSuccess {
		t.Fatalf("expected exit code %d, got %d", ExitCodeSuccess, exec.ExitCode())
	}
}

func TestFromErrorValidationBlocked(t *testing.T) {
	t.Parallel()

	err := ValidationBlockedError()
	exec := FromError(err)
	if exec.Kind != KindValidationBlocked {
		t.Fatalf("expected %q, got %q", KindValidationBlocked, exec.Kind)
	}
	if exec.ExitCode() != ExitCodeValidationBlocked {
		t.Fatalf("expected exit code %d, got %d", ExitCodeValidationBlocked, exec.ExitCode())
	}
}

func TestFromErrorRuntimeFailure(t *testing.T) {
	t.Parallel()

	exec := FromError(errors.New("boom"))
	if exec.Kind != KindRuntimeFailure {
		t.Fatalf("expected %q, got %q", KindRuntimeFailure, exec.Kind)
	}
	if exec.ExitCode() != ExitCodeRuntimeFailure {
		t.Fatalf("expected exit code %d, got %d", ExitCodeRuntimeFailure, exec.ExitCode())
	}
}
