---
title: Iyakokin amfani da sauya asusai
description: Ga amfanin awa 5 da na mako na kowane asusun Claude Code da ccs usage, sannan ka buɗe asusun da ya fi sarari da ccs next --auto.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: a9ee77936668
---

`ccs usage` yana nuna nawa kowane asusu ya yi amfani da shi daga iyakar awa 5 da ta mako. `ccs next --auto` yana buɗe asusun da ya fi sauran sarari, kuma `--continue` yana ɗauko tattaunawar yanzu zuwa can.

```console
$ ccs usage
 계정        5시간         초기화      주간          초기화      기준
 1 personal  ██████░░  72% 1시간 40분  ███░░░░░  41% 3일 6시간   방금
 2 work      █░░░░░░░   8% 4시간 12분  ██░░░░░░  19% 5일 2시간   방금
 3 side      ████████ 100% 22분        █████░░░  68% 1일 9시간   12분 전 (캐시)
→ 여유가 가장 많은 계정: work   (ccs next --auto 로 바로 열기)
```

Ginshiƙan su ne: asusu, amfanin awa 5, lokacin da ya rage kafin a sake saita shi, amfanin mako, lokacin da ya rage kafin a sake saita shi, da kuma tsawon lokacin da lambobin suka yi.

## Daga ina lambobin suke fitowa

ccam yana tambayar endpoint na amfani na Anthropic game da kowane asusu a lokaci ɗaya, ta amfani da bayanan shigar kowane asusu. Yana jira daƙiƙa 3 a mafi yawa ga kowane asusu, kuma yana nuna kowane asusu da zarar ya amsa.

Ba a rubuta bayanin wannan endpoint a bainar jama'a ba, kuma yana iya canzawa. Idan buƙata ta gaza, ccam yana nuna lambobin ƙarshe da aka sani kuma yana yi musu alama a matsayin cache. `ccam doctor` yana gaya maka idan tsarin amsar ya canza.

Ana sake amfani da lambobin da aka samo ƙasa da daƙiƙa 60 da suka wuce. Zaɓuɓɓuka:

| Zaɓi | Tasiri |
| --- | --- |
| `--watch [seconds]` | Sabuntawa koyaushe, kowane daƙiƙa 60 ta tsohuwa |
| `--json` | Sakamakon da na'ura za ta iya karantawa |
| `-v` | Nuna kuma iyakokin mako na kowane model da inda kowace lamba ta fito |
| `--no-cache` | Koyaushe a sake tambaya |
| `--source cache` | Kada a taɓa tambaya; nuna lambobin da aka ajiye kawai |

## Bayanan shiga da suka ƙare

Idan access token na asusu ya ƙare, ccam ba ya sabunta shi, domin sabunta shi yayin da zaman Claude Code ke amfani da bayanan shiga iri ɗaya zai iya fitar da wannan zaman. Fara asusun sau ɗaya da `ccs <name>`, kuma Claude Code zai sabunta shi.

`ccs usage --refresh` yana sabunta token da suka ƙare da kansa, amma ga asusan da ba sa aiki kawai. Ba ya taɓa sabunta bayanan shiga da aka ajiye a Keychain na macOS ko a `~/.claude` da aka shigo da shi.

## Amfani a layin matsayi

```sh
ccam statusline install
```

Sannan Claude Code yana nuna sunan asusun da amfaninsa a layin matsayi (status line), misali `[work] · Opus · 5h 24% · 7d 41%`. Idan kana da umarnin layin matsayi tuni, ccam yana riƙe shi kuma yana saka sunan asusun a gabansa. Layin matsayi yana kuma ajiye lambobinsa, waɗanda `ccs usage` ke amfani da su idan endpoint ɗin bai amsa ba.

## Sauyawa

```sh
ccs next              # the account after the one you used last
ccs next --auto       # the account with the lowest 5-hour usage
ccs next --continue   # the next account that shares your conversation history
```

`ccs next` yana tsallake asusan da ba a shiga da su ba, kuma `--auto` yana tsallake asusan da suka kai 100%. Don ci gaba da tattaunawa, dole ne asusun biyu su raba tarihin tattaunawa; duba [rabawa](/docs/sharing).

[Jagorar matakan aiki idan an kai iyakar amfani](/guides/usage-limit-workflow) tana bin cikakken misali mataki-mataki.
