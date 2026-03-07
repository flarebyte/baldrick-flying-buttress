package diagnostics

type Severity string

const (
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

type Diagnostic struct {
	Code     string   `json:"code"`
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
	Path     string   `json:"path"`
}

type Report struct {
	Diagnostics []Diagnostic `json:"diagnostics"`
}

func (r Report) HasErrors() bool {
	for _, d := range r.Diagnostics {
		if d.Severity == SeverityError {
			return true
		}
	}
	return false
}
