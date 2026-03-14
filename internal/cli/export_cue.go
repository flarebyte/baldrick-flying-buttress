package cli

import (
	"context"
	"fmt"
	"io"
	"slices"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
)

type exportCueAction struct {
	out io.Writer
}

func runExportCue(ctx context.Context, loaderFactory LoaderFactory, configPath *string, reportIDs []string, out io.Writer) error {
	if loaderFactory == nil {
		return fmt.Errorf("loader factory is required")
	}
	if configPath == nil {
		return fmt.Errorf("config path is required")
	}
	loader := loaderFactory(*configPath)
	loader = withReportFilter(loader, reportIDs)
	raw, err := loader.Load(ctx)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return exportCueAction{out: out}.execute(raw)
}

func (a exportCueAction) execute(raw domain.RawApp) error {
	data, err := encodeRawAppToCue(canonicalRawApp(raw))
	if err != nil {
		return err
	}
	_, err = a.out.Write(data)
	return err
}

func encodeRawAppToCue(raw domain.RawApp) ([]byte, error) {
	ctx := cuecontext.New()
	expr, ok := ctx.Encode(raw).Syntax(cue.Final()).(*ast.StructLit)
	if !ok {
		return nil, fmt.Errorf("raw app encoding did not produce a struct literal")
	}
	file := &ast.File{
		Decls: append([]ast.Decl{
			&ast.Package{Name: ast.NewIdent("flyb")},
		}, expr.Elts...),
	}
	out, err := format.Node(file)
	if err != nil {
		return nil, err
	}
	if len(out) == 0 || out[len(out)-1] != '\n' {
		out = append(out, '\n')
	}
	return out, nil
}

func canonicalRawApp(raw domain.RawApp) domain.RawApp {
	raw.Modules = ordering.Strings(raw.Modules)
	raw.Reports = canonicalRawReports(raw.Reports)
	raw.Notes = canonicalRawNotes(raw.Notes)
	raw.Relationships = canonicalRawRelationships(raw.Relationships)
	raw.Registry = canonicalRawRegistry(raw.Registry)
	return raw
}

func canonicalRawReports(in []domain.RawReport) []domain.RawReport {
	out := make([]domain.RawReport, 0, len(in))
	for _, report := range in {
		report.Sections = canonicalRawSections(report.Sections)
		out = append(out, report)
	}
	slices.SortStableFunc(out, func(a, b domain.RawReport) int {
		if v := strings.Compare(a.Filepath, b.Filepath); v != 0 {
			return v
		}
		return strings.Compare(a.Title, b.Title)
	})
	return out
}

func canonicalRawSections(in []domain.RawReportSection) []domain.RawReportSection {
	out := make([]domain.RawReportSection, 0, len(in))
	for _, section := range in {
		section.Arguments = ordering.Strings(section.Arguments)
		section.Notes = ordering.Strings(section.Notes)
		section.Sections = canonicalRawSections(section.Sections)
		out = append(out, section)
	}
	slices.SortStableFunc(out, func(a, b domain.RawReportSection) int {
		if v := strings.Compare(a.Title, b.Title); v != 0 {
			return v
		}
		return strings.Compare(a.Description, b.Description)
	})
	return out
}

func canonicalRawNotes(in []domain.RawNote) []domain.RawNote {
	out := make([]domain.RawNote, 0, len(in))
	for _, note := range in {
		note.Arguments = ordering.Strings(note.Arguments)
		note.Labels = ordering.Strings(note.Labels)
		out = append(out, note)
	}
	slices.SortStableFunc(out, func(a, b domain.RawNote) int {
		if v := strings.Compare(a.Title, b.Title); v != 0 {
			return v
		}
		return strings.Compare(a.Name, b.Name)
	})
	return out
}

func canonicalRawRelationships(in []domain.RawRelationship) []domain.RawRelationship {
	out := make([]domain.RawRelationship, 0, len(in))
	for _, rel := range in {
		rel.Labels = ordering.Strings(rel.Labels)
		out = append(out, rel)
	}
	slices.SortStableFunc(out, func(a, b domain.RawRelationship) int {
		if v := strings.Compare(a.FromID, b.FromID); v != 0 {
			return v
		}
		if v := strings.Compare(a.ToID, b.ToID); v != 0 {
			return v
		}
		return strings.Compare(a.Label, b.Label)
	})
	return out
}

func canonicalRawRegistry(in domain.RawArgumentRegistry) domain.RawArgumentRegistry {
	in.Arguments = canonicalRawArgumentDefinitions(in.Arguments)
	return in
}

func canonicalRawArgumentDefinitions(in []domain.RawArgumentDefinition) []domain.RawArgumentDefinition {
	out := make([]domain.RawArgumentDefinition, 0, len(in))
	for _, arg := range in {
		arg.Scopes = ordering.Strings(arg.Scopes)
		arg.AllowedValues = ordering.Strings(arg.AllowedValues)
		out = append(out, arg)
	}
	slices.SortStableFunc(out, func(a, b domain.RawArgumentDefinition) int {
		return strings.Compare(a.Name, b.Name)
	})
	return out
}
