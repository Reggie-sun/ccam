---
title: Ipinapakilala ang ccam
description: Bakit ginawa ang ccam, ano ang naitutulong nito sa mga may maraming Claude Code account, at paano ito hindi nakakasagabal.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: bbc9d9014f98
---

Ang ccam ay maliit na command-line tool para sa mga gumagamit ng higit sa isang Claude Code account sa iisang computer: isang personal at isang work account, o ilang account para sa isang mahabang proyekto.

## Ang problema

Sinusuportahan ng Claude Code ang maraming account sa pamamagitan ng `CLAUDE_CONFIG_DIR`: isang folder bawat account. Sa praktika, parang hiwalay na isla ang bawat folder. Kailangang i-set up ulit sa bawat folder ang settings, `CLAUDE.md`, mga MCP server, agent at skill, at panatilihing magkakatugma nang mano-mano. Nananatili ang mga usapan kung saan sila nagsimula. At kapag naabot ng isang account ang usage limit nito, malalaman mo lang kapag sinubukan mo.

## Ano ang ginagawa ng ccam

- Ginagawa ng `ccam add` ang folder, pinapatakbo ang login at, kung gusto mo, binibigyan ang account ng shortcut na ikaw mismo ang nagpangalan, gaya ng `ccwork`.
- Sinisimulan ng `ccam work` at `ccwork` ang Claude Code bilang account na iyon at ipinapasa nang walang pagbabago ang bawat argument.
- Shine-share ang mga MCP server, settings, `CLAUDE.md` at extension sa pamamagitan ng mga link, bawat item at bawat account. Hindi kailanman shine-share ang mga login.
- Ipinapakita ng `ccam usage` ang 5-oras at lingguhang usage ng bawat account; itinutuloy ng `ccam next --continue` ang usapan sa susunod na account.

## Hindi nakakasagabal

Ang tool na pinapatakbo mo bago ang bawat session ay hindi dapat magkaroon ng gastos na mapapansin mo. Nagbabasa ang ccam ng isang configuration file, sinusuri ang mga shared na link gamit ang tig-iisang system call, at saka pinapalitan ang sarili nito ng `claude`. Nagdadagdag ito ng mga 4 millisecond, at hindi ito tumatakbo habang nagtatrabaho ka.

Maingat din ito sa iyong mga file. Bawat pagbabago sa sharing ay nagpi-print ng plano bago magsimula at iniimbak ang mga orihinal sa isang backup folder. Sinusuri ng `ccam doctor` ang mga link, permission at login, at inaayos ang kayang ayusin nang ligtas.

## Subukan ito

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add
```

Open source at hindi opisyal ang ccam: hindi ito kaanib o inendorso ng Anthropic. Nakadepende sa terms ng bawat account kung puwede kang gumamit ng maraming account. Malugod naming tinatanggap ang feedback at mga pagsasalin sa [GitHub](https://github.com/jhyunwoo/ccam).
