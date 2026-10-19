---
title: ccam na vibadilisha akaunti vingine vya Claude Code
description: Jinsi ccam inavyolinganishwa na claude-swap, CCS, clauth, CAAM, cux na zana nyingine za kubadili akaunti za Claude Code, na njia ipi inafaa matumizi yapi.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: dc907c88e2bf
---

Zana za kubadili akaunti za Claude Code hufuata mojawapo ya njia mbili. Nyingi hubadilisha taarifa za kuingia ndani ya folda moja ya `~/.claude`, kwa hiyo akaunti moja tu hutumika kwa wakati mmoja. Nyingine, ikiwemo ccam, huipa kila akaunti folda yake kwa `CLAUDE_CONFIG_DIR`, kwa hiyo kila akaunti hubaki imeingia na kadhaa zinaweza kuendeshwa kwa pamoja. ccam huongeza mipangilio na seva za MCP zinazoshirikiwa, na huendeleza mazungumzo katika folda ya akaunti nyingine.

## Kubadilisha kuingia kumoja au folda tofauti?

**Kubadilisha kuingia** huweka seti moja ya mipangilio, seva za MCP na historia, kwa sababu kuna folda moja tu. Kubadili ni haraka, na kipindi kinachoendelea kinaweza kuchukua kuingia mpya. Lakini akaunti moja tu hutumika kwa wakati mmoja, kwa kila terminal na kihariri kwa pamoja.

**Folda tofauti** ndiyo njia ambayo [nyaraka za Claude Code zinaeleza](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Kila akaunti hubaki imeingia, na unaweza kuendesha kipindi cha kazi na cha binafsi sambamba. Gharama yake ni kwamba kila folda huanza na mipangilio, seva za MCP na historia yake yenyewe, isipokuwa zana fulani izishiriki.

## Jinsi zana zinavyolinganishwa

Imechukuliwa kutoka nyaraka za kila mradi tarehe 2026-10-19. "Haijaelezwa" inamaanisha mradi hausemi, si kwamba kipengele hakipo.

| Zana | Njia | Akaunti kwa pamoja | Matumizi ya saa 5 na ya wiki | Kikomo kinapofikiwa | Mifumo |
| --- | --- | --- | --- | --- | --- |
| ccam | Folda tofauti | Ndiyo | Ndiyo, `ccam usage` | Unaendesha `ccam next --auto --continue` ili kuendeleza mazungumzo katika akaunti yenye nafasi kubwa zaidi | Linux, macOS, WSL |
| [claude-swap](https://github.com/realiti4/claude-swap) | Hubadilisha kuingia; folda tofauti kwa `cswap run` | Ndiyo, kwa `cswap run` | Ndiyo | `cswap auto` hubadilisha kuingia kwenye kiwango fulani; kipindi kinachoendelea huendelea | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Folda tofauti | Ndiyo | Kwa sehemu, katika programu ya menu bar ya macOS | Makundi ya akaunti hupeleka maombi kwa akaunti inayofuata kupitia proxy ya ndani | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Hubadilisha kuingia; folda tofauti kwa `clauth start` | Ndiyo, kwa `clauth start` | Ndiyo | Hubadili kufuata mnyororo wa akaunti mbadala | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Hubadilisha kuingia, au profile zilizotengwa | Ndiyo, kwa profile zilizotengwa | Ndiyo | `caam run` hujaribu tena amri zisizo shirikishi kwenye akaunti nyingine | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Hubadilisha kuingia na kuzunguka `claude` | Haijaelezwa | Ndiyo | Hubadilisha kuingia na kuendeleza mazungumzo yale yale | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Hubadilisha kuingia | Hapana | Ndiyo | Haijaelezwa | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Hubadilisha kuingia | Hapana | Ndiyo | Kubadili kwa hiari kwenye 100% | macOS |

ccswap ni fork ya claude-swap ambayo pia husimamia kuingia kwa Codex. Zana kadhaa kati ya hizi pia hushughulikia mawakala wengine wa kuandika msimbo, kama Codex au Gemini.

## ccam inazingatia nini

- Kila akaunti hubaki imeingia. Kubadili hakutoi akaunti yoyote, na `ccam work` na `ccam personal` zinaweza kuendeshwa katika terminal mbili kwa wakati mmoja.
- Seva za MCP, `settings.json`, `CLAUDE.md`, agents, skills na plugins hushirikiwa kati ya folda kupitia viungo, kipengele kwa kipengele na akaunti kwa akaunti. Taarifa za kuingia hazishirikiwi kamwe.
- Historia ikishirikiwa, amri moja, `ccam next --auto --continue`, hufungua mazungumzo yale yale katika folda ya akaunti yenye nafasi kubwa zaidi.
- ccam huanzisha Claude Code na kujibadilisha nayo, kwa hiyo hakuna kinachoendeshwa nyuma. Haifuatilii matumizi yako wala kubadili yenyewe; wewe ndiye unaamua lini.

## Pale zana nyingine zinapofanya zaidi

- claude-swap, CCS, clauth, CAAM na cux huendeshwa kwenye Windows bila WSL. ccam haifanyi hivyo.
- claude-swap, clauth na cux hubadili zenyewe akaunti inapovuka kiwango fulani cha matumizi. ccam hukuachia uamuzi wa kubadili.
- CCSwitcher, Claude Switcher na claude-swap huonyesha matumizi kwenye menu bar ya macOS.
- CCS, CAAM, ccswap na Claude Switcher pia husimamia Codex na mawakala wengine wa kuandika msimbo.

## Ipi inafaa?

Chagua zana ya kubadilisha kuingia ikiwa unatumia akaunti moja kwa wakati mmoja, unataka kubadili kufanyike kiotomatiki, au uko kwenye Windows bila WSL. Chagua ccam ikiwa unataka akaunti sambamba, usanidi mmoja unaoshirikiwa na mazungumzo yale yale katika akaunti inayofuata. Ikiwa huhitaji kushiriki chochote, [alias za kawaida za shell](/compare/aliases) zinatosha.

Iwapo unaruhusiwa kutumia akaunti kadhaa inategemea masharti ya kila akaunti. Yakague kabla ya kufanya hivyo.
