---
title: Limiti di utilizzo
description: ccam usage mostra insieme, come barre, i limiti su 5 ore, settimanali e per modello (Fable) di tutti gli account, e ccam statusline mette il nome dell'account e l'utilizzo in Claude Code.
updated: 2026-10-21
order: 7
translation: machine
sourceHash: "0e9587b79ce8"
---

`ccam usage` mostra quanta parte del limite su 5 ore, del limite settimanale e del limite settimanale proprio di ogni modello (come Fable) ha usato ogni account, con una barra per limite e il tempo mancante al suo azzeramento. Usalo prima di iniziare un'attività lunga, o quando Claude Code ti avvisa che un limite è vicino, per vedere quale account ha ancora margine. `ccam statusline` è il comando per la riga di stato di Claude Code: mostra il nome dell'account e l'utilizzo dentro Claude Code e salva i numeri per `ccam usage`.

## Utilizzo

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam usage --source claude  # let Claude Code look it up: ccam never touches a token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

Ogni account ha una riga con il suo numero, il suo nome e quanto sono vecchi i suoi numeri; `(cached)` significa che provengono da un file salvato e non dall'endpoint. Sotto, una riga per limite mostra una barra di 20 celle, la percentuale e il tempo mancante all'azzeramento della finestra: `5-hour`, `Weekly` (tutti i modelli), poi la finestra settimanale propria di ogni modello, come `Fable`. Una finestra che l'account non ha mostra `-`. Una percentuale diventa gialla dall'80% e rossa al 100%. In un terminale ogni account compare subito come `fetching…`, e le sue righe vengono disegnate quando risponde. L'ultima riga indica l'account che [`ccam next --auto`](/features/next) sceglierebbe con le stesse regole (`auto_max_usage`, `auto_strategy`, le sessioni in esecuzione, solo gli account dell'organizzazione dell'ultimo account usato a meno che non sia impostato `cross_org`, ogni finestra al 100% saltata, così un account con la settimana esaurita non viene mai indicato), e compare solo quando hai più di un account. Un account senza login mostra `not logged in → ccam login <name>`, e un account con chiave API mostra `no limits (API)`.

