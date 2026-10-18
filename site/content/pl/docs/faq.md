---
title: Najczęściej zadawane pytania
description: Odpowiedzi o ccam, kontach Claude Code, limitach użycia, współdzieleniu, prywatności i obsługiwanych systemach.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 12ad0faef2ce
---

## Czy ccam to oficjalne narzędzie Anthropic?

Nie. ccam jest nieoficjalnym projektem o otwartym kodzie i nie jest powiązany z Anthropic ani przez Anthropic popierany. Korzysta z ustawienia `CLAUDE_CONFIG_DIR`, które Claude Code dokumentuje do uruchamiania kilku kont obok siebie.

## Czy wolno mi używać kilku kont Claude Code?

To zależy od warunków każdego konta, a w przypadku kont służbowych także od zasad twojej organizacji. Sprawdź je, zanim zaczniesz używać kilku kont. ccam jedynie ułatwia przełączanie się między kontami, które już masz.

## Czy ccam pozwala obejść limity użycia?

Nie. Każde konto ma własne limity. ccam je pokazuje i pozwala kontynuować pracę na innym twoim koncie, gdy na jednym limit się wyczerpie.

## Jak używać dwóch kont Claude Code na jednym komputerze?

Daj każdemu kontu własny folder konfiguracji przez `CLAUDE_CONFIG_DIR`. ccam robi to za ciebie: uruchom `ccam add` dla każdego konta, a potem uruchamiaj je przez `claude1` i `claude2`. [Poradnik](/guides/multiple-claude-code-accounts) pokazuje też, jak zrobić to ręcznie.

## Jak kontynuować rozmowę na innym koncie?

Włącz współdzielenie historii rozmów między oboma kontami poleceniem `ccam share <name> --history on`, a potem uruchom `ccs next --continue`. Claude Code otworzy najnowszą rozmowę z bieżącego folderu na następnym koncie.

## Czy ccam widzi albo wysyła mój login?

ccam nigdy nie kopiuje ani nie współdzieli loginu. Odczytuje token dostępu konta tylko po to, żeby zapytać endpoint użycia Anthropic o użycie tego konta, gdy uruchamiasz `ccs usage`, `ccs next --auto` lub `ccam doctor`. ccam nie ma telemetrii i nie wysyła niczego nigdzie indziej.

## Gdzie są przechowywane moje konta?

W `~/.local/share/ccam/profiles/<name>/`, jeden folder na konto. Na macOS Claude Code przechowuje sam login w pęku kluczy, pod nazwą utworzoną na podstawie tego folderu.

## Co się stanie z moim istniejącym ~/.claude?

Nic, chyba że uruchomisz `ccam import`. Import rejestruje folder tam, gdzie jest, więc zwykłe `claude` i `ccs <name>` nadal używają tego samego konta i ustawień.

## Jakie systemy obsługuje ccam?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine i inne) oraz macOS. Na Windows używaj go w WSL.

## Dlaczego serwer MCP pojawia się tylko na jednym koncie?

Serwery dodane w trakcie sesji przez `claude mcp add` są zapisywane tylko dla tego konta. ccam poinformuje cię o tym przy następnym uruchomieniu konta; przenieś serwer do współdzielonej grupy poleceniem `ccam mcp import <account> <server>`.

## Czy ccam spowalnia Claude Code?

Uruchamianie Claude Code przez ccam dodaje około 4 milisekund, co zmierzono skryptem `scripts/perf-check.sh` na laptopie. ccam odczytuje jeden mały plik, sprawdza współdzielone linki, a potem zastępuje się procesem `claude`, więc nie działa w czasie twojej pracy.

## Jak odinstalować ccam?

Usuń `ccam`, `ccs` i linki `claude1`… z `~/.local/bin`, a potem `~/.config/ccam`, `~/.local/share/ccam` i `~/.local/state/ccam`, jeśli konta nie są ci już potrzebne. Zobacz [instalacja](/docs/install#uninstall).
