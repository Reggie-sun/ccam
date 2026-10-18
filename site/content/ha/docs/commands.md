---
title: Bayanin umarnai
description: Kowane umarni da zaɓi na ccam, daga ccam add da ccs <name> zuwa ccam usage, ccam next da ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 18ecad9a1104
---

`ccam` da `ccs` shiri ɗaya ne; yi amfani da wanda ka fi so. Saka `--dry-run` kafin kowane umarni don ganin abin da zai yi ba tare da ya yi shi ba, misali `ccs --dry-run work`.

<!-- commands -->

## Fara Claude Code

`ccs <name|number> [arguments]` yana fara Claude Code a matsayin wannan asusun. ccam ba ya karanta arguments da ke bayan sunan; yana miƙa su ga `claude` a tsari ɗaya, kuma standard input, exit code da signals suna aiki daidai kamar da `claude` kansa. `--` bayan sunan ba dole ba ne.

`ccs` tare da arguments amma ba tare da sunan asusu ba, kamar `ccs -p "summarise"`, yana amfani da asusun tsoho. Saita shi da `ccam use <name>`.

Gajerun umarnin `claude1`, `claude2` da sauransu mahaɗai ne zuwa ccam a `~/.local/bin`. `ccam shims install --names` yana kuma ƙirƙirar `claude-<name>` ga kowane asusu.

## Ƙarawa da cire asusai

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` ba tare da folda ba yana neman `~/.claude` da `~/.claude-*` kuma yana tambaya game da kowanne. Ana yin rajistar foldojin da aka shigo da su a inda suke; ccam ba ya taɓa motsawa ko kwafar bayanan shiga.

## Rabawa

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Kowane abu yana karɓar `on` (rukunin `default`), `off`, ko sunan rukuni. Duba [rabawa](/docs/sharing).

## Sabobin MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## Amfani da sauyawa

```sh
ccs usage            # also: ccs u
ccs usage --watch    # refresh every 60 seconds
ccs usage --json
ccs next --continue
ccs next --auto
```

## Kulawa

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## Fayiloli

| Hanya | Abin da ke ciki |
| --- | --- |
| `~/.config/ccam/config.toml` | Asusai da abin da suke rabawa |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` na kowane asusu |
| `~/.local/share/ccam/groups/<group>/` | Saituna, sabobin MCP, extensions da tarihi da aka raba |
| `~/.local/state/ccam/` | Cache na amfani da lokutan amfani na ƙarshe |
| `~/.local/bin/claude1`, `ccs` | Mahaɗai zuwa ccam |

Masu canji (variables) `XDG_CONFIG_HOME`, `XDG_DATA_HOME` da `XDG_STATE_HOME` suna canza wurin waɗannan foldojin.
