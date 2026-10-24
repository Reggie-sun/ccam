---
title: Seva za MCP zinazoshirikiwa
description: Orodha moja ya seva za MCP kwa kila kikundi kwa ccam mcp add, add-json, list, remove na import, hupitishwa kwa --mcp-config au husawazishwa kwa kila akaunti.
updated: 2026-10-23
order: 5
translation: machine
sourceHash: "fefca7f8d857"
---

Claude Code huhifadhi seva za MCP za mtumiaji katika `.claude.json`, faili ile ile inayohifadhi taarifa za kuingia za akaunti, kwa hiyo faili hiyo haiwezi kushirikiwa kati ya akaunti. ccam huhifadhi seva zinazoshirikiwa katika faili ya kikundi, `mcp.json`, na huikabidhi kila akaunti ya kikundi wakati wa kuanza. Tumia `ccam mcp` kuongeza seva mara moja kwa akaunti zako zote, na `ccam mcp import` kuhamisha seva uliyoiongeza ndani ya kipindi kimoja kwenda kwa nyingine.

## Matumizi

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

Badiliko hutumika kuanzia uanzishaji unaofuata wa kila akaunti; vipindi vinavyoendeshwa huendelea na seva vilizoanza nazo:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
only in account side, for /home/me/code/app only: notes
```

Chini ya vikundi, `ccam mcp list` hutaja kila seva isiyofika kwa akaunti nyingine: `only in account work: pg` kwa seva ambayo akaunti iliiongeza kwa ajili yake yenyewe, `only in account side, for /home/me/code/app only: notes` kwa seva iliyoongezwa kwa folda moja (`claude mcp add` ya kawaida, tazama hapa chini), `account work has its own …, used in its sessions in place of the group's` wakati seva ya akaunti yenyewe ina jina la seva ya kikundi, `.mcp.json` ya folda hii pamoja na akaunti ambazo idhini yake inasubiriwa (`waiting for approval in`), na seva zilizo `in no account, only in ~/.claude.json (plain claude's)`, pamoja na amri ya [`ccam migrate --mcp`](/features/migrate) inayozinakili.

Kuongeza jina ambalo kikundi tayari kinalo hubadilisha ufafanuzi wake na huchapisha `updated` badala ya `added`. Seva hufikia akaunti tu ikiwa iko katika kikundi hicho na kikundi kinashiriki MCP (`ccam group set <group> --add mcp`); kikundi kisiposhiriki, `add` husema hivyo.

## Chaguo

| Chaguo | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `-g`, `--group <group>` | kikundi lengwa cha `add`, `add-json` na `remove` | kikundi pekee ambacho akaunti fulani hushiriki MCP kupitia kwacho, vinginevyo `default` |
| `-- <command> [args…]` (`add`) | huendesha amri hii kama seva ya stdio | |
| `--url <url>` (`add`) | URL ya seva ya mbali (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | usafirishaji wa seva ya `--url` | `http` |
| `--header <'Name: value'>` (`add`) | kichwa cha HTTP, kinaweza kurudiwa | hakuna |
| `-e`, `--env <KEY=VALUE>` (`add`) | kigezo cha mazingira kwa seva ya stdio, kinaweza kurudiwa | hakuna |
| `--force` (`import`) | hubadilisha seva ya kikundi yenye jina lile lile inayotofautiana; `mcp.json` ya zamani huhifadhiwa nakala rudufu | imezimwa |
| `--keep-local` (`import`) | huhifadhi vipengele katika `.claude.json` ya akaunti | imezimwa |
| `-h`, `--help` | msaada wa kila amri | |
| `--dry-run` (jumla) | huheshimiwa na `mcp list` tu; `add`, `add-json`, `remove` na `import` huikataa na hazibadilishi chochote | imezimwa |

`add` inahitaji amri baada ya `--` au `--url`; vinginevyo hushindwa na `a command to run (after --) or --url is required`. `--header` hutumika kwa seva za `--url` na `--env` kwa seva za amri. `add-json` hupokea JSON ile ile kama `claude mcp add-json`, na ufafanuzi lazima uwe kitu chenye `command` au `url`. `import` daima hutumia kikundi ambacho akaunti hushiriki MCP kupitia kwacho, na `list` huonyesha kila kikundi, kwa hiyo `-g` haizibadilishi.

## Jinsi inavyofanya kazi

Faili zinazohusika:

- `~/.local/share/ccam/groups/<group>/mcp.json`: seva zinazoshirikiwa, `{"mcpServers": {…}}`, muundo ambao `--mcp-config` husoma. Huandikwa kwa `fsutil.AtomicWrite` na mode `0600`, kwa sababu vipengele vya seva mara nyingi hubeba API keys, na tu chini ya kufuli `<group>/.mcp.json.lock`, ili `ccam mcp add` mbili kwa wakati mmoja zote zihifadhi seva zao.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: faili ya Claude Code yenyewe yenye seva za akaunti. ccam husoma tu thamani yake ya `mcpServers` kwa kisomaji cha `profile.TopLevelValue` na kamwe haifasiri faili nzima kwenye njia ya kuanzisha.
- `~/.local/state/ccam/mcp/<name>/`: faili ya `seen`, `sync-stamp` (pamoja na hash ya kila nakala ambayo ccam iliandika, ili iweze kutofautisha nakala ambayo hakuna aliyeibadilisha na ile uliyoihariri), `sync-checked` (kile ambacho mkakati wa sync uliona mara ya mwisho), `kept` (seva ambazo akaunti ilihifadhi ilipoacha kushiriki MCP) na nakala rudufu za `.claude.json`.

Kuna mikakati miwili, inayowekwa kwa `mcp_strategy` katika `~/.config/ccam/config.toml`. Chaguo-msingi ni mkakati wa flag:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` hupata kikundi cha MCP cha akaunti. `mcp.FlagArgs` hukagua kwa `stat` moja kwamba `mcp.json` ya kikundi ipo, bila kuifasiri.
2. `Plan.InsertArgs` huweka `--mcp-config=<path>` mara tu baada ya `claude`, kabla ya hoja zako mwenyewe, ili `--mcp-config` au `--strict-mcp-config` zako mwenyewe bado zitumike juu yake.
3. `mcp.NewLocal` hukagua kama `.claude.json` ya akaunti imebadilika tangu uanzishaji wa mwisho, na hapo tu ndipo hufasiri seva zake.
4. `syscall.Exec` hubadilisha ccam kwa Claude Code.

