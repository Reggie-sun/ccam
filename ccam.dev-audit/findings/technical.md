# Technical + On-Page SEO findings: https://ccam.dev (audited 2026-10-19)

Scores: Technical SEO 88/100. On-Page SEO 84/100.
Method: raw fetch of the home, /en/, ko, ar pages. Fetched all 30 sitemaps (420 URLs) and parsed every page.

## Verified clean
- All 420 sitemap URLs return 200. There are no duplicates, and every page is self-canonical (canonical == sitemap URL).
- No robots meta noindex on indexable pages. Every page has exactly 1 h1, and there are no missing or duplicate descriptions.
- Hreflang: 31 alternates per page (30 locales + x-default) on all 420 pages. Self-reference and full reciprocity pass (0 failures). Every target is in the sitemap. Case is correct (zh-Hans, zh-Hant, pt-BR). Sitemap xhtml:link alternates match the HTML.
- The homepage DOES have hreflang. It is emitted as `hrefLang` (React casing). HTML is case-insensitive, so it is valid. A case-sensitive grep for "hreflang" misses it.
- RTL: ar, fa and ur pages all carry `lang` plus `dir="rtl"` and an RTL font group.
- Static HTML, fully server-rendered (JSON-LD, content and links in the source). Only one small deferred JS file, so no CSR risk.
- Viewport meta is present and theme-color is set. Pages are about 52-64 KB.
- Trailing slashes are consistent: links, canonicals and sitemaps all use a trailing slash. No-slash URLs 307 to the slash version. `/index.html` 307s to the folder.
- Custom 404 is a real 404 with `noindex`. /sitemap.xml is a 404 (fine, since the robots.txt Sitemap line points at sitemap-index.xml, which is a valid 200 index).
- OG/Twitter tags and per-page OG images are present. JSON-LD includes SoftwareApplication, WebSite and Organization. `/llms.txt` returns 200, and a per-page `.md` alternate is advertised.

## Findings

### HIGH-1: No HTTP to HTTPS redirect, and no HSTS
- Evidence: `curl -I http://ccam.dev/ko/` returns 200 (no 301). There is no `strict-transport-security` header on any response.
- Impact: Duplicate http/https URLs are crawlable. Users can be downgraded.
- Fix: In Cloudflare, enable "Always Use HTTPS" and HSTS (max-age=31536000; includeSubDomains once validated; preload optional). Alternatively add `/_headers` or a Worker rule. No file in site/public currently sets headers, and the CSP is set at the edge or in deploy config. Check /home/jhyunwoo/workspace/ccam/.claude/worktrees/seo-audit/site/deploy/wrangler.jsonc.

### MEDIUM-1: English URL scheme is split, and /en/ is a 301 to /
- Evidence:
  - English home is `/` (canonical, hreflang en, x-default). `/en/` 301s to `/`.
  - English inner pages live at `/en/...` (canonical `/en/docs/install/`). `/docs/install/` is 404.
  - Non-English home pages hreflang "en" to `/` but inner pages to `/en/...`.
  - This is internally consistent, with no duplicate content. But the pattern differs between the home page and inner pages. `/en/` is not in any alternate set, and no page links to it.
  - Home x-default is `/` while inner x-default is `/en/<path>`. Consistent per page set, and it works.
- Impact: Low risk today, but fragile. Link equity to `/en/` is lost if anyone links it, and the rule is easy to break.
- Fix (optional): Serve English inner pages at `/docs/...` as well, with `/en/...` redirecting to them, so the scheme is uniform. Otherwise document the exception. Source: /home/jhyunwoo/workspace/ccam/.claude/worktrees/seo-audit/site/lib/seo.ts and lib/paths.ts.

### MEDIUM-2: 45 titles over 60 characters, 22 descriptions over 160, 78 under 70
- Evidence: the longest title is uk `/guides/multiple-claude-code-accounts/` at 82 characters. Others at 73 are ta and sw. Non-Latin scripts count by chars, so pixel width differs. Many short descriptions are CJK or Indic, which is fine at lower char counts.
- Fix: Trim the 45 long titles to about 60 characters or fewer (check in the locale content frontmatter under site/content). Keep the brand suffix " | ccam". Only fix Latin and Cyrillic locales as priority.

### LOW-1: Duplicate titles in 3 locales
- Evidence: 6 title groups repeat, for example "Changelog | ccam" in de, fil and en. `fil` quick-start equals the English title, and ES "Idiomas | ccam" is duplicated.
- Fix: Localize or differentiate the titles (de/fil changelog, fil quick-start). Make sure each title is unique within a locale.

### LOW-2: og:locale is not exact for some locales
- Evidence: ar uses `ar_AR` (Facebook accepts it). Check fil, sw, ha, ur and fa values.
- Fix: Use valid ll_CC values (ar_AR is acceptable) and add og:locale:alternate only if wanted. Source: site/i18n/locales.ts (`og` field).

### LOW-3: Python-urllib default UA was blocked
- Evidence: the Python default UA got 403 on the first fetch, then worked with a browser-like UA. A later retry with `python-urllib/3` returned 200, so it is intermittent. Bot Fight or WAF may challenge unknown bots.
- Fix: Verify Cloudflare Bot Fight Mode / WAF is not blocking Googlebot, Bingbot or AI crawlers. Allow verified bots.

### INFO
- CSP is strict (hash-based, `frame-ancestors 'none'`). X-Content-Type-Options, Referrer-Policy and Permissions-Policy are present. No X-Frame-Options, but CSP covers it.
- cache-control is `max-age=0, s-maxage=86400, stale-while-revalidate`. That is fine.
- robots.txt is `Allow: /` plus a Sitemap line. It is valid and the sitemap resolves.
- IndexNow is not implemented (optional). Add a key file and ping Bing/Yandex/Naver on deploy.
- www.ccam.dev does not resolve (no DNS record). Add a www to apex redirect if www links exist.
- Core Web Vitals (source-based): the single preloaded font plus a deferred script and no images above the fold suggest good LCP, CLS and INP. Lab testing was not run.

## Score breakdown
Technical 88: crawlability 98, indexability 98, hreflang 100, security 70 (no HSTS or HTTPS redirect), URL structure 85, mobile 95, rendering 100.
On-Page 84: titles and descriptions are present and unique but over-length in places. H1 structure is clean. Internal links are consistent. Home is a hub linking all locales.
