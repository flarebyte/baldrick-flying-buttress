package domain

type RawReport struct {
	Title    string             `json:"title"`
	Filepath string             `json:"filepath"`
	Sections []RawReportSection `json:"sections"`
}

type RawReportSection struct {
	Title     string   `json:"title"`
	Arguments []string `json:"arguments"`
}

type RawNote struct {
	Name      string   `json:"name"`
	Title     string   `json:"title"`
	Arguments []string `json:"arguments"`
}

type RawRelationship struct {
	FromID string `json:"from"`
	ToID   string `json:"to"`
	Label  string `json:"label"`
}

type RawArgumentDefinition struct {
	Name          string   `json:"name"`
	ValueType     string   `json:"valueType"`
	Scopes        []string `json:"scopes"`
	AllowedValues []string `json:"allowedValues"`
	DefaultValue  any      `json:"defaultValue"`
}

type RawArgumentRegistry struct {
	Version   string                  `json:"version"`
	Arguments []RawArgumentDefinition `json:"arguments"`
}

type RawApp struct {
	ConfigPath    string              `json:"-"`
	Source        string              `json:"source"`
	Name          string              `json:"name"`
	Modules       []string            `json:"modules"`
	Reports       []RawReport         `json:"reports"`
	Notes         []RawNote           `json:"notes"`
	Relationships []RawRelationship   `json:"relationships"`
	Registry      RawArgumentRegistry `json:"argumentRegistry"`
}

type Severity string

const (
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

type Diagnostic struct {
	Code         string   `json:"code"`
	Severity     Severity `json:"severity"`
	Source       string   `json:"source"`
	Message      string   `json:"message"`
	Location     string   `json:"location"`
	Path         string   `json:"path"`
	ReportTitle  string   `json:"reportTitle,omitempty"`
	SectionTitle string   `json:"sectionTitle,omitempty"`
	NoteName     string   `json:"noteName,omitempty"`
	ArgumentName string   `json:"argumentName,omitempty"`
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

type ArgumentScope string

const (
	ArgumentScopeH3Section ArgumentScope = "h3-section"
	ArgumentScopeNote      ArgumentScope = "note"
	ArgumentScopeRenderer  ArgumentScope = "renderer"
)

type ArgumentValueType string

const (
	ArgumentValueTypeString  ArgumentValueType = "string"
	ArgumentValueTypeStrings ArgumentValueType = "string[]"
	ArgumentValueTypeBoolean ArgumentValueType = "boolean"
	ArgumentValueTypeInt     ArgumentValueType = "int"
	ArgumentValueTypeFloat   ArgumentValueType = "float"
	ArgumentValueTypeEnum    ArgumentValueType = "enum"
)

type ArgumentDefinition struct {
	Name          string
	ValueType     ArgumentValueType
	Scopes        []ArgumentScope
	AllowedValues []string
	DefaultValue  *string
}

type ArgumentRegistry struct {
	Version   string
	Arguments []ArgumentDefinition
}

type ValidatedApp struct {
	Name          string
	Modules       []string
	Reports       []Report
	Notes         []Note
	Relationships []Relationship
	Registry      ArgumentRegistry
}
