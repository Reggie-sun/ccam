---
title: ccam a confronto con gli alias della shell
description: Alias della shell con CLAUDE_CONFIG_DIR e ccam a confronto, fianco a fianco, per usare più account Claude Code.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: b9fbcb88f111
---

Un alias della shell per ogni account è il modo più semplice per usare più account Claude Code. ccam fa la stessa cosa dietro le quinte e aggiunge condivisione, cambio di account e monitoraggio dell'utilizzo.

## La configurazione con gli alias

Aggiungi una riga per ogni account in più a `~/.bashrc` o `~/.zshrc`, apri un nuovo terminale ed esegui l'alias una volta per fare il login:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

`claude` resta il tuo primo account e `claude-work` avvia il secondo. È il metodo [descritto nella documentazione di Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts).

## La configurazione con ccam

Installa ccam e aggiungi l'account, qui con un comando scorciatoia tutto suo:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add work --shortcut ccwork
```

`ccam work` e `ccwork` avviano l'account, ed entrambi funzionano anche negli script.

## Fianco a fianco

| | Alias della shell | ccam |
| --- | --- | --- |
| Login separati | Sì | Sì |
| Configurazione per account | Modifica di `~/.bashrc`, login | `ccam add` |
| Funziona in script e altri programmi | No, gli alias esistono solo nella shell | Sì, `ccam work` e le scorciatoie come `ccwork` sono veri comandi |
| Impostazioni, `CLAUDE.md` ed estensioni condivisi | Copia o link a mano | Per elemento, con backup |
| Server MCP condivisi | A mano, per ogni account | `ccam mcp add` una sola volta |
| Continuare una conversazione in un altro account | No | `ccam next --continue` |
| Utilizzo di tutti gli account | No | `ccam usage` |
| Scegliere l'account con più margine | No | `ccam next --auto` |
| Rieseguire un job `-p` con un altro account al limite | No | `ccam --fallback` |
| Un account per ogni cartella di progetto | No | `ccam folder` |
| Nome dell'account nella riga di stato | No | `ccam statusline install` |
| Riparare link interrotti | A mano | `ccam doctor --fix` |
| Software aggiuntivo | Nessuno | Un solo binario statico |

Scegli gli alias se hai due account e non vuoi condividere nulla tra loro. Scegli ccam se vuoi che gli account si comportino come un'unica configurazione. Puoi iniziare con gli alias e passare a ccam in seguito: `ccam import` registra le cartelle `~/.claude-*` esistenti dove si trovano.

[Come usare più account su un solo computer](/guides/multiple-claude-code-accounts) spiega entrambi gli approcci passo dopo passo, e [ccam e gli altri strumenti per cambiare account](/compare/account-switchers) confronta ccam con gli strumenti che scambiano un unico login.
