---
title: Sprawdzanie i naprawa
description: Sprawdź linki współdzielenia, uprawnienia plików, loginy i wersję Claude Code poleceniem ccam doctor, a bezpieczne poprawki wykonaj przez ccam doctor --fix.
updated: 2026-10-21
order: 9
translation: machine
sourceHash: "9be63350cb69"
---

`ccam doctor` sprawdza instalację ccam: plik binarny Claude Code, folder, login i linki współdzielenia każdego konta, współdzielone serwery MCP, endpoint użycia i polecenia skrótów. `ccam doctor --fix` naprawia to, co da się naprawić bez utraty danych, a każdy zastępowany plik najpierw przenosi do folderu kopii zapasowej. Uruchom je, gdy `ccam <name>` zatrzymuje się z powodu problemu z linkiem współdzielenia, po ręcznym przywróceniu plików albo po aktualizacji Claude Code.

## Użycie

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

Pierwszy wiersz oznacza naraz każde konto, więc widzisz, gdzie są problemy, zanim przeczytasz cokolwiek innego. Pod nim każde konto dostaje ramkę, a ustalenia, które nie należą do żadnego konta (Claude Code, zmienne środowiskowe, wtyczki, serwery MCP, skróty, reguły folderów, grupy i ustawienia), trafiają do ramki `General`. Tytuł ramki pokazuje najgorszy znacznik konta i liczbę jego problemów. Każdy wiersz w środku to jeden sprawdzony obszar, zawsze `Folder`, `Login` i `Share links` w tej kolejności, a potem pozostałe, ze znacznikiem: `✔` w porządku, `ℹ` informacja, `⚠` ostrzeżenie, `✘` błąd. Folder bez żadnych problemów pokazuje `✔ Folder  OK`.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.302 (Claude Code)                           │
 │                (/home/me/.local/bin/claude)                                │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp               │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ──────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                 │
 │ ℹ              unknown entries (Claude Code layout may have changed):      │
 │                newthing                                                    │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)      │
 │ ✔              login valid (me@example.com)                                │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json: │
 │                set to shared but is a real file, not a link                │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ───────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                          │
 │ ⚠ Login        not logged in → `ccam login work`                           │
 │ ✔ Share links  share links OK                                              │
 ╰────────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

Z `--fix` każde naprawione ustalenie kończy się `→ fixed` i nie jest już liczone jako problem:

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ────────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.302 (Claude Code)                           │
 │                (/home/me/.local/bin/claude)                                │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed       │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ───────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed         │
 │ ℹ              unknown entries (Claude Code layout may have changed):      │
 │                newthing                                                    │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →    │
 │                fixed                                                       │
 │ ✔              login valid (me@example.com)                                │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json: │
 │                set to shared but is a real file, not a link → fixed        │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ───────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                          │
 │ ⚠ Login        not logged in → `ccam login work`                           │
 │ ✔ Share links  share links OK                                              │
 ╰────────────────────────────────────────────────────────────────────────────╯
```

W terminalu ramki dopasowują się do jego szerokości, a długie komunikaty są zawijane, w kolorach samego terminala (`NO_COLOR` je wyłącza). Po przekierowaniu do potoku lub pliku wynik nie ma kolorów, a żaden komunikat nie jest zawijany ani obcinany, więc `ccam doctor | grep settings.json` znajduje cały wiersz; każda ramka ma wtedy szerokość najdłuższego wiersza. Karta Health w panelu pokazuje te same ramki.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--fix` | Naprawia to, co da się bezpiecznie naprawić (zmieniane pliki trafiają do kopii zapasowej) | wyłączone |
| `--dry-run` (globalna) | Wykonuje wszystkie sprawdzenia, ale niczego nie naprawia, nawet z `--fix` | wyłączone |

`ccam doctor` nie przyjmuje argumentów. Respektuje `--dry-run`: `ccam --dry-run doctor --fix` wypisuje te same ustalenia co `ccam doctor` i niczego nie zmienia.

Kod wyjścia to `1`, gdy co najmniej jedno ustalenie `✘` pozostaje nienaprawione, a w przeciwnym razie `0`. Ostrzeżenia i informacje nigdy nie powodują niepowodzenia. Gdy uruchomienie bez `--fix` kończy się niepowodzeniem, ccam wypisuje na końcu pusty wiersz i `` `ccam doctor --fix` repairs what it can. ``. Konfiguracja, której nie da się wczytać, też kończy się kodem `1` i błędem `ccam:`.

## Jak to działa

`internal/cli/doctor.go` wczytuje `~/.config/ccam/config.toml` i wywołuje `doctor.Run` z `internal/doctor/doctor.go`, które wykonuje sprawdzenia w stałej kolejności i zwraca ustalenia, każde z kontem, do którego należy, i obszarem, który sprawdzało. `tui.HealthReport` grupuje je w ramki. `--fix` jest przekazywane jako `Options.Fix`, ustawiane na false przy próbnym uruchomieniu.

