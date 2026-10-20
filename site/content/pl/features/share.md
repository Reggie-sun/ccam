---
title: Współdzielenie przez grupy
description: Grupa ccam współdzieli wybrane elementy między swoimi kontami przez jeden folder; ccam group przenosi dane z planem, kopiami i wyborem przy konfliktach.
updated: 2026-10-20
order: 4
translation: machine
sourceHash: "0c0338362ba2"
---

Każde konto ma własny folder, ale większość jego zawartości nie musi być osobna. **Grupa** określa, co jest współdzielone i przez kogo: każde konto w grupie linkuje te elementy do jednego folderu, `~/.local/share/ccam/groups/<group>/`, więc wszystkie czytają i zapisują tę samą kopię. Tokeny logowania nigdy nie są współdzielone.

Konto może należeć do kilku grup, ale każdy element pochodzi tylko z jednej grupy. Typowa konfiguracja to jedna grupa na ustawienia i rozszerzenia oraz druga na historię rozmów między dwoma kontami:

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

Pierwsze dodane konto tworzy `default`, a kolejne konta do niej dołączają, chyba że wskażesz inaczej.

## Użycie

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # every account in default moves together
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # work's CLAUDE.md goes in first
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # leave the group that shares them now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

Każde polecenie, które przenosi dane, wypisuje plan, pyta o każdy plik, który różni się od kopii grupy, i pyta `Proceed? [y/N]`, zanim cokolwiek zmieni:

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261020-153012
  link     projects, history.jsonl: profiles/work → groups/default
Proceed? [y/N]
```

Te same zmiany są dostępne na karcie **Groups** [panelu](/features/dashboard), a `s` na koncie wybiera jego grupy.

## Elementy

| Element | Co obejmuje |
| --- | --- |
| `mcp` | serwery MCP, przez `mcp.json` grupy ([Współdzielone serwery MCP](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | globalny `CLAUDE.md` |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` i `history.jsonl`, żeby inne konto mogło użyć `--continue` |

Element, którego nie współdzieli żadna grupa konta, zostaje we własnym folderze konta.

## Polecenia i opcje

| Polecenie | Działanie |
| --- | --- |
| `ccam group` lub `ccam group list` | wyświetla grupy |
| `ccam group show <group>` | co grupa współdzieli, jej konta, jej folder i serwery MCP |
| `ccam group add <group> --items <items> [--members <accounts>]` | tworzy grupę; elementy jej kont są do niej przenoszone |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | zmienia, co grupa współdzieli, dla każdego jej konta |
| `ccam group join <group> <account>...` | dodaje konta do grupy |
| `ccam group leave <group> <account>...` | wyjmuje konta z grupy; każde zachowuje kopię |
| `ccam group rename <group> <new-name>` | zmienia nazwę tylko w konfiguracji |
| `ccam group remove <group> [--leave]` | usuwa grupę |

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: konto, którego kopia trafia do grupy jako pierwsza | konto domyślne, potem według numeru |
| `--replace` | `join`: najpierw opuszcza grupy, które współdzielą te same elementy | wyłączone |
| `--fresh` | konta, które przestają współdzielić element, zaczynają bez niego zamiast z kopią | wyłączone |
| `--leave` | `remove`: najpierw wyjmuje konta grupy | wymagane, gdy grupa ma konta |
| `--prefer <group\|mine\|merge>` | jak rozstrzygać konflikty | pytanie; z `--yes` domyślny wybór każdego pliku |
| `-y`, `--yes` | działa bez pytania; konflikty są rozstrzygane według `--prefer` | wyłączone |
| `--dry-run` (globalna) | wypisuje plan, niczego nie zmienia | wyłączone |

Konta podaje się nazwą lub numerem, a elementy jako listę rozdzieloną przecinkami. `ccam add` przyjmuje `--group <group>` (można powtarzać) lub `--no-group`.

## Jak to działa

Pliki, których to dotyczy:

- `~/.config/ccam/config.toml`: wpisy `[[group]]`: `name`, `items`, `members` oraz `dir`, gdy nazwa folderu jest inna.
- `~/.local/share/ccam/profiles/<name>/`: `CLAUDE_CONFIG_DIR` konta.
- `~/.local/share/ccam/groups/<folder>/`: współdzielone kopie.

