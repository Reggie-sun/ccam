---
title: Condivisione con i gruppi
description: Un gruppo di ccam condivide gli elementi scelti tra i suoi account tramite una cartella; ccam group sposta i dati con un piano, backup e scelte sui conflitti.
updated: 2026-10-26
order: 4
translation: machine
sourceHash: 3f94dd13dbf1
---

Ogni account ha la propria cartella, ma la maggior parte del suo contenuto non ha bisogno di essere separata. Un **gruppo** dice cosa viene condiviso e da chi: ogni account del gruppo collega quegli elementi a un'unica cartella, `~/.local/share/ccam/groups/<group>/`, così tutti leggono e scrivono la stessa copia. I token di login non vengono mai condivisi.

Un account può stare in più gruppi, ma ogni elemento proviene da un solo gruppo. Una configurazione tipica è un gruppo per le impostazioni e le estensioni, e un secondo per la cronologia delle conversazioni tra due account:

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

Il primo account che aggiungi crea `default`, e gli account successivi vi si uniscono, a meno che tu non indichi altro.

## Utilizzo

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # all accounts in default start sharing history
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # the group starts from work's copies
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # first leave the group that shares those items now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

Ogni comando che sposta dati stampa un piano, chiede cosa fare per ogni file diverso dalla copia del gruppo e chiede `Proceed? [y/N]` prima di modificare qualsiasi cosa:

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261023-213057
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/work → groups/default
  Note: projects: 1 file differs from the group's copy (-home-me-code-app/s1.jsonl); the group's is kept, the account's stays in profiles/work/.ccam-backup/20261023-213057/projects
Proceed? [y/N]
```

Puoi fare le stesse modifiche nella scheda **Groups** della [dashboard](/features/dashboard); premi `s` su un account per sceglierne i gruppi.

## Elementi

| Elemento | Cosa comprende |
| --- | --- |
| `mcp` | i server MCP, tramite il `mcp.json` del gruppo ([Server MCP condivisi](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | il `CLAUDE.md` globale |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` e `history.jsonl`, così un altro account può usare `--continue` |

Se nessuno dei gruppi di un account condivide un elemento, quell'elemento resta nella cartella propria dell'account.

## Comandi e opzioni

| Comando | Cosa fa |
| --- | --- |
| `ccam group` o `ccam group list` | elenca i gruppi |
| `ccam group show <group>` | cosa condivide il gruppo, i suoi account, la sua cartella e i server MCP |
| `ccam group add <group> --items <items> [--members <accounts>]` | crea un gruppo; gli elementi dei suoi account vi vengono spostati |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | cambia cosa condivide il gruppo, per ogni suo account |
| `ccam group join <group> <account>...` | mette degli account nel gruppo |
| `ccam group leave <group> <account>...` | toglie degli account dal gruppo; ognuno ne tiene una copia, compresi i server MCP |
| `ccam group rename <group> <new-name>` | rinomina solo nella configurazione |
| `ccam group remove <group> [--leave]` | rimuove un gruppo |

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: l'account la cui copia entra per prima nel gruppo | l'account predefinito, poi per numero |
| `--replace` | `join`: lascia prima i gruppi che condividono gli stessi elementi | disattivata |
| `--fresh` | gli account che smettono di condividere un elemento partono senza di esso invece che con una copia (una cartella importata ne tiene comunque una) | disattivata |
| `--leave` | `remove`: toglie prima gli account del gruppo | necessaria quando il gruppo ha account |
| `--prefer <group\|mine\|merge>` | come decidere i conflitti | chiede; con `--yes`, il valore predefinito di ogni file |
| `-y`, `--yes` | esegue senza chiedere; i conflitti seguono `--prefer` | disattivata |
| `--dry-run` (globale) | stampa il piano, non modifica nulla | disattivata |

Gli account si indicano per nome o per numero; gli elementi come elenco separato da virgole. `ccam add` accetta `--group <group>` (ripetibile) oppure `--no-group`.

