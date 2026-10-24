---
title: Domande frequenti
metaTitle: "Domande frequenti: più account Claude Code"
description: Risposte su ccam, account Claude Code, limiti di utilizzo, condivisione, privacy e sistemi supportati.
updated: 2026-10-24
order: 6
translation: machine
sourceHash: "c4e3b89bb4ea"
---

## ccam è uno strumento ufficiale di Anthropic?

No. ccam è un progetto open source non ufficiale e non è affiliato ad Anthropic né approvato da Anthropic. Usa l'impostazione `CLAUDE_CONFIG_DIR`, che [Claude Code documenta](https://code.claude.com/docs/en/env-vars) per eseguire più account in parallelo.

## Posso usare più account Claude Code?

Dipende dai termini di ciascun account e, per gli account di lavoro, dalle policy della tua organizzazione. Controllali prima di usare più account. ccam rende solo più semplice passare da un account all'altro tra quelli che hai già.

**Posti nei piani Team ed Enterprise.** La scelta automatica (`ccam next`, `--fallback`, `ccam pick`, il tasto `n` della dashboard) passa anche da un posto all'altro della stessa organizzazione Claude Team o Enterprise. Il centro assistenza di Anthropic dice: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (ogni posto di un piano Team è assegnato nominativamente a una sola persona, e solo quella persona può usarlo) e "Users can't route their work through another user's account to use that user's remaining usage" (gli utenti non possono far passare il proprio lavoro attraverso l'account di un altro utente per sfruttare l'utilizzo residuo di quell'utente) ([piani Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); lo stesso vale per i [piani Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Usare la scelta automatica tra posti diversi può quindi violare i termini di Anthropic. Verifica con la tua organizzazione e con Anthropic prima di farci affidamento. I piani personali (Pro, Max) sono trattati nel paragrafo qui sopra.

## ccam aggira i limiti di utilizzo?

No. Ogni account mantiene i propri limiti. ccam te li mostra e ti permette di continuare su un altro account di tua proprietà quando uno li esaurisce.

## Come si usano due account Claude Code su un solo computer?

Assegna a ogni account la propria cartella di configurazione con `CLAUDE_CONFIG_DIR`. ccam lo fa per te: esegui `ccam add` per ogni account, poi avviali con `ccam personal` e `ccam work`, oppure con scorciatoie a cui dai tu il nome, come `ccwork`. La [guida](/guides/multiple-claude-code-accounts) mostra anche come farlo a mano.

## Come continuo una conversazione su un altro account?

Esegui `ccam next --continue` nella cartella: porta l'ultima conversazione della cartella nell'account successivo (`--auto` per quello con più margine). Dentro Claude Code, `/ccam:switch` del [plugin di ccam](/docs/install#claude-code-plugin) fa lo stesso per la sessione in cui ti trovi. Quando i due account condividono la cronologia delle conversazioni (`ccam group set default --add history`, o un gruppo tutto loro: `ccam group add pair --items history --members work,side`), il nuovo account apre proprio quella conversazione; altrimenti Claude Code la copia nella cronologia del nuovo account.

## ccam vede o invia il mio login?

ccam non condivide né copia mai il tuo login di Claude tra gli account, e non copia mai gli accessi dei server MCP. Legge il token di accesso di un account solo per chiedere all'endpoint di utilizzo di Anthropic i dati di utilizzo di quell'account, quando esegui `ccam usage`, `ccam next --auto` o `ccam doctor`, e rinnova un token scaduto solo con `ccam usage --refresh`; con `usage_source = "claude"` non legge affatto il token. `ccam ext` e `ccam doctor` leggono i nomi degli accessi MCP di un account, mai i loro token. ccam non ha telemetria. A parte queste consultazioni, l'unica richiesta che fa di sua iniziativa è la domanda della dashboard a GitHub sull'ultima release di ccam, che non contiene nulla su di te o sui tuoi account.

## Come accedo su un server remoto?

Esegui `ccam add` o `ccam login <name>` nella sessione SSH come al solito. Claude Code attende l'accesso sul server, quindi il browser sul tuo computer finisce su una pagina `http://localhost:<port>/callback?code=…` che non si carica. Copia quell'indirizzo dal browser e incollalo nel terminale: ccam lo passa a Claude Code sul server, che completa l'accesso. L'indirizzo contiene un codice monouso, quindi non condividerlo.

## Dove vengono salvati i miei account?

In `~/.local/share/ccam/profiles/<name>/`, una cartella per account. Su macOS, Claude Code conserva il login vero e proprio nel Portachiavi, con un nome derivato da quella cartella, quindi `ccam rename` mantiene lì il vecchio nome della cartella.

## Cosa succede alla mia cartella ~/.claude esistente?

Niente, a meno che tu non esegua `ccam import` o `ccam migrate`. L'importazione registra la cartella dove si trova, quindi `claude` da solo e `ccam <name>` continuano a usare lo stesso account e le stesse impostazioni. `ccam migrate` copia le sue impostazioni, `CLAUDE.md`, le estensioni e i server MCP nelle impostazioni condivise, e si limita a leggere la cartella.

## Quali sistemi supporta ccam?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine e altri), macOS e Windows 10 e 11, in modo nativo o all'interno di WSL. Su Windows, condividere file come `settings.json` richiede la Modalità sviluppatore; le cartelle condivise funzionano anche senza, e senza di essa il primo `ccam add` crea il gruppo `default` con i soli server MCP ed estensioni ([condivisione su Windows](/features/share#on-windows)).

## Perché un server MCP compare in un solo account?

I server aggiunti all'interno di una sessione con `claude mcp add` vengono salvati solo per quell'account e, senza `-s user`, solo per la cartella corrente. ccam te lo segnala la volta successiva che avvii l'account, e `ccam mcp list` e `ccam ext` li mostrano; sposta il server nel gruppo condiviso con `ccam mcp import <account> <server>`.

## ccam rallenta Claude Code?

Avviare Claude Code tramite ccam aggiunge circa 4 millisecondi, misurati con `scripts/perf-check.sh` su un portatile. ccam legge un piccolo file, controlla i link condivisi e poi si sostituisce con `claude`, quindi non è in esecuzione mentre lavori. Windows non può sostituire un processo, quindi lì ccam avvia Claude Code e ne attende la fine, senza fare altro.

## Come si aggiorna ccam?

Esegui `ccam update`. Scarica l'ultima release, ne verifica il checksum e sostituisce il file `ccam`; `ccam update --check` ti dice solo se esiste una release più recente. Anche la dashboard ti avvisa quando esiste una release più recente, e `U` la installa dopo aver chiesto conferma; ccam non si aggiorna mai da solo. Se hai installato il pacchetto `.deb`, installa invece il nuovo `.deb` con apt. Vedi [aggiornamento](/docs/install#update).

## Come si disinstalla ccam?

Rimuovi la scorciatoia di ogni account con `ccam shortcut <account> --remove`. Se hai importato `~/.claude` e fa parte di un gruppo, esegui prima `ccam remove <name>` per quell'account, così riceve una copia di ciò che condivide. Poi elimina `ccam` da `~/.local/bin` e rimuovi `~/.config/ccam`, `~/.local/share/ccam` e `~/.local/state/ccam` se non ti servono più gli account. Vedi [installazione](/docs/install#uninstall).
