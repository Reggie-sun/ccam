# Content Quality and E-E-A-T: ccam.dev

Audit date: 2026-10-19. Framework: Google Search Quality Rater Guidelines (Sept 11, 2025).
Scope: all 14 English URLs, analyzed from the live HTML (render_page.py, raw mode, no SPA) and from `site/content/en` + `site/messages/en.json`. Spot-checked 3 locales (ko, ja, es) on the live site and in source. Parity was checked for all 29 locales.

## Scores

| Metric | Score |
| --- | --- |
| **Content Quality (overall)** | **74 / 100** |
| E-E-A-T composite | 65 / 100 |
| AI citation readiness | 82 / 100 |
| Readability | 95 / 100 (Flesch 67-78, grade 5-8) |
| Translation quality (ko/ja/es sample) | 92 / 100 |
| Templated metadata risk | low (0 / 23 templated, no shared CTA) |

### E-E-A-T breakdown (this skill's weighting, not Google's)

| Factor | Weight | Score | Points | Evidence |
| --- | --- | --- | --- | --- |
| Experience | 20% | 70 | 14.0 | Real console transcripts, "Checked with Claude Code 2.1.302", a measured "about 4 ms" with the named script `scripts/perf-check.sh`, non-obvious findings (macOS Keychain entry `Claude Code-credentials` + path hash; `~/.claude.json` lives outside `~/.claude` without `CLAUDE_CONFIG_DIR`). Weak points: no first-person maintainer voice, no "why I built it" story, no user stories. |
| Expertise | 25% | 68 | 17.0 | Content is technically precise and internally consistent. No author is named anywhere on the site. `TechArticle` JSON-LD has no `author`. The only identity signal is the GitHub handle in URLs. |
| Authoritativeness | 25% | 36 | 9.0 | Repo created 2026-10-18 with 0 stars and 0 forks. No external mentions. Only one outbound link to Anthropic docs (install page). Claims such as "the approach Claude Code's documentation describes" have no citation link. This is expected for launch week. |
| Trustworthiness | 30% | 83 | 25.0 | Unofficial/not-endorsed disclaimer on every page. Terms-of-use caveat. "No telemetry" statement. SHA-256 checksum on install and update. MIT license. "Edit this page on GitHub" on every page. Honest note that the usage endpoint "is not publicly documented and may change". Machine-translation notice. Gaps: no website privacy note (Cloudflare Web Analytics beacon is loaded), no security contact (no SECURITY.md) for a tool that reads OAuth tokens, no named maintainer or contact. |
| **Total** | | | **65** | |

## Word count vs page-type floors (live extracted_text; floors are coverage floors, not targets)

| URL | Words (live) | Type / floor | Verdict |
| --- | --- | --- | --- |
| / | 463 | Homepage / 500 | Slightly under. The "In short" facts block covers the intent. |
| /en/docs/quick-start/ | 505 | Docs | OK |
| /en/docs/install/ | 320 (source 603) | Docs | OK. Live is the stale 0.2.0 build (F2). |
| /en/docs/commands/ | 950 | Docs | OK |
| /en/docs/sharing/ | 988 | Docs | OK. Most thorough page. |
| /en/docs/usage-limits/ | 468 | Docs | OK |
| /en/docs/faq/ | 493 | Docs | OK (13 Q&As) |
| /en/guides/multiple-claude-code-accounts/ | 376 | Guide / 1,500 | **Thin for the site's main acquisition query** (F4) |
| /en/guides/claude-config-dir/ | 313 | Guide / 1,500 | Under. Dense and original, but short. |
| /en/guides/usage-limit-workflow/ | 254 | Guide / 1,500 | Under (F4) |
| /en/compare/aliases/ | 254 | Comparison | Mostly a table. Overlaps the guide (F8). |
| /en/blog/introducing-ccam/ | 318 | Blog / 1,500 | Under. An announcement post with no author (F5). |
| /en/changelog/ | 570 | Utility | Raw commit SHAs (F1) |
| /en/languages/ | 54 | Utility | Thin, x30 locales (F6) |

## Readability (English prose; code, tables, and headings removed)

Every page scores Flesch Reading Ease 69.9-77.8 and FK grade 5.2-6.8. Average sentence length is 11-14 words. The homepage copy scores FRE 67 and grade 7.7. Paragraphs open with the answer. The longest sentence (38 words, sharing.md) is still readable. No action needed.

## AI-generated content quality check (QRG Sept 2025)

