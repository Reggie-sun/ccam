---
title: Fara Claude Code
description: Fara Claude Code a matsayin kowane asusu da ccam <name|number>, ccam run ko dashboard; arguments suna zuwa claude ba canji, --dry-run yana nuna umarnin.
updated: 2026-10-24
order: 1
translation: machine
sourceHash: e279341eb59b
---

`ccam work` yana fara Claude Code a matsayin asusun `work`, kuma duk abin da ke bayan sunan yana zuwa `claude` kamar yadda yake. Yi amfani da shi duk inda za ka rubuta `claude`: a terminal, a rubutun (script), ko tare da shigarwa ta bututu (pipe).

## Amfani

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

Ba da sunan asusu ko lambarsa daga `ccam list`. ccam ba ya karanta arguments da ke bayan sunan, don haka `ccam work --help` yana nuna taimakon Claude Code. `--` bayan sunan ba dole ba ne. Standard input, exit code da signals suna aiki kamar da `claude`.

Idan argument na farko ya fara da `-` kuma ba ɗaya daga cikin flags na ccam kansa ba ne, ccam yana fara asusu ba tare da suna ba, da dukkan arguments. Wannan shi ne asusun da aka saita don foldar yanzu da [`ccam folder`](/features/folder), in ba haka ba asusun tsoho: wanda aka saita da `ccam use`, ko asusun da ke da mafi ƙanƙantar lamba idan ba ka saita ɗaya ba.

### Dashboard

`ccam` ba tare da arguments ba yana buɗe [dashboard](/features/dashboard) idan standard input da output terminal ne, kuma yana buga `ccam list` in ba haka ba. Je kan asusu (da kibiyoyi, ko ta rubuta lambarsa) kuma danna `Enter` don fara shi; `n` yana fara asusun da ya fi sarari, kuma `o` yana ba da `--continue`, `--resume` ko wasu arguments. Dashboard yana rufewa kafin Claude Code ya fara, don haka ana fara maido da terminal.

**Kujeru (seats) na Team da Enterprise.** Zaɓin atomatik (`ccam next`, `--fallback`, `ccam pick`, maɓallin `n` na dashboard) yana kuma sauyawa tsakanin kujerun ƙungiyar Claude Team ko Enterprise iri ɗaya. Cibiyar taimako ta Anthropic ta ce "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (ana ba da kowace kujera a tsarin Team da sunan takamaiman mutum ɗaya, kuma wannan mutumin kaɗai ne aka yarda ya yi amfani da ita) da "Users can't route their work through another user's account to use that user's remaining usage" (masu amfani ba za su iya tura aikinsu ta asusun wani mai amfani ba domin su ci gajiyar ragowar amfanin wannan mai amfanin) ([Tsare-tsaren Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); haka ma ga [tsare-tsaren Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Don haka yin amfani da zaɓin atomatik tsakanin kujeru na iya saɓa wa sharuɗɗan Anthropic. Ka tuntuɓi ƙungiyarka da Anthropic kafin ka dogara da shi. Don tsare-tsaren kai (Pro, Max), duba [bayanin sharuɗɗa a cikin Tambayoyin da ake yawan yi](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

Idan harsashin (shell) yana da `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` ko `ANTHROPIC_AUTH_TOKEN` a saita, Claude Code zai yi amfani da wannan key ko token maimakon shigar asusun kansa, don haka ccam yana faɗin haka a standard error kafin ya fara (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). Yana karanta variables ɗin ne kawai; [asusun API key](/features/api) yana share su maimakon haka.

