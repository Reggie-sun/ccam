---
title: Copiare una configurazione esistente
description: Copia impostazioni, CLAUDE.md, estensioni, server MCP e accessi MCP di una ~/.claude esistente in un gruppo di condivisione di ccam con ccam migrate.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 95ce6ec9e220
---

Se usavi Claude Code prima di ccam, la tua configurazione si trova in `~/.claude` e `~/.claude.json`. `ccam migrate` la copia in un gruppo di condivisione, così ogni account che condivide quegli elementi tramite il gruppo parte con la configurazione che hai già. La cartella di origine viene solo letta, e `claude` da solo continua a funzionare come prima. Eseguilo una volta dopo il tuo primo `ccam add`, e di nuovo ogni volta che un account inizia a condividere i server MCP.

## Utilizzo

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

In un terminale, senza opzioni per gli elementi e senza `--yes`, ccam mostra prima una checklist:

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Spazio (o `x`) attiva o disattiva l'elemento sotto il cursore, i tasti freccia o `j`/`k` spostano la selezione, Invio mostra il piano, ed Esc, `q` o Ctrl+C stampano `Cancelled.`. Gli elementi per cui l'origine non contiene nulla sono indicati con `not found`. Se deselezioni tutto, ccam stampa `Nothing selected.` e si ferma.

Poi ccam stampa il piano e chiede conferma prima di copiare:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        linear → personal; skipped: work (not logged in yet)
  settings    settings.json   new: 61 B
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/, MCP sign-ins for personal
MCP sign-ins were not copied to: work (not logged in yet). Sign in to those servers there with /mcp.
Accounts that share these through group default:
  mcp         personal, work
  settings    personal, work
  memory      personal, work
  extensions  personal, work
```

Un gruppo raggiunge solo gli account che condividono un elemento tramite esso, quindi il riepilogo li elenca per elemento. Per gli account che non lo fanno, ccam stampa i comandi `ccam share <name> --<item> on` che attivano la condivisione.

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `[folder]` | La cartella di Claude Code da cui copiare | `~/.claude` |
| `--settings` | Copia `settings.json` | disattivata |
| `--memory` | Copia il `CLAUDE.md` globale | disattivata |
| `--extensions` | Copia `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | disattivata |
| `--mcp` | Copia i server MCP utente (`mcpServers` in `.claude.json`) | disattivata |
| `--history` | Copia la cronologia delle conversazioni (`projects/`, `history.jsonl`) | disattivata |
| `-g`, `--group string` | Il gruppo di condivisione in cui copiare | `default` |
| `-y`, `--yes` | Copia senza chiedere; in caso di conflitto viene mantenuta la copia del gruppo | disattivata |
| `--dry-run` (globale) | Stampa il piano e si ferma prima di qualsiasi domanda o copia | disattivata |

Senza opzioni per gli elementi, in un terminale decide la checklist. Fuori da un terminale, o con `--yes`, ccam copia tutto tranne la cronologia delle conversazioni. La cronologia è esclusa per impostazione predefinita perché `projects/` può essere grande.

`--group on` indica `default`. `--group off` e `--group local` vengono rifiutati con `--group takes a group name, not "off"`. `ccam migrate` rispetta `--dry-run`, e la checklist compare anche in un dry run, così puoi vedere in anteprima una selezione.

## Come funziona

`internal/cli/migrate.go` sceglie gli elementi, poi `migrate.Build` legge l'origine e il gruppo e restituisce un piano. Non viene scritto nulla fino a `Plan.Apply`, dopo la tua conferma.

1. **Trovare l'origine.** `checkSource` rifiuta un percorso mancante, che non è una cartella o che non contiene nessuna voce nota di Claude Code. `migrate.ClaudeJSON` trova il `.claude.json` che appartiene alla cartella: al suo interno oppure, per una `~/.claude` che non ne contiene uno, `~/.claude.json` accanto a essa. È lì che Claude Code lo tiene quando `CLAUDE_CONFIG_DIR` non è impostata.
2. **Associare gli elementi ai percorsi.** Ogni percorso viene da `share.Layout`, lo stesso elenco di percorsi consentiti usato da `ccam share`. I server MCP vengono letti da `.claude.json` con `mcp.LocalServers`, che usa lo scanner `profile.TopLevelValue` invece di analizzare l'intero file.
3. **Copiare, mai spostare.** I file vanno in `~/.local/share/ccam/groups/<group>/` con `fsutil.AtomicWrite` o `fsutil.CopyFile`; i server MCP vanno nel `mcp.json` del gruppo.

Ogni elemento associa un percorso di origine a un percorso del gruppo:

<!-- motion:migrate-items -->

Le due righe della cronologia sono tratteggiate perché vengono copiate solo con `--history`. `.credentials.json` e `.claude.json` sono elencati in `share.Never` e non arrivano mai al gruppo: `Build` restituisce un errore interno se uno dei due compare nella struttura, e la scansione della cartella li salta per nome.

`buildEntry` assegna a ogni percorso uno stato, che è ciò che descrive ogni riga del piano:

