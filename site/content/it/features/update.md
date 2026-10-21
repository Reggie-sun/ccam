---
title: Aggiornare ccam
description: ccam update sostituisce ccam con l'ultima release su GitHub dopo un controllo SHA-256; --check segnala soltanto e --version installa un tag specifico.
updated: 2026-10-21
order: 10
translation: machine
sourceHash: a92dbf37399f
---

`ccam update` scarica una release di ccam da GitHub, la verifica con il file `checksums.txt` della release e sostituisce il binario in esecuzione con un'unica rinomina. Usalo se hai installato ccam con lo script di installazione o da un archivio. Anche la [dashboard](/features/dashboard) chiede a GitHub l'ultima release ogni volta che si apre e la propone con `U`, che esegue gli stessi passaggi dopo aver chiesto conferma. Nessun altro comando contatta GitHub, quindi funzionano allo stesso modo anche offline.

## Utilizzo

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, also an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

Quando non esiste una release più recente, ccam stampa `ccam 0.4.0 is up to date.` e non modifica nulla. Con `--version` e la versione che stai già usando, stampa `ccam 0.4.0 is already installed.`

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `--check` | Segnala solo se esiste una release più recente | disattivata |
| `--version <tag>` | Installa questo tag di release, come `v0.3.0`, più recente o più vecchio di quello attuale | l'ultima |
| `--dry-run` | Rispettata: cerca l'ultimo tag (a meno che `--version` non ne indichi uno), esegue gli stessi controlli sul binario e stampa i percorsi di download, checksum e sostituzione senza scaricare nulla | disattivata |

`--check` e `--version` non si possono usare insieme. Il tag deve avere la forma `vMAJOR.MINOR.PATCH`, con un eventuale suffisso di pre-release o di build; la `v` iniziale si può omettere. Qualsiasi altro valore fallisce, per esempio `--version takes a release tag such as v0.3.0, not "1.2"`.

La variabile d'ambiente `CCAM_RELEASES_URL` sostituisce `https://github.com/jhyunwoo/ccam/releases`, per un mirror o per i test. Deve avere la stessa struttura: `<url>/latest` reindirizza a `<url>/tag/<tag>`, e i file della release si trovano in `<url>/download/<tag>/<asset>`.

## Come funziona

Il codice si trova in `internal/update`, importato solo da `ccam update`. Il percorso di avvio non lo chiama mai.

1. `cli.updateCmd` trova il binario in esecuzione con `os.Executable` e risolve i link simbolici con `filepath.EvalSymlinks`. È quel percorso reale a essere sostituito.
2. Senza `--version`, `Updater.Latest` richiede `<releases>/latest` senza seguire il reindirizzamento e legge il tag dall'intestazione `Location` (`…/tag/v0.5.0`). In questo modo usa la pagina delle release invece dell'API di GitHub, quindi il limite dell'API di 60 richieste all'ora non si applica. La ricerca ha 10 secondi.
3. `parseVersion` verifica il tag con uno schema rigoroso. Un tag finisce negli URL di download e in un comando di shell stampato, quindi non passa nulla che contenga `/`, `%`, spazi o caratteri della shell. `compareVersions` lo confronta con la versione in esecuzione; una pre-release viene ordinata prima della sua release.
4. `managed()` e `writable()` decidono se ccam può sostituire il binario (vedi sotto).
5. `Updater.Fetch` scarica contemporaneamente `ccam_<os>_<arch>.tar.gz` (su Windows `ccam_windows_<arch>.zip`) e `checksums.txt`, con 2 minuti per entrambi. Ogni download è limitato a 256 MB.
6. Lo SHA-256 dell'archivio deve corrispondere alla sua riga in `checksums.txt` prima che succeda qualsiasi altra cosa. Una riga mancante o un checksum diverso interrompe l'aggiornamento con `nothing was changed`.
7. `extract()` legge il tar e prende il primo file normale chiamato `ccam` (`extractZip` prende `ccam.exe` dallo zip). Cartelle, link e altre voci vengono ignorati, e anche il binario è limitato a 256 MB.
8. `Updater.Install` scrive il nuovo binario con `fsutil.AtomicWrite`: un file temporaneo `.ccam.tmp*` accanto al binario, con i bit di permesso del vecchio file, scritto su disco e rinominato al posto di quello vecchio.

L'intera sequenza, dalla ricerca del tag alla rinomina:

<!-- motion:update-pipeline -->

L'archivio e `checksums.txt` vengono scaricati in parallelo, e non viene scritto nulla finché il checksum non corrisponde. La rinomina è atomica, quindi un crash o un disco pieno lasciano il vecchio binario o quello nuovo, mai metà di uno dei due. I comandi scorciatoia come `~/.local/bin/ccwork` sono link simbolici al binario di ccam, quindi eseguono la nuova versione senza alcuna modifica.

Su Windows il file di un programma in esecuzione non può essere sostituito né eliminato, solo rinominato. `Updater.Install` scrive il nuovo binario in `ccam.exe.new`, rinomina il `ccam.exe` in esecuzione in `ccam.exe.old` e rinomina il nuovo file in `ccam.exe`. Se l'ultimo passaggio fallisce, il vecchio binario riprende il suo nome. L'aggiornamento successivo elimina `ccam.exe.old`. Le scorciatoie che sono hard link invece di link simbolici ([Comandi scorciatoia](/features/shortcuts)) vengono rinnovate subito dopo, così eseguono anch'esse la nuova versione.

### Installazioni gestite da un pacchetto

Un ccam installato dal `.deb` appartiene ad apt, e sostituirlo all'insaputa di apt lascerebbe errato il database dei pacchetti. Prima di scaricare, ccam controlla a chi appartiene il binario:

<!-- motion:update-managed -->

1. Un binario in `/usr/bin`, `/bin`, `/usr/sbin` o `/sbin` è considerato gestito da un pacchetto.
2. Altrove, se `dpkg` è installato e `dpkg -S <path>` riesce, anche in questo caso è considerato gestito da un pacchetto.
3. Altrimenti ccam crea e rimuove un file temporaneo `.ccam-update-*` accanto al binario per assicurarsi di poter scrivere lì.

Per un binario gestito da un pacchetto ccam si ferma e stampa i passaggi per apt:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

Sui sistemi diversi da Linux viene stampata solo la riga dello script di installazione: `irm https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.ps1 | iex` su Windows. Quando la cartella non è scrivibile, ccam suggerisce di eseguire di nuovo lo script di installazione, che installa ccam in `~/.local/bin`.

## Sicurezza e casi limite

- Una build di sviluppo indica come versione `dev` e non ha una release con cui confrontarsi, quindi `ccam update` e `ccam update --check` la rifiutano. `ccam update --version vX.Y.Z` installa una release al suo posto.
- Un download fallito, un file mancante per la tua piattaforma o un checksum che non corrisponde lasciano il binario installato com'era.
- Una release che non esiste fallisce con `there is no release <tag>`. Quando GitHub rifiuta con HTTP 403 o 429, ccam indica che probabilmente si tratta di un limite di richieste e ti chiede di riprovare più tardi.
- `--check` si limita a leggere. Non controlla se il binario è gestito da un pacchetto o scrivibile.
- `--version` con un tag più vecchio è il modo per tornare a una versione precedente.

## Pagine correlate

- [Installazione](/docs/install)
- [Comandi](/docs/commands)
- [Domande frequenti](/docs/faq)
