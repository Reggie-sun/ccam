---
title: Remote Control z ccam
description: "Remote Control na dowolnym koncie przez ccam <name> rc: które opcje działają, czemu ccam proponuje kopię współdzielonych serwerów MCP i jak naprawić problemy."
updated: 2026-10-20
order: 4
translation: machine
sourceHash: 59b802c44bb3
---

`ccam work rc` uruchamia Remote Control z Claude Code jako konto `work`, dzięki czemu możesz sterować sesjami na tym komputerze z claude.ai/code albo z mobilnej aplikacji Claude. Działa tak jak `claude rc`, z jedną różnicą: jeśli konto współdzieli serwery MCP, ccam raz pyta, czy skopiować je do `.claude.json` konta, bo sesje Remote Control nie odczytują serwerów z żadnego innego miejsca.

## Uruchamianie

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

[Skrót](/features/shortcuts) działa tak samo: `ccwork rc` to to samo co `ccam work rc`. Uruchom polecenie w folderze, w którym chcesz pracować, a potem otwórz claude.ai/code albo aplikację Claude zalogowaną na to samo konto Claude co `work`.

## Opcje

Spotykają się tu dwa rodzaje opcji: opcje ccam, które podaje się przed nazwą konta albo ustawia w konfiguracji ccam, oraz opcje Remote Control, które podaje się po `rc` i które trafiają do Claude Code bez zmian.

### Opcje i ustawienia ccam

| Opcja lub ustawienie | Działanie | Domyślnie |
| --- | --- | --- |
| `--dry-run` (przed kontem) | Wypisuje środowisko i polecenie oraz komunikat o MCP, jeśli występuje. Niczego nie pyta ani nie kopiuje | wyłączone |
| `Copy them into .claude.json (a backup is kept)? [y/N]` | Pytanie pojawia się, gdy konto współdzieli serwery MCP, a w jego `.claude.json` brakuje bieżącej listy grupy. `y` je kopiuje, `N` lub Enter uruchamia bez nich | `N` |
| `mcp_strategy = "sync"` w `~/.config/ccam/config.toml` | Kopiuje serwery grupy do `.claude.json` przed każdym uruchomieniem, niezależnie od tego, czy to Remote Control, więc nic nie jest pytane | `"flag"` |
| `ccam share work --mcp off` | Przestaje współdzielić serwery MCP z kontem; Remote Control startuje wtedy tylko z własnymi serwerami konta i o nic nie pyta | włączone |

### Opcje Remote Control

Wszystko po `rc` należy do Claude Code. Oto opcje `claude remote-control` w Claude Code 2.1.304; `ccam work rc --help` wypisuje listę dla Twojej wersji.

| Opcja | Działanie |
| --- | --- |
| `--name <name>` | Nazwa sesji wyświetlana w claude.ai/code |
| `--remote-control-session-name-prefix <prefix>` | Prefiks generowanych nazw sesji (domyślnie: nazwa hosta) |
| `-c`, `--continue` | Ponownie łączy się z sesją, którą Remote Control ostatnio zapisał dla tego folderu, w ciągu mniej więcej 4 godzin |
| `--session-id <id>` | Ponownie łączy się z konkretną sesją |
| `--permission-mode <mode>` | Tryb uprawnień nowych sesji: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` lub `plan` |
| `--spawn <mode>` | `same-dir` (domyślnie), `worktree` (osobny git worktree dla każdej sesji) lub `session` (jedna sesja, kończy działanie razem z nią) |
| `--capacity <N>` | Maksymalna liczba jednoczesnych sesji w trybach `same-dir` i `worktree` (domyślnie 32) |
| `--[no-]create-session-in-dir` | Tworzy na starcie pierwszą sesję w bieżącym folderze (domyślnie włączone) |
| `--[no-]chrome` | Claude in Chrome w nowych sesjach (domyślnie wyłączone) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | Wyjście diagnostyczne i szczegółowe |

Podawaj `rc` zaraz po nazwie konta. ccam traktuje pierwszy argument, który nie zaczyna się od `-`, jako polecenie Claude Code, więc w `ccam work rc --name laptop` widzi `rc`, ale w `ccam work --model opus rc` widzi `opus`.

## Co robi ccam

Uruchomienie przebiega zwykłą ścieżką ([Uruchamianie Claude Code](/features/launch)) z jedną zmianą w kroku MCP:

1. ccam znajduje konto, sprawdza jego linki współdzielenia i odnajduje `claude`, jak przy każdym uruchomieniu.
2. Widzi, że polecenie to `rc` lub `remote-control`, i pomija `--mcp-config`.
3. Jeśli konto współdzieli serwery MCP, porównuje `mcp.json` grupy z `sync-stamp` w `~/.local/state/ccam/mcp/work/`. Gdy `.claude.json` zawiera już bieżące serwery albo grupa nie ma żadnych, przechodzi od razu do kroku 6.
4. W przeciwnym razie wymienia serwery, których zabrakłoby w sesjach:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

5. Po odpowiedzi `y` zapisuje kopię zapasową `.claude.json`, bierze własną blokadę Claude Code, dodaje serwery grupy do `mcpServers` i zachowuje pozostałe serwery konta. Po odpowiedzi `N`, bez terminala albo gdy konto działa, niczego nie zmienia i informuje o tym:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. Rejestruje sesję, ustawia `CLAUDE_CONFIG_DIR` na folder konta i zastępuje się procesem `claude rc` przez `exec`.

Od tej chwili Remote Control działa jako serwer. Każda sesja, którą dla Ciebie uruchamia, jest procesem potomnym tego procesu, więc działa jako `work`: z loginem, ustawieniami i historią rozmów `work` oraz w ramach limitu użycia `work`. Dopóki serwer działa, ccam traktuje konto jako działające.

Próbne uruchomienie pokazuje różnicę. Zwykłe uruchomienie przekazuje serwery grupy jako argument, a Remote Control nie:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude rc
```