<!-- motion:migrate-states -->

- **Missing**: assente nell'origine, oppure un file dove ci si aspetta una cartella. Indicato una volta per elemento come `not found, skipped`.
- **Linked**: il percorso di origine porta già alla copia del gruppo, per esempio quando l'origine è la cartella di un profilo di ccam.
- **Same**: il gruppo ha già un contenuto identico.
- **New**: il gruppo non ce l'ha, quindi viene copiato per intero.
- **Merge**: una cartella presente da entrambe le parti. `copyMissing` aggiunge solo i file che il gruppo non ha; i file diversi restano come li ha il gruppo e il piano indica `group copies kept`. Per `plugins/`, i registri `installed_plugins.json` e `known_marketplaces.json` vengono uniti, e i percorsi assoluti dei plugin vengono riscritti in modo che puntino alla cartella del gruppo.
- **Conflict**: un file presente da entrambe le parti con contenuto diverso. Scegli per ogni file: `keep group`, `use mine`, più `merge JSON` per i file `.json` (un'unione profonda in cui, sui valori in conflitto, prevale il gruppo) o `merge by line` per `history.jsonl`. Anche un server MCP la cui definizione è diversa è un conflitto, con `keep group` o `use mine`. La prima scelta elencata è quella predefinita quando premi Invio; `--yes` mantiene sempre la copia del gruppo.

La checklist stessa è `tui.MigratePicker`, costruita sul componente `checklist` usato da `ccam add`; `migrate.Found` indica gli elementi mancanti con chiamate lstat e l'elenco dei server MCP.

### Accessi MCP

I server MCP che accedono con OAuth conservano i loro token nel file di login, `.credentials.json` sotto `mcpOAuth`, non nell'elenco dei server. Claude Code usa come chiave di ogni voce `<name>|` seguito dai primi 16 caratteri di un hash `sha256` di `type`, `url` e `headers` del server. La chiave non dipende dalla cartella, quindi una voce copiata funziona in un altro account purché la definizione del server sia identica. Per questo `buildLogins` in `internal/migrate/logins.go` copia solo le voci dei server che il gruppo avrà con la definizione dell'origine: nuovi, già uguali, oppure in conflitto quando hai scelto `use mine`. Gli accessi dei server dei plugin (chiavi che iniziano con `plugin:`) vengono copiati quando copi le estensioni e l'origine ha `plugins/`.

<!-- motion:migrate-signins -->

Le voci vanno nel `.credentials.json` di ogni account che le riceve, scritto con `fsutil.AtomicWrite` in modalità `0600`. Non vanno mai nel gruppo. Un account le riceve solo se condivide `mcp` (oppure, per i server dei plugin, `extensions`) tramite il gruppo di destinazione, e vengono aggiunte solo le chiavi che non ha: le voci che ha già, magari più recenti, vengono mantenute. La riga del piano `sign-ins` elenca i server, gli account che li ricevono e quelli saltati. ccam salta un account, indicandone il nome e il motivo, quando:

- è registrato sul posto con `ccam import`, dato che `claude` da solo potrebbe usare il suo file di login;
- è in esecuzione, dato che Claude Code riscrive il suo file di login (controllato di nuovo subito prima della scrittura, segnalato come `started meanwhile`);
- non ha ancora effettuato il login;
- conserva il login nel Portachiavi di macOS, in cui ccam non scrive.

In un account saltato, accedi ai server con `/mcp`. Un account che inizia a condividere i server MCP in seguito riceve gli accessi quando esegui di nuovo `ccam migrate`; il riepilogo lo indica con ``Then run `ccam migrate` again to give them the MCP sign-ins too.``

## Sicurezza e casi limite

- L'origine viene solo letta. Nessun file al suo interno viene spostato, collegato o modificato, e il login di origine non viene mai scritto, nemmeno su macOS, dove potrebbe essere letto dal Portachiavi.
- Prima che un file del gruppo venga sostituito, da `use mine`, da un'unione o da nuovi server MCP in un `mcp.json` esistente, la copia del gruppo viene salvata in `<group>/.ccam-backup/<time>/`. Il nome della cartella è unico per ogni esecuzione, quindi due esecuzioni nello stesso secondo non sovrascrivono i rispettivi backup. Il riepilogo stampa il percorso.
- Con `--yes`, ogni conflitto mantiene la copia del gruppo, e ccam spiega come scegliere invece: esegui lo stesso comando senza `--yes`.
- Se non cambierebbe nulla, ccam stampa `Nothing to copy.`; se ogni conflitto mantiene la copia del gruppo, `Nothing to copy: the group keeps its own copies.`
- Se una copia fallisce a metà, ccam si ferma con `The source folder was not changed. Group files replaced so far are saved in` seguito dalla cartella di backup.
- Un'origine mancante termina con `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Pagine correlate

- [Condividere server MCP, impostazioni e cronologia](/docs/sharing)
- [Controllo e riparazione](/features/doctor)
- [Più account Claude Code](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
