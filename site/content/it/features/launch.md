---
title: Avviare Claude Code
description: Avvia Claude Code con qualsiasi account usando ccam <name|number>, ccam run o il selettore di ccam; gli argomenti arrivano a claude invariati, e --dry-run mostra il comando.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 6c04bb48df22
---

`ccam work` avvia Claude Code con l'account `work`, e tutto ciò che segue il nome arriva a `claude` così com'è. Usalo ovunque scriveresti `claude`: in un terminale, in uno script o con un input passato tramite pipe.

## Utilizzo

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the default account
ccam use work                    # make work the default account
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
```

Indica il nome di un account o il suo numero in `ccam list`. ccam non legge gli argomenti che seguono il nome, quindi `ccam work --help` mostra l'aiuto di Claude Code. Un `--` dopo il nome è facoltativo. Standard input, codice di uscita e segnali si comportano come con `claude`.

Quando il primo argomento inizia con `-` e non è una delle opzioni di ccam, ccam avvia l'account predefinito con tutti gli argomenti. L'account predefinito è quello impostato con `ccam use`, oppure, se non ne hai impostato uno, l'account con il numero più basso.

### Il selettore

`ccam` senza argomenti apre il selettore quando standard input e standard output sono un terminale, altrimenti stampa `ccam list`. Mostra per ogni account il numero, l'email, l'utilizzo su 5 ore, l'ultimo uso e gli elementi condivisi; l'utilizzo compare man mano che ogni account risponde.

| Tasto | Cosa fa |
| --- | --- |
| `Enter` | Avvia Claude Code con l'account selezionato |
| Da `1` a `9` | Avvia l'account con quel numero |
| `Up`/`Down`, `k`/`j` | Sposta la selezione |
| `a` | Chiude il selettore ed esegue la procedura guidata di `ccam add` |
| `s` | Modifica gli elementi condivisi dell'account selezionato (`Space` o `x` attiva o disattiva, `Enter` applica, `Esc` torna indietro) |
| `d` | Rende predefinito l'account selezionato |
| `r` | Recupera di nuovo l'utilizzo, senza la cache breve |
| `q`, `Esc`, `Ctrl+C` | Esce |

Il selettore si chiude prima di agire, così il terminale viene ripristinato prima che Claude Code si avvii. Dopo `d` o una modifica della condivisione si riapre. Una modifica della condivisione mostra il piano e chiede prima conferma, come fa `ccam share`.

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `--dry-run` | Stampa l'ambiente e il comando che ccam eseguirebbe, poi esce | disattivata |
| `-v`, `--version` | Stampa `ccam <version>` (solo come primo argomento) | |
| `-h`, `--help` | L'aiuto di ccam, come `ccam --help` o `ccam run --help` | |

Regole per gli argomenti:

- ccam legge le proprie opzioni solo prima del nome dell'account. Dopo il nome, ogni argomento appartiene a Claude Code, compresi `--dry-run` e `--help`.
- `--dry-run` funziona come `ccam --dry-run work`, `ccam run --dry-run work` e `ccam --dry-run` (l'account predefinito). Viene rispettata, non rifiutata.
- `ccam run <name|number> [claude args]` è la forma estesa di `ccam <name>`. `ccam run` senza un nome usa l'account predefinito.
- I nomi degli account non possono essere nomi di comandi di ccam come `add`, `list`, `run` o `usage`, così `ccam <word>` non è mai ambiguo.

Il dry run stampa una riga che potresti incollare in una shell:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

## Come funziona

Ogni avvio segue gli stessi passaggi, dalla riga di comando fino a `exec`:

<!-- motion:launch -->

Ogni riquadro si illumina mentre viene eseguito il suo passaggio, e le righe sottostanti sono l'ambiente e gli argomenti con cui parte Claude Code.

1. `cmd/ccam/main.go` chiama `app.Main`. Se ccam è stato avviato con un altro nome, si tratta di una [scorciatoia](/features/shortcuts) e si passa direttamente all'avvio.
2. `app.parse` esamina `os.Args` prima che entri in gioco cobra. Solo `ccam` da solo (il selettore), un nome di comando di ccam (`config.Reserved`), `-h`/`--help` e `-v`/`--version` passano a `cli.Execute`. Tutto il resto è un avvio.
3. `launcher.Launch` legge `~/.config/ccam/config.toml`, l'unico file che analizza.
4. `launcher.Prepare` trova l'account per nome o per numero, scarta un eventuale `--` iniziale e individua l'eseguibile `claude` con `ResolveClaude`. La ricerca nel PATH viene salvata in cache in `~/.local/state/ccam/claude-path` insieme all'ora di modifica del file, così ccam non cerca nel `PATH` a ogni avvio. Un file modificato o mancante provoca una nuova ricerca. Un percorso `claude_bin` in `config.toml` viene usato così com'è. ccam rifiuta un `claude` che è ccam stesso, confrontando l'inode, perché si creerebbe un ciclo.
5. `checkLinks` confronta la cartella dell'account `~/.local/share/ccam/profiles/<name>/` con le sue impostazioni di condivisione usando `share.Check`, che chiama solo `lstat` e `readlink`, una o due chiamate di sistema per elemento. Un link mancante viene creato sul momento. Qualsiasi cosa che potrebbe far perdere dati interrompe l'avvio.
6. `applyMCP` aggiunge i server MCP condivisi (vedi sotto).
7. `session.Register` registra la sessione. ccam scrive un marcatore con il nome del suo PID, che contiene l'ora di avvio del processo, in `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` (oppure in `~/.local/state/ccam/locks/` senza `XDG_RUNTIME_DIR`). Per un account che condivide la cronologia delle conversazioni, `historyCheck` cerca prima una sessione in esecuzione in `history-<group>-<hash>`, dove l'hash è quello della cartella corrente, e si registra anche lì. `exec` mantiene il PID, quindi il marcatore indica poi il processo di Claude Code stesso. Un marcatore e non un lock mantenuto, perché i processi avviati da Claude Code erediterebbero il lock e lo conserverebbero dopo la fine della sessione.
8. `Plan.Exec` imposta il titolo del terminale su `claude · <name>` (quando lo standard error è un terminale), scrive l'ora di avvio in `~/.local/state/ccam/last-used/<name>` senza `fsync` e chiama `syscall.Exec`.

Claude Code viene eseguito con `CLAUDE_CONFIG_DIR` impostata sulla cartella dell'account e `CCAM_PROFILE` impostata sul nome dell'account. ccam rimuove gli eventuali valori ereditati di queste variabili. Per una cartella `~/.claude` registrata con `ccam import` su macOS, imposta anche `CLAUDE_SECURESTORAGE_CONFIG_DIR` a un valore vuoto, così l'account e `claude` da solo usano lo stesso login nel Portachiavi.

### Perché exec

`syscall.Exec` sostituisce il processo di ccam con `claude`, quindi ccam non resta in memoria, tasti premuti, segnali e codice di uscita sono quelli di Claude Code, e gli argomenti restano un array, mai una stringa di shell. Il lavoro da fare dopo una sessione, come accorgersi dei server aggiunti con `claude mcp add`, avviene all'avvio successivo: ccam confronta l'ora di modifica di `.claude.json` con `~/.local/state/ccam/mcp/<name>/seen` e analizza il file solo se è cambiato.

### Come viene instradata una riga di comando

Tre righe di comando, tre percorsi attraverso `internal/app`:

<!-- motion:launch-route -->

Una scorciatoia viene riconosciuta da `argv[0]`, `ccam <name>` dal primo argomento, ed entrambe terminano in `syscall.Exec`. Solo un comando di ccam come `ccam usage` arriva a cobra.

### Server MCP condivisi

Con la strategia predefinita `flag`, ccam passa i server MCP del gruppo a Claude Code come argomento, posto prima dei tuoi:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` restituisce un argomento solo se il `mcp.json` del gruppo esiste, e `Plan.InsertArgs` lo mette subito dopo `claude`. Le tue opzioni `--mcp-config` o `--strict-mcp-config` continuano quindi a valere in aggiunta. È un solo argomento, `--mcp-config=<path>`, perché l'opzione `--mcp-config` di Claude Code accetta più valori: scritta come due argomenti, leggerebbe come file di configurazione anche un prompt che segue.

