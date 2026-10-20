---
title: ccam w porównaniu z aliasami powłoki
description: Aliasy powłoki z CLAUDE_CONFIG_DIR i ccam do uruchamiania kilku kont Claude Code, porównane obok siebie.
updated: 2026-10-20
order: 1
translation: machine
sourceHash: ea24bb091caf
---

Alias powłoki dla każdego konta to najprostszy sposób na uruchamianie kilku kont Claude Code. ccam robi pod spodem to samo i dodaje współdzielenie, przełączanie oraz śledzenie użycia.

## Konfiguracja z aliasami

Dodaj po jednym wierszu na każde dodatkowe konto do `~/.bashrc` lub `~/.zshrc`, otwórz nowy terminal i uruchom alias raz, aby się zalogować:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

`claude` pozostaje twoim pierwszym kontem, a `claude-work` uruchamia drugie. To metoda, którą [opisuje dokumentacja Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts).

## Konfiguracja z ccam

Zainstaluj ccam i dodaj konto, tutaj od razu z własnym skrótem:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add work --shortcut ccwork
```

`ccam work` i `ccwork` uruchamiają konto i oba działają też w skryptach.

## Porównanie obok siebie

| | Alias powłoki | ccam |
| --- | --- | --- |
| Osobne loginy | Tak | Tak |
| Konfiguracja każdego konta | Edycja `~/.bashrc`, logowanie | `ccam add` |
| Działa w skryptach i innych programach | Nie, aliasy działają tylko w powłoce | Tak, `ccam work` i skróty takie jak `ccwork` to prawdziwe polecenia |
| Współdzielone ustawienia, `CLAUDE.md`, rozszerzenia | Ręczne kopiowanie lub linkowanie | Dla każdego elementu, z kopiami zapasowymi |
| Współdzielone serwery MCP | Ręcznie, na każdym koncie | Raz przez `ccam mcp add` |
| Kontynuacja rozmowy na innym koncie | Nie | `ccam next --continue` |
| Użycie wszystkich kont | Nie | `ccam usage` |
| Wybór konta z największym zapasem | Nie | `ccam next --auto` |
| Nazwa konta na pasku stanu | Nie | `ccam statusline install` |
| Naprawa zepsutych linków | Ręcznie | `ccam doctor --fix` |
| Dodatkowe oprogramowanie | Brak | Jeden statyczny plik binarny |

Wybierz aliasy, jeśli masz dwa konta i nie chcesz niczego między nimi współdzielić. Wybierz ccam, jeśli chcesz, żeby konta działały jak jedna konfiguracja. Możesz zacząć od aliasów i przejść na ccam później: `ccam import` rejestruje istniejące foldery `~/.claude-*` tam, gdzie są.

[Jak używać kilku kont na jednym komputerze](/guides/multiple-claude-code-accounts) wyjaśnia oba podejścia krok po kroku, a [ccam i inne narzędzia do przełączania kont](/compare/account-switchers) porównuje ccam z narzędziami, które podmieniają jeden login.
