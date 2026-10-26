---
title: Szybki start
metaTitle: "Szybki start: dodaj i przełączaj konta Claude Code"
description: Dodaj dwa konta Claude Code poleceniem ccam add, nadaj im skróty takie jak ccwork, przełączaj się między nimi i wybierz, co mają współdzielić.
updated: 2026-10-26
order: 2
translation: machine
sourceHash: "bd893bbb83ab"
---

Dodaj każde konto poleceniem `ccam add`, a potem uruchamiaj Claude Code przez `ccam <name>` lub przez skrót o nazwie wybranej przez ciebie, na przykład `ccwork`. Większości osób to wystarczy.

Używasz już Claude Code? Zanim zaczniesz, przeczytaj [Używasz już ~/.claude?](#already-using-claude) na końcu tej strony.

## 1. Dodaj pierwsze konto

```sh
ccam add
```

Kreator pyta o nazwę, na przykład `personal` lub `work`, o to, które elementy konto ma współdzielić z pozostałymi kontami, oraz czy utworzyć skrót (zobacz krok 3). Następnie uruchamia `claude auth login`, które otwiera przeglądarkę, żebyś mógł się zalogować. Na zdalnej maszynie twoja przeglądarka może zatrzymać się na stronie `http://localhost:…/callback`, która się nie wczytuje; wklej ten adres do terminala, a ccam dokończy logowanie. Po zakończeniu logowania zobaczysz:

```console
✓ personal (me@example.com) added
  no shortcut (add one with `ccam shortcut personal ccpersonal`)
```

## 2. Przenieś istniejące ustawienia

Jeśli używałeś już Claude Code, twoje ustawienia, `CLAUDE.md` i rozszerzenia są w `~/.claude`, a serwery MCP w `~/.claude.json`. Skopiuj je do współdzielonej grupy `default`, żeby dostało je każde konto:

```sh
ccam migrate
```

ccam pokazuje, co skopiuje, i najpierw pyta o zgodę. Tylko odczytuje `~/.claude`, nigdy nie kopiuje danych logowania i pomija historię rozmów, chyba że dodasz `--history`. Szczegóły znajdziesz w części [kopiowanie istniejącej konfiguracji](/docs/sharing#copy-an-existing-setup).

## 3. Dodaj drugie konto

```sh
ccam add work
```

Kreator pomija nazwę, którą podałeś w wierszu poleceń, i pokazuje grupy do dołączenia z zaznaczoną grupą `default`. Tym razem nadaj kontu własny skrót, na przykład `ccwork`, żeby uruchamiać je jednym słowem. Po wyborze grup kreator pyta:

```console
 Create a shortcut command for work? [y/N]
 A command such as ccwork that opens Claude Code as this account.
```

Naciśnij `y`, a potem wpisz nazwę lub naciśnij Enter, aby użyć proponowanej:

```console
 Shortcut name [ccwork]: █
 a-z, 0-9, '.', '-' and '_'; Enter alone uses ccwork
```

Potem zaloguj się w przeglądarce na drugie konto. ccam potwierdza konto i jego skrót:

```console
✓ work (work@example.com) added
  shortcut: ccwork
```

## 4. Przełączanie kont

`ccam work` uruchamia Claude Code jako konto służbowe. Jeśli pominąłeś pytanie o skrót, skrót możesz dodać w dowolnej chwili:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
```

Każde z tych poleceń uruchamia teraz Claude Code jako konto służbowe (`2` to jego numer w `ccam list`):

```sh
ccwork
ccam work
ccam 2
```

Samo `ccam` otwiera [panel](/features/dashboard): przejdź do konta strzałkami lub jego numerem, a potem naciśnij Enter.

Wszystko po skrócie lub nazwie konta trafia do Claude Code bez zmian, więc nadal używasz znanych opcji:

```sh
ccwork --resume
ccam work -p "fix the failing test" --output-format json
cat error.log | ccwork -p "what caused this?"
```

## 5. Wybierz, co współdzielić

Domyślnie serwery MCP, `settings.json`, globalny `CLAUDE.md` oraz rozszerzenia (agenci, polecenia, umiejętności i wtyczki) są współdzielone między kontami przez grupę `default`. Historia rozmów nie jest. Dodaj ją do grupy, jeśli chcesz kontynuować rozmowy między kontami:

```sh
ccam group set default --add history
```

ccam pokazuje, co przeniesie dla każdego konta, zanim cokolwiek zmieni. Szczegóły znajdziesz w części [współdzielenie przez grupy](/features/share).

## 6. Gdy konto osiągnie limit

```sh
ccam usage
ccam next --continue
```

`ccam usage` pokazuje 5-godzinne i tygodniowe użycie wszystkich kont. `ccam next --continue` otwiera tę samą rozmowę na następnym koncie. W skryptach `ccam --fallback work -p "…"` przełącza konto samo. Zobacz [limity użycia](/docs/usage-limits).

**Miejsca w planach Team i Enterprise.** Automatyczny wybór konta (`ccam next`, `--fallback`, `ccam pick`, klawisz `n` w panelu oraz, gdy go włączysz, [`auto_switch`](/features/next#auto-switch)) przełącza też między miejscami w tej samej organizacji Claude Team lub Enterprise, co może naruszać warunki Anthropic. Przeczytaj [FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) i skonsultuj to ze swoją organizacją i z Anthropic, zanim zaczniesz na nim polegać.

## Używasz już ~/.claude?

Zarejestruj ten folder zamiast logować się ponownie:

```sh
ccam import
```

Folder zostaje tam, gdzie jest, a zwykłe `claude` działa jak wcześniej.
