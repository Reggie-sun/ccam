---
title: Account con chiave API
description: ccam add --api aggiunge un account con chiave API per l'API di Anthropic o un gateway; la chiave si legge all'avvio da variabile o comando, mai salvata.
updated: 2026-10-20
order: 13
translation: machine
sourceHash: de2c5cde2c48
---

Un account con chiave API avvia Claude Code con una chiave API invece che con il login di un abbonamento Claude: l'API di Anthropic, oppure un endpoint che usa lo stesso formato Messages, come un gateway LiteLLM o OpenRouter. ccam non salva la chiave. La legge ogni volta che l'account si avvia, da una variabile d'ambiente o da un comando come un gestore di password, e la passa solo a quel processo di Claude Code.

## Utilizzo

```sh
ccam add api --api --key-from env:ANTHROPIC_KEY_WORK
ccam add gw --api --base-url https://llm.example.com --auth bearer \
  --model claude-sonnet-5-5 --key-from 'command:pass show llm/gateway-key'
ccam add mac --api --key-from 'command:security find-generic-password -s anthropic-api -w'
ccam gw -p "summarise the open issues"
ccam list
```

```console
$ ccam list
   #  Name      Shortcut  Email                    Shared               Last used
*  1  personal  -         me@example.com           mcp settings md ext  2h ago
   2  gw        -         (API: llm.example.com)   mcp settings md ext  -
```

Un account con chiave API condivide impostazioni, `CLAUDE.md`, estensioni e server MCP come qualsiasi altro account, e funziona con le scorciatoie, `ccam folder` e `ccam ext`.

## Opzioni

| Opzione | Cosa fa | Predefinito |
| --- | --- | --- |
| `--api` | Aggiunge un account con chiave API invece di fare il login | disattivata |
| `--key-from env:NAME` | Legge la chiave dalla variabile d'ambiente `NAME` | |
| `--key-from command:COMMAND` | Esegue `COMMAND` (con `sh -c`, oppure `cmd /C` su Windows) e usa la prima riga che stampa | |
| `--base-url <url>` | L'endpoint (`ANTHROPIC_BASE_URL`) | l'API di Anthropic |
| `--auth api-key\|bearer` | Invia la chiave come `ANTHROPIC_API_KEY` (`x-api-key`) oppure come `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, che di solito chiedono i gateway) | `api-key` |
| `--model <model>` | Il modello (`ANTHROPIC_MODEL`) | quello predefinito di Claude Code |

Le impostazioni vengono salvate in `config.toml`:

```toml
[[profile]]
name = 'gw'
slot = 2

[profile.api]
base_url = 'https://llm.example.com'
key_from = 'command:pass show llm/gateway-key'
auth = 'bearer'
model = 'claude-sonnet-5-5'
```

## Come funziona

A ogni avvio dell'account, dopo i soliti controlli e prima che Claude Code parta:

1. ccam legge la chiave: `env:` dall'ambiente, `command:` eseguendo il comando con un limite di 30 secondi. Un comando può chiedere una password nel terminale; anche il suo output di errore va al terminale. Una chiave vuota ferma l'avvio con un messaggio che dice dove era attesa la chiave.
2. Rimuove ogni variabile ereditata che decide come Claude Code accede o quale provider usa: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_BASE_URL`, `ANTHROPIC_MODEL`, `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX` e `CLAUDE_CODE_USE_FOUNDRY`. Una chiave impostata nella tua shell per altro non può arrivare a questo account, e il login di un abbonamento non può prevalere sulla sua chiave.
3. Imposta `ANTHROPIC_BASE_URL` e `ANTHROPIC_MODEL` quando indicati, e la chiave come `ANTHROPIC_API_KEY` o `ANTHROPIC_AUTH_TOKEN`.

Gli account con il login di un abbonamento mantengono l'ambiente che ereditano, come prima. In una sessione interattiva Claude Code può chiedere una volta se usare la chiave che trova nell'ambiente.

## Sicurezza e casi limite

- La chiave non viene mai scritta in `config.toml`, in una cache o in un log, e `ccam --dry-run gw` non la legge: un dry run non esegue una sorgente `command:`.
- `ccam add --api` controlla subito una sorgente `env:` e avvisa se la variabile è vuota; una sorgente `command:` viene eseguita per la prima volta all'avvio.
- `ccam login` e `ccam logout` rifiutano un account con chiave API, perché non ha un login. `ccam doctor` mostra il suo endpoint e la sorgente della chiave, e avvisa quando una variabile `env:` è vuota.
- `ccam usage` mostra `API` per questo account. `ccam next --auto`, `--fallback` e `ccam pick` scelgono solo tra gli account con abbonamento, perché una chiave API non ha limiti su 5 ore o settimanali da confrontare.
- Un `settings.json` condiviso con un blocco `env` o `apiKeyHelper` vale per ogni account del gruppo, compresi gli account con chiave API. Tieni le impostazioni del provider fuori dalle impostazioni condivise.

## Pagine correlate

- [Gestire gli account](/features/accounts)
- [Avviare Claude Code](/features/launch)
- [Riferimento dei comandi](/docs/commands)
