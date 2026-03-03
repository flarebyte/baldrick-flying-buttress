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
      - note: Ensure required fields and cross-reference integrity are valid.
  Generate markdown reports [action.generate.markdown]
    - note: Renders one or more markdown outputs from the validated config.
    Load CUE application data [load.app.data]
      - note: Read notes, relationships, and report definitions from config.
    Validate CUE application data [validate.app.data]
      - note: Ensure required fields and cross-reference integrity are valid.
    Generate markdown sections [action.generate.markdown.sections]
      - note: Build H3 sections from note subsets and renderers.
      Generate a single H3 section [action.generate.markdown.section.h3]
        - note: Compose subgraph, plain content, and file-backed content for one section.
        Extract subgraph using labels [graph.select]
          - note: Filter notes and relationships by labels and optional starting node.
        Render section as a graph [render.section.graph]
          - note: Renders graph-focused content from filtered notes and relationships.
        Render plain section [render.section.plain]
          - note: Render title and markdown body, including markdown links.
        Render section with referenced file content [render.section.file]
          - note: Dispatches file rendering by type (CSV, media, code/diagram).
          Render section with CSV file [render.section.file.csv]
            - note: Render as a markdown table or raw CSV code block.
            Filter CSV rows by column [file.csv.filter]
              - note: Apply include/exclude filters before rendering CSV output.
          Render section with media file [render.section.file.media]
            - note: Embed image previews for supported media types.
          Render section with code or Mermaid snippet [render.section.file.code]
            - note: Preserve fenced-block formatting for code and Mermaid content.
  Generate JSON graph export [action.generate.json]
    - note: Export notes and relationships in machine-readable JSON format.
    Load CUE application data [load.app.data]
      - note: Read notes, relationships, and report definitions from config.
    Validate CUE application data [validate.app.data]
      - note: Ensure required fields and cross-reference integrity are valid.
  Validate the CUE file [action.validate]
    - note: Validate configuration structure and constraints without rendering output.
    Load CUE application data [load.app.data]
      - note: Read notes, relationships, and report definitions from config.
    Validate CUE application data [validate.app.data]
      - note: Ensure required fields and cross-reference integrity are valid.
```

Supported use cases:

  - Generate design reports from configured notes and relationships — This is the primary end-to-end report generation use case.
  - List all configured markdown reports — The CLI exposes a command to enumerate report targets.
  - Declare multiple markdown reports in one config — A single config can drive generation of multiple report files.
  - Export notes and relationships as JSON — JSON export supports machine-readable graph processing.
  - Define labeled relationships between notes in config — CUE can be used as the source format for flexible configuration.
  - Render note title and markdown description — Each note includes a concise title with free-form markdown content.
  - Build a report from a relationship-label subgraph — Report generation can include only edges matching selected labels.
  - Convert note links to markdown links — URL links are rendered with link text in markdown output.
  - Reference a file from a note — Referenced files can be embedded in generated markdown output.
  - Embed CSV content from a referenced file — CSV input can render as a markdown table or as raw CSV.
  - Filter embedded CSV rows by column — Column filters reduce CSV output to the relevant subset.
  - Preview referenced image files in markdown — Image references render as embedded previews in reports.
  - Embed Mermaid diagrams from file content — Mermaid content is emitted in fenced blocks for diagram rendering.


