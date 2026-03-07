package validate

import "github.com/flarebyte/baldrick-flying-buttress/internal/domain"

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
