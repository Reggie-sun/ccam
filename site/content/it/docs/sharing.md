---
title: Condividere server MCP, impostazioni e cronologia
description: Come i gruppi di ccam condividono server MCP, settings.json, CLAUDE.md, skill e cronologia tra account Claude Code, tenendo separato ogni login.
updated: 2026-10-26
order: 4
translation: machine
sourceHash: 2d8029ad6ec8
---

Ogni account mantiene il proprio login. Tutto il resto si può condividere tramite i **gruppi**: un gruppo indica gli elementi che condivide e gli account che ne fanno parte, e ccam collega i file di ciascuno di quegli account alla cartella del gruppo, così una modifica fatta in un account compare in tutti.

## Cosa si può condividere

| Elemento | Cosa viene condiviso |
| --- | --- |
| `mcp` | i server MCP: il `mcp.json` del gruppo, passato a Claude Code all'avvio |
| `settings` | `settings.json`: permessi, hook, modello, riga di stato |
| `memory` | il `CLAUDE.md` globale |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | la cronologia delle conversazioni: `projects/` e `history.jsonl` |

Il tuo primo account crea il gruppo `default`, che condivide tutto tranne la cronologia delle conversazioni, e gli account successivi vi si uniscono.

I file di login `.credentials.json` e `.claude.json` non vengono mai condivisi. ccam collega solo percorsi presi da un elenco fisso, e questi due non ne fanno parte.

Questa pagina tratta le attività di tutti i giorni. Ogni opzione è descritta nelle pagine delle funzionalità: [condivisione con i gruppi](/features/share), [server MCP condivisi](/features/mcp) e [copia di una configurazione esistente](/features/migrate).

## Gruppi

```sh
ccam group                                   # the groups, what they share, their accounts
ccam group set default --add history         # every account in default shares history
ccam group add team --items history --members work,side
```

Un account può stare in più gruppi, ma ogni elemento proviene da un solo gruppo. Quindi, per condividere la cronologia tra `work` e `side` ma non con gli altri tuoi account, lascia la cronologia fuori da `default` e dai ai due un gruppo tutto loro, come sopra. `ccam share work` mostra cosa condivide un account e tramite quale gruppo. La [dashboard](/features/dashboard) fa lo stesso nella sua scheda Groups.

## Quando un account inizia a condividere

