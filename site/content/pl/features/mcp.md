---
title: Współdzielone serwery MCP
description: Jedna lista serwerów MCP na grupę dzięki ccam mcp add, add-json, list, remove i import, przekazywana przez --mcp-config lub synchronizowana do każdego konta.
updated: 2026-10-23
order: 5
translation: machine
sourceHash: "fefca7f8d857"
---

Claude Code przechowuje serwery MCP użytkownika w `.claude.json`, tym samym pliku, w którym jest login konta, więc tego pliku nie można współdzielić między kontami. ccam trzyma współdzielone serwery w pliku grupy, `mcp.json`, i przekazuje go każdemu kontu z grupy przy uruchomieniu. Użyj `ccam mcp`, żeby dodać serwer raz dla wszystkich kont, a `ccam mcp import`, żeby przenieść do pozostałych kont serwer dodany w trakcie sesji na jednym z nich.

## Użycie

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

Zmiana obowiązuje od następnego uruchomienia każdego konta; działające sesje zachowują serwery, z którymi wystartowały:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
only in account side, for /home/me/code/app only: notes
```

Pod grupami `ccam mcp list` wymienia każdy serwer, który nie dociera do pozostałych kont: `only in account work: pg` dla serwera, który konto dodało dla siebie, `only in account side, for /home/me/code/app only: notes` dla serwera dodanego dla jednego folderu (zwykłe `claude mcp add`, zobacz niżej), `account work has its own …, used in its sessions in place of the group's`, gdy własny serwer konta ma nazwę serwera grupy, `.mcp.json` tego folderu z kontami, w których jest `waiting for approval in`, oraz serwery, które są `in no account, only in ~/.claude.json (plain claude's)`, z poleceniem [`ccam migrate --mcp`](/features/migrate), które je kopiuje.

Dodanie nazwy, którą grupa już ma, zastępuje jej definicję i wypisuje `updated` zamiast `added`. Serwery trafiają do konta tylko wtedy, gdy należy ono do tej grupy, a grupa współdzieli MCP (`ccam group set <group> --add mcp`); gdy grupa tego nie robi, `add` o tym informuje.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `-g`, `--group <group>` | grupa docelowa dla `add`, `add-json` i `remove` | jedyna grupa, przez którą jakieś konto współdzieli MCP, w przeciwnym razie `default` |
| `-- <command> [args…]` (`add`) | uruchamia to polecenie jako serwer stdio | |
| `--url <url>` (`add`) | adres URL zdalnego serwera (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | transport serwera z `--url` | `http` |
| `--header <'Name: value'>` (`add`) | nagłówek HTTP, można powtarzać | brak |
| `-e`, `--env <KEY=VALUE>` (`add`) | zmienna środowiskowa dla serwera stdio, można powtarzać | brak |
| `--force` (`import`) | zastępuje różniący się serwer grupy o tej samej nazwie; stary `mcp.json` trafia do kopii zapasowej | wyłączone |
| `--keep-local` (`import`) | zostawia wpisy w `.claude.json` konta | wyłączone |
| `-h`, `--help` | pomoc dla każdego polecenia | |
| `--dry-run` (globalna) | respektowana tylko przez `mcp list`; `add`, `add-json`, `remove` i `import` ją odrzucają i niczego nie zmieniają | wyłączone |

`add` wymaga polecenia po `--` albo `--url`; w przeciwnym razie kończy się błędem `a command to run (after --) or --url is required`. `--header` dotyczy serwerów z `--url`, a `--env` serwerów uruchamianych poleceniem. `add-json` przyjmuje ten sam JSON co `claude mcp add-json`, a definicja musi być obiektem z `command` lub `url`. `import` zawsze używa grupy, przez którą konto współdzieli MCP, a `list` pokazuje wszystkie grupy, więc `-g` niczego w nich nie zmienia.

## Jak to działa

Pliki, których to dotyczy:

- `~/.local/share/ccam/groups/<group>/mcp.json`: współdzielone serwery, `{"mcpServers": {…}}`, w formacie, który czyta `--mcp-config`. Jest zapisywany przez `fsutil.AtomicWrite` z uprawnieniami `0600`, bo wpisy serwerów często zawierają klucze API, i tylko pod blokadą `<group>/.mcp.json.lock`, więc dwa jednoczesne `ccam mcp add` zachowują oba swoje serwery.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: własny plik Claude Code z serwerami konta. ccam odczytuje z niego tylko wartość `mcpServers` skanerem `profile.TopLevelValue` i na ścieżce uruchamiania nigdy nie parsuje całego pliku.
- `~/.local/state/ccam/mcp/<name>/`: plik `seen`, `sync-stamp` (ze skrótem każdej kopii zapisanej przez ccam, żeby odróżnić kopię, której nikt nie zmienił, od tej, którą edytowałeś), `sync-checked` (co strategia sync widziała ostatnio), `kept` (serwery, które konto zachowało, opuszczając współdzielenie MCP) oraz kopie zapasowe `.claude.json`.

Są dwie strategie, ustawiane przez `mcp_strategy` w `~/.config/ccam/config.toml`. Domyślna jest strategia flagi:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` znajduje grupę MCP konta. `mcp.FlagArgs` sprawdza jednym `stat`, czy `mcp.json` grupy istnieje, bez parsowania go.
2. `Plan.InsertArgs` wstawia `--mcp-config=<path>` zaraz po `claude`, przed twoimi argumentami, więc twoje własne `--mcp-config` lub `--strict-mcp-config` nadal obowiązują dodatkowo.
3. `mcp.NewLocal` sprawdza, czy `.claude.json` konta zmienił się od ostatniego uruchomienia, i dopiero wtedy parsuje jego serwery.
4. `syscall.Exec` zastępuje ccam przez Claude Code.

