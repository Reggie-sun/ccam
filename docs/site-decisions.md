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
- **CSP**: hash-based; inline scripts, the inlined stylesheet and the few `style` attributes are hashed at build. The `'inline-speculation-rules'` keyword was not accepted by the tested Chromium, so speculation rules are hashed too.
- **Language coverage**: the plan's 30 languages reach 70.0% (conservative) to 85.6% (optimistic) of the world population by Unicode CLDR data. The conservative estimate cannot reach 80%: adding the 12 best further languages only raises it to 72.3%, because CLDR's per-country language shares are literacy-adjusted and below 100% even for national languages. The site publishes both numbers on `/languages`.
- **Cache purge after deploy**: not needed; Workers Static Assets serves each deployment atomically. HTML `s-maxage` only affects downstream caches.
- **Changelog**: built from GitHub Releases; empty until the first release (the repository is private and has no release yet).
- **Browser tests**: run in CI. The WSL environment used for development does not produce frames in headless Chromium (`requestAnimationFrame` never fires, even on a blank page), so screenshots and frame-dependent checks time out locally.

## Redesign, 2026-10-19

The owner replaced the agent-written direction with `site/DESIGN.md` from `npx getdesign@latest add elevenlabs` (an analysis of ElevenLabs' marketing site). The site follows its tokens; these are the places it departs, and why.

- **Display font**: Waldenburg is licensed. The file's substitute, EB Garamond, has no weight 300, so the display face is Newsreader pinned to weight 300 at optical size 60 (21 KB, Latin-1). Body is Inter 400–600 (27 KB, Latin-1), as the file says. Only Newsreader is preloaded because the hero headline is the largest paint.
- **Contrast**: the file's muted text `#777169` (4.43:1 on `#f5f5f5`) and error `#dc2626` (4.43:1) miss AA, so text uses `#6b655e` (5.28:1) and `#b91c1c` (5.93:1). Body copy that can sit over an orb uses `#44403c`; the worst case, full-strength lavender, still gives 4.5:1 or more for the lighter `#4e4e4e`.
- **Dark mode**: the file has no dark theme; dark tokens come from its own dark surfaces (`#0c0a09`, `#1c1917`) and the orbs drop to 30% opacity. The theme toggle stays.
- **Logo**: the agent-drawn gauge mark is gone from the header and OG images; the wordmark is "ccam" in the display face, as in the file's top nav. The favicon still uses the old mark (antislop review #3).
- **Motion**: the owner asked for striking animation on the landing page. Every animation moves transform, clip-path or filter, never the opacity of text, so contrast checks see full colour mid-animation; scroll reveals use CSS scroll-driven animations where supported and are simply absent elsewhere; everything stops under `prefers-reduced-motion`.
- **Docs**: three columns on wide screens (sections, article, on this page); below 80rem the table of contents folds into the article, below 60rem the section list becomes a sticky disclosure. Closed `details` are shown on desktop with `::details-content`, guarded by `@supports`.
