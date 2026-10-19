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
// Cloudflare rejects the whole deploy when a _headers line is longer than this.
if (existsSync(`${OUT}/_headers`)) {
  readFileSync(`${OUT}/_headers`, "utf8").split("\n").forEach((line, i) => {
    if (line.length > 2000) errors.push(`_headers line ${i + 1} is ${line.length} characters (Cloudflare allows 2000)`);
  });
}
// Cloudflare silently skips _redirects rules past its limits: 2000 static and
// 100 dynamic, where a static rule after a splat or placeholder counts as dynamic.
if (existsSync(`${OUT}/_redirects`)) {
  let statics = 0;
  let dynamics = 0;
  for (const line of readFileSync(`${OUT}/_redirects`, "utf8").split("\n").filter(Boolean)) {
    if (/[*:]/.test(line.split(" ")[0]) || dynamics) dynamics++;
    else statics++;
  }
  if (statics > 2000) errors.push(`_redirects has ${statics} static rules (Cloudflare allows 2000)`);
  if (dynamics > 100) errors.push(`_redirects has ${dynamics} rules after the first splat (Cloudflare allows 100); put static rules first`);
}
fail(errors, "static");
