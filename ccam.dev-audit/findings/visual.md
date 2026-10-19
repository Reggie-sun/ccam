# Visual / Mobile / Images audit - ccam.dev (2026-10-19)
Screenshots: ccam.dev-audit/screenshots/{home,quickstart,guide,ar}-{desktop(1440),mobile(390)}-{fold,full}.png

## Images score: 88/100
The site has zero <img> elements and zero CSS background images. It is text, inline SVG and two woff2 fonts. Payload is tiny, there is no CLS from images, and there are no alt problems. Points are lost for the missing PNG favicon and apple-touch-icon, the PNG OG images, and no OG fallback.

## Findings
1. [Medium] No apple-touch-icon / PNG favicon. Only `<link rel="icon" href="/favicon.svg">`. /apple-touch-icon.png returns 404, and no manifest link. Fix: add a 180x180 apple-touch-icon.png, a 32x32 favicon.ico/png and a web manifest.
2. [Low] OG images are present on all 4 pages (1200x630, og:image:width/height/alt set, twitter:image set, per-locale e.g. /og/ar/home.png). Format is PNG, 26 KB for en/home, which is fine (<300 KB). Cache-Control has a duplicated value (`public, max-age=0, s-maxage=86400, ..., public, max-age=86400`). Fix: dedupe the header. Confirm twitter:card = summary_large_image.
3. [Info] Decorative/inline SVGs: 0 SVGs lack aria-hidden/label/title. No oversized images (no raster images in content).
4. [Info] Fonts: newsreader.woff2 21 KB and inter.woff2 27 KB, both 200. The newsreader font is requested twice on /ar/ and docs pages (two entries), so check for a duplicate preload/fetch.

## Above the fold
- Home 1440: H1 "Several Claude Code accounts. One terminal.", a 3-line sub-value-prop, "curl -fsSL https://ccam.dev/install.sh | sh" with Copy button at y~707, and "Read the quick start" secondary CTA. All visible, and the CTA hierarchy is clear (Copy is dark and primary). Pass.
- Home 390x844: H1, value prop, install box with Copy and quick start CTA all visible. The install command is clipped in the box ("https://ccam.c...") and is only reachable by horizontal scroll inside the box. [Low] Fix: allow wrap or use a shorter display, since Copy works.
- The hero value prop never names "Claude Code account switching" in the first sentence, but H1 covers it.

## Mobile
- No page-level horizontal scroll (scrollWidth == clientWidth = 390 on all 4 pages). Code blocks are the only elements wider than the viewport (right edge 447-548 px) and they sit inside scroll containers. OK.
- Base font 17px, with 11.5-13px for small UI text (footer/labels). [Low] Mostly fine.
- Hamburger present, with language and theme controls. Good.
- Tap targets: [Medium] on mobile, `ccam` logo link is 89x28, and 43px-wide anchors (height 94-238) appear to be the TOC/sidebar or footer links (width below 44). On desktop, nav links are 42px high, with 35-61 small targets per page. The Arabic page's banner link "Open in English" is 118x28 and the close button 33x44. Fix: min-height/width 44px (padding) on logo, banner link and close button.

## RTL (/ar/)
- html dir=rtl lang=ar. Layout mirrors correctly: logo on the right, hamburger on the left, banner "This page is also available in English" is shown. Code blocks stay LTR. No overflow. Pass.
- [Low] The banner text "This page is also available in English." is in English on the Arabic page, so localize it (it may have lang=en set deliberately, so verify lang attribute).
- Arabic H1 breaks as "Claude / .Code" with the period starting line 2 on mobile: acceptable bidi behaviour, but consider a bidi isolate (<bdi>/&lrm;) around "Claude Code".

## Content note
Quick start says "Checked with ccam 0.2.0", while the site is on 0.3.0 per git log. Update to avoid staleness.
