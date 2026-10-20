---
title: Opis poleceń
metaTitle: "Opis poleceń ccam"
description: Wszystkie polecenia i opcje ccam, od ccam add i ccam <name> po ccam usage, ccam next i ccam doctor.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: da2594c672bc
---

Wpisz `--dry-run` przed poleceniem, żeby zobaczyć, co by zrobiło, bez wykonywania go, na przykład `ccam --dry-run work`. Działa to przy uruchamianiu Claude Code (`ccam <name>`, `ccam run`, `ccam next`) oraz z `share`, `migrate`, `shortcut`, `update` i `doctor --fix`. Inne polecenia, które coś zmieniają, odrzucają tę opcję zamiast się wykonać.

<!-- commands -->

## Uruchamianie Claude Code

`ccam <name|number> [arguments]` uruchamia Claude Code jako to konto. ccam nie odczytuje argumentów po nazwie; przekazuje je do `claude` w tej samej kolejności, a standardowe wejście, kod wyjścia i sygnały działają dokładnie tak jak w samym `claude`. `--` po nazwie jest opcjonalne.

`ccam` z argumentami, ale bez nazwy konta, na przykład `ccam -p "summarise"`, używa konta ustawionego dla bieżącego folderu, a jeśli go nie ma, konta domyślnego. Ustawisz je poleceniami `ccam folder <name> [folder]` i `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` przed nazwą konta jest przeznaczone dla przebiegów `-p`/`--print`. Gdy konto osiągnie limit użycia, ccam uruchamia to samo polecenie ponownie na zalogowanym koncie z największym zapasem i wypisuje tylko wyjście przebiegu, który się zakończył. Zobacz [Konta według folderu](/features/folder) i [Ponowienie po osiągnięciu limitu](/features/next#run-again-at-the-limit).

Skrót taki jak `ccwork` uruchamia swoje konto w ten sam sposób: `ccwork --resume` robi to samo co `ccam work --resume`. Zobacz [skróty](/docs/commands#shortcuts).

## Dodawanie i usuwanie kont

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` bez podanego folderu szuka `~/.claude` i `~/.claude-*` i pyta o każdy z nich. Zaimportowane foldery są rejestrowane tam, gdzie są; ccam nigdy nie przenosi ani nie kopiuje loginu.

`ccam add --api` dodaje konto, które zamiast loginu Claude używa klucza API; ccam odczytuje klucz przy każdym uruchomieniu i nigdy go nie zapisuje. Zobacz [Konta z kluczem API](/features/api).

`ccam rename`, `ccam remove` i `ccam share` odmawiają zmiany konta, gdy Claude Code działa jako to konto. Najpierw zakończ sesję.

`ccam list --json` wypisuje każde konto z jego numerem (slotem), skrótem, adresem e-mail i folderem. Jego obiekt `share` przypisuje każdy element (`mcp`, `settings`, `memory`, `extensions`, `history`) do nazwy grupy albo `local`.

## Skróty

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Skrót to polecenie o nazwie wybranej przez ciebie, które uruchamia Claude Code jako jedno konto. Podobnie jak `ccam work` przekazuje każdy argument do Claude Code bez zmian, więc `ccwork --resume` i `ccwork -p "hi"` działają tak, jak się spodziewasz. Jest to link do ccam w `~/.local/bin`.

Skróty są opcjonalne. Kreator `ccam add` pyta, czy utworzyć skrót, `ccam add <name> --shortcut <command>` tworzy go bez pytania, a `ccam list` pokazuje skrót każdego konta albo `-`, gdy konto go nie ma.

Nazwy mogą zawierać znaki a-z, 0-9, `.`, `-` i `_`. Nazwa nie może być poleceniem ccam, skrótem innego konta ani nazwą innego programu w twoim `PATH`.

`ccam rename` zachowuje nazwę skrótu, a `ccam remove` usuwa skrót razem z kontem. `ccam shims install` odtwarza linki wszystkich kont, na przykład po przeniesieniu ccam.

## Współdzielenie

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Każdy element przyjmuje `on` (grupa `default`), `off` lub nazwę grupy. Zobacz [współdzielenie](/docs/sharing).

## Kopiowanie istniejącej konfiguracji

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` kopiuje `settings.json`, `CLAUDE.md`, rozszerzenia i serwery MCP użytkownika z folderu Claude Code do grupy, żeby dostało je każde konto współdzielące te elementy. Folder jest tylko odczytywany. Opcje `--settings`, `--memory`, `--extensions`, `--mcp` i `--history` kopiują tylko wskazane elementy. Bez nich w terminalu pojawia się lista wyboru, na której zaznaczone jest wszystko oprócz historii rozmów; poza terminalem kopiowane jest wszystko oprócz historii rozmów. `-g` lub `--group` wybiera grupę, a `-y` lub `--yes` kopiuje bez pytania i zachowuje kopię grupy, gdy plik się różni. Kopiuje też zapisane logowania serwerów do własnego pliku logowania każdego konta. Zobacz [kopiowanie istniejącej konfiguracji](/docs/sharing#copy-an-existing-setup).

## Serwery MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` przenosi do grupy konta serwery, które ma tylko to konto. Serwer, którego nazwy grupa używa już dla innej definicji, jest pomijany; dodaj `--force`, aby zastąpić kopię w grupie, a ccam najpierw zrobi jej kopię zapasową.

## Użycie i przełączanie

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
```

`next --auto`, `--fallback` i `ccam pick` wybierają w ten sam sposób: najpierw konta poniżej `auto_max_usage` (`config.toml`, domyślnie 100), według `auto_strategy` (`headroom` lub `expiring`), z uwzględnieniem tygodniowego okna `--model`, nigdy konto na poziomie 100%. Gdy nie da się wybrać żadnego, kończą działanie z kodem 75, a z `--wait` czekają na reset pierwszego konta. Zobacz [Przełączanie na następne konto](/features/next) i [Wybór konta dla skryptów](/features/pick).

## Konserwacja

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` sprawdza sumę kontrolną SHA-256 wydania, zanim podmieni ccam, a `--version v0.2.0` instaluje konkretne wydanie. ccam zainstalowany z pakietu `.deb` aktualizujesz przez apt. Szczegóły znajdziesz w części [aktualizacja](/docs/install#update).

## Pliki

| Ścieżka | Zawartość |
| --- | --- |
| `~/.config/ccam/config.toml` | Konta i to, co współdzielą |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` każdego konta |
| `~/.local/share/ccam/groups/<group>/` | Współdzielone ustawienia, serwery MCP, rozszerzenia i historia |
| `~/.local/state/ccam/` | Pamięć podręczna użycia i czasy ostatniego użycia |
| `~/.local/bin/<shortcut>` | Skróty takie jak `ccwork`, czyli linki do ccam |

Zmienne `XDG_CONFIG_HOME`, `XDG_DATA_HOME` i `XDG_STATE_HOME` zmieniają położenie tych folderów. W Windows `~` to `%USERPROFILE%`, a skróty to pliki `<shortcut>.exe`.