## Come funziona

I file coinvolti:

- `~/.config/ccam/config.toml`: le voci `[[group]]`: `name`, `items`, `members` e, quando il nome della cartella è diverso, `dir`.
- `~/.local/share/ccam/profiles/<name>/`: la `CLAUDE_CONFIG_DIR` dell'account.
- `~/.local/share/ccam/groups/<folder>/`: le copie condivise.

Quando `config.toml` viene letto, i gruppi di ogni account vengono risolti in una cartella per elemento. Ci vogliono circa 25 µs per cinque account, quindi un avvio continua a leggere solo `config.toml`. Una configurazione in cui un account riceverebbe un elemento da due gruppi è un errore, e non si avvia nulla finché non viene corretta.

Una modifica lavora su una copia della configurazione:

1. Il comando modifica la copia (`AddGroup`, `SetGroupItems`, `JoinGroup` …), che controlla le regole.
2. ccam confronta la vecchia e la nuova configurazione e, per ogni account i cui elementi cambiano, costruisce un piano (`internal/group`, `share.Build`). I piani vengono costruiti uno dopo l'altro, e ognuno vede cosa i precedenti metteranno nelle cartelle dei gruppi, così solo il primo account vi sposta la propria copia. La prima copia del gruppo viene, in quest'ordine: dalla cartella del gruppo (se ne ha già una), dall'account `--from`, dall'account predefinito, poi dagli altri account per numero.
3. ccam si rifiuta se uno di quegli account è in esecuzione, comunque sia stato avviato Claude Code (``work is running; quit Claude Code there first, then change the groups (`ccam agents` lists the sessions)``), e su Windows senza la Modalità sviluppatore se un file richiederebbe un link.
4. ccam risolve i conflitti (chiedendo, o con `--prefer`), stampa il piano e chiede conferma.
5. Subito prima di salvare, ccam controlla di nuovo. Se un account è stato avviato mentre ccam aspettava le tue risposte, la modifica si ferma; si ferma anche se nel frattempo un altro processo ccam ha cambiato gli stessi gruppi (`the groups changed while this plan waited; nothing was moved, make the change again`). La nuova configurazione viene salvata prima che qualcosa venga spostato, e la modifica viene annotata in `~/.local/state/ccam/group-change.json`. Se un passaggio si ferma a metà, la nota resta: `ccam doctor --fix` completa la modifica nel suo ordine, con le risposte che hai dato, compresi i passaggi MCP, e fino ad allora non parte nessun'altra modifica dei gruppi (``the last group change stopped halfway; finish it with `ccam doctor --fix` first``).
6. I passaggi MCP vengono eseguiti per primi. Quando degli account tolgono le estensioni da un gruppo, gli elenchi dei plugin del gruppo vengono fatti puntare alla cartella del gruppo ovunque i loro percorsi passino per la cartella di un account (`group plugin lists:` nel piano). Poi i piani vengono eseguiti in ordine. Ogni passaggio sposta i dati nel gruppo o nella cartella di backup; solo i file identici alla copia del gruppo vengono eliminati invece di finire nel backup.

ccam tiene un elenco dei percorsi consentiti coperti da ogni elemento (`share.Layout`). Ciò che non vi compare non può mai arrivare a un gruppo. ccam indica anche `.credentials.json` e `.claude.json` come mai condivisi (`share.Never`), e si ferma con un errore interno se uno dei due compare nella struttura (`Build`). L'elemento `mcp` non ha percorsi: l'elenco dei server si trova dentro `.claude.json`, quindi la [condivisione degli MCP](/features/mcp) funziona tramite il `mcp.json` del gruppo.

Ecco cosa succede a una cartella, `projects/`, quando `side` inizia a condividere la cronologia e il gruppo ha già dei contenuti:

<!-- motion:share-on -->

