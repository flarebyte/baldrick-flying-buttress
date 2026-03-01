import type { ImplementationConsideration } from './common.ts';

// Initial implementation suggestions. Keep this list small and actionable.
export const implementations: Record<string, ImplementationConsideration> = {
  ruleNameSeparator: {
    name: 'rules.name.separator.underscore',
    title: 'Use underscore as the rule-name separator',
    description:
      'Use `_` in rule names (for example: `function_map`) instead of `.` to avoid escaping in tools such as jq and JavaScript property access.',
    calls: [
      'cli.analyse',
      'rules.resolve',
      'rules.catalog.list',
      'format.output',
    ],
  },
};
