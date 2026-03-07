# Flow Design Overview

Migrated flow call tree and validation pipeline design notes.

## Flow

Flow call notes.

### Flow Calls

Generated from design/flows.ts.

#### Generate JSON graph export

Export notes and relationships in machine-readable JSON format.

#### Generate markdown reports

Renders one or more markdown outputs from a single validated application model.

#### Generate a single H3 section

Compose subgraph, plain content, and file-backed content with section-level arguments.

#### Generate markdown sections

Build H3 sections from note subsets and renderers with deterministic ordering.

#### Lint note and relationship names

Run naming-style hygiene checks with `--style dot|snake|regex` (default dot), optional `--pattern` for regex style, optional `--prefix` scope, and configurable `--severity warning|error` (default warning).

#### Lint contextual orphans

Run orphan-query lint checks with required `--subject-label`, optional `--edge-label`, optional `--counterpart-label`, optional `--direction in|out|either` (default `either`), and configurable `--severity warning|error` (default warning).

#### List note and relationship names

Print note and relationship identifiers for daily inventory with required `--prefix`, optional `--kind notes|relationships|all`, and `--format table|json` (default table).

#### List configured markdown reports

Enumerate report targets from the validated application model without generating files.

#### Validate the CUE file

Run canonical application validation and emit the same diagnostics that gate generation.

#### Normalize validated application model

Build ValidatedApp with normalized notes, relationships, reports, resolved graph integrity policy, resolved argument registry, and diagnostics. Ordering policy resolution remains generation-time.

#### Coerce arguments to typed values

Coerce validated values to target types (string[], boolean, enum, number).

#### Resolve H3Section free-form arguments

Read flexible section arguments as key/value flags (for example `graph-renderer=mermaid`) and expose candidates for renderer-scoped resolution.

#### Resolve Note free-form arguments

Read note-level rendering options as key/value flags (for example `format-csv=md`) and expose candidates for renderer-scoped resolution with higher precedence than H3Section values.

#### Resolve orphan query arguments from H3 section

Resolve `orphan-subject-label` (required for orphan mode), optional `orphan-edge-label`, optional `orphan-counterpart-label`, and `orphan-direction in|out|either` (default `either`).

#### Resolve argument registry schema

Load known argument definitions (type, default, allowed values, scopes) where valid scopes are `h3-section`, `note`, and `renderer`.

#### Validate argument registry schema consistency

Validate argument definitions, duplicate keys, scopes, defaults, and allowed values.

#### Resolve renderer-scoped arguments

Collect arguments from H3Section and its notes, keep only keys whose registry scope includes `renderer`, apply precedence (`note` overrides `h3-section`, `h3-section` overrides registry defaults), and produce one typed validated renderer argument set.

#### Validate configured free-form arguments

Validate free-form arguments declared in config against registry definitions and scope rules.

#### Validate arguments at runtime

Validate keys and values against a known argument registry and fail fast on invalid input.

#### flyb CLI root command

Entry point for report generation, listing, JSON export, and config validation.

#### Collect validation diagnostics

Collect stable diagnostic codes, severities, sources, canonical machine-readable config `location` paths, and human-readable context fields (`reportTitle`, `sectionTitle`, `noteName`, `relationship`, `argumentName`, `subjectLabel`, `edgeLabel`, `counterpartLabel`).

#### Emit structured diagnostics

Emit diagnostics with code, severity, source, message, canonical machine-readable `location`, and optional human-readable context fields.

#### Export validated graph as JSON

Export notes and relationships from ValidatedApp without re-running validation steps.

#### Filter CSV rows by column

Apply exact-match include/exclude filters before rendering CSV output: `csv-include=column:value` keeps matching rows, `csv-exclude=column:value` removes matching rows, and multiple filters are allowed.

#### Check cross-report references

Validate whether note/edge references across report boundaries are allowed by policy.

#### Check duplicate note names

Detect duplicate note identifiers that can cause ambiguous references.

#### Check missing relationship nodes

Detect relationships that reference notes that do not exist.

#### Check orphan nodes

Detect notes disconnected from report roots/sections.

#### Resolve graph integrity policy

Resolve integrity policy for missing nodes, orphans, duplicates, unknown label references, and cross-report references.

#### Validate graph integrity

Run integrity checks and emit diagnostics according to resolved policy.

#### Resolve H3Section cycle policy argument

Resolve section cycle policy (`disallow` or `allow`): `disallow` requires cycle detection error diagnostics and blocks section graph rendering; `allow` permits cyclic rendering.

#### Extract subgraph using labels

Filter notes and relationships by labels and optional starting node; label references are pre-validated against dataset labels (union of note.labels and relationship.labels).

#### Detect graph shape (tree, DAG, or cyclic)

Classify selected graph as tree, DAG, or cyclic before renderer selection; if shape is cyclic and cycle-policy is `disallow`, emit error diagnostic and prevent graph rendering for that section.

#### Collect dataset labels

Build authoritative labelSet as the union of labels from note.labels and relationship.labels without enforcing a taxonomy.

#### Validate label references

