# Use Cases

Migrated use-case catalog.

## Use Cases

Catalog from design/use_cases.ts.

### Catalog

All use cases.

#### Accept free-form arguments on H3Section and Note

Arguments behave like CLI flags (for example `format-csv=md`) and can carry string, string[], boolean, and similar values.

#### Define an argument registry schema

Registry entries define argument key, type, default, allowed values, and valid scopes (`h3-section`, `note`, `renderer`).

#### Resolve arguments by scope

Apply argument rules by scope (h3-section, note, renderer); for renderer scope, collect from H3Section and note arguments and apply precedence (`note` > `h3-section` > registry default).

#### Validate free-form arguments at runtime

Validate against a known argument registry and fail with clear errors on unknown keys or invalid values.

#### Coerce free-form argument values into typed values

Convert string-like argument inputs into validated typed values before rendering.

#### Keep CUE config compact with argument-driven rendering options

Prefer small composable argument lists over proliferating specialized configuration fields.

#### Define labeled relationships between notes in config

CUE can be used as the source format for flexible configuration; labels on notes and relationships remain free-form.

#### Declare multiple markdown reports in one config

A single config can drive generation of multiple report files.

#### Emit structured diagnostics

Diagnostics include code, severity, message, source, canonical machine-readable `location`, and additional human-readable context fields.

#### Report validation diagnostics with locations

Validation errors and warnings should include canonical index-based config paths plus readable identifiers (report title, section title, note/relationship/argument names).

#### Export notes and relationships as JSON

JSON export supports machine-readable graph processing.

#### Define graph integrity policy beyond cycles

Policy covers missing nodes, orphan nodes, duplicate note names, unknown referenced labels, and cross-report references.

#### Validate graph integrity using policy rules

Integrity checks should emit structured diagnostics tied to offending notes, relationships, and config locations.

#### Lint note and relationship names for style hygiene

The CLI exposes `flyb lint names` to emit structured diagnostics for naming-style violations without introducing label taxonomy requirements.

#### List note and relationship names for daily inventory

The CLI exposes `flyb list names` with `--prefix` filtering and `--format table|json` output.

#### Render names as table or JSON

Default output is human-friendly table; JSON is opt-in and returns `{ notes: [], relationships: [] }`.

#### Filter names by prefix scope

Prefix filtering keeps notes whose name starts with prefix and relationships where `from` or `to` starts with prefix.

#### Define explicit name style policy

Name styles are case-sensitive: dot=`^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$`, snake=`^[a-z][a-z0-9_]*$`, regex=user-provided `--pattern`.

#### Render note title and markdown description

Each note includes a concise title with free-form markdown content.

#### Embed CSV content from a referenced file

CSV input can render as a markdown table or as raw CSV.

#### Filter embedded CSV rows by column

Column filters reduce CSV output to the relevant subset using `csv-include=column:value` and `csv-exclude=column:value` exact-match arguments.

#### Reference a file from a note

Referenced files can be embedded in generated markdown output.

#### Preview referenced image files in markdown

Image references render as embedded previews in reports.

#### Convert note links to markdown links

URL links are rendered with link text in markdown output.

#### Embed Mermaid diagrams from file content

Mermaid content is emitted in fenced blocks for diagram rendering.

#### Lint contextual orphan queries

The CLI exposes `flyb lint orphans` to emit structured diagnostics (`ORPHAN_QUERY_MISSING_LINK`) for notes missing required contextual links, without requiring a label taxonomy.

#### Define contextual orphan query

A subject note (filtered by subject label) is orphan when it has zero matching connections under query filters: relationship label, counterpart note label, and direction in|out|either.

#### Render contextual orphan report section

H3 section arguments can render a deterministic orphan list/table using orphan query filters (`orphan-subject-label`, `orphan-edge-label`, `orphan-counterpart-label`, `orphan-direction`).

#### Guarantee deterministic output ordering

Sort notes, relationships, sections, and arguments with explicit comparators and tie-breakers so repeated runs produce identical output.

#### Define an explicit ordering policy

Ordering policy is part of runtime behavior and contractually defines comparators: notes (primaryLabel, name), relationships (from, to, labelsSortedJoined), sections (case-insensitive title, originalIndex), arguments (name).

#### Select renderer plugin from arguments at runtime

Renderer selection uses one resolved typed renderer argument set sourced from H3Section and note arguments with deterministic precedence and fallback defaults.

#### Register renderers and plugins in a capability registry

A renderer registry maps renderer names to capabilities, supported arguments, and graph-shape compatibility, and defines defaults used by renderer-scoped argument resolution.

#### Generate design reports from configured notes and relationships

This is the primary end-to-end report generation use case.

#### Render graph output as markdown text

Text rendering uses deterministic hierarchy traversal with stable anchors/references for repeated or cyclic nodes plus a short adjacency summary for cyclic graphs.

#### Render graph output as Mermaid diagram

Mermaid output supports visual graph rendering, including cyclic relationships.

#### Render graph output based on graph shape

Renderer behavior adapts to tree, DAG, and cyclic graph structures with deterministic traversal and safe repetition controls.

#### List all configured markdown reports

The CLI exposes a command to enumerate report targets.

#### Build a report from a relationship-label subgraph

Report generation can include only edges matching selected labels, where label references are validated against dataset labels derived from notes and relationships.

#### Allow each H3 section to define cycle policy

H3Section arguments can declare whether cycles are disallowed or allowed.

