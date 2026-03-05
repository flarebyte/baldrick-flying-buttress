import { calls } from './calls.ts';
import {
  appendToReport,
  appendUseCases,
  displayCallsDetailed,
  getSetDifference,
  resetReport,
  toUseCaseSet,
} from './common.ts';
import { cliRoot } from './flows.ts';
import { mustUseCases, useCaseCatalogByName } from './use_cases.ts';

export const generateFlowDesignReport = async () => {
  // Build tree and header
  cliRoot({ level: 0 });
  await resetReport();
  await appendToReport('# FLOW DESIGN OVERVIEW (Generated)\n');
  await appendToReport('## Function calls tree\n');
  await appendToReport('```');
  await displayCallsDetailed(calls);
  await appendToReport('```\n');

  await appendToReport('## Validation Contract\n');
  await appendToReport(
    '- `validate.app.data` is the canonical validation entrypoint for all CLI commands and always includes dataset-based label-reference validation plus graph integrity policy resolution and graph integrity validation.',
  );
  await appendToReport(
    '- Guarantees: schema/structure validation, argument registry validation, configured free-form argument validation, dataset-based label-reference validation, graph integrity checks, and structured diagnostics collection with stable codes/severities/sources/locations.',
  );
  await appendToReport('- Return shape: `ValidatedApp` containing:');
  await appendToReport('  - normalized `notes`, `relationships`, and `reports`');
  await appendToReport(
    '  - resolved `graphIntegrityPolicy` and `argumentRegistry`',
  );
  await appendToReport(
    '  - optional ordering policy resolution (currently deferred to generation-time ordering components)',
  );
  await appendToReport(
    '  - `diagnostics: Diagnostic[]` always present (empty when no issues)',
  );
  await appendToReport(
    '- Generation block rule: any `error` severity diagnostic from `validate.app.data` blocks generation; warnings remain non-blocking by default but are still emitted consistently.',
  );
  await appendToReport(
    '- Label reference rule: labels on notes/relationships are free-form definitions; only label references are validated against dataset `labelSet` and unknown references emit `LABEL_REF_UNKNOWN` with default `warning` severity.',
  );
  await appendToReport('');
  await appendToReport('## Refactor Notes (Pseudo-diff)\n');
  await appendToReport(
    '- Removed direct `graph.integrity.validate` calls from `action.generate.markdown.section.h3` and `action.validate`.',
  );
  await appendToReport(
    '- `validate.app.data` now invokes: `validate.cue.schema`, `args.registry.resolve`, `args.registry.validate`, `args.validate.config`, `labels.dataset.collect`, `labels.reference.validate`, `graph.integrity.policy.resolve`, `graph.integrity.validate`, `diagnostics.collect.validation`, `app.model.normalize`.',
  );
  await appendToReport(
    '- Updated command flows to consume `ValidatedApp`:',
  );
  await appendToReport(
    '  - `action.generate.markdown`: `load.app.data -> validate.app.data -> action.generate.markdown.sections`',
  );
  await appendToReport(
    '  - `action.generate.json`: `load.app.data -> validate.app.data -> export.graph.json`',
  );
  await appendToReport(
    '  - `action.validate`: `load.app.data -> validate.app.data -> diagnostics.emit.structured`',
  );
  await appendToReport(
    '  - `action.list.reports`: `load.app.data -> validate.app.data -> list.reports.output`',
  );
  await appendToReport('');

  // Use-case coverage
  await appendUseCases(
    'Supported use cases:',
    toUseCaseSet(calls),
    useCaseCatalogByName,
  );
  const unsupported = getSetDifference(mustUseCases, toUseCaseSet(calls));
  if (unsupported.size > 0) {
    await appendUseCases(
      'Unsupported use cases (yet):',
      unsupported,
      useCaseCatalogByName,
    );
  }
};
