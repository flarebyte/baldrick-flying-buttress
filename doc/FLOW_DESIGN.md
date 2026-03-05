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
        - note: Validate referenced labels used by config elements (for example graph.select label arguments) against labelSet; emit `LABEL_REF_UNKNOWN` (default severity `warning`) with argument location and referenced label value for unknown references.
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
        - note: Collect stable diagnostic codes, severities, sources, and precise config locations.
      Normalize validated application model [app.model.normalize]
        - note: Build ValidatedApp with normalized notes, relationships, reports, resolved graph integrity policy, resolved argument registry, and diagnostics. Ordering policy resolution remains generation-time.
    List reports from ValidatedApp [list.reports.output]
      - note: Enumerate reports from the normalized validated model with optional strictness behavior handled by validation policy.
  Generate markdown reports [action.generate.markdown]
    - note: Renders one or more markdown outputs from a single validated application model.
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
        - note: Validate referenced labels used by config elements (for example graph.select label arguments) against labelSet; emit `LABEL_REF_UNKNOWN` (default severity `warning`) with argument location and referenced label value for unknown references.
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
        - note: Collect stable diagnostic codes, severities, sources, and precise config locations.
      Normalize validated application model [app.model.normalize]
        - note: Build ValidatedApp with normalized notes, relationships, reports, resolved graph integrity policy, resolved argument registry, and diagnostics. Ordering policy resolution remains generation-time.
    Generate markdown sections [action.generate.markdown.sections]
      - note: Build H3 sections from note subsets and renderers with deterministic ordering.
      Resolve deterministic ordering policy [ordering.policy.resolve]
        - note: Resolve explicit comparators: notes by (primaryLabel, name) where primaryLabel is the lexicographically smallest label; relationships by (from, to, labelsSortedJoined) where labelsSortedJoined is labels sorted lexicographically then joined with `|`; sections by (lowercase(title), originalIndex) for stable tie-breaks; arguments by argument name.
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
        Render section as a graph [render.section.graph]
          - note: Resolve cycle policy and graph shape, then render with selected renderer(s); renderer/runtime diagnostics here must not duplicate graph-integrity diagnostics.
          Resolve renderer/plugin registry [renderer.registry.resolve]
            - note: Load renderer capabilities, supported arguments, and shape compatibility.
          Resolve renderer-scoped arguments [args.renderer.resolve]
            - note: Collect arguments from H3Section and its notes, keep only keys whose registry scope includes `renderer`, apply precedence (`note` overrides `h3-section`, `h3-section` overrides registry defaults), and produce one typed validated renderer argument set.
          Select renderer plugin from arguments [renderer.plugin.select]
            - note: Choose renderer by resolved typed renderer-scoped arguments with deterministic fallback when unspecified, then pass one resolved renderer argument set to the selected plugin.
          Resolve H3Section cycle policy argument [graph.policy.cycle]
            - note: Use section argument to disallow, allow, or collapse cycles.
          Detect graph shape (tree, DAG, or cyclic) [graph.shape.detect]
            - note: Classify graph structure before selecting rendering strategy; emits renderer/runtime warnings only when applicable.
          Render tree or DAG graph [render.graph.tree-or-dag]
            - note: Prefer hierarchical markdown text; Mermaid can be emitted as an additional diagram.
            Render graph as markdown text [render.graph.markdown.text]
              - note: Render adjacency and hierarchy using the same markdown style as FLOW_DESIGN.
            Render graph as Mermaid [render.graph.mermaid]
              - note: Emit Mermaid syntax for visual rendering in markdown consumers.
          Render cyclic graph [render.graph.circular]
            - note: Prefer Mermaid for cycle readability, with markdown text summary as fallback.
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
            - note: Render as a markdown table or raw CSV code block (for example `format-csv=md`).
            Filter CSV rows by column [file.csv.filter]
              - note: Apply include/exclude filters before rendering CSV output.
          Render section with media file [render.section.file.media]
            - note: Embed image previews for supported media types.
          Render section with code or Mermaid snippet [render.section.file.code]
            - note: Preserve fenced-block formatting for code and Mermaid content.
        Apply deterministic ordering [ordering.apply.deterministic]
          - note: Apply resolved comparators exactly and use stable tie-breakers only (including section originalIndex), yielding reproducible output without runtime randomness.
  Generate JSON graph export [action.generate.json]
    - note: Export notes and relationships in machine-readable JSON format.
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
        - note: Validate referenced labels used by config elements (for example graph.select label arguments) against labelSet; emit `LABEL_REF_UNKNOWN` (default severity `warning`) with argument location and referenced label value for unknown references.
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
        - note: Collect stable diagnostic codes, severities, sources, and precise config locations.
      Normalize validated application model [app.model.normalize]
        - note: Build ValidatedApp with normalized notes, relationships, reports, resolved graph integrity policy, resolved argument registry, and diagnostics. Ordering policy resolution remains generation-time.
    Export validated graph as JSON [export.graph.json]
      - note: Export notes and relationships from ValidatedApp without re-running validation steps.
  Validate the CUE file [action.validate]
    - note: Run canonical application validation and emit the same diagnostics that gate generation.
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
        - note: Validate referenced labels used by config elements (for example graph.select label arguments) against labelSet; emit `LABEL_REF_UNKNOWN` (default severity `warning`) with argument location and referenced label value for unknown references.
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
        - note: Collect stable diagnostic codes, severities, sources, and precise config locations.
      Normalize validated application model [app.model.normalize]
        - note: Build ValidatedApp with normalized notes, relationships, reports, resolved graph integrity policy, resolved argument registry, and diagnostics. Ordering policy resolution remains generation-time.
    Emit structured diagnostics [diagnostics.emit.structured]
      - note: Emit diagnostics with code, severity, source, message, and optional location.
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
  - Emit structured diagnostics — Diagnostics include code, severity, message, source, and optional location context.
  - Report validation diagnostics with locations — Validation errors and warnings should point to config paths and offending argument or relationship names.
  - Define an argument registry schema — Registry entries define argument key, type, default, allowed values, and valid scopes (`h3-section`, `note`, `renderer`).
  - Validate free-form arguments at runtime — Validate against a known argument registry and fail with clear errors on unknown keys or invalid values.
  - Accept free-form arguments on H3Section and Note — Arguments behave like CLI flags (for example `format-csv=md`) and can carry string, string[], boolean, and similar values.
  - Define graph integrity policy beyond cycles — Policy covers missing nodes, orphan nodes, duplicate note names, unknown referenced labels, and cross-report references.
  - Validate graph integrity using policy rules — Integrity checks should emit structured diagnostics tied to offending notes, relationships, and config locations.
  - Resolve arguments by scope — Apply argument rules by scope (h3-section, note, renderer); for renderer scope, collect from H3Section and note arguments and apply precedence (`note` > `h3-section` > registry default).
  - Build a report from a relationship-label subgraph — Report generation can include only edges matching selected labels, where label references are validated against dataset labels derived from notes and relationships.
  - Render note title and markdown description — Each note includes a concise title with free-form markdown content.
  - Guarantee deterministic output ordering — Sort notes, relationships, sections, and arguments with explicit comparators and tie-breakers so repeated runs produce identical output.
  - Define an explicit ordering policy — Ordering policy is part of runtime behavior and contractually defines comparators: notes (primaryLabel, name), relationships (from, to, labelsSortedJoined), sections (case-insensitive title, originalIndex), arguments (name).
  - Keep CUE config compact with argument-driven rendering options — Prefer small composable argument lists over proliferating specialized configuration fields.
  - Coerce free-form argument values into typed values — Convert string-like argument inputs into validated typed values before rendering.
  - Allow each H3 section to define cycle policy — H3Section arguments can declare whether cycles are disallowed, allowed, or collapsed.
  - Render graph output based on graph shape — Renderer behavior adapts to tree, DAG, and cyclic graph structures.
  - Register renderers and plugins in a capability registry — A renderer registry maps renderer names to capabilities, supported arguments, and graph-shape compatibility, and defines defaults used by renderer-scoped argument resolution.
  - Select renderer plugin from arguments at runtime — Renderer selection uses one resolved typed renderer argument set sourced from H3Section and note arguments with deterministic precedence and fallback defaults.
  - Render graph output as markdown text — Text rendering supports readable hierarchy and edge summaries in markdown reports.
  - Render graph output as Mermaid diagram — Mermaid output supports visual graph rendering, including cyclic relationships.
  - Embed Mermaid diagrams from file content — Mermaid content is emitted in fenced blocks for diagram rendering.
  - Convert note links to markdown links — URL links are rendered with link text in markdown output.
  - Reference a file from a note — Referenced files can be embedded in generated markdown output.
  - Embed CSV content from a referenced file — CSV input can render as a markdown table or as raw CSV.
  - Filter embedded CSV rows by column — Column filters reduce CSV output to the relevant subset.
  - Preview referenced image files in markdown — Image references render as embedded previews in reports.


