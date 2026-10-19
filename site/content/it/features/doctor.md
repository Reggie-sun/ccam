---
title: Controllo e riparazione
description: Controlla link di condivisione, permessi dei file, login e versione di Claude Code con ccam doctor, e ripara ciò che si può riparare in sicurezza con ccam doctor --fix.
updated: 2026-10-19
order: 9
translation: machine
sourceHash: a1b61efe2b1f
---

`ccam doctor` controlla la tua installazione di ccam: il binario di Claude Code, la cartella, il login e i link di condivisione di ogni account, i server MCP condivisi, l'endpoint di utilizzo e i tuoi comandi scorciatoia. `ccam doctor --fix` ripara ciò che può riparare senza perdere dati e, prima di sostituire un file, lo sposta in una cartella di backup. Eseguilo quando `ccam <name>` si ferma per un problema con i link di condivisione, dopo aver ripristinato file a mano o dopo un aggiornamento di Claude Code.

## Utilizzo

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

Ogni risultato è una riga con un simbolo: `✓` OK, `·` informazione, `!` un avviso, `✗` un errore. I risultati che riguardano un account iniziano con il nome dell'account.

```console
$ ccam doctor
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700)
! personal: .credentials.json permissions are 644 (should be 0600)
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp

`ccam doctor --fix` repairs what it can.
```

Con `--fix`, ogni risultato che è stato riparato termina con `→ fixed`:

```console
$ ccam doctor --fix
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700) → fixed
! personal: .credentials.json permissions are 644 (should be 0600) → fixed
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link → fixed
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp → fixed
```

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `--fix` | Ripara ciò che si può riparare in sicurezza (dei file modificati viene fatto un backup) | disattivata |
| `--dry-run` (globale) | Esegue ogni controllo ma non ripara nulla, anche con `--fix` | disattivata |

`ccam doctor` non accetta argomenti. Rispetta `--dry-run`: `ccam --dry-run doctor --fix` stampa gli stessi risultati di `ccam doctor` e non modifica nulla.

Il codice di uscita è `1` quando resta almeno un risultato `✗` non riparato, altrimenti `0`. Avvisi e informazioni non fanno mai fallire l'esecuzione. Quando l'esecuzione fallisce senza `--fix`, ccam termina con una riga vuota e `` `ccam doctor --fix` repairs what it can. ``. Anche una configurazione che non si riesce a caricare termina con `1` e un errore `ccam:`.

## Come funziona

`internal/cli/doctor.go` carica `~/.config/ccam/config.toml` e chiama `doctor.Run` in `internal/doctor/doctor.go`, che esegue i controlli in un ordine fisso e restituisce i risultati. `--fix` viene passato come `Options.Fix`, impostato a false in un dry run.

<!-- motion:doctor-checks -->

Il blocco centrale viene eseguito una volta per account. I controlli, e cosa fa `--fix` per ciascuno:

1. **Versione di Claude Code** (`checkClaude`). Individua il binario `claude` come fa l'avvio ed esegue `claude --version`. Se la versione non inizia con `2.1.302`, la versione su cui è stata verificata la struttura delle cartelle di ccam (`doctor.LayoutVersion`), aggiunge `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed`. Un binario mancante o un `--version` che fallisce è un errore. Nessuna riparazione.
2. **Credenziali nell'ambiente** (`checkAuthEnv`). Avvisa quando è impostata `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` o `ANTHROPIC_AUTH_TOKEN`: ogni account avviato da quella shell la eredita, e Claude Code la preferisce al login dell'account. Nessuna riparazione; rimuovila tu.
3. **Cartella** (`checkProfile`). Una cartella dell'account mancante è un errore e termina i controlli di quell'account. Permessi con bit di gruppo o di altri utenti impostati sono un avviso; `--fix` esegue `chmod 0700`.
4. **File di login** (`checkCredentials`). `.credentials.json` viene controllato con `lstat`. Un link simbolico è un errore, perché ogni account deve avere il proprio token, e non viene riparato. Una modalità diversa da `0600` è un avviso; `--fix` esegue `chmod 0600`. Poi viene letto il login (su macOS prima dal Portachiavi). L'assenza di login è un avviso che rimanda a `ccam login <name>`. Un refresh token scaduto è l'errore `login expired`; un token di accesso scaduto è solo un'informazione, perché Claude Code lo rinnova al prossimo avvio. Nessuno di questi viene riparato.
5. **Link di condivisione** (`share.Check`). Confronta ogni percorso di `share.Layout` con le impostazioni di condivisione dell'account, usando solo `lstat` e `readlink`. Segnala quattro casi: un link mancante, un link che punta altrove, un file normale dove dovrebbe esserci un link e un link rimasto per un elemento che ora è locale. Un link mancante è un avviso, gli altri sono errori.
6. **Voci sconosciute** (`checkUnknownEntries`). Elenca i nomi di primo livello nella cartella dell'account che ccam non conosce (`share.Known`), come indizio che la struttura di Claude Code è cambiata. Solo informazione.
7. **Percorsi dei plugin** (`checkPlugins`). I registri dei plugin (`installed_plugins.json`, `known_marketplaces.json`) salvano percorsi assoluti. Nel `plugins/` di ogni gruppo, i percorsi che passano per la cartella di un account smettono di funzionare quando quell'account viene rimosso o rinominato. `--fix` li riscrive verso la cartella del gruppo con `profile.RewritePluginPaths`, una scrittura atomica del file di registro.
8. **Server MCP** (`checkMCP`). Carica una volta il `mcp.json` di ogni gruppo. Un file che non si riesce a leggere è un errore; un server che non supera `mcp.Validate` è un avviso. Nessuna riparazione.
9. **Formato della risposta di utilizzo** (`checkUsageFormat`). Interroga l'endpoint di utilizzo, senza cache e con il consueto limite di 3 secondi, per un account dopo l'altro finché uno risponde. Se la risposta non si riesce più ad analizzare, avvisa che `ccam usage` ripiega sulla cache della riga di stato. È l'unico controllo che usa la rete.
10. **Scorciatoie** (`checkShims`). Solo quando un account ha una scorciatoia. I link mancanti in `~/.local/bin` sono un avviso; `--fix` chiama `shim.Install`, che crea solo i link mancanti e non tocca mai altri file. Una cartella bin che non è nel `PATH` è un avviso e non viene riparata.

