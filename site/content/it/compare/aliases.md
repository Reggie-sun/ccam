---
title: ccam a confronto con gli alias della shell
description: Alias della shell con CLAUDE_CONFIG_DIR e ccam a confronto, fianco a fianco, per usare più account Claude Code.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 13da69070650
---

Un alias della shell per ogni account è il modo più semplice per usare più account Claude Code. ccam fa la stessa cosa dietro le quinte e aggiunge condivisione, cambio di account e monitoraggio dell'utilizzo.

| | Alias della shell | ccam |
| --- | --- | --- |
| Login separati | Sì | Sì |
| Configurazione per account | Modifica di `~/.bashrc`, login | `ccam add` |
| Funziona in script e altri programmi | No, gli alias esistono solo nella shell | Sì, `claude2` è un vero comando |
| Impostazioni, `CLAUDE.md` ed estensioni condivisi | Copia o link a mano | Per elemento, con backup |
| Server MCP condivisi | A mano, per ogni account | `ccam mcp add` una sola volta |
| Continuare una conversazione in un altro account | No | `ccam next --continue` |
| Utilizzo di tutti gli account | No | `ccam usage` |
| Scegliere l'account con più margine | No | `ccam next --auto` |
| Nome dell'account nella riga di stato | No | `ccam statusline install` |
| Riparare link interrotti | A mano | `ccam doctor --fix` |
| Software aggiuntivo | Nessuno | Un solo binario statico |

Scegli gli alias se hai due account e non vuoi condividere nulla tra loro. Scegli ccam se vuoi che gli account si comportino come un'unica configurazione. Puoi iniziare con gli alias e passare a ccam in seguito: `ccam import` registra le cartelle `~/.claude-*` esistenti dove si trovano.

[Come usare due account su un solo computer](/guides/multiple-claude-code-accounts) spiega entrambi gli approcci passo dopo passo.
