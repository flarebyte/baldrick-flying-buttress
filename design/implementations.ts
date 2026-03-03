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
  'cli.arguments.typed-models': {
    name: 'cli.arguments.typed-models',
    title: 'Model Note and H3Section arguments with typed structs',
    description:
      'Define small typed argument models for command inputs and section/note filters so argument handling remains explicit and testable.',
    calls: [
      'action.generate.markdown.sections',
      'action.generate.markdown.section.h3',
      'render.section.plain',
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
