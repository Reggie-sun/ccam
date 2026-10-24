---
title: Claude Code ya kai iyakar amfani? Ga abin da za ka yi
description: Jira iyakar ta sake farawa, ko ka ci gaba da tattaunawar a wani asusun Claude Code nan take. Yadda iyakar awa 5 da ta mako ke aiki, mataki-mataki.
updated: 2026-10-24
order: 3
translation: machine
sourceHash: 049a46684c7b
---

Idan wani asusu ya kai iyakar amfaninsa (wanda ake yawan kira rate limit), kana da zaɓi biyu: jira har iyakar ta sake farawa, ko ka ci gaba da wannan tattaunawar yanzu a wani asusun da ka mallaka. Da ccam, zaɓi na biyu yana ɗaukar umarni ɗaya, `ccam next --auto --continue`.

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
ccam group set default --add history
ccam statusline install
```

Yanzu layin matsayi (status line) yana nuna asusun da amfaninsa na awa 5 yayin da kake aiki, don haka kaiwa iyaka ba zai zo maka ba zato ba tsammani.

Idan an kai iyakar, fita daga Claude Code (`/exit`), sannan:

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

Buɗe tattaunawar a asusun da ya fi sarari:

```sh
ccam next --auto --continue
```

Claude Code yana farawa a matsayin `work` da tattaunawa ta baya-bayan nan ta foldar yanzu, har da duk abin da aka faɗa zuwa yanzu. Raba tarihi (`ccam group set default --add history` a sama) yana sa ta zama tattaunawa ɗaya; ba tare da shi ba, Claude Code yana kwafar tattaunawar cikin `work`, kuma ccam yana gaya maka tokens nawa amsar farko za ta sake karantawa.

Da [plugin ɗin ccam](/docs/install#claude-code-plugin) a Claude Code, zaman da ya tsaya a iyakarsa yana nuna sanarwa kuma yana rubuta cikakkiyar tagar, don kada `ccam next --auto` ya sake zaɓar wannan asusun kafin sake saitinsa, kuma `/ccam:switch` (ko `/ccam:switch work`) yana rubuta bayanin tattaunawar don `ccam next --continue` ba tare da barin wani abu da za a yi hasashe ba.

**Kujeru (seats) na Team da Enterprise.** Zaɓin atomatik (`ccam next`, `--fallback`, `ccam pick`, maɓallin `n` na dashboard) yana kuma sauyawa tsakanin kujerun ƙungiyar Claude Team ko Enterprise iri ɗaya. Cibiyar taimako ta Anthropic ta ce "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (ana ba da kowace kujera a tsarin Team da sunan takamaiman mutum ɗaya, kuma wannan mutumin kaɗai ne aka yarda ya yi amfani da ita) da "Users can't route their work through another user's account to use that user's remaining usage" (masu amfani ba za su iya tura aikinsu ta asusun wani mai amfani ba domin su ci gajiyar ragowar amfanin wannan mai amfanin) ([Tsare-tsaren Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); haka ma ga [tsare-tsaren Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Don haka yin amfani da zaɓin atomatik tsakanin kujeru na iya saɓa wa sharuɗɗan Anthropic. Ka tuntuɓi ƙungiyarka da Anthropic kafin ka dogara da shi. Don tsare-tsaren kai (Pro, Max), duba [bayanin sharuɗɗa a cikin Tambayoyin da ake yawan yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## Shin ccam yana sauya asusu ta atomatik?

A'a. ccam ba ya aiki yayin da kake aiki: yana fara Claude Code kuma ya maye gurbin kansa da shi. Kai ne kake yanke shawarar lokacin sauyawa, kuma `ccam next --auto` yana zaɓar asusun da amfaninsa na awa 5 ya fi ƙasa, yana tsallake asusan da suka kai 100%, waɗanda ba a shiga ba, ko na wata ƙungiya (organization) dabam. Da `relay = true` a `config.toml`, ccam yana jira a gefen zaman maimakon haka, don haka `/ccam:switch` yana ci gaba da kansa, yana riƙe zaɓuɓɓukan zaman, yanayin izini (permission mode) da matakin ƙoƙari (effort level). Wasu kayan aikin kuma suna musanya bayanan shiga na zaman da ke gudana maimakon haka; [ccam da sauran masu sauya asusu](/compare/account-switchers) yana bayyana bambancin.

## Abubuwan da za ka tuna

- Gudanar da `ccam next --continue` a foldar da ka yi aiki a ciki; ana ajiye tattaunawa bisa foldar kowane project.
- Ka guji gudanar da asusai biyu a kan tattaunawa ɗaya a lokaci guda. ccam yana gargaɗe ka idan ka gwada.
- Ba tare da tarihin da aka raba ba, ana kwafar tattaunawar zuwa ɗayan asusun kuma kwafin biyu suna bin hanyoyinsu daban; checkpoints da `/rewind` ke maido da lamba daga gare su suna zama tare da asusun farko.
- Saƙon farko a sabon asusun zai iya cinye iyakarsa fiye da yadda aka saba, domin wannan asusun bai riga ya ajiye tattaunawar a cache ba.
- Ko an yarda ka sauya tsakanin asusai ta wannan hanya ya dogara da sharuɗɗan kowane asusu.

Duba [iyakokin amfani](/docs/usage-limits) don kowane zaɓi na `ccam usage` da `ccam next`.
