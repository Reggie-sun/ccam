---
title: Copiare una configurazione esistente
description: Copia impostazioni, CLAUDE.md, estensioni e server MCP di una ~/.claude esistente in un gruppo di ccam con ccam migrate.
updated: 2026-10-26
order: 6
translation: machine
sourceHash: f835fbe84ade
---

Se usavi Claude Code prima di ccam, la tua configurazione si trova in `~/.claude` e `~/.claude.json`. `ccam migrate` la copia in un gruppo di condivisione, così ogni account che condivide quegli elementi tramite il gruppo parte con la configurazione che hai già. La cartella di origine viene solo letta, e `claude` da solo continua a funzionare come prima. Eseguilo una volta dopo il tuo primo `ccam add`; quando `~/.claude.json` ha dei server MCP, il primo `ccam add` li indica nel suo suggerimento.

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
 Choose what to copy (Space to toggle)
 ❯ [✔] MCP servers
   [✔] settings.json
   [✔] global CLAUDE.md
   [✔] skills, agents, commands, plugins
   [ ] conversation history

 space to toggle · enter to show the plan · esc to cancel
```

Spazio (o `x`) attiva o disattiva l'elemento sotto il cursore, i tasti freccia o `j`/`k` spostano la selezione, Invio mostra il piano, ed Esc, `q` o Ctrl+C stampano `Cancelled.`. Gli elementi per cui l'origine non contiene nulla sono indicati con `not found`. Se deselezioni tutto, ccam stampa `Nothing selected.` e si ferma.

Poi ccam stampa il piano e chiede conferma prima di copiare:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

Un gruppo riguarda solo i propri account, e solo gli elementi che condivide, quindi il riepilogo elenca gli account del gruppo. Quando il gruppo non condivide uno degli elementi copiati, o alcuni account non ne fanno parte, ccam stampa i comandi `ccam group set <group> --add <items>` e `ccam group join <group> <accounts>` che lo cambiano.

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `[folder]` | La cartella di Claude Code da cui copiare | `~/.claude` |
| `--settings` | Copia `settings.json` | disattivata |
| `--memory` | Copia il `CLAUDE.md` globale | disattivata |
| `--extensions` | Copia `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | disattivata |
| `--mcp` | Copia i server MCP utente (`mcpServers` in `.claude.json`) | disattivata |
| `--history` | Copia la cronologia delle conversazioni (`projects/`, `history.jsonl`) | disattivata |
| `-g`, `--group <group>` | Il gruppo di condivisione in cui copiare | `default` |
| `-y`, `--yes` | Copia senza chiedere; in caso di conflitto viene mantenuta la copia del gruppo | disattivata |
| `--dry-run` (globale) | Stampa il piano e si ferma prima di qualsiasi domanda o copia | disattivata |

In un terminale, se non indichi opzioni per gli elementi, scegli gli elementi nella checklist. Fuori da un terminale, o con `--yes`, ccam copia tutto tranne la cronologia delle conversazioni. La cronologia è esclusa per impostazione predefinita perché `projects/` può essere grande.

`--group on` indica `default`. `--group off` e `--group local` vengono rifiutati con `--group takes a group name, not "off"`. `ccam migrate` supporta `--dry-run`, e la checklist compare anche in un dry run, così puoi vedere in anteprima una selezione.

## Come funziona

ccam costruisce un piano completo prima di copiare qualsiasi cosa. Sceglie gli elementi (`internal/cli/migrate.go`), poi legge l'origine e il gruppo e costruisce il piano (`migrate.Build`). Non viene scritto nulla fino a `Plan.Apply`, dopo la tua conferma.

1. **Trovare l'origine.** ccam rifiuta un percorso mancante, che non è una cartella o che non contiene nessuna voce nota di Claude Code (`checkSource`). Per i soli server MCP, una `~/.claude` vuota o mancante va bene purché `~/.claude.json` abbia dei server. ccam trova il `.claude.json` che appartiene alla cartella (`migrate.ClaudeJSON`): al suo interno oppure, per una `~/.claude` che non ne contiene uno, `~/.claude.json` accanto a essa. È lì che Claude Code lo tiene quando `CLAUDE_CONFIG_DIR` non è impostata.
2. **Associare gli elementi ai percorsi.** Ogni percorso viene da `share.Layout`, lo stesso elenco di percorsi consentiti usato dai gruppi. I server MCP vengono letti da `.claude.json` con `mcp.LocalServers`, che usa lo scanner `profile.TopLevelValue` invece di analizzare l'intero file.
3. **Copiare, mai spostare.** I file vanno in `~/.local/share/ccam/groups/<group>/` con `fsutil.AtomicWrite` o `fsutil.CopyFile`, che scrive un file temporaneo accanto alla destinazione e lo rinomina, così una copia interrotta non lascia mai un file parziale. I server MCP vanno nel `mcp.json` del gruppo sotto il suo lock: un server che un altro processo ccam aggiunge nel frattempo viene mantenuto, e le copie fatte per [Remote Control](/guides/remote-control) vengono aggiornate dopo.

