---
title: Sprawdzanie i naprawa
description: Sprawdź linki współdzielenia, uprawnienia plików, loginy i wersję Claude Code poleceniem ccam doctor, a bezpieczne poprawki wykonaj przez ccam doctor --fix.
updated: 2026-10-26
order: 9
translation: machine
sourceHash: "282972ac5053"
---

`ccam doctor` sprawdza instalację ccam: plik binarny Claude Code, folder, login i linki współdzielenia każdego konta, współdzielone serwery MCP, endpoint użycia i polecenia skrótów. `ccam doctor --fix` naprawia to, co da się naprawić bez utraty danych, a każdy zastępowany plik najpierw przenosi do folderu kopii zapasowej. Uruchom je, gdy `ccam <name>` zatrzymuje się z powodu problemu z linkiem współdzielenia, po ręcznym przywróceniu plików albo po aktualizacji Claude Code.

## Użycie

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

Pierwszy wiersz oznacza naraz każde konto, więc widzisz, gdzie są problemy, zanim przeczytasz cokolwiek innego. Pod nim każde konto dostaje ramkę, a ustalenia, które nie należą do żadnego konta (Claude Code, zmienne środowiskowe, wtyczki, serwery MCP, skróty, reguły folderów, grupy i ustawienia, `auto_switch`), trafiają do ramki `General`. Tytuł ramki pokazuje najgorszy znacznik konta i liczbę jego problemów. Każdy wiersz w środku to jeden sprawdzony obszar, zawsze `Folder`, `Login` i `Share links` w tej kolejności, a potem pozostałe, ze znacznikiem: `✔` w porządku, `ℹ` informacja, `⚠` ostrzeżenie, `✘` błąd. Folder bez żadnych problemów pokazuje `✔ Folder  OK`.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ───────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp              │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

Z `--fix` każde naprawione ustalenie kończy się `→ fixed` i nie jest już liczone jako problem:

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ───────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed      │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ──────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed        │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →   │
 │                fixed                                                      │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link → fixed     │
 │ ℹ              settings.json differed from the group's copy, which it now │
 │                uses; this account's copy is in                            │
 │                /home/me/.local/share/ccam/profiles/personal/.ccam-        │
 │                backup/20261023-212553/settings.json                       │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯
```

W terminalu ramki dopasowują się do szerokości terminala, długie komunikaty są zawijane, a używane są kolory samego terminala (`NO_COLOR` je wyłącza). Ramka nigdy nie jest węższa niż 29 kolumn, więc w wąskim terminalu żadne słowo nie zostaje ucięte. Po przekierowaniu do potoku lub pliku wynik nie ma kolorów, a żaden komunikat nie jest zawijany ani obcinany, więc `ccam doctor | grep settings.json` znajduje cały wiersz; każda ramka ma wtedy szerokość najdłuższego wiersza. Karta Health w panelu pokazuje te same ramki.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--fix` | Naprawia to, co da się bezpiecznie naprawić (zmieniane pliki trafiają do kopii zapasowej) | wyłączone |
| `--dry-run` (globalna) | Wykonuje wszystkie sprawdzenia, ale niczego nie naprawia, nawet z `--fix` | wyłączone |

`ccam doctor` nie przyjmuje argumentów. Obsługuje `--dry-run`: `ccam --dry-run doctor --fix` wypisuje ustalenia `ccam doctor`, dodaje to, co `--fix` odłożyłby na bok, i niczego nie zmienia:

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

Kod wyjścia to `1`, gdy co najmniej jedno ustalenie `✘` pozostaje nienaprawione, a w przeciwnym razie `0`. Ostrzeżenia i informacje nigdy nie powodują niepowodzenia. Gdy uruchomienie bez `--fix` kończy się niepowodzeniem, ccam wypisuje na końcu pusty wiersz i `` `ccam doctor --fix` repairs what it can. ``. Konfiguracja, której nie da się wczytać, też kończy się kodem `1` i błędem `ccam:`.

## Jak to działa