## Dlaczego to tak działa

**Remote Control nie przyjmuje `--mcp-config`.** Przy zwykłych uruchomieniach ccam przekazuje serwery grupy do Claude Code jako `--mcp-config=<path>` i nigdy nie zapisuje plików Claude Code. `claude remote-control` odrzuca tę opcję przed nazwą swojego polecenia i nie przyjmuje jej po niej, więc przed ccam 0.6.0 `ccam work rc` kończyło się błędem na każdym koncie, które współdzieliło serwery MCP. Teraz ccam pomija tę opcję.

**Jego sesje odczytują serwery tylko z `.claude.json`.** Sesje uruchamiane przez Remote Control nie dziedziczą opcji z wiersza poleceń. Jedynym miejscem, w którym znajdują serwery MCP użytkownika, jest `.claude.json` konta. Ten plik przechowuje też dane logowania konta i stan poszczególnych folderów, więc ccam nie może podlinkować go do współdzielonej kopii, tak jak robi to z `settings.json`. Skopiowanie do niego serwerów to jedyny sposób, by przekazać sesjom współdzieloną listę.

**ccam pyta przed zapisem.** Przy domyślnej strategii flagi ccam obiecuje nie ruszać `.claude.json`. Remote Control to jedyny przypadek, w którym musi tam coś zapisać, więc najpierw pyta, robi kopię zapasową i nie zmienia niczego poza `mcpServers`.

**Nie zapisuje, gdy konto działa.** Claude Code przepisuje `.claude.json` w trakcie sesji, a zmiana wprowadzona z zewnątrz zostałaby utracona albo nadpisana. Dlatego ccam pomija kopiowanie, gdy działa jakakolwiek sesja konta, w tym inny serwer Remote Control.

**Pyta raz przy każdej zmianie.** Po skopiowaniu `sync-stamp` przechowuje czas modyfikacji `mcp.json` grupy. ccam pyta ponownie dopiero wtedy, gdy ten plik się zmieni, na przykład po `ccam mcp add`. Jeśli odpowiesz `N`, kopii jeszcze nie ma, więc pytanie pojawia się przy każdym uruchomieniu Remote Control.

**`ccam mcp` utrzymuje kopie w aktualnym stanie.** Przy kolizji nazw wpis w `.claude.json` wygrywa z `--mcp-config`, więc stara kopia przesłoniłaby nowszą definicję także przy zwykłych uruchomieniach. Dlatego `ccam mcp add`, `add-json`, `remove` i `import` aktualizują kopie na każdym koncie, które nie działa:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**Bez terminala odpowiedź brzmi „nie”.** Remote Control często działa ze skryptu, jako usługa albo przez `nohup`, gdzie nikt nie może odpowiedzieć. ccam uruchamia się wtedy bez serwerów, zamiast czekać.

