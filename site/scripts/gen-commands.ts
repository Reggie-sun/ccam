// Builds data/commands.json from the real `ccam --help` output, so the
// command reference cannot drift from the program. Descriptions come from
// messages/<locale>.json ("commands.<key>"); this only supplies usage lines.
import { execFileSync } from "node:child_process";
import { mkdtempSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";

// Both folders are removed on exit: every site build runs this, and a
// 15 MB binary left in /tmp each time fills a small tmpfs.
const binDir = mkdtempSync(join(tmpdir(), "ccam-"));
const home = mkdtempSync(join(tmpdir(), "ccam-home-"));
process.on("exit", () => {
  for (const d of [binDir, home]) rmSync(d, { recursive: true, force: true });
});
const bin = join(binDir, "ccam");
execFileSync("go", ["build", "-buildvcs=false", "-o", bin, "./cmd/ccam"], { cwd: "..", stdio: "inherit" });
const env = { ...process.env, CCAM_HOME: home };
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
