---
title: "Przełączanie kont Claude Code: ccam i inne narzędzia"
description: Porównanie ccam z claude-swap, CCS, clausona, clauth, CAAM, cux i innymi narzędziami do przełączania kont Claude Code oraz które podejście do czego pasuje.
updated: 2026-10-20
order: 2
translation: machine
sourceHash: "a90ede56b767"
---

Narzędzia do przełączania kont Claude Code stosują jedno z dwóch podejść. Większość podmienia login w jednym folderze `~/.claude`, więc naraz aktywne jest tylko jedno konto. Inne, w tym ccam, dają każdemu kontu własny folder przez `CLAUDE_CONFIG_DIR`, więc każde konto pozostaje zalogowane i kilka może działać jednocześnie. ccam dodaje współdzielone ustawienia i serwery MCP oraz kontynuuje rozmowę w folderze innego konta.

## Podmiana loginu czy osobne foldery?

**Podmiana loginu** zachowuje jeden zestaw ustawień, serwerów MCP i historii, bo jest tylko jeden folder. Przełączanie jest szybkie, a działająca sesja może przejąć nowy login. Ale naraz aktywne jest tylko jedno konto, we wszystkich terminalach i edytorach jednocześnie.

Żeby podmienić login, narzędzie przechowuje kopię loginu każdego innego konta i wpisuje ją na miejsce. Claude Code wymienia token odświeżania w loginie za każdym razem, gdy go odnawia, więc każda zapisana kopia musi nadążać za działającymi sesjami. Konta wylogowane albo wymagające ponownego `/login` po przełączeniu to jedne z problemów, które te projekty poprawiają najczęściej.

