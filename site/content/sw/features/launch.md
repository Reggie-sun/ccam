---
title: Kuanzisha Claude Code
description: Anzisha Claude Code kama akaunti yoyote kwa ccam <name|number>, ccam run au dashibodi; hoja hupita kwa claude bila kubadilishwa, na --dry-run huonyesha amri.
updated: 2026-10-25
order: 1
translation: machine
sourceHash: "afa4942cfed5"
---

`ccam work` huanzisha Claude Code kama akaunti `work`, na kila kitu kinachokuja baada ya jina huenda kwa `claude` kama kilivyo. Itumie popote ambapo ungecharaza `claude`: kwenye terminal, katika script, au ingizo likipitishwa kwa pipe.

## Matumizi

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the dashboard on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

Toa jina la akaunti au namba yake kutoka `ccam list`. ccam haisomi hoja zinazokuja baada ya jina, kwa hiyo `ccam work --help` huonyesha msaada wa Claude Code. `--` baada ya jina ni hiari. Standard input, exit code na signals hufanya kazi kama kwa `claude`.

Hoja ya kwanza inapoanza na `-` na si mojawapo ya chaguo za ccam yenyewe, ccam huanzisha akaunti bila jina, pamoja na hoja zote. Hiyo ni akaunti iliyowekwa kwa folda ya sasa kwa [`ccam folder`](/features/folder), la sivyo ile chaguo-msingi: akaunti iliyowekwa kwa `ccam use`, au akaunti yenye namba ya chini zaidi ikiwa hujaweka moja.

### Dashibodi

`ccam` bila hoja hufungua [dashibodi](/features/dashboard) wakati standard input na output ni terminal, na vinginevyo huchapisha `ccam list`. Nenda kwenye akaunti (kwa mishale, au kwa kuandika namba yake) na ubonyeze `Enter` ili kuianzisha; `n` huanzisha akaunti yenye nafasi kubwa zaidi na `o` hutoa `--continue`, `--resume` au hoja nyingine. Dashibodi hujifunga kabla Claude Code haijaanza, kwa hiyo terminal hurejeshwa kwanza.

