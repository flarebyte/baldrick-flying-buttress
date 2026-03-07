# Implementation Considerations

Migrated implementation guidance.

## Implementations

Implementation notes.

### Catalog

All implementation considerations.

#### Define a typed argument registry schema

Maintain a registry of argument definitions (name, type, default, allowed values, scopes) and use it as the single source of truth for argument behavior; valid scopes are `h3-section`, `note`, and `renderer`.

#### Validate free-form arguments at runtime

Validate against a known argument registry and fail with clear errors on unknown keys or invalid values.

#### Apply scope-aware argument resolution

Resolve and validate arguments by scope (h3-section, note, renderer) so options are accepted only where they are meaningful; renderer-scoped arguments are collected from H3Section and note argument lists with precedence `note` > `h3-section` > registry defaults.

#### Coerce free-form argument values into typed values

Convert string-like argument inputs into validated typed values before rendering.

#### Use free-form key/value arguments with typed coercion

Treat H3Section and Note arguments like CLI-style flags (for example `format-csv=md`) to keep config flexible, then coerce values into typed runtime options per renderer.

#### Use Cobra for CLI command and argument parsing

Use Cobra to model commands, flags, and subcommands (`generate markdown`, `generate json`, `validate`, `list reports`, `list names`, `lint names`) with a consistent command tree.

#### Use arguments to reduce CUE configuration noise

Prefer composable argument lists over adding many specialized CUE fields, so rendering capabilities can evolve without large schema churn.

#### Use CUE as the configuration source of truth

Represent notes, relationships, and report definitions in CUE for schema validation, defaults, and composable configuration.

#### Use a structured diagnostics model

Standardize diagnostics with code, severity, source, message, canonical machine-readable location, and optional human-readable context fields to support CLI UX, CI checks, and future editor integrations.

#### Attach validation diagnostics to precise config locations

Include canonical index-based CUE path plus related report/section titles and note/relationship/argument identifiers in diagnostics so users can quickly fix invalid configuration.

#### Define a graph integrity policy model

Define explicit integrity rules for missing nodes, orphan nodes, duplicate names, and cross-report references with per-rule severity; validate label references separately against dataset-derived labels.

#### Implement graph integrity validation checks

Run focused integrity checks and emit structured diagnostics linked to note names, relationships, arguments, and CUE paths; keep label-definition handling free-form and validate only label references.

#### Implement the CLI in Go

Use Go as the primary implementation language for strong typing, fast startup, and straightforward single-binary distribution.

#### Implement lint names hygiene command

Implement `flyb lint names` with style policy (`dot|snake|regex`), optional regex `--pattern`, optional prefix scope, and configurable severity; emit structured `NAME_STYLE_VIOLATION` diagnostics with canonical config locations and readable context.

#### Implement list names inventory command

Implement `flyb list names` with required `--prefix`, optional `--kind notes|relationships|all`, and `--format table|json` (default table); reuse validated app data and deterministic ordering before filtering/output.

#### Implement contextual orphan lint command

Implement `flyb lint orphans` using orphan-query filters (`subject-label`, optional edge/counterpart labels, direction) and emit deterministic `ORPHAN_QUERY_MISSING_LINK` diagnostics with stable locations/context.

#### Implement contextual orphan report section renderer

Implement H3 orphan section rendering using orphan-query arguments and deterministic row/table output (`name`, `title`, `labels`) so report sections and lint command evaluate the same orphan set.

#### Guarantee deterministic ordering in generated outputs

Apply explicit stable sorting for notes, relationships, sections, and arguments using concrete comparators (notes: primaryLabel/name, relationships: from/to/labelsSortedJoined, sections: case-insensitive title plus originalIndex, arguments: name) so output remains reproducible across runs and machines.

#### Treat ordering policy as a testable contract

Define ordering rules and tie-breakers as a versioned policy (including label normalization and relationship label joining rules) and verify them with golden-file tests.

#### Define a renderer plugin registry contract

Define a small renderer interface (name, supportsGraphShape, supportedArguments, render) and register built-ins (markdown-text, mermaid) in a deterministic lookup map; renderers consume one typed validated renderer-argument set resolved before plugin invocation.

#### Use deterministic renderer selection and fallback policy

Resolve renderer from renderer-scoped arguments sourced from H3Section and notes first, then apply stable defaults by graph shape (Mermaid-first for cyclic graphs, markdown-first for tree/DAG); if cycle-policy is `disallow` and cycles are detected, emit an error diagnostic and skip graph rendering for that section.

#### Use early returns and guard clauses for errors

Handle invalid inputs and failure states first, return immediately, and keep the success path shallow and readable.

#### Keep functions small and single-purpose

Each function should do one thing and remain easy to test in isolation; prefer composition of small steps over large multi-branch handlers.

#### Separate I/O from core logic

Keep parsing, filtering, and rendering logic pure where possible, and isolate file/network/process I/O behind adapter functions.

#### Use tiny structs to avoid long parameter lists

Group related parameters into small intent-revealing structs (for example, render context and filter options) to reduce call-site ambiguity.

#### Replace boolean soup with named predicates

Extract compound conditions into well-named predicate helpers to clarify branching and make tests easier to read.

