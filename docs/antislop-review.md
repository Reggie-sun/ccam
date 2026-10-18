# antislop review: ccam.dev

Mode: after (audit, then fix), chosen by the owner on 2026-10-19. Scope (docs/WEBSITE_PLAN.md 11.1): everything under `site/` that is code, styles, `content/en/` and `messages/en.json`. Translations are checked against the fixed English source.

Priority follows the rule tier: Hard Gate = HIGH, Purpose-Gate = MEDIUM, Quality Locks = LOW.

## Audit 001, 2026-10-19

| # | Priority | Rule | File | Finding | Action |
|---|---|---|---|---|---|
| 1 | HIGH | R-25 | `app/globals.css` | Light-mode amber `#d9921e` is 2.51:1 on the page (step numbers, focus ring) and 2.31:1 against the gauge track, below 3:1. | Fixed: amber `#a86a0c` in light mode (4.29:1 on the page, 3.94:1 on the track); light focus ring uses teal (7.66:1). Dark mode already passed (amber 8.34:1). |
| 2 | HIGH | R-17, R-38 | `components/UsageBoard.tsx` | The hero board shows usage numbers (72%, 8%, 100%) and reset times without saying they are an example. | Fixed: the caption now says it is example output of `ccs usage`. |
| 3 | HIGH | R-23 | `components/Icons.tsx`, `public/favicon.svg`, `scripts/og.ts` | A logo mark (three stacked gauges) was drawn without the owner asking for one, and used for the favicon, OG images and `/icon-512.png`. | Awaiting owner decision: keep the mark, or use the text wordmark only. |
| 4 | HIGH | R-37 | site | No `DESIGN.md` or owner-supplied direction; the palette, type and layout were chosen by the agent. | Fixed as far as the agent can: direction written down in `site/DESIGN.md` and marked agent-supplied. Awaiting owner confirmation. |
| 5 | HIGH | R-24, R-26 | header, footer, docs | Links to `github.com/jhyunwoo/ccam` (nav, footer, edit links, releases, `.deb` downloads) point to a private repository, so visitors get a 404. `ccam.dev/install.sh` needs a public release. | Launch blocker, owner decision: make the repository public and publish a release before deploying. |
| 6 | HIGH | R-38 | `lib/jsonld.ts`, docs meta line | Pages say "Checked with ccam 0.1.0" and JSON-LD lists version 0.1.0, but no 0.1.0 release exists yet. | Launch blocker, owner decision: tag v0.1.0 before deploying (the version is set in one place, `lib/jsonld.ts`). |
| 7 | LOW | R-05 | `components/Landing.tsx` | "How it works" is three numbered steps, and the third ("Share what should be shared") is a feature rather than a step; it is followed by a 2×2 grid of four same-shaped feature blocks. Both are template sections. | Fixed: one section built from the real commands in the order people use them (add, switch, share, check usage, continue, repair), each with what it does. |
| 8 | LOW | frontend direction | `components/Shell.tsx` | Footer links joined with a middle dot ("Source on GitHub · All languages"). | Fixed: separate links. |
| 9 | LOW | R-31 | site | Major decisions had no written one-line reasons. | Fixed: reasons recorded in `site/DESIGN.md`. |
| 10 | LOW | antislop-code | `scripts/finalize.ts`, `scripts/og.ts`, `i18n/locales.ts`, `lib/content.ts` | Several header comments run four to six lines around facts that fit in two. | Fixed: trimmed to the constraint. |
| 11 | MEDIUM | R-36, C-5 | `messages/en.json`, FAQ | "adds about 4 ms" needs a source. | Kept with its evidence: measured with `scripts/perf-check.sh` (p95 4.2–4.4 ms, see docs/PLAN.md decision log); the FAQ now says it was measured. |
| 12 | HIGH | R-35 | site | No recorded click-through of the interactive elements. | Fixed: `tests/site.spec.ts` clicks every control in both themes and on mobile; results below. |

Checked and passing: no em dashes in agent-written text (the one in `quick-start.md` is ccam's own Korean output, quoted verbatim); no R-16 buzzwords; no generic CTAs; FAQ questions are about ccam specifically (R-28); the changelog has an empty state that names the next step (R-27); no fake testimonials, logo bars or statistics; dark mode passes every contrast pair.

## Audit 002, 2026-10-19 (changes after audit 001 and the 30 translations)

| # | Priority | Rule | File | Finding | Action |
|---|---|---|---|---|---|
| 13 | HIGH | R-32 | `components/InstallCommand.tsx`, `lib/content.ts` | Code blocks, tables and the install command scroll sideways on narrow screens but could not be reached by keyboard (axe `scrollable-region-focusable`, found in CI). | Fixed: `tabindex="0"` on scrollable `pre`, `table` and the install command. |
| 14 | HIGH | R-24 | `lib/content.ts` | Heading anchors came from the translated heading text, so `/docs/install#uninstall` led nowhere on 29 of 30 languages (reported by several translators). | Fixed: anchors come from the English heading at the same position; `check-links` now verifies `#fragments`; `i18n-check` requires the same heading structure. |
| 15 | MEDIUM | R-32 | `components/DocPage.tsx` | Sidebar and pager navigation could get the same accessible name in some languages (French, Hindi). | Fixed: the pager is named after its links ("Previous, Next"). |
| 16 | LOW | C-4 | `scripts/client/site.js` | The language banner inherited the page direction, so an English suggestion on an Arabic page put its period on the wrong side. | Fixed: the banner takes the direction of the suggested language. |
| 17 | LOW | R-31 | `scripts/i18n-check.ts` | The command-name rule also matched the translatable table header `docs.commandsTable.command`, so 17 languages kept "Command" in English. | Fixed: the rule only covers `landing.flow[].command/also`; the headers were translated. |

Translations: all 29 were made against the reviewed English source, pass `i18n-check --strict` (code, links, glossary, placeholders, headings, source hash), and each translator back-translated three paragraphs; drifts they found are fixed. Pages are marked `translation: machine` until a native speaker reviews them.
