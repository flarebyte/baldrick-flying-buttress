package cli

import (
	"encoding/json"
	"io"

	"github.com/olivier/baldrick-flying-buttress/internal/app"
	"github.com/olivier/baldrick-flying-buttress/internal/diagnostics"
)

func emitDiagnostics(out io.Writer, report diagnostics.Report) error {
	return emitJSON(out, report)
}

func emitReportList(out io.Writer, validated app.ValidatedApp) error {
	payload := listReportsOutput{Reports: make([]listReport, 0, len(validated.Reports))}
	for _, r := range validated.Reports {
		payload.Reports = append(payload.Reports, listReport{ID: r.ID, Title: r.Title})
	}
	return emitJSON(out, payload)
}

func emitGraphJSON(out io.Writer, validated app.ValidatedApp) error {
	payload := generateJSONOutput{
		Notes:         make([]generateNote, 0, len(validated.Notes)),
		Relationships: make([]generateRelationship, 0, len(validated.Relationships)),
	}
	for _, n := range validated.Notes {
		payload.Notes = append(payload.Notes, generateNote{ID: n.ID, Label: n.Label})
	}
	for _, r := range validated.Relationships {
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
