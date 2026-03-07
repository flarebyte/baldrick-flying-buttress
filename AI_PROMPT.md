# AI Prompt: Model CUE Config for Design Docs in flyb

You are migrating existing design documentation into flyb CUE configuration.

## Goal

Produce idiomatic, deterministic CUE files that can be used by `flyb validate` and `flyb generate markdown`.

Keep the model open: source records may vary by project. Map whatever input shape exists to flyb core concepts.

## Core Concepts to Produce

Always model these top-level collections:

- `reports`
- `notes`
- `relationships`

Also provide:

- `source`
- `name`
- `modules`

Use:

- `package flyb`

## Minimum CUE Shape

Use this structure (extend as needed):

```cue
package flyb

source: "<source-id>"
name:   "<app-name>"
modules: ["<module>"]

reports: [...{
  title: string
  filepath: string
  description?: string
  sections: [...{
    title: string
    description?: string
    sections?: [...{
      title: string
      description?: string
      arguments?: [...string]
      notes?: [...string]
      sections?: [..._]
    }]
    arguments?: [...string]
    notes?: [...string]
  }]
}]

notes: [...{
  name: string
  title: string
  markdown?: string
  filepath?: string
  arguments?: [...string]
  labels?: [...string]
}]

relationships: [...{
  from: string
  to: string
  label?: string
  labels?: [...string]
}]

argumentRegistry?: {
  version?: string
  arguments?: [..._]
}

graphIntegrityPolicy?: _
```

## Mapping Strategy (Open Model)

If source records differ, map by intent:

1. Concepts / entities / items -> `notes`

- stable identifier -> `notes[].name`
- human label -> `notes[].title`
- long description -> `notes[].markdown`
- tags / categories -> `notes[].labels`

2. Links / references / dependencies -> `relationships`

- origin identifier -> `relationships[].from`
- target identifier -> `relationships[].to`
- relation type -> `relationships[].label`
- additional relation tags -> `relationships[].labels`

3. Output documents / views -> `reports`

- one report per generated markdown file
- create H2/H3 sections that group note IDs intentionally
- populate `sections[].sections[].notes` with note names

4. Execution or flow docs

- model step nodes as notes with labels like `flow`, `call`
- model step edges with `contains_step`
- for graph-rendered sections, add H3 `arguments`, e.g.:
  - `graph-subject-label=call`
  - `graph-edge-label=contains_step`
  - `graph-start-node=<root-note-id>`
  - `graph-renderer=markdown-text`
  - `cycle-policy=disallow|allow`

5. Orphan analysis sections (optional)

- add H3 args:
  - `orphan-subject-label=<label>`
  - optional `orphan-edge-label=<label>`
  - optional `orphan-counterpart-label=<label>`
  - optional `orphan-direction=in|out|either`

## Determinism Rules

Always enforce deterministic output:

- stable ordering of arrays (`reports`, `notes`, `relationships`)
- stable ordering of labels and arguments
- no maps in output-facing config
- no random/time-based values
- no environment-dependent text

Use deterministic tie-breakers:

- notes: by `name`
- relationships: by `from`, then `to`, then `label`
- reports: by `filepath` or explicit intended order

## Quality Rules

- keep IDs machine-friendly and stable (`dot.case` is recommended)
- keep `title` human-readable
- keep `markdown` concise but meaningful
- avoid duplicate note names
- avoid duplicate relationships unless semantically required
- ensure every relationship endpoint exists in `notes`
- ensure every report note reference exists in `notes`

## Migration Procedure

1. Inventory source artifacts and record types.
2. Build a normalized intermediate list of notes and links.
3. Deduplicate by canonical IDs.
4. Sort deterministically.
5. Build reports and section note lists.
6. Emit CUE with `package flyb`.
7. Run:

- `flyb validate --config <app.cue>`
- `flyb generate markdown --config <app.cue>`

8. Fix unresolved note references, duplicate IDs, and malformed arguments.

## Output Requirements for the Agent

When generating files:

- produce idiomatic CUE only
- keep files small and readable
- prefer one `app.cue` that is self-contained unless multi-file loading is explicitly supported
- if splitting files, ensure the project loader can consume them

## Optional Labels Convention

Use labels consistently; adapt per project:

- `design`
- `flow`
- `call`
- `usecase`
- `implementation`
- `risk`
- domain-specific tags (`service`, `api`, `db`, etc.)

## Example H3 Graph Section Snippet

```cue
{
  title: "Function calls tree"
  sections: [{
    title: "Flow call graph"
    arguments: [
      "graph-subject-label=call",
      "graph-edge-label=contains_step",
      "graph-start-node=cli.root",
      "graph-renderer=markdown-text",
      "cycle-policy=disallow",
    ]
    notes: [
      "cli.root",
      "action.generate.markdown",
      "validate.app.data",
    ]
  }]
}
```
