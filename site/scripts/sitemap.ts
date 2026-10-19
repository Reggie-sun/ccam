// sitemap-index.xml, one sitemap per locale with xhtml:link alternates, and robots.txt.
import { mkdir, writeFile } from "node:fs/promises";
import { join } from "node:path";
import { getLocale } from "../i18n/locales.ts";
import { abs } from "../lib/paths.ts";
import { builtLocales, OUT, routes } from "./lib.ts";

const esc = (s: string) => s.replace(/&/g, "&amp;");
const all = routes();
await mkdir(join(OUT, "sitemaps"), { recursive: true });

for (const l of builtLocales()) {
  const urls = all.filter((r) => r.locale === l.code).map((r) => {
    const links = Object.entries(r.alternates)
      .map(([code, path]) => `    <xhtml:link rel="alternate" hreflang="${getLocale(code)!.lang}" href="${esc(abs(path))}"/>`)
      .concat(r.alternates.en ? [`    <xhtml:link rel="alternate" hreflang="x-default" href="${esc(abs(r.alternates.en))}"/>`] : []);
    return `  <url>\n    <loc>${esc(abs(r.path))}</loc>\n    <lastmod>${r.updated}</lastmod>\n${links.join("\n")}\n  </url>`;
  });
  await writeFile(
    join(OUT, "sitemaps", `${l.code}.xml`),
    `<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">\n${urls.join("\n")}\n</urlset>\n`,
  );
}

// A locale's sitemap changes when one of its pages does, not on every build.
const latest = (code: string) => all.filter((r) => r.locale === code).map((r) => r.updated).sort().at(-1);
const index = builtLocales().map((l) => `  <sitemap>\n    <loc>${abs(`/sitemaps/${l.code}.xml`)}</loc>\n    <lastmod>${latest(l.code)}</lastmod>\n  </sitemap>`);
await writeFile(join(OUT, "sitemap-index.xml"), `<?xml version="1.0" encoding="UTF-8"?>\n<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n${index.join("\n")}\n</sitemapindex>\n`);

// Search and AI crawlers are all welcome (docs/WEBSITE_PLAN.md 9.4).
await writeFile(join(OUT, "robots.txt"), `User-agent: *\nAllow: /\n\nSitemap: ${abs("/sitemap-index.xml")}\n`);
console.log(`sitemap: ${all.length} URLs in ${builtLocales().length} locale sitemaps`);
