---
title: Paggamit ng Remote Control sa ccam
description: "Simulan ang Remote Control ng Claude Code bilang kahit anong account gamit ang ccam <name> rc, alamin kung aling mga option ang nalalapat, kung bakit nagtatanong ang ccam na kopyahin ang mga shared na MCP server, at kung paano ayusin ang mga karaniwang problema."
updated: 2026-10-20
order: 4
translation: machine
sourceHash: b7b09f9bcb9e
---

Sinisimulan ng `ccam work rc` ang Remote Control ng Claude Code bilang account na `work`, kaya makokontrol mo ang mga session sa computer na ito mula sa claude.ai/code o sa Claude mobile app. Gumagana ito gaya ng `claude rc` na may iisang pagkakaiba: kung nagshe-share ng mga MCP server ang account, magtatanong ang ccam nang isang beses kung kokopyahin ang mga ito sa `.claude.json` ng account, dahil wala nang ibang pinagkukunan ng server ang mga session ng Remote Control.

## Pagsisimula

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

Ganoon din gumagana ang isang [shortcut](/features/shortcuts): ang `ccwork rc` ay `ccam work rc`. Patakbuhin ang command sa folder kung saan mo gustong magtrabaho, saka buksan ang claude.ai/code o ang Claude app na naka-sign in gamit ang parehong Claude account ng `work`.

## Mga option

Dalawang uri ng option ang nagtatagpo rito: ang sariling option ng ccam, na nauuna sa pangalan ng account o nasa settings ng ccam, at ang option ng Remote Control, na kasunod ng `rc` at umaabot sa Claude Code nang walang pagbabago.

### Mga option at setting ng ccam

| Option o setting | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--dry-run` (bago ang account) | I-print ang environment at command, pati ang abiso tungkol sa MCP kung mayroon. Hindi ito nagtatanong o kumokopya | off |
| `Copy them into .claude.json (a backup is kept)? [y/N]` | Itinatanong kapag nagshe-share ng mga MCP server ang account at wala sa `.claude.json` nito ang kasalukuyang listahan ng group. Kinokopya ng `y` ang mga ito; nagsisimula nang wala ang mga ito sa `N` o Enter | `N` |
| `mcp_strategy = "sync"` sa `~/.config/ccam/config.toml` | Kopyahin ang mga server ng group sa `.claude.json` bago ang bawat pagsisimula, Remote Control man o hindi, kaya hindi na kailanman magtatanong | `"flag"` |
| `ccam share work --mcp off` | Itigil ang pag-share ng mga MCP server sa account; pagkatapos ay magsisimula ang Remote Control na may sariling server lang ng account at walang tanong | on |

### Mga option ng Remote Control

Pag-aari ng Claude Code ang lahat ng kasunod ng `rc`. Ito ang mga option ng `claude remote-control` sa Claude Code 2.1.304; ipi-print ng `ccam work rc --help` ang listahan para sa bersyon mo.

| Option | Ano ang ginagawa |
| --- | --- |
| `--name <name>` | Pangalan ng session na ipinapakita sa claude.ai/code |
| `--remote-control-session-name-prefix <prefix>` | Prefix ng mga awtomatikong pangalan ng session (default: ang host name) |
| `-c`, `--continue` | Muling kumonekta sa session na huling itinala ng Remote Control para sa folder na ito, sa loob ng mga 4 na oras |
| `--session-id <id>` | Muling kumonekta sa isang partikular na session |
| `--permission-mode <mode>` | Permission mode ng mga bagong session: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` o `plan` |
| `--spawn <mode>` | `same-dir` (default), `worktree` (isang git worktree bawat session) o `session` (isang session, lalabas kapag natapos ito) |
| `--capacity <N>` | Pinakamaraming sabay-sabay na session sa mode na `same-dir` at `worktree` (default 32) |
| `--[no-]create-session-in-dir` | Gumawa ng unang session sa kasalukuyang folder sa pagsisimula (default on) |
| `--[no-]chrome` | Claude in Chrome para sa mga bagong session (default off) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | Output para sa debug at verbose |

Ilagay ang `rc` kasunod mismo ng pangalan ng account. Itinuturing ng ccam na command ng Claude Code ang unang argument na hindi nagsisimula sa `-`, kaya sa `ccam work rc --name laptop` ay `rc` ang nakikita nito, pero sa `ccam work --model opus rc` ay `opus` ang nakikita nito.

## Ano ang ginagawa ng ccam

Dumadaan ang pagsisimula sa karaniwang landas ([Pagsisimula ng Claude Code](/features/launch)) na may isang pagbabago sa hakbang ng MCP:

