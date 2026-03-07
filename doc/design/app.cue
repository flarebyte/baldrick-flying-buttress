package flyb

source: "design-migration"
name:   "design-migration"
modules: ["design", "documentation"]

reports: [
  {
    title: "Flow Design Overview"
    filepath: "doc/FLOW_DESIGN.md"
    description: "Migrated flow call tree and validation pipeline design notes."
    sections: [{
      title: "Flow"
      description: "Flow call notes."
      sections: [{
        title: "Flow Calls"
        description: "Generated from design/flows.ts."
        notes: ["action.generate.json", "action.generate.markdown", "action.generate.markdown.section.h3", "action.generate.markdown.sections", "action.lint.names", "action.lint.orphans", "action.list.names", "action.list.reports", "action.validate", "app.model.normalize", "args.coerce.typed", "args.h3.resolve", "args.note.resolve", "args.orphan.query.resolve", "args.registry.resolve", "args.registry.validate", "args.renderer.resolve", "args.validate.config", "args.validate.runtime", "cli.root", "diagnostics.collect.validation", "diagnostics.emit.structured", "export.graph.json", "file.csv.filter", "graph.integrity.check.cross-report-references", "graph.integrity.check.duplicate-note-names", "graph.integrity.check.missing-nodes", "graph.integrity.check.orphans", "graph.integrity.policy.resolve", "graph.integrity.validate", "graph.policy.cycle", "graph.select", "graph.shape.detect", "labels.dataset.collect", "labels.reference.validate", "lint.names.notes", "lint.names.policy.resolve", "lint.names.relationships", "lint.orphans.emit", "lint.orphans.query.resolve", "list.reports.output", "load.app.data", "names.filter.kind", "names.filter.prefix", "names.output.json", "names.output.table", "ordering.apply.deterministic", "ordering.policy.resolve", "orphans.query.find", "orphans.render.rows", "render.graph.circular", "render.graph.markdown.text", "render.graph.mermaid", "render.graph.tree-or-dag", "render.section.file", "render.section.file.code", "render.section.file.csv", "render.section.file.media", "render.section.graph", "render.section.orphans", "render.section.plain", "renderer.plugin.select", "renderer.registry.resolve", "validate.app.data", "validate.cue.schema"]
      }]
    }]
  },
  {
    title: "Use Cases"
    filepath: "doc/USE_CASES.md"
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
    filepath: "doc/IMPLEMENTATIONS.md"
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
    filepath: "doc/RISKS.md"
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
