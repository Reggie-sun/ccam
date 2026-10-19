---
title: Kuanza haraka
description: Ongeza akaunti mbili za Claude Code kwa ccam add, zipe njia fupi kama ccwork, badili kati yao, na uchague zinachoshiriki.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: d040e34d2a17
---

Ongeza kila akaunti kwa `ccam add`, kisha anzisha Claude Code kwa `ccam <name>` au kwa njia fupi unayoipa jina mwenyewe, kama `ccwork`. Hivyo ndivyo watu wengi wanavyohitaji.

## 1. Ongeza akaunti yako ya kwanza

```sh
ccam add
```

Mchawi wa usanidi (wizard) huuliza jina kama `personal` au `work`, ni vipengele gani akaunti inapaswa kushiriki na akaunti zako nyingine, na kama uunde amri ya njia fupi (tazama hatua ya 4). Kisha huendesha `claude auth login`, ambayo hufungua kivinjari ili uweze kuingia. Kwenye mashine ya mbali, kivinjari kinaweza kuishia kwenye ukurasa wa `http://localhost:…/callback` ambao haufunguki; bandika anwani hiyo kwenye terminal na ccam hukamilisha kuingia. Ikimaliza utaona:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Leta mipangilio yako iliyopo

Ikiwa umewahi kutumia Claude Code, mipangilio yako, `CLAUDE.md`, viendelezi na seva za MCP viko katika `~/.claude`. Vinakili kwenye mipangilio inayoshirikiwa, ili kila akaunti ivipate:

```sh
ccam migrate
```

ccam huonyesha kitakachonakiliwa na huuliza kwanza. Husoma `~/.claude` tu, kamwe hainakili taarifa za kuingia, na haijumuishi historia ya mazungumzo isipokuwa uongeze `--history`. Tazama [kunakili usanidi uliopo](/docs/sharing#copy-an-existing-setup) kwa maelezo zaidi.

## 3. Ongeza akaunti ya pili

```sh
ccam add work
```

Ingia kwa akaunti nyingine kwenye kivinjari.

## 4. Badili akaunti

`ccam work` huanzisha Claude Code kama akaunti ya kazi. Kwa kitu kifupi zaidi, ipe akaunti amri yake ya njia fupi, kama `ccwork`. Baada ya vipengele vya kushiriki, mchawi wa usanidi huuliza:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Bonyeza `y`, kisha andika jina, au bonyeza Enter kutumia jina lililopendekezwa:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Baada ya kuingia, ccam huthibitisha:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

Ikiwa uliruka swali hilo, ongeza njia fupi wakati wowote:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Sasa yoyote kati ya hizi huanzisha Claude Code kama akaunti ya kazi:

```sh
ccwork
ccam work
ccam 2
```

Kuendesha `ccam` peke yake hufungua kichaguzi: tumia vitufe vya mishale na Enter, au bonyeza namba ya akaunti.

Kila kitu baada ya njia fupi au jina la akaunti hupelekwa kwa Claude Code bila kubadilishwa, kwa hiyo unaendelea kutumia chaguo unazozijua:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Chagua cha kushiriki

Kwa chaguo-msingi, seva za MCP, `settings.json`, `CLAUDE.md` ya jumla, pamoja na agents, commands, skills na plugins hushirikiwa kati ya akaunti. Historia ya mazungumzo haishirikiwi. Iwashe kwa akaunti zote mbili ikiwa unataka kuendeleza mazungumzo kati yao:

```sh
ccam share personal --history on
ccam share work --history on
```

ccam huonyesha kile itakachohamisha kabla ya kubadilisha chochote. Tazama [kushiriki](/docs/sharing) kwa maelezo zaidi.

## 6. Akaunti ikifikia kikomo chake

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
