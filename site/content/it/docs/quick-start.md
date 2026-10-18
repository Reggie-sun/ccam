---
title: Guida rapida
description: Aggiungi due account Claude Code con ccam add, passa dall'uno all'altro con claude1, claude2 o ccs e scegli cosa condividono.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 9462b80d61bc
---

Aggiungi ogni account con `ccam add`, poi avvia Claude Code con `claude1`, `claude2` o `ccs <name>`. Alla maggior parte delle persone non serve altro.

## 1. Aggiungi il primo account

```sh
ccam add
```

La procedura guidata chiede un nome come `personal` o `work` e quali elementi l'account deve condividere con gli altri tuoi account. Poi esegue `claude auth login`, che apre il browser per farti accedere. Al termine vedi:

```console
✓ personal (me@example.com) 추가됨 — 단축 명령: claude1
```

Per ora ccam stampa i propri messaggi in coreano; i comandi e le opzioni sono gli stessi in ogni lingua.

## 2. Aggiungi il secondo account

```sh
ccam add work
```

Accedi con l'altro account nel browser. Diventa `claude2`.

## 3. Cambia account

Ognuno di questi comandi avvia Claude Code con l'account di lavoro:

```sh
claude2
ccs work
ccs 2
```

Tutto ciò che segue il nome dell'account arriva a Claude Code invariato, quindi continui a usare le opzioni che conosci:

```sh
ccs work --resume
ccs work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

Eseguire `ccs` da solo apre un selettore: usa i tasti freccia e Invio, oppure premi il numero dell'account.

## 4. Scegli cosa condividere

Per impostazione predefinita, i server MCP, `settings.json`, il `CLAUDE.md` globale e agenti, comandi, skill e plugin sono condivisi tra gli account. La cronologia delle conversazioni no. Attivala per entrambi gli account se vuoi continuare le conversazioni passando dall'uno all'altro:

```sh
ccam share personal --history on
ccam share work --history on
```

ccam mostra cosa sposterà prima di modificare qualsiasi cosa. Vedi [condivisione](/docs/sharing) per i dettagli.

## 5. Quando un account raggiunge il limite

```sh
ccs usage
ccs next --continue
```

`ccs usage` mostra l'utilizzo su 5 ore e settimanale di ogni account. `ccs next --continue` apre la stessa conversazione nell'account successivo. Vedi [limiti di utilizzo](/docs/usage-limits).

## Usi già ~/.claude?

Registra quella cartella invece di accedere di nuovo:

```sh
ccam import
```

La cartella resta dove si trova e `claude` da solo continua a funzionare come prima.
