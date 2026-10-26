---
title: Sprawdzanie rozszerzeń
description: ccam ext pokazuje skille, polecenia, pluginy, serwery MCP i hooki ładowane przez każde konto, skąd każde z nich pochodzi i co jest zepsute.
updated: 2026-10-26
order: 14
translation: machine
sourceHash: "7ca357a7f38a"
---

`ccam ext` odpowiada na pytanie „Dlaczego ta skill (lub ten serwer MCP) działa na jednym koncie, a na drugim nie?” Dla każdego konta wyświetla rozszerzenia ładowane przez Claude Code i skąd każde z nich pochodzi: z grupy współdzielenia, z własnego folderu konta albo z projektu, w którym jesteś. Pod tabelą wymienia, co przeszkadza rozszerzeniu załadować się tak, jak oczekujesz. Niczego nie zmienia.

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
  deploy  default   default

MCP servers
       personal   work
  ctx  default !  default

Hooks (commands per event)
              personal     work
  PreToolUse  1 (default)  1 (default)

This project (/home/me/code/app):
  commands: ship
  MCP servers: fs

Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.
! personal: MCP server ctx is defined differently in group default's mcp.json and .claude.json
! personal: half: no SKILL.md, so Claude Code does not load it
```

Komórka podaje źródło danego elementu na danym koncie: grupę współdzielenia (`default`), `account`, gdy element leży we własnym folderze konta, `local` dla serwera MCP, który konto ma tylko dla tego folderu, albo `-`, gdy konto go nie ma. Plugin zainstalowany, ale niewłączony, pokazuje `(off)` po swoim źródle. `!` oznacza element z problemem, wyjaśnionym pod tabelą; problemy, które nie należą do żadnego konta, na przykład serwery, które ma tylko zwykłe `claude`, są na końcu bez nazwy. `ccam ext --json` wypisuje ten sam raport (`accounts`, projekt, `problems`), z `remote` przy serwerach MCP, z którymi Claude Code łączy się przez sieć; bez kont wypisuje raport z `accounts: []`.

## Co odczytuje

| Rodzaj | Gdzie | Pokazywane źródło |
| --- | --- | --- |
| Skille | każdy folder w `skills/` z plikiem `SKILL.md` | grupa, do której prowadzi `skills/`, albo `account` |
| Polecenia, agenci | pliki `.md` w `commands/` i `agents/`, podfoldery jako `folder/name` | grupa albo `account` |
| Pluginy | `plugins/installed_plugins.json`, włączane i wyłączane przez `enabledPlugins` w ustawieniach konta i projektu; plugin zainstalowany tylko dla jednego projektu pojawia się tylko w tym projekcie | grupa albo `account`; `(off)`, gdy nie jest włączony |
| Serwery MCP | `mcp.json` grupy (gdy MCP jest współdzielone), `.claude.json` konta (jego własne serwery i te, które ma tylko dla tego folderu) oraz serwery jego pluginów jako `plugin:<plugin>:<server>` | grupa, `account` albo `local` |
| Hooki | `hooks` w `settings.json`, liczone dla każdego zdarzenia | grupa albo `account` |
| Ten projekt | `.claude/skills`, `.claude/commands`, `.claude/agents`, `.claude/settings.json`, `.claude/settings.local.json` i `.mcp.json` w bieżącym folderze | wymienione raz; konto, które nie zaufało folderowi, dostaje problem dla serwerów z jego `.mcp.json` |

## Zgłaszane problemy

- **Zerwane dowiązania**: dowiązanie w `skills/`, `commands/` lub `agents/`, którego cel zniknął, albo współdzielony folder (`skills`, `commands`, `agents`, `plugins`) prowadzący do nieistniejącego miejsca (napraw to przez `ccam doctor --fix`, które odtwarza brakujący folder grupy).
- **Folder skill bez `SKILL.md`**: Claude Code go pomija.
- **Serwer MCP zdefiniowany inaczej w dwóch miejscach**: w `mcp.json` grupy i w `.claude.json` konta albo w jednym z nich i w `.mcp.json` projektu. Identyczne definicje w obu miejscach (takie, jakie zapisuje `mcp_strategy = "sync"`) nie są problemem. Przenieś własny serwer konta do grupy przez `ccam mcp import <account> <server>`.
- **`installed_plugins.json`, który nie jest poprawnym JSON**, plugin, którego `enabledPlugins` nie wymienia, i plugin, którego folderu brakuje.
- **Serwery, które ładuje tylko jedno konto**: serwer tylko w `.claude.json` tego konta, którego pozostałe konta jego grupy nie ładują, oraz serwer w lokalnym zakresie konta dla tego folderu (tak działa `claude mcp add` bez `-s user`). Każdy ma podane polecenie `ccam mcp import`, które go udostępnia. Serwery w `~/.claude.json` zwykłego `claude`, których nie ładuje żadne konto, mają podane polecenie `ccam migrate --mcp`.
- **Serwer projektu czekający na zatwierdzenie**: konto jeszcze nie zaufało folderowi, więc Claude Code zapyta raz, a do tego czasu serwery z `.mcp.json` projektu czekają. Problem podaje, gdzie je zatwierdzić: `"enabledMcpjsonServers"` w `settings.json` kont grupy ustawień albo tego konta. Zaufanie nadane głównemu drzewu roboczemu repozytorium git obowiązuje też w jego pozostałych drzewach roboczych.
- **Logowania i konektory**: serwer MCP, który na tym koncie wymaga zalogowania, podczas gdy inne konto jest do niego zalogowane (logowanie obowiązuje tylko dla jednego konta; zaloguj się przez `/mcp`), oraz konektory claude.ai, których użyło inne konto, a to jeszcze nie.

## Bezpieczeństwo i przypadki brzegowe

- `ccam ext` niczego nie zmienia i nie wysyła zapytań sieciowych. Odczytuje nazwy serwerów i liczbę hooków, nigdy wartości `env` ani nagłówków w definicjach MCP. Na potrzeby sprawdzenia logowań odczytuje nazwy logowań MCP w loginie każdego konta (na macOS z pęku kluczy), nigdy tokeny, i tylko dla kont, które ładują serwer ładowany też przez inne konto.
- Pokazuje to, co jest na dysku. Serwer zablokowany przez ustawienia zarządzane nadal się pojawia; plugin wyłączony w ustawieniach pokazuje `(off)`.
- Uruchom je w folderze projektu, żeby uwzględnić folder `.claude` i `.mcp.json` tego projektu.

## Powiązane

- [Współdzielenie przez grupy](/features/share)
- [Współdzielone serwery MCP](/features/mcp)
- [Sprawdzanie i naprawa](/features/doctor)