**Viti (seats) vya Team na Enterprise.** Chaguo za kiotomatiki (`ccam next`, `--fallback`, `ccam pick`, kitufe `n` cha dashibodi, na [`auto_switch`](/features/next#auto-switch) unapokiwasha) pia hubadili kati ya viti vya shirika lile lile la Claude Team au Enterprise. Kituo cha usaidizi cha Anthropic kinasema "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (kila kiti katika mpango wa Team hutolewa kwa mtu mmoja aliyetajwa kwa jina, na ni mtu huyo pekee anayeruhusiwa kukitumia) na "Users can't route their work through another user's account to use that user's remaining usage" (watumiaji hawaruhusiwi kuelekeza kazi yao kupitia akaunti ya mtumiaji mwingine ili kutumia kiasi cha matumizi kilichobaki cha mtumiaji huyo) ([mipango ya Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); vivyo hivyo kwa [mipango ya Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Kwa hiyo, kutumia chaguo za kiotomatiki kati ya viti kunaweza kukiuka masharti ya Anthropic. Kama utazitumia huko ni uamuzi wako: wasiliana na shirika lako na Anthropic kabla ya kuzitegemea. Mipango ya binafsi (Pro, Max) inaelezwa katika [dokezo la masharti kwenye Maswali yanayoulizwa mara kwa mara](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

Shell ikiwa na `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` au `ANTHROPIC_AUTH_TOKEN` iliyowekwa, Claude Code ingetumia key au tokeni hiyo badala ya taarifa za kuingia za akaunti yenyewe, kwa hiyo ccam husema hivyo kwenye standard error kabla ya kuanza (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). Husoma vigezo hivyo tu; [akaunti ya API key](/features/api) huvifuta badala yake.

Kwa `relay = true`, au [`auto_switch = true`](/features/next#auto-switch) kwenye Linux na macOS, juu kabisa ya `config.toml` (kabla ya `[[group]]`, `[[profile]]` au `[[folder]]` ya kwanza), ccam huanzisha vipindi shirikishi kama mchakato mtoto na hubaki kando yake, ili `/ccam:switch`, au `auto_switch` kwenye kikomo cha matumizi, iweze kuhamisha mazungumzo yenyewe, pamoja na chaguo za kipindi, hali ya ruhusa na kiwango cha juhudi (effort) ([maelezo](/features/next#relay)); vinginevyo ccam hujibadilisha na Claude Code, na hakuna kinachoendeshwa kando yake.

## Chaguo

| Chaguo | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `--dry-run` | Huchapisha mazingira na amri ambayo ccam ingeendesha, kisha hutoka | imezimwa |
| `--fallback` | Kwa uendeshaji wa `-p`/`--print`: akaunti inapofikia kikomo chake cha matumizi, huendesha amri ile ile tena kama akaunti iliyoingia yenye nafasi kubwa zaidi ([maelezo](/features/next#run-again-at-the-limit)) | imezimwa |
| `-v`, `--version` | Huchapisha `ccam <version>` (kama hoja ya kwanza tu) | |
| `-h`, `--help` | Msaada wa ccam, kama `ccam --help` au `ccam run --help` | |

Kanuni za hoja:

- ccam husoma chaguo zake yenyewe kabla ya jina la akaunti tu. Baada ya jina, kila hoja ni ya Claude Code, pamoja na `--dry-run` na `--help`.
- `--dry-run` hufanya kazi kama `ccam --dry-run work`, `ccam run --dry-run work` na `ccam --dry-run` (akaunti bila jina). Inaheshimiwa, haikataliwi. `--dry-run` na `--fallback` zinaweza kuja kwa mpangilio wowote.
- `ccam run <name|number> [claude args]` ni muundo mrefu wa `ccam <name>`. `ccam run` bila jina hutumia akaunti ya folda au ile chaguo-msingi.
- Majina ya akaunti hayawezi kuwa majina ya amri za ccam kama `add`, `list`, `run` au `usage`, kwa hiyo `ccam <word>` kamwe haina utata.

Jaribio la dry run huchapisha mstari mmoja ambao unaweza kubandika kwenye shell, au kwenye PowerShell kwenye Windows. Huanza na kile ambacho uanzishaji hubadilisha katika mazingira: `CLAUDE_CONFIG_DIR` na `CCAM_PROFILE` kila mara, `env -u <name>` kwa kigezo kinachoondolewa (kama `CCAM_RELAY` iliyorithiwa), na kwa [akaunti ya API key](/features/api) endpoint yake, model yake na key kama chanzo chake (`ANTHROPIC_API_KEY="$WORK_KEY"`), kamwe si key yenyewe:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

Mstari wa PowerShell hurudisha vigezo baadaye, kama `VAR=x cmd` inavyofanya kwenye shell, na huipa Claude Code mstari wake wa amri kupitia `Start-Process`, kwa hiyo Windows PowerShell 5.1 na PowerShell 7 huuendesha kwa njia ile ile. Uanzishaji unapoweka kigezo kuwa thamani tupu, mstari wa `#` chini yake husema kwamba PowerShell kabla ya 7.5 huondoa kigezo kama hicho badala yake. Kwa `claude.cmd` iliyosakinishwa na npm, mstari huo huendesha amri ya `cmd.exe` ambayo uanzishaji ungetumia (tazama hapa chini), na hoja ambayo uanzishaji ungeikataa huko hufanya dry run ishindwe pia.

## Jinsi inavyofanya kazi

Kila uanzishaji hupitia hatua zile zile, kuanzia mstari wa amri hadi `exec`:

<!-- motion:launch -->

Kila kisanduku huwaka hatua yake inapoendeshwa, na mistari iliyo chini ni mazingira na hoja ambazo Claude Code huanza nazo.

1. `cmd/ccam/main.go` huita `app.Main`. Ikiwa ccam ilianzishwa kwa jina lingine, ni [njia fupi](/features/shortcuts) na huenda moja kwa moja kwa launcher.
2. `app.parse` huangalia `os.Args` kabla cobra haijahusika. `ccam` peke yake (dashibodi), jina la amri ya ccam (`config.Reserved`), `-h`/`--help` na `-v`/`--version` pekee ndizo huenda kwa `cli.Execute`. Kila kitu kingine ni uanzishaji.
3. `launcher.Launch` husoma `~/.config/ccam/config.toml`, faili pekee ambayo huichanganua.
4. `launcher.Prepare` hupata akaunti kwa jina au namba (bila jina, `HereProfile` huchukua akaunti ya folda, ikilinganisha folda ya kazi na maingizo ya `[[folder]]` ya `config.toml` kama mifuatano ya maandishi, bila kujali herufi kubwa au ndogo kwenye macOS na Windows), huondoa `--` moja ya mwanzo, na hutafuta programu ya `claude` kwa `ResolveClaude`. Utafutaji katika PATH huhifadhiwa kwenye cache katika `~/.local/state/ccam/claude-path` pamoja na muda wa mwisho wa kubadilishwa kwa faili na hash ya `PATH`, kwa hiyo ccam haitafuti `PATH` kila inapoanzisha. Faili iliyobadilika au inayokosekana, au `PATH` iliyobadilika, husababisha utafutaji mpya. Njia ya `claude_bin` katika `config.toml` (`~/` hupanuliwa) lazima iwe faili inayotekelezeka; vinginevyo uanzishaji husimama na `claude_bin X is not an executable file. Check claude_bin in config.toml` kabla ya kitu chochote kuandikwa. ccam hukataa `claude` ambayo ni ccam yenyewe, ikilinganisha kwa inode, kwa sababu hilo lingesababisha mzunguko usioisha.
5. `checkLinks` hulinganisha folda ya akaunti `~/.local/share/ccam/profiles/<name>/` na mipangilio yake ya kushiriki kwa `share.Check`, ambayo huita `lstat` na `readlink` tu, system call moja au mbili kwa kila kipengele. Kiungo kinachokosekana huundwa papo hapo. Chochote kinachoweza kupoteza data husimamisha uanzishaji.
6. `applyMCP` huongeza seva za MCP zinazoshirikiwa (tazama hapa chini).
7. `session.Register` hurekodi kipindi. ccam huandika alama yenye jina la PID yake, ikiwa na muda wa kuanza wa mchakato, kwenye `~/.local/state/ccam/locks/profile-<name>/<pid>`, bila kujali `XDG_RUNTIME_DIR` inasema nini, ili kila shell ione alama zile zile. Kwa akaunti inayoshiriki historia ya mazungumzo, `historyCheck` kwanza hutafuta kipindi kinachoendeshwa katika `history-<group>-<hash>`, ambapo hash ni ya folda ya sasa baada ya symbolic links kufuatwa, na husajili hapo pia. `exec` huhifadhi PID, kwa hiyo baadaye alama hutaja mchakato wa Claude Code wenyewe. Ni alama, si kufuli inayoshikiliwa, kwa sababu michakato ambayo Claude Code huanzisha ingerithi kufuli na kuiendelea kuishikilia baada ya kipindi kuisha. Alama isiyoweza kuandikwa haisimamishi uanzishaji: ccam husema `ccam: could not mark the session as running (…); starting Claude Code anyway`. Ukaguzi mwingine wa ccam wa kama akaunti inaendeshwa (kubadilisha jina, kuondoa, mabadiliko ya vikundi, `doctor --fix`, usawazishaji wa MCP) pia husoma rekodi za Claude Code yenyewe katika folda ya akaunti (`sessions/<pid>.json`, `daemon.lock`), kwa hiyo huona pia Claude Code iliyoanzishwa bila ccam: `claude` ya kawaida, IDE, vipindi vya chinichini na daemon yao.
8. `Plan.Exec` huweka kichwa cha terminal kuwa `claude · <name>` (wakati standard error ni terminal), huandika muda wa uanzishaji kwenye `~/.local/state/ccam/last-used/<name>` bila `fsync`, na huita `syscall.Exec`.

Claude Code huendeshwa ikiwa `CLAUDE_CONFIG_DIR` imewekwa kuwa folda ya akaunti na `CCAM_PROFILE` imewekwa kuwa jina la akaunti. ccam huondoa thamani zozote za vigezo hivi ambazo ilirithi, pamoja na `CCAM_RELAY` iliyorithiwa. Kwa folda ya `~/.claude` iliyosajiliwa kwa `ccam import` kwenye macOS, pia huweka `CLAUDE_SECURESTORAGE_CONFIG_DIR` kuwa thamani tupu, ili akaunti na `claude` ya kawaida zitumie kuingia kule kule kwa Keychain; kwenye Windows hufanya vivyo hivyo Claude Code inapohifadhi kuingia katika Credential Manager.

### Kwa nini exec

`syscall.Exec` hubadilisha mchakato wa ccam kwa `claude`, kwa hiyo ccam haibaki kwenye kumbukumbu, vibonyezo vya vitufe, signals na exit code ni vya Claude Code yenyewe, na hoja hubaki kuwa safu (array), kamwe si mfuatano wa shell.

Windows haina `exec`. Huko, `Plan.Exec` huanzisha `claude.exe` kama mchakato mtoto kwa mazingira, hoja na console zile zile, huusubiri na hutoka kwa exit code yake. `claude.cmd` iliyosakinishwa na npm huendeshwa katika `cmd.exe` kwa njia yake kamili na kila hoja ikiwa ndani ya alama za nukuu ili `cmd.exe` isisome yoyote kama sintaksia yake yenyewe (`&`, `%PATH%`); hoja yenye kivunja mstari haiwezi kupitishwa kwa njia hiyo na hukataliwa, pamoja na ushauri wa kusakinisha Claude Code asilia (native). Ctrl+C hufika kwa Claude Code moja kwa moja kutoka kwenye console, kwa hiyo ccam huipuuza. Alama ya kipindi hubeba PID ya ccam, ambayo huishi muda sawasawa na Claude Code, na huondolewa Claude Code inapotoka. Wakati inasubiri, ccam haifanyi kitu kingine. Kazi baada ya kipindi, kama kugundua seva zilizoongezwa kwa `claude mcp add`, hufanyika katika uanzishaji unaofuata: ccam hulinganisha muda wa kubadilishwa wa `.claude.json` na `~/.local/state/ccam/mcp/<name>/seen` na huichanganua faili tu ilipobadilika.

### Jinsi mstari wa amri unavyoelekezwa

Mistari mitatu ya amri, njia tatu kupitia `internal/app`:

<!-- motion:launch-route -->

Njia fupi hutambuliwa kwa `argv[0]`, `ccam <name>` kwa hoja ya kwanza, na zote mbili huishia katika `syscall.Exec`. Amri ya ccam kama `ccam usage` pekee ndiyo hufika kwa cobra.

### Seva za MCP zinazoshirikiwa

Kwa mkakati wa chaguo-msingi wa `flag`, ccam hupitisha seva za MCP za kikundi kwa Claude Code kama hoja, ikiwekwa kabla ya zako:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` hurudisha hoja tu ikiwa `mcp.json` ya kikundi ipo, na `Plan.InsertArgs` huiweka mara tu baada ya `claude`. `--mcp-config` au `--strict-mcp-config` yako mwenyewe basi bado inatumika juu yake. Ni hoja moja, `--mcp-config=<path>`, kwa sababu `--mcp-config` ya Claude Code hupokea thamani kadhaa: ikiandikwa kama hoja mbili, ingesoma pia prompt inayofuata kama faili ya usanidi.

`claudeVerb` huamua ni wapi flag haiwezi kuwekwa. Husoma amri ndogo (subcommand) ya Claude Code kutoka kwenye hoja kwa njia ile ile ambayo Claude Code hufanya, ikiruka chaguo pamoja na thamani zake. `remote-control` au `rc` haipati flag (tazama hapa chini), wala kitenzi cha kipindi cha chinichini kilichotolewa kwanza (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`). `ccam work mcp list`, `get` na `remove` pia hazipati, kwa kuwa Claude Code huikataa huko; badala yake ccam husema kwenye standard error kwamba seva za kikundi pia hupakiwa katika vipindi vya akaunti:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

Kila amri ndogo nyingine na kila kipindi hubaki na flag.

Kwa `mcp_strategy = "sync"` katika `config.toml`, ccam huandika seva katika `.claude.json` ya akaunti badala yake, na huruka hilo ikiwa akaunti tayari inaendeshwa. Uanzishaji pia huruka kazi hiyo wakati `mcp.json` ya kikundi wala `.claude.json` hazijabadilika tangu ukaguzi wa mwisho. Usawazishaji kamwe haubadilishi seva ambayo akaunti inayo yenyewe kwa jina la seva ya kikundi (husema hivyo mara moja), hurudisha nakala iliyotoweka kutoka `.claude.json` (`put back group default's MCP servers …`), na hubadilisha au kuondoa tu nakala ulizoandika ambazo hakuna aliyezibadilisha tangu hapo.

### Remote Control

`ccam work rc` (au `ccam work remote-control`) huanzisha Remote Control ya Claude Code kama `work`. `claude remote-control` hukataa `--mcp-config` kabla ya kitenzi chake na haipokei yoyote baada yake, na vipindi vyake husoma seva za MCP kutoka `.claude.json` ya akaunti tu. Kwa hiyo mkakati wa flag hufanya kazi tofauti hapa:

<!-- motion:launch-rc -->

1. `claudeVerb` hupata `remote-control` au `rc` kama kitenzi cha Claude Code, hata baada ya chaguo na thamani zake (`ccam work --model opus rc`). Kisha `mcp.FlagArgs` haiongezi chochote.
2. `mcp.Synced` hulinganisha muda wa kubadilishwa wa `mcp.json` ya kikundi na `sync-stamp` ya akaunti, na hukagua kwamba kila seva ya kikundi bado iko katika `.claude.json`. Ikiwa `.claude.json` tayari ina seva za sasa za kikundi, au kikundi hakina seva yoyote, Remote Control huanza bila swali.
3. Vinginevyo ccam hutaja seva ambazo vipindi vingezikosa na huuliza:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; servers of the account's own stay)? [y/N]
   ```

4. Kwa `y`, `mcp.Sync` huzinakili kama mkakati wa `sync` unavyofanya: chini ya kufuli ya Claude Code, baada ya nakala rudufu, ikihifadhi seva za akaunti yenyewe. Kwa `N`, au bila terminal, Remote Control huanza bila hizo, na ccam huelekeza kwa `mcp_strategy = "sync"`, ambayo huzinakili kabla ya kila uanzishaji.

ccam huuliza kwa kila uanzishaji wa Remote Control hadi seva zinakiliwe, na baada ya kunakiliwa huuliza tena tu `mcp.json` ya kikundi inapobadilika tena au seva iliyonakiliwa inapotoweka kutoka `.claude.json`. Nakala iliyo katika `.claude.json` hushinda `--mcp-config` katika uanzishaji wa kawaida unaofuata, kwa hiyo `ccam mcp add`, `remove` na `import` husasisha nakala zinazoshikiliwa na akaunti nyingine ([Seva za MCP zinazoshirikiwa](/features/mcp)).

### Bajeti ya ms 10

ccam huongeza kwa kiwango cha juu ms 10 (p95) kabla Claude Code haijaanza. Njia ya uanzishaji kamwe haijengi mti wa amri wa cobra, haianzishi Bubble Tea wala haifungui muunganisho wa mtandao, na `TestPackageInitIsCheap` hushindwa ikiwa fanksheni za `init()` za vifurushi zinachukua zaidi ya ms 5 kwa jumla. `scripts/perf-check.sh` hupima `ccam --dry-run work`, `ccam work` na njia fupi kwa `hyperfine` na hushindwa juu ya ms 10, na `BenchmarkLaunchPath` katika `internal/launcher` huendesha hatua hizo kwa mzunguko. Ili kuona muda wa kila hatua, weka `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Usalama na hali maalum

- Jina au namba isiyojulikana husimama na `no account "<name>"` na huelekeza kwa `ccam list`. Ikiwa hakuna akaunti kabisa, ccam hukuomba uongeze moja kwa `ccam add <name>`.
- Ikiwa njia inayoshirikiwa imegeuka kuwa faili ya kawaida au inaelekea mahali pasipo sahihi, ccam husimama na `stopped because of share link problems:`, huziorodhesha, na hukuomba uendeshe `ccam doctor --fix`. Haihamishi data kwenye njia ya uanzishaji.
- Ikiwa Claude Code nyingine inayoshiriki historia ya mazungumzo na akaunti hii (kikundi kile kile) inaendeshwa katika folda ile ile, ccam huonya kwamba historia inaweza kuchanganyika na huuliza `Run anyway? [y/N]`. Folda hulinganishwa baada ya symbolic links kufuatwa. Bila terminal (standard input kutoka `/dev/null` huhesabiwa kama hakuna), au kutoka kazi iliyo chinichini, huonyesha onyo na huendelea. Uendeshaji mpya wa `-p` huanzisha mazungumzo yake yenyewe na haupati onyo; `--version`, `--help` na amri ndogo zisizofungua mazungumzo, kama `mcp`, `auth` au `plugin`, hazikaguliwi kabisa.
- `--dry-run` huruka usawazishaji wa MCP, taarifa ya seva mpya, ujumbe wa `mcp list`, ukaguzi wa historia, alama ya kipindi, kichwa na muda wa matumizi ya mwisho, na haifanyi exec. Bado hukagua viungo vya kushiriki, na huunda kiungo kinachokosekana kama uanzishaji wa kawaida ungefanya.
- `ccam work rc` wakati `work` tayari inaendeshwa huacha `.claude.json` bila kuigusa, husema hivyo, na huanzisha Remote Control bila seva za kikundi; iendeshe tena baada ya kipindi kuisha ili kuzinakili. Kwa `--dry-run` huchapisha ujumbe lakini haiulizi wala hainakili.
- Kwa `mcp_strategy = "sync"`, usawazishaji ulioshindwa huchapisha `ccam: MCP sync failed (launching anyway): …` na Claude Code huanza na seva ilizokuwa nazo tayari.

## Zinazohusiana

- [Amri za njia fupi](/features/shortcuts)
- [Kutumia Remote Control na ccam](/guides/remote-control)
- [Kusimamia akaunti](/features/accounts)
- [Marejeo ya amri](/docs/commands)
- [Shiriki seva za MCP, mipangilio na historia](/docs/sharing)
- [Jinsi CLAUDE_CONFIG_DIR inavyofanya kazi](/guides/claude-config-dir)