<!-- motion:doctor-checks -->

Środkowy blok wykonuje się raz dla każdego konta. Sprawdzenia i to, co robi dla każdego z nich `--fix`:

1. **Wersja Claude Code** (`checkClaude`). Odnajduje plik binarny `claude` tak samo jak program uruchamiający i wykonuje `claude --version`. Jeśli wersja nie zaczyna się od `2.1.302`, czyli wersji, względem której zweryfikowano układ folderów ccam (`doctor.LayoutVersion`), dodaje informację `ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed`. Brak pliku binarnego lub nieudane `--version` to błąd. Bez naprawy.
2. **Dane uwierzytelniające w środowisku** (`checkAuthEnv`). Ostrzega, gdy ustawiono `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` lub `ANTHROPIC_AUTH_TOKEN`: dziedziczy je każde konto uruchomione z tej powłoki, a Claude Code przedkłada je nad własny login konta. Bez naprawy; usuń zmienną samodzielnie.
3. **Folder** (`checkProfile`). Brak folderu konta to błąd i koniec sprawdzeń tego konta. Uprawnienia z jakimikolwiek bitami grupy lub innych użytkowników to ostrzeżenie; `--fix` wykonuje `chmod 0700`. Windows nie ma bitów uprawnień, więc tam to sprawdzenie i sprawdzenie `0600` poniżej są pomijane; folder w twoim profilu użytkownika chroni jego lista kontroli dostępu.
4. **Plik logowania** (`checkCredentials`). `.credentials.json` jest sprawdzany przez `lstat`. Dowiązanie symboliczne to błąd, bo każde konto musi mieć własny token, i nie jest naprawiane. Tryb inny niż `0600` to ostrzeżenie; `--fix` wykonuje `chmod 0600`. Następnie odczytywany jest login (na macOS najpierw z pęku kluczy). Brak loginu to ostrzeżenie wskazujące `ccam login <name>`. Wygasły token odświeżania to błąd `login expired`; wygasły token dostępu to tylko informacja, bo Claude Code odnowi go przy następnym uruchomieniu. Ważny login pokazuje adres e-mail i plan (`login valid (dev@company.com, Claude Team · Acme)`). Żadne z nich nie jest naprawiane; `ccam login --expired` loguje każde konto, które tego potrzebuje.
5. **Linki współdzielenia** (`share.Check`). Porównuje każdą ścieżkę z `share.Layout` z grupami konta, używając tylko `lstat` i `readlink`. Zgłasza cztery rodzaje problemów: brakujący link, link wskazujący gdzie indziej, zwykły plik tam, gdzie powinien być link, i link pozostały po elemencie, który jest teraz lokalny. Brakujący link to ostrzeżenie, pozostałe to błędy. Ścieżki dodane do elementu w późniejszej wersji ccam (`file-history/`, `tasks/` i `paste-cache/` dla historii, `keybindings.json` dla ustawień) są tylko ostrzeżeniem, dopóki konto ma jeszcze własną kopię: `is not shared with the group yet (ccam doctor --fix moves it in)`. Uruchomienie nigdy się z ich powodu nie zatrzymuje.
6. **Nieznane wpisy** (`checkUnknownEntries`). Wymienia nazwy najwyższego poziomu w folderze konta, których ccam nie zna (`share.Known`), jako wskazówkę, że układ Claude Code się zmienił. Tylko informacja.
7. **Ścieżki wtyczek** (`checkPlugins`). Rejestry wtyczek (`installed_plugins.json`, `known_marketplaces.json`) przechowują ścieżki bezwzględne. W `plugins/` każdej grupy ścieżki prowadzące przez folder jednego konta przestają działać, gdy to konto zostanie usunięte lub przemianowane. `--fix` przepisuje je na folder grupy przez `profile.RewritePluginPaths`, atomowym zapisem pliku rejestru.
8. **Serwery MCP** (`checkMCP`). Wczytuje `mcp.json` każdej grupy raz. Plik, którego nie da się odczytać, to błąd; serwer, który nie przechodzi `mcp.Validate`, to ostrzeżenie. Bez naprawy.
9. **Format odpowiedzi użycia** (`checkUsageFormat`). Pyta endpoint użycia, z pominięciem pamięci podręcznej i ze zwykłym limitem 3 sekund, o kolejne konta, aż któreś odpowie. Jeśli odpowiedzi nie da się już sparsować, ostrzega, że `ccam usage` korzysta zastępczo z pamięci podręcznej paska stanu. To jedyne sprawdzenie, które używa sieci.
10. **Skróty** (`checkShims`). Tylko gdy któreś konto ma skrót. Brakujące linki w `~/.local/bin` to ostrzeżenie; `--fix` wywołuje `shim.Install`, które tylko tworzy brakujące linki i nigdy nie rusza innych plików. Folder bin, którego nie ma w `PATH`, to ostrzeżenie i nie jest naprawiany.
11. **Foldery** (`checkFolders`). Tylko gdy istnieją wpisy [`ccam folder`](/features/folder). Wpis, którego konto nie istnieje, to błąd, a wpis, którego folderu już nie ma, to ostrzeżenie. Żaden z nich nie jest naprawiany, bo usunięcie wpisu zmieniłoby to, które konto się tam uruchamia.
12. **Grupy i ustawienia** (`checkGroups`). Członek grupy, który nie jest kontem, konto, którego nazwa stała się w międzyczasie poleceniem ccam (uruchom je przez `ccam run <name>` albo zmień jego nazwę), oraz nieznane `theme` lub `usage_source` w `config.toml` to ostrzeżenia. Bez naprawy.
13. **Loginy** (`checkIdentities`). Odczytuje z `.claude.json` każdego folderu, na jakie konto jest zalogowany. Folder zalogowany teraz na inne konto, niż zapisał `ccam login` (na przykład po `/login` w sesji), to ostrzeżenie, podobnie jak folder zalogowany na to samo konto Claude co inny folder: oba dzielą wtedy limity użycia. Bez naprawy.

