import { calls } from './calls';
import { type ComponentCall, type FlowContext, incrContext } from './common';
import { useCases } from './use_cases.ts';

const uc = {
  reportGenerate: useCases['cli.report.generate'].name,
  noteBasicMarkdown: useCases['cli.note.basic-markdown'].name,
  noteFilepathReference: useCases['cli.note.filepath.reference'].name,
  noteCsvEmbed: useCases['cli.note.csv.embed'].name,
  noteCsvFilterColumn: useCases['cli.note.csv.filter-column'].name,
  noteImagePreview: useCases['cli.note.image.preview'].name,
  noteMermaidEmbed: useCases['cli.note.mermaid.embed'].name,
  noteLinkMarkdown: useCases['cli.note.link.markdown'].name,
  configRelationshipsLabeled: useCases['cli.config.relationships.labeled'].name,
  configReportsMultiple: useCases['cli.config.reports.multiple'].name,
  reportList: useCases['cli.report.list'].name,
  namesList: useCases['cli.names.list'].name,
  namesLint: useCases['cli.names.lint'].name,
  namesStylePolicy: useCases['cli.names.style-policy'].name,
  namesPrefixFilter: useCases['cli.names.prefix-filter'].name,
  namesOutputFormats: useCases['cli.names.output-formats'].name,
  exportJsonGraph: useCases['cli.export.json.graph'].name,
  reportSubgraphByLabel: useCases['cli.report.subgraph.by-label'].name,
  sectionH3CyclePolicy: useCases['cli.section.h3.cycle-policy'].name,
  reportGraphShapeAwareRender:
    useCases['cli.report.graph.shape-aware-render'].name,
  reportGraphRendererMarkdownText:
    useCases['cli.report.graph.renderer.markdown-text'].name,
  reportGraphRendererMermaid:
    useCases['cli.report.graph.renderer.mermaid'].name,
  rendererRegistry: useCases['cli.renderer.registry'].name,
  rendererPluginSelection: useCases['cli.renderer.plugin-selection'].name,
  outputDeterministicOrdering:
    useCases['cli.output.deterministic-ordering'].name,
  outputDeterministicOrderingPolicy:
    useCases['cli.output.deterministic-ordering.policy'].name,
  diagnosticsModel: useCases['cli.diagnostics.model'].name,
  diagnosticsValidation: useCases['cli.diagnostics.validation'].name,
  graphIntegrityPolicy: useCases['cli.graph.integrity.policy'].name,
  graphIntegrityValidation: useCases['cli.graph.integrity.validation'].name,
  argumentsFreeForm: useCases['cli.arguments.free-form'].name,
  argumentsRuntimeValidation: useCases['cli.arguments.runtime-validation'].name,
  argumentsRegistrySchema: useCases['cli.arguments.registry.schema'].name,
  argumentsScopeResolution:
    useCases['cli.arguments.registry.scope-resolution'].name,
  argumentsTypeCoercion: useCases['cli.arguments.type-coercion'].name,
  configReduceNoiseWithArgs: useCases['cli.config.reduce-noise.with-args'].name,
};

export const cliRoot = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'cli.root',
    title: 'flyb CLI root command',
    note: 'Entry point for report generation, listing, JSON export, and config validation.',
    level: context.level,
    useCases: [uc.reportGenerate, uc.reportList, uc.configReportsMultiple],
  };
  calls.push(call);
  listReportsAction(incrContext(context));
  listNamesAction(incrContext(context));
  generateMarkdownAction(incrContext(context));
  generateJsonAction(incrContext(context));
  validateAction(incrContext(context));
  lintNamesAction(incrContext(context));
};

export const listReportsAction = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'action.list.reports',
    title: 'List configured markdown reports',
    note: 'Enumerate report targets from the validated application model without generating files.',
    level: context.level,
    useCases: [uc.reportList, uc.configReportsMultiple],
  };
  calls.push(call);
  loadAppData(incrContext(context));
  validateAppData(incrContext(context));
  listValidatedReports(incrContext(context));
};