Przy odczycie `config.toml` grupy każdego konta są rozwiązywane do jednego folderu na element. Przy pięciu kontach zajmuje to około 25 µs, więc uruchomienie nadal czyta tylko `config.toml`. Konfiguracja, w której konto dostawałoby element z dwóch grup, jest błędem i nic się nie uruchomi, dopóki nie zostanie poprawiona.

Zmiana działa na kopii konfiguracji:

1. Polecenie zmienia kopię (`AddGroup`, `SetGroupItems`, `JoinGroup` …), co sprawdza zasady.
2. `internal/group` porównuje stan stary i nowy i dla każdego konta, którego elementy się zmieniają, buduje plan przez `share.Build`. Plany są budowane jeden po drugim i każdy widzi, co wcześniejsze umieszczą w folderach grup, więc tylko pierwsze konto wnosi swoją kopię. Pierwsza kopia pochodzi z folderu grupy, jeśli już jakąś ma, potem z `--from`, potem z konta domyślnego, a potem z pozostałych według numeru.
3. ccam odmawia, jeśli któreś z tych kont jest uruchomione, a w Windows bez Trybu dewelopera także wtedy, gdy plik wymagałby linku.
4. Konflikty są rozstrzygane, przez pytanie lub przez `--prefer`, plan jest wypisywany, a ccam pyta o zgodę.
5. Nowa konfiguracja jest zapisywana, zanim cokolwiek zostanie przeniesione. Jeśli krok przerwie się w połowie, `ccam doctor --fix` i ponowne uruchomienie budują plan względem zapisanych grup i kończą zmianę.
6. Plany są wykonywane po kolei. Każdy krok przenosi dane do grupy albo do folderu kopii zapasowej; nic nie jest usuwane.

`share.Layout` to lista dozwolonych ścieżek, które obejmuje każdy element. Nic spoza niej nigdy nie trafi do grupy. `share.Never` wymienia `.credentials.json` i `.claude.json`, a `Build` zatrzymuje się z błędem wewnętrznym, jeśli któryś z nich kiedykolwiek pojawi się w układzie. Element `mcp` nie ma ścieżek: lista serwerów jest w `.claude.json`, więc [współdzielenie MCP](/features/mcp) działa przez `mcp.json` grupy.

Oto co dzieje się z jednym folderem, `projects/`, gdy `side` zaczyna współdzielić historię, a grupa ma już zawartość:

<!-- motion:share-on -->

Jeśli grupa nie ma jeszcze niczego pod tą ścieżką, konto jest pierwszym, które ją współdzieli, a jego folder zostaje tam przeniesiony bez zmian. W przeciwnym razie każdy plik, którego grupa nie ma, jest przenoszony do grupy, a pliki istniejące po obu stronach zostają w folderze konta, więc zachowywana jest wersja grupy. Folder konta trafia następnie do `.ccam-backup/<time>/` w folderze konta, a na jego miejsce wstawiany jest link. Pliki rozmów mają nazwy według identyfikatora sesji, więc historie z kilku kont łączą się bez konfliktów. Rejestry `plugins/` przechowują ścieżki bezwzględne prowadzące przez folder konta; ccam przestawia je na folder grupy i scala wpisy konta z rejestrem, który grupa już ma.

Opuszczenie grupy odwraca link:

<!-- motion:share-off -->

ccam usuwa link tylko wtedy, gdy wskazuje do folderu grup, a potem kopiuje bieżącą zawartość grupy do konta. Z `--fresh` kopiowanie jest pomijane. Grupa zachowuje swoją zawartość dla pozostałych kont. Konto, które przenosi element do grupy, której folder jest jeszcze pusty, wnosi kopię tego, co widziało, więc zachowuje swoje pliki; serwery MCP są przenoszone tak samo. Konto, które przestaje współdzielić MCP, traci serwery, które ccam skopiował do jego `.claude.json` (kopia zapasowa zostaje), bo pochodziły z grupy.

Pojedynczy plik, który istnieje po obu stronach z różną zawartością, to konflikt:

<!-- motion:share-conflict -->

