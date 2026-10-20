---
title: Kubadili kwenda akaunti inayofuata
description: ccam next hufungua Claude Code kama akaunti nyingine; --auto huchagua nafasi kubwa zaidi, --continue huhamisha mazungumzo, --fallback hurudia kazi za -p.
updated: 2026-10-20
order: 8
translation: machine
sourceHash: "6a5dcb1d5c01"
---

`ccam next` hufungua Claude Code kama akaunti nyingine bila wewe kuitaja. Itumie akaunti inapofikia kikomo chake cha matumizi: `ccam next` ya kawaida huchukua akaunti inayofuata ile uliyotumia mwisho, `--auto` huchukua ile yenye nafasi kubwa zaidi, na `--continue` huhamisha mazungumzo ya mwisho ya folda hii kwenda akaunti mpya, iwe hizo mbili zinashiriki historia au la. Ndani ya Claude Code, [plugin ya ccam](#inside-claude-code) hukujulisha kikomo kinaposimamisha kipindi na `/ccam:switch` hukabidhi mazungumzo. Kwa script zinazoendesha `claude -p`, `ccam --fallback` hufanya ubadilishaji yenyewe (tazama [Endesha tena kwenye kikomo](#run-again-at-the-limit)).

## Matumizi

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the weekly quota that would go unused first
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam huchapisha ubadilishaji kama `→ <from> → <to>` kwenye standard error, pamoja na sababu ya chaguo la `--auto`, kisha huanzisha Claude Code kama akaunti mpya kwa njia ile ile ambayo `ccam <name>` hufanya. `<from>` ni akaunti uliyotumia mwisho, au kwa `--continue` akaunti ambayo mazungumzo ni yake. Hoja zinazokuja baada ya `--` huenda kwa `claude` bila kubadilishwa.

## Chaguo

| Chaguo | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `--auto` | Huchagua akaunti iliyoingia kwa matumizi (tazama hapa chini), ikiruka akaunti zilizofikia 100% | imezimwa |
| `--max-usage <1-100>` | Pamoja na `--auto`: hupendelea akaunti zilizo chini ya matumizi haya | `auto_max_usage`, vinginevyo 100 |
| `--strategy headroom\|expiring\|balance` | Pamoja na `--auto`: matumizi ya chini zaidi ya saa 5, kipindi cha saa 5 kinachoanza upya mapema zaidi, au kiasi cha wiki ambacho kingebaki bila kutumika kwanza | `auto_strategy`, vinginevyo `headroom` |
| `--cross-org` | Pamoja na `--auto`: chagua pia akaunti za shirika jingine (tazama hapa chini) | `cross_org`, vinginevyo imezimwa |
| `--wait[=<duration>]` | Pamoja na `--auto`: hakuna akaunti yenye nafasi, husubiri ya kwanza ianze upya, si zaidi ya muda huu | imezimwa; `5h` ikitolewa peke yake |
| `-c`, `--continue` | Huhamisha mazungumzo ya mwisho ya folda hii kwenda akaunti mpya (tazama hapa chini) | imezimwa |
| `-- <claude args>` | Hoja zinazopitishwa kwa `claude` baada ya zile za ccam yenyewe | hakuna |
| `--dry-run` | Inaheshimiwa: huchapisha `→ <from> → <to>` na mstari wa amri ambao ingeendesha, bila kuanzisha Claude Code wala kurekodi uanzishaji | imezimwa |

`--auto` na `--continue` hufanya kazi pamoja. Kwa `--continue`, hoja za `--resume` huwekwa kabla ya hoja zako mwenyewe.

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

`--auto` huamua kwa matumizi badala ya mpangilio wa slot. `cli.pickAuto` hubakiza akaunti zilizoingia, huomba matumizi yao yote kwa wakati mmoja kwa `Fetcher.FetchAll` ile ile ambayo `ccam usage` hutumia, na huiachia `usage.Policy.Choose()` ichague:

<!-- motion:next-auto -->

Akaunti hupimwa kwa vipindi vyake vya saa 5 na vya wiki, na, `--model` ikiwa miongoni mwa hoja za Claude Code, pia kwa kipindi cha wiki cha modeli hiyo (`opus`, `sonnet`, `haiku` au `fable` katika jina). `Policy.Choose()` hupanga akaunti katika makundi matatu:

1. **Nafasi**: kila kipindi kiko chini ya `auto_max_usage` (`config.toml`, chaguo-msingi 100). Mkakati huchagua kati yao: `headroom` (chaguo-msingi) huchukua matumizi ya chini zaidi ya saa 5, `expiring` huchukua kipindi cha saa 5 kitakachoanza upya mapema zaidi, ili kiasi chake kilichobaki kitumike kabla hakijapotea; `balance` huchukua akaunti yenye kiasi kikubwa zaidi cha wiki kilichobaki kwa kila siku hadi wiki yake ianze upya, ikihesabu kila kipindi cha ccam kinachoendeshwa dhidi yake, ili kiasi cha wiki kisibaki bila kutumika na vipindi sambamba vienee.
2. **Juu**: chini ya 100% lakini kwenye au juu ya `auto_max_usage`. Hutumiwa tu wakati hakuna akaunti yenye nafasi, ile yenye matumizi ya chini zaidi ya saa 5 kwanza.
3. **Imejaa**: kipindi kiko 100%, au matumizi ya saa 5 hayajulikani. Haichaguliwi kamwe.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Mashirika.** Chaguo za kiotomatiki hubaki kwa akaunti zinazofanana na ile unayoiacha: shirika lile lile la team au enterprise, au mipango ya binafsi (Pro, Max) pamoja, ili kazi kutoka akaunti ya kampuni isiendelee kwenye akaunti ya binafsi au kinyume chake. ccam hujua shirika la kila akaunti kutokana na kuingia kwake (`ccam login` hulirekodi; kuingia kwa zamani husomwa kutoka `.claude.json`); akaunti ambayo haiwezi kuitambua hufaa popote. Folda mbili zilizoingia kwenye akaunti ile ile ya Claude hushiriki vikomo vyake, kwa hiyo moja kamwe haichukui nafasi ya nyingine. `--cross-org`, au `cross_org = true` katika `config.toml`, huondoa kanuni ya shirika.

Isipowezekana kuchagua chochote, ccam hutoka na msimbo 75 (`EX_TEMPFAIL`), ili script iweze kutofautisha "jaribu baadaye" na kushindwa. Kwa `--wait`, badala yake hulala hadi kipindi kinachozuia akaunti ya kwanza iliyojaa kianze upya (sekunde 30 baada ya muda wa kuanza upya, endpoint inapoonyesha kipindi kipya), huchapisha `ccam: no account has room; waiting until 15:10`, na huchagua tena; ikiwa hakuna akaunti inayosema itaanza upya lini, huangalia tena baada ya dakika 10. Ikiwa hakuna kinachopata nafasi kabla ya kusubiri kuisha, hutoka na 75. Utafutaji hufuata kanuni za `ccam usage`: namba zenye umri wa chini ya sekunde 60 hutumiwa tena, kila akaunti ina sekunde 3, na akaunti isiyojibu, au ambayo taarifa zake za kuingia zimeisha muda, hupimwa kwa namba zake zilizohifadhiwa. Kipindi kilichohifadhiwa ambacho muda wake wa kuanza upya umepita huhesabiwa kama 0%, ili akaunti ambayo kipindi chake kimeanza upya tangu hapo isirukwe. `--auto` kamwe haifanyi upya tokeni. Tofauti na mzunguko wa kawaida, `--auto` inaweza kuchagua akaunti uliyotumia mwisho ikiwa bado ina nafasi kubwa zaidi.

### --continue

1. ccam hutafuta mazungumzo ya mwisho ya folda hii katika kila akaunti: `<session>.jsonl` mpya zaidi katika `projects/<folder>/` (jina la folda ni folda ya kazi huku kila herufi isiyo herufi ya alfabeti wala tarakimu ikigeuzwa kuwa `-`, kama Claude Code inavyoipa jina). Ikiwa kipindi kiliacha maelezo (kikomo cha matumizi au `/ccam:switch`, tazama hapa chini) na hakuna kipya zaidi kilichoandikwa tangu hapo, kipindi hicho hasa ndicho hutumiwa badala yake.
2. Akaunti ambayo mazungumzo ni yake ndiyo unayobadili kutoka kwayo, na akaunti mpya huchaguliwa kama ilivyo hapo juu, ndani ya shirika lile lile.
3. Akaunti mpya ikishiriki historia hiyo (folda ile ile ya kikundi), ccam huanzisha `claude --resume <session id>`, yale mazungumzo hasa (`--continue` ingefungua yoyote yaliyo mapya zaidi katika folda).
4. Vinginevyo huanzisha `claude --resume <the conversation's file> --fork-session`: Claude Code hunakili mazungumzo yote, chini ya session id mpya, kwenye historia ya akaunti mpya yenyewe, na ya asili hubaki kama ilivyokuwa. Kuanzia hapo `--continue` ya akaunti mpya huyapata.

ccam husema ukubwa wa mazungumzo: akaunti mpya huyasoma yote tena kwa jibu lake la kwanza. Prompt caches kamwe hazishirikiwi kati ya mashirika, kwa hiyo kwa `--cross-org` jibu hilo la kwanza hugharimu ukubwa kamili. Kushiriki `history` katika kikundi huepuka kunakili, huhifadhi mazungumzo mamoja, na pia hushiriki checkpoints ambazo `/rewind` hurejesha msimbo kutoka kwazo (`file-history/`).

## Ndani ya Claude Code

[Plugin ya ccam](/docs/install#claude-code-plugin) huongeza hooks mbili (Claude Code 2.1.78 au baadaye):

- Zamu inaposimama kwenye kikomo cha matumizi (`You've hit your session limit`, kikomo cha wiki au cha modeli), hurekodi mazungumzo na huonyesha arifa ya desktop katika terminal zinazoiunga mkono (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): `ccam: work hit its 5-hour limit (resets 3:45pm). ccam next --continue carries on as another account.` Kizuizi kifupi cha seva au kikomo cha kasi cha API key si kikomo cha matumizi na huachwa. Claude Code 2.1.234 na baadaye pia husubiri ndani ya kipindi na huendelea yenyewe baada ya kipindi cha saa 5 kuanza upya; arifa hukupa chaguo la kubadili badala yake.
- `/ccam:switch [account]` hurekodi mazungumzo na humaliza zamu kabla ombi lolote halijatumwa: andika `/exit`, kisha `ccam next --continue` huyaendeleza, kama akaunti uliyoitaja au kama inayofuata.

Maelezo hayo huhifadhiwa katika `~/.local/state/ccam/handoff/<account>.json` kwa saa 12 na hutumiwa mara moja, kipindi kijacho kinapokaribia kuanza; uanzishaji unaoshindwa huyaacha kwa jaribio linalofuata. Ni ya folda ambamo kipindi kilianzia, hata baada ya `cd` katika shell ya Claude Code. Ikiwa akaunti inayotajwa ndani yake imebadilishwa jina au kuondolewa tangu hapo, ccam husema hivyo na huchagua kama kawaida. Kwa `relay = true` katika `config.toml`, ccam hubaki kando ya vipindi shirikishi badala ya kujibadilisha na Claude Code, na baada ya `/ccam:switch` huendeleza mazungumzo ya kipindi hicho yenyewe (maelezo yaliyoachwa na kipindi kingine hayafanyiwi kazi). Hilo hugharimu mchakato mmoja wa ccam unaosubiri kwa kila kipindi, ndiyo sababu umezimwa kwa chaguo-msingi. Kipindi kinachoendelezwa huanza kwa `ccam next --continue`, bila hoja ambazo cha kwanza kilianzishwa nazo, kama vile `--model`.

## Endesha tena kwenye kikomo

`ccam next` ni kwa kipindi shirikishi kilichosimama. Uendeshaji usio shirikishi, `claude -p` katika script au kazi iliyopangwa, hauwezi kukusubiri, kwa hiyo `--fallback` kabla ya jina la akaunti huiruhusu ccam ibadili kwa ajili yake:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` hufanya kazi tu pamoja na `-p` au `--print`; vinginevyo ccam husimama na huelekeza kwa `ccam next --auto --continue`.
2. Standard input husomwa mara moja na hupewa kila jaribio.
3. ccam huanzisha `claude` kama mchakato mtoto badala ya `exec`, huzuia standard output yake, na hupitisha standard error kama inavyokuja.
4. `claude` inapotoka kwa 0, au inaposhindwa kwa sababu nyingine, output iliyozuiwa na exit code hupitishwa bila kubadilishwa.
5. Kushindwa huhesabiwa kama kikomo cha matumizi ujumbe wa Claude Code unaposema hivyo (`You've hit your limit`, `usage limit reached`), au utafutaji wa matumizi ya akaunti unapoonyesha kipindi chake cha saa 5 au cha wiki kiko 100%.
6. Kisha output ya jaribio hilo hutupwa na amri ile ile huendeshwa kama akaunti nyingine iliyoingia: ile yenye matumizi ya chini zaidi ya saa 5, kisha akaunti ambazo matumizi yao hayawezi kupatikana, kwa mpangilio wa namba. Akaunti zilizofikia 100% hurukwa, na kila akaunti hujaribiwa si zaidi ya mara moja.
7. Kukiwa na `-c`, `--continue`, `-r` au `--resume` miongoni mwa hoja, akaunti zilizo katika kikundi kile kile cha historia pekee ndizo hujaribiwa, ili mazungumzo yaweze kupatikana.
8. Akaunti za shirika lile lile pekee ndizo hujaribiwa, kama ilivyo kwa `--auto` (`cross_org = true` huondoa hili).

Akaunti inayofuata huchaguliwa kwa kanuni zile zile kama `--auto` (`auto_max_usage`, `auto_strategy`, kipindi cha `--model`). Hakuna akaunti inapobaki, ccam huchapisha output ya jaribio la mwisho, husema kwamba hakuna akaunti iliyoingia yenye nafasi iliyobaki, na hutoka kwa code 75. Kwa `--wait` kabla ya jina la akaunti (`ccam --fallback --wait=2h work -p …`), badala yake husubiri akaunti ya kwanza ianze upya, kisha inaweza kujaribu kila akaunti tena; kila mara hulala angalau mara moja kwanza, ili akaunti iliyofikia kikomo chake sasa hivi isiendeshwe tena mara moja. `ccam --dry-run --fallback work -p …` huchapisha amri ya jaribio la kwanza.

Vipindi shirikishi havibadilishwi vyenyewe. Badala yake, kwa `ccam statusline install`, mstari wa hali huongeza `→ ccam next --auto --continue` matumizi ya saa 5 au ya wiki yanapofikia 90% ([Vikomo vya matumizi](/features/usage)).

## Usalama na hali maalum

- Akaunti ambazo hazijaingia hazichaguliwi kamwe. Ikiwa hakuna akaunti nyingine iliyoingia, `ccam next` ya kawaida husimama na `no logged-in account other than <name>`.
- Kwa `--continue` na hakuna mazungumzo katika folda: `no account has a conversation in <folder> to continue`.
- Kanuni ya shirika isipoacha akaunti yoyote: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Kwa `--auto` na hakuna akaunti iliyoingia: `no logged-in account to pick`. Kila akaunti ikiwa imejaa au matumizi yake hayajulikani: `` no account with known usage and headroom (see `ccam usage`) ``. Zote mbili huishia na exit code 75.
- `--max-usage`, `--strategy`, `--wait` na `--cross-org` hufanya kazi tu pamoja na `--auto`, na thamani nje ya 1-100 au mkakati mwingine hukataliwa.
- Fork hunakili mazungumzo, si hali nyingine ya akaunti: checkpoints za `file-history/` ya akaunti ya asili haziji isipokuwa historia inashirikiwa.
- `--auto` hufanya maombi ya mtandao hata kwa `--dry-run`. Husoma tu, na kamwe haiandiki faili za tokeni.
- Akaunti iliyochaguliwa huanza kupitia njia ile ile kama `ccam <name>`, kwa hiyo ukaguzi ule ule hutumika, pamoja na onyo wakati Claude Code nyingine inayoshiriki historia ile ile tayari inaendeshwa katika folda hii.
- `last-used` huandikwa tu Claude Code inapoanza kweli. `--dry-run` huiacha bila kuigusa, kwa hiyo jaribio halisogezi mzunguko mbele.
- `--fallback` huzuia standard output hadi jaribio liishe, kwa hiyo output ya mtiririko kama `--output-format stream-json` hufika kwa kipande kimoja mwishoni. Faili ambazo jaribio la kwanza lilibadilisha kabla ya kufikia kikomo hubaki zimebadilishwa; jaribio linalofuata huanza prompt upya.

## Zinazohusiana

- [Vikomo vya matumizi](/features/usage)
- [Vikomo vya matumizi na kubadili akaunti](/docs/usage-limits)
- [Shiriki seva za MCP, mipangilio na historia](/docs/sharing)
- [Hatua za kufuata kikomo cha matumizi kinapofikiwa](/guides/usage-limit-workflow)
