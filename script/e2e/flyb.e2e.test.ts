import { expect, test } from 'bun:test';
import {
  mkdirSync,
  mkdtempSync,
  readFileSync,
  rmSync,
  writeFileSync,
} from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';

const rootDir = join(import.meta.dir, '..', '..');
const binPath = join(rootDir, '.e2e-bin', 'flyb');
const fixturePath = join(rootDir, 'testdata', 'app.raw.json');
const namesFixturePath = join(rootDir, 'testdata', 'names.raw.json');
const lintFixturePath = join(rootDir, 'testdata', 'lint.raw.json');
const orphansFixturePath = join(rootDir, 'testdata', 'orphans.raw.json');
const markdownFixturePath = join(rootDir, 'testdata', 'markdown.raw.json');
const markdownGraphFixturePath = join(
  rootDir,
  'testdata',
  'markdown.graph.raw.json',
);
const markdownRendererExplicitFixturePath = join(
  rootDir,
  'testdata',
  'markdown.renderer.explicit.raw.json',
);

function runFlyb(args: string[]) {
  const result = Bun.spawnSync({
    cmd: [binPath, ...args],
    stdout: 'pipe',
    stderr: 'pipe',
  });

  if (result.exitCode === null) {
    throw new Error('flyb process did not return an exit code');
  }

  return {
    exitCode: result.exitCode,
    stdout: result.stdout,
    stderr: result.stderr,
  };
}

function readGolden(name: string) {
  return readFileSync(join(import.meta.dir, 'testdata', name));
}

function bytesHex(data: Uint8Array) {
  return Buffer.from(data).toString('hex');
}

function makeTempFixture(sourcePath: string) {
  const dir = mkdtempSync(join(tmpdir(), 'flyb-e2e-'));
  const configPath = join(dir, 'config.raw.json');
  writeFileSync(configPath, readFileSync(sourcePath));
  return { dir, configPath };
}

test('flyb validate stdout matches golden', () => {
  const got = runFlyb(['validate', '--config', fixturePath]);
  const wantStdout = readGolden('validate.stdout.golden');

  expect(got.exitCode).toBe(1);
  expect(bytesHex(got.stdout)).toBe(bytesHex(wantStdout));
  expect(bytesHex(got.stderr)).toBe('');
});

test('flyb list reports stdout matches golden', () => {
  const got = runFlyb(['list', 'reports', '--config', fixturePath]);
  const wantStdout = readGolden('list-reports.stdout.golden');

  expect(got.exitCode).toBe(1);
  expect(bytesHex(got.stdout)).toBe(bytesHex(wantStdout));
  expect(bytesHex(got.stderr)).toBe('');
});

test('flyb generate json stdout matches golden', () => {
  const got = runFlyb(['generate', 'json', '--config', fixturePath]);
  const wantStdout = readGolden('generate-json.stdout.golden');

  expect(got.exitCode).toBe(1);
  expect(bytesHex(got.stdout)).toBe(bytesHex(wantStdout));
  expect(bytesHex(got.stderr)).toBe('');
});

test('flyb list names stdout matches golden', () => {
  const got = runFlyb([
    'list',
    'names',
    '--config',
    namesFixturePath,
    '--prefix',
    'cli.',
  ]);
  const wantStdout = readGolden('list-names.stdout.golden');

  expect(got.exitCode).toBe(0);
  expect(bytesHex(got.stdout)).toBe(bytesHex(wantStdout));
  expect(bytesHex(got.stderr)).toBe('');
});

test('flyb lint names stdout matches golden', () => {
  const got = runFlyb(['lint', 'names', '--config', lintFixturePath]);
  const wantStdout = readGolden('lint-names.stdout.golden');

  expect(got.exitCode).toBe(0);
  expect(bytesHex(got.stdout)).toBe(bytesHex(wantStdout));
  expect(bytesHex(got.stderr)).toBe('');
});

test('flyb lint orphans stdout matches golden', () => {
  const got = runFlyb([
    'lint',
    'orphans',
    '--config',
    orphansFixturePath,
    '--subject-label',
    'ingredient',
  ]);
  const wantStdout = readGolden('lint-orphans.stdout.golden');

  expect(got.exitCode).toBe(0);
  expect(bytesHex(got.stdout)).toBe(bytesHex(wantStdout));
  expect(bytesHex(got.stderr)).toBe('');
});

