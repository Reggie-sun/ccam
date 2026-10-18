// hreflang: valid codes, every alternate exists, links are reciprocal,
// x-default present, canonical points at the page itself.
import { existsSync, readFileSync } from "node:fs";
import { abs, SITE } from "../lib/paths.ts";
import { fail, htmlFile, routes } from "./lib.ts";

const errors: string[] = [];
const alts = new Map<string, Map<string, string>>(); // url → (hreflang → url)
for (const r of routes()) {
  const html = readFileSync(htmlFile(r.path), "utf8");
  const map = new Map<string, string>();
  for (const [, lang, href] of html.matchAll(/<link rel="alternate" hrefLang="([^"]+)" href="([^"]+)"/g)) map.set(lang, href);
  alts.set(abs(r.path), map);
  const canonical = /<link rel="canonical" href="([^"]+)"/.exec(html)?.[1];
  if (canonical !== abs(r.path)) errors.push(`${r.path}: canonical ${canonical}`);
  if (!map.has("x-default")) errors.push(`${r.path}: no x-default`);
  for (const lang of map.keys()) {
    if (lang === "x-default") continue;
    try {
      if (Intl.getCanonicalLocales(lang)[0].toLowerCase() !== lang.toLowerCase()) errors.push(`${r.path}: non-canonical hreflang ${lang}`);
    } catch {
      errors.push(`${r.path}: invalid hreflang ${lang}`);
    }
  }
}
for (const [url, map] of alts) {
  for (const [lang, href] of map) {
    const path = href.slice(SITE.length);
    if (!existsSync(htmlFile(path))) errors.push(`${url}: ${lang} → missing ${href}`);
    else if (lang !== "x-default" && ![...(alts.get(href)?.values() ?? [])].includes(url)) errors.push(`${url}: ${lang} → ${href} does not link back`);
  }
}
fail(errors, "hreflang");
