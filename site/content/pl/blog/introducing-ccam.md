---
title: Przedstawiamy ccam
description: Dlaczego powstał ccam, co daje osobom z kilkoma kontami Claude Code i jak nie przeszkadza w pracy.
updated: 2026-10-19
order: 1
translation: machine
sourceHash: b116c91bb8d4
---

ccam to małe narzędzie wiersza poleceń dla osób, które używają więcej niż jednego konta Claude Code na tym samym komputerze: konta prywatnego i służbowego albo kilku kont do jednego długiego projektu.

## Problem

Claude Code obsługuje kilka kont przez `CLAUDE_CONFIG_DIR`: jeden folder na konto. W praktyce każdy folder jest osobną wyspą. Ustawienia, `CLAUDE.md`, serwery MCP, agentów i umiejętności trzeba skonfigurować od nowa w każdym folderze i ręcznie utrzymywać w zgodzie. Rozmowy zostają tam, gdzie się zaczęły. A o tym, że konto osiągnęło limit użycia, dowiadujesz się dopiero przy próbie użycia.

## Co robi ccam

- `ccam add` tworzy folder, uruchamia logowanie i, jeśli chcesz, nadaje kontu skrót o nazwie wybranej przez ciebie, na przykład `ccwork`.
- `ccam work` i `ccwork` uruchamiają Claude Code jako to konto i przekazują każdy argument bez zmian.
- Serwery MCP, ustawienia, `CLAUDE.md` i rozszerzenia są współdzielone przez linki, element po elemencie i konto po koncie. Loginy nigdy nie są współdzielone.
- `ccam usage` pokazuje 5-godzinne i tygodniowe użycie wszystkich kont; `ccam next --continue` podejmuje rozmowę na następnym koncie.

## Bez wchodzenia w drogę

Narzędzie uruchamiane przed każdą sesją nie powinno kosztować niczego, co da się zauważyć. ccam odczytuje jeden plik konfiguracyjny, sprawdza każdy współdzielony link jednym wywołaniem systemowym, a potem zastępuje się procesem `claude`. Dodaje około 4 milisekund i nie działa w czasie twojej pracy.

Ostrożnie obchodzi się też z twoimi plikami. Każda zmiana współdzielenia najpierw wypisuje plan, a oryginały zachowuje w folderze kopii zapasowych. `ccam doctor` sprawdza linki, uprawnienia i loginy oraz naprawia to, co da się bezpiecznie naprawić.

## Wypróbuj

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add
```

ccam ma otwarty kod i jest nieoficjalny: nie jest powiązany z Anthropic ani przez Anthropic popierany. To, czy możesz używać kilku kont, zależy od warunków każdego konta. Opinie i tłumaczenia są mile widziane na [GitHub](https://github.com/jhyunwoo/ccam).
