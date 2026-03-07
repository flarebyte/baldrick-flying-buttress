package validate

import (
	"fmt"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
)

func validateConfiguredArguments(raw domain.RawApp, registry domain.ArgumentRegistry) []domain.Diagnostic {
	diagnostics := make([]domain.Diagnostic, 0)
	byName := map[string]domain.ArgumentDefinition{}
	for _, arg := range registry.Arguments {
		byName[arg.Name] = arg
	}

	for i, report := range raw.Reports {
		for j, section := range report.Sections {
			for k, entry := range section.Arguments {
				location := fmt.Sprintf("reports[%d].sections[%d].arguments[%d]", i, j, k)
				name, value, ok := parseConfiguredArgument(entry)
				if !ok {
					diagnostics = append(diagnostics, newArgumentDiagnostic("FBC001", "malformed configured argument", location, report.Title, section.Title, "", ""))
					continue
				}
				def, exists := byName[name]
				if !exists {
					diagnostics = append(diagnostics, newArgumentDiagnostic("FBC002", "unknown configured argument key", location, report.Title, section.Title, "", name))
					continue
				}
				if !containsScope(def.Scopes, domain.ArgumentScopeH3Section) && !containsScope(def.Scopes, domain.ArgumentScopeRenderer) {
					diagnostics = append(diagnostics, newArgumentDiagnostic("FBC003", "argument used outside allowed scope", location, report.Title, section.Title, "", name))
				}
				if !valueMatchesType(value, def) {
					diagnostics = append(diagnostics, newArgumentDiagnostic("FBC004", "argument value does not match declared type", location, report.Title, section.Title, "", name))
				}
			}
		}
	}

	for i, note := range raw.Notes {
		for j, entry := range note.Arguments {
			location := fmt.Sprintf("notes[%d].arguments[%d]", i, j)
			name, value, ok := parseConfiguredArgument(entry)
			if !ok {
				diagnostics = append(diagnostics, newArgumentDiagnostic("FBC001", "malformed configured argument", location, "", "", note.Name, ""))
				continue
			}
			def, exists := byName[name]
			if !exists {
				diagnostics = append(diagnostics, newArgumentDiagnostic("FBC002", "unknown configured argument key", location, "", "", note.Name, name))
				continue
			}
			if !containsScope(def.Scopes, domain.ArgumentScopeNote) {
				diagnostics = append(diagnostics, newArgumentDiagnostic("FBC003", "argument used outside allowed scope", location, "", "", note.Name, name))
			}
			if !valueMatchesType(value, def) {
				diagnostics = append(diagnostics, newArgumentDiagnostic("FBC004", "argument value does not match declared type", location, "", "", note.Name, name))
			}
		}
	}

	return diagnostics
}
