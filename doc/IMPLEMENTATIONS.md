# Implementation Considerations (Generated)

This document summarizes suggested implementation choices.

## Summary
- Use underscore as the rule-name separator [rules.name.separator.underscore]

## Use underscore as the rule-name separator [rules.name.separator.underscore]

- Description: Use `_` in rule names (for example: `function_map`) instead of `.` to avoid escaping in tools such as jq and JavaScript property access.
- Calls: cli.analyse, rules.resolve, rules.catalog.list, format.output
