---
title: Uruchamianie Claude Code
description: Uruchom Claude Code na dowolnym koncie przez ccam <name|number>, ccam run lub panel; argumenty idą do claude bez zmian, --dry-run pokazuje polecenie.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: "4812b197f143"
---

`ccam work` uruchamia Claude Code jako konto `work`, a wszystko po nazwie trafia do `claude` bez zmian. Używaj go wszędzie tam, gdzie wpisałbyś `claude`: w terminalu, w skrypcie albo z wejściem przekazanym potokiem.

## Użycie

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the folder's account, else the default
ccam use work                    # make work the default account
ccam folder work ~/code/company  # use work in that folder when no name is given
ccam run work --continue         # same as ccam work --continue
ccam                             # the dashboard on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

Podaj nazwę konta albo jego numer z `ccam list`. ccam nie odczytuje argumentów po nazwie, więc `ccam work --help` pokazuje pomoc Claude Code. `--` po nazwie jest opcjonalne. Standardowe wejście, kod wyjścia i sygnały działają tak jak w `claude`.

Gdy pierwszy argument zaczyna się od `-` i nie jest żadną z własnych opcji ccam, ccam uruchamia konto bez nazwy, ze wszystkimi argumentami. Jest to konto ustawione dla bieżącego folderu przez [`ccam folder`](/features/folder), a jeśli go nie ma, konto domyślne: ustawione przez `ccam use`, a jeśli go nie ustawiłeś, konto z najniższym numerem.

### Panel

`ccam` bez argumentów otwiera [panel](/features/dashboard), gdy standardowe wejście i wyjście są terminalem, a w przeciwnym razie wypisuje `ccam list`. Przejdź do konta (strzałkami albo wpisując jego numer) i naciśnij `Enter`, żeby je uruchomić; `n` uruchamia konto z największym zapasem, a `o` proponuje `--continue`, `--resume` lub inne argumenty. Panel zamyka się przed startem Claude Code, więc terminal zostaje najpierw przywrócony.

Gdy w powłoce ustawiono `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` lub `ANTHROPIC_AUTH_TOKEN`, Claude Code użyłby tego klucza lub tokena zamiast własnego loginu konta, więc ccam mówi o tym na standardowym wyjściu błędów przed startem (`ccam: $ANTHROPIC_API_KEY is set, so Claude Code may use that key instead of work's login…`). Tylko odczytuje te zmienne; [konto z kluczem API](/features/api) zamiast tego je czyści.

