---
title: Co zrobić, gdy Claude Code osiągnie limit użycia
description: Sprawdź, ile zapasu zostało na każdym koncie Claude Code, i kontynuuj tę samą rozmowę na innym koncie, gdy jedno osiągnie limit.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 64fe0b561bd1
---

Gdy jedno konto osiągnie limit, kontynuuj tę samą rozmowę na innym swoim koncie: raz włącz współdzielenie historii rozmów, a potem uruchom `ccs next --continue`. Jeśli masz tylko jedno konto, poczekaj do czasu resetu, który pokazuje Claude Code.

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
ccs usage
```

```console
 1 personal  ████████ 100% 22분        █████░░░  68% 1일 9시간   방금
 2 work      █░░░░░░░   8% 4시간 12분  ██░░░░░░  19% 5일 2시간   방금
```

Otwórz rozmowę na koncie z największym zapasem:

```sh
ccs next --auto --continue
```

Claude Code uruchomi się jako `work` z najnowszą rozmową z bieżącego folderu, łącznie ze wszystkim, co zostało w niej dotąd powiedziane.

## O czym pamiętać

- Uruchamiaj `ccs next --continue` w tym samym folderze, w którym pracowałeś; rozmowy są przechowywane osobno dla każdego folderu projektu.
- Nie uruchamiaj dwóch kont na tej samej rozmowie jednocześnie. ccam ostrzeże cię, jeśli spróbujesz.
- Drugie konto też musi współdzielić historię rozmów, w tej samej grupie.
- To, czy wolno ci w ten sposób przełączać się między kontami, zależy od warunków każdego konta.

Wszystkie opcje `ccs usage` i `ccs next` znajdziesz w części [limity użycia](/docs/usage-limits).
