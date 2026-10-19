---
title: "Przełączanie kont Claude Code: ccam i inne narzędzia"
description: Porównanie ccam z claude-swap, CCS, clauth, CAAM, cux i innymi narzędziami do przełączania kont Claude Code oraz które podejście pasuje do czego.
updated: 2026-10-19
order: 2
translation: machine
sourceHash: dc907c88e2bf
---

Narzędzia do przełączania kont Claude Code stosują jedno z dwóch podejść. Większość podmienia login w jednym folderze `~/.claude`, więc naraz aktywne jest tylko jedno konto. Inne, w tym ccam, dają każdemu kontu własny folder przez `CLAUDE_CONFIG_DIR`, więc każde konto pozostaje zalogowane i kilka może działać jednocześnie. ccam dodaje współdzielone ustawienia i serwery MCP oraz kontynuuje rozmowę w folderze innego konta.

## Podmiana loginu czy osobne foldery?

**Podmiana loginu** zachowuje jeden zestaw ustawień, serwerów MCP i historii, bo jest tylko jeden folder. Przełączanie jest szybkie, a działająca sesja może przejąć nowy login. Ale naraz aktywne jest tylko jedno konto, we wszystkich terminalach i edytorach jednocześnie.

**Osobne foldery** to metoda, którą [opisuje dokumentacja Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Każde konto pozostaje zalogowane i możesz uruchomić sesję służbową i prywatną obok siebie. Ceną jest to, że każdy folder zaczyna z własnymi ustawieniami, serwerami MCP i historią, chyba że jakieś narzędzie je współdzieli.

## Porównanie narzędzi

Na podstawie dokumentacji każdego projektu z 2026-10-19. „Nieudokumentowane” oznacza, że projekt o tym nie wspomina, a nie że brakuje funkcji.

| Narzędzie | Podejście | Kilka kont jednocześnie | Użycie 5-godzinne i tygodniowe | Po osiągnięciu limitu | Platformy |
| --- | --- | --- | --- | --- | --- |
| ccam | Osobne foldery | Tak | Tak, `ccam usage` | Uruchamiasz `ccam next --auto --continue`, żeby kontynuować rozmowę na koncie z największym zapasem | Linux, macOS, WSL |
| [claude-swap](https://github.com/realiti4/claude-swap) | Podmienia login; osobne foldery z `cswap run` | Tak, z `cswap run` | Tak | `cswap auto` podmienia login po przekroczeniu progu; działająca sesja trwa dalej | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Osobne foldery | Tak | Częściowo, w aplikacji na pasku menu macOS | Pule kont przekazują żądania do następnego konta przez lokalne proxy | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Podmienia login; osobne foldery z `clauth start` | Tak, z `clauth start` | Tak | Przełącza według łańcucha zapasowych kont | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Podmienia login lub izolowane profile | Tak, z izolowanymi profilami | Tak | `caam run` ponawia nieinteraktywne polecenia na innym koncie | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Podmienia login i opakowuje `claude` | Nieudokumentowane | Tak | Podmienia login i wznawia tę samą rozmowę | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Podmienia login | Nie | Tak | Nieudokumentowane | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Podmienia login | Nie | Tak | Opcjonalne przełączenie przy 100% | macOS |

ccswap to fork claude-swap, który zarządza też loginami Codex. Kilka z tych narzędzi obsługuje też innych agentów kodujących, takich jak Codex czy Gemini.

## Na czym skupia się ccam

- Każde konto pozostaje zalogowane. Przełączenie nigdy nie wylogowuje żadnego z nich, a `ccam work` i `ccam personal` mogą działać jednocześnie w dwóch terminalach.
- Serwery MCP, `settings.json`, `CLAUDE.md`, agenci, umiejętności i wtyczki są współdzielone między folderami przez linki, element po elemencie i konto po koncie. Loginy nigdy nie są współdzielone.
- Przy współdzielonej historii jedno polecenie, `ccam next --auto --continue`, otwiera tę samą rozmowę w folderze konta z największym zapasem.
- ccam uruchamia Claude Code i zastępuje się nim, więc nic nie działa w tle. Nie śledzi twojego użycia ani nie przełącza kont sam; to ty decydujesz kiedy.

## W czym inne narzędzia robią więcej

- claude-swap, CCS, clauth, CAAM i cux działają w Windows bez WSL. ccam nie.
- claude-swap, clauth i cux przełączają się same, gdy konto przekroczy próg użycia. ccam zostawia tę decyzję tobie.
- CCSwitcher, Claude Switcher i claude-swap pokazują użycie na pasku menu macOS.
- CCS, CAAM, ccswap i Claude Switcher zarządzają też Codex i innymi agentami kodującymi.

## Co wybrać?

Wybierz narzędzie podmieniające login, jeśli używasz jednego konta naraz, chcesz, żeby przełączanie działo się automatycznie, albo pracujesz w Windows bez WSL. Wybierz ccam, jeśli chcesz mieć konta obok siebie, jedną wspólną konfigurację i tę samą rozmowę na następnym koncie. Jeśli nie potrzebujesz niczego współdzielić, wystarczą [zwykłe aliasy powłoki](/compare/aliases).

To, czy wolno ci używać kilku kont, zależy od warunków każdego konta. Sprawdź je, zanim zaczniesz.
