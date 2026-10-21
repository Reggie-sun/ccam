---
title: "Cambiare account Claude Code: ccam e le alternative"
description: Come si confronta ccam con claude-swap, CCS, clausona, clauth, CAAM, cux e altri strumenti per cambiare account Claude Code, e quale approccio scegliere.
updated: 2026-10-21
order: 2
translation: machine
sourceHash: "c74be64d2978"
---

Gli strumenti per cambiare account Claude Code seguono uno di due approcci. La maggior parte scambia il login dentro l'unica cartella `~/.claude`, quindi è attivo un solo account alla volta. Altri, tra cui ccam, danno a ogni account la sua cartella con `CLAUDE_CONFIG_DIR`, quindi ogni account resta connesso e più account possono girare insieme. ccam aggiunge impostazioni e server MCP condivisi, e continua una conversazione nella cartella di un altro account.

## Scambiare un login o usare cartelle separate?

**Scambiare un login** mantiene un solo insieme di impostazioni, server MCP e cronologia, perché c'è una sola cartella. Il cambio è rapido, e una sessione in corso può adottare il nuovo login. Però è attivo un solo account alla volta, per tutti i terminali e gli editor insieme.

Per fare lo scambio, lo strumento conserva una copia del login di ogni altro account e la rimette al suo posto. Claude Code sostituisce il refresh token di un login ogni volta che lo rinnova, quindi ogni copia salvata deve restare allineata con le sessioni in corso. Account disconnessi o che chiedono di nuovo `/login` dopo un cambio sono tra i problemi che questi progetti correggono più spesso.

