package domain

type RawReport struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type RawNote struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type RawRelationship struct {
	FromID string `json:"from"`
	ToID   string `json:"to"`
	Label  string `json:"label"`
}

type RawApp struct {
	ConfigPath    string            `json:"-"`
	Source        string            `json:"source"`
	Name          string            `json:"name"`
	Modules       []string          `json:"modules"`
	Reports       []RawReport       `json:"reports"`
	Notes         []RawNote         `json:"notes"`
	Relationships []RawRelationship `json:"relationships"`
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
