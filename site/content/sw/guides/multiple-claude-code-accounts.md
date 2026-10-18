---
title: Jinsi ya kutumia akaunti mbili za Claude Code kwenye kompyuta moja
description: Endesha akaunti kadhaa za Claude Code sambamba kwa CLAUDE_CONFIG_DIR, kwa mkono kwa alias za shell au kiotomatiki kwa ccam.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: d82129ddcb4c
---

Ipe kila akaunti folda yake ya usanidi. Claude Code husoma folda kutoka kigezo cha mazingira `CLAUDE_CONFIG_DIR`, na kila folda huhifadhi taarifa zake za kuingia, mipangilio na historia. Unaweza kusanidi hili kwa mkono kwa alias ya shell, au uiachie ccam ifanye.

## Kwa mkono, kwa alias

Hii ndiyo njia inayoelezwa na nyaraka za Claude Code. Ongeza mstari mmoja kwa kila akaunti ya ziada katika `~/.bashrc` au `~/.zshrc`:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Fungua terminal mpya na uendeshe `claude-work`. Claude Code hukuomba uingie, kwa sababu folda ni mpya. Kuanzia hapo, `claude` ni akaunti yako ya kwanza na `claude-work` ni ya pili.

Hii hufanya kazi vizuri kwa akaunti mbili. Mambo ya kujua:

- Kila folda huanza ikiwa tupu: hakuna mipangilio, hakuna `CLAUDE.md`, hakuna seva za MCP, hakuna agents wala skills. Unanakili au kuunganisha unachohitaji, na unahakikisha nakala zinalingana wewe mwenyewe.
- Alias zipo katika shell shirikishi (interactive) tu. Skripti, Makefiles na programu nyingine hazizioni.
- Seva za MCP ziko katika `.claude.json`, pamoja na taarifa za kuingia, kwa hiyo haziwezi kuunganishwa tu kati ya folda.
- Mazungumzo hubaki katika akaunti yalimoanzia, kwa hiyo huwezi kuendeleza mazungumzo kwenye akaunti nyingine.
- Hakuna kitu kinachoonyesha matumizi ya akaunti zote kwa pamoja.

## Kwa ccam

ccam hufanya wazo hilo hilo kiotomatiki na huongeza sehemu zinazokosekana:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
```

Sasa `claude1` na `claude2`, au `ccs personal` na `ccs work`, huanzisha akaunti hizo mbili. Kwa kuwa ni programu halisi badala ya alias, zinafanya kazi pia ndani ya skripti. Mipangilio, `CLAUDE.md`, seva za MCP na viendelezi hushirikiwa tangu mwanzo, na unaweza kuwasha historia ya mazungumzo inayoshirikiwa ili kuendeleza mazungumzo kwenye akaunti nyingine kwa `ccs next --continue`.

Ikiwa tayari unatumia `~/.claude`, endesha `ccam import` badala ya kuongeza akaunti hiyo tena.

## Ipi uchague

Kwa akaunti mbili bila kushiriki chochote, alias inatosha. Unapotaka mipangilio na seva za MCP zinazoshirikiwa, akaunti zaidi, kuona matumizi kwa haraka au kuendeleza mazungumzo, ccam hukuepushia kazi ya kufuatilia mambo hayo. [Ulinganisho](/compare/aliases) unaorodhesha tofauti sambamba.