### Naprawa linków współdzielenia

Ścieżka uruchamiania od razu tworzy brakujący link, ale zatrzymuje się przy wszystkim, co mogłoby spowodować utratę danych, i odsyła tutaj:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Przy każdym problemie z linkiem współdzielenia `--fix` wywołuje `share.Build` z ustawieniami współdzielenia konta jako stanem bieżącym i docelowym, a potem `Plan.Apply`. `Build` wyprowadza każdy krok od nowa z tego, co jest na dysku, czyli tworzy ten sam plan, którego używa zmiana przez `ccam group`, więc naprawa działa według tych samych zasad co dołączanie do grupy.

<!-- motion:doctor-repair -->

Tutaj `settings.json` jest zwykłym plikiem. Ponieważ grupa ma już inny `settings.json`, konflikt przyjmuje domyślny wybór `KeepGroup`, więc kopia grupy zostaje. Plik konta trafia do `.ccam-backup/<time>/settings.json` w folderze konta, a `os.Symlink` linkuje ścieżkę do grupy. Inne przypadki:

- Jeśli grupa nie ma jeszcze kopii, plik konta przechodzi do grupy i staje się jej pierwszą zawartością.
- W przypadku folderu pliki, których grupa nie ma, przechodzą do grupy; reszta trafia do kopii zapasowej.
- Dla `history.jsonl` domyślny wybór to `merge by line`: najpierw wiersze grupy, potem twoje nowe. Poprzednia kopia grupy jest zapisywana w `.ccam-backup/<time>/group/`.
- Link do innej grupy jest zastępowany. Link do folderu poza ccam trafia do kopii zapasowej, a jego zawartość staje się początkową zawartością grupy, jeśli grupa jest pusta.
- Link pozostały po elemencie, który jest teraz lokalny, jest usuwany, a konto dostaje kopię zawartości grupy.

Poza usuwaniem linków każdy krok przenosi dane do grupy albo do folderu kopii zapasowej. Nic nie jest usuwane.

## Bezpieczeństwo i przypadki brzegowe

- Samo `ccam doctor` nigdy niczego nie zmienia. W skryptach używaj `--dry-run` razem z `--fix`, żeby mieć pewność.
- Dla uruchomionego konta `doctor --fix` tylko tworzy brakujące linki i dodaje `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest ``, bo sesja dalej zapisuje do plików, które naprawa by przeniosła.
- Jeśli naprawa współdzielenia się nie powiedzie, ccam zgłasza dla tego konta `could not repair share links: <error>` i pomija jego sprawdzenie nieznanych wpisów. Pozostałe konta są nadal sprawdzane.
- Jeśli zmiana przez `ccam group` przerwie się w połowie, nowe grupy są już zapisane, a ccam każe dokończyć poleceniem `ccam doctor --fix`.
- Dowiązany symbolicznie `.credentials.json`, wygasły login i dane uwierzytelniające w środowisku są zgłaszane, ale nigdy zmieniane.
- Przepisania rejestru wtyczek są zapisywane atomowo, ale nie są kopiowane do `.ccam-backup/`.

## Powiązane

- [Współdzielenie serwerów MCP, ustawień i historii](/docs/sharing)
- [Kopiowanie istniejącej konfiguracji](/features/migrate)
- [Limity użycia i przełączanie kont](/docs/usage-limits)
- [Polecenia](/docs/commands)
