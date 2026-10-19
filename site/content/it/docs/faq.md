---
title: Domande frequenti
description: Risposte su ccam, account Claude Code, limiti di utilizzo, condivisione, privacy e sistemi supportati.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: e07d468d453b
---

## ccam è uno strumento ufficiale di Anthropic?

No. ccam è un progetto open source non ufficiale e non è affiliato ad Anthropic né approvato da Anthropic. Usa l'impostazione `CLAUDE_CONFIG_DIR`, che Claude Code documenta per eseguire più account in parallelo.

## Posso usare più account Claude Code?

Dipende dai termini di ciascun account e, per gli account di lavoro, dalle policy della tua organizzazione. Controllali prima di usare più account. ccam rende solo più semplice passare da un account all'altro tra quelli che hai già.

## ccam aggira i limiti di utilizzo?

No. Ogni account mantiene i propri limiti. ccam te li mostra e ti permette di continuare su un altro account di tua proprietà quando uno li esaurisce.

## Come si usano due account Claude Code su un solo computer?

Assegna a ogni account la propria cartella di configurazione con `CLAUDE_CONFIG_DIR`. ccam lo fa per te: esegui `ccam add` per ogni account, poi avviali con `ccam personal` e `ccam work`, oppure con scorciatoie a cui dai tu il nome, come `ccwork`. La [guida](/guides/multiple-claude-code-accounts) mostra anche come farlo a mano.

## Come continuo una conversazione su un altro account?

Condividi la cronologia delle conversazioni tra i due account con `ccam share <name> --history on`, poi esegui `ccam next --continue`. Claude Code apre la conversazione più recente della cartella corrente nell'account successivo.

## ccam vede o invia il mio login?

ccam non copia né condivide mai un login. Legge il token di accesso di un account solo per chiedere all'endpoint di utilizzo di Anthropic i dati di utilizzo di quell'account, quando esegui `ccam usage`, `ccam next --auto` o `ccam doctor`. ccam non ha telemetria e non invia nulla da nessun'altra parte.

## Dove vengono salvati i miei account?

In `~/.local/share/ccam/profiles/<name>/`, una cartella per account. Su macOS, Claude Code conserva il login vero e proprio nel Portachiavi, con un nome derivato da quella cartella.

## Cosa succede alla mia cartella ~/.claude esistente?

Niente, a meno che tu non esegua `ccam import` o `ccam migrate`. L'importazione registra la cartella dove si trova, quindi `claude` da solo e `ccam <name>` continuano a usare lo stesso account e le stesse impostazioni. `ccam migrate` copia le sue impostazioni, `CLAUDE.md`, le estensioni e i server MCP nelle impostazioni condivise, e si limita a leggere la cartella.

## Quali sistemi supporta ccam?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine e altri) e macOS. Su Windows, usalo all'interno di WSL.

## Perché un server MCP compare in un solo account?

I server aggiunti all'interno di una sessione con `claude mcp add` vengono salvati solo per quell'account. ccam te lo segnala la volta successiva che avvii l'account; sposta il server nel gruppo condiviso con `ccam mcp import <account> <server>`.

## ccam rallenta Claude Code?

Avviare Claude Code tramite ccam aggiunge circa 4 millisecondi, misurati con `scripts/perf-check.sh` su un portatile. ccam legge un piccolo file, controlla i link condivisi e poi si sostituisce con `claude`, quindi non è in esecuzione mentre lavori.

## Come si disinstalla ccam?

Rimuovi la scorciatoia di ogni account con `ccam shortcut <account> --remove`, elimina `ccam` da `~/.local/bin`, poi `~/.config/ccam`, `~/.local/share/ccam` e `~/.local/state/ccam` se non ti servono più gli account. Vedi [installazione](/docs/install#uninstall).
