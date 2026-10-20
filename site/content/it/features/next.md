---
title: Passare all'account successivo
description: ccam next apre Claude Code con un altro account; --auto sceglie quello con più margine, --continue porta con sé la conversazione e --fallback riesegue i job -p.
updated: 2026-10-20
order: 8
translation: machine
sourceHash: "6a5dcb1d5c01"
---

`ccam next` apre Claude Code con un altro account senza che tu debba indicarne il nome. Usalo quando un account raggiunge il suo limite di utilizzo: `ccam next` da solo prende l'account che segue quello usato per ultimo, `--auto` prende quello con più margine e `--continue` porta l'ultima conversazione di questa cartella nel nuovo account, che i due condividano o no la cronologia. Dentro Claude Code, il [plugin di ccam](#inside-claude-code) ti avvisa quando un limite ferma una sessione e `/ccam:switch` passa la conversazione a un altro account. Per gli script che eseguono `claude -p`, `ccam --fallback` fa il cambio da solo (vedi [Rieseguire al limite](#run-again-at-the-limit)).

## Utilizzo

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the weekly quota that would go unused first
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam stampa il cambio come `→ <from> → <to>` sullo standard error, con il motivo di una scelta `--auto`, poi avvia Claude Code con il nuovo account come fa `ccam <name>`. `<from>` è l'account usato per ultimo, oppure, con `--continue`, l'account a cui appartiene la conversazione. Gli argomenti dopo `--` arrivano a `claude` invariati.

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `--auto` | Sceglie l'account con login attivo in base all'utilizzo (vedi sotto), saltando gli account al 100% | disattivata |
| `--max-usage <1-100>` | Con `--auto`: preferisce gli account sotto questo utilizzo | `auto_max_usage`, altrimenti 100 |
| `--strategy headroom\|expiring\|balance` | Con `--auto`: l'utilizzo su 5 ore più basso, l'azzeramento su 5 ore più vicino, oppure la quota settimanale che andrebbe persa per prima | `auto_strategy`, altrimenti `headroom` |
| `--cross-org` | Con `--auto`: sceglie anche account di un'altra organizzazione (vedi sotto) | `cross_org`, altrimenti disattivata |
| `--wait[=<duration>]` | Con `--auto`: quando nessun account ha margine, attende l'azzeramento del primo, al massimo per questo tempo | disattivata; `5h` se indicata da sola |
| `-c`, `--continue` | Porta l'ultima conversazione di questa cartella nel nuovo account (vedi sotto) | disattivata |
| `-- <claude args>` | Argomenti passati a `claude` dopo quelli di ccam | nessuno |
| `--dry-run` | Rispettata: stampa `→ <from> → <to>` e la riga di comando che eseguirebbe, senza avviare Claude Code né registrare un avvio | disattivata |

`--auto` e `--continue` funzionano insieme. Con `--continue`, gli argomenti di `--resume` vengono messi prima dei tuoi.

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

`--auto` decide in base all'utilizzo invece che all'ordine degli slot. `cli.pickAuto` tiene gli account con login, chiede l'utilizzo di tutti contemporaneamente con lo stesso `Fetcher.FetchAll` usato da `ccam usage` e lascia scegliere a `usage.Policy.Choose()`:

<!-- motion:next-auto -->

Un account viene valutato in base alle finestre di 5 ore e settimanale e, quando `--model` è tra gli argomenti per Claude Code, anche in base alla finestra settimanale di quel modello (`opus`, `sonnet`, `haiku` o `fable` nel nome). `Policy.Choose()` divide gli account in tre gruppi:

1. **Con margine**: ogni finestra sotto `auto_max_usage` (`config.toml`, predefinito 100). La strategia sceglie tra questi: `headroom` (predefinita) prende l'utilizzo su 5 ore più basso; `expiring` prende la finestra di 5 ore che si azzera prima, così la sua quota residua viene usata prima di andare persa; `balance` prende l'account con più quota settimanale residua per ogni giorno che manca all'azzeramento della sua settimana, contando a suo carico ogni sessione ccam in esecuzione, così la quota settimanale non resta inutilizzata e le sessioni in parallelo si distribuiscono.
2. **Oltre la soglia**: sotto il 100% ma pari o superiore a `auto_max_usage`. Usati solo quando nessun account ha margine, partendo dall'utilizzo su 5 ore più basso.
3. **Pieni**: una finestra al 100%, oppure utilizzo su 5 ore sconosciuto. Mai scelti.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organizzazioni.** Le scelte automatiche restano tra account simili a quello che lasci: la stessa organizzazione team o enterprise, oppure i piani personali (Pro, Max) tra loro, così il lavoro di un account aziendale non prosegue su uno privato né viceversa. ccam conosce l'organizzazione di ogni account dal suo login (`ccam login` la registra; per i login più vecchi viene letta da `.claude.json`); un account di cui non riesce a saperla va bene ovunque. Due cartelle con il login nello stesso account Claude ne condividono i limiti, quindi una non subentra mai all'altra. `--cross-org`, oppure `cross_org = true` in `config.toml`, toglie la regola dell'organizzazione.

Quando non si può scegliere nulla, ccam esce con codice 75 (`EX_TEMPFAIL`), così uno script può distinguere "riprova più tardi" da un errore. Con `--wait`, invece, attende finché la finestra bloccante del primo account pieno non si azzera (30 secondi dopo l'orario di azzeramento, quando l'endpoint mostra la nuova finestra), stampa `ccam: no account has room; waiting until 15:10` e sceglie di nuovo; quando nessun account dice quando si azzera, ricontrolla dopo 10 minuti. Se nulla si libera prima della fine dell'attesa, esce con 75. La richiesta segue le regole di `ccam usage`: i numeri ottenuti meno di 60 secondi prima vengono riutilizzati, ogni account ha 3 secondi e un account che non risponde, o il cui login è scaduto, viene valutato in base ai numeri salvati. Una finestra salvata il cui orario di azzeramento è passato conta come 0%, così un account la cui finestra si è nel frattempo azzerata non viene scartato. `--auto` non rinnova mai i token. A differenza della semplice rotazione, `--auto` può scegliere l'account usato per ultimo se è ancora quello con più margine.

### --continue

1. ccam cerca l'ultima conversazione di questa cartella in ogni account: il `<session>.jsonl` più recente in `projects/<folder>/` (il nome della cartella è la cartella di lavoro con ogni carattere che non sia una lettera o una cifra trasformato in `-`, come la chiama Claude Code). Se la sessione ha lasciato una nota (un limite di utilizzo o `/ccam:switch`, vedi sotto) e da allora non è stato scritto nulla di più recente, viene usata proprio quella sessione.
2. L'account a cui appartiene è quello da cui passi a un altro, e il nuovo account viene scelto come sopra, dentro la stessa organizzazione.
3. Se il nuovo account condivide quella cronologia (la stessa cartella di gruppo), ccam avvia `claude --resume <session id>`, proprio quella conversazione (`--continue` aprirebbe la più recente della cartella, qualunque sia).
4. Altrimenti avvia `claude --resume <the conversation's file> --fork-session`: Claude Code copia l'intera conversazione, con un nuovo id di sessione, nella cronologia propria del nuovo account, e l'originale resta com'era. Da quel momento il `--continue` del nuovo account la trova.

ccam dice quanto è grande la conversazione: il nuovo account la rilegge tutta con la sua prima risposta. Le cache dei prompt non vengono mai condivise tra organizzazioni, quindi con `--cross-org` quella prima risposta costa l'intera dimensione. Condividere `history` in un gruppo evita la copia, mantiene un'unica conversazione e condivide anche i checkpoint da cui `/rewind` ripristina il codice (`file-history/`).

## Dentro Claude Code

Il [plugin di ccam](/docs/install#claude-code-plugin) aggiunge due hook (Claude Code 2.1.78 o successivo):

- Quando un turno si ferma a un limite di utilizzo (`You've hit your session limit`, il limite settimanale o quello di un modello), annota la conversazione e mostra una notifica desktop nei terminali che la supportano (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): `ccam: work hit its 5-hour limit (resets 3:45pm). ccam next --continue carries on as another account.` Un breve rallentamento del server o il limite di frequenza di una chiave API non è un limite di utilizzo e viene ignorato. Claude Code 2.1.234 e successivi inoltre attendono nella sessione e proseguono da soli dopo l'azzeramento delle 5 ore; la notifica ti dà la possibilità di cambiare account invece di aspettare.
- `/ccam:switch [account]` annota la conversazione e chiude il turno prima che venga fatta qualsiasi richiesta: digita `/exit`, poi `ccam next --continue` la porta avanti, con l'account che hai indicato o con quello successivo.

La nota resta in `~/.local/state/ccam/handoff/<account>.json` per 12 ore e viene usata una sola volta, quando la sessione successiva sta per partire; un avvio che fallisce la lascia per il tentativo successivo. Appartiene alla cartella in cui è partita la sessione, anche dopo un `cd` nella shell di Claude Code. Se nel frattempo l'account che indica è stato rinominato o rimosso, ccam lo segnala e sceglie come al solito. Con `relay = true` in `config.toml`, ccam resta accanto alle sessioni interattive invece di sostituirsi con Claude Code, e dopo `/ccam:switch` porta avanti da solo la conversazione di quella sessione (una nota lasciata da un'altra sessione non viene presa in carico). Questo costa un processo ccam in attesa per ogni sessione, ed è per questo che è disattivato per impostazione predefinita. La sessione che prosegue la conversazione parte con `ccam next --continue`, senza gli argomenti con cui era stata avviata la prima, come `--model`.

## Rieseguire al limite

`ccam next` serve per una sessione interattiva che si è fermata. Un'esecuzione non interattiva, `claude -p` in uno script o in un job pianificato, non può aspettarti, quindi `--fallback` prima del nome dell'account permette a ccam di cambiare account al posto suo:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` funziona solo insieme a `-p` o `--print`; altrimenti ccam si ferma e rimanda a `ccam next --auto --continue`.
2. Lo standard input viene letto una volta sola e passato a ogni tentativo.
3. ccam avvia `claude` come processo figlio invece di usare `exec`, trattiene il suo standard output e lascia passare lo standard error man mano che arriva.
4. Quando `claude` esce con 0, o fallisce per un altro motivo, l'output trattenuto e il codice di uscita vengono passati invariati.
5. Un errore conta come limite di utilizzo quando il messaggio di Claude Code lo dice (`You've hit your limit`, `usage limit reached`), oppure quando la consultazione dell'utilizzo dell'account mostra la finestra di 5 ore o quella settimanale al 100%.
6. A quel punto l'output di quel tentativo viene scartato e lo stesso comando viene eseguito con un altro account connesso: quello con l'utilizzo su 5 ore più basso, poi gli account il cui utilizzo non si può consultare, in ordine di numero. Gli account al 100% vengono saltati, e ogni account viene provato al massimo una volta.
7. Con `-c`, `--continue`, `-r` o `--resume` tra gli argomenti vengono provati solo gli account dello stesso gruppo di cronologia, così la conversazione può essere trovata.
8. Vengono provati solo gli account della stessa organizzazione, come con `--auto` (`cross_org = true` toglie questa regola).

L'account successivo viene scelto con le stesse regole di `--auto` (`auto_max_usage`, `auto_strategy`, la finestra di `--model`). Quando non resta alcun account, ccam stampa l'output dell'ultimo tentativo, dice che nessun account connesso ha ancora margine ed esce con codice 75. Con `--wait` prima del nome dell'account (`ccam --fallback --wait=2h work -p …`), attende invece l'azzeramento del primo account, poi può provare di nuovo ogni account; attende sempre almeno una volta prima, così un account che ha appena raggiunto il limite non viene rieseguito subito. `ccam --dry-run --fallback work -p …` stampa il comando del primo tentativo.

Le sessioni interattive non cambiano account da sole. Invece, con `ccam statusline install`, la riga di stato aggiunge `→ ccam next --auto --continue` quando l'utilizzo su 5 ore o settimanale arriva al 90% ([Limiti di utilizzo](/features/usage)).

## Sicurezza e casi limite

- Gli account senza login non vengono mai scelti. Se non c'è un altro account con login, `ccam next` da solo si ferma con `no logged-in account other than <name>`.
- Con `--continue` e nessuna conversazione nella cartella: `no account has a conversation in <folder> to continue`.
- Quando la regola dell'organizzazione non lascia alcun account: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Con `--auto` e nessun account con login: `no logged-in account to pick`. Quando ogni account è pieno o ha un utilizzo sconosciuto: `` no account with known usage and headroom (see `ccam usage`) ``. In entrambi i casi il codice di uscita è 75.
- `--max-usage`, `--strategy`, `--wait` e `--cross-org` funzionano solo con `--auto`, e un valore fuori da 1-100 o un'altra strategia vengono rifiutati.
- Un fork copia la conversazione, non il resto dello stato dell'account: i checkpoint nel `file-history/` dell'account originale non vengono con essa, a meno che la cronologia non sia condivisa.
- `--auto` effettua richieste di rete anche con `--dry-run`. Si limita a leggere e non scrive mai i file dei token.
- L'account scelto si avvia tramite lo stesso percorso di `ccam <name>`, quindi valgono gli stessi controlli, compreso l'avviso quando in questa cartella è già in esecuzione un altro Claude Code che condivide la stessa cronologia.
- `last-used` viene scritto solo quando Claude Code si avvia davvero. `--dry-run` non lo tocca, quindi un dry run non fa avanzare la rotazione.
- `--fallback` trattiene lo standard output finché un tentativo non termina, quindi un output in streaming come `--output-format stream-json` arriva tutto insieme alla fine. I file che il primo tentativo ha modificato prima di raggiungere il limite restano modificati; il tentativo successivo riparte dall'inizio del prompt.

## Pagine correlate

- [Limiti di utilizzo](/features/usage)
- [Limiti di utilizzo e cambio di account](/docs/usage-limits)
- [Condividere server MCP, impostazioni e cronologia](/docs/sharing)
- [Flusso di lavoro con i limiti di utilizzo](/guides/usage-limit-workflow)
