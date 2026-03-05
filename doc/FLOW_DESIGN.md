# FLOW DESIGN OVERVIEW (Generated)

## Function calls tree

```
flyb CLI root command [cli.root]
  - note: Entry point for report generation, listing, JSON export, and config validation.
  List configured markdown reports [action.list.reports]
    - note: Enumerate report targets from the validated application model without generating files.
    Load CUE application data [load.app.data]
      - note: Read notes, relationships, and report definitions from config.
    Validate CUE application data [validate.app.data]
      - note: Canonical validation pipeline: schema checks, argument registry and free-form argument validation, dataset-based label reference validation, graph integrity policy resolution and graph integrity checks, diagnostic collection, and normalized ValidatedApp output.
      Validate CUE schema and structure [validate.cue.schema]
        - note: Validate required fields, types, and cross-references and attach precise config locations to diagnostics.
      Resolve argument registry schema [args.registry.resolve]
        - note: Load known argument definitions (type, default, allowed values, scopes) where valid scopes are `h3-section`, `note`, and `renderer`.
      Validate argument registry schema consistency [args.registry.validate]
        - note: Validate argument definitions, duplicate keys, scopes, defaults, and allowed values.
      Validate configured free-form arguments [args.validate.config]
        - note: Validate free-form arguments declared in config against registry definitions and scope rules.
      Collect dataset labels [labels.dataset.collect]
        - note: Build authoritative labelSet as the union of labels from note.labels and relationship.labels without enforcing a taxonomy.
      Validate label references [labels.reference.validate]
        - note: Validate referenced labels used by config elements (for example graph.select and orphan-query label arguments) against labelSet; emit `LABEL_REF_UNKNOWN` (default severity `warning`) with argument location and referenced label value for unknown references.
      Resolve graph integrity policy [graph.integrity.policy.resolve]
        - note: Resolve integrity policy for missing nodes, orphans, duplicates, unknown label references, and cross-report references.
      Validate graph integrity [graph.integrity.validate]
        - note: Run integrity checks and emit diagnostics according to resolved policy.
        Check missing relationship nodes [graph.integrity.check.missing-nodes]
          - note: Detect relationships that reference notes that do not exist.
        Check orphan nodes [graph.integrity.check.orphans]
          - note: Detect notes disconnected from report roots/sections.
        Check duplicate note names [graph.integrity.check.duplicate-note-names]
          - note: Detect duplicate note identifiers that can cause ambiguous references.
        Check cross-report references [graph.integrity.check.cross-report-references]
          - note: Validate whether note/edge references across report boundaries are allowed by policy.
      Collect validation diagnostics [diagnostics.collect.validation]
        - note: Collect stable diagnostic codes, severities, sources, canonical machine-readable config `location` paths, and human-readable context fields (`reportTitle`, `sectionTitle`, `noteName`, `relationship`, `argumentName`, `subjectLabel`, `edgeLabel`, `counterpartLabel`).
      Normalize validated application model [app.model.normalize]
        - note: Build ValidatedApp with normalized notes, relationships, reports, resolved graph integrity policy, resolved argument registry, and diagnostics. Ordering policy resolution remains generation-time.
    List reports from ValidatedApp [list.reports.output]
      - note: Enumerate reports from the normalized validated model with optional strictness behavior handled by validation policy.
  List note and relationship names [action.list.names]
    - note: Print note and relationship identifiers for daily inventory with required `--prefix`, optional `--kind notes|relationships|all`, and `--format table|json` (default table).
    Load CUE application data [load.app.data]
      - ref: see first occurrence above for full subtree
    Validate CUE application data [validate.app.data]
      - ref: see first occurrence above for full subtree
    Resolve deterministic ordering policy [ordering.policy.resolve]
      - note: Resolve explicit comparators: notes by (primaryLabel, name) where primaryLabel is the lexicographically smallest label; relationships by (from, to, labelsSortedJoined) where labelsSortedJoined is labels sorted lexicographically then joined with `|`; sections by (lowercase(title), originalIndex) for stable tie-breaks; arguments by argument name.
    Apply deterministic ordering [ordering.apply.deterministic]
      - note: Apply resolved comparators exactly and use stable tie-breakers only (including section originalIndex), yielding reproducible output without runtime randomness.
    Filter names by prefix [names.filter.prefix]
      - note: Apply required `--prefix` filter: keep notes where `name` starts with prefix and relationships where `from` or `to` starts with prefix.
    Filter names by kind [names.filter.kind]
      - note: Apply optional `--kind notes|relationships|all` filter (default `all`) to reduce output noise.
    Output names as table [names.output.table]
      - note: Default output: notes table rows `name | title | labels` and relationship rows `from | to | labels`.
    Output names as JSON [names.output.json]
      - note: Optional `--format json` output as `{ notes: [], relationships: [] }` with the same fields used in table mode.
  Generate markdown reports [action.generate.markdown]
    - note: Renders one or more markdown outputs from a single validated application model.
    Load CUE application data [load.app.data]
      - ref: see first occurrence above for full subtree
    Validate CUE application data [validate.app.data]
      - ref: see first occurrence above for full subtree
    Generate markdown sections [action.generate.markdown.sections]
      - note: Build H3 sections from note subsets and renderers with deterministic ordering.
      Resolve deterministic ordering policy [ordering.policy.resolve]
        - ref: see first occurrence above for full subtree
      Generate a single H3 section [action.generate.markdown.section.h3]
        - note: Compose subgraph, plain content, and file-backed content with section-level arguments.
        Resolve H3Section free-form arguments [args.h3.resolve]
          - note: Read flexible section arguments as key/value flags (for example `graph-renderer=mermaid`) and expose candidates for renderer-scoped resolution.
        Resolve argument registry schema [args.registry.resolve]
          - note: Load known argument definitions (type, default, allowed values, scopes) where valid scopes are `h3-section`, `note`, and `renderer`.
        Validate arguments at runtime [args.validate.runtime]
          - note: Validate keys and values against a known argument registry and fail fast on invalid input.
        Coerce arguments to typed values [args.coerce.typed]
          - note: Coerce validated values to target types (string[], boolean, enum, number).
        Extract subgraph using labels [graph.select]
          - note: Filter notes and relationships by labels and optional starting node; label references are pre-validated against dataset labels (union of note.labels and relationship.labels).
        Render section as contextual orphans [render.section.orphans]
          - note: Resolve orphan query arguments and emit deterministic orphan list/table rows (`name`, `title`, `labels`) for subject notes missing required contextual links.
          Resolve orphan query arguments from H3 section [args.orphan.query.resolve]
            - note: Resolve `orphan-subject-label` (required for orphan mode), optional `orphan-edge-label`, optional `orphan-counterpart-label`, and `orphan-direction in|out|either` (default `either`).
          Find contextual orphans [orphans.query.find]
            - note: For each subject note, require at least one matching relationship under query filters (edge label, counterpart label, direction). Notes with zero matches are contextual orphans.
          Render orphan rows [orphans.render.rows]
            - note: Render deterministic orphan output rows/table with `name`, `title`, and `labels`.
        Render section as a graph [render.section.graph]
          - note: Resolve cycle policy and graph shape, then render with selected renderer(s); renderer/runtime diagnostics here must not duplicate graph-integrity diagnostics.
          Resolve renderer/plugin registry [renderer.registry.resolve]
            - note: Load renderer capabilities, supported arguments, and shape compatibility.
          Resolve renderer-scoped arguments [args.renderer.resolve]
            - note: Collect arguments from H3Section and its notes, keep only keys whose registry scope includes `renderer`, apply precedence (`note` overrides `h3-section`, `h3-section` overrides registry defaults), and produce one typed validated renderer argument set.
          Select renderer plugin from arguments [renderer.plugin.select]
            - note: Choose renderer by resolved typed renderer-scoped arguments with deterministic fallback when unspecified, then pass one resolved renderer argument set to the selected plugin.
          Resolve H3Section cycle policy argument [graph.policy.cycle]
            - note: Resolve section cycle policy (`disallow` or `allow`): `disallow` requires cycle detection error diagnostics and blocks section graph rendering; `allow` permits cyclic rendering.
          Detect graph shape (tree, DAG, or cyclic) [graph.shape.detect]
            - note: Classify selected graph as tree, DAG, or cyclic before renderer selection; if shape is cyclic and cycle-policy is `disallow`, emit error diagnostic and prevent graph rendering for that section.
          Render tree or DAG graph [render.graph.tree-or-dag]
            - note: Prefer hierarchical markdown text; Mermaid can be emitted as an additional diagram.
            Render graph as markdown text [render.graph.markdown.text]
              - note: Render adjacency and hierarchy using the same markdown style as FLOW_DESIGN.
            Render graph as Mermaid [render.graph.mermaid]
              - note: Emit Mermaid syntax for visual rendering in markdown consumers.
          Render cyclic graph [render.graph.circular]
            - note: Render only when cycle-policy is `allow`; prefer Mermaid for deterministic cycle readability, with markdown text summary fallback.
            Render graph as Mermaid [render.graph.mermaid]
              - note: Emit Mermaid syntax for visual rendering in markdown consumers.
            Render graph as markdown text [render.graph.markdown.text]
              - note: Render adjacency and hierarchy using the same markdown style as FLOW_DESIGN.
        Render plain section [render.section.plain]
          - note: Render title and markdown body, including markdown links and note-level argument options.
        Render section with referenced file content [render.section.file]
          - note: Dispatches file rendering by type (CSV, media, code/diagram).
          Resolve Note free-form arguments [args.note.resolve]
            - note: Read note-level rendering options as key/value flags (for example `format-csv=md`) and expose candidates for renderer-scoped resolution with higher precedence than H3Section values.
          Resolve argument registry schema [args.registry.resolve]
            - note: Load known argument definitions (type, default, allowed values, scopes) where valid scopes are `h3-section`, `note`, and `renderer`.
          Validate arguments at runtime [args.validate.runtime]
            - note: Validate keys and values against a known argument registry and fail fast on invalid input.
          Coerce arguments to typed values [args.coerce.typed]
            - note: Coerce validated values to target types (string[], boolean, enum, number).
          Render section with CSV file [render.section.file.csv]
            - note: Render as a markdown table or raw CSV code block (for example `format-csv=md`) and apply note-scoped CSV filters (`csv-include` / `csv-exclude`) using `column:value` exact-match rules.
            Filter CSV rows by column [file.csv.filter]
              - note: Apply exact-match include/exclude filters before rendering CSV output: `csv-include=column:value` keeps matching rows, `csv-exclude=column:value` removes matching rows, and multiple filters are allowed.
          Render section with media file [render.section.file.media]
            - note: Embed image previews for supported media types.
          Render section with code or Mermaid snippet [render.section.file.code]
            - note: Preserve fenced-block formatting for code and Mermaid content.
        Apply deterministic ordering [ordering.apply.deterministic]
          - ref: see first occurrence above for full subtree
  Generate JSON graph export [action.generate.json]
    - note: Export notes and relationships in machine-readable JSON format.
    Load CUE application data [load.app.data]
      - ref: see first occurrence above for full subtree
    Validate CUE application data [validate.app.data]
      - ref: see first occurrence above for full subtree
    Export validated graph as JSON [export.graph.json]
      - note: Export notes and relationships from ValidatedApp without re-running validation steps.
  Validate the CUE file [action.validate]
    - note: Run canonical application validation and emit the same diagnostics that gate generation.
    Load CUE application data [load.app.data]
      - ref: see first occurrence above for full subtree
    Validate CUE application data [validate.app.data]
      - ref: see first occurrence above for full subtree
    Emit structured diagnostics [diagnostics.emit.structured]
      - note: Emit diagnostics with code, severity, source, message, canonical machine-readable `location`, and optional human-readable context fields.
  Lint note and relationship names [action.lint.names]
    - note: Run naming-style hygiene checks with `--style dot|snake|regex` (default dot), optional `--pattern` for regex style, optional `--prefix` scope, and configurable `--severity warning|error` (default warning).
    Load CUE application data [load.app.data]
      - ref: see first occurrence above for full subtree
    Validate CUE application data [validate.app.data]
      - ref: see first occurrence above for full subtree
    Resolve deterministic ordering policy [ordering.policy.resolve]
      - ref: see first occurrence above for full subtree
    Apply deterministic ordering [ordering.apply.deterministic]
      - ref: see first occurrence above for full subtree
    Resolve name style policy [lint.names.policy.resolve]
      - note: Resolve style matcher as case-sensitive policy: `dot`=`^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$`, `snake`=`^[a-z][a-z0-9_]*$`, `regex`=user-provided `--pattern`.
    Filter names by prefix [names.filter.prefix]
      - note: Apply required `--prefix` filter: keep notes where `name` starts with prefix and relationships where `from` or `to` starts with prefix.
    Lint note names [lint.names.notes]
      - note: Check note names against resolved style; emit `NAME_STYLE_VIOLATION` diagnostics with canonical config location and human-readable context for each violation.
    Lint relationship endpoint names [lint.names.relationships]
      - note: Check relationship `from` and `to` endpoint names against resolved style; emit `NAME_STYLE_VIOLATION` diagnostics with canonical config location and relationship context for each violation.
    Emit structured diagnostics [diagnostics.emit.structured]
      - note: Emit diagnostics with code, severity, source, message, canonical machine-readable `location`, and optional human-readable context fields.
  Lint contextual orphans [action.lint.orphans]
    - note: Run orphan-query lint checks with required `--subject-label`, optional `--edge-label`, optional `--counterpart-label`, optional `--direction in|out|either` (default `either`), and configurable `--severity warning|error` (default warning).
    Load CUE application data [load.app.data]
      - ref: see first occurrence above for full subtree
    Validate CUE application data [validate.app.data]
      - ref: see first occurrence above for full subtree
    Resolve deterministic ordering policy [ordering.policy.resolve]
      - ref: see first occurrence above for full subtree
    Apply deterministic ordering [ordering.apply.deterministic]
      - ref: see first occurrence above for full subtree
    Resolve orphan query from CLI flags [lint.orphans.query.resolve]
      - note: Resolve lint flags into orphan query context: `--subject-label` required, `--edge-label` optional, `--counterpart-label` optional, `--direction` default `either`.
    Find contextual orphans [orphans.query.find]
      - note: For each subject note, require at least one matching relationship under query filters (edge label, counterpart label, direction). Notes with zero matches are contextual orphans.
    Emit contextual orphan diagnostics [lint.orphans.emit]
      - note: Emit `ORPHAN_QUERY_MISSING_LINK` diagnostics for each orphan note with canonical config/CLI context location and human-readable fields (`noteName`, `subjectLabel`, `edgeLabel`, `counterpartLabel`).
    Emit structured diagnostics [diagnostics.emit.structured]
      - note: Emit diagnostics with code, severity, source, message, canonical machine-readable `location`, and optional human-readable context fields.
```

