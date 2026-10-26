---
title: Współdzielenie przez grupy
description: Grupa ccam współdzieli wybrane elementy między swoimi kontami przez jeden folder; ccam group przenosi dane z planem, kopiami i wyborem przy konfliktach.
updated: 2026-10-26
order: 4
translation: machine
sourceHash: "3f94dd13dbf1"
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
ccam group set default --add history              # all accounts in default start sharing history
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # the group starts from work's copies
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # first leave the group that shares those items now
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
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261023-213057
  link     projects, history.jsonl, file-history, tasks, paste-cache: profiles/work → groups/default
  Note: projects: 1 file differs from the group's copy (-home-me-code-app/s1.jsonl); the group's is kept, the account's stays in profiles/work/.ccam-backup/20261023-213057/projects
Proceed? [y/N]
```

Te same zmiany możesz wprowadzić na karcie **Groups** [panelu](/features/dashboard); naciśnij `s` na koncie, żeby wybrać jego grupy.

## Elementy

| Element | Co obejmuje |
| --- | --- |
| `mcp` | serwery MCP, przez `mcp.json` grupy ([Współdzielone serwery MCP](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | globalny `CLAUDE.md` |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` i `history.jsonl`, żeby inne konto mogło użyć `--continue` |

Jeśli żadna z grup konta nie współdzieli danego elementu, ten element zostaje we własnym folderze konta.

## Polecenia i opcje

| Polecenie | Działanie |
| --- | --- |
| `ccam group` lub `ccam group list` | wyświetla grupy |
| `ccam group show <group>` | co grupa współdzieli, jej konta, jej folder i serwery MCP |
| `ccam group add <group> --items <items> [--members <accounts>]` | tworzy grupę; elementy jej kont są do niej przenoszone |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | zmienia, co grupa współdzieli, dla każdego jej konta |
| `ccam group join <group> <account>...` | dodaje konta do grupy |
| `ccam group leave <group> <account>...` | wyjmuje konta z grupy; każde zachowuje kopię, łącznie z serwerami MCP |
| `ccam group rename <group> <new-name>` | zmienia nazwę tylko w konfiguracji |
| `ccam group remove <group> [--leave]` | usuwa grupę |

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: konto, którego kopia trafia do grupy jako pierwsza | konto domyślne, potem według numeru |
| `--replace` | `join`: najpierw opuszcza grupy, które współdzielą te same elementy | wyłączone |
| `--fresh` | konta, które przestają współdzielić element, zaczynają bez niego zamiast z kopią (zaimportowany folder i tak ją zachowuje) | wyłączone |
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
2. ccam porównuje stan stary i nowy i dla każdego konta, którego elementy się zmieniają, buduje plan (`internal/group`, `share.Build`). Plany są budowane jeden po drugim i każdy widzi, co wcześniejsze umieszczą w folderach grup, więc tylko pierwsze konto wnosi swoją kopię. Pierwsza kopia grupy pochodzi, w tej kolejności: z folderu grupy (jeśli już jakąś ma), z konta `--from`, z konta domyślnego, a potem z pozostałych kont według numeru.
3. ccam odmawia, jeśli któreś z tych kont jest uruchomione, niezależnie od tego, jak uruchomiono Claude Code (``work is running; quit Claude Code there first, then change the groups (`ccam agents` lists the sessions)``), a w Windows bez Trybu dewelopera także wtedy, gdy plik wymagałby linku.
4. ccam rozstrzyga konflikty (pytając albo według `--prefer`), wypisuje plan i prosi o potwierdzenie.
5. Tuż przed zapisem ccam sprawdza ponownie. Jeśli jakieś konto zostało uruchomione, gdy ccam czekał na twoje odpowiedzi, zmiana się zatrzymuje; zatrzymuje się też, jeśli inny proces ccam zmienił w międzyczasie te same grupy (`the groups changed while this plan waited; nothing was moved, make the change again`). Nowa konfiguracja jest zapisywana, zanim cokolwiek zostanie przeniesione, a zmiana jest odnotowywana w `~/.local/state/ccam/group-change.json`. Jeśli krok przerwie się w połowie, notatka zostaje: `ccam doctor --fix` kończy zmianę w jej kolejności, z twoimi odpowiedziami, łącznie z krokami MCP, a do tego czasu nie rusza żadna inna zmiana grup (``the last group change stopped halfway; finish it with `ccam doctor --fix` first``).
6. Najpierw wykonywane są kroki MCP. Gdy konta wyjmują rozszerzenia z grupy, listy wtyczek grupy są przestawiane na folder grupy wszędzie tam, gdzie ich ścieżki prowadzą przez folder konta (`group plugin lists:` w planie). Potem plany są wykonywane po kolei. Każdy krok przenosi dane do grupy albo do folderu kopii zapasowej; tylko pliki identyczne z kopią grupy są usuwane zamiast trafiać do kopii zapasowej.

