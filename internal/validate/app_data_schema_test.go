package validate

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestAppDataValidatorSchemaDiagnostics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		raw       domain.RawApp
		wantCode  string
		wantPath  string
		wantCount int
	}{
		{
			name:     "missing report title",
			raw:      domain.RawApp{Source: "app", Reports: []domain.RawReport{{Title: "", Filepath: "reports/r1.md", Sections: []domain.RawReportSection{{Title: "S"}}}}},
			wantCode: "FBV101", wantPath: "reports[0].title", wantCount: 3,
		},
		{
			name:     "missing report filepath",
			raw:      domain.RawApp{Source: "app", Reports: []domain.RawReport{{Title: "R1", Filepath: "", Sections: []domain.RawReportSection{{Title: "S"}}}}},
			wantCode: "FBV102", wantPath: "reports[0].filepath", wantCount: 3,
		},
		{
			name:     "missing note name",
			raw:      domain.RawApp{Source: "app", Notes: []domain.RawNote{{Name: "", Title: "N"}}},
			wantCode: "FBV201", wantPath: "notes[0].name", wantCount: 3,
		},
		{
			name:     "missing note title",
			raw:      domain.RawApp{Source: "app", Notes: []domain.RawNote{{Name: "n1", Title: ""}}},
			wantCode: "FBV202", wantPath: "notes[0].title", wantCount: 3,
		},
		{
			name:     "missing relationship from",
			raw:      domain.RawApp{Source: "app", Relationships: []domain.RawRelationship{{FromID: "", ToID: "n2", Label: "depends_on"}}},
			wantCode: "FBV301", wantPath: "relationships[0].from", wantCount: 3,
		},
		{
			name:     "missing relationship to",
			raw:      domain.RawApp{Source: "app", Relationships: []domain.RawRelationship{{FromID: "n1", ToID: "", Label: "depends_on"}}},
			wantCode: "FBV302", wantPath: "relationships[0].to", wantCount: 3,
		},
		{
			name:     "missing report sections shape",
			raw:      domain.RawApp{Source: "app", Reports: []domain.RawReport{{Title: "R1", Filepath: "reports/r1.md", Sections: nil}}},
			wantCode: "FBV103", wantPath: "reports[0].sections", wantCount: 3,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, report, err := AppDataValidator{}.Validate(tc.raw)
			if err != nil {
				t.Fatalf("validate failed: %v", err)
			}
			if len(report.Diagnostics) != tc.wantCount {
				t.Fatalf("diagnostic count mismatch: got %d want %d", len(report.Diagnostics), tc.wantCount)
			}
			assertHasDiagnostics(t, report.Diagnostics, []diagExpectation{{code: tc.wantCode, path: tc.wantPath}}, "schema")
		})
	}
}
func TestAppDataValidatorCanonicalDiagnosticLocationAndMetadata(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{
		Source:  "app",
		Reports: []domain.RawReport{{Title: "", Filepath: "", Sections: []domain.RawReportSection{{Title: ""}}}},
	}

	_, report, err := AppDataValidator{}.Validate(raw)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}

	want := []domain.Diagnostic{
		{Code: "FBV101", Severity: domain.SeverityError, Source: schemaValidationSource, Message: "missing required field: report title", Location: "reports[0].title", Path: "reports[0].title"},
		{Code: "FBV102", Severity: domain.SeverityError, Source: schemaValidationSource, Message: "missing required field: report filepath", Location: "reports[0].filepath", Path: "reports[0].filepath"},
		{Code: "FBV104", Severity: domain.SeverityError, Source: schemaValidationSource, Message: "missing required field: section title", Location: "reports[0].sections[0].title", Path: "reports[0].sections[0].title"},
		{Code: "FBV200", Severity: domain.SeverityError, Source: schemaValidationSource, Message: "missing required collection: notes", Location: "notes", Path: "notes"},
		{Code: "FBV300", Severity: domain.SeverityError, Source: schemaValidationSource, Message: "missing required collection: relationships", Location: "relationships", Path: "relationships"},
	}
	if len(report.Diagnostics) != len(want) {
		t.Fatalf("diagnostic count mismatch: got %d want %d", len(report.Diagnostics), len(want))
	}
	for i := range want {
		got := report.Diagnostics[i]
		if got.Code != want[i].Code || got.Severity != want[i].Severity || got.Source != want[i].Source || got.Message != want[i].Message || got.Location != want[i].Location || got.Path != want[i].Path {
			t.Fatalf("diagnostic %d mismatch: got %#v want %#v", i, got, want[i])
		}
	}
}
