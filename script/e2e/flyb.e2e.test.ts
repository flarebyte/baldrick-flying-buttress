import { expect, test } from 'bun:test';
import {
  mkdirSync,
  mkdtempSync,
  readFileSync,
  rmSync,
  writeFileSync,
} from 'node:fs';
import { tmpdir } from 'node:os';
import { dirname, join } from 'node:path';

const rootDir = join(import.meta.dir, '..', '..');
const binPath = join(rootDir, '.e2e-bin', 'flyb');
const fixturePath = join(rootDir, 'testdata', 'app.raw.json');
const cueFixturePath = join(rootDir, 'testdata', 'app.cue');
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
const markdownOrphansFixturePath = join(
  rootDir,
  'testdata',
  'markdown.orphans.raw.json',
);
const markdownFileFixturePath = join(
  rootDir,
  'testdata',
  'markdown.file.raw.json',
);
const markdownFileMultiFilterFixturePath = join(
  rootDir,
  'testdata',
  'markdown.file.multifilter.raw.json',
);
const markdownFileRawCSVFixturePath = join(
  rootDir,
  'testdata',
  'markdown.file.rawcsv.raw.json',
);
const designMetaDirPath = join(rootDir, 'doc', 'design-meta');
const designMetaAppCuePath = join(designMetaDirPath, 'app.cue');

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

function makeTempFixtureWithFiles(sourcePath: string, relativePaths: string[]) {
  const fixture = makeTempFixture(sourcePath);
  for (const relPath of relativePaths) {
    const srcPath = join(rootDir, 'testdata', relPath);
    const dstPath = join(fixture.dir, relPath);
    mkdirSync(dirname(dstPath), { recursive: true });
    writeFileSync(dstPath, readFileSync(srcPath));
  }
  return fixture;
}

function runGenerateMarkdown(fixtureConfigPath: string) {
  return runFlyb(['generate', 'markdown', '--config', fixtureConfigPath]);
}

function assertGenerateMarkdownOutput(
  fixtureDir: string,
  reportPath: string,
  goldenName: string,
) {
  const got = runGenerateMarkdown(join(fixtureDir, 'config.raw.json'));
  const want = readGolden(goldenName);
  const gotOutput = readFileSync(join(fixtureDir, 'out', reportPath));

  expect(got.exitCode).toBe(0);
  expect(bytesHex(got.stdout)).toBe('');
  expect(bytesHex(got.stderr)).toBe('');
  expect(bytesHex(gotOutput)).toBe(bytesHex(want));
}

test('flyb validate stdout matches golden', () => {
  const got = runFlyb(['validate', '--config', fixturePath]);
  const wantStdout = readGolden('validate.stdout.golden');

  expect(got.exitCode).toBe(1);
  expect(bytesHex(got.stdout)).toBe(bytesHex(wantStdout));
  expect(bytesHex(got.stderr)).toBe('');
});

test('flyb validate with cue fixture stdout matches golden', () => {
  const got = runFlyb(['validate', '--config', cueFixturePath]);
  const wantStdout = readGolden('validate-cue.stdout.golden');

  expect(got.exitCode).toBe(1);
  expect(bytesHex(got.stdout)).toBe(bytesHex(wantStdout));
  expect(bytesHex(got.stderr)).toBe('');
});

test('flyb directory config matches equivalent app.cue package config', () => {
  const validateFromDir = runFlyb(['validate', '--config', designMetaDirPath]);
  const validateFromFile = runFlyb([
    'validate',
    '--config',
    designMetaAppCuePath,
  ]);

  expect(validateFromDir.exitCode).toBe(validateFromFile.exitCode);
  expect(bytesHex(validateFromDir.stdout)).toBe(
    bytesHex(validateFromFile.stdout),
  );
  expect(bytesHex(validateFromDir.stderr)).toBe(
    bytesHex(validateFromFile.stderr),
  );

  const listFromDir = runFlyb([
    'list',
    'reports',
    '--config',
    designMetaDirPath,
  ]);
  const listFromFile = runFlyb([
    'list',
    'reports',
    '--config',
    designMetaAppCuePath,
  ]);

  expect(listFromDir.exitCode).toBe(listFromFile.exitCode);
  expect(bytesHex(listFromDir.stdout)).toBe(bytesHex(listFromFile.stdout));
  expect(bytesHex(listFromDir.stderr)).toBe(bytesHex(listFromFile.stderr));
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
    assertGenerateMarkdownOutput(
      fixture.dir,
      'renderer-explicit.md',
      'generate-markdown-renderer-explicit.golden',
    );
  } finally {
    rmSync(fixture.dir, { recursive: true, force: true });
  }
});