Con `mcp_strategy = "sync"` in `config.toml`, ccam scrive invece i server nel `.claude.json` dell'account, e salta questo passaggio se l'account è già in esecuzione.

### Remote Control

`ccam work rc` (o `ccam work remote-control`) avvia il Remote Control di Claude Code come `work`. `claude remote-control` rifiuta `--mcp-config` prima del suo verbo e non ne accetta dopo, e le sue sessioni leggono i server MCP solo dal `.claude.json` dell'account. Per questo qui la strategia flag funziona in modo diverso:

<!-- motion:launch-rc -->

1. `remoteControl` prende come verbo di Claude Code il primo argomento dopo l'account che non inizia con `-`. Se è `remote-control` o `rc`, `mcp.FlagArgs` non aggiunge nulla.
2. `mcp.Synced` confronta l'ora di modifica del `mcp.json` del gruppo con il `sync-stamp` dell'account. Se `.claude.json` contiene già i server attuali del gruppo, o il gruppo non ne ha, Remote Control si avvia senza domande.
3. Altrimenti ccam elenca i server che mancherebbero alle sessioni e chiede:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. Con `y`, `mcp.Sync` li copia come fa la strategia `sync`: con il lock di Claude Code, dopo un backup, mantenendo i server propri dell'account. Con `N`, o senza un terminale, Remote Control si avvia senza di loro, e ccam suggerisce `mcp_strategy = "sync"`, che li copia prima di ogni avvio.

