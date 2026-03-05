import { describe, expect, test } from 'bun:test';
import { readFileSync } from 'node:fs';
import { join } from 'node:path';
import { exampleArgumentRegistry } from './examples.ts';

const root = process.cwd();

describe('argument scope consistency', () => {
  test('example registry uses only supported scopes', () => {
    const registry = JSON.parse(exampleArgumentRegistry) as {
      arguments: Array<{ scopes: string[] }>;
    };

    const allowed = new Set(['h3-section', 'note', 'renderer']);
    for (const arg of registry.arguments) {
      for (const scope of arg.scopes) {
        expect(allowed.has(scope)).toBe(true);
      }
    }
  });

  test('spec source has no global/h2-section argument scope references', () => {
    const files = [
      'design/common.ts',
      'design/examples.ts',
      'design/use_cases.ts',
      'design/implementations.ts',
      'design/flows.ts',
    ];

    for (const rel of files) {
      const content = readFileSync(join(root, rel), 'utf8');
      expect(content.includes('global')).toBe(false);
      expect(content.includes('h2-section')).toBe(false);
      expect(content.includes('(global, h2, h3, note, renderer)')).toBe(false);
    }
  });
});
