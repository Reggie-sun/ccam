---
title: Sanggunian ng mga command
description: Bawat command at option ng ccam, mula sa ccam add at ccs <name> hanggang ccam usage, ccam next at ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 18ecad9a1104
---

Iisang program ang `ccam` at `ccs`; gamitin kung alin ang gusto mo. Ilagay ang `--dry-run` bago ang anumang command para makita kung ano ang gagawin nito nang hindi ito talagang ginagawa, halimbawa `ccs --dry-run work`.

<!-- commands -->

## Pagsisimula ng Claude Code

Sinisimulan ng `ccs <name|number> [arguments]` ang Claude Code bilang account na iyon. Hindi binabasa ng ccam ang mga argument na kasunod ng pangalan; ipinapasa nito ang mga iyon sa `claude` sa parehong pagkakasunod-sunod, at kumikilos ang standard input, exit code at mga signal nang eksaktong gaya ng sa `claude` mismo. Opsyonal ang `--` pagkatapos ng pangalan.

Ang `ccs` na may mga argument pero walang pangalan ng account, gaya ng `ccs -p "summarise"`, ay gumagamit ng default na account. Itakda ito gamit ang `ccam use <name>`.

Ang mga shortcut na `claude1`, `claude2` at iba pa ay mga link sa ccam sa `~/.local/bin`. Gumagawa rin ang `ccam shims install --names` ng `claude-<name>` para sa bawat account.

## Pagdagdag at pag-alis ng mga account

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

Kapag walang folder, hinahanap ng `ccam import` ang `~/.claude` at `~/.claude-*` at nagtatanong tungkol sa bawat isa. Nirerehistro ang mga in-import na folder kung nasaan sila; hindi kailanman naglilipat o kumokopya ng login ang ccam.

## Sharing

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Tumatanggap ang bawat item ng `on` (ang `default` group), `off`, o pangalan ng group. Tingnan ang [sharing](/docs/sharing).

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
ccs usage            # also: ccs u
ccs usage --watch    # refresh every 60 seconds
ccs usage --json
ccs next --continue
ccs next --auto
```

## Maintenance

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## Mga file

| Path | Nilalaman |
| --- | --- |
| `~/.config/ccam/config.toml` | Mga account at kung ano ang shine-share nila |
| `~/.local/share/ccam/profiles/<name>/` | Ang `CLAUDE_CONFIG_DIR` ng bawat account |
| `~/.local/share/ccam/groups/<group>/` | Mga shared na setting, MCP server, extension at kasaysayan |
| `~/.local/state/ccam/` | Usage cache at mga oras ng huling paggamit |
| `~/.local/bin/claude1`, `ccs` | Mga link sa ccam |

Inililipat ng mga variable na `XDG_CONFIG_HOME`, `XDG_DATA_HOME` at `XDG_STATE_HOME` ang mga folder na ito.