test('flyb generate markdown can target a single report', () => {
  const fixture = makeTempFixture(markdownFixturePath);
  try {
    const got = runFlyb([
      'generate',
      'markdown',
      '--config',
      fixture.configPath,
      '--report',
      'alpha',
    ]);
    const wantAlpha = readGolden('generate-markdown-alpha.golden');
    const gotAlpha = readFileSync(join(fixture.dir, 'out', 'alpha.md'));

    expect(got.exitCode).toBe(0);
    expect(bytesHex(got.stdout)).toBe('');
    expect(bytesHex(got.stderr)).toBe('');
    expect(bytesHex(gotAlpha)).toBe(bytesHex(wantAlpha));
    expect(() => readFileSync(join(fixture.dir, 'out', 'beta.md'))).toThrow();
  } finally {
    rmSync(fixture.dir, { recursive: true, force: true });
  }
});

test('flyb generate markdown renders orphan sections', () => {
  const fixture = makeTempFixture(markdownOrphansFixturePath);
  try {
    assertGenerateMarkdownOutput(
      fixture.dir,
      'orphans-subject.md',
      'generate-markdown-orphans.golden',
    );
  } finally {
    rmSync(fixture.dir, { recursive: true, force: true });
  }
});

test('flyb generate markdown renders file-backed sections', () => {
  const fixture = makeTempFixtureWithFiles(markdownFileFixturePath, [
    'fixtures/data.csv',
    'fixtures/diagram.png',
    'fixtures/flow.mmd',
  ]);
  try {
    assertGenerateMarkdownOutput(
      fixture.dir,
      'file.md',
      'generate-markdown-file.golden',
    );
  } finally {
    rmSync(fixture.dir, { recursive: true, force: true });
  }
});

test('flyb generate markdown renders csv table with multiple filters', () => {
  const fixture = makeTempFixtureWithFiles(markdownFileMultiFilterFixturePath, [
    'fixtures/data.filters.csv',
  ]);
  try {
    assertGenerateMarkdownOutput(
      fixture.dir,
      'file-multifilter.md',
      'generate-markdown-file-multifilter.golden',
    );
  } finally {
    rmSync(fixture.dir, { recursive: true, force: true });
  }
});

test('flyb generate markdown renders raw csv mode', () => {
  const fixture = makeTempFixtureWithFiles(markdownFileRawCSVFixturePath, [
    'fixtures/data.filters.csv',
  ]);
  try {
    assertGenerateMarkdownOutput(
      fixture.dir,
      'file-rawcsv.md',
      'generate-markdown-file-rawcsv.golden',
    );
  } finally {
    rmSync(fixture.dir, { recursive: true, force: true });
  }
});

test('flyb generate markdown file-backed output is deterministic across runs', () => {
  const fixture = makeTempFixtureWithFiles(markdownFileFixturePath, [
    'fixtures/data.csv',
    'fixtures/diagram.png',
    'fixtures/flow.mmd',
  ]);
  try {
    const first = runGenerateMarkdown(fixture.configPath);
    const firstOutput = readFileSync(join(fixture.dir, 'out', 'file.md'));
    const second = runGenerateMarkdown(fixture.configPath);
    const secondOutput = readFileSync(join(fixture.dir, 'out', 'file.md'));

    expect(second.exitCode).toBe(first.exitCode);
    expect(bytesHex(second.stdout)).toBe(bytesHex(first.stdout));
    expect(bytesHex(second.stderr)).toBe(bytesHex(first.stderr));
    expect(bytesHex(secondOutput)).toBe(bytesHex(firstOutput));
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

test('flyb validate oversized config fails deterministically', () => {
  const fixture = makeTempFixture(fixturePath);
  try {
    const oversized = 'a'.repeat(1024 * 1024 + 1);
    writeFileSync(fixture.configPath, oversized);

    const first = runFlyb(['validate', '--config', fixture.configPath]);
    const second = runFlyb(['validate', '--config', fixture.configPath]);

    expect(first.exitCode).toBe(2);
    expect(second.exitCode).toBe(2);
    expect(bytesHex(first.stdout)).toBe('');
    expect(bytesHex(second.stdout)).toBe('');
    expect(bytesHex(second.stderr)).toBe(bytesHex(first.stderr));
  } finally {
    rmSync(fixture.dir, { recursive: true, force: true });
  }
});
