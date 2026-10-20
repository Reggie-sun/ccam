---
title: Limity użycia
description: ccam usage pokazuje naraz 5-godzinne i tygodniowe limity każdego konta, a ccam statusline wyświetla w Claude Code nazwę konta i jego użycie.
updated: 2026-10-20
order: 7
translation: machine
sourceHash: 525174ceed0d
---

`ccam usage` pokazuje, jaką część 5-godzinnego i tygodniowego limitu wykorzystało każde konto, wraz z czasem do resetu każdego okna. Używaj go przed rozpoczęciem długiego zadania albo gdy Claude Code informuje, że limit jest blisko, żeby zobaczyć, które konto ma jeszcze zapas. `ccam statusline` to polecenie paska stanu dla Claude Code: pokazuje nazwę konta i użycie w Claude Code i zapisuje liczby dla `ccam usage`.

## Użycie

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add per-model weekly limits and the data source
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ░░░░░░░░   0% -         ██████░░  70% 2d 16h    9h 36m ago (cached)
 2 work      ██░░░░░░  24% 2h 29m    ███░░░░░  41% 3d 11h    just now (cached)
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

W terminalu tabela pojawia się od razu z `fetching…` w każdym wierszu, a każdy wiersz jest odświeżany, gdy jego konto odpowie. `As of` mówi, jak stare są liczby; `(cached)` oznacza, że pochodzą z zapisanego pliku, a nie z endpointu. Ostatni wiersz wskazuje konto z najniższym użyciem 5-godzinnym poniżej 100% i pojawia się tylko wtedy, gdy masz więcej niż jedno konto. Konto, które nie jest zalogowane, pokazuje `not logged in → ccam login <name>`, a konto z kluczem API pokazuje `no limits (API)`.

Po `ccam statusline install` Claude Code pokazuje wiersz taki jak `[work] · Opus · 5h 24% · 7d 41%`.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--json` | Wypisuje raporty jako JSON: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | wyłączone |
| `--watch[=N]` | Odświeża co `N` sekund, aż naciśniesz Ctrl+C. Samo `--watch` albo wartość 0 lub mniejsza oznacza 60 | wyłączone |
| `-v`, `--verbose` | Pokazuje tygodniowe limity dla poszczególnych modeli i wiersz `source:` (`live`, `cache` lub `statusline`) | wyłączone |
| `--no-cache` | Ignoruje 60-sekundową pamięć podręczną i pyta endpoint ponownie | wyłączone |
| `--source live\|cache` | `cache` nigdy nie pyta endpointu i pokazuje tylko zapisane liczby; `live` to zachowanie domyślne | live |
| `--refresh` | Odnawia wygasłe tokeny dostępu, tylko dla kont, które nie są uruchomione | wyłączone |
| `--dry-run` | Akceptowana: `ccam usage` nie zmienia niczego poza własną pamięcią podręczną. Odrzucana razem z `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | wyłączone |

Każda inna wartość `--source` kończy się błędem `--source must be live or cache`. `ccam usage` nie przyjmuje argumentów.

