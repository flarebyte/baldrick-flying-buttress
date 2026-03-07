package validate

import (
	"context"
	"fmt"

	"github.com/flarebyte/baldrick-flying-buttress/internal/domain"
	"github.com/flarebyte/baldrick-flying-buttress/internal/safety"
)

const (
	schemaValidationSource         = "validate.cue.schema"
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

func (v AppDataValidator) Validate(ctx context.Context, raw domain.RawApp) (domain.ValidatedApp, domain.ValidationReport, error) {
	if err := ctx.Err(); err != nil {
		return domain.ValidatedApp{}, domain.ValidationReport{}, err
	}
	v.step("raw_model_normalization_precheck")
	rawModel := normalizeRaw(raw)
	if err := enforceRawModelLimits(rawModel); err != nil {
		return domain.ValidatedApp{}, domain.ValidationReport{}, err
	}

	v.step("validate_cue_schema")
	diagnostics := validateStructure(rawModel)
	if err := ctx.Err(); err != nil {
		return domain.ValidatedApp{}, domain.ValidationReport{}, err
	}

	v.step("args_registry_resolve")
	registry := resolveRegistry(rawModel.Registry)

	v.step("args_registry_validate")
	diagnostics = append(diagnostics, validateRegistry(rawModel.Registry)...)
	if err := ctx.Err(); err != nil {
		return domain.ValidatedApp{}, domain.ValidationReport{}, err
	}

	v.step("args_validate_config")
	diagnostics = append(diagnostics, validateConfiguredArguments(rawModel, registry)...)
	if err := ctx.Err(); err != nil {
		return domain.ValidatedApp{}, domain.ValidationReport{}, err
	}

	v.step("labels_dataset_collect")
	datasetLabels := collectDatasetLabels(rawModel)

	v.step("labels_reference_validate")
	diagnostics = append(diagnostics, validateLabelReferences(rawModel, datasetLabels)...)
	if err := ctx.Err(); err != nil {
		return domain.ValidatedApp{}, domain.ValidationReport{}, err
	}

	v.step("graph_integrity_policy_resolve")
	graphPolicy := resolveGraphIntegrityPolicy(rawModel.GraphIntegrityPolicy)

	v.step("graph_integrity_validate")
	diagnostics = append(diagnostics, validateGraphIntegrity(rawModel, graphPolicy)...)
	if err := ctx.Err(); err != nil {
		return domain.ValidatedApp{}, domain.ValidationReport{}, err
	}

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

func enforceRawModelLimits(raw domain.RawApp) error {
	if err := safety.CheckReportsCount(len(raw.Reports)); err != nil {
		return fmt.Errorf("raw app limit exceeded: %w", err)
	}
	if err := safety.CheckNotesCount(len(raw.Notes)); err != nil {
		return fmt.Errorf("raw app limit exceeded: %w", err)
	}
	if err := safety.CheckRelationshipsCount(len(raw.Relationships)); err != nil {
		return fmt.Errorf("raw app limit exceeded: %w", err)
	}
	return nil
}
