---
title: Condivisione tra account
description: Come ccam share collega impostazioni, CLAUDE.md, estensioni e cronologia a una cartella di gruppo, con piano e backup, e come off ricopia i dati.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 16de29bbdd57
---

Ogni account ha la propria cartella, ma la maggior parte del suo contenuto non ha bisogno di essere separata. `ccam share` attiva o disattiva la condivisione per un account e un elemento alla volta: il file o la cartella dell'account viene spostato in una cartella di gruppo e sostituito da un link simbolico, così ogni account del gruppo legge e scrive la stessa copia. Usalo dopo `ccam add` per cambiare cosa condivide un account, oppure per mettere alcuni account in un gruppo a parte.

## Utilizzo

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

Senza opzioni per gli elementi, `ccam share` stampa lo stato attuale:

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

Con le opzioni, ccam stampa un piano, chiede cosa fare per ogni file diverso dalla copia del gruppo e chiede `Proceed? [y/N]` prima di modificare qualsiasi cosa. Un nuovo account condivide tutto tranne la cronologia delle conversazioni; [Condividere server MCP, impostazioni e cronologia](/docs/sharing) descrive ogni elemento dal punto di vista dell'utente.

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | Server MCP | invariato |
| `--settings <on\|off\|group>` | `settings.json` | invariato |
| `--memory <on\|off\|group>` | il `CLAUDE.md` globale | invariato |
| `--extensions <on\|off\|group>` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | invariato |
| `--history <on\|off\|group>` | `projects/` e `history.jsonl` | invariato |
| `--fresh` | quando disattivi la condivisione, parte da vuoto invece di copiare il contenuto del gruppo | disattivata |
| `--prefer <group\|mine\|merge>` | come decidere i conflitti | chiede; con `--yes`, il valore predefinito di ogni elemento |
| `-y`, `--yes` | esegue senza chiedere; i conflitti seguono `--prefer` | disattivata |
| `-h`, `--help` | aiuto per `share` | |
| `--dry-run` (globale) | rispettata: chiede dei conflitti, stampa il piano, non modifica nulla | disattivata |

Il valore `on` indica il gruppo `default`, `off` indica locale (non condiviso) e qualsiasi altro valore è il nome di un gruppo. I nomi dei gruppi diventano nomi di cartelle sotto `~/.local/share/ccam/groups/`, quindi `ccam share side --history ../x` fallisce con `invalid group name: "../x"`. L'account si indica per nome o per numero, e cambiano solo gli elementi che indichi.

## Come funziona

I file coinvolti:

- `~/.config/ccam/config.toml`: la tabella `[profile.share]` di ogni account, un valore per elemento.
- `~/.local/share/ccam/profiles/<name>/`: la `CLAUDE_CONFIG_DIR` dell'account.
- `~/.local/share/ccam/groups/<group>/`: le copie condivise.

1. `internal/cli/share.go` legge la configurazione, applica ogni valore `--<item>` a una copia delle impostazioni dell'account con `config.Share.Set` e si rifiuta se l'account è in esecuzione.
2. `share.Build` confronta le impostazioni attuali con quelle di destinazione ed elenca i passaggi. Legge solo dal disco: `os.Lstat` e `os.Readlink` per ogni percorso di `share.Layout`. Il piano è un elenco di operazioni: spostare in una posizione vuota del gruppo (adopt), spostare i file che il gruppo non ha (merge dir), scrivere nel gruppo un conflitto risolto, fare un backup, rimuovere un link, creare un link, ricopiare, riscrivere i percorsi dei plugin e unire il registro dei plugin.
3. I conflitti vengono decisi dalla domanda o da `--prefer`.
4. `Plan.Describe` stampa il piano e ccam chiede conferma.
5. Le impostazioni di destinazione vengono scritte in `config.toml` prima che qualcosa venga spostato. Se `Plan.Apply` si ferma a metà, `ccam doctor --fix` e una nuova esecuzione costruiscono il loro piano verso la destinazione salvata e completano la modifica.
6. `Plan.Apply` esegue i passaggi in ordine. Ogni passaggio sposta i dati nel gruppo o nella cartella di backup; non viene eliminato nulla.

`share.Layout` è un elenco dei percorsi consentiti coperti da ogni elemento. Ciò che non vi compare non può mai arrivare a un gruppo. `share.Never` indica `.credentials.json` e `.claude.json`, e `Build` si ferma con un errore interno se uno dei due compare nella struttura. L'elemento `mcp` non ha percorsi: l'elenco dei server si trova dentro `.claude.json`, quindi la [condivisione degli MCP](/features/mcp) funziona invece tramite il `mcp.json` del gruppo, e `--mcp` cambia solo l'impostazione.

Ecco cosa succede a una cartella, `projects/`, quando `side` attiva la condivisione della cronologia e il gruppo ha già dei contenuti:

<!-- motion:share-on -->

Se il gruppo non ha ancora nulla in quel percorso, l'account è il primo a condividerlo e la sua cartella viene spostata lì così com'è, senza backup. Altrimenti `mergeDir` sposta nel gruppo ogni file che il gruppo non ha, e i file presenti da entrambe le parti restano nella cartella dell'account, così viene mantenuta la versione del gruppo. La cartella dell'account va poi in `.ccam-backup/<time>/` dentro la cartella dell'account, e `os.Symlink` mette un link al suo posto. I file delle conversazioni prendono il nome dal loro ID di sessione, quindi le cronologie di più account si uniscono senza conflitti.

