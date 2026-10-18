// Per-page size budgets (docs/WEBSITE_PLAN.md 6.4), all gzip:
// landing JS ≤ 25 KB, other pages' JS ≤ 15 KB, HTML+CSS ≤ 40 KB,
// at most one font request and one high-priority image.
import { readFileSync } from "node:fs";
import { join } from "node:path";
import { gzipSync } from "node:zlib";
import { fail, htmlFile, OUT, routes } from "./lib.ts";

const gz = (s: string | Buffer) => gzipSync(s).length;
const errors: string[] = [];
let worst = { path: "", js: 0, html: 0 };
for (const r of routes()) {
  const html = readFileSync(htmlFile(r.path), "utf8");
  let js = 0;
  for (const [, src] of html.matchAll(/<script[^>]*\bsrc="([^"]+)"/g)) js += gz(readFileSync(join(OUT, src)));
  let css = 0;
  for (const [, href] of html.matchAll(/<link rel="stylesheet" href="([^"]+)"/g)) css += gz(readFileSync(join(OUT, href)));
  const htmlCss = gz(html) + css;
  const jsLimit = r.key === "home" ? 25 * 1024 : 15 * 1024;
  if (js > jsLimit) errors.push(`${r.path}: JS ${js} B > ${jsLimit} B`);
  if (htmlCss > 40 * 1024) errors.push(`${r.path}: HTML+CSS ${htmlCss} B > 40960 B`);
  const fonts = (html.match(/rel="preload"[^>]*as="font"/g) ?? []).length;
  if (fonts > 1) errors.push(`${r.path}: ${fonts} font preloads`);
  if ((html.match(/fetchpriority="high"/gi) ?? []).length > 1) errors.push(`${r.path}: more than one high-priority image`);
  if (htmlCss > worst.html) worst = { path: r.path, js, html: htmlCss };
}
console.log(`budget: largest page ${worst.path} HTML+CSS ${worst.html} B, JS ${worst.js} B (gzip)`);
fail(errors, "budget");
