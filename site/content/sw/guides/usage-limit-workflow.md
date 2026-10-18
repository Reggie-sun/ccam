---
title: Cha kufanya Claude Code inapofikia kikomo chake cha matumizi
description: Kagua kila akaunti ya Claude Code ina nafasi kiasi gani na uendeleze mazungumzo yale yale kwenye akaunti nyingine pale moja inapofikia kikomo chake.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 64fe0b561bd1
---

Akaunti moja inapofikia kikomo chake cha matumizi (usage limit), endeleza mazungumzo yale yale katika akaunti nyingine unayoimiliki: shiriki historia ya mazungumzo mara moja, kisha endesha `ccs next --continue`. Ikiwa una akaunti moja tu, subiri muda wa kuanza upya ambao Claude Code inaonyesha.

## Kabla halijatokea

Sanidi akaunti ili kubadili kuchukue amri moja:

```sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

Sasa mstari wa hali (status line) huonyesha akaunti na matumizi yake ya saa 5 wakati unafanya kazi, kwa hiyo kikomo hakikushtui.

## Kikomo kinapofikiwa

Toka kwenye Claude Code (`/exit`), kisha:

```sh
ccs usage
```

```console
 1 personal  ████████ 100% 22분        █████░░░  68% 1일 9시간   방금
 2 work      █░░░░░░░   8% 4시간 12분  ██░░░░░░  19% 5일 2시간   방금
```

Fungua mazungumzo katika akaunti yenye nafasi kubwa zaidi:

```sh
ccs next --auto --continue
```

Claude Code huanza kama `work` ikiwa na mazungumzo ya hivi karibuni zaidi ya folda ya sasa, pamoja na kila kitu kilichosemwa hadi sasa.

## Mambo ya kukumbuka

- Endesha `ccs next --continue` katika folda ile ile uliyofanyia kazi; mazungumzo huhifadhiwa kwa kila folda ya mradi.
- Epuka kuendesha akaunti mbili kwenye mazungumzo yale yale kwa wakati mmoja. ccam hukuonya ukijaribu.
- Akaunti nyingine lazima pia ishiriki historia ya mazungumzo, katika kikundi kile kile.
- Iwapo unaruhusiwa kubadili kati ya akaunti kwa namna hii inategemea masharti ya kila akaunti.

Tazama [vikomo vya matumizi](/docs/usage-limits) kwa kila chaguo la `ccs usage` na `ccs next`.
