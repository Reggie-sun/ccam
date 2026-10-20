# SXO Findings: ccam.dev

Date: 2026-10-19. Analyst: SXO subagent. This file reports an **SXO Gap Score**, which is separate from the SEO Health Score.

Method: I fetched six pages with `render_page.py --mode auto`. All six were static HTML (`is_spa=false`, raw mode), so the rendered DOM is the same as the raw HTML. I parsed each page with `parse_html.py`. For the SERP pass I ran 10 web searches (the 5 target queries plus 5 supporting queries) and opened 6 competitor pages to check their type and depth. Taxonomy, user-story and persona rules come from the `seo-sxo/references/` files.

---

## 0. Headline finding (read first)

**CRITICAL: ccam.dev does not appear in any observed SERP, and a `site:ccam.dev` search returns no ccam.dev URLs.**
- Evidence: none of the 5 target queries returned a ccam.dev result in the top 9-10. `site:ccam.dev` returned CSIRO's "Conformal Cubic Atmospheric Model" and the EU "CCAM Association" instead. The brand query "ccam claude code account manager" returned CAAM, CCS, ccswitch, ccdaddy and ccm-cli, and the search engine asked "did you mean CAAM?".
- What it means: the page-type fixes below only pay off once the site is indexed. Part of the problem is outside SXO (indexing). Part of it is SXO: the brand name "ccam" collides with two established entities and with a competitor, "caam". The expansion "Claude Code Account Manager" appears only on header hover (per changelog commit 5b483c9). It is not in the title, the H1 or the schema `name`.
- Fix: (1) Check indexing in Search Console with URL Inspection and submit the sitemap. Hand this to `/seo technical`. (2) Put "Claude Code Account Manager" in visible text: the homepage title (`ccam: Claude Code Account Manager for multiple accounts`), the hero subline, and `SoftwareApplication.alternateName`. (3) The GitHub repo will probably rank before the site does, because tool queries are dominated by GitHub and PyPI. Make sure the repo description and topics contain "claude code", "multiple accounts", "account switcher" and "usage limit".

---

## 1. Query-to-page map and page-type mismatch

SERP consensus uses the 8-type taxonomy. "Tool" covers software project pages (GitHub repo, PyPI, docs.rs), because those results serve the "I want something to DO this" intent.

