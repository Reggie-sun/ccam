---
title: Gabatar da ccam
description: Dalilin samuwar ccam, abin da yake yi wa masu asusun Claude Code da dama, da yadda yake guje wa kawo cikas.
updated: 2026-10-26
order: 1
translation: machine
sourceHash: e4370f29b818
---

ccam ƙaramin kayan aikin command-line ne ga mutanen da ke amfani da asusun Claude Code fiye da ɗaya a kwamfuta ɗaya: asusu na kai da na aiki, ko asusai da dama don dogon project ɗaya.

## Matsalar

Claude Code yana tallafa wa asusai da dama ta [`CLAUDE_CONFIG_DIR`](/guides/claude-config-dir): folda ɗaya ga kowane asusu. A aikace, kowace folda kamar tsibiri take. Dole ne a sake saita saituna, `CLAUDE.md`, sabobin MCP, agents da skills a kowace folda, kuma a riƙa daidaita su da hannu. Tattaunawa suna zama inda aka fara su. Kuma idan wani asusu ya kai iyakar amfaninsa, sai ka gwada sannan za ka gane.

## Abin da ccam ke yi

- `ccam add` yana ƙirƙirar foldar, yana gudanar da shiga kuma, idan kana so, yana ba asusun gajeren umarni da kai da kanka ka sa masa suna, kamar `ccwork`.
- `ccam work` da `ccwork` suna fara Claude Code a matsayin wannan asusun kuma suna miƙa kowane argument ba tare da canji ba.
- Ana raba sabobin MCP, saituna, `CLAUDE.md` da extensions ta mahaɗai, abu bayan abu da asusu bayan asusu. Ba a taɓa raba bayanan shiga ba.
- `ccam usage` yana nuna amfanin awa 5 da na mako na kowane asusu; `ccam next --continue` yana ɗauko tattaunawar a asusu na gaba.

## Guje wa kawo cikas

Kayan aikin da kake gudanarwa kafin kowane zama bai kamata ya ci wani abu da za ka lura da shi ba. ccam yana karanta fayil ɗin saitawa ɗaya, yana duba kowane mahaɗin da aka raba da system call ɗaya ko biyu, sannan yana maye gurbin kansa da `claude`. Yana ƙara kusan millisecond 4, kuma ta tsohuwa ba ya aiki yayin da kake aiki.

Yana kuma kula da fayilolinka. Kowane canji ga rabawa yana buga shiri kafin ya fara, kuma yana ajiye ainihin fayiloli a foldar kwafin ajiya. `ccam doctor` yana duba mahaɗai, izini da bayanan shiga, kuma yana gyara abin da zai iya gyarawa lafiya.

## Gwada shi

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add
```

[Farawa da sauri](/docs/quick-start) yana nuna yadda ake ƙara asusai biyu, kuma [yadda ake amfani da asusun Claude Code da dama](/guides/multiple-claude-code-accounts) yana bayyana hanyar tare da ccam da kuma ba tare da shi ba.

ccam open source ne kuma ba na hukuma ba ne: ba shi da alaƙa da Anthropic, kuma Anthropic ba ta amince da shi ba. Ko an yarda ka yi amfani da asusu fiye da ɗaya ya dogara da sharuɗɗan da suka shafi kowanne. Ana maraba da ra'ayoyi da fassarori a [GitHub](https://github.com/jhyunwoo/ccam).
