---
title: Opis poleceń
description: Wszystkie polecenia i opcje ccam, od ccam add i ccs <name> po ccam usage, ccam next i ccam doctor.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 18ecad9a1104
---

`ccam` i `ccs` to ten sam program; używaj tego, które wolisz. Wpisz `--dry-run` przed dowolnym poleceniem, żeby zobaczyć, co by zrobiło, bez wykonywania go, na przykład `ccs --dry-run work`.

<!-- commands -->

## Uruchamianie Claude Code

`ccs <name|number> [arguments]` uruchamia Claude Code jako to konto. ccam nie odczytuje argumentów po nazwie; przekazuje je do `claude` w tej samej kolejności, a standardowe wejście, kod wyjścia i sygnały działają dokładnie tak jak w samym `claude`. `--` po nazwie jest opcjonalne.

`ccs` z argumentami, ale bez nazwy konta, na przykład `ccs -p "summarise"`, używa domyślnego konta. Ustawisz je poleceniem `ccam use <name>`.

Skróty `claude1`, `claude2` i kolejne to linki do ccam w `~/.local/bin`. `ccam shims install --names` tworzy też `claude-<name>` dla każdego konta.

## Dodawanie i usuwanie kont

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` bez podanego folderu szuka `~/.claude` i `~/.claude-*` i pyta o każdy z nich. Zaimportowane foldery są rejestrowane tam, gdzie są; ccam nigdy nie przenosi ani nie kopiuje loginu.

## Współdzielenie

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Każdy element przyjmuje `on` (grupa `default`), `off` lub nazwę grupy. Zobacz [współdzielenie](/docs/sharing).

## Serwery MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## Użycie i przełączanie

```sh
ccs usage            # also: ccs u
ccs usage --watch    # refresh every 60 seconds
ccs usage --json
ccs next --continue
ccs next --auto
```

## Konserwacja

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## Pliki

| Ścieżka | Zawartość |
| --- | --- |
| `~/.config/ccam/config.toml` | Konta i to, co współdzielą |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` każdego konta |
| `~/.local/share/ccam/groups/<group>/` | Współdzielone ustawienia, serwery MCP, rozszerzenia i historia |
| `~/.local/state/ccam/` | Pamięć podręczna użycia i czasy ostatniego użycia |
| `~/.local/bin/claude1`, `ccs` | Linki do ccam |

Zmienne `XDG_CONFIG_HOME`, `XDG_DATA_HOME` i `XDG_STATE_HOME` zmieniają położenie tych folderów.
