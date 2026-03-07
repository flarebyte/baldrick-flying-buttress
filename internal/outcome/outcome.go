package outcome

import "errors"

type Kind string

const (
	KindSuccess           Kind = "success"
	KindValidationBlocked Kind = "validation_blocked"
	KindRuntimeFailure    Kind = "runtime_failure"
)

const (
	ExitCodeSuccess           = 0
	ExitCodeValidationBlocked = 1
	ExitCodeRuntimeFailure    = 2
)

var errValidationBlocked = errors.New("validation blocked")

func ValidationBlockedError() error {
	return errValidationBlocked
}

func IsValidationBlocked(err error) bool {
	return errors.Is(err, errValidationBlocked)
}

type Execution struct {
	Kind Kind
	Err  error
}

func FromError(err error) Execution {
	if err == nil {
		return Execution{Kind: KindSuccess}
	}
	if IsValidationBlocked(err) {
		return Execution{Kind: KindValidationBlocked, Err: err}
	}
	return Execution{Kind: KindRuntimeFailure, Err: err}
}

func (e Execution) ExitCode() int {
	switch e.Kind {
	case KindSuccess:
		return ExitCodeSuccess
	case KindValidationBlocked:
		return ExitCodeValidationBlocked
	default:
		return ExitCodeRuntimeFailure
	}
}
