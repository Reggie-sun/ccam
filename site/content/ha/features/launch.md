---
title: Fara Claude Code
description: Fara Claude Code a matsayin kowane asusu da ccam <name|number>, ccam run ko mai zaɓi na ccam; arguments suna zuwa claude ba tare da canji ba, kuma --dry-run yana nuna umarnin.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: d033da86f44b
---

`ccam work` yana fara Claude Code a matsayin asusun `work`, kuma duk abin da ke bayan sunan yana zuwa `claude` kamar yadda yake. Yi amfani da shi duk inda za ka rubuta `claude`: a terminal, a rubutun (script), ko tare da shigarwa ta bututu (pipe).

## Amfani

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the default account
ccam use work                    # make work the default account
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
```

Ba da sunan asusu ko lambarsa daga `ccam list`. ccam ba ya karanta arguments da ke bayan sunan, don haka `ccam work --help` yana nuna taimakon Claude Code. `--` bayan sunan ba dole ba ne. Standard input, exit code da signals suna aiki kamar da `claude`.

Idan argument na farko ya fara da `-` kuma ba ɗaya daga cikin flags na ccam kansa ba ne, ccam yana fara asusun tsoho da dukkan arguments. Tsohon asusu shi ne wanda aka saita da `ccam use`, ko asusun da ke da mafi ƙanƙantar lamba idan ba ka saita ɗaya ba.

### Mai zaɓi

`ccam` ba tare da arguments ba yana buɗe mai zaɓi (picker) idan standard input da output terminal ne, kuma yana buga `ccam list` in ba haka ba. Yana nuna lambar kowane asusu, imel, amfanin awa 5, amfani na ƙarshe da abubuwan da aka raba; amfani yana cika yayin da kowane asusu ya amsa.

| Maɓalli | Abin da yake yi |
| --- | --- |
| `Enter` | Fara Claude Code a matsayin asusun da aka zaɓa |
| `1` zuwa `9` | Fara asusun da ke da wannan lambar |
| `Up`/`Down`, `k`/`j` | Matsar da zaɓin |
| `a` | Rufe mai zaɓin kuma gudanar da mataimakin `ccam add` |
| `s` | Gyara abubuwan da asusun da aka zaɓa ke rabawa (`Space` ko `x` yana sauyawa, `Enter` yana aiwatarwa, `Esc` yana komawa baya) |
| `d` | Mai da asusun da aka zaɓa tsoho |
| `r` | Sake ɗauko amfani, ba tare da gajeren cache ba |
| `q`, `Esc`, `Ctrl+C` | Fita |

Mai zaɓin yana rufewa kafin ya yi aiki, don haka ana maido da terminal kafin Claude Code ya fara. Bayan `d` ko canjin rabawa yana sake buɗewa. Canjin rabawa yana nuna shirin kuma yana fara tambaya, kamar yadda `ccam share` ke yi.

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--dry-run` | Buga muhalli (environment) da umarnin da ccam zai gudanar, sannan fita | a kashe |
| `-v`, `--version` | Buga `ccam <version>` (sai a matsayin argument na farko) | |
| `-h`, `--help` | Taimakon ccam, a matsayin `ccam --help` ko `ccam run --help` | |

Ƙa'idodin arguments:

- ccam yana karanta flags nasa ne kawai kafin sunan asusun. Bayan sunan, kowane argument na Claude Code ne, har da `--dry-run` da `--help`.
- `--dry-run` yana aiki a matsayin `ccam --dry-run work`, `ccam run --dry-run work` da `ccam --dry-run` (asusun tsoho). Ana karɓar sa, ba a ƙin sa.
- `ccam run <name|number> [claude args]` cikakken nau'i ne na `ccam <name>`. `ccam run` ba tare da suna ba yana amfani da asusun tsoho.
- Sunayen asusai ba za su iya zama sunayen umarnin ccam kamar `add`, `list`, `run` ko `usage` ba, don haka `ccam <word>` ba ya taɓa yin shubuha.

Dry run yana buga layi ɗaya da za ka iya liƙawa a shell:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

## Yadda yake aiki

Kowane farawa yana bin matakai iri ɗaya, daga layin umarni zuwa `exec`:

<!-- motion:launch -->

Kowane akwati yana haskakawa yayin da matakinsa ke gudana, kuma layukan da ke ƙasa su ne muhalli da arguments da Claude Code ke farawa da su.

