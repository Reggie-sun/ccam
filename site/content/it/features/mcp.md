---
title: Server MCP condivisi
description: Mantieni un unico elenco di server MCP per gruppo con ccam mcp add, add-json, list, remove e import, passato con --mcp-config o sincronizzato in ogni account.
updated: 2026-10-26
order: 5
translation: machine
sourceHash: 27a06ad35295
---

Claude Code salva i server MCP utente in `.claude.json`, lo stesso file che contiene il login dell'account, quindi quel file non può essere condiviso tra account. ccam tiene i server condivisi in un file del gruppo, `mcp.json`, e lo passa all'avvio a ogni account del gruppo. Usa `ccam mcp` per aggiungere un server una sola volta per tutti i tuoi account, e `ccam mcp import` per portare agli altri un server che hai aggiunto all'interno di una sessione.

## Utilizzo

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's gh if it differs
```

Una modifica vale dal successivo avvio di ogni account; le sessioni in esecuzione mantengono i server con cui sono partite:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
only in account side, for /home/me/code/app only: notes
```

Sotto i gruppi, `ccam mcp list` indica ogni server che non arriva agli altri account:

- `only in account work: pg` per un server che un account ha aggiunto per sé;
- `only in account side, for /home/me/code/app only: notes` per uno aggiunto per una sola cartella (il semplice `claude mcp add`, vedi sotto);
- `account work has its own …, used in its sessions in place of the group's` quando un server proprio di un account ha il nome di un server del gruppo;
- il `.mcp.json` di questa cartella, con gli account la cui approvazione è `waiting for approval in`;
- i server che sono `in no account, only in ~/.claude.json (plain claude's)`, con il comando [`ccam migrate --mcp`](/features/migrate) che li copia.

