import type { UseCase } from './common.ts';

// Use cases for parsing a single source file (Go, Dart, TypeScript).
export const useCases: Record<string, UseCase> = {
  'cli.report.generate': {
    name: 'cli.report.generate',
    title: 'Generate design reports from configured notes and relationships',
    note: 'This is the primary end-to-end report generation use case.',
  },
  'cli.note.basic-markdown': {
    name: 'cli.note.basic-markdown',
    title: 'Render note title and markdown description',
    note: 'Each note includes a concise title with free-form markdown content.',
  },
  'cli.note.filepath.reference': {
    name: 'cli.note.filepath.reference',
    title: 'Reference a file from a note',
    note: 'Referenced files can be embedded in generated markdown output.',
  },
  'cli.note.csv.embed': {
    name: 'cli.note.csv.embed',
    title: 'Embed CSV content from a referenced file',
    note: 'CSV input can render as a markdown table or as raw CSV.',
  },
  'cli.note.csv.filter-column': {
    name: 'cli.note.csv.filter-column',
    title: 'Filter embedded CSV rows by column',
    note: 'Column filters reduce CSV output to the relevant subset.',
  },
  'cli.note.image.preview': {
    name: 'cli.note.image.preview',
    title: 'Preview referenced image files in markdown',
    note: 'Image references render as embedded previews in reports.',
  },
  'cli.note.mermaid.embed': {
    name: 'cli.note.mermaid.embed',
    title: 'Embed Mermaid diagrams from file content',
    note: 'Mermaid content is emitted in fenced blocks for diagram rendering.',
  },
  'cli.note.link.markdown': {
    name: 'cli.note.link.markdown',
    title: 'Convert note links to markdown links',
    note: 'URL links are rendered with link text in markdown output.',
  },
  'cli.config.relationships.labeled': {
    name: 'cli.config.relationships.labeled',
    title: 'Define labeled relationships between notes in config',
    note: 'CUE can be used as the source format for flexible configuration.',
  },
  'cli.config.reports.multiple': {
    name: 'cli.config.reports.multiple',
    title: 'Declare multiple markdown reports in one config',
    note: 'A single config can drive generation of multiple report files.',
  },
  'cli.report.list': {
    name: 'cli.report.list',
    title: 'List all configured markdown reports',
    note: 'The CLI exposes a command to enumerate report targets.',
  },
  'cli.export.json.graph': {
    name: 'cli.export.json.graph',
    title: 'Export notes and relationships as JSON',
    note: 'JSON export supports machine-readable graph processing.',
  },
  'cli.report.subgraph.by-label': {
    name: 'cli.report.subgraph.by-label',
    title: 'Build a report from a relationship-label subgraph',
    note: 'Report generation can include only edges matching selected labels.',
  },
  'cli.section.h3.cycle-policy': {
    name: 'cli.section.h3.cycle-policy',
    title: 'Allow each H3 section to define cycle policy',
    note: 'H3Section arguments can declare whether cycles are disallowed, allowed, or collapsed.',
  },
  'cli.report.graph.shape-aware-render': {
    name: 'cli.report.graph.shape-aware-render',
    title: 'Render graph output based on graph shape',
    note: 'Renderer behavior adapts to tree, DAG, and cyclic graph structures.',
  },
  'cli.report.graph.renderer.markdown-text': {
    name: 'cli.report.graph.renderer.markdown-text',
    title: 'Render graph output as markdown text',
    note: 'Text rendering supports readable hierarchy and edge summaries in markdown reports.',
  },
  'cli.report.graph.renderer.mermaid': {
    name: 'cli.report.graph.renderer.mermaid',
    title: 'Render graph output as Mermaid diagram',
    note: 'Mermaid output supports visual graph rendering, including cyclic relationships.',
  },
  'cli.renderer.registry': {
    name: 'cli.renderer.registry',
    title: 'Register renderers and plugins in a capability registry',
    note: 'A renderer registry maps renderer names to capabilities, supported arguments, and graph-shape compatibility.',
  },
  'cli.renderer.plugin-selection': {
    name: 'cli.renderer.plugin-selection',
    title: 'Select renderer plugin from arguments at runtime',
    note: 'Renderer selection is resolved from section and note arguments with fallback defaults.',
  },
  'cli.output.deterministic-ordering': {
    name: 'cli.output.deterministic-ordering',
    title: 'Guarantee deterministic output ordering',
    note: 'Sort notes, relationships, sections, and arguments with stable rules so repeated runs produce identical output.',
  },
  'cli.output.deterministic-ordering.policy': {
    name: 'cli.output.deterministic-ordering.policy',
    title: 'Define an explicit ordering policy',
    note: 'Ordering policy is part of runtime behavior and can be documented/tested as a contract.',
  },
  'cli.diagnostics.model': {
    name: 'cli.diagnostics.model',
    title: 'Emit structured diagnostics',
    note: 'Diagnostics include code, severity, message, source, and optional location context.',
  },
  'cli.diagnostics.validation': {
    name: 'cli.diagnostics.validation',
    title: 'Report validation diagnostics with locations',
    note: 'Validation errors and warnings should point to config paths and offending argument or relationship names.',
  },
  'cli.graph.integrity.policy': {
    name: 'cli.graph.integrity.policy',
    title: 'Define graph integrity policy beyond cycles',
    note: 'Policy covers missing nodes, orphan nodes, duplicate note names, unknown relationship labels, and cross-report references.',
  },
  'cli.graph.integrity.validation': {
    name: 'cli.graph.integrity.validation',
    title: 'Validate graph integrity using policy rules',
    note: 'Integrity checks should emit structured diagnostics tied to offending notes, relationships, and config locations.',
  },
  'cli.arguments.free-form': {
    name: 'cli.arguments.free-form',
    title: 'Accept free-form arguments on H3Section and Note',
    note: 'Arguments behave like CLI flags (for example `format-csv=md`) and can carry string, string[], boolean, and similar values.',
  },
  'cli.arguments.runtime-validation': {
    name: 'cli.arguments.runtime-validation',
    title: 'Validate free-form arguments at runtime',
    note: 'Validate against a known argument registry and fail with clear errors on unknown keys or invalid values.',
  },
  'cli.arguments.registry.schema': {
    name: 'cli.arguments.registry.schema',
    title: 'Define an argument registry schema',
    note: 'Registry entries define argument key, type, default, allowed values, and valid scopes (`h3-section`, `note`, `renderer`).',
  },
  'cli.arguments.registry.scope-resolution': {
    name: 'cli.arguments.registry.scope-resolution',
    title: 'Resolve arguments by scope',
    note: 'Apply argument rules by scope (h3-section, note, renderer) to prevent invalid combinations.',
  },
  'cli.arguments.type-coercion': {
    name: 'cli.arguments.type-coercion',
    title: 'Coerce free-form argument values into typed values',
    note: 'Convert string-like argument inputs into validated typed values before rendering.',
  },
  'cli.config.reduce-noise.with-args': {
    name: 'cli.config.reduce-noise.with-args',
    title: 'Keep CUE config compact with argument-driven rendering options',
    note: 'Prefer small composable argument lists over proliferating specialized configuration fields.',
  },
};

export const getByName = (expectedName: string) =>
  Object.values(useCases).find(({ name }) => name === expectedName);

export const mustUseCases = new Set([
  ...Object.values(useCases).map(({ name }) => name),
]);

export const useCaseCatalogByName: Record<
  string,
  { name: string; title: string; note?: string }
> = Object.fromEntries(Object.values(useCases).map((u) => [u.name, u]));
