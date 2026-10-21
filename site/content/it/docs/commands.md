---
title: Riferimento dei comandi
metaTitle: Riferimento dei comandi di ccam
description: Tutti i comandi e le opzioni di ccam, da ccam add e ccam <name> fino a ccam usage, ccam next e ccam doctor.
updated: 2026-10-21
order: 3
translation: machine
sourceHash: "338ac11a0b5a"
---

Metti `--dry-run` prima di un comando per vedere cosa farebbe senza eseguirlo, per esempio `ccam --dry-run work`. Funziona quando avvii Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`) e con `group`, `migrate`, `shortcut`, `update` e `doctor --fix`. Gli altri comandi che modificano qualcosa lo rifiutano invece di essere eseguiti.

<!-- commands -->

## Avviare Claude Code

`ccam <name|number> [arguments]` avvia Claude Code con quell'account. ccam non legge gli argomenti che seguono il nome; li passa a `claude` nello stesso ordine, e standard input, codice di uscita e segnali si comportano esattamente come con `claude` stesso. Un `--` dopo il nome è facoltativo.

`ccam` con argomenti ma senza nome di account, come `ccam -p "summarise"`, usa l'account impostato per la cartella corrente, altrimenti l'account predefinito. Impostali con `ccam folder <name> [folder]` e `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` prima del nome dell'account serve per le esecuzioni `-p`/`--print`. Quando l'account raggiunge il suo limite di utilizzo, ccam riesegue lo stesso comando con l'account connesso con più margine e stampa solo l'output dell'esecuzione portata a termine. Vedi [Account per cartella](/features/folder) e [Rieseguire al limite](/features/next#run-again-at-the-limit).

Una scorciatoia come `ccwork` avvia il suo account allo stesso modo: `ccwork --resume` fa ciò che fa `ccam work --resume`. Vedi [scorciatoie](/docs/commands#shortcuts).

## Aggiungere e rimuovere account

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` senza una cartella cerca `~/.claude` e `~/.claude-*` e chiede conferma per ciascuna. Le cartelle importate vengono registrate dove si trovano; ccam non sposta né copia mai un login.

`ccam add --api` aggiunge un account che usa una chiave API invece di un login di Claude; ccam legge la chiave a ogni avvio e non la salva mai. Vedi [Account con chiave API](/features/api).

`ccam rename`, `ccam remove` e i comandi `ccam group` si rifiutano di modificare un account mentre Claude Code è in esecuzione con quell'account. Chiudi prima la sessione.

`ccam list --json` stampa ogni account con il suo slot, la scorciatoia, l'email e la cartella. Il suo array `groups` elenca i gruppi dell'account, e il suo oggetto `share` associa ogni elemento (`mcp`, `settings`, `memory`, `extensions`, `history`) al gruppo da cui proviene, oppure a `local`.

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
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Un gruppo condivide degli elementi tra i suoi account; un account può stare in più gruppi, ma ogni elemento proviene da un solo gruppo. Vedi [condivisione con i gruppi](/features/share).

## Copiare una configurazione esistente

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` copia `settings.json`, `CLAUDE.md`, le estensioni e i server MCP utente da una cartella di Claude Code in un gruppo, così ogni account che condivide quegli elementi li riceve. La cartella viene solo letta. Le opzioni `--settings`, `--memory`, `--extensions`, `--mcp` e `--history` copiano solo gli elementi indicati. Senza di esse, in un terminale una checklist ti permette di scegliere, con tutto selezionato tranne la cronologia delle conversazioni; altrove viene copiato tutto tranne la cronologia delle conversazioni. `-g` o `--group` sceglie il gruppo, e `-y` o `--yes` copia senza chiedere e mantiene la copia del gruppo quando un file è diverso. Copia anche gli accessi salvati dei server nel file di login di ciascun account. Vedi [copiare una configurazione esistente](/docs/sharing#copy-an-existing-setup).

## Server MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` sposta nel suo gruppo i server che ha solo quell'account. Un server il cui nome è già usato dal gruppo per una definizione diversa viene saltato; aggiungi `--force` per sostituire la copia del gruppo, di cui ccam fa prima un backup.

## Utilizzo e cambio di account

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam never touches a token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it
```

`next --auto`, `--fallback` e `ccam pick` scelgono allo stesso modo: prima gli account sotto `auto_max_usage` (`config.toml`, predefinito 100), secondo `auto_strategy` (`headroom`, `expiring` o `balance`), contando la finestra settimanale di `--model`, mai un account al 100%, e dentro l'organizzazione dell'account da cui partono (`--cross-org` o `cross_org = true` per cercare oltre). Quando non possono scegliere nessun account escono con codice 75, oppure con `--wait` attendono il reset del primo account. `ccam next --continue` porta una conversazione nell'altro account anche tra account che non condividono la cronologia: in quel caso Claude Code la copia nel nuovo account. A differenza di `ccam <name>`, `ccam next` legge le proprie opzioni, quindi le opzioni di Claude Code vanno dopo `--`; senza, ccam si ferma con `unknown flag`. Dentro Claude Code, il comando `/ccam:switch [account]` del plugin passa una conversazione a un altro account, e `ccam hook` è ciò che eseguono i suoi hook. Vedi [Passare all'account successivo](/features/next), [Scegliere un account per gli script](/features/pick) e [Sessioni in background di ogni account](/features/agents).

## Manutenzione

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

`ccam update` verifica il checksum SHA-256 della release prima di sostituire ccam, e `--version v0.2.0` installa una release specifica. Un ccam installato dal pacchetto `.deb` si aggiorna invece con apt. Vedi [aggiornamento](/docs/install#update).

## File

| Percorso | Contenuto |
| --- | --- |
| `~/.config/ccam/config.toml` | Account, gruppi e impostazioni (`theme`, `usage_source`, `relay`, `cross_org` …) |
| `~/.local/share/ccam/profiles/<name>/` | La `CLAUDE_CONFIG_DIR` di ogni account |
| `~/.local/share/ccam/groups/<group>/` | Impostazioni, server MCP, estensioni e cronologia condivisi |
| `~/.local/state/ccam/` | Cache dell'utilizzo, orari dell'ultimo uso e le note lasciate da `/ccam:switch` e dall'hook del limite (`handoff/`) |
| `~/.local/bin/<shortcut>` | Scorciatoie come `ccwork`, link a ccam |

Le variabili `XDG_CONFIG_HOME`, `XDG_DATA_HOME` e `XDG_STATE_HOME` spostano queste cartelle. Su Windows, `~` è `%USERPROFILE%` e le scorciatoie sono `<shortcut>.exe`.
