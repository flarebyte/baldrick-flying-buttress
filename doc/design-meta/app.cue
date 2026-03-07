package flyb

source: "design-migration"
name:   "design-migration"
modules: ["design", "documentation"]

reports: [
  {
    title: "Flow Design Overview"
    filepath: "../design/FLOW_DESIGN.md"
    description: "Migrated flow call tree and validation pipeline design notes."
    sections: [{
      title: "Function calls tree"
      description: ""
      sections: [{
        title: "Flow call graph"
        description: "Generated from design/flows.ts."
        arguments: [
          "graph-subject-label=call",
          "graph-edge-label=contains_step",
          "graph-start-node=cli.root",
          "graph-renderer=markdown-text",
          "cycle-policy=disallow",
        ]
        notes: ["action.generate.json", "action.generate.markdown", "action.generate.markdown.section.h3", "action.generate.markdown.sections", "action.lint.names", "action.lint.orphans", "action.list.names", "action.list.reports", "action.validate", "app.model.normalize", "args.coerce.typed", "args.h3.resolve", "args.note.resolve", "args.orphan.query.resolve", "args.registry.resolve", "args.registry.validate", "args.renderer.resolve", "args.validate.config", "args.validate.runtime", "cli.root", "diagnostics.collect.validation", "diagnostics.emit.structured", "export.graph.json", "file.csv.filter", "graph.integrity.check.cross-report-references", "graph.integrity.check.duplicate-note-names", "graph.integrity.check.missing-nodes", "graph.integrity.check.orphans", "graph.integrity.policy.resolve", "graph.integrity.validate", "graph.policy.cycle", "graph.select", "graph.shape.detect", "labels.dataset.collect", "labels.reference.validate", "lint.names.notes", "lint.names.policy.resolve", "lint.names.relationships", "lint.orphans.emit", "lint.orphans.query.resolve", "list.reports.output", "load.app.data", "names.filter.kind", "names.filter.prefix", "names.output.json", "names.output.table", "ordering.apply.deterministic", "ordering.policy.resolve", "orphans.query.find", "orphans.render.rows", "render.graph.circular", "render.graph.markdown.text", "render.graph.mermaid", "render.graph.tree-or-dag", "render.section.file", "render.section.file.code", "render.section.file.csv", "render.section.file.media", "render.section.graph", "render.section.orphans", "render.section.plain", "renderer.plugin.select", "renderer.registry.resolve", "validate.app.data", "validate.cue.schema"]
      }]
    }]
  },
  {
    title: "Use Cases"
    filepath: "../design/USE_CASES.md"
    description: "Migrated use-case catalog."
    sections: [{
      title: "Use Cases"
      description: "Catalog from design/use_cases.ts."
      sections: [{
        title: "Catalog"
        description: "All use cases."
        notes: ["cli.arguments.free-form", "cli.arguments.registry.schema", "cli.arguments.registry.scope-resolution", "cli.arguments.runtime-validation", "cli.arguments.type-coercion", "cli.config.reduce-noise.with-args", "cli.config.relationships.labeled", "cli.config.reports.multiple", "cli.diagnostics.model", "cli.diagnostics.validation", "cli.export.json.graph", "cli.graph.integrity.policy", "cli.graph.integrity.validation", "cli.names.lint", "cli.names.list", "cli.names.output-formats", "cli.names.prefix-filter", "cli.names.style-policy", "cli.note.basic-markdown", "cli.note.csv.embed", "cli.note.csv.filter-column", "cli.note.filepath.reference", "cli.note.image.preview", "cli.note.link.markdown", "cli.note.mermaid.embed", "cli.orphans.lint", "cli.orphans.query.contextual", "cli.orphans.report.section", "cli.output.deterministic-ordering", "cli.output.deterministic-ordering.policy", "cli.renderer.plugin-selection", "cli.renderer.registry", "cli.report.generate", "cli.report.graph.renderer.markdown-text", "cli.report.graph.renderer.mermaid", "cli.report.graph.shape-aware-render", "cli.report.list", "cli.report.subgraph.by-label", "cli.section.h3.cycle-policy"]
      }]
    }]
  },
  {
    title: "Implementation Considerations"
    filepath: "../design/IMPLEMENTATIONS.md"
    description: "Migrated implementation guidance."
    sections: [{
      title: "Implementations"
      description: "Implementation notes."
      sections: [{
        title: "Catalog"
        description: "All implementation considerations."
        notes: ["cli.arguments.registry-schema", "cli.arguments.runtime-validation", "cli.arguments.scope-resolution", "cli.arguments.type-coercion", "cli.arguments.typed-models", "cli.cobra", "config.arguments.reduce-noise", "config.cue", "diagnostics.structured-model", "diagnostics.validation-location", "graph.integrity.policy-model", "graph.integrity.validation-engine", "lang.go", "names.lint.command", "names.list.command", "orphans.lint.command", "orphans.report.section", "output.ordering.deterministic", "output.ordering.policy-contract", "renderer.registry.contract", "renderer.selection.fallback-policy", "style.errors.guard-clauses", "style.functions.small-single-purpose", "style.io-separate-from-logic", "style.parameters.tiny-structs", "style.predicates.named"]
      }]
    }]
  },
  {
    title: "Risks Overview"
    filepath: "../design/RISKS.md"
    description: "Migrated risk catalog and mitigations."
    sections: [{
      title: "Risks"
      description: "Risk notes."
      sections: [{
        title: "Catalog"
        description: "All risks."
        notes: ["authoring.cue-ai-assistance-gap", "collaboration.cue-merge-conflicts", "graph.circular-dependency", "maintenance.single-cue-file-size", "performance.report-generation-scale"]
      }]
    }]
  },
]

argumentRegistry: {
  version: "1"
  arguments: []
}


