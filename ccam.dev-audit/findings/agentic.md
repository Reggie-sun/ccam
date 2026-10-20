# Agent readiness: ccam.dev (checked 2026-10-19)

Lighthouse Agentic Browsing: not measured. PSI returned HTTP 429 (daily quota, no API key) for mobile and desktop. Retry with a Google API key.
Agent-UX heuristic (separate from Lighthouse): 100/100, complete, on / and /en/docs/quick-start/ (0 unnamed interactive nodes, 0 div-onclick widgets, 9 landmarks).

## Access policy
- Training: allowed (robots.txt `User-agent: *` / `Allow: /`, no AI groups, no Content-Signal)
- Search: allowed (same group)
- User-triggered agents: allowed (nothing blocked)

## Findings
- P0 pass: server-rendered, 689 words without JS; robots.txt 200.
- P1 pass: llms.txt 200, follows Lighthouse rules; links to .md pages that resolve (200, text/markdown on /en/docs/quick-start.md).
- P1 info/opportunity: no Content-Signal (draft/proposal, checked 2026-09-23). Fix (optional): add `Content-Signal: search=yes, ai-input=yes, ai-train=no` per policy.
- P2 low: Accept: text/markdown on HTML URLs returns text/html, no Vary: Accept, no rel=alternate markdown link. Fix: add `<link rel="alternate" type="text/markdown" href="...md">` on docs pages. Homepage /index.md is 404.
- P2/P3 opportunity: no WebMCP (draft), ai-catalog.json (draft) 404, no api-catalog/agent-card/OAuth well-known. Not defects for a static docs site.