`ccam statusline` odczytuje JSON paska stanu Claude Code ze standardowego wejścia. Jego jedyna opcja to `--wrap <command>`: istniejące polecenie paska stanu, przed którego pierwszym wierszem wyniku pojawia się nazwa konta. `ccam statusline install` nie przyjmuje opcji. Oba odrzucają `--dry-run` komunikatem `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Jak to działa

Liczby pochodzą z dwóch miejsc: z nieudokumentowanego endpointu użycia Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, wywoływanego z własnym loginem każdego konta, oraz z wartości, które Claude Code przekazuje poleceniu paska stanu. Cały ten kod jest w `internal/usage`.

1. `cli.usageCmd` wczytuje `config.toml` i tworzy jeden `usage.Account` na konto: jego folder `~/.local/share/ccam/profiles/<name>/`, folder sesji (żeby sprawdzić, czy konto działa) i informację, czy zostało zaimportowane w miejscu.
2. `Fetcher.FetchAll` uruchamia jedną gorutynę na konto. Wszystkie współdzielą jeden `http.Client`, więc połączenia z endpointem są używane ponownie. Nawiązanie połączenia może trwać 2 s, a każde konto ma łącznie 3 s. Klient nie podąża za przekierowaniami, bo przekierowanie wysłałoby token gdzie indziej, i odczytuje najwyżej 1 MB odpowiedzi.
3. Gdy każde konto kończy, jego wiersz jest odświeżany. Konto, którego zapytanie się nie powiedzie lub które nie odpowie na czas, korzysta z zapisanych liczb, więc jedno wolne konto nigdy nie wstrzymuje pozostałych.

Każde konto odpytuje endpoint we własnej gorutynie; to, które przekroczy czas, pokazuje zamiast tego zapisane liczby.

<!-- motion:usage-fetch -->

`personal` i `work` odpowiadają na czas i pokazują aktualne liczby. `side` potrzebuje więcej niż 3 s, więc `fallback()` odczytuje jego dwa pliki pamięci podręcznej i używa tego, który zapisano później. `usage.Table()` rysuje wiersze, a `usage.Best()` wybiera konto wskazane w ostatnim wierszu.

Dla jednego konta `Fetcher.Fetch` wykonuje te sprawdzenia w kolejności:

<!-- motion:usage-cache -->

1. Z `--source cache` przechodzi od razu do `fallback()`.
2. Jeśli nie ustawiono `--no-cache`, aktualny wynik w `~/.local/state/ccam/usage/<name>.json` młodszy niż 60 sekund (`CacheTTL`) jest zwracany bez zmian. Dwukrotne uruchomienie `ccam usage` pod rząd kosztuje jedno zapytanie na konto.
3. `readCredentials()` odczytuje `.credentials.json` (lub pęk kluczy macOS) bez zmieniania go i dekoduje tylko klucz `claudeAiOauth`. Brak loginu subskrypcji oznacza `not logged in` albo `no limits (API)`, gdy `.claude.json` zawiera klucz API.
4. Wygasły token dostępu nie jest używany. Bez `--refresh` konto pokazuje swoje zapisane liczby oznaczone jako wygasłe.
5. W przeciwnym razie ccam wywołuje endpoint. Kod 401 lub 403 jest traktowany jako wygasły login, każdy inny status lub treść, której nie da się odczytać, jako błąd, i w obu przypadkach używane są zapisane liczby. Kod 429 oznacza, że endpoint dostaje zbyt wiele zapytań: ccam zapisuje `<name>.pause` z czasem kolejnego zapytania (`Retry-After` albo 5 minut, najwyżej godzina), pokazuje zapisane liczby z komunikatem `rate limited by the usage endpoint; asking again after 15:10` i do tego czasu nie wywołuje endpointu dla tego konta, nawet z `--no-cache`. Natychmiastowe ponawianie tylko podtrzymywałoby limit.
6. Poprawna odpowiedź jest zapisywana w `<name>.json` razem z czasem pomiaru.

`fallback()` odczytuje zarówno `<name>.json`, jak i `<name>.statusline.json` i zachowuje nowszy. Zanim pokaże zapisane liczby, sprawdza czas resetu każdego okna: okno, które zresetowało się od zapisania liczb, pokazuje 0%, bo stary procent już nie obowiązuje. Dlatego `personal` powyżej pokazuje `0%` bez czasu resetu.

Endpoint ma dwa formaty odpowiedzi. Nowszy ma tablicę `limits` z wpisami `session` (5-godzinny), `weekly_all` i `weekly_scoped` (dla poszczególnych modeli). Starszy ma obiekty `five_hour`, `seven_day`, `seven_day_opus` i `seven_day_sonnet`. `ParseLive` odczytuje oba: tablica `limits` ma pierwszeństwo, a stare obiekty uzupełniają to, czego ona nie obejmuje. Treść bez żadnych limitów jest zgłaszana jako zmieniony format, o czym informuje `ccam doctor`.

### Odnawianie tokenów

Tokeny odświeżania są rotowane: odnowienie daje nowy token odświeżania i unieważnia stary. Gdyby ccam odnowił login, gdy Claude Code działa jako to samo konto, ta sesja zostałaby wylogowana. Dlatego `ccam usage` odczytuje tokeny i nigdy ich nie zapisuje, chyba że podasz `--refresh`. Wtedy `Fetcher.refresh`:

1. odrzuca konta zarejestrowane w miejscu przez `ccam import`, bo zwykłe `claude` może używać tego loginu bez wiedzy ccam;
2. odrzuca loginy przechowywane w pęku kluczy macOS;
3. bierze plik blokady `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock` (w `~/.local/state/ccam/locks/`, gdy ta zmienna nie jest ustawiona), żeby dwa procesy ccam nigdy nie odnowiły tego samego tokenu dwa razy;
4. pod blokadą odmawia, jeśli działa sesja tego konta, i pomija zapytanie, jeśli inny proces już odnowił token;
5. wysyła token odświeżania do endpointu tokenów i zapisuje nowe tokeny z powrotem w `.credentials.json` przez `fsutil.AtomicWrite`, zachowując wszystkie pozostałe klucze w pliku.

### Zbieranie danych z paska stanu

Claude Code uruchamia polecenie paska stanu i przekazuje mu potokiem dokument JSON. ccam odczytuje z niego nazwę modelu i `rate_limits`:

<!-- motion:usage-statusline -->

1. `runStatusline` znajduje konto na podstawie `CCAM_PROFILE`, które ccam ustawia przy uruchamianiu Claude Code, a w przeciwnym razie przez dopasowanie `CLAUDE_CONFIG_DIR` do folderów kont.
2. `ParseStatusLine` odczytuje `rate_limits.five_hour` i `rate_limits.seven_day`. Dla kont z kluczem API oraz przed pierwszą odpowiedzią w sesji JSON nie zawiera limitów; wtedy nic nie jest zapisywane.
3. `SaveStatusLine` zapisuje `~/.local/state/ccam/usage/<name>.statusline.json`. `ccam usage` korzysta z niego, gdy endpoint nie odpowiada.
4. ccam wypisuje `[name] · model · 5h N% · 7d N%`. Z `--wrap` uruchamia twoje polecenie przez `sh -c` z tym samym wejściem i zamiast tego wypisuje `[name]`, a po nim pierwszy wiersz wyniku polecenia. W Windows uruchamia polecenie tak jak Claude Code: przez Git Bash, jeśli jest zainstalowany, a w przeciwnym razie przez PowerShell.
5. Gdy 5-godzinne lub tygodniowe użycie wynosi 90% lub więcej i istnieje inne konto, wiersz kończy się `→ ccam next --auto --continue`, na przykład `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

