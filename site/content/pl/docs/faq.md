---
title: Najczęściej zadawane pytania
metaTitle: "FAQ ccam: kilka kont Claude Code"
description: Odpowiedzi o ccam, kontach Claude Code, limitach użycia, współdzieleniu, prywatności i obsługiwanych systemach.
updated: 2026-10-20
order: 6
translation: machine
sourceHash: "e7d4a2c877ac"
---

## Czy ccam to oficjalne narzędzie Anthropic?

Nie. ccam jest nieoficjalnym projektem o otwartym kodzie i nie jest powiązany z Anthropic ani przez Anthropic popierany. Korzysta z ustawienia `CLAUDE_CONFIG_DIR`, które [Claude Code dokumentuje](https://code.claude.com/docs/en/env-vars) do uruchamiania kilku kont obok siebie.

## Czy wolno mi używać kilku kont Claude Code?

To zależy od warunków każdego konta, a w przypadku kont służbowych także od zasad twojej organizacji. Sprawdź je, zanim zaczniesz używać kilku kont. ccam jedynie ułatwia przełączanie się między kontami, które już masz.

## Czy ccam pozwala obejść limity użycia?

Nie. Każde konto ma własne limity. ccam je pokazuje i pozwala kontynuować pracę na innym twoim koncie, gdy na jednym limit się wyczerpie.

## Jak używać dwóch kont Claude Code na jednym komputerze?

Daj każdemu kontu własny folder konfiguracji przez `CLAUDE_CONFIG_DIR`. ccam robi to za ciebie: uruchom `ccam add` dla każdego konta, a potem uruchamiaj je przez `ccam personal` i `ccam work` albo przez skróty o nazwach wybranych przez ciebie, na przykład `ccwork`. [Poradnik](/guides/multiple-claude-code-accounts) pokazuje też, jak zrobić to ręcznie.

## Jak kontynuować rozmowę na innym koncie?

Włącz współdzielenie historii rozmów między kontami poleceniem `ccam group set default --add history` (albo w osobnej grupie tylko dla nich: `ccam group add pair --items history --members work,side`), a potem uruchom `ccam next --continue`. Claude Code otworzy najnowszą rozmowę z bieżącego folderu na następnym koncie.

## Czy ccam widzi albo wysyła mój login?

ccam nigdy nie współdzieli twojego loginu Claude między kontami. Odczytuje token dostępu konta tylko po to, żeby zapytać endpoint użycia Anthropic o użycie tego konta, gdy uruchamiasz `ccam usage`, `ccam next --auto` lub `ccam doctor`. `ccam migrate` kopiuje zapisane logowania serwerów MCP z `~/.claude` do własnych plików logowania twoich kont na tym samym komputerze. ccam nie ma telemetrii i nie wysyła niczego nigdzie indziej.

## Jak zalogować się na zdalnym serwerze?

Uruchom `ccam add` lub `ccam login <name>` w sesji SSH jak zwykle. Claude Code czeka na logowanie na serwerze, więc przeglądarka na twoim komputerze kończy na stronie `http://localhost:<port>/callback?code=…`, która się nie wczytuje. Skopiuj ten adres z przeglądarki i wklej go do terminala: ccam przekaże go do Claude Code na serwerze, który dokończy logowanie. Adres zawiera jednorazowy kod, więc nikomu go nie udostępniaj.

## Gdzie są przechowywane moje konta?

W `~/.local/share/ccam/profiles/<name>/`, jeden folder na konto. Na macOS Claude Code przechowuje sam login w pęku kluczy, pod nazwą utworzoną na podstawie tego folderu.

## Co się stanie z moim istniejącym ~/.claude?

Nic, chyba że uruchomisz `ccam import` lub `ccam migrate`. Import rejestruje folder tam, gdzie jest, więc zwykłe `claude` i `ccam <name>` nadal używają tego samego konta i ustawień. `ccam migrate` kopiuje jego ustawienia, `CLAUDE.md`, rozszerzenia i serwery MCP do współdzielonych ustawień, a folder tylko odczytuje.

## Jakie systemy obsługuje ccam?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine i inne), macOS oraz Windows 10 i 11, natywnie lub w WSL. W Windows współdzielenie plików takich jak `settings.json` wymaga Trybu dewelopera; współdzielone foldery działają bez niego ([współdzielenie w Windows](/features/share#on-windows)).

## Dlaczego serwer MCP pojawia się tylko na jednym koncie?

Serwery dodane w trakcie sesji przez `claude mcp add` są zapisywane tylko dla tego konta. ccam poinformuje cię o tym przy następnym uruchomieniu konta; przenieś serwer do współdzielonej grupy poleceniem `ccam mcp import <account> <server>`.

## Czy ccam spowalnia Claude Code?

Uruchamianie Claude Code przez ccam dodaje około 4 milisekund, co zmierzono skryptem `scripts/perf-check.sh` na laptopie. ccam odczytuje jeden mały plik, sprawdza współdzielone linki, a potem zastępuje się procesem `claude`, więc nie działa w czasie twojej pracy. Windows nie potrafi zastąpić procesu, więc tam ccam uruchamia Claude Code i czeka na jego zakończenie, nie robiąc nic więcej.

## Jak zaktualizować ccam?

Uruchom `ccam update`. To polecenie pobiera najnowsze wydanie, sprawdza jego sumę kontrolną i podmienia plik `ccam`, a `ccam update --check` tylko informuje, czy jest dostępne nowsze wydanie. ccam nigdy nie aktualizuje się sam. Jeśli ccam pochodzi z pakietu `.deb`, zainstaluj zamiast tego nowy `.deb` przez apt. Szczegóły znajdziesz w części [aktualizacja](/docs/install#update).

## Jak odinstalować ccam?

Usuń skrót każdego konta poleceniem `ccam shortcut <account> --remove`, usuń `ccam` z `~/.local/bin`, a potem `~/.config/ccam`, `~/.local/share/ccam` i `~/.local/state/ccam`, jeśli konta nie są ci już potrzebne. Zobacz [instalacja](/docs/install#uninstall).
