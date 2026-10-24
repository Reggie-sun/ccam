---
title: Współdzielenie serwerów MCP, ustawień i historii
description: Jak grupy ccam współdzielą serwery MCP, settings.json, CLAUDE.md, umiejętności i historię rozmów między kontami Claude Code, przy osobnych loginach.
updated: 2026-10-23
order: 4
translation: machine
sourceHash: "27e0a0e5abbd"
---

Każde konto ma własny login. Wszystko inne można współdzielić przez **grupy**: grupa wymienia elementy, które współdzieli, i należące do niej konta, a ccam linkuje pliki każdego z tych kont do folderu grupy, więc zmiana wprowadzona na jednym koncie pojawia się na wszystkich.

## Co można współdzielić

| Element | Co jest współdzielone |
| --- | --- |
| `mcp` | serwery MCP: `mcp.json` grupy, przekazywany do Claude Code przy starcie |
| `settings` | `settings.json`: uprawnienia, hooki, model, pasek stanu |
| `memory` | globalny `CLAUDE.md` |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | historia rozmów: `projects/` i `history.jsonl` |

Pierwsze konto tworzy grupę `default`, która współdzieli wszystko oprócz historii rozmów, a kolejne konta do niej dołączają.

Pliki logowania `.credentials.json` i `.claude.json` nigdy nie są współdzielone. ccam linkuje tylko ścieżki ze stałej listy, a tych dwóch na niej nie ma.

Ta strona opisuje codzienne zadania. Wszystkie opcje są opisane na stronach funkcji: [współdzielenie przez grupy](/features/share), [współdzielone serwery MCP](/features/mcp) i [kopiowanie istniejącej konfiguracji](/features/migrate).

## Grupy

```sh
ccam group                                   # the groups, what they share, their accounts
ccam group set default --add history         # every account in default shares history
ccam group add team --items history --members work,side
```

Konto może należeć do kilku grup, ale każdy element pochodzi tylko z jednej grupy. Żeby więc współdzielić historię między `work` i `side`, ale nie z pozostałymi kontami, nie dodawaj historii do `default` i daj tym dwóm kontom osobną grupę, jak powyżej. `ccam share work` pokazuje, co konto współdzieli i przez którą grupę. [Panel](/features/dashboard) robi to samo na karcie Groups.

## Gdy konto zaczyna współdzielić

Gdy konto zaczyna współdzielić element (dołącza do grupy albo jego grupa dodaje ten element), ccam:

1. wypisuje plan wszystkich plików, które przeniesie, skopiuje lub podlinkuje dla każdego konta, i pyta o zgodę przed rozpoczęciem;
2. przenosi pliki konta do folderu grupy, jeśli grupa jest jeszcze pusta, albo dodaje pliki, których grupa jeszcze nie ma;
3. gdy ten sam plik się różni, pyta, czy zachować wersję grupy, użyć twojej, czy scalić JSON;
4. zachowuje oryginał w `.ccam-backup/<time>/` w folderze konta;
5. zastępuje oryginał linkiem do grupy.

Pliki rozmów mają nazwy według identyfikatora sesji, więc historie z kilku kont łączą się bez konfliktów.

ccam nie zmienia konta, gdy Claude Code działa jako to konto, niezależnie od tego, jak go uruchomiono (także sesje w tle), bo sesja dalej zapisuje do przenoszonych plików. Najpierw ją zakończ. ccam sprawdza to ponownie tuż przed zapisem, więc wychwytuje też sesję uruchomioną w czasie, gdy plan czekał.

## Gdy konto przestaje współdzielić

```sh
ccam group leave team work                   # work leaves the group
ccam group set default --remove history      # nobody in default shares history any more
```

Link zostaje usunięty, a konto dostaje własną kopię bieżącej zawartości grupy; serwery MCP grupy trafiają do jego własnego `.claude.json`, obok jego własnych serwerów. Dodaj `--fresh`, żeby zamiast tego zacząć od pustego folderu; plan wymienia wtedy serwery MCP, które przestaną się ładować. Zaimportowany folder, taki jak `~/.claude`, zawsze zachowuje kopię. Grupa zachowuje swoją zawartość dla pozostałych kont.

## Kopiowanie istniejącej konfiguracji

Jeśli używałeś Claude Code przed ccam, twoje ustawienia są w `~/.claude`. `ccam migrate` kopiuje je do grupy, żeby dostało je każde konto współdzielące te elementy:

```sh
ccam migrate
```

W terminalu ccam najpierw pokazuje listę elementów do zaznaczenia. Zaznaczone jest wszystko oprócz historii rozmów; naciśnij Spację, żeby zmienić element, i Enter, żeby zobaczyć plan. Elementy, dla których folder nic nie zawiera, są oznaczone jako `not found`.

Co jest kopiowane:

- `settings.json` i globalny `CLAUDE.md`;
- rozszerzenia `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/` i `rules/`;
- twoje serwery MCP użytkownika, czyli wpis `mcpServers` z `~/.claude.json`, do pliku grupy `mcp.json`;
- historia rozmów (`projects/` i `history.jsonl`), ale tylko gdy dodasz `--history`.

