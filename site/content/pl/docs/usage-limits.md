---
title: Limity użycia i przełączanie kont
description: Sprawdź 5-godzinne i tygodniowe użycie każdego konta Claude Code przez ccam usage i otwórz konto z największym zapasem przez ccam next --auto.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: 9aebb54d40e2
---

`ccam usage` pokazuje, jaką część 5-godzinnego i tygodniowego limitu wykorzystało każde konto. `ccam next --auto` otwiera konto z największym zapasem, a `--continue` przenosi bieżącą rozmowę.

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ██████░░  72% 1h 40m    ███░░░░░  41% 3d 6h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
 3 side      ████████ 100% 22m       █████░░░  68% 1d 9h     12m ago (cached)
→ Most headroom: work   (open it with ccam next --auto)
```

Kolumny to: konto, użycie 5-godzinne, czas do resetu, użycie tygodniowe, czas do resetu oraz wiek danych.

## Skąd pochodzą liczby

ccam pyta endpoint użycia Anthropic o wszystkie konta jednocześnie, korzystając z własnego loginu każdego konta. Czeka najwyżej 3 sekundy na każde konto i pokazuje każde konto, gdy tylko odpowie.

Ten endpoint nie jest publicznie udokumentowany i może się zmienić. Gdy zapytanie się nie powiedzie, ccam pokazuje ostatnie znane liczby i oznacza je jako pochodzące z pamięci podręcznej. Okno, które zresetowało się od zapisania tych liczb, pokazuje 0%. `ccam doctor` informuje, jeśli format odpowiedzi się zmienił.

Liczby pobrane mniej niż 60 sekund temu są używane ponownie. Opcje:

| Opcja | Działanie |
| --- | --- |
| `--watch [seconds]` | Odświeżaj stale, domyślnie co 60 sekund |
| `--json` | Wynik czytelny dla maszyn |
| `-v` | Pokaż też tygodniowe limity dla poszczególnych modeli i źródło każdej liczby |
| `--no-cache` | Zawsze pytaj ponownie |
| `--source cache` | Nigdy nie pytaj; pokazuj tylko zapisane liczby |

## Wygasłe loginy

Gdy token dostępu konta wygaśnie, ccam go nie odnawia, ponieważ odnowienie go w czasie, gdy sesja Claude Code korzysta z tego samego loginu, mogłoby tę sesję wylogować. Uruchom konto raz przez `ccam <name>`, a Claude Code je odnowi.

`ccam usage --refresh` samo odnawia wygasłe tokeny, ale tylko dla kont, które nie są uruchomione. Nigdy nie odnawia loginów przechowywanych w pęku kluczy macOS ani w zaimportowanym `~/.claude`.

## Użycie na pasku stanu

```sh
ccam statusline install
```

Claude Code pokazuje wtedy nazwę konta i jego użycie na pasku stanu, na przykład `[work] · Opus · 5h 24% · 7d 41%`. Jeśli masz już polecenie paska stanu, ccam je zachowuje i dodaje przed nim nazwę konta. Pasek stanu zapisuje też swoje liczby, z których `ccam usage` korzysta, gdy endpoint nie odpowiada.

## Przełączanie

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # the next account that shares your conversation history
```

`ccam next` pomija konta, które nie są zalogowane, a `--auto` pomija konta na poziomie 100%. Aby kontynuować rozmowę, oba konta muszą współdzielić historię rozmów; zobacz [współdzielenie](/docs/sharing).

[Poradnik o pracy z limitem użycia](/guides/usage-limit-workflow) pokazuje pełny przykład.
