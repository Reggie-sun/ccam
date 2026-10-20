---
title: Polecenia skrótów
description: Daj kontu własne polecenie, takie jak ccwork, przez ccam shortcut, a linki w ~/.local/bin odtwórz poleceniem ccam shims install.
updated: 2026-10-20
order: 3
translation: machine
sourceHash: 939256e7a67e
---

Skrót to polecenie o nazwie wybranej przez ciebie, które uruchamia Claude Code jako jedno konto: `ccwork --resume` robi to samo co `ccam work --resume`. Jest to link do ccam w `~/.local/bin`, więc działa w bash, zsh, skryptach i wszędzie tam, gdzie programy są uruchamiane z `PATH`. Skróty są opcjonalne; utwórz je dla kont, które uruchamiasz często.

## Użycie

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing links of all accounts
ccwork -p "summarise" --output-format json
```

ccam potwierdza każdą zmianę:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

Konto ma najwyżej jeden skrót, zapisany jako `shortcut` w `~/.config/ccam/config.toml`. Kreator `ccam add` pyta `Create a shortcut command for <name>? [y/N]`, a potem `Shortcut name [cc<name>]:`. Poza kreatorem `ccam add` tworzy skrót tylko z `--shortcut`. `ccam list` pokazuje skrót każdego konta albo `-`.

Skrót przekazuje każdy argument do Claude Code, tak samo jak `ccam <name>` po nazwie. Dotyczy to także `--dry-run` i `--help`: aby podejrzeć uruchomienie, użyj `ccam --dry-run work`.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--remove` | Usuwa skrót konta i jego link | wyłączone |
| `-h`, `--help` | Pomoc dla polecenia | |
| `--dry-run` | `ccam shortcut` ją respektuje: sprawdza nazwę i wypisuje `Would set the shortcut of work to ccwork.` lub `Would remove ccwork, the shortcut of work.` bez zapisywania i linkowania | wyłączone |