Folder źródłowy jest tylko odczytywany: nic w nim nie jest przenoszone, linkowane ani zmieniane, a zwykłe `claude` działa jak dotąd. Pliki logowania `.credentials.json` i `.claude.json` nigdy nie są kopiowane do grupy.

Serwery MCP logujące się przez OAuth przechowują swoje tokeny w loginie, a nie na liście serwerów. ccam nie kopiuje tych logowań: skopiowane logowanie zawiera ten sam token odświeżania, a pierwsze odnowienie wylogowuje z tego serwera wszystkich pozostałych posiadaczy, w tym `~/.claude`. Plan i podsumowanie wymieniają konta, które same logują się do serwerów przez `/mcp`.

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
Group default: personal
To share them:
  ccam group join default work
```

Na końcu ccam wypisuje konta grupy. Jeśli grupa nie współdzieli jeszcze któregoś z elementów albo niektóre konta do niej nie należą, podaje polecenia `ccam group`, które to zmieniają.

Gdy grupa ma już inną kopię pliku, ccam pyta tak samo jak przy włączaniu współdzielenia: zachować kopię grupy czy użyć twojej. `settings.json` i inne pliki JSON można też scalić, a `history.jsonl` scalić wiersz po wierszu. Poprzednie kopie grupy są zapisywane w `.ccam-backup/<time>/` w folderze grupy.

| Opcja | Działanie |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Kopiuje tylko wskazane elementy, bez listy wyboru; bez nich i poza terminalem wszystko oprócz historii rozmów |
| `-g`, `--group <group>` | Kopiuje do tej grupy zamiast do `default` |
| `-y`, `--yes` | Kopiuje bez pytania; gdy plik się różni, zachowuje kopię grupy |

Umieść `--dry-run` przed poleceniem, jak w `ccam --dry-run migrate`, aby zobaczyć plan bez kopiowania czegokolwiek. Aby skopiować z innego folderu, podaj go: `ccam migrate ~/.claude-work --group work`.

## Serwery MCP

Claude Code przechowuje serwery MCP w `.claude.json` razem z loginem, więc tego pliku nie można współdzielić. ccam trzyma współdzielone serwery w `mcp.json` grupy i uruchamia Claude Code z `--mcp-config=<group>/mcp.json`. Twoje własne opcje `--mcp-config` nadal obowiązują i są stosowane dodatkowo.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

Serwer dodany w trakcie sesji przez `claude mcp add` jest zapisywany tylko dla tego konta i dla folderu, w którym byłeś (zakres lokalny Claude Code; `-s user` przypisuje go kontu we wszystkich folderach). Przy następnym uruchomieniu tego konta ccam poda nazwę serwera wraz z jego folderem i polecenie, które przeniesie go do grupy:

```sh
ccam mcp import work gh
```

Jeśli grupa ma już inny serwer o tej nazwie, ccam go pomija i o tym informuje. `--force` zastępuje serwer grupy i zachowuje stary `mcp.json` w folderze `.ccam-backup/` grupy.

Serwery MCP logujące się przez OAuth przechowują swoje tokeny razem z loginem konta, więc autoryzuj je raz na każdym koncie przez `/mcp`. `ccam ext` i `ccam doctor` wymieniają serwery, do których konto nie jest zalogowane, a inne konto jest.

Jeśli wolisz, żeby serwery pojawiały się w `claude mcp list` i `/mcp` jako zwykłe serwery użytkownika, ustaw `mcp_strategy = "sync"` na górze `~/.config/ccam/config.toml`, nad pierwszym `[[group]]`, `[[profile]]` lub `[[folder]]`. ccam zapisuje je wtedy w `.claude.json` każdego konta przed startem, z użyciem blokady pliku samego Claude Code i z kopią zapasową. Własny serwer konta o tej samej nazwie pozostaje serwerem konta, a ccam o tym informuje.

## Dwa konta w tym samym folderze

Konta, które współdzielą historię rozmów, zapisują do tych samych plików. Jeśli uruchomisz drugi Claude Code w folderze, w którym działa już inne konto z tej samej grupy, ccam ostrzeże cię i zapyta, czy kontynuować.

## Gdy coś się zepsuje

```sh
ccam doctor
ccam doctor --fix
```

Jeśli współdzielona ścieżka stała się zwykłym plikiem, na przykład po ręcznym przywróceniu kopii zapasowej, ccam zatrzymuje się przed uruchomieniem Claude Code i wskazuje `ccam doctor --fix`, które ponownie łączy go z grupą: plik staje się kopią grupy, gdy grupa nie ma własnej, a w przeciwnym razie grupa zachowuje swoją, a twój trafia do kopii zapasowej, którą `--fix` wskazuje. Zmianę grup, która zatrzymała się w połowie, także kończy `ccam doctor --fix`.
