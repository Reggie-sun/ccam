// Strips React's client runtime from the static export (pages have no client
// components), inlines small CSS, and writes _headers (cache, CSP) and _redirects.
import { createHash } from "node:crypto";
import { readdir, readFile, rm, writeFile, mkdir, stat } from "node:fs/promises";
import { join, relative } from "node:path";
import { gzipSync } from "node:zlib";
import { locales, localeAliases } from "../i18n/locales.ts";

export const OUT = "dist/client";
const INLINE_CSS_LIMIT = 15 * 1024; // gzip bytes (docs/WEBSITE_PLAN.md 6.2)

export async function walk(dir: string): Promise<string[]> {
  const out: string[] = [];
  for (const e of await readdir(dir, { withFileTypes: true })) {
    const p = join(dir, e.name);
    if (e.isDirectory()) out.push(...(await walk(p)));
    else out.push(p);
  }
  return out;
}

const sha256 = (s: string) => createHash("sha256").update(s).digest("base64");
const short = (s: string) => createHash("sha256").update(s).digest("hex").slice(0, 10);

async function main() {
  const files = await walk(OUT);
  const htmlFiles = files.filter((f) => f.endsWith(".html"));

  const siteJs = await readFile("scripts/client/site.js", "utf8");
  const siteJsPath = `/assets/site.${short(siteJs)}.js`;
  await mkdir(join(OUT, "assets"), { recursive: true });
  await writeFile(join(OUT, siteJsPath), siteJs);

  let cssHref = "";
  let css = "";
  const inlineScripts = new Set<string>();
  const styleHashes = new Set<string>();

  for (const file of htmlFiles) {
    let html = await readFile(file, "utf8");
    // React client runtime: module scripts, modulepreloads and RSC bootstrap/payload scripts.
    html = html.replace(/<link rel="modulepreload"[^>]*\/?>/g, "");
    html = html.replace(/<script[^>]*\bsrc="\/_next\/[^"]*"[^>]*><\/script>/g, "");
    for (const [, rules] of html.matchAll(/<script type="speculationrules">([\s\S]*?)<\/script>/g)) inlineScripts.add(rules);
    html = html.replace(/<script>([\s\S]*?)<\/script>/g, (whole, body: string) => {
      if (body.startsWith("try{var t=localStorage")) {
        inlineScripts.add(body);
        return whole;
      }
      return ""; // vinext navigation runtime and RSC payload
    });

    const pageStyles = new Set<string>();
    let inline = false;
    const link = /<link rel="stylesheet"[^>]*href="(\/_next\/static\/css\/[^"]+\.css)"[^>]*\/?>/.exec(html);
    if (link) {
      if (!css) {
        cssHref = link[1];
        css = await readFile(join(OUT, cssHref), "utf8");
      }
      inline = gzipSync(css).length <= INLINE_CSS_LIMIT;
      if (inline) {
        styleHashes.add(sha256(css));
        pageStyles.add(sha256(css));
      }
      html = html.replace(link[0], inline ? `<style>${css}</style>` : `<link rel="stylesheet" href="${cssHref}"/>`);
    }
    // style="..." attributes need 'unsafe-hashes'.
    for (const [, value] of html.matchAll(/ style="([^"]*)"/g)) pageStyles.add(sha256(value.replace(/&quot;/g, '"')));
    // Style hashes go in each page's own meta CSP: the motion diagrams have
    // too many for _headers, where a line may not exceed 2000 characters.
    // The browser enforces this policy and the header's, so the stricter one wins.
    const styleSrc = [inline ? "" : "'self'", ...[...pageStyles].map((h) => `'sha256-${h}'`), "'unsafe-hashes'"].filter(Boolean).join(" ");
    const charset = '<meta charSet="utf-8"/>';
    if (!html.includes(charset)) throw new Error(`${file}: no ${charset} to put the CSP after`);
    html = html.replace(charset, `${charset}<meta http-equiv="Content-Security-Policy" content="style-src ${styleSrc}"/>`);
    html = html.replace("</head>", `<script src="${siteJsPath}" defer></script></head>`);
    await writeFile(file, html);
  }

  // Drop what only the client runtime needed.
  for (const f of files) {
    const rel = relative(OUT, f);
    if (rel.endsWith(".txt") && rel !== "robots.txt" && !rel.startsWith("llms")) await rm(f);
    if (rel.startsWith("_next/static/chunks/") || rel.endsWith("_buildManifest.js") || rel.endsWith("_ssgManifest.js")) await rm(f);
  }
  await rm(join(OUT, ".vite"), { recursive: true, force: true });
  await rm(join(OUT, "vinext-client-entry-manifest.json"), { force: true });
  if (styleHashes.size) await rm(join(OUT, cssHref), { force: true });

  const scriptHashes = [...inlineScripts].map((s) => `'sha256-${sha256(s)}'`).join(" ");
  const csp = [
    "default-src 'self'",
    // Cloudflare Web Analytics, injected at the edge (plan section 9); it reports to same-origin /cdn-cgi/rum.
    `script-src 'self' https://static.cloudflareinsights.com ${scriptHashes}`,
    // Styles are limited by each page's meta CSP; a hash here would make browsers ignore 'unsafe-inline'.
    "style-src 'self' 'unsafe-inline'",
    "img-src 'self' data:",
    "font-src 'self'",
    "connect-src 'self'",
    "base-uri 'none'",
    "form-action 'none'",
    "frame-ancestors 'none'",
    "object-src 'none'",
  ].join("; ");

  await writeFile(
    join(OUT, "_headers"),
    `/*
  Content-Security-Policy: ${csp}
  Strict-Transport-Security: max-age=31536000; includeSubDomains
  X-Content-Type-Options: nosniff
  Referrer-Policy: strict-origin-when-cross-origin
  Permissions-Policy: interest-cohort=(), browsing-topics=()
  Cache-Control: public, max-age=0, s-maxage=86400, stale-while-revalidate=604800
  Link: </fonts/newsreader.woff2>; rel=preload; as=font; type=font/woff2; crossorigin

/assets/*
  Cache-Control: public, max-age=31536000, immutable
/_next/static/*
  Cache-Control: public, max-age=31536000, immutable
/fonts/*
  Cache-Control: public, max-age=31536000, immutable
/og/*
  Cache-Control: public, max-age=86400
/*.md
  Content-Type: text/markdown; charset=utf-8
/llms.txt
  Content-Type: text/plain; charset=utf-8
/llms-full.txt
  Content-Type: text/plain; charset=utf-8
`,
  );

  // Locale aliases and the duplicate English landing page.
  const codes = new Set(locales.map((l) => l.code));
  const lines = ["/en/ / 301", "/en / 301", "/index.html / 301"];
  // html_handling adds a missing trailing slash with a 307; search engines
  // should see a permanent redirect instead. Static rules go before the
  // splats: any rule after a splat counts against the 100 dynamic rules.
  for (const f of htmlFiles) {
    const dir = relative(OUT, f).replace(/\/?index\.html$/, "");
    if (!f.endsWith("index.html") || !dir || dir === "en" || dir === "404") continue;
    lines.push(`/${dir} /${dir}/ 301`);
  }
  for (const [alias, code] of Object.entries(localeAliases)) {
    if (codes.has(code)) lines.push(`/${alias}/* /${code}/:splat 301`);
  }
  for (const l of locales) {
    if (l.code !== l.code.toLowerCase()) continue;
    const upper = l.code.replace(/-([a-z]+)$/, (_, r: string) => "-" + r.toUpperCase());
    if (upper !== l.code) lines.push(`/${upper}/* /${l.code}/:splat 301`);
  }
  await writeFile(join(OUT, "_redirects"), lines.join("\n") + "\n");

  // IndexNow proves site ownership with a key file at the root.
  const key = process.env.INDEXNOW_KEY;
  if (key && /^[a-zA-Z0-9-]{8,128}$/.test(key)) await writeFile(join(OUT, `${key}.txt`), key);

  const cssSize = css ? gzipSync(css).length : 0;
  console.log(`finalize: ${htmlFiles.length} pages, CSS ${cssSize} B gzip (${styleHashes.size ? "inlined" : "linked"}), ${siteJsPath} ${gzipSync(siteJs).length} B gzip`);
  void stat;
}

if (import.meta.main) await main();