## Validation Contract

- `validate.app.data` is the canonical validation entrypoint for all CLI commands and always includes dataset-based label-reference validation plus graph integrity policy resolution and graph integrity validation.
- Guarantees: schema/structure validation, argument registry validation, configured free-form argument validation, dataset-based label-reference validation, graph integrity checks, and structured diagnostics collection with stable codes/severities/sources/locations.
- Return shape: `ValidatedApp` containing:
  - normalized `notes`, `relationships`, and `reports`
  - resolved `graphIntegrityPolicy` and `argumentRegistry`
  - optional ordering policy resolution (currently deferred to generation-time ordering components)
  - `diagnostics: Diagnostic[]` always present (empty when no issues)
- Generation block rule: any `error` severity diagnostic from `validate.app.data` blocks generation; warnings remain non-blocking by default but are still emitted consistently.
- Label reference rule: labels on notes/relationships are free-form definitions; only label references are validated against dataset `labelSet` and unknown references emit `LABEL_REF_UNKNOWN` with default `warning` severity.
- Diagnostic location contract: `location` is the canonical machine-readable index path (Report -> H2 -> H3 -> notes/relationships -> arguments); optional context fields (`reportTitle`, `sectionTitle`, `noteName`, `relationship`, `argumentName`, `subjectLabel`, `edgeLabel`, `counterpartLabel`) provide human-readable debugging context.