Ogni elemento associa un percorso di origine a un percorso del gruppo:

<!-- motion:migrate-items -->

Le due righe della cronologia sono tratteggiate perché vengono copiate solo con `--history`. `.credentials.json` e `.claude.json` sono elencati in `share.Never` e non arrivano mai al gruppo: ccam si ferma con un errore interno se uno dei due compare nella struttura (`Build`), e la scansione della cartella li salta per nome.

ccam assegna a ogni percorso uno stato, che è ciò che descrive ogni riga del piano (`buildEntry`):

<!-- motion:migrate-states -->

- **Missing**: assente nell'origine, oppure un file dove ci si aspetta una cartella. Indicato una volta per elemento come `not found, skipped`.
- **Linked**: il percorso di origine porta già alla copia del gruppo, per esempio quando l'origine è la cartella di un profilo di ccam.
- **Same**: il gruppo ha già un contenuto identico.
- **New**: il gruppo non ce l'ha, quindi viene copiato per intero.
- **Merge**: una cartella presente da entrambe le parti. ccam aggiunge solo i file che il gruppo non ha (`copyMissing`); i file diversi restano come li ha il gruppo e il piano indica `group copies kept`. Per `plugins/`, i registri `installed_plugins.json` e `known_marketplaces.json` vengono uniti, e i percorsi assoluti dei plugin vengono riscritti in modo che puntino alla cartella del gruppo.
- **Conflict**: un file presente da entrambe le parti con contenuto diverso. Scegli per ogni file: `keep group`, `use mine`, più `merge JSON` per i file `.json` (un'unione profonda in cui, sui valori in conflitto, prevale il gruppo) o `merge by line` per `history.jsonl`. Anche un server MCP la cui definizione è diversa è un conflitto, con `keep group` o `use mine`. La prima scelta elencata è quella predefinita quando premi Invio; `--yes` mantiene sempre la copia del gruppo.

La checklist stessa è `tui.MigratePicker`, costruita sul componente `checklist` usato da `ccam add`; ccam indica gli elementi mancanti con chiamate lstat e l'elenco dei server MCP (`migrate.Found`).

### Accessi MCP

I server MCP che accedono con OAuth conservano i loro token nel file di login, `.credentials.json` sotto `mcpOAuth`, non nell'elenco dei server. ccam non li copia. Un accesso copiato contiene lo stesso refresh token dell'originale, e il server ruota quel token al primo rinnovo: tutti gli altri che lo possiedono, compresa la `~/.claude` di origine, vengono allora disconnessi da quel server. Per questo, per ogni server copiato a cui l'origine ha effettuato l'accesso, la riga `sign-ins` del piano elenca gli account del gruppo che devono accedere da sé con `/mcp`, e il riepilogo lo ripete. Una voce senza token non è un accesso e non viene elencata, e nemmeno un account che ha già un proprio accesso o un server di cui il gruppo mantiene la definizione. Non viene scritto alcun file di login.

## Sicurezza e casi limite

- L'origine viene solo letta. Nessun file al suo interno viene spostato, collegato o modificato, e il login di origine non viene mai scritto, nemmeno su macOS, dove potrebbe essere letto dal Portachiavi.
- Prima che un file del gruppo venga sostituito, da `use mine`, da un'unione o da nuovi server MCP in un `mcp.json` esistente, la copia del gruppo viene salvata in `<group>/.ccam-backup/<time>/`. Il nome della cartella è unico per ogni esecuzione, quindi due esecuzioni nello stesso secondo non sovrascrivono i rispettivi backup. Il riepilogo stampa il percorso.
- Con `--yes`, ogni conflitto mantiene la copia del gruppo, e ccam spiega come scegliere invece: esegui lo stesso comando senza `--yes`.
- Se non cambierebbe nulla, ccam stampa `Nothing to copy.`; se ogni conflitto mantiene la copia del gruppo, `Nothing to copy: the group keeps its own copies.`
- Se una copia fallisce a metà, ccam si ferma con `The source folder was not changed. Group files replaced so far are saved in` seguito dalla cartella di backup. Eseguirlo di nuovo copia i file che non sono stati copiati per intero.
- Quando più gruppi condividono gli MCP, i comandi suggeriti da `ccam doctor`, `ccam ext` e `ccam mcp list` sono `ccam migrate --mcp --group <a|b>`.
- Un'origine mancante termina con `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Pagine correlate

- [Condividere server MCP, impostazioni e cronologia](/docs/sharing)
- [Controllo e riparazione](/features/doctor)
- [Usare più account Claude Code su un solo computer](/guides/multiple-claude-code-accounts)
- [Come funziona CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
