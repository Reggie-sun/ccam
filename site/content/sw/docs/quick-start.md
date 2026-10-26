---
title: Kuanza haraka
metaTitle: "Kuanza haraka: ongeza, badili akaunti za Claude Code"
description: Ongeza akaunti mbili za Claude Code kwa ccam add, zipe njia fupi kama ccwork, badili kati yao, na uchague zinachoshiriki.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: "bd893bbb83ab"
---

Ongeza kila akaunti kwa `ccam add`, kisha anzisha Claude Code kwa `ccam <name>` au kwa njia fupi unayoipa jina mwenyewe, kama `ccwork`. Hivyo ndivyo watu wengi wanavyohitaji.

Tayari unatumia Claude Code? Soma [Tayari unatumia ~/.claude?](#already-using-claude) mwishoni mwa ukurasa huu kabla ya kuanza.

## 1. Ongeza akaunti yako ya kwanza

```sh
ccam add
```

Mchawi wa usanidi (wizard) huuliza jina kama `personal` au `work`, ni vipengele gani akaunti inapaswa kushiriki na akaunti zako nyingine, na kama uunde amri ya njia fupi (tazama hatua ya 3). Kisha huendesha `claude auth login`, ambayo hufungua kivinjari ili uweze kuingia. Kwenye mashine ya mbali, kivinjari chako kinaweza kuishia kwenye ukurasa wa `http://localhost:…/callback` ambao haufunguki; bandika anwani hiyo kwenye terminal na ccam hukamilisha kuingia. Kuingia kukikamilika, utaona:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Leta mipangilio yako iliyopo

Ikiwa umewahi kutumia Claude Code, mipangilio yako, `CLAUDE.md` na viendelezi viko katika `~/.claude`, na seva zako za MCP ziko katika `~/.claude.json`. Vinakili kwenye kikundi cha `default` kinachoshirikiwa ili kila akaunti ivipate:

```sh
ccam migrate
```

ccam huonyesha kitakachonakiliwa na huuliza kwanza. Husoma `~/.claude` tu, kamwe hainakili taarifa za kuingia, na haijumuishi historia ya mazungumzo isipokuwa uongeze `--history`. Tazama [kunakili usanidi uliopo](/docs/sharing#copy-an-existing-setup) kwa maelezo zaidi.

## 3. Ongeza akaunti ya pili

```sh
ccam add work
```

Mchawi wa usanidi huruka swali la jina, ulilolitoa kwenye mstari wa amri, na huonyesha vikundi vya kujiunga huku `default` ikiwa imeteuliwa. Safari hii, ipe akaunti amri yake ya njia fupi, kama `ccwork`, ili uweze kuianzisha kwa neno moja. Baada ya vikundi, mchawi wa usanidi huuliza:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Bonyeza `y`, kisha andika jina, au bonyeza Enter kutumia jina lililopendekezwa:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Kisha ingia kwa akaunti nyingine kwenye kivinjari. ccam huthibitisha akaunti na njia yake fupi:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

## 4. Badili akaunti

`ccam work` huanzisha Claude Code kama akaunti ya kazi. Ikiwa uliruka swali la njia fupi, ongeza njia fupi wakati wowote:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Sasa kila moja ya amri hizi huanzisha Claude Code kama akaunti ya kazi (`2` ni namba yake katika `ccam list`):

```sh
ccwork
ccam work
ccam 2
```

Kuendesha `ccam` peke yake hufungua [dashibodi](/features/dashboard): nenda kwenye akaunti kwa vitufe vya mishale au namba yake, kisha bonyeza Enter.

Kila kitu baada ya njia fupi au jina la akaunti hupelekwa kwa Claude Code bila kubadilishwa, kwa hiyo unaendelea kutumia chaguo unazozijua:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Chagua cha kushiriki

Kwa chaguo-msingi, seva za MCP, `settings.json`, `CLAUDE.md` ya jumla na viendelezi (agents, commands, skills na plugins) hushirikiwa kati ya akaunti kupitia kikundi cha `default`. Historia ya mazungumzo haishirikiwi. Iongeze kwenye kikundi ikiwa unataka kuendeleza mazungumzo kati ya akaunti:

```sh
ccam group set default --add history
```

ccam huonyesha kile itakachohamisha kwa kila akaunti kabla ya kubadilisha chochote. Tazama [kushiriki kwa vikundi](/features/share) kwa maelezo zaidi.

## 6. Akaunti ikifikia kikomo chake

```sh
ccam usage
ccam next --continue
```

`ccam usage` huonyesha matumizi ya saa 5 na ya wiki ya kila akaunti. `ccam next --continue` hufungua mazungumzo yale yale kwenye akaunti inayofuata. Katika script, `ccam --fallback work -p "…"` hubadili yenyewe. Tazama [vikomo vya matumizi](/docs/usage-limits).

**Viti (seats) vya Team na Enterprise.** Uchaguzi wa akaunti wa kiotomatiki (`ccam next`, `--fallback`, `ccam pick`, kitufe `n` cha dashibodi na, ukishakiwasha, [`auto_switch`](/features/next#auto-switch)) pia hubadili kati ya viti vya shirika lile lile la Claude Team au Enterprise, jambo ambalo linaweza kukiuka masharti ya Anthropic. Soma [Maswali yanayoulizwa mara kwa mara](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) na uwasiliane na shirika lako na Anthropic kabla ya kulitegemea.

## Tayari unatumia ~/.claude?

Isajili badala ya kuingia tena:

```sh
ccam import
```

Folda hubaki pale ilipo, na `claude` ya kawaida inaendelea kufanya kazi kama awali.
