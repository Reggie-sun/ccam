---
title: Współdzielenie między kontami
description: Jak ccam share linkuje ustawienia, CLAUDE.md, rozszerzenia i historię do folderu grupy, z planem, kopiami zapasowymi i wyborem przy konfliktach, oraz jak off kopiuje je z powrotem.
updated: 2026-10-19
order: 4
translation: machine
sourceHash: b30b7cdb2b15
---

Każde konto ma własny folder, ale większość jego zawartości nie musi być osobna. `ccam share` włącza lub wyłącza współdzielenie dla jednego konta i jednego elementu naraz: plik lub folder konta jest przenoszony do folderu grupy i zastępowany dowiązaniem symbolicznym, więc każde konto w grupie czyta i zapisuje tę samą kopię. Używaj go po `ccam add`, żeby zmienić, co konto współdzieli, albo żeby umieścić kilka kont w osobnej grupie.

## Użycie

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

Bez opcji elementów `ccam share` wypisuje bieżący stan:

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

Z opcjami ccam wypisuje plan, pyta o każdy plik, który różni się od kopii grupy, i pyta `Proceed? [y/N]`, zanim cokolwiek zmieni. Nowe konto współdzieli wszystko oprócz historii rozmów; [Współdzielenie serwerów MCP, ustawień i historii](/docs/sharing) opisuje każdy element z perspektywy użytkownika.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | serwery MCP | bez zmian |
| `--settings <on\|off\|group>` | `settings.json` | bez zmian |
| `--memory <on\|off\|group>` | globalny `CLAUDE.md` | bez zmian |
| `--extensions <on\|off\|group>` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | bez zmian |
| `--history <on\|off\|group>` | `projects/` i `history.jsonl` | bez zmian |
| `--fresh` | przy wyłączaniu współdzielenia zaczyna od pustego stanu zamiast kopiować zawartość grupy | wyłączone |
| `--prefer <group\|mine\|merge>` | jak rozstrzygać konflikty | pytanie; z `--yes` domyślny wybór każdego elementu |
| `-y`, `--yes` | działa bez pytania; konflikty są rozstrzygane według `--prefer` | wyłączone |
| `-h`, `--help` | pomoc dla `share` | |
| `--dry-run` (globalna) | respektowana: pyta o konflikty, wypisuje plan, niczego nie zmienia | wyłączone |

Wartość `on` oznacza grupę `default`, `off` oznacza lokalnie (bez współdzielenia), a każda inna wartość to nazwa grupy. Nazwy grup stają się nazwami folderów w `~/.local/share/ccam/groups/`, więc `ccam share side --history ../x` kończy się błędem `invalid group name: "../x"`. Konto podaje się nazwą lub numerem, a zmieniają się tylko wymienione elementy.

## Jak to działa

Pliki, których to dotyczy:

- `~/.config/ccam/config.toml`: tabela `[profile.share]` każdego konta, jedna wartość na element.
- `~/.local/share/ccam/profiles/<name>/`: `CLAUDE_CONFIG_DIR` konta.
- `~/.local/share/ccam/groups/<group>/`: współdzielone kopie.

1. `internal/cli/share.go` odczytuje konfigurację, stosuje każdą wartość `--<item>` do kopii ustawień konta przez `config.Share.Set` i odmawia, jeśli konto jest uruchomione.
2. `share.Build` porównuje bieżące ustawienia z docelowymi i wypisuje kroki. Tylko odczytuje dysk: `os.Lstat` i `os.Readlink` dla każdej ścieżki z `share.Layout`. Plan to lista operacji: przeniesienie do pustego miejsca w grupie (przejęcie), przeniesienie plików, których grupa nie ma (scalenie folderu), zapisanie rozstrzygniętego konfliktu w grupie, kopia zapasowa, usunięcie linku, utworzenie linku, skopiowanie z powrotem, przepisanie ścieżek wtyczek i scalenie rejestru wtyczek.
3. Konflikty są rozstrzygane, przez pytanie lub przez `--prefer`.
4. `Plan.Describe` wypisuje plan, a ccam pyta o zgodę.
5. Docelowe ustawienia są zapisywane w `config.toml`, zanim cokolwiek zostanie przeniesione. Jeśli `Plan.Apply` przerwie się w połowie, `ccam doctor --fix` i ponowne uruchomienie budują plan względem zapisanego celu i kończą zmianę.
6. `Plan.Apply` wykonuje kroki po kolei. Każdy krok przenosi dane do grupy albo do folderu kopii zapasowej; nic nie jest usuwane.

`share.Layout` to lista dozwolonych ścieżek, które obejmuje każdy element. Nic spoza niej nigdy nie trafi do grupy. `share.Never` wymienia `.credentials.json` i `.claude.json`, a `Build` zatrzymuje się z błędem wewnętrznym, jeśli któryś z nich kiedykolwiek pojawi się w układzie. Element `mcp` nie ma ścieżek: lista serwerów jest w `.claude.json`, więc [współdzielenie MCP](/features/mcp) działa zamiast tego przez `mcp.json` grupy, a `--mcp` zmienia tylko ustawienie.

Oto co dzieje się z jednym folderem, `projects/`, gdy `side` włącza historię, a grupa ma już zawartość:

<!-- motion:share-on -->

Jeśli grupa nie ma jeszcze niczego pod tą ścieżką, konto jest pierwszym, które ją współdzieli, a jego folder zostaje tam przeniesiony bez zmian, bez kopii zapasowej. W przeciwnym razie `mergeDir` przenosi do grupy każdy plik, którego grupa nie ma, a pliki istniejące po obu stronach zostają w folderze konta, więc zachowywana jest wersja grupy. Folder konta trafia następnie do `.ccam-backup/<time>/` w folderze konta, a `os.Symlink` wstawia na jego miejsce link. Pliki rozmów mają nazwy według identyfikatora sesji, więc historie z kilku kont łączą się bez konfliktów.

