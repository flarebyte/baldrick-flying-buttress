package validate

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
)

const (
	schemaValidationSource   = "validate.app.data.schema"
	registryValidationSource = "validate.app.data.args.registry"
)

type AppDataValidator struct {
	stepHook func(string)
}

func (v AppDataValidator) Validate(raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
	v.step("raw_model_normalization_precheck")
	rawModel := normalizeRaw(raw)

	v.step("schema_structure_validation")
	diagnostics := validateStructure(rawModel)

	v.step("args_registry_resolve")
	registry := resolveRegistry(rawModel.Registry)

	v.step("args_registry_validate")
	diagnostics = append(diagnostics, validateRegistry(rawModel.Registry)...)

	v.step("diagnostics_collection")
	diagnostics = collectDiagnostics(diagnostics)

	v.step("validated_app_normalization")
	validated := normalizeValidatedApp(rawModel, registry)

	return validated, domain.ValidationReport{Diagnostics: diagnostics}.Canonical(), nil
}

func (v AppDataValidator) step(name string) {
	if v.stepHook != nil {
		v.stepHook(name)
	}
}

func normalizeRaw(raw domain.RawApp) domain.RawApp {
	if raw.Modules == nil {
		raw.Modules = []string{}
	}
	return raw
}

func validateStructure(raw domain.RawApp) []domain.Diagnostic {
	diagnostics := make([]domain.Diagnostic, 0)

	if raw.Source == "" {
		diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV000", "missing required field: source", "source"))
	}

	if raw.Reports == nil {
		diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV100", "missing required collection: reports", "reports"))
	}
	if raw.Notes == nil {
		diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV200", "missing required collection: notes", "notes"))
	}
	if raw.Relationships == nil {
		diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV300", "missing required collection: relationships", "relationships"))
	}

	for i, report := range raw.Reports {
		if report.Title == "" {
			diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV101", "missing required field: report title", reportLocation(i, "title")))
		}
		if report.Filepath == "" {
			diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV102", "missing required field: report filepath", reportLocation(i, "filepath")))
		}
		if report.Sections == nil {
			diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV103", "missing required field: report sections", reportLocation(i, "sections")))
		}
		for j, section := range report.Sections {
			if strings.TrimSpace(section.Title) == "" {
				diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV104", "missing required field: section title", reportSectionLocation(i, j, "title")))
			}
		}
	}

	for i, note := range raw.Notes {
		if note.Name == "" {
			diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV201", "missing required field: note name", noteLocation(i, "name")))
		}
		if note.Title == "" {
			diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV202", "missing required field: note title", noteLocation(i, "title")))
		}
	}

	for i, relationship := range raw.Relationships {
		if relationship.FromID == "" {
			diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV301", "missing required field: relationship from", relationshipLocation(i, "from")))
		}
		if relationship.ToID == "" {
			diagnostics = append(diagnostics, newDiagnostic(schemaValidationSource, "FBV302", "missing required field: relationship to", relationshipLocation(i, "to")))
		}
	}

	return diagnostics
}

func resolveRegistry(raw domain.RawArgumentRegistry) domain.ArgumentRegistry {
	resolved := domain.ArgumentRegistry{
		Version:   raw.Version,
		Arguments: make([]domain.ArgumentDefinition, 0, len(raw.Arguments)),
	}

	for _, arg := range raw.Arguments {
		resolved.Arguments = append(resolved.Arguments, domain.ArgumentDefinition{
			Name:          strings.TrimSpace(arg.Name),
			ValueType:     domain.ArgumentValueType(strings.TrimSpace(arg.ValueType)),
			Scopes:        normalizeScopes(arg.Scopes),
			AllowedValues: normalizeAllowedValues(arg.AllowedValues),
			DefaultValue:  normalizeDefaultValue(arg.DefaultValue),
		})
	}

	slices.SortStableFunc(resolved.Arguments, func(a, b domain.ArgumentDefinition) int {
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		if a.ValueType < b.ValueType {
			return -1
		}
		if a.ValueType > b.ValueType {
			return 1
		}
		return 0
	})

	return resolved
}

func validateRegistry(raw domain.RawArgumentRegistry) []domain.Diagnostic {
	diagnostics := make([]domain.Diagnostic, 0)
	seenByName := map[string]int{}

	for i, arg := range raw.Arguments {
		name := strings.TrimSpace(arg.Name)
		locationBase := registryArgLocation(name, i)

		if name == "" {
			diagnostics = append(diagnostics, newDiagnostic(registryValidationSource, "FBR001", "missing argument name", locationBase+".name"))
		} else {
			seenByName[name]++
		}

		valueType := strings.TrimSpace(arg.ValueType)
		if valueType == "" {
			diagnostics = append(diagnostics, newDiagnostic(registryValidationSource, "FBR002", "missing argument value type", locationBase+".valueType"))
		} else if !isValidValueType(valueType) {
			diagnostics = append(diagnostics, newDiagnostic(registryValidationSource, "FBR003", "invalid argument value type", locationBase+".valueType"))
		}

		normalizedScopes := normalizeScopes(arg.Scopes)
		if len(normalizedScopes) == 0 {
			diagnostics = append(diagnostics, newDiagnostic(registryValidationSource, "FBR004", "missing argument scopes", locationBase+".scopes"))
		}
		for _, scope := range arg.Scopes {
			if !isValidScope(scope) {
				diagnostics = append(diagnostics, newDiagnostic(registryValidationSource, "FBR005", "invalid argument scope", locationBase+".scopes"))
			}
		}

		if valueType == string(domain.ArgumentValueTypeEnum) {
			normalizedAllowed, hadDuplicateAllowed := normalizeAllowedValuesWithDuplicateInfo(arg.AllowedValues)
			if hadDuplicateAllowed {
				diagnostics = append(diagnostics, newDiagnostic(registryValidationSource, "FBR007", "duplicate enum allowed values", locationBase+".allowedValues"))
			}
			if len(normalizedAllowed) == 0 {
				diagnostics = append(diagnostics, newDiagnostic(registryValidationSource, "FBR006", "invalid enum default/allowed-values combination", locationBase+".allowedValues"))
			}
			if arg.DefaultValue != nil {
				defaultValue, ok := arg.DefaultValue.(string)
				if !ok || !containsString(normalizedAllowed, defaultValue) {
					diagnostics = append(diagnostics, newDiagnostic(registryValidationSource, "FBR006", "invalid enum default/allowed-values combination", locationBase+".defaultValue"))
				}
			}
		}
	}

	duplicateNames := make([]string, 0)
	for name, count := range seenByName {
		if count > 1 {
			duplicateNames = append(duplicateNames, name)
		}
	}
	slices.Sort(duplicateNames)
	for _, name := range duplicateNames {
		diagnostics = append(diagnostics, newDiagnostic(registryValidationSource, "FBR000", "duplicate argument name", registryArgNameLocation(name)))
	}

	return diagnostics
}

