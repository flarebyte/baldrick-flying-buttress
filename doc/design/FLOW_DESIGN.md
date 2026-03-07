# Flow Design Overview

Migrated flow call tree and validation pipeline design notes.

## Function calls tree

### Flow call graph

Generated from design/flows.ts.

- <a id="graph-node-cli-root"></a> flyb CLI root command: Entry point for report generation, listing, JSON export, and config validation.
  - <a id="graph-node-action-generate-json"></a> Generate JSON graph export: Export notes and relationships in machine-readable JSON format.
    - <a id="graph-node-export-graph-json"></a> Export validated graph as JSON: Export notes and relationships from ValidatedApp without re-running validation steps.
    - <a id="graph-node-load-app-data"></a> Load CUE application data: Read notes, relationships, and report definitions from config.
    - <a id="graph-node-validate-app-data"></a> Validate CUE application data: Canonical validation pipeline: schema checks, argument registry and free-form argument validation, dataset-based label reference validation, graph integrity policy resolution and graph integrity checks, diagnostic collection, and normalized ValidatedApp output.
      - <a id="graph-node-app-model-normalize"></a> Normalize validated application model: Build ValidatedApp with normalized notes, relationships, reports, resolved graph integrity policy, resolved argument registry, and diagnostics. Ordering policy resolution remains generation-time.
      - <a id="graph-node-args-registry-resolve"></a> Resolve argument registry schema: Load known argument definitions (type, default, allowed values, scopes) where valid scopes are `h3-section`, `note`, and `renderer`.
      - <a id="graph-node-args-registry-validate"></a> Validate argument registry schema consistency: Validate argument definitions, duplicate keys, scopes, defaults, and allowed values.
      - <a id="graph-node-args-validate-config"></a> Validate configured free-form arguments: Validate free-form arguments declared in config against registry definitions and scope rules.
      - <a id="graph-node-diagnostics-collect-validation"></a> Collect validation diagnostics: Collect stable diagnostic codes, severities, sources, canonical machine-readable config `location` paths, and human-readable context fields (`reportTitle`, `sectionTitle`, `noteName`, `relationship`, `argumentName`, `subjectLabel`, `edgeLabel`, `counterpartLabel`).
      - <a id="graph-node-graph-integrity-policy-resolve"></a> Resolve graph integrity policy: Resolve integrity policy for missing nodes, orphans, duplicates, unknown label references, and cross-report references.
      - <a id="graph-node-graph-integrity-validate"></a> Validate graph integrity: Run integrity checks and emit diagnostics according to resolved policy.
        - <a id="graph-node-graph-integrity-check-cross-report-references"></a> Check cross-report references: Validate whether note/edge references across report boundaries are allowed by policy.
        - <a id="graph-node-graph-integrity-check-duplicate-note-names"></a> Check duplicate note names: Detect duplicate note identifiers that can cause ambiguous references.
        - <a id="graph-node-graph-integrity-check-missing-nodes"></a> Check missing relationship nodes: Detect relationships that reference notes that do not exist.
        - <a id="graph-node-graph-integrity-check-orphans"></a> Check orphan nodes: Detect notes disconnected from report roots/sections.
      - <a id="graph-node-labels-dataset-collect"></a> Collect dataset labels: Build authoritative labelSet as the union of labels from note.labels and relationship.labels without enforcing a taxonomy.
      - <a id="graph-node-labels-reference-validate"></a> Validate label references: Validate referenced labels used by config elements (for example graph.select and orphan-query label arguments) against labelSet; emit `LABEL_REF_UNKNOWN` (default severity `warning`) with argument location and referenced label value for unknown references.
      - <a id="graph-node-validate-cue-schema"></a> Validate CUE schema and structure: Validate required fields, types, and cross-references and attach precise config locations to diagnostics.
  - <a id="graph-node-action-generate-markdown"></a> Generate markdown reports: Renders one or more markdown outputs from a single validated application model.
    - <a id="graph-node-action-generate-markdown-sections"></a> Generate markdown sections: Build H3 sections from note subsets and renderers with deterministic ordering.
      - <a id="graph-node-action-generate-markdown-section-h3"></a> Generate a single H3 section: Compose subgraph, plain content, and file-backed content with section-level arguments.
        - <a id="graph-node-args-coerce-typed"></a> Coerce arguments to typed values: Coerce validated values to target types (string[], boolean, enum, number).
        - <a id="graph-node-args-h3-resolve"></a> Resolve H3Section free-form arguments: Read flexible section arguments as key/value flags (for example `graph-renderer=mermaid`) and expose candidates for renderer-scoped resolution.
        - *(see [Resolve argument registry schema](#graph-node-args-registry-resolve))*
        - <a id="graph-node-args-validate-runtime"></a> Validate arguments at runtime: Validate keys and values against a known argument registry and fail fast on invalid input.
        - <a id="graph-node-graph-select"></a> Extract subgraph using labels: Filter notes and relationships by labels and optional starting node; label references are pre-validated against dataset labels (union of note.labels and relationship.labels).
        - <a id="graph-node-ordering-apply-deterministic"></a> Apply deterministic ordering: Apply resolved comparators exactly and use stable tie-breakers only (including section originalIndex), yielding reproducible output without runtime randomness.
        - <a id="graph-node-render-section-file"></a> Render section with referenced file content: Dispatches file rendering by type (CSV, media, code/diagram).
          - *(see [Coerce arguments to typed values](#graph-node-args-coerce-typed))*
          - <a id="graph-node-args-note-resolve"></a> Resolve Note free-form arguments: Read note-level rendering options as key/value flags (for example `format-csv=md`) and expose candidates for renderer-scoped resolution with higher precedence than H3Section values.
          - *(see [Resolve argument registry schema](#graph-node-args-registry-resolve))*
          - *(see [Validate arguments at runtime](#graph-node-args-validate-runtime))*
          - <a id="graph-node-render-section-file-code"></a> Render section with code or Mermaid snippet: Preserve fenced-block formatting for code and Mermaid content.
          - <a id="graph-node-render-section-file-csv"></a> Render section with CSV file: Render as a markdown table or raw CSV code block (for example `format-csv=md`) and apply note-scoped CSV filters (`csv-include` / `csv-exclude`) using `column:value` exact-match rules.
            - <a id="graph-node-file-csv-filter"></a> Filter CSV rows by column: Apply exact-match include/exclude filters before rendering CSV output: `csv-include=column:value` keeps matching rows, `csv-exclude=column:value` removes matching rows, and multiple filters are allowed.
          - <a id="graph-node-render-section-file-media"></a> Render section with media file: Embed image previews for supported media types.
        - <a id="graph-node-render-section-graph"></a> Render section as a graph: Resolve cycle policy and graph shape, then render with selected renderer(s); renderer/runtime diagnostics here must not duplicate graph-integrity diagnostics.
          - <a id="graph-node-args-renderer-resolve"></a> Resolve renderer-scoped arguments: Collect arguments from H3Section and its notes, keep only keys whose registry scope includes `renderer`, apply precedence (`note` overrides `h3-section`, `h3-section` overrides registry defaults), and produce one typed validated renderer argument set.
          - <a id="graph-node-graph-policy-cycle"></a> Resolve H3Section cycle policy argument: Resolve section cycle policy (`disallow` or `allow`): `disallow` requires cycle detection error diagnostics and blocks section graph rendering; `allow` permits cyclic rendering.
          - <a id="graph-node-graph-shape-detect"></a> Detect graph shape (tree, DAG, or cyclic): Classify selected graph as tree, DAG, or cyclic before renderer selection; if shape is cyclic and cycle-policy is `disallow`, emit error diagnostic and prevent graph rendering for that section.
          - <a id="graph-node-render-graph-circular"></a> Render cyclic graph: Render only when cycle-policy is `allow`; markdown traversal expands each node once and when revisiting a node emits `*(cycle back to <node>)*` linking to first anchor, then appends a short deterministic adjacency summary (`A -> B (labels)`). Mermaid remains preferred for cycle readability.
            - <a id="graph-node-render-graph-markdown-text"></a> Render graph as markdown text: Render tree/DAG/cyclic graphs in plain markdown with deterministic traversal order, stable note anchors derived from note names, and reference links to first occurrence anchors for repeated nodes or cycle backs.
            - <a id="graph-node-render-graph-mermaid"></a> Render graph as Mermaid: Emit Mermaid syntax for visual rendering in markdown consumers.
          - <a id="graph-node-render-graph-tree-or-dag"></a> Render tree or DAG graph: Tree: render full hierarchy as nested markdown lists (`**name** — title` plus optional short description). DAG: use stable DFS by ordering policy, expand first encounter, and on repeated encounters allow repetition only when children<=3 and depth<=2; otherwise emit `*(see above)*` reference linking to first anchor.
            - *(see [Render graph as markdown text](#graph-node-render-graph-markdown-text))*
            - *(see [Render graph as Mermaid](#graph-node-render-graph-mermaid))*
          - <a id="graph-node-renderer-plugin-select"></a> Select renderer plugin from arguments: Choose renderer by resolved typed renderer-scoped arguments with deterministic fallback when unspecified, then pass one resolved renderer argument set to the selected plugin.
          - <a id="graph-node-renderer-registry-resolve"></a> Resolve renderer/plugin registry: Load renderer capabilities, supported arguments, and shape compatibility.
        - <a id="graph-node-render-section-orphans"></a> Render section as contextual orphans: Resolve orphan query arguments and emit deterministic orphan list/table rows (`name`, `title`, `labels`) for subject notes missing required contextual links.
          - <a id="graph-node-args-orphan-query-resolve"></a> Resolve orphan query arguments from H3 section: Resolve `orphan-subject-label` (required for orphan mode), optional `orphan-edge-label`, optional `orphan-counterpart-label`, and `orphan-direction in|out|either` (default `either`).
          - <a id="graph-node-orphans-query-find"></a> Find contextual orphans: For each subject note, require at least one matching relationship under query filters (edge label, counterpart label, direction). Notes with zero matches are contextual orphans.
          - <a id="graph-node-orphans-render-rows"></a> Render orphan rows: Render deterministic orphan output rows/table with `name`, `title`, and `labels`.
        - <a id="graph-node-render-section-plain"></a> Render plain section: Render title and markdown body, including markdown links and note-level argument options.
      - <a id="graph-node-ordering-policy-resolve"></a> Resolve deterministic ordering policy: Resolve explicit comparators: notes by (primaryLabel, name) where primaryLabel is the lexicographically smallest label; relationships by (from, to, labelsSortedJoined) where labelsSortedJoined is labels sorted lexicographically then joined with `|`; sections by (lowercase(title), originalIndex) for stable tie-breaks; arguments by argument name.
    - *(see [Load CUE application data](#graph-node-load-app-data))*
    - *(see [Validate CUE application data](#graph-node-validate-app-data))*
  - <a id="graph-node-action-lint-names"></a> Lint note and relationship names: Run naming-style hygiene checks with `--style dot|snake|regex` (default dot), optional `--pattern` for regex style, optional `--prefix` scope, and configurable `--severity warning|error` (default warning).
    - <a id="graph-node-diagnostics-emit-structured"></a> Emit structured diagnostics: Emit diagnostics with code, severity, source, message, canonical machine-readable `location`, and optional human-readable context fields.
    - <a id="graph-node-lint-names-notes"></a> Lint note names: Check note names against resolved style; emit `NAME_STYLE_VIOLATION` diagnostics with canonical config location and human-readable context for each violation.
    - <a id="graph-node-lint-names-policy-resolve"></a> Resolve name style policy: Resolve style matcher as case-sensitive policy: `dot`=`^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$`, `snake`=`^[a-z][a-z0-9_]*$`, `regex`=user-provided `--pattern`.
    - <a id="graph-node-lint-names-relationships"></a> Lint relationship endpoint names: Check relationship `from` and `to` endpoint names against resolved style; emit `NAME_STYLE_VIOLATION` diagnostics with canonical config location and relationship context for each violation.
    - *(see [Load CUE application data](#graph-node-load-app-data))*
    - <a id="graph-node-names-filter-prefix"></a> Filter names by prefix: Apply required `--prefix` filter: keep notes where `name` starts with prefix and relationships where `from` or `to` starts with prefix.
    - *(see [Apply deterministic ordering](#graph-node-ordering-apply-deterministic))*
    - *(see [Resolve deterministic ordering policy](#graph-node-ordering-policy-resolve))*
    - *(see [Validate CUE application data](#graph-node-validate-app-data))*
  - <a id="graph-node-action-lint-orphans"></a> Lint contextual orphans: Run orphan-query lint checks with required `--subject-label`, optional `--edge-label`, optional `--counterpart-label`, optional `--direction in|out|either` (default `either`), and configurable `--severity warning|error` (default warning).
    - *(see [Emit structured diagnostics](#graph-node-diagnostics-emit-structured))*
    - <a id="graph-node-lint-orphans-emit"></a> Emit contextual orphan diagnostics: Emit `ORPHAN_QUERY_MISSING_LINK` diagnostics for each orphan note with canonical config/CLI context location and human-readable fields (`noteName`, `subjectLabel`, `edgeLabel`, `counterpartLabel`).
    - <a id="graph-node-lint-orphans-query-resolve"></a> Resolve orphan query from CLI flags: Resolve lint flags into orphan query context: `--subject-label` required, `--edge-label` optional, `--counterpart-label` optional, `--direction` default `either`.
    - *(see [Load CUE application data](#graph-node-load-app-data))*
    - *(see [Apply deterministic ordering](#graph-node-ordering-apply-deterministic))*
    - *(see [Resolve deterministic ordering policy](#graph-node-ordering-policy-resolve))*
    - *(see [Find contextual orphans](#graph-node-orphans-query-find))*
    - *(see [Validate CUE application data](#graph-node-validate-app-data))*
  - <a id="graph-node-action-list-names"></a> List note and relationship names: Print note and relationship identifiers for daily inventory with required `--prefix`, optional `--kind notes|relationships|all`, and `--format table|json` (default table).
    - *(see [Load CUE application data](#graph-node-load-app-data))*
    - <a id="graph-node-names-filter-kind"></a> Filter names by kind: Apply optional `--kind notes|relationships|all` filter (default `all`) to reduce output noise.
    - *(see [Filter names by prefix](#graph-node-names-filter-prefix))*
    - <a id="graph-node-names-output-json"></a> Output names as JSON: Optional `--format json` output as `{ notes: [], relationships: [] }` with the same fields used in table mode.
    - <a id="graph-node-names-output-table"></a> Output names as table: Default output: notes table rows `name | title | labels` and relationship rows `from | to | labels`.
    - *(see [Apply deterministic ordering](#graph-node-ordering-apply-deterministic))*
    - *(see [Resolve deterministic ordering policy](#graph-node-ordering-policy-resolve))*
    - *(see [Validate CUE application data](#graph-node-validate-app-data))*
  - <a id="graph-node-action-list-reports"></a> List configured markdown reports: Enumerate report targets from the validated application model without generating files.
    - <a id="graph-node-list-reports-output"></a> List reports from ValidatedApp: Enumerate reports from the normalized validated model with optional strictness behavior handled by validation policy.
    - *(see [Load CUE application data](#graph-node-load-app-data))*
    - *(see [Validate CUE application data](#graph-node-validate-app-data))*
  - <a id="graph-node-action-validate"></a> Validate the CUE file: Run canonical application validation and emit the same diagnostics that gate generation.
    - *(see [Emit structured diagnostics](#graph-node-diagnostics-emit-structured))*
    - *(see [Load CUE application data](#graph-node-load-app-data))*
    - *(see [Validate CUE application data](#graph-node-validate-app-data))*

