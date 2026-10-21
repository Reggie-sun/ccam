---
title: Limiti di utilizzo e cambio di account
description: Visualizza come barre l'utilizzo su 5 ore, settimanale e di Fable di ogni account Claude Code con ccam usage e apri quello con più margine con ccam next --auto.
updated: 2026-10-21
order: 5
translation: machine
sourceHash: "0e64e3251c69"
---

`ccam usage` mostra quanta parte del limite su 5 ore, del limite settimanale e del limite settimanale proprio di ogni modello (come Fable) ha usato ogni account. `ccam next --auto` apre l'account con più margine, e `--continue` porta con sé la conversazione corrente.

```console
$ ccam usage
 1 personal  just now
   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h
   Fable    ██████████████████░░  88%  resets in 1d 6h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
   Fable    █░░░░░░░░░░░░░░░░░░░   3%  resets in 5d 2h
 3 side      12m ago (cached)
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
   Fable    ████████░░░░░░░░░░░░  40%  resets in 1d 9h
→ Most headroom: work   (open it with ccam next --auto)
```

Ogni account ha una riga che indica quanto sono recenti i suoi dati, poi una barra per limite con il tempo mancante all'azzeramento: la finestra di 5 ore, la finestra settimanale per tutti i modelli e la finestra settimanale propria di ogni modello, come Fable. Dall'80% una barra diventa gialla, e al 100% rossa. L'ultima riga indica l'account che `ccam next --auto` aprirebbe.

## Da dove vengono i numeri

ccam interroga l'endpoint di utilizzo di Anthropic per tutti gli account contemporaneamente, usando il login di ciascun account. Attende al massimo 3 secondi per account e mostra ogni account appena risponde.

Questo endpoint non è documentato pubblicamente e potrebbe cambiare. Quando una richiesta non va a buon fine, ccam mostra gli ultimi numeri noti e li segna come presi dalla cache. Una finestra che si è azzerata dopo il salvataggio di quei numeri mostra 0%. `ccam doctor` ti avvisa se il formato della risposta è cambiato.

I numeri ottenuti meno di 60 secondi prima vengono riutilizzati. Opzioni:

| Opzione | Effetto |
| --- | --- |
| `--watch [seconds]` | Aggiorna di continuo, ogni 60 secondi per impostazione predefinita |
| `--json` | Output leggibile da una macchina |
| `-v` | Mostra anche da dove vengono i numeri di ogni account |
| `--no-cache` | Chiede sempre di nuovo |
| `--source cache` | Non chiede mai; mostra solo i numeri salvati |

## Login scaduti

Quando il token di accesso di un account è scaduto, ccam non lo rinnova, perché rinnovarlo mentre una sessione di Claude Code usa lo stesso login potrebbe disconnettere quella sessione. Avvia l'account una volta con `ccam <name>` e sarà Claude Code a rinnovarlo.

`ccam usage --refresh` rinnova da solo i token scaduti, ma solo per gli account che non sono in esecuzione. Non rinnova mai i login salvati nel Portachiavi di macOS o in una cartella `~/.claude` importata.

## Utilizzo nella riga di stato

```sh
ccam statusline install
```

Claude Code mostra quindi il nome dell'account e il suo utilizzo nella riga di stato, per esempio `[work] · Opus · 5h 24% · 7d 41%`. Dal 90% in su, la riga suggerisce anche `ccam next --auto --continue`. Se hai già un comando per la riga di stato, ccam lo mantiene e mette davanti il nome dell'account. La riga di stato salva anche i propri numeri, che `ccam usage` usa quando l'endpoint non risponde.

## Cambiare account

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # carry this folder's conversation on as the next account
```

`ccam next` salta gli account che non hanno effettuato il login, e `--auto` salta gli account al 100% e quelli di un'altra organizzazione. Quando i due account condividono la cronologia delle conversazioni, `--continue` apre la stessa conversazione; altrimenti Claude Code la copia nel nuovo account. Vedi [condivisione](/docs/sharing) e [passare all'account successivo](/features/next#continue).

Per script e job pianificati che eseguono `claude -p`, metti `--fallback` prima del nome dell'account. Quando quell'account raggiunge il limite, ccam riesegue lo stesso comando con l'account con più margine, e viene stampato solo l'output dell'esecuzione portata a termine:

```sh
ccam --fallback work -p "review the diff" --output-format json
```

Vedi [Rieseguire al limite](/features/next#run-again-at-the-limit) per come decide ccam.

La [guida al flusso di lavoro per i limiti di utilizzo](/guides/usage-limit-workflow) illustra un esempio completo.
