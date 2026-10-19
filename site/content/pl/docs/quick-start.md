---
title: Szybki start
description: Dodaj dwa konta Claude Code poleceniem ccam add, przełączaj się między nimi przez claude1, claude2 lub ccam <name> i wybierz, co mają współdzielić.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: c6b1c1c3ebaf
---

Dodaj każde konto poleceniem `ccam add`, a potem uruchamiaj Claude Code przez `claude1`, `claude2` lub `ccam <name>`. Większości osób to wystarczy.

## 1. Dodaj pierwsze konto

```sh
ccam add
```

Kreator pyta o nazwę, na przykład `personal` lub `work`, oraz o to, które elementy konto ma współdzielić z pozostałymi kontami. Następnie uruchamia `claude auth login`, które otwiera przeglądarkę, żebyś mógł się zalogować. Po zakończeniu zobaczysz:

```console
✓ personal (me@example.com) 추가됨 — 단축 명령: claude1
```

Na razie ccam wypisuje własne komunikaty po koreańsku; polecenia i opcje są takie same w każdym języku.

## 2. Dodaj drugie konto

```sh
ccam add work
```

Zaloguj się w przeglądarce na drugie konto. Stanie się ono kontem `claude2`.

## 3. Przełączanie

Każde z tych poleceń uruchamia Claude Code jako konto służbowe:

```sh
claude2
ccam work
ccam 2
```

Wszystko po nazwie konta trafia do Claude Code bez zmian, więc nadal używasz znanych ci opcji:

```sh
ccam work --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | claude2 -p "what caused this?"
```

Samo `ccam` otwiera wybór konta: użyj strzałek i Enter albo naciśnij numer konta.

## 4. Wybierz, co współdzielić

Domyślnie serwery MCP, `settings.json`, globalny `CLAUDE.md` oraz agenci, polecenia, umiejętności i wtyczki są współdzielone między kontami. Historia rozmów nie jest. Włącz ją dla obu kont, jeśli chcesz kontynuować rozmowy między nimi:

```sh
ccam share personal --history on
ccam share work --history on
```

ccam pokazuje, co przeniesie, zanim cokolwiek zmieni. Szczegóły znajdziesz w części [współdzielenie](/docs/sharing).

## 5. Gdy konto osiągnie limit

```sh
ccam usage
ccam next --continue
```

`ccam usage` pokazuje 5-godzinne i tygodniowe użycie wszystkich kont. `ccam next --continue` otwiera tę samą rozmowę na następnym koncie. Zobacz [limity użycia](/docs/usage-limits).

## Używasz już ~/.claude?

Zarejestruj ten folder zamiast logować się ponownie:

```sh
ccam import
```

Folder zostaje tam, gdzie jest, a zwykłe `claude` działa jak wcześniej.
