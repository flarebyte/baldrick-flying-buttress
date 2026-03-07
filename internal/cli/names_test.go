package cli

import (
	"bytes"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestFilterNamesPrefixNotes(t *testing.T) {
	t.Parallel()

	app := domain.ValidatedApp{Notes: []domain.Note{{ID: "cli.root", Label: "cli.root"}, {ID: "app.db", Label: "app.db"}}}
	notes, relationships, err := filterNames(app, "cli.", namesKindNotes)
	if err != nil {
		t.Fatalf("filter failed: %v", err)
	}
	if len(notes) != 1 || notes[0].ID != "cli.root" {
		t.Fatalf("unexpected notes: %#v", notes)
	}
	if len(relationships) != 0 {
		t.Fatalf("expected no relationships, got %#v", relationships)
	}
}

func TestFilterNamesPrefixRelationships(t *testing.T) {
	t.Parallel()

	app := domain.ValidatedApp{Relationships: []domain.Relationship{{FromID: "cli.root", ToID: "app.db", Label: "depends_on"}, {FromID: "app.db", ToID: "cli.worker", Label: "depends_on"}, {FromID: "app.db", ToID: "infra.cache", Label: "depends_on"}}}
	notes, relationships, err := filterNames(app, "cli.", namesKindRelationships)
	if err != nil {
		t.Fatalf("filter failed: %v", err)
	}
	if len(notes) != 0 {
		t.Fatalf("expected no notes, got %#v", notes)
	}
	if len(relationships) != 2 {
		t.Fatalf("expected 2 relationships, got %#v", relationships)
	}
}

func TestFilterNamesKindAllNotesRelationships(t *testing.T) {
	t.Parallel()

	app := listNamesValidatedApp()

	notesAll, relAll, err := filterNames(app, "cli.", namesKindAll)
	if err != nil {
		t.Fatalf("filter all failed: %v", err)
	}
	if len(notesAll) != 2 || len(relAll) != 2 {
		t.Fatalf("unexpected all filter result notes=%d relationships=%d", len(notesAll), len(relAll))
	}

	notesOnly, relOnly, err := filterNames(app, "cli.", namesKindNotes)
	if err != nil {
		t.Fatalf("filter notes failed: %v", err)
	}
	if len(notesOnly) != 2 || len(relOnly) != 0 {
		t.Fatalf("unexpected notes filter result notes=%d relationships=%d", len(notesOnly), len(relOnly))
	}

	notesRel, relRel, err := filterNames(app, "cli.", namesKindRelationships)
	if err != nil {
		t.Fatalf("filter relationships failed: %v", err)
	}
	if len(notesRel) != 0 || len(relRel) != 2 {
		t.Fatalf("unexpected relationships filter result notes=%d relationships=%d", len(notesRel), len(relRel))
	}
}

func TestEmitNamesTableDeterministic(t *testing.T) {
	t.Parallel()

	notes, relationships, err := filterNames(listNamesValidatedApp(), "cli.", namesKindAll)
	if err != nil {
		t.Fatalf("filter failed: %v", err)
	}

	var first bytes.Buffer
	if err := emitNamesTable(&first, notes, relationships); err != nil {
		t.Fatalf("emit first failed: %v", err)
	}
	var second bytes.Buffer
	if err := emitNamesTable(&second, notes, relationships); err != nil {
		t.Fatalf("emit second failed: %v", err)
	}

	want := readGolden(t, "list_names_table_output.golden")
	if first.String() != want {
		t.Fatalf("table mismatch\nwant: %q\n got: %q", want, first.String())
	}
	if second.String() != first.String() {
		t.Fatalf("non-deterministic table output\nfirst: %q\nsecond: %q", first.String(), second.String())
	}
}

func TestEmitNamesTableEmptyLabelsDeterministic(t *testing.T) {
	t.Parallel()

	notes := []domain.Note{{ID: "x.note", Title: "X Note"}}
	relationships := []domain.Relationship{{FromID: "x.a", ToID: "x.b"}}

	var out bytes.Buffer
	if err := emitNamesTable(&out, notes, relationships); err != nil {
		t.Fatalf("emit failed: %v", err)
	}
	want := readGolden(t, "list_names_empty_labels_table_output.golden")
	if out.String() != want {
		t.Fatalf("table mismatch\nwant: %q\n got: %q", want, out.String())
	}
}

func TestJoinSortedLabelsDeterministic(t *testing.T) {
	t.Parallel()

	got := joinSortedLabels("zeta,alpha, beta")
	if got != "alpha, beta, zeta" {
		t.Fatalf("unexpected labels %q", got)
	}
}

func TestEmitNamesJSONDeterministic(t *testing.T) {
	t.Parallel()

	notes, relationships, err := filterNames(listNamesValidatedApp(), "cli.", namesKindAll)
	if err != nil {
		t.Fatalf("filter failed: %v", err)
	}

	var first bytes.Buffer
	if err := emitNamesJSON(&first, notes, relationships); err != nil {
		t.Fatalf("emit first failed: %v", err)
	}
	var second bytes.Buffer
	if err := emitNamesJSON(&second, notes, relationships); err != nil {
		t.Fatalf("emit second failed: %v", err)
	}

	want := "{\"notes\":[{\"name\":\"cli.root\"},{\"name\":\"cli.worker\"}],\"relationships\":[{\"from\":\"app.db\",\"to\":\"cli.worker\"},{\"from\":\"cli.root\",\"to\":\"app.db\"}]}\n"
	if first.String() != want {
		t.Fatalf("json mismatch\nwant: %q\n got: %q", want, first.String())
	}
	if second.String() != first.String() {
		t.Fatalf("non-deterministic json output\nfirst: %q\nsecond: %q", first.String(), second.String())
	}
}

func listNamesValidatedApp() domain.ValidatedApp {
	return domain.ValidatedApp{
		Notes: []domain.Note{
			{ID: "cli.worker", Label: "cli.worker", Title: "CLI Worker", LabelsCSV: "worker,cli"},
			{ID: "app.db", Label: "app.db", Title: "Application DB", LabelsCSV: "database,persistent"},
			{ID: "cli.root", Label: "cli.root", Title: "CLI Root", LabelsCSV: "gateway,cli"},
		},
		Relationships: []domain.Relationship{
			{FromID: "cli.root", ToID: "app.db", Label: "depends_on", LabelsCSV: "gateway,depends_on"},
			{FromID: "app.db", ToID: "cli.worker", Label: "depends_on", LabelsCSV: "persistent,depends_on"},
			{FromID: "app.db", ToID: "infra.cache", Label: "depends_on", LabelsCSV: "cache_link,depends_on"},
		},
	}
}
