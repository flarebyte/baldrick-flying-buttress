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
	Code             string          `json:"code"`
	Severity         domain.Severity `json:"severity"`
	Message          string          `json:"message"`
	Path             string          `json:"path"`
	NormalizedPath   string          `json:"normalizedPath,omitempty"`
	ConfigPath       string          `json:"configPath,omitempty"`
	ConfigPathAbs    string          `json:"configPathAbsolute,omitempty"`
	ReportTitle      string          `json:"reportTitle,omitempty"`
	ReportID         string          `json:"reportId,omitempty"`
	SectionTitle     string          `json:"sectionTitle,omitempty"`
	NoteName         string          `json:"noteName,omitempty"`
	NoteTitle        string          `json:"noteTitle,omitempty"`
	ArgumentName     string          `json:"argumentName,omitempty"`
	LabelValue       string          `json:"labelValue,omitempty"`
	SubjectLabel     string          `json:"subjectLabel,omitempty"`
	EdgeLabel        string          `json:"edgeLabel,omitempty"`
	CounterpartLabel string          `json:"counterpartLabel,omitempty"`
	RelationshipFrom string          `json:"relationshipFrom,omitempty"`
	RelationshipTo   string          `json:"relationshipTo,omitempty"`
	RelatedNodes     []string        `json:"relatedNodes,omitempty"`
	SuggestedFixes   []string        `json:"suggestedFixes,omitempty"`
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
			Code:             d.Code,
			Severity:         d.Severity,
			Message:          d.Message,
			Path:             d.Path,
			NormalizedPath:   d.NormalizedPath,
			ConfigPath:       d.ConfigPath,
			ConfigPathAbs:    d.ConfigPathAbs,
			ReportTitle:      d.ReportTitle,
			ReportID:         d.ReportID,
			SectionTitle:     d.SectionTitle,
			NoteName:         d.NoteName,
			NoteTitle:        d.NoteTitle,
			ArgumentName:     d.ArgumentName,
			LabelValue:       d.LabelValue,
			SubjectLabel:     d.SubjectLabel,
			EdgeLabel:        d.EdgeLabel,
			CounterpartLabel: d.CounterpartLabel,
			RelationshipFrom: d.RelationshipFrom,
			RelationshipTo:   d.RelationshipTo,
			RelatedNodes:     d.RelatedNodes,
			SuggestedFixes:   d.SuggestedFixes,
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
