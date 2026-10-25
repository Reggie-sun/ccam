---
title: Limity użycia
description: ccam usage pokazuje naraz 5-godzinne i tygodniowe limity każdego konta, a ccam statusline wyświetla w Claude Code nazwę konta i jego użycie.
updated: 2026-10-25
order: 7
translation: machine
sourceHash: "91ad20327289"
---

`ccam usage` pokazuje, jaką część 5-godzinnego limitu, tygodniowego limitu i własnego tygodniowego limitu każdego modelu (takiego jak Fable) wykorzystało każde konto, jako pasek dla każdego limitu wraz z czasem do jego resetu. Używaj go przed rozpoczęciem długiego zadania albo gdy Claude Code informuje, że limit jest blisko, żeby zobaczyć, które konto ma jeszcze zapas. `ccam statusline` to polecenie paska stanu dla Claude Code: pokazuje nazwę konta i użycie w Claude Code i zapisuje liczby dla `ccam usage`.

## Użycie

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam usage --source claude  # let Claude Code look it up: ccam sends no token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

Każde konto ma wiersz ze swoim numerem, nazwą i informacją, jak stare są jego liczby; `(cached)` oznacza, że pochodzą z zapisanego pliku, a nie z endpointu. Pod nim każdy limit ma własny wiersz z 20-polowym paskiem, procentem i czasem do resetu okna: `5-hour`, `Weekly` (wszystkie modele), a potem własne tygodniowe okno każdego modelu, takie jak `Fable`. Okno, którego konto nie ma, pokazuje `-`. Procent zmienia kolor na żółty od 80% i na czerwony przy 100%. Okno, które Claude Code zgłosił jako pełne, gdy odrzucił żądanie na limicie użycia, pokazuje 100% aż do resetu, bez względu na zapisane liczby (zobacz [zapis limitu](#the-limit-record)). W terminalu każde konto pojawia się od razu jako `fetching…`, a jego wiersze są rysowane, gdy odpowie. Ostatni wiersz wskazuje konto, które [`ccam next --auto`](/features/next) wybrałby według tych samych zasad (`auto_max_usage`, `auto_strategy`, uruchomione sesje, tylko konta organizacji ostatnio używanego konta, chyba że ustawiono `cross_org`, z pominięciem każdego okna na poziomie 100%, więc konto z wyczerpanym tygodniem nigdy nie jest wskazywane), i pojawia się tylko wtedy, gdy masz więcej niż jedno konto. Przy domyślnej strategii `headroom` brzmi `→ Most headroom: <name>`; przy `expiring` lub `balance` podaje przyczynę, na przykład `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**Miejsca w planach Team i Enterprise.** Automatyczny wybór konta (`ccam next`, `--fallback`, `ccam pick`, klawisz `n` w panelu oraz [`auto_switch`](/features/next#auto-switch), gdy go włączysz) przełącza też między miejscami w tej samej organizacji Claude Team lub Enterprise. Centrum pomocy Anthropic podaje: „Each seat on a Team plan is assigned to one named individual, and only that person may use it” (każde miejsce w planie Team jest przypisane do jednej imiennie wskazanej osoby i tylko ta osoba może z niego korzystać) oraz „Users can't route their work through another user's account to use that user's remaining usage” (użytkownicy nie mogą kierować swojej pracy przez konto innego użytkownika, aby wykorzystać pozostały limit użycia tego użytkownika) ([plany Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); podobnie w przypadku [planów Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Dlatego korzystanie z automatycznego wyboru między miejscami może naruszać warunki Anthropic. Czy korzystać z niego w tej sytuacji, decydujesz ty: zanim zaczniesz na nim polegać, skonsultuj to ze swoją organizacją i z Anthropic. Planów prywatnych (Pro, Max) dotyczy [uwaga o warunkach w FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

Konto, które nie jest zalogowane, pokazuje `not logged in → ccam login <name>`, a konto z kluczem API pokazuje `no limits (API)`. Wygasły token dostępu pokazuje ``access token expired; `ccam <name>` renews it when it starts``, nad zapisanymi liczbami, jeśli jakieś są: Claude Code odnowi go przy następnym uruchomieniu. Login, który wygasł w całości, bo wygasł jego token odświeżania albo Claude Code wyczyścił go po odrzuceniu przez serwer, pokazuje `` login expired → `ccam login <name>` ``.

Po `ccam statusline install` Claude Code pokazuje wiersz taki jak `[work] · Opus · 5h 24% · 7d 41%`.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--json` | Wypisuje raporty jako JSON: `name`, `slot`, `status`, `limits` (z `credits`, gdy konto ma włączone kredyty użycia), `source`, `at`, `error` oraz `limited` i `limited_windows`, gdy zapis limitu ustawia okna na 100%. Bez kont wypisuje `[]` | wyłączone |
| `--watch[=N]` | Odświeża co `N` sekund, aż naciśniesz Ctrl+C. Samo `--watch` albo wartość 0 lub mniejsza oznacza 60. Wartość podaje się po `=`: `--watch 10` jest odrzucane z podpowiedzią | wyłączone |
| `-v`, `--verbose` | Dodaje pod każdym kontem wiersz `source:`: skąd pochodzą jego liczby (endpoint użycia, Claude Code, pasek stanu, zapisany plik albo limit użycia zgłoszony przez Claude Code). Gdy nad liczbami leży zapis limitu, pojawia się po nim wiersz taki jak `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` | wyłączone |
| `--no-cache` | Ignoruje 60-sekundową pamięć podręczną i pyta endpoint ponownie | wyłączone |
| `--source live\|claude\|cache` | `live` pyta endpoint z tokenem konta; `claude` zleca sprawdzenie Claude Code (poniżej); `cache` pokazuje tylko zapisane liczby | `usage_source`, w przeciwnym razie `live` |
| `--refresh` | Odnawia wygasłe tokeny dostępu, tylko dla kont, które nie są uruchomione. Odnowienie może potrwać do 30 sekund | wyłączone |
| `--dry-run` | Akceptowana: `ccam usage` nie zmienia niczego poza własną pamięcią podręczną. Odrzucana razem z `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | wyłączone |

Każda inna wartość `--source` kończy się błędem `--source must be live, claude or cache`.

### Użycie z samego Claude Code

Z `usage_source = "claude"` na górze `config.toml` (nad pierwszym `[[group]]`, `[[profile]]` lub `[[folder]]`) albo z `--source claude` ccam w ogóle nie odczytuje tokena logowania. Dla każdego konta uruchamia `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` w folderze konta i wysyła na standardowe wejście żądanie `get_usage` Claude Code (Claude Code 2.1.270 lub nowszy); Claude Code sprawdza użycie z własnym loginem i odpowiada oknami 5-godzinnym, tygodniowym i modeli. Dzięki `--safe-mode` hooki, serwery MCP i pluginy konta się nie uruchamiają, a żadna sesja nie jest zapisywana. Każde konto to jedno uruchomienie Claude Code, około sekundy, więc to źródło jest wolniejsze niż `live`; konta są pytane naraz, każde z limitem 20 sekund, a 60-sekundowa pamięć podręczna działa jak zwykle. Login, którego token dostępu wygasł lub wygaśnie w ciągu 5 minut, dostaje 50 sekund, bo Claude Code odnawia go przy starcie, a przerwane odnowienie może stracić login. Zapytanie, z którego ccam rezygnuje, prosi Claude Code o zakończenie sygnałem SIGTERM, więc ten dokończy rozpoczęte odnowienie (Windows nie ma takiego sygnału; tam proces jest zatrzymywany). Login, który wygasł w całości, pokazuje `` login expired → `ccam login <name>` `` bez uruchamiania Claude Code. Claude Code oznacza to żądanie jako eksperymentalne, więc jeśli przestanie odpowiadać, ccam pokazuje zapisane liczby, a `ccam doctor` to zgłasza. Gdy Claude Code odpowie, że jego własne sprawdzenie się nie powiodło, ccam pokazuje zapisane liczby z komunikatem `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)`, a nie plan bez limitów. `--dry-run` nie uruchamia Claude Code (bo ten zapisuje w folderze konta) i pokazuje zapisane liczby. `usage_source` przyjmuje `live` lub `claude`; `cache` to wyłącznie wartość `--source`. `ccam usage` nie przyjmuje argumentów.

`ccam statusline` odczytuje JSON paska stanu Claude Code ze standardowego wejścia. Jego jedyna opcja to `--wrap <command>`: istniejące polecenie paska stanu, przed którego pierwszym wierszem wyniku pojawia się nazwa konta. `ccam statusline install` nie przyjmuje opcji. Oba odrzucają `--dry-run` komunikatem `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Jak to działa

Liczby pochodzą z dwóch miejsc: z nieudokumentowanego endpointu użycia Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, wywoływanego z własnym loginem każdego konta, oraz z wartości, które Claude Code przekazuje poleceniu paska stanu. Cały ten kod jest w `internal/usage`.

1. `cli.usageCmd` wczytuje `config.toml` i tworzy jeden `usage.Account` na konto: jego folder `~/.local/share/ccam/profiles/<name>/`, znaczniki sesji (razem z własnymi zapisami Claude Code w folderze konta mówią, czy konto działa) i informację, czy zostało zaimportowane w miejscu.
2. `Fetcher.FetchAll` uruchamia jedną gorutynę na konto. Wszystkie współdzielą jeden `http.Client`, więc połączenia z endpointem są używane ponownie. Nawiązanie połączenia może trwać 2 s, a każde konto ma łącznie 3 s. Klient nie podąża za przekierowaniami, bo przekierowanie wysłałoby token gdzie indziej, i odczytuje najwyżej 1 MB odpowiedzi.
3. Gdy każde konto kończy, jego wiersz jest odświeżany. Konto, którego zapytanie się nie powiedzie lub które nie odpowie na czas, korzysta z zapisanych liczb, więc jedno wolne konto nigdy nie wstrzymuje pozostałych.

Każde konto odpytuje endpoint we własnej gorutynie; to, które przekroczy czas, pokazuje zamiast tego zapisane liczby.

<!-- motion:usage-fetch -->

`personal` i `work` odpowiadają na czas i pokazują aktualne liczby. `side` potrzebuje więcej niż 3 s, więc `fallback()` pokazuje zamiast tego jego zapisane liczby. `usage.Table()` rysuje blok dla każdego konta, a `usage.Best()` wybiera konto wskazane w ostatnim wierszu.

Dla jednego konta `Fetcher.Fetch` wykonuje te sprawdzenia w kolejności:

<!-- motion:usage-cache -->

1. Z `--source cache` przechodzi od razu do `fallback()`.
2. Jeśli nie ustawiono `--no-cache`, aktualny wynik w `~/.local/state/ccam/usage/<name>.json` zmierzony mniej niż 60 sekund temu (`CacheTTL`) jest używany i oceniany według bieżącej chwili: okno, które się w międzyczasie zresetowało, pokazuje 0%, a plik z datą późniejszą niż czas zegara nie liczy się jako świeży. Dwukrotne uruchomienie `ccam usage` pod rząd kosztuje jedno zapytanie na konto.
3. `readCredentials()` odczytuje `.credentials.json` (lub pęk kluczy macOS) bez zmieniania go i dekoduje tylko klucz `claudeAiOauth`. Brak loginu subskrypcji oznacza `not logged in` albo `no limits (API)`, gdy `.claude.json` zawiera klucz API. Login, którego token odświeżania wygasł albo który wyczyścił Claude Code, to `` login expired → `ccam login <name>` `` i nie jest sprawdzany.
4. Wygasły token dostępu nie jest używany. Bez `--refresh` konto pokazuje swoje zapisane liczby z ``access token expired; `ccam <name>` renews it when it starts``.
5. W przeciwnym razie ccam wywołuje endpoint. Kod 401 lub 403 jest traktowany jako wygasły login (`login expired: authentication failed (401)`), każdy inny status lub treść, której nie da się odczytać, jako błąd, i w obu przypadkach używane są zapisane liczby. Zapytanie przerwane przez Ctrl+C wyświetla `cancelled`. Kod 429 oznacza, że endpoint dostaje zbyt wiele zapytań: ccam zapisuje `<name>.pause` z czasem kolejnego zapytania (`Retry-After` albo 5 minut, najwyżej godzina), pokazuje zapisane liczby z komunikatem `rate limited by the usage endpoint; asking again after 15:10` i do tego czasu nie wywołuje endpointu dla tego konta, nawet z `--no-cache`. Natychmiastowe ponawianie tylko podtrzymywałoby limit.
6. Poprawna odpowiedź jest zapisywana w `<name>.json` razem z czasem pomiaru.

`fallback()` bierze najnowsze zapisane sprawdzenie (`<name>.json` albo `<name>.claude.json` ze sprawdzeń przez Claude Code), a gdy `<name>.statusline.json` zapisano później, nakłada na nie 5-godzinne i tygodniowe okno z paska stanu; okna modeli, których pasek stanu nie ma, zostają takie, jak zapisało je sprawdzenie. Zanim pokaże zapisane liczby, sprawdza czas resetu każdego okna: okno, które zresetowało się od zapisania liczb, pokazuje 0%, bo stary procent już nie obowiązuje. Dlatego `personal` powyżej pokazuje `0%` bez czasu resetu.

Endpoint ma dwa formaty odpowiedzi. Nowszy ma tablicę `limits` z wpisami `session` (5-godzinny), `weekly_all` i `weekly_scoped` (dla poszczególnych modeli). Starszy ma obiekty `five_hour`, `seven_day`, `seven_day_opus` i `seven_day_sonnet`. `ParseLive` odczytuje oba: tablica `limits` ma pierwszeństwo, a stare obiekty uzupełniają to, czego ona nie obejmuje. Treść bez żadnych limitów jest zgłaszana jako zmieniony format, o czym informuje `ccam doctor`.

### Zapis limitu

Gdy Claude Code zatrzyma turę na limicie użycia, hook [pluginu ccam](/docs/install#claude-code-plugin) odczytuje z rozmowy odrzucone żądanie: które okno było pełne (5-godzinne, tygodniowe, Opus lub Sonnet) i kiedy się resetuje. Zapisuje to w `~/.local/state/ccam/usage/<name>.limit.json`. `Fetch` nakłada ten zapis na wszystko inne, co wie, z każdego źródła: każde okno w zapisie liczy się jako 100% aż do resetu, bez względu na to, jak stare lub niskie są zapisane liczby, więc `ccam next --auto`, `ccam pick`, `--fallback` i [`auto_switch`](/features/next#auto-switch) nie wybiorą tego konta wcześniej. Gdy nic innego nie wiadomo, konto pokazuje sam zapis (`source: the usage limit Claude Code reported` z `-v`). Sesje, przy których nikt nie siedzi (`-p`, sesje w tle), nie tworzą zapisu, podobnie jak zatrzymanie przez organizację (brak użycia, grupa ustawiona na $0), które jej administratorzy mogą w każdej chwili cofnąć. `ccam login` na inne konto Claude usuwa zapis; usunięcie pliku także go czyści.

### Odnawianie tokenów

Tokeny odświeżania są rotowane: odnowienie daje nowy token odświeżania i unieważnia stary. Gdyby ccam odnowił login, gdy Claude Code działa jako to samo konto, ta sesja zostałaby wylogowana. Dlatego `ccam usage` odczytuje tokeny i nigdy ich nie zapisuje, chyba że podasz `--refresh`. Wtedy `Fetcher.refresh`:

1. odrzuca konta zarejestrowane w miejscu przez `ccam import`, bo zwykłe `claude` może używać tego loginu bez wiedzy ccam;
2. odrzuca loginy przechowywane w pęku kluczy macOS;
3. bierze plik blokady `~/.local/state/ccam/locks/profile-<name>.refresh.lock`, żeby dwa procesy ccam nigdy nie odnowiły tego samego tokenu dwa razy; dopóki jest trzymany, konto pokazuje `token not renewed: another ccam is using this login right now (a refresh, rename, remove or logout)`. `ccam rename`, `remove` i `logout` też go biorą i odmawiają, gdy trwa odnowienie;
4. pod blokadą odmawia, jeśli działa sesja tego konta, niezależnie od sposobu uruchomienia (przez ccam, zwykłe `claude`, IDE albo sesja w tle i stojący za nią demon, który sam odnawia login), i pomija zapytanie, jeśli inny proces już odnowił token;
5. nic nie wysyła, dopóki Claude Code trzyma własne blokady odnowienia w folderze konta (`.oauth_refresh.lock`, `<folder>.lock`): Claude Code właśnie odnawia ten login;
6. wysyła token odświeżania do endpointu tokenów i czeka na odpowiedź do 30 sekund, tyle, ile Claude Code daje własnym odnowieniom. Wysłane żądanie nigdy nie jest przerywane, nawet przez Ctrl+C, bo serwer mógł już wymienić token odświeżania, a login utrzymuje tylko jego odpowiedź: Ctrl+C lub rozłączenie kończy pozostałe zapytania, a ccam kończy działanie (130 po Ctrl+C), gdy odpowiedź zostanie zapisana;
7. zapisuje pod blokadą zapisu Claude Code (`.storage-write.lock`) i zmienia tylko pola tokenów w `.credentials.json` w jego bieżącej postaci: to, co Claude Code w międzyczasie zapisał (logowanie MCP), zostaje, login odnowiony w międzyczasie przez Claude Code jest zachowywany, a wylogowanie nigdy nie jest cofane. Gdy odpowiedź podaje czas życia nowego tokenu odświeżania, jest on zapisywany jako `refreshTokenExpiresAt`, tak jak robi to Claude Code.

Odnowienie, które się nie udało, podaje przyczynę w wierszu konta albo nad jego zapisanymi liczbami: na przykład `token not renewed: a session is running` albo `login may be lost: …`, gdy żądanie zostało wysłane, ale jego odpowiedzi nie dało się użyć ani zapisać, z radą, żeby uruchomić `ccam login <name>`, jeśli Claude Code potem poprosi o zalogowanie.

### Zbieranie danych z paska stanu

Claude Code uruchamia polecenie paska stanu i przekazuje mu potokiem dokument JSON. ccam odczytuje z niego nazwę modelu i `rate_limits`:

<!-- motion:usage-statusline -->

1. `runStatusline` znajduje konto na podstawie `CCAM_PROFILE`, które ccam ustawia przy uruchamianiu Claude Code, a w przeciwnym razie przez dopasowanie `CLAUDE_CONFIG_DIR` do folderów kont.
2. `ParseStatusLine` odczytuje `rate_limits.five_hour` i `rate_limits.seven_day`. Dla kont z kluczem API oraz przed pierwszą odpowiedzią w sesji JSON nie zawiera limitów; wtedy nic nie jest zapisywane.
3. `SaveStatusLine` zapisuje `~/.local/state/ccam/usage/<name>.statusline.json`. `ccam usage` korzysta z niego, gdy endpoint nie odpowiada.
4. ccam wypisuje `[name] · model · 5h N% · 7d N%`. Z `--wrap` uruchamia twoje polecenie przez `sh -c` z tym samym wejściem i wypisuje jego wynik tak, jak Claude Code pokazuje pasek stanu: każdy niepusty wiersz, z `[name]` przed pierwszym. Polecenie, które się nie powiedzie, pokazuje samą nazwę konta. W Windows uruchamia polecenie tak jak Claude Code: przez Git Bash, jeśli jest zainstalowany, a w przeciwnym razie przez PowerShell.
5. Gdy 5-godzinne lub tygodniowe użycie wynosi 90% lub więcej i istnieje inne konto, wiersz kończy się `→ ccam next --auto --continue`, na przykład `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

Opakowane polecenie nie ma własnego limitu czasu: Claude Code daje poleceniu paska stanu 10 minut i zatrzymuje przebieg, gdy zaczyna się nowsza aktualizacja, więc wolne polecenie wygląda tak, jak bez ccam. Gdy Claude Code kończy ccam, polecenie też się kończy (na Linux i macOS razem z każdym procesem, który uruchomiło).

`ccam statusline install` zapisuje `statusLine` w `settings.json` każdego konta i zachowuje pozostałe klucze obiektu, takie jak `padding`, `refreshInterval` i `hideVimModeIndicator`. Gdy konto ma już polecenie paska stanu, ccam zachowuje je jako `--wrap '<old command>'`; w Windows jako `--wrap-base64 <base64>`, które w żadnej powłoce nie wymaga cudzysłowów. Pasek stanu, który już uruchamia ccam, z innej ścieżki albo w cudzysłowach zapisanych przez starszy ccam, jest zastępowany zamiast opakowywany, z zachowaniem tego, co opakowuje. W Windows ścieżka do ccam jest zapisywana z ukośnikami zwykłymi i `~` w miejsce folderu domowego (`~/.local/bin/ccam.exe statusline`), co rozumieją zarówno Git Bash, jak i PowerShell. Ścieżka wymagająca cudzysłowów jest zapisywana w całości dla powłoki, której używa Claude Code: `'…'` dla Git Bash, `& '…'` dla PowerShell. `settings.json` współdzielony przez grupę jest zapisywany raz, przez swój link, więc dostaje go każde konto w grupie. Ponowne uruchomienie wypisuje `(already set)` dla plików, które już używają ccam.

## Bezpieczeństwo i przypadki brzegowe

- `ccam usage` nigdy nie zmienia pliku z tokenami, chyba że podasz `--refresh`. Gdy token dostępu wygaśnie, uruchom konto raz przez `ccam <name>`, a Claude Code go odnowi. Login, który wygasł w całości, wymaga `ccam login <name>`.
- `--refresh` nigdy nie rusza działających kont (także gdy Claude Code uruchomiono bez ccam), zaimportowanych kont ani loginów z pęku kluczy. Odrzucone odnowienie pokazuje zapisane liczby i przyczynę (`token not renewed: …`).
- Endpoint użycia nie jest publicznie udokumentowany i może się zmienić. Gdy tak się stanie, nadal widzisz zapisane liczby, a `ccam doctor` zgłasza zmianę formatu.
- Gdy nie ma zapisanych liczb ani odpowiedzi, wiersz konta pokazuje `lookup failed:` i przyczynę, na przykład `no saved usage` z `--source cache`.
- Pliki pamięci podręcznej zawierają tylko procenty i czasy resetu, nigdy tokeny.

## Powiązane

- [Limity użycia i przełączanie kont](/docs/usage-limits)
- [Przełączanie na następne konto](/features/next)
- [Praca z limitem użycia](/guides/usage-limit-workflow)
- [Polecenia](/docs/commands)
