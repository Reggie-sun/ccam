# Performance audit: https://ccam.dev (lab only, no CrUX)

Lighthouse 13.5.0, single run per page (Chromium 1248 headless, simulated throttling on mobile). Date 2026-10-19.
**Performance score: 92/100** (lab scores 94-99; deductions for the /ar/ CLS, conflicting cache headers, no field data).

| Page | Form | Perf | LCP | FCP | CLS | TBT | TTFB (LH) | Transfer |
|---|---|---|---|---|---|---|---|---|
| / | mobile | 98 | 2.1s | 1.9s | 0 | 0 | 95ms | 82KB |
| / | desktop | 99 | 0.8s | 0.7s | 0 | 0 | 172ms | 82KB |
| /en/docs/quick-start/ | mobile | 97 | 2.2s | 2.0s | 0 | 0 | 86ms | 83KB |
| /en/docs/quick-start/ | desktop | 98 | 1.0s | 0.8s | 0.001 | 0 | 260ms | 82KB |
| /en/docs/commands/ | mobile | 98 | 2.1s | 1.9s | 0.002 | 0 | 210ms | 84KB |
| /en/docs/commands/ | desktop | 99 | 0.8s | 0.7s | 0 | 0 | 102ms | 83KB |
| /ar/ (RTL) | mobile | 94 | 2.1s | 1.9s | **0.111** | 30ms | 87ms | 56KB |
| /ar/ | desktop | 98 | 0.9s | 0.9s | 0.044 | 0 | 506ms | 55KB |

CWV status (lab): LCP pass everywhere (<2.5s). INP proxy TBT 0-30ms, pass. CLS pass except /ar/ mobile (0.111, "needs improvement", just over 0.1).

## Findings
1. /ar/ mobile CLS 0.111: whole `main#main` shifts (single shift, score 0.111). Likely font swap (Inter/Newsreader or Arabic fallback metrics differ) reflowing the page. Fix: size-adjust/ascent-override fallback for the Arabic stack, or preload the font actually used for Arabic, `font-display: optional` for the non-LCP font. Re-test with several runs.
2. Cache-Control is duplicated/conflicting on static assets (/assets/site.*.js, /fonts/*.woff2): `public, max-age=0, s-maxage=86400, stale-while-revalidate=604800, public, max-age=31536000, immutable`. The Cloudflare default header is merged with the intended immutable header; browsers may honor the first max-age=0 and revalidate hashed files every visit. Fix: remove the generic header rule for /assets/* and /fonts/* (Worker/_headers: set only one value, with `Cache-Control` overwritten not appended). Not flagged by Lighthouse because the first run is cold; repeat-visit cost only.
3. HTML: `max-age=0, s-maxage=86400, swr=604800`. Good edge caching (cf-cache-status HIT, TTFB ~0.2s curl from here). Cache MISS on /ar/ cost TTFB 0.56s (cold edge); the 506ms desktop /ar/ TTFB reflects this. Consider cache warming or longer tiered cache for the less-visited locales.
4. Compression: Brotli on HTML (home 15.7KB wire, 18KB with headers). OK.
5. Render-blocking: none (render-blocking insight score 1). Inline CSS (hash-pinned CSP), one `defer` script (3.9KB), speculationrules. No render-blocking issues.
6. Fonts: newsreader.woff2 (22KB) preloaded with crossorigin; inter.woff2 (28KB) discovered late (not preloaded, requested at ~590ms). LCP element is the hero `p.lede`; LCP breakdown: TTFB ~190-270ms, element render delay ~230-255ms (font-gated). Fixes: preload inter.woff2 if used above the fold, or make it `font-display: swap/optional`; the render delay is the main LCP lever on mobile (LCP 2.1s is mostly simulated 4G latency).
7. Images: no raster images; only favicon.svg (1.1KB). No image weight issue. Total page 82KB.
8. Third party: Cloudflare Web Analytics beacon.min.js (10.5KB, 24h cache, ships legacy JS polyfills ~10KB wasted) plus /cdn-cgi/rum request. It sits on the request chain (longest chain 670ms) but is non-blocking. Optional: drop it or rely on edge-side analytics to save a connection and 10KB.
9. Main thread: TBT 0-30ms, tiny JS. INP risk is low.
10. Lighthouse CI config (site/lighthouserc.json) requires perf >=0.95 on mobile for 4 pages; current results meet it (97-98). /ar/ is not covered by it; add /ar/ (and assert `cumulative-layout-shift` <= 0.1) to catch RTL regressions.

## Priority fixes
1. Fix /ar/ font-swap CLS (high, only CWV failure).
2. Fix duplicated Cache-Control on hashed assets/fonts (medium, repeat-visit speed).
3. Preload inter.woff2 or trim to one font above the fold; reduce element render delay (medium, LCP).
4. Warm edge cache for non-default locales (low-medium, TTFB).
5. Consider removing CF beacon (low).