Dopo `ccam statusline install`, Claude Code mostra una riga come `[work] · Opus · 5h 24% · 7d 41%`.

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `--json` | Stampa i risultati in JSON: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | disattivata |
| `--watch[=N]` | Ridisegna ogni `N` secondi finché non premi Ctrl+C. `--watch` da solo, o un valore pari o inferiore a 0, indica 60 | disattivata |
| `-v`, `--verbose` | Aggiunge sotto ogni account una riga `source:`: da dove vengono i suoi numeri (l'endpoint di utilizzo, Claude Code, la riga di stato o un file salvato) | disattivata |
| `--no-cache` | Ignora la cache di 60 secondi e interroga di nuovo l'endpoint | disattivata |
| `--source live\|claude\|cache` | `live` interroga l'endpoint con il token dell'account; `claude` fa consultare l'utilizzo a Claude Code (vedi sotto); `cache` mostra solo i numeri salvati | `usage_source`, altrimenti `live` |
| `--refresh` | Rinnova i token di accesso scaduti, solo per gli account che non sono in esecuzione | disattivata |
| `--dry-run` | Accettata: `ccam usage` non modifica nulla tranne la propria cache. Rifiutata insieme a `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | disattivata |

Qualsiasi altro valore di `--source` fallisce con `--source must be live, claude or cache`.

### L'utilizzo letto da Claude Code stesso

Con `usage_source = "claude"` in `config.toml` (oppure `--source claude`), ccam non legge affatto il token di login. Per ogni account avvia `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` nella cartella dell'account e invia sullo standard input la richiesta `get_usage` di Claude Code (Claude Code 2.1.270 o successivo); Claude Code consulta l'utilizzo con il proprio login e risponde con le finestre su 5 ore, settimanale e per modello. `--safe-mode` impedisce l'avvio di hook, server MCP e plugin dell'account, e nessuna sessione viene salvata. Ogni account richiede un avvio di Claude Code, circa un secondo, quindi questa fonte è più lenta di `live`; gli account vengono interrogati tutti insieme, ciascuno entro 20 secondi, e la cache di 60 secondi vale come al solito. Claude Code indica la richiesta come sperimentale, quindi se smette di rispondere, ccam mostra i numeri salvati e `ccam doctor` lo segnala. Un `--dry-run` non avvia Claude Code (scrive nella cartella dell'account) e mostra i numeri salvati. `usage_source` accetta `live` o `claude`; `cache` è solo una scelta di `--source`. `ccam usage` non accetta argomenti.

`ccam statusline` legge dallo standard input il JSON della riga di stato di Claude Code. La sua unica opzione è `--wrap <command>`: un comando esistente per la riga di stato, davanti alla cui prima riga di output viene messo il nome dell'account. `ccam statusline install` non accetta opzioni. Entrambi rifiutano `--dry-run` con `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Come funziona

I numeri provengono da due fonti: l'endpoint di utilizzo non documentato di Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, chiamato con il login di ciascun account, e i valori che Claude Code passa al suo comando per la riga di stato. Tutto si trova in `internal/usage`.

1. `cli.usageCmd` carica `config.toml` e costruisce un `usage.Account` per ogni account: la sua cartella `~/.local/share/ccam/profiles/<name>/`, la sua cartella di sessione (per sapere se è in esecuzione) e se è stato importato sul posto.
2. `Fetcher.FetchAll` avvia una goroutine per account. Tutte condividono un solo `http.Client`, così le connessioni all'endpoint vengono riutilizzate. La connessione può richiedere 2 s, e ogni account ha in totale 3 s. Il client non segue i reindirizzamenti, perché un reindirizzamento invierebbe di nuovo il token altrove, e legge al massimo 1 MB della risposta.
3. Quando ogni account termina, la sua riga viene ridisegnata. Un account che fallisce o non risponde in tempo ripiega sui numeri salvati, così un account lento non blocca mai gli altri.

Ogni account interroga l'endpoint nella propria goroutine; quello che supera il tempo limite mostra invece i numeri salvati.

<!-- motion:usage-fetch -->

`personal` e `work` rispondono in tempo e mostrano numeri aggiornati. `side` impiega più di 3 s, quindi `fallback()` legge i suoi due file di cache e usa quello scritto per ultimo. `usage.Table()` disegna un blocco per account e `usage.Best()` sceglie l'account indicato nell'ultima riga.

Per un singolo account, `Fetcher.Fetch` esegue questi controlli in ordine:

<!-- motion:usage-cache -->

1. Con `--source cache`, passa direttamente a `fallback()`.
2. A meno che non sia impostato `--no-cache`, un risultato aggiornato in `~/.local/state/ccam/usage/<name>.json` più recente di 60 secondi (`CacheTTL`) viene restituito così com'è. Eseguire `ccam usage` due volte di seguito costa una richiesta per account.
3. `readCredentials()` legge `.credentials.json` (o il Portachiavi di macOS) senza modificarlo, e viene decodificata solo la chiave `claudeAiOauth`. Nessun login con abbonamento significa `not logged in`, oppure `no limits (API)` quando `.claude.json` contiene una chiave API.
4. Un token di accesso scaduto non viene usato. Senza `--refresh`, l'account mostra i numeri salvati segnati come scaduti.
5. Altrimenti ccam chiama l'endpoint. Un 401 o un 403 conta come login scaduto, qualsiasi altro stato o un corpo che non riesce a leggere conta come errore, e in entrambi i casi si ripiega sui numeri salvati. Un 429 significa che l'endpoint riceve troppe richieste: ccam scrive `<name>.pause` con l'ora in cui chiedere di nuovo (`Retry-After`, oppure 5 minuti, al massimo un'ora), mostra i numeri salvati con `rate limited by the usage endpoint; asking again after 15:10` e non chiama l'endpoint per quell'account fino ad allora, nemmeno con `--no-cache`. Riprovare subito non farebbe che prolungare il limite.
6. Una risposta valida viene scritta in `<name>.json` con l'ora in cui è stata misurata.

`fallback()` legge sia `<name>.json` sia `<name>.statusline.json` e tiene il più recente. Prima di mostrare i numeri salvati, controlla l'orario di azzeramento di ogni finestra: una finestra che si è azzerata dopo il salvataggio dei numeri mostra 0%, perché la vecchia percentuale non vale più. Per questo `personal` qui sopra mostra `0%` senza orario di azzeramento.

L'endpoint ha due formati di risposta. Quello più recente ha un array `limits` con le voci `session` (5 ore), `weekly_all` e `weekly_scoped` (per modello). Quello più vecchio ha gli oggetti `five_hour`, `seven_day`, `seven_day_opus` e `seven_day_sonnet`. `ParseLive` li legge entrambi: prevale l'array `limits`, e i vecchi oggetti completano ciò che non copre. Un corpo senza alcun limite viene segnalato come cambio di formato, di cui ti avvisa `ccam doctor`.