Opakowane polecenie ma 2 sekundy. Po tym czasie jest zabijane, a wiersz pokazuje samą nazwę konta, więc wolne polecenie nie może zablokować paska stanu Claude Code.

`ccam statusline install` zapisuje `statusLine` w `settings.json` każdego konta. Gdy konto ma już polecenie paska stanu, ccam zachowuje je jako `--wrap '<old command>'`. W Windows ścieżka do ccam jest zapisywana z ukośnikami zwykłymi i `~` w miejsce folderu domowego (`~/.local/bin/ccam.exe statusline`), co rozumieją zarówno Git Bash, jak i PowerShell. `settings.json` współdzielony przez grupę jest zapisywany raz, przez swój link, więc dostaje go każde konto w grupie. Ponowne uruchomienie wypisuje `(already set)` dla plików, które już używają ccam.

## Bezpieczeństwo i przypadki brzegowe

- `ccam usage` nigdy nie zmienia pliku z tokenami, chyba że podasz `--refresh`. Gdy login wygaśnie, uruchom konto raz przez `ccam <name>`, a Claude Code go odnowi.
- `--refresh` nigdy nie rusza działających kont, zaimportowanych kont ani loginów z pęku kluczy. Odrzucone odnowienie pokazuje zapisane liczby i przyczynę.
- Endpoint użycia nie jest publicznie udokumentowany i może się zmienić. Gdy tak się stanie, nadal widzisz zapisane liczby, a `ccam doctor` zgłasza zmianę formatu.
- Gdy nie ma zapisanych liczb ani odpowiedzi, wiersz pokazuje `lookup failed:` i przyczynę, na przykład `no saved usage` z `--source cache`.
- Pliki pamięci podręcznej zawierają tylko procenty i czasy resetu, nigdy tokeny.

## Powiązane

- [Limity użycia i przełączanie kont](/docs/usage-limits)
- [Przełączanie na następne konto](/features/next)
- [Praca z limitem użycia](/guides/usage-limit-workflow)
- [Polecenia](/docs/commands)