export const listNamesAction = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'action.list.names',
    title: 'List note and relationship names',
    note: 'Print note and relationship identifiers for daily inventory with required `--prefix`, optional `--kind notes|relationships|all`, and `--format table|json` (default table).',
    level: context.level,
    useCases: [
      uc.namesList,
      uc.namesPrefixFilter,
      uc.namesOutputFormats,
      uc.outputDeterministicOrdering,
      uc.outputDeterministicOrderingPolicy,
    ],
  };
  calls.push(call);
  loadAppData(incrContext(context));
  validateAppData(incrContext(context));
  resolveOrderingPolicy(incrContext(context));
  applyDeterministicOrdering(incrContext(context));
  filterNamesByPrefix(incrContext(context));
  filterNamesByKind(incrContext(context));
  outputNamesTable(incrContext(context));
  outputNamesJson(incrContext(context));
};

export const generateMarkdownAction = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'action.generate.markdown',
    title: 'Generate markdown reports',
    note: 'Renders one or more markdown outputs from a single validated application model.',
    level: context.level,
    useCases: [uc.reportGenerate, uc.configReportsMultiple],
  };
  calls.push(call);
  loadAppData(incrContext(context));
  validateAppData(incrContext(context));
  generateMarkdownSections(incrContext(context));
};

export const generateMarkdownSections = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'action.generate.markdown.sections',
    title: 'Generate markdown sections',
    note: 'Build H3 sections from note subsets and renderers with deterministic ordering.',
    level: context.level,
    useCases: [
      uc.reportGenerate,
      uc.noteBasicMarkdown,
      uc.outputDeterministicOrdering,
      uc.outputDeterministicOrderingPolicy,
    ],
  };
  calls.push(call);
  resolveOrderingPolicy(incrContext(context));
  generateSingleH3Section(incrContext(context));
};

export const renderPlainSection = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'render.section.plain',
    title: 'Render plain section',
    note: 'Render title and markdown body, including markdown links and note-level argument options.',
    level: context.level,
    useCases: [
      uc.noteBasicMarkdown,
      uc.noteLinkMarkdown,
      uc.argumentsFreeForm,
      uc.argumentsRuntimeValidation,
    ],
  };
  calls.push(call);
};

export const renderGraphSection = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'render.section.graph',
    title: 'Render section as a graph',
    note: 'Resolve cycle policy and graph shape, then render with selected renderer(s); renderer/runtime diagnostics here must not duplicate graph-integrity diagnostics.',
    level: context.level,
    useCases: [
      uc.configRelationshipsLabeled,
      uc.reportSubgraphByLabel,
      uc.sectionH3CyclePolicy,
      uc.reportGraphShapeAwareRender,
      uc.rendererRegistry,
      uc.rendererPluginSelection,
      uc.outputDeterministicOrdering,
    ],
  };
  calls.push(call);
  resolveRendererRegistry(incrContext(context));
  resolveRendererScopedArguments(incrContext(context));
  selectRendererPlugin(incrContext(context));
  resolveH3GraphCyclePolicy(incrContext(context));
  detectGraphShape(incrContext(context));
  renderGraphTreeOrDag(incrContext(context));
  renderGraphCircular(incrContext(context));
};
export const renderSectionWithFile = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'render.section.file',
    title: 'Render section with referenced file content',
    note: 'Dispatches file rendering by type (CSV, media, code/diagram).',
    level: context.level,
    useCases: [
      uc.noteFilepathReference,
      uc.argumentsFreeForm,
      uc.outputDeterministicOrdering,
    ],
  };
  calls.push(call);
  resolveNoteRenderArguments(incrContext(context));
  resolveArgumentRegistry(incrContext(context));
  validateRenderArguments(incrContext(context));
  coerceRenderArguments(incrContext(context));
  renderSectionWithCsvFile(incrContext(context));
  renderSectionWithMedia(incrContext(context));
  renderSectionWithCodeSnippet(incrContext(context));
};

export const renderSectionWithCsvFile = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'render.section.file.csv',
    title: 'Render section with CSV file',
    note: 'Render as a markdown table or raw CSV code block (for example `format-csv=md`) and apply note-scoped CSV filters (`csv-include` / `csv-exclude`) using `column:value` exact-match rules.',
    level: context.level,
    useCases: [
      uc.noteFilepathReference,
      uc.noteCsvEmbed,
      uc.argumentsFreeForm,
      uc.argumentsRuntimeValidation,
    ],
  };
  calls.push(call);
  filterCsvFile(incrContext(context));
};

