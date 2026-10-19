---
title: Gestire gli account
description: Aggiungi, collega, elenca, rinomina, rimuovi e importa account Claude Code con ccam add, ccam login, ccam list, ccam rename, ccam remove e ccam import.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 99c931574262
---

Ogni account di ccam è una cartella di configurazione di Claude Code separata, con il proprio login. Questi comandi creano gli account, ne effettuano il login, li elencano, li rinominano e li rimuovono, e registrano le cartelle che hai già.

## Utilizzo

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Dove un comando accetta `<name|number>`, puoi indicare il nome dell'account oppure il suo numero di slot, cioè la colonna `#` di `ccam list`.

`ccam add` in un terminale senza opzioni apre una procedura guidata. Chiede il nome, poi mostra una checklist degli elementi condivisi (è tutto selezionato tranne la cronologia delle conversazioni), poi chiede `Create a shortcut command for work? [y/N]` e, se rispondi `y`, `Shortcut name [ccwork]:`. Una qualsiasi tra `--share`, `--shortcut` o `--no-login`, oppure una pipe al posto di un terminale, salta la procedura guidata; in quel caso il nome è obbligatorio e non viene creata nessuna scorciatoia, a meno che `--shortcut` non ne indichi una.

`ccam list` mostra lo slot, il nome, la scorciatoia, l'email, gli elementi condivisi e l'ultimo uso, con `*` davanti all'account predefinito:

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

Con `--json`, ogni account è un oggetto con `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (omesso se l'account non è mai stato avviato) e `share`. `share` associa ogni elemento (`mcp`, `settings`, `memory`, `extensions`, `history`) al nome del suo gruppo oppure a `"local"`.

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `ccam add --no-login` | Crea l'account senza effettuare il login; accedi in seguito con `ccam login <name>` | disattivata |
| `ccam add --share <item=group>` | Condivide un elemento con un gruppo, oppure `item=local` per tenerlo privato. Ripetibile | `mcp`, `settings`, `memory`, `extensions` in `default`; `history` locale |
| `ccam add --shortcut <command>` | Crea anche questo comando scorciatoia in `~/.local/bin` | nessuna |
| `ccam list --json` | Stampa gli account in JSON | disattivata |
| `ccam remove -y`, `--yes` | Rimuove senza chiedere | disattivata |
| `ccam import --name <name>` | Nome dell'account per la cartella indicata | ricavato dal nome della cartella |
| `ccam import -y`, `--yes` | Importa ogni cartella trovata senza chiedere | disattivata |
| `--dry-run` | Rispettata da `ccam list` e `ccam completion`, che non modificano nulla. Rifiutata da `add`, `login`, `logout`, `use`, `rename`, `remove` e `import`, con `` --dry-run is not supported by `ccam add`; nothing was changed `` (con il nome del comando) | disattivata |

`ccam login`, `ccam logout`, `ccam use` e `ccam rename` non hanno opzioni proprie. `ccam completion` accetta `bash` o `zsh`; qualsiasi altro valore fallisce con `only bash and zsh are supported`.

I nomi degli account sono lunghi da 1 a 32 caratteri tra `a-z`, `0-9`, `-` e `_`, e iniziano con una lettera o una cifra. Un nome non può essere composto solo da cifre, perché `ccam 2` indica lo slot 2. Un nome non può nemmeno essere un comando di ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`), così `ccam <name>` non è mai ambiguo. L'elenco è `config.Reserved` in `internal/config`.

## Come funziona

