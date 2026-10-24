---
title: Maswali yanayoulizwa mara kwa mara
metaTitle: "Maswali kuhusu ccam: akaunti kadhaa za Claude Code"
description: Majibu kuhusu ccam, akaunti za Claude Code, vikomo vya matumizi, kushiriki, faragha na mifumo inayotumika.
updated: 2026-10-24
order: 6
translation: machine
sourceHash: "c4e3b89bb4ea"
---

## Je, ccam ni zana rasmi ya Anthropic?

Hapana. ccam ni mradi usio rasmi wa open source na hauhusiani na Anthropic wala haujaidhinishwa nayo. Inatumia mpangilio wa `CLAUDE_CONFIG_DIR` ambao [Claude Code inauelezea katika nyaraka zake](https://code.claude.com/docs/en/env-vars) kwa ajili ya kuendesha akaunti kadhaa sambamba.

## Je, ninaruhusiwa kutumia akaunti kadhaa za Claude Code?

Hilo linategemea masharti ya kila akaunti na, kwa akaunti za kazi, sera ya shirika lako. Yakague kabla ya kutumia akaunti kadhaa. ccam hurahisisha tu kubadili kati ya akaunti ambazo tayari unazo.

**Viti (seats) vya Team na Enterprise.** Chaguo za kiotomatiki (`ccam next`, `--fallback`, `ccam pick`, kitufe `n` cha dashibodi) pia hubadili kati ya viti vya shirika lile lile la Claude Team au Enterprise. Kituo cha usaidizi cha Anthropic kinasema "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (kila kiti katika mpango wa Team hutolewa kwa mtu mmoja aliyetajwa kwa jina, na ni mtu huyo pekee anayeruhusiwa kukitumia) na "Users can't route their work through another user's account to use that user's remaining usage" (watumiaji hawaruhusiwi kuelekeza kazi yao kupitia akaunti ya mtumiaji mwingine ili kutumia kiasi cha matumizi kilichobaki cha mtumiaji huyo) ([mipango ya Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); vivyo hivyo kwa [mipango ya Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Kwa hiyo, kutumia chaguo za kiotomatiki kati ya viti kunaweza kukiuka masharti ya Anthropic. Wasiliana na shirika lako na Anthropic kabla ya kuzitegemea. Mipango ya binafsi (Pro, Max) inaelezwa katika aya iliyo hapo juu.

## Je, ccam hukwepa vikomo vya matumizi?

Hapana. Kila akaunti huhifadhi vikomo vyake. ccam huvionyesha na hukuwezesha kuendelea kwenye akaunti nyingine unayoimiliki pale moja inapoishiwa.

## Nitatumiaje akaunti mbili za Claude Code kwenye kompyuta moja?

Ipe kila akaunti folda yake ya usanidi kwa `CLAUDE_CONFIG_DIR`. ccam hukufanyia hivyo: endesha `ccam add` kwa kila akaunti, kisha uzianzishe kwa `ccam personal` na `ccam work`, au kwa njia fupi unazozipa majina mwenyewe, kama `ccwork`. [Mwongozo](/guides/multiple-claude-code-accounts) pia unaonyesha jinsi ya kufanya hivyo kwa mkono.

## Nitaendeleaje na mazungumzo kwenye akaunti nyingine?

Endesha `ccam next --continue` ndani ya folda: huhamisha mazungumzo ya mwisho ya folda hiyo kwenda akaunti inayofuata (`--auto` kwa ile yenye nafasi kubwa zaidi). Ndani ya Claude Code, `/ccam:switch` kutoka [plugin ya ccam](/docs/install#claude-code-plugin) hufanya vivyo hivyo kwa kipindi ulichomo. Akaunti hizo mbili zinaposhiriki historia ya mazungumzo (`ccam group set default --add history`, au kikundi chao wenyewe: `ccam group add pair --items history --members work,side`), akaunti mpya hufungua mazungumzo hayo hayo; vinginevyo Claude Code huyanakili kwenye historia ya akaunti mpya.

## Je, ccam huona au kutuma taarifa zangu za kuingia?

ccam kamwe haishiriki wala kunakili taarifa zako za kuingia za Claude kati ya akaunti, na kamwe hainakili kuingia kwa seva za MCP. Husoma tokeni ya ufikiaji ya akaunti ili tu kuuliza endpoint ya matumizi ya Anthropic kuhusu matumizi ya akaunti hiyo, unapoendesha `ccam usage`, `ccam next --auto` au `ccam doctor`, na hufanya upya tokeni iliyokwisha muda kwa `ccam usage --refresh` pekee; kwa `usage_source = "claude"` haisomi tokeni hata kidogo. `ccam ext` na `ccam doctor` husoma majina ya kuingia kwa MCP ya akaunti, kamwe si tokeni zake. ccam haina telemetry. Mbali na maswali hayo, ombi pekee inalotuma yenyewe ni swali la dashibodi kwa GitHub kuhusu toleo jipya zaidi la ccam, ambalo halibebi chochote kukuhusu wewe au akaunti zako.

## Nitaingiaje kwenye seva ya mbali?

Endesha `ccam add` au `ccam login <name>` kwenye kipindi cha SSH kama kawaida. Claude Code husubiri kuingia kwenye seva, kwa hiyo kivinjari kwenye kompyuta yako huishia kwenye ukurasa wa `http://localhost:<port>/callback?code=…` ambao haufunguki. Nakili anwani hiyo kutoka kwenye kivinjari na uibandike kwenye terminal: ccam huipitisha kwa Claude Code kwenye seva, ambayo hukamilisha kuingia. Anwani hiyo ina msimbo wa matumizi ya mara moja, kwa hiyo usiishiriki.

## Akaunti zangu zimehifadhiwa wapi?

Katika `~/.local/share/ccam/profiles/<name>/`, folda moja kwa kila akaunti. Kwenye macOS, Claude Code huhifadhi taarifa za kuingia zenyewe katika Keychain, chini ya jina linalotokana na folda hiyo, kwa hiyo huko `ccam rename` huhifadhi jina la zamani la folda.

## Nini kinatokea kwa ~/.claude yangu iliyopo?

Hakuna kitu, isipokuwa uendeshe `ccam import` au `ccam migrate`. Import husajili folda pale ilipo, kwa hiyo `claude` ya kawaida na `ccam <name>` huendelea kutumia akaunti na mipangilio ile ile. `ccam migrate` hunakili mipangilio yake, `CLAUDE.md`, viendelezi na seva za MCP kwenye mipangilio inayoshirikiwa, na husoma folda tu.

## ccam inatumika kwenye mifumo gani?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine na mingineyo), macOS na Windows 10 na 11, moja kwa moja au ndani ya WSL. Kwenye Windows, kushiriki faili kama `settings.json` kunahitaji Developer Mode; folda zinazoshirikiwa hufanya kazi bila hiyo, na bila hiyo `ccam add` ya kwanza huunda kikundi `default` chenye seva za MCP na viendelezi pekee ([kushiriki kwenye Windows](/features/share#on-windows)).

## Kwa nini seva ya MCP inaonekana katika akaunti moja tu?

Seva zilizoongezwa ndani ya kipindi kwa `claude mcp add` huhifadhiwa kwa akaunti hiyo tu, na bila `-s user` kwa folda ya sasa tu. ccam hukujulisha wakati ujao utakapoanzisha akaunti hiyo, na `ccam mcp list` na `ccam ext` huzionyesha; hamishia seva kwenye kikundi kinachoshiriki kwa `ccam mcp import <account> <server>`.

## Je, ccam hupunguza kasi ya Claude Code?

Kuanzisha Claude Code kupitia ccam huongeza takriban milisekunde 4, zilizopimwa kwa `scripts/perf-check.sh` kwenye kompyuta mpakato. ccam husoma faili moja ndogo, hukagua viungo vinavyoshirikiwa kisha hujibadilisha yenyewe kuwa `claude`, kwa hiyo haiendeshwi wakati unafanya kazi. Windows haiwezi kubadilisha mchakato, kwa hiyo huko ccam huanzisha Claude Code na kuisubiri, bila kufanya kitu kingine.

## Nitasasishaje ccam?

Endesha `ccam update`. Hupakua toleo jipya zaidi, hukagua checksum yake na hubadilisha faili ya `ccam`, huku `ccam update --check` ikikuambia tu kama kuna toleo jipya zaidi. Dashibodi pia hukuambia kunapokuwa na toleo jipya zaidi, na `U` hulisakinisha baada ya kuuliza; ccam kamwe haijisasishi yenyewe. Ikiwa ulisakinisha kifurushi cha `.deb`, badala yake sakinisha `.deb` mpya kwa apt. Tazama [kusasisha ccam](/docs/install#update).

## Nitaondoaje ccam?

Ondoa njia fupi ya kila akaunti kwa `ccam shortcut <account> --remove`. Ikiwa uliingiza `~/.claude` na iko katika kikundi, kwanza endesha `ccam remove <name>` kwa ajili yake, ili ipate tena nakala ya kile inachoshiriki. Kisha futa `ccam` kutoka `~/.local/bin`, na uondoe `~/.config/ccam`, `~/.local/share/ccam` na `~/.local/state/ccam` ikiwa huhitaji tena akaunti hizo. Tazama [usakinishaji](/docs/install#uninstall).