1. Hinahanap ng ccam ang account, sinusuri ang mga share link nito at hinahanap ang `claude`, gaya sa anumang pagsisimula.
2. Nakikita nitong `rc` o `remote-control` ang command, kaya hindi nito isinasama ang `--mcp-config`.
3. Kung nagshe-share ng mga MCP server ang account, ikinukumpara nito ang `mcp.json` ng group sa `sync-stamp` sa `~/.local/state/ccam/mcp/work/`. Kapag nasa `.claude.json` na ang kasalukuyang mga server, o walang server ang group, dumidiretso ito sa hakbang 6.
4. Kung hindi, pinapangalanan nito ang mga server na mawawala sa mga session:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

5. Sa `y`, nagse-save ito ng backup ng `.claude.json`, kinukuha ang sariling lock ng Claude Code, idinadagdag ang mga server ng group sa `mcpServers` at pinapanatili ang iba pang server ng account. Sa `N`, kapag walang terminal, o habang tumatakbo ang account, wala itong binabago at sinasabi iyon:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. Itinatala nito ang session, itinatakda ang `CLAUDE_CONFIG_DIR` sa folder ng account at pinapalitan ang sarili nito ng `claude rc` sa pamamagitan ng `exec`.

Pagkatapos ay patuloy na tumatakbo ang Remote Control bilang server. Child process ng prosesong iyon ang bawat session na sinisimulan nito para sa iyo, kaya tumatakbo ito bilang `work`: gamit ang login, settings at kasaysayan ng usapan ng `work`, at ibinabawas sa usage limit ng `work`. Habang tumatakbo ang server, itinuturing ng ccam na tumatakbo ang account.

