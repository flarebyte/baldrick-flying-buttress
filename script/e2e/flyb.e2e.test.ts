import { expect, test } from "bun:test";
import { readFileSync } from "node:fs";
import { join } from "node:path";

const rootDir = join(import.meta.dir, "..", "..");
const binPath = join(rootDir, ".e2e-bin", "flyb");

function runFlyb(args: string[]) {
  const result = Bun.spawnSync({
    cmd: [binPath, ...args],
    stdout: "pipe",
    stderr: "pipe",
  });

  if (result.exitCode === null) {
    throw new Error("flyb process did not return an exit code");
  }

  return {
    exitCode: result.exitCode,
    stdout: result.stdout,
    stderr: result.stderr,
  };
}

function readGolden(name: string) {
  return readFileSync(join(import.meta.dir, "testdata", name));
}

function bytesHex(data: Uint8Array) {
  return Buffer.from(data).toString("hex");
}

test("flyb validate stdout matches golden", () => {
  const got = runFlyb(["validate"]);
  const wantStdout = readGolden("validate.stdout.golden");

  expect(got.exitCode).toBe(1);
  expect(bytesHex(got.stdout)).toBe(bytesHex(wantStdout));
  expect(bytesHex(got.stderr)).toBe("");
});

test("flyb list reports stdout matches golden", () => {
  const got = runFlyb(["list", "reports"]);
  const wantStdout = readGolden("list-reports.stdout.golden");

  expect(got.exitCode).toBe(1);
  expect(bytesHex(got.stdout)).toBe(bytesHex(wantStdout));
  expect(bytesHex(got.stderr)).toBe("");
});

test("flyb generate json stdout matches golden", () => {
  const got = runFlyb(["generate", "json"]);
  const wantStdout = readGolden("generate-json.stdout.golden");

  expect(got.exitCode).toBe(1);
  expect(bytesHex(got.stdout)).toBe(bytesHex(wantStdout));
  expect(bytesHex(got.stderr)).toBe("");
});

test("flyb outputs are deterministic across repeated runs", () => {
  const commands: string[][] = [
    ["validate"],
    ["list", "reports"],
    ["generate", "json"],
  ];

  for (const args of commands) {
    const first = runFlyb(args);
    const second = runFlyb(args);

    expect(second.exitCode).toBe(first.exitCode);
    expect(bytesHex(second.stdout)).toBe(bytesHex(first.stdout));
    expect(bytesHex(second.stderr)).toBe(bytesHex(first.stderr));
  }
});
