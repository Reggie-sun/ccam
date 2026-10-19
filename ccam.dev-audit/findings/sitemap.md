# Sitemap audit: https://ccam.dev (2026-10-19)

## Passes
- robots.txt references https://ccam.dev/sitemap-index.xml (200, application/xml). Index lists 30 sitemaps; all 30 fetch and parse as well-formed XML.
- 420 URLs total, 14 per locale, 0 duplicates. Far under the 50k / 50MB limits. No priority/changefreq tags.
- Every URL carries 31 xhtml:link alternates (30 locales plus x-default -> English). Namespace declared. Alternates are generated from one route table (scripts/sitemap.ts), so they are reciprocal.
- Sampled 42 sitemap URLs across all 30 locales (every 10th URL): 42/42 return 200, have a self-referencing canonical, and no noindex.
- Root https://ccam.dev/ is the English home (listed, canonical = itself, x-default). /en/ 301s to / and is NOT listed, so there is no duplicate.
- Coverage: 330 content files / 30 locales = 11 pages, plus home, changelog, languages = 14 per locale. Internal links sampled from / and /en/docs/faq/ all point to sitemap URLs.

## Findings

### 1. [Medium] lastmod is a uniform build-time value, not real modification dates
Evidence: all 420 URL lastmod = 2026-10-19; all 330 content front matter `updated:` = 2026-10-19 (site/content/<locale>/...md). Home/changelog/languages use the max `updated` of all pages (site/scripts/lib.ts:56-62). The index lastmod uses `new Date()` at build (site/scripts/sitemap.ts, `today`), so it changes every build.
Impact: Google learns to ignore lastmod when it is uniformly inaccurate; loses recrawl signal.
Fix: bump `updated` only when a page really changes (or derive from `git log -1 --format=%cs -- <file>`); give index entries the max lastmod of their locale rather than build date; keep home/changelog/languages dates tied to their own content.

### 2. [Low] /sitemap.xml returns 404
Evidence: https://ccam.dev/sitemap.xml -> 404 with noindex. Discovery via robots.txt works.
Fix: optionally 301 /sitemap.xml to /sitemap-index.xml; submit the index in Search Console and Bing Webmaster Tools.

### 3. [Info] Home URL asymmetry
English home is `/` while other English pages are `/en/...`. Consistent across canonical, hreflang and sitemap, and /en/ redirects. Keep internal links and the language switcher pointing to / (none to /en/ found on the home page).

### 4. [Info] Index cached s-maxage=86400
cache-control: public, max-age=0, s-maxage=86400 (Cloudflare HIT). New builds may take up to a day to show to crawlers; acceptable.

### 5. [Info] Scale
30 locales x 14 pages. Not location pages, so no doorway gate triggers; translation quality and uniqueness are a separate indexing risk.
