---
title: Tambayoyin da ake yawan yi
description: Amsoshi game da ccam, asusun Claude Code, iyakokin amfani, rabawa, sirri da tsarukan da ake tallafawa.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: e07d468d453b
---

## Shin ccam kayan aikin Anthropic ne na hukuma?

A'a. ccam aikin open-source ne da ba na hukuma ba, ba shi da alaƙa da Anthropic, kuma Anthropic ba ta amince da shi ba. Yana amfani da saitin `CLAUDE_CONFIG_DIR` wanda Claude Code ya rubuta bayaninsa don gudanar da asusai da dama gefe da gefe.

## Shin an yarda in yi amfani da asusun Claude Code da dama?

Wannan ya dogara da sharuɗɗan kowane asusu kuma, ga asusan aiki, da manufofin ƙungiyarka. Ka duba su kafin ka yi amfani da asusai da dama. ccam yana sauƙaƙa sauyawa ne kawai tsakanin asusan da kake da su tuni.

## Shin ccam yana kauce wa iyakokin amfani?

A'a. Kowane asusu yana riƙe iyakokinsa. ccam yana nuna su kuma yana ba ka damar ci gaba a wani asusun da ka mallaka idan iyakar ɗaya ta ƙare.

## Ta yaya zan yi amfani da asusun Claude Code biyu a kwamfuta ɗaya?

Ba kowane asusu foldar saitunsa ta kansa da `CLAUDE_CONFIG_DIR`. ccam yana yi maka wannan: gudanar da `ccam add` ga kowane asusu, sannan ka fara su da `ccam personal` da `ccam work`, ko da gajerun umarni da kai da kanka ka sa musu suna, kamar `ccwork`. [Jagorar](/guides/multiple-claude-code-accounts) kuma tana nuna yadda ake yi da hannu.

## Ta yaya zan ci gaba da tattaunawa a wani asusu?

Raba tarihin tattaunawa tsakanin asusun biyu da `ccam share <name> --history on`, sannan ka gudanar da `ccam next --continue`. Claude Code yana buɗe tattaunawa ta baya-bayan nan ta foldar yanzu a asusu na gaba.

## Shin ccam yana ganin bayanan shigata ko yana aika su?

ccam ba ya taɓa kwafa ko raba bayanan shiga. Yana karanta access token na asusu ne kawai don tambayar endpoint na amfani na Anthropic game da amfanin wannan asusun, lokacin da ka gudanar da `ccam usage`, `ccam next --auto` ko `ccam doctor`. ccam ba shi da telemetry kuma ba ya aika komai zuwa wani wuri dabam.

## A ina ake ajiye asusuna?

A `~/.local/share/ccam/profiles/<name>/`, folda ɗaya ga kowane asusu. A macOS, Claude Code yana ajiye bayanan shigar kansu a Keychain, ƙarƙashin sunan da aka samo daga wannan foldar.

## Me zai faru da ~/.claude da nake da shi?

Babu komai, sai idan ka gudanar da `ccam import` ko `ccam migrate`. Import yana yin rajistar foldar a inda take, don haka `claude` na yau da kullum da `ccam <name>` suna ci gaba da amfani da asusu da saituna iri ɗaya. `ccam migrate` yana kwafar saitunanta, `CLAUDE.md`, extensions da sabobin MCP zuwa saitunan da aka raba, kuma yana karanta foldar ne kawai.

## Waɗanne tsaruka ne ccam yake tallafawa?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine da sauransu) da macOS. A Windows, yi amfani da shi a cikin WSL.

## Me yasa sabar MCP take bayyana a asusu ɗaya kawai?

Sabobin da aka ƙara a cikin zama (session) da `claude mcp add` ana ajiye su ga wannan asusun kawai. ccam yana gaya maka a gaba da ka fara asusun; motsa sabar zuwa rukunin da aka raba da `ccam mcp import <account> <server>`.

## Shin ccam yana rage saurin Claude Code?

Fara Claude Code ta ccam yana ƙara kusan millisecond 4, wanda aka auna da `scripts/perf-check.sh` a kwamfutar tafi-da-gidanka. ccam yana karanta ƙaramin fayil ɗaya, yana duba mahaɗan da aka raba, sannan yana maye gurbin kansa da `claude`, don haka ba ya aiki yayin da kake aiki.

## Ta yaya zan cire ccam?

Cire gajeren umarnin kowane asusu da `ccam shortcut <account> --remove`, goge `ccam` daga `~/.local/bin`, sannan ka cire `~/.config/ccam`, `~/.local/share/ccam` da `~/.local/state/ccam` idan ba ka buƙatar asusan kuma. Duba [shigarwa](/docs/install#uninstall).
