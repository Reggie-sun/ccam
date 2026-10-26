---
title: Farawa da sauri
metaTitle: "Farawa da sauri: ƙara da sauya asusun Claude Code"
description: Ƙara asusun Claude Code biyu da ccam add, ba su gajerun umarni kamar ccwork, sauya tsakaninsu, kuma zaɓi abin da za su raba.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: bd893bbb83ab
---

Ƙara kowane asusu da `ccam add`, sannan ka fara Claude Code da `ccam <name>` ko da gajeren umarni da kai da kanka ka sa masa suna, kamar `ccwork`. Wannan ne kawai abin da yawancin mutane suke buƙata.

Kana amfani da Claude Code tuni? Karanta [Kana amfani da ~/.claude tuni?](#already-using-claude) a ƙarshen wannan shafin kafin ka fara.

## 1. Ƙara asusunka na farko

```sh
ccam add
```

Mataimakin saitawa (wizard) yana tambayar suna kamar `personal` ko `work`, abubuwan da asusun zai raba da sauran asusunka, da kuma ko a ƙirƙiri gajeren umarni (duba mataki na 3). Sannan yana gudanar da `claude auth login`, wanda ke buɗe burauza don ka shiga. A na'ura mai nisa, burauzarka na iya ƙarewa a shafin `http://localhost:…/callback` wanda ba ya buɗewa; liƙa wannan adireshin a terminal kuma ccam zai kammala shigarwar. Idan shigarwar ta kammala, za ka ga:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Kawo saitunanka da ke akwai

Idan ka taɓa amfani da Claude Code a baya, saitunanka, `CLAUDE.md` da extensions suna cikin `~/.claude`, kuma sabobin MCP naka suna cikin `~/.claude.json`. Kwafe su zuwa rukunin `default` da aka raba don kowane asusu ya same su:

```sh
ccam migrate
```

ccam yana nuna abin da zai kwafa kuma yana tambaya tukuna. Yana karanta `~/.claude` ne kawai, ba ya taɓa kwafar bayanan shiga, kuma ba ya haɗa tarihin tattaunawa sai ka ƙara `--history`. Duba [kwafar saitin da ke akwai](/docs/sharing#copy-an-existing-setup) don cikakken bayani.

## 3. Ƙara asusu na biyu

```sh
ccam add work
```

Mataimakin saitawa yana tsallake sunan, wanda ka bayar a layin umarni, kuma yana nuna rukunonin da za a shiga tare da `default` an yi masa alama. A wannan karon, ba asusun nasa gajeren umarni, kamar `ccwork`, don ka iya fara shi da kalma ɗaya. Bayan rukunonin, mataimakin saitawa yana tambaya:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Danna `y`, sannan ka rubuta suna, ko ka danna Enter don amfani da wanda aka ba da shawara:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Sannan ka shiga da ɗayan asusun a burauza. ccam yana tabbatar da asusun da gajeren umarninsa:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

## 4. Sauya asusu

`ccam work` yana fara Claude Code a matsayin asusun work. Idan ka tsallake tambayar gajeren umarni, za ka iya ƙara gajeren umarni a kowane lokaci:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Yanzu kowane ɗaya daga cikin waɗannan umarnai yana fara Claude Code a matsayin asusun work (`2` lambarsa ce a `ccam list`):

```sh
ccwork
ccam work
ccam 2
```

Gudanar da `ccam` shi kaɗai yana buɗe [dashboard](/features/dashboard): je kan asusu da maɓallan kibiya ko lambarsa, sannan danna Enter.

Duk abin da ke bayan gajeren umarnin ko sunan asusun yana zuwa Claude Code ba tare da canji ba, don haka za ka ci gaba da amfani da zaɓuɓɓukan da ka sani:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Zaɓi abin da za a raba

Ta tsohuwa, ana raba sabobin MCP, `settings.json`, `CLAUDE.md` na gaba ɗaya (global) da extensions (agents, commands, skills da plugins) tsakanin asusai ta rukunin `default`. Ba a raba tarihin tattaunawa ba. Ƙara shi ga rukunin idan kana son ci gaba da tattaunawa tsakanin asusai:

```sh
ccam group set default --add history
```

ccam yana nuna abin da zai motsa ga kowane asusu kafin ya canza komai. Duba [rabawa da rukunoni](/features/share) don cikakken bayani.

## 6. Idan asusu ya kai iyakarsa

```sh
ccam usage
ccam next --continue
```

`ccam usage` yana nuna amfanin awa 5 da na mako na kowane asusu. `ccam next --continue` yana buɗe wannan tattaunawar a asusu na gaba. A cikin rubutun umarni (scripts), `ccam --fallback work -p "…"` yana sauyawa da kansa. Duba [iyakokin amfani](/docs/usage-limits).

**Kujeru (seats) na Team da Enterprise.** Zaɓin asusu na atomatik (`ccam next`, `--fallback`, `ccam pick`, maɓallin `n` na dashboard kuma, da zarar ka kunna shi, [`auto_switch`](/features/next#auto-switch)) yana kuma sauyawa tsakanin kujerun ƙungiyar Claude Team ko Enterprise iri ɗaya, wanda na iya saɓa wa sharuɗɗan Anthropic. Ka karanta [Tambayoyin da ake yawan yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) kuma ka tuntuɓi ƙungiyarka da Anthropic kafin ka dogara da shi.

## Kana amfani da ~/.claude tuni?

Yi masa rajista maimakon sake shiga:

```sh
ccam import
```

Foldar tana nan a inda take, kuma `claude` kai tsaye yana ci gaba da aiki kamar da.
