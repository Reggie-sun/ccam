// data/changelog.json from GitHub Releases. Keeps the previous file when the
// API is unreachable (offline builds, private repository without a token).
import { writeFileSync } from "node:fs";
import { Marked } from "marked";

const headers: Record<string, string> = { Accept: "application/vnd.github+json", "User-Agent": "ccam-site" };
if (process.env.GITHUB_TOKEN) headers.Authorization = `Bearer ${process.env.GITHUB_TOKEN}`;

/** goreleaser lists every commit as "* <sha> type(scope): message". Readers
 *  want what changed in ccam: features, fixes and speed-ups, without commit
 *  hashes, tests, chores or changes to this site. Other text is kept. */
function userFacing(body: string): string {
  const groups: Record<string, string[]> = { feat: [], fix: [], perf: [] };
  const rest: string[] = [];
  for (const line of body.split(/\r?\n/)) {
    if (/^[*-] (?:[0-9a-f]{7,40} )?Merge /.test(line)) continue; // merge commits say nothing new
    const c = /^[*-] (?:[0-9a-f]{7,40} )?(\w+)(?:\(([^)]*)\))?!?: (.+)$/.exec(line);
    if (!c) {
      if (!/^## Changelog\s*$/.test(line)) rest.push(line);
      continue;
    }
    const [, type, scope, message] = c;
    const sentence = /^ccam\b/.test(message) ? message : message[0].toUpperCase() + message.slice(1);
    if (type in groups && scope !== "site") groups[type].push(`- ${sentence}`);
  }
  const titles: Record<string, string> = { feat: "New", fix: "Fixed", perf: "Faster" };
  const lists = Object.entries(groups).filter(([, items]) => items.length).map(([type, items]) => `### ${titles[type]}\n\n${items.join("\n")}`);
  return [rest.join("\n").trim(), ...lists].filter(Boolean).join("\n\n");
}

try {
  const res = await fetch("https://api.github.com/repos/jhyunwoo/ccam/releases?per_page=50", { headers, signal: AbortSignal.timeout(8000) });
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  const releases = (await res.json()) as { tag_name: string; published_at: string; html_url: string; body: string; draft: boolean; prerelease: boolean }[];
  const marked = new Marked();
  const out = releases
    .filter((r) => !r.draft)
    .map((r) => ({ tag: r.tag_name, date: r.published_at.slice(0, 10), url: r.html_url, notes: marked.parse(userFacing(r.body ?? ""), { async: false }) as string }));
  writeFileSync("data/changelog.json", JSON.stringify(out, null, 2) + "\n");
  console.log(`changelog: ${out.length} releases`);
} catch (e) {
  console.warn(`changelog: kept the existing data/changelog.json (${(e as Error).message})`);
}
