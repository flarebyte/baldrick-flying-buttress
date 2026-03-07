package output

import (
	"encoding/json"
	"io"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
)

type DiagnosticsDTO struct {
	Diagnostics []DiagnosticDTO `json:"diagnostics"`
}

type DiagnosticDTO struct {
	Code     string          `json:"code"`
	Severity domain.Severity `json:"severity"`
	Message  string          `json:"message"`
	Path     string          `json:"path"`
}

type ReportListDTO struct {
	Reports []ReportDTO `json:"reports"`
}

type ReportDTO struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type GraphDTO struct {
	Notes         []NoteDTO         `json:"notes"`
	Relationships []RelationshipDTO `json:"relationships"`
}

type NoteDTO struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type RelationshipDTO struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Label string `json:"label"`
}

func EmitDiagnostics(w io.Writer, report domain.ValidationReport) error {
	canonical := report.Canonical()
	ordered := ordering.Diagnostics(canonical.Diagnostics)
	dto := DiagnosticsDTO{Diagnostics: make([]DiagnosticDTO, 0, len(ordered))}
	for _, d := range ordered {
		dto.Diagnostics = append(dto.Diagnostics, DiagnosticDTO{
			Code:     d.Code,
			Severity: d.Severity,
			Message:  d.Message,
			Path:     d.Path,
		})
	}
	return emitJSON(w, dto)
}

func EmitReportList(w io.Writer, app domain.ValidatedApp) error {
	ordered := ordering.Reports(app.Reports)
	dto := ReportListDTO{Reports: make([]ReportDTO, 0, len(ordered))}
	for _, r := range ordered {
		dto.Reports = append(dto.Reports, ReportDTO{ID: r.ID, Title: r.Title})
	}
	return emitJSON(w, dto)
}

func EmitGraphJSON(w io.Writer, app domain.ValidatedApp) error {
	orderedNotes := ordering.Notes(app.Notes)
	orderedRelationships := ordering.Relationships(app.Relationships)
	dto := GraphDTO{
		Notes:         make([]NoteDTO, 0, len(orderedNotes)),
		Relationships: make([]RelationshipDTO, 0, len(orderedRelationships)),
	}
	for _, n := range orderedNotes {
		dto.Notes = append(dto.Notes, NoteDTO{ID: n.ID, Label: n.Label})
	}
	for _, r := range orderedRelationships {
		dto.Relationships = append(dto.Relationships, RelationshipDTO{From: r.FromID, To: r.ToID, Label: r.Label})
	}
	return emitJSON(w, dto)
}

func emitJSON(w io.Writer, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = w.Write(append(data, '\n'))
	return err
}
