# antislop review: ccam.dev

Mode: after (audit, then fix), chosen by the owner on 2026-10-19. Scope (docs/WEBSITE_PLAN.md 11.1): everything under `site/` that is code, styles, `content/en/` and `messages/en.json`. Translations are checked against the fixed English source.

Priority follows the rule tier: Hard Gate = HIGH, Purpose-Gate = MEDIUM, Quality Locks = LOW.

## Audit 001, 2026-10-19

| # | Priority | Rule | File | Finding | Action |
|---|---|---|---|---|---|
| 1 | HIGH | R-25 | `app/globals.css` | Light-mode amber `#d9921e` is 2.51:1 on the page (step numbers, focus ring) and 2.31:1 against the gauge track, below 3:1. | Fixed: amber `#a86a0c` in light mode (4.29:1 on the page, 3.94:1 on the track); light focus ring uses teal (7.66:1). Dark mode already passed (amber 8.34:1). |
| 2 | HIGH | R-17, R-38 | `components/UsageBoard.tsx` | The hero board shows usage numbers (72%, 8%, 100%) and reset times without saying they are an example. | Fixed: the caption now says it is example output of `ccs usage`. |
| 3 | HIGH | R-23 | `components/Icons.tsx`, `public/favicon.svg`, `scripts/og.ts` | A logo mark (three stacked gauges) was drawn without the owner asking for one, and used for the favicon, OG images and `/icon-512.png`. | Partly resolved by the 2026-10-19 redesign: header, footer and OG images use the text wordmark (as the owner's DESIGN.md does). On 2026-10-19 the owner chose the favicon: three left-aligned gauges, top and bottom longest, the middle one shorter, in Claude orange `#d97757` (3.12:1 on a white tab, 5.16:1 on a dark one); `/icon-512.png` is rendered from it. Resolved. |
| 4 | HIGH | R-37 | site | No `DESIGN.md` or owner-supplied direction; the palette, type and layout were chosen by the agent. | Fixed as far as the agent can: direction written down in `site/DESIGN.md` and marked agent-supplied. Awaiting owner confirmation. |
| 5 | HIGH | R-24, R-26 | header, footer, docs | Links to `github.com/jhyunwoo/ccam` (nav, footer, edit links, releases, `.deb` downloads) point to a private repository, so visitors get a 404. `ccam.dev/install.sh` needs a public release. | Resolved 2026-10-19: the owner made the repository public; v0.1.0 is released and `ccam.dev/install.sh` installs it (checked on Linux x64/arm64, macOS arm64/x64 and the .deb). |
| 6 | HIGH | R-38 | `lib/jsonld.ts`, docs meta line | Pages say "Checked with ccam 0.1.0" and JSON-LD lists version 0.1.0, but no 0.1.0 release exists yet. | Resolved 2026-10-19: v0.1.0 is tagged and released. |
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

## Audit 003, 2026-10-19 (redesign from the owner's DESIGN.md)

Mode: after (audit, then fix), chosen by the owner for this redesign. Direction: `site/DESIGN.md`, generated by the owner with `npx getdesign@latest add elevenlabs`. Scope: every file the redesign changed (`app/globals.css`, `components/`, `scripts/client/site.js`, `scripts/og.ts`, `scripts/subset-font.ts`, `scripts/finalize.ts`, new UI strings).

| # | Priority | Rule | File | Finding | Action |
|---|---|---|---|---|---|
| 18 | HIGH | R-37, R-01, R-29 | `site/DESIGN.md`, `components/Orbs.tsx` | The owner's direction asks for pastel radial gradient orbs (mint, peach, lavender, sky, rose), a pattern R-01 names, and five orb colours exceed R-29's palette cap. | Applied as owner direction, and only where the file allows: atmosphere behind the hero, the demo stage and the "How it works" panel, never behind a control or as a component fill. Resolved 2026-10-19: the owner dropped the orbs; they are removed from every page and the OG images, and DESIGN.md records the override. |
| 19 | HIGH | R-25 | `app/globals.css` | DESIGN.md's muted text `#777169` (4.43:1 on `#f5f5f5`) and error `#dc2626` (4.43:1) miss AA. | Fixed: `#6b655e` (5.28:1) and `#b91c1c` (5.93:1). Body copy over orbs uses `#44403c`; worst case for `#4e4e4e` on full lavender is 4.51:1. Dark: muted `#a8a29e` 6.93:1 on cards, 7.83:1 on the canvas. |
| 20 | HIGH | R-32, R-34 | `app/globals.css` | Code blocks and the terminal are dark in both themes, so the light theme's ink focus ring (`#0c0a09`) was invisible on the copy buttons inside them. | Fixed: focus rings inside `.code` and `.term` use `--term-ink` (15.7:1 on `#0c0a09`). |
| 21 | HIGH | R-03 | `app/globals.css` | Sidebar, table-of-contents and language-list links were 30 to 36 px tall on touch screens. | Fixed: 44 px targets under `(pointer: coarse)`. Header controls, pills and footer links were already 44 px. |
| 22 | HIGH | R-26, R-32 | `components/Shell.tsx`, `scripts/client/site.js` | New controls: the mobile menu, the language popover, the docs section and on-this-page disclosures. | Pass: each is a native `details`, keyboard operable; header popovers close on Escape, a click elsewhere, or opening the other one; docs disclosures close on Escape and after following a link. The mobile menu and the on-this-page links have browser tests; the language picker and theme toggle keep theirs. |
| 23 | MEDIUM | R-01, DESIGN.md | `components/Landing.tsx`, `scripts/client/site.js` | Guide cards lit a lavender orb under the pointer. DESIGN.md: orbs are never component backgrounds. | Fixed: removed (CSS, the `data-spot` hook and its script). Cards keep the lift and arrow on hover. |
| 24 | MEDIUM | R-19 | `app/globals.css`, `scripts/client/site.js` | The owner asked for striking motion; each animation needs a purpose. | Reasons: headline rise (the one entrance moment); orb drift and pointer lean (DESIGN.md's atmosphere); the terminal replay (shows the product's job: see usage, open the account with room); the sticky "How it works" panel (keeps the command in view while its explanation scrolls); scroll reveal (paces long sections); caret blink (a prompt waiting for input). All are transform or opacity on non-text layers, all stop under reduced motion. |
| 25 | MEDIUM | R-10, R-12 | `app/globals.css` | Backdrop blur and lifted shadows. | Blur on two elements only (sticky header, sticky docs section bar on phones), because content scrolls under them. Raised shadow only on the install command, the terminal, popovers and hovered cards, the things that float above the page. |
| 26 | LOW | R-05 | `components/UsageBoard.tsx` | A terminal panel as a feature visual is a named pattern. | Kept with evidence: no window chrome; the commands are real; the accounts, percentages, reset times and "most room left: work" match the documented real output of `ccam usage` (`content/en/docs/usage-limits.md`); the caption labels it as example output. |
| 27 | LOW | R-06 | `app/globals.css`, `scripts/subset-font.ts` | Inter is on the default-font roster. | Kept: DESIGN.md names Inter for body. Display is Newsreader 300 because the file's Waldenburg is licensed and its suggested EB Garamond has no 300 weight (`docs/site-decisions.md`). |
| 28 | LOW | R-14 | `components/Landing.tsx` | Guide cards share one shape. | Kept: the three guides are equal-rank articles; the card is a large link target, not a feature tile. |
| 29 | LOW | antislop-code | `components/UsageBoard.tsx`, `scripts/subset-font.ts` | Two header comments ran three lines around one fact. | Fixed: one and two lines. |
| 30 | LOW | Lighthouse (C-4) | landing | The first CI run scored 0.88 on mobile performance: CLS 0.187 (orb layer resized when the web font swapped in) and 220 ms TBT (blur, clip-path and border animations on the main thread, Intl set-up at start). | Fixed: orb layer sized by the viewport, two orbs fewer on phones, transform-only entrances and reveals, opacity header hairline, formatter created on first play. |

Checked and passing: no em dashes in agent-written text; no buzzwords or generic CTAs ("Read the quick start", "Copy"); no badges, glows, grid backgrounds or illustrations; the ink pill is the only primary action per view; dark and light both verified (screenshots, axe in both themes); the 100% gauge shows its percentage, never colour alone; RTL (Arabic) checked: code stays left to right with the copy button beside it.

## Audit 004, 2026-10-19 (SEO audit fixes)

Scope: the English pages rewritten or added after the SEO audit (`guides/multiple-claude-code-accounts`, `guides/usage-limit-workflow`, `guides/claude-config-dir`, new `compare/account-switchers`, the last paragraph of `compare/aliases`, the FAQ link), new UI strings (`site.title`, `site.description`, `changelog.description`, `languages.description`, `docs.byline`, `footer.analytics`), `SECURITY.md`, and the code comments added in `lib/`, `scripts/`, `components/` and `app/globals.css`.

| # | Priority | Rule | File | Finding | Action |
|---|---|---|---|---|---|
| 31 | HIGH | R-36, R-38 | `content/en/compare/account-switchers.md` | Claims about other tools must be checkable. | Every row comes from the tool's own README or docs, read on 2026-10-19 (sources kept with the audit); the table says so and marks undocumented features as "Not documented", not "No". A draft line saying no other tool continues a conversation across folders was cut: CAAM's shallow profiles share `projects/`, so the claim could not be shown. Star counts left out because they go stale. |
| 32 | HIGH | R-36 | `content/en/guides/*.md` | Statements about Claude Code ("Claude Code's documentation describes", the limits, `/login`, `/status`, waiting for a reset, Windows support) had no source. | Each now links the official page it comes from (code.claude.com, support.claude.com). `/usage` was described as showing reset times; the documentation does not say that, so the sentence now says only what it does. |
| 33 | MEDIUM | R-36 | `content/en/compare/account-switchers.md` | Bold lead-ins ("**Windows without WSL.**", "**Menu bar apps.**") restated each item. | Fixed: plain sentences that name the tools. |
| 34 | MEDIUM | R-36 | `content/en/guides/usage-limit-workflow.md` | "starts Claude Code and gets out of the way" is a stock phrase. | Fixed: "replaces itself with it", which is what `syscall.Exec` does. |
| 35 | LOW | R-17 | `SECURITY.md` | A draft promised a reply within a week, a commitment for the owner to make. | Removed. The checksum sentence was checked against `install.sh` and `internal/update`. |

Checked and passing: no em or en dashes in the new English text or strings; no buzzwords; headings are the questions people search for; the comparison names where other tools do more (native Windows, automatic switching, menu bar apps, other agents) as plainly as where ccam does.

## Audit 005, 2026-10-20 (merging the SEO fixes)

Scope: the English text added on 2026-10-20 (`compare/aliases` setup sections, the feature-page paragraph in `docs/sharing`, two links in `blog/introducing-ccam`, the shorter `guides/remote-control` description, `docs.hubs`), the new hub page `app/[locale]/[section]/page.tsx`, and comments added in `components/`, `lib/width.ts` and `scripts/lib.ts`. The antislop skill is not installed in this checkout, so the rules were applied by reading.

| # | Priority | Rule | File | Finding | Action |
|---|---|---|---|---|---|
| 36 | MEDIUM | R-36 | `content/en/compare/aliases.md` | The alias recipe must match what the guide and Claude Code's documentation say, and `ccam add --shortcut` must exist. | Checked: the alias line is the one in `guides/multiple-claude-code-accounts`, the link is the documentation page it cites, and `--shortcut` is defined in `internal/cli/profiles.go`. |
| 37 | LOW | R-36 | `messages/en.json` (`docs.hubs`) | Hub descriptions list what each section holds; no claims beyond the pages listed. | Kept as written. |
| 38 | LOW | antislop-code | `components/DocPage.tsx`, `components/Landing.tsx` | The two maps (command to feature page, landing step to feature page) each carry a one-line comment saying what they map. | Kept: one line each, no restating of the code. |

Checked and passing: no em or en dashes in the new English text or strings; no buzzwords; translations keep code blocks, inline code and link targets (i18n-check `--strict` passes for all 29 languages).

## Audit 006, 2026-10-25 (auto_switch, --brief, the Agents tab)

Scope: the English text changed on 2026-10-25 (`features/next`, `launch`, `doctor`, `dashboard`, `usage`, `pick`; `docs/faq`, `commands`, `usage-limits`, `install`, `quick-start`; `guides/usage-limit-workflow`, `multiple-claude-code-accounts`; `compare/account-switchers`) and `commands.ccam` in `messages/en.json`. The rules were applied by reading, as in audit 005.

| # | Priority | Rule | File | Finding | Action |
|---|---|---|---|---|---|
| 39 | HIGH | R-38 | the 14 pages above, `lib/jsonld.ts` | The pages describe `auto_switch`, `ccam next --continue --brief` and the Agents tab, which no release has yet, while every page says "Checked with ccam 0.11.0" and the site deploys on each push to main. | Open: publish these pages with the release that contains the features (the release flow sets the version). Until then `auto_switch` is called experimental and off by default wherever it is named. |
| 40 | HIGH | R-36 | `features/next.md`, `docs/faq.md` | Quotes from Anthropic must be checkable. | Checked: the seat sentences are the help center's, as already on the site (2026-10-24); the Usage Policy sentence is in the English text of the version effective November 12, 2026 at https://www.anthropic.com/legal/aup (2026-10-25). |
| 41 | MEDIUM | R-36 | `features/next.md`, `features/dashboard.md`, `features/doctor.md` | Example output must be real. | The `auto_switch` and `--brief` lines and the Agents tab screen come from the built binary run with a fake `claude` and synthetic accounts in a scratch folder (the tab on an 80-column tmux terminal), the doctor lines from `ccam doctor` run the same way, and the notifications from the Go source. |

Checked and passing: no em or en dashes in the new English text; no buzzwords; the new FAQ question is about ccam; `auto_switch` is always written as code, so translations keep it.

## Audit 007, 2026-10-26 (English copyedit)

Scope: every English page in `site/content/en/` and the changed strings in `site/messages/en.json` (MSG, INS … BLG items and the shared rules P-01 to P-23 of the copyedit brief), the CLI messages the pages quote (C-01 to C-11 in `internal/`), and the plugin's descriptions. The rules were applied by reading, as in audit 005.

| # | Priority | Rule | File | Finding | Action |
|---|---|---|---|---|---|
| 42 | HIGH | R-36 | `content/en/docs/commands.md` | The `--dry-run` paragraph said every other command that changes something refuses the flag, but `folder`, `pick` and `ccusage` support it (their own pages say so). | Fixed from `dryRunCommands` in `internal/cli/root.go`: the supported commands are listed, read-only commands accept it, and `group rename` is named among the refusals. |
| 43 | HIGH | R-36 | `content/en/features/migrate.md` | The checklist example no longer matched the TUI (`?` prompt, `x` marks, `Enter show the plan`). | Replaced with what `tui.MigratePicker.View()` draws, rendered from the code. |
| 44 | MEDIUM | R-36 | `messages/en.json` (`landing.flow[5].body`) | "Every change to your files keeps a backup" overstated: plugin registry rewrites are not backed up (`features/doctor.md`). | Now "Files it moves out of the way are kept in a backup." |
| 45 | MEDIUM | R-36 | `content/en/docs/quick-start.md` | MCP servers were said to live in `~/.claude`; they are in `~/.claude.json`. The shortcut question appeared in step 4 although `ccam add` asks it in step 3. | Both fixed; step 3 shows the question and the confirmation. |
| 46 | MEDIUM | R-17 | 11 pages with the Team and Enterprise seats note | The same long quote from the help center repeated on 12 pages and broke each page's flow. | Two sentences on 11 pages pointing to the FAQ, which keeps the quote; `features/next.md` keeps its own terms section. |
| 47 | MEDIUM | readability | `features/*.md` How it works | Function and package names were the subjects of sentences, so the behaviour a reader needs was buried. 26 paragraphs ran over 150 words. | ccam is the subject and the names follow in parentheses (each page keeps the same inline code). Long paragraphs became lists, steps, tables or shorter paragraphs; none is over 150 words now. |
| 48 | LOW | spelling, idiom | all pages, CLI | British spellings (honour, organisation, greyed), "carry … on", "change of groups", "another ccam", "beside" for processes, "else" in prose, "nobody sits at". | US spelling in prose; idioms replaced in the pages and in the CLI strings the pages quote. Code blocks and CLI output keep their own words (`"summarise"`, `cancelled`). |
| 49 | LOW | R-36 | `content/en/compare/account-switchers.md` | The ccam row of the comparison table was several times longer than the others. | Shortened to the same facts, with the same inline code. |

Checked and passing: no em or en dashes added (the 4 in the English pages are in quoted CLI output, as before); no serial commas added except where an item already contains "and" or "or"; no buzzwords; heading count, levels and order unchanged (three headings reworded, no links point at them); every description 70 to 160 wide; i18n-check `--strict` passes for all 29 languages.
