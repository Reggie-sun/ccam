---
title: Kutumia Remote Control na ccam
description: "Remote Control ya Claude Code kama akaunti yoyote kwa ccam <name> rc: chaguo zinazotumika, kwa nini ccam hutoa kunakili seva za MCP zinazoshirikiwa, na utatuzi."
updated: 2026-10-23
order: 4
translation: machine
sourceHash: "3e17d89f98c8"
---

`ccam work rc` huanzisha Remote Control ya Claude Code kama akaunti `work`, ili uweze kuendesha vipindi kwenye kompyuta hii kutoka claude.ai/code au programu ya simu ya Claude. Hufanya kazi kama `claude rc` kwa tofauti moja tu: ikiwa akaunti inashiriki seva za MCP, ccam huuliza mara moja kama izinakili kwenda `.claude.json` ya akaunti, kwa sababu vipindi vya Remote Control havisomi seva kutoka mahali pengine popote.

## Kuianzisha

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

[Amri ya njia fupi](/features/shortcuts) hufanya kazi vivyo hivyo: `ccwork rc` ni `ccam work rc`. Endesha amri katika folda unayotaka kufanyia kazi, kisha fungua claude.ai/code au programu ya Claude, ukiwa umeingia kwa akaunti ile ile ya Claude kama `work`.

## Chaguo

Aina mbili za chaguo hukutana hapa: chaguo za ccam yenyewe, ambazo huja kabla ya jina la akaunti au huwa katika mipangilio ya ccam, na chaguo za Remote Control, ambazo huja baada ya `rc` na hufika kwa Claude Code bila kubadilishwa.

### Chaguo na mipangilio ya ccam

| Chaguo au mpangilio | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `--dry-run` (kabla ya akaunti) | Huchapisha mazingira na amri, pamoja na taarifa ya MCP ikiwa ipo. Haiulizi wala hainakili | imezimwa |
| `Copy them into .claude.json (a backup is kept; servers of the account's own stay)? [y/N]` | Huulizwa wakati akaunti inashiriki seva za MCP na `.claude.json` yake haina orodha ya sasa ya kikundi, au seva iliyonakiliwa imetoweka ndani yake. `y` huzinakili, `N` au Enter huanza bila hizo | `N` |
| `mcp_strategy = "sync"` juu kabisa ya `~/.config/ccam/config.toml` | Hunakili seva za kikundi kwenda `.claude.json` kabla ya kila uanzishaji, iwe Remote Control au la, kwa hiyo hakuna kinachoulizwa kamwe | `"flag"` |
| `ccam group leave <group> work` | Hutoa akaunti kwenye kikundi kinachoshiriki MCP nayo. Seva za kikundi hunakiliwa katika `.claude.json` yake kama seva zake yenyewe, kwa hiyo Remote Control huwa nazo na hakuna kinachoulizwa; kwa `--fresh` huanza na seva za akaunti yenyewe tu | katika `default` |

### Chaguo za Remote Control

Kila kitu baada ya `rc` ni mali ya Claude Code. Hizi ni chaguo za `claude remote-control` katika Claude Code 2.1.304; `ccam work rc --help` huchapisha orodha ya toleo lako.