export const renderSectionWithCodeSnippet = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'render.section.file.code',
    title: 'Render section with code or Mermaid snippet',
    note: 'Preserve fenced-block formatting for code and Mermaid content.',
    level: context.level,
    useCases: [
      uc.noteFilepathReference,
      uc.noteMermaidEmbed,
      uc.argumentsFreeForm,
      uc.argumentsRuntimeValidation,
    ],
  };
  calls.push(call);
};
export const filterCsvFile = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'file.csv.filter',
    title: 'Filter CSV rows by column',
    note: 'Apply exact-match include/exclude filters before rendering CSV output: `csv-include=column:value` keeps matching rows, `csv-exclude=column:value` removes matching rows, and multiple filters are allowed.',
    level: context.level,
    useCases: [uc.noteCsvFilterColumn, uc.argumentsFreeForm],
  };
  calls.push(call);
};

export const selectSubGraph = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'graph.select',
    title: 'Extract subgraph using labels',
    note: 'Filter notes and relationships by labels and optional starting node; label references are pre-validated against dataset labels (union of note.labels and relationship.labels).',
    level: context.level,
    useCases: [uc.configRelationshipsLabeled, uc.reportSubgraphByLabel],
  };
  calls.push(call);
};

export const resolveH3GraphCyclePolicy = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'graph.policy.cycle',
    title: 'Resolve H3Section cycle policy argument',
    note: 'Resolve section cycle policy (`disallow` or `allow`): `disallow` requires cycle detection error diagnostics and blocks section graph rendering; `allow` permits cyclic rendering.',
    level: context.level,
    useCases: [uc.sectionH3CyclePolicy, uc.argumentsScopeResolution],
  };
  calls.push(call);
};

export const detectGraphShape = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'graph.shape.detect',
    title: 'Detect graph shape (tree, DAG, or cyclic)',
    note: 'Classify selected graph as tree, DAG, or cyclic before renderer selection; if shape is cyclic and cycle-policy is `disallow`, emit error diagnostic and prevent graph rendering for that section.',
    level: context.level,
    useCases: [uc.reportGraphShapeAwareRender, uc.sectionH3CyclePolicy],
  };
  calls.push(call);
};

export const resolveRendererRegistry = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'renderer.registry.resolve',
    title: 'Resolve renderer/plugin registry',
    note: 'Load renderer capabilities, supported arguments, and shape compatibility.',
    level: context.level,
    useCases: [uc.rendererRegistry],
  };
  calls.push(call);
};

export const selectRendererPlugin = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'renderer.plugin.select',
    title: 'Select renderer plugin from arguments',
    note: 'Choose renderer by resolved typed renderer-scoped arguments with deterministic fallback when unspecified, then pass one resolved renderer argument set to the selected plugin.',
    level: context.level,
    useCases: [uc.rendererPluginSelection, uc.argumentsFreeForm],
  };
  calls.push(call);
};

export const resolveRendererScopedArguments = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'args.renderer.resolve',
    title: 'Resolve renderer-scoped arguments',
    note: 'Collect arguments from H3Section and its notes, keep only keys whose registry scope includes `renderer`, apply precedence (`note` overrides `h3-section`, `h3-section` overrides registry defaults), and produce one typed validated renderer argument set.',
    level: context.level,
    useCases: [
      uc.argumentsScopeResolution,
      uc.argumentsRuntimeValidation,
      uc.argumentsTypeCoercion,
      uc.rendererPluginSelection,
    ],
  };
  calls.push(call);
};

export const renderGraphTreeOrDag = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'render.graph.tree-or-dag',
    title: 'Render tree or DAG graph',
    note: 'Prefer hierarchical markdown text; Mermaid can be emitted as an additional diagram.',
    level: context.level,
    useCases: [
      uc.reportGraphShapeAwareRender,
      uc.reportGraphRendererMarkdownText,
      uc.reportGraphRendererMermaid,
    ],
  };
  calls.push(call);
  renderGraphAsMarkdownText(incrContext(context));
  renderGraphAsMermaid(incrContext(context));
};

export const renderGraphCircular = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'render.graph.circular',
    title: 'Render cyclic graph',
    note: 'Render only when cycle-policy is `allow`; prefer Mermaid for deterministic cycle readability, with markdown text summary fallback.',
    level: context.level,
    useCases: [
      uc.sectionH3CyclePolicy,
      uc.reportGraphShapeAwareRender,
      uc.reportGraphRendererMermaid,
      uc.reportGraphRendererMarkdownText,
    ],
  };
  calls.push(call);
  renderGraphAsMermaid(incrContext(context));
  renderGraphAsMarkdownText(incrContext(context));
};

