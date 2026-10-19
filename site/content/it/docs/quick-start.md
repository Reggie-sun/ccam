---
title: Guida rapida
metaTitle: "Guida rapida: aggiungi e cambia account Claude Code"
description: Aggiungi due account Claude Code con ccam add, dai loro scorciatoie come ccwork, passa dall'uno all'altro e scegli cosa condividono.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 745a9904e32a
---

Aggiungi ogni account con `ccam add`, poi avvia Claude Code con `ccam <name>` o con una scorciatoia a cui dai tu il nome, come `ccwork`. È tutto ciò che serve alla maggior parte delle persone.

## 1. Aggiungi il primo account

```sh
ccam add
```

La procedura guidata chiede un nome come `personal` o `work`, quali elementi l'account deve condividere con gli altri tuoi account e se creare una scorciatoia (vedi il passo 4). Poi esegue `claude auth login`, che apre il browser per farti accedere. Al termine vedi:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Porta con te le impostazioni esistenti

Se usavi già Claude Code, le tue impostazioni, `CLAUDE.md`, le estensioni e i server MCP si trovano in `~/.claude`. Copiali nelle impostazioni condivise, così ogni account li riceve:

```sh
ccam migrate
```

ccam mostra cosa copierà e chiede conferma prima. Si limita a leggere `~/.claude`, non copia mai un login e lascia fuori la cronologia delle conversazioni, a meno che tu non aggiunga `--history`. Vedi [copiare una configurazione esistente](/docs/sharing#copy-an-existing-setup) per i dettagli.

## 3. Aggiungi il secondo account

```sh
ccam add work
```

Accedi con l'altro account nel browser.

## 4. Cambia account

`ccam work` avvia Claude Code con l'account work. Per qualcosa di più breve, dai all'account una scorciatoia tutta sua, come `ccwork`. Dopo gli elementi condivisi, la procedura guidata chiede:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Premi `y`, poi digita un nome, oppure premi Invio per usare quello suggerito:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Dopo il login, ccam la conferma:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

Se hai saltato la domanda, puoi aggiungere una scorciatoia in qualsiasi momento:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Ora uno qualsiasi di questi comandi avvia Claude Code con l'account work:

```sh
ccwork
ccam work
ccam 2
```

Eseguire `ccam` da solo apre un selettore: usa i tasti freccia e Invio, oppure premi il numero dell'account.

Tutto ciò che segue la scorciatoia o il nome dell'account arriva a Claude Code invariato, quindi continui a usare le opzioni che conosci:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Scegli cosa condividere

Per impostazione predefinita, i server MCP, `settings.json`, il `CLAUDE.md` globale e agenti, comandi, skill e plugin sono condivisi tra gli account. La cronologia delle conversazioni no. Attivala per entrambi gli account se vuoi continuare le conversazioni passando dall'uno all'altro:

```sh
ccam share personal --history on
ccam share work --history on
```

ccam mostra cosa sposterà prima di modificare qualsiasi cosa. Vedi [condivisione](/docs/sharing) per i dettagli.

## 6. Quando un account raggiunge il limite

```sh
ccam usage
ccam next --continue
```

`ccam usage` mostra l'utilizzo su 5 ore e settimanale di ogni account. `ccam next --continue` apre la stessa conversazione nell'account successivo. Vedi [limiti di utilizzo](/docs/usage-limits).

## Usi già ~/.claude?

Registra quella cartella invece di accedere di nuovo:

```sh
ccam import
```

La cartella resta dove si trova e `claude` da solo continua a funzionare come prima.
