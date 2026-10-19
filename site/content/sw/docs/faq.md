---
title: Maswali yanayoulizwa mara kwa mara
metaTitle: "Maswali kuhusu ccam: akaunti kadhaa za Claude Code"
description: Majibu kuhusu ccam, akaunti za Claude Code, vikomo vya matumizi, kushiriki, faragha na mifumo inayotumika.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 8876e4f00bd2
---

## Je, ccam ni zana rasmi ya Anthropic?

Hapana. ccam ni mradi usio rasmi wa open source na hauhusiani na Anthropic wala haujaidhinishwa nayo. Inatumia mpangilio wa `CLAUDE_CONFIG_DIR` ambao [Claude Code inauelezea katika nyaraka zake](https://code.claude.com/docs/en/env-vars) kwa ajili ya kuendesha akaunti kadhaa sambamba.

## Je, ninaruhusiwa kutumia akaunti kadhaa za Claude Code?

Hilo linategemea masharti ya kila akaunti na, kwa akaunti za kazi, sera ya shirika lako. Yakague kabla ya kutumia akaunti kadhaa. ccam hurahisisha tu kubadili kati ya akaunti ambazo tayari unazo.

## Je, ccam hukwepa vikomo vya matumizi?

Hapana. Kila akaunti huhifadhi vikomo vyake. ccam huvionyesha na hukuwezesha kuendelea kwenye akaunti nyingine unayoimiliki pale moja inapoishiwa.

## Nitatumiaje akaunti mbili za Claude Code kwenye kompyuta moja?

Ipe kila akaunti folda yake ya usanidi kwa `CLAUDE_CONFIG_DIR`. ccam hukufanyia hivyo: endesha `ccam add` kwa kila akaunti, kisha uzianzishe kwa `ccam personal` na `ccam work`, au kwa njia fupi unazozipa majina mwenyewe, kama `ccwork`. [Mwongozo](/guides/multiple-claude-code-accounts) pia unaonyesha jinsi ya kufanya hivyo kwa mkono.

## Nitaendeleaje na mazungumzo kwenye akaunti nyingine?

Shiriki historia ya mazungumzo kati ya akaunti hizo mbili kwa `ccam share <name> --history on`, kisha endesha `ccam next --continue`. Claude Code hufungua mazungumzo ya hivi karibuni zaidi ya folda ya sasa kwenye akaunti inayofuata.

## Je, ccam huona au kutuma taarifa zangu za kuingia?

ccam kamwe hainakili wala haishiriki taarifa za kuingia. Husoma tokeni ya ufikiaji ya akaunti ili tu kuuliza endpoint ya matumizi ya Anthropic kuhusu matumizi ya akaunti hiyo, unapoendesha `ccam usage`, `ccam next --auto` au `ccam doctor`. ccam haina telemetry na haitumi chochote mahali pengine popote.

## Akaunti zangu zimehifadhiwa wapi?

Katika `~/.local/share/ccam/profiles/<name>/`, folda moja kwa kila akaunti. Kwenye macOS, Claude Code huhifadhi taarifa za kuingia zenyewe katika Keychain, chini ya jina linalotokana na folda hiyo.

## Nini kinatokea kwa ~/.claude yangu iliyopo?

Hakuna kitu, isipokuwa uendeshe `ccam import` au `ccam migrate`. Import husajili folda pale ilipo, kwa hiyo `claude` ya kawaida na `ccam <name>` huendelea kutumia akaunti na mipangilio ile ile. `ccam migrate` hunakili mipangilio yake, `CLAUDE.md`, viendelezi na seva za MCP kwenye mipangilio inayoshirikiwa, na husoma folda tu.

## ccam inatumika kwenye mifumo gani?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine na mingineyo) na macOS. Kwenye Windows, itumie ndani ya WSL.

## Kwa nini seva ya MCP inaonekana katika akaunti moja tu?

Seva zilizoongezwa ndani ya kipindi kwa `claude mcp add` huhifadhiwa kwa akaunti hiyo tu. ccam hukujulisha wakati ujao utakapoanzisha akaunti hiyo; hamishia seva kwenye kikundi kinachoshiriki kwa `ccam mcp import <account> <server>`.

## Je, ccam hupunguza kasi ya Claude Code?

Kuanzisha Claude Code kupitia ccam huongeza takriban milisekunde 4, zilizopimwa kwa `scripts/perf-check.sh` kwenye kompyuta mpakato. ccam husoma faili moja ndogo, hukagua viungo vinavyoshirikiwa kisha hujibadilisha yenyewe kuwa `claude`, kwa hiyo haiendeshwi wakati unafanya kazi.

## Nitasasishaje ccam?

Endesha `ccam update`. Hupakua toleo jipya zaidi, hukagua checksum yake na hubadilisha faili ya `ccam`, huku `ccam update --check` ikikuambia tu kama kuna toleo jipya zaidi. ccam kamwe haijisasishi yenyewe. Ikiwa ulisakinisha kifurushi cha `.deb`, badala yake sakinisha `.deb` mpya kwa apt. Tazama [kusasisha ccam](/docs/install#update).

## Nitaondoaje ccam?

Ondoa njia fupi ya kila akaunti kwa `ccam shortcut <account> --remove`, futa `ccam` kutoka `~/.local/bin`, kisha ondoa `~/.config/ccam`, `~/.local/share/ccam` na `~/.local/state/ccam` ikiwa huhitaji tena akaunti hizo. Tazama [usakinishaji](/docs/install#uninstall).