| Query | Page that should rank | SERP dominant type (confidence) | ccam page type | Mismatch |
|---|---|---|---|---|
| multiple claude code accounts | /en/guides/multiple-claude-code-accounts/ | Blog Post how-to (~78%: fireup.pro, agentsroom blog x5 locales, skillshare blog, codeminer42, aident, dev.to); Landing ~22% (agentsroom /features x2) | Blog Post / Hybrid (TechArticle docs guide) | **MEDIUM** (type fits; depth and format do not) |
| claude code multiple accounts same computer | same guide | Blog Post how-to (~78%: codeminer42 x2, dev.to, aident, agentsroom, fireup); Tool ~11% (claude-switch README); Landing ~11% | same | **MEDIUM** |
| claude code switch account | / (homepage), backed by /en/docs/quick-start/ | Tool (~100%: claude-swap GitHub + 6 PyPI versions, claude-account-switcher GitHub, clauth docs.rs) | Landing Page (SoftwareApplication) | **HIGH** (product page is the right family, but it never says "switch" or "switcher" and doesn't answer the native `/login` intent) |
| claude code usage limit switch account | /en/guides/usage-limit-workflow/ (or /en/docs/usage-limits/) | Tool (~100%: claude-swap GitHub/PyPI x6, ccswap PyPI, clauth docs.rs) | Blog Post (308-word TechArticle) | **HIGH** |
| CLAUDE_CONFIG_DIR | /en/guides/claude-config-dir/ | Reference / Blog Post (~60%: Medium explainer, pasqualepillitteri guide, official code.claude.com debug page, Nix option docs); GitHub issues + mirrors ~40% | Blog Post / reference explainer | **ALIGNED** |
| /en/compare/aliases/ (supporting) | for "claude code account switcher" / "alternatives" intents | Comparison intent is implied by the 8+ competing tools in the SERPs | Comparison page, but against aliases only | **HIGH** (compares against the wrong alternative) |
| /en/docs/quick-start/ | brand and navigational queries only | n/a | Docs | Not a SERP target. It should support the homepage and should not compete with it. |

### F1. HIGH: "usage limit switch account" rewards tool pages; ccam sends this intent to a 308-word guide
- Evidence: 9 of 9 results are software projects (claude-swap 3.6k stars, ccswap, clauth). Their copy repeats the same vocabulary: "auto-switch", "rate limit", "5-hour / 7-day window", "threshold (default 90%)", "fallback account". A grep of `site/content/en/**` finds 0 occurrences of "rate limit", 0 of "switcher", 0 of "auto-switch". The ccam guide has no install command (it starts at `ccam add`) and no image or recording.
- Fix: Retitle the guide around the SERP vocabulary: `Claude Code usage limit reached? Switch to another account and keep the conversation`. Add a "Install (one line)" block above "Before it happens". Add a 3-row table: Wait for reset / Manual login switch / `ccam next --auto --continue`. State plainly whether ccam switches automatically at a threshold (it doesn't; `next --auto` picks the account with the most room left when you run it), because searchers are comparing against claude-swap's `cswap auto`. Link to /en/docs/usage-limits/. Grow the page to about 800 words.

### F2. HIGH: the homepage does not match the "switch account" query wording or its two-part intent
- Evidence: Title `ccam: use several Claude Code accounts on one computer`, H1 `Several Claude Code accounts. One terminal.` Neither contains "switch" or "switcher". The SERP shows two intents. (a) A tool to switch accounts: claude-swap and others. (b) How to switch with the built-in `/login`, `/logout` and `/status`: how2shout, usagebar, anthropics/claude-code #42775, #23906. No English page mentions `/login` or `/status`.
- Fix: Change the hero subline to "A Claude Code account switcher: add each account once, switch with one command". Add one sentence plus a link: "Only need to switch occasionally? `/login` inside Claude Code works; ccam is for running accounts side by side." This captures intent (b) instead of losing it.

### F3. HIGH: the comparison page compares against shell aliases, but searchers compare tools
- Evidence: the SERPs show claude-swap, CCS (Claude Code Switch), CAAM, ccswitch, clauth, claude-switch, maxx-headroom and AgentsRoom. /en/compare/ has one page, `aliases`, at 239 words. Its BreadcrumbList repeats the same URL at positions 2 and 3.
- Fix: Add `/en/compare/claude-swap/` (the most visible competitor). Optionally add a hub, `/en/compare/`, titled "Claude Code account switchers compared". Columns: how logins are stored (copied tokens vs separate CLAUDE_CONFIG_DIR, which is ccam's trust advantage), parallel sessions, shared MCP and settings, continuing a conversation, usage view, auto-switch, OS (claude-swap supports Windows; ccam doesn't), runtime (Python 3.12 vs a static binary), and startup overhead (~4 ms). Keep the comparison factual and dated. Fix the duplicate breadcrumb (`/seo schema`).

### F4. MEDIUM: the multiple-accounts guide is the right page type but shallow next to the SERP
- Evidence: ranking how-tos run about 1,100 words (fireup.pro, with author, date, a 5-question FAQ and a "Windows Users" section), about 1,600 words (aident.ai, with 6 numbered steps and a "Common Failure Modes" table), and about 400 words (codeminer42, with author and date). The ccam guide has 440 words, 3 content H2s, 0 images, no author, and no FAQ block. The title and H1 say "two", while the query says "multiple". The guide doesn't cover: whether separate accounts give more quota, hooks hard-coded to ~/.claude cross-contaminating accounts (aident), CLAUDE.md loading from both folders (issue #30230), or Windows/WSL.
- Fix: Title `How to use multiple Claude Code accounts on one computer (2026)`. Add H2s "Does this give me more usage?", "Is it allowed?" (link /en/docs/faq/), "Common problems" (table: symptom / cause / fix), and "Windows" (WSL note; ccam is Linux/macOS). Embed 3-4 Q&As from /en/docs/faq/ at the end. Target about 900-1,100 words.

### F5. LOW (ALIGNED): the CLAUDE_CONFIG_DIR explainer has unique, high-value content
- Evidence: the folder-contents table, the `~/.claude.json` location gotcha, and the macOS Keychain naming ("Claude Code-credentials" plus a path hash) don't appear in any SERP result I sampled. The first sentence works as a definition snippet. Gaps: no link to the official code.claude.com env-vars or debug-your-config docs, no "known issues" (#30230 double CLAUDE.md; the variable does not move project-level `.claude/`; project and local settings cannot set it), and no Windows path. 349 words.
- Fix: add "Official documentation" links, a "Known issues" H2 with issue links, and a "Set it permanently" section (shell rc, `settings.json` env with an absolute path, as the official env-vars page describes).

### F6. MEDIUM: trust signals are thin everywhere
- Evidence: no author or maintainer entity in the TechArticle schema or the visible bylines, no GitHub star count or release badge, and no screenshots. The SERP raises these trust concerns: "avoid tools that store raw credentials", Anthropic terms/policy notes (clauder-navi), and issue #12786 ("rate limit carried over when switching Max accounts on the same device"). ccam answers some of this in /en/docs/faq/ ("Does ccam see or send my login?", "Am I allowed…?", "Does ccam get around usage limits?"), but the guides don't surface those answers. The usage guide puts its terms caveat as the last bullet.
- Fix: add a 3-line "Safe by design" callout to the homepage, the multi-accounts guide and the usage guide: logins stay in each account's own folder or Keychain entry and are never copied, linked or sent; the tool is open source with a link to the source; it is unofficial and not affiliated. Add `author` (Person: jhyunwoo, `url` to GitHub) and `datePublished` to the TechArticle schema (`/seo schema`, `/seo content` for E-E-A-T).

### F7. LOW: freshness signal is inconsistent
- Evidence: the live homepage schema has `softwareVersion: "0.2.0"` and the guides say "Checked with ccam 0.2.0". The repo's `site/lib/jsonld.ts` says `VERSION = "0.3.0"`, and install.md shows a 0.3.0 update. The deploy is behind the repo. `htmldate` reads the homepage date as 2024-11-01, which is a wrong inference; the homepage has no visible date.
- Fix: redeploy. Consider adding a visible "Latest: v0.3.0 (date)" line to the homepage hero.

### F8. LOW: site nav headings come before the content outline
- Evidence: every docs and guide page lists H2s `Documentation`, `Comparisons`, `Blog` (the sidebar) before the content H2s. Passage parsers and AI-answer extractors see these first.
- Fix: render sidebar group labels as non-heading elements (for example a `<p>` or `<span>` with `role="heading"` kept out of the main outline, or `aria-labelledby` on the nav).

---

## 2. User stories (from SERP signals)

1. **Decision stage.** As a developer paying for Pro and Max, I want to continue the same conversation on my other account the moment one hits its 5-hour limit, because waiting for the reset breaks my flow. I'm blocked by **no native in-session switch** and **manual context transfer**.
   *Signals: anthropics/claude-code #64376 ("no way to switch to a secondary account within the same session… manual context transfer"); claude-swap `cswap auto` at 90%; clauth "auto-switches to a fallback account at the 5-hour limit".*
2. **Consideration stage.** As a developer with a company (Team) account and a personal account, I want both isolated on one laptop without logging out and back in, because mixing work and personal history is a compliance worry. I'm blocked by **fear of cross-contamination** (hooks, CLAUDE.md loaded twice) and buggy `/login`.
   *Signals: claudeissues #22872 "Switching between personal and work Claude Code accounts is abysmal"; #23906 `/login` keeps the old account; aident.ai "…Accounts Safely"; issue #30230.*
3. **Awareness stage.** As a newcomer who just bought a second subscription, I want to know whether multiple accounts are possible at all and the simplest way to switch, because sources contradict each other. I'm blocked by an **information gap** (logout vs `/login`; "documented" vs "undocumented").
   *Signals: usagebar "Can I Use Multiple Accounts with Claude Code?"; how2shout logout guide; the SERP summary notes conflicting claims about whether CLAUDE_CONFIG_DIR is documented.*
4. **Consideration stage.** As someone comparing account switchers, I want to choose between claude-swap, CCS, CAAM, ccswitch and clauth, because they all claim the same thing. I'm blocked by **comparison fatigue** and **credential-handling distrust**.
   *Signals: 9/9 "switch account" results are switcher tools; the brand query surfaces 6 competitors; "avoid tools that store raw credentials".*
5. **Consideration stage.** As a cautious subscriber, I want to know whether using several accounts is allowed and whether it gets me flagged, because I can't afford to lose my account. I'm blocked by a **trust gap**.
   *Signals: clauder-navi note on Anthropic's terms; issue #12786 on device-level rate limit tracking.*

---

## 3. SXO Gap Score per target page (100 points)

Dimensions: Page Type /15, Content Depth /15, UX Signals /15, Schema /15, Media /15, Authority /15, Freshness /10.

| Page | Type | Depth | UX | Schema | Media | Auth | Fresh | **Total** |
|---|---|---|---|---|---|---|---|---|
| / (for "switch account") | 11 | 9 | 13 | 12 | 9 | 4 | 7 | **65** |
| /en/guides/claude-config-dir/ | 13 | 9 | 12 | 10 | 4 | 5 | 9 | **62** |
| /en/guides/multiple-claude-code-accounts/ | 11 | 7 | 11 | 10 | 3 | 5 | 9 | **56** |
| /en/guides/usage-limit-workflow/ | 6 | 6 | 12 | 9 | 5 | 4 | 9 | **51** |
| /en/compare/aliases/ | 9 | 5 | 12 | 8 | 2 | 4 | 9 | **49** |

What the scores rest on:
- **Homepage.** Install one-liner above the fold, an animated terminal demo (SVG usage bars), and SoftwareApplication + Offer(0) + Organization + WebSite schema. Weak on authority: no stars, users or testimonials. Stale version (0.2.0). 602 words.
- **config-dir.** Strong on type and uniqueness thanks to the table. Weak on media, and it lacks official-doc links.
- **multi-accounts.** 440 words against SERP pages of roughly 1,100-1,600. 0 images. No FAQ. "two" instead of "multiple".
- **usage-limit.** Wrong type for a SERP made of tools. 308 words. No install step. The text usage-bar mock is a good UX element.
- **aliases.** Compares against the wrong alternative. 239 words. Duplicate breadcrumb URL.

Shared across all pages: `Last updated October 19, 2026` and "Checked with Claude Code 2.1.302" are visible, which is good. No author. 0 `<img>` elements. Clean typography, an "On this page" TOC, and a prev/next pager.

---

## 4. Persona scores (each persona scored on the page it most likely lands on)

Weighting: tool intent dominates 2 of 5 queries and how-to intent dominates the other 3.

| Persona | Entry page | Relevance | Clarity | Trust | Action | Total | Rating |
|---|---|---|---|---|---|---|---|
| Tool evaluator comparing switchers (SERP-derived) | /en/compare/aliases/ | 8 | 20 | 10 | 14 | **52** | Needs Work |
| Newcomer (second subscription, unsure how) | / and multi-accounts guide | 15 | 14 | 12 | 18 | **59** | Needs Work |
| Power user, Pro + Max, hits limits | /en/guides/usage-limit-workflow/ | 22 | 20 | 12 | 14 | **68** | Good |
| Developer with team (work) + personal accounts | /en/guides/multiple-claude-code-accounts/ | 18 | 19 | 13 | 18 | **68** | Good |

### Weakest: tool evaluator (52)
- Top issue: the evaluator arrives from a SERP full of claude-swap, CCS and CAAM, and gets a comparison against shell aliases. Nothing explains why to pick ccam over claude-swap.
- Fix: see F3. Create `/en/compare/claude-swap/` with a "How each tool handles your login" row first; that's ccam's main difference. Put a "Install ccam" CTA under the table and a "See all differences vs aliases" link.

### Newcomer (59)
- Top issues: the H1 "Several Claude Code accounts. One terminal." is clever but doesn't confirm the question ("Can I use two accounts?"). Jargon (CLAUDE_CONFIG_DIR, MCP, symlinks) shows up in the first scroll. The native `/login` route is never mentioned.
- Fix: add a "Can I use more than one Claude Code account?" answer box near the top of the multi-accounts guide: Yes, with 3 options (`/login` to swap, an alias per account, or ccam), each with one line on when to use it. On the homepage, put "Claude Code account switcher" in the subline (F2).

### Power user (68)
- Top issues: no install line on the usage guide. No statement on auto-switching versus claude-swap's threshold auto-switch. The terms caveat is buried at the end.
- Fix: add an install block plus a "Before it happens" quick checklist. Add one sentence: "ccam doesn't switch in the background; `ccam next --auto --continue` picks the account with most room when you run it." Move "Is this allowed?" (link /en/docs/faq/) to directly under the intro.

### Developer with team + personal (68)
- Top issue: the guide never addresses work/personal isolation explicitly: history is not shared by default, logins are never shared, and hooks that hard-code ~/.claude can leak. No Team/Enterprise SSO note.
- Fix: add an H2 "Keeping work and personal separate" listing what ccam shares by default (MCP, settings, CLAUDE.md, extensions), what it never shares (logins), what is opt-in (history), and how to turn sharing off per account (`ccam share work …`).

### Systemic issues
- **Trust is the lowest dimension for every persona (10-13/25).** There is no maintainer identity, no adoption proof and no "Safe by design" statement on entry pages, even though the FAQ already has the answers. See F6.
- **The SERP's vocabulary is missing from the site.** "switcher", "rate limit", "auto-switch", `/login` and Windows/WSL are absent from the English guides.

### Priority actions
1. Fix indexing and brand disambiguation (Section 0). This is a prerequisite for everything else.
2. Build `/en/compare/claude-swap/` and retitle the usage-limit guide in tool/switcher vocabulary (F1, F3). These target the weakest persona and the two HIGH mismatches.
3. Add a newcomer answer box and the `/login` mention to the homepage and the multi-accounts guide (F2, F4).
4. Add a "Safe by design" callout and author schema across the guides (F6).
5. Deepen the multi-accounts and CLAUDE_CONFIG_DIR guides with "Common problems", "Known issues" and "Windows" sections (F4, F5).

---

## 5. Structured findings (for audit-data.json, category "Search Experience")

```json
[
  {"id":"sxo-00","severity":"critical","title":"ccam.dev absent from all target SERPs; site: query returns no ccam.dev URLs; brand collides with CSIRO CCAM / EU CCAM / CAAM","pages":["/"],"fix":"Verify indexing in GSC, submit sitemap; add 'Claude Code Account Manager' to title/hero/schema alternateName; optimize GitHub repo description/topics"},
  {"id":"sxo-01","severity":"high","title":"Usage-limit query rewards tool pages; ccam serves a 308-word guide lacking switcher vocabulary and install step","pages":["/en/guides/usage-limit-workflow/"],"fix":"Retitle with 'usage limit reached / switch account', add install block, options table, auto-switch clarification, ~800 words"},
  {"id":"sxo-02","severity":"high","title":"Homepage lacks 'switch/switcher' wording and ignores native /login intent","pages":["/"],"fix":"Hero subline 'Claude Code account switcher'; one-line /login mention with link"},
  {"id":"sxo-03","severity":"high","title":"Only comparison is vs shell aliases; SERP shows 8+ competing switcher tools; duplicate breadcrumb URL","pages":["/en/compare/aliases/"],"fix":"Add /en/compare/claude-swap/ (and optional hub); fix BreadcrumbList position 2"},
  {"id":"sxo-04","severity":"medium","title":"Multi-accounts guide is right type but 440 words vs 1,100-1,600 in SERP; 'two' not 'multiple'; no FAQ/problems/Windows","pages":["/en/guides/multiple-claude-code-accounts/"],"fix":"Retitle 'multiple', add quota/allowed/common-problems/Windows sections, embed FAQ, ~1,000 words"},
  {"id":"sxo-05","severity":"low","title":"CLAUDE_CONFIG_DIR explainer aligned and unique; missing official links and known issues","pages":["/en/guides/claude-config-dir/"],"fix":"Add official doc links, Known issues (#30230), persistence via settings.json env"},
  {"id":"sxo-06","severity":"medium","title":"Thin trust signals: no author entity, no adoption proof, FAQ trust answers not surfaced on entry pages","pages":["all guides","/"],"fix":"'Safe by design' callout; TechArticle author + datePublished"},
  {"id":"sxo-07","severity":"low","title":"Live site shows ccam 0.2.0 while repo/site source says 0.3.0","pages":["/","all guides"],"fix":"Redeploy; visible latest-version line"},
  {"id":"sxo-08","severity":"low","title":"Sidebar group labels are H2s preceding content headings","pages":["all docs/guides"],"fix":"Render nav group labels as non-heading elements"}
]
```

---

## 6. Limitations

- I used the WebSearch tool as a stand-in for Google. It returns result lists and summaries but not Google's SERP features, so I couldn't observe People Also Ask, featured snippets, AI Overviews, ads or related searches. User stories therefore rest on result titles, result content and issue threads, not on PAA. Rankings vary by location and personalization; I saw no Google data from Search Console.
- "Not indexed" is inferred from search-tool results, not confirmed in Search Console.
- Some competitor pages couldn't be fetched (agentsroom: HTTP 413/500; dev.to: 404; Medium: 403). Their classification comes from title and URL only.
- No search-volume data, so persona weighting is qualitative.
- I didn't check Core Web Vitals, real-user engagement, mobile layout or screenshots. All pages are static HTML, so rendered and raw content match.
- Only the English pages were assessed. The 29 other locales share the same templates, so structural findings probably apply to them, but I didn't verify that.
- Wireframes were not requested and are not included.

Cross-skill: `/seo technical` (indexing), `/seo schema` (author, breadcrumb, alternateName), `/seo content` (E-E-A-T, depth), `/seo page` (thin compare/usage pages). Generate a PDF report? Use `/seo google report`.

Sources (SERP/competitors): fireup.pro/news/multiple-claude-code-accounts, agentsroom.dev/blog/how-to-use-multiple-claude-code-accounts, blog.codeminer42.com/ai-tip-running-two-claude-code-accounts/, aident.ai/blog/use-multiple-claude-code-accounts-safely, usagebar.com/blog/can-i-use-multiple-accounts-with-claude-code, github.com/realiti4/claude-swap, pypi.org/project/claude-swap/, pypi.org/project/ccswap/, docs.rs/crate/clauth, code.claude.com/docs/en/env-vars, code.claude.com/docs/en/debug-your-config, github.com/anthropics/claude-code/issues/30230, /28808, /42775, /64376, /12786, claudeissues.com/issue/22872, medium.com/@ddasdocs/where-claude-code-keeps-its-configuration-1c108dde1c20.
