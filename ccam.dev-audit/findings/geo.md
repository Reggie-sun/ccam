# GEO / AI Search Readiness: ccam.dev

Audited 2026-10-19 against the live site (https://ccam.dev) plus read-only review of `site/`.
Pages checked: `/`, `/en/docs/quick-start/`, `/en/docs/faq/`, `/en/guides/multiple-claude-code-accounts/`,
`/en/guides/claude-config-dir/`, `/en/guides/usage-limit-workflow/`, `/en/compare/aliases/`.

## AI Search Readiness score: 66 / 100

| Dimension | Weight | Score | Weighted | Notes |
|---|---|---|---|---|
| Citability | 25% | 78 | 19.5 | Direct first-sentence answers, self-contained FAQ answers, tables. Few outbound citations, little data. |
| Structural readability | 20% | 85 | 17.0 | Clean H1/H2, FAQ uses question H2s, tables, code blocks, Markdown mirrors. Guide H2s are not question-shaped. |
| Multi-modal content | 15% | 45 | 6.8 | Code blocks, tables, OG images, HTML terminal demo. No screenshots, GIF/asciinema or video. |
| Authority & brand signals | 20% | 20 | 4.0 | Repo created 2026-10-18, 0 stars/forks, no third-party mentions, no author entity, ambiguous name. |
| Technical accessibility | 20% | 95 | 19.0 | Static SSR HTML, all AI user agents get 200, sitemap, canonical, hreflang, llms.txt + .md mirrors. |

The site itself is in good shape. The score is held back mainly by **off-site entity weakness**.
The site went live and the repo was created within the last 48 hours, so this is expected. Content changes alone will not fix it.

### Platform-specific estimates

| Platform | Score | Reason |
|---|---|---|
| Google AI Overviews | 55 | Follows Googlebot indexing and ranking. On-page content answers the target queries well, but a new domain with no links has little ranking authority. |
| ChatGPT Search | 45 | Uses OAI-SearchBot (allowed) plus the Bing index. Brand-mention signals are close to zero. Bing indexing has not been verified. |
| Perplexity | 60 | Prefers fresh, direct, extractable answers and often cites GitHub and Reddit. The on-page content suits it well. |
| Bing Copilot | 50 | Depends on the Bing index. Submit to Bing Webmaster Tools / IndexNow. |

## AI crawler access

robots.txt: `User-agent: *` / `Allow: /` plus `Sitemap: https://ccam.dev/sitemap-index.xml`. No per-bot groups, so every crawler is allowed.
User-agent spoof test against `/en/docs/faq/`: every UA below returned HTTP 200, so the Cloudflare "Block AI bots" setting is not active for these UAs.
Caveat: Cloudflare verifies real bots by IP, so also check that AI Crawl Control / managed robots.txt is off in the dashboard.

| Crawler | Governs | Status |
|---|---|---|
| OAI-SearchBot | ChatGPT Search citations | Allowed (200) |
| ChatGPT-User | User-triggered fetches in ChatGPT | Allowed (200) |
| Claude-SearchBot | Claude search citations | Allowed (200) |
| PerplexityBot / Perplexity-User | Perplexity index / user fetches | Allowed (200) |
| Googlebot | Google Search and AI Overviews inclusion | Allowed |
| bingbot | Bing, Copilot, and the index behind ChatGPT search | Allowed (200) |
| GPTBot | OpenAI model training only | Allowed (200) |
| ClaudeBot | Anthropic model training only | Allowed (200) |
| Google-Extended | Gemini/Vertex training and grounding; not Search/AIO inclusion | Allowed (200) |
| Applebot-Extended | Apple Intelligence training only | Allowed |
| CCBot, cohere-ai | Training corpora | Allowed (200) |

Allowing the training crawlers is a reasonable choice for an open-source tool, because it helps future models know that ccam exists. Nothing needs to change.

## llms.txt (optional; Google Search ignores it)

Status: **present and well formed.** `https://ccam.dev/llms.txt` returns 200 `text/plain`. It has an H1, a blockquote summary, a disclaimer, and sections for Documentation, Guides, Comparisons and Blog.
Every linked `.md` mirror returns 200 `text/markdown`. `/llms-full.txt` is also present (36.6 KB). HTML pages link their mirror with `<link rel="alternate" type="text/markdown">`.
No RSL 1.0 license file was found. That is fine to skip.
Minor: the expansion "Claude Code Account Manager" is missing from llms.txt (see F2). There is no `.md` mirror for the homepage (`/index.md` returns 404), but the llms.txt summary covers it. Low value; don't overweight.

## Query coverage

| Query | Best page | Assessment |
|---|---|---|
| how to use multiple Claude Code accounts | /en/guides/multiple-claude-code-accounts/ and the FAQ entry | Strong. The first sentence answers directly ("Give each account its own configuration folder ... `CLAUDE_CONFIG_DIR`"). It covers both the manual alias and ccam, and it is honest about trade-offs. That honesty makes it citable even for answers that are not about ccam. |
| CLAUDE_CONFIG_DIR | /en/guides/claude-config-dir/ | Strong and distinctive. It has a folder contents table ("Checked with Claude Code 2.1.302"), the `~/.claude.json` location gotcha, and the macOS Keychain naming. This is the most original information on the site. Missing: a link to the official Anthropic docs, Windows/WSL notes, how to confirm which account is active, and troubleshooting for "CLAUDE_CONFIG_DIR not working". |
| switch Claude Code account when usage limit reached | /en/guides/usage-limit-workflow/ | Good. The first sentence answers it. But the title and H1 say "reaches its usage limit" without "switch account", and there is no non-ccam path (`/logout`, `/login`) for readers who don't want to install a tool. At 294 words it is thin. |

## Findings

### F1. Almost no off-site entity footprint. Severity: High
- Evidence: GitHub API for `jhyunwoo/ccam` shows created_at 2026-10-18, 0 stars, 0 forks, 0 watchers. HN Algolia search for "ccam.dev" returns 0 hits. DuckDuckGo shows no third-party pages for "ccam.dev" or "jhyunwoo/ccam". No Reddit or YouTube presence was found.
- Impact: AI answer engines rely heavily on third-party mentions (YouTube and Reddit correlate most strongly with citations). Right now an LLM has no outside source to corroborate that ccam exists.
- Fix (effort: medium, ongoing):
  - Post a Show HN.
  - Post a genuinely useful write-up in r/ClaudeAI or r/ClaudeCode. The CLAUDE_CONFIG_DIR internals are the interesting part, not the plug.
  - Record a 1-2 minute YouTube or asciinema demo of the limit, then `ccam next --auto --continue` flow.
  - Open a PR to awesome-claude-code style lists.
  - Answer existing GitHub issues and discussions about multiple accounts where it is relevant, without spamming.

### F2. The name "ccam" is ambiguous and is not expanded where machines read it. Severity: Medium
- Evidence: Wikipedia has a "CCAM" disambiguation page. The acronym also means EU Connected, Cooperative and Automated Mobility and the French medical-procedure classification. On the site, "Claude Code Account Manager" appears only in the header logo text. It is missing from the H1, the meta description, the JSON-LD (no `alternateName`), llms.txt and the GitHub description.
- Fix (effort: low): Use "ccam (Claude Code Account Manager)" once in the homepage "In short" paragraph, in llms.txt line 1, and in the GitHub repo description. Add `"alternateName": "Claude Code Account Manager"` to the `SoftwareApplication` and `Organization` JSON-LD in `site/lib/jsonld.ts`.

### F3. No outbound citations for claims about Claude Code behavior. Severity: Medium
- Evidence: None of the audited pages links outside ccam.dev except to github.com/jhyunwoo/ccam. The multi-account guide says "This is the approach Claude Code's documentation describes" with no link. The FAQ says ccam "uses the `CLAUDE_CONFIG_DIR` setting that Claude Code documents" with no link. The usage-limit pages give no source for the 5-hour/weekly limits.
- Fix (effort: low): Link the official Claude Code docs for `CLAUDE_CONFIG_DIR` / environment variables and for usage limits. Sourced claims are easier to cite and make the "unofficial" positioning more credible.

### F4. No author or maintainer entity, and no datePublished. Severity: Medium
- Evidence: `TechArticle` JSON-LD has headline, description, dateModified and url, but no `author`, `publisher` or `datePublished`. `SoftwareApplication` has no `author`. `Organization` "ccam" has `sameAs` set to only the GitHub repo. No byline appears on any page.
- Fix (effort: low):
  - Add `author: {"@type":"Person","name":"...","url":"https://github.com/jhyunwoo"}` to `TechArticle` and `SoftwareApplication`.
  - Add `publisher` pointing at the Organization.
  - Add `datePublished` from front matter.
  - Extend `sameAs` as profiles appear (YouTube, Reddit, HN, X).

### F5. Gaps in query coverage and thin guides. Severity: Medium
- Evidence: Guides run 281-408 words (aliases 281, usage-limit 294, claude-config-dir 355, multi-accounts 408). Nothing covers:
  - switching accounts in a single config with `/logout` and `/login`
  - Windows and WSL specifics, beyond one FAQ line
  - checking which account is active (`/status`)
  - API-key vs subscription accounts
  - "CLAUDE_CONFIG_DIR not working" troubleshooting
- Fix (effort: medium):
  - Add an H2 to the usage-limit guide, e.g. "Can I switch Claude Code accounts without logging out?", that answers by hand first, then with ccam.
  - Add H2s to the CLAUDE_CONFIG_DIR guide: "How do I check which account Claude Code is using?" and "Why is CLAUDE_CONFIG_DIR not working?".
  - Consider a short guide titled "How to switch Claude Code accounts".

### F6. Guide H2s are labels, not questions. Severity: Low
- Evidence: Example headings: "By hand, with an alias", "With ccam", "Which one to choose", "Two details that catch people out", "Before it happens". The FAQ already uses question H2s.
- Fix (effort: low): Optional. Rephrase a few headings to mirror real queries, e.g. "How do I set up a second account with an alias?" or "Where does Claude Code store the login on macOS?". Google says content does not need to be chunked for AI, so treat this as a light polish.

### F7. Stale `softwareVersion` in JSON-LD. Severity: Low
- Evidence: The live homepage JSON-LD has `"softwareVersion": "0.2.0"`, but GitHub release v0.3.0 was published 2026-10-19T08:33Z. The latest commit `f0b0957` ("site says ccam 0.3.0") is probably not deployed yet.
- Fix (effort: low): Redeploy, and make deployment part of the release step so `VERSION` (`site/lib/jsonld.ts`) never lags behind.

### F8. Little visual or video content. Severity: Low
- Evidence: Content pages have 0 `<img>`; only icon SVGs. The homepage terminal demo is HTML text, which is good for extraction but provides no visual asset. No video exists.
- Fix (effort: medium): Add an asciinema cast or GIF of the limit-to-continue flow, with a descriptive caption and alt text. Publish the same demo on YouTube and link it from the homepage and README. YouTube mentions are the strongest citation-correlated signal.

### F9. Breadcrumb puts guides under "Documentation". Severity: Info
- Evidence: The guide's BreadcrumbList is ccam > Documentation (`/en/docs/install/`) > guide.
- Fix: Use a "Guides" crumb, or drop the middle level.

### Strengths (keep)
- Every page opens with a one-to-two-sentence direct answer. The homepage "In short" block is a ready-made 90-word entity definition.
- FAQ answers are self-contained (30-70 words each) and backed by `FAQPage` JSON-LD. Google limits FAQ rich results, but the markup is harmless and useful to other engines.
- The comparison table on `/en/compare/aliases/` and the folder-contents table on `/en/guides/claude-config-dir/` are highly extractable.
- Honest disclaimers ("unofficial", "does not get around usage limits", "depends on each account's terms") build trust and make the answers safe for engines to quote.
- Specific, verifiable facts: "about 4 milliseconds", "Checked with Claude Code 2.1.302", the Keychain entry name `Claude Code-credentials`.
- Static SSR HTML (Next.js export on Cloudflare), canonical tags, full hreflang alternates for 30 locales, per-locale sitemaps.

## Top 5 changes by impact

| # | Change | Effort | Dimension |
|---|---|---|---|
| 1 | Seed third-party mentions: Show HN, Reddit write-up, YouTube/asciinema demo, awesome-list PRs (F1, F8) | Medium | Authority, Multi-modal |
| 2 | Expand "ccam (Claude Code Account Manager)" in the homepage, llms.txt, GitHub description and JSON-LD `alternateName` (F2) | Low | Authority, Citability |
| 3 | Link the official Claude Code docs for CLAUDE_CONFIG_DIR and usage limits (F3) | Low | Citability, Authority |
| 4 | Add question H2s covering switching without ccam, the active-account check, Windows/WSL and troubleshooting (F5, F6) | Medium | Citability, Structure |
| 5 | Add author/publisher/datePublished to the JSON-LD and redeploy on each release (F4, F7) | Low | Authority, Technical |

## Structured findings (audit-data.json, category "AI Search Readiness")

```json
{
  "category": "AI Search Readiness",
  "score": 66,
  "dimensions": {"citability": 78, "structural_readability": 85, "multi_modal": 45, "authority_brand": 20, "technical_accessibility": 95},
  "platforms": {"google_aio": 55, "chatgpt": 45, "perplexity": 60, "bing_copilot": 50},
  "llms_txt": "present_valid",
  "robots_ai_access": "all_allowed",
  "findings": [
    {"id": "GEO-1", "severity": "high", "title": "Near-zero off-site entity footprint", "evidence": "Repo created 2026-10-18, 0 stars/forks; no HN/Reddit/YouTube mentions", "fix": "Show HN, Reddit write-up, YouTube/asciinema demo, awesome-list PRs"},
    {"id": "GEO-2", "severity": "medium", "title": "Ambiguous name not expanded in machine-readable places", "evidence": "Wikipedia CCAM disambiguation; expansion only in header logo; no alternateName", "fix": "Use 'ccam (Claude Code Account Manager)' in home copy, llms.txt, repo description, JSON-LD alternateName"},
    {"id": "GEO-3", "severity": "medium", "title": "No outbound citations to official Claude Code docs", "evidence": "Only external link is the GitHub repo; 'Claude Code's documentation describes' is unlinked", "fix": "Link official CLAUDE_CONFIG_DIR and usage-limit docs"},
    {"id": "GEO-4", "severity": "medium", "title": "No author/publisher/datePublished in structured data", "evidence": "TechArticle has only dateModified; Organization sameAs = GitHub only", "fix": "Add Person author, publisher, datePublished, broader sameAs"},
    {"id": "GEO-5", "severity": "medium", "title": "Query coverage gaps; guides 281-408 words", "evidence": "No /logout-/login switch path, Windows/WSL, active-account check, troubleshooting", "fix": "Add question H2 sections or a 'switch Claude Code accounts' guide"},
    {"id": "GEO-6", "severity": "low", "title": "Guide H2s are labels, not questions", "evidence": "'By hand, with an alias', 'Before it happens'", "fix": "Rephrase a few H2s as real queries"},
    {"id": "GEO-7", "severity": "low", "title": "Stale softwareVersion in JSON-LD", "evidence": "Live JSON-LD 0.2.0 vs release v0.3.0", "fix": "Redeploy site as part of the release"},
    {"id": "GEO-8", "severity": "low", "title": "No images or video", "evidence": "0 <img> on content pages; no video", "fix": "Add an asciinema/GIF demo and a YouTube video"},
    {"id": "GEO-9", "severity": "info", "title": "Guide breadcrumb points to Documentation/install", "evidence": "BreadcrumbList position 2 = /en/docs/install/", "fix": "Use a Guides crumb"}
  ]
}
```
