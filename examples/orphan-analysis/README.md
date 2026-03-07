# Orphan analysis example

This example demonstrates:
- orphan query section arguments on an H3 section
- contextual orphan detection against a label-based relationship query

## Run

```bash
flyb validate --config app.cue
flyb generate markdown --config app.cue
```

## Expected output

Generates:
- `out/orphans.md`

The report includes an orphan-rendered section driven by:
- `orphan-subject-label`
- `orphan-edge-label`
- `orphan-counterpart-label`
- `orphan-direction`
