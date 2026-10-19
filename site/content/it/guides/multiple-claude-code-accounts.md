---
title: Come usare due account Claude Code su un solo computer
description: Usa più account Claude Code in parallelo con CLAUDE_CONFIG_DIR, a mano con gli alias della shell o in automatico con ccam.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 6a4b7d961e10
---

Per usare più account Claude Code sullo stesso computer, assegna a ogni account la propria cartella di configurazione. Claude Code legge la cartella dalla variabile d'ambiente `CLAUDE_CONFIG_DIR`, e ogni cartella contiene il proprio login, le proprie impostazioni e la propria cronologia. Puoi configurarlo a mano con un alias della shell oppure lasciarlo fare a ccam.

## A mano, con un alias

È l'approccio descritto nella documentazione di Claude Code. Aggiungi a `~/.bashrc` o `~/.zshrc` una riga per ogni account in più:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Apri un nuovo terminale ed esegui `claude-work`. Claude Code ti chiede di effettuare il login, perché la cartella è nuova. Da quel momento `claude` è il tuo primo account e `claude-work` il secondo.

Con due account funziona bene. Cose da sapere:

- Ogni cartella parte vuota: nessuna impostazione, nessun `CLAUDE.md`, nessun server MCP, nessun agente o skill. Copi o colleghi ciò che ti serve, e tieni allineate le copie da solo.
- Gli alias esistono solo nelle shell interattive. Script, Makefile e altri programmi non li vedono.
- I server MCP si trovano in `.claude.json`, accanto al login, quindi non si possono semplicemente collegare tra cartelle.
- Le conversazioni restano nell'account in cui sono iniziate, quindi non puoi continuarne una in un altro account.
- Non c'è nulla che mostri l'utilizzo di tutti gli account in una volta.

## Con ccam

ccam automatizza la stessa idea e aggiunge le parti mancanti:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

Ora `ccam personal` e `ccam work` avviano i due account, e lo stesso fa `ccwork`, la scorciatoia che hai dato all'account work. Poiché sono veri programmi e non alias, funzionano anche negli script. Impostazioni, `CLAUDE.md`, server MCP ed estensioni sono condivisi fin dall'inizio, e puoi attivare la cronologia delle conversazioni condivisa per continuare una conversazione sull'altro account con `ccam next --continue`.

Se usi già `~/.claude`, esegui `ccam import` invece di aggiungere di nuovo quell'account, oppure `ccam migrate` per copiarne le impostazioni, i server MCP e le estensioni nelle impostazioni condivise.

## Quale scegliere

Per due account senza nulla di condiviso basta un alias. Quando vuoi impostazioni e server MCP condivisi, più account, l'utilizzo a colpo d'occhio o la possibilità di continuare le conversazioni, ccam ti risparmia la gestione manuale. Il [confronto](/compare/aliases) mette le differenze una accanto all'altra.