Njia hupitishwa kama hoja moja, `--mcp-config=<path>`, si mbili. `--mcp-config` ya Claude Code hupokea thamani kadhaa, kwa hiyo `claude --mcp-config a.json "prompt"` ingesoma prompt kama faili ya pili ya usanidi na kushindwa na `MCP config file not found`. Muundo wa hoja moja husimama hapo, na prompt au amri ndogo inayofuata hufanya kazi. `ccam work mcp list`, `get` na `remove` hazipati flag, kwa kuwa Claude Code huikataa huko; badala yake ccam hutaja seva za kikundi kwenye standard error ([maelezo](/features/launch#shared-mcp-servers)). Claude Code hupakia seva za `--mcp-config` pamoja na zile za akaunti yenyewe; majina yakigongana, kipengele katika `.claude.json` hushinda. Mkakati wa flag huandika `.claude.json` tu kwa ajili ya Remote Control, akaunti inapoacha kushiriki MCP au inapojiunga tena ([hapa chini](#leaving-and-joining)), na kwa `ccam mcp import`.

`claude remote-control` (`rc`) haipokei `--mcp-config`, na vipindi vyake husoma seva kutoka `.claude.json` tu. Kwa `ccam work rc`, ccam huacha chaguo hilo na, kwenye terminal, huuliza kama inakili seva za kikundi kwenda `.claude.json`, pamoja na nakala rudufu, kama mkakati wa `sync` unavyofanya. Huuliza hadi nakala iwepo, kisha tena tu baada ya `mcp.json` ya kikundi kubadilika au seva iliyonakiliwa kutoweka ([Remote Control](/features/launch#remote-control)). Kwa sababu nakala kama hiyo hushinda `--mcp-config` majina yanapogongana, `ccam mcp add`, `add-json`, `remove` na `import` husasisha nakala zinazoshikiliwa na akaunti nyingine za kikundi:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

Mstari wa `✓ Updated` huonekana tu wakati nakala imebadilika. Nakala iliyotengenezwa na ccam ya zamani husasishwa kwa njia ile ile, vivyo hivyo nakala iliyopo baada ya `ccam migrate`.

Seva iliyoongezwa ndani ya kipindi kwa `claude mcp add` huingia kwenye `.claude.json` ya akaunti hiyo tu. `claude mcp add` ya kawaida hutumia local scope: seva hupakiwa kwa akaunti hiyo moja katika folda hiyo moja (`-s user` huifanya iwe ya akaunti katika kila folda). ccam haiwezi kuiona wakati huo, kwa sababu Claude Code imeichukua nafasi yake, kwa hiyo hukagua wakati wa uanzishaji unaofuata. Faili ya `seen` huhifadhi muda wa kubadilishwa wa `.claude.json` na majina ya seva yaliyokwisha ripotiwa. Muda unapotofautiana, ccam hulinganisha majina na kikundi na huchapisha dokezo kwa kila seva mpya mara moja, bila kusubiri ingizo, kwa mikakati yote miwili:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

Dokezo pia hutaja seva ya akaunti yenyewe yenye jina la seva ya kikundi, kwa kuwa ile ya akaunti hushinda katika vipindi vyake.

Kwa `mcp_strategy = "sync"`, seva huonekana katika `claude mcp list` na `/mcp` kama seva za kawaida za mtumiaji, na ccam huziandika katika `.claude.json` ya kila akaunti kabla ya kuanza:

<!-- motion:mcp-sync -->

1. `mcp.Sync` hulinganisha muda wa kubadilishwa wa `mcp.json` ya kikundi na `sync-stamp`, na `.claude.json` ya akaunti na kile ilichoona mara ya mwisho. Ikiwa hakuna kilichobadilika, hakuna kingine kinachofanyika. Vinginevyo pia hukagua kwamba kila seva ya kikundi bado iko katika `.claude.json`, na hurudisha nakala iliyotoweka (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`).
2. Huchukua kufuli ya Claude Code yenyewe, folda `.claude.json.lock` kando ya faili, inayoundwa kwa `mkdir`. Husubiri hadi sekunde 2 na huchukulia kufuli yenye umri wa zaidi ya sekunde 10 kuwa imepitwa na wakati, kama maktaba ya kufuli ya Claude Code inavyofanya.
3. Huhifadhi `.claude.json` ya sasa katika `backups/` ndani ya folda ya hali na huhifadhi 5 za karibuni zaidi.
4. Huweka kila seva ya kikundi katika `mcpServers`, huondoa seva ilizoandika awali ambazo kikundi hakina tena, na huandika faili kwa atomiki ikiwa na mode yake asili, tu wakati kitu kimebadilika. Hubadilisha au huondoa tu nakala ilizoandika ambazo hakuna aliyezibadilisha tangu hapo. Seva ambayo akaunti inayo yenyewe kwa jina la seva ya kikundi hubaki kuwa ya akaunti, na ccam husema hivyo mara moja: `account work has its own MCP servers fs, so group default's of those names are not copied there`.
5. Hurekodi stamp mpya pamoja na hash ya kila nakala, na Claude Code huanza bila `--mcp-config`.

Usawazishaji hurukwa, kwa `ccam: account work is already running, skipping MCP sync`, wakati kipindi cha akaunti kinaendeshwa, bila kujali Claude Code ilianzishwa vipi, kwa sababu Claude Code huandika upya `.claude.json` inapoendeshwa. Usawazishaji ulioshindwa huchapisha `ccam: MCP sync failed (launching anyway): …` na huanzisha hata hivyo.

### Kuacha na kujiunga

Akaunti inayoacha kushiriki MCP (`ccam group leave`, `ccam group set --remove mcp`, `ccam group remove --leave`) huhifadhi seva za kikundi: mpango huzinakili katika `.claude.json` yake, pamoja na nakala rudufu, na seva zake yenyewe zenye majina hayo hubaki. Kwa `--fresh` haipati yoyote, na mpango huorodhesha seva zitakazoacha kupakiwa. Inapojiunga tena na kikundi chenye seva hizo, nakala zilizohifadhiwa ambazo hakuna aliyezibadilisha huondolewa, ili mabadiliko ya baadaye ya kikundi yaifikie.

`ccam mcp import` huhamisha seva za akaunti yenyewe kwenda kwenye kikundi chake:

<!-- motion:mcp-import -->

ccam husoma seva za akaunti, zake yenyewe na zile ilizo nazo kwa folda moja tu, na hulinganisha kila moja na seva ya kikundi yenye jina lile lile kwa kutumia `mcp.SameDef`, ambayo hupuuza mpangilio wa funguo na nafasi. Seva ya folda huwa seva ya kikundi kwa kila folda (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`); jina moja lenye fasili mbili tofauti za folda hukataliwa. Seva za plugins hazitolewi, kwa kuwa plugin huzileta yenyewe. Seva ambayo kikundi hakina, au kinayo kwa ufafanuzi ule ule, hunakiliwa. Iliyo tofauti hurukwa, kwa sababu ufafanuzi wa kikundi ndio ambao kila akaunti nyingine hutumia:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

Kwa `--force` seva ya kikundi hubadilishwa, na `mcp.json` ya awali kwanza huhifadhiwa katika `.ccam-backup/<time>/mcp.json` ndani ya folda ya kikundi. Kisha `mcp.Remove` huondoa vipengele vilivyonakiliwa kutoka `.claude.json` ya akaunti, vya akaunti yenyewe na vya folda kwa uandishi mmoja, chini ya kufuli ile ile na baada ya nakala rudufu katika folda ya hali. `mcp.json` ya kikundi hubaki imefungwa kwa kufuli tangu inapopakiwa hadi vipengele hivyo viondoke, ili hakuna anayeondoa seva hiyo kutoka kwenye kikundi katikati. Hili linahitajika kwa sababu kipengele cha akaunti kingeshinda kile cha kikundi majina yakigongana. `--keep-local` huviacha.

## Usalama na hali maalum

- **Taarifa za kuingia.** Seva zinazoingia kwa OAuth huhifadhi tokeni zao katika `.credentials.json` ya kila akaunti (`mcpOAuth`), si kwenye orodha ya seva, kwa hiyo kushiriki seva hakushiriki kuingia kwake. Ingia mara moja kwa kila akaunti kwa `/mcp`. ccam kamwe hainakili kuingia, hata [`ccam migrate`](/features/migrate) haifanyi hivyo: refresh token iliyonakiliwa huwatoa nje ya akaunti wengine wote wenye nakala yake. [`ccam ext`](/features/ext) na `ccam doctor` hutaja seva ambazo akaunti haijaingia wakati akaunti nyingine imeingia, na `ccam next --continue` husema ni zipi ambazo akaunti mpya inakosa. Kikundi kamwe hakipokei faili ya kuingia.
- **Akaunti zinazoendeshwa.** `import` hunakili kwenda kwenye kikundi lakini huacha vipengele vya akaunti inayoendeshwa na husema hivyo; iendeshe tena baada ya kipindi kuisha. Usawazishaji huruka akaunti zinazoendeshwa. Zote mbili huhesabu pia Claude Code iliyoanzishwa bila ccam.
- **Nakala kwa ajili ya Remote Control.** Nakala ya akaunti inayoendeshwa haiguswi: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Kwa `mcp_strategy = "sync"` hakuna kinachohitaji kusasishwa, kwa sababu kila uanzishaji husawazisha.
- **Faili zilizoharibika.** Ikiwa `.claude.json` si JSON halali, usawazishaji huiacha bila kuigusa na huripoti kosa. Ikiwa kufuli haiwezi kuchukuliwa ndani ya sekunde 2, ccam huripoti `could not lock … (Claude Code is writing the file)`.
- **Akaunti zisizoshiriki MCP.** `import` hukataa na hukuambia uweke akaunti kwanza katika kikundi kinachoshiriki MCP (`ccam group join <group> side`).
- **Majina ya vikundi.** `-g` hupokea jina la kikundi tu; huwa folda chini ya `groups/`, kwa hiyo majina kama `../x` hukataliwa.

## Zinazohusiana

- [Kushiriki kati ya akaunti](/features/share)
- [Kutumia Remote Control na ccam](/guides/remote-control)
- [Shiriki seva za MCP, mipangilio na historia](/docs/sharing)
- [Nakili usanidi uliopo](/features/migrate)
- [Amri zote](/docs/commands)