Przenoszenie używa `rename(2)`, a między systemami plików przechodzi na kopiowanie i usuwanie. Względny link wewnątrz przenoszonego drzewa, który wskazywał poza nie, jest przepisywany na link bezwzględny (`fsutil.PortableLink`), więc nadal działa z nowego miejsca. Rejestry `plugins/` (`installed_plugins.json`, `known_marketplaces.json`) przechowują ścieżki bezwzględne prowadzące przez folder konta; ccam przestawia je na folder grupy, a gdy grupa ma już rejestr, scala z nim wpisy konta.

Wyłączenie elementu odwraca link:

<!-- motion:share-off -->

ccam usuwa link tylko wtedy, gdy wskazuje do folderu grup, a potem kopiuje bieżącą zawartość grupy do konta przez `CopyTree`. Z `--fresh` kopiowanie jest pomijane, a konto zaczyna bez niczego pod tą ścieżką. Grupa zachowuje swoją zawartość dla pozostałych kont. Przeniesienie elementu z jednej grupy do innej usuwa link i tworzy link do nowej grupy, zostawiając starą grupę bez zmian; plan informuje o tym w wierszu `Note:`.

Pojedynczy plik, który istnieje po obu stronach z różną zawartością, to konflikt:

<!-- motion:share-conflict -->

`share.NewConflict` wymienia wybory pasujące do pliku. `history.jsonl` oferuje `merge by line`, `keep group` i `use mine`, ze scalaniem na pierwszym miejscu. Inne pliki `.json` oferują `keep group`, `use mine` i `merge JSON`. Wszystkie pozostałe oferują `keep group` i `use mine`. Pierwszy wybór jest domyślny:

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON` (`MergeJSONBytes`) scala obiekty rekurencyjnie, łączy tablice bez duplikatów i zachowuje wartość grupy, gdy obie strony ustawiają ten sam klucz inaczej, więc pozostałe konta nie widzą zmiany w działaniu. `merge by line` zachowuje wiersze grupy, a potem dodaje twoje wiersze, których grupa nie ma. Każdy wybór oprócz `keep group` najpierw zapisuje kopię grupy w `.ccam-backup/<time>/group/` w folderze konta, a potem zapisuje wynik przez `fsutil.AtomicWrite`. Twój własny plik zawsze trafia do kopii zapasowej.

`--prefer group`, `mine` lub `merge` rozstrzyga każdy konflikt bez pytania; `merge` oznacza `merge JSON` lub `merge by line`, w zależności od tego, co oferuje plik, a plik, który nie oferuje żadnego z nich, zachowuje wybór domyślny. Samo `--yes` przyjmuje domyślny wybór każdego pliku.

Przy każdym uruchomieniu `share.Check` porównuje każdą ścieżkę z układu z ustawieniem, używając tylko `lstat` i `readlink`, dzięki czemu jest wystarczająco tanie dla ścieżki uruchamiania. Brakujący link jest tworzony od razu, bo nie ma nic do stracenia. Zwykły plik tam, gdzie powinien być link, link w złe miejsce albo nieaktualny link zatrzymują uruchomienie:

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix` uruchamia `share.Build` z zapisanymi ustawieniami po obu stronach, co wyprowadza każdy krok z tego, co jest na dysku, zachowuje zawartość grupy przy konfliktach i robi kopię zapasową wszystkiego, co zastępuje.

## Bezpieczeństwo i przypadki brzegowe

- **Działające konta.** ccam odmawia zmiany współdzielenia, gdy Claude Code działa jako to konto (`account work is running; quit it first, then change its sharing`), bo sesja dalej zapisuje do przenoszonych plików, a wiersze dopisane w tym czasie zostałyby utracone. `--dry-run` nadal wypisuje plan.
- **Kopie zapasowe.** Każdy plan dostaje własny folder `.ccam-backup/<time>/` w folderze konta, z numerowanym przyrostkiem, gdy dwa plany zaczynają się w tej samej sekundzie. ccam nigdy ich nie usuwa.
- **Częściowe niepowodzenia.** Jeśli krok się nie powiedzie, ccam wypisuje błąd i miejsce, w którym są oryginały, i każe dokończyć poleceniem `ccam doctor --fix`. Cel jest już zapisany, więc naprawa kontynuuje w jego kierunku.
- **Twoje własne linki.** Jeśli ścieżka jest już linkiem, który sam utworzyłeś, na przykład do repozytorium dotfiles, a grupa jest pusta, ccam kopiuje cel linku do grupy, a sam link przenosi do kopii zapasowej. Cel nie jest ruszany.
- **Pliki logowania.** `.credentials.json` i `.claude.json` nigdy nie są linkowane, przenoszone ani kopiowane do grupy.
- **Dwa konta w jednym folderze.** Konta, które współdzielą historię, zapisują do tych samych plików. Gdy drugie konto z tej samej grupy historii startuje w folderze, w którym jedno już działa, ccam ostrzega i pyta, zanim będzie kontynuować.

## Powiązane

- [Współdzielenie serwerów MCP, ustawień i historii](/docs/sharing)
- [Współdzielone serwery MCP](/features/mcp)
- [Kopiowanie istniejącej konfiguracji](/features/migrate)
- [Wszystkie polecenia](/docs/commands)
