---
title: Guida rapida
metaTitle: "Guida rapida: aggiungi e cambia account Claude Code"
description: Aggiungi due account Claude Code con ccam add, dai loro scorciatoie come ccwork, passa dall'uno all'altro e scegli cosa condividono.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: bd893bbb83ab
---

Aggiungi ogni account con `ccam add`, poi avvia Claude Code con `ccam <name>` o con una scorciatoia a cui dai tu il nome, come `ccwork`. È tutto ciò che serve alla maggior parte delle persone.

Usi già Claude Code? Prima di iniziare, leggi [Usi già ~/.claude?](#already-using-claude) alla fine di questa pagina.

## 1. Aggiungi il primo account

```sh
ccam add
```

La procedura guidata chiede un nome come `personal` o `work`, quali elementi l'account deve condividere con gli altri tuoi account e se creare una scorciatoia (vedi il passo 3). Poi esegue `claude auth login`, che apre il browser per farti accedere. Su una macchina remota il browser può finire su una pagina `http://localhost:…/callback` che non si carica; incolla quell'indirizzo nel terminale e ccam completa l'accesso. Quando l'accesso è completato, vedi:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Porta con te le impostazioni esistenti

Se hai già usato Claude Code, le tue impostazioni, `CLAUDE.md` e le estensioni si trovano in `~/.claude`, e i tuoi server MCP in `~/.claude.json`. Copiali nel gruppo condiviso `default`, così ogni account li riceve:

```sh
ccam migrate
```

ccam mostra cosa copierà e chiede conferma prima. Si limita a leggere `~/.claude`, non copia mai un login e lascia fuori la cronologia delle conversazioni, a meno che tu non aggiunga `--history`. Vedi [copiare una configurazione esistente](/docs/sharing#copy-an-existing-setup) per i dettagli.

## 3. Aggiungi il secondo account

```sh
ccam add work
```

La procedura guidata salta il nome, che hai indicato nella riga di comando, e mostra i gruppi a cui unirsi con `default` selezionato. Questa volta, dai all'account una scorciatoia tutta sua, come `ccwork`, così puoi avviarlo con una sola parola. Dopo i gruppi, la procedura guidata chiede:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Premi `y`, poi digita un nome, oppure premi Invio per usare quello suggerito:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Poi accedi con l'altro account nel browser. ccam conferma l'account e la sua scorciatoia:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

## 4. Cambia account

`ccam work` avvia Claude Code con l'account work. Se hai saltato la domanda sulla scorciatoia, puoi aggiungerne una in qualsiasi momento:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Ora ciascuno di questi comandi avvia Claude Code con l'account work (`2` è il suo numero in `ccam list`):

```sh
ccwork
ccam work
ccam 2
```

Eseguire `ccam` da solo apre la [dashboard](/features/dashboard): spostati su un account con i tasti freccia o con il suo numero, poi premi Invio.

Tutto ciò che segue la scorciatoia o il nome dell'account arriva a Claude Code invariato, quindi continui a usare le opzioni che conosci:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Scegli cosa condividere

Per impostazione predefinita, i server MCP, `settings.json`, il `CLAUDE.md` globale e le estensioni (agenti, comandi, skill e plugin) sono condivisi tra gli account tramite il gruppo `default`. La cronologia delle conversazioni no. Aggiungila al gruppo se vuoi continuare le conversazioni passando da un account all'altro:

```sh
ccam group set default --add history
```

ccam mostra cosa sposterà per ogni account prima di modificare qualsiasi cosa. Vedi [condivisione con i gruppi](/features/share) per i dettagli.

## 6. Quando un account raggiunge il limite

```sh
ccam usage
ccam next --continue
```

`ccam usage` mostra l'utilizzo su 5 ore e settimanale di ogni account. `ccam next --continue` apre la stessa conversazione nell'account successivo. Negli script, `ccam --fallback work -p "…"` cambia account da solo. Vedi [limiti di utilizzo](/docs/usage-limits).

**Posti nei piani Team ed Enterprise.** La scelta automatica dell'account (`ccam next`, `--fallback`, `ccam pick`, il tasto `n` della dashboard e, una volta attivato, [`auto_switch`](/features/next#auto-switch)) passa anche da un posto all'altro della stessa organizzazione Claude Team o Enterprise, e questo può violare i termini di Anthropic. Leggi le [domande frequenti](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) e verifica con la tua organizzazione e con Anthropic prima di farci affidamento.

## Usi già ~/.claude?

Registra quella cartella invece di accedere di nuovo:

```sh
ccam import
```

La cartella resta dove si trova e `claude` da solo continua a funzionare come prima.
