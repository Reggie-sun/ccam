---
title: Pagsisimula ng Claude Code
description: Simulan ang Claude Code bilang anumang account gamit ang ccam <name|number>, ccam run o ang picker ng ccam; ipinapasa nang walang pagbabago sa claude ang mga argument, at ipinapakita ng --dry-run ang command.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 6c04bb48df22
---

Sinisimulan ng `ccam work` ang Claude Code bilang account na `work`, at ipinapasa sa `claude` nang walang pagbabago ang lahat ng kasunod ng pangalan. Gamitin ito saanman mo ita-type ang `claude`: sa terminal, sa script, o kapag may input na ipina-pipe.

## Paggamit

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

Ibigay ang pangalan ng account o ang numero nito mula sa `ccam list`. Hindi binabasa ng ccam ang mga argument na kasunod ng pangalan, kaya ipinapakita ng `ccam work --help` ang help ng Claude Code. Opsyonal ang `--` pagkatapos ng pangalan. Kumikilos ang standard input, exit code at mga signal gaya ng sa `claude`.

Kapag nagsisimula sa `-` ang unang argument at hindi ito isa sa sariling flag ng ccam, sinisimulan ng ccam ang default na account kasama ang lahat ng argument. Ang default ay ang account na itinakda gamit ang `ccam use`, o ang account na may pinakamababang numero kung wala kang itinakda.

### Ang picker

Binubuksan ng `ccam` na walang argument ang picker kapag terminal ang standard input at output, at ipi-print ang `ccam list` kung hindi. Ipinapakita nito ang numero, email, 5-oras na usage, huling paggamit at mga shared na item ng bawat account; napupunan ang usage habang sumasagot ang bawat account.

| Key | Ano ang ginagawa |
| --- | --- |
| `Enter` | Simulan ang Claude Code bilang ang napiling account |
| `1` hanggang `9` | Simulan ang account na may numerong iyon |
| `Up`/`Down`, `k`/`j` | Ilipat ang pagpili |
| `a` | Isara ang picker at patakbuhin ang wizard ng `ccam add` |
| `s` | Baguhin ang mga shared na item ng napiling account (nagto-toggle ang `Space` o `x`, naglalapat ang `Enter`, bumabalik ang `Esc`) |
| `d` | Gawing default ang napiling account |
| `r` | Kunin ulit ang usage, nang hindi ginagamit ang maikling cache |
| `q`, `Esc`, `Ctrl+C` | Lumabas |

