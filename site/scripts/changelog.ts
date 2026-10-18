// data/changelog.json from GitHub Releases. Keeps the previous file when the
// API is unreachable (offline builds, private repository without a token).
import { writeFileSync } from "node:fs";
import { Marked } from "marked";

const headers: Record<string, string> = { Accept: "application/vnd.github+json", "User-Agent": "ccam-site" };
if (process.env.GITHUB_TOKEN) headers.Authorization = `Bearer ${process.env.GITHUB_TOKEN}`;

try {
  const res = await fetch("https://api.github.com/repos/jhyunwoo/ccam/releases?per_page=50", { headers, signal: AbortSignal.timeout(8000) });
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  const releases = (await res.json()) as { tag_name: string; published_at: string; html_url: string; body: string; draft: boolean; prerelease: boolean }[];
  const marked = new Marked();
  const out = releases
    .filter((r) => !r.draft)
    .map((r) => ({ tag: r.tag_name, date: r.published_at.slice(0, 10), url: r.html_url, notes: marked.parse(r.body ?? "", { async: false }) as string }));
  writeFileSync("data/changelog.json", JSON.stringify(out, null, 2) + "\n");
  console.log(`changelog: ${out.length} releases`);
} catch (e) {
  console.warn(`changelog: kept the existing data/changelog.json (${(e as Error).message})`);
}
