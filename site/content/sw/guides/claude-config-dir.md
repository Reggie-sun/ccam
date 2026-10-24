---
title: CLAUDE_CONFIG_DIR imefafanuliwa
description: CLAUDE_CONFIG_DIR hufanya nini katika Claude Code, nini huhifadhiwa kwenye folda, macOS huweka wapi taarifa za kuingia, na faili zipi zinaweza kushirikiwa.
updated: 2026-10-23
order: 2
translation: machine
sourceHash: "2cf08478085f"
---

`CLAUDE_CONFIG_DIR` huiambia Claude Code itumie folda ipi badala ya `~/.claude`. Kila kitu kinachohusu akaunti mahususi kiko humo: taarifa za kuingia, mipangilio, seva za MCP, plugins na historia ya mazungumzo. Kuielekeza kwenye folda tofauti hukupa akaunti tofauti.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## Je, CLAUDE_CONFIG_DIR inaungwa mkono rasmi?

Ndiyo. Claude Code huiorodhesha miongoni mwa [vigezo vyake vya mazingira](https://code.claude.com/docs/en/env-vars) kama njia ya kubadilisha folda ya usanidi, "useful for running multiple accounts side by side", na [mwongozo wake wa uthibitishaji](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) huitumia kuingia kwa akaunti kadhaa. Unaweza kuiweka kwenye shell yako au kwenye mipangilio ya mtumiaji, lakini si kwenye mipangilio ya mradi.

## Kuna nini ndani ya folda?

Imekaguliwa kwa Claude Code 2.1.309:

| Njia | Yaliyomo | Inaweza kushirikiwa |
| --- | --- | --- |
| `.credentials.json` | Tokeni za kuingia (Linux, Windows), tokeni za MCP OAuth | Hapana |
| `.claude.json` | Maelezo ya akaunti, seva za MCP za mtumiaji, hali ya kila mradi | Hapana (seva za MCP kando) |
| `settings.json` | Ruhusa, hooks, model, mstari wa hali (status line), plugins | Ndiyo |
| `CLAUDE.md` | Maagizo ya jumla | Ndiyo |
| `agents/`, `commands/`, `skills/`, `plugins/` | Viendelezi | Ndiyo |
| `projects/` | Mazungumzo kama faili za `.jsonl`, kumbukumbu ya kila mradi | Ndiyo |
| `history.jsonl` | Historia ya maagizo (prompts) | Ndiyo |
| `sessions/`, `shell-snapshots/`, `cache/` na mengineyo | Hali ya wakati wa kuendesha | Hapana |

## Ni nini huwachanganya watu?

**Bila `CLAUDE_CONFIG_DIR`, `.claude.json` haimo ndani ya `~/.claude`.** Hapo Claude Code huiweka katika `~/.claude.json`, kando ya folda. Ukianza kutumia `CLAUDE_CONFIG_DIR=~/.claude`, Claude Code hutafuta `~/.claude/.claude.json` badala yake na haipati akaunti yako. `ccam import` hushughulikia hili kwa kuziunganisha hizo mbili. Ikiwa `~/.claude` tayari ina `.claude.json` yake yenyewe, kutoka uendeshaji wa awali kwa `CLAUDE_CONFIG_DIR=~/.claude`, ccam huziacha faili zote mbili na husema hivyo, kwa kuwa ni wewe tu unayejua ipi ni sahihi.

**Kwenye macOS taarifa za kuingia ziko katika Keychain, si kwenye folda.** Claude Code huzihifadhi chini ya jina `Claude Code-credentials`, likifuatiwa na `-` na tarakimu 8 za kwanza za hex za SHA-256 ya njia ya folda wakati `CLAUDE_CONFIG_DIR` imewekwa, kwa akaunti inayotajwa na `$USER`. Kwa hiyo folda tofauti hupata maingizo tofauti ya Keychain, na folda inayohamishwa au kubadilishwa jina hupoteza taarifa zake za kuingia; ndiyo sababu `ccam rename` huhifadhi jina la folda kwenye macOS. Ikiwa Keychain itakataa kuandika, Claude Code hurudi kwenye `.credentials.json` ndani ya folda, kama [nyaraka zake](https://code.claude.com/docs/en/authentication#credential-management) zinavyoeleza.

**Si kila aina ya kuingia huhifadhiwa kwenye folda.** Kuingia kwa Claude Console bila API key huhifadhiwa nje yake, kwa hiyo folda tofauti haziwezi kutenganisha kuingia mbili za aina hiyo. Kuingia kwa claude.ai na API keys huhifadhiwa kwa kila folda.

## Folda iko wapi kwenye Windows?

Kwa chaguo-msingi ni `%USERPROFILE%\.claude`, na taarifa za kuingia ziko katika `.credentials.json` ndani yake, au katika Credential Manager ya Windows Claude Code inapowekwa kuitumia, chini ya jina linalotokana na njia ya folda kama kwenye macOS. Weka `CLAUDE_CONFIG_DIR` kwenye folda nyingine kwa kila akaunti ya ziada, kama kwenye Linux na macOS.

## Je, folda zinaweza kushiriki faili?

Viungo vya ishara (symbolic links) hufanya kazi: Claude Code husoma na kuandika kupitia `settings.json` au folda ya `projects/` iliyounganishwa na huacha kiungo kama kilivyo. Seva za MCP ndizo zisizofuata kanuni hii, kwa sababu zinashiriki `.claude.json` na taarifa za kuingia; badala yake zipitishe kwa `--mcp-config`.

ccam imejengwa juu ya kanuni hizi hizi. Tazama [kushiriki](/docs/sharing) kuona jinsi inavyounganisha folda, na [mwongozo wa akaunti nyingi](/guides/multiple-claude-code-accounts) kwa picha kamili.