test('flyb generate markdown writes deterministic report files', () => {
  const fixture = makeTempFixture(markdownFixturePath);
  try {
    const got = runFlyb([
      'generate',
      'markdown',
      '--config',
      fixture.configPath,
    ]);
    const wantAlpha = readGolden('generate-markdown-alpha.golden');
    const wantBeta = readGolden('generate-markdown-beta.golden');
    const gotAlpha = readFileSync(join(fixture.dir, 'out', 'alpha.md'));
    const gotBeta = readFileSync(join(fixture.dir, 'out', 'beta.md'));

    expect(got.exitCode).toBe(0);
    expect(bytesHex(got.stdout)).toBe('');
    expect(bytesHex(got.stderr)).toBe('');
    expect(bytesHex(gotAlpha)).toBe(bytesHex(wantAlpha));
    expect(bytesHex(gotBeta)).toBe(bytesHex(wantBeta));
  } finally {
    rmSync(fixture.dir, { recursive: true, force: true });
  }
});

test('flyb generate markdown renders graph sections', () => {
  const fixture = makeTempFixture(markdownGraphFixturePath);
  try {
    const got = runFlyb([
      'generate',
      'markdown',
      '--config',
      fixture.configPath,
    ]);
    const want = readGolden('generate-markdown-graph.golden');
    const gotGraph = readFileSync(join(fixture.dir, 'out', 'graph.md'));

    expect(got.exitCode).toBe(0);
    expect(bytesHex(got.stdout)).toBe('');
    expect(bytesHex(got.stderr)).toBe('');
    expect(bytesHex(gotGraph)).toBe(bytesHex(want));
  } finally {
    rmSync(fixture.dir, { recursive: true, force: true });
  }
});

test('flyb generate markdown supports explicit mermaid renderer', () => {
  const fixture = makeTempFixture(markdownRendererExplicitFixturePath);
  try {
    const got = runFlyb([
      'generate',
      'markdown',
      '--config',
      fixture.configPath,
    ]);
    const want = readGolden('generate-markdown-renderer-explicit.golden');
    const gotOutput = readFileSync(
      join(fixture.dir, 'out', 'renderer-explicit.md'),
    );

    expect(got.exitCode).toBe(0);
    expect(bytesHex(got.stdout)).toBe('');
    expect(bytesHex(got.stderr)).toBe('');
    expect(bytesHex(gotOutput)).toBe(bytesHex(want));
  } finally {
    rmSync(fixture.dir, { recursive: true, force: true });
  }
});

test('flyb outputs are deterministic across repeated runs', () => {
  const commands: string[][] = [
    ['validate', '--config', fixturePath],
    ['list', 'reports', '--config', fixturePath],
    ['generate', 'json', '--config', fixturePath],
    ['list', 'names', '--config', namesFixturePath, '--prefix', 'cli.'],
    ['lint', 'names', '--config', lintFixturePath],
    [
      'lint',
      'orphans',
      '--config',
      orphansFixturePath,
      '--subject-label',
      'ingredient',
    ],
  ];

  for (const args of commands) {
    const first = runFlyb(args);
    const second = runFlyb(args);

    expect(second.exitCode).toBe(first.exitCode);
    expect(bytesHex(second.stdout)).toBe(bytesHex(first.stdout));
    expect(bytesHex(second.stderr)).toBe(bytesHex(first.stderr));
  }
});

test('flyb generate markdown is deterministic across repeated runs', () => {
  const fixture = makeTempFixture(markdownFixturePath);
  try {
    const first = runFlyb([
      'generate',
      'markdown',
      '--config',
      fixture.configPath,
    ]);
    const firstAlpha = readFileSync(join(fixture.dir, 'out', 'alpha.md'));
    const firstBeta = readFileSync(join(fixture.dir, 'out', 'beta.md'));

    mkdirSync(join(fixture.dir, 'out'), { recursive: true });
    const second = runFlyb([
      'generate',
      'markdown',
      '--config',
      fixture.configPath,
    ]);
    const secondAlpha = readFileSync(join(fixture.dir, 'out', 'alpha.md'));
    const secondBeta = readFileSync(join(fixture.dir, 'out', 'beta.md'));

    expect(second.exitCode).toBe(first.exitCode);
    expect(bytesHex(second.stdout)).toBe(bytesHex(first.stdout));
    expect(bytesHex(second.stderr)).toBe(bytesHex(first.stderr));
    expect(bytesHex(secondAlpha)).toBe(bytesHex(firstAlpha));
    expect(bytesHex(secondBeta)).toBe(bytesHex(firstBeta));
  } finally {
    rmSync(fixture.dir, { recursive: true, force: true });
  }
});
