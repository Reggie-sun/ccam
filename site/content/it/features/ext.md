---
title: Ispezionare le estensioni
description: ccam ext mostra skill, comandi, plugin, server MCP e hook che ogni account carica, da dove arriva ciascuno e cosa non funziona.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

`ccam ext` risponde alla domanda "perché questa skill (o questo server MCP) funziona in un account e non nell'altro?". Elenca, per ogni account, le estensioni che Claude Code carica e da dove arriva ciascuna: un gruppo di condivisione, la cartella dell'account o il progetto in cui ti trovi. Sotto la tabella elenca cosa impedisce a un'estensione di caricarsi come ti aspetti. Si limita a leggere.

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
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

Una cella indica la provenienza di quell'elemento in quell'account: il gruppo di condivisione (`default`), `account` quando si trova nella cartella dell'account, oppure `-` quando l'account non ce l'ha. `!` segnala un elemento con un problema, spiegato sotto la tabella.

## Cosa legge

| Tipo | Dove | Provenienza mostrata |
| --- | --- | --- |
| Skill | ogni cartella in `skills/` con un `SKILL.md` | il gruppo a cui punta `skills/`, oppure `account` |
| Comandi, agenti | i file `.md` sotto `commands/` e `agents/`, le sottocartelle come `folder/name` | il gruppo, oppure `account` |
| Plugin | `plugins/installed_plugins.json` | il gruppo, oppure `account` |
| Server MCP | il `mcp.json` del gruppo (quando MCP è condiviso) e il `.claude.json` dell'account | il gruppo, oppure `account` |
| Hook | `hooks` in `settings.json`, contati per evento | il gruppo, oppure `account` |
| Questo progetto | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` e `.mcp.json` nella cartella corrente | elencati una volta, per ogni account |

## Problemi che segnala

- **Link interrotti**: un link in `skills/`, `commands/` o `agents/` la cui destinazione non esiste più, oppure una cartella condivisa (`skills`, `commands`, `agents`, `plugins`) che punta a una posizione mancante (riparala con `ccam doctor --fix`).
- **Una cartella di skill senza `SKILL.md`**: Claude Code la salta.
- **Un server MCP definito in modo diverso in due posti**: il `mcp.json` del gruppo e il `.claude.json` dell'account, oppure uno dei due e il `.mcp.json` del progetto. Le stesse impostazioni in entrambi i posti (quelle che scrive `mcp_strategy = "sync"`) non sono un problema. Sposta un server proprio dell'account nel gruppo con `ccam mcp import <account> <server>`.
- **Un `installed_plugins.json` che non è JSON valido.**

## Sicurezza e casi limite

- `ccam ext` non modifica nulla e non fa richieste di rete. Legge i nomi dei server e il numero di hook, mai i valori di `env` o delle intestazioni nelle definizioni MCP.
- Mostra ciò che c'è sul disco. Un plugin installato ma disattivato nelle impostazioni, o un server bloccato dalle impostazioni gestite, compare comunque.
- Eseguilo nella cartella del progetto per includere la cartella `.claude` e il `.mcp.json` di quel progetto.

## Pagine correlate

- [Condivisione tra account](/features/share)
- [Server MCP condivisi](/features/mcp)
- [Controllo e riparazione](/features/doctor)