## Problemy i rozwiązania

### W sesjach Remote Control brakuje współdzielonych serwerów MCP

Na pytanie odpowiedziano `N`, nie było terminala albo konto działało. Zatrzymaj Remote Control, zamknij pozostałe sesje konta i ponownie uruchom `ccam work rc` w terminalu; odpowiedz `y`. Aby nigdy nie widzieć tego pytania, ustaw `mcp_strategy = "sync"` w `~/.config/ccam/config.toml`. Serwery w sesji nie zmieniają się w trakcie jej działania, więc po skopiowaniu rozpocznij nową sesję.

### "Account work is running, so its .claude.json is left alone"

Sesja `work` jest wciąż otwarta: zwykłe `ccam work`, skrót albo inny serwer Remote Control. Zakończ ją i ponownie uruchom `ccam work rc`. Sesje innego konta się nie liczą.

### ccam pyta za każdym razem

Serwery nie zostały jeszcze skopiowane. Odpowiedz raz `y` albo ustaw `mcp_strategy = "sync"`. Jeśli nigdy nie chcesz mieć tych serwerów w Remote Control, przestań współdzielić MCP z kontem: `ccam share work --mcp off`.

### `ccam work rc` od razu kończy się błędem dotyczącym `--mcp-config`

Twój ccam jest starszy niż 0.6.0 i przekazuje `--mcp-config` do Remote Control. Uruchom `ccam update`. Ten sam błąd w aktualnym ccam oznacza, że przed `rc` znalazła się opcja z wartością, jak w `ccam work --model opus rc`, więc ccam uznał `opus` za polecenie; podaj `rc` zaraz po nazwie konta.

### Zmiana wprowadzona przez `ccam mcp` nie dociera do Remote Control

`ccam mcp add` i `remove` aktualizują kopie tylko na kontach, które nie działają, i informują o tym: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Zatrzymaj Remote Control i uruchom go ponownie; ccam zauważy zmieniony `mcp.json` i znów zapyta.

### Kopia zastąpiła mój własny serwer

Jeśli konto miało własny serwer o tej samej nazwie co serwer grupy, kopia go zastąpiła. Poprzedni `.claude.json` jest w `~/.local/state/ccam/mcp/work/backups/`, gdzie przechowywanych jest pięć najnowszych kopii. Aby zachować swoją wersję dla wszystkich kont, przenieś ją do grupy przez `ccam mcp import work <server> --force` przed następnym uruchomieniem Remote Control.

### Chcę usunąć skopiowane serwery

`ccam mcp remove <server>` usuwa serwer z grupy i ze wszystkich kopii. Aby wrócić do stanu sprzed kopiowania, przywróć najnowszą kopię zapasową z `~/.local/state/ccam/mcp/work/backups/`, gdy nie działa żadna sesja konta.

### Remote Control pyta, czy ufać folderowi, albo kończy działanie bez terminala

Zaufanie do folderów jest zapisywane osobno dla każdego konta w `.claude.json`, więc folder zaufany jako `personal` jest dla `work` nowy. Uruchom raz `ccam work` w tym folderze i zaakceptuj, a potem uruchom tam Remote Control.

### "You must be logged in with a Claude account that has a subscription"

Remote Control wymaga logowania na konto z subskrypcją Claude. Sprawdź adres e-mail konta w `ccam list` i w razie potrzeby zaloguj się ponownie przez `ccam login work`.

### Sesja nie pojawia się na telefonie

Aplikacja Claude lub claude.ai/code jest zalogowana na inne konto Claude. Porównaj je z adresem e-mail, który `ccam list` pokazuje dla `work`.

### ccam ostrzega, że inny Claude Code współdzieli historię rozmów w tym folderze

W tym samym folderze działa inne konto z tej samej grupy historii. Sesje kontynuujące tę samą rozmowę z obu kont mogą wymieszać historię. Odpowiedz `y` tylko wtedy, gdy będą pracować nad różnymi rozmowami, albo uruchom Remote Control w innym folderze.

## Powiązane

- [Uruchamianie Claude Code](/features/launch)
- [Współdzielone serwery MCP](/features/mcp)
- [Współdzielenie między kontami](/features/share)
- [Zarządzanie kontami](/features/accounts)
