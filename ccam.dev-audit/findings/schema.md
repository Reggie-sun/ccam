# Schema / Structured Data Audit: ccam.dev

Score: 72/100

## Detection (live, JSON-LD only, no Microdata/RDFa)
| Page | Blocks |
|---|---|
| / and /ko/ | SoftwareApplication+Offer; [WebSite, Organization] |
| /en/docs/quick-start/, guides, blog, compare/aliases, ja faq | [TechArticle, BreadcrumbList]; [WebSite, Organization] |
| /en/docs/faq/, /ja/docs/faq/ | above + FAQPage (13 Q&A) |
| /en/changelog/ | only [WebSite, Organization]; no page-level schema |

All blocks use https://schema.org, absolute URLs, ISO dates, no placeholders. Syntax valid.

## Findings

1. HIGH - Live site is stale vs worktree. Live softwareVersion "0.2.0"; worktree lib/jsonld.ts VERSION = "0.3.0". Deploy/rebuild needed, or schema reports the wrong version.
2. MEDIUM - Breadcrumb position 2 and 3 share the same URL on blog and compare pages (e.g. /en/blog/introducing-ccam/ twice; /en/compare/aliases/ twice). Position 2 "Blog"/"Comparisons" should point to a section index (if one exists), or be dropped; the last crumb is allowed to omit `item`. Docs/guides pages are OK-ish but guide's "Documentation" crumb points to /en/docs/install/ (acceptable if no docs index).
3. MEDIUM - Blog post uses TechArticle with no datePublished. Use BlogPosting with datePublished and dateModified. Also dateModified is 2026-10-19 (build date) on every page, which is not a true modification date; derive it from content frontmatter/git. Misleading dates can be ignored by Google.
4. MEDIUM - Article lacks author and publisher. Recommended by Google Article guidelines. Add author/publisher as Organization with @id.
5. MEDIUM - FAQPage answer text is corrupted: "CLAUDECONFIGDIR" (should be CLAUDE_CONFIG_DIR; underscores eaten by markdown stripping) and `ccam share <name --history on` / `~/.local/share/ccam/profiles/<name/` (closing ">" lost) in both en and ja. Fix the markdown-to-text conversion in the FAQ extractor (faqLd input). Affects AI/GEO text quality.
6. LOW - Organization/WebSite are emitted on every page without @id, so entities are not linked across pages; TechArticle isPartOf redefines WebSite inline. Add @id and reference. Organization logo is 512px PNG: fine. Add `description`, and consider `publisher` link.
7. LOW - Homepage SoftwareApplication is not eligible for the Google software-app rich result (needs aggregateRating or review besides offers). Do not fabricate ratings. Keep as entity markup. `operatingSystem`, `applicationCategory` (DeveloperApplication), `offers` OK. Add `author`/`publisher`, `screenshot`/`image`, `sameAs`, `softwareRequirements` optional. Consider SoftwareSourceCode link to repo via `codeRepository` (not Google-rich but valid).
8. LOW - WebSite lacks `potentialAction` SearchAction: only add if a site search exists (Google no longer shows the sitelinks searchbox, so skip).
9. LOW - /en/changelog/ has no page schema. Add WebPage/CollectionPage + BreadcrumbList.
10. LOW - WebSite `inLanguage` and Organization name repeat per-locale; WebSite url always "/" rather than locale root. Acceptable; optionally use locale URL for url on non-en.
11. INFO - FAQPage present on /docs/faq/ in all locales. No Google SERP benefit since May 7, 2026; AI benefit unconfirmed. Not recommending add or remove for SERP; harmless to keep.
12. INFO - Deprecated types (HowTo, SpecialAnnouncement, etc.): none found. Quick-start and guide pages correctly do not use HowTo.

## Recommended JSON-LD

Blog post (replace TechArticle on /blog/ pages):
```json
{
  "@context": "https://schema.org",
  "@type": "BlogPosting",
  "headline": "Introducing ccam",
  "description": "Why ccam exists, what it does for people with several Claude Code accounts, and how it stays out of the way.",
  "inLanguage": "en",
  "mainEntityOfPage": "https://ccam.dev/en/blog/introducing-ccam/",
  "datePublished": "YYYY-MM-DD (real publish date from frontmatter)",
  "dateModified": "YYYY-MM-DD (real edit date)",
  "author": { "@type": "Person", "name": "jhyunwoo", "url": "https://github.com/jhyunwoo" },
  "publisher": { "@id": "https://ccam.dev/#organization" },
  "isPartOf": { "@id": "https://ccam.dev/#website" }
}
```
(Replace the date placeholders with real values before shipping.)

Linked entities (homepage and every page):
```json
[
  { "@context": "https://schema.org", "@type": "WebSite", "@id": "https://ccam.dev/#website", "name": "ccam", "url": "https://ccam.dev/", "inLanguage": "en", "publisher": { "@id": "https://ccam.dev/#organization" } },
  { "@context": "https://schema.org", "@type": "Organization", "@id": "https://ccam.dev/#organization", "name": "ccam", "url": "https://ccam.dev/", "logo": "https://ccam.dev/icon-512.png", "sameAs": ["https://github.com/jhyunwoo/ccam"] }
]
```

Corrected breadcrumb for blog post (omit duplicate; last item without `item` is allowed):
```json
{
  "@context": "https://schema.org",
  "@type": "BreadcrumbList",
  "itemListElement": [
    { "@type": "ListItem", "position": 1, "name": "ccam", "item": "https://ccam.dev/en/" },
    { "@type": "ListItem", "position": 2, "name": "Blog" },
    { "@type": "ListItem", "position": 3, "name": "Introducing ccam" }
  ]
}
```
Better: position 2 gets `item` only if a /en/blog/ index exists.

Changelog:
```json
{
  "@context": "https://schema.org",
  "@type": "WebPage",
  "name": "Changelog",
  "url": "https://ccam.dev/en/changelog/",
  "inLanguage": "en",
  "isPartOf": { "@id": "https://ccam.dev/#website" }
}
```