### Riparare i link di condivisione

Il percorso di avvio crea subito un link mancante, ma si ferma per qualsiasi cosa che potrebbe far perdere dati e rimanda qui:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Per qualsiasi problema con i link di condivisione, `--fix` chiama `share.Build` con le impostazioni di condivisione dell'account sia come stato attuale sia come stato di destinazione, poi `Plan.Apply`. `Build` ricava di nuovo ogni passaggio da ciò che si trova su disco, lo stesso piano usato da `ccam share`, quindi la riparazione segue le stesse regole dell'attivazione della condivisione.

<!-- motion:doctor-repair -->

Qui `settings.json` è un file normale. Poiché il gruppo ha già un `settings.json` diverso, per il conflitto viene usata la scelta predefinita, `KeepGroup`, quindi la copia del gruppo resta. Il file dell'account viene spostato in `.ccam-backup/<time>/settings.json` dentro la cartella dell'account e `os.Symlink` collega il percorso al gruppo. Altri casi:

- Se il gruppo non ha ancora una copia, il file dell'account viene spostato nel gruppo e ne diventa il primo contenuto.
- Per una cartella, i file che il gruppo non ha vengono spostati nel gruppo; il resto va nel backup.
- Per `history.jsonl`, il valore predefinito è `merge by line`: le righe del gruppo, poi le tue nuove righe. La copia precedente del gruppo viene salvata in `.ccam-backup/<time>/group/`.
- Un link a un altro gruppo viene sostituito. Un link a una cartella esterna a ccam va nel backup, e il suo contenuto diventa il contenuto iniziale del gruppo se il gruppo è vuoto.
- Un link rimasto per un elemento che ora è locale viene rimosso, e l'account riceve una copia del contenuto del gruppo.

A parte la rimozione dei link, ogni passaggio sposta i dati nel gruppo o nella cartella di backup. Non viene eliminato nulla.

## Sicurezza e casi limite

- Un semplice `ccam doctor` non modifica mai nulla. Negli script usa `--dry-run` insieme a `--fix` per esserne certo.
- `doctor --fix` non controlla se un account è in esecuzione. Chiudi Claude Code per un account prima di ripararne i link di condivisione, perché la sessione continua a scrivere in quei file.
- Se una riparazione della condivisione fallisce, ccam segnala `could not repair share links: <error>` per quell'account e salta il controllo delle voci sconosciute. Gli altri account vengono comunque controllati.
- Se `ccam share` fallisce a metà, salva prima le impostazioni di destinazione e ti dice di completare con `ccam doctor --fix`.
- Un `.credentials.json` collegato con un link simbolico, un login scaduto e le credenziali nell'ambiente vengono segnalati ma mai modificati.
- Le riscritture dei registri dei plugin sono atomiche ma non vengono copiate in `.ccam-backup/`.

## Pagine correlate

- [Condividere server MCP, impostazioni e cronologia](/docs/sharing)
- [Copiare una configurazione esistente](/features/migrate)
- [Limiti di utilizzo e cambio di account](/docs/usage-limits)
- [Comandi](/docs/commands)