`history.jsonl` oferuje `merge by line`, `keep group` i `use mine`, ze scalaniem na pierwszym miejscu. Inne pliki `.json` oferują `keep group`, `use mine` i `merge JSON`. Wszystkie pozostałe oferują `keep group` i `use mine`. `merge JSON` scala obiekty rekurencyjnie, łączy tablice bez duplikatów i zachowuje wartość grupy, gdy obie strony ustawiają ten sam klucz inaczej, więc pozostałe konta nie widzą żadnej zmiany. Każdy wybór oprócz `keep group` najpierw zapisuje kopię grupy w `.ccam-backup/<time>/group/` konta. Twój własny plik zawsze trafia do kopii zapasowej.

Przy każdym uruchomieniu `share.Check` porównuje każdą ścieżkę z układu z grupami konta, używając tylko `lstat` i `readlink`. Brakujący link jest tworzony od razu. Zwykły plik tam, gdzie powinien być link, link w złe miejsce albo nieaktualny link zatrzymują uruchomienie i odsyłają do `ccam doctor --fix`.

### Zmiana nazwy, usuwanie i stary format

- `ccam group rename` zmienia tylko konfigurację. Folder zachowuje swoją nazwę, a grupa zapisuje ją jako `dir`, więc nic się nie przenosi i nic nie może zginąć.
- `ccam group remove` odmawia usunięcia grupy, która ma konta, chyba że dodasz `--leave`. Każde konto zachowuje wtedy kopię, a folder trafia do `~/.local/share/ccam/trash/`, chyba że używa go inna grupa.
- Przed wprowadzeniem grup `config.toml` przechowywał tabelę `[profile.share]` dla każdego konta. ccam odczytuje ten format jako grupy, niczego nie przenosząc. Przy pierwszym zapisie zapisuje nowy format, a stary plik zachowuje jako `config.toml.bak`. Elementy współdzielone przez te same konta stają się jedną grupą nazwaną od folderu; pozostałe dostają nazwy takie jak `default-history` i dalej używają tego samego folderu.

### W Windows

Windows tworzy dowiązania symboliczne tylko w Trybie dewelopera (Ustawienia > System > Dla deweloperów) albo dla administratora. Bez tego uprawnienia współdzielony folder staje się punktem połączenia katalogu, który dla Claude Code działa tak samo. Współdzielony plik (`settings.json`, `CLAUDE.md`, `history.jsonl`) może być tylko dowiązaniem symbolicznym, więc zmiana grupy, która go wymaga, jest odrzucana, zanim cokolwiek zostanie przeniesione.

## Bezpieczeństwo i przypadki brzegowe

- **Działające konta.** Zmiana jest odrzucana, gdy Claude Code działa jako jedno z kont, które przenosi, bo sesja dalej zapisuje do tych plików. Opuszczenie grupy i dołączenie do niej dotyczą tylko wskazanego konta. `--dry-run` nadal wypisuje plan.
- **Kopie zapasowe.** Każdy plan dostaje własny folder `.ccam-backup/<time>/` w folderze konta. ccam nigdy ich nie usuwa.
- **Twoje własne linki.** Jeśli ścieżka jest już linkiem, który sam utworzyłeś, na przykład do repozytorium dotfiles, a grupa jest pusta, ccam kopiuje cel linku do grupy, a sam link przenosi do kopii zapasowej. Cel nie jest ruszany.
- **Zaimportowane foldery.** Konto zarejestrowane w miejscu przez `ccam import` (na przykład `~/.claude`) zmienia się w miejscu, więc zwykłe `claude` też widzi zmianę; plan o tym informuje.
- **Dwa konta w jednym folderze.** Konta, które współdzielą historię, zapisują do tych samych plików. Gdy drugie z nich startuje w folderze, w którym jedno już działa, ccam ostrzega i pyta.
- **Ręczne zmiany.** Członek grupy, który nie jest kontem, jest pomijany i zgłaszany przez `ccam doctor`. Nazwy `local`, `on`, `off` i `none` nie mogą być nazwami grup.

## Powiązane

- [Panel](/features/dashboard)
- [Współdzielenie serwerów MCP, ustawień i historii](/docs/sharing)
- [Współdzielone serwery MCP](/features/mcp)
- [Kopiowanie istniejącej konfiguracji](/features/migrate)
- [Wszystkie polecenia](/docs/commands)
