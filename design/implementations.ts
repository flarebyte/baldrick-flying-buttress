import type { ImplementationConsideration } from './common.ts';

// Initial implementation suggestions. Keep this list small and actionable.
export const implementations: Record<string, ImplementationConsideration> = {
  'lang.go': {
    name: 'lang.go',
    title: 'Implement the CLI in Go',
    description:
      'Use Go as the primary implementation language for strong typing, fast startup, and straightforward single-binary distribution.',
    calls: ['cli.root'],
  },
  'cli.cobra': {
    name: 'cli.cobra',
    title: 'Use Cobra for CLI command and argument parsing',
    description:
      'Use Cobra to model commands, flags, and subcommands (`generate markdown`, `generate json`, `validate`, `list reports`, `list names`, `lint names`) with a consistent command tree.',
    calls: [
      'cli.root',
      'action.generate.markdown',
      'action.generate.json',
      'action.validate',
      'action.list.reports',
      'action.list.names',
      'action.lint.names',
    ],
  },
  'names.list.command': {
    name: 'names.list.command',
    title: 'Implement list names inventory command',
    description:
      'Implement `flyb list names` with required `--prefix`, optional `--kind notes|relationships|all`, and `--format table|json` (default table); reuse validated app data and deterministic ordering before filtering/output.',
    calls: [
      'action.list.names',
      'load.app.data',
      'validate.app.data',
      'ordering.policy.resolve',
      'ordering.apply.deterministic',
      'names.filter.prefix',
      'names.filter.kind',
      'names.output.table',
      'names.output.json',
    ],
  },
  'names.lint.command': {
    name: 'names.lint.command',
    title: 'Implement lint names hygiene command',
    description:
      'Implement `flyb lint names` with style policy (`dot|snake|regex`), optional regex `--pattern`, optional prefix scope, and configurable severity; emit structured `NAME_STYLE_VIOLATION` diagnostics with canonical config locations and readable context.',
    calls: [
      'action.lint.names',
      'load.app.data',
      'validate.app.data',
      'ordering.policy.resolve',
      'ordering.apply.deterministic',
      'lint.names.policy.resolve',
      'names.filter.prefix',
      'lint.names.notes',
      'lint.names.relationships',
      'diagnostics.emit.structured',
    ],
  },
  'orphans.lint.command': {
    name: 'orphans.lint.command',
    title: 'Implement contextual orphan lint command',
    description:
      'Implement `flyb lint orphans` using orphan-query filters (`subject-label`, optional edge/counterpart labels, direction) and emit deterministic `ORPHAN_QUERY_MISSING_LINK` diagnostics with stable locations/context.',
    calls: [
      'action.lint.orphans',
      'load.app.data',
      'validate.app.data',
      'ordering.policy.resolve',
      'ordering.apply.deterministic',
      'lint.orphans.query.resolve',
      'orphans.query.find',
      'lint.orphans.emit',
      'diagnostics.emit.structured',
    ],
  },
  'orphans.report.section': {
    name: 'orphans.report.section',
    title: 'Implement contextual orphan report section renderer',
    description:
      'Implement H3 orphan section rendering using orphan-query arguments and deterministic row/table output (`name`, `title`, `labels`) so report sections and lint command evaluate the same orphan set.',
    calls: [
      'render.section.orphans',
      'args.orphan.query.resolve',
      'orphans.query.find',
      'orphans.render.rows',
      'ordering.apply.deterministic',
    ],
  },
  'renderer.registry.contract': {
    name: 'renderer.registry.contract',
    title: 'Define a renderer plugin registry contract',
    description:
      'Define a small renderer interface (name, supportsGraphShape, supportedArguments, render) and register built-ins (markdown-text, mermaid) in a deterministic lookup map; renderers consume one typed validated renderer-argument set resolved before plugin invocation.',
    calls: [
      'renderer.registry.resolve',
      'renderer.plugin.select',
      'render.section.graph',
      'render.graph.markdown.text',
      'render.graph.mermaid',
    ],
  },
  'renderer.selection.fallback-policy': {
    name: 'renderer.selection.fallback-policy',
    title: 'Use deterministic renderer selection and fallback policy',
    description:
      'Resolve renderer from renderer-scoped arguments sourced from H3Section and notes first, then apply stable defaults by graph shape (Mermaid-first for cyclic graphs, markdown-first for tree/DAG); if cycle-policy is `disallow` and cycles are detected, emit an error diagnostic and skip graph rendering for that section.',
    calls: [
      'renderer.plugin.select',
      'graph.shape.detect',
      'render.graph.tree-or-dag',
      'render.graph.circular',
    ],
  },
  'output.ordering.deterministic': {
    name: 'output.ordering.deterministic',
    title: 'Guarantee deterministic ordering in generated outputs',
    description:
      'Apply explicit stable sorting for notes, relationships, sections, and arguments using concrete comparators (notes: primaryLabel/name, relationships: from/to/labelsSortedJoined, sections: case-insensitive title plus originalIndex, arguments: name) so output remains reproducible across runs and machines.',
    calls: [
      'ordering.policy.resolve',
      'ordering.apply.deterministic',
      'action.generate.markdown.sections',
      'action.generate.markdown.section.h3',
    ],
  },
  'output.ordering.policy-contract': {
    name: 'output.ordering.policy-contract',
    title: 'Treat ordering policy as a testable contract',
    description:
      'Define ordering rules and tie-breakers as a versioned policy (including label normalization and relationship label joining rules) and verify them with golden-file tests.',
    calls: [
      'ordering.policy.resolve',
      'ordering.apply.deterministic',
      'action.generate.markdown',
    ],
  },
  'diagnostics.structured-model': {
    name: 'diagnostics.structured-model',
    title: 'Use a structured diagnostics model',
    description:
      'Standardize diagnostics with code, severity, source, message, canonical machine-readable location, and optional human-readable context fields to support CLI UX, CI checks, and future editor integrations.',
    calls: [
      'validate.app.data',
      'args.validate.runtime',
      'diagnostics.emit.structured',
    ],
  },
  'diagnostics.validation-location': {
    name: 'diagnostics.validation-location',
    title: 'Attach validation diagnostics to precise config locations',
    description:
      'Include canonical index-based CUE path plus related report/section titles and note/relationship/argument identifiers in diagnostics so users can quickly fix invalid configuration.',
    calls: [
      'validate.app.data',
      'action.validate',
      'diagnostics.emit.structured',
    ],
  },
  'graph.integrity.policy-model': {
    name: 'graph.integrity.policy-model',
    title: 'Define a graph integrity policy model',
    description:
      'Define explicit integrity rules for missing nodes, orphan nodes, duplicate names, and cross-report references with per-rule severity; validate label references separately against dataset-derived labels.',
    calls: [
      'graph.integrity.policy.resolve',
      'graph.integrity.validate',
      'validate.app.data',
    ],
  },
  'graph.integrity.validation-engine': {
    name: 'graph.integrity.validation-engine',
    title: 'Implement graph integrity validation checks',
    description:
      'Run focused integrity checks and emit structured diagnostics linked to note names, relationships, arguments, and CUE paths; keep label-definition handling free-form and validate only label references.',
    calls: [
      'graph.integrity.validate',
      'graph.integrity.check.missing-nodes',
      'graph.integrity.check.orphans',
      'graph.integrity.check.duplicate-note-names',
      'graph.integrity.check.cross-report-references',
      'labels.dataset.collect',
      'labels.reference.validate',
      'diagnostics.emit.structured',
    ],
  },
  'cli.arguments.typed-models': {
    name: 'cli.arguments.typed-models',
    title: 'Use free-form key/value arguments with typed coercion',
    description:
      'Treat H3Section and Note arguments like CLI-style flags (for example `format-csv=md`) to keep config flexible, then coerce values into typed runtime options per renderer.',
    calls: [
      'action.generate.markdown.sections',
      'action.generate.markdown.section.h3',
      'args.h3.resolve',
      'args.note.resolve',
      'render.section.plain',
      'render.section.file',
    ],
  },
  'cli.arguments.runtime-validation': {
    name: 'cli.arguments.runtime-validation',
    title: 'Validate free-form arguments with Cobra-style validators',
    description:
      'Validate argument keys and values at runtime against a known registry (enums, booleans, repeated values) using familiar CLI validation patterns and clear error messages.',
    calls: [
      'args.validate.runtime',
      'args.registry.resolve',
      'action.generate.markdown.section.h3',
      'render.section.file.csv',
      'render.section.graph',
    ],
  },
  'cli.arguments.registry-schema': {
    name: 'cli.arguments.registry-schema',
    title: 'Define a typed argument registry schema',
    description:
      'Maintain a registry of argument definitions (name, type, default, allowed values, scopes) and use it as the single source of truth for argument behavior; valid scopes are `h3-section`, `note`, and `renderer`.',
    calls: [
      'args.registry.resolve',
      'args.validate.runtime',
      'args.coerce.typed',
    ],
  },
  'cli.arguments.scope-resolution': {
    name: 'cli.arguments.scope-resolution',
    title: 'Apply scope-aware argument resolution',
    description:
      'Resolve and validate arguments by scope (h3-section, note, renderer) so options are accepted only where they are meaningful; renderer-scoped arguments are collected from H3Section and note argument lists with precedence `note` > `h3-section` > registry defaults.',
    calls: [
      'args.h3.resolve',
      'args.note.resolve',
      'args.validate.runtime',
      'graph.policy.cycle',
      'renderer.plugin.select',
    ],
  },
  'cli.arguments.type-coercion': {
    name: 'cli.arguments.type-coercion',
    title: 'Use explicit type coercion for free-form arguments',
    description:
      'After validation, coerce argument values to typed runtime options before rendering to avoid stringly-typed behavior in renderers.',
    calls: [
      'args.coerce.typed',
      'render.section.file.csv',
      'render.graph.mermaid',
      'render.graph.markdown.text',
    ],
  },
  'config.arguments.reduce-noise': {
    name: 'config.arguments.reduce-noise',
    title: 'Use arguments to reduce CUE configuration noise',
    description:
      'Prefer composable argument lists over adding many specialized CUE fields, so rendering capabilities can evolve without large schema churn.',
    calls: [
      'args.h3.resolve',
      'args.note.resolve',
      'action.generate.markdown.section.h3',
      'render.section.file',
    ],
  },
  'style.functions.small-single-purpose': {
    name: 'style.functions.small-single-purpose',
    title: 'Keep functions small and single-purpose',
    description:
      'Each function should do one thing and remain easy to test in isolation; prefer composition of small steps over large multi-branch handlers.',
    calls: [
      'action.generate.markdown.sections',
      'action.generate.markdown.section.h3',
      'render.section.file',
    ],
  },
  'style.io-separate-from-logic': {
    name: 'style.io-separate-from-logic',
    title: 'Separate I/O from core logic',
    description:
      'Keep parsing, filtering, and rendering logic pure where possible, and isolate file/network/process I/O behind adapter functions.',
    calls: [
      'load.app.data',
      'validate.app.data',
      'render.section.file.csv',
      'render.section.file.media',
      'action.generate.markdown',
    ],
  },
  'style.errors.guard-clauses': {
    name: 'style.errors.guard-clauses',
    title: 'Use early returns and guard clauses for errors',
    description:
      'Handle invalid inputs and failure states first, return immediately, and keep the success path shallow and readable.',
    calls: [
      'action.generate.markdown',
      'action.generate.json',
      'action.validate',
      'validate.app.data',
    ],
  },
  'style.parameters.tiny-structs': {
    name: 'style.parameters.tiny-structs',
    title: 'Use tiny structs to avoid long parameter lists',
    description:
      'Group related parameters into small intent-revealing structs (for example, render context and filter options) to reduce call-site ambiguity.',
    calls: [
      'action.generate.markdown.section.h3',
      'graph.select',
      'file.csv.filter',
      'render.section.file',
    ],
  },
  'style.predicates.named': {
    name: 'style.predicates.named',
    title: 'Replace boolean soup with named predicates',
    description:
      'Extract compound conditions into well-named predicate helpers to clarify branching and make tests easier to read.',
    calls: [
      'graph.select',
      'file.csv.filter',
      'render.section.file',
      'validate.app.data',
    ],
  },
  'config.cue': {
    name: 'config.cue',
    title: 'Use CUE as the configuration source of truth',
    description:
      'Represent notes, relationships, and report definitions in CUE for schema validation, defaults, and composable configuration.',
    calls: ['load.app.data', 'validate.app.data', 'action.validate'],
  },
};
