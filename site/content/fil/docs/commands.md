---
title: Sanggunian ng mga command
metaTitle: Sanggunian ng mga command ng ccam
description: Bawat command at option ng ccam, mula sa ccam add at ccam <name> hanggang ccam usage, ccam next at ccam doctor.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: da2594c672bc
---

Ilagay ang `--dry-run` bago ang isang command para makita kung ano ang gagawin nito nang hindi ito talagang ginagawa, halimbawa `ccam --dry-run work`. Gumagana ito kapag sinisimulan ang Claude Code (`ccam <name>`, `ccam run`, `ccam next`) at sa `share`, `migrate`, `shortcut`, `update` at `doctor --fix`. Tinatanggihan ito ng ibang command na may binabago sa halip na tumakbo.

<!-- commands -->

## Pagsisimula ng Claude Code

Sinisimulan ng `ccam <name|number> [arguments]` ang Claude Code bilang account na iyon. Hindi binabasa ng ccam ang mga argument na kasunod ng pangalan; ipinapasa nito ang mga iyon sa `claude` sa parehong pagkakasunod-sunod, at kumikilos ang standard input, exit code at mga signal nang eksaktong gaya ng sa `claude` mismo. Opsyonal ang `--` pagkatapos ng pangalan.

Ang `ccam` na may mga argument pero walang pangalan ng account, gaya ng `ccam -p "summarise"`, ay gumagamit ng account na itinakda para sa kasalukuyang folder, o kung wala, ng default na account. Itakda ang mga ito gamit ang `ccam folder <name> [folder]` at `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

Para sa mga `-p`/`--print` run ang `--fallback` bago ang pangalan ng account. Kapag naabot ng account ang usage limit nito, pinapatakbo ulit ng ccam ang parehong command bilang naka-log in na account na may pinakamaraming natitira at ipi-print lang ang output ng run na natapos. Tingnan ang [Mga account ayon sa folder](/features/folder) at [Patakbuhin ulit sa limit](/features/next#run-again-at-the-limit).

Sa parehong paraan sinisimulan ng isang shortcut gaya ng `ccwork` ang account nito: ginagawa ng `ccwork --resume` ang ginagawa ng `ccam work --resume`. Tingnan ang [mga shortcut](/docs/commands#shortcuts).

## Pagdagdag at pag-alis ng mga account

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

Kapag walang folder, hinahanap ng `ccam import` ang `~/.claude` at `~/.claude-*` at nagtatanong tungkol sa bawat isa. Nirerehistro ang mga in-import na folder kung nasaan sila; hindi kailanman naglilipat o kumokopya ng login ang ccam.

Nagdadagdag ang `ccam add --api` ng account na gumagamit ng API key sa halip na Claude login; binabasa ng ccam ang key sa bawat pagsisimula at hindi ito kailanman iniimbak. Tingnan ang [Mga API key account](/features/api).

Tumatanggi ang `ccam rename`, `ccam remove` at `ccam share` na baguhin ang isang account habang tumatakbo ang Claude Code bilang account na iyon. Isara muna ang session.

Ipinapakita ng `ccam list --json` ang bawat account kasama ang slot, shortcut, email at folder nito. Ang `share` object nito ay nagtutugma sa bawat item (`mcp`, `settings`, `memory`, `extensions`, `history`) sa pangalan ng group o sa `local`.

## Mga shortcut

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Ang shortcut ay isang command na ikaw mismo ang nagpangalan at nagsisimula ng Claude Code bilang isang account. Gaya ng `ccam work`, ipinapasa nito nang walang pagbabago sa Claude Code ang bawat argument, kaya gumagana ang `ccwork --resume` at `ccwork -p "hi"` gaya ng inaasahan mo. Isa itong link sa ccam sa `~/.local/bin`.

Opsyonal ang mga shortcut. Nagtatanong ang wizard ng `ccam add` kung gagawa ng isa, gumagawa ng isa ang `ccam add <name> --shortcut <command>` nang hindi nagtatanong, at ipinapakita ng `ccam list` ang shortcut ng bawat account, o `-` kapag wala.

Puwedeng gamitin sa mga pangalan ang mga character na a-z, 0-9, `.`, `-` at `_`. Hindi puwedeng maging ccam command ang pangalan, shortcut ng ibang account, o pangalan ng ibang program sa iyong `PATH`.

Pinapanatili ng `ccam rename` ang pangalan ng shortcut, at binubura ng `ccam remove` ang shortcut kasama ng account. Ginagawa muli ng `ccam shims install` ang mga link ng lahat ng account, halimbawa pagkatapos mong ilipat ang ccam.

## Sharing

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Tumatanggap ang bawat item ng `on` (ang `default` group), `off`, o pangalan ng group. Tingnan ang [sharing](/docs/sharing).

## Pagkopya ng dati nang setup

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

Kinokopya ng `ccam migrate` ang `settings.json`, `CLAUDE.md`, mga extension at user MCP server mula sa isang Claude Code folder papunta sa isang group, kaya makukuha ang mga ito ng bawat account na nagshe-share ng mga item na iyon. Binabasa lang ang folder. Kinokopya lang ng mga option na `--settings`, `--memory`, `--extensions`, `--mcp` at `--history` ang mga item na pinangalanan mo. Kung wala ang mga ito, sa terminal ay may checklist para makapili ka, na naka-check ang lahat maliban sa kasaysayan ng usapan; sa ibang lugar, kinokopya ang lahat maliban sa kasaysayan ng usapan. Pinipili ng `-g` o `--group` ang group, at kumokopya ang `-y` o `--yes` nang hindi nagtatanong at pinapanatili ang kopya ng group kapag magkaiba ang isang file. Kinokopya rin nito ang mga naka-save na sign-in ng mga server sa sariling login file ng bawat account. Tingnan ang [pagkopya ng dati nang setup](/docs/sharing#copy-an-existing-setup).

## Mga MCP server

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

Inililipat ng `ccam mcp import <account>` sa group nito ang mga server na ang account na iyon lang ang mayroon. Nilalaktawan ang server na ang pangalan ay ginagamit na ng group para sa ibang definition; idagdag ang `--force` para palitan ang kopya ng group, na bina-back up muna ng ccam.

## Usage at paglipat

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
```

