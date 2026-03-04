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
  generateMarkdownAction(incrContext(context));
  generateJsonAction(incrContext(context));
  validateAction(incrContext(context));
};

export const listReportsAction = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'action.list.reports',
    title: 'List configured markdown reports',
    note: 'Enumerate report targets without generating files.',
    level: context.level,
    useCases: [uc.reportList, uc.configReportsMultiple],
  };
  calls.push(call);
  loadAppData(incrContext(context));
  validateAppData(incrContext(context));
};

export const generateMarkdownAction = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'action.generate.markdown',
    title: 'Generate markdown reports',
    note: 'Renders one or more markdown outputs from the validated config.',
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
    note: 'Build H3 sections from note subsets and renderers.',
    level: context.level,
    useCases: [uc.reportGenerate, uc.noteBasicMarkdown],
  };
  calls.push(call);
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
    note: 'Resolve cycle policy and graph shape, then render with selected renderer(s).',
    level: context.level,
    useCases: [
      uc.configRelationshipsLabeled,
      uc.reportSubgraphByLabel,
      uc.sectionH3CyclePolicy,
      uc.reportGraphShapeAwareRender,
      uc.rendererRegistry,
      uc.rendererPluginSelection,
    ],
  };
  calls.push(call);
  resolveRendererRegistry(incrContext(context));
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
    useCases: [uc.noteFilepathReference, uc.argumentsFreeForm],
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
    note: 'Render as a markdown table or raw CSV code block (for example `format-csv=md`).',
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
    note: 'Apply include/exclude filters before rendering CSV output.',
    level: context.level,
    useCases: [uc.noteCsvFilterColumn, uc.argumentsFreeForm],
  };
  calls.push(call);
};

export const selectSubGraph = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'graph.select',
    title: 'Extract subgraph using labels',
    note: 'Filter notes and relationships by labels and optional starting node.',
    level: context.level,
    useCases: [uc.configRelationshipsLabeled, uc.reportSubgraphByLabel],
  };
  calls.push(call);
};

export const resolveH3GraphCyclePolicy = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'graph.policy.cycle',
    title: 'Resolve H3Section cycle policy argument',
    note: 'Use section argument to disallow, allow, or collapse cycles.',
    level: context.level,
    useCases: [uc.sectionH3CyclePolicy, uc.argumentsScopeResolution],
  };
  calls.push(call);
};

export const detectGraphShape = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'graph.shape.detect',
    title: 'Detect graph shape (tree, DAG, or cyclic)',
    note: 'Classify graph structure before selecting rendering strategy.',
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
    note: 'Choose renderer by arguments with deterministic fallback when unspecified.',
    level: context.level,
    useCases: [uc.rendererPluginSelection, uc.argumentsFreeForm],
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
    note: 'Prefer Mermaid for cycle readability, with markdown text summary as fallback.',
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
};

export const resolveH3SectionArguments = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'args.h3.resolve',
    title: 'Resolve H3Section free-form arguments',
    note: 'Read flexible section arguments as key/value flags (for example `graph-renderer=mermaid`).',
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
    note: 'Read note-level rendering options as key/value flags (for example `format-csv=md`).',
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
    note: 'Load known argument definitions (type, default, allowed values, scopes).',
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
};

export const validateAction = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'action.validate',
    title: 'Validate the CUE file',
    note: 'Validate configuration structure and constraints without rendering output.',
    level: context.level,
    useCases: [uc.configRelationshipsLabeled, uc.configReportsMultiple],
  };
  calls.push(call);
  loadAppData(incrContext(context));
  validateAppData(incrContext(context));
};

export const loadAppData = (context: FlowContext) => {
  const call: ComponentCall = {
    name: 'load.app.data',
    title: 'Load CUE application data',
    note: 'Read notes, relationships, and report definitions from config.',
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
    note: 'Ensure required fields and cross-reference integrity are valid.',
    level: context.level,
    useCases: [uc.configRelationshipsLabeled, uc.configReportsMultiple],
  };
  calls.push(call);
};
