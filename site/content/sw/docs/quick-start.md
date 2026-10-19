---
title: Kuanza haraka
description: Ongeza akaunti mbili za Claude Code kwa ccam add, badili kati yao kwa claude1, claude2 au ccam <name>, na uchague zinachoshiriki.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 6baf0bfa5029
---

Ongeza kila akaunti kwa `ccam add`, kisha anzisha Claude Code kwa `claude1`, `claude2` au `ccam <name>`. Hivyo ndivyo watu wengi wanavyohitaji.

## 1. Ongeza akaunti yako ya kwanza

```sh
ccam add
```

Mchawi wa usanidi (wizard) huuliza jina kama `personal` au `work`, na ni vipengele gani akaunti inapaswa kushiriki na akaunti zako nyingine. Kisha huendesha `claude auth login`, ambayo hufungua kivinjari ili uweze kuingia. Ikimaliza utaona:

```console
✓ personal (me@example.com) added — shortcut: claude1
```

ccam huchapisha ujumbe wake kwa Kiingereza; amri na chaguo ni zile zile katika kila lugha.

## 2. Ongeza akaunti ya pili

```sh
ccam add work
```

Ingia kwa akaunti nyingine kwenye kivinjari. Itakuwa `claude2`.

## 3. Badili

Yoyote kati ya hizi huanzisha Claude Code kama akaunti ya kazi:

```sh
claude2
ccam work
ccam 2
```

Kila kitu baada ya jina la akaunti hupelekwa kwa Claude Code bila kubadilishwa, kwa hiyo unaendelea kutumia chaguo unazozijua:

```sh
ccam work --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

Kuendesha `ccam` peke yake hufungua kichaguzi: tumia vitufe vya mishale na Enter, au bonyeza namba ya akaunti.

## 4. Chagua cha kushiriki

Kwa chaguo-msingi, seva za MCP, `settings.json`, `CLAUDE.md` ya jumla, pamoja na agents, commands, skills na plugins hushirikiwa kati ya akaunti. Historia ya mazungumzo haishirikiwi. Iwashe kwa akaunti zote mbili ikiwa unataka kuendeleza mazungumzo kati yao:

```sh
ccam share personal --history on
ccam share work --history on
```

ccam huonyesha kile itakachohamisha kabla ya kubadilisha chochote. Tazama [kushiriki](/docs/sharing) kwa maelezo zaidi.

## 5. Akaunti inapofikia kikomo chake

```sh
ccam usage
ccam next --continue
```

`ccam usage` huonyesha matumizi ya saa 5 na ya wiki ya kila akaunti. `ccam next --continue` hufungua mazungumzo yale yale kwenye akaunti inayofuata. Tazama [vikomo vya matumizi](/docs/usage-limits).

## Tayari unatumia ~/.claude?

Isajili badala ya kuingia tena:

```sh
ccam import
```

Folda hubaki pale ilipo, na `claude` ya kawaida inaendelea kufanya kazi kama awali.
