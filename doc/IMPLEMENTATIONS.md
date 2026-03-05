# Implementation Considerations (Generated)

This document summarizes suggested implementation choices.

## Summary
- Define a typed argument registry schema [cli.arguments.registry-schema]
- Validate free-form arguments with Cobra-style validators [cli.arguments.runtime-validation]
- Apply scope-aware argument resolution [cli.arguments.scope-resolution]
- Use explicit type coercion for free-form arguments [cli.arguments.type-coercion]
- Use free-form key/value arguments with typed coercion [cli.arguments.typed-models]
- Use Cobra for CLI command and argument parsing [cli.cobra]
- Use arguments to reduce CUE configuration noise [config.arguments.reduce-noise]
- Use CUE as the configuration source of truth [config.cue]
- Use a structured diagnostics model [diagnostics.structured-model]
- Attach validation diagnostics to precise config locations [diagnostics.validation-location]
- Define a graph integrity policy model [graph.integrity.policy-model]
- Implement graph integrity validation checks [graph.integrity.validation-engine]
- Implement the CLI in Go [lang.go]
- Guarantee deterministic ordering in generated outputs [output.ordering.deterministic]
- Treat ordering policy as a testable contract [output.ordering.policy-contract]
- Define a renderer plugin registry contract [renderer.registry.contract]
- Use deterministic renderer selection and fallback policy [renderer.selection.fallback-policy]
- Use early returns and guard clauses for errors [style.errors.guard-clauses]
- Keep functions small and single-purpose [style.functions.small-single-purpose]
- Separate I/O from core logic [style.io-separate-from-logic]
- Use tiny structs to avoid long parameter lists [style.parameters.tiny-structs]
- Replace boolean soup with named predicates [style.predicates.named]

## Define a typed argument registry schema [cli.arguments.registry-schema]

- Description: Maintain a registry of argument definitions (name, type, default, allowed values, scopes) and use it as the single source of truth for argument behavior; valid scopes are `h3-section`, `note`, and `renderer`.
- Calls: args.registry.resolve, args.validate.runtime, args.coerce.typed

## Validate free-form arguments with Cobra-style validators [cli.arguments.runtime-validation]

- Description: Validate argument keys and values at runtime against a known registry (enums, booleans, repeated values) using familiar CLI validation patterns and clear error messages.
- Calls: args.validate.runtime, args.registry.resolve, action.generate.markdown.section.h3, render.section.file.csv, render.section.graph

## Apply scope-aware argument resolution [cli.arguments.scope-resolution]

- Description: Resolve and validate arguments by scope (h3-section, note, renderer) so options are accepted only where they are meaningful; renderer-scoped arguments are collected from H3Section and note argument lists with precedence `note` > `h3-section` > registry defaults.
- Calls: args.h3.resolve, args.note.resolve, args.validate.runtime, graph.policy.cycle, renderer.plugin.select

## Use explicit type coercion for free-form arguments [cli.arguments.type-coercion]

- Description: After validation, coerce argument values to typed runtime options before rendering to avoid stringly-typed behavior in renderers.
- Calls: args.coerce.typed, render.section.file.csv, render.graph.mermaid, render.graph.markdown.text

## Use free-form key/value arguments with typed coercion [cli.arguments.typed-models]

- Description: Treat H3Section and Note arguments like CLI-style flags (for example `format-csv=md`) to keep config flexible, then coerce values into typed runtime options per renderer.
- Calls: action.generate.markdown.sections, action.generate.markdown.section.h3, args.h3.resolve, args.note.resolve, render.section.plain, render.section.file

## Use Cobra for CLI command and argument parsing [cli.cobra]

- Description: Use Cobra to model commands, flags, and subcommands (`generate markdown`, `generate json`, `validate`, `list`) with a consistent command tree.
- Calls: cli.root, action.generate.markdown, action.generate.json, action.validate, action.list.reports

## Use arguments to reduce CUE configuration noise [config.arguments.reduce-noise]

- Description: Prefer composable argument lists over adding many specialized CUE fields, so rendering capabilities can evolve without large schema churn.
- Calls: args.h3.resolve, args.note.resolve, action.generate.markdown.section.h3, render.section.file

## Use CUE as the configuration source of truth [config.cue]

- Description: Represent notes, relationships, and report definitions in CUE for schema validation, defaults, and composable configuration.
- Calls: load.app.data, validate.app.data, action.validate

## Use a structured diagnostics model [diagnostics.structured-model]

