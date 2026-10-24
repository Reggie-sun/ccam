---
title: Usare Remote Control con ccam
description: "Avvia il Remote Control di Claude Code con qualsiasi account usando ccam <name> rc: opzioni valide, perché ccam propone di copiare i server MCP e soluzioni."
updated: 2026-10-23
order: 4
translation: machine
sourceHash: "3e17d89f98c8"
---

`ccam work rc` avvia il Remote Control di Claude Code con l'account `work`, così puoi guidare le sessioni su questo computer da claude.ai/code o dall'app mobile di Claude. Funziona come `claude rc` con una sola differenza: se l'account condivide server MCP, ccam chiede una volta se copiarli nel `.claude.json` dell'account, perché le sessioni di Remote Control non leggono i server da nessun'altra parte.

## Avviarlo

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

Una [scorciatoia](/features/shortcuts) funziona allo stesso modo: `ccwork rc` equivale a `ccam work rc`. Esegui il comando nella cartella in cui vuoi lavorare, poi apri claude.ai/code o l'app di Claude, con l'accesso fatto con lo stesso account Claude di `work`.

## Opzioni

Qui si incontrano due tipi di opzioni: quelle di ccam, che vanno prima del nome dell'account o stanno nelle impostazioni di ccam, e quelle di Remote Control, che vanno dopo `rc` e arrivano a Claude Code invariate.

### Opzioni e impostazioni di ccam

