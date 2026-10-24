---
title: Shiriki seva za MCP, mipangilio na historia
description: Jinsi vikundi vya ccam vinavyoshiriki seva za MCP, settings.json, CLAUDE.md, skills na historia kati ya akaunti za Claude Code; kuingia hubaki tofauti.
updated: 2026-10-23
order: 4
translation: machine
sourceHash: "27e0a0e5abbd"
---

Kila akaunti huhifadhi taarifa zake za kuingia. Kila kitu kingine kinaweza kushirikiwa kupitia **vikundi**: kikundi hutaja vipengele kinachoshiriki na akaunti zilizomo, na ccam huunganisha faili za kila moja ya akaunti hizo na folda ya kikundi, ili badiliko lililofanywa katika akaunti moja lionekane katika zote.

## Kinachoweza kushirikiwa

| Kipengele | Kinachoshirikiwa |
| --- | --- |
| `mcp` | seva za MCP: `mcp.json` ya kikundi, hupitishwa kwa Claude Code wakati wa kuanza |
| `settings` | `settings.json`: ruhusa, hooks, model, mstari wa hali (status line) |
| `memory` | `CLAUDE.md` ya jumla |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | historia ya mazungumzo: `projects/` na `history.jsonl` |

Akaunti yako ya kwanza huunda kikundi `default`, ambacho hushiriki kila kitu isipokuwa historia ya mazungumzo, na akaunti zinazofuata hujiunga nacho.

Faili za kuingia `.credentials.json` na `.claude.json` hazishirikiwi kamwe. ccam huunganisha njia kutoka kwenye orodha maalum tu, na faili hizi mbili hazimo ndani yake.

Ukurasa huu unashughulikia kazi za kila siku. Kila chaguo kimeelezwa kwenye kurasa za vipengele: [kushiriki kwa vikundi](/features/share), [seva za MCP zinazoshirikiwa](/features/mcp) na [kunakili usanidi uliopo](/features/migrate).

## Vikundi

```sh
ccam group                                   # the groups, what they share, their accounts
ccam group set default --add history         # every account in default shares history
ccam group add team --items history --members work,side
```

Akaunti inaweza kuwa katika vikundi kadhaa, lakini kila kipengele hutoka kwenye kikundi kimoja tu. Kwa hiyo ili kushiriki historia kati ya `work` na `side` lakini si na akaunti zako nyingine, iache historia nje ya `default` na uzipe akaunti hizo mbili kikundi chao wenyewe, kama ilivyo hapo juu. `ccam share work` huonyesha kile ambacho akaunti inashiriki na kupitia kikundi kipi. [Dashibodi](/features/dashboard) hufanya vivyo hivyo kwenye kichupo chake cha Groups.

## Akaunti inapoanza kushiriki

Akaunti inapoanza kushiriki kipengele (inajiunga na kikundi, au kikundi chake kinaongeza kipengele hicho), ccam:

1. huchapisha mpango wa kila faili itakayohamisha, kunakili au kuunganisha kwa kila akaunti, na huuliza kabla ya kuanza;
2. huhamisha faili za akaunti kwenda folda ya kikundi ikiwa kikundi bado kiko tupu, au huongeza faili ambazo kikundi bado hakina;
3. faili ile ile inapotofautiana, huuliza kama ubaki na toleo la kikundi, utumie lako, au uunganishe JSON;
4. huhifadhi faili yako asili katika `.ccam-backup/<time>/` ndani ya folda ya akaunti;
5. hubadilisha faili asili kwa kiungo kinachoelekea kwenye kikundi.

Faili za mazungumzo hupewa majina kwa kufuata session ID yake, kwa hiyo historia kutoka akaunti kadhaa huungana bila migongano.

ccam haibadilishi akaunti wakati Claude Code inaendeshwa kama akaunti hiyo, bila kujali Claude Code ilianzishwa vipi (pamoja na vipindi vya chinichini), kwa sababu kipindi kinaendelea kuandika kwenye faili zinazohamishwa. Kifunge kwanza. ccam hukagua tena kabla tu ya kuhifadhi, kwa hiyo kipindi kilichoanzishwa wakati mpango ukisubiri nacho hugunduliwa.

## Akaunti inapoacha kushiriki

```sh
ccam group leave team work                   # work leaves the group
ccam group set default --remove history      # nobody in default shares history any more
```

