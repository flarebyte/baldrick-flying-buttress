# baldrick-flying-buttress (`flyb`)

`flyb` is a config-driven CLI for validating architecture knowledge and generating deterministic documentation from it.

## Project overview

`flyb` models architecture as:
- notes
- relationships between notes
- reports made of sections

It solves the common drift between design docs and reality by keeping design as structured data in source control and generating outputs from that model.

Key design principles:
- deterministic output
- config-driven reports
- machine-readable command output

## Quick start

Install with Homebrew:

```bash
brew install flarebyte/tap/baldrick-flying-buttress
flyb --help
```

Build locally for development:

```bash
make build-dev
./.e2e-bin/flyb --version
```

Use a minimal example project:

```bash
cd examples/minimal
```

`--config` accepts a raw JSON file, a standalone CUE file, or a CUE config directory containing `app.cue`. For packaged CUE configs, `flyb` loads the full package from the directory, including sibling `.cue` files and standard CUE imports.

Validate config:

```bash
flyb validate --config .
```

Generate markdown:

```bash
flyb generate markdown --config .
```

Target a subset of reports while validating or generating:

```bash
flyb validate --config . --report minimal
flyb generate markdown --config . --report minimal
```

## Example CUE configuration

```cue
package flyb

source: "minimal"
name:   "minimal"
modules: ["core"]

reports: [{
	title:    "Minimal Report"
	filepath: "out/minimal.md"
	sections: [{
		title: "Overview"
		sections: [{
			title: "Core"
			notes: ["app.api", "app.db"]
		}]
	}]
}]

notes: [
	{name: "app.api", title: "API", markdown: "Public API."},
	{name: "app.db", title: "Database", markdown: "Primary database."},
]

relationships: [{
	from:   "app.api"
	to:     "app.db"
	label:  "depends_on"
	labels: ["depends_on"]
}]
```

## CLI commands overview

- `flyb validate`
- `flyb list reports`
- `flyb list names`
- `flyb lint names`
- `flyb lint orphans`
- `flyb generate markdown`
- `flyb generate json`

## Project structure explanation

- `notes`: named design nodes (with title/markdown/labels)
- `relationships`: directed links (`from`, `to`, `label`, `labels`)
- `reports`: generated markdown targets (`title`, `filepath`)
- `sections`: nested report structure (`H2` and `H3`) for plain, graph, or orphan views

Validation note:
Notes referenced by report sections, including nested subsections, are treated as intentional documentation nodes and are not flagged as graph orphans just because they have no relationships.

## Graph documentation flows

Graph sections are configured through free-form H3 arguments, so graph rendering can be tuned without changing the core config model. A typical section looks like this:

```cue
{
  title: "Flows"
  arguments: [
    "graph-subject-label=flow",
    "graph-edge-label=contains_step",
    "graph-start-node=cli.root",
    "graph-child-order=title",
    "graph-max-depth=3",
    "graph-branch-priority-label=main",
    "graph-branch-priority-label=future",
    "graph-exclude-label=helper",
    "graph-show-helper-nodes=false",
    "graph-renderer=markdown-text",
  ]
}
```

Useful graph arguments:

- `graph-child-order=id|title`: stable child ordering for rendered branches
- `graph-max-depth=<n>`: limit traversal depth from the selected root/start node
- `graph-include-label=...`: include only notes carrying one of the selected labels
- `graph-exclude-label=...`: remove notes carrying selected labels
- `graph-branch-priority-label=...`: prioritize important branches ahead of others
- `graph-show-helper-nodes=false`: hide helper nodes while keeping the main path readable
- `graph-helper-label=helper`: define which label marks helper nodes when hiding them

These controls are especially useful for large flow-oriented documents where the main path should stay prominent and future or helper branches should be shown selectively.

For plain note sections, `show-labels=true` adds a `Labels: ...` line under each note heading so generated markdown can expose release scope or intent without overloading note titles.

## Determinism guarantees

- ordering policy is stable for diagnostics, reports, notes, and relationships
- machine-readable JSON output is deterministic
- repeated runs with the same config produce reproducible outputs

## Machine-friendly diagnostics

Validation and lint commands emit structured diagnostics designed to be consumable by scripts and AI agents. In addition to `code`, `severity`, `message`, and `path`, diagnostics may now include:

- `normalizedPath`: canonicalized model path for stable matching
- `configPath` and `configPathAbsolute`: normalized config source paths
- dedicated fields such as `reportTitle`, `reportId`, `sectionTitle`, `noteName`, and `noteTitle`
- `relatedNodes`: directly related note or relationship endpoints when applicable
- `suggestedFixes`: concrete next actions for common failures

These fields are optional and appear only when they can be derived reliably from the config and validation context.

## Examples

See starter configurations in [`examples/`](./examples):
- [`examples/minimal`](./examples/minimal)
- [`examples/microservice-architecture`](./examples/microservice-architecture)
- [`examples/orphan-analysis`](./examples/orphan-analysis)
