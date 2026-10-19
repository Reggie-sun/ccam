---
title: Jak używać dwóch kont Claude Code na jednym komputerze
description: Uruchamiaj kilka kont Claude Code obok siebie dzięki CLAUDE_CONFIG_DIR, ręcznie z aliasami powłoki albo automatycznie z ccam.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: 6a4b7d961e10
---

Aby mieć kilka kont Claude Code na jednym komputerze, daj każdemu kontu własny folder konfiguracji. Claude Code odczytuje folder ze zmiennej środowiskowej `CLAUDE_CONFIG_DIR`, a każdy folder przechowuje własny login, ustawienia i historię. Możesz to skonfigurować ręcznie za pomocą aliasu powłoki albo zostawić to ccam.

## Ręcznie, z aliasem

To podejście opisuje dokumentacja Claude Code. Dodaj do `~/.bashrc` lub `~/.zshrc` jeden wiersz na każde dodatkowe konto:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Otwórz nowy terminal i uruchom `claude-work`. Claude Code poprosi o zalogowanie, bo folder jest nowy. Od tej chwili `claude` to twoje pierwsze konto, a `claude-work` drugie.

Przy dwóch kontach to dobrze działa. Warto wiedzieć:

- Każdy folder jest na początku pusty: bez ustawień, bez `CLAUDE.md`, bez serwerów MCP, bez agentów i umiejętności. Kopiujesz lub linkujesz to, czego potrzebujesz, i sam dbasz o zgodność kopii.
- Aliasy istnieją tylko w interaktywnych powłokach. Skrypty, pliki Makefile i inne programy ich nie widzą.
- Serwery MCP są w `.claude.json`, obok loginu, więc nie da się ich po prostu podlinkować między folderami.
- Rozmowy zostają na koncie, na którym się zaczęły, więc nie możesz kontynuować ich na innym koncie.
- Nic nie pokazuje użycia wszystkich kont naraz.

## Z ccam

ccam automatyzuje ten sam pomysł i dodaje brakujące elementy:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

Teraz `ccam personal` i `ccam work` uruchamiają oba konta, podobnie jak `ccwork`, skrót nadany kontu służbowemu. Ponieważ są to prawdziwe programy, a nie aliasy, działają też w skryptach. Ustawienia, `CLAUDE.md`, serwery MCP i rozszerzenia są współdzielone od początku, a po włączeniu współdzielonej historii rozmów możesz kontynuować rozmowę na drugim koncie przez `ccam next --continue`.

Jeśli już używasz `~/.claude`, uruchom `ccam import` zamiast dodawać to konto ponownie albo `ccam migrate`, aby skopiować jego ustawienia, serwery MCP i rozszerzenia do współdzielonych ustawień.

## Co wybrać

Dla dwóch kont bez niczego współdzielonego wystarczy alias. Gdy chcesz współdzielić ustawienia i serwery MCP, mieć więcej kont, widzieć użycie na pierwszy rzut oka albo kontynuować rozmowy, ccam oszczędza ci ręcznego pilnowania wszystkiego. [Porównanie](/compare/aliases) zestawia różnice obok siebie.
