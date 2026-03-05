import { describe, expect, test } from 'bun:test';
import { calls } from './calls.ts';
import { exampleDiagnostic } from './examples.ts';
import { cliRoot } from './flows.ts';

type IndexedCall = {
  index: number;
  level: number;
  name: string;
  parentIndex: number | null;
  parentName: string | null;
};

const buildFlow = () => {
  calls.length = 0;
  cliRoot({ level: 0 });

  const stack: number[] = [];
  const indexed: IndexedCall[] = [];

  for (let i = 0; i < calls.length; i += 1) {
    const call = calls[i];
    while (stack.length > 0 && calls[stack[stack.length - 1]].level >= call.level) {
      stack.pop();
    }

    const parentIndex = stack.length > 0 ? stack[stack.length - 1] : null;
    indexed.push({
      index: i,
      level: call.level,
      name: call.name,
      parentIndex,
      parentName: parentIndex === null ? null : calls[parentIndex].name,
    });

    stack.push(i);
  }

  return indexed;
};

const directChildren = (indexed: IndexedCall[], parentName: string) =>
  indexed.filter((item) => item.parentName === parentName).map((item) => item.name);

const firstChildByName = (
  indexed: IndexedCall[],
  parentIndex: number,
  childName: string,
) =>
  indexed.find(
    (item) => item.parentIndex === parentIndex && item.name === childName,
  );

const subtreeNames = (indexed: IndexedCall[], rootIndex: number) => {
  const rootLevel = indexed[rootIndex].level;
  const names: string[] = [indexed[rootIndex].name];

  for (let i = rootIndex + 1; i < indexed.length; i += 1) {
    if (indexed[i].level <= rootLevel) {
      break;
    }
    names.push(indexed[i].name);
  }

  return names;
};

describe('validation entrypoint refactor', () => {
  test('uses validate.app.data consistently across command flows', () => {
    const indexed = buildFlow();

    expect(directChildren(indexed, 'action.generate.markdown')).toEqual([
      'load.app.data',
      'validate.app.data',
      'action.generate.markdown.sections',
    ]);

    expect(directChildren(indexed, 'action.generate.json')).toEqual([
      'load.app.data',
      'validate.app.data',
      'export.graph.json',
    ]);

    expect(directChildren(indexed, 'action.validate')).toEqual([
      'load.app.data',
      'validate.app.data',
      'diagnostics.emit.structured',
    ]);

    expect(directChildren(indexed, 'action.list.reports')).toEqual([
      'load.app.data',
      'validate.app.data',
      'list.reports.output',
    ]);
  });

  test('golden: validate and generate markdown share the same validation diagnostics pipeline', () => {
    const indexed = buildFlow();

    const generateMarkdown = indexed.find(
      (item) => item.name === 'action.generate.markdown',
    );
    const validate = indexed.find((item) => item.name === 'action.validate');

    if (!generateMarkdown || !validate) {
      throw new Error('Expected action.generate.markdown and action.validate');
    }

    const generateValidation = firstChildByName(
      indexed,
      generateMarkdown.index,
      'validate.app.data',
    );
    const validateValidation = firstChildByName(
      indexed,
      validate.index,
      'validate.app.data',
    );

    if (!generateValidation || !validateValidation) {
      throw new Error('Expected validate.app.data under generate markdown and validate');
    }

    const generatePipeline = subtreeNames(indexed, generateValidation.index);
    const validatePipeline = subtreeNames(indexed, validateValidation.index);

    const expectedPipeline = [
      'validate.app.data',
      'validate.cue.schema',
      'args.registry.resolve',
      'args.registry.validate',
      'args.validate.config',
      'graph.integrity.policy.resolve',
      'graph.integrity.validate',
      'graph.integrity.check.missing-nodes',
      'graph.integrity.check.orphans',
      'graph.integrity.check.duplicate-note-names',
      'graph.integrity.check.unknown-labels',
      'graph.integrity.check.cross-report-references',
      'diagnostics.collect.validation',
      'app.model.normalize',
    ];

    expect(generatePipeline).toEqual(expectedPipeline);
    expect(validatePipeline).toEqual(expectedPipeline);
  });

  test('graph.integrity.validate is only invoked inside validate.app.data', () => {
    const indexed = buildFlow();

    const integrityCalls = indexed.filter(
      (item) => item.name === 'graph.integrity.validate',
    );
    expect(integrityCalls.length).toBeGreaterThan(0);

    for (const item of integrityCalls) {
      expect(item.parentName).toBe('validate.app.data');
    }
  });

  test('diagnostic codes and locations remain stable', () => {
    const diagnostic = JSON.parse(exampleDiagnostic) as {
      code: string;
      severity: string;
      source: string;
      location?: string;
    };

    expect(diagnostic.code).toBe('GRAPH_INTEGRITY_MISSING_NODE');
    expect(diagnostic.severity).toBe('error');
    expect(diagnostic.source).toBe('graph.integrity.validate');
    expect(diagnostic.location).toBe('relationships[0].to');
    expect(diagnostic.location).toMatch(/^relationships\[\d+\]\.(from|to)$/);
  });

  test('renderer arguments are resolved deterministically before plugin selection', () => {
    const indexed = buildFlow();

    expect(directChildren(indexed, 'render.section.graph')).toEqual([
      'renderer.registry.resolve',
      'args.renderer.resolve',
      'renderer.plugin.select',
      'graph.policy.cycle',
      'graph.shape.detect',
      'render.graph.tree-or-dag',
      'render.graph.circular',
    ]);

    const h3Resolve = calls.find((call) => call.name === 'args.h3.resolve');
    const noteResolve = calls.find((call) => call.name === 'args.note.resolve');
    const rendererResolve = calls.find(
      (call) => call.name === 'args.renderer.resolve',
    );
    const pluginSelect = calls.find(
      (call) => call.name === 'renderer.plugin.select',
    );

    if (!h3Resolve || !noteResolve || !rendererResolve || !pluginSelect) {
      throw new Error('Expected renderer argument resolution calls to exist');
    }

    expect(h3Resolve.note).toContain('renderer-scoped resolution');
    expect(noteResolve.note).toContain('higher precedence');
    expect(rendererResolve.note).toContain('`note` overrides `h3-section`');
    expect(rendererResolve.note).toContain('typed validated renderer argument set');
    expect(pluginSelect.note).toContain('pass one resolved renderer argument set');
  });

  test('ordering policy defines explicit comparators and tie-breakers', () => {
    buildFlow();

    const resolveOrdering = calls.find(
      (call) => call.name === 'ordering.policy.resolve',
    );
    const applyOrdering = calls.find(
      (call) => call.name === 'ordering.apply.deterministic',
    );

    if (!resolveOrdering || !applyOrdering) {
      throw new Error('Expected ordering policy calls to exist');
    }

    expect(resolveOrdering.note).toContain('(primaryLabel, name)');
    expect(resolveOrdering.note).toContain('(from, to, labelsSortedJoined)');
    expect(resolveOrdering.note).toContain('joined with `|`');
    expect(resolveOrdering.note).toContain('(lowercase(title), originalIndex)');
    expect(resolveOrdering.note).toContain('arguments by argument name');
    expect(applyOrdering.note).toContain('stable tie-breakers only');
    expect(applyOrdering.note).toContain('without runtime randomness');
  });
});
