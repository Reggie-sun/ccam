---
title: Dashboard
description: ccam a terminal yana buɗe dashboard mai cike allo na asusai, rukunoni, ƙa'idodin folda, extensions da dubawa; kowane canji yana fara nuna shirinsa.
updated: 2026-10-20
order: 15
translation: machine
sourceHash: "decdd5fa54c9"
---

`ccam` ba tare da arguments ba yana buɗe dashboard mai cike allo idan standard input da output terminal ne; in ba haka ba yana buga `ccam list`. Dashboard yana nuna kowane asusu tare da amfaninsa da rukunoninsa, kuma yana yin abin da umarnan ke yi: farawa, ƙarawa, shiga, canza rukunoni, saita ƙa'idodin folda, dubawa da gyarawa. Farawa yana rufe shi; ƙara asusu ko shiga yana barin sa zuwa terminal sannan ya dawo.

```console
 ccam v0.9.0 · 3 accounts · 2 groups         ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health

    # Account     Email             5-hour         7-day Last  Groups
 ❯ 1. personal ✔   me@example.com    █████▊░░  72%   41% 2h    default
   2. work      ▸• dev@company.com   ██░░░░░░  25%  100% now   default team
   3. side         -                 ✘ not logged in         -     no group

 ✔ personal is the default account
 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

Firam na farko yana zana daga `config.toml` kaɗai, cikin ƙasa da 100 ms sosai; amfani yana cika yayin da kowane asusu ya amsa. Kashi yana zama rawaya daga 80% kuma ja a 100%, kuma ƙimomin da suka fito daga cache bayan neman da ya gaza ana zana su a dusashe.

## Shafuka

| Shafi | Abin da yake nunawa | Maɓallai |
| --- | --- | --- |
| Accounts | lamba, suna, imel, sandar awa 5, amfanin kwana 7, amfani na ƙarshe, rukunoni; `✔` tsoho, `▸` asusun wannan foldar, `•` mai aiki | a ƙasa |
| Groups | kowane rukuni, abin da yake rabawa da asusansa, da asusan da ba sa cikin kowane rukuni | `enter` gyara, `n` sabo, `e` sake suna, `x` cire |
| Folders | ƙa'idodin [`ccam folder`](/features/folder) da wanda ke aiki a nan | `n` yi amfani da asusu a nan, `enter` canza, `x` cire |
| Extensions | abin da [`ccam ext`](/features/ext) ke nunawa, kowane extension sau ɗaya tare da rukunonin da yake fitowa daga gare su | `r` sake lodawa |
| Health | abin da [`ccam doctor`](/features/doctor) ke bayar da rahoto | `f` gyara (`doctor --fix`) |

`←`/`→` ko `Tab` suna sauya shafuka, `?` yana lissafa maɓallan shafin yanzu, `r` yana sake neman amfani ba tare da gajeren cache ba, kuma `q` ko `Esc` yana fita.

## Maɓallai a shafin Accounts

| Maɓalli | Abin da yake yi |
| --- | --- |
| `Enter` | Fara Claude Code a matsayin asusun da aka zaɓa |
| `1`–`9` | Je kan asusun da ke da wannan lambar; rubuta `1` `2` da sauri don 12. `Enter` ne kaɗai ke farawa |
| `↑`/`↓`, `k`/`j` | Matsawa |
| `n` | Fara asusun da [`ccam next --auto`](/features/next) zai zaɓa |
| `o` | Fara da zaɓuɓɓuka: `--continue`, `--resume`, wasu arguments, ko dry run |
| `u` | Cikakken bayanin amfani: tagogin awa 5, kwana 7 da na model tare da lokutan sake saitinsu da inda lambobin suka fito |
| `a` | Ƙara asusu (mataimakin da ke ƙasa) |
| `l` | Sake shiga |
| `s` | Zaɓi rukunonin asusun |
| `d` | Mai da shi asusun tsoho |
| `f` | Yi amfani da shi a foldar yanzu |
| `e` | Sake masa suna |
| `c` | Saita ko cire gajeren umarninsa |
| `x` | Cire shi (yana fara tambaya, yana farawa a kan No) |

## Canza rukunoni

A shafin Groups, `Enter` yana buɗe mai gyarawa: abubuwan a hagu, asusan a dama, `Tab` don sauyawa, `Space` don kunnawa ko kashewa. Abin da ɗaya daga cikin asusan ke samu tuni daga wani rukuni ana dusashe shi tare da dalilin, haka ma asusun da zai sami abu sau biyu; abu zai iya fitowa daga rukuni ɗaya kawai. `Enter` yana nuna shirin:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Ana fara tambaya game da fayil da ya bambanta da kwafin rukunin (riƙe na rukunin, yi amfani da nawa, ko haɗa su). Sannan canjin yana gudana a bayan fage da lamba iri ɗaya da ta [`ccam group`](/features/share): ana fara ajiye saitin, bayanan kowane asusu suna motsawa a tsari, kuma fayilolin asali suna zuwa `.ccam-backup/`. Ana ƙin canjin da zai matsar da bayanan asusun da ke aiki, tare da sunan asusun.

## Ƙara asusu

`a` yana buɗe mataimakin da `ccam add` ma ke amfani da shi a terminal:

1. Sunan asusun.
2. Yadda yake shiga: biyan kuɗin Claude yanzu (a burauza), [API key](/features/api), ko biyan kuɗin Claude daga baya.
3. Ga API key: inda ccam ke karanta key ɗin a kowane farawa (`env:NAME` ko `command:COMMAND`), endpoint, yadda ake aika key ɗin da model.
4. Rukunonin da za a shiga, tare da alama a kan `default`; kafin a sami rukunoni, abin da rukunin `default` ke rabawa.
5. Gajeren umarni, idan kana son ɗaya.
6. Taƙaitawa don tabbatarwa.

`Esc` yana komawa mataki ɗaya baya, kuma a mataki na farko yana rufe mataimakin.

## Launuka

Launukan suna bin matsayin launuka na UI na terminal na Claude Code, don haka ccam yana dacewa a gefensa; sunayen rukunoni suna samun launi daga sunansu. ccam yana tambayar terminal bangonsa kuma yana zaɓar jigo mai haske ko mai duhu. Saita ɗaya da kanka a `config.toml`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Jigogin `-daltonized` na makantar launin ja da kore ne, kuma jigogin `-ansi` suna amfani da launuka 16 na terminal kansa. `NO_COLOR` yana kashe launuka.

## Tsaro da yanayi na musamman

- Dashboard yana kiran ayyuka (functions) iri ɗaya da umarnan, don haka kowace ƙa'ida tana aiki: shiri kafin bayanai su motsa, ƙin canza asusan da ke aiki, kwafin ajiya, da tokens na shiga da ba sa taɓa barin asusu.
- Dashboard ba shi da dry run: `ccam --dry-run` ba tare da suna ba yana buga yadda asusun wannan foldar zai fara, kuma `o` a kan asusu yana ba da dry run na wannan farawa ɗaya kawai.
- Ana maido da terminal kafin Claude Code ya fara ko shiga ta buɗe.
- Allon da ya yi ƙanƙanta ga akwatin tattaunawa (dialog) yana fara cire layukan shafin; ana yanke layuka da `…` maimakon naɗa su.

## Masu alaƙa

- [Rabawa da rukunoni](/features/share)
- [Fara Claude Code](/features/launch)
- [Sarrafa asusai](/features/accounts)
- [Iyakokin amfani](/features/usage)