| Chaguo | Inachofanya |
| --- | --- |
| `--name <name>` | Jina la kipindi linaloonyeshwa katika claude.ai/code |
| `--remote-control-session-name-prefix <prefix>` | Kiambishi awali cha majina ya vipindi yanayotengenezwa (chaguo-msingi: jina la mwenyeji) |
| `-c`, `--continue` | Huunganisha tena kwenye kipindi ambacho Remote Control ilikirekodi mwisho kwa folda hii, ndani ya takriban saa 4 |
| `--session-id <id>` | Huunganisha tena kwenye kipindi mahususi |
| `--permission-mode <mode>` | Hali ya ruhusa kwa vipindi vipya: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` au `plan` |
| `--spawn <mode>` | `same-dir` (chaguo-msingi), `worktree` (git worktree moja kwa kila kipindi) au `session` (kipindi kimoja, hutoka kinapoisha) |
| `--capacity <N>` | Idadi kubwa zaidi ya vipindi kwa wakati mmoja katika hali ya `same-dir` na `worktree` (chaguo-msingi 32) |
| `--[no-]create-session-in-dir` | Huunda kipindi cha kwanza katika folda ya sasa wakati wa kuanza (chaguo-msingi imewashwa) |
| `--[no-]chrome` | Claude in Chrome kwa vipindi vipya (chaguo-msingi imezimwa) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | Matokeo ya utatuzi na ya kina |

ccam husoma amri ya Claude Code kwa njia ile ile ambayo Claude Code hufanya, ikiruka chaguo pamoja na thamani zake, kwa hiyo `ccam work rc --name laptop` na `ccam work --model opus rc` zote huanzisha Remote Control. Chaguo ambalo ccam hailijui huhesabiwa kama lisilo na thamani; ili kuwa salama, weka `rc` mara moja baada ya jina la akaunti.

## ccam inachofanya

Uanzishaji hufuata njia ya kawaida ([Kuanzisha Claude Code](/features/launch)) ukiwa na badiliko moja katika hatua ya MCP:

1. ccam hupata akaunti, hukagua viungo vyake vya kushiriki na hutafuta `claude`, kama ilivyo kwa uanzishaji wowote.
2. Huona kwamba amri ni `rc` au `remote-control` na huacha `--mcp-config` nje.
3. Ikiwa akaunti inashiriki seva za MCP, hulinganisha `mcp.json` ya kikundi na `sync-stamp` iliyo katika `~/.local/state/ccam/mcp/work/`, na hukagua kwamba kila seva ya kikundi bado iko katika `.claude.json`. Ikiwa `.claude.json` tayari ina seva za sasa, au kikundi hakina seva yoyote, huenda moja kwa moja hadi hatua ya 6.
4. Vinginevyo hutaja seva ambazo vipindi vingezikosa:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; servers of the account's own stay)? [y/N]
   ```

5. Kwa `y` huhifadhi nakala rudufu ya `.claude.json`, huchukua kufuli ya Claude Code yenyewe, huongeza seva za kikundi kwenye `mcpServers` na huhifadhi seva nyingine za akaunti; seva ya akaunti yenyewe yenye jina la seva ya kikundi hubaki kuwa ya akaunti, na ccam husema hivyo. Kwa `N`, bila terminal, au wakati akaunti inaendeshwa, haibadilishi chochote na husema hivyo:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. Hurekodi kipindi, huweka `CLAUDE_CONFIG_DIR` kuwa folda ya akaunti na hujibadilisha yenyewe kuwa `claude rc` kupitia `exec`.

Kisha Remote Control huendelea kufanya kazi kama seva. Kila kipindi inachokuanzishia ni mchakato mtoto wa mchakato huo, kwa hiyo huendeshwa kama `work`: kwa taarifa za kuingia, mipangilio na historia ya mazungumzo ya `work`, na huhesabiwa dhidi ya kikomo cha matumizi cha `work`. Wakati seva iliyoanzishwa na ccam inafanya kazi, ccam huichukulia akaunti kuwa inaendeshwa, vivyo hivyo kwa seva iliyoanzishwa kwa `claude remote-control` ya kawaida inapokuwa na vipindi; seva kama hiyo ambayo bado haina kipindi haionekani.