Nagsasara ang picker bago ito kumilos, kaya naibabalik ang terminal bago magsimula ang Claude Code. Pagkatapos ng `d` o ng pagbabago sa sharing, bumubukas itong muli. Ipinapakita muna ng pagbabago sa sharing ang plano at nagtatanong, gaya ng `ccam share`.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--dry-run` | I-print ang environment at command na patatakbuhin ng ccam, saka lumabas | off |
| `-v`, `--version` | I-print ang `ccam <version>` (bilang unang argument lang) | |
| `-h`, `--help` | Ang help ng ccam, bilang `ccam --help` o `ccam run --help` | |

Mga patakaran sa argument:

- Binabasa lang ng ccam ang sarili nitong mga flag bago ang pangalan ng account. Pagkatapos ng pangalan, sa Claude Code na ang bawat argument, kasama ang `--dry-run` at `--help`.
- Gumagana ang `--dry-run` bilang `ccam --dry-run work`, `ccam run --dry-run work` at `ccam --dry-run` (ang default na account). Sinusunod ito, hindi tinatanggihan.
- Ang `ccam run <name|number> [claude args]` ay ang mahabang anyo ng `ccam <name>`. Ginagamit ng `ccam run` na walang pangalan ang default na account.
- Hindi puwedeng pangalan ng ccam command ang mga pangalan ng account, gaya ng `add`, `list`, `run` o `usage`, kaya hindi kailanman malabo ang `ccam <word>`.

Nagpi-print ang dry run ng isang linya na puwede mong idikit sa shell:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

## Paano ito gumagana

Pareho ang mga hakbang ng bawat pagsisimula, mula sa command line hanggang `exec`:

<!-- motion:launch -->

Umiilaw ang bawat kahon habang tumatakbo ang hakbang nito, at ang mga linya sa ilalim ay ang environment at mga argument na kasama ng pagsisimula ng Claude Code.

1. Tinatawag ng `cmd/ccam/main.go` ang `app.Main`. Kung sinimulan ang ccam sa ibang pangalan, isa itong [shortcut](/features/shortcuts) at direktang napupunta sa launcher.
2. Tinitingnan ng `app.parse` ang `os.Args` bago pa makisangkot ang cobra. Ang `ccam` lang na mag-isa (ang picker), ang pangalan ng ccam command (`config.Reserved`), `-h`/`--help` at `-v`/`--version` ang napupunta sa `cli.Execute`. Pagsisimula ang lahat ng iba pa.
3. Binabasa ng `launcher.Launch` ang `~/.config/ccam/config.toml`, ang tanging file na pina-parse nito.
4. Hinahanap ng `launcher.Prepare` ang account ayon sa pangalan o numero, inaalis ang isang `--` sa unahan, at hinahanap ang `claude` executable gamit ang `ResolveClaude`. Naka-cache ang paghahanap sa PATH sa `~/.local/state/ccam/claude-path` kasama ang modification time ng file, kaya hindi hinahanap ng ccam ang `PATH` sa bawat pagsisimula. Naghahanap ulit kapag nagbago o nawala ang file. Ginagamit nang walang pagbabago ang path na `claude_bin` sa `config.toml`. Tinatanggihan ng ccam ang `claude` na ccam mismo, na ikinukumpara ayon sa inode, dahil magdudulot iyon ng loop.
5. Ikinukumpara ng `checkLinks` ang folder ng account na `~/.local/share/ccam/profiles/<name>/` sa sharing settings nito gamit ang `share.Check`, na tumatawag lang ng `lstat` at `readlink`, isa o dalawang system call bawat item. Ginagawa agad ang nawawalang link. Pinapahinto ang pagsisimula ng anumang puwedeng magpawala ng data.
6. Idinadagdag ng `applyMCP` ang mga shared na MCP server (tingnan sa ibaba).
7. Itinatala ng `session.Register` ang session. Nagsusulat ang ccam ng marker na ipinangalan sa PID nito, na naglalaman ng oras ng pagsisimula ng proseso, sa `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` (o `~/.local/state/ccam/locks/` kapag walang `XDG_RUNTIME_DIR`). Para sa account na nagshe-share ng kasaysayan ng usapan, hinahanap muna ng `historyCheck` ang tumatakbong session sa `history-<group>-<hash>`, kung saan ang hash ay sa kasalukuyang folder, at nagrerehistro rin doon. Pinapanatili ng `exec` ang PID, kaya ang marker ay tumutukoy na sa mismong proseso ng Claude Code. Marker, hindi hawak na lock, dahil mamanahin ng mga prosesong sinisimulan ng Claude Code ang lock at hahawakan ito kahit tapos na ang session.
8. Itinatakda ng `Plan.Exec` ang pamagat ng terminal sa `claude · <name>` (kapag terminal ang standard error), isinusulat ang oras ng pagsisimula sa `~/.local/state/ccam/last-used/<name>` nang walang `fsync`, at tinatawag ang `syscall.Exec`.

Tumatakbo ang Claude Code na nakatakda ang `CLAUDE_CONFIG_DIR` sa folder ng account at ang `CCAM_PROFILE` sa pangalan ng account. Inaalis ng ccam ang anumang value ng mga ito na namana nito. Para sa `~/.claude` folder na nirehistro gamit ang `ccam import` sa macOS, itinatakda rin nito sa walang lamang value ang `CLAUDE_SECURESTORAGE_CONFIG_DIR`, para iisang Keychain login ang gamitin ng account at ng karaniwang `claude`.

### Bakit exec

Pinapalitan ng `syscall.Exec` ang proseso ng ccam ng `claude`, kaya hindi nananatili sa memory ang ccam, sa Claude Code mismo ang mga keystroke, signal at exit code, at nananatiling array ang mga argument, hindi kailanman shell string. Ang mga gawain pagkatapos ng session, gaya ng pagpansin sa mga server na idinagdag gamit ang `claude mcp add`, ay nangyayari sa susunod na pagsisimula: ikinukumpara ng ccam ang modification time ng `.claude.json` sa `~/.local/state/ccam/mcp/<name>/seen` at pina-parse lang ang file kapag nagbago ito.

### Paano dinadala ang isang command line

Tatlong command line, tatlong ruta sa `internal/app`:

<!-- motion:launch-route -->

Itinutugma ang shortcut sa `argv[0]`, ang `ccam <name>` sa unang argument, at parehong nagtatapos sa `syscall.Exec`. Ang ccam command lang gaya ng `ccam usage` ang umaabot sa cobra.

### Mga shared na MCP server

Sa default na `flag` strategy, ipinapasa ng ccam sa Claude Code ang mga MCP server ng group bilang argument, na inilalagay bago ang sa iyo:

<!-- motion:launch-mcp -->

Nagbabalik lang ng argument ang `mcp.FlagArgs` kung mayroong `mcp.json` ang group, at inilalagay ito ng `Plan.InsertArgs` kasunod mismo ng `claude`. Sa ganitong paraan, nalalapat pa rin sa ibabaw nito ang sarili mong `--mcp-config` o `--strict-mcp-config`. Isa itong argument, `--mcp-config=<path>`, dahil tumatanggap ng ilang value ang `--mcp-config` ng Claude Code: kapag isinulat bilang dalawang argument, babasahin din nito bilang config file ang kasunod na prompt.

Kapag may `mcp_strategy = "sync"` sa `config.toml`, isinusulat na lang ng ccam ang mga server sa `.claude.json` ng account, at nilalaktawan iyon kung tumatakbo na ang account.

### Remote Control

Sinisimulan ng `ccam work rc` (o `ccam work remote-control`) ang Remote Control ng Claude Code bilang `work`. Tinatanggihan ng `claude remote-control` ang `--mcp-config` bago ang verb nito at wala itong tinatanggap pagkatapos nito, at binabasa lang ng mga session nito ang mga MCP server mula sa `.claude.json` ng account. Kaya iba ang takbo ng flag strategy dito:

<!-- motion:launch-rc -->

1. Itinuturing ng `remoteControl` na verb ng Claude Code ang unang argument pagkatapos ng account na hindi nagsisimula sa `-`. Kapag `remote-control` o `rc` ito, walang idinadagdag ang `mcp.FlagArgs`.
2. Ikinukumpara ng `mcp.Synced` ang modification time ng `mcp.json` ng group sa `sync-stamp` ng account. Kung nasa `.claude.json` na ang kasalukuyang mga server ng group, o walang server ang group, nagsisimula ang Remote Control nang walang tanong.
3. Kung hindi, pinapangalanan ng ccam ang mga server na mawawala sa mga session at nagtatanong:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. Sa `y`, kinokopya ng `mcp.Sync` ang mga ito gaya ng ginagawa ng `sync` strategy: sa ilalim ng lock ng Claude Code, pagkatapos ng backup, at pinapanatili ang sariling server ng account. Sa `N`, o kapag walang terminal, nagsisimula ang Remote Control nang wala ang mga ito, at itinuturo ng ccam ang `mcp_strategy = "sync"`, na kumokopya sa mga ito bago ang bawat pagsisimula.

Nagtatanong ang ccam sa bawat pagsisimula ng Remote Control hanggang makopya ang mga server, at pagkatapos ng kopya ay kapag nagbago lang ulit ang `mcp.json` ng group. Nananaig ang kopya sa `.claude.json` sa `--mcp-config` sa mga susunod na normal na pagsisimula, kaya ina-update ng `ccam mcp add`, `remove` at `import` ang mga kopyang hawak ng ibang account ([Mga shared na MCP server](/features/mcp)).

### Ang 10 ms na budget

Hindi hihigit sa 10 ms (p95) ang idinadagdag ng ccam bago magsimula ang Claude Code. Hindi kailanman binubuo ng launch path ang command tree ng cobra, sinisimulan ang Bubble Tea o nagbubukas ng network connection, at pumapalya ang `TestPackageInitIsCheap` kapag lumampas sa kabuuang 5 ms ang mga function na `init()` ng package. Sinusukat ng `scripts/perf-check.sh` ang `ccam --dry-run work`, `ccam work` at isang shortcut gamit ang `hyperfine` at pumapalya kapag lampas sa 10 ms, at pinapatakbo nang paulit-ulit ng `BenchmarkLaunchPath` sa `internal/launcher` ang mga hakbang. Para makita ang oras ng bawat hakbang, itakda ang `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Kaligtasan at mga espesyal na kaso