**Le cartelle separate** sono il metodo che [descrive la documentazione di Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Ogni account resta connesso, e puoi eseguire una sessione di lavoro e una personale in parallelo. Ogni login resta nella sua cartella, dove Claude Code lo rinnova come sempre. Il prezzo è che ogni cartella parte con impostazioni, server MCP e cronologia propri, a meno che uno strumento non li condivida.

## Gli strumenti a confronto

Tratto dalla documentazione e dal codice sorgente di ciascun progetto il 2026-10-20 (clausona: v0.12.0-beta). "Non documentato" significa che il progetto non lo dice, non che la funzione manchi.

| Strumento | Approccio | Più account insieme | Utilizzo su 5 ore e settimanale | Quando si raggiunge un limite | Piattaforme |
| --- | --- | --- | --- | --- | --- |
| ccam | Cartelle separate | Sì | Sì, `ccam usage` | Una notifica nella sessione; `ccam next --auto --continue` o `/ccam:switch` porta la conversazione nell'account con più margine, con o senza cronologia condivisa; `ccam --fallback` riesegue lì da solo i job `-p`, e `ccam pick` distribuisce i worker | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | Scambia il login; cartelle separate con `cswap run` | Sì, con `cswap run` | Sì | `cswap auto` cambia il login al superamento di una soglia; la sessione in corso continua | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Cartelle separate | Sì | In parte | I pool di account passano le richieste all'account successivo tramite un proxy locale | Linux, macOS, Windows |
| [clausona](https://github.com/larcane97/clausona) | Cartelle separate | Sì | Sì | Le route scelgono un account con quota prima di ogni esecuzione, possono attendere un reset e assegnano account diversi ai worker | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Scambia il login; cartelle separate con `clauth start` | Sì, con `clauth start` | Sì | Cambia seguendo una catena di fallback | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Scambia il login, oppure profili isolati | Sì, con i profili isolati | Sì | `caam run` riprova i comandi non interattivi su un altro account | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Scambia il login e fa da wrapper a `claude` | No, un solo account attivo | Sì | Scambia il login e riprende la stessa conversazione | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Scambia il login | No | Sì | Cambio facoltativo al superamento di una soglia | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Scambia il login | No | Sì | Cambio facoltativo al 100% | macOS |

ccswap è un fork di claude-swap che gestisce anche i login di Codex. Anche [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) e [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) danno a ogni account la sua cartella. Diversi di questi strumenti gestiscono anche altri agenti di programmazione, come Codex o Gemini.

## Come sono fatti

| Strumento | Scritto in | Mentre Claude Code è in esecuzione | Copia i login da un posto all'altro | Cosa possono condividere gli account |
| --- | --- | --- | --- | --- |
| ccam | Go, un unico binario statico | Niente: ccam si sostituisce con `claude` (su Windows, che non lo permette, ne attende la fine; un relay facoltativo resta accanto) | No | Gruppi: ognuno indica gli elementi che condivide e i suoi account; un account può stare in più gruppi |
| claude-swap | Python 3.12+, installato con uv o pipx | Niente con `cswap run`; `cswap auto` e l'app della barra dei menu restano in esecuzione | Sì, quando scambia il login | Con `cswap run`: impostazioni, `CLAUDE.md` ed estensioni insieme, cronologia facoltativa; server MCP copiati a ogni avvio |
| CCS | Node.js 18+ o Bun | CCS resta come processo padre; i pool aggiungono un proxy locale | No per le sue cartelle di account; i pool hanno login propri | Condiviso o locale per account, cronologia in gruppi di contesto; server MCP copiati |
| clausona | TypeScript, Node.js 20+ | `claude` semplice passa per una funzione della shell; `clausona run` tiene Node come processo padre; un processo Node in background dopo ogni uscita | No, ma rinnova da sé i token degli account inattivi per leggerne l'utilizzo | Tutto tranne i file propri di ogni account, da un unico account principale; cronologia facoltativa; i server aggiunti con `claude mcp add` non vengono condivisi |
| clauth | Rust | clauth resta come processo padre | Sì, quando scambia il login | Tutto `~/.claude`, oppure niente |
| CAAM | Go | Dipende dalla modalità; un daemon facoltativo | Sì, nella modalità vault | Fissato dalla modalità; server MCP copiati a ogni avvio |
| cux | Go | Un wrapper resta attorno a `claude` | Sì | Tutto: una sola cartella |
| CCSwitcher, Claude Switcher | App macOS (Swift, Python) | L'app della barra dei menu | Sì | Tutto: una sola cartella |

## Cosa ccam fa bene

### I login restano dove li ha messi Claude Code

ccam non copia mai un login di Claude da un posto all'altro. Il login di ogni account resta nella sua cartella, o su macOS nell'elemento del portachiavi proprio di quella cartella, ed è Claude Code il programma che lo rinnova. Non esiste una copia salvata che possa disallinearsi. Un `~/.claude` esistente registrato con `ccam import` mantiene l'elemento del portachiavi usato da `claude` da solo, quindi resta connesso in entrambi i casi.

ccam rinnova da sé un token scaduto solo quando esegui `ccam usage --refresh`. Prima acquisisce un lock e si rifiuta finché quell'account è in esecuzione, quindi non ruota mai un token sotto una sessione attiva.

### Niente gira tra te e Claude Code

ccam legge un file, controlla i suoi link e si sostituisce con `claude` (`exec`). Mentre lavori non c'è alcun processo ccam, wrapper, proxy o daemon: standard input, segnali e codice di uscita appartengono a Claude Code, e le sue richieste vanno direttamente ad Anthropic. Windows non può sostituire un processo, quindi lì ccam avvia Claude Code con la stessa console e ne attende il codice di uscita, senza fare altro.

Avviare un account così aggiunge circa 4 ms, sotto i 6 ms al 95° percentile, con 8 MB di memoria di picco ([`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) con hyperfine, 300 esecuzioni). Sullo stesso portatile, il solo avvio dell'interprete Python richiede 13 ms e quello di Node.js 31 ms, prima che un programma faccia qualsiasi cosa.

### I gruppi decidono cosa si condivide

Server MCP, `settings.json`, `CLAUDE.md`, estensioni (agenti, comandi, skill, plugin) e cronologia delle conversazioni sono cinque elementi separati. Un gruppo indica gli elementi che condivide e gli account che ne fanno parte, così due account di lavoro possono condividere la cronologia in un gruppo `work` mentre un account personale tiene la propria, e tutti e tre condividono gli stessi server MCP tramite `default`. `ccam group set work --add memory` sposta tutti gli account del gruppo con un unico piano; la [dashboard](/features/dashboard) fa lo stesso con le caselle di spunta e disattiva un elemento che un account riceve già da un altro gruppo.

- I server MCP vengono passati con `--mcp-config` all'avvio, quindi per impostazione predefinita ccam non modifica mai `.claude.json`, il file che Claude Code riscrive mentre è in esecuzione. `ccam mcp add` aggiunge un server una sola volta per tutto il gruppo.
- Gli altri elementi sono link alla cartella del gruppo, quindi una modifica fatta in un account compare subito negli altri. All'avvio non viene copiato nulla.
- Login, `.claude.json` e stato di esecuzione non vengono mai condivisi.

### Ogni modifica viene mostrata prima e salvata in backup

`ccam group` e `ccam migrate` mostrano un piano e chiedono conferma prima di spostare qualsiasi cosa. Quando un file differisce tra un account e il suo gruppo, scegli tu: tenere quello del gruppo, usare il tuo o unire il JSON. Gli originali vanno in `.ccam-backup/<time>/`, la cartella di un account rimosso finisce in una cartella cestino, e i file vengono sostituiti con un'unica rinomina, quindi una scrittura interrotta lascia il file precedente. `--dry-run` mostra cosa farebbero avvio, condivisione, migrazione, scorciatoie, aggiornamenti e riparazioni. `ccam doctor --fix` ripara link e permessi danneggiati.

### La stessa conversazione nell'account successivo

`ccam next --auto --continue` chiede a tutti gli account il loro utilizzo contemporaneamente, apre quello con più margine e continua la conversazione. Con la cronologia condivisa è la stessa conversazione; senza, Claude Code la copia nel nuovo account (`--resume <file> --fork-session`), e ccam dice quanti token rilegge la prima risposta. Quando un limite ferma una sessione, l'hook del plugin ccam mostra una notifica, e `/ccam:switch` passa la conversazione a un altro account. Decidi tu quando cambiare; niente controlla il tuo utilizzo in background, e le scelte automatiche restano dentro una sola organizzazione, quindi una conversazione aziendale non finisce mai per sbaglio su un account privato. La riga di stato suggerisce il comando quando un account arriva al 90%.

Gli script non possono aspettarti, quindi `ccam --fallback work -p "…"` cambia account al posto loro: quando l'account raggiunge il limite, lo stesso comando viene rieseguito con l'account con più margine, e viene stampato solo l'output dell'esecuzione portata a termine. Un errore che non è un limite di utilizzo viene passato così com'è. Per i worker in parallelo, `ccam pick` assegna a ciascuno un account diverso con margine, ed entrambi escono con codice 75 quando nessun account ne ha più.

### L'account giusto in ogni cartella

`ccam folder work ~/code/company` fa sì che ogni avvio senza nome di account dentro quella cartella usi `work`, mentre il resto dei tuoi progetti continua a usare l'account predefinito. L'impostazione sta nello stesso `config.toml` che ogni avvio legge, quindi non costa tempo.

### Un solo binario, download verificati

ccam è un binario statico di circa 10 MB, senza runtime da installare, per Linux, macOS e Windows. Gli script di installazione (`install.sh`, `install.ps1`) e `ccam update` verificano ogni download con i checksum SHA-256 della release, e ccam non si aggiorna mai senza chiedere. Non ha telemetria: le uniche richieste che fa sono le consultazioni dell'utilizzo dei tuoi account presso Anthropic, in parallelo con un limite di 3 secondi e una cache di 60 secondi, la domanda della dashboard a GitHub sull'ultima release quando si apre, e i download che chiedi tu. Con `usage_source = "claude"`, anche queste passano per Claude Code, e ccam non tocca mai un token.

## Dove altri strumenti fanno di più

- claude-swap, clauth, cux, CAAM e CCSwitcher possono cambiare da soli l'account di una sessione interattiva in corso quando un account supera una soglia di utilizzo. ccam cambia le esecuzioni `-p` con `--fallback` e lascia a te le sessioni interattive, con una notifica al limite e `/ccam:switch`.
- Su Windows senza la Modalità sviluppatore, ccam condivide le cartelle ma non i singoli file come `settings.json`; claude-swap e CCS in quel caso li copiano.
- CCSwitcher, Claude Switcher e claude-swap mostrano l'utilizzo nella barra dei menu di macOS.
- CCS, clausona, CAAM, clauth, ccswap e Claude Switcher gestiscono anche Codex e altri agenti di programmazione.
- clausona attiva e disattiva anche le estensioni per account (`ccam ext` di ccam si limita a mostrarle), tiene traccia del costo locale dei token e distribuisce plugin di Claude Code che eseguono flotte di worker in tmux e altri terminali.
- CAAM, CCSwitcher e Claude Switcher si installano con Homebrew, e clauth e CAAM firmano i checksum delle loro release. ccam si installa con il suo script, un `.deb` o `go install`, e verifica i checksum senza firma.

## Quale fa per te?

Scegli uno strumento che scambia il login se usi un account alla volta o vuoi che le sessioni in corso cambino account in automatico. Scegli ccam se vuoi account in parallelo, un'unica configurazione condivisa, la stessa conversazione sull'account successivo e niente in esecuzione accanto a Claude Code, su Linux, macOS o Windows. Se non ti serve condividere nulla, bastano [semplici alias della shell](/compare/aliases).

Se puoi usare più account dipende dai termini di ciascun account. Controllali prima di farlo.
