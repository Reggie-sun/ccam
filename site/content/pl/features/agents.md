---
title: Sesje w tle wszystkich kont
description: ccam agents wyświetla w jednej tabeli agentów w tle i otwarte sesje Claude Code ze wszystkich twoich kont, a ccam attach otwiera sesję po jej identyfikatorze.
updated: 2026-10-26
order: 16
translation: machine
sourceHash: "06281fb13119"
---

Claude Code przechowuje sesje w tle (`claude --bg`, widok agentów) osobno dla każdego konta: `claude agents` pokazuje tylko sesje konta, jako które działa. `ccam agents` pyta wszystkie konta naraz i wyświetla ich sesje razem, a `ccam attach <id|name>` otwiera sesję w tle jako konto, które ją uruchomiło.

## Użycie

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique prefix of the ID is enough
ccam attach probe                     # or part of the session's name
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` to stan sesji w tle (`working`, `blocked`, `done`, `failed`, `stopped`) i to, co robi jej proces (`busy`, `waiting`, `idle`); otwarta sesja interaktywna pokazuje `open`.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--all` | Wyświetla też sesje w tle, które się zakończyły, nie powiodły lub zostały zatrzymane | wyłączone |
| `--json` | Wypisuje sesje jako JSON, każdą z polem `account`, czyli kontem, do którego należy | wyłączone |

`ccam attach <id|name> [claude args]` przekazuje wszystko, co podano po identyfikatorze, do `claude attach` bez zmian, łącznie z flagami: `ccam attach 49ec --dry-run` przekazuje `--dry-run` do Claude Code. Aby zamiast tego wypisać polecenie, umieść `--dry-run` przed identyfikatorem (`ccam --dry-run attach 49ec` lub `ccam attach --dry-run 49ec`).

## Jak to działa

1. Dla każdego konta ccam uruchamia `claude agents --json` (z `--all`, jeśli podano) z `CLAUDE_CONFIG_DIR` tego konta, równolegle, każde z limitem czasu 15 sekund. Odczytuje to własną listę sesji konta i nie uruchamia usługi Claude Code działającej w tle.
2. Każdy wpis dostaje nazwę konta, a lista jest sortowana według czasu uruchomienia, od najnowszych. Konto, które nie odpowiada, jest wymieniane na standardowym wyjściu błędów i pomijane, a ccam kończy wtedy działanie z kodem 1, żeby skrypt nie pomylił błędu z pustą listą.
3. `ccam attach` szuka `jobs/<id>` w folderze każdego konta, gdzie Claude Code przechowuje sesję w tle, i uruchamia `claude attach <id>` jako to konto, tak samo jak `ccam <name>` uruchamia Claude Code. Słowo, które nie jest prefiksem identyfikatora, jest porównywane z nazwami sesji bez względu na wielkość liter; słowo wyglądające jak identyfikator nigdy nie jest brane za nazwę.

## Bezpieczeństwo i przypadki brzegowe

- Identyfikatory sesji w tle należą do jednego konta: `claude attach` uruchomione jako inne konto ich nie znajdzie, i właśnie to rozwiązuje `ccam attach`.
- Jeśli sesje na dwóch kontach mają wspólny prefiks identyfikatora lub słowo w nazwie, `ccam attach` wymienia je i prosi o dłuższy identyfikator lub nazwę.
- ccam przekazuje swoje współdzielone serwery MCP do Claude Code przez `--mcp-config=<file>` w jednym argumencie, co akceptują `claude agents` i `claude attach`; sesje w tle dalej używają tego pliku, gdy się restartują.
- Usługa w tle konta z kluczem API bierze klucz z powłoki, która ją uruchomiła. Pierwszą sesję w tle takiego konta uruchom przez ccam (`ccam gw --bg …`), żeby usługa miała klucz.

## Powiązane

- [Uruchamianie Claude Code](/features/launch)
- [Wybór konta dla skryptów](/features/pick)
- [Panel](/features/dashboard)
