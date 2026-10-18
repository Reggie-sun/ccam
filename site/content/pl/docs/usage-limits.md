---
title: Limity użycia i przełączanie kont
description: Sprawdź 5-godzinne i tygodniowe użycie każdego konta Claude Code przez ccs usage i otwórz konto z największym zapasem przez ccs next --auto.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: a9ee77936668
---

`ccs usage` pokazuje, jaką część 5-godzinnego i tygodniowego limitu wykorzystało każde konto. `ccs next --auto` otwiera konto z największym zapasem, a `--continue` przenosi bieżącą rozmowę.

```console
$ ccs usage
 계정        5시간         초기화      주간          초기화      기준
 1 personal  ██████░░  72% 1시간 40분  ███░░░░░  41% 3일 6시간   방금
 2 work      █░░░░░░░   8% 4시간 12분  ██░░░░░░  19% 5일 2시간   방금
 3 side      ████████ 100% 22분        █████░░░  68% 1일 9시간   12분 전 (캐시)
→ 여유가 가장 많은 계정: work   (ccs next --auto 로 바로 열기)
```

Kolumny to: konto, użycie 5-godzinne, czas do resetu, użycie tygodniowe, czas do resetu oraz wiek danych.

## Skąd pochodzą liczby

ccam pyta endpoint użycia Anthropic o wszystkie konta jednocześnie, korzystając z własnego loginu każdego konta. Czeka najwyżej 3 sekundy na każde konto i pokazuje każde konto, gdy tylko odpowie.

Ten endpoint nie jest publicznie udokumentowany i może się zmienić. Gdy zapytanie się nie powiedzie, ccam pokazuje ostatnie znane liczby i oznacza je jako pochodzące z pamięci podręcznej. `ccam doctor` informuje, jeśli format odpowiedzi się zmienił.

Liczby pobrane mniej niż 60 sekund temu są używane ponownie. Opcje:

| Opcja | Działanie |
| --- | --- |
| `--watch [seconds]` | Odświeżaj stale, domyślnie co 60 sekund |
| `--json` | Wynik czytelny dla maszyn |
| `-v` | Pokaż też tygodniowe limity dla poszczególnych modeli i źródło każdej liczby |
| `--no-cache` | Zawsze pytaj ponownie |
| `--source cache` | Nigdy nie pytaj; pokazuj tylko zapisane liczby |

## Wygasłe loginy

Gdy token dostępu konta wygaśnie, ccam go nie odnawia, ponieważ odnowienie go w czasie, gdy sesja Claude Code korzysta z tego samego loginu, mogłoby tę sesję wylogować. Uruchom konto raz przez `ccs <name>`, a Claude Code je odnowi.

`ccs usage --refresh` samo odnawia wygasłe tokeny, ale tylko dla kont, które nie są uruchomione. Nigdy nie odnawia loginów przechowywanych w pęku kluczy macOS ani w zaimportowanym `~/.claude`.

## Użycie na pasku stanu

```sh
ccam statusline install
```

Claude Code pokazuje wtedy nazwę konta i jego użycie na pasku stanu, na przykład `[work] · Opus · 5h 24% · 7d 41%`. Jeśli masz już polecenie paska stanu, ccam je zachowuje i dodaje przed nim nazwę konta. Pasek stanu zapisuje też swoje liczby, z których `ccs usage` korzysta, gdy endpoint nie odpowiada.

## Przełączanie

```sh
ccs next              # the account after the one you used last
ccs next --auto       # the account with the lowest 5-hour usage
ccs next --continue   # the next account that shares your conversation history
```

`ccs next` pomija konta, które nie są zalogowane, a `--auto` pomija konta na poziomie 100%. Aby kontynuować rozmowę, oba konta muszą współdzielić historię rozmów; zobacz [współdzielenie](/docs/sharing).

[Poradnik o pracy z limitem użycia](/guides/usage-limit-workflow) pokazuje pełny przykład.