ccam przechowuje listę dozwolonych ścieżek, które obejmuje każdy element (`share.Layout`). Nic spoza niej nigdy nie trafi do grupy. ccam wymienia też `.credentials.json` i `.claude.json` jako nigdy niewspółdzielone (`share.Never`) i zatrzymuje się z błędem wewnętrznym, jeśli któryś z nich kiedykolwiek pojawi się w układzie (`Build`). Element `mcp` nie ma ścieżek: lista serwerów jest w `.claude.json`, więc [współdzielenie MCP](/features/mcp) działa przez `mcp.json` grupy.

Oto co dzieje się z jednym folderem, `projects/`, gdy `side` zaczyna współdzielić historię, a grupa ma już zawartość:

<!-- motion:share-on -->

Jeśli grupa nie ma jeszcze niczego pod tą ścieżką, konto jest pierwszym, które ją współdzieli, a jego folder zostaje tam przeniesiony bez zmian. W przeciwnym razie każdy plik, którego grupa nie ma, jest przenoszony do grupy, a pliki istniejące po obu stronach zostają w folderze konta, więc zachowywana jest wersja grupy. Pliki identyczne z plikami grupy są usuwane, a plan wymienia te, które się różnią (`Note: projects: 1 file differs from the group's copy (…); the group's is kept, the account's stays in …`).

Folder konta trafia następnie do `.ccam-backup/<time>/` w folderze konta, a na jego miejsce wstawiany jest link. Jeśli wpis grupy sam jest linkiem do własnego wpisu dołączającego konta, zawartość konta trafia do grupy na jego miejsce. Pliki rozmów mają nazwy według identyfikatora sesji, więc historie z kilku kont łączą się bez konfliktów. Rejestry `plugins/` przechowują ścieżki bezwzględne prowadzące przez folder konta; ccam przestawia je na folder grupy i scala wpisy konta z rejestrem, który grupa już ma.

Opuszczenie grupy odwraca link:

<!-- motion:share-off -->

ccam najpierw kopiuje bieżącą zawartość grupy do `.ccam-backup/incomplete-<name>` w folderze konta i zastępuje link dopiero wtedy, gdy kopia jest kompletna. Jeśli kopiowanie się przerwie, link zostaje na miejscu, a `ccam doctor --fix` kopiuje ponownie (plan takiego dokończenia pokazuje `restore <entry>: … (copied back from the group before a stop)`). Grupa zachowuje swoją zawartość dla pozostałych kont.

- ccam usuwa link tylko wtedy, gdy nadal wskazuje do folderu grup tak jak przy tworzeniu planu. Ścieżka, która się od tego czasu zmieniła, na przykład zwykły plik w miejscu linku, zostaje bez zmian, a zmiana się zatrzymuje.
- Z `--fresh` kopiowanie jest pomijane.
- Konto, które przenosi element do grupy, której folder jest jeszcze pusty, wnosi kopię tego, co widziało, więc zachowuje swoje pliki. Ta kopia jest przygotowywana w `<group>/.ccam-incomplete-<name>`, dopóki nie jest kompletna.
- Konto, które przestaje współdzielić MCP, zachowuje serwery grupy: trafiają one do jego własnego `.claude.json` (kopia zapasowa zostaje), z wyjątkiem sytuacji, gdy konto ma już serwer o tej samej nazwie, który pozostaje bez zmian. Z `--fresh` plan wymienia zamiast tego serwery, które przestaną się ładować. Ponowne dołączenie usuwa kopie, których nikt nie zmienił.
- Opuszczenie `extensions` przestawia też listy wtyczek grupy na folder grupy tam, gdzie ich ścieżki prowadziły przez folder opuszczającego konta; dla pozostałych kont przepisane ścieżki prowadzą do tych samych wtyczek.

Pojedynczy plik, który istnieje po obu stronach z różną zawartością, to konflikt:

<!-- motion:share-conflict -->

Wybierasz dla pojedynczych plików, takich jak `settings.json`; wewnątrz folderów zostaje kopia grupy, a plan wymienia pliki, które się różnią. Dostępne wybory zależą od pliku; pierwszy jest domyślny:

| Plik | Wybory |
| --- | --- |
| `history.jsonl` | `merge by line`, `keep group`, `use mine` |
| Inne pliki `.json` | `keep group`, `use mine`, `merge JSON` |
| Wszystko inne | `keep group`, `use mine` |

`merge JSON` scala obiekty rekurencyjnie, łączy tablice bez duplikatów i zachowuje wartość grupy, gdy obie strony ustawiają ten sam klucz inaczej, więc pozostałe konta nie widzą żadnej zmiany. Każdy wybór oprócz `keep group` najpierw zapisuje kopię grupy w `.ccam-backup/<time>/group/` konta. Twój własny plik zawsze trafia do kopii zapasowej.

