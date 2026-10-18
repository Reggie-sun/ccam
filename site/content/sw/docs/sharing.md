---
title: Shiriki seva za MCP, mipangilio na historia
description: Jinsi ccam inavyoshiriki seva za MCP, settings.json, CLAUDE.md, agents, skills na historia ya mazungumzo kati ya akaunti za Claude Code; kuingia hubaki tofauti.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: e4b31e5098e3
---

Kila akaunti huhifadhi taarifa zake za kuingia. Kila kitu kingine kinaweza kushirikiwa kipengele kwa kipengele: ccam huweka kiungo kutoka faili au folda ya akaunti kwenda folda ya kikundi kinachoshirikiwa, ili badiliko lililofanywa katika akaunti moja lionekane katika akaunti zote za kikundi hicho.

## Kinachoweza kushirikiwa

| Kipengele | Chaguo | Kinachoshirikiwa |
| --- | --- | --- |
| Seva za MCP | `--mcp` | `mcp.json` ya kikundi, hupitishwa kwa Claude Code wakati wa kuanza |
| Mipangilio | `--settings` | `settings.json`: ruhusa, hooks, model, mstari wa hali (status line) |
| Kumbukumbu ya jumla | `--memory` | `CLAUDE.md` |
| Viendelezi | `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| Historia ya mazungumzo | `--history` | `projects/` na `history.jsonl` |

Akaunti mpya hushiriki kila kitu isipokuwa historia ya mazungumzo.

Faili za kuingia `.credentials.json` na `.claude.json` hazishirikiwi kamwe. ccam huunganisha njia kutoka kwenye orodha maalum tu, na faili hizi mbili hazimo ndani yake.

## Vikundi

`on` inamaanisha kikundi cha `default`. Unaweza pia kutoa jina, ili akaunti fulani tu zishiriki kipengele:

```sh
ccam share work --history team
ccam share side --history team
```

Sasa `work` na `side` zinashiriki mazungumzo kati yao, lakini si na akaunti zako nyingine.

## Kuwasha kushiriki

Akaunti inapoanza kushiriki kipengele, ccam:

1. huchapisha mpango wa kila faili itakayohamisha, kunakili au kuunganisha, na huuliza kabla ya kuanza;
2. huhamisha faili za akaunti kwenda folda ya kikundi ikiwa kikundi bado kiko tupu, au huongeza faili ambazo kikundi bado hakina;
3. faili ile ile inapotofautiana, huuliza kama ubaki na toleo la kikundi, utumie lako, au uunganishe JSON;
4. huhifadhi faili yako asili katika `.ccam-backup/<time>/` ndani ya folda ya akaunti;
5. hubadilisha faili asili kwa kiungo kinachoelekea kwenye kikundi.

Faili za mazungumzo hupewa majina kwa kufuata session ID yake, kwa hiyo historia kutoka akaunti kadhaa huungana bila migongano.

## Kuzima kushiriki

```sh
ccam share work --history off
```

Kiungo huondolewa na akaunti hupata nakala yake yenyewe ya maudhui ya sasa ya kikundi. Ongeza `--fresh` ili kuanza ikiwa tupu badala yake. Kikundi huendelea kuhifadhi maudhui yake kwa ajili ya akaunti nyingine.

## Seva za MCP

Claude Code huhifadhi seva za MCP katika `.claude.json`, pamoja na taarifa za kuingia, kwa hiyo faili hiyo haiwezi kushirikiwa. ccam huhifadhi seva zinazoshirikiwa katika `mcp.json` ya kikundi na huanzisha Claude Code kwa `--mcp-config=<group>/mcp.json`. Chaguo zako mwenyewe za `--mcp-config` bado zinatumika juu yake.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

Seva unayoongeza ndani ya kipindi (session) kwa `claude mcp add` huhifadhiwa kwa akaunti hiyo tu. Wakati ujao utakapoanzisha akaunti hiyo, ccam hutaja seva hiyo na amri inayoihamisha kwenda kwenye kikundi:

```sh
ccam mcp import work gh
```

Seva za MCP zinazoingia kwa OAuth huhifadhi tokeni zao pamoja na taarifa za kuingia za akaunti, kwa hiyo ziidhinishe mara moja katika kila akaunti.

Ikiwa unapendelea seva zionekane katika `claude mcp list` na `/mcp` kama seva za kawaida za mtumiaji, weka `mcp_strategy = "sync"` katika `~/.config/ccam/config.toml`. Kisha ccam huziandika katika `.claude.json` ya kila akaunti kabla ya kuanza, chini ya kufuli ya faili ya Claude Code yenyewe na kwa nakala rudufu.

## Akaunti mbili katika folda moja

Akaunti zinazoshiriki historia ya mazungumzo huandika kwenye faili zile zile. Ukianzisha Claude Code ya pili katika folda ambamo akaunti nyingine ya kikundi hicho hicho tayari inaendeshwa, ccam hukuonya na kuuliza kama uendelee.

## Kitu kikiharibika

```sh
ccam doctor
ccam doctor --fix
```

Ikiwa njia inayoshirikiwa imegeuka kuwa faili ya kawaida, kwa mfano baada ya kurejesha nakala rudufu kwa mkono, ccam husimama kabla ya kuanzisha Claude Code na huelekeza kwa `ccam doctor --fix`, ambayo huunganisha faili hiyo na kikundi na huhifadhi nakala rudufu.
