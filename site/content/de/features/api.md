---
title: API-Schlüssel-Konten
description: "ccam add --api legt ein Konto mit API-Schlüssel für die API von Anthropic oder ein Gateway an; der Schlüssel kommt beim Start aus Variable oder Befehl."
updated: 2026-10-20
order: 13
translation: machine
sourceHash: "a903c92f2c50"
---

Ein API-Schlüssel-Konto startet Claude Code mit einem API-Schlüssel statt mit der Anmeldung eines Claude-Abos: die API von Anthropic oder ein Endpunkt, der dasselbe Messages-Format spricht, etwa ein LiteLLM- oder OpenRouter-Gateway. ccam speichert den Schlüssel nicht. Es liest ihn bei jedem Start des Kontos aus einer Umgebungsvariable oder über einen Befehl, etwa einen Passwortmanager, und gibt ihn nur an diesen Claude Code-Prozess weiter.

## Verwendung

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
   #  Name      Shortcut  Email                   Groups   Last used
*  1  personal  -         me@example.com          default  2h ago
   2  gw        -         (API: llm.example.com)  default  -
```

Ein API-Schlüssel-Konto tritt Gruppen bei wie jedes andere Konto (standardmäßig `default`: Einstellungen, `CLAUDE.md`, Erweiterungen und MCP-Server) und funktioniert mit Kurzbefehlen, `ccam folder` und `ccam ext`.

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `--api` | Ein API-Schlüssel-Konto hinzufügen, statt sich anzumelden | aus |
| `--key-from env:NAME` | Den Schlüssel aus der Umgebungsvariable `NAME` lesen | |
| `--key-from command:COMMAND` | `COMMAND` ausführen (mit `sh -c` oder unter Windows `cmd /C`) und die erste ausgegebene Zeile verwenden | |
| `--base-url <url>` | Der Endpunkt (`ANTHROPIC_BASE_URL`) | Die API von Anthropic |
| `--auth api-key\|bearer` | Den Schlüssel als `ANTHROPIC_API_KEY` (`x-api-key`) oder als `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, was Gateways meist erwarten) senden | `api-key` |
| `--model <model>` | Das Modell (`ANTHROPIC_MODEL`) | Standard von Claude Code |

Die Einstellungen werden in `config.toml` gespeichert:

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

## So funktioniert es

Bei jedem Start des Kontos, nach den üblichen Prüfungen und bevor Claude Code startet:

1. ccam liest den Schlüssel: `env:` aus der Umgebung, `command:` durch Ausführen des Befehls mit einem Limit von 30 Sekunden. Ein Befehl kann im Terminal nach einem Passwort fragen; auch seine Fehlerausgabe geht an das Terminal. Ein leerer Schlüssel stoppt den Start mit einer Meldung, wo der Schlüssel erwartet wurde.
2. Es entfernt jede geerbte Variable, die bestimmt, wie sich Claude Code anmeldet oder welchen Anbieter es verwendet: die Schlüssel und Tokens (`ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN_FILE_DESCRIPTOR`, `CLAUDE_CODE_API_KEY_FILE_DESCRIPTOR`), den Endpunkt (`ANTHROPIC_BASE_URL`, `ANTHROPIC_UNIX_SOCKET`, `ANTHROPIC_MODEL`) und die Anbieterschalter (`CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX`, `CLAUDE_CODE_USE_FOUNDRY`, `CLAUDE_CODE_USE_ANTHROPIC_AWS`, `CLAUDE_CODE_USE_ANTHROPIC_GOOGLE_CLOUD`, `CLAUDE_CODE_USE_MANTLE`, `CLAUDE_CODE_USE_GATEWAY`). Ein Schlüssel, der in deiner Shell für etwas anderes gesetzt ist, kann dieses Konto nicht erreichen, und eine Abo-Anmeldung kann sich nicht gegen seinen Schlüssel durchsetzen.
3. Es setzt `ANTHROPIC_BASE_URL` und `ANTHROPIC_MODEL`, wenn angegeben, und den Schlüssel als `ANTHROPIC_API_KEY` oder `ANTHROPIC_AUTH_TOKEN`.

Konten mit Abo-Anmeldung behalten wie bisher die Umgebung, die sie erben. In einer interaktiven Sitzung fragt Claude Code möglicherweise einmal, ob es den in der Umgebung gefundenen Schlüssel verwenden soll.

## Sicherheit und Sonderfälle

- Der Schlüssel wird nie in `config.toml`, einen Cache oder ein Protokoll geschrieben, und `ccam --dry-run gw` liest ihn nicht: Eine `command:`-Quelle wird bei einem Probelauf nicht ausgeführt.
- `ccam add --api` prüft eine `env:`-Quelle sofort und warnt, wenn die Variable leer ist; eine `command:`-Quelle wird erst beim Start ausgeführt.
- `ccam login` und `ccam logout` lehnen ein API-Schlüssel-Konto ab, da es keine Anmeldung hat. `ccam doctor` zeigt seinen Endpunkt und seine Schlüsselquelle und warnt, wenn eine `env:`-Variable leer ist.
- `ccam usage` zeigt dafür `API`. `ccam next --auto`, `--fallback` und `ccam pick` wählen nur unter Abo-Konten, weil ein API-Schlüssel kein 5-Stunden- oder Wochenlimit zum Vergleichen hat.
- Ein geteiltes `settings.json` mit einem `env`-Block oder `apiKeyHelper` gilt für jedes Konto der Gruppe, API-Schlüssel-Konten eingeschlossen. Halte Anbietereinstellungen aus geteilten Einstellungen heraus, oder nimm das API-Schlüssel-Konto in eine Gruppe ohne `settings` auf.

## Verwandte Themen

- [Konten verwalten](/features/accounts)
- [Claude Code starten](/features/launch)
- [Befehlsreferenz](/docs/commands)