Pareho ang paraan ng pagpili ng `next --auto`, `--fallback` at `ccam pick`: inuuna ang mga account na mas mababa sa `auto_max_usage` (`config.toml`, default 100), ayon sa `auto_strategy` (`headroom` o `expiring`), binibilang ang lingguhang window ng `--model`, at hindi kailanman pumipili ng account na nasa 100%. Kapag walang mapipili, lumalabas sila na may code 75, o sa `--wait` ay naghihintay sa unang account na mag-reset. Tingnan ang [Paglipat sa susunod na account](/features/next) at ang [Pagpili ng account para sa mga script](/features/pick).

## Maintenance

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

Sinusuri ng `ccam update` ang SHA-256 checksum ng release bago nito palitan ang ccam, at nag-i-install ang `--version v0.2.0` ng partikular na release. Ang ccam na in-install mula sa `.deb` package ay ina-update na lang gamit ang apt. Tingnan ang [pag-update](/docs/install#update).

## Mga file

| Path | Nilalaman |
| --- | --- |
| `~/.config/ccam/config.toml` | Mga account at kung ano ang shine-share nila |
| `~/.local/share/ccam/profiles/<name>/` | Ang `CLAUDE_CONFIG_DIR` ng bawat account |
| `~/.local/share/ccam/groups/<group>/` | Mga shared na setting, MCP server, extension at kasaysayan |
| `~/.local/state/ccam/` | Usage cache at mga oras ng huling paggamit |
| `~/.local/bin/<shortcut>` | Mga shortcut command gaya ng `ccwork`, mga link sa ccam |

Inililipat ng mga variable na `XDG_CONFIG_HOME`, `XDG_DATA_HOME` at `XDG_STATE_HOME` ang mga folder na ito. Sa Windows, ang `~` ay `%USERPROFILE%` at ang mga shortcut ay `<shortcut>.exe`.