**Pass.** Phrasing is specific, with no filler or generic intros. Most pages contain claims you cannot find elsewhere: measured latency, Keychain naming, the `.claude.json` location gotcha, the MCP sharing workaround via `--mcp-config`, a token-refresh safety rationale, and real CLI output. Page structure differs by page type, so there is no repeated template. The repo has an antislop review process (`docs/antislop-review.md`).

## Templated metadata (metadata_template.py, 23 title/description pairs, en + ko/ja/es samples)

`site_risk: low`, `templated_ratio: 0.0`, `shared_cta_phrases: {}`. Twelve pages are flagged `brand-suffix-in-description`. These are **false positives**: "ccam" is the product name and appears naturally in descriptions. No action.

## Translation spot-check (ko, ja, es) and parity (all 29 locales)

- **Freshness:** for all 29 locales x 11 Markdown pages, every `sourceHash` matches the current English file (0 stale, 0 missing).
- **Structure:** in ko, ja, and es, H2/H3 counts and code-fence counts match English exactly on all 11 pages. Code and commands are untranslated, as the glossary requires.
- **Length ratio (prose chars vs English):** ko 0.52-0.63, ja 0.65-0.77, es 1.06-1.16. These ratios are normal for these languages, so nothing is truncated or thin.
- **Quality:** Korean reads idiomatically rather than literally. Japanese uses consistent です/ます polite form and natural term choices. Spanish is fluent neutral Latin American ("computadora"). Internal links are rewritten per locale (`/ko/docs/sharing/#copy-an-existing-setup`). The machine-translation notice and edit link are present.
- **Minor drift:** the ko homepage lede ("한도에 도달하면 바로 다른 계정으로 변경해 보세요") drops the "same conversation" point that the English lede and the demo make. It also uses 변경 where the docs use 전환 for "switch".
- **Not translated:** the body of /{locale}/changelog/ is the English commit list on all 30 locales (F1).

---

## Findings

### F1. Changelog is a raw commit dump, duplicated on 30 locale URLs. Severity: High
- **Evidence:** `/en/changelog/` lists entries such as `10bd4ec7f325a9c4790fdcd1aa56074e0eb74d39 test: compare shortcut links by resolved path; site says ccam 0.2.0`, `fix(site): smaller Korean hero headline and intro`, `feat: Phase 1 skeleton, profiles, launcher and shims`. These are 40-char SHAs, website-internal and test commits, and dev-phase jargon. `scripts/changelog.ts` copies the GitHub release bodies, which goreleaser auto-generated. `/ko/changelog/` (and the 28 other locales) shows the same English body. All 30 are in the sitemap (30 of 420 URLs).
- **Impact:** low-value, near-duplicate pages. Users get no answer to "what changed for me?" The page is also a poor citation source for AI answers.
- **Fix:** curate release notes at the source. In `.goreleaser.yaml`, set `changelog: use: github`, `abbrev: 7`, add `groups` (Features / Fixes), and add `filters.exclude` for `^test`, `^chore`, `^ci`, `^docs\(site\)`, `^fix\(site\)`, `^feat\(site\)`, `^perf\(site\)`. Better still, write a 3-5 bullet user-facing summary per release. Until notes are localized, either noindex the 29 non-English changelog pages and drop them from the locale sitemaps, or render a translated one-line summary per release.

### F2. Live site is a release behind; source changelog data lacks v0.3.0. Severity: Medium
- **Evidence:** GitHub shows `v0.3.0 Latest 2026-10-19T08:33:53Z`. The live docs say "Checked with ccam 0.2.0". The live `/en/docs/install/` has no "Update" section: 320 words live vs 603 in source, and the description lacks "keep it up to date with ccam update". The live `/llms.txt` has the old install description. In source, `lib/jsonld.ts` has `VERSION = "0.3.0"`, but `data/changelog.json` has only v0.2.0 and v0.1.0. `changelog.ts` "keeps the existing file when the API is unreachable", so an offline build would ship "0.3.0" next to a changelog that ends at 0.2.0.
- **Fix:** redeploy. Make the build fail, or at least warn loudly, when `VERSION` has no matching tag in `data/changelog.json`.

### F3. No identifiable author or maintainer anywhere. Severity: Medium
- **Evidence:** no byline, About page, or maintainer name on any page. `articleLd()` emits `TechArticle` without `author` or `datePublished`. `Organization.name` is "ccam", but the project is run by one person (GitHub `jhyunwoo`).
- **Fix:** add a short maintainer line (name, GitHub link, optional site) to the blog post and footer, or create a small /about page. Add `author: {"@type":"Person", name, url: "https://github.com/jhyunwoo"}` to the article JSON-LD and the blog post. Use the real maintainer identity; do not invent credentials.

