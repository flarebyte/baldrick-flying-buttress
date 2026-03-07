# Microservice architecture example

This example demonstrates:
- a medium-sized service graph (5 notes)
- multiple relationship paths
- graph section rendering
- explicit Mermaid renderer usage

## Run

```bash
flyb validate --config app.cue
flyb generate markdown --config app.cue
```

## Expected output

Generates:
- `out/microservice.md`

The report includes:
- an H3 section rendered as Mermaid (`graph-renderer=mermaid`)
- an H3 section rendered as markdown-text fallback/tree view
