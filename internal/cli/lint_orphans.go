package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/orphans"
	"github.com/flarebyte/baldrick-flying-buttress/internal/outcome"
	clioutput "github.com/flarebyte/baldrick-flying-buttress/internal/output"
)

const defaultOrphanDirection = orphans.DirectionEither

type lintOrphansAction struct {
	out      io.Writer
	query    orphans.Query
	severity domain.Severity
}

func (a lintOrphansAction) Execute(validated domain.ValidatedApp, report domain.ValidationReport) error {
	_ = report
	diagnostics := a.diagnostics(validated)
	lintReport := domain.ValidationReport{Diagnostics: diagnostics}
	if err := clioutput.EmitDiagnostics(a.out, lintReport); err != nil {
		return err
	}
	if lintReport.HasErrors() {
		return outcome.ValidationBlockedError()
	}
	return nil
}

func (lintOrphansAction) AllowOnValidationErrors() bool {
	return false
}

func (a lintOrphansAction) diagnostics(validated domain.ValidatedApp) []domain.Diagnostic {
	orphanNotes := orphans.Find(validated, a.query)
	diagnostics := make([]domain.Diagnostic, 0, len(orphanNotes))
	for _, note := range orphanNotes {
		diagnostics = append(diagnostics, domain.Diagnostic{
			Code:             "ORPHAN_QUERY_MISSING_LINK",
			Severity:         a.severity,
			Source:           "orphans.query.find",
			Message:          "subject note has no relationships matching orphan query",
			Location:         fmt.Sprintf("notes[name=%q]", note.ID),
			Path:             fmt.Sprintf("notes[name=%q]", note.ID),
			NoteName:         note.ID,
			SubjectLabel:     a.query.SubjectLabel,
			EdgeLabel:        a.query.EdgeLabel,
			CounterpartLabel: a.query.CounterpartLabel,
		})
	}
	return diagnostics
}

func resolveLintOrphansQuery(subjectLabel, edgeLabel, counterpartLabel, direction, severity string) (orphans.Query, domain.Severity, error) {
	query := orphans.Query{
		SubjectLabel:     strings.TrimSpace(subjectLabel),
		EdgeLabel:        strings.TrimSpace(edgeLabel),
		CounterpartLabel: strings.TrimSpace(counterpartLabel),
		Direction:        orphans.Direction(strings.TrimSpace(direction)),
	}
	if query.Direction == "" {
		query.Direction = defaultOrphanDirection
	}
	if err := query.Validate(); err != nil {
		return orphans.Query{}, "", err
	}

	var diagSeverity domain.Severity
	switch strings.TrimSpace(severity) {
	case "", string(domain.SeverityWarning):
		diagSeverity = domain.SeverityWarning
	case string(domain.SeverityError):
		diagSeverity = domain.SeverityError
	default:
		return orphans.Query{}, "", fmt.Errorf("invalid severity: %s", severity)
	}

	return query, diagSeverity, nil
}
