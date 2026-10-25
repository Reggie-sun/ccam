---
title: Controllo e riparazione
description: Controlla link di condivisione, permessi dei file, login e versione di Claude Code con ccam doctor, e ripara ciò che è sicuro riparare con ccam doctor --fix.
updated: 2026-10-25
order: 9
translation: machine
sourceHash: "0f5ae4ab9a54"
---

`ccam doctor` controlla la tua installazione di ccam: il binario di Claude Code, la cartella, il login e i link di condivisione di ogni account, i server MCP condivisi, l'endpoint di utilizzo e i tuoi comandi scorciatoia. `ccam doctor --fix` ripara ciò che può riparare senza perdere dati e, prima di sostituire un file, lo sposta in una cartella di backup. Eseguilo quando `ccam <name>` si ferma per un problema con i link di condivisione, dopo aver ripristinato file a mano o dopo un aggiornamento di Claude Code.

## Utilizzo

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

La prima riga segna tutti gli account in una volta, così vedi dove sono i problemi prima di leggere qualsiasi altra cosa. Sotto, ogni account ha un riquadro, e i risultati che non appartengono a nessun account (Claude Code, variabili d'ambiente, plugin, server MCP, scorciatoie, regole delle cartelle, gruppi e impostazioni, `auto_switch`) finiscono in un riquadro `General`. Il titolo del riquadro mostra il simbolo peggiore dell'account e quanti problemi ha. Ogni riga all'interno è un'area controllata, sempre `Folder`, `Login` e `Share links` in quest'ordine, poi le eventuali altre, con un simbolo: `✔` OK, `ℹ` informazione, `⚠` un avviso, `✘` un errore. Una cartella senza problemi mostra `✔ Folder  OK`.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ───────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp              │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

Con `--fix`, ogni risultato che è stato riparato termina con `→ fixed` e non conta più come problema:

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ───────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed      │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ──────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed        │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →   │
 │                fixed                                                      │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link → fixed     │
 │ ℹ              settings.json differed from the group's copy, which it now │
 │                uses; this account's copy is in                            │
 │                /home/me/.local/share/ccam/profiles/personal/.ccam-        │
 │                backup/20261023-212553/settings.json                       │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯
```

In un terminale i riquadri si adattano alla sua larghezza e i messaggi lunghi vanno a capo, con i colori del terminale (`NO_COLOR` li disattiva). Un riquadro non è mai più stretto di 29 colonne, così in un terminale stretto nessuna parola viene tagliata. Inviato a una pipe o a un file, l'output non ha colori e nessun messaggio va a capo o viene tagliato, così `ccam doctor | grep settings.json` trova la riga intera; ogni riquadro è allora largo quanto la riga più lunga. La scheda Health della dashboard mostra gli stessi riquadri.

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `--fix` | Ripara ciò che si può riparare in sicurezza (dei file modificati viene fatto un backup) | disattivata |
| `--dry-run` (globale) | Esegue ogni controllo ma non ripara nulla, anche con `--fix` | disattivata |

`ccam doctor` non accetta argomenti. Rispetta `--dry-run`: `ccam --dry-run doctor --fix` stampa i risultati di `ccam doctor`, aggiunge ciò che `--fix` metterebbe da parte e non modifica nulla:

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

Il codice di uscita è `1` quando resta almeno un risultato `✘` non riparato, altrimenti `0`. Avvisi e informazioni non fanno mai fallire l'esecuzione. Quando l'esecuzione fallisce senza `--fix`, ccam termina con una riga vuota e `` `ccam doctor --fix` repairs what it can. ``. Anche una configurazione che non si riesce a caricare termina con `1` e un errore `ccam:`.

## Come funziona

`internal/cli/doctor.go` carica `~/.config/ccam/config.toml` e chiama `doctor.Run` in `internal/doctor/doctor.go`, che esegue i controlli in un ordine fisso e restituisce i risultati, ciascuno con l'account a cui appartiene e l'area che ha controllato. `tui.HealthReport` li raggruppa nei riquadri. `--fix` viene passato come `Options.Fix`, impostato a false in un dry run.

<!-- motion:doctor-checks -->

Il blocco centrale viene eseguito una volta per account. I controlli, e cosa fa `--fix` per ciascuno:

1. **Versione di Claude Code** (`checkClaude`). Individua il binario `claude` come fa l'avvio ed esegue `claude --version`. Un binario mancante, un `claude_bin` che non è un file eseguibile o un `--version` che fallisce è un errore. Nessuna riparazione.
2. **Credenziali nell'ambiente** (`checkAuthEnv`). Avvisa quando è impostata `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` o `ANTHROPIC_AUTH_TOKEN`: ogni account avviato da quella shell la eredita, e Claude Code la preferisce al login dell'account. Nessuna riparazione; rimuovila tu.
3. **Cartella** (`checkProfile`). Una cartella dell'account mancante è un errore e termina i controlli di quell'account. Permessi con bit di gruppo o di altri utenti impostati sono un avviso; `--fix` esegue `chmod 0700`. Una `~/.claude` importata con un proprio `.claude.json` accanto a `~/.claude.json` è un'informazione: `ccam main` e `claude` da solo tengono allora separati i dati del login e i server MCP. Windows non ha bit di permesso, quindi lì questo controllo e quello su `0600` qui sotto vengono saltati; la cartella nel tuo profilo utente è protetta dalla sua lista di controllo di accesso.
4. **File di login** (`checkCredentials`). `.credentials.json` viene controllato con `lstat`. Un link simbolico è un errore, perché ogni account deve avere il proprio token, e non viene riparato. Una modalità diversa da `0600` è un avviso; `--fix` esegue `chmod 0600`. Poi viene letto il login (su macOS prima dal Portachiavi). L'assenza di login, o un file di login che non ne contiene uno, è un avviso che rimanda a `ccam login <name>`. Un refresh token scaduto, o un login cancellato da Claude Code dopo che il server lo ha rifiutato, è l'errore `` login expired → `ccam login <name>` ``; un token di accesso scaduto è solo un'informazione, perché Claude Code lo rinnova al prossimo avvio. Un login conservato in Gestione credenziali di Windows è un'informazione, perché ccam non può leggerlo. Un login valido mostra l'email con cui la cartella ha effettuato il login adesso e il piano (`login valid (dev@company.com, Claude Team · Acme)`). Nessuno di questi viene riparato; `ccam login --expired` effettua il login in ogni account che ne ha bisogno.
5. **Link di condivisione** (`share.Check`). Confronta ogni percorso di `share.Layout` con i gruppi dell'account, usando solo `lstat` e `readlink`. Segnala cinque casi: un link mancante, un link che punta altrove, un file normale dove dovrebbe esserci un link, un link rimasto per un elemento che ora è locale, e un link corretto la cui cartella del gruppo manca (`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)`; un avvio non lo controlla). Un link mancante e una cartella del gruppo mancante sono avvisi, gli altri sono errori. Viene segnalata anche una copia di ritorno dal gruppo rimasta a metà: `--fix` mette al suo posto una copia preparata completa, ma quando il percorso è stato ricreato nel frattempo (da un avvio), indica la cartella con la copia completa perché tu possa sistemarla. I percorsi aggiunti a un elemento in una versione successiva di ccam (`file-history/`, `tasks/` e `paste-cache/` per la cronologia, `keybindings.json` per le impostazioni) sono solo un avviso finché l'account ne ha ancora una copia propria: `is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`. Un avvio non si ferma mai per questi.
6. **Voci sconosciute** (`checkUnknownEntries`). Elenca i nomi di primo livello nella cartella dell'account che ccam non conosce (`share.Known`), come indizio che la struttura di Claude Code è cambiata. Solo informazione.
7. **Percorsi dei plugin** (`checkPlugins`). I registri dei plugin (`installed_plugins.json`, `known_marketplaces.json`) salvano percorsi assoluti. Nel `plugins/` di ogni gruppo, i percorsi che passano per la cartella di un account smettono di funzionare quando quell'account viene rimosso o rinominato. `--fix` li riscrive verso la cartella del gruppo con `profile.RewritePluginPaths`, una scrittura atomica del file di registro.
8. **Server MCP** (`checkMCP`). Carica una volta il `mcp.json` di ogni gruppo. Un file che non si riesce a leggere è un errore; un server che non supera `mcp.Validate` è un avviso. I server nel `~/.claude.json` di `claude` da solo che nessun account carica sono un avviso che indica il comando `ccam migrate --mcp` che li copia (con `--group <a|b>` quando più gruppi condividono MCP). I server che ha un solo account di un gruppo sono un'informazione con il comando `ccam mcp import` che li condivide, e un account che non ha effettuato l'accesso a server MCP a cui gli altri account hanno accesso riceve un avviso per accedere con `/mcp` (ccam non copia mai gli accessi). Nessuna riparazione.
9. **Formato della risposta di utilizzo** (`checkUsageFormat`). Interroga l'endpoint di utilizzo, senza cache, per un account dopo l'altro finché uno risponde, il tutto entro un unico limite di 4 secondi. Con `usage_source = "claude"` chiede invece a Claude Code per un account. Se la risposta non si riesce più ad analizzare, avvisa che `ccam usage` ripiega sulla cache della riga di stato; quando nessun account ha risposto, lo dice come informazione (`could not check the usage response format: …`). È l'unico controllo che usa la rete.
10. **Scorciatoie** (`checkShims`). Solo quando un account ha una scorciatoia. I link mancanti in `~/.local/bin` sono un avviso, e lo sono anche le scorciatoie create da ccam che eseguono ancora un ccam più vecchio, uno sostituito o spostato (`shortcuts in <bin> that run another ccam (one replaced or moved): …`); `--fix` chiama `shim.Install`, che crea e rinnova solo i link creati da ccam e registrati in `~/.local/state/ccam/shortcuts.json`, e non tocca mai altri file. L'avviso termina con `→ fixed` solo quando ciascuno di quei link è stato creato; se uno non è stato possibile, viene aggiunto il motivo. Una scorciatoia che esegue un altro ccam ancora installato funziona, ed è un'informazione; `--fix` la fa puntare a questo ccam. Una cartella bin che non è nel `PATH` è un avviso e non viene riparata.
11. **Cartelle** (`checkFolders`). Solo quando esistono voci di [`ccam folder`](/features/folder). Una voce il cui account non esiste è un errore, e una voce la cui cartella non c'è più è un avviso. Nessuna delle due viene riparata, perché rimuovere una voce cambierebbe l'account che parte lì.
12. **Gruppi e impostazioni** (`checkGroups`). Un membro di un gruppo che non è un account, un account il cui nome nel frattempo è diventato un comando di ccam (avvialo con `ccam run <name>` o rinominalo), un `theme` o un `usage_source` sconosciuto, e le righe di `config.toml` che ccam non legge sono avvisi. Queste ultime di solito sono un'impostazione scritta sotto una tabella `[[group]]`, `[[profile]]` o `[[folder]]`, dove appartiene a quella tabella (`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`); i comandi che modificano la configurazione si rifiutano finché non viene corretta. Nessuna riparazione.
13. **Login** (`checkIdentities`). Legge da `.claude.json` con quale account ha effettuato il login ogni cartella. Una cartella ora collegata a un account diverso da quello registrato da `ccam login` (per esempio dopo `/login` dentro una sessione) è un avviso; un account senza registrazione del login, come uno importato o creato da un ccam più vecchio, viene confrontato con la sua email registrata. Anche una cartella collegata allo stesso account Claude di un'altra è un avviso: le due condividono i limiti di utilizzo. Una cartella che contiene lo stesso token di login di un'altra, una cartella copiata (compresa la `~/.claude` di `claude` da solo), è un errore, perché chi lo rinnova per primo disconnette l'altra; esegui `ccam login <name>` per darle il suo. Nessuna riparazione.
14. **`auto_switch`** (`checkAutoSwitch`). Solo con `auto_switch = true`. Un'informazione nel riquadro `General` dice cosa fa `auto_switch` e rimanda alla [sua descrizione](/features/next#auto-switch). Per ogni piano con posti tra gli account (Claude Team, Claude Enterprise), una riga li elenca, dice che `auto_switch` sposta una sessione tra i posti di una stessa organizzazione, cita la regola di Anthropic sui posti con il link al centro assistenza e lascia a te la scelta. Un Claude Code più vecchio della 2.1.214, che non comunica agli hook il proprio processo (`CLAUDE_PID`), è lì un avviso, perché nessuna sessione può essere spostata. Un account senza il plugin ccam, o con il plugin disattivato, riceve un avviso nel proprio riquadro con il comando che lo risolve in Claude Code: `auto_switch is on, but work does not have the ccam plugin, so its sessions are not moved. In Claude Code as work: /plugin install ccam --marketplace jhyunwoo/ccam` (oppure `/plugin enable ccam@ccam`). Gli account con chiave API restano esclusi. Su Windows l'unica riga è `auto_switch is ignored on Windows; with relay = true, /ccam:switch and then /exit carry a conversation on.` Nessuna riparazione.

Prima degli account, un **cambio di gruppi** rimasto a metà è un errore nel riquadro `General`: ``the last change of groups stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``. `--fix` lo completa nell'ordine in cui era stato pianificato, con le risposte date ai suoi conflitti; fino ad allora non parte nessun altro cambio di gruppi, rinomina o rimozione. Gli account coinvolti ricevono solo i link mancanti finché non è completato. Se non si può completare, l'errore dice che rimuovere quel file rinuncia al resto del cambio. `ccam --dry-run doctor --fix` mostra l'ordine (`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`) e cosa metterebbe da parte la parte di ciascun account. Un cambio che un altro ccam sta eseguendo in questo momento è un'informazione, e mentre è in corso `--fix` non ripara i link di condivisione (`its share links are not repaired while another ccam changes the groups; …`). `--fix` tiene `~/.local/state/ccam/group-change.lock` mentre ripara, e rilegge `config.toml` dopo aver completato un cambio.

### Riparare i link di condivisione

Il percorso di avvio crea subito un link mancante, ma si ferma per qualsiasi cosa che potrebbe far perdere dati e rimanda qui:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Per qualsiasi problema con i link di condivisione, `--fix` chiama `share.Build` con le impostazioni di condivisione dell'account sia come stato attuale sia come stato di destinazione, poi `Plan.Apply`. `Build` ricava di nuovo ogni passaggio da ciò che si trova su disco, lo stesso piano usato da una modifica con `ccam group`, quindi la riparazione segue le stesse regole dell'ingresso in un gruppo.

<!-- motion:doctor-repair -->

Qui `settings.json` è un file normale. Poiché il gruppo ha già un `settings.json` diverso, per il conflitto viene usata la scelta predefinita, `KeepGroup`, quindi la copia del gruppo resta. Il file dell'account viene spostato in `.ccam-backup/<time>/settings.json` dentro la cartella dell'account e `os.Symlink` collega il percorso al gruppo. `--fix` lo segnala: `settings.json differed from the group's copy, which it now uses; this account's copy is in …`. Altri casi:

- Se il gruppo non ha ancora una copia, il file dell'account viene spostato nel gruppo e ne diventa il primo contenuto.
- Per una cartella, i file che il gruppo non ha vengono spostati nel gruppo. I file che il gruppo ha identici byte per byte vengono scartati, e quelli diversi vanno nel backup, indicati in una nota (`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`).
- Una cartella del gruppo mancante viene creata per prima.
- Per `history.jsonl`, il valore predefinito è `merge by line`: le righe del gruppo, poi le tue nuove righe. La copia precedente del gruppo viene salvata in `.ccam-backup/<time>/group/`.
- Un link a un altro gruppo viene sostituito. Un link a una cartella esterna a ccam va nel backup, e il suo contenuto diventa il contenuto iniziale del gruppo se il gruppo è vuoto.
- Un link rimasto per un elemento che ora è locale viene rimosso, e l'account riceve una copia del contenuto del gruppo.

A parte la rimozione dei link e lo scarto dei file che il gruppo ha già identici byte per byte, ogni passaggio sposta i dati nel gruppo o nella cartella di backup. Non si perde nulla.

## Sicurezza e casi limite

- Un semplice `ccam doctor` non modifica mai nulla. Negli script usa `--dry-run` insieme a `--fix` per esserne certo.
- Per un account in esecuzione, comunque sia stato avviato Claude Code, `doctor --fix` crea solo i link mancanti e le cartelle del gruppo mancanti, il che non fa perdere nulla. Quando resta altro, aggiunge `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest ``, perché la sessione continua a scrivere nei file che una riparazione sposterebbe.
- Se una riparazione della condivisione fallisce, ccam segnala `could not repair share links: <error>` per quell'account e salta il controllo delle voci sconosciute. Gli altri account vengono comunque controllati.
- Se una modifica con `ccam group` fallisce a metà, i nuovi gruppi sono già salvati e il cambio viene annotato in `~/.local/state/ccam/group-change.json`; ccam ti dice di completare con `ccam doctor --fix`, che esegue il resto, compresi i passaggi MCP.
- Un `.credentials.json` collegato con un link simbolico, un login scaduto e le credenziali nell'ambiente vengono segnalati ma mai modificati.
- Le riscritture dei registri dei plugin sono atomiche ma non vengono copiate in `.ccam-backup/`.

## Pagine correlate

- [Condividere server MCP, impostazioni e cronologia](/docs/sharing)
- [Copiare una configurazione esistente](/features/migrate)
- [Limiti di utilizzo e cambio di account](/docs/usage-limits)
- [Comandi](/docs/commands)
