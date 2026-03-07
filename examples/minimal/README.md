# Minimal example

This example demonstrates the smallest practical flyb setup:
- 2 notes
- 1 relationship
- 1 report
- 1 H3 section

## Run

```bash
flyb validate --config app.cue
flyb generate markdown --config app.cue
```

## Expected output

Generates:
- `out/minimal.md`

The markdown report includes one section with inline rendering of the two notes in deterministic order.
