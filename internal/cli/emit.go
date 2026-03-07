package cli

import (
	"encoding/json"
	"io"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
)

func emitDiagnostics(out io.Writer, report domain.ValidationReport) error {
	canonical := report.Canonical()
	canonical.Diagnostics = ordering.Diagnostics(canonical.Diagnostics)
	return emitJSON(out, canonical)
}

func emitReportList(out io.Writer, validated domain.ValidatedApp) error {
	orderedReports := ordering.Reports(validated.Reports)
	payload := listReportsOutput{Reports: make([]listReport, 0, len(orderedReports))}
	for _, r := range orderedReports {
		payload.Reports = append(payload.Reports, listReport{ID: r.ID, Title: r.Title})
	}
	return emitJSON(out, payload)
}

func emitGraphJSON(out io.Writer, validated domain.ValidatedApp) error {
	orderedNotes := ordering.Notes(validated.Notes)
	orderedRelationships := ordering.Relationships(validated.Relationships)
	payload := generateJSONOutput{
		Notes:         make([]generateNote, 0, len(orderedNotes)),
		Relationships: make([]generateRelationship, 0, len(orderedRelationships)),
	}
	for _, n := range orderedNotes {
		payload.Notes = append(payload.Notes, generateNote{ID: n.ID, Label: n.Label})
	}
	for _, r := range orderedRelationships {
		payload.Relationships = append(payload.Relationships, generateRelationship{From: r.FromID, To: r.ToID, Label: r.Label})
	}
	return emitJSON(out, payload)
}

func emitJSON(out io.Writer, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = out.Write(append(data, '\n'))
	return err
}