Sprawdzenia wykonują się w stałej kolejności, a ustalenia pojawiają się w osobnej ramce dla każdego konta. ccam wczytuje `~/.config/ccam/config.toml` (`internal/cli/doctor.go`) i wykonuje sprawdzenia, a każde ustalenie wskazuje konto, do którego należy, i obszar, który sprawdzało (`doctor.Run` w `internal/doctor/doctor.go`). Potem grupuje ustalenia w ramki (`tui.HealthReport`). `--fix` jest przekazywane jako `Options.Fix`, ustawiane na false przy próbnym uruchomieniu.

<!-- motion:doctor-checks -->

Środkowy blok wykonuje się raz dla każdego konta. Sprawdzenia i to, co robi dla każdego z nich `--fix`:

1. **Wersja Claude Code** (`checkClaude`). Odnajduje plik binarny `claude` tak samo jak program uruchamiający i wykonuje `claude --version`. Brak pliku binarnego, `claude_bin`, który nie jest plikiem wykonywalnym, lub nieudane `--version` to błąd. Bez naprawy.
2. **Dane uwierzytelniające w środowisku** (`checkAuthEnv`). Ostrzega, gdy ustawiono `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` lub `ANTHROPIC_AUTH_TOKEN`: dziedziczy je każde konto uruchomione z tej powłoki, a Claude Code przedkłada je nad własny login konta. Bez naprawy; usuń zmienną samodzielnie.
3. **Folder** (`checkProfile`). Brak folderu konta to błąd i koniec sprawdzeń tego konta. Uprawnienia z jakimikolwiek bitami grupy lub innych użytkowników to ostrzeżenie; `--fix` wykonuje `chmod 0700`. Zaimportowany `~/.claude` z własnym `.claude.json` obok `~/.claude.json` to informacja: `ccam main` i zwykłe `claude` mają wtedy osobne dane logowania i serwery MCP. Windows nie ma bitów uprawnień, więc tam to sprawdzenie i sprawdzenie `0600` poniżej są pomijane; folder w twoim profilu użytkownika chroni jego lista kontroli dostępu.
4. **Plik logowania** (`checkCredentials`). `.credentials.json` jest sprawdzany przez `lstat`. Dowiązanie symboliczne to błąd, bo każde konto musi mieć własny token, i nie jest naprawiane. Tryb inny niż `0600` to ostrzeżenie; `--fix` wykonuje `chmod 0600`.

   Następnie odczytywany jest login (na macOS najpierw z pęku kluczy). Brak loginu albo plik logowania, w którym nie ma loginu, to ostrzeżenie wskazujące `ccam login <name>`. Wygasły token odświeżania albo login, który Claude Code wyczyścił po odrzuceniu przez serwer, to błąd `` login expired → `ccam login <name>` ``; wygasły token dostępu to tylko informacja, bo Claude Code odnowi go przy następnym uruchomieniu. Login przechowywany w Menedżerze poświadczeń Windows to informacja, bo ccam nie może go odczytać. Ważny login pokazuje adres e-mail, na który folder jest teraz zalogowany, i plan (`login valid (dev@company.com, Claude Team · Acme)`). Żadne z nich nie jest naprawiane; `ccam login --expired` loguje każde konto, które tego potrzebuje.
5. **Linki współdzielenia** (`share.Check`). Porównuje każdą ścieżkę z `share.Layout` z grupami konta, używając tylko `lstat` i `readlink`. Zgłasza pięć rodzajów problemów: brakujący link, link wskazujący gdzie indziej, zwykły plik tam, gdzie powinien być link, link pozostały po elemencie, który jest teraz lokalny, i poprawny link, którego folderu grupy brakuje (`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)`; uruchomienie tego nie sprawdza). Brakujący link i brakujący folder grupy to ostrzeżenia, pozostałe to błędy.

   Zgłaszane jest też kopiowanie z grupy z powrotem, które zatrzymało się w połowie: `--fix` wstawia na miejsce kompletną przygotowaną kopię, ale jeśli od tego czasu uruchomienie utworzyło ścieżkę na nowo, wskazuje folder z kompletną kopią, żebyś mógł sam to uporządkować. Ścieżki dodane do elementu w późniejszej wersji ccam (`file-history/`, `tasks/` i `paste-cache/` dla historii, `keybindings.json` dla ustawień) są tylko ostrzeżeniem, dopóki konto ma jeszcze własną kopię: `is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`. Uruchomienie nigdy się z ich powodu nie zatrzymuje.
