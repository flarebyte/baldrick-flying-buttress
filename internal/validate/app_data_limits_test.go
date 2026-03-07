package validate

import (
	"context"
	"strings"
	"testing"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/safety"
)

func TestAppDataValidatorReportsLimitExceeded(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{Source: "app", Reports: make([]domain.RawReport, safety.MaxReportsCount+1), Notes: []domain.RawNote{}, Relationships: []domain.RawRelationship{}}
	_, _, err := AppDataValidator{}.Validate(context.Background(), raw)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "reports count") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAppDataValidatorNotesLimitExceeded(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{Source: "app", Reports: []domain.RawReport{}, Notes: make([]domain.RawNote, safety.MaxNotesCount+1), Relationships: []domain.RawRelationship{}}
	_, _, err := AppDataValidator{}.Validate(context.Background(), raw)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "notes count") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAppDataValidatorRelationshipsLimitExceeded(t *testing.T) {
	t.Parallel()

	raw := domain.RawApp{Source: "app", Reports: []domain.RawReport{}, Notes: []domain.RawNote{}, Relationships: make([]domain.RawRelationship, safety.MaxRelationshipsCount+1)}
	_, _, err := AppDataValidator{}.Validate(context.Background(), raw)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "relationships count") {
		t.Fatalf("unexpected error: %v", err)
	}
}