Da `relay = true` a saman `config.toml` (kafin `[[group]]`, `[[profile]]` ko `[[folder]]` na farko), ccam yana fara zaman mu'amala a matsayin process na ƙasa (child) kuma yana zama a gefensu, don `/ccam:switch` ya iya kai tattaunawa gaba da kansa, tare da zaɓuɓɓukan zaman, yanayin izini (permission mode) da matakin ƙoƙari (effort level) ([cikakken bayani](/features/next#relay)); in ba haka ba ccam yana maye gurbin kansa da Claude Code, kuma babu abin da ke gudana a gefensa.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--dry-run` | Buga muhalli (environment) da umarnin da ccam zai gudanar, sannan fita | a kashe |
| `--fallback` | Don gudun `-p`/`--print`: idan asusun ya kai iyakar amfaninsa, sake gudanar da umarnin iri ɗaya a matsayin asusun da aka shiga da shi mafi sarari ([cikakken bayani](/features/next#run-again-at-the-limit)) | a kashe |
| `-v`, `--version` | Buga `ccam <version>` (sai a matsayin argument na farko) | |
| `-h`, `--help` | Taimakon ccam, a matsayin `ccam --help` ko `ccam run --help` | |

Ƙa'idodin arguments:

- ccam yana karanta flags nasa ne kawai kafin sunan asusun. Bayan sunan, kowane argument na Claude Code ne, har da `--dry-run` da `--help`.
- `--dry-run` yana aiki a matsayin `ccam --dry-run work`, `ccam run --dry-run work` da `ccam --dry-run` (asusun da ba tare da suna ba). Ana karɓar sa, ba a ƙin sa. `--dry-run` da `--fallback` za su iya zuwa a kowane tsari.
- `ccam run <name|number> [claude args]` cikakken nau'i ne na `ccam <name>`. `ccam run` ba tare da suna ba yana amfani da asusun foldar ko na tsoho.
- Sunayen asusai ba za su iya zama sunayen umarnin ccam kamar `add`, `list`, `run` ko `usage` ba, don haka `ccam <word>` ba ya taɓa yin shubuha.

Dry run yana buga layi ɗaya da za ka iya liƙawa a shell, ko a PowerShell a Windows. Yana farawa da abin da farawar ke canzawa a muhalli: `CLAUDE_CONFIG_DIR` da `CCAM_PROFILE` koyaushe, `env -u <name>` ga variable da take cirewa (kamar `CCAM_RELAY` da aka gada), kuma ga [asusun API key](/features/api) endpoint ɗinsa da model ɗinsa da key ɗin a matsayin tushensa (`ANTHROPIC_API_KEY="$WORK_KEY"`), ba key ɗin kansa ba:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
& { $saved = @{}; foreach ($k in 'CLAUDE_CONFIG_DIR','CCAM_PROFILE') { $saved[$k] = [Environment]::GetEnvironmentVariable($k) }; $env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; $env:CCAM_PROFILE='work'; try { $p = Start-Process -FilePath 'C:\Users\me\.local\bin\claude.exe' -ArgumentList '--mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p "fix the tests"' -NoNewWindow -Wait -PassThru; $global:LASTEXITCODE = $p.ExitCode } finally { foreach ($k in $saved.Keys) { [Environment]::SetEnvironmentVariable($k, $saved[$k]) } } }
```

Layin PowerShell yana mayar da variables yadda suke bayan haka, kamar yadda `VAR=x cmd` ke yi a shell, kuma yana miƙa wa Claude Code layin umarninsa ta `Start-Process`, don haka Windows PowerShell 5.1 da PowerShell 7 suna gudanar da shi ta hanya ɗaya. Idan farawa ta saita variable zuwa ƙima mara komai, layin `#` da ke ƙasa da shi yana cewa PowerShell kafin 7.5 yana cire irin wannan variable maimakon haka. Ga `claude.cmd` da npm ya girka, layin yana gudanar da umarnin `cmd.exe` da farawar za ta yi amfani da shi (duba ƙasa), kuma argument da farawar za ta ƙi a can yana sa dry run ɗin ya gaza ma.

## Yadda yake aiki

Kowane farawa yana bin matakai iri ɗaya, daga layin umarni zuwa `exec`:

<!-- motion:launch -->

Kowane akwati yana haskakawa yayin da matakinsa ke gudana, kuma layukan da ke ƙasa su ne muhalli da arguments da Claude Code ke farawa da su.

1. `cmd/ccam/main.go` yana kiran `app.Main`. Idan an fara ccam da wani suna, [gajeren umarni](/features/shortcuts) ne kuma yana zuwa mai farawa (launcher) kai tsaye.
2. `app.parse` yana duba `os.Args` kafin cobra ya shiga. `ccam` shi kaɗai (dashboard), sunan umarnin ccam (`config.Reserved`), `-h`/`--help` da `-v`/`--version` ne kawai ke zuwa `cli.Execute`. Duk sauran farawa ne.
3. `launcher.Launch` yana karanta `~/.config/ccam/config.toml`, fayil ɗaya tilo da yake karantawa.
4. `launcher.Prepare` yana nemo asusun ta suna ko lamba (idan babu, `HereProfile` yana ɗaukar asusun foldar, yana kwatanta foldar aiki da shigarwar `[[folder]]` na `config.toml` a matsayin strings, ba tare da la'akari da manya ko ƙananan haruffa ba a macOS da Windows), yana cire `--` ɗaya na farko, kuma yana nemo shirin `claude` da `ResolveClaude`. Ana ajiye sakamakon neman PATH a cache a `~/.local/state/ccam/claude-path` tare da lokacin gyaran fayil ɗin da hash na `PATH`, don haka ccam ba ya bincika `PATH` a kowane farawa. Fayil da ya canza ko ya ɓace, ko `PATH` da ya canza, yana sa a sake nema. Hanyar `claude_bin` a `config.toml` (ana faɗaɗa `~/`) dole ta zama fayil mai gudana (executable); in ba haka ba farawar tana tsayawa da `claude_bin X is not an executable file. Check claude_bin in config.toml` kafin a rubuta komai. ccam yana ƙin `claude` wanda ccam kansa ne, ana kwatanta su ta inode, domin hakan zai yi ta zagayawa.
5. `checkLinks` yana kwatanta foldar asusun `~/.local/share/ccam/profiles/<name>/` da saitunan rabawarta ta amfani da `share.Check`, wanda ke kiran `lstat` da `readlink` kawai, kiran tsarin (system call) ɗaya ko biyu ga kowane abu. Ana ƙirƙirar mahaɗin da ya ɓace nan take. Duk abin da zai iya rasa bayanai yana dakatar da farawar.
6. `applyMCP` yana ƙara sabobin MCP da aka raba (duba ƙasa).
7. `session.Register` yana rubuta zaman. ccam yana rubuta alama mai sunan PID ɗinsa, mai ɗauke da lokacin farawar process ɗin, zuwa `~/.local/state/ccam/locks/profile-<name>/<pid>`, ko me `XDG_RUNTIME_DIR` ya ce, don kowane shell ya ga alamomi iri ɗaya. Ga asusun da ke raba tarihin tattaunawa, `historyCheck` yana fara neman zaman da ke aiki a `history-<group>-<hash>`, inda hash ɗin na foldar yanzu ne bayan warware mahaɗan alama (symbolic links), kuma yana yin rajista a can ma. `exec` yana riƙe PID ɗin, don haka alamar tana nufin process ɗin Claude Code kansa. Alama ce, ba kulle da aka riƙe ba, domin processes da Claude Code ke farawa za su gaji kulle kuma su riƙe shi bayan zaman ya ƙare. Idan ba a iya rubuta alamar ba, hakan ba ya dakatar da farawar: ccam yana cewa `ccam: could not mark the session as running (…); starting Claude Code anyway`. Sauran binciken ccam na ko asusu yana aiki (sake suna, cirewa, canjin rukunoni, `doctor --fix`, MCP sync) suna karanta bayanan Claude Code kansa a foldar asusun ma (`sessions/<pid>.json`, `daemon.lock`), don haka suna ganin Claude Code da aka fara ba ta ccam ba ma: `claude` na yau da kullum, IDE, zaman bayan fage da daemon ɗinsu.
8. `Plan.Exec` yana saita taken terminal zuwa `claude · <name>` (idan standard error terminal ne), yana rubuta lokacin farawa zuwa `~/.local/state/ccam/last-used/<name>` ba tare da `fsync` ba, kuma yana kiran `syscall.Exec`.

Claude Code yana gudana tare da `CLAUDE_CONFIG_DIR` da aka saita zuwa foldar asusun da `CCAM_PROFILE` da aka saita zuwa sunan asusun. ccam yana cire duk ƙimar waɗannan da ya gada, da kuma `CCAM_RELAY` da aka gada. Ga foldar `~/.claude` da aka yi wa rajista da `ccam import` a macOS, yana kuma saita `CLAUDE_SECURESTORAGE_CONFIG_DIR` zuwa ƙima mara komai, don asusun da `claude` na yau da kullum su yi amfani da shigar Keychain ɗaya; a Windows yana yin haka idan Claude Code yana ajiye shigar a Credential Manager.

### Me ya sa exec

`syscall.Exec` yana maye gurbin process ɗin ccam da `claude`, don haka ccam ba ya zama a ƙwaƙwalwa, danna maɓallai, signals da exit code na Claude Code ne kansa, kuma arguments suna zama array, ba taɓa zama string na shell ba.

Windows ba shi da `exec`. A can, `Plan.Exec` yana fara `claude.exe` a matsayin child process da muhalli, arguments da console iri ɗaya, yana jiran sa kuma yana fita da exit code ɗinsa. `claude.cmd` da npm ya girka yana gudana a `cmd.exe` tare da cikakkiyar hanyarsa da kowane argument a cikin alamomin zance, don kada `cmd.exe` ya karanta ko ɗaya a matsayin ƙa'idar rubutunsa (`&`, `%PATH%`); argument mai karya layi ba za a iya miƙa shi ta wannan hanya ba kuma ana ƙin sa, tare da shawarar girka Claude Code na asali (native). Ctrl+C yana kai wa Claude Code daga console kai tsaye, don haka ccam yana watsi da shi. Alamar zaman tana ɗauke da PID na ccam, wanda ke rayuwa daidai tsawon Claude Code, kuma ana cire ta idan Claude Code ya fita. Yayin da yake jira, ccam ba ya yin wani abu. Aikin bayan zama, kamar lura da sabobin da aka ƙara da `claude mcp add`, yana faruwa a farawa ta gaba: ccam yana kwatanta lokacin gyaran `.claude.json` da `~/.local/state/ccam/mcp/<name>/seen` kuma yana karanta fayil ɗin ne kawai idan ya canza.

### Yadda ake tura layin umarni

Layukan umarni uku, hanyoyi uku ta `internal/app`:

<!-- motion:launch-route -->

Ana gane gajeren umarni ta `argv[0]`, `ccam <name>` ta argument na farko, kuma duka biyun suna ƙarewa a `syscall.Exec`. Umarnin ccam kamar `ccam usage` ne kawai ke kai wa cobra.

### Sabobin MCP da aka raba

Da tsohuwar dabarar `flag`, ccam yana miƙa sabobin MCP na rukunin ga Claude Code a matsayin argument, wanda aka sa kafin naka:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` yana mayar da argument ne kawai idan `mcp.json` na rukunin yana nan, kuma `Plan.InsertArgs` yana saka shi kai tsaye bayan `claude`. Sannan `--mcp-config` ko `--strict-mcp-config` naka har yanzu yana aiki a saman sa. Argument ɗaya ne, `--mcp-config=<path>`, domin `--mcp-config` na Claude Code yana karɓar ƙima da dama: idan aka rubuta shi a matsayin arguments biyu, zai kuma karanta prompt da ke biye a matsayin fayil ɗin saiti.

`claudeVerb` yana yanke inda flag ɗin ba zai iya zuwa ba. Yana karanta ƙaramin umarnin (subcommand) Claude Code daga arguments kamar yadda Claude Code ke yi, yana tsallake zaɓuɓɓuka tare da ƙimominsu. `remote-control` ko `rc` ba ya samun flag (a ƙasa), haka ma aikatau na zaman bayan fage da aka ba da shi da farko (`attach`, `logs`, `stop`, `kill`, `respawn`, `rm`). `ccam work mcp list`, `get` da `remove` su ma ba sa samun ko ɗaya, tunda Claude Code yana ƙin shi a can; maimakon haka ccam yana faɗa a standard error cewa sabobin rukunin suna lodawa a zaman asusun ma:

```console
ccam: `claude mcp list` sees only account work's own MCP servers; its sessions also load group default's: ctx, fs
      See all of them with `ccam mcp list`; change the group's with `ccam mcp add|remove`.
```

Kowane sauran ƙaramin umarni da kowane zama suna riƙe flag ɗin.

Da `mcp_strategy = "sync"` a `config.toml`, ccam yana rubuta sabobin cikin `.claude.json` na asusun maimakon haka, kuma yana tsallake hakan idan asusun yana aiki tuni. Farawa tana tsallake aikin ma idan babu `mcp.json` na rukunin ko `.claude.json` da ya canza tun dubawa ta ƙarshe. Sync ba ya taɓa maye gurbin sabar da asusun ke da ita ta kansa mai sunan sabar rukuni (yana faɗin haka sau ɗaya), yana mayar da kwafin da ya ɓace daga `.claude.json` (`put back group default's MCP servers …`), kuma yana canzawa ko cire kwafin da shi ya rubuta waɗanda ba wanda ya canza tun daga lokacin kawai.

### Remote Control

`ccam work rc` (ko `ccam work remote-control`) yana fara Remote Control na Claude Code a matsayin `work`. `claude remote-control` yana ƙin `--mcp-config` kafin aikatau (verb) ɗinsa kuma ba ya karɓar ko ɗaya bayansa, kuma zamansa suna karanta sabobin MCP ne kawai daga `.claude.json` na asusun. Don haka dabarar flag tana aiki daban a nan:

<!-- motion:launch-rc -->

1. `claudeVerb` yana gano `remote-control` ko `rc` a matsayin aikatau na Claude Code, har ma bayan zaɓuɓɓuka da ƙimominsu (`ccam work --model opus rc`). Sannan `mcp.FlagArgs` ba ya ƙara komai.
2. `mcp.Synced` yana kwatanta lokacin gyaran `mcp.json` na rukunin da `sync-stamp` na asusun, kuma yana duba cewa kowace sabar rukuni har yanzu tana cikin `.claude.json`. Idan `.claude.json` yana ɗauke da sabobin rukunin na yanzu tuni, ko rukunin ba shi da ko ɗaya, Remote Control yana farawa ba tare da tambaya ba.
3. In ba haka ba, ccam yana faɗin sunayen sabobin da zaman za su rasa kuma yana tambaya:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; servers of the account's own stay)? [y/N]
   ```

4. Idan `y`, `mcp.Sync` yana kwafar su kamar yadda dabarar `sync` ke yi: ƙarƙashin kullen Claude Code, bayan ajiya, kuma yana riƙe sabobin asusun kansa. Idan `N`, ko babu terminal, Remote Control yana farawa ba tare da su ba, kuma ccam yana nuna `mcp_strategy = "sync"`, wanda ke kwafar su kafin kowace farawa.

ccam yana tambaya a kowace farawar Remote Control har sai an kwafi sabobin, kuma bayan kwafi sai kawai idan `mcp.json` na rukunin ya sake canzawa ko sabar da aka kwafa ta ɓace daga `.claude.json`. Kwafin da ke `.claude.json` yana nasara a kan `--mcp-config` a farawa ta yau da kullum na gaba, don haka `ccam mcp add`, `remove` da `import` suna sabunta kwafin da sauran asusan ke riƙe da su ([Sabobin MCP da aka raba](/features/mcp)).

### Kasafin 10 ms

ccam yana ƙara 10 ms (p95) a mafi yawa kafin Claude Code ya fara. Hanyar farawa ba ta taɓa gina bishiyar umarnin cobra ba, ba ta fara Bubble Tea ba ko buɗe haɗin network, kuma `TestPackageInitIsCheap` yana gazawa idan ayyukan `init()` na packages sun ɗauki fiye da 5 ms gaba ɗaya. `scripts/perf-check.sh` yana auna `ccam --dry-run work`, `ccam work` da gajeren umarni da `hyperfine` kuma yana gazawa sama da 10 ms, kuma `BenchmarkLaunchPath` a `internal/launcher` yana gudanar da matakan a madauki. Don ganin lokacin kowane mataki, saita `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Tsaro da yanayi na musamman

- Sunan ko lambar da ba a sani ba yana tsayawa da `no account "<name>"` kuma yana nuna `ccam list`. Idan babu asusu ko ɗaya, ccam yana roƙon ka ka ƙara ɗaya da `ccam add <name>`.
- Idan hanyar da aka raba ta zama fayil na yau da kullum ko tana nuna wurin da bai dace ba, ccam yana tsayawa da `stopped because of share link problems:`, yana lissafa su, kuma yana roƙon ka ka gudanar da `ccam doctor --fix`. Ba ya matsar da bayanai a hanyar farawa.
- Idan wani Claude Code da ke raba tarihin tattaunawa da wannan asusun (rukuni ɗaya) yana aiki a folda ɗaya, ccam yana gargaɗin cewa tarihin zai iya cakuɗewa kuma yana tambaya `Run anyway? [y/N]`. Ana daidaita foldar bayan warware mahaɗan alama (symbolic links). Ba tare da terminal ba (standard input daga `/dev/null` ana ƙirga shi a matsayin babu), ko daga aiki a bayan fage, yana nuna gargaɗin kuma yana ci gaba. Sabon gudun `-p` yana fara tattaunawarsa ta kansa kuma ba ya samun gargaɗi; `--version`, `--help` da ƙananan umarnin da ba sa buɗe tattaunawa, kamar `mcp`, `auth` ko `plugin`, ba a duba su kwata-kwata.
- `--dry-run` yana tsallake daidaita MCP (sync), sanarwar sabuwar sabar, bayanin `mcp list`, duba tarihi, alamar zama, take da lokacin amfani na ƙarshe, kuma ba ya yin exec. Har yanzu yana duba mahaɗan rabawa, kuma yana ƙirƙirar mahaɗin da ya ɓace kamar yadda farawa ta yau da kullum za ta yi.
- `ccam work rc` yayin da `work` yana aiki tuni yana barin `.claude.json` yadda yake, yana faɗin haka, kuma yana fara Remote Control ba tare da sabobin rukunin ba; sake gudanar da shi bayan zaman ya ƙare don kwafar su. Da `--dry-run` yana buga saƙon amma ba ya tambaya kuma ba ya kwafi.
- Da `mcp_strategy = "sync"`, sync da ya gaza yana buga `ccam: MCP sync failed (launching anyway): …` kuma Claude Code yana farawa da sabobin da yake da su tuni.

## Masu alaƙa

- [Gajerun umarni](/features/shortcuts)
- [Amfani da Remote Control tare da ccam](/guides/remote-control)
- [Sarrafa asusai](/features/accounts)
- [Bayanin umarnai](/docs/commands)
- [Raba sabobin MCP, saituna da tarihi](/docs/sharing)
- [Yadda CLAUDE_CONFIG_DIR ke aiki](/guides/claude-config-dir)