export const renderGraphAsMarkdownText = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'render.graph.markdown.text',
    title: 'Render graph as markdown text',
    note: 'Render adjacency and hierarchy using the same markdown style as FLOW_DESIGN.',
    level: context.level,
    useCases: [uc.reportGraphRendererMarkdownText, uc.rendererRegistry],
  };
  calls.push(call);
};

export const renderGraphAsMermaid = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'render.graph.mermaid',
    title: 'Render graph as Mermaid',
    note: 'Emit Mermaid syntax for visual rendering in markdown consumers.',
    level: context.level,
    useCases: [
      uc.reportGraphRendererMermaid,
      uc.noteMermaidEmbed,
      uc.rendererRegistry,
    ],
  };
  calls.push(call);
};

export const renderSectionWithMedia = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'render.section.file.media',
    title: 'Render section with media file',
    note: 'Embed image previews for supported media types.',
    level: context.level,
    useCases: [
      uc.noteFilepathReference,
      uc.noteImagePreview,
      uc.argumentsFreeForm,
    ],
  };
  calls.push(call);
};

export const generateSingleH3Section = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'action.generate.markdown.section.h3',
    title: 'Generate a single H3 section',
    note: 'Compose subgraph, plain content, and file-backed content with section-level arguments.',
    level: context.level,
    useCases: [
      uc.reportGenerate,
      uc.reportSubgraphByLabel,
      uc.argumentsFreeForm,
      uc.configReduceNoiseWithArgs,
      uc.outputDeterministicOrdering,
    ],
  };
  calls.push(call);
  resolveH3SectionArguments(incrContext(context));
  resolveArgumentRegistry(incrContext(context));
  validateRenderArguments(incrContext(context));
  coerceRenderArguments(incrContext(context));
  selectSubGraph(incrContext(context));
  renderGraphSection(incrContext(context));
  renderPlainSection(incrContext(context));
  renderSectionWithFile(incrContext(context));
  applyDeterministicOrdering(incrContext(context));
};

export const resolveH3SectionArguments = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'args.h3.resolve',
    title: 'Resolve H3Section free-form arguments',
    note: 'Read flexible section arguments as key/value flags (for example `graph-renderer=mermaid`) and expose candidates for renderer-scoped resolution.',
    level: context.level,
    useCases: [
      uc.argumentsFreeForm,
      uc.configReduceNoiseWithArgs,
      uc.argumentsScopeResolution,
    ],
  };
  calls.push(call);
};

export const resolveNoteRenderArguments = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'args.note.resolve',
    title: 'Resolve Note free-form arguments',
    note: 'Read note-level rendering options as key/value flags (for example `format-csv=md`) and expose candidates for renderer-scoped resolution with higher precedence than H3Section values.',
    level: context.level,
    useCases: [
      uc.argumentsFreeForm,
      uc.configReduceNoiseWithArgs,
      uc.argumentsScopeResolution,
    ],
  };
  calls.push(call);
};

export const validateRenderArguments = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'args.validate.runtime',
    title: 'Validate arguments at runtime',
    note: 'Validate keys and values against a known argument registry and fail fast on invalid input.',
    level: context.level,
    useCases: [
      uc.argumentsRuntimeValidation,
      uc.argumentsRegistrySchema,
      uc.argumentsScopeResolution,
    ],
  };
  calls.push(call);
};

export const resolveArgumentRegistry = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'args.registry.resolve',
    title: 'Resolve argument registry schema',
    note: 'Load known argument definitions (type, default, allowed values, scopes) where valid scopes are `h3-section`, `note`, and `renderer`.',
    level: context.level,
    useCases: [uc.argumentsRegistrySchema],
  };
  calls.push(call);
};

export const coerceRenderArguments = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'args.coerce.typed',
    title: 'Coerce arguments to typed values',
    note: 'Coerce validated values to target types (string[], boolean, enum, number).',
    level: context.level,
    useCases: [uc.argumentsTypeCoercion, uc.argumentsRegistrySchema],
  };
  calls.push(call);
};

export const generateJsonAction = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'action.generate.json',
    title: 'Generate JSON graph export',
    note: 'Export notes and relationships in machine-readable JSON format.',
    level: context.level,
    useCases: [uc.exportJsonGraph],
  };
  calls.push(call);
  loadAppData(incrContext(context));
  validateAppData(incrContext(context));
  exportJsonGraph(incrContext(context));
};

