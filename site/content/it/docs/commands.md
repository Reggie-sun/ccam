---
title: Riferimento dei comandi
description: Tutti i comandi e le opzioni di ccam, da ccam add e ccs <name> fino a ccam usage, ccam next e ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 18ecad9a1104
---

`ccam` e `ccs` sono lo stesso programma; usa quello che preferisci. Metti `--dry-run` prima di qualsiasi comando per vedere cosa farebbe senza eseguirlo, per esempio `ccs --dry-run work`.

<!-- commands -->

## Avviare Claude Code

`ccs <name|number> [arguments]` avvia Claude Code con quell'account. ccam non legge gli argomenti che seguono il nome; li passa a `claude` nello stesso ordine, e standard input, codice di uscita e segnali si comportano esattamente come con `claude` stesso. Un `--` dopo il nome è facoltativo.

`ccs` con argomenti ma senza nome di account, come `ccs -p "summarise"`, usa l'account predefinito. Impostalo con `ccam use <name>`.

Le scorciatoie `claude1`, `claude2` e così via sono link a ccam in `~/.local/bin`. `ccam shims install --names` crea anche `claude-<name>` per ogni account.

## Aggiungere e rimuovere account

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` senza una cartella cerca `~/.claude` e `~/.claude-*` e chiede conferma per ciascuna. Le cartelle importate vengono registrate dove si trovano; ccam non sposta né copia mai un login.

## Condivisione

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Ogni elemento accetta `on` (il gruppo `default`), `off` oppure il nome di un gruppo. Vedi [condivisione](/docs/sharing).

## Server MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## Utilizzo e cambio di account

```sh
ccs usage            # also: ccs u
ccs usage --watch    # refresh every 60 seconds
ccs usage --json
ccs next --continue
ccs next --auto
```

## Manutenzione

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## File

| Percorso | Contenuto |
| --- | --- |
| `~/.config/ccam/config.toml` | Gli account e cosa condividono |
| `~/.local/share/ccam/profiles/<name>/` | La `CLAUDE_CONFIG_DIR` di ogni account |
| `~/.local/share/ccam/groups/<group>/` | Impostazioni, server MCP, estensioni e cronologia condivisi |
| `~/.local/state/ccam/` | Cache dell'utilizzo e orari dell'ultimo uso |
| `~/.local/bin/claude1`, `ccs` | Link a ccam |

Le variabili `XDG_CONFIG_HOME`, `XDG_DATA_HOME` e `XDG_STATE_HOME` spostano queste cartelle.
