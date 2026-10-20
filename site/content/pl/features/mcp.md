---
title: Współdzielone serwery MCP
description: Jedna lista serwerów MCP na grupę dzięki ccam mcp add, add-json, list, remove i import, przekazywana przez --mcp-config lub synchronizowana do każdego konta.
updated: 2026-10-20
order: 5
translation: machine
sourceHash: f1de7db4b7be
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
ccam mcp list                                                     # every group, and servers only one account has
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
```

Dodanie nazwy, którą grupa już ma, zastępuje jej definicję i wypisuje `updated` zamiast `added`. Serwery trafiają do konta tylko wtedy, gdy współdzieli ono MCP z tą grupą (`ccam share <name> --mcp on`); gdy żadne konto jeszcze nie współdzieli MCP, `add` o tym informuje.

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

- `~/.local/share/ccam/groups/<group>/mcp.json`: współdzielone serwery, `{"mcpServers": {…}}`, w formacie, który czyta `--mcp-config`. Jest zapisywany przez `fsutil.AtomicWrite` z uprawnieniami `0600`, bo wpisy serwerów często zawierają klucze API.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: własny plik Claude Code z serwerami konta. ccam odczytuje z niego tylko wartość `mcpServers` skanerem `profile.TopLevelValue` i na ścieżce uruchamiania nigdy nie parsuje całego pliku.
- `~/.local/state/ccam/mcp/<name>/`: plik `seen`, `sync-stamp` i kopie zapasowe `.claude.json`.

Są dwie strategie, ustawiane przez `mcp_strategy` w `~/.config/ccam/config.toml`. Domyślna jest strategia flagi:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` znajduje grupę MCP konta. `mcp.FlagArgs` sprawdza jednym `stat`, czy `mcp.json` grupy istnieje, bez parsowania go.
2. `Plan.InsertArgs` wstawia `--mcp-config=<path>` zaraz po `claude`, przed twoimi argumentami, więc twoje własne `--mcp-config` lub `--strict-mcp-config` nadal obowiązują dodatkowo.
3. `mcp.NewLocal` sprawdza, czy `.claude.json` konta zmienił się od ostatniego uruchomienia, i dopiero wtedy parsuje jego serwery.
4. `syscall.Exec` zastępuje ccam przez Claude Code.

Ścieżka jest przekazywana jako jeden argument, `--mcp-config=<path>`, a nie dwa. `--mcp-config` w Claude Code przyjmuje kilka wartości, więc `claude --mcp-config a.json "prompt"` odczytałoby prompt jako drugi plik konfiguracji i zakończyło się błędem `MCP config file not found`. Forma jednoargumentowa kończy się na sobie, więc następujący po niej prompt lub podpolecenie takie jak `mcp list` działa. Claude Code wczytuje serwery z `--mcp-config` oprócz własnych serwerów konta; przy kolizji nazw wygrywa wpis w `.claude.json`. Strategia flagi zapisuje `.claude.json` tylko na potrzeby Remote Control.