export const validateAction = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'action.validate',
    title: 'Validate the CUE file',
    note: 'Run canonical application validation and emit the same diagnostics that gate generation.',
    level: context.level,
    useCases: [
      uc.configRelationshipsLabeled,
      uc.configReportsMultiple,
      uc.diagnosticsModel,
      uc.diagnosticsValidation,
    ],
  };
  calls.push(call);
  loadAppData(incrContext(context));
  validateAppData(incrContext(context));
  emitDiagnostics(incrContext(context));
};

export const lintNamesAction = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'action.lint.names',
    title: 'Lint note and relationship names',
    note: 'Run naming-style hygiene checks with `--style dot|snake|regex` (default dot), optional `--pattern` for regex style, optional `--prefix` scope, and configurable `--severity warning|error` (default warning).',
    level: context.level,
    useCases: [
      uc.namesLint,
      uc.namesStylePolicy,
      uc.namesPrefixFilter,
      uc.diagnosticsModel,
      uc.diagnosticsValidation,
      uc.outputDeterministicOrdering,
      uc.outputDeterministicOrderingPolicy,
    ],
  };
  calls.push(call);
  loadAppData(incrContext(context));
  validateAppData(incrContext(context));
  resolveOrderingPolicy(incrContext(context));
  applyDeterministicOrdering(incrContext(context));
  resolveNameStylePolicy(incrContext(context));
  filterNamesByPrefix(incrContext(context));
  lintNoteNames(incrContext(context));
  lintRelationshipEndpointNames(incrContext(context));
  emitDiagnostics(incrContext(context));
};

export const loadAppData = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'load.app.data',
    title: 'Load CUE application data',
    note: 'Read notes, relationships, and report definitions from config.',
    displayOnce: true,
    level: context.level,
    useCases: [
      uc.reportGenerate,
      uc.exportJsonGraph,
      uc.configRelationshipsLabeled,
      uc.configReportsMultiple,
      uc.reportList,
    ],
  };
  calls.push(call);
};

export const validateAppData = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'validate.app.data',
    title: 'Validate CUE application data',
    note: 'Canonical validation pipeline: schema checks, argument registry and free-form argument validation, dataset-based label reference validation, graph integrity policy resolution and graph integrity checks, diagnostic collection, and normalized ValidatedApp output.',
    displayOnce: true,
    level: context.level,
    useCases: [
      uc.configRelationshipsLabeled,
      uc.configReportsMultiple,
      uc.diagnosticsModel,
      uc.diagnosticsValidation,
      uc.argumentsRegistrySchema,
      uc.argumentsRuntimeValidation,
      uc.argumentsFreeForm,
      uc.graphIntegrityPolicy,
      uc.graphIntegrityValidation,
    ],
  };
  calls.push(call);
  validateCueSchema(incrContext(context));
  resolveArgumentRegistry(incrContext(context));
  validateArgumentRegistry(incrContext(context));
  validateConfigArguments(incrContext(context));
  collectDatasetLabels(incrContext(context));
  validateLabelReferences(incrContext(context));
  resolveGraphIntegrityPolicy(incrContext(context));
  validateGraphIntegrity(incrContext(context));
  collectValidationDiagnostics(incrContext(context));
  normalizeValidatedAppModel(incrContext(context));
};

export const validateCueSchema = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'validate.cue.schema',
    title: 'Validate CUE schema and structure',
    note: 'Validate required fields, types, and cross-references and attach precise config locations to diagnostics.',
    level: context.level,
    useCases: [uc.diagnosticsValidation, uc.configReportsMultiple],
  };
  calls.push(call);
};

export const validateArgumentRegistry = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'args.registry.validate',
    title: 'Validate argument registry schema consistency',
    note: 'Validate argument definitions, duplicate keys, scopes, defaults, and allowed values.',
    level: context.level,
    useCases: [uc.argumentsRegistrySchema, uc.diagnosticsValidation],
  };
  calls.push(call);
};

export const validateConfigArguments = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'args.validate.config',
    title: 'Validate configured free-form arguments',
    note: 'Validate free-form arguments declared in config against registry definitions and scope rules.',
    level: context.level,
    useCases: [
      uc.argumentsFreeForm,
      uc.argumentsRuntimeValidation,
      uc.argumentsScopeResolution,
      uc.diagnosticsValidation,
    ],
  };
  calls.push(call);
};

