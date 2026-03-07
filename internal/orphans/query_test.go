package orphans

import (
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func TestFindSubjectLabelFiltering(t *testing.T) {
	t.Parallel()
	app := queryFixtureApp()
	query := Query{SubjectLabel: "ingredient", Direction: DirectionEither}
	orphans := Find(app, query)
	if len(orphans) != 1 || orphans[0].ID != "n.ingredient.orphan" {
		t.Fatalf("unexpected orphans: %#v", orphans)
	}
}

func TestFindEdgeLabelFiltering(t *testing.T) {
	t.Parallel()
	app := queryFixtureApp()
	query := Query{SubjectLabel: "ingredient", EdgeLabel: "uses", Direction: DirectionEither}
	orphans := Find(app, query)
	if len(orphans) != 3 {
		t.Fatalf("expected 3 orphans, got %#v", orphans)
	}
}

func TestFindCounterpartLabelFiltering(t *testing.T) {
	t.Parallel()
	app := queryFixtureApp()
	query := Query{SubjectLabel: "ingredient", CounterpartLabel: "tool", Direction: DirectionEither}
	orphans := Find(app, query)
	if len(orphans) != 1 || orphans[0].ID != "n.ingredient.orphan" {
		t.Fatalf("unexpected orphans: %#v", orphans)
	}
}

func TestFindDirectionFiltering(t *testing.T) {
	t.Parallel()
	app := queryFixtureApp()

	inQuery := Query{SubjectLabel: "ingredient", Direction: DirectionIn}
	inOrphans := Find(app, inQuery)
	if len(inOrphans) != 2 {
		t.Fatalf("expected 2 in-direction orphans, got %#v", inOrphans)
	}

	outQuery := Query{SubjectLabel: "ingredient", Direction: DirectionOut}
	outOrphans := Find(app, outQuery)
	if len(outOrphans) != 2 || outOrphans[0].ID != "n.ingredient.b" || outOrphans[1].ID != "n.ingredient.orphan" {
		t.Fatalf("unexpected out-direction orphans: %#v", outOrphans)
	}

	eitherQuery := Query{SubjectLabel: "ingredient", Direction: DirectionEither}
	eitherOrphans := Find(app, eitherQuery)
	if len(eitherOrphans) != 1 || eitherOrphans[0].ID != "n.ingredient.orphan" {
		t.Fatalf("unexpected either-direction orphans: %#v", eitherOrphans)
	}
}

func TestFindMultipleOrphansDetected(t *testing.T) {
	t.Parallel()
	app := queryFixtureApp()
	query := Query{SubjectLabel: "ingredient", EdgeLabel: "missing", Direction: DirectionEither}
	orphans := Find(app, query)
	if len(orphans) != 3 {
		t.Fatalf("expected 3 orphans, got %#v", orphans)
	}
}

func queryFixtureApp() domain.ValidatedApp {
	return domain.ValidatedApp{
		Notes: []domain.Note{
			{ID: "n.ingredient.a", Label: "ingredient"},
			{ID: "n.ingredient.b", Label: "ingredient"},
			{ID: "n.ingredient.orphan", Label: "ingredient"},
			{ID: "n.tool.a", Label: "tool"},
		},
		Relationships: []domain.Relationship{
			{FromID: "n.ingredient.a", ToID: "n.tool.a", Label: "depends_on"},
			{FromID: "n.tool.a", ToID: "n.ingredient.b", Label: "feeds"},
		},
	}
}
