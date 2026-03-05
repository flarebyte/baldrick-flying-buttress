# baldrick-flying-buttress

**Divert the weight of design away from the code with a flying buttress.**

`baldrick-flying-buttress` (**flyb**) is a CLI that turns **structured design notes and relationships** into **living architecture documentation**.

Instead of writing large design documents that quickly become outdated, flyb lets you describe your system as **small connected notes**. From those notes, the CLI generates **deterministic Markdown reports and diagrams** that stay consistent with the model.

Think of it as **“architecture-as-data”**.

---

## The idea

You define:

- **notes** — small pieces of design knowledge
- **relationships** — how those pieces connect
- **reports** — curated views of the graph

Example:

```
cli.root ──satisfies-usecase──► usecase.io.calls.count
```

From this simple structure, flyb can generate:

- architecture documentation
- call flow graphs
- dependency maps
- use-case mappings
- Markdown reports with Mermaid diagrams

---

## Why flyb exists

Documentation tends to fail in one of two ways:

- it becomes **too large and hard to maintain**
- it becomes **out of sync with the code**

flyb solves this by making documentation:

- **modular** — small notes instead of large documents
- **structured** — relationships instead of free text
- **deterministic** — generated output works well with Git
- **close to the code** — configuration lives in the repository

---

## What flyb helps teams do

With flyb you can:

- generate architecture reports automatically
- visualize relationships between system components
- validate design models for consistency
- enforce naming conventions across design identifiers
- detect unused or disconnected concepts
- embed diagrams, tables, and references in documentation

---

## Example CLI workflow

Generate documentation:

```
flyb generate markdown
```

Inspect notes and relationships:

```
flyb list names
```

Validate configuration:

```
flyb validate
```

Check naming consistency:

```
flyb lint names --style dot
```

---

## What makes flyb different

Most documentation tools start with **documents**.

flyb starts with a **graph of knowledge**.

Documentation becomes a **generated view** of that graph.

This keeps design knowledge **structured, reusable, and maintainable** as systems grow.
