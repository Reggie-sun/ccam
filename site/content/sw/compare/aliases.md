---
title: ccam ikilinganishwa na alias za shell
description: Alias za shell zenye CLAUDE_CONFIG_DIR dhidi ya ccam kwa kuendesha akaunti kadhaa za Claude Code, zikilinganishwa sambamba.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 13da69070650
---

Alias ya shell kwa kila akaunti ndiyo njia rahisi zaidi ya kuendesha akaunti kadhaa za Claude Code. ccam hufanya jambo hilo hilo kwa ndani na huongeza kushiriki, kubadili na kufuatilia matumizi.

| | Alias ya shell | ccam |
| --- | --- | --- |
| Kuingia tofauti kwa kila akaunti | Ndiyo | Ndiyo |
| Usanidi wa kila akaunti | Hariri `~/.bashrc`, ingia | `ccam add` |
| Hufanya kazi katika skripti na programu nyingine | Hapana, alias ni za shell tu | Ndiyo, `claude2` ni amri halisi |
| Mipangilio, `CLAUDE.md`, viendelezi vinavyoshirikiwa | Nakili au unganisha kwa mkono | Kwa kila kipengele, na nakala rudufu |
| Seva za MCP zinazoshirikiwa | Kwa mkono, kwa kila akaunti | `ccam mcp add` mara moja |
| Endeleza mazungumzo kwenye akaunti nyingine | Hapana | `ccam next --continue` |
| Matumizi ya kila akaunti | Hapana | `ccam usage` |
| Chagua akaunti yenye nafasi kubwa zaidi | Hapana | `ccam next --auto` |
| Jina la akaunti kwenye mstari wa hali (status line) | Hapana | `ccam statusline install` |
| Rekebisha viungo vilivyoharibika | Kwa mkono | `ccam doctor --fix` |
| Programu ya ziada | Hakuna | Binary moja ya static |

Chagua alias ikiwa una akaunti mbili na hutaki zishiriki chochote. Chagua ccam ikiwa unataka akaunti zihisike kama usanidi mmoja. Unaweza kuanza na alias na kuhamia ccam baadaye: `ccam import` husajili folda zilizopo za `~/.claude-*` pale zilipo.

[Jinsi ya kutumia akaunti mbili kwenye kompyuta moja](/guides/multiple-claude-code-accounts) inaeleza njia zote mbili hatua kwa hatua.
