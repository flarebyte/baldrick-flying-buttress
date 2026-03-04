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
      'Use Cobra to model commands, flags, and subcommands (`generate markdown`, `generate json`, `validate`, `list`) with a consistent command tree.',
    calls: [
      'cli.root',
      'action.generate.markdown',
      'action.generate.json',
      'action.validate',
      'action.list.reports',
    ],
  },
  'renderer.registry.contract': {
    name: 'renderer.registry.contract',
    title: 'Define a renderer plugin registry contract',
    description:
      'Define a small renderer interface (name, supportsGraphShape, supportedArguments, render) and register built-ins (markdown-text, mermaid) in a deterministic lookup map.',
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
      'Resolve renderer from arguments first, then apply stable defaults by graph shape (for example Mermaid-first for cycles, markdown-first for tree/DAG).',
    calls: [
      'renderer.plugin.select',
      'graph.shape.detect',
      'render.graph.tree-or-dag',
      'render.graph.circular',
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
      'Maintain a registry of argument definitions (name, type, default, allowed values, scopes) and use it as the single source of truth for argument behavior.',
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
      'Resolve and validate arguments by scope (global, h2, h3, note, renderer) so options are accepted only where they are meaningful.',
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