notes: [
  {
    name: "action.generate.json"
    title: "Generate JSON graph export"
    markdown: "Export notes and relationships in machine-readable JSON format."
    labels: ["flow", "call"]
  },
  {
    name: "action.generate.markdown"
    title: "Generate markdown reports"
    markdown: "Renders one or more markdown outputs from a single validated application model."
    labels: ["flow", "call"]
  },
  {
    name: "action.generate.markdown.section.h3"
    title: "Generate a single H3 section"
    markdown: "Compose subgraph, plain content, and file-backed content with section-level arguments."
    labels: ["flow", "call"]
  },
  {
    name: "action.generate.markdown.sections"
    title: "Generate markdown sections"
    markdown: "Build H3 sections from note subsets and renderers with deterministic ordering."
    labels: ["flow", "call"]
  },
  {
    name: "action.lint.names"
    title: "Lint note and relationship names"
    markdown: "Run naming-style hygiene checks with `--style dot|snake|regex` (default dot), optional `--pattern` for regex style, optional `--prefix` scope, and configurable `--severity warning|error` (default warning)."
    labels: ["flow", "call"]
  },
  {
    name: "action.lint.orphans"
    title: "Lint contextual orphans"
    markdown: "Run orphan-query lint checks with required `--subject-label`, optional `--edge-label`, optional `--counterpart-label`, optional `--direction in|out|either` (default `either`), and configurable `--severity warning|error` (default warning)."
    labels: ["flow", "call"]
  },
  {
    name: "action.list.names"
    title: "List note and relationship names"
    markdown: "Print note and relationship identifiers for daily inventory with required `--prefix`, optional `--kind notes|relationships|all`, and `--format table|json` (default table)."
    labels: ["flow", "call"]
  },
  {
    name: "action.list.reports"
    title: "List configured markdown reports"
    markdown: "Enumerate report targets from the validated application model without generating files."
    labels: ["flow", "call"]
  },
  {
    name: "action.validate"
    title: "Validate the CUE file"
    markdown: "Run canonical application validation and emit the same diagnostics that gate generation."
    labels: ["flow", "call"]
  },
  {
    name: "app.model.normalize"
    title: "Normalize validated application model"
    markdown: "Build ValidatedApp with normalized notes, relationships, reports, resolved graph integrity policy, resolved argument registry, and diagnostics. Ordering policy resolution remains generation-time."
    labels: ["flow", "call"]
  },
  {
    name: "args.coerce.typed"
    title: "Coerce arguments to typed values"
    markdown: "Coerce validated values to target types (string[], boolean, enum, number)."
    labels: ["flow", "call"]
  },
  {
    name: "args.h3.resolve"
    title: "Resolve H3Section free-form arguments"
    markdown: "Read flexible section arguments as key/value flags (for example `graph-renderer=mermaid`) and expose candidates for renderer-scoped resolution."
    labels: ["flow", "call"]
  },
  {
    name: "args.note.resolve"
    title: "Resolve Note free-form arguments"
    markdown: "Read note-level rendering options as key/value flags (for example `format-csv=md`) and expose candidates for renderer-scoped resolution with higher precedence than H3Section values."
    labels: ["flow", "call"]
  },
  {
    name: "args.orphan.query.resolve"
    title: "Resolve orphan query arguments from H3 section"
    markdown: "Resolve `orphan-subject-label` (required for orphan mode), optional `orphan-edge-label`, optional `orphan-counterpart-label`, and `orphan-direction in|out|either` (default `either`)."
    labels: ["flow", "call"]
  },
  {
    name: "args.registry.resolve"
    title: "Resolve argument registry schema"
    markdown: "Load known argument definitions (type, default, allowed values, scopes) where valid scopes are `h3-section`, `note`, and `renderer`."
    labels: ["flow", "call"]
  },
  {
    name: "args.registry.validate"
    title: "Validate argument registry schema consistency"
    markdown: "Validate argument definitions, duplicate keys, scopes, defaults, and allowed values."
    labels: ["flow", "call"]
  },
  {
    name: "args.renderer.resolve"
    title: "Resolve renderer-scoped arguments"
    markdown: "Collect arguments from H3Section and its notes, keep only keys whose registry scope includes `renderer`, apply precedence (`note` overrides `h3-section`, `h3-section` overrides registry defaults), and produce one typed validated renderer argument set."
    labels: ["flow", "call"]
  },
  {
    name: "args.validate.config"
    title: "Validate configured free-form arguments"
    markdown: "Validate free-form arguments declared in config against registry definitions and scope rules."
    labels: ["flow", "call"]
  },
  {
    name: "args.validate.runtime"
    title: "Validate arguments at runtime"
    markdown: "Validate keys and values against a known argument registry and fail fast on invalid input."
    labels: ["flow", "call"]
  },
  {
    name: "authoring.cue-ai-assistance-gap"
    title: "Limited practical AI assistance for CUE authoring"
    markdown: "Description: Generative AI tools often have weaker CUE support than for mainstream formats, which can reduce productivity and increase hand-written config errors.\n\nMitigation: Provide templates, examples, and lintable conventions in-repo; rely on strong validation and developer documentation rather than AI-generated CUE."
    labels: ["risk", "design"]
  },
  {
    name: "cli.arguments.free-form"
    title: "Accept free-form arguments on H3Section and Note"
    markdown: "Arguments behave like CLI flags (for example `format-csv=md`) and can carry string, string[], boolean, and similar values."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.arguments.registry-schema"
    title: "Define a typed argument registry schema"
    markdown: "Maintain a registry of argument definitions (name, type, default, allowed values, scopes) and use it as the single source of truth for argument behavior; valid scopes are `h3-section`, `note`, and `renderer`."
    labels: ["implementation", "design"]
  },
  {
    name: "cli.arguments.registry.schema"
    title: "Define an argument registry schema"
    markdown: "Registry entries define argument key, type, default, allowed values, and valid scopes (`h3-section`, `note`, `renderer`)."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.arguments.registry.scope-resolution"
    title: "Resolve arguments by scope"
    markdown: "Apply argument rules by scope (h3-section, note, renderer); for renderer scope, collect from H3Section and note arguments and apply precedence (`note` > `h3-section` > registry default)."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.arguments.runtime-validation"
    title: "Validate free-form arguments at runtime"
    markdown: "Validate against a known argument registry and fail with clear errors on unknown keys or invalid values."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.arguments.scope-resolution"
    title: "Apply scope-aware argument resolution"
    markdown: "Resolve and validate arguments by scope (h3-section, note, renderer) so options are accepted only where they are meaningful; renderer-scoped arguments are collected from H3Section and note argument lists with precedence `note` > `h3-section` > registry defaults."
    labels: ["implementation", "design"]
  },
  {
    name: "cli.arguments.type-coercion"
    title: "Coerce free-form argument values into typed values"
    markdown: "Convert string-like argument inputs into validated typed values before rendering."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.arguments.typed-models"
    title: "Use free-form key/value arguments with typed coercion"
    markdown: "Treat H3Section and Note arguments like CLI-style flags (for example `format-csv=md`) to keep config flexible, then coerce values into typed runtime options per renderer."
    labels: ["implementation", "design"]
  },
  {
    name: "cli.cobra"
    title: "Use Cobra for CLI command and argument parsing"
    markdown: "Use Cobra to model commands, flags, and subcommands (`generate markdown`, `generate json`, `validate`, `list reports`, `list names`, `lint names`) with a consistent command tree."
    labels: ["implementation", "design"]
  },
  {
    name: "cli.config.reduce-noise.with-args"
    title: "Keep CUE config compact with argument-driven rendering options"
    markdown: "Prefer small composable argument lists over proliferating specialized configuration fields."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.config.relationships.labeled"
    title: "Define labeled relationships between notes in config"
    markdown: "CUE can be used as the source format for flexible configuration; labels on notes and relationships remain free-form."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.config.reports.multiple"
    title: "Declare multiple markdown reports in one config"
    markdown: "A single config can drive generation of multiple report files."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.diagnostics.model"
    title: "Emit structured diagnostics"
    markdown: "Diagnostics include code, severity, message, source, canonical machine-readable `location`, and additional human-readable context fields."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.diagnostics.validation"
    title: "Report validation diagnostics with locations"
    markdown: "Validation errors and warnings should include canonical index-based config paths plus readable identifiers (report title, section title, note/relationship/argument names)."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.export.json.graph"
    title: "Export notes and relationships as JSON"
    markdown: "JSON export supports machine-readable graph processing."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.graph.integrity.policy"
    title: "Define graph integrity policy beyond cycles"
    markdown: "Policy covers missing nodes, orphan nodes, duplicate note names, unknown referenced labels, and cross-report references."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.graph.integrity.validation"
    title: "Validate graph integrity using policy rules"
    markdown: "Integrity checks should emit structured diagnostics tied to offending notes, relationships, and config locations."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.names.lint"
    title: "Lint note and relationship names for style hygiene"
    markdown: "The CLI exposes `flyb lint names` to emit structured diagnostics for naming-style violations without introducing label taxonomy requirements."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.names.list"
    title: "List note and relationship names for daily inventory"
    markdown: "The CLI exposes `flyb list names` with `--prefix` filtering and `--format table|json` output."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.names.output-formats"
    title: "Render names as table or JSON"
    markdown: "Default output is human-friendly table; JSON is opt-in and returns `{ notes: [], relationships: [] }`."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.names.prefix-filter"
    title: "Filter names by prefix scope"
    markdown: "Prefix filtering keeps notes whose name starts with prefix and relationships where `from` or `to` starts with prefix."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.names.style-policy"
    title: "Define explicit name style policy"
    markdown: "Name styles are case-sensitive: dot=`^[a-z][a-z0-9]*(\\.[a-z][a-z0-9]*)*$`, snake=`^[a-z][a-z0-9_]*$`, regex=user-provided `--pattern`."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.note.basic-markdown"
    title: "Render note title and markdown description"
    markdown: "Each note includes a concise title with free-form markdown content."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.note.csv.embed"
    title: "Embed CSV content from a referenced file"
    markdown: "CSV input can render as a markdown table or as raw CSV."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.note.csv.filter-column"
    title: "Filter embedded CSV rows by column"
    markdown: "Column filters reduce CSV output to the relevant subset using `csv-include=column:value` and `csv-exclude=column:value` exact-match arguments."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.note.filepath.reference"
    title: "Reference a file from a note"
    markdown: "Referenced files can be embedded in generated markdown output."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.note.image.preview"
    title: "Preview referenced image files in markdown"
    markdown: "Image references render as embedded previews in reports."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.note.link.markdown"
    title: "Convert note links to markdown links"
    markdown: "URL links are rendered with link text in markdown output."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.note.mermaid.embed"
    title: "Embed Mermaid diagrams from file content"
    markdown: "Mermaid content is emitted in fenced blocks for diagram rendering."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.orphans.lint"
    title: "Lint contextual orphan queries"
    markdown: "The CLI exposes `flyb lint orphans` to emit structured diagnostics (`ORPHAN_QUERY_MISSING_LINK`) for notes missing required contextual links, without requiring a label taxonomy."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.orphans.query.contextual"
    title: "Define contextual orphan query"
    markdown: "A subject note (filtered by subject label) is orphan when it has zero matching connections under query filters: relationship label, counterpart note label, and direction in|out|either."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.orphans.report.section"
    title: "Render contextual orphan report section"
    markdown: "H3 section arguments can render a deterministic orphan list/table using orphan query filters (`orphan-subject-label`, `orphan-edge-label`, `orphan-counterpart-label`, `orphan-direction`)."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.output.deterministic-ordering"
    title: "Guarantee deterministic output ordering"
    markdown: "Sort notes, relationships, sections, and arguments with explicit comparators and tie-breakers so repeated runs produce identical output."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.output.deterministic-ordering.policy"
    title: "Define an explicit ordering policy"
    markdown: "Ordering policy is part of runtime behavior and contractually defines comparators: notes (primaryLabel, name), relationships (from, to, labelsSortedJoined), sections (case-insensitive title, originalIndex), arguments (name)."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.renderer.plugin-selection"
    title: "Select renderer plugin from arguments at runtime"
    markdown: "Renderer selection uses one resolved typed renderer argument set sourced from H3Section and note arguments with deterministic precedence and fallback defaults."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.renderer.registry"
    title: "Register renderers and plugins in a capability registry"
    markdown: "A renderer registry maps renderer names to capabilities, supported arguments, and graph-shape compatibility, and defines defaults used by renderer-scoped argument resolution."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.report.generate"
    title: "Generate design reports from configured notes and relationships"
    markdown: "This is the primary end-to-end report generation use case."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.report.graph.renderer.markdown-text"
    title: "Render graph output as markdown text"
    markdown: "Text rendering uses deterministic hierarchy traversal with stable anchors/references for repeated or cyclic nodes plus a short adjacency summary for cyclic graphs."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.report.graph.renderer.mermaid"
    title: "Render graph output as Mermaid diagram"
    markdown: "Mermaid output supports visual graph rendering, including cyclic relationships."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.report.graph.shape-aware-render"
    title: "Render graph output based on graph shape"
    markdown: "Renderer behavior adapts to tree, DAG, and cyclic graph structures with deterministic traversal and safe repetition controls."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.report.list"
    title: "List all configured markdown reports"
    markdown: "The CLI exposes a command to enumerate report targets."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.report.subgraph.by-label"
    title: "Build a report from a relationship-label subgraph"
    markdown: "Report generation can include only edges matching selected labels, where label references are validated against dataset labels derived from notes and relationships."
    labels: ["usecase", "design"]
  },
  {
    name: "cli.root"
    title: "flyb CLI root command"
    markdown: "Entry point for report generation, listing, JSON export, and config validation."
    labels: ["flow", "call"]
  },
  {
    name: "cli.section.h3.cycle-policy"
    title: "Allow each H3 section to define cycle policy"
    markdown: "H3Section arguments can declare whether cycles are disallowed or allowed."
    labels: ["usecase", "design"]
  },
  {
    name: "collaboration.cue-merge-conflicts"
    title: "Higher merge-conflict risk in shared CUE files"
    markdown: "Description: When multiple developers edit the same large CUE file, concurrent changes can frequently overlap and create conflict-heavy pull requests.\n\nMitigation: Reduce shared hotspots with file partitioning, stable key ordering, and ownership boundaries; add CI validation to catch conflicts and schema drift early."
    labels: ["risk", "design"]
  },
  {
    name: "config.arguments.reduce-noise"
    title: "Use arguments to reduce CUE configuration noise"
    markdown: "Prefer composable argument lists over adding many specialized CUE fields, so rendering capabilities can evolve without large schema churn."
    labels: ["implementation", "design"]
  },
  {
    name: "config.cue"
    title: "Use CUE as the configuration source of truth"
    markdown: "Represent notes, relationships, and report definitions in CUE for schema validation, defaults, and composable configuration."
    labels: ["implementation", "design"]
  },
  {
    name: "diagnostics.collect.validation"
    title: "Collect validation diagnostics"
    markdown: "Collect stable diagnostic codes, severities, sources, canonical machine-readable config `location` paths, and human-readable context fields (`reportTitle`, `sectionTitle`, `noteName`, `relationship`, `argumentName`, `subjectLabel`, `edgeLabel`, `counterpartLabel`)."
    labels: ["flow", "call"]
  },
  {
    name: "diagnostics.emit.structured"
    title: "Emit structured diagnostics"
    markdown: "Emit diagnostics with code, severity, source, message, canonical machine-readable `location`, and optional human-readable context fields."
    labels: ["flow", "call"]
  },
  {
    name: "diagnostics.structured-model"
    title: "Use a structured diagnostics model"
    markdown: "Standardize diagnostics with code, severity, source, message, canonical machine-readable location, and optional human-readable context fields to support CLI UX, CI checks, and future editor integrations."
    labels: ["implementation", "design"]
  },
  {
    name: "diagnostics.validation-location"
    title: "Attach validation diagnostics to precise config locations"
    markdown: "Include canonical index-based CUE path plus related report/section titles and note/relationship/argument identifiers in diagnostics so users can quickly fix invalid configuration."
    labels: ["implementation", "design"]
  },
  {
    name: "export.graph.json"
    title: "Export validated graph as JSON"
    markdown: "Export notes and relationships from ValidatedApp without re-running validation steps."
    labels: ["flow", "call"]
  },
  {
    name: "file.csv.filter"
    title: "Filter CSV rows by column"
    markdown: "Apply exact-match include/exclude filters before rendering CSV output: `csv-include=column:value` keeps matching rows, `csv-exclude=column:value` removes matching rows, and multiple filters are allowed."
    labels: ["flow", "call"]
  },
  {
    name: "graph.circular-dependency"
    title: "Circular dependencies in note relationships"
    markdown: "Description: Tree or DAG-like relationship graphs are usually straightforward, but circular dependencies can break assumptions in traversal, filtering, and report assembly.\n\nMitigation: Add explicit cycle detection and policy controls (`disallow` to emit error and skip section graph rendering, `allow` to render cyclic graphs), and test traversal logic with cyclic graph fixtures."
    labels: ["risk", "design"]
  },
  {
    name: "graph.integrity.check.cross-report-references"
    title: "Check cross-report references"
    markdown: "Validate whether note/edge references across report boundaries are allowed by policy."
    labels: ["flow", "call"]
  },
  {
    name: "graph.integrity.check.duplicate-note-names"
    title: "Check duplicate note names"
    markdown: "Detect duplicate note identifiers that can cause ambiguous references."
    labels: ["flow", "call"]
  },
  {
    name: "graph.integrity.check.missing-nodes"
    title: "Check missing relationship nodes"
    markdown: "Detect relationships that reference notes that do not exist."
    labels: ["flow", "call"]
  },
  {
    name: "graph.integrity.check.orphans"
    title: "Check orphan nodes"
    markdown: "Detect notes disconnected from report roots/sections."
    labels: ["flow", "call"]
  },
  {
    name: "graph.integrity.policy-model"
    title: "Define a graph integrity policy model"
    markdown: "Define explicit integrity rules for missing nodes, orphan nodes, duplicate names, and cross-report references with per-rule severity; validate label references separately against dataset-derived labels."
    labels: ["implementation", "design"]
  },
  {
    name: "graph.integrity.policy.resolve"
    title: "Resolve graph integrity policy"
    markdown: "Resolve integrity policy for missing nodes, orphans, duplicates, unknown label references, and cross-report references."
    labels: ["flow", "call"]
  },
  {
    name: "graph.integrity.validate"
    title: "Validate graph integrity"
    markdown: "Run integrity checks and emit diagnostics according to resolved policy."
    labels: ["flow", "call"]
  },
  {
    name: "graph.integrity.validation-engine"
    title: "Implement graph integrity validation checks"
    markdown: "Run focused integrity checks and emit structured diagnostics linked to note names, relationships, arguments, and CUE paths; keep label-definition handling free-form and validate only label references."
    labels: ["implementation", "design"]
  },
  {
    name: "graph.policy.cycle"
    title: "Resolve H3Section cycle policy argument"
    markdown: "Resolve section cycle policy (`disallow` or `allow`): `disallow` requires cycle detection error diagnostics and blocks section graph rendering; `allow` permits cyclic rendering."
    labels: ["flow", "call"]
  },
  {
    name: "graph.select"
    title: "Extract subgraph using labels"
    markdown: "Filter notes and relationships by labels and optional starting node; label references are pre-validated against dataset labels (union of note.labels and relationship.labels)."
    labels: ["flow", "call"]
  },
  {
    name: "graph.shape.detect"
    title: "Detect graph shape (tree, DAG, or cyclic)"
    markdown: "Classify selected graph as tree, DAG, or cyclic before renderer selection; if shape is cyclic and cycle-policy is `disallow`, emit error diagnostic and prevent graph rendering for that section."
    labels: ["flow", "call"]
  },
  {
    name: "labels.dataset.collect"
    title: "Collect dataset labels"
    markdown: "Build authoritative labelSet as the union of labels from note.labels and relationship.labels without enforcing a taxonomy."
    labels: ["flow", "call"]
  },
  {
    name: "labels.reference.validate"
    title: "Validate label references"
    markdown: "Validate referenced labels used by config elements (for example graph.select and orphan-query label arguments) against labelSet; emit `LABEL_REF_UNKNOWN` (default severity `warning`) with argument location and referenced label value for unknown references."
    labels: ["flow", "call"]
  },
  {
    name: "lang.go"
    title: "Implement the CLI in Go"
    markdown: "Use Go as the primary implementation language for strong typing, fast startup, and straightforward single-binary distribution."
    labels: ["implementation", "design"]
  },
  {
    name: "lint.names.notes"
    title: "Lint note names"
    markdown: "Check note names against resolved style; emit `NAME_STYLE_VIOLATION` diagnostics with canonical config location and human-readable context for each violation."
    labels: ["flow", "call"]
  },
  {
    name: "lint.names.policy.resolve"
    title: "Resolve name style policy"
    markdown: "Resolve style matcher as case-sensitive policy: `dot`=`^[a-z][a-z0-9]*(\\.[a-z][a-z0-9]*)*$`, `snake`=`^[a-z][a-z0-9_]*$`, `regex`=user-provided `--pattern`."
    labels: ["flow", "call"]
  },
  {
    name: "lint.names.relationships"
    title: "Lint relationship endpoint names"
    markdown: "Check relationship `from` and `to` endpoint names against resolved style; emit `NAME_STYLE_VIOLATION` diagnostics with canonical config location and relationship context for each violation."
    labels: ["flow", "call"]
  },
  {
    name: "lint.orphans.emit"
    title: "Emit contextual orphan diagnostics"
    markdown: "Emit `ORPHAN_QUERY_MISSING_LINK` diagnostics for each orphan note with canonical config/CLI context location and human-readable fields (`noteName`, `subjectLabel`, `edgeLabel`, `counterpartLabel`)."
    labels: ["flow", "call"]
  },
  {
    name: "lint.orphans.query.resolve"
    title: "Resolve orphan query from CLI flags"
    markdown: "Resolve lint flags into orphan query context: `--subject-label` required, `--edge-label` optional, `--counterpart-label` optional, `--direction` default `either`."
    labels: ["flow", "call"]
  },
  {
    name: "list.reports.output"
    title: "List reports from ValidatedApp"
    markdown: "Enumerate reports from the normalized validated model with optional strictness behavior handled by validation policy."
    labels: ["flow", "call"]
  },
  {
    name: "load.app.data"
    title: "Load CUE application data"
    markdown: "Read notes, relationships, and report definitions from config."
    labels: ["flow", "call"]
  },
  {
    name: "maintenance.single-cue-file-size"
    title: "Single large CUE file becomes hard to maintain"
    markdown: "Description: Packing too many notes and relationships into one CUE file increases cognitive load, makes reviews difficult, and raises the chance of accidental breakage.\n\nMitigation: Split configuration into modular CUE packages/files by domain or report, then compose them through imports and shared schema constraints."
    labels: ["risk", "design"]
  },
  {
    name: "names.filter.kind"
    title: "Filter names by kind"
    markdown: "Apply optional `--kind notes|relationships|all` filter (default `all`) to reduce output noise."
    labels: ["flow", "call"]
  },
  {
    name: "names.filter.prefix"
    title: "Filter names by prefix"
    markdown: "Apply required `--prefix` filter: keep notes where `name` starts with prefix and relationships where `from` or `to` starts with prefix."
    labels: ["flow", "call"]
  },
  {
    name: "names.lint.command"
    title: "Implement lint names hygiene command"
    markdown: "Implement `flyb lint names` with style policy (`dot|snake|regex`), optional regex `--pattern`, optional prefix scope, and configurable severity; emit structured `NAME_STYLE_VIOLATION` diagnostics with canonical config locations and readable context."
    labels: ["implementation", "design"]
  },
  {
    name: "names.list.command"
    title: "Implement list names inventory command"
    markdown: "Implement `flyb list names` with required `--prefix`, optional `--kind notes|relationships|all`, and `--format table|json` (default table); reuse validated app data and deterministic ordering before filtering/output."
    labels: ["implementation", "design"]
  },
  {
    name: "names.output.json"
    title: "Output names as JSON"
    markdown: "Optional `--format json` output as `{ notes: [], relationships: [] }` with the same fields used in table mode."
    labels: ["flow", "call"]
  },
  {
    name: "names.output.table"
    title: "Output names as table"
    markdown: "Default output: notes table rows `name | title | labels` and relationship rows `from | to | labels`."
    labels: ["flow", "call"]
  },
  {
    name: "ordering.apply.deterministic"
    title: "Apply deterministic ordering"
    markdown: "Apply resolved comparators exactly and use stable tie-breakers only (including section originalIndex), yielding reproducible output without runtime randomness."
    labels: ["flow", "call"]
  },
  {
    name: "ordering.policy.resolve"
    title: "Resolve deterministic ordering policy"
    markdown: "Resolve explicit comparators: notes by (primaryLabel, name) where primaryLabel is the lexicographically smallest label; relationships by (from, to, labelsSortedJoined) where labelsSortedJoined is labels sorted lexicographically then joined with `|`; sections by (lowercase(title), originalIndex) for stable tie-breaks; arguments by argument name."
    labels: ["flow", "call"]
  },
  {
    name: "orphans.lint.command"
    title: "Implement contextual orphan lint command"
    markdown: "Implement `flyb lint orphans` using orphan-query filters (`subject-label`, optional edge/counterpart labels, direction) and emit deterministic `ORPHAN_QUERY_MISSING_LINK` diagnostics with stable locations/context."
    labels: ["implementation", "design"]
  },
  {
    name: "orphans.query.find"
    title: "Find contextual orphans"
    markdown: "For each subject note, require at least one matching relationship under query filters (edge label, counterpart label, direction). Notes with zero matches are contextual orphans."
    labels: ["flow", "call"]
  },
  {
    name: "orphans.render.rows"
    title: "Render orphan rows"
    markdown: "Render deterministic orphan output rows/table with `name`, `title`, and `labels`."
    labels: ["flow", "call"]
  },
  {
    name: "orphans.report.section"
    title: "Implement contextual orphan report section renderer"
    markdown: "Implement H3 orphan section rendering using orphan-query arguments and deterministic row/table output (`name`, `title`, `labels`) so report sections and lint command evaluate the same orphan set."
    labels: ["implementation", "design"]
  },
  {
    name: "output.ordering.deterministic"
    title: "Guarantee deterministic ordering in generated outputs"
    markdown: "Apply explicit stable sorting for notes, relationships, sections, and arguments using concrete comparators (notes: primaryLabel/name, relationships: from/to/labelsSortedJoined, sections: case-insensitive title plus originalIndex, arguments: name) so output remains reproducible across runs and machines."
    labels: ["implementation", "design"]
  },
  {
    name: "output.ordering.policy-contract"
    title: "Treat ordering policy as a testable contract"
    markdown: "Define ordering rules and tie-breakers as a versioned policy (including label normalization and relationship label joining rules) and verify them with golden-file tests."
    labels: ["implementation", "design"]
  },
  {
    name: "performance.report-generation-scale"
    title: "Report generation performance at scale"
    markdown: "Description: Generating many sections from large note graphs and file-backed content can increase CPU, memory, and I/O cost, causing slow report builds and degraded CLI responsiveness.\n\nMitigation: Use bounded concurrency, lazy file loading, and optional caching of parsed CUE and graph selections; add profiling baselines and fail-fast limits for oversized runs."
    labels: ["risk", "design"]
  },
  {
    name: "render.graph.circular"
    title: "Render cyclic graph"
    markdown: "Render only when cycle-policy is `allow`; markdown traversal expands each node once and when revisiting a node emits `*(cycle back to <node>)*` linking to first anchor, then appends a short deterministic adjacency summary (`A -> B (labels)`). Mermaid remains preferred for cycle readability."
    labels: ["flow", "call"]
  },
  {
    name: "render.graph.markdown.text"
    title: "Render graph as markdown text"
    markdown: "Render tree/DAG/cyclic graphs in plain markdown with deterministic traversal order, stable note anchors derived from note names, and reference links to first occurrence anchors for repeated nodes or cycle backs."
    labels: ["flow", "call"]
  },
  {
    name: "render.graph.mermaid"
    title: "Render graph as Mermaid"
    markdown: "Emit Mermaid syntax for visual rendering in markdown consumers."
    labels: ["flow", "call"]
  },
  {
    name: "render.graph.tree-or-dag"
    title: "Render tree or DAG graph"
    markdown: "Tree: render full hierarchy as nested markdown lists (`**name** — title` plus optional short description). DAG: use stable DFS by ordering policy, expand first encounter, and on repeated encounters allow repetition only when children<=3 and depth<=2; otherwise emit `*(see above)*` reference linking to first anchor."
    labels: ["flow", "call"]
  },
  {
    name: "render.section.file"
    title: "Render section with referenced file content"
    markdown: "Dispatches file rendering by type (CSV, media, code/diagram)."
    labels: ["flow", "call"]
  },
  {
    name: "render.section.file.code"
    title: "Render section with code or Mermaid snippet"
    markdown: "Preserve fenced-block formatting for code and Mermaid content."
    labels: ["flow", "call"]
  },
  {
    name: "render.section.file.csv"
    title: "Render section with CSV file"
    markdown: "Render as a markdown table or raw CSV code block (for example `format-csv=md`) and apply note-scoped CSV filters (`csv-include` / `csv-exclude`) using `column:value` exact-match rules."
    labels: ["flow", "call"]
  },
  {
    name: "render.section.file.media"
    title: "Render section with media file"
    markdown: "Embed image previews for supported media types."
    labels: ["flow", "call"]
  },
  {
    name: "render.section.graph"
    title: "Render section as a graph"
    markdown: "Resolve cycle policy and graph shape, then render with selected renderer(s); renderer/runtime diagnostics here must not duplicate graph-integrity diagnostics."
    labels: ["flow", "call"]
  },
  {
    name: "render.section.orphans"
    title: "Render section as contextual orphans"
    markdown: "Resolve orphan query arguments and emit deterministic orphan list/table rows (`name`, `title`, `labels`) for subject notes missing required contextual links."
    labels: ["flow", "call"]
  },
  {
    name: "render.section.plain"
    title: "Render plain section"
    markdown: "Render title and markdown body, including markdown links and note-level argument options."
    labels: ["flow", "call"]
  },
  {
    name: "renderer.plugin.select"
    title: "Select renderer plugin from arguments"
    markdown: "Choose renderer by resolved typed renderer-scoped arguments with deterministic fallback when unspecified, then pass one resolved renderer argument set to the selected plugin."
    labels: ["flow", "call"]
  },
  {
    name: "renderer.registry.contract"
    title: "Define a renderer plugin registry contract"
    markdown: "Define a small renderer interface (name, supportsGraphShape, supportedArguments, render) and register built-ins (markdown-text, mermaid) in a deterministic lookup map; renderers consume one typed validated renderer-argument set resolved before plugin invocation."
    labels: ["implementation", "design"]
  },
  {
    name: "renderer.registry.resolve"
    title: "Resolve renderer/plugin registry"
    markdown: "Load renderer capabilities, supported arguments, and shape compatibility."
    labels: ["flow", "call"]
  },
  {
    name: "renderer.selection.fallback-policy"
    title: "Use deterministic renderer selection and fallback policy"
    markdown: "Resolve renderer from renderer-scoped arguments sourced from H3Section and notes first, then apply stable defaults by graph shape (Mermaid-first for cyclic graphs, markdown-first for tree/DAG); if cycle-policy is `disallow` and cycles are detected, emit an error diagnostic and skip graph rendering for that section."
    labels: ["implementation", "design"]
  },
  {
    name: "style.errors.guard-clauses"
    title: "Use early returns and guard clauses for errors"
    markdown: "Handle invalid inputs and failure states first, return immediately, and keep the success path shallow and readable."
    labels: ["implementation", "design"]
  },
  {
    name: "style.functions.small-single-purpose"
    title: "Keep functions small and single-purpose"
    markdown: "Each function should do one thing and remain easy to test in isolation; prefer composition of small steps over large multi-branch handlers."
    labels: ["implementation", "design"]
  },
  {
    name: "style.io-separate-from-logic"
    title: "Separate I/O from core logic"
    markdown: "Keep parsing, filtering, and rendering logic pure where possible, and isolate file/network/process I/O behind adapter functions."
    labels: ["implementation", "design"]
  },
  {
    name: "style.parameters.tiny-structs"
    title: "Use tiny structs to avoid long parameter lists"
    markdown: "Group related parameters into small intent-revealing structs (for example, render context and filter options) to reduce call-site ambiguity."
    labels: ["implementation", "design"]
  },
  {
    name: "style.predicates.named"
    title: "Replace boolean soup with named predicates"
    markdown: "Extract compound conditions into well-named predicate helpers to clarify branching and make tests easier to read."
    labels: ["implementation", "design"]
  },
  {
    name: "validate.app.data"
    title: "Validate CUE application data"
    markdown: "Canonical validation pipeline: schema checks, argument registry and free-form argument validation, dataset-based label reference validation, graph integrity policy resolution and graph integrity checks, diagnostic collection, and normalized ValidatedApp output."
    labels: ["flow", "call"]
  },
  {
    name: "validate.cue.schema"
    title: "Validate CUE schema and structure"
    markdown: "Validate required fields, types, and cross-references and attach precise config locations to diagnostics."
    labels: ["flow", "call"]
  },
]


