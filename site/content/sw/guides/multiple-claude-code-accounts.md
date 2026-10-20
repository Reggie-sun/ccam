---
title: Akaunti kadhaa za Claude Code kwenye kompyuta moja
description: Endesha akaunti kadhaa za Claude Code sambamba kwa CLAUDE_CONFIG_DIR, kwa mkono kwa alias za shell au kiotomatiki kwa ccam, na ubadili kati yao.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: "ee92078b81c3"
---

Ipe kila akaunti folda yake ya usanidi. Claude Code husoma folda kutoka kigezo cha mazingira `CLAUDE_CONFIG_DIR`, na kila folda huhifadhi taarifa zake za kuingia, mipangilio na historia, kwa hiyo kila akaunti hubaki imeingia na unaweza kuziendesha kwa wakati mmoja. Unaweza kusanidi hili kwa mkono kwa alias ya shell, au uiachie ccam ifanye.

## Ninasanidije kwa mkono?

Hii ndiyo njia ambayo [nyaraka za Claude Code zinaeleza](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Ongeza mstari mmoja kwa kila akaunti ya ziada katika `~/.bashrc` au `~/.zshrc`:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Fungua terminal mpya na uendeshe `claude-work`. Claude Code hukuomba uingie, kwa sababu folda ni mpya. Kuanzia hapo, `claude` ni akaunti yako ya kwanza na `claude-work` ni ya pili. Ongeza alias moja kwa kila akaunti, kila moja ikiwa na folda yake.

Kwenye Windows bila WSL, function ya PowerShell katika profile yako (`notepad $PROFILE`) hufanya vivyo hivyo:

```powershell
function claude-work { $env:CLAUDE_CONFIG_DIR = "$HOME\.claude-work"; try { claude @args } finally { Remove-Item Env:CLAUDE_CONFIG_DIR } }
```

## Njia ya alias inakosa nini?

Hufanya kazi vizuri kwa akaunti mbili. Mambo ya kujua:

- Kila folda huanza ikiwa tupu: hakuna mipangilio, hakuna `CLAUDE.md`, hakuna seva za MCP, hakuna agents wala skills. Unanakili au kuunganisha unachohitaji, na unahakikisha nakala zinalingana wewe mwenyewe.
- Alias zipo katika shell shirikishi (interactive) tu. Skripti, Makefiles na programu nyingine hazizioni.
- Seva za MCP ziko katika `.claude.json`, pamoja na taarifa za kuingia, kwa hiyo haziwezi kuunganishwa tu kati ya folda.
- Mazungumzo hubaki katika akaunti yalimoanzia, kwa hiyo huwezi kuendeleza mazungumzo kwenye akaunti nyingine.
- Hakuna kitu kinachoonyesha matumizi ya akaunti zote kwa pamoja.

## Je, naweza kubadili akaunti kwa /login badala yake?

Ndiyo, ikiwa unabadili mara chache tu. Kuandika `/login` ndani ya Claude Code [hubadili akaunti](https://code.claude.com/docs/en/quickstart) ya folda uliyomo. Akaunti moja tu huwa imeingia kwa wakati mmoja, kwa hiyo unaingia upya kila unapobadili, na huwezi kuendesha akaunti mbili sambamba. Folda tofauti huepuka matatizo yote mawili.

## Ninasanidije kwa ccam?

ccam hufanya wazo hilo hilo kiotomatiki na huongeza sehemu zinazokosekana:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

Sasa `ccam personal` na `ccam work` huanzisha akaunti hizo mbili, na vivyo hivyo `ccwork`, njia fupi uliyoipa akaunti ya kazi. Kwa kuwa ni programu halisi badala ya alias, hufanya kazi pia katika skripti. Endesha `ccam` peke yake ili kuchagua akaunti kutoka kwenye orodha, na `ccam use work` ili kuifanya akaunti moja kuwa chaguo-msingi.

Mipangilio, `CLAUDE.md`, seva za MCP na viendelezi hushirikiwa tangu mwanzo, na unaweza kuwasha historia ya mazungumzo inayoshirikiwa ili kuendeleza mazungumzo kwenye akaunti nyingine kwa `ccam next --continue`. `ccam usage` huonyesha matumizi ya saa 5 na ya wiki ya kila akaunti kwa pamoja.

Ikiwa tayari unatumia `~/.claude`, endesha `ccam import` badala ya kuongeza akaunti hiyo tena, au `ccam migrate` ili kunakili mipangilio yake, seva za MCP na viendelezi kwenye mipangilio inayoshirikiwa.

## Nitajuaje ni akaunti ipi inaendeshwa?

Ndani ya Claude Code, `/status` huonyesha akaunti uliyoingia nayo. Kwa ccam, endesha `ccam statusline install` mara moja na mstari wa hali (status line) utaonyesha jina la akaunti na matumizi yake wakati wote, kwa mfano `[work] · Opus · 5h 24% · 7d 41%`.

## Je, hii inafanya kazi kwenye Windows na WSL?

Claude Code huendeshwa [moja kwa moja kwenye Windows na ndani ya WSL](https://code.claude.com/docs/en/setup). `CLAUDE_CONFIG_DIR` hufanya kazi katika zote mbili, kwa hiyo alias au function ya PowerShell iliyo hapo juu inatosha kwenye Windows bila WSL. ccam nayo huendeshwa moja kwa moja kwenye Windows: isakinishe katika PowerShell kwa `irm https://ccam.dev/install.ps1 | iex`, au ndani ya WSL kwa amri ya Linux. Njia fupi huwa `ccwork.exe`, na kushiriki faili kama `settings.json` kunahitaji Developer Mode ([Sakinisha ccam](/docs/install)).

## Kutatua matatizo

**Claude Code huniomba niingie kila mara.** Mwanzo wa kwanza katika folda mpya huomba kila wakati. Ikiendelea kuomba, muda wa kuingia umeisha au taarifa za kuingia zimeondolewa: endesha `ccam login work` (au anzisha alias mara moja na uingie).

**Seva ya MCP inaonekana katika akaunti moja tu.** Seva zinazoongezwa ndani ya kipindi kwa `claude mcp add` ni za akaunti hiyo tu. Ihamishe kwa akaunti zote kwa `ccam mcp import work <server>`.

**Mipangilio inatofautiana kati ya akaunti.** `ccam share work` huonyesha kile ambacho akaunti inashiriki na kupitia kikundi kipi; `ccam group join default work` huifanya ishiriki kile ambacho `default` hushiriki, baada ya kuonyesha kitakachobadilika. `ccam doctor --fix` hurekebisha viungo vilivyoharibika na huhifadhi nakala rudufu ya kila faili inayoibadilisha.

**Akaunti mbili zinafungua mazungumzo yale yale.** Epuka kufanya kazi kwenye mazungumzo moja katika akaunti mbili kwa wakati mmoja. ccam hukuonya ukijaribu.

## Nichague ipi?

Kwa akaunti mbili bila kushiriki chochote, alias inatosha. Unapotaka mipangilio na seva za MCP zinazoshirikiwa, akaunti zaidi, kuona matumizi kwa haraka au kuendeleza mazungumzo, ccam hukuepushia kazi ya kufuatilia mambo hayo. [Ulinganisho na alias](/compare/aliases) unaorodhesha tofauti sambamba, na [ccam na zana nyingine za kubadili akaunti](/compare/account-switchers) inaeleza zana zinazobadilisha taarifa moja ya kuingia badala yake.

Iwapo unaruhusiwa kutumia akaunti kadhaa inategemea masharti ya kila akaunti. Yakague kabla ya kufanya hivyo.
