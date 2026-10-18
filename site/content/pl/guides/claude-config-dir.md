---
title: CLAUDE_CONFIG_DIR wyjaśnione
description: Co robi CLAUDE_CONFIG_DIR w Claude Code, co jest w tym folderze, gdzie macOS trzyma login i które pliki można współdzielić między kontami.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 62f178595254
---

`CLAUDE_CONFIG_DIR` mówi Claude Code, którego folderu użyć zamiast `~/.claude`. Znajduje się w nim wszystko, co dotyczy konta: login, ustawienia, serwery MCP, wtyczki i historia rozmów. Wskazanie różnych folderów daje osobne konta.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## Co jest w folderze

Sprawdzono z Claude Code 2.1.302:

| Ścieżka | Zawartość | Można współdzielić |
| --- | --- | --- |
| `.credentials.json` | Tokeny logowania (Linux), tokeny OAuth dla MCP | Nie |
| `.claude.json` | Dane konta, serwery MCP użytkownika, stan poszczególnych projektów | Nie (serwery MCP osobno) |
| `settings.json` | Uprawnienia, hooki, model, pasek stanu, wtyczki | Tak |
| `CLAUDE.md` | Globalne instrukcje | Tak |
| `agents/`, `commands/`, `skills/`, `plugins/` | Rozszerzenia | Tak |
| `projects/` | Rozmowy jako pliki `.jsonl`, pamięć poszczególnych projektów | Tak |
| `history.jsonl` | Historia promptów | Tak |
| `sessions/`, `shell-snapshots/`, `cache/` i inne | Stan w czasie działania | Nie |

## Dwa szczegóły, które często zaskakują

**Bez `CLAUDE_CONFIG_DIR` plik `.claude.json` nie znajduje się w `~/.claude`.** Claude Code trzyma go wtedy w `~/.claude.json`, obok folderu. Jeśli zaczniesz używać `CLAUDE_CONFIG_DIR=~/.claude`, Claude Code będzie szukać `~/.claude/.claude.json` i nie znajdzie twojego konta. `ccam import` rozwiązuje to, linkując oba pliki.

**Na macOS login jest w pęku kluczy, a nie w folderze.** Claude Code zapisuje go pod nazwą `Claude Code-credentials`, po której, gdy ustawiono `CLAUDE_CONFIG_DIR`, następuje krótki skrót ścieżki folderu. Osobne foldery dostają więc osobne wpisy w pęku kluczy. Jeśli pęk kluczy odmówi zapisu, Claude Code zapisuje login w `.credentials.json` w folderze.

## Współdzielenie między folderami

Dowiązania symboliczne działają: Claude Code czyta i zapisuje przez podlinkowany `settings.json` lub folder `projects/` i zachowuje link. Wyjątkiem są serwery MCP, bo dzielą `.claude.json` z loginem; zamiast tego przekaż je przez `--mcp-config`.

ccam opiera się dokładnie na tych zasadach. Zobacz [współdzielenie](/docs/sharing), żeby dowiedzieć się, jak linkuje foldery, oraz [poradnik o wielu kontach](/guides/multiple-claude-code-accounts), żeby zobaczyć szerszy obraz.