Jaribio (dry run) huonyesha tofauti. Uanzishaji wa kawaida hubeba seva za kikundi kama hoja, Remote Control haibebi:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude rc
```

## Kwa nini inafanya kazi hivi

**Remote Control haipokei `--mcp-config`.** Katika uanzishaji wa kawaida ccam hukabidhi seva za kikundi kwa Claude Code kama `--mcp-config=<path>` na kamwe haiandiki faili za Claude Code. `claude remote-control` hukataa chaguo hilo kabla ya jina la amri yake na haipokei yoyote baada yake, kwa hiyo kabla ya ccam 0.6.0, `ccam work rc` ilishindwa kwa kila akaunti iliyoshiriki seva za MCP. Sasa ccam huacha chaguo hilo nje.

**Vipindi vyake husoma seva kutoka `.claude.json` tu.** Vipindi ambavyo Remote Control huanzisha havirithi chaguo kutoka kwenye mstari wa amri. Mahali pekee vinapopata seva za MCP za mtumiaji ni `.claude.json` ya akaunti. Faili hiyo pia hushikilia taarifa za kuingia za akaunti na hali ya kila folda, kwa hiyo ccam haiwezi kuiunganisha na nakala inayoshirikiwa kama inavyofanya kwa `settings.json`. Kunakili seva ndani yake ndiyo njia pekee ya kuvipa vipindi orodha inayoshirikiwa.

**ccam huuliza kabla ya kuandika.** Kwa mkakati wa flag ambao ni chaguo-msingi, ccam huacha `.claude.json` bila kuigusa wakati wa uanzishaji. Remote Control ndiyo uanzishaji pekee ambapo inalazimika kuandika humo, kwa hiyo huuliza kwanza, huweka nakala rudufu, na haigusi chochote isipokuwa `mcpServers`. (Nje ya uanzishaji, kuacha au kujiunga na kushiriki MCP na `ccam mcp import` nazo huandika humo, kama mipango yake inavyosema.)

**Haiandiki wakati akaunti inaendeshwa.** Claude Code huandika upya `.claude.json` wakati wa kipindi, na badiliko lililofanywa kutoka nje lingepotea au kuandikwa juu yake. Ndiyo maana ccam huruka kunakili wakati kipindi chochote cha akaunti kinaendeshwa, ikiwemo seva nyingine ya Remote Control na Claude Code iliyoanzishwa bila ccam.

**Huuliza mara moja kwa kila badiliko.** Baada ya kunakili, `sync-stamp` hushikilia muda wa kubadilishwa wa `mcp.json` ya kikundi. ccam huuliza tena tu faili hiyo inapobadilika, kwa mfano baada ya `ccam mcp add`, au seva iliyonakiliwa inapotoweka kutoka `.claude.json`. Ukijibu `N`, bado hakuna nakala, kwa hiyo huuliza kwa kila uanzishaji wa Remote Control.

**`ccam mcp` huweka nakala zikiwa za sasa.** Kipengele kilicho katika `.claude.json` hushinda `--mcp-config` majina yanapogongana, kwa hiyo nakala ya zamani ingeficha ufafanuzi mpya zaidi hata katika uanzishaji wa kawaida. Kwa sababu hiyo `ccam mcp add`, `add-json`, `remove` na `import`, pamoja na `ccam migrate`, husasisha nakala za kila akaunti isiyoendeshwa, na husema hivyo nakala inapobadilika:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**Bila terminal jibu ni hapana.** Remote Control mara nyingi huendeshwa kutoka kwenye hati (script), huduma au `nohup`, ambapo hakuna anayeweza kujibu. Hapo ccam huanza bila seva hizo badala ya kusubiri.

## Matatizo na suluhisho

### Seva za MCP zinazoshirikiwa hazipo katika vipindi vya Remote Control

Swali lilijibiwa `N`, hapakuwa na terminal, au akaunti ilikuwa inaendeshwa. Simamisha Remote Control, funga vipindi vingine vya akaunti, na uendeshe `ccam work rc` tena kwenye terminal; jibu `y`. Ili usiulizwe kamwe, weka `mcp_strategy = "sync"` katika `~/.config/ccam/config.toml`. Seva ambazo tayari ziko katika kipindi hazibadiliki wakati kinaendeshwa, kwa hiyo anzisha kipindi kipya baada ya kunakili.

### "Account work is running, so its .claude.json is left alone"

Kipindi cha `work` bado kiko wazi: `ccam work` ya kawaida, amri ya njia fupi, au seva nyingine ya Remote Control. Toka humo na uendeshe `ccam work rc` tena. Vipindi vinavyoendeshwa na akaunti nyingine havihesabiwi.

### ccam huuliza kila mara

Bado hujanakili seva. Jibu `y` mara moja, au weka `mcp_strategy = "sync"`. Ikiwa hutaki kamwe seva hizo katika Remote Control, acha kushiriki MCP na akaunti bila kuzihifadhi: `ccam group leave default work --fresh`, au weka MCP katika kikundi chake chenyewe ambacho akaunti haimo. `ccam group leave` ya kawaida hunakili seva za kikundi katika `.claude.json` ya akaunti, kwa hiyo Remote Control huwa nazo.

### `ccam work rc` hushindwa mara moja kwa kosa kuhusu `--mcp-config`

ccam yako ni ya zamani kuliko 0.6.0 na hupitisha `--mcp-config` kwa Remote Control. Endesha `ccam update`. Kosa lile lile kwa ccam ya sasa linamaanisha kwamba chaguo ambalo ccam hailijui, lenye thamani, lilikuja kabla ya `rc`, kwa hiyo ccam ilichukua thamani hiyo kama amri; weka `rc` mara moja baada ya jina la akaunti.

### Badiliko lililofanywa kwa `ccam mcp` halifiki kwa Remote Control

`ccam mcp add` na `remove` husasisha nakala katika akaunti zisizoendeshwa tu, na husema hivyo: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Simamisha Remote Control na uianzishe tena; ccam huona `mcp.json` iliyobadilika na huuliza tena.

### Seva yangu mwenyewe hutumika badala ya ile ya kikundi

Ikiwa akaunti ina seva yake yenye jina sawa na seva ya kikundi, unakili huiacha mahali pake na husema hivyo: `account work has its own MCP servers fs, so group default's of those names are not copied there`. Ili kutumia ya kikundi, ondoa ya akaunti (`claude mcp remove -s user <name>`); ili kuwapa akaunti zote toleo lako, lihamishie kwenye kikundi kwa `ccam mcp import work <server> --force`. Nakala rudufu za `.claude.json` ziko katika `~/.local/state/ccam/mcp/work/backups/`, ambayo huhifadhi tano za karibuni zaidi.

