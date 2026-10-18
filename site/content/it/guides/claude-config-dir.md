---
title: Come funziona CLAUDE_CONFIG_DIR
description: Cosa fa CLAUDE_CONFIG_DIR in Claude Code, cosa contiene la cartella, dove macOS conserva il login e quali file si possono condividere tra account.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

`CLAUDE_CONFIG_DIR` indica a Claude Code quale cartella usare al posto di `~/.claude`. Tutto ciò che è specifico dell'account si trova lì: il login, le impostazioni, i server MCP, i plugin e la cronologia delle conversazioni. Se la fai puntare a cartelle diverse, ottieni account separati.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## Cosa c'è nella cartella

Verificato con Claude Code 2.1.302:

| Percorso | Contenuto | Condivisibile |
| --- | --- | --- |
| `.credentials.json` | Token di login (Linux), token OAuth dei server MCP | No |
| `.claude.json` | Dati dell'account, server MCP dell'utente, stato per progetto | No (i server MCP a parte) |
| `settings.json` | Permessi, hook, modello, riga di stato, plugin | Sì |
| `CLAUDE.md` | Istruzioni globali | Sì |
| `agents/`, `commands/`, `skills/`, `plugins/` | Estensioni | Sì |
| `projects/` | Conversazioni come file `.jsonl`, memoria per progetto | Sì |
| `history.jsonl` | Cronologia dei prompt | Sì |
| `sessions/`, `shell-snapshots/`, `cache/` e altri | Stato di runtime | No |

## Due dettagli che colgono di sorpresa

**Senza `CLAUDE_CONFIG_DIR`, `.claude.json` non si trova dentro `~/.claude`.** In quel caso Claude Code lo tiene in `~/.claude.json`, accanto alla cartella. Se inizi a usare `CLAUDE_CONFIG_DIR=~/.claude`, Claude Code cerca invece `~/.claude/.claude.json` e non trova il tuo account. `ccam import` gestisce questo caso collegando i due file.

**Su macOS il login si trova nel Portachiavi, non nella cartella.** Claude Code lo salva con il nome `Claude Code-credentials`, seguito da un breve hash del percorso della cartella quando `CLAUDE_CONFIG_DIR` è impostata. Cartelle separate hanno quindi voci separate nel Portachiavi. Se il Portachiavi rifiuta la scrittura, Claude Code ripiega su `.credentials.json` nella cartella.

## Condivisione tra cartelle

I link simbolici funzionano: Claude Code legge e scrive attraverso un `settings.json` o una cartella `projects/` collegati e mantiene il link. I server MCP sono l'eccezione, perché condividono `.claude.json` con il login; passali invece con `--mcp-config`.

ccam si basa esattamente su queste regole. Vedi [condivisione](/docs/sharing) per sapere come collega le cartelle e [la guida agli account multipli](/guides/multiple-claude-code-accounts) per il quadro generale.