### F4. The primary acquisition guides are thin for their queries. Severity: Medium
- **Evidence:** `/en/guides/multiple-claude-code-accounts/` (376 words) targets the site's core query ("use two/multiple Claude Code accounts"). `/en/guides/usage-limit-workflow/` has 254 words. `/en/guides/claude-config-dir/` has 313 words. All are well under the 1,500-word guide floor (a coverage floor, not a target).
- **Fix (coverage, not padding):** in the multiple-accounts guide, add a short "how to tell which account you are in" section (`ccam list`, status line), Windows/WSL notes, macOS Keychain implications, a troubleshooting section (wrong account signed in, aliases not seen in scripts), and a migration path from aliases (`ccam import`). In the workflow guide, add what happens to an in-progress session, the weekly vs 5-hour limit distinction, and a worked status-line example. Verify every Claude Code behavior you add against the current version, as the config-dir guide already does.

### F5. Blog post lacks experience signals and publication metadata. Severity: Medium
- **Evidence:** "Introducing ccam" (318 words) is written impersonally and has no author. It shows only "Last updated" and uses the `TechArticle` schema with no `datePublished`.
- **Fix:** add a first-person paragraph on the concrete problem the maintainer hit (for example, running out of the 5-hour limit mid-task and re-creating MCP config per account). Add a byline and publish date. Use `BlogPosting` with `author`, `datePublished`, and `dateModified`.

### F6. /languages/ is a 54-word page x 30 locales. Severity: Low
- **Evidence:** extracted text is a title plus three sentences (the language links are navigation). All 30 locale copies are in the sitemap.
- **Fix:** it's a utility page, so noindex it, or keep only `/en/languages/` indexable. Alternatively, add value such as native names, translation status (machine vs reviewed), and how to contribute corrections.

### F7. Unlinked references to Anthropic documentation. Severity: Low-Medium (authority)
- **Evidence:** the FAQ says ccam "uses the `CLAUDE_CONFIG_DIR` setting that Claude Code documents". The guide says "This is the approach Claude Code's documentation describes". Neither links to a source.
- **Fix:** link each claim to the specific Anthropic docs page or anchor. Citing primary sources is the cheapest authority signal available to a new project.

### F8. Overlap between the compare page, the guide, and the blog. Severity: Low
- **Evidence:** the list of alias limitations (empty folders, aliases not in scripts, MCP in `.claude.json`, no cross-account conversations, no usage view) appears in the guide, the compare table, the blog "The problem" section, and the FAQ. `/en/compare/aliases/` is 254 words, mostly that table.
- **Fix:** keep the table as the canonical comparison. Enrich it with a script/Makefile example showing aliases failing, setup time for N accounts, and an "alias to ccam" migration snippet. Shorten the duplicate list in the guide to one sentence plus a link.

### F9. Generic doc titles. Severity: Low
- **Evidence:** "Quick start | ccam", "Install ccam | ccam" (brand appears twice), "Command reference | ccam", "Frequently asked questions | ccam", "Changelog | ccam" (description is 42 chars), "Languages | ccam" (38 chars). None of these says what ccam is for.
- **Fix:** examples: "Install ccam on Linux and macOS (Claude Code account manager)", "ccam quick start: add and switch Claude Code accounts", "ccam FAQ: multiple Claude Code accounts". Extend the changelog description to about 120 characters.

### F10. No website privacy or security contact. Severity: Low-Medium (trust)
- **Evidence:** the CSP and HTML load the `static.cloudflareinsights.com/beacon.min.js` Cloudflare Web Analytics beacon. The site's only privacy statement ("ccam has no telemetry") is about the CLI. The repo root has no `SECURITY.md`, yet the tool reads OAuth access tokens.
- **Fix:** add one footer line or FAQ entry stating that the site uses Cloudflare Web Analytics and what it collects. Add `SECURITY.md` with a reporting channel and link it from the FAQ "Does ccam see or send my login?" answer.

### F11. Homepage has no date signal; parsers misread the beacon version. Severity: Low
- **Evidence:** render_page reported `publication_date: 2024-11-01` for `/` and `/en/languages/`. That date is parsed from the Cloudflare beacon's `"version":"2024.11.0"`, which predates the project.
- **Fix:** add `dateModified` (and `datePublished` = first release date) to the homepage `SoftwareApplication` / `WebPage` JSON-LD, or a visible `<time>` for "Latest release vX.Y.Z (date)".

