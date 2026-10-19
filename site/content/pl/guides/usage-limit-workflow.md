---
title: Co zrobić, gdy Claude Code osiągnie limit użycia
description: Sprawdź, ile zapasu zostało na każdym koncie Claude Code, i kontynuuj tę samą rozmowę na innym koncie, gdy jedno osiągnie limit.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 767a731b912e
---

Gdy jedno konto osiągnie limit, kontynuuj tę samą rozmowę na innym swoim koncie: raz włącz współdzielenie historii rozmów, a potem uruchom `ccam next --continue`. Jeśli masz tylko jedno konto, poczekaj do czasu resetu, który pokazuje Claude Code.

## Zanim to nastąpi

Skonfiguruj konta tak, żeby przełączenie wymagało jednego polecenia:

```sh
ccam add personal
ccam add work
ccam share personal --history on
ccam share work --history on
ccam statusline install
```

Pasek stanu pokazuje teraz podczas pracy konto i jego 5-godzinne użycie, więc limit nie będzie zaskoczeniem.

## Gdy limit zostanie osiągnięty

Wyjdź z Claude Code (`/exit`), a następnie:

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

Claude Code uruchomi się jako `work` z najnowszą rozmową z bieżącego folderu, łącznie ze wszystkim, co zostało w niej dotąd powiedziane.

## O czym pamiętać

- Uruchamiaj `ccam next --continue` w tym samym folderze, w którym pracowałeś; rozmowy są przechowywane osobno dla każdego folderu projektu.
- Nie uruchamiaj dwóch kont na tej samej rozmowie jednocześnie. ccam ostrzeże cię, jeśli spróbujesz.
- Drugie konto też musi współdzielić historię rozmów, w tej samej grupie.
- To, czy wolno ci w ten sposób przełączać się między kontami, zależy od warunków każdego konta.

Wszystkie opcje `ccam usage` i `ccam next` znajdziesz w części [limity użycia](/docs/usage-limits).
