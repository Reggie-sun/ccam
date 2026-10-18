---
title: Condividere server MCP, impostazioni e cronologia
description: Come ccam condivide server MCP, settings.json, CLAUDE.md, agenti, skill e cronologia tra account Claude Code, tenendo separato ogni login.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: e4b31e5098e3
---

Ogni account mantiene il proprio login. Tutto il resto si può condividere elemento per elemento: ccam collega il file o la cartella dell'account a una cartella di gruppo condivisa, così una modifica fatta in un account compare in tutti gli account di quel gruppo.

## Cosa si può condividere

| Elemento | Opzione | Cosa viene condiviso |
| --- | --- | --- |
| Server MCP | `--mcp` | Il `mcp.json` del gruppo, passato a Claude Code all'avvio |
| Impostazioni | `--settings` | `settings.json`: permessi, hook, modello, riga di stato |
| Memoria globale | `--memory` | `CLAUDE.md` |
| Estensioni | `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| Cronologia delle conversazioni | `--history` | `projects/` e `history.jsonl` |

I nuovi account condividono tutto tranne la cronologia delle conversazioni.

I file di login `.credentials.json` e `.claude.json` non vengono mai condivisi. ccam collega solo percorsi presi da un elenco fisso, e questi due non ne fanno parte.

## Gruppi

`on` indica il gruppo `default`. Puoi anche indicare un nome, in modo che solo alcuni account condividano un elemento:

```sh
ccam share work --history team
ccam share side --history team
```

Ora `work` e `side` condividono le conversazioni tra loro, ma non con gli altri tuoi account.

## Attivare la condivisione

Quando un account inizia a condividere un elemento, ccam:

1. stampa un piano con ogni file che sposterà, copierà o collegherà, e chiede conferma prima di iniziare;
2. sposta i file dell'account nella cartella del gruppo se il gruppo è ancora vuoto, oppure aggiunge i file che il gruppo non ha ancora;
3. quando lo stesso file è diverso, chiede se mantenere la versione del gruppo, usare la tua o unire il JSON;
4. conserva l'originale in `.ccam-backup/<time>/` dentro la cartella dell'account;
5. sostituisce l'originale con un link al gruppo.

I file delle conversazioni prendono il nome dal loro ID di sessione, quindi le cronologie di più account si uniscono senza conflitti.

## Disattivare la condivisione

```sh
ccam share work --history off
```

Il link viene rimosso e l'account riceve una propria copia del contenuto attuale del gruppo. Aggiungi `--fresh` per partire invece da vuoto. Il gruppo mantiene il suo contenuto per gli altri account.

## Server MCP

Claude Code salva i server MCP in `.claude.json`, insieme al login, quindi quel file non può essere condiviso. ccam tiene i server condivisi nel `mcp.json` del gruppo e avvia Claude Code con `--mcp-config=<group>/mcp.json`. Le tue opzioni `--mcp-config` continuano a valere in aggiunta.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

Un server che aggiungi all'interno di una sessione con `claude mcp add` viene salvato solo per quell'account. La volta successiva che avvii quell'account, ccam indica il server e il comando che lo sposta nel gruppo:

```sh
ccam mcp import work gh
```

I server MCP che accedono con OAuth salvano i loro token insieme al login dell'account, quindi autorizzali una volta in ciascun account.

Se preferisci che i server compaiano in `claude mcp list` e `/mcp` come normali server utente, imposta `mcp_strategy = "sync"` in `~/.config/ccam/config.toml`. ccam li scrive allora nel `.claude.json` di ogni account prima dell'avvio, rispettando il lock sul file usato da Claude Code stesso e con un backup.

## Due account nella stessa cartella

Gli account che condividono la cronologia delle conversazioni scrivono negli stessi file. Se avvii un secondo Claude Code in una cartella in cui è già in esecuzione un altro account dello stesso gruppo, ccam ti avvisa e chiede se continuare.

## Quando qualcosa si rompe

```sh
ccam doctor
ccam doctor --fix
```

Se un percorso condiviso è diventato un file normale, per esempio dopo che hai ripristinato un backup a mano, ccam si ferma prima di avviare Claude Code e rimanda a `ccam doctor --fix`, che unisce il file al gruppo e conserva un backup.