export const collectDatasetLabels = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'labels.dataset.collect',
    title: 'Collect dataset labels',
    note: 'Build authoritative labelSet as the union of labels from note.labels and relationship.labels without enforcing a taxonomy.',
    level: context.level,
    useCases: [uc.configRelationshipsLabeled, uc.diagnosticsValidation],
  };
  calls.push(call);
};

export const validateLabelReferences = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'labels.reference.validate',
    title: 'Validate label references',
    note: 'Validate referenced labels used by config elements (for example graph.select label arguments) against labelSet; emit `LABEL_REF_UNKNOWN` (default severity `warning`) with argument location and referenced label value for unknown references.',
    level: context.level,
    useCases: [
      uc.reportSubgraphByLabel,
      uc.configRelationshipsLabeled,
      uc.argumentsScopeResolution,
      uc.diagnosticsValidation,
    ],
  };
  calls.push(call);
};

export const collectValidationDiagnostics = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'diagnostics.collect.validation',
    title: 'Collect validation diagnostics',
    note: 'Collect stable diagnostic codes, severities, sources, canonical machine-readable config `location` paths, and human-readable context fields (`reportTitle`, `sectionTitle`, `noteName`, `relationship`, `argumentName`).',
    level: context.level,
    useCases: [uc.diagnosticsModel, uc.diagnosticsValidation],
  };
  calls.push(call);
};

export const normalizeValidatedAppModel = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'app.model.normalize',
    title: 'Normalize validated application model',
    note: 'Build ValidatedApp with normalized notes, relationships, reports, resolved graph integrity policy, resolved argument registry, and diagnostics. Ordering policy resolution remains generation-time.',
    level: context.level,
    useCases: [
      uc.configReportsMultiple,
      uc.configRelationshipsLabeled,
      uc.graphIntegrityPolicy,
      uc.argumentsRegistrySchema,
      uc.diagnosticsModel,
    ],
  };
  calls.push(call);
};

export const resolveGraphIntegrityPolicy = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'graph.integrity.policy.resolve',
    title: 'Resolve graph integrity policy',
    note: 'Resolve integrity policy for missing nodes, orphans, duplicates, unknown label references, and cross-report references.',
    level: context.level,
    useCases: [uc.graphIntegrityPolicy],
  };
  calls.push(call);
};

export const validateGraphIntegrity = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'graph.integrity.validate',
    title: 'Validate graph integrity',
    note: 'Run integrity checks and emit diagnostics according to resolved policy.',
    level: context.level,
    useCases: [uc.graphIntegrityValidation, uc.diagnosticsValidation],
  };
  calls.push(call);
  checkMissingRelationshipNodes(incrContext(context));
  checkOrphanNodes(incrContext(context));
  checkDuplicateNoteNames(incrContext(context));
  checkCrossReportReferences(incrContext(context));
};

export const checkMissingRelationshipNodes = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'graph.integrity.check.missing-nodes',
    title: 'Check missing relationship nodes',
    note: 'Detect relationships that reference notes that do not exist.',
    level: context.level,
    useCases: [uc.graphIntegrityValidation],
  };
  calls.push(call);
};

export const checkOrphanNodes = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'graph.integrity.check.orphans',
    title: 'Check orphan nodes',
    note: 'Detect notes disconnected from report roots/sections.',
    level: context.level,
    useCases: [uc.graphIntegrityValidation],
  };
  calls.push(call);
};

export const checkDuplicateNoteNames = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'graph.integrity.check.duplicate-note-names',
    title: 'Check duplicate note names',
    note: 'Detect duplicate note identifiers that can cause ambiguous references.',
    level: context.level,
    useCases: [uc.graphIntegrityValidation],
  };
  calls.push(call);
};

export const checkCrossReportReferences = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'graph.integrity.check.cross-report-references',
    title: 'Check cross-report references',
    note: 'Validate whether note/edge references across report boundaries are allowed by policy.',
    level: context.level,
    useCases: [uc.graphIntegrityValidation, uc.graphIntegrityPolicy],
  };
  calls.push(call);
};

