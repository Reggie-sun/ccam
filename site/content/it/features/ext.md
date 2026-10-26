---
title: Ispezionare le estensioni
description: ccam ext mostra skill, comandi, plugin, server MCP e hook che ogni account carica, da dove arriva ciascuno e cosa non funziona.
updated: 2026-10-26
order: 14
translation: machine
sourceHash: 7ca357a7f38a
---

`ccam ext` risponde alla domanda "Perché questa skill (o questo server MCP) funziona in un account ma non nell'altro?" Elenca, per ogni account, le estensioni che Claude Code carica e da dove arriva ciascuna: un gruppo di condivisione, la cartella dell'account o il progetto in cui ti trovi. Sotto la tabella elenca cosa impedisce a un'estensione di caricarsi come ti aspetti. Non modifica nulla.

## Utilizzo

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  default   default

MCP servers
       personal   work
  ctx  default !  default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app):
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
! personal: half: no SKILL.md, so Claude Code does not load it
```

Una cella indica la provenienza di quell'elemento in quell'account: il gruppo di condivisione (`default`), `account` quando si trova nella cartella dell'account, `local` per un server MCP che l'account ha solo per questa cartella, oppure `-` quando l'account non ce l'ha. Un plugin installato ma non attivato mostra `(off)` dopo la sua provenienza. `!` segnala un elemento con un problema, spiegato sotto la tabella; i problemi che non appartengono a nessun account, come i server che ha solo `claude` da solo, vengono per ultimi senza un nome. `ccam ext --json` stampa lo stesso rapporto (`accounts`, il progetto, `problems`), con `remote` sui server MCP che Claude Code raggiunge in rete; senza account stampa il rapporto con `accounts: []`.

## Cosa legge

| Tipo | Dove | Provenienza mostrata |
| --- | --- | --- |
| Skill | ogni cartella in `skills/` con un `SKILL.md` | il gruppo a cui punta `skills/`, oppure `account` |
| Comandi, agenti | i file `.md` sotto `commands/` e `agents/`, le sottocartelle come `folder/name` | il gruppo, oppure `account` |
| Plugin | `plugins/installed_plugins.json`, attivati o disattivati da `enabledPlugins` nelle impostazioni dell'account e del progetto; un plugin installato per un solo progetto compare solo in quel progetto | il gruppo, oppure `account`; `(off)` quando non è attivato |
| Server MCP | il `mcp.json` del gruppo (quando MCP è condiviso), il `.claude.json` dell'account (i suoi server, e quelli che ha solo per questa cartella) e i server dei suoi plugin, come `plugin:<plugin>:<server>` | il gruppo, `account` o `local` |
| Hook | `hooks` in `settings.json`, contati per evento | il gruppo, oppure `account` |
| Questo progetto | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` e `.mcp.json` nella cartella corrente | elencati una volta; un account che non ha dato fiducia alla cartella riceve un problema per i server del suo `.mcp.json` |

## Problemi che segnala

- **Link interrotti**: un link in `skills/`, `commands/` o `agents/` la cui destinazione non esiste più, oppure una cartella condivisa (`skills`, `commands`, `agents`, `plugins`) che punta a una posizione mancante (riparala con `ccam doctor --fix`, che ricrea una cartella del gruppo mancante).
- **Una cartella di skill senza `SKILL.md`**: Claude Code la salta.
- **Un server MCP definito in modo diverso in due posti**: il `mcp.json` del gruppo e il `.claude.json` dell'account, oppure uno dei due e il `.mcp.json` del progetto. Definizioni identiche in entrambi i posti (quelle che scrive `mcp_strategy = "sync"`) non sono un problema. Sposta un server proprio dell'account nel gruppo con `ccam mcp import <account> <server>`.
- **Un `installed_plugins.json` che non è JSON valido**, un plugin che `enabledPlugins` non nomina e un plugin la cui cartella manca.
- **Server che carica un solo account**: un server presente solo nel `.claude.json` di questo account, che gli altri account del suo gruppo non caricano, e uno nell'ambito locale dell'account per questa cartella (ciò che fa `claude mcp add` se non gli si passa `-s user`). Ciascuno viene indicato con il comando `ccam mcp import` che lo condivide. I server nel `~/.claude.json` di `claude` da solo che nessun account carica vengono indicati con il comando `ccam migrate --mcp`.
- **Un server di progetto in attesa di approvazione**: l'account non ha ancora dato fiducia alla cartella, quindi Claude Code lo chiede una volta, e fino ad allora i server del `.mcp.json` del progetto restano in attesa. Il problema dice dove approvarli: `"enabledMcpjsonServers"` nel `settings.json` degli account del gruppo delle impostazioni, o di questo account. La fiducia data al worktree principale di un repository git vale anche negli altri suoi worktree.
- **Accessi e connettori**: un server MCP che richiede un accesso in questo account mentre un altro account vi ha già effettuato l'accesso (un accesso vale per un solo account; accedi con `/mcp`), e i connettori di claude.ai che un altro account ha usato e questo no.

## Sicurezza e casi limite

- `ccam ext` non modifica nulla e non fa richieste di rete. Legge i nomi dei server e il numero di hook, mai i valori di `env` o delle intestazioni nelle definizioni MCP. Per il controllo degli accessi legge i nomi degli accessi MCP nel login di ogni account (su macOS dal portachiavi), mai i token, e solo per gli account che caricano un server che carica anche un altro account.
- Mostra ciò che c'è sul disco. Un server bloccato dalle impostazioni gestite compare comunque; un plugin disattivato nelle impostazioni mostra `(off)`.
- Eseguilo nella cartella del progetto per includere la cartella `.claude` e il `.mcp.json` di quel progetto.

## Pagine correlate

- [Condivisione con i gruppi](/features/share)
- [Server MCP condivisi](/features/mcp)
- [Controllo e riparazione](/features/doctor)
