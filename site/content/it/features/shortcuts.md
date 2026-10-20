---
title: Comandi scorciatoia
description: Assegna a un account un comando tutto suo, come ccwork, con ccam shortcut, e ricrea i link in ~/.local/bin con ccam shims install.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 939256e7a67e
---

Una scorciatoia è un comando a cui dai tu il nome e che avvia Claude Code con un determinato account: `ccwork --resume` fa ciò che fa `ccam work --resume`. È un link a ccam in `~/.local/bin`, quindi funziona in bash, zsh, negli script e in qualsiasi altra cosa che esegua programmi dal `PATH`. Le scorciatoie sono facoltative; creane una per gli account che avvii spesso.

## Utilizzo

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing links of all accounts
ccwork -p "summarise" --output-format json
```

ccam conferma ogni modifica:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

Un account ha al massimo una scorciatoia, salvata come `shortcut` in `~/.config/ccam/config.toml`. La procedura guidata di `ccam add` chiede `Create a shortcut command for <name>? [y/N]` e poi `Shortcut name [cc<name>]:`. Fuori dalla procedura guidata, `ccam add` ne crea una solo con `--shortcut`. `ccam list` mostra la scorciatoia di ogni account, oppure `-`.

Una scorciatoia passa ogni argomento a Claude Code, come fa `ccam <name>` dopo il nome. Questo vale anche per `--dry-run` e `--help`: per vedere in anteprima un avvio, usa `ccam --dry-run work`.

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `--remove` | Elimina la scorciatoia dell'account e il suo link | disattivata |
| `-h`, `--help` | Aiuto per il comando | |
| `--dry-run` | `ccam shortcut` la rispetta: controlla il nome e stampa `Would set the shortcut of work to ccwork.` o `Would remove ccwork, the shortcut of work.` senza salvare né creare link | disattivata |

`ccam shims install` non accetta argomenti e rifiuta `--dry-run` con `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``.

Regole per i nomi, controllate in questo ordine:

- da 1 a 64 caratteri tra `a-z`, `0-9`, `.`, `-` e `_`, iniziando con una lettera o una cifra;
- non `ccam`, `claude` o un comando di ccam come `run` o `usage`;
- non la scorciatoia di un altro account;
- non un file in `~/.local/bin` che non sia già un link a ccam;
- non il nome di un altro eseguibile nel tuo `PATH`.

Indicare insieme un comando e `--remove` si ferma con `give a command or --remove, not both`. Impostare il nome che un account ha già stampa `work: ccwork (unchanged)` e ricrea il suo link se manca.

## Come funziona

Impostare un nome esegue `checkShortcut`, che chiama `shim.Check` con la cartella bin e il percorso reale del binario di ccam:

<!-- motion:shortcuts-check -->

Il nome scende lungo i controlli in ordine. Accanto a ogni controllo c'è un nome che rifiuta: `Work` non rispetta lo schema, `run` è un comando di ccam, `ccside` appartiene a un altro account, `ccfile` è un file esistente in `~/.local/bin` e `git` nasconderebbe `/usr/bin/git`. Solo un nome che supera tutti i controlli viene salvato e collegato.

1. `shim.Check` confronta il nome con `^[a-z0-9][a-z0-9._-]{0,63}$`, poi con `config.Reserved`, poi con `cfg.FindShortcut`. Chiama `os.Lstat` su `~/.local/bin/<name>`, e per ogni altra cartella del `PATH` chiama `os.Stat` e rifiuta un file eseguibile che non sia un link a ccam.
2. ccam salva `config.toml` con `fsutil.AtomicWrite` (file temporaneo, `fsync`, rename).
3. `shim.Remove` elimina il link della vecchia scorciatoia, ma solo se è un link a ccam.
4. `shim.Install` crea il link di ogni scorciatoia configurata che manca, con `os.Symlink` verso il binario di ccam, individuato con `os.Executable` e `filepath.EvalSymlinks`. Non rimuove mai nulla.

