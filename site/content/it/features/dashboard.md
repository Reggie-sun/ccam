---
title: La dashboard
description: ccam in un terminale apre una dashboard a schermo intero con account, gruppi, regole per cartella, estensioni e controlli; ogni modifica mostra prima il piano.
updated: 2026-10-20
order: 15
translation: machine
sourceHash: "decdd5fa54c9"
---

`ccam` senza argomenti apre una dashboard a schermo intero quando standard input e standard output sono un terminale; altrimenti stampa `ccam list`. La dashboard mostra ogni account con il suo utilizzo e i suoi gruppi, e fa ciò che fanno i comandi: avviare, aggiungere, effettuare il login, cambiare i gruppi, impostare le regole delle cartelle, controllare e riparare. Avviare un account la chiude; aggiungere un account o effettuare il login la lascia per il terminale e poi vi ritorna.

```console
 ccam v0.9.0 · 3 accounts · 2 groups         ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health

    # Account     Email             5-hour         7-day Last  Groups
 ❯ 1. personal ✔   me@example.com    █████▊░░  72%   41% 2h    default
   2. work      ▸• dev@company.com   ██░░░░░░  25%  100% now   default team
   3. side         -                 ✘ not logged in         -     no group

 ✔ personal is the default account
 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

Il primo fotogramma viene disegnato a partire dal solo `config.toml`, in molto meno di 100 ms; l'utilizzo compare man mano che ogni account risponde. Una percentuale diventa gialla dall'80% e rossa al 100%, e i valori presi da una cache dopo una consultazione fallita vengono mostrati attenuati.

## Schede

| Scheda | Cosa mostra | Tasti |
| --- | --- | --- |
| Accounts | numero, nome, email, barra delle 5 ore, utilizzo su 7 giorni, ultimo uso, gruppi; `✔` l'account predefinito, `▸` l'account di questa cartella, `•` in esecuzione | vedi sotto |
| Groups | ogni gruppo, cosa condivide e i suoi account, e gli account che non fanno parte di alcun gruppo | `enter` modifica, `n` nuovo, `e` rinomina, `x` rimuovi |
| Folders | le regole di [`ccam folder`](/features/folder) e quella che vale qui | `n` usa un account qui, `enter` cambia, `x` rimuovi |
| Extensions | ciò che mostra [`ccam ext`](/features/ext), ogni estensione una sola volta con i gruppi da cui proviene | `r` ricarica |
| Health | ciò che segnala [`ccam doctor`](/features/doctor) | `f` ripara (`doctor --fix`) |

`←`/`→` o `Tab` cambiano scheda, `?` elenca i tasti della scheda corrente, `r` consulta di nuovo l'utilizzo senza la cache breve, e `q` o `Esc` escono.

## Tasti della scheda Accounts

| Tasto | Cosa fa |
| --- | --- |
| `Enter` | Avvia Claude Code con l'account selezionato |
| `1`–`9` | Si sposta sull'account con quel numero; digita `1` `2` in rapida successione per il 12. Solo `Enter` avvia |
| `↑`/`↓`, `k`/`j` | Sposta la selezione |
| `n` | Avvia l'account che sceglierebbe [`ccam next --auto`](/features/next) |
| `o` | Avvia con opzioni: `--continue`, `--resume`, altri argomenti, oppure un dry run |
| `u` | Dettagli dell'utilizzo: le finestre su 5 ore, su 7 giorni e per modello, con i loro azzeramenti e la provenienza dei numeri |
| `a` | Aggiunge un account (la procedura guidata qui sotto) |
| `l` | Esegue di nuovo il login |
| `s` | Sceglie i gruppi dell'account |
| `d` | Lo rende l'account predefinito |
| `f` | Lo usa nella cartella corrente |
| `e` | Lo rinomina |
| `c` | Imposta o rimuove il suo comando scorciatoia |
| `x` | Lo rimuove (chiede prima, partendo da No) |

## Cambiare i gruppi

Nella scheda Groups, `Enter` apre l'editor: gli elementi a sinistra, gli account a destra, `Tab` per passare dall'uno all'altro, `Space` per attivare o disattivare. Un elemento che uno degli account riceve già da un altro gruppo viene mostrato disattivato con il motivo, così come un account che riceverebbe un elemento due volte; un elemento può provenire da un solo gruppo. `Enter` mostra il piano:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Per un file diverso dalla copia del gruppo viene chiesto prima cosa fare (tenere quella del gruppo, usare la mia o unire). Poi la modifica viene eseguita in background con lo stesso codice di [`ccam group`](/features/share): prima viene salvata la configurazione, poi i dati di ogni account vengono spostati in ordine, e gli originali vanno in `.ccam-backup/`. Una modifica che sposterebbe i dati di un account in esecuzione viene rifiutata indicando il nome dell'account.

## Aggiungere un account

`a` apre la procedura guidata che `ccam add` usa anche in un terminale:

1. Il nome dell'account.
2. Come accede: un abbonamento Claude subito (nel browser), una [chiave API](/features/api), oppure un abbonamento Claude più tardi.
3. Per una chiave API: da dove ccam legge la chiave a ogni avvio (`env:NAME` o `command:COMMAND`), l'endpoint, come viene inviata la chiave e il modello.
4. I gruppi a cui unirsi, con `default` selezionato; prima che esistano gruppi, cosa condivide il gruppo `default`.
5. Un comando scorciatoia, se ne vuoi uno.
6. Un riepilogo da confermare.

`Esc` torna indietro di un passo e, al primo passo, chiude la procedura guidata.

## Colori

I colori seguono i ruoli dell'interfaccia da terminale di Claude Code, così ccam si integra bene accanto a esso; i nomi dei gruppi ricevono un colore ricavato dal loro nome. ccam chiede al terminale il suo colore di sfondo e sceglie un tema chiaro o scuro. Puoi impostarne uno tu in `config.toml`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

I temi `-daltonized` sono pensati per il daltonismo rosso-verde e i temi `-ansi` usano i 16 colori del terminale stesso. `NO_COLOR` disattiva i colori.

## Sicurezza e casi limite

- La dashboard chiama le stesse funzioni dei comandi, quindi valgono tutte le regole: piani prima di spostare i dati, rifiuti per gli account in esecuzione, backup, e token di login che non lasciano mai un account.
- La dashboard non ha un dry run: `ccam --dry-run` senza un nome stampa come partirebbe l'account di questa cartella, e `o` su un account offre un dry run di quel solo avvio.
- Il terminale viene ripristinato prima che Claude Code si avvii o che si apra il login.
- Uno schermo troppo piccolo per una finestra di dialogo toglie prima le righe della scheda; le righe vengono tagliate con `…` invece di andare a capo.

## Pagine correlate

- [Condivisione con i gruppi](/features/share)
- [Avviare Claude Code](/features/launch)
- [Gestire gli account](/features/accounts)
- [Limiti di utilizzo](/features/usage)
