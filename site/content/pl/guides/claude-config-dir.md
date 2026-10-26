---
title: CLAUDE_CONFIG_DIR wyjaśnione
description: Co robi CLAUDE_CONFIG_DIR w Claude Code, co jest w tym folderze, gdzie macOS trzyma login i które pliki można współdzielić między kontami.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: "dc2fe42e11e1"
---

`CLAUDE_CONFIG_DIR` mówi Claude Code, którego folderu użyć zamiast `~/.claude`. Znajduje się w nim wszystko, co dotyczy konta: login, ustawienia, serwery MCP, wtyczki i historia rozmów. Wskazanie różnych folderów daje osobne konta.

```sh
CLAUDE_CONFIG_DIR=~/.claude-work claude
```

## Czy CLAUDE_CONFIG_DIR jest oficjalnie obsługiwane?

Tak. Claude Code wymienia tę zmienną wśród swoich [zmiennych środowiskowych](https://code.claude.com/docs/en/env-vars) jako sposób na zmianę folderu konfiguracji, „przydatny do uruchamiania kilku kont obok siebie”, a jego [przewodnik po uwierzytelnianiu](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts) używa jej do logowania na kilka kont. Możesz ją ustawić w powłoce lub w ustawieniach użytkownika, ale nie w ustawieniach projektu.

## Co jest w folderze?

Sprawdzono z Claude Code 2.1.309:

| Ścieżka | Zawartość | Można współdzielić |
| --- | --- | --- |
| `.credentials.json` | Tokeny logowania (Linux, Windows), tokeny OAuth dla MCP | Nie |
| `.claude.json` | Dane konta, serwery MCP użytkownika, stan poszczególnych projektów | Nie (serwery MCP osobno) |
| `settings.json` | Uprawnienia, hooki, model, pasek stanu, wtyczki | Tak |
| `CLAUDE.md` | Globalne instrukcje | Tak |
| `agents/`, `commands/`, `skills/`, `plugins/` | Rozszerzenia | Tak |
| `projects/` | Rozmowy jako pliki `.jsonl`, pamięć poszczególnych projektów | Tak |
| `history.jsonl` | Historia promptów | Tak |
| `sessions/`, `shell-snapshots/`, `cache/` i inne | Stan w czasie działania | Nie |

## Co najczęściej zaskakuje?

**Bez `CLAUDE_CONFIG_DIR` plik `.claude.json` nie znajduje się w `~/.claude`.** Claude Code trzyma go wtedy w `~/.claude.json`, obok folderu. Jeśli zaczniesz używać `CLAUDE_CONFIG_DIR=~/.claude`, Claude Code będzie szukać `~/.claude/.claude.json` i nie znajdzie twojego konta. `ccam import` rozwiązuje to, linkując oba pliki. Jeśli `~/.claude` ma już własny `.claude.json`, z wcześniejszego uruchomienia z `CLAUDE_CONFIG_DIR=~/.claude`, ccam zostawia oba pliki i o tym informuje, bo tylko ty wiesz, który jest właściwy.

**Na macOS login jest w pęku kluczy, a nie w folderze.** Claude Code zapisuje go pod nazwą `Claude Code-credentials`. Gdy ustawiono `CLAUDE_CONFIG_DIR`, nazwa kończy się dodatkowo `-` i pierwszymi 8 cyframi szesnastkowymi skrótu SHA-256 ścieżki folderu. Nazwa konta elementu pęku kluczy to `$USER`. Osobne foldery dostają więc osobne wpisy w pęku kluczy, a folder przeniesiony lub przemianowany traci login; dlatego `ccam rename` na macOS zachowuje nazwę folderu. Jeśli pęk kluczy odmówi zapisu, Claude Code zapisuje login w `.credentials.json` w folderze, jak opisuje [jego dokumentacja](https://code.claude.com/docs/en/authentication#credential-management).

**Nie każdy login jest przechowywany w folderze.** Logowanie do Claude Console bez klucza API jest zapisywane poza nim, więc osobne foldery nie rozdzielą dwóch takich logowań. Loginy claude.ai i klucze API są przechowywane osobno dla każdego folderu.

## Gdzie jest folder w Windows?

Domyślnie jest to `%USERPROFILE%\.claude`. Login jest w pliku `.credentials.json` wewnątrz albo w Menedżerze poświadczeń Windows, jeśli Claude Code jest ustawiony, żeby go używać, pod nazwą wyprowadzoną ze ścieżki folderu, tak jak na macOS. Dla każdego dodatkowego konta ustaw `CLAUDE_CONFIG_DIR` na inny folder, tak samo jak na Linux i macOS.

## Czy foldery mogą współdzielić pliki?

Dowiązania symboliczne działają: Claude Code czyta i zapisuje przez podlinkowany `settings.json` lub folder `projects/` i zachowuje link. Wyjątkiem są serwery MCP, bo dzielą `.claude.json` z loginem; zamiast tego przekaż je przez `--mcp-config`.

ccam opiera się dokładnie na tych zasadach. Zobacz [współdzielenie](/docs/sharing), żeby dowiedzieć się, jak linkuje foldery, oraz [poradnik o wielu kontach](/guides/multiple-claude-code-accounts), żeby zobaczyć szerszy obraz.