1. `cmd/ccam/main.go` yana kiran `app.Main`. Idan an fara ccam da wani suna, [gajeren umarni](/features/shortcuts) ne kuma yana zuwa mai farawa (launcher) kai tsaye.
2. `app.parse` yana duba `os.Args` kafin cobra ya shiga. `ccam` shi kaɗai (mai zaɓi), sunan umarnin ccam (`config.Reserved`), `-h`/`--help` da `-v`/`--version` ne kawai ke zuwa `cli.Execute`. Duk sauran farawa ne.
3. `launcher.Launch` yana karanta `~/.config/ccam/config.toml`, fayil ɗaya tilo da yake karantawa.
4. `launcher.Prepare` yana nemo asusun ta suna ko lamba, yana cire `--` ɗaya na farko, kuma yana nemo shirin `claude` da `ResolveClaude`. Ana ajiye sakamakon neman PATH a cache a `~/.local/state/ccam/claude-path` tare da lokacin gyaran fayil ɗin, don haka ccam ba ya bincika `PATH` a kowane farawa. Fayil da ya canza ko ya ɓace yana sa a sake nema. Ana amfani da hanyar `claude_bin` a `config.toml` kamar yadda take. ccam yana ƙin `claude` wanda ccam kansa ne, ana kwatanta su ta inode, domin hakan zai yi ta zagayawa.
5. `checkLinks` yana kwatanta foldar asusun `~/.local/share/ccam/profiles/<name>/` da saitunan rabawarta ta amfani da `share.Check`, wanda ke kiran `lstat` da `readlink` kawai, kiran tsarin (system call) ɗaya ko biyu ga kowane abu. Ana ƙirƙirar mahaɗin da ya ɓace nan take. Duk abin da zai iya rasa bayanai yana dakatar da farawar.
6. `applyMCP` yana ƙara sabobin MCP da aka raba (duba ƙasa).
7. `session.Register` yana rubuta zaman. ccam yana rubuta alama mai sunan PID ɗinsa, mai ɗauke da lokacin farawar process ɗin, zuwa `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` (ko `~/.local/state/ccam/locks/` idan babu `XDG_RUNTIME_DIR`). Ga asusun da ke raba tarihin tattaunawa, `historyCheck` yana fara neman zaman da ke aiki a `history-<group>-<hash>`, inda hash ɗin na foldar yanzu ne, kuma yana yin rajista a can ma. `exec` yana riƙe PID ɗin, don haka alamar tana nufin process ɗin Claude Code kansa. Alama ce, ba kulle da aka riƙe ba, domin processes da Claude Code ke farawa za su gaji kulle kuma su riƙe shi bayan zaman ya ƙare.
8. `Plan.Exec` yana saita taken terminal zuwa `claude · <name>` (idan standard error terminal ne), yana rubuta lokacin farawa zuwa `~/.local/state/ccam/last-used/<name>` ba tare da `fsync` ba, kuma yana kiran `syscall.Exec`.

Claude Code yana gudana tare da `CLAUDE_CONFIG_DIR` da aka saita zuwa foldar asusun da `CCAM_PROFILE` da aka saita zuwa sunan asusun. ccam yana cire duk ƙimar waɗannan da ya gada. Ga foldar `~/.claude` da aka yi wa rajista da `ccam import` a macOS, yana kuma saita `CLAUDE_SECURESTORAGE_CONFIG_DIR` zuwa ƙima mara komai, don asusun da `claude` na yau da kullum su yi amfani da shigar Keychain ɗaya.

### Me ya sa exec

`syscall.Exec` yana maye gurbin process ɗin ccam da `claude`, don haka ccam ba ya zama a ƙwaƙwalwa, danna maɓallai, signals da exit code na Claude Code ne kansa, kuma arguments suna zama array, ba taɓa zama string na shell ba. Aikin bayan zama, kamar lura da sabobin da aka ƙara da `claude mcp add`, yana faruwa a farawa ta gaba: ccam yana kwatanta lokacin gyaran `.claude.json` da `~/.local/state/ccam/mcp/<name>/seen` kuma yana karanta fayil ɗin ne kawai idan ya canza.

### Yadda ake tura layin umarni

Layukan umarni uku, hanyoyi uku ta `internal/app`:

<!-- motion:launch-route -->

Ana gane gajeren umarni ta `argv[0]`, `ccam <name>` ta argument na farko, kuma duka biyun suna ƙarewa a `syscall.Exec`. Umarnin ccam kamar `ccam usage` ne kawai ke kai wa cobra.

### Sabobin MCP da aka raba

Da tsohuwar dabarar `flag`, ccam yana miƙa sabobin MCP na rukunin ga Claude Code a matsayin argument, wanda aka sa kafin naka:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` yana mayar da argument ne kawai idan `mcp.json` na rukunin yana nan, kuma `Plan.InsertArgs` yana saka shi kai tsaye bayan `claude`. Sannan `--mcp-config` ko `--strict-mcp-config` naka har yanzu yana aiki a saman sa. Argument ɗaya ne, `--mcp-config=<path>`, domin `--mcp-config` na Claude Code yana karɓar ƙima da dama: idan aka rubuta shi a matsayin arguments biyu, zai kuma karanta prompt da ke biye a matsayin fayil ɗin saiti.

Da `mcp_strategy = "sync"` a `config.toml`, ccam yana rubuta sabobin cikin `.claude.json` na asusun maimakon haka, kuma yana tsallake hakan idan asusun yana aiki tuni.

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
- Idan wani Claude Code da ke raba tarihin tattaunawa da wannan asusun (rukuni ɗaya) yana aiki a folda ɗaya, ccam yana gargaɗin cewa tarihin zai iya cakuɗewa kuma yana tambaya `Run anyway? [y/N]`. Ba tare da terminal ba yana nuna gargaɗin kuma yana ci gaba.
- `--dry-run` yana tsallake daidaita MCP (sync), sanarwar sabuwar sabar, duba tarihi, alamar zama, take da lokacin amfani na ƙarshe, kuma ba ya yin exec. Har yanzu yana duba mahaɗan rabawa, kuma yana ƙirƙirar mahaɗin da ya ɓace kamar yadda farawa ta yau da kullum za ta yi.
- Da `mcp_strategy = "sync"`, sync da ya gaza yana buga `ccam: MCP sync failed (launching anyway): …` kuma Claude Code yana farawa da sabobin da yake da su tuni.

## Masu alaƙa

- [Gajerun umarni](/features/shortcuts)
- [Sarrafa asusai](/features/accounts)
- [Bayanin umarnai](/docs/commands)
- [Raba sabobin MCP, saituna da tarihi](/docs/sharing)
- [Yadda CLAUDE_CONFIG_DIR ke aiki](/guides/claude-config-dir)
