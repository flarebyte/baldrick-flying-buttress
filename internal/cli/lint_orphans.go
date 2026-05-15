// purpose: Implements CLI behavior for lint_orphans.go so commands expose deterministic, machine-friendly output surfaces.
// responsibilities: parse command inputs; call pipeline/domain services; render structured outputs or diagnostics; enforce deterministic CLI behavior
// architecture_notes: CLI logic is split into focused files per command area to keep Cobra wiring thin and to isolate rendering from validation and domain logic.
package cli

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/orphans"
)

const defaultOrphanDirection = orphans.DirectionEither

type lintOrphansAction struct {
	out      io.Writer
	query    orphans.Query
	severity domain.Severity
}

func (a lintOrphansAction) Execute(_ context.Context, validated domain.ValidatedApp, report domain.ValidationReport) error {
	_ = report
	return emitDiagnosticsOutcome(a.out, a.diagnostics(validated))
}

func (lintOrphansAction) AllowOnValidationErrors() bool {
	return false
}

func (a lintOrphansAction) diagnostics(validated domain.ValidatedApp) []domain.Diagnostic {
	orphanNotes := orphans.Find(validated, a.query)
	diagnostics := make([]domain.Diagnostic, 0, len(orphanNotes))
	for _, note := range orphanNotes {
		diagnostics = append(diagnostics, withDiagnosticContextMessage(domain.Diagnostic{
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
		}))
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
	diagSeverity, err := resolveSeverity(severity)
	if err != nil {
		return orphans.Query{}, "", err
	}
	return query, diagSeverity, nil
}
