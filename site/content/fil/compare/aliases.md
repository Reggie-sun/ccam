---
title: Ang ccam kumpara sa mga shell alias
description: Mga shell alias na may CLAUDE_CONFIG_DIR laban sa ccam para sa pagpapatakbo ng maraming Claude Code account, magkatabing inihambing.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 13da69070650
---

Ang shell alias bawat account ang pinakasimpleng paraan para magpatakbo ng maraming Claude Code account. Ganoon din ang ginagawa ng ccam sa ilalim, at nagdadagdag ito ng sharing, paglipat at pagsubaybay sa usage.

| | Shell alias | ccam |
| --- | --- | --- |
| Magkakahiwalay na login | Oo | Oo |
| Setup bawat account | I-edit ang `~/.bashrc`, mag-log in | `ccam add` |
| Gumagana sa mga script at ibang program | Hindi, sa shell lang ang mga alias | Oo, totoong command ang `claude2` |
| Shared na settings, `CLAUDE.md`, extension | Kopyahin o i-link nang mano-mano | Bawat item, may backup |
| Shared na mga MCP server | Mano-mano, bawat account | `ccam mcp add` nang isang beses |
| Ituloy ang usapan sa ibang account | Hindi | `ccam next --continue` |
| Usage ng bawat account | Hindi | `ccam usage` |
| Piliin ang account na may pinakamaraming natitira | Hindi | `ccam next --auto` |
| Pangalan ng account sa status line | Hindi | `ccam statusline install` |
| Ayusin ang mga sirang link | Mano-mano | `ccam doctor --fix` |
| Karagdagang software | Wala | Isang static na binary |

Piliin ang mga alias kung may dalawa kang account at ayaw mong may i-share sa pagitan nila. Piliin ang ccam kung gusto mong maging parang iisang setup ang mga account. Puwede kang magsimula sa mga alias at lumipat sa ccam sa ibang pagkakataon: nirerehistro ng `ccam import` ang mga umiiral na `~/.claude-*` folder kung nasaan sila.

Ipinapaliwanag nang hakbang-hakbang ng [Paano gumamit ng dalawang account sa iisang computer](/guides/multiple-claude-code-accounts) ang dalawang paraan.
