---
title: Farawa da sauri
description: Ƙara asusun Claude Code biyu da ccam add, sauya tsakaninsu da claude1, claude2 ko ccam <name>, sannan ka zaɓi abin da za su raba.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 6baf0bfa5029
---

Ƙara kowane asusu da `ccam add`, sannan ka fara Claude Code da `claude1`, `claude2` ko `ccam <name>`. Wannan ne kawai abin da yawancin mutane suke buƙata.

## 1. Ƙara asusunka na farko

```sh
ccam add
```

Mataimakin saitawa (wizard) yana tambayar suna kamar `personal` ko `work`, da kuma abubuwan da asusun zai raba da sauran asusunka. Sannan yana gudanar da `claude auth login`, wanda ke buɗe burauza don ka shiga. Idan ya gama, za ka ga:

```console
✓ personal (me@example.com) added — shortcut: claude1
```

ccam yana buga saƙonninsa da harshen Turanci; umarni da zaɓuɓɓuka iri ɗaya ne a kowane harshe.

## 2. Ƙara asusu na biyu

```sh
ccam add work
```

Shiga da ɗayan asusun a burauza. Zai zama `claude2`.

## 3. Sauyawa

Kowane ɗaya daga cikin waɗannan yana fara Claude Code a matsayin asusun work:

```sh
claude2
ccam work
ccam 2
```

Duk abin da ke bayan sunan asusun yana zuwa Claude Code ba tare da canji ba, don haka za ka ci gaba da amfani da zaɓuɓɓukan da ka sani:

```sh
ccam work --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

Gudanar da `ccam` shi kaɗai yana buɗe mai zaɓi: yi amfani da maɓallan kibiya da Enter, ko danna lambar asusun.

## 4. Zaɓi abin da za a raba

Ta tsohuwa, ana raba sabobin MCP, `settings.json`, `CLAUDE.md` na gaba ɗaya (global), da agents, commands, skills da plugins tsakanin asusai. Ba a raba tarihin tattaunawa ba. Kunna shi ga asusun biyu idan kana son ci gaba da tattaunawa a tsakaninsu:

```sh
ccam share personal --history on
ccam share work --history on
```

ccam yana nuna abin da zai motsa kafin ya canza komai. Duba [rabawa](/docs/sharing) don cikakken bayani.

## 5. Idan asusu ya kai iyakarsa

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