ccam chiede a ogni avvio di Remote Control finché i server non vengono copiati e, dopo una copia, solo quando il `mcp.json` del gruppo cambia di nuovo. Una copia in `.claude.json` prevale su `--mcp-config` nei successivi avvii normali, quindi `ccam mcp add`, `remove` e `import` aggiornano le copie che hanno gli altri account ([Server MCP condivisi](/features/mcp)).

### Il limite di 10 ms

ccam aggiunge al massimo 10 ms (p95) prima che Claude Code si avvii. Il percorso di avvio non costruisce mai l'albero dei comandi di cobra, non avvia Bubble Tea e non apre connessioni di rete, e `TestPackageInitIsCheap` fallisce se le funzioni `init()` dei pacchetti richiedono in totale più di 5 ms. `scripts/perf-check.sh` misura `ccam --dry-run work`, `ccam work` e una scorciatoia con `hyperfine` e fallisce sopra i 10 ms, e `BenchmarkLaunchPath` in `internal/launcher` esegue i passaggi in un ciclo. Per vedere il tempo di ogni passaggio, imposta `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Sicurezza e casi limite

- Un nome o un numero sconosciuto si ferma con `no account "<name>"` e rimanda a `ccam list`. Se non ci sono account, ccam ti chiede di aggiungerne uno con `ccam add <name>`.
- Se un percorso condiviso è diventato un file normale o punta al posto sbagliato, ccam si ferma con `stopped because of share link problems:`, elenca i problemi e ti chiede di eseguire `ccam doctor --fix`. Sul percorso di avvio non sposta dati.
- Se un altro Claude Code che condivide la cronologia delle conversazioni con questo account (stesso gruppo) è in esecuzione nella stessa cartella, ccam avvisa che la cronologia potrebbe mescolarsi e chiede `Run anyway? [y/N]`. Senza un terminale mostra l'avviso e continua.
- `--dry-run` salta la sincronizzazione MCP, l'avviso sui nuovi server, il controllo della cronologia, il marcatore di sessione, il titolo e l'ora dell'ultimo uso, e non esegue exec. Controlla comunque i link di condivisione e crea un link mancante come farebbe un avvio normale.
- `ccam work rc` mentre `work` è già in esecuzione lascia stare `.claude.json`, lo segnala e avvia Remote Control senza i server del gruppo; eseguilo di nuovo dopo la fine della sessione per copiarli. Con `--dry-run` stampa il messaggio ma non chiede e non copia.
- Con `mcp_strategy = "sync"`, una sincronizzazione fallita stampa `ccam: MCP sync failed (launching anyway): …` e Claude Code si avvia con i server che aveva già.

## Pagine correlate

- [Comandi scorciatoia](/features/shortcuts)
- [Gestire gli account](/features/accounts)
- [Riferimento dei comandi](/docs/commands)
- [Condividere server MCP, impostazioni e cronologia](/docs/sharing)
- [Come funziona CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
