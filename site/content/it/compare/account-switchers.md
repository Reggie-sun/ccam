---
title: "Cambiare account Claude Code: ccam e le alternative"
description: Come si confronta ccam con claude-swap, CCS, clauth, CAAM, cux e altri strumenti per cambiare account Claude Code, e quale approccio scegliere per ogni uso.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: dc907c88e2bf
---

Gli strumenti per cambiare account Claude Code seguono uno di due approcci. La maggior parte scambia il login dentro l'unica cartella `~/.claude`, quindi è attivo un solo account alla volta. Altri, tra cui ccam, danno a ogni account la sua cartella con `CLAUDE_CONFIG_DIR`, quindi ogni account resta connesso e più account possono girare insieme. ccam aggiunge impostazioni e server MCP condivisi, e continua una conversazione nella cartella di un altro account.

## Scambiare un login o usare cartelle separate?

**Scambiare un login** mantiene un solo insieme di impostazioni, server MCP e cronologia, perché c'è una sola cartella. Il cambio è rapido, e una sessione in corso può adottare il nuovo login. Però è attivo un solo account alla volta, per tutti i terminali e gli editor insieme.

**Le cartelle separate** sono il metodo che [descrive la documentazione di Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Ogni account resta connesso, e puoi eseguire una sessione di lavoro e una personale in parallelo. Il prezzo è che ogni cartella parte con impostazioni, server MCP e cronologia propri, a meno che uno strumento non li condivida.

## Gli strumenti a confronto

Tratto dalla documentazione di ciascun progetto il 2026-10-19. "Non documentato" significa che il progetto non lo dice, non che la funzione manchi.

| Strumento | Approccio | Più account insieme | Utilizzo su 5 ore e settimanale | Quando si raggiunge un limite | Piattaforme |
| --- | --- | --- | --- | --- | --- |
| ccam | Cartelle separate | Sì | Sì, `ccam usage` | Esegui `ccam next --auto --continue` per continuare la conversazione nell'account con più margine | Linux, macOS, WSL |
| [claude-swap](https://github.com/realiti4/claude-swap) | Scambia il login; cartelle separate con `cswap run` | Sì, con `cswap run` | Sì | `cswap auto` cambia il login al superamento di una soglia; la sessione in corso continua | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Cartelle separate | Sì | In parte, in un'app della barra dei menu di macOS | I pool di account passano le richieste all'account successivo tramite un proxy locale | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Scambia il login; cartelle separate con `clauth start` | Sì, con `clauth start` | Sì | Cambia seguendo una catena di fallback | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Scambia il login, oppure profili isolati | Sì, con i profili isolati | Sì | `caam run` riprova i comandi non interattivi su un altro account | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Scambia il login e fa da wrapper a `claude` | Non documentato | Sì | Scambia il login e riprende la stessa conversazione | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Scambia il login | No | Sì | Non documentato | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Scambia il login | No | Sì | Cambio facoltativo al 100% | macOS |

ccswap è un fork di claude-swap che gestisce anche i login di Codex. Diversi di questi strumenti gestiscono anche altri agenti di programmazione, come Codex o Gemini.

## Su cosa si concentra ccam

- Ogni account resta connesso. Cambiare non disconnette mai un account, e `ccam work` e `ccam personal` possono girare in due terminali contemporaneamente.
- Server MCP, `settings.json`, `CLAUDE.md`, agenti, skill e plugin sono condivisi tra le cartelle tramite link, elemento per elemento e account per account. I login non vengono mai condivisi.
- Con la cronologia condivisa, un solo comando, `ccam next --auto --continue`, apre la stessa conversazione nella cartella dell'account con più margine.
- ccam avvia Claude Code e si sostituisce con esso, quindi nulla resta in esecuzione in background. Non tiene d'occhio il tuo utilizzo e non cambia account da solo; decidi tu quando.

## Dove altri strumenti fanno di più

- claude-swap, CCS, clauth, CAAM e cux girano su Windows senza WSL. ccam no.
- claude-swap, clauth e cux cambiano da soli quando un account supera una soglia di utilizzo. ccam lascia a te il cambio.
- CCSwitcher, Claude Switcher e claude-swap mostrano l'utilizzo nella barra dei menu di macOS.
- CCS, CAAM, ccswap e Claude Switcher gestiscono anche Codex e altri agenti di programmazione.

## Quale fa per te?

Scegli uno strumento che scambia il login se usi un account alla volta, vuoi che il cambio avvenga in automatico o sei su Windows senza WSL. Scegli ccam se vuoi account in parallelo, un'unica configurazione condivisa e la stessa conversazione sull'account successivo. Se non ti serve condividere nulla, bastano [semplici alias della shell](/compare/aliases).

Se puoi usare più account dipende dai termini di ciascun account. Controllali prima di farlo.