| Opzione o impostazione | Cosa fa | Predefinito |
| --- | --- | --- |
| `--dry-run` (prima dell'account) | Stampa l'ambiente e il comando, più l'avviso MCP se c'è. Non chiede e non copia nulla | disattivata |
| `Copy them into .claude.json (a backup is kept; servers of the account's own stay)? [y/N]` | Viene chiesto quando l'account condivide server MCP e al suo `.claude.json` manca l'elenco attuale del gruppo, o ne è sparito un server copiato. `y` li copia, `N` o Invio avvia senza di loro | `N` |
| `mcp_strategy = "sync"` in cima a `~/.config/ccam/config.toml` | Copia i server del gruppo in `.claude.json` prima di ogni avvio, che sia Remote Control o no, così non viene mai chiesto nulla | `"flag"` |
| `ccam group leave <group> work` | Toglie l'account dal gruppo che condivide con esso i server MCP. I server del gruppo vengono copiati nel suo `.claude.json` come server propri, così Remote Control li ha e non chiede nulla; con `--fresh` si avvia solo con i server propri dell'account | in `default` |

### Opzioni di Remote Control

Tutto ciò che segue `rc` appartiene a Claude Code. Queste sono le opzioni di `claude remote-control` in Claude Code 2.1.304; `ccam work rc --help` stampa l'elenco per la tua versione.

| Opzione | Cosa fa |
| --- | --- |
| `--name <name>` | Nome della sessione mostrato in claude.ai/code |
| `--remote-control-session-name-prefix <prefix>` | Prefisso per i nomi di sessione generati (predefinito: il nome host) |
| `-c`, `--continue` | Si ricollega alla sessione che Remote Control ha registrato per ultima per questa cartella, entro circa 4 ore |
| `--session-id <id>` | Si ricollega a una sessione specifica |
| `--permission-mode <mode>` | Modalità dei permessi per le nuove sessioni: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` o `plan` |
| `--spawn <mode>` | `same-dir` (predefinito), `worktree` (un git worktree per sessione) o `session` (una sola sessione, esce quando termina) |
| `--capacity <N>` | Numero massimo di sessioni contemporanee in modalità `same-dir` e `worktree` (predefinito 32) |
| `--[no-]create-session-in-dir` | Crea una prima sessione nella cartella corrente all'avvio (attiva per impostazione predefinita) |
| `--[no-]chrome` | Claude in Chrome per le nuove sessioni (disattivata per impostazione predefinita) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | Output di debug e dettagliato |

ccam legge il comando di Claude Code come fa Claude Code, saltando le opzioni insieme ai loro valori, quindi `ccam work rc --name laptop` e `ccam work --model opus rc` avviano entrambi Remote Control. Un'opzione che ccam non conosce conta come opzione senza valore; per sicurezza, metti `rc` subito dopo il nome dell'account.

## Cosa fa ccam

L'avvio segue il percorso consueto ([Avviare Claude Code](/features/launch)) con una sola modifica nel passaggio MCP:

1. ccam trova l'account, controlla i suoi link di condivisione e risolve `claude`, come per qualsiasi avvio.
2. Vede che il comando è `rc` o `remote-control` e omette `--mcp-config`.
3. Se l'account condivide server MCP, confronta il `mcp.json` del gruppo con il `sync-stamp` in `~/.local/state/ccam/mcp/work/`, e controlla che ogni server del gruppo sia ancora in `.claude.json`. Se `.claude.json` contiene già i server attuali, o il gruppo non ne ha, passa direttamente al punto 6.
4. Altrimenti elenca i server che mancherebbero alle sessioni:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; servers of the account's own stay)? [y/N]
   ```

5. Con `y` salva un backup di `.claude.json`, prende il lock di Claude Code, aggiunge i server del gruppo a `mcpServers` e mantiene gli altri server dell'account; un server proprio dell'account con il nome di un server del gruppo resta dell'account, e ccam lo segnala. Con `N`, senza un terminale o mentre l'account è in esecuzione, non cambia nulla e lo segnala:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. Registra la sessione, imposta `CLAUDE_CONFIG_DIR` sulla cartella dell'account e si sostituisce con `claude rc` tramite `exec`.

Da lì Remote Control resta in esecuzione come server. Ogni sessione che avvia per te è un processo figlio di quel processo, quindi gira come `work`: con il login, le impostazioni e la cronologia delle conversazioni di `work`, e conta sul limite di utilizzo di `work`. Finché è in esecuzione un server avviato da ccam, ccam considera l'account in esecuzione, e lo stesso vale per un server avviato con il semplice `claude remote-control` finché ha sessioni; un server del genere senza ancora alcuna sessione non viene visto.

Il dry run mostra la differenza. Un avvio normale passa i server del gruppo come argomento, Remote Control no:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude rc
```

## Perché funziona così

**Remote Control non accetta `--mcp-config`.** Negli avvii normali ccam passa i server del gruppo a Claude Code come `--mcp-config=<path>` e non scrive mai nei file di Claude Code. `claude remote-control` rifiuta quell'opzione prima del nome del suo comando e non ne accetta dopo, quindi prima di ccam 0.6.0 `ccam work rc` falliva per ogni account che condivideva server MCP. Ora ccam omette l'opzione.

**Le sue sessioni leggono i server solo da `.claude.json`.** Le sessioni avviate da Remote Control non ereditano le opzioni dalla riga di comando. L'unico posto in cui trovano i server MCP utente è il `.claude.json` dell'account. Quel file contiene anche i dati di login dell'account e lo stato per cartella, quindi ccam non può collegarlo a una copia condivisa come fa con `settings.json`. Copiarvi i server è l'unico modo per dare alle sessioni l'elenco condiviso.

**ccam chiede prima di scrivere.** Con la strategia flag predefinita, ccam lascia stare `.claude.json` all'avvio. Remote Control è l'unico avvio in cui deve scriverci, quindi chiede prima, conserva un backup e non tocca nient'altro che `mcpServers`. (Al di fuori degli avvii, anche lasciare o riprendere la condivisione MCP e `ccam mcp import` vi scrivono, come dicono i loro piani.)

**Non scrive mentre l'account è in esecuzione.** Claude Code riscrive `.claude.json` durante una sessione, e una modifica fatta dall'esterno andrebbe persa o sovrascritta. Per questo ccam salta la copia quando è in esecuzione una qualsiasi sessione dell'account, compreso un altro server di Remote Control e Claude Code avviato senza ccam.

**Chiede una volta per ogni modifica.** Dopo una copia, `sync-stamp` contiene l'ora di modifica del `mcp.json` del gruppo. ccam chiede di nuovo solo quando quel file cambia, per esempio dopo `ccam mcp add`, o quando un server copiato è sparito da `.claude.json`. Se rispondi `N` non c'è ancora una copia, quindi chiede a ogni avvio di Remote Control.

**`ccam mcp` tiene aggiornate le copie.** Una voce in `.claude.json` prevale su `--mcp-config` in caso di nomi uguali, quindi una copia vecchia nasconderebbe una definizione più recente anche negli avvii normali. Per questo `ccam mcp add`, `add-json`, `remove` e `import`, e `ccam migrate`, aggiornano le copie di ogni account che non è in esecuzione, e lo segnalano quando una copia è cambiata:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**Senza un terminale la risposta è no.** Remote Control viene spesso eseguito da uno script, da un servizio o con `nohup`, dove nessuno può rispondere. In quel caso ccam si avvia senza i server invece di restare in attesa.

## Problemi e soluzioni

### Nelle sessioni di Remote Control mancano i server MCP condivisi

Alla domanda è stato risposto `N`, non c'era un terminale oppure l'account era in esecuzione. Ferma Remote Control, chiudi le altre sessioni dell'account ed esegui di nuovo `ccam work rc` in un terminale; rispondi `y`. Per non ricevere mai la domanda, imposta `mcp_strategy = "sync"` in `~/.config/ccam/config.toml`. I server di una sessione non cambiano mentre è in esecuzione, quindi avvia una nuova sessione dopo la copia.

### "Account work is running, so its .claude.json is left alone"

Una sessione di `work` è ancora aperta: un normale `ccam work`, una scorciatoia o un altro server di Remote Control. Chiudila ed esegui di nuovo `ccam work rc`. Le sessioni eseguite da un altro account non contano.

### ccam chiede ogni volta

Non hai ancora copiato i server. Rispondi `y` una volta, oppure imposta `mcp_strategy = "sync"`. Se non vuoi mai i server in Remote Control, smetti di condividere MCP con l'account senza mantenerli: `ccam group leave default work --fresh`, oppure tieni gli MCP in un gruppo a parte di cui l'account non fa parte. Un semplice `ccam group leave` copia i server del gruppo nel `.claude.json` dell'account, così Remote Control li ha.

### `ccam work rc` fallisce subito con un errore su `--mcp-config`

Il tuo ccam è precedente alla 0.6.0 e passa `--mcp-config` a Remote Control. Esegui `ccam update`. Lo stesso errore con un ccam aggiornato significa che un'opzione che ccam non conosce, con un valore, precedeva `rc`, e ccam ha preso il valore per il comando; metti `rc` subito dopo il nome dell'account.

### Una modifica fatta con `ccam mcp` non arriva a Remote Control

`ccam mcp add` e `remove` aggiornano le copie solo negli account che non sono in esecuzione, e lo segnalano: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Ferma Remote Control e avvialo di nuovo; ccam vede il `mcp.json` modificato e chiede di nuovo.

### Viene usato un mio server invece di quello del gruppo

Se l'account ha un proprio server con lo stesso nome di un server del gruppo, la copia lo lascia al suo posto e lo segnala: `account work has its own MCP servers fs, so group default's of those names are not copied there`. Per usare quello del gruppo, rimuovi quello dell'account (`claude mcp remove -s user <name>`); per dare la tua versione a tutti gli account, spostala nel gruppo con `ccam mcp import work <server> --force`. I backup di `.claude.json` si trovano in `~/.local/state/ccam/mcp/work/backups/`, che conserva i cinque più recenti.

### Voglio rimuovere i server copiati

`ccam mcp remove <server>` rimuove un server dal gruppo e da ogni copia. Per tornare allo stato precedente alla copia, ripristina il backup più recente da `~/.local/state/ccam/mcp/work/backups/` mentre non è in esecuzione nessuna sessione dell'account.

### Remote Control chiede se considerare attendibile la cartella, oppure esce senza un terminale

L'attendibilità delle cartelle è salvata per account in `.claude.json`, quindi una cartella che hai considerato attendibile come `personal` è nuova per `work`. Esegui `ccam work` una volta in quella cartella e accetta, poi avvia lì Remote Control.

### "You must be logged in with a Claude account that has a subscription"

Remote Control richiede un login con un abbonamento Claude. Controlla l'email dell'account in `ccam list` e, se serve, accedi di nuovo con `ccam login work`.

### La sessione non compare sul telefono

L'app di Claude o claude.ai/code ha l'accesso fatto con un altro account Claude. Confrontalo con l'email che `ccam list` mostra per `work`.

### ccam avvisa che un altro Claude Code condivide la cronologia delle conversazioni in questa cartella

Un altro account dello stesso gruppo di cronologia è in esecuzione nella stessa cartella. Le sessioni che proseguono la stessa conversazione da entrambi gli account possono mescolare la cronologia. Rispondi `y` solo se lavoreranno su conversazioni diverse, oppure avvia Remote Control in un'altra cartella.

## Pagine correlate

- [Avviare Claude Code](/features/launch)
- [Server MCP condivisi](/features/mcp)
- [Condivisione tra account](/features/share)
- [Gestire gli account](/features/accounts)