`claude remote-control` (`rc`) nie przyjmuje `--mcp-config`, a jego sesje odczytują serwery tylko z `.claude.json`. Dla `ccam work rc` ccam pomija flagę i, jeśli działa w terminalu, pyta, czy skopiować serwery grupy do `.claude.json`, z kopią zapasową, tak jak robi to strategia `sync`. Pyta, dopóki kopia nie istnieje, a potem ponownie dopiero po zmianie `mcp.json` grupy ([Remote Control](/features/launch#remote-control)). Ponieważ taka kopia wygrywa z `--mcp-config` przy kolizji nazw, `ccam mcp add`, `add-json`, `remove` i `import` aktualizują kopie przechowywane przez inne konta grupy:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

Serwer dodany w trakcie sesji przez `claude mcp add` trafia tylko do `.claude.json` tego konta. ccam nie może go wtedy zauważyć, bo Claude Code go zastąpił, więc sprawdza to przy następnym uruchomieniu. Plik `seen` przechowuje czas modyfikacji `.claude.json` i nazwy serwerów, które już zgłoszono. Gdy czas się różni, ccam porównuje nazwy z grupą i dla każdego nowego serwera raz wypisuje wskazówkę, nie czekając na odpowiedź:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
```

Z `mcp_strategy = "sync"` serwery pojawiają się w `claude mcp list` i `/mcp` jako zwykłe serwery użytkownika, a ccam przed uruchomieniem zapisuje je w `.claude.json` każdego konta:

<!-- motion:mcp-sync -->

1. `mcp.Sync` porównuje czas modyfikacji `mcp.json` grupy z `sync-stamp`. Jeśli się nie zmienił, `.claude.json` konta nie jest otwierany.
2. Bierze własną blokadę Claude Code, folder `.claude.json.lock` obok pliku, tworzony przez `mkdir`. Czeka do 2 sekund i traktuje blokadę starszą niż 10 sekund jako nieaktualną, tak jak robi to biblioteka blokad Claude Code.
3. Zapisuje bieżący `.claude.json` w `backups/` w folderze stanu i zachowuje 5 najnowszych kopii.
4. Ustawia każdy serwer grupy w `mcpServers`, usuwa serwery, które wcześniej zapisał, a których grupa już nie ma, zachowuje własne serwery konta i zapisuje plik atomowo z jego pierwotnymi uprawnieniami.
5. Zapisuje nowy znacznik czasu i nazwy, którymi zarządza, a Claude Code startuje bez `--mcp-config`.

Synchronizacja jest pomijana z komunikatem `ccam: account work is already running, skipping MCP sync`, gdy działa sesja tego konta, bo Claude Code przepisuje `.claude.json` w trakcie działania. Nieudana synchronizacja wypisuje `ccam: MCP sync failed (launching anyway): …` i mimo to uruchamia konto. Wskazówka o serwerach dodanych w trakcie sesji dotyczy tylko strategii flagi.

`ccam mcp import` przenosi własne serwery konta do jego grupy:

<!-- motion:mcp-import -->

ccam odczytuje serwery konta przez `mcp.LocalServers` i porównuje każdy z serwerem grupy o tej samej nazwie przez `mcp.SameDef`, które ignoruje kolejność kluczy i odstępy. Serwer, którego grupa nie ma albo ma z tą samą definicją, jest kopiowany. Inny jest pomijany, bo definicji grupy używają wszystkie pozostałe konta:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

Z `--force` serwer grupy jest zastępowany, a poprzedni `mcp.json` najpierw trafia do `.ccam-backup/<time>/mcp.json` w folderze grupy. Następnie `mcp.RemoveLocal` usuwa skopiowane wpisy z `.claude.json` konta, pod tą samą blokadą i po wykonaniu kopii zapasowej w folderze stanu. Jest to potrzebne, bo przy kolizji nazw wpis konta wygrałby z wpisem grupy. `--keep-local` je zostawia.

## Bezpieczeństwo i przypadki brzegowe

- **Logowania.** Serwery logujące się przez OAuth przechowują swoje tokeny w `.credentials.json` każdego konta (`mcpOAuth`), a nie na liście serwerów, więc współdzielenie serwera nie współdzieli jego logowania. Zaloguj się raz na każdym koncie przez `/mcp` albo skopiuj logowania z `~/.claude` poleceniem [`ccam migrate`](/features/migrate). Grupa nigdy nie dostaje pliku logowania.
- **Działające konta.** `import` kopiuje do grupy, ale zostawia własne wpisy działającego konta i o tym informuje; uruchom polecenie ponownie po zakończeniu sesji. Synchronizacja pomija działające konta.
- **Kopie dla Remote Control.** Kopia działającego konta nie jest zmieniana: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Z `mcp_strategy = "sync"` nic nie trzeba aktualizować, bo każde uruchomienie synchronizuje.
- **Uszkodzone pliki.** Jeśli `.claude.json` nie jest poprawnym JSON-em, synchronizacja go nie rusza i zgłasza błąd. Jeśli blokady nie da się uzyskać w ciągu 2 sekund, ccam zgłasza `could not lock … (Claude Code is writing the file)`.
- **Konta, które nie współdzielą MCP.** `import` odmawia i każe najpierw uruchomić `ccam share side --mcp on`.
- **Nazwy grup.** `-g` przyjmuje tylko nazwę grupy; staje się ona folderem w `groups/`, więc nazwy takie jak `../x` są odrzucane.

## Powiązane

- [Współdzielenie między kontami](/features/share)
- [Remote Control z ccam](/guides/remote-control)
- [Współdzielenie serwerów MCP, ustawień i historii](/docs/sharing)
- [Kopiowanie istniejącej konfiguracji](/features/migrate)
- [Wszystkie polecenia](/docs/commands)
