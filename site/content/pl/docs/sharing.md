---
title: Współdzielenie serwerów MCP, ustawień i historii
description: Jak ccam współdzieli serwery MCP, settings.json, CLAUDE.md, agentów, umiejętności i historię rozmów między kontami Claude Code, zachowując osobne loginy.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: 0864ea87c3e9
---

Każde konto ma własny login. Wszystko inne można współdzielić element po elemencie: ccam linkuje plik lub folder konta do współdzielonego folderu grupy, więc zmiana wprowadzona na jednym koncie pojawia się na wszystkich kontach tej grupy.

## Co można współdzielić

| Element | Opcja | Co jest współdzielone |
| --- | --- | --- |
| Serwery MCP | `--mcp` | `mcp.json` grupy, przekazywany do Claude Code przy starcie |
| Ustawienia | `--settings` | `settings.json`: uprawnienia, hooki, model, pasek stanu |
| Pamięć globalna | `--memory` | `CLAUDE.md` |
| Rozszerzenia | `--extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| Historia rozmów | `--history` | `projects/` i `history.jsonl` |

Nowe konta współdzielą wszystko oprócz historii rozmów.

Pliki logowania `.credentials.json` i `.claude.json` nigdy nie są współdzielone. ccam linkuje tylko ścieżki ze stałej listy, a tych dwóch na niej nie ma.

## Grupy

`on` oznacza grupę `default`. Możesz też podać nazwę, żeby tylko niektóre konta współdzieliły dany element:

```sh
ccam share work --history team
ccam share side --history team
```

`work` i `side` współdzielą teraz rozmowy ze sobą, ale nie z pozostałymi kontami.

## Włączanie współdzielenia

Gdy konto zaczyna współdzielić element, ccam:

1. wypisuje plan wszystkich plików, które przeniesie, skopiuje lub podlinkuje, i pyta o zgodę przed rozpoczęciem;
2. przenosi pliki konta do folderu grupy, jeśli grupa jest jeszcze pusta, albo dodaje pliki, których grupa jeszcze nie ma;
3. gdy ten sam plik się różni, pyta, czy zachować wersję grupy, użyć twojej, czy scalić JSON;
4. zachowuje oryginał w `.ccam-backup/<time>/` w folderze konta;
5. zastępuje oryginał linkiem do grupy.

Pliki rozmów mają nazwy według identyfikatora sesji, więc historie z kilku kont łączą się bez konfliktów.

ccam nie zmienia współdzielenia, gdy Claude Code działa jako to konto, bo sesja dalej zapisuje do przenoszonych plików. Najpierw ją zakończ.

## Wyłączanie współdzielenia

```sh
ccam share work --history off
```

Link zostaje usunięty, a konto dostaje własną kopię bieżącej zawartości grupy. Dodaj `--fresh`, żeby zamiast tego zacząć od pustego folderu. Grupa zachowuje swoją zawartość dla pozostałych kont.

## Kopiowanie istniejącej konfiguracji

Jeśli używałeś Claude Code przed ccam, twoje ustawienia są w `~/.claude`. `ccam migrate` kopiuje je do grupy, żeby dostało je każde konto współdzielące te elementy:

```sh
ccam migrate
```

Co jest kopiowane:

- `settings.json` i globalny `CLAUDE.md`;
- rozszerzenia `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/` i `rules/`;
- twoje serwery MCP użytkownika, czyli wpis `mcpServers` z `~/.claude.json`, do pliku grupy `mcp.json`;
- historia rozmów (`projects/` i `history.jsonl`), ale tylko gdy dodasz `--history`.

Folder źródłowy jest tylko odczytywany: nic w nim nie jest przenoszone, linkowane ani zmieniane, a zwykłe `claude` działa jak dotąd. Pliki logowania `.credentials.json` i `.claude.json` nigdy nie są kopiowane.

Zanim cokolwiek skopiuje, ccam wypisuje plan i pyta o zgodę:

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Accounts that share these through group default:
  mcp         personal
  settings    personal
  memory      personal
  extensions  personal
```

Na końcu ccam wypisuje konta, które współdzielą każdy element przez grupę. Jeśli niektóre konta tego nie robią, podaje polecenia `ccam share`, które włączają dla nich współdzielenie.

Gdy grupa ma już inną kopię pliku, ccam pyta tak samo jak przy włączaniu współdzielenia: zachować kopię grupy czy użyć twojej. `settings.json` i inne pliki JSON można też scalić, a `history.jsonl` scalić wiersz po wierszu. Poprzednie kopie grupy są zapisywane w `.ccam-backup/<time>/` w folderze grupy.

| Opcja | Działanie |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Kopiuje tylko wskazane elementy; bez nich wszystko oprócz historii rozmów |
| `-g`, `--group <group>` | Kopiuje do tej grupy zamiast do `default` |
| `-y`, `--yes` | Kopiuje bez pytania; gdy plik się różni, zachowuje kopię grupy |

Umieść `--dry-run` przed poleceniem, jak w `ccam --dry-run migrate`, aby zobaczyć plan bez kopiowania czegokolwiek. Aby skopiować z innego folderu, podaj go: `ccam migrate ~/.claude-work --group work`.

## Serwery MCP

Claude Code przechowuje serwery MCP w `.claude.json` razem z loginem, więc tego pliku nie można współdzielić. ccam trzyma współdzielone serwery w `mcp.json` grupy i uruchamia Claude Code z `--mcp-config=<group>/mcp.json`. Twoje własne opcje `--mcp-config` nadal obowiązują i są stosowane dodatkowo.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

Serwer dodany w trakcie sesji przez `claude mcp add` jest zapisywany tylko dla tego konta. Przy następnym uruchomieniu tego konta ccam poda nazwę serwera i polecenie, które przeniesie go do grupy:

```sh
ccam mcp import work gh
```

Jeśli grupa ma już inny serwer o tej nazwie, ccam go pomija i o tym informuje. `--force` zastępuje serwer grupy i zachowuje stary `mcp.json` w folderze `.ccam-backup/` grupy.

Serwery MCP logujące się przez OAuth przechowują swoje tokeny razem z loginem konta, więc autoryzuj je raz na każdym koncie.

Jeśli wolisz, żeby serwery pojawiały się w `claude mcp list` i `/mcp` jako zwykłe serwery użytkownika, ustaw `mcp_strategy = "sync"` w `~/.config/ccam/config.toml`. ccam zapisuje je wtedy w `.claude.json` każdego konta przed startem, z użyciem blokady pliku samego Claude Code i z kopią zapasową.

## Dwa konta w tym samym folderze

Konta, które współdzielą historię rozmów, zapisują do tych samych plików. Jeśli uruchomisz drugi Claude Code w folderze, w którym działa już inne konto z tej samej grupy, ccam ostrzeże cię i zapyta, czy kontynuować.

## Gdy coś się zepsuje

```sh
ccam doctor
ccam doctor --fix
```

Jeśli współdzielona ścieżka stała się zwykłym plikiem, na przykład po ręcznym przywróceniu kopii zapasowej, ccam zatrzymuje się przed uruchomieniem Claude Code i wskazuje `ccam doctor --fix`, które scala plik z grupą i zachowuje kopię zapasową.
