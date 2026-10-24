---
title: Cosa fare al limite di utilizzo di Claude Code
description: Aspetta il reset o continua subito la stessa conversazione su un altro account Claude Code. Come funzionano i limiti di 5 ore e settimanale, passo dopo passo.
updated: 2026-10-24
order: 3
translation: machine
sourceHash: "049a46684c7b"
---

Quando un account raggiunge il limite di utilizzo (spesso chiamato rate limit), hai due possibilità: aspettare che il limite si azzeri oppure continuare subito la stessa conversazione su un altro account che possiedi. Con ccam la seconda richiede un solo comando, `ccam next --auto --continue`.

## Quali limiti ha Claude Code?

Con i piani Pro e Max, Claude Code ha un [limite di sessione di cinque ore e un limite settimanale](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Entrambi sono condivisi con Claude sul web e nelle app, e l'utilizzo conta su entrambi contemporaneamente. Quando uno si esaurisce, Claude Code [mostra un messaggio](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) come "You've hit your session limit · resets 3:45pm". Cambiare modello non serve, perché i limiti valgono per tutti i modelli.

## Opzione 1: aspettare il reset

Se hai un solo account, aspetta. Le versioni recenti di Claude Code [aspettano nella sessione aperta](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) e riprendono il lavoro da sole quando il limite si azzera. `/usage` dentro Claude Code mostra quanta parte dei limiti del tuo piano hai usato.

## Opzione 2: continuare subito su un altro account

Se hai un secondo account, puoi andare avanti subito con la stessa conversazione. Configuralo una volta:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

Ora la riga di stato mostra l'account e il suo utilizzo su 5 ore mentre lavori, così il limite non ti coglie di sorpresa.

Quando raggiungi il limite, esci da Claude Code (`/exit`), poi:

```sh
ccam usage
```

```console
 1 personal  just now
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
→ Most headroom: work   (open it with ccam next --auto)
```

Apri la conversazione nell'account con più margine:

```sh
ccam next --auto --continue
```

Claude Code si avvia come `work` con la conversazione più recente della cartella corrente, compreso tutto ciò che è stato detto finora. Condividere la cronologia (`ccam group set default --add history` qui sopra) la mantiene un'unica conversazione; senza, Claude Code copia la conversazione in `work`, e ccam ti dice quanti token rilegge quella prima risposta.

Con il [plugin di ccam](/docs/install#claude-code-plugin) in Claude Code, una sessione che si ferma al limite mostra una notifica e registra la finestra piena, così `ccam next --auto` non sceglie di nuovo quell'account prima del suo reset, e `/ccam:switch` (oppure `/ccam:switch work`) annota la conversazione per `ccam next --continue` senza lasciare nulla da indovinare.

**Posti nei piani Team ed Enterprise.** La scelta automatica (`ccam next`, `--fallback`, `ccam pick`, il tasto `n` della dashboard) passa anche da un posto all'altro della stessa organizzazione Claude Team o Enterprise. Il centro assistenza di Anthropic dice: "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (ogni posto di un piano Team è assegnato nominativamente a una sola persona, e solo quella persona può usarlo) e "Users can't route their work through another user's account to use that user's remaining usage" (gli utenti non possono far passare il proprio lavoro attraverso l'account di un altro utente per sfruttare l'utilizzo residuo di quell'utente) ([piani Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); lo stesso vale per i [piani Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Usare la scelta automatica tra posti diversi può quindi violare i termini di Anthropic. Verifica con la tua organizzazione e con Anthropic prima di farci affidamento. I piani personali (Pro, Max) sono trattati nella [nota sui termini nelle domande frequenti](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## ccam cambia account automaticamente?

No. ccam non è in esecuzione mentre lavori: avvia Claude Code e si sostituisce con esso. Decidi tu quando cambiare, e `ccam next --auto` sceglie l'account con l'utilizzo su 5 ore più basso e salta gli account al 100%, senza login o di un'altra organizzazione. Con `relay = true` in `config.toml`, ccam invece resta in attesa accanto alla sessione, così `/ccam:switch` prosegue da solo, mantenendo le opzioni, la modalità dei permessi e il livello di effort della sessione. Alcuni altri strumenti invece cambiano il login di una sessione in corso; [ccam e gli altri strumenti per cambiare account](/compare/account-switchers) spiega la differenza.

## Cose da tenere a mente

- Esegui `ccam next --continue` nella stessa cartella in cui stavi lavorando; le conversazioni sono salvate per cartella di progetto.
- Evita di eseguire due account sulla stessa conversazione contemporaneamente. ccam ti avvisa se ci provi.
- Senza cronologia condivisa, la conversazione viene copiata nell'altro account e le due copie proseguono ciascuna per conto proprio; i checkpoint da cui `/rewind` ripristina il codice restano nel primo account.
- Il primo messaggio sul nuovo account può consumare più limite del solito, perché quell'account non ha ancora la conversazione in cache.
- Se puoi passare da un account all'altro in questo modo dipende dai termini di ciascun account.

Vedi [limiti di utilizzo](/docs/usage-limits) per tutte le opzioni di `ccam usage` e `ccam next`.
