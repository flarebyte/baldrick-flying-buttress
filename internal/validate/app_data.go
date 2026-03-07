package validate

import (
	"fmt"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/ordering"
)

const (
	schemaValidationSource         = "validate.app.data.schema"
	registryValidationSource       = "validate.app.data.args.registry"
	configArgsValidationSource     = "validate.app.data.args.config"
	labelReferenceValidationSource = "labels.reference.validate"
	graphMissingNodesSource        = "graph.integrity.check.missing-nodes"
	graphOrphansSource             = "graph.integrity.check.orphans"
	graphDuplicateNoteNamesSource  = "graph.integrity.check.duplicate-note-names"
	graphCrossReportSource         = "graph.integrity.check.cross-report-references"
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

	v.step("args_validate_config")
	diagnostics = append(diagnostics, validateConfiguredArguments(rawModel, registry)...)

	v.step("labels_dataset_collect")
	datasetLabels := collectDatasetLabels(rawModel)

	v.step("labels_reference_validate")
	diagnostics = append(diagnostics, validateLabelReferences(rawModel, datasetLabels)...)

	v.step("graph_integrity_policy_resolve")
	graphPolicy := resolveGraphIntegrityPolicy(rawModel.GraphIntegrityPolicy)

	v.step("graph_integrity_validate")
	diagnostics = append(diagnostics, validateGraphIntegrity(rawModel, graphPolicy)...)

	v.step("diagnostics_collection")
	diagnostics = collectDiagnostics(diagnostics)

	v.step("validated_app_normalization")
	validated := normalizeValidatedApp(rawModel, registry, datasetLabels, graphPolicy)

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

func collectDatasetLabels(raw domain.RawApp) []string {
	labels := make([]string, 0)
	for _, note := range raw.Notes {
		labels = append(labels, note.Labels...)
	}
	for _, relationship := range raw.Relationships {
		labels = append(labels, relationship.Labels...)
		if relationship.Label != "" {
			labels = append(labels, relationship.Label)
		}
	}
	return normalizeAllowedValues(labels)
}

func validateLabelReferences(raw domain.RawApp, labelSet []string) []domain.Diagnostic {
	diagnostics := make([]domain.Diagnostic, 0)

	for i, report := range raw.Reports {
		for j, section := range report.Sections {
			for k, entry := range section.Arguments {
				name, value, ok := parseConfiguredArgument(entry)
				if !ok || !isLabelReferenceArgument(name) {
					continue
				}
				if containsString(labelSet, value) {
					continue
				}
				location := fmt.Sprintf("reports[%d].sections[%d].arguments[%d]", i, j, k)
				diagnostics = append(diagnostics, newLabelReferenceDiagnostic("LABEL_REF_UNKNOWN", "unknown label reference", location, report.Title, section.Title, "", name, value))
			}
		}
	}

	for i, note := range raw.Notes {
		for j, entry := range note.Arguments {
			name, value, ok := parseConfiguredArgument(entry)
			if !ok || !isLabelReferenceArgument(name) {
				continue
			}
			if containsString(labelSet, value) {
				continue
			}
			location := fmt.Sprintf("notes[%d].arguments[%d]", i, j)
			diagnostics = append(diagnostics, newLabelReferenceDiagnostic("LABEL_REF_UNKNOWN", "unknown label reference", location, "", "", note.Name, name, value))
		}
	}

	return diagnostics
}

func resolveGraphIntegrityPolicy(raw domain.RawGraphIntegrityPolicy) domain.GraphIntegrityPolicy {
	return domain.GraphIntegrityPolicy{
		MissingNode:          normalizePolicySeverity(raw.MissingNode, domain.PolicySeverityIgnore),
		OrphanNode:           normalizePolicySeverity(raw.OrphanNode, domain.PolicySeverityIgnore),
		DuplicateNoteName:    normalizePolicySeverity(raw.DuplicateNoteName, domain.PolicySeverityIgnore),
		CrossReportReference: normalizeCrossReportPolicy(raw.CrossReportReference),
	}
}

func normalizePolicySeverity(value string, fallback domain.PolicySeverity) domain.PolicySeverity {
	switch domain.PolicySeverity(strings.TrimSpace(value)) {
	case domain.PolicySeverityError, domain.PolicySeverityWarning, domain.PolicySeverityIgnore:
		return domain.PolicySeverity(strings.TrimSpace(value))
	default:
		return fallback
	}
}

func normalizeCrossReportPolicy(value string) domain.CrossReportPolicy {
	switch domain.CrossReportPolicy(strings.TrimSpace(value)) {
	case domain.CrossReportPolicyAllow, domain.CrossReportPolicyDisallow:
		return domain.CrossReportPolicy(strings.TrimSpace(value))
	default:
		return domain.CrossReportPolicyAllow
	}
}

func validateGraphIntegrity(raw domain.RawApp, policy domain.GraphIntegrityPolicy) []domain.Diagnostic {
	diagnostics := make([]domain.Diagnostic, 0)
	diagnostics = append(diagnostics, checkMissingNodes(raw, policy)...)
	diagnostics = append(diagnostics, checkOrphans(raw, policy.OrphanNode)...)
	diagnostics = append(diagnostics, checkDuplicateNoteNames(raw, policy.DuplicateNoteName)...)
	diagnostics = append(diagnostics, checkCrossReportReferences(raw, policy)...)
	return diagnostics
}

func checkMissingNodes(raw domain.RawApp, policy domain.GraphIntegrityPolicy) []domain.Diagnostic {
	if policy.MissingNode == domain.PolicySeverityIgnore {
		return []domain.Diagnostic{}
	}
	known := map[string]struct{}{}
	for _, note := range raw.Notes {
		name := strings.TrimSpace(note.Name)
		if name != "" {
			known[name] = struct{}{}
		}
	}
	diagnostics := make([]domain.Diagnostic, 0)
	for i, rel := range raw.Relationships {
		if rel.FromID != "" {
			if _, ok := known[rel.FromID]; !ok {
				diagnostics = append(diagnostics, newGraphDiagnostic(policy.MissingNode, graphMissingNodesSource, "GRAPH_MISSING_NODE", "relationship references missing from node", relationshipLocation(i, "from"), "", "", "", rel.FromID, rel.ToID))
			}
		}
		if rel.ToID != "" {
			if _, ok := known[rel.ToID]; !ok {
				diagnostics = append(diagnostics, newGraphDiagnostic(policy.MissingNode, graphMissingNodesSource, "GRAPH_MISSING_NODE", "relationship references missing to node", relationshipLocation(i, "to"), "", "", "", rel.FromID, rel.ToID))
			}
		}
	}
	return diagnostics
}

func checkOrphans(raw domain.RawApp, policy domain.PolicySeverity) []domain.Diagnostic {
	if policy == domain.PolicySeverityIgnore {
		return []domain.Diagnostic{}
	}
	connected := map[string]struct{}{}
	for _, rel := range raw.Relationships {
		if rel.FromID != "" {
			connected[rel.FromID] = struct{}{}
		}
		if rel.ToID != "" {
			connected[rel.ToID] = struct{}{}
		}
	}
	for _, report := range raw.Reports {
		for _, section := range report.Sections {
			for _, noteName := range section.Notes {
				n := strings.TrimSpace(noteName)
				if n != "" {
					connected[n] = struct{}{}
				}
			}
		}
	}
	diagnostics := make([]domain.Diagnostic, 0)
	for i, note := range raw.Notes {
		name := strings.TrimSpace(note.Name)
		if name == "" {
			continue
		}
		if _, ok := connected[name]; ok {
			continue
		}
		diagnostics = append(diagnostics, newGraphDiagnostic(policy, graphOrphansSource, "GRAPH_ORPHAN_NODE", "note is orphaned from relationships and report sections", noteLocation(i, "name"), "", "", name, "", ""))
	}
	return diagnostics
}

func checkDuplicateNoteNames(raw domain.RawApp, policy domain.PolicySeverity) []domain.Diagnostic {
	if policy == domain.PolicySeverityIgnore {
		return []domain.Diagnostic{}
	}
	counts := map[string]int{}
	for _, note := range raw.Notes {
		name := strings.TrimSpace(note.Name)
		if name != "" {
			counts[name]++
		}
	}
	names := make([]string, 0)
	for name, count := range counts {
		if count > 1 {
			names = append(names, name)
		}
	}
	slices.Sort(names)
	diagnostics := make([]domain.Diagnostic, 0, len(names))
	for _, name := range names {
		location := fmt.Sprintf("notes[name=%q]", name)
		diagnostics = append(diagnostics, newGraphDiagnostic(policy, graphDuplicateNoteNamesSource, "GRAPH_DUPLICATE_NOTE_NAME", "duplicate note name", location, "", "", name, "", ""))
	}
	return diagnostics
}

func checkCrossReportReferences(raw domain.RawApp, policy domain.GraphIntegrityPolicy) []domain.Diagnostic {
	if policy.CrossReportReference == domain.CrossReportPolicyAllow {
		return []domain.Diagnostic{}
	}
	membership := collectNoteReportMembership(raw)
	diagnostics := make([]domain.Diagnostic, 0)
	for i, rel := range raw.Relationships {
		if rel.FromID == "" || rel.ToID == "" {
			continue
		}
		fromReports := membership[rel.FromID]
		toReports := membership[rel.ToID]
		if len(fromReports) == 0 || len(toReports) == 0 {
			continue
		}
		if hasReportOverlap(fromReports, toReports) {
			continue
		}
		diagnostics = append(diagnostics, newGraphDiagnostic(domain.PolicySeverityError, graphCrossReportSource, "GRAPH_CROSS_REPORT_REFERENCE", "cross-report relationship reference is disallowed", fmt.Sprintf("relationships[%d]", i), "", "", "", rel.FromID, rel.ToID))
	}
	return diagnostics
}

func collectNoteReportMembership(raw domain.RawApp) map[string][]string {
	membership := map[string]map[string]struct{}{}
	for _, report := range raw.Reports {
		reportID := reportIDFromFilepath(report.Filepath)
		if reportID == "" {
			reportID = report.Title
		}
		for _, section := range report.Sections {
			for _, noteName := range section.Notes {
				name := strings.TrimSpace(noteName)
				if name == "" {
					continue
				}
				if _, ok := membership[name]; !ok {
					membership[name] = map[string]struct{}{}
				}
				membership[name][reportID] = struct{}{}
			}
		}
	}

	out := map[string][]string{}
	for noteName, reportSet := range membership {
		reports := make([]string, 0, len(reportSet))
		for reportID := range reportSet {
			reports = append(reports, reportID)
		}
		slices.Sort(reports)
		out[noteName] = reports
	}
	return out
}

func hasReportOverlap(a, b []string) bool {
	for _, left := range a {
		for _, right := range b {
			if left == right {
				return true
			}
		}
	}
	return false
}

func collectDiagnostics(diagnostics []domain.Diagnostic) []domain.Diagnostic {
	if diagnostics == nil {
		diagnostics = []domain.Diagnostic{}
	}
	return ordering.Diagnostics(diagnostics)
}

func normalizeValidatedApp(raw domain.RawApp, registry domain.ArgumentRegistry, datasetLabels []string, graphPolicy domain.GraphIntegrityPolicy) domain.ValidatedApp {
	reports := make([]domain.Report, 0, len(raw.Reports))
	for _, report := range raw.Reports {
		reports = append(reports, domain.Report{
			ID:    reportIDFromFilepath(report.Filepath),
			Title: report.Title,
		})
	}

	notes := make([]domain.Note, 0, len(raw.Notes))
	for _, note := range raw.Notes {
		notes = append(notes, domain.Note{
			ID:       note.Name,
			Label:    note.Title,
			Title:    note.Title,
			Markdown: note.Markdown,
		})
	}

	relationships := make([]domain.Relationship, 0, len(raw.Relationships))
	for _, relationship := range raw.Relationships {
		relationships = append(relationships, domain.Relationship{
			FromID: relationship.FromID,
			ToID:   relationship.ToID,
			Label:  relationship.Label,
		})
	}

	configDir := "."
	if raw.ConfigPath != "" {
		configDir = filepath.Dir(raw.ConfigPath)
	}

	return domain.ValidatedApp{
		Name:                 raw.Name,
		ConfigDir:            configDir,
		Modules:              raw.Modules,
		Reports:              ordering.Reports(reports),
		MarkdownReports:      normalizeMarkdownReports(raw),
		Notes:                ordering.Notes(notes),
		Relationships:        ordering.Relationships(relationships),
		Registry:             registry,
		DatasetLabels:        datasetLabels,
		GraphIntegrityPolicy: graphPolicy,
	}
}

func normalizeMarkdownReports(raw domain.RawApp) []domain.MarkdownReport {
	reports := make([]domain.MarkdownReport, 0, len(raw.Reports))
	for _, rawReport := range raw.Reports {
		report := domain.MarkdownReport{
			Title:       rawReport.Title,
			Filepath:    rawReport.Filepath,
			Description: rawReport.Description,
			Sections:    make([]domain.MarkdownH2Section, 0, len(rawReport.Sections)),
		}
		for _, rawH2 := range rawReport.Sections {
			h2 := domain.MarkdownH2Section{
				Title:       rawH2.Title,
				Description: rawH2.Description,
				Sections:    make([]domain.MarkdownH3Section, 0, len(rawH2.Sections)),
			}
			for _, rawH3 := range rawH2.Sections {
				h2.Sections = append(h2.Sections, domain.MarkdownH3Section{
					Title:       rawH3.Title,
					Description: rawH3.Description,
					NoteIDs:     ordering.Strings(rawH3.Notes),
				})
			}
			h2.Sections = ordering.MarkdownH3Sections(h2.Sections)
			report.Sections = append(report.Sections, h2)
		}
		report.Sections = ordering.MarkdownH2Sections(report.Sections)
		reports = append(reports, report)
	}
	return ordering.MarkdownReports(reports)
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

func containsScope(scopes []domain.ArgumentScope, target domain.ArgumentScope) bool {
	for _, scope := range scopes {
		if scope == target {
			return true
		}
	}
	return false
}

func parseConfiguredArgument(entry string) (string, string, bool) {
	parts := strings.SplitN(entry, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	if key == "" || value == "" {
		return "", "", false
	}
	return key, value, true
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

func newArgumentDiagnostic(code, message, location, reportTitle, sectionTitle, noteName, argumentName string) domain.Diagnostic {
	d := newDiagnostic(configArgsValidationSource, code, message, location)
	d.ReportTitle = reportTitle
	d.SectionTitle = sectionTitle
	d.NoteName = noteName
	d.ArgumentName = argumentName
	return d
}

func newLabelReferenceDiagnostic(code, message, location, reportTitle, sectionTitle, noteName, argumentName, labelValue string) domain.Diagnostic {
	d := newDiagnostic(labelReferenceValidationSource, code, message, location)
	d.Severity = domain.SeverityWarning
	d.ReportTitle = reportTitle
	d.SectionTitle = sectionTitle
	d.NoteName = noteName
	d.ArgumentName = argumentName
	d.LabelValue = labelValue
	return d
}

func newGraphDiagnostic(policy domain.PolicySeverity, source, code, message, location, reportTitle, sectionTitle, noteName, relationshipFrom, relationshipTo string) domain.Diagnostic {
	if policy == domain.PolicySeverityIgnore {
		return domain.Diagnostic{}
	}
	d := newDiagnostic(source, code, message, location)
	if policy == domain.PolicySeverityWarning {
		d.Severity = domain.SeverityWarning
	} else {
		d.Severity = domain.SeverityError
	}
	d.ReportTitle = reportTitle
	d.SectionTitle = sectionTitle
	d.NoteName = noteName
	d.RelationshipFrom = relationshipFrom
	d.RelationshipTo = relationshipTo
	return d
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
