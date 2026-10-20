# ccam.dev decisions

Decisions taken while building the site, where they differ from or add to docs/WEBSITE_PLAN.md.

## Phase 0: checked on 2026-10-19

| Question | Finding | Decision |
| --- | --- | --- |
| vinext version | 1.2.1 (2026-10-15), React 19.3.1, Vite 8 | Use it; checked by building a three-locale export. |
| App Router `[locale]` + `output: "export"` | Works: every locale and page is pre-rendered to `dist/client/<path>/index.html` with `trailingSlash: true`. | Use `generateStaticParams` + `dynamicParams = false` (the plan's choice). |
| Client JavaScript | vinext's export ships the React runtime and client router on every page: about 135 KB gzip, five times the 25 KB budget. | `scripts/finalize.ts` removes React's runtime from the exported HTML. Pages have no client components; the plan's three islands are one 0.9 KB script. The plan already uses Speculation Rules and view transitions instead of a client router. |
| Deployment | `wrangler deploy` with Static Assets and a small Worker (the plan's default). The vinext static-export example keeps `wrangler.jsonc` outside the Vite root, so it lives in `site/deploy/`. | `html_handling: "force-trailing-slash"` (matches the export) instead of `auto-trailing-slash`. |
| `@vinext/cloudflare deploy` comparison | Not run: needs Cloudflare login, which had expired. | Revisit after login; nothing in the site depends on the choice. |
| Locale aliases (`/zh-CN/` etc.) | Workers Static Assets reads `_redirects`. | Generated `_redirects` instead of Worker code; the Worker only serves `/install.sh` and `/latest`. |
| `/en/` vs `/` | Both would carry the same landing page. | `/en/` redirects to `/` (301); `/` is the English page and `x-default`. |
| Search keywords per language | Drafted from the plan's examples for en, ko, ja, zh-hans, es, vi, fr. Not checked against live search results. | Refine with Search Console data after launch. |
| Cloudflare AI crawler settings | Needs the dashboard (login). | Apply plan 9.4 when deploying. |

## Later decisions

- **Content format**: Markdown (`.md`) with front matter instead of MDX; pages need no components, and the same files become the `.md` copies for LLMs.
- **OG images**: SVG rendered with resvg instead of satori, because satori does not shape Arabic, Urdu or Indic scripts. Fonts are Noto, downloaded at build.
- **Fonts**: superseded by the 2026-10-19 redesign below. Latin-script locales use Latin-1 subsets; Vietnamese, Turkish, Polish, Filipino and Hausa use the system font so letters outside the subset never fall back mid-word.
- **CSP**: hash-based; inline scripts, the inlined stylesheet and the few `style` attributes are hashed at build. The `'inline-speculation-rules'` keyword was not accepted by the tested Chromium, so speculation rules are hashed too. Script hashes sit in `_headers`; style hashes sit in each page's own `<meta http-equiv="Content-Security-Policy">`, because the motion diagrams' 200-odd `style` hashes made the header line about 12,000 characters and Cloudflare rejects a `_headers` line over 2000 (check-static enforces it). The header's `style-src 'self' 'unsafe-inline'` carries no hash, so browsers honour it, and since both policies apply, the page's hashes are what limit styles.
- **Language coverage**: the plan's 30 languages reach 70.0% (conservative) to 85.6% (optimistic) of the world population by Unicode CLDR data. The conservative estimate cannot reach 80%: adding the 12 best further languages only raises it to 72.3%, because CLDR's per-country language shares are literacy-adjusted and below 100% even for national languages. The site publishes both numbers on `/languages`.
- **Cache purge after deploy**: not needed; Workers Static Assets serves each deployment atomically. HTML `s-maxage` only affects downstream caches.
- **Changelog**: built from GitHub Releases at build time; v0.1.0 (2026-10-19) is the first.
- **SEO audit fixes (2026-10-19)**: findings in `ccam.dev-audit/`.
  - The language banner is fixed to the bottom of the screen: shown after first paint at the top, it pushed the page down (CLS 0.111 on `/ar/` for an English browser).
  - `published` lives in English front matter only; translations read it from English. `metaTitle` (optional) is the search-result title when the short title used in headings and the sidebar says too little.
  - The changelog shows only `feat`, `fix` and `perf` commits outside the site, without hashes; prose a maintainer writes above goreleaser's list is kept. `.goreleaser.yaml` filters the same way for future release notes.
  - `check-html` fails titles wider than 60 and descriptions outside 70 to 160 (CJK counts double), breadcrumbs that repeat a URL, articles without author, publisher or dates, and FAQ answers that lose inline code.
  - No `Content-Signal` line in robots.txt: it is a draft, and Lighthouse's robots.txt audit may flag unknown directives.
- **SEO audit, 2026-10-20**: `html_handling` adds a missing trailing slash with a 307, so `_redirects` lists a 301 for every page without its slash. These static rules come before the locale-alias splats, because Cloudflare counts any rule after a splat against the 100 dynamic rules and skips the rest (check-static enforces it). Page titles read "<page> | ccam for Claude Code" (`site.titleSuffix`, localized) unless the title already names Claude Code or would pass 60 wide. `_headers` sends HSTS without `preload`. `apple-touch-icon.png` and `favicon.ico` are rendered from `favicon.svg` beside `icon-512.png`. Redirecting `http://` to `https://` is the zone's "Always Use HTTPS" setting, not something the deploy can set.
- **Merging the two SEO fixes, 2026-10-20**: the 2026-10-19 branch and the 2026-10-20 commit fixed overlapping findings.
  - Titles: `pageTitle` adds `site.titleSuffix` to `metaTitle` (or `title`) and drops it when the result is wider than 60. It measures width like `check-html` (`lib/width.ts`, CJK counts double), so a built title never fails the check.
  - HSTS stays without `preload`, as decided on 2026-10-20; the branch's second header with `preload` was dropped.
  - Section hubs: `/<locale>/<section>/` lists that section's pages with their descriptions (`docs.hubs` in messages). The crumb above a page title and its `BreadcrumbList` point to the hub, and a guide's crumb says "Guides" rather than its sidebar group.
  - The commands table links each command to the feature page that documents it, and each "How it works" step on the landing page links to its feature page; the install label links to the install page.
  - Sidebar group titles are `<p>`, so the page's `<h1>` is its first heading; the sidebar `nav` has its own label.
  - The landing headline keeps "Claude Code" on one line: its words rise one by one and could break between the two.
  - Unchanged on purpose: Newsreader is preloaded on every page because the header wordmark uses it in every language; Inter is not preloaded (see "Display font" below); IndexNow needs `INDEXNOW_KEY` in CI (`scripts/indexnow.ts` exists).
- **Browser tests**: run in CI. The WSL environment used for development does not produce frames in headless Chromium (`requestAnimationFrame` never fires, even on a blank page), so screenshots and frame-dependent checks time out locally.

## Redesign, 2026-10-19

The owner replaced the agent-written direction with `site/DESIGN.md` from `npx getdesign@latest add elevenlabs` (an analysis of ElevenLabs' marketing site). The site follows its tokens; these are the places it departs, and why.

- **Display font**: Waldenburg is licensed. The file's substitute, EB Garamond, has no weight 300, so the display face is Newsreader pinned to weight 300 at optical size 60 (21 KB, Latin-1). Body is Inter 400–600 (27 KB, Latin-1), as the file says. Only Newsreader is preloaded because the hero headline is the largest paint. Both use `font-display: optional`: a font that misses the first paint is used from the next page on, so a late font never shifts the layout (the first redesign build lost Lighthouse points to exactly that). Two font files exceed the plan's one-font budget (6.4); the owner's direction names two typefaces, and together they are 48 KB.
- **Contrast**: the file's muted text `#777169` (4.43:1 on `#f5f5f5`) and error `#dc2626` (4.43:1) miss AA, so text uses `#6b655e` (5.28:1) and `#b91c1c` (5.93:1).
- **Dark mode**: the file has no dark theme; dark tokens come from its own dark surfaces (`#0c0a09`, `#1c1917`) The theme toggle stays and shows a sun in light mode, a moon in dark mode.
- **Logo**: the agent-drawn gauge mark is gone from the header and OG images; the wordmark is "ccam" in the display face, as in the file's top nav. The favicon is the owner's choice: three left-aligned gauges, top and bottom longest, middle shorter, in Claude orange `#d97757`, the only place that colour appears.
- **Orbs**: removed at the owner's request on 2026-10-19 (DESIGN.md "Owner overrides"); the site is flat canvas, white cards and hairlines.
- **Docs and guides**: one list in the sidebar and one "Docs" entry in the header; URLs stay `/docs/` and `/guides/` so links and search results keep working. The list is always open from 768 px up; on phones it is in the header menu.
- **Motion**: the owner asked for striking animation on the landing page. Every animation moves transform, clip-path or filter, never the opacity of text, so contrast checks see full colour mid-animation; scroll reveals use CSS scroll-driven animations where supported and are simply absent elsewhere; everything stops under `prefers-reduced-motion`.
- **Docs**: three columns on wide screens (sections, article, on this page); below 80rem the table of contents folds into the article, below 60rem the section list becomes a sticky disclosure. Closed `details` are shown on desktop with `::details-content`, guarded by `@supports`.

## Feature pages, 2026-10-19

- **Section**: one page per ccam feature under `/features/<slug>/`, listed as a third sidebar group ("Features", in every locale's `docs.sections`). The existing docs keep their task-oriented pages; feature pages hold every option and how the feature works inside, for developers.
- **Motion diagrams**: inline SVG drawn at build time by `lib/diagrams/<name>.ts` with the kit in `lib/motion.ts`, placed with `<!-- motion:<name> -->`. No script: the pieces share one CSS cycle (`--t`) and start at `--d`, moving only by transform and clip-path, so they follow the motion rule above and stand complete under `prefers-reduced-motion`. Timing travels in `style` attributes, which the CSP already hashes; a `<style>` per diagram would not be. The kit's CSS is about 1 KB in the shared stylesheet.
- **Diagram text is code only** (commands, paths, Go functions, syscalls): one diagram serves all 30 languages, and the page text around it explains it. `i18n-check` requires a translation to keep the English page's diagrams in order; the Markdown copies for LLMs drop the markers.
- **At most 400 units wide, flowing down**, so a phone shows a diagram at about full size; `figure()` refuses wider ones.
- **Review without frames**: `node scripts/motion-preview.ts [name] [--dark]` renders a diagram with resvg (dots at their destinations, highlights on), since headless Chromium under WSL draws no frames. The browser test that checks playback and reduced motion runs in CI.
