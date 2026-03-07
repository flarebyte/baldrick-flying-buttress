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

Validate config:

```bash
../../.e2e-bin/flyb validate --config app.cue
```

Generate markdown:

```bash
../../.e2e-bin/flyb generate markdown --config app.cue
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

## Determinism guarantees

- ordering policy is stable for diagnostics, reports, notes, and relationships
- machine-readable JSON output is deterministic
- repeated runs with the same config produce reproducible outputs

## Examples

See starter configurations in [`examples/`](./examples):
- [`examples/minimal`](./examples/minimal)
- [`examples/microservice-architecture`](./examples/microservice-architecture)
- [`examples/orphan-analysis`](./examples/orphan-analysis)
