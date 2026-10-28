---
title: Installare ccam
metaTitle: Installare ccam su Linux, macOS e Windows
description: Installa ccam su Linux, macOS o Windows con un solo comando, un pacchetto .deb o go install, e aggiornalo con ccam update. Serve Claude Code già installato.
updated: 2026-10-28
order: 1
translation: machine
sourceHash: 193df30c9d9a
---

ccam avvia Claude Code per te, quindi installa prima [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) e verifica che `claude --version` funzioni.

Esegui questo comando in un terminale su Linux o macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Lo script scarica la release per il tuo sistema, ne verifica il checksum SHA-256 e mette `ccam` in `~/.local/bin`, dove ccam crea anche le scorciatoie che scegli in seguito, come `ccwork`. Se quella cartella non è nel tuo `PATH`, lo script stampa la riga da aggiungere a `~/.bashrc` o `~/.zshrc`.

Su Windows, esegui questo in PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Scarica `ccam_windows_<arch>.zip`, ne verifica il checksum SHA-256 e mette `ccam.exe` in `%USERPROFILE%\.local\bin`, la cartella in cui l'installer di Claude Code mette `claude.exe`. Se quella cartella non è ancora nel `PATH` del tuo utente, lo script la aggiunge e lascia invariate voci come `%USERPROFILE%\bin`. Dopo, apri un nuovo terminale.

## Sistemi supportati