Se il gruppo non ha ancora nulla in quel percorso, l'account è il primo a condividerlo e la sua cartella viene spostata lì così com'è. Altrimenti ogni file che il gruppo non ha viene spostato nel gruppo, e i file presenti da entrambe le parti restano nella cartella dell'account, così viene mantenuta la versione del gruppo. I file identici a quelli del gruppo vengono eliminati, e il piano indica quelli che differiscono (`Note: projects: 1 file differs from the group's copy (…); the group's is kept, the account's stays in …`).

La cartella dell'account va poi in `.ccam-backup/<time>/` dentro la cartella dell'account, e un link prende il suo posto. Se la voce del gruppo è a sua volta un link alla voce propria dell'account che si unisce, il contenuto dell'account entra nel gruppo al suo posto. I file delle conversazioni prendono il nome dal loro ID di sessione, quindi le cronologie di più account si uniscono senza conflitti. I registri di `plugins/` salvano percorsi assoluti che passano per la cartella dell'account; ccam li fa puntare alla cartella del gruppo e unisce le voci dell'account in un registro che il gruppo ha già.

Uscire da un gruppo inverte il link:

<!-- motion:share-off -->

ccam copia prima il contenuto attuale del gruppo in `.ccam-backup/incomplete-<name>` nella cartella dell'account, e sostituisce il link solo quando la copia è completa. Se la copia si interrompe, il link resta al suo posto e `ccam doctor --fix` copia di nuovo (il piano di un completamento del genere mostra `restore <entry>: … (copied back from the group before a stop)`). Il gruppo mantiene il suo contenuto per gli altri account.

- ccam rimuove il link solo se punta ancora all'interno della cartella dei gruppi come quando è stato creato il piano. Un percorso cambiato nel frattempo, come un file normale dove c'era il link, viene lasciato com'è e la modifica si ferma.
- Con `--fresh`, la copia viene saltata.
- Un account che sposta un elemento in un gruppo la cui cartella è ancora vuota porta con sé una copia di ciò che vedeva, così mantiene i suoi file. Quella copia viene preparata in `<group>/.ccam-incomplete-<name>` finché non è completa.
- Un account che smette di condividere gli MCP mantiene i server del gruppo: finiscono nel suo `.claude.json` (viene conservato un backup), tranne dove l'account ha già un server con lo stesso nome, che resta com'è. Con `--fresh` il piano elenca invece i server che smettono di caricarsi. Unirsi di nuovo rimuove le copie che nessuno ha cambiato.
- Lasciare `extensions` fa anche puntare gli elenchi dei plugin del gruppo alla cartella del gruppo dove i loro percorsi passavano per la cartella dell'account che esce; per gli altri account, i percorsi riscritti portano agli stessi plugin.

Un singolo file presente da entrambe le parti con contenuto diverso è un conflitto:

<!-- motion:share-conflict -->

Per i singoli file come `settings.json` scegli tu; dentro le cartelle resta la copia del gruppo, e il piano indica i file che differiscono. Le scelte dipendono dal file; la prima è quella predefinita:

| File | Scelte |
| --- | --- |
| `history.jsonl` | `merge by line`, `keep group`, `use mine` |
| Altri file `.json` | `keep group`, `use mine`, `merge JSON` |
| Tutto il resto | `keep group`, `use mine` |

`merge JSON` unisce gli oggetti in modo ricorsivo, unisce gli array senza duplicati e mantiene il valore del gruppo quando le due parti impostano la stessa chiave in modo diverso, così gli altri account non notano cambiamenti. Ogni scelta tranne `keep group` salva prima la copia del gruppo in `.ccam-backup/<time>/group/` dell'account. Il tuo file va sempre nel backup.

A ogni avvio, ccam confronta ogni percorso della struttura con i gruppi dell'account usando solo `lstat` e `readlink` (`share.Check`). Un link mancante viene creato sul momento. Un file normale dove dovrebbe esserci un link, un link verso il posto sbagliato o un link obsoleto interrompono l'avvio e rimandano a `ccam doctor --fix`.

### Rinominare, rimuovere e il vecchio formato

- `ccam group rename` cambia solo la configurazione. La cartella mantiene il suo nome e il gruppo lo registra come `dir`, quindi non si sposta nulla e non si può perdere nulla.
- `ccam group remove` rifiuta un gruppo con account, a meno che tu non aggiunga `--leave`. In quel caso ogni account tiene una copia, e la cartella va in `~/.local/share/ccam/trash/`, a meno che un altro gruppo non la usi.
- Prima dei gruppi, `config.toml` salvava una tabella `[profile.share]` per ogni account. ccam legge quel formato come gruppi senza spostare nulla. Al primo salvataggio scrive il nuovo formato e conserva il vecchio file come `config.toml.bak`. Gli elementi condivisi dagli stessi account diventano un unico gruppo che prende il nome dalla cartella; gli altri ricevono nomi come `default-history` e continuano a usare la stessa cartella.

### Su Windows

Windows crea link simbolici solo in Modalità sviluppatore (Impostazioni > Sistema > Per sviluppatori) o per un amministratore. Senza quel diritto, una cartella condivisa diventa una giunzione di directory, che per Claude Code funziona allo stesso modo. Un file condiviso (`settings.json`, `CLAUDE.md`, `history.jsonl`) può essere solo un link simbolico, quindi una modifica di gruppo che ne richiede uno viene rifiutata prima che si sposti qualcosa. Per lo stesso motivo il primo `ccam add` lì crea il gruppo `default` con i soli `mcp` ed `extensions` e lo segnala, e `ccam add --group` verso un gruppo che condivide impostazioni, memoria o cronologia viene rifiutato.

## Sicurezza e casi limite

- **Account in esecuzione.** Una modifica viene rifiutata mentre Claude Code è in esecuzione con uno degli account che sposta, perché la sessione continua a scrivere in quei file, che l'abbia avviata ccam o no (`claude` da solo, un IDE, le sessioni in background). Il controllo viene ripetuto subito prima di salvare. Uscire da un gruppo o unirsi a uno riguarda solo l'account indicato; lasciare `extensions` riscrive anche i percorsi negli elenchi dei plugin del gruppo, senza cambiarne il significato per gli altri. `--dry-run` stampa comunque il piano.
- **Backup.** Ogni piano ha la propria cartella `.ccam-backup/<time>/` nella cartella dell'account. ccam non li elimina mai.
- **I tuoi link.** Se un percorso è già un link creato da te, per esempio verso un repository di dotfile, e il gruppo è vuoto, ccam copia nel gruppo la destinazione del link e sposta il link stesso nel backup. La destinazione non viene toccata.
- **Cartelle importate.** Un account registrato dove si trova con `ccam import` (come `~/.claude`) cambia sul posto, quindi anche `claude` da solo vede la modifica; il piano lo indica. Quando lascia un gruppo ne tiene sempre una copia, anche con `--fresh`, e `ccam remove` lo toglie dai suoi gruppi in quel modo prima di dimenticare l'account. Finché fa parte di un gruppo, i suoi file condivisi stanno nella cartella del gruppo, quindi rimuovi l'account (o lascia i gruppi) prima di eliminare `~/.local/share/ccam`.
- **Due account nella stessa cartella.** Gli account che condividono la cronologia scrivono negli stessi file. Quando un secondo account si avvia in una cartella in cui un altro è già in esecuzione, ccam avvisa e chiede conferma.
- **Modifiche manuali.** Se `members` elenca un nome che non è un account, ccam lo ignora e `ccam doctor` lo segnala. I nomi `local`, `on`, `off` e `none` non possono essere nomi di gruppi.

## Pagine correlate

- [La dashboard](/features/dashboard)
- [Condividere server MCP, impostazioni e cronologia](/docs/sharing)
- [Server MCP condivisi](/features/mcp)
- [Copiare una configurazione esistente](/features/migrate)
- [Riferimento dei comandi](/docs/commands)
