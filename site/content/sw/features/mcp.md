---
title: Seva za MCP zinazoshirikiwa
description: Hifadhi orodha moja ya seva za MCP kwa kila kikundi kwa ccam mcp add, add-json, list, remove na import, inayopitishwa kwa --mcp-config au kusawazishwa kwenye kila akaunti.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 46a7aa716f99
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
ccam mcp list                                                     # every group, and servers only one account has
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
```

Kuongeza jina ambalo kikundi tayari kinalo hubadilisha ufafanuzi wake na huchapisha `updated` badala ya `added`. Seva hufikia akaunti tu ikiwa inashiriki MCP na kikundi hicho (`ccam share <name> --mcp on`); ikiwa bado hakuna akaunti inayoshiriki MCP, `add` husema hivyo.

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

- `~/.local/share/ccam/groups/<group>/mcp.json`: seva zinazoshirikiwa, `{"mcpServers": {…}}`, muundo ambao `--mcp-config` husoma. Huandikwa kwa `fsutil.AtomicWrite` na mode `0600`, kwa sababu vipengele vya seva mara nyingi hubeba API keys.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: faili ya Claude Code yenyewe yenye seva za akaunti. ccam husoma tu thamani yake ya `mcpServers` kwa kisomaji cha `profile.TopLevelValue` na kamwe haifasiri faili nzima kwenye njia ya kuanzisha.
- `~/.local/state/ccam/mcp/<name>/`: faili ya `seen`, `sync-stamp` na nakala rudufu za `.claude.json`.

Kuna mikakati miwili, inayowekwa kwa `mcp_strategy` katika `~/.config/ccam/config.toml`. Chaguo-msingi ni mkakati wa flag:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` hupata kikundi cha MCP cha akaunti. `mcp.FlagArgs` hukagua kwa `stat` moja kwamba `mcp.json` ya kikundi ipo, bila kuifasiri.
2. `Plan.InsertArgs` huweka `--mcp-config=<path>` mara tu baada ya `claude`, kabla ya hoja zako mwenyewe, ili `--mcp-config` au `--strict-mcp-config` zako mwenyewe bado zitumike juu yake.
3. `mcp.NewLocal` hukagua kama `.claude.json` ya akaunti imebadilika tangu uanzishaji wa mwisho, na hapo tu ndipo hufasiri seva zake.
4. `syscall.Exec` hubadilisha ccam kwa Claude Code.

Njia hupitishwa kama hoja moja, `--mcp-config=<path>`, si mbili. `--mcp-config` ya Claude Code hupokea thamani kadhaa, kwa hiyo `claude --mcp-config a.json "prompt"` ingesoma prompt kama faili ya pili ya usanidi na kushindwa na `MCP config file not found`. Muundo wa hoja moja husimama hapo, na prompt au amri ndogo kama `mcp list` inayofuata hufanya kazi. Claude Code hupakia seva za `--mcp-config` pamoja na zile za akaunti yenyewe; majina yakigongana, kipengele katika `.claude.json` hushinda. Mkakati wa flag kamwe hauandiki `.claude.json`.

Seva iliyoongezwa ndani ya kipindi kwa `claude mcp add` huingia kwenye `.claude.json` ya akaunti hiyo tu. ccam haiwezi kuiona wakati huo, kwa sababu Claude Code imeichukua nafasi yake, kwa hiyo hukagua wakati wa uanzishaji unaofuata. Faili ya `seen` huhifadhi muda wa kubadilishwa wa `.claude.json` na majina ya seva yaliyokwisha ripotiwa. Muda unapotofautiana, ccam hulinganisha majina na kikundi na huchapisha dokezo kwa kila seva mpya mara moja, bila kusubiri ingizo:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

Kwa `mcp_strategy = "sync"`, seva huonekana katika `claude mcp list` na `/mcp` kama seva za kawaida za mtumiaji, na ccam huziandika katika `.claude.json` ya kila akaunti kabla ya kuanza:

<!-- motion:mcp-sync -->

