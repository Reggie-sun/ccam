# ccam.dev design direction

Status: written by the agent, not yet confirmed by the owner (docs/antislop-review.md #4). Change anything here and the site follows.

Reading this as: a documentation-led product site for developers in many countries, in a calm technical-reference style with one bold data visual. Dial: ENERGY 2 / RHYTHM 2 / MOTION 2.

## Identity

- Subject: a terminal tool that runs several Claude Code accounts and shows how much of each account's limit is left.
- Motif: the usage gauge. The hero board, the brand mark (pending owner decision, review #3) and the OG image strip are all gauges, because quota is what ccam makes visible.
- Voice: plain, specific, second person, no marketing words. Every page answers its question in the first paragraph.

## Palette

| Token | Light | Dark | Reason |
| --- | --- | --- | --- |
| paper | `#fafbfd` | `#121a22` | Cool white and blue-black: reads as documentation, avoids the cream and near-black defaults. |
| ink | `#18212b` | `#e6ebf0` | Text, borders of the two framed elements (install command, usage board). |
| ink-2 | `#4a5866` | `#a5b1bd` | Secondary text; 6.5:1 or better everywhere. |
| amber | `#a86a0c` | `#e9a63a` | The one accent: gauge fills and the translation notice only, like an analog meter. |
| teal | `#0e5a5f` | `#5cc3c8` | Links and the focus ring; separate from the accent so links never look like gauges. |
| full | `#b4442f` | `#e07a62` | A gauge at 100%; always shown with its percentage, never colour alone. |

All text and component pairs are checked with the antislop contrast checker (review #1).

## Type

- Atkinson Hyperlegible Next, subset to Latin (18 KB): designed for legibility, fits a site read in many second languages.
- Every non-Latin script uses the system font, so no large web font downloads.
- Commands use the system monospace; it marks literal input, not a "terminal look".

## Layout

- Left-aligned single column, 68 characters per line for prose.
- Landing: hero (headline, install command, usage board) → how it works as the real commands in order → guides → the four key facts.
- Docs: sidebar of every page grouped by section, article, previous/next.
- Reason: the content is commands and answers, so the layout is a reference layout, not a card grid.

## Rules for this site

- One framed element per screen at most (the install command or the board); everything else sits on the page.
- The amber stripe on the machine-translation notice marks a real state; no decorative stripes elsewhere.
- Motion: the gauges fill once on load and pages cross-fade with view transitions; both off under reduced motion. Nothing loops.
- Numbers on the page are either measured (launch time) or labelled as an example (the usage board).