### Nataka kuondoa seva zilizonakiliwa

`ccam mcp remove <server>` huondoa seva kutoka kwenye kikundi na kutoka kila nakala. Ili kurudi kwenye hali ya kabla ya kunakili, rejesha nakala rudufu ya karibuni zaidi kutoka `~/.local/state/ccam/mcp/work/backups/` wakati hakuna kipindi cha akaunti kinachoendeshwa.

### Remote Control huuliza kama iamini folda, au hutoka bila terminal

Uaminifu wa folda huhifadhiwa kwa kila akaunti katika `.claude.json`, kwa hiyo folda uliyoiamini kama `personal` ni mpya kwa `work`. Endesha `ccam work` mara moja katika folda hiyo na ukubali, kisha anzisha Remote Control hapo.

### "You must be logged in with a Claude account that has a subscription"

Remote Control inahitaji kuingia kwa akaunti ya Claude yenye usajili. Kagua barua pepe ya akaunti katika `ccam list`, na uingie tena kwa `ccam login work` ikihitajika.

### Kipindi hakionekani kwenye simu

Programu ya Claude au claude.ai/code imeingia kwa akaunti tofauti ya Claude. Ilinganishe na barua pepe ambayo `ccam list` huonyesha kwa `work`.

### ccam huonya kwamba Claude Code nyingine inashiriki historia ya mazungumzo katika folda hii

Akaunti nyingine katika kikundi kile kile cha historia inaendeshwa katika folda ile ile. Vipindi vinavyoendeleza mazungumzo yale yale kutoka akaunti zote mbili vinaweza kuchanganya historia. Jibu `y` tu ikiwa zitafanya kazi kwenye mazungumzo tofauti, au anzisha Remote Control katika folda nyingine.

## Zinazohusiana

- [Kuanzisha Claude Code](/features/launch)
- [Seva za MCP zinazoshirikiwa](/features/mcp)
- [Kushiriki kati ya akaunti](/features/share)
- [Kusimamia akaunti](/features/accounts)