1. `mcp.Sync` hulinganisha muda wa kubadilishwa wa `mcp.json` ya kikundi na `sync-stamp`. Ikiwa haujabadilika, `.claude.json` ya akaunti haifunguliwi.
2. Huchukua kufuli ya Claude Code yenyewe, folda `.claude.json.lock` kando ya faili, inayoundwa kwa `mkdir`. Husubiri hadi sekunde 2 na huchukulia kufuli yenye umri wa zaidi ya sekunde 10 kuwa imepitwa na wakati, kama maktaba ya kufuli ya Claude Code inavyofanya.
3. Huhifadhi `.claude.json` ya sasa katika `backups/` ndani ya folda ya hali na huhifadhi 5 za karibuni zaidi.
4. Huweka kila seva ya kikundi katika `mcpServers`, huondoa seva ilizoandika awali ambazo kikundi hakina tena, huhifadhi seva za akaunti yenyewe, na huandika faili kwa atomiki ikiwa na mode yake asili.
5. Hurekodi stamp mpya na majina inayoyasimamia, na Claude Code huanza bila `--mcp-config`.

Usawazishaji hurukwa, kwa `ccam: account work is already running, skipping MCP sync`, wakati kipindi cha akaunti kinaendeshwa, kwa sababu Claude Code huandika upya `.claude.json` inapoendeshwa. Usawazishaji ulioshindwa huchapisha `ccam: MCP sync failed (launching anyway): …` na huanzisha hata hivyo. Dokezo la seva zilizoongezwa ndani ya kipindi ni la mkakati wa flag tu.

`ccam mcp import` huhamisha seva za akaunti yenyewe kwenda kwenye kikundi chake:

<!-- motion:mcp-import -->

ccam husoma seva za akaunti kwa `mcp.LocalServers` na hulinganisha kila moja na seva ya kikundi yenye jina lile lile kwa kutumia `mcp.SameDef`, ambayo hupuuza mpangilio wa funguo na nafasi. Seva ambayo kikundi hakina, au kinayo kwa ufafanuzi ule ule, hunakiliwa. Iliyo tofauti hurukwa, kwa sababu ufafanuzi wa kikundi ndio ambao kila akaunti nyingine hutumia:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

Kwa `--force` seva ya kikundi hubadilishwa, na `mcp.json` ya awali kwanza huhifadhiwa katika `.ccam-backup/<time>/mcp.json` ndani ya folda ya kikundi. Kisha `mcp.RemoveLocal` huondoa vipengele vilivyonakiliwa kutoka `.claude.json` ya akaunti, chini ya kufuli ile ile na baada ya nakala rudufu katika folda ya hali. Hili linahitajika kwa sababu kipengele cha akaunti kingeshinda kile cha kikundi majina yakigongana. `--keep-local` huviacha.

## Usalama na hali maalum

- **Taarifa za kuingia.** Seva zinazoingia kwa OAuth huhifadhi tokeni zao katika `.credentials.json` ya kila akaunti (`mcpOAuth`), si kwenye orodha ya seva, kwa hiyo kushiriki seva hakushiriki kuingia kwake. Ingia mara moja kwa kila akaunti kwa `/mcp`, au nakili taarifa za kuingia kutoka `~/.claude` kwa [`ccam migrate`](/features/migrate). Kikundi kamwe hakipokei faili ya kuingia.
- **Akaunti zinazoendeshwa.** `import` hunakili kwenda kwenye kikundi lakini huacha vipengele vya akaunti inayoendeshwa na husema hivyo; iendeshe tena baada ya kipindi kuisha. Usawazishaji huruka akaunti zinazoendeshwa.
- **Faili zilizoharibika.** Ikiwa `.claude.json` si JSON halali, usawazishaji huiacha bila kuigusa na huripoti kosa. Ikiwa kufuli haiwezi kuchukuliwa ndani ya sekunde 2, ccam huripoti `could not lock … (Claude Code is writing the file)`.
- **Akaunti zisizoshiriki MCP.** `import` hukataa na hukuambia uendeshe `ccam share side --mcp on` kwanza.
- **Majina ya vikundi.** `-g` hupokea jina la kikundi tu; huwa folda chini ya `groups/`, kwa hiyo majina kama `../x` hukataliwa.

## Zinazohusiana

- [Kushiriki kati ya akaunti](/features/share)
- [Shiriki seva za MCP, mipangilio na historia](/docs/sharing)
- [Nakili usanidi uliopo](/features/migrate)
- [Amri zote](/docs/commands)
