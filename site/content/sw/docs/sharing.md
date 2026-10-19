---
title: Shiriki seva za MCP, mipangilio na historia
description: Jinsi ccam inavyoshiriki seva za MCP, settings.json, CLAUDE.md, agents, skills na historia ya mazungumzo kati ya akaunti za Claude Code; kuingia hubaki tofauti.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: 0df8a7d51c87
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

ccam haibadilishi kushiriki wakati Claude Code inaendeshwa kama akaunti hiyo, kwa sababu kipindi kinaendelea kuandika kwenye faili zinazohamishwa. Kifunge kwanza.

## Kuzima kushiriki

```sh
ccam share work --history off
```

Kiungo huondolewa na akaunti hupata nakala yake yenyewe ya maudhui ya sasa ya kikundi. Ongeza `--fresh` ili kuanza ikiwa tupu badala yake. Kikundi huendelea kuhifadhi maudhui yake kwa ajili ya akaunti nyingine.

## Nakili usanidi uliopo

Ikiwa ulitumia Claude Code kabla ya ccam, mipangilio yako iko katika `~/.claude`. `ccam migrate` huinakili kwenye kikundi, ili kila akaunti inayoshiriki vipengele hivyo ivipate:

```sh
ccam migrate
```

Kwenye terminal, ccam huonyesha kwanza orodha ya kuteua ya vipengele. Kila kitu isipokuwa historia ya mazungumzo kimeteuliwa; bonyeza Space kubadilisha kipengele na Enter kuona mpango. Vipengele ambavyo folda haina chochote kwavyo huwekwa alama `not found`.

Kinachonakiliwa:

- `settings.json` na `CLAUDE.md` ya jumla;
- viendelezi `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/` na `rules/`;
- seva zako za MCP za mtumiaji, yaani kipengele `mcpServers` cha `~/.claude.json`, kwenda `mcp.json` ya kikundi;
- taarifa za kuingia zilizohifadhiwa (OAuth) za seva za MCP inazonakili, kwenda kwenye faili ya kuingia ya kila akaunti yenyewe;
- historia ya mazungumzo (`projects/` na `history.jsonl`), lakini tu ukiongeza `--history`.

Folda chanzo husomwa tu: hakuna kitu ndani yake kinachohamishwa, kuunganishwa au kubadilishwa, na `claude` ya kawaida inaendelea kufanya kazi kama awali. Faili za kuingia `.credentials.json` na `.claude.json` hazinakiliwi kamwe kwenye kikundi.

Seva za MCP zinazoingia kwa OAuth huhifadhi tokeni zao kwenye taarifa za kuingia, si kwenye orodha ya seva. Ili usilazimike kuingia kwenye seva hizo tena, ccam hunakili taarifa hizo za kuingia kwenye faili ya kuingia ya kila akaunti inayopata seva kupitia kikundi, isipokuwa akaunti tayari ina zake. Huruka akaunti zinazoendeshwa, ambazo bado hazijaingia, zilizosajiliwa mahali zilipo kwa `ccam import`, au zinazohifadhi taarifa zao za kuingia katika Keychain ya macOS, na huzitaja; ingia kwenye seva huko kwa `/mcp`. Akaunti inayoanza kushiriki seva za MCP baadaye hupata taarifa hizo za kuingia unapoendesha `ccam migrate` tena.

ccam huchapisha mpango na huuliza kabla ya kunakili chochote:

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Accounts that share these through group default:
  mcp         personal
  settings    personal
  memory      personal
  extensions  personal
```

Mwishoni, ccam huorodhesha akaunti zinazoshiriki kila kipengele kupitia kikundi. Ikiwa baadhi ya akaunti hazishiriki, huchapisha amri za `ccam share` zinazowasha kushiriki kwa akaunti hizo.

Kikundi kikiwa tayari na nakala tofauti ya faili, ccam huuliza, kama inavyofanya unapowasha kushiriki: kuhifadhi nakala ya kikundi au kutumia yako. `settings.json` na faili nyingine za JSON zinaweza pia kuunganishwa, na `history.jsonl` mstari kwa mstari. Nakala za awali za kikundi huhifadhiwa katika `.ccam-backup/<time>/` ndani ya folda ya kikundi.

| Chaguo | Athari |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Nakili vipengele unavyovitaja tu, bila orodha ya kuteua; bila hivyo na nje ya terminal, kila kitu isipokuwa historia ya mazungumzo |
| `-g`, `--group <group>` | Nakili kwenye kikundi hiki badala ya `default` |
| `-y`, `--yes` | Nakili bila kuuliza; faili ikitofautiana, hifadhi nakala ya kikundi |

Weka `--dry-run` kabla ya amri, kama katika `ccam --dry-run migrate`, ili kuona mpango bila kunakili chochote. Ili kunakili kutoka folda nyingine, itaje: `ccam migrate ~/.claude-work --group work`.

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

Ikiwa kikundi tayari kina seva tofauti yenye jina hilo, ccam huiruka na kukujulisha. `--force` hubadilisha seva ya kikundi na huhifadhi `mcp.json` ya zamani kwenye folda ya `.ccam-backup/` ya kikundi.

Seva za MCP zinazoingia kwa OAuth huhifadhi tokeni zao pamoja na taarifa za kuingia za akaunti, kwa hiyo ziidhinishe mara moja katika kila akaunti, au nakili taarifa za kuingia kutoka `~/.claude` kwa `ccam migrate`.

Ikiwa unapendelea seva zionekane katika `claude mcp list` na `/mcp` kama seva za kawaida za mtumiaji, weka `mcp_strategy = "sync"` katika `~/.config/ccam/config.toml`. Kisha ccam huziandika katika `.claude.json` ya kila akaunti kabla ya kuanza, chini ya kufuli ya faili ya Claude Code yenyewe na kwa nakala rudufu.

## Akaunti mbili katika folda moja

Akaunti zinazoshiriki historia ya mazungumzo huandika kwenye faili zile zile. Ukianzisha Claude Code ya pili katika folda ambamo akaunti nyingine ya kikundi hicho hicho tayari inaendeshwa, ccam hukuonya na kuuliza kama uendelee.

## Kitu kikiharibika

```sh
ccam doctor
ccam doctor --fix
```

Ikiwa njia inayoshirikiwa imegeuka kuwa faili ya kawaida, kwa mfano baada ya kurejesha nakala rudufu kwa mkono, ccam husimama kabla ya kuanzisha Claude Code na huelekeza kwa `ccam doctor --fix`, ambayo huunganisha faili hiyo na kikundi na huhifadhi nakala rudufu.
