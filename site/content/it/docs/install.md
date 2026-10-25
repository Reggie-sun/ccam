---
title: Installare ccam
metaTitle: Installare ccam su Linux, macOS e Windows
description: Installa ccam su Linux, macOS o Windows con un solo comando, un pacchetto .deb o dai sorgenti, e aggiornalo con ccam update. Serve Claude Code già installato.
updated: 2026-10-25
order: 1
translation: machine
sourceHash: "a7f9dac24f79"
---

Esegui questo comando in un terminale su Linux o macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Lo script scarica la release per il tuo sistema, ne verifica il checksum SHA-256 e mette `ccam` in `~/.local/bin`, dove ccam crea anche le scorciatoie che scegli in seguito, come `ccwork`. Se quella cartella non è nel tuo `PATH`, lo script stampa la riga da aggiungere a `~/.bashrc` o `~/.zshrc`.

Su Windows, esegui questo in PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

Scarica `ccam_windows_<arch>.zip`, ne verifica il checksum SHA-256, mette `ccam.exe` in `%USERPROFILE%\.local\bin`, la cartella in cui l'installer di Claude Code mette `claude.exe`, e aggiunge quella cartella al `PATH` del tuo utente se non c'è ancora, lasciando invariate voci come `%USERPROFILE%\bin`. Dopo, apri un nuovo terminale. Eseguirlo di nuovo sostituisce ccam, anche mentre una sessione di Claude Code gira ancora dal vecchio binario (che viene rinominato in `ccam.exe.old`, o `ccam.exe.old-<number>` quando uno più vecchio è ancora in uso), e rinnova le scorciatoie che sono hard link: `OK 2 shortcuts now run the new ccam`.

ccam avvia Claude Code per te, quindi installa prima [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) e verifica che `claude --version` funzioni.

## Sistemi supportati

- Linux su x86-64 e ARM64: Ubuntu, Debian, Fedora, Arch, Alpine e altre distribuzioni. Il binario è statico e non ha dipendenze.
- macOS su Apple silicon e Intel.
- Windows 10 e 11 su x86-64 e ARM64, in PowerShell, Windows Terminal o Git Bash. Le cartelle condivise funzionano per ogni utente. I file condivisi (`settings.json`, `CLAUDE.md`, `history.jsonl`) richiedono la Modalità sviluppatore (Impostazioni > Sistema > Per sviluppatori), perché Windows consente di creare link a file solo in quella modalità o a un amministratore; vedi [condivisione su Windows](/features/share#on-windows). Senza di essa, il primo `ccam add` crea il gruppo `default` con i soli server MCP ed estensioni, e spiega come aggiungere il resto una volta attivata la Modalità sviluppatore.
- All'interno di WSL, ccam funziona come su qualsiasi sistema Linux.

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

Una build di un tag di release (`@latest` o `@vX.Y.Z`) riporta quella versione, quindi `ccam update` funziona; una build di `@main` o di un checkout locale successivo a un tag riporta `dev`.

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

Aggiunge una skill, i comandi `/ccam:usage`, `/ccam:ext` e `/ccam:switch`, e due hook: quando una sessione si ferma a un limite di utilizzo ricevi una notifica nel terminale (e su Linux e macOS, con [`auto_switch`](/features/next#auto-switch) attivo, ccam sposta la sessione su un altro account con margine disponibile), e `/ccam:switch [account]` passa la conversazione a un altro account ([dettagli](/features/next#inside-claude-code)). Gli hook eseguono `ccam`, quindi deve trovarsi nel `PATH` che vede Claude Code. La skill non esegue mai comandi che fanno il login o spostano file; li suggerisce invece.

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

Anche la dashboard (`ccam` in un terminale) cerca una nuova release ogni volta che si apre e, se ce n'è una, la propone: premi `U` e conferma. Nessun altro comando la cerca, e ccam non si aggiorna mai senza chiedere. Per installare una release specifica, anche una precedente, passa `--version`. `ccam --dry-run update` mostra cosa verrebbe scaricato e sostituito, senza modificare nulla.

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

Anche eseguire di nuovo lo script di installazione aggiorna la copia in `~/.local/bin`. Se hai compilato ccam con `go install` da un tag di release, funziona anche `ccam update`, oppure esegui di nuovo lo stesso comando.

**Windows.** `ccam update` scarica il `.zip` per Windows. Su Windows un programma in esecuzione non può essere sostituito, quindi il vecchio binario viene rinominato in `ccam.exe.old`, che l'aggiornamento successivo elimina; mentre una sessione gira ancora da un `ccam.exe.old` precedente, il nome è `ccam.exe.old-<time>`. Subito dopo vengono rinnovate le scorciatoie che sono hard link.

**Pacchetto per Debian e Ubuntu.** `ccam update` non sostituisce un `/usr/bin/ccam` installato con apt e al suo posto stampa questi passaggi. Scarica il nuovo `ccam_<version>_linux_amd64.deb` dalla [pagina delle release](https://github.com/jhyunwoo/ccam/releases) e installalo sopra quello vecchio:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## Disinstallazione

Rimuovi prima le scorciatoie degli account, finché ccam è ancora installato. `ccam list` le mostra nella colonna Shortcut; rimuovile una per una con `ccam shortcut <account> --remove`, oppure elimina tu quei link da `~/.local/bin`. Poi rimuovi il comando e, se non ti servono più gli account, le cartelle dei dati:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Su Windows, in PowerShell:

```powershell
ccam shortcut work --remove   # once for each account that has a shortcut
Remove-Item ~\.local\bin\ccam.exe
Remove-Item -Recurse ~\.config\ccam, ~\.local\share\ccam, ~\.local\state\ccam
```

Eliminare `~/.local/share/ccam` elimina anche i login e la cronologia delle conversazioni salvati lì. Una cartella `~/.claude` che hai [importato](/docs/commands) resta dove si trova, ma finché fa parte di un gruppo, i file che condivide stanno in `~/.local/share/ccam/groups/`. Esegui prima `ccam remove <name>` per quell'account (o `ccam group leave`): lascia i suoi gruppi con una copia di tutto ciò che condivideva.
