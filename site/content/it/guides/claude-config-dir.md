---
title: Come funziona CLAUDE_CONFIG_DIR
description: Cosa fa CLAUDE_CONFIG_DIR in Claude Code, cosa contiene la cartella, dove macOS conserva il login e quali file si possono condividere tra account.
updated: 2026-10-23
order: 2
translation: machine
sourceHash: "2cf08478085f"
---

`CLAUDE_CONFIG_DIR` indica a Claude Code quale cartella usare al posto di `~/.claude`. Tutto ciò che è specifico dell'account sta lì: il login, le impostazioni, i server MCP, i plugin e la cronologia delle conversazioni. Puntandola a cartelle diverse ottieni account separati.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## CLAUDE_CONFIG_DIR è supportata ufficialmente?

Sì. Claude Code la elenca tra le sue [variabili d'ambiente](https://code.claude.com/docs/en/env-vars) come il modo per cambiare la cartella di configurazione, "utile per eseguire più account in parallelo", e la sua [guida all'autenticazione](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) la usa per accedere con più account. Puoi impostarla nella shell o nelle impostazioni utente, ma non nelle impostazioni di un progetto.

## Cosa c'è nella cartella?

Verificato con Claude Code 2.1.309:

| Percorso | Contenuto | Condivisibile |
| --- | --- | --- |
| `.credentials.json` | Token di login (Linux, Windows), token OAuth MCP | No |
| `.claude.json` | Dati dell'account, server MCP dell'utente, stato per progetto | No (i server MCP a parte) |
| `settings.json` | Permessi, hook, modello, riga di stato, plugin | Sì |
| `CLAUDE.md` | Istruzioni globali | Sì |
| `agents/`, `commands/`, `skills/`, `plugins/` | Estensioni | Sì |
| `projects/` | Conversazioni come file `.jsonl`, memoria per progetto | Sì |
| `history.jsonl` | Cronologia dei prompt | Sì |
| `sessions/`, `shell-snapshots/`, `cache/` e altri | Stato di runtime | No |

## Cosa trae in inganno?

**Senza `CLAUDE_CONFIG_DIR`, `.claude.json` non si trova dentro `~/.claude`.** In quel caso Claude Code lo tiene in `~/.claude.json`, accanto alla cartella. Se inizi a usare `CLAUDE_CONFIG_DIR=~/.claude`, Claude Code cerca invece `~/.claude/.claude.json` e non trova il tuo account. `ccam import` gestisce questo caso collegando i due file. Se `~/.claude` ha già un proprio `.claude.json`, da un'esecuzione precedente con `CLAUDE_CONFIG_DIR=~/.claude`, ccam lascia entrambi i file e lo segnala, perché solo tu sai quale è quello giusto.

**Su macOS il login sta nel Portachiavi, non nella cartella.** Claude Code lo salva con il nome `Claude Code-credentials`, seguito da `-` e dalle prime 8 cifre esadecimali dello SHA-256 del percorso della cartella quando `CLAUDE_CONFIG_DIR` è impostata, per l'account indicato da `$USER`. Cartelle separate ottengono quindi voci separate nel Portachiavi, e una cartella spostata o rinominata perde il suo login; per questo su macOS `ccam rename` mantiene il nome della cartella. Se il Portachiavi rifiuta la scrittura, Claude Code ripiega su `.credentials.json` nella cartella, come descrive [la sua documentazione](https://code.claude.com/docs/en/authentication#credential-management).

**Non tutti i login stanno nella cartella.** Un login di Claude Console senza chiave API viene salvato fuori dalla cartella, quindi cartelle separate non ne tengono distinti due. I login di claude.ai e le chiavi API sono salvati per cartella.

## Dove si trova la cartella su Windows?

Per impostazione predefinita in `%USERPROFILE%\.claude`, con il login in `.credentials.json` al suo interno, oppure in Gestione credenziali di Windows quando Claude Code è impostato per usarla, con un nome ricavato dal percorso della cartella come su macOS. Imposta `CLAUDE_CONFIG_DIR` su un'altra cartella per ogni account in più, come su Linux e macOS.

## Le cartelle possono condividere file?

I link simbolici funzionano: Claude Code legge e scrive attraverso un `settings.json` o una cartella `projects/` collegati e mantiene il link. I server MCP sono l'eccezione, perché condividono `.claude.json` con il login; passali invece con `--mcp-config`.

ccam si basa proprio su queste regole. Vedi [condivisione](/docs/sharing) per come collega le cartelle, e [la guida a più account](/guides/multiple-claude-code-accounts) per il quadro generale.
