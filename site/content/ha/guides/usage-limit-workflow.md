---
title: Abin da za ka yi idan Claude Code ya kai iyakar amfaninsa
description: Duba sauran sararin da kowane asusun Claude Code yake da shi, kuma ka ci gaba da wannan tattaunawar a wani asusu idan ɗaya ya kai iyakarsa.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 767a731b912e
---

Idan wani asusu ya kai iyakarsa, ci gaba da wannan tattaunawar a wani asusun da ka mallaka: raba tarihin tattaunawa sau ɗaya, sannan ka gudanar da `ccam next --continue`. Idan kana da asusu ɗaya kawai, jira lokacin sake saitawa da Claude Code yake nunawa.

## Kafin hakan ta faru

Saita asusan ta yadda sauyawar za ta ɗauki umarni ɗaya kawai:

```sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

Yanzu layin matsayi (status line) yana nuna asusun da amfaninsa na awa 5 yayin da kake aiki, don haka kaiwa iyaka ba zai zo maka ba zato ba tsammani.

## Idan an kai iyakar

Fita daga Claude Code (`/exit`), sannan:

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22m       █████░░░  68% 1d 9h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
```

Buɗe tattaunawar a asusun da ya fi sarari:

```sh
ccam next --auto --continue
```

Claude Code yana farawa a matsayin `work` da tattaunawa ta baya-bayan nan ta foldar yanzu, har da duk abin da aka faɗa zuwa yanzu.

## Abubuwan da za ka tuna

- Gudanar da `ccam next --continue` a foldar da ka yi aiki a ciki; ana ajiye tattaunawa bisa foldar kowane project.
- Ka guji gudanar da asusai biyu a kan tattaunawa ɗaya a lokaci guda. ccam yana gargaɗe ka idan ka gwada.
- Dole ne ɗayan asusun ma ya raba tarihin tattaunawa, a rukuni ɗaya.
- Ko an yarda ka sauya tsakanin asusai ta wannan hanya ya dogara da sharuɗɗan kowane asusu.

Duba [iyakokin amfani](/docs/usage-limits) don kowane zaɓi na `ccam usage` da `ccam next`.
