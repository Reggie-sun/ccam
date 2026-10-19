---
title: Installare ccam
description: Installa ccam su Linux o macOS con un solo comando, da un pacchetto .deb o dai sorgenti. Claude Code deve essere già installato.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: fc63978e8f9e
---

Esegui questo comando in un terminale su Linux o macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

Lo script scarica la release per il tuo sistema, ne verifica il checksum SHA-256 e mette `ccam` in `~/.local/bin`, dove in seguito `ccam add` crea le scorciatoie `claude1`, `claude2` e così via. Se quella cartella non è nel tuo `PATH`, lo script stampa la riga da aggiungere a `~/.bashrc` o `~/.zshrc`.

ccam avvia Claude Code per te, quindi installa prima [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) e verifica che `claude --version` funzioni.

## Sistemi supportati

- Linux su x86-64 e ARM64: Ubuntu, Debian, Fedora, Arch, Alpine e altre distribuzioni. Il binario è statico e non ha dipendenze.
- macOS su Apple silicon e Intel.
- Windows non è supportato direttamente. All'interno di WSL, ccam funziona come su qualsiasi sistema Linux.

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

**Una versione specifica.** Imposta `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## Completamento della shell

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

Il completamento copre tutti i comandi e le opzioni di `ccam`, compresi i nomi degli account.

## Verificare l'installazione

```sh
ccam version
ccam doctor
```

`ccam doctor` controlla che Claude Code venga trovato e che ogni cartella di account, link e login sia in ordine. Prosegui con la [guida rapida](/docs/quick-start).

## Disinstallazione

Rimuovi il comando e le scorciatoie degli account, poi le cartelle dei dati se non ti servono più gli account:

```sh
rm ~/.local/bin/ccam ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

Eliminare `~/.local/share/ccam` elimina anche i login e la cronologia delle conversazioni salvati lì. Una cartella `~/.claude` che hai [importato](/docs/commands) resta dove si trova.
