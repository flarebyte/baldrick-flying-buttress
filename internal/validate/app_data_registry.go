package validate

import (
	"slices"
	"strconv"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/textutil"
)

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
			diagnostics = append(diagnostics, newRegistryDiagnostic("FBR001", "missing argument name", locationBase+".name", name))
		} else {
			seenByName[name]++
		}

		valueType := strings.TrimSpace(arg.ValueType)
		if valueType == "" {
			diagnostics = append(diagnostics, newRegistryDiagnostic("FBR002", "missing argument value type", locationBase+".valueType", name))
		} else if !isValidValueType(valueType) {
			diagnostics = append(diagnostics, newRegistryDiagnostic("FBR003", "invalid argument value type", locationBase+".valueType", name))
		}

		normalizedScopes := normalizeScopes(arg.Scopes)
		if len(normalizedScopes) == 0 {
			diagnostics = append(diagnostics, newRegistryDiagnostic("FBR004", "missing argument scopes", locationBase+".scopes", name))
		}
		for _, scope := range arg.Scopes {
			if !isValidScope(scope) {
				diagnostics = append(diagnostics, newRegistryDiagnostic("FBR005", "invalid argument scope", locationBase+".scopes", name))
			}
		}

		if valueType == string(domain.ArgumentValueTypeEnum) {
			normalizedAllowed, hadDuplicateAllowed := normalizeAllowedValuesWithDuplicateInfo(arg.AllowedValues)
			if hadDuplicateAllowed {
				diagnostics = append(diagnostics, newRegistryDiagnostic("FBR007", "duplicate enum allowed values", locationBase+".allowedValues", name))
			}
			if len(normalizedAllowed) == 0 {
				diagnostics = append(diagnostics, newRegistryDiagnostic("FBR006", "invalid enum default/allowed-values combination", locationBase+".allowedValues", name))
			}
			if arg.DefaultValue != nil {
				defaultValue, ok := arg.DefaultValue.(string)
				if !ok || !containsString(normalizedAllowed, defaultValue) {
					diagnostics = append(diagnostics, newRegistryDiagnostic("FBR006", "invalid enum default/allowed-values combination", locationBase+".defaultValue", name))
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
		diagnostics = append(diagnostics, newRegistryDiagnostic("FBR000", "duplicate argument name", registryArgNameLocation(name), name))
	}

	return diagnostics
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

func containsScope(scopes []domain.ArgumentScope, target domain.ArgumentScope) bool {
	for _, scope := range scopes {
		if scope == target {
			return true
		}
	}
	return false
}

func parseConfiguredArgument(entry string) (string, string, bool) {
	return textutil.ParseKeyValue(entry)
}

func valueMatchesType(value string, def domain.ArgumentDefinition) bool {
	switch def.ValueType {
	case domain.ArgumentValueTypeString:
		return value != ""
	case domain.ArgumentValueTypeStrings:
		parts := strings.Split(value, ",")
		if len(parts) == 0 {
			return false
		}
		for _, part := range parts {
			if strings.TrimSpace(part) == "" {
				return false
			}
		}
		return true
	case domain.ArgumentValueTypeBoolean:
		return value == "true" || value == "false"
	case domain.ArgumentValueTypeInt:
		_, err := strconv.Atoi(value)
		return err == nil
	case domain.ArgumentValueTypeFloat:
		_, err := strconv.ParseFloat(value, 64)
		return err == nil
	case domain.ArgumentValueTypeEnum:
		return containsString(def.AllowedValues, value)
	default:
		return false
	}
}

func isLabelReferenceArgument(name string) bool {
	return strings.HasSuffix(name, "-label")
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