Kiungo huondolewa na akaunti hupata nakala yake yenyewe ya maudhui ya sasa ya kikundi; seva za MCP za kikundi huingia katika `.claude.json` yake yenyewe, kando ya seva zake yenyewe. Ongeza `--fresh` ili kuanza ikiwa tupu badala yake; hapo mpango huorodhesha seva za MCP zitakazoacha kupakiwa. Folda iliyoingizwa kama `~/.claude` daima hubaki na nakala. Kikundi huendelea kuhifadhi maudhui yake kwa ajili ya akaunti nyingine.

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
- historia ya mazungumzo (`projects/` na `history.jsonl`), lakini tu ukiongeza `--history`.

Folda chanzo husomwa tu: hakuna kitu ndani yake kinachohamishwa, kuunganishwa au kubadilishwa, na `claude` ya kawaida inaendelea kufanya kazi kama awali. Faili za kuingia `.credentials.json` na `.claude.json` hazinakiliwi kamwe kwenye kikundi.

Seva za MCP zinazoingia kwa OAuth huhifadhi tokeni zao kwenye taarifa za kuingia, si kwenye orodha ya seva. ccam hainakili kuingia huko: kuingia kulikonakiliwa hubeba refresh token ile ile, na kufanywa upya kwa mara ya kwanza huwatoa nje ya seva hiyo wengine wote walio na nakala yake, ikiwemo `~/.claude`. Mpango na muhtasari hutaja akaunti zinazopaswa kuingia kwenye seva hizo zenyewe kwa `/mcp`.

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
Group default: personal
To share them:
  ccam group join default work
```

Mwishoni, ccam huorodhesha akaunti za kikundi. Ikiwa kikundi bado hakishiriki mojawapo ya vipengele, au baadhi ya akaunti hazimo ndani yake, huchapisha amri za `ccam group` zinazobadilisha hilo.

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

Seva unayoongeza ndani ya kipindi (session) kwa `claude mcp add` huhifadhiwa kwa akaunti hiyo tu, na kwa folda uliyokuwamo (local scope ya Claude Code; `-s user` huifanya iwe ya akaunti katika kila folda). Wakati ujao utakapoanzisha akaunti hiyo, ccam hutaja seva hiyo, pamoja na folda yake, na amri inayoihamisha kwenda kwenye kikundi:

```sh
ccam mcp import work gh
```

Ikiwa kikundi tayari kina seva tofauti yenye jina hilo, ccam huiruka na kukujulisha. `--force` hubadilisha seva ya kikundi na huhifadhi `mcp.json` ya zamani kwenye folda ya `.ccam-backup/` ya kikundi.

Seva za MCP zinazoingia kwa OAuth huhifadhi tokeni zao pamoja na taarifa za kuingia za akaunti, kwa hiyo ziidhinishe mara moja katika kila akaunti kwa `/mcp`. `ccam ext` na `ccam doctor` hutaja seva ambazo akaunti haijaingia wakati akaunti nyingine imeingia.

Ikiwa unapendelea seva zionekane katika `claude mcp list` na `/mcp` kama seva za kawaida za mtumiaji, weka `mcp_strategy = "sync"` juu kabisa ya `~/.config/ccam/config.toml`, kabla ya `[[group]]`, `[[profile]]` au `[[folder]]` ya kwanza. Kisha ccam huziandika katika `.claude.json` ya kila akaunti kabla ya kuanza, chini ya kufuli ya faili ya Claude Code yenyewe na kwa nakala rudufu. Seva ambayo akaunti inayo yenyewe kwa jina lile lile hubaki kuwa ya akaunti, na ccam husema hivyo.

## Akaunti mbili katika folda moja

Akaunti zinazoshiriki historia ya mazungumzo huandika kwenye faili zile zile. Ukianzisha Claude Code ya pili katika folda ambamo akaunti nyingine ya kikundi hicho hicho tayari inaendeshwa, ccam hukuonya na kuuliza kama uendelee.

## Kitu kikiharibika

```sh
ccam doctor
ccam doctor --fix
```

Ikiwa njia inayoshirikiwa imegeuka kuwa faili ya kawaida, kwa mfano baada ya kurejesha nakala rudufu kwa mkono, ccam husimama kabla ya kuanzisha Claude Code na huelekeza kwa `ccam doctor --fix`, ambayo huiunganisha tena na kikundi: faili hiyo huwa nakala ya kikundi ikiwa kikundi hakina, vinginevyo kikundi hubaki na yake na yako huenda kwenye nakala rudufu, ambayo `--fix` huitaja. Badiliko la vikundi lililosimama nusu njia pia hukamilishwa na `ccam doctor --fix`.
