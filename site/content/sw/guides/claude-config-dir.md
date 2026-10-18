---
title: CLAUDE_CONFIG_DIR imefafanuliwa
description: CLAUDE_CONFIG_DIR hufanya nini katika Claude Code, nini huhifadhiwa kwenye folda, macOS huweka wapi taarifa za kuingia, na faili zipi zinaweza kushirikiwa.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

`CLAUDE_CONFIG_DIR` huiambia Claude Code itumie folda ipi badala ya `~/.claude`. Kila kitu kinachohusu akaunti mahususi kiko humo: taarifa za kuingia, mipangilio, seva za MCP, plugins na historia ya mazungumzo. Kuielekeza kwenye folda tofauti hukupa akaunti tofauti.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## Kilicho ndani ya folda

Imekaguliwa kwa Claude Code 2.1.302:

| Njia | Yaliyomo | Inaweza kushirikiwa |
| --- | --- | --- |
| `.credentials.json` | Tokeni za kuingia (Linux), tokeni za MCP OAuth | Hapana |
| `.claude.json` | Maelezo ya akaunti, seva za MCP za mtumiaji, hali ya kila mradi | Hapana (seva za MCP kando) |
| `settings.json` | Ruhusa, hooks, model, mstari wa hali (status line), plugins | Ndiyo |
| `CLAUDE.md` | Maagizo ya jumla | Ndiyo |
| `agents/`, `commands/`, `skills/`, `plugins/` | Viendelezi | Ndiyo |
| `projects/` | Mazungumzo kama faili za `.jsonl`, kumbukumbu ya kila mradi | Ndiyo |
| `history.jsonl` | Historia ya maagizo (prompts) | Ndiyo |
| `sessions/`, `shell-snapshots/`, `cache/` na mengineyo | Hali ya wakati wa kuendesha | Hapana |

## Mambo mawili yanayowachanganya watu

**Bila `CLAUDE_CONFIG_DIR`, `.claude.json` haimo ndani ya `~/.claude`.** Hapo Claude Code huiweka katika `~/.claude.json`, kando ya folda. Ukianza kutumia `CLAUDE_CONFIG_DIR=~/.claude`, Claude Code hutafuta `~/.claude/.claude.json` badala yake na haipati akaunti yako. `ccam import` hushughulikia hili kwa kuziunganisha hizo mbili.

**Kwenye macOS taarifa za kuingia ziko katika Keychain, si kwenye folda.** Claude Code huzihifadhi chini ya jina `Claude Code-credentials`, likifuatiwa na hash fupi ya njia ya folda wakati `CLAUDE_CONFIG_DIR` imewekwa. Kwa hiyo folda tofauti hupata maingizo tofauti ya Keychain. Ikiwa Keychain itakataa kuandika, Claude Code hurudi kwenye `.credentials.json` ndani ya folda.

## Kushiriki kati ya folda

Viungo vya ishara (symbolic links) hufanya kazi: Claude Code husoma na kuandika kupitia `settings.json` au folda ya `projects/` iliyounganishwa na huacha kiungo kama kilivyo. Seva za MCP ndizo zisizofuata kanuni hii, kwa sababu zinashiriki `.claude.json` na taarifa za kuingia; badala yake zipitishe kwa `--mcp-config`.

ccam imejengwa juu ya kanuni hizi hizi. Tazama [kushiriki](/docs/sharing) kuona jinsi inavyounganisha folda, na [mwongozo wa akaunti nyingi](/guides/multiple-claude-code-accounts) kwa picha kamili.
