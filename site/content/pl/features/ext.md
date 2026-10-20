---
title: Sprawdzanie rozszerzeń
description: ccam ext pokazuje skille, polecenia, pluginy, serwery MCP i hooki ładowane przez każde konto, skąd każde z nich pochodzi i co jest zepsute.
updated: 2026-10-20
order: 14
translation: machine
sourceHash: 78d83272963c
---

`ccam ext` odpowiada na pytanie „dlaczego ta skill (lub ten serwer MCP) działa na jednym koncie, a na drugim nie?”. Dla każdego konta wyświetla rozszerzenia ładowane przez Claude Code i skąd każde z nich pochodzi: z grupy współdzielenia, z własnego folderu konta albo z projektu, w którym jesteś. Pod tabelą wymienia, co przeszkadza rozszerzeniu załadować się tak, jak oczekujesz. Tylko odczytuje.

## Użycie

```sh
ccam ext              # every account
ccam ext work         # one account
ccam ext --json       # for scripts
```

```console
$ ccam ext
Skills
          personal   work
  half    account !  -
  review  -          default

Commands
          personal  work
  deploy  -         default

MCP servers
       personal  work
  ctx  account   default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app), for every account:
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), - (not loaded).
! personal: half: no SKILL.md, so Claude Code does not load it
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
```

Komórka podaje źródło danego elementu na danym koncie: grupę współdzielenia (`default`), `account`, gdy element leży we własnym folderze konta, albo `-`, gdy konto go nie ma. `!` oznacza element z problemem, wyjaśnionym pod tabelą.

## Co odczytuje

| Rodzaj | Gdzie | Pokazywane źródło |
| --- | --- | --- |
| Skille | każdy folder w `skills/` z plikiem `SKILL.md` | grupa, do której prowadzi `skills/`, albo `account` |
| Polecenia, agenci | pliki `.md` w `commands/` i `agents/`, podfoldery jako `folder/name` | grupa albo `account` |
| Pluginy | `plugins/installed_plugins.json` | grupa albo `account` |
| Serwery MCP | `mcp.json` grupy (gdy MCP jest współdzielone) i `.claude.json` konta | grupa albo `account` |
| Hooki | `hooks` w `settings.json`, liczone dla każdego zdarzenia | grupa albo `account` |
| Ten projekt | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` i `.mcp.json` w bieżącym folderze | wymienione raz, dla wszystkich kont |

## Zgłaszane problemy

- **Zerwane dowiązania**: dowiązanie w `skills/`, `commands/` lub `agents/`, którego cel zniknął, albo współdzielony folder (`skills`, `commands`, `agents`, `plugins`) prowadzący do nieistniejącego miejsca (napraw to przez `ccam doctor --fix`).
- **Folder skill bez `SKILL.md`**: Claude Code go pomija.
- **Serwer MCP zdefiniowany inaczej w dwóch miejscach**: w `mcp.json` grupy i w `.claude.json` konta albo w jednym z nich i w `.mcp.json` projektu. Te same ustawienia w obu miejscach (takie, jakie zapisuje `mcp_strategy = "sync"`) nie są problemem. Przenieś własny serwer konta do grupy przez `ccam mcp import <account> <server>`.
- **`installed_plugins.json`, który nie jest poprawnym JSON.**

## Bezpieczeństwo i przypadki brzegowe

- `ccam ext` niczego nie zmienia i nie wysyła zapytań sieciowych. Odczytuje nazwy serwerów i liczbę hooków, nigdy wartości `env` ani nagłówków w definicjach MCP.
- Pokazuje to, co jest na dysku. Plugin zainstalowany, ale wyłączony w ustawieniach, albo serwer zablokowany przez ustawienia zarządzane nadal się pojawia.
- Uruchom je w folderze projektu, żeby uwzględnić folder `.claude` i `.mcp.json` tego projektu.

## Powiązane

- [Współdzielenie między kontami](/features/share)
- [Współdzielone serwery MCP](/features/mcp)
- [Sprawdzanie i naprawa](/features/doctor)