6. **Nieznane wpisy** (`checkUnknownEntries`). Wymienia nazwy najwyższego poziomu w folderze konta, których ccam nie zna (`share.Known`), jako wskazówkę, że układ Claude Code się zmienił. Tylko informacja.
7. **Ścieżki wtyczek** (`checkPlugins`). Rejestry wtyczek (`installed_plugins.json`, `known_marketplaces.json`) przechowują ścieżki bezwzględne. W `plugins/` każdej grupy ścieżki prowadzące przez folder jednego konta przestają działać, gdy to konto zostanie usunięte lub przemianowane. `--fix` przepisuje je na folder grupy przez `profile.RewritePluginPaths`, atomowym zapisem pliku rejestru.
8. **Serwery MCP** (`checkMCP`). Wczytuje `mcp.json` każdej grupy raz. Plik, którego nie da się odczytać, to błąd; serwer, który nie przechodzi `mcp.Validate`, to ostrzeżenie. Serwery w `~/.claude.json` zwykłego `claude`, których nie ładuje żadne konto, to ostrzeżenie z poleceniem `ccam migrate --mcp`, które je kopiuje (z `--group <a|b>`, gdy MCP współdzieli kilka grup). Serwery, które ma tylko jedno konto grupy, to informacja z poleceniem `ccam mcp import`, które je udostępnia, a konto niezalogowane do serwerów MCP, do których zalogowane są inne konta, dostaje ostrzeżenie, żeby zalogować się przez `/mcp` (ccam nigdy nie kopiuje logowań). Bez naprawy.
9. **Format odpowiedzi użycia** (`checkUsageFormat`). Pyta endpoint użycia, z pominięciem pamięci podręcznej, o kolejne konta, aż któreś odpowie, w ramach jednego limitu czasu 4 sekund. Z `usage_source = "claude"` pyta zamiast tego Claude Code o jedno konto. Jeśli odpowiedzi nie da się już sparsować, ostrzega, że `ccam usage` korzysta zastępczo z pamięci podręcznej paska stanu; gdy żadne konto nie odpowiedziało, podaje to jako informację (`could not check the usage response format: …`). To jedyne sprawdzenie, które używa sieci.
10. **Skróty** (`checkShims`). Tylko gdy któreś konto ma skrót. Brakujące linki w `~/.local/bin` to ostrzeżenie, podobnie jak skróty utworzone przez ccam, które wciąż uruchamiają starszą wersję ccam, podmienioną lub przeniesioną (`shortcuts in <bin> that run another ccam (one replaced or moved): …`); `--fix` wywołuje `shim.Install`, które tworzy i odnawia tylko linki utworzone przez ccam i zapisane w `~/.local/state/ccam/shortcuts.json`, i nigdy nie rusza innych plików. Ostrzeżenie kończy się `→ fixed` tylko wtedy, gdy każdy z tych linków powstał; jeśli któregoś nie dało się utworzyć, dodawana jest przyczyna. Skrót, który uruchamia inny, wciąż zainstalowany ccam, działa, więc `ccam doctor` wymienia go tylko jako informację; `--fix` kieruje go na ten plik binarny ccam. Folder bin, którego nie ma w `PATH`, to ostrzeżenie i nie jest naprawiany.
11. **Foldery** (`checkFolders`). Tylko gdy istnieją wpisy [`ccam folder`](/features/folder). Wpis, którego konto nie istnieje, to błąd, a wpis, którego folderu już nie ma, to ostrzeżenie. Żaden z nich nie jest naprawiany, bo usunięcie wpisu zmieniłoby to, które konto się tam uruchamia.
12. **Grupy i ustawienia** (`checkGroups`). Członek grupy, który nie jest kontem, konto, którego nazwa stała się w międzyczasie poleceniem ccam (uruchom je przez `ccam run <name>` albo zmień jego nazwę), nieznane `theme` lub `usage_source` oraz wiersze `config.toml`, których ccam nie odczytuje, to ostrzeżenia. Ten ostatni przypadek to zwykle ustawienie umieszczone pod tabelą `[[group]]`, `[[profile]]` lub `[[folder]]`, przez co staje się częścią tej tabeli (`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`); polecenia zmieniające konfigurację odmawiają, dopóki nie zostanie to poprawione. Bez naprawy.
13. **Loginy** (`checkIdentities`). Odczytuje z `.claude.json` każdego folderu, na jakie konto jest zalogowany. Folder zalogowany teraz na inne konto niż to, które zapisał `ccam login`, (na przykład po `/login` w sesji), to ostrzeżenie; konto bez zapisu logowania, na przykład zaimportowane albo utworzone przez starszą wersję ccam, jest porównywane według zapisanego adresu e-mail. Folder zalogowany na to samo konto Claude co inny folder to także ostrzeżenie: oba dzielą wtedy limity użycia. Folder, który ma ten sam token logowania co inny, czyli skopiowany folder (w tym `~/.claude` zwykłego `claude`), to błąd, bo ten, który pierwszy go odnowi, wyloguje drugi; uruchom `ccam login <name>`, żeby dać mu własny. Bez naprawy.
14. **`auto_switch`** (`checkAutoSwitch`). Tylko gdy `auto_switch = true`. Informacja w ramce `General` mówi, co robi `auto_switch`, i prowadzi do [jego opisu](/features/next#auto-switch). Dla każdego planu z miejscami wśród kont (Claude Team, Claude Enterprise) osobny wiersz wymienia te konta, mówi, że `auto_switch` przenosi sesję między miejscami jednej organizacji, cytuje zasadę Anthropic dotyczącą miejsc z linkiem do centrum pomocy i zostawia wybór tobie. Wersje Claude Code starsze niż 2.1.214, które nie przekazują hookom swojego procesu (`CLAUDE_PID`), to tam ostrzeżenie, bo wtedy żadnej sesji nie da się przenieść.

    Konto bez pluginu ccam albo z wyłączonym pluginem dostaje ostrzeżenie we własnej ramce z poleceniem, które naprawia to w Claude Code: `auto_switch is on, but work does not have the ccam plugin, so its sessions are not moved. In Claude Code as work: /plugin install ccam --marketplace jhyunwoo/ccam` (albo `/plugin enable ccam@ccam`). Konta z kluczem API są pomijane. W systemie Windows jedynym wierszem jest `auto_switch is ignored on Windows; with relay = true, /ccam:switch followed by /exit still moves a conversation to another account.` Bez naprawy.

Przed kontami **zmiana grup**, która zatrzymała się w połowie, to błąd w ramce `General`: ``the last group change stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``. `--fix` kończy ją w zaplanowanej kolejności, z odpowiedziami udzielonymi przy jej konfliktach; do tego czasu nie rusza żadna inna zmiana grup, zmiana nazwy ani usunięcie. Dopóki nie zostanie dokończona, ccam tworzy dla kont, których dotyczy, tylko brakujące linki. Jeśli nie da się jej dokończyć, błąd mówi, że usunięcie tego pliku porzuca resztę zmiany.

`ccam --dry-run doctor --fix` pokazuje kolejność (`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`) i to, co część każdego konta odłożyłaby na bok. Zmiana grup, którą właśnie wykonuje inny proces ccam, to informacja, a dopóki trwa, `--fix` nie naprawia linków współdzielenia (`its share links are not repaired while another ccam changes the groups; …`). `--fix` trzyma `~/.local/state/ccam/group-change.lock` podczas naprawy i po dokończeniu zmiany ponownie wczytuje `config.toml`.

### Naprawa linków współdzielenia

Ścieżka uruchamiania od razu tworzy brakujący link, ale zatrzymuje się przy wszystkim, co mogłoby spowodować utratę danych, i odsyła tutaj:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Przy każdym problemie z linkiem współdzielenia `--fix` wywołuje `share.Build` z ustawieniami współdzielenia konta jako stanem bieżącym i docelowym, a potem `Plan.Apply`. ccam wyprowadza każdy krok od nowa z tego, co jest na dysku (`Build`), i otrzymuje ten sam plan, którego używa zmiana przez `ccam group`, więc naprawa działa według tych samych zasad co dołączanie do grupy.

<!-- motion:doctor-repair -->

Tutaj `settings.json` jest zwykłym plikiem. Ponieważ grupa ma już inny `settings.json`, konflikt przyjmuje domyślny wybór `keep group` (`KeepGroup` na diagramie), więc kopia grupy zostaje. Plik konta trafia do `.ccam-backup/<time>/settings.json` w folderze konta, a ccam linkuje ścieżkę do grupy (`os.Symlink`). `--fix` o tym informuje: `settings.json differed from the group's copy, which it now uses; this account's copy is in …`. Inne przypadki:

- Jeśli grupa nie ma jeszcze kopii, plik konta przechodzi do grupy i staje się jej pierwszą zawartością.
- W przypadku folderu pliki, których grupa nie ma, przechodzą do grupy. Pliki identyczne z kopią grupy są usuwane, a te, które się różnią, trafiają do kopii zapasowej wskazanej w notatce (`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`).
- Brakujący folder grupy jest tworzony najpierw.
- Dla `history.jsonl` domyślny wybór to `merge by line`: najpierw wiersze grupy, potem twoje nowe. Poprzednia kopia grupy jest zapisywana w `.ccam-backup/<time>/group/`.
- Link do innej grupy jest zastępowany. Link do folderu poza ccam trafia do kopii zapasowej, a jego zawartość staje się początkową zawartością grupy, jeśli grupa jest pusta.
- Link pozostały po elemencie, który jest teraz lokalny, jest usuwany, a konto dostaje kopię zawartości grupy.

Poza usuwaniem linków i usuwaniem plików identycznych z kopią grupy, każdy krok przenosi dane do grupy albo do folderu kopii zapasowej. Nic nie ginie.

## Bezpieczeństwo i przypadki brzegowe

- Samo `ccam doctor` nigdy niczego nie zmienia. W skryptach używaj `--dry-run` razem z `--fix`, żeby mieć pewność.
- Dla uruchomionego konta, niezależnie od tego, jak uruchomiono Claude Code, `doctor --fix` tylko tworzy brakujące linki i brakujące foldery grup, co niczego nie traci. Gdy zostaje coś jeszcze, dodaje `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest ``, bo sesja dalej zapisuje do plików, które naprawa by przeniosła.
- Jeśli naprawa współdzielenia się nie powiedzie, ccam zgłasza dla tego konta `could not repair share links: <error>` i pomija jego sprawdzenie nieznanych wpisów. Pozostałe konta są nadal sprawdzane.
- Jeśli zmiana przez `ccam group` przerwie się w połowie, nowe grupy są już zapisane, a zmiana jest odnotowana w `~/.local/state/ccam/group-change.json`; ccam każe dokończyć poleceniem `ccam doctor --fix`, które wykonuje jej resztę, łącznie z krokami MCP.
- Dowiązany symbolicznie `.credentials.json`, wygasły login i dane uwierzytelniające w środowisku są zgłaszane, ale nigdy zmieniane.
- Przepisania rejestru wtyczek są zapisywane atomowo, ale nie są kopiowane do `.ccam-backup/`.

## Powiązane

- [Współdzielenie serwerów MCP, ustawień i historii](/docs/sharing)
- [Kopiowanie istniejącej konfiguracji](/features/migrate)
- [Limity użycia i przełączanie kont](/docs/usage-limits)
- [Opis poleceń](/docs/commands)