Przy każdym uruchomieniu ccam porównuje każdą ścieżkę z układu z grupami konta, używając tylko `lstat` i `readlink` (`share.Check`). Brakujący link jest tworzony od razu. Zwykły plik tam, gdzie powinien być link, link w złe miejsce albo nieaktualny link zatrzymują uruchomienie i odsyłają do `ccam doctor --fix`.

### Zmiana nazwy, usuwanie i stary format

- `ccam group rename` zmienia tylko konfigurację. Folder zachowuje swoją nazwę, a grupa zapisuje ją jako `dir`, więc nic się nie przenosi i nic nie może zginąć.
- `ccam group remove` odmawia usunięcia grupy, która ma konta, chyba że dodasz `--leave`. Każde konto zachowuje wtedy kopię, a folder trafia do `~/.local/share/ccam/trash/`, chyba że używa go inna grupa.
- Przed wprowadzeniem grup `config.toml` przechowywał tabelę `[profile.share]` dla każdego konta. ccam odczytuje ten format jako grupy, niczego nie przenosząc. Przy pierwszym zapisie zapisuje nowy format, a stary plik zachowuje jako `config.toml.bak`. Elementy współdzielone przez te same konta stają się jedną grupą nazwaną od folderu; pozostałe dostają nazwy takie jak `default-history` i dalej używają tego samego folderu.

### W Windows

Windows tworzy dowiązania symboliczne tylko w Trybie dewelopera (Ustawienia > System > Dla deweloperów) albo dla administratora. Bez tego uprawnienia współdzielony folder staje się punktem połączenia katalogu, który dla Claude Code działa tak samo. Współdzielony plik (`settings.json`, `CLAUDE.md`, `history.jsonl`) może być tylko dowiązaniem symbolicznym, więc zmiana grupy, która go wymaga, jest odrzucana, zanim cokolwiek zostanie przeniesione. Z tego samego powodu pierwsze `ccam add` tworzy tam grupę `default` tylko z `mcp` i `extensions` i o tym informuje, a `ccam add --group` do grupy, która współdzieli ustawienia, pamięć lub historię, jest odrzucane.

## Bezpieczeństwo i przypadki brzegowe

- **Działające konta.** Zmiana jest odrzucana, gdy Claude Code działa jako jedno z kont, które przenosi, bo sesja dalej zapisuje do tych plików, niezależnie od tego, czy uruchomił ją ccam (zwykłe `claude`, IDE, sesje w tle). Sprawdzenie jest powtarzane tuż przed zapisem. Opuszczenie grupy i dołączenie do niej dotyczą tylko wskazanego konta; opuszczenie `extensions` przepisuje też ścieżki w listach wtyczek grupy, nie zmieniając tego, co znaczą dla pozostałych. `--dry-run` nadal wypisuje plan.
- **Kopie zapasowe.** Każdy plan dostaje własny folder `.ccam-backup/<time>/` w folderze konta. ccam nigdy ich nie usuwa.
- **Twoje własne linki.** Jeśli ścieżka jest już linkiem, który sam utworzyłeś, na przykład do repozytorium dotfiles, a grupa jest pusta, ccam kopiuje cel linku do grupy, a sam link przenosi do kopii zapasowej. Cel nie jest ruszany.
- **Zaimportowane foldery.** Konto zarejestrowane w miejscu przez `ccam import` (na przykład `~/.claude`) zmienia się w miejscu, więc zwykłe `claude` też widzi zmianę; plan o tym informuje. Gdy opuszcza grupę, zawsze zachowuje kopię, nawet z `--fresh`, a `ccam remove` wyjmuje je w ten sposób z jego grup, zanim zapomni konto. Dopóki jest w grupie, jego współdzielone pliki leżą w folderze grupy, więc usuń konto (albo opuść grupy), zanim usuniesz `~/.local/share/ccam`.
- **Dwa konta w jednym folderze.** Konta, które współdzielą historię, zapisują do tych samych plików. Gdy drugie z nich startuje w folderze, w którym jedno już działa, ccam ostrzega i pyta.
- **Ręczne zmiany.** Jeśli `members` zawiera nazwę, która nie jest kontem, ccam ją pomija, a `ccam doctor` ją zgłasza. Nazwy `local`, `on`, `off` i `none` nie mogą być nazwami grup.

## Powiązane

- [Panel](/features/dashboard)
- [Współdzielenie serwerów MCP, ustawień i historii](/docs/sharing)
- [Współdzielone serwery MCP](/features/mcp)
- [Kopiowanie istniejącej konfiguracji](/features/migrate)
- [Opis poleceń](/docs/commands)
