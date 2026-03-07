package ordering

import (
	"slices"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func Diagnostics(in []domain.Diagnostic) []domain.Diagnostic {
	out := slices.Clone(in)
	slices.SortStableFunc(out, func(a, b domain.Diagnostic) int {
		if v := cmpString(a.Code, b.Code); v != 0 {
			return v
		}
		if v := cmpString(string(a.Severity), string(b.Severity)); v != 0 {
			return v
		}
		if v := cmpString(a.Source, b.Source); v != 0 {
			return v
		}
		if v := cmpString(a.Location, b.Location); v != 0 {
			return v
		}
		if v := cmpString(a.Path, b.Path); v != 0 {
			return v
		}
		if v := cmpString(a.ArgumentName, b.ArgumentName); v != 0 {
			return v
		}
		return cmpString(a.Message, b.Message)
	})
	return out
}

func Reports(in []domain.Report) []domain.Report {
	out := slices.Clone(in)
	slices.SortStableFunc(out, func(a, b domain.Report) int {
		if v := cmpString(a.ID, b.ID); v != 0 {
			return v
		}
		return cmpString(a.Title, b.Title)
	})
	return out
}

func Notes(in []domain.Note) []domain.Note {
	out := slices.Clone(in)
	slices.SortStableFunc(out, func(a, b domain.Note) int {
		if v := cmpString(a.Label, b.Label); v != 0 {
			return v
		}
		return cmpString(a.ID, b.ID)
	})
	return out
}

func Relationships(in []domain.Relationship) []domain.Relationship {
	out := slices.Clone(in)
	slices.SortStableFunc(out, func(a, b domain.Relationship) int {
		if v := cmpString(a.FromID, b.FromID); v != 0 {
			return v
		}
		if v := cmpString(a.ToID, b.ToID); v != 0 {
			return v
		}
		return cmpString(a.Label, b.Label)
	})
	return out
}

func cmpString(a, b string) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
