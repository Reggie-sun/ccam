# ccam.dev Action Plan

Each item: why it matters (first principle), what it unblocks, how we'd know it failed, and a leading indicator to watch.

## Phase 1: Critical fixes (week 1)

### 1. Get indexed and own the brand name
- Do: verify ccam.dev in Google Search Console and Bing Webmaster Tools, submit /sitemap-index.xml, request indexing for `/`. Add "Claude Code Account Manager" to the homepage title or H1 subtext and as `alternateName` in Organization/SoftwareApplication JSON-LD. Set the GitHub repo description and topics (claude-code, cli, multi-account).
- Why: search engines cannot rank pages they have not indexed, and "CCAM" is an ambiguous acronym owned by older entities.
- Unblocks: every other ranking item, and real data for re-auditing (GSC queries, CrUX).
- Failed if: GSC "Pages" report shows fewer than about 400 indexed URLs after 4 weeks, or `site:ccam.dev` still shows other CCAMs.
- Watch: GSC indexed page count; impressions for "ccam".

### 2. Force HTTPS and add HSTS
- Do: Cloudflare "Always Use HTTPS" (or a redirect rule in site/deploy), then `Strict-Transport-Security: max-age=31536000; includeSubDomains`.
- Why: http and https are separate URLs to crawlers, and http responses are open to tampering, which matters for a site that serves `curl | sh`.
- Failed if: `curl -sI http://ccam.dev/ko/` returns anything other than 301 to https.
- Watch: no http URLs in GSC coverage.

## Phase 2: High-impact improvements (weeks 2-3)

### 3. Fix the FAQ JSON-LD text stripping
- Do: in site/lib/content.ts:119, strip Markdown syntax only (inline code backticks, emphasis markers, heading hashes, blockquote markers at line start), not every `_` and `>`.
- Failed if: `curl -s https://ccam.dev/en/docs/faq/ | grep CLAUDECONFIGDIR` still matches.

### 4. Fix duplicate breadcrumbs and add author/dates
- Do: section crumb should point to the section index (or be dropped) on blog and compare pages. Add `author` (Person with GitHub URL), `publisher` (Organization @id), and `datePublished` from front matter to TechArticle; use BlogPosting for the blog post. Ready-made JSON-LD is in findings/schema.md.
- Failed if: Rich Results Test reports breadcrumb or article warnings.

### 5. Match search intent for "switch account" queries
- Do: use the words people search ("switch", "switcher", "usage limit", "rate limit") on the homepage and usage-limit guide; add an install step to the guide; add a comparison page against other switchers (e.g. claude-swap) alongside /compare/aliases/.
- Depends on: item 1 (needs indexing to measure).
- Failed if: after 8 weeks, GSC shows impressions for these queries but CTR under 1%, or no impressions at all.
- Watch: GSC queries containing "switch".

### 6. Turn the changelog into release notes
- Do: per-version human summaries; drop SHAs and internal `fix(site)`/`test` commits. Consider marking non-English changelog URLs noindex or translating them, since all 30 show identical English text.
- Failed if: GSC flags changelog URLs as "Duplicate, Google chose different canonical".

## Phase 3: Content and authority (month 2)

### 7. Deepen the three guides
- Do: grow multiple-claude-code-accounts toward 1,000+ words with the `/logout` + `/login` route, Windows/WSL, checking which account is active, troubleshooting. Link to official Claude Code docs where cited. Use question-style headings.
- Failed if: average position for "multiple claude code accounts" does not improve within 8 weeks of indexing.

### 8. Earn the first links
- Do: GitHub README links, Homebrew tap / AUR / Scoop / Nix / pkg.go.dev listings, awesome-list PRs, then Show HN and r/ClaudeAI once there is a demo. Details in findings/backlinks.md.
- Watch: GitHub stars, referring domains in GSC "Links".

### 9. Real lastmod dates
- Do: bump `updated` only on real changes (or derive from git log); give each sitemap index entry its locale's latest page date instead of `new Date()`.

### 10. Small technical cleanups
- One Cache-Control header for /assets/* and /fonts/* (`public, max-age=31536000, immutable`).
- Fix /ar/ mobile CLS (size-adjusted fallback font or `font-display: optional`); add /ar/ to lighthouserc.json with CLS <= 0.1.
- Preload inter.woff2 or remove it from above the fold.
- Trim titles over 60 chars and descriptions over 160; make doc titles specific.
- Add apple-touch-icon.png; enlarge tap targets to 44 px.

## Phase 4: Monitoring (ongoing)

- Run `/seo drift baseline https://ccam.dev` now, then `/seo drift compare` after each deploy.
- Re-audit in 3 months with Search Console and a free Moz key for field CWV, index coverage and backlinks.
