package domain

type RawApp struct {
	ConfigPath string
	Source     string
}

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

type ValidationReport struct {
	Diagnostics []Diagnostic `json:"diagnostics"`
}

func (r ValidationReport) Canonical() ValidationReport {
	if r.Diagnostics == nil {
		r.Diagnostics = []Diagnostic{}
	}
	return r
}

func (r ValidationReport) HasErrors() bool {
	for _, d := range r.Diagnostics {
		if d.Severity == SeverityError {
			return true
		}
	}
	return false
}

type Note struct {
	ID    string
	Label string
}

type Relationship struct {
	FromID string
	ToID   string
	Label  string
}

type Report struct {
	ID    string
	Title string
}

type ValidatedApp struct {
	Name          string
	Modules       []string
	Reports       []Report
	Notes         []Note
	Relationships []Relationship
}
