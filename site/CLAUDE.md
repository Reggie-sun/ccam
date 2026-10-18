# ccam.dev site

Read ../docs/WEBSITE_PLAN.md first. Decisions: ../docs/site-decisions.md. Design direction: DESIGN.md. antislop findings: ../docs/antislop-review.md.

## Rules
- Framework is vinext with Next.js public APIs only (no vinext-specific APIs)
- Every page is generated at build time; nothing may need request-time rendering
- No React client components: the static HTML ships without React's runtime (scripts/finalize.ts removes it). The three behaviours (copy button, language banner, theme toggle) live in scripts/client/site.js
- UI strings are server-only (messages/<locale>.json); never ship them to the client
- CSS uses logical properties only (RTL)
- Images, OG images, sitemaps and llms.txt come from the scripts in scripts/
- Terms in i18n/glossary.yaml are never translated; code blocks and inline code are copied exactly
- Every page footer says the site is unofficial
- Translations: set `translation: machine` and `sourceHash` (node scripts/i18n-check.ts --hash content/en/<section>/<slug>.md)

## Commands
- dev: pnpm dev
- build: pnpm build (vite build, then finalize, sitemap, llms, og, images)
- check: pnpm check (static, budget, hreflang, links, html, i18n)
- browser tests: pnpm exec playwright test (needs a display that renders frames; WSL headless does not)
- worker tests: node --test tests/worker.test.ts
- deploy: pnpm exec wrangler deploy --config deploy/wrangler.jsonc

## Working
- Write the check first, then make the page pass it
- Before finishing a phase, review changed files with the antislop skills and record findings in ../docs/antislop-review.md; rerun pnpm check afterwards
