---
title: Zarządzanie kontami
description: Dodawaj, loguj, wyświetlaj, zmieniaj nazwy, usuwaj i importuj konta Claude Code poleceniami ccam add, login, list, rename, remove i import.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 99c931574262
---

Każde konto ccam to osobny folder konfiguracji Claude Code z własnym loginem. Te polecenia tworzą konta, logują do nich, wyświetlają je, zmieniają ich nazwy i je usuwają, a także rejestrują foldery, które już masz.

## Użycie

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Gdy polecenie przyjmuje `<name|number>`, możesz podać nazwę konta albo jego numer slotu, czyli kolumnę `#` w wyniku `ccam list`.

`ccam add` uruchomione w terminalu bez opcji otwiera kreator. Pyta o nazwę, potem pokazuje listę współdzielonych elementów do zaznaczenia (zaznaczone jest wszystko oprócz historii rozmów), następnie pyta `Create a shortcut command for work? [y/N]`, a po odpowiedzi `y` także `Shortcut name [ccwork]:`. Każda z opcji `--share`, `--shortcut` lub `--no-login`, albo potok zamiast terminala, pomija kreator; wtedy nazwa jest wymagana, a skrót powstaje tylko wtedy, gdy podasz go w `--shortcut`.

`ccam list` pokazuje slot, nazwę, skrót, adres e-mail, współdzielone elementy i ostatnie użycie, z `*` przed domyślnym kontem:

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

Z `--json` każde konto jest obiektem z polami `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (pomijanym, gdy konto nigdy nie było uruchomione) i `share`. `share` przypisuje każdy element (`mcp`, `settings`, `memory`, `extensions`, `history`) do nazwy jego grupy albo do `"local"`.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `ccam add --no-login` | Tworzy konto bez logowania; zalogujesz się później przez `ccam login <name>` | wyłączone |
| `ccam add --share <item=group>` | Współdzieli element z grupą albo, z `item=local`, zostawia go prywatnym. Można powtarzać | `mcp`, `settings`, `memory`, `extensions` w `default`; `history` lokalnie |
| `ccam add --shortcut <command>` | Tworzy też to polecenie skrótu w `~/.local/bin` | brak |
| `ccam list --json` | Wypisuje konta jako JSON | wyłączone |
| `ccam remove -y`, `--yes` | Usuwa bez pytania | wyłączone |
| `ccam import --name <name>` | Nazwa konta dla podanego folderu | wyprowadzona z nazwy folderu |
| `ccam import -y`, `--yes` | Importuje każdy znaleziony folder bez pytania | wyłączone |
| `--dry-run` | Respektowana przez `ccam list` i `ccam completion`, które niczego nie zmieniają. Odrzucana przez `add`, `login`, `logout`, `use`, `rename`, `remove` i `import` komunikatem `` --dry-run is not supported by `ccam add`; nothing was changed `` (z nazwą polecenia) | wyłączone |

`ccam login`, `ccam logout`, `ccam use` i `ccam rename` nie mają własnych opcji. `ccam completion` przyjmuje `bash` lub `zsh`; wszystko inne kończy się błędem `only bash and zsh are supported`.

Nazwy kont mają od 1 do 32 znaków spośród `a-z`, `0-9`, `-` i `_` i zaczynają się literą lub cyfrą. Nazwa nie może składać się z samych cyfr, bo `ccam 2` oznacza slot 2. Nie może też być poleceniem ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`), więc `ccam <name>` nigdy nie jest niejednoznaczne. Ta lista to `config.Reserved` w `internal/config`.

## Jak to działa

