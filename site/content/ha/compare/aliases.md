---
title: ccam idan aka kwatanta da shell alias
description: Shell alias da CLAUDE_CONFIG_DIR idan aka kwatanta da ccam don gudanar da asusun Claude Code da dama, fasali bayan fasali.
updated: 2026-10-26
order: 1
translation: machine
sourceHash: c30d84bceba0
---

Shell alias ga kowane asusu ita ce hanya mafi sauƙi ta gudanar da asusun Claude Code da dama. ccam yana yin abu ɗaya a ƙarƙashi kuma yana ƙara rabawa, sauyawa da bin diddigin amfani.

## Saitin alias

Ƙara layi ɗaya ga kowane ƙarin asusu a `~/.bashrc` ko `~/.zshrc`, buɗe sabon terminal kuma ka gudanar da alias ɗin sau ɗaya don shiga:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

`claude` yana ci gaba da zama asusunka na farko, kuma `claude-work` yana fara na biyu. Wannan ita ce hanyar da [takardun Claude Code suka bayyana](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts).

## Saitin ccam

Girka ccam kuma ka ƙara asusun, a nan tare da gajeren umarni nasa:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add work --shortcut ccwork
```

`ccam work` da `ccwork` suna fara asusun, kuma dukansu suna aiki a cikin scripts ma.

## Gefe da gefe

| | Shell alias | ccam |
| --- | --- | --- |
| Bayanan shiga daban-daban | Ee | Ee |
| Saitawa ga kowane asusu | Gyara `~/.bashrc`, shiga | `ccam add` |
| Yana aiki a scripts da sauran shirye-shirye | A'a, alias na shell ne kawai | Ee, `ccam work` da gajerun umarni kamar `ccwork` umarni ne na gaske |
| Saituna, `CLAUDE.md`, extensions da aka raba | Kwafa ko haɗa da hannu | Ee, ana zaɓa abu bayan abu, tare da kwafin ajiya |
| Sabobin MCP da aka raba | Da hannu, ga kowane asusu | `ccam mcp add` sau ɗaya |
| Ci gaba da tattaunawa a wani asusu | A'a | `ccam next --continue` |
| Amfanin kowane asusu | A'a | `ccam usage` |
| Zaɓar asusun da ya fi sarari | A'a | `ccam next --auto` |
| Sake gudanar da aikin `-p` a matsayin wani asusu a iyaka | A'a | `ccam --fallback` |
| Asusu ga kowace foldar project | A'a | `ccam folder` |
| Sunan asusu a layin matsayi | A'a | `ccam statusline install` |
| Gyara mahaɗan da suka lalace | Da hannu | `ccam doctor --fix` |
| Ƙarin software | Babu | Static binary ɗaya |

Zaɓi alias idan kana da asusai biyu kuma ba ka son a raba komai a tsakaninsu. Zaɓi ccam idan kana son asusan su ji kamar saiti ɗaya. Za ka iya farawa da alias ka koma ccam daga baya: `ccam import` yana yin rajistar foldojin `~/.claude-*` da ke akwai a inda suke.

[Yadda ake amfani da asusai da dama a kwamfuta ɗaya](/guides/multiple-claude-code-accounts) yana bayyana hanyoyin biyu mataki-mataki, kuma [ccam da sauran masu sauya asusu](/compare/account-switchers) yana kwatanta ccam da kayan aikin da ke musanya bayanan shiga ɗaya.
