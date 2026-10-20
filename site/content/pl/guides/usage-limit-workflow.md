---
title: Co zrobić, gdy Claude Code osiągnie limit użycia
description: Poczekaj na reset albo od razu kontynuuj tę samą rozmowę na innym koncie Claude Code. Jak działają limity 5-godzinne i tygodniowe, krok po kroku.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: "38b80d780bba"
---

Gdy konto osiągnie limit użycia (często nazywany rate limitem), masz dwie możliwości: poczekać, aż limit się zresetuje, albo od razu kontynuować tę samą rozmowę na innym swoim koncie. Z ccam to drugie wymaga jednego polecenia, `ccam next --auto --continue`, gdy historia rozmów jest już współdzielona między kontami.

## Jakie limity ma Claude Code?

W planach Pro i Max Claude Code ma [5-godzinny limit sesji i limit tygodniowy](https://support.claude.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan). Oba są wspólne z Claude w przeglądarce i w aplikacjach, a użycie liczy się do obu jednocześnie. Gdy jeden się wyczerpie, Claude Code [wyświetla komunikat](https://code.claude.com/docs/en/errors#youve-hit-your-session-limit) w rodzaju „You've hit your session limit · resets 3:45pm”. Zmiana modelu nie pomaga, bo limity obejmują wszystkie modele.

## Opcja 1: poczekaj na reset

Jeśli masz jedno konto, poczekaj. Nowsze wersje Claude Code [czekają w otwartej sesji](https://code.claude.com/docs/en/interactive-mode#wait-for-a-usage-limit-to-reset) i same wznawiają zadanie po zresetowaniu limitu. `/usage` w Claude Code pokazuje, jaką część limitów twojego planu już wykorzystałeś.

## Opcja 2: kontynuuj od razu na innym koncie

Jeśli masz drugie konto, możesz od razu kontynuować tę samą rozmowę. Skonfiguruj to raz:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work
ccam group set default --add history
ccam statusline install
```

Pasek stanu pokazuje teraz podczas pracy konto i jego 5-godzinne użycie, więc limit nie będzie zaskoczeniem.

Gdy limit zostanie osiągnięty, wyjdź z Claude Code (`/exit`), a następnie:

```sh
ccam usage
```

```console
 1 personal  ████████ 100% 22m       █████░░░  68% 1d 9h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
```

Otwórz rozmowę na koncie z największym zapasem:

```sh
ccam next --auto --continue
```

Claude Code uruchomi się jako `work` z najnowszą rozmową z bieżącego folderu, łącznie ze wszystkim, co zostało w niej dotąd powiedziane. Współdzielona historia (`ccam group set default --add history` powyżej) sprawia, że pozostaje to jedna rozmowa; bez niej Claude Code kopiuje rozmowę do `work`, a ccam podaje, ile tokenów ta pierwsza odpowiedź przeczyta ponownie.

Z [pluginem ccam](/docs/install#claude-code-plugin) w Claude Code sesja zatrzymana na limicie pokazuje powiadomienie, a `/ccam:switch` (lub `/ccam:switch work`) zapisuje rozmowę dla `ccam next --continue`, tak że nie trzeba niczego zgadywać.

## Czy ccam przełącza konta automatycznie?

Nie. ccam nie działa w trakcie twojej pracy: uruchamia Claude Code i zastępuje się nim. To ty decydujesz, kiedy przełączyć, a `ccam next --auto` wybiera konto z najniższym 5-godzinnym użyciem i pomija konta, które mają 100%, nie są zalogowane albo należą do innej organizacji. Z `relay = true` w `config.toml` ccam zamiast tego czeka obok sesji, więc `/ccam:switch` sam kontynuuje rozmowę. Niektóre inne narzędzia zamiast tego podmieniają login działającej sesji; [ccam i inne narzędzia do przełączania kont](/compare/account-switchers) wyjaśnia różnicę.

## O czym pamiętać

- Uruchamiaj `ccam next --continue` w tym samym folderze, w którym pracowałeś; rozmowy są przechowywane osobno dla każdego folderu projektu.
- Nie uruchamiaj dwóch kont na tej samej rozmowie jednocześnie. ccam ostrzeże cię, jeśli spróbujesz.
- Bez współdzielonej historii rozmowa jest kopiowana na drugie konto i obie kopie toczą się dalej osobno; punkty kontrolne, z których `/rewind` przywraca kod, zostają na pierwszym koncie.
- Pierwsza wiadomość na nowym koncie może zużyć więcej jego limitu niż zwykle, bo to konto nie ma jeszcze rozmowy w pamięci podręcznej.
- To, czy wolno ci w ten sposób przełączać się między kontami, zależy od warunków każdego konta.

Wszystkie opcje `ccam usage` i `ccam next` znajdziesz w części [limity użycia](/docs/usage-limits).
