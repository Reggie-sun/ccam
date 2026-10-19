---
title: Farawa da sauri
description: Ƙara asusun Claude Code biyu da ccam add, ba su gajerun umarni kamar ccwork, sauya tsakaninsu, kuma zaɓi abin da za su raba.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: d040e34d2a17
---

Ƙara kowane asusu da `ccam add`, sannan ka fara Claude Code da `ccam <name>` ko da gajeren umarni da kai da kanka ka sa masa suna, kamar `ccwork`. Wannan ne kawai abin da yawancin mutane suke buƙata.

## 1. Ƙara asusunka na farko

```sh
ccam add
```

Mataimakin saitawa (wizard) yana tambayar suna kamar `personal` ko `work`, abubuwan da asusun zai raba da sauran asusunka, da kuma ko a ƙirƙiri gajeren umarni (duba mataki na 4). Sannan yana gudanar da `claude auth login`, wanda ke buɗe burauza don ka shiga. A na'ura mai nisa, burauza na iya ƙarewa a shafin `http://localhost:…/callback` wanda ba ya buɗewa; liƙa wannan adireshin a terminal kuma ccam zai kammala shigarwar. Idan ya gama, za ka ga:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Kawo saitunanka da ke akwai

Idan ka taɓa amfani da Claude Code a baya, saitunanka, `CLAUDE.md`, extensions da sabobin MCP suna cikin `~/.claude`. Kwafe su zuwa saitunan da aka raba, don kowane asusu ya same su:

```sh
ccam migrate
```

ccam yana nuna abin da zai kwafa kuma yana tambaya tukuna. Yana karanta `~/.claude` ne kawai, ba ya taɓa kwafar bayanan shiga, kuma ba ya haɗa tarihin tattaunawa sai ka ƙara `--history`. Duba [kwafar saitin da ke akwai](/docs/sharing#copy-an-existing-setup) don cikakken bayani.

## 3. Ƙara asusu na biyu

```sh
ccam add work
```

Shiga da ɗayan asusun a burauza.

## 4. Sauya asusu

`ccam work` yana fara Claude Code a matsayin asusun work. Don wani abu mafi gajarta, ba asusun nasa gajeren umarni, kamar `ccwork`. Bayan abubuwan da za a raba, mataimakin saitawa yana tambaya:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Danna `y`, sannan ka rubuta suna, ko ka danna Enter don amfani da wanda aka ba da shawara:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Bayan shiga, ccam yana tabbatar da shi:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

Idan ka tsallake tambayar, za ka iya ƙara gajeren umarni a kowane lokaci:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Yanzu kowane ɗaya daga cikin waɗannan yana fara Claude Code a matsayin asusun work:

```sh
ccwork
ccam work
ccam 2
```

Gudanar da `ccam` shi kaɗai yana buɗe mai zaɓi: yi amfani da maɓallan kibiya da Enter, ko danna lambar asusun.

Duk abin da ke bayan gajeren umarnin ko sunan asusun yana zuwa Claude Code ba tare da canji ba, don haka za ka ci gaba da amfani da zaɓuɓɓukan da ka sani:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Zaɓi abin da za a raba

Ta tsohuwa, ana raba sabobin MCP, `settings.json`, `CLAUDE.md` na gaba ɗaya (global), da agents, commands, skills da plugins tsakanin asusai. Ba a raba tarihin tattaunawa ba. Kunna shi ga asusun biyu idan kana son ci gaba da tattaunawa a tsakaninsu:

```sh
ccam share personal --history on
ccam share work --history on
```

ccam yana nuna abin da zai motsa kafin ya canza komai. Duba [rabawa](/docs/sharing) don cikakken bayani.

## 6. Idan asusu ya kai iyakarsa

```sh
ccam usage
ccam next --continue
```

`ccam usage` yana nuna amfanin awa 5 da na mako na kowane asusu. `ccam next --continue` yana buɗe wannan tattaunawar a asusu na gaba. Duba [iyakokin amfani](/docs/usage-limits).

## Kana amfani da ~/.claude tuni?

Yi masa rajista maimakon sake shiga:

```sh
ccam import
```

Foldar tana nan a inda take, kuma `claude` kai tsaye yana ci gaba da aiki kamar da.