Aggiungere un nome che il gruppo ha già ne sostituisce la definizione e stampa `updated` invece di `added`. I server arrivano a un account solo se fa parte di quel gruppo e il gruppo condivide gli MCP (`ccam group set <group> --add mcp`); quando il gruppo non li condivide, `add` lo segnala.

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `-g`, `--group <group>` | gruppo di destinazione per `add`, `add-json` e `remove` | l'unico gruppo tramite cui qualche account condivide gli MCP, altrimenti `default` |
| `-- <command> [args…]` (`add`) | esegue questo comando come server stdio | |
| `--url <url>` (`add`) | URL del server remoto (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | trasporto di un server `--url` | `http` |
| `--header <'Name: value'>` (`add`) | intestazione HTTP, ripetibile | nessuna |
| `-e`, `--env <KEY=VALUE>` (`add`) | variabile d'ambiente per un server stdio, ripetibile | nessuna |
| `--force` (`import`) | sostituisce un server del gruppo con lo stesso nome ma diverso; del vecchio `mcp.json` viene fatto un backup | disattivata |
| `--keep-local` (`import`) | mantiene le voci nel `.claude.json` dell'account | disattivata |
| `-h`, `--help` | aiuto per ogni comando | |
| `--dry-run` (globale) | supportata solo da `mcp list`; `add`, `add-json`, `remove` e `import` la rifiutano e non modificano nulla | disattivata |

`add` richiede un comando dopo `--` oppure `--url`; altrimenti fallisce con `a command to run (after --) or --url is required`. `--header` vale per i server `--url` e `--env` per i server a comando. `add-json` accetta lo stesso JSON di `claude mcp add-json`, e la definizione deve essere un oggetto con un `command` o un `url`. `import` usa sempre il gruppo tramite cui l'account condivide gli MCP, e `list` mostra tutti i gruppi, quindi `-g` non li modifica.

## Come funziona

I file coinvolti:

- `~/.local/share/ccam/groups/<group>/mcp.json`: i server condivisi, `{"mcpServers": {…}}`, il formato letto da `--mcp-config`. Viene scritto con `fsutil.AtomicWrite` e modalità `0600`, perché le voci dei server contengono spesso chiavi API, e solo sotto il lock `<group>/.mcp.json.lock`, così due comandi `ccam mcp add` eseguiti contemporaneamente mantengono entrambi i loro server.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: il file di Claude Code con i server dell'account. ccam legge solo il suo valore `mcpServers` con lo scanner `profile.TopLevelValue` e sul percorso di avvio non analizza mai l'intero file.
- `~/.local/state/ccam/mcp/<name>/`: il file `seen`, il `sync-stamp` (con un hash di ogni copia scritta da ccam, così distingue una copia che nessuno ha cambiato da una che hai modificato), `sync-checked` (ciò che la strategia sync ha visto l'ultima volta), `kept` (i server che l'account ha mantenuto quando ha lasciato la condivisione MCP) e i backup di `.claude.json`.

Ci sono due strategie, impostate da `mcp_strategy` in `~/.config/ccam/config.toml`. Quella predefinita è la strategia flag:

<!-- motion:mcp-flag -->

1. ccam trova il gruppo MCP dell'account (`launcher.applyMCP`) e verifica con un solo `stat` che il `mcp.json` del gruppo esista, senza analizzarlo (`mcp.FlagArgs`).
2. ccam mette `--mcp-config=<path>` subito dopo `claude`, prima dei tuoi argomenti, così le tue opzioni `--mcp-config` o `--strict-mcp-config` continuano a valere in aggiunta (`Plan.InsertArgs`).
3. ccam controlla se il `.claude.json` dell'account è cambiato dall'ultimo avvio, e solo in quel caso ne analizza i server (`mcp.NewLocal`).
4. ccam si sostituisce con Claude Code (`syscall.Exec`).

Il percorso viene passato come un solo argomento, `--mcp-config=<path>`, non come due. L'opzione `--mcp-config` di Claude Code accetta più valori, quindi `claude --mcp-config a.json "prompt"` leggerebbe il prompt come secondo file di configurazione e fallirebbe con `MCP config file not found`. La forma con un solo argomento evita il problema, quindi un prompt o un sottocomando che la segue funziona. `ccam work mcp list`, `get` e `remove` non ricevono l'opzione, perché lì Claude Code la rifiuta; ccam indica invece i server del gruppo sullo standard error ([dettagli](/features/launch#shared-mcp-servers)). Claude Code carica i server di `--mcp-config` in aggiunta a quelli dell'account; in caso di nomi uguali prevale la voce in `.claude.json`. La strategia flag scrive `.claude.json` solo per Remote Control, quando un account lascia la condivisione MCP o vi rientra ([sotto](#leaving-and-joining)), e per `ccam mcp import`.

`claude remote-control` (`rc`) non accetta `--mcp-config`, e le sue sessioni leggono i server solo da `.claude.json`. Per `ccam work rc`, ccam omette l'opzione e, in un terminale, chiede se copiare i server del gruppo in `.claude.json`, con un backup, come fa la strategia `sync`. Chiede finché non esiste una copia, poi di nuovo solo dopo che il `mcp.json` del gruppo cambia o un server copiato è sparito ([Remote Control](/features/launch#remote-control)). Poiché una copia del genere prevale su `--mcp-config` in caso di nomi uguali, `ccam mcp add`, `add-json`, `remove` e `import` aggiornano le copie che hanno gli altri account del gruppo:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

La riga `✓ Updated` compare solo quando una copia è cambiata. Una copia creata da una versione più vecchia di ccam viene aggiornata allo stesso modo, e così anche una dopo `ccam migrate`.

Un server aggiunto all'interno di una sessione con `claude mcp add` finisce solo nel `.claude.json` di quell'account. Il semplice `claude mcp add` usa l'ambito locale: il server si carica per quel solo account in quella sola cartella (`-s user` lo rende disponibile all'account in ogni cartella). ccam non può vederlo nel momento in cui succede, perché Claude Code ha già sostituito il processo di ccam, quindi controlla all'avvio successivo. Il file `seen` contiene l'ora di modifica di `.claude.json` e i nomi dei server già segnalati. Quando l'ora è diversa, ccam confronta i nomi con il gruppo e stampa una sola volta un suggerimento per ogni nuovo server, senza attendere input, con entrambe le strategie:

```console
ccam: MCP servers that only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers that account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

Il suggerimento indica anche un server proprio di un account che ha il nome di un server del gruppo, perché nelle sue sessioni prevale quello dell'account.

Con `mcp_strategy = "sync"`, i server compaiono in `claude mcp list` e `/mcp` come normali server utente, e ccam li scrive nel `.claude.json` di ogni account prima dell'avvio:

<!-- motion:mcp-sync -->

1. ccam confronta l'ora di modifica del `mcp.json` del gruppo con `sync-stamp`, e il `.claude.json` dell'account con ciò che ha visto l'ultima volta (`mcp.Sync`). Se nessuno dei due è cambiato, non succede nient'altro. Altrimenti controlla anche che ogni server del gruppo sia ancora in `.claude.json`, e rimette una copia sparita (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`).
2. Prende il lock di Claude Code stesso, la cartella `.claude.json.lock` accanto al file, creata con `mkdir`. Attende fino a 2 secondi e considera scaduto un lock più vecchio di 10 secondi, come fa la libreria di lock di Claude Code.
3. Salva il `.claude.json` attuale in `backups/` nella cartella di stato e conserva i 5 più recenti.
4. Imposta ogni server del gruppo in `mcpServers`, rimuove i server che aveva scritto in precedenza e che il gruppo non ha più, e scrive il file in modo atomico con la sua modalità originale, solo quando qualcosa è cambiato. Modifica o rimuove solo le copie che ha scritto e che nessuno ha cambiato da allora. Se l'account ha un proprio server con lo stesso nome di un server del gruppo, ccam mantiene il server dell'account e lo segnala una volta: `account work has its own MCP servers fs, so group default's servers with those names are not copied there`.
5. Registra il nuovo stamp con un hash di ogni copia, e Claude Code si avvia senza `--mcp-config`.

La sincronizzazione viene saltata, con `ccam: account work is already running, skipping MCP sync`, quando è in esecuzione una sessione dell'account, comunque sia stato avviato Claude Code, perché Claude Code riscrive `.claude.json` mentre è in esecuzione. Una sincronizzazione fallita stampa `ccam: MCP sync failed (launching anyway): …` e avvia comunque.

### Uscire e rientrare

Un account che smette di condividere gli MCP (`ccam group leave`, `ccam group set --remove mcp`, `ccam group remove --leave`) mantiene i server del gruppo: il piano li copia nel suo `.claude.json`, con un backup, e i suoi server con gli stessi nomi restano. Con `--fresh` non ne riceve nessuno, e il piano elenca i server che smettono di caricarsi. Quando si unisce di nuovo a un gruppo che ha quei server, le copie mantenute che nessuno ha cambiato vengono rimosse, così le modifiche successive del gruppo gli arrivano.

`ccam mcp import` sposta i server propri di un account nel suo gruppo:

<!-- motion:mcp-import -->

ccam legge i server dell'account, quelli propri e quelli che ha per una sola cartella, e confronta ciascuno con il server del gruppo con lo stesso nome usando `mcp.SameDef`, che ignora l'ordine delle chiavi e gli spazi. Un server di cartella diventa un server del gruppo per ogni cartella (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`); un nome con due definizioni di cartella diverse viene rifiutato. I server dei plugin non vengono proposti, perché li porta il plugin. Un server che il gruppo non ha, o che ha con la stessa definizione, viene copiato. Uno diverso viene saltato, perché la definizione del gruppo è quella usata da tutti gli altri account:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

Con `--force` il server del gruppo viene sostituito, e il `mcp.json` precedente viene prima salvato in `.ccam-backup/<time>/mcp.json` dentro la cartella del gruppo. Poi ccam rimuove le voci copiate dal `.claude.json` dell'account, quelle proprie e quelle della cartella in un'unica scrittura, con lo stesso lock e dopo un backup nella cartella di stato (`mcp.Remove`). ccam tiene bloccato il `mcp.json` del gruppo dal momento in cui legge il file finché quelle voci non sono state rimosse, così nessuno toglie il server dal gruppo nel frattempo. Serve perché, in caso di nomi uguali, la voce dell'account prevarrebbe su quella del gruppo. Con `--keep-local`, le voci restano nel `.claude.json` dell'account.

## Sicurezza e casi limite

- **Accessi.** I server che accedono con OAuth conservano i loro token nel `.credentials.json` di ogni account (`mcpOAuth`), non nell'elenco dei server, quindi condividere un server non ne condivide l'accesso. Accedi una volta per account con `/mcp`. ccam non copia mai gli accessi, e non lo fa nemmeno [`ccam migrate`](/features/migrate): un refresh token copiato disconnette tutti gli altri che lo possiedono. [`ccam ext`](/features/ext) e `ccam doctor` indicano i server a cui un account non ha effettuato l'accesso mentre un altro account sì, e `ccam next --continue` dice quali mancano al nuovo account. Il gruppo non riceve mai un file di login.
- **Account in esecuzione.** `import` copia nel gruppo ma lascia le voci proprie di un account in esecuzione e lo segnala; eseguilo di nuovo dopo la fine della sessione. La sincronizzazione salta gli account in esecuzione. Entrambi rilevano anche le sessioni di Claude Code avviate senza ccam.
- **Copie per Remote Control.** La copia di un account in esecuzione non viene toccata: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Con `mcp_strategy = "sync"` non c'è niente da aggiornare, perché ogni avvio sincronizza.
- **File danneggiati.** Se `.claude.json` non è un JSON valido, la sincronizzazione non lo tocca e segnala l'errore. Se non riesce a prendere il lock entro 2 secondi, ccam segnala `could not lock … (Claude Code is writing the file)`.
- **Account che non condividono gli MCP.** `import` si rifiuta e ti dice di mettere prima l'account in un gruppo che condivide gli MCP (`ccam group join <group> side`).
- **Nomi dei gruppi.** `-g` accetta solo un nome di gruppo; diventa una cartella sotto `groups/`, quindi nomi come `../x` vengono rifiutati.

## Pagine correlate

- [Condivisione con i gruppi](/features/share)
- [Usare Remote Control con ccam](/guides/remote-control)
- [Condividere server MCP, impostazioni e cronologia](/docs/sharing)
- [Copiare una configurazione esistente](/features/migrate)
- [Riferimento dei comandi](/docs/commands)