`ccam add work --shortcut ccwork` wykonuje te kroki w `addProfile` (`internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. `config.ValidateName` sprawdza nazwę, a `shim.Check` sprawdza skrót, zanim cokolwiek zostanie utworzone.
2. `profile.Add` tworzy `~/.local/share/ccam/profiles/work/` z uprawnieniami `0700` i nadaje kontu najniższy wolny numer slotu (`Config.NextSlot`), więc numer usuniętego konta jest używany ponownie. Pierwsze konto staje się domyślnym.
3. `applyShare` linkuje nowy, pusty folder do jego grup współdzielenia przez `share.Build`, a potem `Apply`, na przykład `settings.json` do `groups/default/settings.json`. Jeśli to się nie uda, nowy folder jest usuwany, bo pozostawiony folder blokowałby każdą kolejną próbę błędem `folder already exists`.
4. Konto jest zapisywane w `~/.config/ccam/config.toml` przez `fsutil.AtomicWrite`: plik tymczasowy, `fsync`, potem `rename`, więc awaria nigdy nie zostawia połowy konfiguracji.
5. `runLogin` uruchamia `claude auth login` z `CLAUDE_CONFIG_DIR` ustawionym na folder konta. Jeśli logowanie się nie powiedzie, konto zostaje, a ccam wypisuje błąd i każe później uruchomić `ccam login work`.
6. `profile.ReadEmail` odczytuje `oauthAccount.emailAddress` z `profiles/work/.claude.json`. Ten plik może mieć kilka megabajtów, więc `TopLevelValue` wyszukuje jeden klucz zamiast parsować cały plik.
7. `shim.Install` linkuje `~/.local/bin/ccwork` do pliku binarnego ccam, który przy starcie szuka `argv[0]` w `config.toml`.

Logowanie może trwać minuty, a w tym czasie inny ccam może zapisać konfigurację. Dlatego polecenia, które czekają, zapisują przez `env.update`, które ponownie wczytuje `config.toml`, stosuje jedną zmianę (tutaj adres e-mail) i zapisuje. Z tego samego powodu `ccam remove` wczytuje konfigurację ponownie po pytaniu.

`ccam login` wykonuje kroki 5 i 6 dla istniejącego konta. `ccam logout` uruchamia `claude auth logout` na koncie i czyści zapisany adres e-mail. `ccam use` zmienia tylko `default` w `config.toml`.

### Logowanie na zdalnej maszynie

Przez SSH część `claude auth login` działająca w przeglądarce może zawieść. Claude Code czeka na zalogowanie na serwerze zwrotnym pod `http://localhost:<port>/callback` na serwerze, ale przeglądarka, która kończy logowanie, działa na twoim komputerze, gdzie ten adres prowadzi donikąd, więc przeglądarka pokazuje błąd taki jak `ERR_CONNECTION_REFUSED`. Adres nadal zawiera jednorazowy kod, więc ccam może go dostarczyć:

<!-- motion:accounts-relay -->

`runLogin` nie oddaje twojego terminala Claude Code. Podłącza standardowe wejście Claude Code do potoku, a `relayLogin` kopiuje do niego każdy wpisany przez ciebie wiersz, więc kod wklejony przy monicie `Paste code here` w Claude Code nadal dociera. Gdy wiersz jest adresem zwrotnym, `callbackURL` akceptuje go tylko wtedy, gdy host to `localhost`, `127.0.0.1` lub `[::1]`, podano port, ścieżka to `/callback` i obecne są zarówno `code`, jak i `state`; `http://` można pominąć. `deliverCallback` wysyła wtedy żądanie pod ten adres na serwerze z limitem 10 sekund i bez podążania za przekierowaniami, tak jak zrobiłaby to przeglądarka. Żądanie trafia do serwera zwrotnego oczekującego `claude auth login` (linia przerywana na diagramie). Claude Code otrzymuje kod i zapisuje login, na Linux w `.credentials.json` w folderze konta. Wklejony wiersz nigdy nie może sprawić, że ccam połączy się z innym hostem.

Przed rozpoczęciem logowania ccam wypisuje:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

Po powodzeniu wypisuje `→ Passed the sign-in to Claude Code on this machine.`

Gdy logowanie czeka na dane w terminalu, ccam wyłącza echo i tryb liniowy terminala i dla każdego wpisanego znaku pokazuje `*`, więc wklejony kod lub adres zwrotny nie pojawia się na ekranie. `readLine` sam obsługuje edycję wiersza: `Enter` wysyła wiersz, `Backspace` usuwa jeden znak, a `Ctrl+D` w pustym wierszu kończy wprowadzanie. Terminal jest przywracany po zakończeniu logowania, a także po `Ctrl+C`, `SIGTERM` lub `SIGHUP`, które nadal kończą ccam. Gdy standardowe wejście i standardowe wyjście błędów nie są terminalem, wiersze są kopiowane na bieżąco.

### Zmiana nazwy, usuwanie i import

`ccam rename` i `ccam remove` działają na folderze konta i na tym, co ccam przechowuje pod nazwą konta w innych miejscach. `nameState` wymienia te ścieżki: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` oraz pamięci podręczne użycia `~/.local/state/ccam/usage/<name>.json` i `<name>.statusline.json`. `ccam remove work` obsługuje je w tej kolejności:

<!-- motion:accounts-remove -->

`ccam remove` najpierw odmawia, jeśli konto jest uruchomione, a potem pyta. Jeśli konto współdzieli rozszerzenia, `RewritePluginPaths` przepisuje `installed_plugins.json` i `known_marketplaces.json` grupy, bo Claude Code przechowuje bezwzględne ścieżki wtyczek, które mogą prowadzić przez folder tego konta. `profile.Trash` przenosi folder do `~/.local/share/ccam/trash/<name>-<time>/`, konto znika z `config.toml`, ścieżki `nameState` są usuwane, żeby późniejsze konto o tej samej nazwie ich nie odziedziczyło, a link skrótu jest usuwany. Aby odzyskać konto, przenieś folder z kosza z powrotem.

`ccam rename` zmienia nazwę folderu, przepisuje ścieżki wtyczek w jego `plugins/` ze starego folderu na nowy, przenosi status domyślnego konta i przenosi ścieżki `nameState` pod nową nazwę. Polecenie skrótu zachowuje swoją nazwę. Jeśli przepisanie ścieżek wtyczek się nie powiedzie, folder wraca do starej nazwy, bo konfiguracja nie zostaje zapisana.

`ccam import` rejestruje folder tam, gdzie jest, jako `dir` w `config.toml`. Nigdy go nie przenosi ani nie kopiuje: skopiowany token odświeżania sprawiłby, że kopie nawzajem by się wylogowywały. Bez argumentu z folderem szuka niezarejestrowanych folderów `~/.claude` i `~/.claude-*` i pyta o każdy; nazwy pochodzą od folderu (`~/.claude` to `main`, `~/.claude-work` to `work`, `~/.claude-2` to `acct2`), a folder, z którego nie da się wyprowadzić nazwy, jest wypisywany razem z poleceniem `--name` do użycia. Bez `CLAUDE_CONFIG_DIR` Claude Code trzyma `.claude.json` w folderze domowym obok `~/.claude`, więc gdy importujesz `~/.claude`, a nie ma on własnego `.claude.json`, ccam tworzy link `~/.claude/.claude.json` do `~/.claude.json`. Zwykłe `claude` i `ccam main` widzą wtedy to samo konto. Zaimportowane konta niczego nie współdzielą, dopóki nie uruchomisz `ccam share`, i nie powstaje dla nich skrót.

`ccam completion bash|zsh` wypisuje skrypt uzupełniania, który uzupełnia też nazwy kont.

## Bezpieczeństwo i przypadki brzegowe

- `ccam rename` i `ccam remove` odmawiają działania, gdy konto jest uruchomione (`refuseRunning`), z komunikatem `work is running; quit its Claude Code session first`, bo sesja dalej zapisuje do folderu.
- `ccam remove` na zaimportowanym koncie tylko je wyrejestrowuje: `✓ main unregistered (folder ... left in place)`. Twój `~/.claude` nie trafia do kosza.
- `ccam add` odrzuca nazwę, której folder już istnieje w `profiles/`, i wskazuje `ccam import`.
- `ccam import` odrzuca folder, który jest już zarejestrowany, ścieżkę, która nie jest folderem, i zajętą nazwę.
- `ccam logout` i `ccam login` nie sprawdzają, czy konto jest uruchomione.
- Jeśli przekazanie nie może dotrzeć do Claude Code, na przykład po zakończeniu logowania, ccam wypisuje `ccam: could not pass the sign-in to Claude Code:` z przyczyną. Błąd HTTP zwykle oznacza, że kod wygasł; zacznij logowanie od nowa.

## Powiązane

- [Uruchamianie Claude Code](/features/launch)
- [Polecenia skrótów](/features/shortcuts)
- [Współdzielenie serwerów MCP, ustawień i historii](/docs/sharing)
- [Kilka kont Claude Code](/guides/multiple-claude-code-accounts)
