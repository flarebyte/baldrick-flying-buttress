# Risks Overview

Migrated risk catalog and mitigations.

## Risks

Risk notes.

### Catalog

All risks.

#### Limited practical AI assistance for CUE authoring

Description: Generative AI tools often have weaker CUE support than for mainstream formats, which can reduce productivity and increase hand-written config errors.

Mitigation: Provide templates, examples, and lintable conventions in-repo; rely on strong validation and developer documentation rather than AI-generated CUE.

#### Higher merge-conflict risk in shared CUE files

Description: When multiple developers edit the same large CUE file, concurrent changes can frequently overlap and create conflict-heavy pull requests.

Mitigation: Reduce shared hotspots with file partitioning, stable key ordering, and ownership boundaries; add CI validation to catch conflicts and schema drift early.

#### Circular dependencies in note relationships

Description: Tree or DAG-like relationship graphs are usually straightforward, but circular dependencies can break assumptions in traversal, filtering, and report assembly.

Mitigation: Add explicit cycle detection and policy controls (`disallow` to emit error and skip section graph rendering, `allow` to render cyclic graphs), and test traversal logic with cyclic graph fixtures.

#### Single large CUE file becomes hard to maintain

Description: Packing too many notes and relationships into one CUE file increases cognitive load, makes reviews difficult, and raises the chance of accidental breakage.

Mitigation: Split configuration into modular CUE packages/files by domain or report, then compose them through imports and shared schema constraints.

#### Report generation performance at scale

Description: Generating many sections from large note graphs and file-backed content can increase CPU, memory, and I/O cost, causing slow report builds and degraded CLI responsiveness.

Mitigation: Use bounded concurrency, lazy file loading, and optional caching of parsed CUE and graph selections; add profiling baselines and fail-fast limits for oversized runs.