Ścieżka jest przekazywana jako jeden argument, `--mcp-config=<path>`, a nie dwa. `--mcp-config` w Claude Code przyjmuje kilka wartości, więc `claude --mcp-config a.json "prompt"` odczytałoby prompt jako drugi plik konfiguracji i zakończyło się błędem `MCP config file not found`. Forma jednoargumentowa kończy się na sobie, więc następujący po niej prompt lub podpolecenie działa. `ccam work mcp list`, `get` i `remove` nie dostają flagi, bo Claude Code ją tam odrzuca; zamiast tego ccam wymienia serwery grupy na standardowym wyjściu błędów ([szczegóły](/features/launch#shared-mcp-servers)). Claude Code wczytuje serwery z `--mcp-config` oprócz własnych serwerów konta; przy kolizji nazw wygrywa wpis w `.claude.json`. Strategia flagi zapisuje `.claude.json` tylko na potrzeby Remote Control, gdy konto opuszcza współdzielenie MCP lub ponownie do niego dołącza ([niżej](#leaving-and-joining)), oraz dla `ccam mcp import`.

`claude remote-control` (`rc`) nie przyjmuje `--mcp-config`, a jego sesje odczytują serwery tylko z `.claude.json`. Dla `ccam work rc` ccam pomija flagę i, jeśli działa w terminalu, pyta, czy skopiować serwery grupy do `.claude.json`, z kopią zapasową, tak jak robi to strategia `sync`. Pyta, dopóki kopia nie istnieje, a potem ponownie dopiero po zmianie `mcp.json` grupy albo gdy skopiowany serwer zniknie ([Remote Control](/features/launch#remote-control)). Ponieważ taka kopia wygrywa z `--mcp-config` przy kolizji nazw, `ccam mcp add`, `add-json`, `remove` i `import` aktualizują kopie przechowywane przez inne konta grupy:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

Wiersz `✓ Updated` pojawia się tylko wtedy, gdy kopia się zmieniła. Kopia utworzona przez starszy ccam jest aktualizowana w ten sam sposób, podobnie jak kopia po `ccam migrate`.

Serwer dodany w trakcie sesji przez `claude mcp add` trafia tylko do `.claude.json` tego konta. Zwykłe `claude mcp add` używa zakresu lokalnego: serwer ładuje się dla tego jednego konta w tym jednym folderze (`-s user` przypisuje go kontu we wszystkich folderach). ccam nie może go wtedy zauważyć, bo Claude Code go zastąpił, więc sprawdza to przy następnym uruchomieniu. Plik `seen` przechowuje czas modyfikacji `.claude.json` i nazwy serwerów, które już zgłoszono. Gdy czas się różni, ccam porównuje nazwy z grupą i dla każdego nowego serwera raz wypisuje wskazówkę, nie czekając na odpowiedź, przy obu strategiach:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

Wskazówka wymienia też własny serwer konta, który ma nazwę serwera grupy, bo w sesjach konta wygrywa serwer konta.

Z `mcp_strategy = "sync"` serwery pojawiają się w `claude mcp list` i `/mcp` jako zwykłe serwery użytkownika, a ccam przed uruchomieniem zapisuje je w `.claude.json` każdego konta:

<!-- motion:mcp-sync -->

1. `mcp.Sync` porównuje czas modyfikacji `mcp.json` grupy z `sync-stamp`, a `.claude.json` konta z tym, co widział ostatnio. Jeśli nic się nie zmieniło, nic więcej się nie dzieje. W przeciwnym razie sprawdza też, czy każdy serwer grupy nadal jest w `.claude.json`, i przywraca kopię, która zniknęła (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`).
2. Bierze własną blokadę Claude Code, folder `.claude.json.lock` obok pliku, tworzony przez `mkdir`. Czeka do 2 sekund i traktuje blokadę starszą niż 10 sekund jako nieaktualną, tak jak robi to biblioteka blokad Claude Code.
3. Zapisuje bieżący `.claude.json` w `backups/` w folderze stanu i zachowuje 5 najnowszych kopii.
4. Ustawia każdy serwer grupy w `mcpServers`, usuwa serwery, które wcześniej zapisał, a których grupa już nie ma, i zapisuje plik atomowo z jego pierwotnymi uprawnieniami, tylko gdy coś się zmieniło. Zmienia lub usuwa tylko te zapisane przez siebie kopie, których nikt od tego czasu nie zmienił. Własny serwer konta o nazwie serwera grupy pozostaje serwerem konta, a ccam mówi o tym raz: `account work has its own MCP servers fs, so group default's of those names are not copied there`.
5. Zapisuje nowy znacznik czasu ze skrótem każdej kopii, a Claude Code startuje bez `--mcp-config`.

Synchronizacja jest pomijana z komunikatem `ccam: account work is already running, skipping MCP sync`, gdy działa sesja tego konta, niezależnie od tego, jak uruchomiono Claude Code, bo Claude Code przepisuje `.claude.json` w trakcie działania. Nieudana synchronizacja wypisuje `ccam: MCP sync failed (launching anyway): …` i mimo to uruchamia konto.

### Opuszczanie i dołączanie

Konto, które przestaje współdzielić MCP (`ccam group leave`, `ccam group set --remove mcp`, `ccam group remove --leave`), zachowuje serwery grupy: plan kopiuje je do jego `.claude.json`, z kopią zapasową, a jego własne serwery o tych samych nazwach zostają. Z `--fresh` nie dostaje żadnych, a plan wymienia serwery, które przestaną się ładować. Gdy ponownie dołącza do grupy, która ma te serwery, zachowane kopie, których nikt nie zmienił, są usuwane, więc późniejsze zmiany grupy do niego docierają.

`ccam mcp import` przenosi własne serwery konta do jego grupy:

<!-- motion:mcp-import -->

ccam odczytuje serwery konta, jego własne i te, które ma tylko dla jednego folderu, i porównuje każdy z serwerem grupy o tej samej nazwie przez `mcp.SameDef`, które ignoruje kolejność kluczy i odstępy. Serwer folderu staje się serwerem grupy dla każdego folderu (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`); jedna nazwa z dwiema różnymi definicjami folderów jest odrzucana. Serwery pluginów nie są proponowane, bo dostarcza je plugin. Serwer, którego grupa nie ma albo ma z tą samą definicją, jest kopiowany. Inny jest pomijany, bo definicji grupy używają wszystkie pozostałe konta:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

Z `--force` serwer grupy jest zastępowany, a poprzedni `mcp.json` najpierw trafia do `.ccam-backup/<time>/mcp.json` w folderze grupy. Następnie `mcp.Remove` usuwa skopiowane wpisy z `.claude.json` konta, jego własne i te folderu w jednym zapisie, pod tą samą blokadą i po wykonaniu kopii zapasowej w folderze stanu. `mcp.json` grupy pozostaje zablokowany od wczytania aż do usunięcia tych wpisów, więc nikt w międzyczasie nie usunie serwera z grupy. Jest to potrzebne, bo przy kolizji nazw wpis konta wygrałby z wpisem grupy. `--keep-local` je zostawia.

## Bezpieczeństwo i przypadki brzegowe

- **Logowania.** Serwery logujące się przez OAuth przechowują swoje tokeny w `.credentials.json` każdego konta (`mcpOAuth`), a nie na liście serwerów, więc współdzielenie serwera nie współdzieli jego logowania. Zaloguj się raz na każdym koncie przez `/mcp`. ccam nigdy nie kopiuje logowań, nawet [`ccam migrate`](/features/migrate): skopiowany token odświeżania wylogowuje wszystkich pozostałych posiadaczy. [`ccam ext`](/features/ext) i `ccam doctor` wymieniają serwery, do których konto nie jest zalogowane, a inne konto jest, a `ccam next --continue` mówi, których brakuje nowemu kontu. Grupa nigdy nie dostaje pliku logowania.
- **Działające konta.** `import` kopiuje do grupy, ale zostawia własne wpisy działającego konta i o tym informuje; uruchom polecenie ponownie po zakończeniu sesji. Synchronizacja pomija działające konta. Oba uwzględniają też Claude Code uruchomiony bez ccam.
- **Kopie dla Remote Control.** Kopia działającego konta nie jest zmieniana: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Z `mcp_strategy = "sync"` nic nie trzeba aktualizować, bo każde uruchomienie synchronizuje.
- **Uszkodzone pliki.** Jeśli `.claude.json` nie jest poprawnym JSON-em, synchronizacja go nie rusza i zgłasza błąd. Jeśli blokady nie da się uzyskać w ciągu 2 sekund, ccam zgłasza `could not lock … (Claude Code is writing the file)`.
- **Konta, które nie współdzielą MCP.** `import` odmawia i każe najpierw dodać konto do grupy, która współdzieli MCP (`ccam group join <group> side`).
- **Nazwy grup.** `-g` przyjmuje tylko nazwę grupy; staje się ona folderem w `groups/`, więc nazwy takie jak `../x` są odrzucane.

## Powiązane

- [Współdzielenie między kontami](/features/share)
- [Remote Control z ccam](/guides/remote-control)
- [Współdzielenie serwerów MCP, ustawień i historii](/docs/sharing)
- [Kopiowanie istniejącej konfiguracji](/features/migrate)
- [Wszystkie polecenia](/docs/commands)
