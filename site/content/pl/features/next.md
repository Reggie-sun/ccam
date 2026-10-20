---
title: Przełączanie na następne konto
description: ccam next otwiera Claude Code jako konto po ostatnio używanym; --auto wybiera konto z największym zapasem, --continue wznawia, --fallback ponawia zadania -p.
updated: 2026-10-20
order: 8
translation: machine
sourceHash: 784ae4d1f6d9
---

`ccam next` otwiera Claude Code jako inne konto bez podawania jego nazwy. Używaj go, gdy konto osiągnie limit użycia: samo `ccam next` wybiera konto następne po ostatnio używanym, `--auto` wybiera konto z największym zapasem w 5-godzinnym oknie, a `--continue` kontynuuje ostatnią rozmowę na nowym koncie. W skryptach, które uruchamiają `claude -p`, `ccam --fallback` przełącza konto samo (zobacz [Ponowienie po osiągnięciu limitu](#run-again-at-the-limit)).

## Użycie

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next
→ personal → work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude
```

ccam wypisuje przełączenie jako `→ <from> → <to>` na standardowe wyjście błędów, a potem uruchamia Claude Code jako nowe konto tak samo jak `ccam <name>`. `<from>` to konto, którego używałeś ostatnio. Argumenty po `--` trafiają do `claude` bez zmian.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--auto` | Wybiera zalogowane konto według użycia (zobacz niżej), pomijając konta na poziomie 100% | wyłączone |
| `--max-usage <1-100>` | Z `--auto`: preferuje konta poniżej tego użycia | `auto_max_usage`, w przeciwnym razie 100 |
| `--strategy headroom\|expiring` | Z `--auto`: najniższe użycie 5-godzinne albo najbliższy reset okna 5-godzinnego | `auto_strategy`, w przeciwnym razie `headroom` |
| `--wait[=<duration>]` | Z `--auto`: gdy żadne konto nie ma zapasu, czeka na reset pierwszego, najwyżej tyle | wyłączone; `5h`, gdy podano samą opcję |
| `-c`, `--continue` | Bierze pod uwagę tylko konta, które współdzielą historię rozmów z ostatnio używanym kontem, i uruchamia `claude --continue` | wyłączone |
| `-- <claude args>` | Argumenty przekazywane do `claude` po argumentach ccam | brak |
| `--dry-run` | Respektowana: wypisuje `→ <from> → <to>` i wiersz polecenia, który zostałby uruchomiony, bez startu Claude Code i bez zapisywania uruchomienia | wyłączone |

`--auto` i `--continue` działają razem. Z `--continue` opcja `--continue` jest umieszczana przed twoimi argumentami.

## Jak to działa

Każde uruchomienie przez ccam zapisuje czas w `~/.local/state/ccam/last-used/<name>`. `ccam next` odczytuje te pliki, sprawdza, które konta są zalogowane, i wybiera.

1. `cli.candidates` tworzy jeden `next.Candidate` na konto z `config.toml`: konto, czas jego ostatniego uruchomienia (`launcher.LastUsed`) i informację, czy ma login (`creds.Has`: czytelny `.credentials.json` albo, na macOS, wpis w pęku kluczy).
2. Bez `--auto` `next.Rotate` sortuje konta według numeru slotu i znajduje to uruchomione najpóźniej. To konto, z którego się przełączasz. Jeśli żadne konto nie było jeszcze uruchamiane, jest to najniższy slot.
3. Od niego idzie dalej w kolejności slotów, wracając od ostatniego slotu do pierwszego, i zwraca pierwsze zalogowane konto.
4. `launcher.LaunchWith` przygotowuje i uruchamia wybrane konto dokładnie tak jak `ccam <name>`: sprawdza linki konta, stosuje współdzielenie MCP, rejestruje sesję i zastępuje ccam przez `claude` za pomocą `syscall.Exec`.

Przejście zaczyna się od ostatnio używanego konta i zatrzymuje na pierwszym zalogowanym koncie po nim:

<!-- motion:next-rotate -->

`work` uruchomiono najpóźniej, więc przejście zaczyna się po nim. `side` nie ma loginu i zostaje pominięte, a wybrane zostaje `lab`. Gdyby `lab` też było wylogowane, przejście wróciłoby na początek do `personal` (`% len` na diagramie). Rotacja nigdy nie zwraca konta, od którego zaczęła: gdy nie ma innego zalogowanego konta, kończy się błędem.

### --auto

`--auto` decyduje według użycia zamiast kolejności slotów. `cli.pickAuto` zostawia zalogowane konta, pyta o użycie wszystkich naraz tym samym `Fetcher.FetchAll`, którego używa `ccam usage`, i pozwala wybrać `usage.Policy.Choose()`:

<!-- motion:next-auto -->

Konto jest oceniane według okna 5-godzinnego i tygodniowego, a gdy wśród argumentów dla Claude Code jest `--model`, także według tygodniowego okna tego modelu (`opus`, `sonnet`, `haiku` lub `fable` w nazwie). `Policy.Choose()` dzieli konta na trzy grupy:

1. **Z zapasem**: każde okno poniżej `auto_max_usage` (`config.toml`, domyślnie 100). Spośród nich wybiera strategia: `headroom` (domyślna) bierze najniższe użycie 5-godzinne, `expiring` bierze okno 5-godzinne, które zresetuje się najwcześniej, żeby pozostały limit został zużyty, zanim przepadnie.
2. **Powyżej progu**: poniżej 100%, ale na poziomie `auto_max_usage` lub wyżej. Używane tylko wtedy, gdy żadne konto nie ma zapasu, zaczynając od najniższego użycia 5-godzinnego.
3. **Pełne**: okno na poziomie 100% albo nieznane użycie 5-godzinne. Nigdy niewybierane.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

Gdy nie da się nic wybrać, ccam kończy działanie z kodem 75 (`EX_TEMPFAIL`), więc skrypt może odróżnić „spróbuj później” od błędu. Z `--wait` zamiast tego czeka, aż zresetuje się okno blokujące pierwsze pełne konto (30 sekund po czasie resetu, gdy endpoint pokazuje już nowe okno), wypisuje `ccam: no account has room; waiting until 15:10` i wybiera ponownie; gdy żadne konto nie podaje czasu resetu, sprawdza ponownie po 10 minutach. Jeśli nic się nie zwolni przed końcem oczekiwania, kończy działanie z kodem 75. Zapytanie działa według zasad `ccam usage`: liczby młodsze niż 60 sekund są używane ponownie, każde konto ma 3 sekundy, a konto, które nie odpowiada lub którego login wygasł, jest oceniane na podstawie zapisanych liczb. Zapisane okno, którego czas resetu minął, liczy się jako 0%, więc konto, którego okno w międzyczasie się zresetowało, nie zostaje pominięte. `--auto` nigdy nie odnawia tokenów. W odróżnieniu od zwykłej rotacji `--auto` może wybrać ostatnio używane konto, jeśli nadal ma największy zapas.

### --continue

`--continue` w Claude Code otwiera najnowszą rozmowę w bieżącym folderze, co działa tylko wtedy, gdy nowe konto widzi tę rozmowę. Dlatego z `--continue` ccam bierze pod uwagę tylko konta z tej samej grupy `history` co ostatnio używane konto, a potem umieszcza `--continue` przed twoimi argumentami. Konta z historią `local` nie współdzielą jej z nikim, więc `--continue` z takiego konta zawsze się nie powiedzie. Najpierw włącz współdzielenie przez `ccam share <name> --history on` (zobacz [współdzielenie](/docs/sharing)).

## Ponowienie po osiągnięciu limitu

`ccam next` służy do sesji interaktywnej, która się zatrzymała. Przebieg nieinteraktywny, `claude -p` w skrypcie lub zaplanowanym zadaniu, nie może na ciebie czekać, więc `--fallback` przed nazwą konta pozwala ccam przełączyć konto za niego:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` działa tylko razem z `-p` lub `--print`; w przeciwnym razie ccam zatrzymuje się i odsyła do `ccam next --auto --continue`.
2. Standardowe wejście jest odczytywane raz i przekazywane każdej próbie.
3. ccam uruchamia `claude` jako proces potomny zamiast przez `exec`, wstrzymuje jego standardowe wyjście, a standardowe wyjście błędów przekazuje na bieżąco.
4. Gdy `claude` zakończy się kodem 0 albo nie powiedzie się z innego powodu, wstrzymane wyjście i kod wyjścia są przekazywane bez zmian.
5. Błąd jest uznawany za limit użycia, gdy tak mówi komunikat Claude Code (`You've hit your limit`, `usage limit reached`) albo gdy sprawdzenie użycia konta pokazuje okno 5-godzinne lub tygodniowe na poziomie 100%.
6. Wtedy wyjście tej próby jest odrzucane, a to samo polecenie uruchamia się na innym zalogowanym koncie: najpierw na tym z najniższym użyciem 5-godzinnym, potem na kontach, których użycia nie da się sprawdzić, według numerów. Konta na poziomie 100% są pomijane, a każde konto jest próbowane najwyżej raz.
7. Gdy wśród argumentów jest `-c`, `--continue`, `-r` lub `--resume`, próbowane są tylko konta z tej samej grupy historii, żeby dało się znaleźć rozmowę.

Następne konto jest wybierane według tych samych zasad co w `--auto` (`auto_max_usage`, `auto_strategy`, okno `--model`). Gdy nie zostanie żadne konto, ccam wypisuje wyjście ostatniej próby, informuje, że żadne zalogowane konto nie ma już zapasu, i kończy działanie z kodem 75. Z `--wait` przed nazwą konta (`ccam --fallback --wait=2h work -p …`) zamiast tego czeka na reset pierwszego konta, a potem może ponownie spróbować każdego konta; zawsze najpierw czeka co najmniej raz, żeby konto, które właśnie osiągnęło limit, nie zostało od razu uruchomione ponownie. `ccam --dry-run --fallback work -p …` wypisuje polecenie pierwszej próby.

Sesje interaktywne nie przełączają się same. Zamiast tego, po `ccam statusline install`, pasek stanu dodaje `→ ccam next --auto --continue`, gdy użycie 5-godzinne lub tygodniowe osiągnie 90% ([Limity użycia](/features/usage)).

## Bezpieczeństwo i przypadki brzegowe

- Konta bez loginu nigdy nie są wybierane. Gdy nie ma innego zalogowanego konta, samo `ccam next` kończy się komunikatem `no logged-in account other than <name>`.
- Z `--continue`, gdy żadne konto nie współdzieli historii: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- Z `--auto`, gdy nie ma zalogowanego konta: `no logged-in account to pick`. Gdy każde konto jest pełne lub ma nieznane użycie: `` no account with known usage and headroom (see `ccam usage`) ``. W obu przypadkach kod wyjścia to 75.
- `--max-usage`, `--strategy` i `--wait` działają tylko z `--auto`, a wartość spoza zakresu 1-100 lub inna strategia jest odrzucana.
- `--auto` wysyła zapytania sieciowe nawet z `--dry-run`. Tylko odczytuje i nigdy nie zapisuje plików z tokenami.
- Wybrane konto startuje tą samą ścieżką co `ccam <name>`, więc obowiązują te same sprawdzenia, w tym ostrzeżenie, gdy w tym folderze działa już inny Claude Code współdzielący tę samą historię.
- `last-used` jest zapisywany tylko wtedy, gdy Claude Code rzeczywiście wystartuje. `--dry-run` go nie zmienia, więc próbne uruchomienie nie przesuwa rotacji.
- `--fallback` wstrzymuje standardowe wyjście do końca próby, więc wyjście strumieniowe, takie jak `--output-format stream-json`, przychodzi w całości na końcu. Pliki zmienione przez pierwszą próbę przed osiągnięciem limitu pozostają zmienione; następna próba zaczyna prompt od nowa.

## Powiązane

- [Limity użycia](/features/usage)
- [Limity użycia i przełączanie kont](/docs/usage-limits)
- [Współdzielenie serwerów MCP, ustawień i historii](/docs/sharing)
- [Praca z limitem użycia](/guides/usage-limit-workflow)
