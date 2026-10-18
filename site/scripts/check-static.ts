// Every expected page exists as a static file and nothing needs a server:
// no route templates in the output, no client runtime left behind
// (stands in for Next.js's ensureStatic until vinext supports it).
import { existsSync, readFileSync } from "node:fs";
import { relative } from "node:path";
import { fail, htmlFile, OUT, routes, walk } from "./lib.ts";

const errors: string[] = [];
for (const r of routes()) if (!existsSync(htmlFile(r.path))) errors.push(`missing ${r.path}`);
for (const f of await walk(OUT)) {
  const rel = relative(OUT, f);
  if (/\[[^\]]+\]/.test(rel)) errors.push(`route template in output: ${rel}`);
  if (rel.startsWith("_next/static/chunks/")) errors.push(`client runtime chunk left: ${rel}`);
  if (rel.endsWith(".html")) {
    const html = readFileSync(f, "utf8");
    if (/src="\/_next\//.test(html) || html.includes("vinext.navigationRuntime")) errors.push(`client runtime referenced in ${rel}`);
  }
}
for (const f of ["404.html", "_headers", "_redirects", "robots.txt", "sitemap-index.xml", "llms.txt", "llms-full.txt"]) {
  if (!existsSync(`${OUT}/${f}`)) errors.push(`missing ${f}`);
}
fail(errors, "static");
