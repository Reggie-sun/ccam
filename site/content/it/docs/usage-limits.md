---
title: Limiti di utilizzo e cambio di account
description: Visualizza l'utilizzo su 5 ore e settimanale di ogni account Claude Code con ccam usage e apri quello con più margine con ccam next --auto.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 8ab41501c548
---

`ccam usage` mostra quanta parte dei limiti su 5 ore e settimanali ha usato ogni account. `ccam next --auto` apre l'account con più margine, e `--continue` porta con sé la conversazione corrente.

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ██████░░  72% 1h 40m    ███░░░░░  41% 3d 6h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
 3 side      ████████ 100% 22m       █████░░░  68% 1d 9h     12m ago (cached)
→ Most headroom: work   (open it with ccam next --auto)
```

Le colonne sono: account, utilizzo su 5 ore, tempo mancante all'azzeramento, utilizzo settimanale, tempo mancante all'azzeramento e quanto sono recenti i dati.

## Da dove vengono i numeri

ccam interroga l'endpoint di utilizzo di Anthropic per tutti gli account contemporaneamente, usando il login di ciascun account. Attende al massimo 3 secondi per account e mostra ogni account appena risponde.

Questo endpoint non è documentato pubblicamente e potrebbe cambiare. Quando una richiesta non va a buon fine, ccam mostra gli ultimi numeri noti e li segna come presi dalla cache. `ccam doctor` ti avvisa se il formato della risposta è cambiato.

I numeri ottenuti meno di 60 secondi prima vengono riutilizzati. Opzioni:

| Opzione | Effetto |
| --- | --- |
| `--watch [seconds]` | Aggiorna di continuo, ogni 60 secondi per impostazione predefinita |
| `--json` | Output leggibile da una macchina |
| `-v` | Mostra anche i limiti settimanali specifici per modello e da dove viene ogni numero |
| `--no-cache` | Chiede sempre di nuovo |
| `--source cache` | Non chiede mai; mostra solo i numeri salvati |

## Login scaduti

Quando il token di accesso di un account è scaduto, ccam non lo rinnova, perché rinnovarlo mentre una sessione di Claude Code usa lo stesso login potrebbe disconnettere quella sessione. Avvia l'account una volta con `ccam <name>` e sarà Claude Code a rinnovarlo.

`ccam usage --refresh` rinnova da solo i token scaduti, ma solo per gli account che non sono in esecuzione. Non rinnova mai i login salvati nel Portachiavi di macOS o in una cartella `~/.claude` importata.

## Utilizzo nella riga di stato

```sh
ccam statusline install
```

Claude Code mostra quindi il nome dell'account e il suo utilizzo nella riga di stato, per esempio `[work] · Opus · 5h 24% · 7d 41%`. Se hai già un comando per la riga di stato, ccam lo mantiene e mette davanti il nome dell'account. La riga di stato salva anche i propri numeri, che `ccam usage` usa quando l'endpoint non risponde.

## Cambiare account

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # the next account that shares your conversation history
```

`ccam next` salta gli account che non hanno effettuato il login, e `--auto` salta gli account al 100%. Per continuare una conversazione, entrambi gli account devono condividere la cronologia delle conversazioni; vedi [condivisione](/docs/sharing).

La [guida al flusso di lavoro per i limiti di utilizzo](/guides/usage-limit-workflow) illustra un esempio completo.
