---
title: Presentiamo ccam
description: Perché esiste ccam, cosa fa per chi ha più account Claude Code e come evita di essere d'intralcio.
updated: 2026-10-26
order: 1
translation: machine
sourceHash: e4370f29b818
---

ccam è un piccolo strumento a riga di comando per chi usa più di un account Claude Code sullo stesso computer: un account personale e uno di lavoro, oppure più account per un unico progetto lungo.

## Il problema

Claude Code supporta più account tramite [`CLAUDE_CONFIG_DIR`](/guides/claude-config-dir): una cartella per account. In pratica, ogni cartella è un'isola. Impostazioni, `CLAUDE.md`, server MCP, agenti e skill vanno configurati di nuovo in ogni cartella e tenuti allineati a mano. Le conversazioni restano dove sono iniziate. E quando un account raggiunge il limite di utilizzo, lo scopri provando.

## Cosa fa ccam

- `ccam add` crea la cartella, esegue il login e, se vuoi, assegna all'account una scorciatoia a cui dai tu il nome, come `ccwork`.
- `ccam work` e `ccwork` avviano Claude Code con quell'account e passano ogni argomento invariato.
- Server MCP, impostazioni, `CLAUDE.md` ed estensioni sono condivisi tramite link, elemento per elemento e account per account. I login non vengono mai condivisi.
- `ccam usage` mostra l'utilizzo su 5 ore e settimanale di ogni account; `ccam next --continue` riprende la conversazione nell'account successivo.

## Senza essere d'intralcio

Uno strumento che esegui prima di ogni sessione non dovrebbe avere un costo che tu possa notare. ccam legge un file di configurazione, controlla i link condivisi con una o due chiamate di sistema per ciascuno e poi si sostituisce con `claude`. Aggiunge circa 4 millisecondi e, per impostazione predefinita, non è in esecuzione mentre lavori.

Inoltre tratta i tuoi file con cura. Ogni modifica alla condivisione stampa un piano prima di iniziare e conserva gli originali in una cartella di backup. `ccam doctor` controlla link, permessi e login, e ripara ciò che può riparare in sicurezza.

## Provalo

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add
```

La [guida rapida](/docs/quick-start) mostra passo passo come aggiungere due account, e [come usare più account Claude Code](/guides/multiple-claude-code-accounts) spiega l'approccio con e senza ccam.

ccam è open source e non ufficiale: non è affiliato ad Anthropic né approvato da Anthropic. La possibilità di usare più di un account dipende dai termini che si applicano a ciascuno. Feedback e traduzioni sono benvenuti su [GitHub](https://github.com/jhyunwoo/ccam).
