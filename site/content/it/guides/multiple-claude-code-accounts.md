---
title: Usare più account Claude Code su un solo computer
description: Esegui più account Claude Code in parallelo con CLAUDE_CONFIG_DIR, a mano con gli alias della shell o in automatico con ccam, e passa dall'uno all'altro.
updated: 2026-10-25
order: 1
translation: machine
sourceHash: "5c5dbd9fcd23"
---

Dai a ogni account la sua cartella di configurazione. Claude Code legge la cartella dalla variabile d'ambiente `CLAUDE_CONFIG_DIR`, e ogni cartella contiene il proprio login, le impostazioni e la cronologia, quindi ogni account resta connesso e puoi eseguirli contemporaneamente. Puoi configurarlo a mano con un alias della shell oppure lasciare che se ne occupi ccam.

## Come lo configuro a mano?

Questo è l'approccio che [descrive la documentazione di Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Aggiungi una riga per ogni account in più a `~/.bashrc` o `~/.zshrc`:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Apri un nuovo terminale ed esegui `claude-work`. Claude Code ti chiede di effettuare il login, perché la cartella è nuova. Da quel momento `claude` è il tuo primo account e `claude-work` il secondo. Aggiungi un alias per ogni account, ciascuno con la sua cartella.

Su Windows senza WSL, una funzione PowerShell nel tuo profilo (`notepad $PROFILE`) fa la stessa cosa:

```powershell
function claude-work { $env:CLAUDE_CONFIG_DIR = "$HOME\.claude-work"; try { claude @args } finally { Remove-Item Env:CLAUDE_CONFIG_DIR } }
```

## Cosa manca con gli alias?

Funziona bene per due account. Cose da sapere:

- Ogni cartella parte vuota: niente impostazioni, niente `CLAUDE.md`, niente server MCP, niente agenti o skill. Copi o colleghi quello che ti serve e tieni allineate le copie da solo.
- Gli alias esistono solo nelle shell interattive. Script, Makefile e altri programmi non li vedono.
- I server MCP stanno in `.claude.json`, accanto al login, quindi non si possono semplicemente collegare tra cartelle.
- Le conversazioni restano nell'account in cui sono iniziate, quindi non puoi continuarne una in un altro account.
- Non c'è nulla che mostri l'utilizzo di tutti gli account insieme.

## Posso cambiare account con /login?

Sì, se cambi solo ogni tanto. Digitare `/login` dentro Claude Code [cambia l'account](https://code.claude.com/docs/en/quickstart) della cartella in cui ti trovi. È connesso un solo account alla volta, quindi rifai il login a ogni cambio e non puoi eseguire due account in parallelo. Le cartelle separate evitano entrambi i problemi.

## Come lo configuro con ccam?

ccam automatizza la stessa idea e aggiunge le parti mancanti:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

Ora `ccam personal` e `ccam work` avviano i due account, e lo stesso fa `ccwork`, la scorciatoia che hai dato all'account work. Poiché sono veri programmi e non alias, funzionano anche negli script. Esegui `ccam` da solo per scegliere un account da un elenco, e `ccam use work` per renderne uno quello predefinito.

Impostazioni, `CLAUDE.md`, server MCP ed estensioni sono condivisi fin dall'inizio. `ccam next --continue` porta una conversazione sull'altro account, e la cronologia delle conversazioni condivisa, che puoi attivare, la mantiene un'unica conversazione. `ccam usage` mostra l'utilizzo su 5 ore e settimanale di tutti gli account insieme.

Se usi già `~/.claude`, esegui `ccam import` invece di aggiungere di nuovo quell'account, oppure `ccam migrate` per copiarne impostazioni, server MCP ed estensioni nelle impostazioni condivise.

## Come faccio a sapere quale account è in esecuzione?

Dentro Claude Code, `/status` mostra l'account con cui hai effettuato il login. Con ccam, esegui una volta `ccam statusline install` e la riga di stato mostra sempre il nome dell'account e il suo utilizzo, per esempio `[work] · Opus · 5h 24% · 7d 41%`.

## Funziona su Windows e WSL?

Claude Code gira [nativamente su Windows e dentro WSL](https://code.claude.com/docs/en/setup). `CLAUDE_CONFIG_DIR` funziona in entrambi, quindi su Windows nativo bastano l'alias o la funzione PowerShell vista sopra. Anche ccam gira in modo nativo su Windows: installalo in PowerShell con `irm https://ccam.dev/install.ps1 | iex`, oppure dentro WSL con il comando per Linux. Le scorciatoie diventano `ccwork.exe`, e condividere file come `settings.json` richiede la Modalità sviluppatore ([Installare ccam](/docs/install)).

## Risoluzione dei problemi

**Claude Code mi chiede il login ogni volta.** Il primo avvio in una cartella nuova lo chiede sempre. Se continua a chiederlo, il login è scaduto o è stato rimosso: esegui `ccam login work` (oppure avvia l'alias una volta ed effettua il login).

**Un server MCP compare in un solo account.** I server aggiunti dentro una sessione con `claude mcp add` appartengono a quell'account e, senza `-s user`, solo a quella cartella. Spostane uno su tutti gli account con `ccam mcp import work <server>`.

**Le impostazioni sono diverse tra gli account.** `ccam share work` mostra cosa condivide l'account e tramite quale gruppo; `ccam group join default work` gli fa condividere ciò che condivide `default`, dopo aver mostrato cosa cambierà. `ccam doctor --fix` ripara i link interrotti e conserva un backup di ogni file che modifica.

**Due account aprono la stessa conversazione.** Evita di lavorare sulla stessa conversazione da due account contemporaneamente. ccam ti avvisa se ci provi.

## Quale scegliere?

Per due account senza nulla di condiviso basta un alias. Se vuoi impostazioni e server MCP condivisi, più account, l'utilizzo a colpo d'occhio o continuare le conversazioni, ccam ti risparmia il lavoro di manutenzione. Il [confronto con gli alias](/compare/aliases) elenca le differenze fianco a fianco, e [ccam e gli altri strumenti per cambiare account](/compare/account-switchers) tratta gli strumenti che scambiano un unico login.

Se puoi usare più account dipende dai termini di ciascun account. Controllali prima di farlo.

**Posti nei piani Team ed Enterprise.** La scelta automatica (`ccam next`, `--fallback`, `ccam pick`, il tasto `n` della dashboard e, se lo attivi, [`auto_switch`](/features/next#auto-switch)) passa anche da un posto all'altro della stessa organizzazione Claude Team o Enterprise. Il centro assistenza di Anthropic dice: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (ogni posto di un piano Team è assegnato nominativamente a una sola persona, e solo quella persona può usarlo) e "Users can't route their work through another user's account to use that user's remaining usage" (gli utenti non possono far passare il proprio lavoro attraverso l'account di un altro utente per sfruttare l'utilizzo residuo di quell'utente) ([piani Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); lo stesso vale per i [piani Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Usare la scelta automatica tra posti diversi può quindi violare i termini di Anthropic. Decidi tu se usarla lì: verifica con la tua organizzazione e con Anthropic prima di farci affidamento. I piani personali (Pro, Max) sono trattati nella nota qui sopra.
