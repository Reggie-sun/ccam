// Internal links, assets and #fragments resolve in the build. External links
// are checked separately (lychee in CI) to keep this check offline.
import { existsSync, readFileSync } from "node:fs";
import { join, relative } from "node:path";
import { fail, OUT, walk } from "./lib.ts";

const errors: string[] = [];
const ok = (path: string) => {
  const p = decodeURIComponent(path.split("#")[0].split("?")[0]);
  if (!p) return true;
  return existsSync(join(OUT, p)) && !p.endsWith("/") ? true : existsSync(join(OUT, p, "index.html"));
};
const ids = new Map<string, Set<string>>();
const idsOf = (path: string) => {
  const file = join(OUT, decodeURIComponent(path), "index.html");
  if (!ids.has(file)) ids.set(file, existsSync(file) ? new Set([...readFileSync(file, "utf8").matchAll(/\sid="([^"]+)"/g)].map((m) => m[1])) : new Set());
  return ids.get(file)!;
};
for (const f of (await walk(OUT)).filter((f) => f.endsWith(".html"))) {
  const html = readFileSync(f, "utf8");
  for (const [, attr, url] of html.matchAll(/\s(href|src)="(\/[^"]*)"/g)) {
    if (url.startsWith("//")) continue;
    if (!ok(url)) errors.push(`${relative(OUT, f)}: ${attr}=${url}`);
    const [path, hash] = url.split("#");
    if (hash && path.endsWith("/") && !idsOf(path).has(decodeURIComponent(hash))) errors.push(`${relative(OUT, f)}: no #${hash} on ${path}`);
  }
}
fail(errors, "links");
