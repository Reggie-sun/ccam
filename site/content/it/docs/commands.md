---
title: Riferimento dei comandi
description: Tutti i comandi e le opzioni di ccam, da ccam add e ccam <name> fino a ccam usage, ccam next e ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 0c869c92fa63
---

Metti `--dry-run` prima di qualsiasi comando per vedere cosa farebbe senza eseguirlo, per esempio `ccam --dry-run work`.

<!-- commands -->

## Avviare Claude Code

`ccam <name|number> [arguments]` avvia Claude Code con quell'account. ccam non legge gli argomenti che seguono il nome; li passa a `claude` nello stesso ordine, e standard input, codice di uscita e segnali si comportano esattamente come con `claude` stesso. Un `--` dopo il nome è facoltativo.

`ccam` con argomenti ma senza nome di account, come `ccam -p "summarise"`, usa l'account predefinito. Impostalo con `ccam use <name>`.

Una scorciatoia come `ccwork` avvia il suo account allo stesso modo: `ccwork --resume` fa ciò che fa `ccam work --resume`. Vedi [scorciatoie](/docs/commands#shortcuts).

## Aggiungere e rimuovere account

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` senza una cartella cerca `~/.claude` e `~/.claude-*` e chiede conferma per ciascuna. Le cartelle importate vengono registrate dove si trovano; ccam non sposta né copia mai un login.

## Scorciatoie

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Una scorciatoia è un comando a cui dai tu il nome e che avvia Claude Code con un determinato account. Come `ccam work`, passa ogni argomento a Claude Code invariato, quindi `ccwork --resume` e `ccwork -p "hi"` funzionano come ti aspetti. È un link a ccam in `~/.local/bin`.

Le scorciatoie sono facoltative. La procedura guidata di `ccam add` chiede se crearne una, `ccam add <name> --shortcut <command>` ne crea una senza chiedere e `ccam list` mostra la scorciatoia di ogni account, oppure `-` se non ne ha.

I nomi usano i caratteri a-z, 0-9, `.`, `-` e `_`. Un nome non può essere un comando di ccam, la scorciatoia di un altro account o il nome di un altro programma nel tuo `PATH`.

`ccam rename` mantiene il nome della scorciatoia, e `ccam remove` elimina la scorciatoia insieme all'account. `ccam shims install` ricrea i link di tutti gli account, per esempio dopo che hai spostato ccam.

## Condivisione

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Ogni elemento accetta `on` (il gruppo `default`), `off` oppure il nome di un gruppo. Vedi [condivisione](/docs/sharing).

## Copiare una configurazione esistente

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` copia `settings.json`, `CLAUDE.md`, le estensioni e i server MCP utente da una cartella di Claude Code in un gruppo, così ogni account che condivide quegli elementi li riceve. La cartella viene solo letta. Le opzioni `--settings`, `--memory`, `--extensions`, `--mcp` e `--history` copiano solo gli elementi indicati; senza di esse viene copiato tutto tranne la cronologia delle conversazioni. `-g` o `--group` sceglie il gruppo, e `-y` o `--yes` copia senza chiedere e mantiene la copia del gruppo quando un file è diverso. Vedi [copiare una configurazione esistente](/docs/sharing#copy-an-existing-setup).

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
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## Manutenzione

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install   # recreate the shortcut links
```

## File

| Percorso | Contenuto |
| --- | --- |
| `~/.config/ccam/config.toml` | Gli account e cosa condividono |
| `~/.local/share/ccam/profiles/<name>/` | La `CLAUDE_CONFIG_DIR` di ogni account |
| `~/.local/share/ccam/groups/<group>/` | Impostazioni, server MCP, estensioni e cronologia condivisi |
| `~/.local/state/ccam/` | Cache dell'utilizzo e orari dell'ultimo uso |
| `~/.local/bin/<shortcut>` | Scorciatoie come `ccwork`, link a ccam |

Le variabili `XDG_CONFIG_HOME`, `XDG_DATA_HOME` e `XDG_STATE_HOME` spostano queste cartelle.
