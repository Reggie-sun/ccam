---
title: Kubadili kwenda akaunti inayofuata
description: ccam next hufungua Claude Code kama akaunti inayofuata ile uliyotumia mwisho; --auto huchagua matumizi madogo zaidi ya saa 5, --continue huendeleza mazungumzo.
updated: 2026-10-19
order: 8
translation: machine
sourceHash: 92352f532e74
---

`ccam next` hufungua Claude Code kama akaunti nyingine bila wewe kuitaja. Itumie akaunti inapofikia kikomo chake cha matumizi: `ccam next` ya kawaida huchukua akaunti inayofuata ile uliyotumia mwisho, `--auto` huchukua ile yenye nafasi kubwa zaidi katika kipindi chake cha saa 5, na `--continue` huendeleza mazungumzo ya mwisho katika akaunti mpya.

## Matumizi

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next
→ personal → work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude
```

ccam huchapisha ubadilishaji kama `→ <from> → <to>` kwenye standard error, kisha huanzisha Claude Code kama akaunti mpya kwa njia ile ile ambayo `ccam <name>` hufanya. `<from>` ni akaunti uliyotumia mwisho. Hoja zinazokuja baada ya `--` huenda kwa `claude` bila kubadilishwa.

## Chaguo

| Chaguo | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `--auto` | Huchagua akaunti iliyoingia yenye matumizi ya chini zaidi ya saa 5, ikiruka akaunti zilizofikia 100% | imezimwa |
| `-c`, `--continue` | Huzingatia tu akaunti zinazoshiriki historia ya mazungumzo na akaunti uliyotumia mwisho, na huanzisha `claude --continue` | imezimwa |
| `-- <claude args>` | Hoja zinazopitishwa kwa `claude` baada ya zile za ccam yenyewe | hakuna |
| `--dry-run` | Inaheshimiwa: huchapisha `→ <from> → <to>` na mstari wa amri ambao ingeendesha, bila kuanzisha Claude Code wala kurekodi uanzishaji | imezimwa |

`--auto` na `--continue` hufanya kazi pamoja. Kwa `--continue`, `--continue` huwekwa kabla ya hoja zako mwenyewe.

## Jinsi inavyofanya kazi

Kila uanzishaji kupitia ccam huandika muda katika `~/.local/state/ccam/last-used/<name>`. `ccam next` husoma faili hizi, hukagua ni akaunti zipi zimeingia, na huchagua.

1. `cli.candidates` huunda `next.Candidate` moja kwa kila akaunti kutoka `config.toml`: akaunti, muda wake wa mwisho wa kuanzishwa (`launcher.LastUsed`) na kama ina taarifa za kuingia (`creds.Has`: `.credentials.json` inayosomeka, au kwenye macOS kipengee cha Keychain).
2. Bila `--auto`, `next.Rotate` hupanga akaunti kwa namba ya slot na hutafuta ile iliyoanzishwa hivi karibuni zaidi. Hiyo ndiyo akaunti unayoiacha. Ikiwa hakuna akaunti iliyowahi kuanzishwa, ni slot ya chini zaidi.
3. Kutoka hapo huendelea mbele kwa mpangilio wa slot, ikirudi kutoka slot ya mwisho hadi ya kwanza, na hurudisha akaunti ya kwanza iliyoingia.
4. `launcher.LaunchWith` huandaa na kuanzisha akaunti iliyochaguliwa sawasawa kama `ccam <name>` ingefanya: hukagua viungo vya akaunti, hutumia kushiriki kwa MCP, husajili kipindi na hubadilisha ccam kwa `claude` kupitia `syscall.Exec`.

Mzunguko huanza kwenye akaunti iliyotumiwa mwisho na husimama kwenye akaunti ya kwanza iliyoingia baada yake:

<!-- motion:next-rotate -->

`work` ndiyo iliyoanzishwa hivi karibuni zaidi, kwa hiyo mzunguko huanza baada yake. `side` haijaingia na inarukwa, na `lab` huchaguliwa. Kama `lab` nayo ingekuwa haijaingia, mzunguko ungerudi hadi `personal` (`% len` kwenye mchoro). Mzunguko kamwe haurudishi akaunti uliyoanzia: ikiwa hakuna akaunti nyingine iliyoingia, husimama na kosa.

### --auto

`--auto` huamua kwa matumizi badala ya mpangilio wa slot. `cli.pickAuto` hubakiza akaunti zilizoingia, huomba matumizi yao yote kwa wakati mmoja kwa `Fetcher.FetchAll` ile ile ambayo `ccam usage` hutumia, na huiachia `usage.Best()` ichague:

<!-- motion:next-auto -->

`usage.Best()` huchukua asilimia ya chini zaidi ya saa 5 na huruka akaunti zilizofikia 100% au ambazo matumizi yao ya saa 5 hayajulikani. Utafutaji hufuata kanuni za `ccam usage`: namba zenye umri wa chini ya sekunde 60 hutumiwa tena, kila akaunti ina sekunde 3, na akaunti isiyojibu, au ambayo taarifa zake za kuingia zimeisha muda, hupimwa kwa namba zake zilizohifadhiwa. Kipindi kilichohifadhiwa ambacho muda wake wa kuanza upya umepita huhesabiwa kama 0%, ili akaunti ambayo kipindi chake kimeanza upya tangu hapo isirukwe. `--auto` kamwe haifanyi upya tokeni. Tofauti na mzunguko wa kawaida, `--auto` inaweza kuchagua akaunti uliyotumia mwisho ikiwa bado ina nafasi kubwa zaidi.

### --continue

`--continue` ya Claude Code hufungua mazungumzo ya hivi karibuni zaidi katika folda ya sasa, jambo ambalo linafanya kazi tu ikiwa akaunti mpya inaweza kuona mazungumzo hayo. Kwa hiyo kwa `--continue`, ccam huzingatia tu akaunti zilizo katika kikundi kile kile cha `history` kama akaunti uliyotumia mwisho, kisha huweka `--continue` mbele ya hoja zako. Akaunti ambazo historia yao ni `local` hazishiriki na yeyote, kwa hiyo `--continue` kutoka akaunti kama hiyo hushindwa kila mara. Washa kushiriki kwanza kwa `ccam share <name> --history on` (tazama [kushiriki](/docs/sharing)).

## Usalama na hali maalum

- Akaunti ambazo hazijaingia hazichaguliwi kamwe. Ikiwa hakuna akaunti nyingine iliyoingia, `ccam next` ya kawaida husimama na `no logged-in account other than <name>`.
- Kwa `--continue` na hakuna akaunti inayoshiriki historia: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- Kwa `--auto` na hakuna akaunti iliyoingia: `no logged-in account to pick`. Kila akaunti ikiwa imejaa au matumizi yake hayajulikani: `` no account with known usage and headroom (see `ccam usage`) ``.
- `--auto` hufanya maombi ya mtandao hata kwa `--dry-run`. Husoma tu, na kamwe haiandiki faili za tokeni.
- Akaunti iliyochaguliwa huanza kupitia njia ile ile kama `ccam <name>`, kwa hiyo ukaguzi ule ule hutumika, pamoja na onyo wakati Claude Code nyingine inayoshiriki historia ile ile tayari inaendeshwa katika folda hii.
- `last-used` huandikwa tu Claude Code inapoanza kweli. `--dry-run` huiacha bila kuigusa, kwa hiyo jaribio halisogezi mzunguko mbele.

## Zinazohusiana

- [Vikomo vya matumizi](/features/usage)
- [Vikomo vya matumizi na kubadili akaunti](/docs/usage-limits)
- [Shiriki seva za MCP, mipangilio na historia](/docs/sharing)
- [Hatua za kufuata kikomo cha matumizi kinapofikiwa](/guides/usage-limit-workflow)
