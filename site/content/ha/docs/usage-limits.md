---
title: Iyakokin amfani da sauya asusai
description: Ga amfanin awa 5, na mako da na Fable na kowane asusun Claude Code a matsayin sanduna da ccam usage, sannan ka buɗe asusun da ya fi sarari da ccam next --auto.
updated: 2026-10-21
order: 5
translation: machine
sourceHash: 0e64e3251c69
---

`ccam usage` yana nuna nawa kowane asusu ya yi amfani da shi daga iyakar awa 5, iyakar mako da iyakar mako ta kowane model (kamar Fable). `ccam next --auto` yana buɗe asusun da ya fi sauran sarari, kuma `--continue` yana ɗauko tattaunawar yanzu zuwa can.

```console
$ ccam usage
 1 personal  just now
   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h
   Fable    ██████████████████░░  88%  resets in 1d 6h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
   Fable    █░░░░░░░░░░░░░░░░░░░   3%  resets in 5d 2h
 3 side      12m ago (cached)
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
   Fable    ████████░░░░░░░░░░░░  40%  resets in 1d 9h
→ Most headroom: work   (open it with ccam next --auto)
```

Kowane asusu yana samun layi mai nuna tsawon lokacin da lambobinsa suka yi, sannan sanda ɗaya ga kowace iyaka tare da lokacin da ya rage kafin a sake saita ta: taga ta awa 5, taga ta mako ga dukkan model, da taga ta mako ta kowane model, kamar Fable. Daga 80% sanda tana zama rawaya, a 100% kuma ja. Layin ƙarshe yana faɗin asusun da `ccam next --auto` zai buɗe.

## Daga ina lambobin suke fitowa

ccam yana tambayar endpoint na amfani na Anthropic game da kowane asusu a lokaci ɗaya, ta amfani da bayanan shigar kowane asusu. Yana jira daƙiƙa 3 a mafi yawa ga kowane asusu, kuma yana nuna kowane asusu da zarar ya amsa.

Ba a rubuta bayanin wannan endpoint a bainar jama'a ba, kuma yana iya canzawa. Idan buƙata ta gaza, ccam yana nuna lambobin ƙarshe da aka sani kuma yana yi musu alama a matsayin cache. Taga (window) da ya sake saita tun bayan da aka ajiye waɗannan lambobin yana nuna 0%. `ccam doctor` yana gaya maka idan tsarin amsar ya canza.

Ana sake amfani da lambobin da aka samo ƙasa da daƙiƙa 60 da suka wuce. Zaɓuɓɓuka:

| Zaɓi | Tasiri |
| --- | --- |
| `--watch [seconds]` | Sabuntawa koyaushe, kowane daƙiƙa 60 ta tsohuwa |
| `--json` | Sakamakon da na'ura za ta iya karantawa |
| `-v` | Nuna kuma inda lambobin kowane asusu suka fito |
| `--no-cache` | Koyaushe a sake tambaya |
| `--source cache` | Kada a taɓa tambaya; nuna lambobin da aka ajiye kawai |

## Bayanan shiga da suka ƙare

Idan access token na asusu ya ƙare, ccam ba ya sabunta shi, domin sabunta shi yayin da zaman Claude Code ke amfani da bayanan shiga iri ɗaya zai iya fitar da wannan zaman. Fara asusun sau ɗaya da `ccam <name>`, kuma Claude Code zai sabunta shi.

`ccam usage --refresh` yana sabunta token da suka ƙare da kansa, amma ga asusan da ba sa aiki kawai. Ba ya taɓa sabunta bayanan shiga da aka ajiye a Keychain na macOS ko a `~/.claude` da aka shigo da shi.

## Amfani a layin matsayi

```sh
ccam statusline install
```

Sannan Claude Code yana nuna sunan asusun da amfaninsa a layin matsayi (status line), misali `[work] · Opus · 5h 24% · 7d 41%`. Daga 90% zuwa sama, layin yana kuma ba da shawarar `ccam next --auto --continue`. Idan kana da umarnin layin matsayi tuni, ccam yana riƙe shi kuma yana saka sunan asusun a gabansa. Layin matsayi yana kuma ajiye lambobinsa, waɗanda `ccam usage` ke amfani da su idan endpoint ɗin bai amsa ba.

## Sauyawa

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # carry this folder's conversation on as the next account
```

`ccam next` yana tsallake asusan da ba a shiga da su ba, kuma `--auto` yana tsallake asusan da suka kai 100% da asusan wata ƙungiya (organization) dabam. Idan asusun biyu suna raba tarihin tattaunawa, `--continue` yana buɗe tattaunawa ɗaya; in ba haka ba Claude Code yana kwafar ta cikin sabon asusun. Duba [rabawa](/docs/sharing) da [sauyawa zuwa asusu na gaba](/features/next#continue).

Don rubutun umarni (scripts) da ayyukan da aka tsara lokacinsu masu gudanar da `claude -p`, saka `--fallback` kafin sunan asusun. Idan wannan asusun ya kai iyakarsa, ccam yana sake gudanar da umarnin iri ɗaya a matsayin asusun da ya fi sarari, kuma sakamakon gudun da ya kammala kaɗai ake bugawa:

```sh
ccam --fallback work -p "review the diff" --output-format json
```

Duba [Sake gudanarwa a iyaka](/features/next#run-again-at-the-limit) don yadda ccam yake yanke shawara.

[Jagorar matakan aiki idan an kai iyakar amfani](/guides/usage-limit-workflow) tana bin cikakken misali mataki-mataki.