### Positive / no action
- Answer-first openings on every page, plus the homepage "In short" facts block. These are highly quotable for AI answers.
- `llms.txt` plus `.md` mirrors of every doc (HTTP 200), with tables and real CLI transcripts.
- Freshness: every page shows "Last updated" and "Checked with ccam X and Claude Code Y". This is an excellent pattern once it's redeployed (F2).
- FAQPage JSON-LD is present on /docs/faq/. FAQ rich results were retired on May 7, 2026, so it gives no SERP benefit. It's harmless; keep or remove it as you like. Do not invest in it for search.
- Keyword use is natural, with no stuffing. Variants such as "several/two/multiple Claude Code accounts" and "usage limit" appear where they fit.

## AI citation readiness: 82 / 100
Strengths: definitional first sentences, the facts block, llms.txt with Markdown mirrors, comparison and option tables, specific numbers (5-hour and weekly limits, 3 s timeout, 60 s cache, about 4 ms), question-form H2s in the FAQ, stable anchors across locales. Deductions: no author or `datePublished` (F3/F5), the stale live llms.txt and version (F2), the changelog as a source (F1), guide H2s that are not phrased as questions (for example, "By hand, with an alias" could be "How do I run two Claude Code accounts with an alias?"), and no primary-source citations (F7).

---

## Structured findings (audit-data.json, category "Content Quality")

```json
{
  "category": "Content Quality",
  "score": 74,
  "eeat": {"experience": 70, "expertise": 68, "authoritativeness": 36, "trustworthiness": 83, "composite": 65},
  "ai_citation_readiness": 82,
  "readability": {"flesch_reading_ease_range": [67.2, 77.8], "fk_grade_range": [5.2, 7.7]},
  "templated_metadata": {"site_risk": "low", "templated_ratio": 0.0, "shared_cta_phrases": {}, "note": "brand-suffix flags are false positives (product name)"},
  "translations": {"locales_checked_for_staleness": 29, "stale_pages": 0, "spot_checked": ["ko", "ja", "es"], "quality": "fluent; structural parity exact"},
  "findings": [
    {"id": "F1", "severity": "high", "title": "Changelog is raw commit SHAs incl. site/test commits; English body duplicated on 30 locale URLs", "url": "/en/changelog/", "fix": "Curate goreleaser changelog (filters, groups, abbrev) or hand-written notes; noindex non-en changelog until localized"},
    {"id": "F2", "severity": "medium", "title": "Live site on 0.2.0 while v0.3.0 is Latest; update docs not live; data/changelog.json lacks v0.3.0", "url": "/en/docs/install/", "fix": "Redeploy; fail build when VERSION has no tag in changelog.json"},
    {"id": "F3", "severity": "medium", "title": "No author/maintainer identity; TechArticle has no author", "url": "site-wide", "fix": "Maintainer byline/about; Person author in JSON-LD"},
    {"id": "F4", "severity": "medium", "title": "Main acquisition guides thin (376/313/254 words)", "url": "/en/guides/", "fix": "Add which-account, WSL, Keychain, troubleshooting, migration sections"},
    {"id": "F5", "severity": "medium", "title": "Blog post impersonal, no byline or datePublished, typed TechArticle", "url": "/en/blog/introducing-ccam/", "fix": "First-person origin story; BlogPosting with author/datePublished"},
    {"id": "F6", "severity": "low", "title": "/languages/ 54 words x30 locales indexed", "url": "/en/languages/", "fix": "noindex non-en or enrich"},
    {"id": "F7", "severity": "low", "title": "Claims about Claude Code docs not linked", "url": "/en/docs/faq/, /en/guides/multiple-claude-code-accounts/", "fix": "Link to Anthropic docs anchors"},
    {"id": "F8", "severity": "low", "title": "Alias-limitations content repeated across 4 pages; compare page thin", "url": "/en/compare/aliases/", "fix": "Make compare page canonical and richer; trim duplicates"},
    {"id": "F9", "severity": "low", "title": "Generic doc titles; 'Install ccam | ccam' double brand; short changelog/languages descriptions", "url": "docs", "fix": "Add purpose keywords to titles"},
    {"id": "F10", "severity": "low", "title": "No site privacy note for Cloudflare Web Analytics; no SECURITY.md", "url": "site-wide", "fix": "Footer privacy line; SECURITY.md linked from FAQ"},
    {"id": "F11", "severity": "low", "title": "Homepage lacks date signal; parsers read beacon version 2024.11.0 as publish date", "url": "/", "fix": "dateModified/datePublished in homepage JSON-LD or visible <time>"}
  ]
}
```
