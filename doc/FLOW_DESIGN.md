# FLOW DESIGN OVERVIEW (Generated)

## Function calls tree

```
flyb CLI root command [cli.root]
  - note: Entry point for report generation, listing, JSON export, and config validation.
  List configured markdown reports [action.list.reports]
    - note: Enumerate report targets without generating files.
    Load CUE application data [load.app.data]
      - note: Read notes, relationships, and report definitions from config.
    Validate CUE application data [validate.app.data]
      - note: Ensure required fields and cross-reference integrity are valid and diagnostics are attached to config locations.
  Generate markdown reports [action.generate.markdown]
    - note: Renders one or more markdown outputs from the validated config.
    Load CUE application data [load.app.data]
      - note: Read notes, relationships, and report definitions from config.
    Validate CUE application data [validate.app.data]
      - note: Ensure required fields and cross-reference integrity are valid and diagnostics are attached to config locations.
    Generate markdown sections [action.generate.markdown.sections]
      - note: Build H3 sections from note subsets and renderers with deterministic ordering.
      Resolve deterministic ordering policy [ordering.policy.resolve]
        - note: Resolve stable ordering rules for notes, relationships, sections, and arguments.
      Generate a single H3 section [action.generate.markdown.section.h3]
        - note: Compose subgraph, plain content, and file-backed content with section-level arguments.
        Resolve H3Section free-form arguments [args.h3.resolve]
          - note: Read flexible section arguments as key/value flags (for example `graph-renderer=mermaid`).
        Resolve argument registry schema [args.registry.resolve]
          - note: Load known argument definitions (type, default, allowed values, scopes).
        Validate arguments at runtime [args.validate.runtime]
          - note: Validate keys and values against a known argument registry and fail fast on invalid input.
        Coerce arguments to typed values [args.coerce.typed]
          - note: Coerce validated values to target types (string[], boolean, enum, number).
        Extract subgraph using labels [graph.select]
          - note: Filter notes and relationships by labels and optional starting node.
        Render section as a graph [render.section.graph]
          - note: Resolve cycle policy and graph shape, then render with selected renderer(s).
          Resolve renderer/plugin registry [renderer.registry.resolve]
            - note: Load renderer capabilities, supported arguments, and shape compatibility.
          Select renderer plugin from arguments [renderer.plugin.select]
            - note: Choose renderer by arguments with deterministic fallback when unspecified.
          Resolve H3Section cycle policy argument [graph.policy.cycle]
            - note: Use section argument to disallow, allow, or collapse cycles.
          Detect graph shape (tree, DAG, or cyclic) [graph.shape.detect]
            - note: Classify graph structure before selecting rendering strategy.
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
            - note: Read note-level rendering options as key/value flags (for example `format-csv=md`).
          Resolve argument registry schema [args.registry.resolve]
            - note: Load known argument definitions (type, default, allowed values, scopes).
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
          - note: Sort entities and edges with stable tie-breakers before rendering output.
  Generate JSON graph export [action.generate.json]
    - note: Export notes and relationships in machine-readable JSON format.
    Load CUE application data [load.app.data]
      - note: Read notes, relationships, and report definitions from config.
    Validate CUE application data [validate.app.data]
      - note: Ensure required fields and cross-reference integrity are valid and diagnostics are attached to config locations.
  Validate the CUE file [action.validate]
    - note: Validate configuration structure and constraints and emit structured diagnostics.
    Load CUE application data [load.app.data]
      - note: Read notes, relationships, and report definitions from config.
    Validate CUE application data [validate.app.data]
      - note: Ensure required fields and cross-reference integrity are valid and diagnostics are attached to config locations.
    Emit structured diagnostics [diagnostics.emit.structured]
      - note: Emit diagnostics with code, severity, source, message, and optional location.
```

Supported use cases:

  - Generate design reports from configured notes and relationships — This is the primary end-to-end report generation use case.
  - List all configured markdown reports — The CLI exposes a command to enumerate report targets.
  - Declare multiple markdown reports in one config — A single config can drive generation of multiple report files.
  - Export notes and relationships as JSON — JSON export supports machine-readable graph processing.
  - Define labeled relationships between notes in config — CUE can be used as the source format for flexible configuration.
  - Emit structured diagnostics — Diagnostics include code, severity, message, source, and optional location context.
  - Report validation diagnostics with locations — Validation errors and warnings should point to config paths and offending argument or relationship names.
  - Render note title and markdown description — Each note includes a concise title with free-form markdown content.
  - Guarantee deterministic output ordering — Sort notes, relationships, sections, and arguments with stable rules so repeated runs produce identical output.
  - Define an explicit ordering policy — Ordering policy is part of runtime behavior and can be documented/tested as a contract.
  - Build a report from a relationship-label subgraph — Report generation can include only edges matching selected labels.
  - Accept free-form arguments on H3Section and Note — Arguments behave like CLI flags (for example `format-csv=md`) and can carry string, string[], boolean, and similar values.
  - Keep CUE config compact with argument-driven rendering options — Prefer small composable argument lists over proliferating specialized configuration fields.
  - Resolve arguments by scope — Apply argument rules by scope (global, h2, h3, note, renderer) to prevent invalid combinations.
  - Define an argument registry schema — Registry entries define argument key, type, default, allowed values, and valid scopes.
  - Validate free-form arguments at runtime — Validate against a known argument registry and fail with clear errors on unknown keys or invalid values.
  - Coerce free-form argument values into typed values — Convert string-like argument inputs into validated typed values before rendering.
  - Allow each H3 section to define cycle policy — H3Section arguments can declare whether cycles are disallowed, allowed, or collapsed.
  - Render graph output based on graph shape — Renderer behavior adapts to tree, DAG, and cyclic graph structures.
  - Register renderers and plugins in a capability registry — A renderer registry maps renderer names to capabilities, supported arguments, and graph-shape compatibility.
  - Select renderer plugin from arguments at runtime — Renderer selection is resolved from section and note arguments with fallback defaults.
  - Render graph output as markdown text — Text rendering supports readable hierarchy and edge summaries in markdown reports.
  - Render graph output as Mermaid diagram — Mermaid output supports visual graph rendering, including cyclic relationships.
  - Embed Mermaid diagrams from file content — Mermaid content is emitted in fenced blocks for diagram rendering.
  - Convert note links to markdown links — URL links are rendered with link text in markdown output.
  - Reference a file from a note — Referenced files can be embedded in generated markdown output.
  - Embed CSV content from a referenced file — CSV input can render as a markdown table or as raw CSV.
  - Filter embedded CSV rows by column — Column filters reduce CSV output to the relevant subset.
  - Preview referenced image files in markdown — Image references render as embedded previews in reports.