Makikita ang pagkakaiba sa dry run. Dala ng karaniwang pagsisimula ang mga server ng group bilang argument, pero hindi ng Remote Control:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude rc
```

## Bakit ganito ito gumagana

**Hindi tumatanggap ng `--mcp-config` ang Remote Control.** Sa karaniwang pagsisimula, ibinibigay ng ccam ang mga server ng group sa Claude Code bilang `--mcp-config=<path>` at hindi kailanman nagsusulat sa mga file ng Claude Code. Tinatanggihan ng `claude remote-control` ang option na iyon bago ang pangalan ng command nito at wala itong tinatanggap pagkatapos, kaya bago ang ccam 0.6.0, pumapalya ang `ccam work rc` sa bawat account na nagshe-share ng mga MCP server. Hindi na ngayon isinasama ng ccam ang option.

**Sa `.claude.json` lang kumukuha ng server ang mga session nito.** Hindi namamana ng mga session na sinisimulan ng Remote Control ang mga option mula sa command line. Ang `.claude.json` ng account ang tanging lugar kung saan nila nakikita ang mga user MCP server. Hawak din ng file na iyon ang mga detalye ng login ng account at ang state ng bawat folder, kaya hindi ito mai-link ng ccam sa isang shared na kopya gaya ng ginagawa nito sa `settings.json`. Ang pagkopya ng mga server dito ang tanging paraan para maibigay sa mga session ang shared na listahan.

**Nagtatanong ang ccam bago magsulat.** Sa default na flag strategy, nangangako ang ccam na hindi gagalawin ang `.claude.json`. Ang Remote Control ang tanging kaso kung saan kailangan nitong magsulat doon, kaya nagtatanong muna ito, nagtatabi ng backup, at walang ginagalaw maliban sa `mcpServers`.

**Hindi ito nagsusulat habang tumatakbo ang account.** Isinusulat ulit ng Claude Code ang `.claude.json` habang may session, at mawawala o mapapatungan ang pagbabagong ginawa mula sa labas. Kaya nilalaktawan ng ccam ang pagkopya kapag may tumatakbong session ng account, kasama ang isa pang Remote Control server.

**Isang beses itong nagtatanong sa bawat pagbabago.** Pagkatapos ng kopya, hawak ng `sync-stamp` ang modification time ng `mcp.json` ng group. Muling magtatanong ang ccam kapag nagbago lang ang file na iyon, halimbawa pagkatapos ng `ccam mcp add`. Kung `N` ang sagot mo, wala pang kopya, kaya magtatanong ito sa bawat pagsisimula ng Remote Control.

**Pinapanatiling bago ng `ccam mcp` ang mga kopya.** Nananaig ang entry sa `.claude.json` sa `--mcp-config` kapag nagkapareho ang pangalan, kaya matatakpan ng lumang kopya ang mas bagong definition kahit sa karaniwang pagsisimula. Kaya ina-update ng `ccam mcp add`, `add-json`, `remove` at `import` ang mga kopya ng bawat account na hindi tumatakbo:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**Kapag walang terminal, hindi ang sagot.** Madalas tumatakbo ang Remote Control mula sa script, service o `nohup`, kung saan walang makakasagot. Sa ganoong kaso, nagsisimula ang ccam nang wala ang mga server sa halip na maghintay.

## Mga problema at solusyon

### Wala ang mga shared na MCP server sa mga session ng Remote Control

`N` ang naging sagot sa tanong, walang terminal, o tumatakbo ang account. Itigil ang Remote Control, isara ang iba pang session ng account, at patakbuhin ulit ang `ccam work rc` sa isang terminal; sumagot ng `y`. Para hindi na kailanman tanungin, itakda ang `mcp_strategy = "sync"` sa `~/.config/ccam/config.toml`. Hindi nagbabago ang mga server na nasa session na habang tumatakbo ito, kaya magsimula ng bagong session pagkatapos ng kopya.

### "Account work is running, so its .claude.json is left alone"

May bukas pang session ng `work`: isang karaniwang `ccam work`, isang shortcut, o isa pang Remote Control server. Lumabas dito at patakbuhin ulit ang `ccam work rc`. Hindi binibilang ang mga session na pinapatakbo ng ibang account.

### Nagtatanong ang ccam sa bawat pagkakataon

Hindi mo pa nakokopya ang mga server. Sumagot ng `y` nang isang beses, o itakda ang `mcp_strategy = "sync"`. Kung ayaw mo talaga ng mga server sa Remote Control, itigil ang pag-share ng MCP sa account: `ccam share work --mcp off`.

### Pumapalya agad ang `ccam work rc` na may error tungkol sa `--mcp-config`

Mas luma sa 0.6.0 ang ccam mo at ipinapasa nito ang `--mcp-config` sa Remote Control. Patakbuhin ang `ccam update`. Kung pareho ang error sa kasalukuyang ccam, ibig sabihin ay may option na may value na nauna sa `rc`, gaya ng sa `ccam work --model opus rc`, kaya inakala ng ccam na command ang `opus`; ilagay ang `rc` kasunod mismo ng pangalan ng account.

### Hindi umaabot sa Remote Control ang pagbabagong ginawa gamit ang `ccam mcp`

Ina-update lang ng `ccam mcp add` at `remove` ang mga kopya sa mga account na hindi tumatakbo, at sinasabi nila iyon: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Itigil ang Remote Control at simulan ulit; makikita ng ccam ang binagong `mcp.json` at magtatanong ulit.

### Pinalitan ng kopya ang sarili kong server

Kung may sariling server ang account na kapareho ng pangalan ng isang server ng group, pinalitan ito ng kopya. Nasa `~/.local/state/ccam/mcp/work/backups/` ang dating `.claude.json`, na nagtatabi ng limang pinakabago. Para mapanatili ang bersyon mo para sa bawat account, ilipat ito sa group gamit ang `ccam mcp import work <server> --force` bago ang susunod na pagsisimula ng Remote Control.

### Gusto kong alisin ang mga kinopyang server

Inaalis ng `ccam mcp remove <server>` ang isang server mula sa group at mula sa bawat kopya. Para bumalik sa estado bago ang kopya, i-restore ang pinakabagong backup mula sa `~/.local/state/ccam/mcp/work/backups/` habang walang tumatakbong session ng account.

### Nagtatanong ang Remote Control kung pagkakatiwalaan ang folder, o lumalabas ito kapag walang terminal

Iniimbak bawat account sa `.claude.json` ang pagtitiwala sa folder, kaya bago para sa `work` ang folder na pinagkatiwalaan mo bilang `personal`. Patakbuhin ang `ccam work` nang isang beses sa folder na iyon at tanggapin, saka simulan doon ang Remote Control.

### "You must be logged in with a Claude account that has a subscription"

Kailangan ng Remote Control ng login na may Claude subscription. Tingnan ang email ng account sa `ccam list`, at mag-sign in ulit gamit ang `ccam login work` kung kailangan.

### Hindi lumalabas ang session sa telepono

Naka-sign in ang Claude app o claude.ai/code gamit ang ibang Claude account. Ikumpara ito sa email na ipinapakita ng `ccam list` para sa `work`.

### Nagbababala ang ccam na may ibang Claude Code na nagshe-share ng kasaysayan ng usapan sa folder na ito

May ibang account sa parehong history group na tumatakbo sa parehong folder. Puwedeng maghalo ang kasaysayan kapag ipinagpatuloy ng dalawang account ang iisang usapan. Sumagot ng `y` kung magkaibang usapan lang ang hahawakan nila, o simulan ang Remote Control sa ibang folder.

## Kaugnay

- [Pagsisimula ng Claude Code](/features/launch)
- [Mga shared na MCP server](/features/mcp)
- [Sharing sa pagitan ng mga account](/features/share)
- [Pamamahala ng mga account](/features/accounts)
