---
title: Przełączanie na następne konto
description: ccam next otwiera Claude Code jako inne konto; --auto wybiera konto z największym zapasem, --continue przenosi rozmowę, --fallback ponawia zadania -p.
updated: 2026-10-21
order: 8
translation: machine
sourceHash: "afee5cce02cd"
---

`ccam next` otwiera Claude Code jako inne konto bez podawania jego nazwy. Używaj go, gdy konto osiągnie limit użycia: samo `ccam next` wybiera konto następne po ostatnio używanym, `--auto` wybiera konto z największym zapasem, a `--continue` przenosi ostatnią rozmowę z tego folderu na nowe konto, niezależnie od tego, czy oba konta współdzielą historię. W Claude Code [plugin ccam](#inside-claude-code) informuje, gdy limit zatrzyma sesję, a `/ccam:switch` przekazuje rozmowę dalej. W skryptach, które uruchamiają `claude -p`, `ccam --fallback` przełącza konto samo (zobacz [Ponowienie po osiągnięciu limitu](#run-again-at-the-limit)).

## Użycie

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the weekly quota that would go unused first
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam wypisuje przełączenie jako `→ <from> → <to>` na standardowe wyjście błędów, przy wyborze `--auto` z powodem, a potem uruchamia Claude Code jako nowe konto tak samo jak `ccam <name>`. `<from>` to konto, którego używałeś ostatnio, a z `--continue` konto, do którego należy rozmowa. Argumenty po `--` trafiają do `claude` bez zmian.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--auto` | Wybiera zalogowane konto według użycia (zobacz niżej), pomijając konta na poziomie 100% | wyłączone |
| `--max-usage <1-100>` | Z `--auto`: preferuje konta poniżej tego użycia | `auto_max_usage`, w przeciwnym razie 100 |
| `--strategy headroom\|expiring\|balance` | Z `--auto`: najniższe użycie 5-godzinne, najbliższy reset okna 5-godzinnego albo tygodniowy limit, który najwcześniej przepadłby niewykorzystany | `auto_strategy`, w przeciwnym razie `headroom` |
| `--cross-org` | Z `--auto`: wybiera też konta innej organizacji (zobacz niżej) | `cross_org`, w przeciwnym razie wyłączone |
| `--wait[=<duration>]` | Z `--auto`: gdy żadne konto nie ma zapasu, czeka na reset pierwszego, najwyżej tyle | wyłączone; `5h`, gdy podano samą opcję |
| `-c`, `--continue` | Przenosi ostatnią rozmowę z tego folderu na nowe konto (zobacz niżej) | wyłączone |
| `-- <claude args>` | Argumenty przekazywane do `claude` po argumentach ccam | brak |
| `--dry-run` | Respektowana: wypisuje `→ <from> → <to>` i wiersz polecenia, który zostałby uruchomiony, bez startu Claude Code i bez zapisywania uruchomienia | wyłączone |

`--auto` i `--continue` działają razem. Z `--continue` argumenty `--resume` są umieszczane przed twoimi.

Własne opcje Claude Code, takie jak `--dangerously-skip-permissions` czy `--model`, podaje się po `--`. W przeciwieństwie do `ccam <name>`, które przekazuje dalej wszystko po nazwie, `ccam next` samo odczytuje swoje argumenty, więc bez `--` zatrzymuje się z `ccam: unknown flag: --dangerously-skip-permissions`.

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

Konto jest oceniane według okna 5-godzinnego i tygodniowego, a gdy wśród argumentów dla Claude Code (po `--`) jest `--model`, także według tygodniowego okna tego modelu (`opus`, `sonnet`, `haiku` lub `fable` w nazwie). `Policy.Choose()` dzieli konta na trzy grupy:

1. **Z zapasem**: każde okno poniżej `auto_max_usage` (`config.toml`, domyślnie 100). Spośród nich wybiera strategia: `headroom` (domyślna) bierze najniższe użycie 5-godzinne; `expiring` bierze okno 5-godzinne, które zresetuje się najwcześniej, żeby pozostały limit został zużyty, zanim przepadnie; `balance` bierze konto z największym tygodniowym limitem pozostałym na każdy dzień do resetu tygodnia, licząc przeciw niemu każdą uruchomioną sesję ccam, żeby tygodniowy limit nie przepadał niewykorzystany, a równoległe sesje się rozkładały.
2. **Powyżej progu**: poniżej 100%, ale na poziomie `auto_max_usage` lub wyżej. Używane tylko wtedy, gdy żadne konto nie ma zapasu, zaczynając od najniższego użycia 5-godzinnego.
3. **Pełne**: okno na poziomie 100% albo nieznane użycie 5-godzinne. Nigdy niewybierane.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Organizacje.** Automatyczne wybory pozostają przy kontach podobnych do tego, które opuszczasz: z tej samej organizacji zespołowej lub firmowej albo razem wśród planów prywatnych (Pro, Max), więc praca z konta firmowego nie jest kontynuowana na prywatnym ani odwrotnie. ccam zna organizację każdego konta z jego loginu (`ccam login` ją zapisuje; starsze loginy są odczytywane z `.claude.json`); konto, którego organizacji nie da się ustalić, pasuje wszędzie. Dwa foldery zalogowane na to samo konto Claude dzielą jego limity, więc jeden nigdy nie przejmuje pracy po drugim. `--cross-org` albo `cross_org = true` w `config.toml` znosi zasadę organizacji.

Gdy nie da się nic wybrać, ccam kończy działanie z kodem 75 (`EX_TEMPFAIL`), więc skrypt może odróżnić „spróbuj później” od błędu. Z `--wait` zamiast tego czeka, aż zresetuje się okno blokujące pierwsze pełne konto (30 sekund po czasie resetu, gdy endpoint pokazuje już nowe okno), wypisuje `ccam: no account has room; waiting until 15:10` i wybiera ponownie; gdy żadne konto nie podaje czasu resetu, sprawdza ponownie po 10 minutach. Jeśli nic się nie zwolni przed końcem oczekiwania, kończy działanie z kodem 75. Zapytanie działa według zasad `ccam usage`: liczby młodsze niż 60 sekund są używane ponownie, każde konto ma 3 sekundy, a konto, które nie odpowiada lub którego login wygasł, jest oceniane na podstawie zapisanych liczb. Zapisane okno, którego czas resetu minął, liczy się jako 0%, więc konto, którego okno w międzyczasie się zresetowało, nie zostaje pominięte. `--auto` nigdy nie odnawia tokenów. W odróżnieniu od zwykłej rotacji `--auto` może wybrać ostatnio używane konto, jeśli nadal ma największy zapas.

### --continue

1. ccam szuka ostatniej rozmowy z tego folderu na każdym koncie: najnowszego `<session>.jsonl` w `projects/<folder>/` (nazwa folderu to folder roboczy, w którym każdy znak inny niż litera lub cyfra zamieniono na `-`, tak jak nazywa go Claude Code). Jeśli sesja zostawiła notatkę (limit użycia lub `/ccam:switch`, zobacz niżej) i od tego czasu nie zapisano nic nowszego, używana jest dokładnie ta sesja.
2. Konto, do którego ona należy, jest kontem, z którego się przełączasz, a nowe konto jest wybierane jak powyżej, w obrębie tej samej organizacji.
3. Jeśli nowe konto współdzieli tę historię (ten sam folder grupy), ccam uruchamia `claude --resume <session id>`, czyli właśnie tę rozmowę (`--continue` otworzyłoby najnowszą rozmowę w folderze, jakakolwiek by to była).
4. W przeciwnym razie uruchamia `claude --resume <the conversation's file> --fork-session`: Claude Code kopiuje całą rozmowę, pod nowym identyfikatorem sesji, do własnej historii nowego konta, a oryginał pozostaje bez zmian. Od tej chwili `--continue` nowego konta ją znajduje.

ccam podaje, jak duża jest rozmowa: nowe konto czyta ją całą ponownie przy pierwszej odpowiedzi. Pamięć podręczna promptów nigdy nie jest współdzielona między organizacjami, więc z `--cross-org` ta pierwsza odpowiedź kosztuje pełny rozmiar. Współdzielenie `history` w grupie pozwala uniknąć kopii, zachowuje jedną rozmowę i współdzieli też punkty kontrolne, z których `/rewind` przywraca kod (`file-history/`).

## W Claude Code

[Plugin ccam](/docs/install#claude-code-plugin) dodaje dwa hooki (Claude Code 2.1.78 lub nowszy):

- Gdy tura zatrzyma się na limicie użycia (`You've hit your session limit`, limit tygodniowy lub limit modelu), zapisuje rozmowę i pokazuje powiadomienie na pulpicie w terminalach, które to obsługują (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): `ccam: work hit its 5-hour limit (resets 3:45pm). ccam next --continue carries on as another account.` Krótkie ograniczenie po stronie serwera albo limit szybkości klucza API nie jest limitem użycia i zostaje pominięte. Claude Code 2.1.234 i nowsze także czekają w sesji i po resecie okna 5-godzinnego same kontynuują; powiadomienie daje ci wybór, żeby zamiast tego przełączyć konto.
- `/ccam:switch [account]` zapisuje rozmowę i kończy turę, zanim zostanie wysłane jakiekolwiek żądanie: wpisz `/exit`, a potem `ccam next --continue` kontynuuje ją jako konto, które podałeś, albo jako następne.

Notatka jest przechowywana w `~/.local/state/ccam/handoff/<account>.json` przez 12 godzin i używana jeden raz, gdy ma wystartować następna sesja; uruchomienie, które się nie powiedzie, zostawia ją na następną próbę. Należy ona do folderu, w którym sesja się zaczęła, nawet po `cd` w powłoce Claude Code. Jeśli wskazane w niej konto zostało w międzyczasie przemianowane lub usunięte, ccam o tym informuje i wybiera jak zwykle. Z `relay = true` w `config.toml` ccam pozostaje obok sesji interaktywnych, zamiast zastępować się procesem Claude Code, i po `/ccam:switch` sam przenosi rozmowę tej sesji dalej (na notatkę zostawioną przez inną sesję nie reaguje). Kosztuje to jeden czekający proces ccam na sesję, dlatego domyślnie jest wyłączone. Przekazana sesja startuje z `ccam next --continue`, bez argumentów, z którymi uruchomiono pierwszą, takich jak `--model`.

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
8. Próbowane są tylko konta z tej samej organizacji, tak jak w `--auto` (`cross_org = true` znosi to ograniczenie).

Następne konto jest wybierane według tych samych zasad co w `--auto` (`auto_max_usage`, `auto_strategy`, okno `--model`). Gdy nie zostanie żadne konto, ccam wypisuje wyjście ostatniej próby, informuje, że żadne zalogowane konto nie ma już zapasu, i kończy działanie z kodem 75. Z `--wait` przed nazwą konta (`ccam --fallback --wait=2h work -p …`) zamiast tego czeka na reset pierwszego konta, a potem może ponownie spróbować każdego konta; zawsze najpierw czeka co najmniej raz, żeby konto, które właśnie osiągnęło limit, nie zostało od razu uruchomione ponownie. `ccam --dry-run --fallback work -p …` wypisuje polecenie pierwszej próby.

Sesje interaktywne nie przełączają się same. Zamiast tego, po `ccam statusline install`, pasek stanu dodaje `→ ccam next --auto --continue`, gdy użycie 5-godzinne lub tygodniowe osiągnie 90% ([Limity użycia](/features/usage)).

## Bezpieczeństwo i przypadki brzegowe

- Konta bez loginu nigdy nie są wybierane. Gdy nie ma innego zalogowanego konta, samo `ccam next` kończy się komunikatem `no logged-in account other than <name>`.
- Z `--continue`, gdy w folderze nie ma żadnej rozmowy: `no account has a conversation in <folder> to continue`.
- Gdy zasada organizacji nie zostawia żadnego konta: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Z `--auto`, gdy nie ma zalogowanego konta: `no logged-in account to pick`. Gdy każde konto jest pełne lub ma nieznane użycie: `` no account with known usage and headroom (see `ccam usage`) ``. W obu przypadkach kod wyjścia to 75.
- `--max-usage`, `--strategy`, `--wait` i `--cross-org` działają tylko z `--auto`, a wartość spoza zakresu 1-100 lub inna strategia jest odrzucana.
- Fork kopiuje rozmowę, a nie pozostały stan konta: punkty kontrolne z `file-history/` pierwotnego konta nie są przenoszone, chyba że historia jest współdzielona.
- `--auto` wysyła zapytania sieciowe nawet z `--dry-run`. Tylko odczytuje i nigdy nie zapisuje plików z tokenami.
- Wybrane konto startuje tą samą ścieżką co `ccam <name>`, więc obowiązują te same sprawdzenia, w tym ostrzeżenie, gdy w tym folderze działa już inny Claude Code współdzielący tę samą historię.
- `last-used` jest zapisywany tylko wtedy, gdy Claude Code rzeczywiście wystartuje. `--dry-run` go nie zmienia, więc próbne uruchomienie nie przesuwa rotacji.
- `--fallback` wstrzymuje standardowe wyjście do końca próby, więc wyjście strumieniowe, takie jak `--output-format stream-json`, przychodzi w całości na końcu. Pliki zmienione przez pierwszą próbę przed osiągnięciem limitu pozostają zmienione; następna próba zaczyna prompt od nowa.

## Powiązane

- [Limity użycia](/features/usage)
- [Limity użycia i przełączanie kont](/docs/usage-limits)
- [Współdzielenie serwerów MCP, ustawień i historii](/docs/sharing)
- [Praca z limitem użycia](/guides/usage-limit-workflow)