func collectDiagnostics(diagnostics []domain.Diagnostic) []domain.Diagnostic {
	if diagnostics == nil {
		diagnostics = []domain.Diagnostic{}
	}
	return ordering.Diagnostics(diagnostics)
}

func normalizeValidatedApp(raw domain.RawApp, registry domain.ArgumentRegistry) domain.ValidatedApp {
	reports := make([]domain.Report, 0, len(raw.Reports))
	for _, report := range raw.Reports {
		reports = append(reports, domain.Report{
			ID:    reportIDFromFilepath(report.Filepath),
			Title: report.Title,
		})
	}

	notes := make([]domain.Note, 0, len(raw.Notes))
	for _, note := range raw.Notes {
		notes = append(notes, domain.Note{ID: note.Name, Label: note.Title})
	}

	relationships := make([]domain.Relationship, 0, len(raw.Relationships))
	for _, relationship := range raw.Relationships {
		relationships = append(relationships, domain.Relationship(relationship))
	}

	return domain.ValidatedApp{
		Name:          raw.Name,
		Modules:       raw.Modules,
		Reports:       ordering.Reports(reports),
		Notes:         ordering.Notes(notes),
		Relationships: ordering.Relationships(relationships),
		Registry:      registry,
	}
}

func reportIDFromFilepath(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

func normalizeScopes(scopes []string) []domain.ArgumentScope {
	seen := map[domain.ArgumentScope]struct{}{}
	out := make([]domain.ArgumentScope, 0)
	for _, raw := range scopes {
		scope := domain.ArgumentScope(strings.TrimSpace(raw))
		if !isValidScope(string(scope)) {
			continue
		}
		if _, exists := seen[scope]; exists {
			continue
		}
		seen[scope] = struct{}{}
		out = append(out, scope)
	}
	slices.Sort(out)
	return out
}

func normalizeAllowedValues(values []string) []string {
	normalized, _ := normalizeAllowedValuesWithDuplicateInfo(values)
	return normalized
}

func normalizeAllowedValuesWithDuplicateInfo(values []string) ([]string, bool) {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	hadDuplicate := false
	for _, v := range values {
		value := strings.TrimSpace(v)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			hadDuplicate = true
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	slices.Sort(out)
	return out, hadDuplicate
}

func normalizeDefaultValue(v any) *string {
	if v == nil {
		return nil
	}
	s, ok := v.(string)
	if !ok {
		return nil
	}
	s = strings.TrimSpace(s)
	return &s
}

func containsString(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func isValidValueType(valueType string) bool {
	switch valueType {
	case string(domain.ArgumentValueTypeString),
		string(domain.ArgumentValueTypeStrings),
		string(domain.ArgumentValueTypeBoolean),
		string(domain.ArgumentValueTypeInt),
		string(domain.ArgumentValueTypeFloat),
		string(domain.ArgumentValueTypeEnum):
		return true
	default:
		return false
	}
}

func isValidScope(scope string) bool {
	switch strings.TrimSpace(scope) {
	case string(domain.ArgumentScopeH3Section),
		string(domain.ArgumentScopeNote),
		string(domain.ArgumentScopeRenderer):
		return true
	default:
		return false
	}
}

func newDiagnostic(source, code, message, location string) domain.Diagnostic {
	return domain.Diagnostic{
		Code:     code,
		Severity: domain.SeverityError,
		Source:   source,
		Message:  message,
		Location: location,
		Path:     location,
	}
}

func reportLocation(i int, field string) string {
	return fmt.Sprintf("reports[%d].%s", i, field)
}

func reportSectionLocation(reportIndex, sectionIndex int, field string) string {
	return fmt.Sprintf("reports[%d].sections[%d].%s", reportIndex, sectionIndex, field)
}

func noteLocation(i int, field string) string {
	return fmt.Sprintf("notes[%d].%s", i, field)
}

func relationshipLocation(i int, field string) string {
	return fmt.Sprintf("relationships[%d].%s", i, field)
}

func registryArgLocation(name string, index int) string {
	if name != "" {
		return fmt.Sprintf("argumentRegistry.arguments[name=%q]", name)
	}
	return fmt.Sprintf("argumentRegistry.arguments[%d]", index)
}

func registryArgNameLocation(name string) string {
	return fmt.Sprintf("argumentRegistry.arguments[name=%q]", name)
}
