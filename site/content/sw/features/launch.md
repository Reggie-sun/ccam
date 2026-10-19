---
title: Kuanzisha Claude Code
description: Anzisha Claude Code kama akaunti yoyote kwa ccam <name|number>, ccam run au kichaguzi cha ccam; hoja huenda kwa claude bila kubadilishwa, na --dry-run huonyesha amri.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: d033da86f44b
---

`ccam work` huanzisha Claude Code kama akaunti `work`, na kila kitu kinachokuja baada ya jina huenda kwa `claude` kama kilivyo. Itumie popote ambapo ungecharaza `claude`: kwenye terminal, katika script, au ingizo likipitishwa kwa pipe.

## Matumizi

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

Toa jina la akaunti au namba yake kutoka `ccam list`. ccam haisomi hoja zinazokuja baada ya jina, kwa hiyo `ccam work --help` huonyesha msaada wa Claude Code. `--` baada ya jina ni hiari. Standard input, exit code na signals hufanya kazi kama kwa `claude`.

Hoja ya kwanza inapoanza na `-` na si mojawapo ya chaguo za ccam yenyewe, ccam huanzisha akaunti chaguo-msingi pamoja na hoja zote. Chaguo-msingi ni akaunti iliyowekwa kwa `ccam use`, au akaunti yenye namba ya chini zaidi ikiwa hujaweka moja.

### Kichaguzi

`ccam` bila hoja hufungua kichaguzi wakati standard input na output ni terminal, na vinginevyo huchapisha `ccam list`. Huonyesha namba ya kila akaunti, barua pepe, matumizi ya saa 5, matumizi ya mwisho na vipengele vinavyoshirikiwa; matumizi hujazwa kila akaunti inapojibu.

| Kitufe | Inachofanya |
| --- | --- |
| `Enter` | Huanzisha Claude Code kama akaunti iliyochaguliwa |
| `1` hadi `9` | Huanzisha akaunti yenye namba hiyo |
| `Up`/`Down`, `k`/`j` | Husogeza uteuzi |
| `a` | Hufunga kichaguzi na huendesha mchawi wa `ccam add` |
| `s` | Huhariri vipengele vinavyoshirikiwa vya akaunti iliyochaguliwa (`Space` au `x` hubadilisha, `Enter` hutekeleza, `Esc` hurudi nyuma) |
| `d` | Hufanya akaunti iliyochaguliwa kuwa chaguo-msingi |
| `r` | Huleta matumizi tena, bila cache fupi |
| `q`, `Esc`, `Ctrl+C` | Hutoka |

Kichaguzi hujifunga kabla ya kutenda, kwa hiyo terminal hurejeshwa kabla Claude Code haijaanza. Baada ya `d` au badiliko la kushiriki hufunguka tena. Badiliko la kushiriki huonyesha mpango na huuliza kwanza, kama `ccam share` inavyofanya.

## Chaguo

| Chaguo | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `--dry-run` | Huchapisha mazingira na amri ambayo ccam ingeendesha, kisha hutoka | imezimwa |
| `-v`, `--version` | Huchapisha `ccam <version>` (kama hoja ya kwanza tu) | |
| `-h`, `--help` | Msaada wa ccam, kama `ccam --help` au `ccam run --help` | |

Kanuni za hoja:

- ccam husoma chaguo zake yenyewe kabla ya jina la akaunti tu. Baada ya jina, kila hoja ni ya Claude Code, pamoja na `--dry-run` na `--help`.
- `--dry-run` hufanya kazi kama `ccam --dry-run work`, `ccam run --dry-run work` na `ccam --dry-run` (akaunti chaguo-msingi). Inaheshimiwa, haikataliwi.
- `ccam run <name|number> [claude args]` ni muundo mrefu wa `ccam <name>`. `ccam run` bila jina hutumia akaunti chaguo-msingi.
- Majina ya akaunti hayawezi kuwa majina ya amri za ccam kama `add`, `list`, `run` au `usage`, kwa hiyo `ccam <word>` kamwe haina utata.

Jaribio la dry run huchapisha mstari mmoja ambao unaweza kubandika kwenye shell:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

## Jinsi inavyofanya kazi

Kila uanzishaji hupitia hatua zile zile, kuanzia mstari wa amri hadi `exec`:

<!-- motion:launch -->

Kila kisanduku huwaka hatua yake inapoendeshwa, na mistari iliyo chini ni mazingira na hoja ambazo Claude Code huanza nazo.

