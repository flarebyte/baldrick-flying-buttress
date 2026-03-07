package cli

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
	clioutput "github.com/flarebyte/baldrick-flying-buttress/internal/output"
)

const (
	lintStyleDot   = "dot"
	lintStyleSnake = "snake"
	lintStyleRegex = "regex"
)

var (
	dotNameRE   = regexp.MustCompile(`^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$`)
	snakeNameRE = regexp.MustCompile(`^[a-z][a-z0-9]*(_[a-z0-9]+)*$`)
)

type lintNamesPolicy struct {
	style    string
	severity domain.Severity
	matcher  func(string) bool
}

type lintNamesAction struct {
	out    io.Writer
	prefix string
	policy lintNamesPolicy
}

func (a lintNamesAction) Execute(validated domain.ValidatedApp, report domain.ValidationReport) error {
	_ = report
	diagnostics := lintNames(validated, a.prefix, a.policy)
	lintReport := domain.ValidationReport{Diagnostics: diagnostics}
	if err := clioutput.EmitDiagnostics(a.out, lintReport); err != nil {
		return err
	}
	if lintReport.HasErrors() {
		return outcome.ValidationBlockedError()
	}
	return nil
}

func (lintNamesAction) AllowOnValidationErrors() bool {
	return false
}

func resolveLintNamesPolicy(style string, pattern string, severity string) (lintNamesPolicy, error) {
	pol := lintNamesPolicy{style: style, severity: domain.SeverityWarning}
	switch severity {
	case "warning", "":
		pol.severity = domain.SeverityWarning
	case "error":
		pol.severity = domain.SeverityError
	default:
		return lintNamesPolicy{}, fmt.Errorf("invalid severity: %s", severity)
	}

	switch style {
	case lintStyleDot, "":
		pol.style = lintStyleDot
		pol.matcher = dotNameRE.MatchString
	case lintStyleSnake:
		pol.matcher = snakeNameRE.MatchString
	case lintStyleRegex:
		if strings.TrimSpace(pattern) == "" {
			return lintNamesPolicy{}, fmt.Errorf("pattern is required when style=regex")
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return lintNamesPolicy{}, err
		}
		pol.matcher = re.MatchString
	default:
		return lintNamesPolicy{}, fmt.Errorf("invalid style: %s", style)
	}

	return pol, nil
}

func lintNames(app domain.ValidatedApp, prefix string, policy lintNamesPolicy) []domain.Diagnostic {
	orderedNotes := ordering.Notes(app.Notes)
	orderedRelationships := ordering.Relationships(app.Relationships)
	diagnostics := make([]domain.Diagnostic, 0)

	for _, note := range orderedNotes {
		if !matchesNotePrefix(note.ID, prefix) {
			continue
		}
		if policy.matcher(note.ID) {
			continue
		}
		diagnostics = append(diagnostics, domain.Diagnostic{
			Code:     "NAME_STYLE_VIOLATION",
			Severity: policy.severity,
			Source:   "lint.names.notes",
			Message:  "note name violates style policy",
			Location: fmt.Sprintf("notes[name=%q]", note.ID),
			Path:     fmt.Sprintf("notes[name=%q]", note.ID),
			NoteName: note.ID,
		})
	}

	for i, rel := range orderedRelationships {
		if !matchesRelationshipPrefix(rel, prefix) {
			continue
		}
		if !policy.matcher(rel.FromID) {
			diagnostics = append(diagnostics, domain.Diagnostic{
				Code:             "NAME_STYLE_VIOLATION",
				Severity:         policy.severity,
				Source:           "lint.names.relationships",
				Message:          "relationship from endpoint violates style policy",
				Location:         fmt.Sprintf("relationships[%d].from", i),
				Path:             fmt.Sprintf("relationships[%d].from", i),
				RelationshipFrom: rel.FromID,
				RelationshipTo:   rel.ToID,
			})
		}
		if !policy.matcher(rel.ToID) {
			diagnostics = append(diagnostics, domain.Diagnostic{
				Code:             "NAME_STYLE_VIOLATION",
				Severity:         policy.severity,
				Source:           "lint.names.relationships",
				Message:          "relationship to endpoint violates style policy",
				Location:         fmt.Sprintf("relationships[%d].to", i),
				Path:             fmt.Sprintf("relationships[%d].to", i),
				RelationshipFrom: rel.FromID,
				RelationshipTo:   rel.ToID,
			})
		}
	}

	return ordering.Diagnostics(diagnostics)
}