export const resolveOrderingPolicy = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'ordering.policy.resolve',
    title: 'Resolve deterministic ordering policy',
    note: 'Resolve explicit comparators: notes by (primaryLabel, name) where primaryLabel is the lexicographically smallest label; relationships by (from, to, labelsSortedJoined) where labelsSortedJoined is labels sorted lexicographically then joined with `|`; sections by (lowercase(title), originalIndex) for stable tie-breaks; arguments by argument name.',
    displayOnce: true,
    level: context.level,
    useCases: [uc.outputDeterministicOrderingPolicy],
  };
  calls.push(call);
};

export const applyDeterministicOrdering = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'ordering.apply.deterministic',
    title: 'Apply deterministic ordering',
    note: 'Apply resolved comparators exactly and use stable tie-breakers only (including section originalIndex), yielding reproducible output without runtime randomness.',
    displayOnce: true,
    level: context.level,
    useCases: [
      uc.outputDeterministicOrdering,
      uc.outputDeterministicOrderingPolicy,
    ],
  };
  calls.push(call);
};

export const filterNamesByPrefix = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'names.filter.prefix',
    title: 'Filter names by prefix',
    note: 'Apply required `--prefix` filter: keep notes where `name` starts with prefix and relationships where `from` or `to` starts with prefix.',
    level: context.level,
    useCases: [uc.namesPrefixFilter, uc.namesList, uc.namesLint],
  };
  calls.push(call);
};

export const filterNamesByKind = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'names.filter.kind',
    title: 'Filter names by kind',
    note: 'Apply optional `--kind notes|relationships|all` filter (default `all`) to reduce output noise.',
    level: context.level,
    useCases: [uc.namesList, uc.namesOutputFormats],
  };
  calls.push(call);
};

export const outputNamesTable = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'names.output.table',
    title: 'Output names as table',
    note: 'Default output: notes table rows `name | title | labels` and relationship rows `from | to | labels`.',
    level: context.level,
    useCases: [uc.namesList, uc.namesOutputFormats],
  };
  calls.push(call);
};

export const outputNamesJson = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'names.output.json',
    title: 'Output names as JSON',
    note: 'Optional `--format json` output as `{ notes: [], relationships: [] }` with the same fields used in table mode.',
    level: context.level,
    useCases: [uc.namesList, uc.namesOutputFormats],
  };
  calls.push(call);
};

export const resolveNameStylePolicy = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'lint.names.policy.resolve',
    title: 'Resolve name style policy',
    note: 'Resolve style matcher as case-sensitive policy: `dot`=`^[a-z][a-z0-9]*(\\.[a-z][a-z0-9]*)*$`, `snake`=`^[a-z][a-z0-9_]*$`, `regex`=user-provided `--pattern`.',
    level: context.level,
    useCases: [uc.namesLint, uc.namesStylePolicy],
  };
  calls.push(call);
};

export const lintNoteNames = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'lint.names.notes',
    title: 'Lint note names',
    note: 'Check note names against resolved style; emit `NAME_STYLE_VIOLATION` diagnostics with canonical config location and human-readable context for each violation.',
    level: context.level,
    useCases: [uc.namesLint, uc.namesStylePolicy, uc.diagnosticsValidation],
  };
  calls.push(call);
};

export const lintRelationshipEndpointNames = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'lint.names.relationships',
    title: 'Lint relationship endpoint names',
    note: 'Check relationship `from` and `to` endpoint names against resolved style; emit `NAME_STYLE_VIOLATION` diagnostics with canonical config location and relationship context for each violation.',
    level: context.level,
    useCases: [uc.namesLint, uc.namesStylePolicy, uc.diagnosticsValidation],
  };
  calls.push(call);
};

export const emitDiagnostics = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'diagnostics.emit.structured',
    title: 'Emit structured diagnostics',
    note: 'Emit diagnostics with code, severity, source, message, canonical machine-readable `location`, and optional human-readable context fields.',
    level: context.level,
    useCases: [uc.diagnosticsModel, uc.diagnosticsValidation],
  };
  calls.push(call);
};

export const listValidatedReports = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'list.reports.output',
    title: 'List reports from ValidatedApp',
    note: 'Enumerate reports from the normalized validated model with optional strictness behavior handled by validation policy.',
    level: context.level,
    useCases: [uc.reportList, uc.configReportsMultiple],
  };
  calls.push(call);
};

export const exportJsonGraph = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'export.graph.json',
    title: 'Export validated graph as JSON',
    note: 'Export notes and relationships from ValidatedApp without re-running validation steps.',
    level: context.level,
    useCases: [uc.exportJsonGraph],
  };
  calls.push(call);
};
