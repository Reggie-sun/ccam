---
title: Vikomo vya matumizi na kubadili akaunti
description: Ona matumizi ya saa 5 na ya wiki ya kila akaunti ya Claude Code kwa ccam usage, na ufungue akaunti yenye nafasi kubwa zaidi kwa ccam next --auto.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 9aebb54d40e2
---

`ccam usage` huonyesha kiasi cha vikomo vya saa 5 na vya wiki ambacho kila akaunti imetumia. `ccam next --auto` hufungua akaunti yenye nafasi kubwa zaidi, na `--continue` huhamisha mazungumzo ya sasa.

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ██████░░  72% 1h 40m    ███░░░░░  41% 3d 6h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
 3 side      ████████ 100% 22m       █████░░░  68% 1d 9h     12m ago (cached)
→ Most headroom: work   (open it with ccam next --auto)
```

Safu ni: akaunti, matumizi ya saa 5, muda hadi yaanze upya, matumizi ya wiki, muda hadi yaanze upya, na umri wa namba hizo.

## Namba zinatoka wapi

ccam huuliza endpoint ya matumizi ya Anthropic kwa ajili ya kila akaunti kwa wakati mmoja, ikitumia taarifa za kuingia za kila akaunti yenyewe. Husubiri sekunde 3 kwa akaunti kwa kiwango cha juu na huonyesha kila akaunti mara tu inapojibu.

Endpoint hii haijaandikiwa nyaraka hadharani na inaweza kubadilika. Ombi linaposhindwa, ccam huonyesha namba za mwisho zinazojulikana na kuziwekea alama kuwa ni za cache. Kipindi cha muda kilichoanza upya tangu namba hizo zihifadhiwe huonyesha 0%. `ccam doctor` hukujulisha ikiwa muundo wa jibu umebadilika.

Namba zilizopatikana chini ya sekunde 60 zilizopita hutumiwa tena. Chaguo:

| Chaguo | Athari |
| --- | --- |
| `--watch [seconds]` | Onyesha upya mfululizo, kila sekunde 60 kwa chaguo-msingi |
| `--json` | Matokeo yanayosomeka na mashine |
| `-v` | Onyesha pia vikomo vya wiki vya kila model na kila namba ilikotoka |
| `--no-cache` | Uliza upya kila mara |
| `--source cache` | Usiulize kamwe; onyesha namba zilizohifadhiwa tu |

## Taarifa za kuingia zilizoisha muda

Tokeni ya ufikiaji ya akaunti inapoisha muda, ccam haiifanyi upya, kwa sababu kuifanya upya wakati kipindi cha Claude Code kinatumia taarifa zile zile za kuingia kunaweza kusababisha kipindi hicho kitolewe nje ya akaunti (sign out). Anzisha akaunti mara moja kwa `ccam <name>` na Claude Code itaifanya upya.

`ccam usage --refresh` hufanya upya tokeni zilizoisha muda yenyewe, lakini kwa akaunti ambazo haziendeshwi tu. Kamwe haifanyi upya taarifa za kuingia zilizohifadhiwa katika Keychain ya macOS au katika `~/.claude` iliyoingizwa.

## Matumizi kwenye mstari wa hali

```sh
ccam statusline install
```

Kisha Claude Code huonyesha jina la akaunti na matumizi yake kwenye mstari wa hali (status line), kwa mfano `[work] · Opus · 5h 24% · 7d 41%`. Ikiwa tayari una amri ya mstari wa hali, ccam huihifadhi na huweka jina la akaunti mbele yake. Mstari wa hali pia huhifadhi namba zake, ambazo `ccam usage` hutumia wakati endpoint haijibu.

## Kubadili

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # the next account that shares your conversation history
```

`ccam next` huruka akaunti ambazo hazijaingia, na `--auto` huruka akaunti zilizofikia 100%. Ili kuendeleza mazungumzo, akaunti zote mbili lazima zishiriki historia ya mazungumzo; tazama [kushiriki](/docs/sharing).

[Mwongozo wa hatua za kufuata kikomo cha matumizi kinapofikiwa](/guides/usage-limit-workflow) unapitia mfano kamili.