## Refactor Notes (Pseudo-diff)

- Removed direct `graph.integrity.validate` calls from `action.generate.markdown.section.h3` and `action.validate`.
- `validate.app.data` now invokes: `validate.cue.schema`, `args.registry.resolve`, `args.registry.validate`, `args.validate.config`, `labels.dataset.collect`, `labels.reference.validate`, `graph.integrity.policy.resolve`, `graph.integrity.validate`, `diagnostics.collect.validation`, `app.model.normalize`.
- Updated command flows to consume `ValidatedApp`:
  - `action.generate.markdown`: `load.app.data -> validate.app.data -> action.generate.markdown.sections`
  - `action.generate.json`: `load.app.data -> validate.app.data -> export.graph.json`
  - `action.validate`: `load.app.data -> validate.app.data -> diagnostics.emit.structured`
  - `action.list.reports`: `load.app.data -> validate.app.data -> list.reports.output`

Supported use cases:

  - Generate design reports from configured notes and relationships — This is the primary end-to-end report generation use case.
  - List all configured markdown reports — The CLI exposes a command to enumerate report targets.
  - Declare multiple markdown reports in one config — A single config can drive generation of multiple report files.
  - Export notes and relationships as JSON — JSON export supports machine-readable graph processing.
  - Define labeled relationships between notes in config — CUE can be used as the source format for flexible configuration; labels on notes and relationships remain free-form.
  - Emit structured diagnostics — Diagnostics include code, severity, message, source, canonical machine-readable `location`, and additional human-readable context fields.
  - Report validation diagnostics with locations — Validation errors and warnings should include canonical index-based config paths plus readable identifiers (report title, section title, note/relationship/argument names).
  - Define an argument registry schema — Registry entries define argument key, type, default, allowed values, and valid scopes (`h3-section`, `note`, `renderer`).
  - Validate free-form arguments at runtime — Validate against a known argument registry and fail with clear errors on unknown keys or invalid values.
  - Accept free-form arguments on H3Section and Note — Arguments behave like CLI flags (for example `format-csv=md`) and can carry string, string[], boolean, and similar values.
  - Define graph integrity policy beyond cycles — Policy covers missing nodes, orphan nodes, duplicate note names, unknown referenced labels, and cross-report references.
  - Validate graph integrity using policy rules — Integrity checks should emit structured diagnostics tied to offending notes, relationships, and config locations.
  - Resolve arguments by scope — Apply argument rules by scope (h3-section, note, renderer); for renderer scope, collect from H3Section and note arguments and apply precedence (`note` > `h3-section` > registry default).
  - Build a report from a relationship-label subgraph — Report generation can include only edges matching selected labels, where label references are validated against dataset labels derived from notes and relationships.
  - List note and relationship names for daily inventory — The CLI exposes `flyb list names` with `--prefix` filtering and `--format table|json` output.
  - Filter names by prefix scope — Prefix filtering keeps notes whose name starts with prefix and relationships where `from` or `to` starts with prefix.
  - Render names as table or JSON — Default output is human-friendly table; JSON is opt-in and returns `{ notes: [], relationships: [] }`.
  - Guarantee deterministic output ordering — Sort notes, relationships, sections, and arguments with explicit comparators and tie-breakers so repeated runs produce identical output.
  - Define an explicit ordering policy — Ordering policy is part of runtime behavior and contractually defines comparators: notes (primaryLabel, name), relationships (from, to, labelsSortedJoined), sections (case-insensitive title, originalIndex), arguments (name).
  - Lint note and relationship names for style hygiene — The CLI exposes `flyb lint names` to emit structured diagnostics for naming-style violations without introducing label taxonomy requirements.
  - Render note title and markdown description — Each note includes a concise title with free-form markdown content.
  - Keep CUE config compact with argument-driven rendering options — Prefer small composable argument lists over proliferating specialized configuration fields.
  - Coerce free-form argument values into typed values — Convert string-like argument inputs into validated typed values before rendering.
  - Render contextual orphan report section — H3 section arguments can render a deterministic orphan list/table using orphan query filters (`orphan-subject-label`, `orphan-edge-label`, `orphan-counterpart-label`, `orphan-direction`).
  - Define contextual orphan query — A subject note (filtered by subject label) is orphan when it has zero matching connections under query filters: relationship label, counterpart note label, and direction in|out|either.
  - Lint contextual orphan queries — The CLI exposes `flyb lint orphans` to emit structured diagnostics (`ORPHAN_QUERY_MISSING_LINK`) for notes missing required contextual links, without requiring a label taxonomy.
  - Allow each H3 section to define cycle policy — H3Section arguments can declare whether cycles are disallowed or allowed.
  - Render graph output based on graph shape — Renderer behavior adapts to tree, DAG, and cyclic graph structures.
  - Register renderers and plugins in a capability registry — A renderer registry maps renderer names to capabilities, supported arguments, and graph-shape compatibility, and defines defaults used by renderer-scoped argument resolution.
  - Select renderer plugin from arguments at runtime — Renderer selection uses one resolved typed renderer argument set sourced from H3Section and note arguments with deterministic precedence and fallback defaults.
  - Render graph output as markdown text — Text rendering supports readable hierarchy and edge summaries in markdown reports.
  - Render graph output as Mermaid diagram — Mermaid output supports visual graph rendering, including cyclic relationships.
  - Embed Mermaid diagrams from file content — Mermaid content is emitted in fenced blocks for diagram rendering.
  - Convert note links to markdown links — URL links are rendered with link text in markdown output.
  - Reference a file from a note — Referenced files can be embedded in generated markdown output.
  - Embed CSV content from a referenced file — CSV input can render as a markdown table or as raw CSV.
  - Filter embedded CSV rows by column — Column filters reduce CSV output to the relevant subset using `csv-include=column:value` and `csv-exclude=column:value` exact-match arguments.
  - Preview referenced image files in markdown — Image references render as embedded previews in reports.
  - Define explicit name style policy — Name styles are case-sensitive: dot=`^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$`, snake=`^[a-z][a-z0-9_]*$`, regex=user-provided `--pattern`.


