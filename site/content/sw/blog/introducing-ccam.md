---
title: Tunakuletea ccam
description: Kwa nini ccam ipo, inawafanyia nini watu wenye akaunti kadhaa za Claude Code, na jinsi inavyokaa nje ya njia yako.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 5e4121938d8c
---

ccam ni zana ndogo ya mstari wa amri kwa watu wanaotumia zaidi ya akaunti moja ya Claude Code kwenye kompyuta moja: akaunti binafsi na ya kazi, au akaunti kadhaa kwa mradi mmoja mrefu.

## Tatizo

Claude Code inaruhusu akaunti kadhaa kupitia `CLAUDE_CONFIG_DIR`: folda moja kwa kila akaunti. Kwa vitendo, kila folda ni kisiwa. Mipangilio, `CLAUDE.md`, seva za MCP, agents na skills lazima visanidiwe upya katika kila folda na vilinganishwe kwa mkono. Mazungumzo hubaki pale yalipoanzia. Na akaunti moja inapofikia kikomo chake cha matumizi, unagundua hilo kwa kujaribu.

## ccam hufanya nini

- `ccam add` huunda folda, huendesha mchakato wa kuingia na huipa akaunti njia fupi: `claude1`, `claude2`, na kadhalika.
- `ccs work` na `claude2` huanzisha Claude Code kama akaunti hiyo na hupitisha kila hoja bila kuibadilisha.
- Seva za MCP, mipangilio, `CLAUDE.md` na viendelezi hushirikiwa kupitia viungo, kipengele kwa kipengele na akaunti kwa akaunti. Taarifa za kuingia hazishirikiwi kamwe.
- `ccs usage` huonyesha matumizi ya saa 5 na ya wiki ya kila akaunti; `ccs next --continue` huendeleza mazungumzo kwenye akaunti inayofuata.

## Kukaa nje ya njia yako

Zana unayoiendesha kabla ya kila kipindi haipaswi kugharimu chochote unachoweza kukigundua. ccam husoma faili moja ya usanidi, hukagua viungo vinavyoshirikiwa kwa system call moja kwa kila kiungo, kisha hujibadilisha yenyewe kuwa `claude`. Huongeza takriban milisekunde 4, na haiendeshwi wakati unafanya kazi.

Pia ni makini na faili zako. Kila badiliko la kushiriki huchapisha mpango kabla ya kuanza na huhifadhi faili asili katika folda ya nakala rudufu. `ccam doctor` hukagua viungo, ruhusa na taarifa za kuingia, na hurekebisha kile kinachoweza kurekebishwa kwa usalama.

## Ijaribu

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add
```

ccam ni open source na si rasmi: haihusiani na Anthropic wala haijaidhinishwa nayo. Iwapo unaruhusiwa kutumia akaunti kadhaa inategemea masharti ya kila akaunti. Maoni na tafsiri yanakaribishwa kwenye [GitHub](https://github.com/jhyunwoo/ccam).