### Rinnovare i token

I refresh token ruotano: rinnovarne uno ti dà un nuovo refresh token e invalida quello vecchio. Se ccam rinnovasse il login mentre Claude Code è in esecuzione con lo stesso account, quella sessione verrebbe disconnessa. Per questo `ccam usage` legge i token e non li scrive mai, a meno che tu non passi `--refresh`. In quel caso `Fetcher.refresh`:

1. rifiuta gli account registrati sul posto con `ccam import`, dato che `claude` da solo potrebbe usare quel login a insaputa di ccam;
2. rifiuta i login salvati nel Portachiavi di macOS;
3. prende un file di lock, `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock` (sotto `~/.local/state/ccam/locks/` quando quella variabile non è impostata), così due processi di ccam non rinnovano mai due volte lo stesso token;
4. con il lock acquisito, si rifiuta se è in esecuzione una sessione dell'account, e salta la richiesta se un altro processo ha già rinnovato il token;
5. invia il refresh token all'endpoint dei token e riscrive i nuovi token in `.credentials.json` con `fsutil.AtomicWrite`, mantenendo tutte le altre chiavi del file.

### Il raccoglitore della riga di stato

Claude Code esegue il comando della riga di stato e gli passa tramite pipe un documento JSON. ccam ne legge il nome del modello e `rate_limits`:

<!-- motion:usage-statusline -->

1. `runStatusline` trova l'account da `CCAM_PROFILE`, che ccam imposta quando avvia Claude Code, oppure confrontando `CLAUDE_CONFIG_DIR` con le cartelle degli account.
2. `ParseStatusLine` legge `rate_limits.five_hour` e `rate_limits.seven_day`. Il JSON non contiene limiti per gli account con chiave API, né prima della prima risposta della sessione; in quel caso non viene salvato nulla.
3. `SaveStatusLine` scrive `~/.local/state/ccam/usage/<name>.statusline.json`. `ccam usage` lo usa quando l'endpoint non risponde.
4. ccam stampa `[name] · model · 5h N% · 7d N%`. Con `--wrap`, esegue invece il tuo comando tramite `sh -c` con lo stesso input e stampa `[name]` seguito dalla prima riga del comando. Su Windows esegue il comando come fa Claude Code: con Git Bash se è installato, altrimenti con PowerShell.
5. Quando l'utilizzo su 5 ore o settimanale è al 90% o oltre e c'è un altro account, la riga termina con `→ ccam next --auto --continue`, per esempio `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

Il comando avvolto ha 2 secondi. Dopo viene terminato e la riga mostra solo il nome dell'account, così un comando lento non può bloccare la riga di stato di Claude Code.

`ccam statusline install` scrive `statusLine` nel `settings.json` di ogni account. Quando un account ha già un comando per la riga di stato, ccam lo mantiene come `--wrap '<old command>'`. Su Windows il percorso di ccam viene scritto con barre normali e `~` per la cartella home (`~/.local/bin/ccam.exe statusline`), che sia Git Bash sia PowerShell capiscono. Un `settings.json` condiviso tramite un gruppo viene scritto una sola volta, attraverso il suo link, così lo riceve ogni account del gruppo. Eseguirlo di nuovo stampa `(already set)` per i file che usano già ccam.

## Sicurezza e casi limite

- `ccam usage` non modifica mai un file di token a meno che tu non passi `--refresh`. Quando un login è scaduto, avvia l'account una volta con `ccam <name>` e sarà Claude Code a rinnovarlo.
- `--refresh` non tocca mai gli account in esecuzione, gli account importati o i login nel Portachiavi. Un rinnovo rifiutato mostra i numeri salvati e il motivo.
- L'endpoint di utilizzo non è documentato pubblicamente e potrebbe cambiare. Se succede, vedi comunque i numeri salvati, e `ccam doctor` segnala il cambio di formato.
- Senza numeri salvati e senza risposta, la riga dell'account mostra `lookup failed:` e il motivo, per esempio `no saved usage` con `--source cache`.
- I file di cache contengono solo percentuali e orari di azzeramento, mai token.

## Pagine correlate

- [Limiti di utilizzo e cambio di account](/docs/usage-limits)
- [Passare all'account successivo](/features/next)
- [Flusso di lavoro con i limiti di utilizzo](/guides/usage-limit-workflow)
- [Comandi](/docs/commands)