Quando un account inizia a condividere un elemento (si unisce a un gruppo, oppure il suo gruppo aggiunge l'elemento), ccam:

1. stampa un piano con ogni file che sposterà, copierà o collegherà per ogni account, e chiede conferma prima di iniziare;
2. sposta i file dell'account nella cartella del gruppo se il gruppo è ancora vuoto, oppure aggiunge i file che il gruppo non ha ancora;
3. quando lo stesso file è diverso, chiede se mantenere la versione del gruppo, usare la tua o unire il JSON;
4. conserva l'originale in `.ccam-backup/<time>/` dentro la cartella dell'account;
5. sostituisce l'originale con un link al gruppo.

I file delle conversazioni prendono il nome dal loro ID di sessione, quindi le cronologie di più account si uniscono senza conflitti.

ccam non modifica un account mentre Claude Code è in esecuzione con quell'account, comunque sia stato avviato (comprese le sessioni in background), perché la sessione continua a scrivere nei file che vengono spostati. Chiudila prima. ccam controlla di nuovo subito prima di salvare, così intercetta anche una sessione avviata mentre ccam aspettava la tua risposta.

## Quando un account smette di condividere

```sh
ccam group leave team work                   # work leaves the group
ccam group set default --remove history      # nobody in default shares history any more
```

Il link viene rimosso e l'account riceve una propria copia del contenuto attuale del gruppo; i server MCP del gruppo finiscono nel suo `.claude.json`, accanto ai server propri. Aggiungi `--fresh` per partire invece da vuoto; il piano elenca allora i server MCP che smettono di caricarsi. Una cartella importata come `~/.claude` conserva sempre una copia. Il gruppo mantiene il suo contenuto per gli altri account.

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
- la cronologia delle conversazioni (`projects/` e `history.jsonl`), ma solo se aggiungi `--history`.

La cartella di origine viene solo letta: niente al suo interno viene spostato, collegato o modificato, e `claude` da solo continua a funzionare come prima. I file di login `.credentials.json` e `.claude.json` non vengono mai copiati nel gruppo.

I server MCP che accedono con OAuth conservano i loro token nel login, non nell'elenco dei server. ccam non copia quegli accessi: un accesso copiato contiene lo stesso refresh token, e il primo rinnovo disconnette da quel server tutti gli altri che lo possiedono, `~/.claude` compresa. Il piano e il riepilogo indicano gli account che devono accedere da sé a quei server con `/mcp`.

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
Group default: personal
To share them:
  ccam group join default work
```

Alla fine, ccam elenca gli account del gruppo. Se il gruppo non condivide ancora uno degli elementi, o alcuni account non ne fanno parte, stampa i comandi `ccam group` che lo cambiano.

Quando il gruppo ha già una copia diversa di un file, ccam chiede, come quando attivi la condivisione, se mantenere la copia del gruppo o usare la tua. `settings.json` e gli altri file JSON si possono anche unire, e `history.jsonl` riga per riga. Le copie precedenti del gruppo vengono salvate in `.ccam-backup/<time>/` all'interno della cartella del gruppo.

| Opzione | Effetto |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Copia solo gli elementi indicati, senza la checklist. Fuori da un terminale, senza opzioni di elemento, ccam copia tutto tranne la cronologia delle conversazioni. |
| `-g`, `--group <group>` | Copia in questo gruppo invece che in `default` |
| `-y`, `--yes` | Copia senza chiedere; quando un file è diverso, mantiene la copia del gruppo |

Metti `--dry-run` prima del comando, come in `ccam --dry-run migrate`, per vedere il piano senza copiare nulla. Per copiare da un'altra cartella, indicala: `ccam migrate ~/.claude-work --group work`.

## Server MCP

Claude Code salva i server MCP in `.claude.json`, insieme al login, quindi quel file non può essere condiviso. ccam tiene i server condivisi nel `mcp.json` del gruppo e avvia Claude Code con `--mcp-config=<group>/mcp.json`. Le tue opzioni `--mcp-config` continuano a valere in aggiunta.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

Un server che aggiungi all'interno di una sessione con `claude mcp add` viene salvato solo per quell'account, e per la cartella in cui ti trovavi (l'ambito locale di Claude Code; `-s user` lo rende disponibile all'account in ogni cartella). La volta successiva che avvii quell'account, ccam indica il server, con la sua cartella, e il comando che lo sposta nel gruppo:

```sh
ccam mcp import work gh
```

Se il gruppo ha già un server diverso con quel nome, ccam lo salta e te lo segnala. `--force` sostituisce il server del gruppo e conserva il vecchio `mcp.json` nella cartella `.ccam-backup/` del gruppo.

I server MCP che accedono con OAuth salvano i loro token insieme al login dell'account, quindi accedi a ciascuno di essi una volta in ogni account con `/mcp`. `ccam ext` e `ccam doctor` indicano i server a cui un account non ha effettuato l'accesso mentre un altro account sì.

Se preferisci che i server compaiano in `claude mcp list` e `/mcp` come normali server utente, imposta `mcp_strategy = "sync"` in cima a `~/.config/ccam/config.toml`, sopra il primo `[[group]]`, `[[profile]]` o `[[folder]]`. ccam li scrive allora nel `.claude.json` di ogni account prima dell'avvio, rispettando il lock sul file usato da Claude Code stesso e con un backup. Se l'account ha un proprio server con lo stesso nome di un server del gruppo, ccam mantiene il server dell'account e lo segnala.

## Due account nella stessa cartella

Gli account che condividono la cronologia delle conversazioni scrivono negli stessi file. Se avvii una seconda sessione di Claude Code in una cartella in cui è già in esecuzione un altro account dello stesso gruppo, ccam ti avvisa e chiede se continuare.

## Quando qualcosa si rompe

```sh
ccam doctor
ccam doctor --fix
```

Se un percorso condiviso è diventato un file normale, per esempio dopo che hai ripristinato un backup a mano, ccam si ferma prima di avviare Claude Code e rimanda a `ccam doctor --fix`. La correzione ricollega il percorso al gruppo: il file diventa la copia del gruppo quando il gruppo non ne ha una; altrimenti il gruppo tiene la propria copia e la tua va in un backup che `--fix` indica. Anche una modifica dei gruppi rimasta a metà viene completato da `ccam doctor --fix`.