- Description: Standardize diagnostics with code, severity, source, message, and location to support CLI UX, CI checks, and future editor integrations.
- Calls: validate.app.data, args.validate.runtime, diagnostics.emit.structured

## Attach validation diagnostics to precise config locations [diagnostics.validation-location]

- Description: Include CUE path, related note/relationship, and argument key in diagnostics so users can quickly fix invalid configuration.
- Calls: validate.app.data, action.validate, diagnostics.emit.structured

## Define a graph integrity policy model [graph.integrity.policy-model]

- Description: Define explicit integrity rules for missing nodes, orphan nodes, duplicate names, and cross-report references with per-rule severity; validate label references separately against dataset-derived labels.
- Calls: graph.integrity.policy.resolve, graph.integrity.validate, validate.app.data

## Implement graph integrity validation checks [graph.integrity.validation-engine]

- Description: Run focused integrity checks and emit structured diagnostics linked to note names, relationships, arguments, and CUE paths; keep label-definition handling free-form and validate only label references.
- Calls: graph.integrity.validate, graph.integrity.check.missing-nodes, graph.integrity.check.orphans, graph.integrity.check.duplicate-note-names, graph.integrity.check.cross-report-references, labels.dataset.collect, labels.reference.validate, diagnostics.emit.structured

## Implement the CLI in Go [lang.go]

- Description: Use Go as the primary implementation language for strong typing, fast startup, and straightforward single-binary distribution.
- Calls: cli.root

## Guarantee deterministic ordering in generated outputs [output.ordering.deterministic]

- Description: Apply explicit stable sorting for notes, relationships, sections, and arguments using concrete comparators (notes: primaryLabel/name, relationships: from/to/labelsSortedJoined, sections: case-insensitive title plus originalIndex, arguments: name) so output remains reproducible across runs and machines.
- Calls: ordering.policy.resolve, ordering.apply.deterministic, action.generate.markdown.sections, action.generate.markdown.section.h3

## Treat ordering policy as a testable contract [output.ordering.policy-contract]

- Description: Define ordering rules and tie-breakers as a versioned policy (including label normalization and relationship label joining rules) and verify them with golden-file tests.
- Calls: ordering.policy.resolve, ordering.apply.deterministic, action.generate.markdown

## Define a renderer plugin registry contract [renderer.registry.contract]

- Description: Define a small renderer interface (name, supportsGraphShape, supportedArguments, render) and register built-ins (markdown-text, mermaid) in a deterministic lookup map; renderers consume one typed validated renderer-argument set resolved before plugin invocation.
- Calls: renderer.registry.resolve, renderer.plugin.select, render.section.graph, render.graph.markdown.text, render.graph.mermaid

## Use deterministic renderer selection and fallback policy [renderer.selection.fallback-policy]

- Description: Resolve renderer from renderer-scoped arguments sourced from H3Section and notes first, then apply stable defaults by graph shape (Mermaid-first for cyclic graphs, markdown-first for tree/DAG); if cycle-policy is `disallow` and cycles are detected, emit an error diagnostic and skip graph rendering for that section.
- Calls: renderer.plugin.select, graph.shape.detect, render.graph.tree-or-dag, render.graph.circular

## Use early returns and guard clauses for errors [style.errors.guard-clauses]

- Description: Handle invalid inputs and failure states first, return immediately, and keep the success path shallow and readable.
- Calls: action.generate.markdown, action.generate.json, action.validate, validate.app.data

## Keep functions small and single-purpose [style.functions.small-single-purpose]

- Description: Each function should do one thing and remain easy to test in isolation; prefer composition of small steps over large multi-branch handlers.
- Calls: action.generate.markdown.sections, action.generate.markdown.section.h3, render.section.file

## Separate I/O from core logic [style.io-separate-from-logic]

- Description: Keep parsing, filtering, and rendering logic pure where possible, and isolate file/network/process I/O behind adapter functions.
- Calls: load.app.data, validate.app.data, render.section.file.csv, render.section.file.media, action.generate.markdown

## Use tiny structs to avoid long parameter lists [style.parameters.tiny-structs]

- Description: Group related parameters into small intent-revealing structs (for example, render context and filter options) to reduce call-site ambiguity.
- Calls: action.generate.markdown.section.h3, graph.select, file.csv.filter, render.section.file

## Replace boolean soup with named predicates [style.predicates.named]

- Description: Extract compound conditions into well-named predicate helpers to clarify branching and make tests easier to read.
- Calls: graph.select, file.csv.filter, render.section.file, validate.app.data
