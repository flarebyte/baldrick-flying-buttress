package flyb

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
