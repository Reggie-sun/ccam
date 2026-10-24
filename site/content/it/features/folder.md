---
title: Account per cartella
description: ccam folder fa usare un account a una cartella e a tutto ciò che contiene quando avvii Claude Code senza nome di account, come work per il codice dell'azienda.
updated: 2026-10-23
order: 11
translation: machine
sourceHash: "48de98d42b99"
---

`ccam folder work ~/code/company` fa sì che ogni avvio senza nome di account dentro `~/code/company`, e in qualsiasi cartella al suo interno, avvii Claude Code con `work`. Fuori da lì, l'account predefinito si avvia come prima. Indicare un account, con `ccam personal` o con una scorciatoia come `ccpersonal`, avvia sempre quell'account.

## Utilizzo

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

Questi avvii usano l'account della cartella:

- `ccam -p "…"` e qualsiasi altro avvio il cui primo argomento è un'opzione per Claude Code
- `ccam run` senza nome, e `ccam --dry-run`
- `ccam --fallback -p "…"` per il suo primo tentativo
- la dashboard, `ccam` da solo, che parte con il cursore sull'account della cartella, anche quando sei arrivato nella cartella tramite un link simbolico

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `<account>` | Nome o numero dell'account da usare | |
| `[folder]` | La cartella, relativa o assoluta | la cartella corrente |
| `--remove` | Rimuove la voce esattamente di quella cartella | disattivata |
| `--dry-run` | Rispettata: stampa la modifica senza salvarla | disattivata |

## Come funziona

Ogni cartella è una voce `[[folder]]` in `~/.config/ccam/config.toml`, il file che ogni avvio legge comunque:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` risolve la cartella in un percorso assoluto, con i link simbolici risolti, e lo salva. Impostare di nuovo una cartella ne sostituisce l'account.
2. Un avvio senza nome chiama `launcher.HereProfile`. Se non ci sono voci, prende subito l'account predefinito.
3. Altrimenti confronta la directory di lavoro con ogni `path` come stringhe: la cartella stessa o una al suo interno, mai una cartella che inizia soltanto con le stesse lettere (`/code/company` non copre `/code/companyx`). Vince la voce più profonda, quindi `~/code` può essere `personal` mentre `~/code/company` è `work`.
4. Solo se la directory di lavoro così com'è non corrisponde a nulla, ccam ne risolve i link simbolici e confronta ancora una volta.
5. Su Windows e macOS il confronto ignora maiuscole e minuscole, come di solito fanno i loro file system.

Nessun file oltre a `config.toml` viene letto, quindi un avvio dentro una cartella non costa un tempo misurabile (4,0 ms con o senza voci, `hyperfine`, 300 esecuzioni).

## Sicurezza e casi limite

- `ccam rename` porta le voci delle cartelle sul nuovo nome, e `ccam remove` elimina le voci dell'account rimosso.
- Una voce il cui account non esiste più blocca un avvio senza nome in quella cartella: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` segnala queste voci e quelle la cui cartella non c'è più.
- `ccam folder --remove` dentro una cartella che riceve il suo account da una cartella superiore indica quella cartella invece di rimuovere qualcosa.
- `ccam next` sceglie per rotazione o per utilizzo e non guarda le cartelle.

## Pagine correlate

- [Avviare Claude Code](/features/launch)
- [Comandi scorciatoia](/features/shortcuts)
- [Gestire gli account](/features/accounts)
- [Riferimento dei comandi](/docs/commands)
