package cli

import (
	"strings"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestEncodeRawAppToCueDeterministicAndNormalized(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{
		Source:  "app",
		Name:    "demo",
		Modules: []string{"edge", "core"},
		Reports: []domain.RawReport{{
			Title:    "Zeta",
			Filepath: "out/zeta.md",
			Sections: []domain.RawReportSection{{
				Title:     "B",
				Notes:     []string{"n.b", "n.a"},
				Arguments: []string{"show-labels=true", "x=1"},
			}},
		}, {
			Title:    "Alpha",
			Filepath: "out/alpha.md",
			Sections: []domain.RawReportSection{{
				Title:       "A",
				Description: "overview",
				Sections: []domain.RawReportSection{{
					Title: "Nested",
					Notes: []string{"n.c", "n.a"},
				}},
			}},
		}},
		Notes: []domain.RawNote{
			{Name: "n.b", Title: "Beta", Labels: []string{"v2", "future"}},
			{Name: "n.a", Title: "Alpha", Labels: []string{"v1", "implementation"}},
		},
		Relationships: []domain.RawRelationship{
			{FromID: "n.b", ToID: "n.a", Label: "depends_on", Labels: []string{"core", "depends_on"}},
			{FromID: "n.a", ToID: "n.b", Label: "depends_on", Labels: []string{"depends_on"}},
		},
	}

	got, err := encodeRawAppToCue(canonicalRawApp(raw))
	if err != nil {
		t.Fatalf("encode cue failed: %v", err)
	}

	text := string(got)
	if !strings.HasPrefix(text, "package flyb\n\n") {
		t.Fatalf("expected package header, got %q", text)
	}
	if !strings.Contains(text, "modules: [\"core\", \"edge\"]") {
		t.Fatalf("expected sorted modules, got %q", text)
	}
	if strings.Index(text, "filepath: \"out/alpha.md\"") > strings.Index(text, "filepath: \"out/zeta.md\"") {
		t.Fatalf("expected reports sorted by filepath, got %q", text)
	}
	if !strings.Contains(text, "arguments: [\"show-labels=true\", \"x=1\"]") {
		t.Fatalf("expected sorted arguments, got %q", text)
	}
	if !strings.Contains(text, "labels: [\"future\", \"v2\"]") {
		t.Fatalf("expected sorted note labels, got %q", text)
	}
}
