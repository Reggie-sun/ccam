# ccam.dev SEO Audit

Date: 2026-10-19 · URL: https://ccam.dev · Business type: open-source developer tool (SaaS-like docs site, no pricing)

## SEO Health Score: 81 / 100

| Category | Weight | Score | Source |
|---|---|---|---|
| Technical SEO | 22% | 88 | findings/technical.md |
| Content Quality | 23% | 74 | findings/content.md |
| On-Page SEO | 20% | 84 | findings/technical.md |
| Schema / Structured Data | 10% | 72 | findings/schema.md |
| Performance (CWV, lab only) | 10% | 92 | findings/performance.md |
| AI Search Readiness | 10% | 66 | findings/geo.md, findings/agentic.md |
| Images | 5% | 88 | findings/visual.md |

Not scored: backlinks (INSUFFICIENT DATA, only Common Crawl configured), SXO (diagnostic, see findings/sxo.md).

## Data limits

- No Google Search Console, CrUX, GA4, Matomo, Moz or DataForSEO credentials. Performance is lab-only (Lighthouse 13.5). Indexing status is inferred, not confirmed.
- PageSpeed Insights returned 429, so the Lighthouse Agentic Browsing fraction was not measured.
- The domain and repo are one day old (repo created 2026-10-18). Low authority and zero rankings are expected at this age.
- The site was redeployed during the audit. Findings about a stale 0.2.0 deploy (schema softwareVersion, "Checked with ccam 0.2.0", missing v0.3.0 changelog) were re-checked at the end and are resolved on the live site.

## What the site does well

- Static, server-rendered HTML. All 420 URLs (30 locales x 14 pages) return 200, are self-canonical and indexable, with exactly one h1 each.
- Hreflang is complete and reciprocal: 31 alternates per page (30 locales + x-default), in both HTML and sitemaps.
- RTL locales (ar, fa, ur) set `dir="rtl"` and render correctly.
- Lighthouse 94-99, LCP 2.1-2.2 s mobile, TBT 0-30 ms, pages 55-84 KB.
- Strong security headers (CSP etc.), real 404 with noindex.
- llms.txt plus `.md` mirrors served as `text/markdown`. All AI crawlers allowed.
- Readable, direct writing (Flesch 67-78). Pages answer in their first sentence. The CLAUDE_CONFIG_DIR guide has content no ranking page has.

## Critical

1. **ccam.dev is not findable for its own name.** `site:ccam.dev` and "ccam claude code" searches return CSIRO's climate model, the EU CCAM association and a competing tool (CAAM). "Claude Code Account Manager" appears only in the header. Indexing not yet confirmed (no Search Console). Evidence: findings/sxo.md.

## High

2. **No http to https redirect, no HSTS.** Verified: `http://ccam.dev/ko/` returns 200; no `Strict-Transport-Security` header. Fix in Cloudflare (Always Use HTTPS + HSTS) or site/deploy. findings/technical.md.
3. **Almost no outside footprint.** 0 stars, 0 forks, no HN/Reddit/YouTube mentions; the GitHub repo is the only known inbound link. findings/backlinks.md, findings/geo.md.
4. **Changelog is a raw goreleaser commit dump** with 40-char SHAs and internal `fix(site)`/`test` commits, identical English text on all 30 locale URLs, all in the sitemap. findings/content.md.
5. **Query-to-page mismatches.** "claude code switch account" and "usage limit switch account" SERPs reward switcher-tool pages; the homepage never says "switch"/"switcher", and the usage-limit guide (254-308 words) has no install step. /compare/aliases/ compares only against shell aliases while searchers compare 8+ tools. findings/sxo.md.

## Medium

6. **FAQ JSON-LD answer text is corrupted.** Verified live: "CLAUDE_CONFIG_DIR" appears as "CLAUDECONFIGDIR". Cause: site/lib/content.ts:119 strips every `` ` * _ > # `` character, not just Markdown syntax. FAQ rich results are retired (May 2026), so this is about text quality for parsers, not a SERP feature.
7. **Breadcrumb duplicates.** Verified: on /en/blog/introducing-ccam/ positions 2 and 3 have the same URL; same on compare pages. findings/schema.md.
8. **No named author or publish dates.** TechArticle has no `author`, `publisher` or `datePublished`; the blog post has no byline. findings/content.md, findings/schema.md.
9. **Thin guides.** multiple-claude-code-accounts 376-440 words vs ranking how-tos of about 1,100-1,600; usage-limit-workflow 254; claude-config-dir 313. Missing: `/logout` + `/login` route, Windows/WSL, checking the active account, troubleshooting. findings/content.md, findings/geo.md.
10. **lastmod is a uniform build stamp.** All 420 URLs and all 330 content files say 2026-10-19; the index uses `new Date()`. site/scripts/lib.ts:56-62, site/scripts/sitemap.ts. findings/sitemap.md.
11. **Conflicting Cache-Control on hashed assets.** Verified: fonts send `public, max-age=0, s-maxage=86400, stale-while-revalidate=604800, public, max-age=31536000, immutable`. findings/performance.md.
12. **/ar/ mobile CLS 0.111** (limit 0.1), one shift of `main#main`, likely font swap. findings/performance.md.
13. **Title and description lengths.** 45 titles over 60 chars (uk up to 82), 22 descriptions over 160, 78 under 70. Doc titles are generic ("Install ccam | ccam"). findings/technical.md.
14. **Tap targets under 44 px** (logo link 89x28, /ar/ "Open in English" 28 px tall). `/apple-touch-icon.png` is 404, no PNG/ICO fallback. findings/visual.md.
15. **No links to official Claude Code docs** where the text cites them. findings/geo.md.

## Low

- English URL scheme is split: home at `/`, inner pages at `/en/...`, `/en/` 301s to `/`. Consistent today, fragile. site/lib/seo.ts, site/lib/paths.ts.
- Six repeated title groups (e.g. "Changelog | ccam" in de and fil).
- Cloudflare intermittently 403s the default Python user agent; confirm Bot Fight Mode / WAF does not touch verified crawlers.
- /sitemap.xml is 404 (discovery via robots.txt works).
- inter.woff2 is discovered late (about 590 ms); Newsreader loads twice on /ar/ and docs pages.
- /languages/ is 54 words on 30 locales. Arabic "also available in English" banner is in English.
- Organization/WebSite have no `@id`; changelog has no page-level schema.
- No privacy note for the Cloudflare Web Analytics beacon; no SECURITY.md.

## Info

- Existing FAQPage on /docs/faq/ in all locales: kept as Info only. No SERP benefit since May 7, 2026; no recommendation to add or remove.
- SoftwareApplication is not eligible for the software rich result without ratings or reviews. Do not fabricate them.
- Optional agent extras not present: Content-Signal, `rel="alternate" type="text/markdown"`, ai-catalog.json, WebMCP, IndexNow. Agent-UX heuristic 100/100.

See ACTION-PLAN.md for the prioritized fix order.
