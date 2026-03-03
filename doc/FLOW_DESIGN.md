# FLOW DESIGN OVERVIEW (Generated)

## Function calls tree

```
flyb CLI root command [cli.root]
  Generate the markdown reports [action.generate.markdown]
    Load CLUE application data [load.app.data]
    Validate CLUE application data [validate.app.data]
    Generate the markdown sections [action.generate.markdown.sections]
      Generate the H3 section [action.generate.markdown.section.h3]
  Generate as json [action.generate.json]
    Load CLUE application data [load.app.data]
    Validate CLUE application data [validate.app.data]
  Validate the CUE file [action.validate]
    Load CLUE application data [load.app.data]
    Validate CLUE application data [validate.app.data]
```

Supported use cases:

  - Generate design reports from configured notes and relationships — This is the primary end-to-end report generation use case.


Unsupported use cases (yet):

  - Render note title and markdown description — Each note includes a concise title with free-form markdown content.
  - Reference a file from a note — Referenced files can be embedded in generated markdown output.
  - Embed CSV content from a referenced file — CSV input can render as a markdown table or as raw CSV.
  - Filter embedded CSV rows by column — Column filters reduce CSV output to the relevant subset.
  - Preview referenced image files in markdown — Image references render as embedded previews in reports.
  - Embed Mermaid diagrams from file content — Mermaid content is emitted in fenced blocks for diagram rendering.
  - Convert note links to markdown links — URL links are rendered with link text in markdown output.
  - Define labeled relationships between notes in config — CUE can be used as the source format for flexible configuration.
  - Declare multiple markdown reports in one config — A single config can drive generation of multiple report files.
  - List all configured markdown reports — The CLI exposes a command to enumerate report targets.
  - Export notes and relationships as JSON — JSON export supports machine-readable graph processing.
  - Build a report from a relationship-label subgraph — Report generation can include only edges matching selected labels.