`ccam shims install` nie przyjmuje argumentów i odrzuca `--dry-run` komunikatem `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``.

Zasady nazw, sprawdzane w tej kolejności:

- od 1 do 64 znaków spośród `a-z`, `0-9`, `.`, `-` i `_`, zaczynając od litery lub cyfry;
- nie `ccam`, `claude` ani polecenie ccam, takie jak `run` czy `usage`;
- nie skrót innego konta;
- nie plik w `~/.local/bin`, który nie jest już linkiem do ccam;
- nie nazwa innego pliku wykonywalnego w twoim `PATH`.

Podanie polecenia razem z `--remove` kończy się błędem `give a command or --remove, not both`. Ustawienie nazwy, którą konto już ma, wypisuje `work: ccwork (unchanged)` i odtwarza link, jeśli go brakuje.

## Jak to działa

Ustawienie nazwy uruchamia `checkShortcut`, które wywołuje `shim.Check` z folderem bin i rzeczywistą ścieżką pliku binarnego ccam:

<!-- motion:shortcuts-check -->

Nazwa przechodzi kolejno przez sprawdzenia. Obok każdego sprawdzenia jest nazwa, którą ono odrzuca: `Work` nie pasuje do wzorca, `run` jest poleceniem ccam, `ccside` należy do innego konta, `ccfile` to istniejący plik w `~/.local/bin`, a `git` przesłoniłby `/usr/bin/git`. Zapisywana i linkowana jest tylko nazwa, która przejdzie wszystkie sprawdzenia.

1. `shim.Check` dopasowuje nazwę do `^[a-z0-9][a-z0-9._-]{0,63}$`, potem do `config.Reserved`, potem do `cfg.FindShortcut`. Wywołuje `os.Lstat` na `~/.local/bin/<name>`, a dla każdego innego folderu z `PATH` wywołuje `os.Stat` i odrzuca plik wykonywalny, który nie jest linkiem do ccam.
2. ccam zapisuje `config.toml` przez `fsutil.AtomicWrite` (plik tymczasowy, `fsync`, zmiana nazwy).
3. `shim.Remove` usuwa link starego skrótu, ale tylko jeśli jest linkiem do ccam.
4. `shim.Install` tworzy brakujący link każdego skonfigurowanego skrótu przez `os.Symlink` do pliku binarnego ccam, odnalezionego przez `os.Executable` i `filepath.EvalSymlinks`. Nigdy niczego nie usuwa.

Link uznaje się za należący do ccam, gdy `os.Readlink` zwraca ścieżkę ccam albo gdy link i ścieżka ccam prowadzą do tego samego pliku. Każde usunięcie i każde sprawdzenie przed nadpisaniem przechodzi przez ten test, więc ccam nie rusza plików, których nie utworzył.

Gdy uruchamiasz skrót, wpisana nazwa to jedyne, czego ccam potrzebuje:

<!-- motion:shortcuts-link -->

Powłoka znajduje `ccwork` w `PATH`, jądro podąża za linkiem i uruchamia plik binarny ccam, a `argv[0]` nadal ma wartość `ccwork`. W `app.Main` każda nazwa inna niż `ccam` raz wczytuje `config.toml` i wywołuje `shim.Lookup`, które znajduje konto, którego `shortcut` jest równy `filepath.Base(argv[0])`. `launcher.LaunchWith` wykonuje następnie te same kroki co [`ccam <name>`](/features/launch) z już wczytaną konfiguracją i przekazuje `argv[1:]` do `claude` bez zmian.

Dlaczego dowiązanie symboliczne, a nie alias powłoki czy skrypt opakowujący: dowiązanie działa w każdej powłoce i w skryptach, nie dodaje uruchomienia kolejnej powłoki (około 1 do 2 ms), a uruchomienie nadal odczytuje jeden plik. `scripts/perf-check.sh` mierzy skrót za pomocą `hyperfine` względem tego samego limitu 10 ms (p95) co `ccam <name>`.

### W Windows

Skrót to `ccwork.exe` w `%USERPROFILE%\.local\bin`, a `shim.Name` przed wyszukaniem obcina `.exe` (bez względu na wielkość liter), więc `ccwork`, `ccwork.exe` i `CCWORK.EXE` uruchamiają `work`. ccam najpierw tworzy dowiązanie symboliczne, na które Windows pozwala w Trybie dewelopera albo administratorowi. Bez tego uprawnienia tworzy dowiązanie twarde: drugą nazwę tego samego pliku `ccam.exe`, która nie wymaga uprawnień. Dowiązanie twarde uznawane jest za należące do ccam, gdy jest tym samym plikiem co `ccam.exe`. Ponieważ `ccam update` wstawia nowy plik w miejsce `ccam.exe`, odnawia potem każdy skrót będący dowiązaniem twardym przez `shim.Relink` i wypisuje `✓ shortcut ccwork renewed`. Sprawdzenie nazwy szuka `<name>` z każdym rozszerzeniem z `PATHEXT` (`.exe`, `.cmd` i tak dalej) w pozostałych folderach `PATH`.

## Bezpieczeństwo i przypadki brzegowe

- **Link, do którego nie przyznaje się żadne konto.** Jeśli link do ccam w `~/.local/bin` ma nazwę, której nie używa żadne konto, na przykład po ręcznej edycji `shortcut` w `config.toml`, ccam nie przechodzi na domyślne konto. Zatrzymuje się z komunikatem `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) ``. Ustaw nazwę ponownie przez `ccam shortcut` albo usuń link.
- **Istniejące pliki.** ccam nigdy nie zastępuje pliku, którego nie utworzył. `ccam shims install` zatrzymuje się z komunikatem `<path> already exists and is not a ccam link, so it was left alone`.
- **Przeniesiony plik binarny ccam.** Stare linki wskazują wtedy ścieżkę, która już nie istnieje, i ccam nie uznaje ich za swoje. Usuń je, a potem uruchom `ccam shims install`.
- **Zmiana nazwy i usuwanie.** `ccam rename` zachowuje nazwę skrótu, a link dalej działa, bo wskazuje na ccam, a ccam szuka nazwy w `config.toml`. `ccam remove` usuwa link konta. `ccam shortcut <account> <new>` usuwa stary link, zanim utworzy nowy.
- **PATH.** Gdy ccam tworzy link, a `~/.local/bin` nie ma w twoim `PATH`, wypisuje `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:`, a po tym wiersz `export PATH=...` do dodania. W Windows wypisuje zamiast tego wiersz PowerShell, który dodaje folder do `PATH` twojego użytkownika.
- **Późniejsze sprawdzenie.** `ccam doctor` zgłasza `missing shortcuts in /home/me/.local/bin: ccwork` i ostrzega, gdy folderu nie ma w `PATH`. `ccam doctor --fix` tworzy brakujące linki.

## Powiązane

- [Uruchamianie Claude Code](/features/launch)
- [Zarządzanie kontami](/features/accounts)
- [Opis poleceń](/docs/commands#shortcuts)
- [Instalacja ccam](/docs/install)