- Humihinto ang hindi kilalang pangalan o numero na may `no account "<name>"` at itinuturo ang `ccam list`. Kapag wala pang kahit isang account, hihilingin ng ccam na magdagdag ka gamit ang `ccam add <name>`.
- Kung naging karaniwang file ang isang shared na path o tumuturo ito sa maling lugar, humihinto ang ccam na may `stopped because of share link problems:`, inililista ang mga ito, at hinihiling na patakbuhin mo ang `ccam doctor --fix`. Hindi ito naglilipat ng data sa launch path.
- Kung may ibang Claude Code na nagshe-share ng kasaysayan ng usapan sa account na ito (parehong group) at tumatakbo sa parehong folder, nagbababala ang ccam na puwedeng maghalo ang kasaysayan at nagtatanong ng `Run anyway? [y/N]`. Kapag walang terminal, ipinapakita nito ang babala at nagpapatuloy.
- Nilalaktawan ng `--dry-run` ang MCP sync, ang abiso tungkol sa bagong server, ang pagsusuri ng kasaysayan, ang session marker, ang pamagat at ang oras ng huling paggamit, at hindi ito nag-e-exec. Sinusuri pa rin nito ang mga share link, at ginagawa ang nawawalang link gaya ng sa karaniwang pagsisimula.
- Kapag pinatakbo ang `ccam work rc` habang tumatakbo na ang `work`, hindi nito ginagalaw ang `.claude.json`, sinasabi iyon, at sinisimulan ang Remote Control nang wala ang mga server ng group; patakbuhin ulit pagkatapos lumabas ng session para makopya ang mga ito. Sa `--dry-run`, ipi-print nito ang mensahe pero hindi ito magtatanong o kokopya.
- Sa `mcp_strategy = "sync"`, nagpi-print ang pumalyang sync ng `ccam: MCP sync failed (launching anyway): …` at nagsisimula ang Claude Code kasama ang mga server na mayroon na ito.

## Kaugnay

- [Mga shortcut command](/features/shortcuts)
- [Pamamahala ng mga account](/features/accounts)
- [Sanggunian ng mga command](/docs/commands)
- [I-share ang mga MCP server, settings at kasaysayan](/docs/sharing)
- [Paano gumagana ang CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
