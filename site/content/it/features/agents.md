---
title: Sessioni in background di ogni account
description: ccam agents elenca in una tabella gli agenti in background e le sessioni aperte di Claude Code di tutti i tuoi account; ccam attach ne apre una per id.
updated: 2026-10-20
order: 16
translation: machine
sourceHash: "1f459764d66d"
---

Claude Code tiene le sessioni in background (`claude --bg`, la vista degli agenti) per account: `claude agents` mostra solo le sessioni dell'account con cui è in esecuzione. `ccam agents` interroga tutti gli account contemporaneamente e le elenca insieme, e `ccam attach <id>` apre una sessione in background con l'account che l'ha avviata, qualunque esso sia.

## Utilizzo

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique start of the id is enough
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` è lo stato della sessione in background (`working`, `blocked`, `done`, `failed`, `stopped`) e ciò che sta facendo il suo processo (`busy`, `waiting`, `idle`); una sessione interattiva aperta mostra `open`.

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `--all` | Elenca anche le sessioni in background terminate, fallite o fermate | disattivata |
| `--json` | Stampa le sessioni in JSON, ciascuna con l'`account` a cui appartiene | disattivata |

`ccam attach <id> [claude args]` passa a `claude attach` gli argomenti che seguono l'id. Con `--dry-run` stampa il comando.

## Come funziona

1. Per ogni account, ccam esegue `claude agents --json` (con `--all` se indicata) con la `CLAUDE_CONFIG_DIR` dell'account, tutti contemporaneamente, ciascuno entro 15 secondi. Questo legge l'elenco delle sessioni proprio dell'account e non avvia il servizio in background di Claude Code.
2. Ogni voce riceve il nome dell'account, e l'elenco viene ordinato per ora di avvio, dalla più recente. Un account che non risponde viene indicato sullo standard error ed escluso, e poi ccam esce con codice 1, così uno script non scambia un errore per un elenco vuoto.
3. `ccam attach` cerca `jobs/<id>` nella cartella di ogni account, dove Claude Code tiene una sessione in background, ed esegue `claude attach <id>` con quell'account, nello stesso modo in cui `ccam <name>` avvia Claude Code.

## Sicurezza e casi limite

- Gli id delle sessioni in background appartengono a un solo account: `claude attach` eseguito con un altro account non riesce a trovarli, ed è proprio questo che `ccam attach` risolve.
- Due account potrebbero avere sessioni i cui id iniziano allo stesso modo; in quel caso `ccam attach` le indica entrambe e chiede una parte più lunga dell'id.
- ccam passa i suoi server MCP condivisi a Claude Code con `--mcp-config=<file>` in un solo argomento, che `claude agents` e `claude attach` accettano; le sessioni in background continuano a usare quel file quando si riavviano.
- Il servizio in background di un account con chiave API prende la chiave dalla shell che lo ha avviato. Avvia tramite ccam la prima sessione in background di un account del genere (`ccam gw --bg …`), così il servizio ha la chiave.

## Pagine correlate

- [Avviare Claude Code](/features/launch)
- [Scegliere un account per gli script](/features/pick)
- [La dashboard](/features/dashboard)