Un link è considerato di ccam quando `os.Readlink` restituisce il percorso di ccam, oppure quando il link e il percorso di ccam portano allo stesso file. Ogni rimozione e ogni controllo di sovrascrittura passa per questa verifica, quindi ccam non tocca file che non ha creato.

Quando esegui una scorciatoia, il nome che hai digitato è l'unica cosa di cui ccam ha bisogno:

<!-- motion:shortcuts-link -->

La shell trova `ccwork` nel `PATH`, il kernel segue il link e avvia il binario di ccam, e `argv[0]` resta `ccwork`. In `app.Main`, qualsiasi nome diverso da `ccam` carica una volta `config.toml` e chiama `shim.Lookup`, che trova l'account il cui `shortcut` è uguale a `filepath.Base(argv[0])`. `launcher.LaunchWith` segue poi gli stessi passaggi di [`ccam <name>`](/features/launch) con la configurazione già caricata, e passa `argv[1:]` a `claude` senza modificarli.

Perché un link simbolico e non un alias della shell o uno script wrapper: un link simbolico funziona in ogni shell e negli script, non aggiunge l'avvio di un'altra shell (circa da 1 a 2 ms) e l'avvio legge comunque un solo file. `scripts/perf-check.sh` misura una scorciatoia con `hyperfine` rispetto allo stesso limite di 10 ms (p95) di `ccam <name>`.

### Su Windows

Una scorciatoia è `ccwork.exe` in `%USERPROFILE%\.local\bin`, e `shim.Name` toglie il `.exe` (maiuscolo o minuscolo) prima della ricerca, quindi `ccwork`, `ccwork.exe` e `CCWORK.EXE` avviano tutti `work`. ccam crea prima un link simbolico, che Windows consente in Modalità sviluppatore o a un amministratore. Senza quel diritto crea un hard link: un secondo nome per lo stesso file `ccam.exe`, che non richiede privilegi. Un hard link conta come di ccam quando è lo stesso file di `ccam.exe`. Poiché `ccam update` mette un nuovo file al posto di `ccam.exe`, rinnova poi ogni scorciatoia creata come hard link con `shim.Relink` e stampa `✓ shortcut ccwork renewed`. Il controllo del nome cerca `<name>` con ogni estensione di `PATHEXT` (`.exe`, `.cmd` e così via) nelle altre cartelle del `PATH`.

## Sicurezza e casi limite

- **Un link che nessun account rivendica.** Se un link a ccam in `~/.local/bin` ha un nome che nessun account usa, per esempio dopo che hai modificato `shortcut` in `config.toml` a mano, ccam non ripiega sull'account predefinito. Si ferma con `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) ``. Imposta di nuovo il nome con `ccam shortcut`, oppure elimina il link.
- **File esistenti.** ccam non sostituisce mai un file che non ha creato. `ccam shims install` si ferma con `<path> already exists and is not a ccam link, so it was left alone`.
- **Un binario di ccam spostato.** I vecchi link puntano allora a un percorso che non esiste più, e ccam non li considera suoi. Eliminali, poi esegui `ccam shims install`.
- **Rinomina e rimozione.** `ccam rename` mantiene il nome della scorciatoia, e il link continua a funzionare, perché punta a ccam e ccam cerca il nome in `config.toml`. `ccam remove` elimina il link dell'account. `ccam shortcut <account> <new>` elimina il vecchio link prima di creare quello nuovo.
- **PATH.** Quando ccam crea un link e `~/.local/bin` non è nel tuo `PATH`, stampa `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` seguito dalla riga `export PATH=...` da aggiungere. Su Windows stampa invece la riga PowerShell che aggiunge la cartella al `PATH` del tuo utente.
- **Controllare in seguito.** `ccam doctor` segnala `missing shortcuts in /home/me/.local/bin: ccwork` e avvisa quando la cartella non è nel `PATH`. `ccam doctor --fix` crea i link mancanti.

## Pagine correlate

- [Avviare Claude Code](/features/launch)
- [Gestire gli account](/features/accounts)
- [Riferimento dei comandi](/docs/commands#shortcuts)
- [Installare ccam](/docs/install)