Validate referenced labels used by config elements (for example graph.select and orphan-query label arguments) against labelSet; emit `LABEL_REF_UNKNOWN` (default severity `warning`) with argument location and referenced label value for unknown references.

#### Lint note names

Check note names against resolved style; emit `NAME_STYLE_VIOLATION` diagnostics with canonical config location and human-readable context for each violation.

#### Resolve name style policy

Resolve style matcher as case-sensitive policy: `dot`=`^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$`, `snake`=`^[a-z][a-z0-9_]*$`, `regex`=user-provided `--pattern`.

#### Lint relationship endpoint names

Check relationship `from` and `to` endpoint names against resolved style; emit `NAME_STYLE_VIOLATION` diagnostics with canonical config location and relationship context for each violation.

#### Emit contextual orphan diagnostics

Emit `ORPHAN_QUERY_MISSING_LINK` diagnostics for each orphan note with canonical config/CLI context location and human-readable fields (`noteName`, `subjectLabel`, `edgeLabel`, `counterpartLabel`).

#### Resolve orphan query from CLI flags

Resolve lint flags into orphan query context: `--subject-label` required, `--edge-label` optional, `--counterpart-label` optional, `--direction` default `either`.

#### List reports from ValidatedApp

Enumerate reports from the normalized validated model with optional strictness behavior handled by validation policy.

#### Load CUE application data

Read notes, relationships, and report definitions from config.

#### Filter names by kind

Apply optional `--kind notes|relationships|all` filter (default `all`) to reduce output noise.

#### Filter names by prefix

Apply required `--prefix` filter: keep notes where `name` starts with prefix and relationships where `from` or `to` starts with prefix.

#### Output names as JSON

Optional `--format json` output as `{ notes: [], relationships: [] }` with the same fields used in table mode.

#### Output names as table

Default output: notes table rows `name | title | labels` and relationship rows `from | to | labels`.

#### Apply deterministic ordering

Apply resolved comparators exactly and use stable tie-breakers only (including section originalIndex), yielding reproducible output without runtime randomness.

#### Resolve deterministic ordering policy

Resolve explicit comparators: notes by (primaryLabel, name) where primaryLabel is the lexicographically smallest label; relationships by (from, to, labelsSortedJoined) where labelsSortedJoined is labels sorted lexicographically then joined with `|`; sections by (lowercase(title), originalIndex) for stable tie-breaks; arguments by argument name.

#### Find contextual orphans

For each subject note, require at least one matching relationship under query filters (edge label, counterpart label, direction). Notes with zero matches are contextual orphans.

#### Render orphan rows

Render deterministic orphan output rows/table with `name`, `title`, and `labels`.

#### Render cyclic graph

Render only when cycle-policy is `allow`; markdown traversal expands each node once and when revisiting a node emits `*(cycle back to <node>)*` linking to first anchor, then appends a short deterministic adjacency summary (`A -> B (labels)`). Mermaid remains preferred for cycle readability.

#### Render graph as markdown text

Render tree/DAG/cyclic graphs in plain markdown with deterministic traversal order, stable note anchors derived from note names, and reference links to first occurrence anchors for repeated nodes or cycle backs.

#### Render graph as Mermaid

Emit Mermaid syntax for visual rendering in markdown consumers.

#### Render tree or DAG graph

Tree: render full hierarchy as nested markdown lists (`**name** — title` plus optional short description). DAG: use stable DFS by ordering policy, expand first encounter, and on repeated encounters allow repetition only when children<=3 and depth<=2; otherwise emit `*(see above)*` reference linking to first anchor.

#### Render section with referenced file content

Dispatches file rendering by type (CSV, media, code/diagram).

#### Render section with code or Mermaid snippet

Preserve fenced-block formatting for code and Mermaid content.

#### Render section with CSV file

Render as a markdown table or raw CSV code block (for example `format-csv=md`) and apply note-scoped CSV filters (`csv-include` / `csv-exclude`) using `column:value` exact-match rules.

#### Render section with media file

Embed image previews for supported media types.

#### Render section as a graph

Resolve cycle policy and graph shape, then render with selected renderer(s); renderer/runtime diagnostics here must not duplicate graph-integrity diagnostics.

#### Render section as contextual orphans

Resolve orphan query arguments and emit deterministic orphan list/table rows (`name`, `title`, `labels`) for subject notes missing required contextual links.

#### Render plain section

Render title and markdown body, including markdown links and note-level argument options.

#### Select renderer plugin from arguments

Choose renderer by resolved typed renderer-scoped arguments with deterministic fallback when unspecified, then pass one resolved renderer argument set to the selected plugin.

#### Resolve renderer/plugin registry

Load renderer capabilities, supported arguments, and shape compatibility.

#### Validate CUE application data

Canonical validation pipeline: schema checks, argument registry and free-form argument validation, dataset-based label reference validation, graph integrity policy resolution and graph integrity checks, diagnostic collection, and normalized ValidatedApp output.

#### Validate CUE schema and structure

Validate required fields, types, and cross-references and attach precise config locations to diagnostics.

