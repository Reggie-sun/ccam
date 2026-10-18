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

// ccam's own help is Korean; the site shows English placeholders.
const words: Record<string, string> = {
  "새 이름": "new name", "이름": "name", "번호": "number", "폴더": "folder", "그룹": "group",
  "계정": "account", "서버": "server", "명령 인자…": "command…", "주소": "url", "claude 인자": "claude args",
};
const english = (usage: string) => Object.entries(words).reduce((u, [ko, en]) => u.replaceAll(ko, en), usage);

const out = [
  { key: "ccs", usage: "ccs" },
  { key: "name", usage: "ccs <name|number> [claude arguments]" },
];
for (const name of names) {
  const usage = /Usage:\n\s+(.+)/.exec(help(name))?.[1].replace(/^ccam/, "ccam").trim() ?? `ccam ${name}`;
  out.push({ key: name, usage: english(usage.replace(/\s+\[flags\]$/, "")) });
}
writeFileSync("data/commands.json", JSON.stringify(out, null, 2) + "\n");
console.log(`commands: ${out.length} entries`);
