---
title: Cosa fare quando Claude Code raggiunge il limite di utilizzo
description: Controlla quanto margine resta a ogni account Claude Code e continua la stessa conversazione su un altro account quando uno raggiunge il limite.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 767a731b912e
---

Quando un account raggiunge il limite di utilizzo di Claude Code, continua la stessa conversazione in un altro account di tua proprietà: condividi una volta la cronologia delle conversazioni, poi esegui `ccam next --continue`. Se hai un solo account, attendi l'orario di azzeramento indicato da Claude Code.

## Prima che succeda

Configura gli account in modo che il cambio richieda un solo comando:

```sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

Ora la riga di stato mostra l'account e il suo utilizzo su 5 ore mentre lavori, così il limite non ti coglie di sorpresa.

## Quando raggiungi il limite

Esci da Claude Code (`/exit`), poi:

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22m       █████░░░  68% 1d 9h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
```

Apri la conversazione nell'account con più margine:

```sh
ccam next --auto --continue
```

Claude Code si avvia come `work` con la conversazione più recente della cartella corrente, compreso tutto ciò che è stato detto finora.

## Cose da tenere a mente

- Esegui `ccam next --continue` nella stessa cartella in cui stavi lavorando; le conversazioni sono salvate per cartella di progetto.
- Evita di usare due account sulla stessa conversazione nello stesso momento. ccam ti avvisa se ci provi.
- Anche l'altro account deve condividere la cronologia delle conversazioni, nello stesso gruppo.
- La possibilità di passare da un account all'altro in questo modo dipende dai termini di ciascun account.

Vedi [limiti di utilizzo](/docs/usage-limits) per tutte le opzioni di `ccam usage` e `ccam next`.