**Osobne foldery** to metoda, którą [opisuje dokumentacja Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Każde konto pozostaje zalogowane i możesz uruchomić sesję służbową i prywatną obok siebie. Każdy login zostaje we własnym folderze, gdzie Claude Code odnawia go jak zwykle. Ceną jest to, że każdy folder zaczyna z własnymi ustawieniami, serwerami MCP i historią, chyba że jakieś narzędzie je współdzieli.

## Porównanie narzędzi

Na podstawie dokumentacji i kodu źródłowego każdego projektu z 2026-10-20 (clausona: v0.12.0-beta). „Nieudokumentowane” oznacza, że projekt o tym nie wspomina, a nie że brakuje funkcji.

| Narzędzie | Podejście | Kilka kont jednocześnie | Użycie 5-godzinne i tygodniowe | Po osiągnięciu limitu | Platformy |
| --- | --- | --- | --- | --- | --- |
| ccam | Osobne foldery | Tak | Tak, `ccam usage` | Powiadomienie w sesji; `ccam next --auto --continue` lub `/ccam:switch` przenosi rozmowę na konto z największym zapasem, ze współdzieloną historią lub bez niej; `ccam --fallback` sam ponawia tam zadania `-p`, a `ccam pick` rozdziela workery | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | Podmienia login; osobne foldery z `cswap run` | Tak, z `cswap run` | Tak | `cswap auto` podmienia login po przekroczeniu progu; działająca sesja trwa dalej | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Osobne foldery | Tak | Częściowo | Pule kont przekazują żądania do następnego konta przez lokalne proxy | Linux, macOS, Windows |
| [clausona](https://github.com/larcane97/clausona) | Osobne foldery | Tak | Tak | Trasy wybierają konto z limitem przed każdym uruchomieniem, mogą czekać na reset i przydzielają workerom różne konta | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Podmienia login; osobne foldery z `clauth start` | Tak, z `clauth start` | Tak | Przełącza według łańcucha zapasowych kont | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Podmienia login lub izolowane profile | Tak, z izolowanymi profilami | Tak | `caam run` ponawia nieinteraktywne polecenia na innym koncie | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Podmienia login i opakowuje `claude` | Nie, jedno aktywne konto | Tak | Podmienia login i wznawia tę samą rozmowę | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Podmienia login | Nie | Tak | Opcjonalne przełączenie po przekroczeniu progu | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Podmienia login | Nie | Tak | Opcjonalne przełączenie przy 100% | macOS |

ccswap to fork claude-swap, który zarządza też loginami Codex. [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) i [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) również dają każdemu kontu własny folder. Kilka z tych narzędzi obsługuje też innych agentów kodujących, takich jak Codex czy Gemini.

## Jak są zbudowane

| Narzędzie | Język | Podczas działania Claude Code | Kopiuje loginy między miejscami | Co konta mogą współdzielić |
| --- | --- | --- | --- | --- |
| ccam | Go, jeden statyczny plik binarny | Nic: ccam zastępuje się procesem `claude` (w Windows, który na to nie pozwala, czeka na jego zakończenie; obok może działać opcjonalny przekaźnik) | Nie | Grupy: każda wymienia elementy, które współdzieli, i swoje konta; konto może należeć do kilku |
| claude-swap | Python 3.12+, instalowany przez uv lub pipx | Nic przy `cswap run`; `cswap auto` i aplikacja na pasku menu działają dalej | Tak, gdy podmienia login | Przy `cswap run`: ustawienia, `CLAUDE.md` i rozszerzenia razem, historia opcjonalnie; serwery MCP kopiowane przy każdym uruchomieniu |
| CCS | Node.js 18+ lub Bun | CCS zostaje procesem nadrzędnym; pule dodają lokalne proxy | Nie dla swoich folderów kont; pule trzymają własne loginy | Współdzielone lub lokalne dla każdego konta, historia w grupach kontekstu; serwery MCP kopiowane |
| clausona | TypeScript, Node.js 20+ | Zwykłe `claude` przechodzi przez funkcję powłoki; `clausona run` zostawia Node jako proces nadrzędny; proces Node w tle po każdym wyjściu | Nie, ale sam odnawia tokeny nieużywanych kont, żeby odczytać ich użycie | Wszystko poza plikami należącymi do konta, z jednego konta głównego; historia opcjonalnie; serwery dodane przez `claude mcp add` nie są współdzielone |
| clauth | Rust | clauth zostaje procesem nadrzędnym | Tak, gdy podmienia login | Całe `~/.claude` albo nic |
| CAAM | Go | Zależy od trybu; opcjonalny demon | Tak, w trybie sejfu | Ustalone przez tryb; serwery MCP kopiowane przy każdym uruchomieniu |
| cux | Go | Wrapper pozostaje wokół `claude` | Tak | Wszystko: jeden folder |
| CCSwitcher, Claude Switcher | Aplikacje macOS (Swift, Python) | Aplikacja na pasku menu | Tak | Wszystko: jeden folder |

## Mocne strony ccam

### Loginy zostają tam, gdzie umieścił je Claude Code

ccam nigdy nie kopiuje loginu Claude z jednego miejsca do drugiego. Login każdego konta zostaje we własnym folderze, a na macOS we własnym elemencie pęku kluczy tego folderu, i to Claude Code go odnawia. Nie ma zapisanej kopii, która mogłaby się rozjechać z oryginałem. Istniejący `~/.claude` zarejestrowany przez `ccam import` zachowuje element pęku kluczy, którego używa zwykły `claude`, więc pozostaje zalogowany w obu przypadkach.

ccam sam odnawia wygasły token tylko wtedy, gdy uruchomisz `ccam usage --refresh`. Najpierw zakłada blokadę i odmawia, gdy to konto działa, więc nigdy nie wymienia tokenu pod działającą sesją.

### Nic nie działa między tobą a Claude Code

ccam czyta jeden plik, sprawdza swoje linki i zastępuje się procesem `claude` (`exec`). Podczas pracy nie ma procesu ccam, wrappera, proxy ani demona: standardowe wejście, sygnały i kod wyjścia należą do Claude Code, a jego żądania trafiają prosto do Anthropic. Windows nie potrafi zastąpić procesu, więc tam ccam uruchamia Claude Code w tej samej konsoli i czeka na jego kod wyjścia, nie robiąc nic więcej.

Uruchomienie konta w ten sposób dodaje około 4 ms, poniżej 6 ms w 95. percentylu, przy szczytowym zużyciu pamięci 8 MB ([`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) z hyperfine, 300 uruchomień). Na tym samym laptopie samo uruchomienie interpretera Python zajmuje 13 ms, a Node.js 31 ms, zanim program cokolwiek zrobi.

### Grupy decydują, co jest współdzielone

Serwery MCP, `settings.json`, `CLAUDE.md`, rozszerzenia (agenci, polecenia, umiejętności, wtyczki) i historia rozmów to pięć osobnych elementów. Grupa wymienia elementy, które współdzieli, i należące do niej konta, więc dwa konta służbowe mogą współdzielić historię w grupie `work`, konto prywatne zachowuje własną, a wszystkie trzy współdzielą te same serwery MCP przez `default`. `ccam group set work --add memory` przenosi wszystkie konta grupy w jednym planie; [panel](/features/dashboard) robi to samo za pomocą pól wyboru i wyszarza element, który konto dostaje już z innej grupy.

- Serwery MCP są przekazywane przy uruchomieniu przez `--mcp-config`, więc domyślnie ccam nigdy nie edytuje `.claude.json`, pliku, który Claude Code nadpisuje w trakcie działania. `ccam mcp add` dodaje serwer raz dla całej grupy.
- Pozostałe elementy to linki do folderu grupy, więc zmiana wprowadzona na jednym koncie od razu pojawia się na pozostałych. Przy uruchomieniu nic nie jest kopiowane.
- Loginy, `.claude.json` i stan działania nigdy nie są współdzielone.

### Każda zmiana jest najpierw pokazana i ma kopię zapasową

`ccam group` i `ccam migrate` wypisują plan i pytają, zanim cokolwiek przeniosą. Gdy plik różni się między kontem a jego grupą, wybierasz: zachować wersję grupy, użyć swojej albo scalić JSON. Oryginały trafiają do `.ccam-backup/<time>/`, folder usuniętego konta trafia do folderu kosza, a pliki są podmieniane jedną zmianą nazwy, więc przerwany zapis zostawia stary plik. `--dry-run` pokazuje, co zrobiłoby uruchomienie, współdzielenie, migracja, skróty, aktualizacje i naprawy. `ccam doctor --fix` naprawia zepsute linki i uprawnienia.

### Ta sama rozmowa na następnym koncie

`ccam next --auto --continue` pyta wszystkie konta naraz o użycie, otwiera to z największym zapasem i kontynuuje rozmowę. Przy współdzielonej historii to ta sama rozmowa; bez niej Claude Code kopiuje ją na nowe konto (`--resume <file> --fork-session`), a ccam podaje, ile tokenów pierwsza odpowiedź przeczyta ponownie. Gdy limit zatrzyma sesję, hook pluginu ccam pokazuje powiadomienie, a `/ccam:switch` przekazuje rozmowę dalej. To ty decydujesz, kiedy przełączyć; nic nie śledzi twojego użycia w tle, a automatyczne wybory nie wychodzą poza jedną organizację, więc rozmowa firmowa nigdy przypadkiem nie trafi na konto prywatne. Pasek stanu podpowiada to polecenie, gdy konto osiągnie 90%.

Skrypty nie mogą na ciebie czekać, więc `ccam --fallback work -p "…"` przełącza za nie: gdy konto osiągnie limit, to samo polecenie uruchamia się ponownie na koncie z największym zapasem, a wypisywane jest tylko wyjście przebiegu, który się zakończył. Błąd, który nie jest limitem użycia, jest przekazywany bez zmian. Dla równoległych workerów `ccam pick` daje każdemu inne konto z zapasem, a oba polecenia kończą się kodem 75, gdy żadne konto nie ma już zapasu.

### Właściwe konto w każdym folderze

`ccam folder work ~/code/company` sprawia, że każde uruchomienie bez nazwy konta w tym folderze używa `work`, a pozostałe projekty nadal korzystają z konta domyślnego. To ustawienie jest w tym samym `config.toml`, który czyta każde uruchomienie, więc nie kosztuje czasu.

### Jeden plik binarny, sprawdzane pobrania

ccam to statyczny plik binarny o rozmiarze około 10 MB, bez środowiska uruchomieniowego do instalowania, dla systemów Linux, macOS i Windows. Skrypty instalacyjne (`install.sh`, `install.ps1`) i `ccam update` sprawdzają każde pobranie z sumami kontrolnymi SHA-256 wydania, a ccam nigdy sam nie sprawdza aktualizacji. Nie ma telemetrii: jedyne żądania, jakie wysyła, to zapytania do Anthropic o użycie twoich kont, równoległe, z limitem 3 sekund i 60-sekundową pamięcią podręczną, oraz pobrania, o które sam poprosisz. Z `usage_source = "claude"` nawet te zapytania przechodzą przez Claude Code, a ccam nigdy nie dotyka tokena.

## W czym inne narzędzia robią więcej

- claude-swap, clauth, cux, CAAM i CCSwitcher potrafią same przełączyć działającą sesję interaktywną, gdy konto przekroczy próg użycia. ccam przełącza przebiegi `-p` przez `--fallback`, a sesje interaktywne zostawia tobie, z powiadomieniem przy limicie i `/ccam:switch`.
- W Windows bez Trybu dewelopera ccam współdzieli foldery, ale nie pojedyncze pliki, takie jak `settings.json`; claude-swap i CCS w takim przypadku je kopiują.
- CCSwitcher, Claude Switcher i claude-swap pokazują użycie na pasku menu macOS.
- CCS, clausona, CAAM, clauth, ccswap i Claude Switcher zarządzają też Codex i innymi agentami kodującymi.
- clausona potrafi też włączać i wyłączać rozszerzenia dla każdego konta (`ccam ext` w ccam tylko je pokazuje), śledzi lokalny koszt tokenów i dostarcza pluginy Claude Code, które uruchamiają floty workerów w tmux i innych terminalach.
- CAAM, CCSwitcher i Claude Switcher instalują się przez Homebrew, a clauth i CAAM podpisują sumy kontrolne swoich wydań. ccam instaluje się własnym skryptem, pakietem `.deb` lub przez `go install` i sprawdza sumy kontrolne bez podpisu.

## Co wybrać?

Wybierz narzędzie podmieniające login, jeśli używasz jednego konta naraz albo chcesz, żeby działające sesje przełączały się automatycznie. Wybierz ccam, jeśli chcesz mieć konta obok siebie, jedną wspólną konfigurację, tę samą rozmowę na następnym koncie i nic, co działa obok Claude Code, w systemie Linux, macOS lub Windows. Jeśli nie potrzebujesz niczego współdzielić, wystarczą [zwykłe aliasy powłoki](/compare/aliases).

To, czy wolno ci używać kilku kont, zależy od warunków każdego konta. Sprawdź je, zanim zaczniesz.
