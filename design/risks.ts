import type { Risk } from './common.ts';

// Central risk catalogue. Extend as the design evolves.
export const risks: Record<string, Risk> = {
  'performance.report-generation-scale': {
    name: 'performance.report-generation-scale',
    title: 'Report generation performance at scale',
    description:
      'Generating many sections from large note graphs and file-backed content can increase CPU, memory, and I/O cost, causing slow report builds and degraded CLI responsiveness.',
    mitigation:
      'Use bounded concurrency, lazy file loading, and optional caching of parsed CUE and graph selections; add profiling baselines and fail-fast limits for oversized runs.',
    calls: [
      'action.generate.markdown',
      'action.generate.markdown.sections',
      'action.generate.markdown.section.h3',
      'render.section.file',
      'render.section.file.csv',
      'graph.select',
    ],
  },
  'maintenance.single-cue-file-size': {
    name: 'maintenance.single-cue-file-size',
    title: 'Single large CUE file becomes hard to maintain',
    description:
      'Packing too many notes and relationships into one CUE file increases cognitive load, makes reviews difficult, and raises the chance of accidental breakage.',
    mitigation:
      'Split configuration into modular CUE packages/files by domain or report, then compose them through imports and shared schema constraints.',
    calls: ['load.app.data', 'validate.app.data', 'action.validate'],
  },
  'collaboration.cue-merge-conflicts': {
    name: 'collaboration.cue-merge-conflicts',
    title: 'Higher merge-conflict risk in shared CUE files',
    description:
      'When multiple developers edit the same large CUE file, concurrent changes can frequently overlap and create conflict-heavy pull requests.',
    mitigation:
      'Reduce shared hotspots with file partitioning, stable key ordering, and ownership boundaries; add CI validation to catch conflicts and schema drift early.',
    calls: ['load.app.data', 'validate.app.data', 'action.validate'],
  },
  'graph.circular-dependency': {
    name: 'graph.circular-dependency',
    title: 'Circular dependencies in note relationships',
    description:
      'Tree or DAG-like relationship graphs are usually straightforward, but circular dependencies can break assumptions in traversal, filtering, and report assembly.',
    mitigation:
      'Add explicit cycle detection and policy controls (`disallow` to emit error and skip section graph rendering, `allow` to render cyclic graphs), and test traversal logic with cyclic graph fixtures.',
    calls: [
      'graph.select',
      'render.section.graph',
      'action.generate.markdown.section.h3',
      'action.validate',
    ],
  },
  'authoring.cue-ai-assistance-gap': {
    name: 'authoring.cue-ai-assistance-gap',
    title: 'Limited practical AI assistance for CUE authoring',
    description:
      'Generative AI tools often have weaker CUE support than for mainstream formats, which can reduce productivity and increase hand-written config errors.',
    mitigation:
      'Provide templates, examples, and lintable conventions in-repo; rely on strong validation and developer documentation rather than AI-generated CUE.',
    calls: ['load.app.data', 'validate.app.data', 'action.validate'],
  },
};
