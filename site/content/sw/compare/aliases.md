---
title: ccam ikilinganishwa na alias za shell
description: Alias za shell zenye CLAUDE_CONFIG_DIR dhidi ya ccam kwa kuendesha akaunti kadhaa za Claude Code, zikilinganishwa sambamba.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: b9fbcb88f111
---

Alias ya shell kwa kila akaunti ndiyo njia rahisi zaidi ya kuendesha akaunti kadhaa za Claude Code. ccam hufanya jambo hilo hilo kwa ndani na huongeza kushiriki, kubadili na kufuatilia matumizi.

## Usanidi wa alias

Ongeza mstari mmoja kwa kila akaunti ya ziada kwenye `~/.bashrc` au `~/.zshrc`, fungua terminal mpya na uendeshe alias mara moja ili kuingia:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

`claude` inabaki kuwa akaunti yako ya kwanza na `claude-work` huanzisha ya pili. Hii ndiyo njia ambayo [nyaraka za Claude Code zinaeleza](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts).

## Usanidi wa ccam

Sakinisha ccam na uongeze akaunti, hapa ikiwa na amri yake ya njia fupi:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add work --shortcut ccwork
```

`ccam work` na `ccwork` huanzisha akaunti, na zote mbili hufanya kazi kwenye skripti pia.

## Kando kwa kando

| | Alias ya shell | ccam |
| --- | --- | --- |
| Kuingia tofauti kwa kila akaunti | Ndiyo | Ndiyo |
| Usanidi wa kila akaunti | Hariri `~/.bashrc`, ingia | `ccam add` |
| Hufanya kazi katika skripti na programu nyingine | Hapana, alias ni za shell tu | Ndiyo, `ccam work` na njia fupi kama `ccwork` ni amri halisi |
| Mipangilio, `CLAUDE.md`, viendelezi vinavyoshirikiwa | Nakili au unganisha kwa mkono | Kwa kila kipengele, na nakala rudufu |
| Seva za MCP zinazoshirikiwa | Kwa mkono, kwa kila akaunti | `ccam mcp add` mara moja |
| Endeleza mazungumzo kwenye akaunti nyingine | Hapana | `ccam next --continue` |
| Matumizi ya kila akaunti | Hapana | `ccam usage` |
| Chagua akaunti yenye nafasi kubwa zaidi | Hapana | `ccam next --auto` |
| Endesha tena kazi ya `-p` kama akaunti nyingine kwenye kikomo | Hapana | `ccam --fallback` |
| Akaunti kwa kila folda ya mradi | Hapana | `ccam folder` |
| Jina la akaunti kwenye mstari wa hali (status line) | Hapana | `ccam statusline install` |
| Rekebisha viungo vilivyoharibika | Kwa mkono | `ccam doctor --fix` |
| Programu ya ziada | Hakuna | Binary moja ya static |

Chagua alias ikiwa una akaunti mbili na hutaki zishiriki chochote. Chagua ccam ikiwa unataka akaunti zihisike kama usanidi mmoja. Unaweza kuanza na alias na kuhamia ccam baadaye: `ccam import` husajili folda zilizopo za `~/.claude-*` pale zilipo.

[Jinsi ya kutumia akaunti kadhaa kwenye kompyuta moja](/guides/multiple-claude-code-accounts) inaeleza njia zote mbili hatua kwa hatua, na [ccam na zana nyingine za kubadili akaunti](/compare/account-switchers) inailinganisha ccam na zana zinazobadilisha taarifa moja ya kuingia.