relationships: [
  {
    from: "action.generate.json"
    to: "cli.export.json.graph"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.generate.json"
    to: "export.graph.json"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.generate.json"
    to: "load.app.data"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.generate.json"
    to: "validate.app.data"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.generate.markdown"
    to: "action.generate.markdown.sections"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.generate.markdown"
    to: "cli.config.reports.multiple"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.generate.markdown"
    to: "cli.report.generate"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.generate.markdown"
    to: "load.app.data"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.generate.markdown"
    to: "validate.app.data"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.generate.markdown.section.h3"
    to: "args.coerce.typed"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.generate.markdown.section.h3"
    to: "args.h3.resolve"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.generate.markdown.section.h3"
    to: "args.registry.resolve"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.generate.markdown.section.h3"
    to: "args.validate.runtime"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.generate.markdown.section.h3"
    to: "cli.arguments.free-form"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.generate.markdown.section.h3"
    to: "cli.config.reduce-noise.with-args"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.generate.markdown.section.h3"
    to: "cli.output.deterministic-ordering"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.generate.markdown.section.h3"
    to: "cli.report.generate"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.generate.markdown.section.h3"
    to: "cli.report.subgraph.by-label"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.generate.markdown.section.h3"
    to: "graph.select"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.generate.markdown.section.h3"
    to: "ordering.apply.deterministic"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.generate.markdown.section.h3"
    to: "render.section.file"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.generate.markdown.section.h3"
    to: "render.section.graph"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.generate.markdown.section.h3"
    to: "render.section.orphans"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.generate.markdown.section.h3"
    to: "render.section.plain"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.generate.markdown.sections"
    to: "action.generate.markdown.section.h3"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.generate.markdown.sections"
    to: "cli.note.basic-markdown"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.generate.markdown.sections"
    to: "cli.output.deterministic-ordering"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.generate.markdown.sections"
    to: "cli.output.deterministic-ordering.policy"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.generate.markdown.sections"
    to: "cli.report.generate"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.generate.markdown.sections"
    to: "ordering.policy.resolve"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.lint.names"
    to: "cli.diagnostics.model"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.lint.names"
    to: "cli.diagnostics.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.lint.names"
    to: "cli.names.lint"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.lint.names"
    to: "cli.names.prefix-filter"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.lint.names"
    to: "cli.names.style-policy"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.lint.names"
    to: "cli.output.deterministic-ordering"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.lint.names"
    to: "cli.output.deterministic-ordering.policy"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.lint.names"
    to: "diagnostics.emit.structured"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.lint.names"
    to: "lint.names.notes"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.lint.names"
    to: "lint.names.policy.resolve"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.lint.names"
    to: "lint.names.relationships"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.lint.names"
    to: "load.app.data"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.lint.names"
    to: "names.filter.prefix"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.lint.names"
    to: "ordering.apply.deterministic"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.lint.names"
    to: "ordering.policy.resolve"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.lint.names"
    to: "validate.app.data"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.lint.orphans"
    to: "cli.diagnostics.model"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.lint.orphans"
    to: "cli.diagnostics.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.lint.orphans"
    to: "cli.orphans.lint"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.lint.orphans"
    to: "cli.orphans.query.contextual"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.lint.orphans"
    to: "cli.output.deterministic-ordering"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.lint.orphans"
    to: "cli.output.deterministic-ordering.policy"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.lint.orphans"
    to: "diagnostics.emit.structured"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.lint.orphans"
    to: "lint.orphans.emit"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.lint.orphans"
    to: "lint.orphans.query.resolve"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.lint.orphans"
    to: "load.app.data"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.lint.orphans"
    to: "ordering.apply.deterministic"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.lint.orphans"
    to: "ordering.policy.resolve"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.lint.orphans"
    to: "orphans.query.find"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.lint.orphans"
    to: "validate.app.data"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.list.names"
    to: "cli.names.list"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.list.names"
    to: "cli.names.output-formats"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.list.names"
    to: "cli.names.prefix-filter"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.list.names"
    to: "cli.output.deterministic-ordering"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.list.names"
    to: "cli.output.deterministic-ordering.policy"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.list.names"
    to: "load.app.data"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.list.names"
    to: "names.filter.kind"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.list.names"
    to: "names.filter.prefix"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.list.names"
    to: "names.output.json"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.list.names"
    to: "names.output.table"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.list.names"
    to: "ordering.apply.deterministic"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.list.names"
    to: "ordering.policy.resolve"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.list.names"
    to: "validate.app.data"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.list.reports"
    to: "cli.config.reports.multiple"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.list.reports"
    to: "cli.report.list"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.list.reports"
    to: "list.reports.output"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.list.reports"
    to: "load.app.data"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.list.reports"
    to: "validate.app.data"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.validate"
    to: "cli.config.relationships.labeled"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.validate"
    to: "cli.config.reports.multiple"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.validate"
    to: "cli.diagnostics.model"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.validate"
    to: "cli.diagnostics.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "action.validate"
    to: "diagnostics.emit.structured"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.validate"
    to: "load.app.data"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "action.validate"
    to: "validate.app.data"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "app.model.normalize"
    to: "cli.arguments.registry.schema"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "app.model.normalize"
    to: "cli.config.relationships.labeled"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "app.model.normalize"
    to: "cli.config.reports.multiple"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "app.model.normalize"
    to: "cli.diagnostics.model"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "app.model.normalize"
    to: "cli.graph.integrity.policy"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.coerce.typed"
    to: "cli.arguments.registry.schema"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.coerce.typed"
    to: "cli.arguments.type-coercion"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.h3.resolve"
    to: "cli.arguments.free-form"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.h3.resolve"
    to: "cli.arguments.registry.scope-resolution"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.h3.resolve"
    to: "cli.config.reduce-noise.with-args"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.note.resolve"
    to: "cli.arguments.free-form"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.note.resolve"
    to: "cli.arguments.registry.scope-resolution"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.note.resolve"
    to: "cli.config.reduce-noise.with-args"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.orphan.query.resolve"
    to: "cli.arguments.registry.scope-resolution"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.orphan.query.resolve"
    to: "cli.arguments.runtime-validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.orphan.query.resolve"
    to: "cli.orphans.query.contextual"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.orphan.query.resolve"
    to: "cli.orphans.report.section"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.registry.resolve"
    to: "cli.arguments.registry.schema"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.registry.validate"
    to: "cli.arguments.registry.schema"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.registry.validate"
    to: "cli.diagnostics.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.renderer.resolve"
    to: "cli.arguments.registry.scope-resolution"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.renderer.resolve"
    to: "cli.arguments.runtime-validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.renderer.resolve"
    to: "cli.arguments.type-coercion"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.renderer.resolve"
    to: "cli.renderer.plugin-selection"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.validate.config"
    to: "cli.arguments.free-form"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.validate.config"
    to: "cli.arguments.registry.scope-resolution"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.validate.config"
    to: "cli.arguments.runtime-validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.validate.config"
    to: "cli.diagnostics.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.validate.runtime"
    to: "cli.arguments.registry.schema"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.validate.runtime"
    to: "cli.arguments.registry.scope-resolution"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "args.validate.runtime"
    to: "cli.arguments.runtime-validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "authoring.cue-ai-assistance-gap"
    to: "action.validate"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "authoring.cue-ai-assistance-gap"
    to: "load.app.data"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "authoring.cue-ai-assistance-gap"
    to: "validate.app.data"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "cli.arguments.registry-schema"
    to: "args.coerce.typed"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.registry-schema"
    to: "args.registry.resolve"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.registry-schema"
    to: "args.validate.runtime"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.runtime-validation"
    to: "action.generate.markdown.section.h3"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.runtime-validation"
    to: "args.registry.resolve"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.runtime-validation"
    to: "args.validate.runtime"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.runtime-validation"
    to: "render.section.file.csv"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.runtime-validation"
    to: "render.section.graph"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.scope-resolution"
    to: "args.h3.resolve"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.scope-resolution"
    to: "args.note.resolve"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.scope-resolution"
    to: "args.validate.runtime"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.scope-resolution"
    to: "graph.policy.cycle"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.scope-resolution"
    to: "renderer.plugin.select"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.type-coercion"
    to: "args.coerce.typed"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.type-coercion"
    to: "render.graph.markdown.text"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.type-coercion"
    to: "render.graph.mermaid"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.type-coercion"
    to: "render.section.file.csv"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.typed-models"
    to: "action.generate.markdown.section.h3"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.typed-models"
    to: "action.generate.markdown.sections"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.typed-models"
    to: "args.h3.resolve"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.typed-models"
    to: "args.note.resolve"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.typed-models"
    to: "render.section.file"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.arguments.typed-models"
    to: "render.section.plain"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.cobra"
    to: "action.generate.json"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.cobra"
    to: "action.generate.markdown"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.cobra"
    to: "action.lint.names"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.cobra"
    to: "action.list.names"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.cobra"
    to: "action.list.reports"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.cobra"
    to: "action.validate"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.cobra"
    to: "cli.root"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "cli.root"
    to: "action.generate.json"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "cli.root"
    to: "action.generate.markdown"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "cli.root"
    to: "action.lint.names"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "cli.root"
    to: "action.lint.orphans"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "cli.root"
    to: "action.list.names"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "cli.root"
    to: "action.list.reports"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "cli.root"
    to: "action.validate"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "cli.root"
    to: "cli.config.reports.multiple"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "cli.root"
    to: "cli.report.generate"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "cli.root"
    to: "cli.report.list"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "collaboration.cue-merge-conflicts"
    to: "action.validate"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "collaboration.cue-merge-conflicts"
    to: "load.app.data"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "collaboration.cue-merge-conflicts"
    to: "validate.app.data"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "config.arguments.reduce-noise"
    to: "action.generate.markdown.section.h3"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "config.arguments.reduce-noise"
    to: "args.h3.resolve"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "config.arguments.reduce-noise"
    to: "args.note.resolve"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "config.arguments.reduce-noise"
    to: "render.section.file"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "config.cue"
    to: "action.validate"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "config.cue"
    to: "load.app.data"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "config.cue"
    to: "validate.app.data"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "diagnostics.collect.validation"
    to: "cli.diagnostics.model"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "diagnostics.collect.validation"
    to: "cli.diagnostics.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "diagnostics.emit.structured"
    to: "cli.diagnostics.model"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "diagnostics.emit.structured"
    to: "cli.diagnostics.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "diagnostics.structured-model"
    to: "args.validate.runtime"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "diagnostics.structured-model"
    to: "diagnostics.emit.structured"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "diagnostics.structured-model"
    to: "validate.app.data"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "diagnostics.validation-location"
    to: "action.validate"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "diagnostics.validation-location"
    to: "diagnostics.emit.structured"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "diagnostics.validation-location"
    to: "validate.app.data"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "export.graph.json"
    to: "cli.export.json.graph"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "file.csv.filter"
    to: "cli.arguments.free-form"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "file.csv.filter"
    to: "cli.note.csv.filter-column"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "graph.circular-dependency"
    to: "action.generate.markdown.section.h3"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "graph.circular-dependency"
    to: "action.validate"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "graph.circular-dependency"
    to: "graph.select"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "graph.circular-dependency"
    to: "render.section.graph"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "graph.integrity.check.cross-report-references"
    to: "cli.graph.integrity.policy"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "graph.integrity.check.cross-report-references"
    to: "cli.graph.integrity.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "graph.integrity.check.duplicate-note-names"
    to: "cli.graph.integrity.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "graph.integrity.check.missing-nodes"
    to: "cli.graph.integrity.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "graph.integrity.check.orphans"
    to: "cli.graph.integrity.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "graph.integrity.policy-model"
    to: "graph.integrity.policy.resolve"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "graph.integrity.policy-model"
    to: "graph.integrity.validate"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "graph.integrity.policy-model"
    to: "validate.app.data"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "graph.integrity.policy.resolve"
    to: "cli.graph.integrity.policy"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "graph.integrity.validate"
    to: "cli.diagnostics.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "graph.integrity.validate"
    to: "cli.graph.integrity.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "graph.integrity.validate"
    to: "graph.integrity.check.cross-report-references"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "graph.integrity.validate"
    to: "graph.integrity.check.duplicate-note-names"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "graph.integrity.validate"
    to: "graph.integrity.check.missing-nodes"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "graph.integrity.validate"
    to: "graph.integrity.check.orphans"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "graph.integrity.validation-engine"
    to: "diagnostics.emit.structured"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "graph.integrity.validation-engine"
    to: "graph.integrity.check.cross-report-references"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "graph.integrity.validation-engine"
    to: "graph.integrity.check.duplicate-note-names"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "graph.integrity.validation-engine"
    to: "graph.integrity.check.missing-nodes"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "graph.integrity.validation-engine"
    to: "graph.integrity.check.orphans"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "graph.integrity.validation-engine"
    to: "graph.integrity.validate"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "graph.integrity.validation-engine"
    to: "labels.dataset.collect"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "graph.integrity.validation-engine"
    to: "labels.reference.validate"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "graph.policy.cycle"
    to: "cli.arguments.registry.scope-resolution"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "graph.policy.cycle"
    to: "cli.section.h3.cycle-policy"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "graph.select"
    to: "cli.config.relationships.labeled"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "graph.select"
    to: "cli.report.subgraph.by-label"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "graph.shape.detect"
    to: "cli.report.graph.shape-aware-render"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "graph.shape.detect"
    to: "cli.section.h3.cycle-policy"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "labels.dataset.collect"
    to: "cli.config.relationships.labeled"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "labels.dataset.collect"
    to: "cli.diagnostics.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "labels.reference.validate"
    to: "cli.arguments.registry.scope-resolution"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "labels.reference.validate"
    to: "cli.config.relationships.labeled"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "labels.reference.validate"
    to: "cli.diagnostics.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "labels.reference.validate"
    to: "cli.report.subgraph.by-label"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "lang.go"
    to: "cli.root"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "lint.names.notes"
    to: "cli.diagnostics.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "lint.names.notes"
    to: "cli.names.lint"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "lint.names.notes"
    to: "cli.names.style-policy"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "lint.names.policy.resolve"
    to: "cli.names.lint"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "lint.names.policy.resolve"
    to: "cli.names.style-policy"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "lint.names.relationships"
    to: "cli.diagnostics.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "lint.names.relationships"
    to: "cli.names.lint"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "lint.names.relationships"
    to: "cli.names.style-policy"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "lint.orphans.emit"
    to: "cli.diagnostics.model"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "lint.orphans.emit"
    to: "cli.diagnostics.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "lint.orphans.emit"
    to: "cli.orphans.lint"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "lint.orphans.query.resolve"
    to: "cli.orphans.lint"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "lint.orphans.query.resolve"
    to: "cli.orphans.query.contextual"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "list.reports.output"
    to: "cli.config.reports.multiple"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "list.reports.output"
    to: "cli.report.list"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "load.app.data"
    to: "cli.config.relationships.labeled"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "load.app.data"
    to: "cli.config.reports.multiple"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "load.app.data"
    to: "cli.export.json.graph"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "load.app.data"
    to: "cli.report.generate"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "load.app.data"
    to: "cli.report.list"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "maintenance.single-cue-file-size"
    to: "action.validate"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "maintenance.single-cue-file-size"
    to: "load.app.data"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "maintenance.single-cue-file-size"
    to: "validate.app.data"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "names.filter.kind"
    to: "cli.names.list"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "names.filter.kind"
    to: "cli.names.output-formats"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "names.filter.prefix"
    to: "cli.names.lint"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "names.filter.prefix"
    to: "cli.names.list"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "names.filter.prefix"
    to: "cli.names.prefix-filter"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "names.lint.command"
    to: "action.lint.names"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.lint.command"
    to: "diagnostics.emit.structured"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.lint.command"
    to: "lint.names.notes"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.lint.command"
    to: "lint.names.policy.resolve"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.lint.command"
    to: "lint.names.relationships"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.lint.command"
    to: "load.app.data"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.lint.command"
    to: "names.filter.prefix"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.lint.command"
    to: "ordering.apply.deterministic"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.lint.command"
    to: "ordering.policy.resolve"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.lint.command"
    to: "validate.app.data"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.list.command"
    to: "action.list.names"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.list.command"
    to: "load.app.data"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.list.command"
    to: "names.filter.kind"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.list.command"
    to: "names.filter.prefix"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.list.command"
    to: "names.output.json"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.list.command"
    to: "names.output.table"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.list.command"
    to: "ordering.apply.deterministic"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.list.command"
    to: "ordering.policy.resolve"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.list.command"
    to: "validate.app.data"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "names.output.json"
    to: "cli.names.list"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "names.output.json"
    to: "cli.names.output-formats"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "names.output.table"
    to: "cli.names.list"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "names.output.table"
    to: "cli.names.output-formats"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "ordering.apply.deterministic"
    to: "cli.output.deterministic-ordering"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "ordering.apply.deterministic"
    to: "cli.output.deterministic-ordering.policy"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "ordering.policy.resolve"
    to: "cli.output.deterministic-ordering.policy"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "orphans.lint.command"
    to: "action.lint.orphans"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "orphans.lint.command"
    to: "diagnostics.emit.structured"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "orphans.lint.command"
    to: "lint.orphans.emit"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "orphans.lint.command"
    to: "lint.orphans.query.resolve"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "orphans.lint.command"
    to: "load.app.data"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "orphans.lint.command"
    to: "ordering.apply.deterministic"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "orphans.lint.command"
    to: "ordering.policy.resolve"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "orphans.lint.command"
    to: "orphans.query.find"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "orphans.lint.command"
    to: "validate.app.data"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "orphans.query.find"
    to: "cli.orphans.lint"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "orphans.query.find"
    to: "cli.orphans.query.contextual"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "orphans.query.find"
    to: "cli.orphans.report.section"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "orphans.query.find"
    to: "cli.output.deterministic-ordering"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "orphans.render.rows"
    to: "cli.orphans.report.section"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "orphans.render.rows"
    to: "cli.output.deterministic-ordering"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "orphans.report.section"
    to: "args.orphan.query.resolve"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "orphans.report.section"
    to: "ordering.apply.deterministic"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "orphans.report.section"
    to: "orphans.query.find"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "orphans.report.section"
    to: "orphans.render.rows"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "orphans.report.section"
    to: "render.section.orphans"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "output.ordering.deterministic"
    to: "action.generate.markdown.section.h3"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "output.ordering.deterministic"
    to: "action.generate.markdown.sections"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "output.ordering.deterministic"
    to: "ordering.apply.deterministic"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "output.ordering.deterministic"
    to: "ordering.policy.resolve"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "output.ordering.policy-contract"
    to: "action.generate.markdown"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "output.ordering.policy-contract"
    to: "ordering.apply.deterministic"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "output.ordering.policy-contract"
    to: "ordering.policy.resolve"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "performance.report-generation-scale"
    to: "action.generate.markdown"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "performance.report-generation-scale"
    to: "action.generate.markdown.section.h3"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "performance.report-generation-scale"
    to: "action.generate.markdown.sections"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "performance.report-generation-scale"
    to: "graph.select"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "performance.report-generation-scale"
    to: "render.section.file"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "performance.report-generation-scale"
    to: "render.section.file.csv"
    label: "affects_call"
    labels: ["affects_call"]
  },
  {
    from: "render.graph.circular"
    to: "cli.report.graph.renderer.markdown-text"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.graph.circular"
    to: "cli.report.graph.renderer.mermaid"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.graph.circular"
    to: "cli.report.graph.shape-aware-render"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.graph.circular"
    to: "cli.section.h3.cycle-policy"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.graph.circular"
    to: "render.graph.markdown.text"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.graph.circular"
    to: "render.graph.mermaid"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.graph.markdown.text"
    to: "cli.renderer.registry"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.graph.markdown.text"
    to: "cli.report.graph.renderer.markdown-text"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.graph.mermaid"
    to: "cli.note.mermaid.embed"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.graph.mermaid"
    to: "cli.renderer.registry"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.graph.mermaid"
    to: "cli.report.graph.renderer.mermaid"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.graph.tree-or-dag"
    to: "cli.report.graph.renderer.markdown-text"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.graph.tree-or-dag"
    to: "cli.report.graph.renderer.mermaid"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.graph.tree-or-dag"
    to: "cli.report.graph.shape-aware-render"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.graph.tree-or-dag"
    to: "render.graph.markdown.text"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.graph.tree-or-dag"
    to: "render.graph.mermaid"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.file"
    to: "args.coerce.typed"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.file"
    to: "args.note.resolve"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.file"
    to: "args.registry.resolve"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.file"
    to: "args.validate.runtime"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.file"
    to: "cli.arguments.free-form"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.file"
    to: "cli.note.filepath.reference"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.file"
    to: "cli.output.deterministic-ordering"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.file"
    to: "render.section.file.code"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.file"
    to: "render.section.file.csv"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.file"
    to: "render.section.file.media"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.file.code"
    to: "cli.arguments.free-form"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.file.code"
    to: "cli.arguments.runtime-validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.file.code"
    to: "cli.note.filepath.reference"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.file.code"
    to: "cli.note.mermaid.embed"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.file.csv"
    to: "cli.arguments.free-form"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.file.csv"
    to: "cli.arguments.runtime-validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.file.csv"
    to: "cli.note.csv.embed"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.file.csv"
    to: "cli.note.filepath.reference"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.file.csv"
    to: "file.csv.filter"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.file.media"
    to: "cli.arguments.free-form"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.file.media"
    to: "cli.note.filepath.reference"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.file.media"
    to: "cli.note.image.preview"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.graph"
    to: "args.renderer.resolve"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.graph"
    to: "cli.config.relationships.labeled"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.graph"
    to: "cli.output.deterministic-ordering"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.graph"
    to: "cli.renderer.plugin-selection"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.graph"
    to: "cli.renderer.registry"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.graph"
    to: "cli.report.graph.shape-aware-render"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.graph"
    to: "cli.report.subgraph.by-label"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.graph"
    to: "cli.section.h3.cycle-policy"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.graph"
    to: "graph.policy.cycle"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.graph"
    to: "graph.shape.detect"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.graph"
    to: "render.graph.circular"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.graph"
    to: "render.graph.tree-or-dag"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.graph"
    to: "renderer.plugin.select"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.graph"
    to: "renderer.registry.resolve"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.orphans"
    to: "args.orphan.query.resolve"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.orphans"
    to: "cli.orphans.query.contextual"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.orphans"
    to: "cli.orphans.report.section"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.orphans"
    to: "cli.output.deterministic-ordering"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.orphans"
    to: "orphans.query.find"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.orphans"
    to: "orphans.render.rows"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "render.section.plain"
    to: "cli.arguments.free-form"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.plain"
    to: "cli.arguments.runtime-validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.plain"
    to: "cli.note.basic-markdown"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "render.section.plain"
    to: "cli.note.link.markdown"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "renderer.plugin.select"
    to: "cli.arguments.free-form"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "renderer.plugin.select"
    to: "cli.renderer.plugin-selection"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "renderer.registry.contract"
    to: "render.graph.markdown.text"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "renderer.registry.contract"
    to: "render.graph.mermaid"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "renderer.registry.contract"
    to: "render.section.graph"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "renderer.registry.contract"
    to: "renderer.plugin.select"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "renderer.registry.contract"
    to: "renderer.registry.resolve"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "renderer.registry.resolve"
    to: "cli.renderer.registry"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "renderer.selection.fallback-policy"
    to: "graph.shape.detect"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "renderer.selection.fallback-policy"
    to: "render.graph.circular"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "renderer.selection.fallback-policy"
    to: "render.graph.tree-or-dag"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "renderer.selection.fallback-policy"
    to: "renderer.plugin.select"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.errors.guard-clauses"
    to: "action.generate.json"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.errors.guard-clauses"
    to: "action.generate.markdown"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.errors.guard-clauses"
    to: "action.validate"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.errors.guard-clauses"
    to: "validate.app.data"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.functions.small-single-purpose"
    to: "action.generate.markdown.section.h3"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.functions.small-single-purpose"
    to: "action.generate.markdown.sections"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.functions.small-single-purpose"
    to: "render.section.file"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.io-separate-from-logic"
    to: "action.generate.markdown"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.io-separate-from-logic"
    to: "load.app.data"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.io-separate-from-logic"
    to: "render.section.file.csv"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.io-separate-from-logic"
    to: "render.section.file.media"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.io-separate-from-logic"
    to: "validate.app.data"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.parameters.tiny-structs"
    to: "action.generate.markdown.section.h3"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.parameters.tiny-structs"
    to: "file.csv.filter"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.parameters.tiny-structs"
    to: "graph.select"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.parameters.tiny-structs"
    to: "render.section.file"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.predicates.named"
    to: "file.csv.filter"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.predicates.named"
    to: "graph.select"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.predicates.named"
    to: "render.section.file"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "style.predicates.named"
    to: "validate.app.data"
    label: "applies_to_call"
    labels: ["applies_to_call"]
  },
  {
    from: "validate.app.data"
    to: "app.model.normalize"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "validate.app.data"
    to: "args.registry.resolve"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "validate.app.data"
    to: "args.registry.validate"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "validate.app.data"
    to: "args.validate.config"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "validate.app.data"
    to: "cli.arguments.free-form"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "validate.app.data"
    to: "cli.arguments.registry.schema"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "validate.app.data"
    to: "cli.arguments.runtime-validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "validate.app.data"
    to: "cli.config.relationships.labeled"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "validate.app.data"
    to: "cli.config.reports.multiple"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "validate.app.data"
    to: "cli.diagnostics.model"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "validate.app.data"
    to: "cli.diagnostics.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "validate.app.data"
    to: "cli.graph.integrity.policy"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "validate.app.data"
    to: "cli.graph.integrity.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "validate.app.data"
    to: "diagnostics.collect.validation"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "validate.app.data"
    to: "graph.integrity.policy.resolve"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "validate.app.data"
    to: "graph.integrity.validate"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "validate.app.data"
    to: "labels.dataset.collect"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "validate.app.data"
    to: "labels.reference.validate"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "validate.app.data"
    to: "validate.cue.schema"
    label: "contains_step"
    labels: ["contains_step"]
  },
  {
    from: "validate.cue.schema"
    to: "cli.config.reports.multiple"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
  {
    from: "validate.cue.schema"
    to: "cli.diagnostics.validation"
    label: "satisfies_usecase"
    labels: ["satisfies_usecase"]
  },
]
