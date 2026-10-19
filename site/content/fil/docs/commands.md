---
title: Sanggunian ng mga command
metaTitle: Sanggunian ng mga command ng ccam
description: Bawat command at option ng ccam, mula sa ccam add at ccam <name> hanggang ccam usage, ccam next at ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 65fb2b8f6617
---

Ilagay ang `--dry-run` bago ang anumang command para makita kung ano ang gagawin nito nang hindi ito talagang ginagawa, halimbawa `ccam --dry-run work`.

<!-- commands -->

## Pagsisimula ng Claude Code

Sinisimulan ng `ccam <name|number> [arguments]` ang Claude Code bilang account na iyon. Hindi binabasa ng ccam ang mga argument na kasunod ng pangalan; ipinapasa nito ang mga iyon sa `claude` sa parehong pagkakasunod-sunod, at kumikilos ang standard input, exit code at mga signal nang eksaktong gaya ng sa `claude` mismo. Opsyonal ang `--` pagkatapos ng pangalan.

Ang `ccam` na may mga argument pero walang pangalan ng account, gaya ng `ccam -p "summarise"`, ay gumagamit ng default na account. Itakda ito gamit ang `ccam use <name>`.

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
```

Kapag walang folder, hinahanap ng `ccam import` ang `~/.claude` at `~/.claude-*` at nagtatanong tungkol sa bawat isa. Nirerehistro ang mga in-import na folder kung nasaan sila; hindi kailanman naglilipat o kumokopya ng login ang ccam.

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

Kinokopya ng `ccam migrate` ang `settings.json`, `CLAUDE.md`, mga extension at user MCP server mula sa isang Claude Code folder papunta sa isang group, kaya makukuha ang mga ito ng bawat account na nagshe-share ng mga item na iyon. Binabasa lang ang folder. Kinokopya lang ng mga option na `--settings`, `--memory`, `--extensions`, `--mcp` at `--history` ang mga item na pinangalanan mo; kung wala ang mga ito, kinokopya ang lahat maliban sa kasaysayan ng usapan. Pinipili ng `-g` o `--group` ang group, at kumokopya ang `-y` o `--yes` nang hindi nagtatanong at pinapanatili ang kopya ng group kapag magkaiba ang isang file. Tingnan ang [pagkopya ng dati nang setup](/docs/sharing#copy-an-existing-setup).

## Mga MCP server

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## Usage at paglipat

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## Maintenance

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
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

Inililipat ng mga variable na `XDG_CONFIG_HOME`, `XDG_DATA_HOME` at `XDG_STATE_HOME` ang mga folder na ito.