`ccam add work --shortcut ccwork` esegue questi passaggi in `addProfile` (`internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. `config.ValidateName` controlla il nome e `shim.Check` controlla la scorciatoia prima che venga creato qualsiasi cosa.
2. `profile.Add` crea `~/.local/share/ccam/profiles/work/` con modalità `0700` e assegna all'account il numero di slot libero più basso (`Config.NextSlot`), così il numero di un account rimosso viene riutilizzato. Il primo account diventa quello predefinito.
3. `applyShare` collega la nuova cartella vuota ai suoi gruppi di condivisione con `share.Build` e poi `Apply`, per esempio `settings.json` a `groups/default/settings.json`. Se questo passaggio fallisce, la nuova cartella viene rimossa, perché una cartella rimasta bloccherebbe ogni nuovo tentativo con `folder already exists`.
4. L'account viene salvato in `~/.config/ccam/config.toml` tramite `fsutil.AtomicWrite`: un file temporaneo, `fsync`, poi `rename`, così un crash non lascia mai una configurazione scritta a metà.
5. `runLogin` avvia `claude auth login` con `CLAUDE_CONFIG_DIR` impostata sulla cartella dell'account. Se il login fallisce, l'account resta e ccam stampa l'errore e ti dice di eseguire `ccam login work` più tardi.
6. `profile.ReadEmail` estrae `oauthAccount.emailAddress` da `profiles/work/.claude.json`. Quel file può pesare diversi megabyte, quindi `TopLevelValue` cerca quell'unica chiave invece di analizzare l'intero file.
7. `shim.Install` collega `~/.local/bin/ccwork` al binario di ccam, che all'avvio cerca `argv[0]` in `config.toml`.

Un login può richiedere minuti, e nel frattempo un altro ccam potrebbe salvare la configurazione. Per questo i comandi che restano in attesa salvano tramite `env.update`, che ricarica `config.toml`, applica l'unica modifica (qui, l'email) e salva. Per lo stesso motivo `ccam remove` ricarica la configurazione dopo la sua domanda.

`ccam login` esegue i passaggi 5 e 6 per un account esistente. `ccam logout` esegue `claude auth logout` nell'account e cancella l'email salvata. `ccam use` riscrive solo `default` in `config.toml`.

### Accedere su una macchina remota

Tramite SSH, la parte di `claude auth login` che passa dal browser può non funzionare. Claude Code attende il login su un server di callback all'indirizzo `http://localhost:<port>/callback` sul server, ma il browser che completa il login è in esecuzione sul tuo computer, dove quell'indirizzo non porta da nessuna parte e il browser mostra un errore come `ERR_CONNECTION_REFUSED`. L'indirizzo contiene comunque il codice monouso, quindi ccam può consegnarlo:

<!-- motion:accounts-relay -->

`runLogin` non cede il tuo terminale a Claude Code. Collega lo standard input di Claude Code a una pipe e `relayLogin` vi copia ogni riga che digiti, così un codice incollato al prompt `Paste code here` di Claude Code arriva comunque. Quando una riga è un indirizzo di callback, `callbackURL` lo accetta solo se l'host è `localhost`, `127.0.0.1` o `[::1]`, è indicata una porta, il percorso è `/callback` e sono presenti sia `code` sia `state`; `http://` si può omettere. `deliverCallback` richiede poi quell'indirizzo sul server con un limite di 10 secondi e senza seguire i reindirizzamenti, come avrebbe fatto il browser. La richiesta raggiunge il server di callback del `claude auth login` in attesa (la linea tratteggiata nel diagramma). Claude Code riceve il codice e salva il login, su Linux in `.credentials.json` nella cartella dell'account. Una riga incollata non può mai far contattare a ccam un altro host.

ccam stampa questo prima che inizi il login:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

Se riesce, stampa `→ Passed the sign-in to Claude Code on this machine.`

Mentre il login attende un input in un terminale, ccam disattiva l'eco e la modalità di riga del terminale e mostra un `*` per ogni carattere che digiti, così un codice o un indirizzo di callback incollato non compare sullo schermo. `readLine` gestisce da sé la modifica della riga: `Enter` invia la riga, `Backspace` cancella un carattere e `Ctrl+D` su una riga vuota termina l'input. Il terminale viene ripristinato quando il login finisce, e anche con `Ctrl+C`, `SIGTERM` o `SIGHUP`, che continuano a terminare ccam. Senza un terminale su standard input e standard error, le righe vengono copiate così come arrivano.

### Rinominare, rimuovere e importare

`ccam rename` e `ccam remove` agiscono sulla cartella dell'account e su ciò che ccam conserva altrove con il nome dell'account. `nameState` elenca questi percorsi: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` e le cache dell'utilizzo `~/.local/state/ccam/usage/<name>.json` e `<name>.statusline.json`. `ccam remove work` li gestisce in questo ordine:

<!-- motion:accounts-remove -->

`ccam remove` per prima cosa si rifiuta se l'account è in esecuzione, poi chiede conferma. Se l'account condivide le estensioni, `RewritePluginPaths` riscrive i file `installed_plugins.json` e `known_marketplaces.json` del gruppo, perché Claude Code salva percorsi assoluti dei plugin che possono passare per la cartella di questo account. `profile.Trash` sposta la cartella in `~/.local/share/ccam/trash/<name>-<time>/`, l'account esce da `config.toml`, i percorsi di `nameState` vengono eliminati, così un account creato in seguito con lo stesso nome non li eredita, e il link della scorciatoia viene rimosso. Per recuperarlo, sposta di nuovo la cartella fuori dal cestino.

`ccam rename` rinomina la cartella, riscrive i percorsi dei plugin nel suo `plugins/` dalla vecchia cartella alla nuova, mantiene l'eventuale impostazione di account predefinito e sposta i percorsi di `nameState` sul nuovo nome. Il comando scorciatoia mantiene il suo nome. Se la riscrittura dei plugin fallisce, la cartella torna al nome precedente, dato che la configurazione non viene salvata.

`ccam import` registra una cartella dove si trova, come `dir` in `config.toml`. Non la sposta né la copia mai: un refresh token copiato farebbe sì che le copie si disconnettano a vicenda. Senza una cartella come argomento cerca le cartelle `~/.claude` e `~/.claude-*` non registrate e chiede conferma per ciascuna; i nomi derivano dalla cartella (`~/.claude` diventa `main`, `~/.claude-work` diventa `work`, `~/.claude-2` diventa `acct2`), e una cartella da cui non si ricava un nome viene elencata con il comando `--name` da usare. Senza `CLAUDE_CONFIG_DIR`, Claude Code tiene `.claude.json` nella tua cartella home accanto a `~/.claude`, quindi quando importi `~/.claude` e non ha un proprio `.claude.json`, ccam crea il link `~/.claude/.claude.json` verso `~/.claude.json`. In questo modo `claude` da solo e `ccam main` vedono lo stesso account. Gli account importati non condividono nulla finché non esegui `ccam share`, e non viene creata nessuna scorciatoia.

`ccam completion bash|zsh` stampa uno script di completamento che completa anche i nomi degli account.

## Sicurezza e casi limite

- `ccam rename` e `ccam remove` si rifiutano mentre l'account è in esecuzione (`refuseRunning`), con `work is running; quit its Claude Code session first`, perché la sessione continua a scrivere nella cartella.
- `ccam remove` su un account importato si limita ad annullarne la registrazione: `✓ main unregistered (folder ... left in place)`. La tua `~/.claude` non viene spostata nel cestino.
- `ccam add` rifiuta un nome la cui cartella esiste già in `profiles/` e rimanda a `ccam import`.
- `ccam import` rifiuta una cartella già registrata, un percorso che non è una cartella e un nome già in uso.
- `ccam logout` e `ccam login` non controllano se l'account è in esecuzione.
- Se l'inoltro non riesce a raggiungere Claude Code, per esempio dopo che il login è terminato, stampa `ccam: could not pass the sign-in to Claude Code:` con il motivo. Un errore HTTP di solito significa che il codice è scaduto; avvia di nuovo il login.

## Pagine correlate

- [Avviare Claude Code](/features/launch)
- [Comandi scorciatoia](/features/shortcuts)
- [Condividere server MCP, impostazioni e cronologia](/docs/sharing)
- [Più account Claude Code](/guides/multiple-claude-code-accounts)
