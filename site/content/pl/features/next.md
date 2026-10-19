---
title: Przełączanie na następne konto
description: ccam next otwiera Claude Code jako konto następne po ostatnio używanym; --auto wybiera najniższe użycie 5-godzinne, a --continue wznawia rozmowę.
updated: 2026-10-19
order: 8
translation: machine
sourceHash: 92352f532e74
---

`ccam next` otwiera Claude Code jako inne konto bez podawania jego nazwy. Używaj go, gdy konto osiągnie limit użycia: samo `ccam next` wybiera konto następne po ostatnio używanym, `--auto` wybiera konto z największym zapasem w 5-godzinnym oknie, a `--continue` kontynuuje ostatnią rozmowę na nowym koncie.

## Użycie

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next
→ personal → work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude
```

ccam wypisuje przełączenie jako `→ <from> → <to>` na standardowe wyjście błędów, a potem uruchamia Claude Code jako nowe konto tak samo jak `ccam <name>`. `<from>` to konto, którego używałeś ostatnio. Argumenty po `--` trafiają do `claude` bez zmian.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--auto` | Wybiera zalogowane konto z najniższym użyciem 5-godzinnym, pomijając konta na poziomie 100% | wyłączone |
| `-c`, `--continue` | Bierze pod uwagę tylko konta, które współdzielą historię rozmów z ostatnio używanym kontem, i uruchamia `claude --continue` | wyłączone |
| `-- <claude args>` | Argumenty przekazywane do `claude` po argumentach ccam | brak |
| `--dry-run` | Respektowana: wypisuje `→ <from> → <to>` i wiersz polecenia, który zostałby uruchomiony, bez startu Claude Code i bez zapisywania uruchomienia | wyłączone |

`--auto` i `--continue` działają razem. Z `--continue` opcja `--continue` jest umieszczana przed twoimi argumentami.

## Jak to działa

Każde uruchomienie przez ccam zapisuje czas w `~/.local/state/ccam/last-used/<name>`. `ccam next` odczytuje te pliki, sprawdza, które konta są zalogowane, i wybiera.

1. `cli.candidates` tworzy jeden `next.Candidate` na konto z `config.toml`: konto, czas jego ostatniego uruchomienia (`launcher.LastUsed`) i informację, czy ma login (`creds.Has`: czytelny `.credentials.json` albo, na macOS, wpis w pęku kluczy).
2. Bez `--auto` `next.Rotate` sortuje konta według numeru slotu i znajduje to uruchomione najpóźniej. To konto, z którego się przełączasz. Jeśli żadne konto nie było jeszcze uruchamiane, jest to najniższy slot.
3. Od niego idzie dalej w kolejności slotów, wracając od ostatniego slotu do pierwszego, i zwraca pierwsze zalogowane konto.
4. `launcher.LaunchWith` przygotowuje i uruchamia wybrane konto dokładnie tak jak `ccam <name>`: sprawdza linki konta, stosuje współdzielenie MCP, rejestruje sesję i zastępuje ccam przez `claude` za pomocą `syscall.Exec`.

Przejście zaczyna się od ostatnio używanego konta i zatrzymuje na pierwszym zalogowanym koncie po nim:

<!-- motion:next-rotate -->

`work` uruchomiono najpóźniej, więc przejście zaczyna się po nim. `side` nie ma loginu i zostaje pominięte, a wybrane zostaje `lab`. Gdyby `lab` też było wylogowane, przejście wróciłoby na początek do `personal` (`% len` na diagramie). Rotacja nigdy nie zwraca konta, od którego zaczęła: gdy nie ma innego zalogowanego konta, kończy się błędem.

### --auto

`--auto` decyduje według użycia zamiast kolejności slotów. `cli.pickAuto` zostawia zalogowane konta, pyta o użycie wszystkich naraz tym samym `Fetcher.FetchAll`, którego używa `ccam usage`, i pozwala wybrać `usage.Best()`:

<!-- motion:next-auto -->

`usage.Best()` wybiera najniższy procent 5-godzinny i pomija konta na poziomie 100% lub z nieznanym użyciem 5-godzinnym. Zapytanie działa według zasad `ccam usage`: liczby młodsze niż 60 sekund są używane ponownie, każde konto ma 3 sekundy, a konto, które nie odpowiada lub którego login wygasł, jest oceniane na podstawie zapisanych liczb. Zapisane okno, którego czas resetu minął, liczy się jako 0%, więc konto, którego okno w międzyczasie się zresetowało, nie zostaje pominięte. `--auto` nigdy nie odnawia tokenów. W odróżnieniu od zwykłej rotacji `--auto` może wybrać ostatnio używane konto, jeśli nadal ma największy zapas.

### --continue

`--continue` w Claude Code otwiera najnowszą rozmowę w bieżącym folderze, co działa tylko wtedy, gdy nowe konto widzi tę rozmowę. Dlatego z `--continue` ccam bierze pod uwagę tylko konta z tej samej grupy `history` co ostatnio używane konto, a potem umieszcza `--continue` przed twoimi argumentami. Konta z historią `local` nie współdzielą jej z nikim, więc `--continue` z takiego konta zawsze się nie powiedzie. Najpierw włącz współdzielenie przez `ccam share <name> --history on` (zobacz [współdzielenie](/docs/sharing)).

## Bezpieczeństwo i przypadki brzegowe

- Konta bez loginu nigdy nie są wybierane. Gdy nie ma innego zalogowanego konta, samo `ccam next` kończy się komunikatem `no logged-in account other than <name>`.
- Z `--continue`, gdy żadne konto nie współdzieli historii: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- Z `--auto`, gdy nie ma zalogowanego konta: `no logged-in account to pick`. Gdy każde konto jest pełne lub ma nieznane użycie: `` no account with known usage and headroom (see `ccam usage`) ``.
- `--auto` wysyła zapytania sieciowe nawet z `--dry-run`. Tylko odczytuje i nigdy nie zapisuje plików z tokenami.
- Wybrane konto startuje tą samą ścieżką co `ccam <name>`, więc obowiązują te same sprawdzenia, w tym ostrzeżenie, gdy w tym folderze działa już inny Claude Code współdzielący tę samą historię.
- `last-used` jest zapisywany tylko wtedy, gdy Claude Code rzeczywiście wystartuje. `--dry-run` go nie zmienia, więc próbne uruchomienie nie przesuwa rotacji.

## Powiązane

- [Limity użycia](/features/usage)
- [Limity użycia i przełączanie kont](/docs/usage-limits)
- [Współdzielenie serwerów MCP, ustawień i historii](/docs/sharing)
- [Praca z limitem użycia](/guides/usage-limit-workflow)
