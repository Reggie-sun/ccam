---
title: Limity użycia i przełączanie kont
description: Sprawdź 5-godzinne i tygodniowe użycie każdego konta Claude Code przez ccam usage i otwórz konto z największym zapasem przez ccam next --auto.
updated: 2026-10-24
order: 5
translation: machine
sourceHash: "ccaed084207a"
---

`ccam usage` pokazuje, jaką część 5-godzinnego limitu, tygodniowego limitu i własnego tygodniowego limitu każdego modelu (takiego jak Fable) wykorzystało każde konto. `ccam next --auto` otwiera konto z największym zapasem, a `--continue` przenosi bieżącą rozmowę.

```console
$ ccam usage
 1 personal  just now
   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h
   Fable    ██████████████████░░  88%  resets in 1d 6h
 2 work      just now
   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m
   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h
   Fable    █░░░░░░░░░░░░░░░░░░░   3%  resets in 5d 2h
 3 side      12m ago (cached)
   5-hour   ████████████████████ 100%  resets in 22m
   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h
   Fable    ████████░░░░░░░░░░░░  40%  resets in 1d 9h
→ Most headroom: work   (open it with ccam next --auto)
```

Każde konto ma wiersz z informacją, jak stare są jego liczby, a pod nim jeden pasek na każdy limit z czasem do resetu: okno 5-godzinne, tygodniowe okno dla wszystkich modeli i własne tygodniowe okno każdego modelu, takie jak Fable. Od 80% pasek zmienia kolor na żółty, a przy 100% na czerwony. Ostatni wiersz wskazuje konto, które otworzyłby `ccam next --auto`; z `auto_strategy = "expiring"` lub `"balance"` brzmi on zamiast tego `→ next --auto picks <name>: <reason>`.

**Miejsca w planach Team i Enterprise.** Automatyczny wybór konta (`ccam next`, `--fallback`, `ccam pick`, klawisz `n` w panelu) przełącza też między miejscami w tej samej organizacji Claude Team lub Enterprise. Centrum pomocy Anthropic podaje: „Each seat on a Team plan is assigned to one named individual, and only that person may use it” (każde miejsce w planie Team jest przypisane do jednej imiennie wskazanej osoby i tylko ta osoba może z niego korzystać) oraz „Users can't route their work through another user's account to use that user's remaining usage” (użytkownicy nie mogą kierować swojej pracy przez konto innego użytkownika, aby wykorzystać pozostały limit użycia tego użytkownika) ([plany Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); podobnie w przypadku [planów Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Dlatego korzystanie z automatycznego wyboru między miejscami może naruszać warunki Anthropic. Zanim zaczniesz na nim polegać, skonsultuj to ze swoją organizacją i z Anthropic. Planów prywatnych (Pro, Max) dotyczy [uwaga o warunkach w FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

## Skąd pochodzą liczby

ccam pyta endpoint użycia Anthropic o wszystkie konta jednocześnie, korzystając z własnego loginu każdego konta. Czeka najwyżej 3 sekundy na każde konto i pokazuje każde konto, gdy tylko odpowie.

Ten endpoint nie jest publicznie udokumentowany i może się zmienić. Gdy zapytanie się nie powiedzie, ccam pokazuje ostatnie znane liczby i oznacza je jako pochodzące z pamięci podręcznej. Okno, które zresetowało się od zapisania tych liczb, pokazuje 0%. `ccam doctor` informuje, jeśli format odpowiedzi się zmienił.

Gdy Claude Code zatrzyma sesję na limicie użycia, [plugin ccam](/docs/install#claude-code-plugin) zapisuje, które okno było pełne i kiedy się resetuje. Do tego czasu to okno pokazuje 100% bez względu na pozostałe liczby, więc `ccam next --auto` i `--fallback` nie wybiorą tego konta ponownie zbyt wcześnie.

Liczby pobrane mniej niż 60 sekund temu są używane ponownie. Opcje:

| Opcja | Działanie |
| --- | --- |
| `--watch[=seconds]` | Odświeżaj stale, domyślnie co 60 sekund |
| `--json` | Wynik czytelny dla maszyn |
| `-v` | Pokaż też, skąd pochodzą liczby każdego konta |
| `--no-cache` | Zawsze pytaj ponownie |
| `--source cache` | Nigdy nie pytaj; pokazuj tylko zapisane liczby |

## Wygasłe loginy

Gdy token dostępu konta wygaśnie, ccam go nie odnawia, ponieważ odnowienie go w czasie, gdy sesja Claude Code korzysta z tego samego loginu, mogłoby tę sesję wylogować. Konto pokazuje ``access token expired; `ccam <name>` renews it when it starts``: uruchom je raz przez `ccam <name>`, a Claude Code odnowi token. Login, który wygasł w całości albo który wyczyścił Claude Code, pokazuje `` login expired → `ccam login <name>` `` i wymaga ponownego zalogowania.

`ccam usage --refresh` samo odnawia wygasłe tokeny, ale tylko dla kont, które nie są uruchomione, niezależnie od tego, jak uruchomiono Claude Code. Nigdy nie odnawia loginów przechowywanych w pęku kluczy macOS ani w zaimportowanym `~/.claude`, a odnowienie może potrwać do 30 sekund.

## Użycie na pasku stanu

```sh
ccam statusline install
```

Claude Code pokazuje wtedy nazwę konta i jego użycie na pasku stanu, na przykład `[work] · Opus · 5h 24% · 7d 41%`. Od 90% pasek podpowiada też `ccam next --auto --continue`. Jeśli masz już polecenie paska stanu, ccam je zachowuje i dodaje przed nim nazwę konta. Pasek stanu zapisuje też swoje liczby, z których `ccam usage` korzysta, gdy endpoint nie odpowiada.

## Przełączanie

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # carry this folder's conversation on as the next account
```

`ccam next` pomija konta, które nie są zalogowane, a `--auto` pomija konta na poziomie 100% (także z oknem, które Claude Code zgłosił jako pełne, aż do jego resetu) oraz konta innej organizacji. Gdy oba konta współdzielą historię rozmów, `--continue` otwiera tę samą rozmowę; w przeciwnym razie Claude Code kopiuje ją na nowe konto. Zobacz [współdzielenie](/docs/sharing) i [przełączanie na następne konto](/features/next#continue).

W skryptach i zaplanowanych zadaniach, które uruchamiają `claude -p`, umieść `--fallback` przed nazwą konta. Gdy to konto osiągnie limit, ccam uruchamia to samo polecenie ponownie na koncie z największym zapasem, a wypisywane jest tylko wyjście przebiegu, który się zakończył:

```sh
ccam --fallback work -p "review the diff" --output-format json
```

Zobacz [Ponowienie po osiągnięciu limitu](/features/next#run-again-at-the-limit), jak ccam podejmuje decyzję.

[Poradnik o pracy z limitem użycia](/guides/usage-limit-workflow) pokazuje pełny przykład.