Lo spostamento usa `rename(2)` e, tra file system diversi, ripiega su copia ed eliminazione. Un link relativo all'interno dell'albero spostato che puntava al di fuori di esso viene riscritto come link assoluto (`fsutil.PortableLink`), così funziona ancora dalla nuova posizione. I registri di `plugins/` (`installed_plugins.json`, `known_marketplaces.json`) salvano percorsi assoluti che passano per la cartella dell'account; ccam li fa puntare alla cartella del gruppo e, quando il gruppo ha già un registro, vi unisce le voci dell'account.

Disattivare un elemento inverte il link:

<!-- motion:share-off -->

ccam rimuove il link solo se punta all'interno della cartella dei gruppi, poi copia nell'account il contenuto attuale del gruppo con `CopyTree`. Con `--fresh` la copia viene saltata e l'account parte senza nulla in quel percorso. Il gruppo mantiene il suo contenuto per gli altri account. Spostare un elemento da un gruppo a un altro rimuove il link e lo ricrea verso il nuovo gruppo, lasciando il vecchio gruppo com'è; il piano lo indica in una riga `Note:`.

Un singolo file presente da entrambe le parti con contenuto diverso è un conflitto:

<!-- motion:share-conflict -->

`share.NewConflict` elenca le scelte adatte al file. `history.jsonl` offre `merge by line`, `keep group` e `use mine`, con l'unione per prima. Gli altri file `.json` offrono `keep group`, `use mine` e `merge JSON`. Tutto il resto offre `keep group` e `use mine`. La prima scelta è quella predefinita:

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON` (`MergeJSONBytes`) unisce gli oggetti in modo ricorsivo, unisce gli array senza duplicati e mantiene il valore del gruppo quando le due parti impostano la stessa chiave in modo diverso, così gli altri account non notano cambiamenti di comportamento. `merge by line` mantiene le righe del gruppo, poi aggiunge le tue righe che il gruppo non ha. Ogni scelta tranne `keep group` salva prima la copia del gruppo in `.ccam-backup/<time>/group/` nella cartella dell'account, poi scrive il risultato con `fsutil.AtomicWrite`. Il tuo file va sempre nel backup.

`--prefer group`, `mine` o `merge` decide ogni conflitto senza chiedere; `merge` indica `merge JSON` o `merge by line`, a seconda di cosa offre il file, e un file che non offre nessuno dei due mantiene il valore predefinito. `--yes` da solo usa il valore predefinito di ogni file.

A ogni avvio, `share.Check` confronta ogni percorso della struttura con l'impostazione usando solo `lstat` e `readlink`, il che lo mantiene abbastanza leggero per il percorso di avvio. Un link mancante viene creato sul momento, perché non c'è nulla da perdere. Un file normale dove dovrebbe esserci un link, un link verso il posto sbagliato o un link obsoleto interrompono l'avvio:

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix` esegue `share.Build` con le impostazioni salvate da entrambe le parti, il che ricava ogni passaggio da ciò che si trova su disco, mantiene il contenuto del gruppo in caso di conflitti e fa un backup di tutto ciò che sostituisce.

### Su Windows

Windows crea link simbolici solo in Modalità sviluppatore (Impostazioni > Sistema > Per sviluppatori) o per un amministratore. `fsutil.Link` prova prima un link simbolico. Senza quel diritto, una cartella condivisa (`agents`, `skills`, `projects` e le altre) diventa una giunzione di directory, che qualsiasi utente può creare, e per Claude Code funziona allo stesso modo. Un file condiviso (`settings.json`, `CLAUDE.md`, `history.jsonl`) può essere solo un link simbolico, quindi ccam si ferma con `Windows links files only in Developer Mode` e suggerisce di attivarla o di tenere quell'elemento locale. `share.Check` legge le giunzioni come link (`fsutil.LinkTarget`), e rimuovere la cartella di un account rimuove la giunzione stessa, mai il contenuto del gruppo a cui punta.

## Sicurezza e casi limite

- **Account in esecuzione.** ccam non modifica la condivisione mentre Claude Code è in esecuzione con l'account (`account work is running; quit it first, then change its sharing`), perché la sessione continua a scrivere nei file che vengono spostati e le righe aggiunte nel frattempo andrebbero perse. `--dry-run` stampa comunque il piano.
- **Backup.** Ogni piano ha la propria cartella `.ccam-backup/<time>/` nella cartella dell'account, con un suffisso numerato quando due piani iniziano nello stesso secondo. ccam non li elimina mai.
- **Errori a metà.** Se un passaggio fallisce, ccam stampa l'errore e dove si trovano gli originali, e ti dice di completare con `ccam doctor --fix`. La destinazione è già salvata, quindi la riparazione prosegue verso di essa.
- **I tuoi link.** Se un percorso è già un link creato da te, per esempio verso un repository di dotfile, e il gruppo è vuoto, ccam copia nel gruppo la destinazione del link e sposta il link stesso nel backup. La destinazione non viene toccata.
- **File di login.** `.credentials.json` e `.claude.json` non vengono mai collegati, spostati o copiati in un gruppo.
- **Due account nella stessa cartella.** Gli account che condividono la cronologia scrivono negli stessi file. Quando un secondo account dello stesso gruppo di cronologia si avvia in una cartella in cui un altro è già in esecuzione, ccam avvisa e chiede conferma prima di continuare.

## Pagine correlate

- [Condividere server MCP, impostazioni e cronologia](/docs/sharing)
- [Server MCP condivisi](/features/mcp)
- [Copiare una configurazione esistente](/features/migrate)
- [Tutti i comandi](/docs/commands)
