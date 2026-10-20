---
title: Condividere server MCP, impostazioni e cronologia
description: Come ccam condivide server MCP, settings.json, CLAUDE.md, agenti, skill e cronologia tra account Claude Code, tenendo separato ogni login.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 888d466b29e9
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

Questa pagina tratta le attività di tutti i giorni. Ogni opzione è descritta nelle pagine delle funzionalità: [condivisione tra account](/features/share), [server MCP condivisi](/features/mcp) e [copia di una configurazione esistente](/features/migrate).

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

ccam non modifica la condivisione mentre Claude Code è in esecuzione con l'account, perché la sessione continua a scrivere nei file che vengono spostati. Chiudila prima.

## Disattivare la condivisione

```sh
ccam share work --history off
```

Il link viene rimosso e l'account riceve una propria copia del contenuto attuale del gruppo. Aggiungi `--fresh` per partire invece da vuoto. Il gruppo mantiene il suo contenuto per gli altri account.

## Copiare una configurazione esistente

Se usavi Claude Code prima di ccam, le tue impostazioni si trovano in `~/.claude`. `ccam migrate` le copia in un gruppo, così ogni account che condivide quegli elementi le riceve:

```sh
ccam migrate
```

In un terminale, ccam mostra prima una checklist degli elementi. È tutto selezionato tranne la cronologia delle conversazioni; premi Spazio per cambiare un elemento e Invio per vedere il piano. Gli elementi per cui la cartella non contiene nulla sono indicati con `not found`.

Cosa copia:

- `settings.json` e il `CLAUDE.md` globale;
- le estensioni `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/` e `rules/`;
- i tuoi server MCP utente, cioè la voce `mcpServers` di `~/.claude.json`, nel `mcp.json` del gruppo;
- gli accessi salvati (OAuth) dei server MCP che copia, nel file di login di ciascun account;
- la cronologia delle conversazioni (`projects/` e `history.jsonl`), ma solo se aggiungi `--history`.

La cartella di origine viene solo letta: niente al suo interno viene spostato, collegato o modificato, e `claude` da solo continua a funzionare come prima. I file di login `.credentials.json` e `.claude.json` non vengono mai copiati nel gruppo.

I server MCP che accedono con OAuth conservano i loro token nel login, non nell'elenco dei server. Per evitarti di accedere di nuovo, ccam copia quegli accessi nel file di login di ogni account che riceve i server tramite il gruppo, a meno che l'account non abbia già i propri. Salta gli account in esecuzione, quelli in cui non hai ancora effettuato il login, quelli registrati sul posto con `ccam import` e quelli che conservano il login nel Portachiavi di macOS, e li elenca per nome; in quegli account accedi ai server con `/mcp`. Un account che inizia a condividere i server MCP in seguito riceve gli accessi quando esegui di nuovo `ccam migrate`.

ccam stampa il piano e chiede conferma prima di copiare qualsiasi cosa:

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Accounts that share these through group default:
  mcp         personal
  settings    personal
  memory      personal
  extensions  personal
```

Alla fine, ccam elenca gli account che condividono ogni elemento tramite il gruppo. Se alcuni account non lo fanno, stampa i comandi `ccam share` che attivano la condivisione per loro.

Quando il gruppo ha già una copia diversa di un file, ccam chiede, come quando attivi la condivisione, se mantenere la copia del gruppo o usare la tua. `settings.json` e gli altri file JSON si possono anche unire, e `history.jsonl` riga per riga. Le copie precedenti del gruppo vengono salvate in `.ccam-backup/<time>/` all'interno della cartella del gruppo.

| Opzione | Effetto |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Copia solo gli elementi indicati, senza la checklist; senza di esse e fuori da un terminale, tutto tranne la cronologia delle conversazioni |
| `-g`, `--group <group>` | Copia in questo gruppo invece che in `default` |
| `-y`, `--yes` | Copia senza chiedere; quando un file è diverso, mantiene la copia del gruppo |

Metti `--dry-run` prima del comando, come in `ccam --dry-run migrate`, per vedere il piano senza copiare nulla. Per copiare da un'altra cartella, indicala: `ccam migrate ~/.claude-work --group work`.

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

Se il gruppo ha già un server diverso con quel nome, ccam lo salta e te lo segnala. `--force` sostituisce il server del gruppo e conserva il vecchio `mcp.json` nella cartella `.ccam-backup/` del gruppo.

I server MCP che accedono con OAuth salvano i loro token insieme al login dell'account, quindi autorizzali una volta in ciascun account, oppure copia gli accessi da `~/.claude` con `ccam migrate`.

Se preferisci che i server compaiano in `claude mcp list` e `/mcp` come normali server utente, imposta `mcp_strategy = "sync"` in `~/.config/ccam/config.toml`. ccam li scrive allora nel `.claude.json` di ogni account prima dell'avvio, rispettando il lock sul file usato da Claude Code stesso e con un backup.

## Due account nella stessa cartella

Gli account che condividono la cronologia delle conversazioni scrivono negli stessi file. Se avvii un secondo Claude Code in una cartella in cui è già in esecuzione un altro account dello stesso gruppo, ccam ti avvisa e chiede se continuare.

## Quando qualcosa si rompe

```sh
ccam doctor
ccam doctor --fix
```

Se un percorso condiviso è diventato un file normale, per esempio dopo che hai ripristinato un backup a mano, ccam si ferma prima di avviare Claude Code e rimanda a `ccam doctor --fix`, che unisce il file al gruppo e conserva un backup.
