---
title: Claude Code ya kai iyakar amfani? Ga abin da za ka yi
description: Jira iyakar ta sake farawa, ko ka ci gaba da tattaunawar a wani asusun Claude Code nan take. Yadda iyakar awa 5 da ta mako ke aiki, mataki-mataki.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: ff765a189332
---

Idan wani asusu ya kai iyakar amfaninsa (wanda ake yawan kira rate limit), kana da zaɓi biyu: jira har iyakar ta sake farawa, ko ka ci gaba da wannan tattaunawar yanzu a wani asusun da ka mallaka. Da ccam, zaɓi na biyu yana ɗaukar umarni ɗaya, `ccam next --auto --continue`, da zarar an raba tarihin tattaunawa tsakanin asusan.

## Waɗanne iyakoki Claude Code yake da su?

A tsare-tsaren Pro da Max, Claude Code yana da [iyakar zama ta awa biyar da iyakar mako](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Duka biyun ana raba su da Claude a yanar gizo da kuma a manhajoji, kuma amfani yana ƙirguwa a kan duka biyun a lokaci guda. Idan ɗaya ya ƙare, Claude Code [yana nuna saƙo](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) kamar "You've hit your session limit · resets 3:45pm". Sauya model ba ya taimaka, domin iyakokin sun shafi duk model.

## Zaɓi na 1: jira iyakar ta sake farawa

Idan kana da asusu ɗaya, ka jira. Sababbin nau'ikan Claude Code [suna jira a cikin zaman da ke buɗe](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) kuma suna ci gaba da aikin da kansu da zarar iyakar ta sake farawa. `/usage` a cikin Claude Code yana nuna nawa ka yi amfani da shi daga iyakokin tsarinka.

## Zaɓi na 2: ci gaba a wani asusu yanzu

Idan kana da asusu na biyu, za ka iya ci gaba nan take da wannan tattaunawar. Saita shi sau ɗaya:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

Yanzu layin matsayi (status line) yana nuna asusun da amfaninsa na awa 5 yayin da kake aiki, don haka kaiwa iyaka ba zai zo maka ba zato ba tsammani.

Idan an kai iyakar, fita daga Claude Code (`/exit`), sannan:

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

## Shin ccam yana sauya asusu ta atomatik?

A'a. ccam ba ya aiki yayin da kake aiki: yana fara Claude Code kuma ya maye gurbin kansa da shi. Kai ne kake yanke shawarar lokacin sauyawa, kuma `ccam next --auto` yana zaɓar asusun da amfaninsa na awa 5 ya fi ƙasa, yana tsallake asusan da suka kai 100% ko waɗanda ba a shiga ba. Wasu kayan aikin kuma suna musanya bayanan shiga na zaman da ke gudana maimakon haka; [ccam da sauran masu sauya asusu](/compare/account-switchers) yana bayyana bambancin.

## Abubuwan da za ka tuna

- Gudanar da `ccam next --continue` a foldar da ka yi aiki a ciki; ana ajiye tattaunawa bisa foldar kowane project.
- Ka guji gudanar da asusai biyu a kan tattaunawa ɗaya a lokaci guda. ccam yana gargaɗe ka idan ka gwada.
- Dole ne ɗayan asusun ma ya raba tarihin tattaunawa, a rukuni ɗaya.
- Saƙon farko a sabon asusun zai iya cinye iyakarsa fiye da yadda aka saba, domin wannan asusun bai riga ya ajiye tattaunawar a cache ba.
- Ko an yarda ka sauya tsakanin asusai ta wannan hanya ya dogara da sharuɗɗan kowane asusu.

Duba [iyakokin amfani](/docs/usage-limits) don kowane zaɓi na `ccam usage` da `ccam next`.
