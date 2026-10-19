---
title: Passare all'account successivo
description: ccam next apre Claude Code con l'account che segue quello usato per ultimo; --auto sceglie l'utilizzo su 5 ore più basso e --continue riprende la conversazione.
updated: 2026-10-19
order: 8
translation: machine
sourceHash: 92352f532e74
---

`ccam next` apre Claude Code con un altro account senza che tu debba indicarne il nome. Usalo quando un account raggiunge il suo limite di utilizzo: `ccam next` da solo prende l'account che segue quello usato per ultimo, `--auto` prende quello con più margine nella finestra di 5 ore e `--continue` prosegue l'ultima conversazione nel nuovo account.

## Utilizzo

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next
→ personal → work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude
```

ccam stampa il cambio come `→ <from> → <to>` sullo standard error, poi avvia Claude Code con il nuovo account come fa `ccam <name>`. `<from>` è l'account usato per ultimo. Gli argomenti dopo `--` arrivano a `claude` invariati.

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `--auto` | Sceglie l'account con login attivo che ha l'utilizzo su 5 ore più basso, saltando gli account al 100% | disattivata |
| `-c`, `--continue` | Considera solo gli account che condividono la cronologia delle conversazioni con l'account usato per ultimo e avvia `claude --continue` | disattivata |
| `-- <claude args>` | Argomenti passati a `claude` dopo quelli di ccam | nessuno |
| `--dry-run` | Rispettata: stampa `→ <from> → <to>` e la riga di comando che eseguirebbe, senza avviare Claude Code né registrare un avvio | disattivata |

`--auto` e `--continue` funzionano insieme. Con `--continue`, `--continue` viene messo prima dei tuoi argomenti.

## Come funziona

Ogni avvio tramite ccam scrive l'ora in `~/.local/state/ccam/last-used/<name>`. `ccam next` legge questi file, controlla quali account hanno effettuato il login e sceglie.

1. `cli.candidates` costruisce un `next.Candidate` per ogni account a partire da `config.toml`: l'account, l'ora del suo ultimo avvio (`launcher.LastUsed`) e se ha un login (`creds.Has`: un `.credentials.json` leggibile oppure, su macOS, una voce nel Portachiavi).
2. Senza `--auto`, `next.Rotate` ordina gli account per numero di slot e trova quello avviato più di recente. È l'account da cui stai passando a un altro. Se nessun account è ancora stato avviato, è lo slot più basso.
3. Da lì procede in avanti nell'ordine degli slot, ricominciando dal primo dopo l'ultimo, e restituisce il primo account con un login.
4. `launcher.LaunchWith` prepara e avvia l'account scelto esattamente come farebbe `ccam <name>`: controlla i link dell'account, applica la condivisione degli MCP, registra la sessione e sostituisce ccam con `claude` tramite `syscall.Exec`.

Il percorso parte dall'account usato per ultimo e si ferma al primo account con login che lo segue:

<!-- motion:next-rotate -->

`work` è stato avviato più di recente, quindi il percorso parte da quello successivo. `side` non ha un login e viene saltato, e viene scelto `lab`. Se anche `lab` fosse stato disconnesso, il percorso sarebbe ricominciato da `personal` (`% len` nel diagramma). La rotazione non restituisce mai l'account da cui sei partito: se non c'è un altro account con login, si ferma con un errore.

### --auto

`--auto` decide in base all'utilizzo invece che all'ordine degli slot. `cli.pickAuto` tiene gli account con login, chiede l'utilizzo di tutti contemporaneamente con lo stesso `Fetcher.FetchAll` usato da `ccam usage` e lascia scegliere a `usage.Best()`:

<!-- motion:next-auto -->

`usage.Best()` prende la percentuale su 5 ore più bassa e salta gli account al 100% o con utilizzo su 5 ore sconosciuto. La richiesta segue le regole di `ccam usage`: i numeri ottenuti meno di 60 secondi prima vengono riutilizzati, ogni account ha 3 secondi e un account che non risponde, o il cui login è scaduto, viene valutato in base ai numeri salvati. Una finestra salvata il cui orario di azzeramento è passato conta come 0%, così un account la cui finestra si è nel frattempo azzerata non viene scartato. `--auto` non rinnova mai i token. A differenza della semplice rotazione, `--auto` può scegliere l'account usato per ultimo se è ancora quello con più margine.

### --continue

L'opzione `--continue` di Claude Code apre la conversazione più recente nella cartella corrente, il che funziona solo se il nuovo account può vedere quella conversazione. Per questo, con `--continue`, ccam considera solo gli account nello stesso gruppo `history` dell'account usato per ultimo, poi mette `--continue` davanti ai tuoi argomenti. Gli account la cui cronologia è `local` non condividono con nessuno, quindi `--continue` da un account del genere fallisce sempre. Attiva prima la condivisione con `ccam share <name> --history on` (vedi [condivisione](/docs/sharing)).

## Sicurezza e casi limite

- Gli account senza login non vengono mai scelti. Se non c'è un altro account con login, `ccam next` da solo si ferma con `no logged-in account other than <name>`.
- Con `--continue` e nessun account che condivide la cronologia: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- Con `--auto` e nessun account con login: `no logged-in account to pick`. Quando ogni account è pieno o ha un utilizzo sconosciuto: `` no account with known usage and headroom (see `ccam usage`) ``.
- `--auto` effettua richieste di rete anche con `--dry-run`. Si limita a leggere e non scrive mai i file dei token.
- L'account scelto si avvia tramite lo stesso percorso di `ccam <name>`, quindi valgono gli stessi controlli, compreso l'avviso quando in questa cartella è già in esecuzione un altro Claude Code che condivide la stessa cronologia.
- `last-used` viene scritto solo quando Claude Code si avvia davvero. `--dry-run` non lo tocca, quindi un dry run non fa avanzare la rotazione.

## Pagine correlate

- [Limiti di utilizzo](/features/usage)
- [Limiti di utilizzo e cambio di account](/docs/usage-limits)
- [Condividere server MCP, impostazioni e cronologia](/docs/sharing)
- [Flusso di lavoro con i limiti di utilizzo](/guides/usage-limit-workflow)
