---
title: Iyakokin amfani da sauya asusai
description: Ga amfanin awa 5 da na mako na kowane asusun Claude Code da ccam usage, sannan ka buɗe asusun da ya fi sarari da ccam next --auto.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 8ab41501c548
---

`ccam usage` yana nuna nawa kowane asusu ya yi amfani da shi daga iyakar awa 5 da ta mako. `ccam next --auto` yana buɗe asusun da ya fi sauran sarari, kuma `--continue` yana ɗauko tattaunawar yanzu zuwa can.

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ██████░░  72% 1h 40m    ███░░░░░  41% 3d 6h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
 3 side      ████████ 100% 22m       █████░░░  68% 1d 9h     12m ago (cached)
→ Most headroom: work   (open it with ccam next --auto)
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

Idan access token na asusu ya ƙare, ccam ba ya sabunta shi, domin sabunta shi yayin da zaman Claude Code ke amfani da bayanan shiga iri ɗaya zai iya fitar da wannan zaman. Fara asusun sau ɗaya da `ccam <name>`, kuma Claude Code zai sabunta shi.

`ccam usage --refresh` yana sabunta token da suka ƙare da kansa, amma ga asusan da ba sa aiki kawai. Ba ya taɓa sabunta bayanan shiga da aka ajiye a Keychain na macOS ko a `~/.claude` da aka shigo da shi.

## Amfani a layin matsayi

```sh
ccam statusline install
```

Sannan Claude Code yana nuna sunan asusun da amfaninsa a layin matsayi (status line), misali `[work] · Opus · 5h 24% · 7d 41%`. Idan kana da umarnin layin matsayi tuni, ccam yana riƙe shi kuma yana saka sunan asusun a gabansa. Layin matsayi yana kuma ajiye lambobinsa, waɗanda `ccam usage` ke amfani da su idan endpoint ɗin bai amsa ba.

## Sauyawa

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # the next account that shares your conversation history
```

`ccam next` yana tsallake asusan da ba a shiga da su ba, kuma `--auto` yana tsallake asusan da suka kai 100%. Don ci gaba da tattaunawa, dole ne asusun biyu su raba tarihin tattaunawa; duba [rabawa](/docs/sharing).

[Jagorar matakan aiki idan an kai iyakar amfani](/guides/usage-limit-workflow) tana bin cikakken misali mataki-mataki.