1. `cmd/ccam/main.go` huita `app.Main`. Ikiwa ccam ilianzishwa kwa jina lingine, ni [njia fupi](/features/shortcuts) na huenda moja kwa moja kwa launcher.
2. `app.parse` huangalia `os.Args` kabla cobra haijahusika. `ccam` peke yake (kichaguzi), jina la amri ya ccam (`config.Reserved`), `-h`/`--help` na `-v`/`--version` pekee ndizo huenda kwa `cli.Execute`. Kila kitu kingine ni uanzishaji.
3. `launcher.Launch` husoma `~/.config/ccam/config.toml`, faili pekee ambayo huichanganua.
4. `launcher.Prepare` hupata akaunti kwa jina au namba, huondoa `--` moja ya mwanzo, na hutafuta programu ya `claude` kwa `ResolveClaude`. Utafutaji katika PATH huhifadhiwa kwenye cache katika `~/.local/state/ccam/claude-path` pamoja na muda wa mwisho wa kubadilishwa kwa faili, kwa hiyo ccam haitafuti `PATH` kila inapoanzisha. Faili iliyobadilika au inayokosekana husababisha utafutaji mpya. Njia ya `claude_bin` katika `config.toml` hutumiwa kama ilivyo. ccam hukataa `claude` ambayo ni ccam yenyewe, ikilinganisha kwa inode, kwa sababu hilo lingesababisha mzunguko usioisha.
5. `checkLinks` hulinganisha folda ya akaunti `~/.local/share/ccam/profiles/<name>/` na mipangilio yake ya kushiriki kwa `share.Check`, ambayo huita `lstat` na `readlink` tu, system call moja au mbili kwa kila kipengele. Kiungo kinachokosekana huundwa papo hapo. Chochote kinachoweza kupoteza data husimamisha uanzishaji.
6. `applyMCP` huongeza seva za MCP zinazoshirikiwa (tazama hapa chini).
7. `session.Register` hurekodi kipindi. ccam huandika alama yenye jina la PID yake, ikiwa na muda wa kuanza wa mchakato, kwenye `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` (au `~/.local/state/ccam/locks/` bila `XDG_RUNTIME_DIR`). Kwa akaunti inayoshiriki historia ya mazungumzo, `historyCheck` kwanza hutafuta kipindi kinachoendeshwa katika `history-<group>-<hash>`, ambapo hash ni ya folda ya sasa, na husajili hapo pia. `exec` huhifadhi PID, kwa hiyo baadaye alama hutaja mchakato wa Claude Code wenyewe. Ni alama, si kufuli inayoshikiliwa, kwa sababu michakato ambayo Claude Code huanzisha ingerithi kufuli na kuiendelea kuishikilia baada ya kipindi kuisha.
8. `Plan.Exec` huweka kichwa cha terminal kuwa `claude · <name>` (wakati standard error ni terminal), huandika muda wa uanzishaji kwenye `~/.local/state/ccam/last-used/<name>` bila `fsync`, na huita `syscall.Exec`.

Claude Code huendeshwa ikiwa `CLAUDE_CONFIG_DIR` imewekwa kuwa folda ya akaunti na `CCAM_PROFILE` imewekwa kuwa jina la akaunti. ccam huondoa thamani zozote za vigezo hivi ambazo ilirithi. Kwa folda ya `~/.claude` iliyosajiliwa kwa `ccam import` kwenye macOS, pia huweka `CLAUDE_SECURESTORAGE_CONFIG_DIR` kuwa thamani tupu, ili akaunti na `claude` ya kawaida zitumie kuingia kule kule kwa Keychain.

### Kwa nini exec

`syscall.Exec` hubadilisha mchakato wa ccam kwa `claude`, kwa hiyo ccam haibaki kwenye kumbukumbu, vibonyezo vya vitufe, signals na exit code ni vya Claude Code yenyewe, na hoja hubaki kuwa safu (array), kamwe si mfuatano wa shell. Kazi baada ya kipindi, kama kugundua seva zilizoongezwa kwa `claude mcp add`, hufanyika katika uanzishaji unaofuata: ccam hulinganisha muda wa kubadilishwa wa `.claude.json` na `~/.local/state/ccam/mcp/<name>/seen` na huichanganua faili tu ilipobadilika.

### Jinsi mstari wa amri unavyoelekezwa

Mistari mitatu ya amri, njia tatu kupitia `internal/app`:

<!-- motion:launch-route -->

Njia fupi hutambuliwa kwa `argv[0]`, `ccam <name>` kwa hoja ya kwanza, na zote mbili huishia katika `syscall.Exec`. Amri ya ccam kama `ccam usage` pekee ndiyo hufika kwa cobra.

### Seva za MCP zinazoshirikiwa

Kwa mkakati wa chaguo-msingi wa `flag`, ccam hupitisha seva za MCP za kikundi kwa Claude Code kama hoja, ikiwekwa kabla ya zako:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` hurudisha hoja tu ikiwa `mcp.json` ya kikundi ipo, na `Plan.InsertArgs` huiweka mara tu baada ya `claude`. `--mcp-config` au `--strict-mcp-config` yako mwenyewe basi bado inatumika juu yake. Ni hoja moja, `--mcp-config=<path>`, kwa sababu `--mcp-config` ya Claude Code hupokea thamani kadhaa: ikiandikwa kama hoja mbili, ingesoma pia prompt inayofuata kama faili ya usanidi.

Kwa `mcp_strategy = "sync"` katika `config.toml`, ccam huandika seva katika `.claude.json` ya akaunti badala yake, na huruka hilo ikiwa akaunti tayari inaendeshwa.

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
- Ikiwa Claude Code nyingine inayoshiriki historia ya mazungumzo na akaunti hii (kikundi kile kile) inaendeshwa katika folda ile ile, ccam huonya kwamba historia inaweza kuchanganyika na huuliza `Run anyway? [y/N]`. Bila terminal huonyesha onyo na huendelea.
- `--dry-run` huruka usawazishaji wa MCP, taarifa ya seva mpya, ukaguzi wa historia, alama ya kipindi, kichwa na muda wa matumizi ya mwisho, na haifanyi exec. Bado hukagua viungo vya kushiriki, na huunda kiungo kinachokosekana kama uanzishaji wa kawaida ungefanya.
- Kwa `mcp_strategy = "sync"`, usawazishaji ulioshindwa huchapisha `ccam: MCP sync failed (launching anyway): …` na Claude Code huanza na seva ilizokuwa nazo tayari.

## Zinazohusiana

- [Amri za njia fupi](/features/shortcuts)
- [Kusimamia akaunti](/features/accounts)
- [Marejeo ya amri](/docs/commands)
- [Shiriki seva za MCP, mipangilio na historia](/docs/sharing)
- [Jinsi CLAUDE_CONFIG_DIR inavyofanya kazi](/guides/claude-config-dir)