Z `relay = true` w `config.toml` ccam uruchamia sesje interaktywne jako proces potomny i pozostaje obok nich, żeby `/ccam:switch` mogło samo przenieść rozmowę dalej ([szczegóły](/features/next#inside-claude-code)); w przeciwnym razie ccam zastępuje się procesem Claude Code i nic nie działa obok niego.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--dry-run` | Wypisuje środowisko i polecenie, które ccam by uruchomił, i kończy działanie | wyłączone |
| `--fallback` | Dla przebiegów `-p`/`--print`: gdy konto osiągnie limit użycia, uruchamia to samo polecenie ponownie na zalogowanym koncie z największym zapasem ([szczegóły](/features/next#run-again-at-the-limit)) | wyłączone |
| `-v`, `--version` | Wypisuje `ccam <version>` (tylko jako pierwszy argument) | |
| `-h`, `--help` | Pomoc ccam, jako `ccam --help` lub `ccam run --help` | |

Zasady dotyczące argumentów:

- ccam odczytuje własne opcje tylko przed nazwą konta. Po nazwie każdy argument należy do Claude Code, w tym `--dry-run` i `--help`.
- `--dry-run` działa jako `ccam --dry-run work`, `ccam run --dry-run work` i `ccam --dry-run` (konto bez nazwy). Jest respektowana, a nie odrzucana. `--dry-run` i `--fallback` mogą występować w dowolnej kolejności.
- `ccam run <name|number> [claude args]` to długa forma `ccam <name>`. `ccam run` bez nazwy używa konta folderu albo domyślnego.
- Nazwy kont nie mogą być nazwami poleceń ccam, takimi jak `add`, `list`, `run` czy `usage`, więc `ccam <word>` nigdy nie jest niejednoznaczne.

Próbne uruchomienie wypisuje jeden wiersz, który możesz wkleić do powłoki, a w Windows do PowerShell:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

```console
PS> ccam --dry-run work -p "fix the tests"
$env:CLAUDE_CONFIG_DIR='C:\Users\me\.local\share\ccam\profiles\work'; & C:\Users\me\.local\bin\claude.exe --mcp-config=C:\Users\me\.local\share\ccam\groups\default\mcp.json -p 'fix the tests'
```

## Jak to działa

Każde uruchomienie przechodzi te same kroki, od wiersza poleceń do `exec`:

<!-- motion:launch -->

Każde pole podświetla się, gdy wykonuje się jego krok, a wiersze pod spodem to środowisko i argumenty, z którymi startuje Claude Code.

1. `cmd/ccam/main.go` wywołuje `app.Main`. Jeśli ccam uruchomiono pod inną nazwą, jest to [skrót](/features/shortcuts) i trafia od razu do programu uruchamiającego.
2. `app.parse` analizuje `os.Args`, zanim w ogóle zaangażuje się cobra. Tylko samo `ccam` (panel), nazwa polecenia ccam (`config.Reserved`), `-h`/`--help` i `-v`/`--version` trafiają do `cli.Execute`. Wszystko inne to uruchomienie.
3. `launcher.Launch` odczytuje `~/.config/ccam/config.toml`, jedyny plik, który parsuje.
4. `launcher.Prepare` znajduje konto po nazwie lub numerze (bez nazwy `HereProfile` bierze konto folderu, porównując katalog roboczy z wpisami `[[folder]]` w `config.toml` jako ciągi znaków), odrzuca jedno początkowe `--` i odnajduje plik wykonywalny `claude` przez `ResolveClaude`. Wynik wyszukiwania w PATH jest zapisywany w pamięci podręcznej `~/.local/state/ccam/claude-path` razem z czasem modyfikacji pliku, więc ccam nie przeszukuje `PATH` przy każdym uruchomieniu. Zmieniony lub brakujący plik wywołuje nowe wyszukiwanie. Ścieżka `claude_bin` w `config.toml` jest używana bez zmian. ccam odrzuca `claude`, który jest samym ccam (porównanie po i-węźle), bo powstałaby pętla.
5. `checkLinks` porównuje folder konta `~/.local/share/ccam/profiles/<name>/` z jego ustawieniami współdzielenia przez `share.Check`, które wywołuje tylko `lstat` i `readlink`, jedno lub dwa wywołania systemowe na element. Brakujący link jest tworzony od razu. Wszystko, co mogłoby spowodować utratę danych, zatrzymuje uruchomienie.
6. `applyMCP` dodaje współdzielone serwery MCP (zobacz niżej).
7. `session.Register` rejestruje sesję. ccam zapisuje znacznik nazwany od swojego PID, zawierający czas startu procesu, w `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` (albo w `~/.local/state/ccam/locks/`, gdy nie ma `XDG_RUNTIME_DIR`). Dla konta, które współdzieli historię rozmów, `historyCheck` najpierw szuka działającej sesji w `history-<group>-<hash>`, gdzie skrót jest liczony z bieżącego folderu, i tam też się rejestruje. `exec` zachowuje PID, więc znacznik wskazuje potem sam proces Claude Code. To znacznik, a nie trzymana blokada, bo procesy uruchamiane przez Claude Code odziedziczyłyby blokadę i trzymały ją po zakończeniu sesji.
8. `Plan.Exec` ustawia tytuł terminala na `claude · <name>` (gdy standardowe wyjście błędów jest terminalem), zapisuje czas uruchomienia w `~/.local/state/ccam/last-used/<name>` bez `fsync` i wywołuje `syscall.Exec`.

Claude Code działa z `CLAUDE_CONFIG_DIR` ustawionym na folder konta i `CCAM_PROFILE` ustawionym na nazwę konta. ccam usuwa wszelkie odziedziczone wartości tych zmiennych. Dla folderu `~/.claude` zarejestrowanego przez `ccam import` na macOS ustawia też `CLAUDE_SECURESTORAGE_CONFIG_DIR` na pustą wartość, żeby konto i zwykłe `claude` używały tego samego loginu z pęku kluczy.

### Dlaczego exec

`syscall.Exec` zastępuje proces ccam procesem `claude`, więc ccam nie zostaje w pamięci, naciśnięcia klawiszy, sygnały i kod wyjścia należą do Claude Code, a argumenty pozostają tablicą, nigdy ciągiem dla powłoki.

Windows nie ma `exec`. Tam `Plan.Exec` uruchamia `claude.exe` jako proces potomny z tym samym środowiskiem, argumentami i konsolą, czeka na niego i kończy działanie z jego kodem wyjścia. Ctrl+C trafia do Claude Code bezpośrednio z konsoli, więc ccam go ignoruje. Znacznik sesji zawiera PID ccam, który żyje dokładnie tak długo jak Claude Code, i jest usuwany, gdy Claude Code się zakończy. Podczas czekania ccam nie robi nic innego. Prace po sesji, takie jak wykrycie serwerów dodanych przez `claude mcp add`, odbywają się przy następnym uruchomieniu: ccam porównuje czas modyfikacji `.claude.json` z `~/.local/state/ccam/mcp/<name>/seen` i parsuje plik tylko wtedy, gdy się zmienił.

### Jak kierowany jest wiersz poleceń

Trzy wiersze poleceń, trzy drogi przez `internal/app`:

<!-- motion:launch-route -->

Skrót jest rozpoznawany po `argv[0]`, `ccam <name>` po pierwszym argumencie, a obie drogi kończą się na `syscall.Exec`. Tylko polecenie ccam, takie jak `ccam usage`, trafia do cobra.

### Współdzielone serwery MCP

Przy domyślnej strategii `flag` ccam przekazuje serwery MCP grupy do Claude Code jako argument, umieszczony przed twoimi:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` zwraca argument tylko wtedy, gdy `mcp.json` grupy istnieje, a `Plan.InsertArgs` wstawia go zaraz po `claude`. Twoje własne `--mcp-config` lub `--strict-mcp-config` nadal obowiązują dodatkowo. To jeden argument, `--mcp-config=<path>`, bo `--mcp-config` w Claude Code przyjmuje kilka wartości: zapisany jako dwa argumenty odczytałby też następujący po nim prompt jako plik konfiguracji.

Z `mcp_strategy = "sync"` w `config.toml` ccam zamiast tego zapisuje serwery w `.claude.json` konta i pomija ten krok, jeśli konto już działa.

### Remote Control

`ccam work rc` (lub `ccam work remote-control`) uruchamia Remote Control z Claude Code jako konto `work`. `claude remote-control` odrzuca `--mcp-config` przed swoim poleceniem i nie przyjmuje go po nim, a jego sesje odczytują serwery MCP tylko z `.claude.json` konta. Dlatego strategia flagi działa tu inaczej:

<!-- motion:launch-rc -->

1. `remoteControl` traktuje pierwszy argument po koncie, który nie zaczyna się od `-`, jako polecenie Claude Code. Gdy jest to `remote-control` lub `rc`, `mcp.FlagArgs` niczego nie dodaje.
2. `mcp.Synced` porównuje czas modyfikacji `mcp.json` grupy z `sync-stamp` konta. Jeśli `.claude.json` zawiera już bieżące serwery grupy albo grupa nie ma żadnych, Remote Control startuje bez pytania.
3. W przeciwnym razie ccam wymienia serwery, których zabrakłoby w sesjach, i pyta:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. Po odpowiedzi `y` funkcja `mcp.Sync` kopiuje je tak jak strategia `sync`: pod blokadą Claude Code, po wykonaniu kopii zapasowej i z zachowaniem własnych serwerów konta. Po odpowiedzi `N` albo bez terminala Remote Control startuje bez nich, a ccam wskazuje `mcp_strategy = "sync"`, które kopiuje je przed każdym uruchomieniem.

ccam pyta przy każdym uruchomieniu Remote Control, dopóki serwery nie zostaną skopiowane, a po skopiowaniu dopiero wtedy, gdy `mcp.json` grupy znów się zmieni. Kopia w `.claude.json` wygrywa z `--mcp-config` przy późniejszych zwykłych uruchomieniach, dlatego `ccam mcp add`, `remove` i `import` aktualizują kopie przechowywane przez inne konta ([Współdzielone serwery MCP](/features/mcp)).

### Budżet 10 ms

ccam dodaje najwyżej 10 ms (p95) przed startem Claude Code. Ścieżka uruchamiania nigdy nie buduje drzewa poleceń cobra, nie uruchamia Bubble Tea ani nie otwiera połączenia sieciowego, a `TestPackageInitIsCheap` kończy się niepowodzeniem, jeśli funkcje `init()` pakietów zajmują łącznie więcej niż 5 ms. `scripts/perf-check.sh` mierzy `ccam --dry-run work`, `ccam work` i skrót za pomocą `hyperfine` i zgłasza błąd powyżej 10 ms, a `BenchmarkLaunchPath` w `internal/launcher` wykonuje kroki w pętli. Aby zobaczyć czas każdego kroku, ustaw `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## Bezpieczeństwo i przypadki brzegowe

- Nieznana nazwa lub numer kończy się komunikatem `no account "<name>"` i wskazaniem `ccam list`. Gdy nie ma żadnych kont, ccam prosi o dodanie konta przez `ccam add <name>`.
- Jeśli współdzielona ścieżka stała się zwykłym plikiem lub wskazuje złe miejsce, ccam zatrzymuje się z komunikatem `stopped because of share link problems:`, wymienia problemy i prosi o uruchomienie `ccam doctor --fix`. Na ścieżce uruchamiania nie przenosi danych.
- Jeśli w tym samym folderze działa inny Claude Code, który współdzieli historię rozmów z tym kontem (ta sama grupa), ccam ostrzega, że historia może się wymieszać, i pyta `Run anyway? [y/N]`. Bez terminala pokazuje ostrzeżenie i kontynuuje.
- `--dry-run` pomija synchronizację MCP, informację o nowych serwerach, sprawdzenie historii, znacznik sesji, tytuł i czas ostatniego użycia, i nie wykonuje exec. Nadal sprawdza linki współdzielenia i tworzy brakujący link tak jak zwykłe uruchomienie.
- `ccam work rc`, gdy `work` już działa, nie rusza `.claude.json`, informuje o tym i uruchamia Remote Control bez serwerów grupy; uruchom polecenie ponownie po zakończeniu sesji, aby je skopiować. Z `--dry-run` wypisuje komunikat, ale ani nie pyta, ani nie kopiuje.
- Z `mcp_strategy = "sync"` nieudana synchronizacja wypisuje `ccam: MCP sync failed (launching anyway): …`, a Claude Code startuje z serwerami, które już miał.

## Powiązane

- [Polecenia skrótów](/features/shortcuts)
- [Remote Control z ccam](/guides/remote-control)
- [Zarządzanie kontami](/features/accounts)
- [Opis poleceń](/docs/commands)
- [Współdzielenie serwerów MCP, ustawień i historii](/docs/sharing)
- [Jak działa CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