- Linux su x86-64 e ARM64: Ubuntu, Debian, Fedora, Arch, Alpine e altre distribuzioni. Il binario è statico e non ha dipendenze.
- macOS su Apple silicon e Intel.
- Windows 10 e 11 su x86-64 e ARM64, in PowerShell, Windows Terminal o Git Bash:
  - Le cartelle condivise funzionano senza la Modalità sviluppatore.
  - I file condivisi (`settings.json`, `CLAUDE.md`, `history.jsonl`) richiedono la Modalità sviluppatore (Impostazioni > Sistema > Per sviluppatori), perché Windows consente di creare link a file solo in quella modalità o a un amministratore; vedi [condivisione su Windows](/features/share#on-windows).
  - Senza la Modalità sviluppatore, il primo `ccam add` crea il gruppo `default` con i soli server MCP ed estensioni, e spiega come aggiungere il resto una volta attivata la Modalità sviluppatore.
- All'interno di WSL, ccam funziona come su qualsiasi altro sistema Linux.

## Altri modi per installare

**Pacchetto per Debian e Ubuntu.** Scarica `ccam_<version>_linux_amd64.deb` dalla [pagina delle release](https://github.com/jhyunwoo/ccam/releases) e installalo:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

Il pacchetto installa `/usr/bin/ccam`.

**Dai sorgenti.** Con Go 1.26 o successivo:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

Se installi un tag di release (`@latest` o `@vX.Y.Z`), ccam riporta quella versione, e `ccam update` e la proposta di aggiornamento della dashboard funzionano. Una build di `@main`, o di un checkout locale con commit successivi a un tag, riporta `dev`, e `ccam update` la rifiuta a meno che tu non passi `--version`.

**Una versione specifica.** Imposta `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## Completamento della shell

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

Il completamento copre tutti i comandi e le opzioni di `ccam`, compresi i nomi degli account.

## Plugin per Claude Code

Il plugin di ccam insegna a Claude Code stesso a usare ccam: a controllare quale account ha margine, ad avviare worker `claude -p` in parallelo su account diversi con `ccam pick` e a spiegare una skill o un server MCP mancante con `ccam ext`. Installalo dentro Claude Code:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

Aggiunge:

- una skill che suggerisce i comandi che fanno il login o spostano file, ma non li esegue mai;
- i comandi `/ccam:usage`, `/ccam:ext` e `/ccam:switch`;
- due hook ([dettagli](/features/next#inside-claude-code)). Uno mostra una notifica nel terminale quando una sessione si ferma a un limite di utilizzo; su Linux e macOS, con [`auto_switch`](/features/next#auto-switch) attivo, ccam sposta poi la sessione su un altro account con margine disponibile. L'altro permette a `/ccam:switch [account]` di passare la conversazione a un altro account.

Gli hook eseguono `ccam`, quindi deve trovarsi nel `PATH` che vede Claude Code.

## Verificare l'installazione

```sh
ccam version
ccam doctor
```

`ccam doctor` controlla che Claude Code venga trovato e che ogni cartella di account, link e login sia in ordine. Prosegui con la [guida rapida](/docs/quick-start).

## Aggiornamento

```sh
ccam update --check
ccam update
```

`ccam update --check` ti dice solo se esiste una release più recente. `ccam update` scarica la release per il tuo sistema, ne verifica il checksum SHA-256 confrontandolo con il file `checksums.txt` della release e poi sostituisce il file `ccam` in un solo passaggio. Se il checksum non corrisponde, non viene modificato nulla. Le scorciatoie come `ccwork` sono link a ccam, quindi usano subito la nuova versione.

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

Anche la dashboard (`ccam` senza argomenti, in un terminale) cerca una nuova release ogni volta che si apre e, se ce n'è una, la propone: premi `U` e conferma. Nessun altro comando cerca aggiornamenti, e ccam non si aggiorna mai senza chiedere. Per installare una release specifica, anche una precedente, passa `--version`. `ccam --dry-run update` mostra cosa verrebbe scaricato e sostituito, senza modificare nulla.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Anche eseguire di nuovo lo script di installazione aggiorna la copia in `~/.local/bin`. Se hai compilato ccam con `go install` da un tag di release, puoi eseguire `ccam update` o di nuovo lo stesso comando `go install`.

**Windows.** `ccam update` scarica il `.zip` per Windows. Su Windows un programma in esecuzione non può essere sostituito, quindi il vecchio binario viene rinominato in `ccam.exe.old`, che l'aggiornamento successivo elimina; mentre una sessione è ancora in esecuzione da un `ccam.exe.old` precedente, il nome è `ccam.exe.old-<time>`. Subito dopo vengono rinnovate le scorciatoie che sono hard link. Eseguire di nuovo `install.ps1` funziona allo stesso modo, anche mentre una sessione di Claude Code è ancora in esecuzione dal vecchio binario, e stampa `OK 2 shortcuts now run the new ccam` quando rinnova le scorciatoie che sono hard link.

**Pacchetto per Debian e Ubuntu.** `ccam update` non sostituisce un `/usr/bin/ccam` installato con apt e al suo posto stampa questi passaggi. Scarica il nuovo `ccam_<version>_linux_amd64.deb` dalla [pagina delle release](https://github.com/jhyunwoo/ccam/releases) e installalo sopra quello vecchio:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Disinstallazione

`ccam uninstall` rimuove ccam e tutto ciò che ha creato. Prima elenca cosa verrà rimosso, e chiede conferma:

```sh
ccam --dry-run uninstall   # only show what would be removed
ccam uninstall
```

- Ogni account viene rimosso come fa `ccam remove`. Un account nelle cartelle di ccam viene prima disconnesso e poi eliminato insieme alla sua cartella, compresi la cronologia delle conversazioni e le impostazioni.
- Una cartella `~/.claude` che hai [importato](/docs/commands) resta dove si trova per il semplice `claude`. Se fa parte di un gruppo, riceve una copia di tutto ciò che condivideva. Nel suo `settings.json`, la riga di stato di ccam lascia il posto alla riga di stato che ccam avvolgeva, se ce n'era una. Se ha il plugin ccam, ccam stampa il comando `/plugin uninstall` da eseguire in Claude Code.
- I comandi scorciatoia, `~/.config/ccam`, `~/.local/share/ccam` (con il cestino di ccam) e `~/.local/state/ccam` vengono eliminati, poi il comando `ccam` stesso.

ccam si rifiuta finché un account è in esecuzione; chiudi prima le sue sessioni di Claude Code. Un `/usr/bin/ccam` installato da apt resta; rimuovilo dopo con `sudo apt remove ccam`. Su Windows un programma in esecuzione non può eliminare se stesso, quindi `ccam uninstall` termina stampando il comando `Remove-Item` per `ccam.exe`.
