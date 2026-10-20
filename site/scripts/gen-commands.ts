// Builds data/commands.json from the real `ccam --help` output, so the
// command reference cannot drift from the program. Descriptions come from
// messages/<locale>.json ("commands.<key>"); this only supplies usage lines.
import { execFileSync } from "node:child_process";
import { mkdtempSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";

const bin = join(mkdtempSync(join(tmpdir(), "ccam-")), "ccam");
execFileSync("go", ["build", "-buildvcs=false", "-o", bin, "./cmd/ccam"], { cwd: "..", stdio: "inherit" });
const env = { ...process.env, CCAM_HOME: mkdtempSync(join(tmpdir(), "ccam-home-")) };
const help = (...args: string[]) => execFileSync(bin, [...args, "--help"], { env, encoding: "utf8" });

const skip = new Set(["help", "version", "run"]);
const names = [...help().split("Available Commands:")[1].split(/\n\n/)[0].matchAll(/^\s{2}(\S+)/gm)].map((m) => m[1]).filter((n) => !skip.has(n));

const out = [
  { key: "ccam", usage: "ccam" },
  { key: "name", usage: "ccam <name|number> [claude arguments]" },
  { key: "fallback", usage: "ccam --fallback [name|number] -p [claude arguments]" },
];
for (const name of names) {
  const usage = /Usage:\n\s+(.+)/.exec(help(name))?.[1].replace(/^ccam/, "ccam").trim() ?? `ccam ${name}`;
  out.push({ key: name, usage: usage.replace(/\s+\[flags\]$/, "") });
}
writeFileSync("data/commands.json", JSON.stringify(out, null, 2) + "\n");
console.log(`commands: ${out.length} entries`);
