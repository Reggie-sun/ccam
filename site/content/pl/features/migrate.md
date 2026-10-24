---
title: Kopiowanie istniejącej konfiguracji
description: Skopiuj ustawienia, CLAUDE.md, rozszerzenia i serwery MCP z istniejącego ~/.claude do grupy ccam poleceniem ccam migrate.
updated: 2026-10-23
order: 6
translation: machine
sourceHash: "7bca9efb6eb0"
---

Jeśli używałeś Claude Code przed ccam, twoja konfiguracja jest w `~/.claude` i `~/.claude.json`. `ccam migrate` kopiuje ją do grupy współdzielenia, więc każde konto, które współdzieli te elementy przez grupę, zaczyna z konfiguracją, którą już masz. Folder źródłowy jest tylko odczytywany, a zwykłe `claude` działa jak dotąd. Uruchom polecenie raz po pierwszym `ccam add`; gdy `~/.claude.json` zawiera serwery MCP, pierwsze `ccam add` wymienia je w swojej wskazówce.

## Użycie

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

W terminalu, bez opcji elementów i bez `--yes`, ccam najpierw pokazuje listę wyboru:

```console
 Copy into group default

 From /home/me/.claude
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Spacja (lub `x`) przełącza element pod kursorem, strzałki lub `j`/`k` przesuwają kursor, Enter pokazuje plan, a Esc, `q` lub Ctrl+C wypisuje `Cancelled.`. Elementy, dla których źródło nic nie zawiera, są oznaczone jako `not found`. Jeśli odznaczysz wszystko, ccam wypisuje `Nothing selected.` i kończy działanie.

Następnie ccam wypisuje plan i pyta o zgodę przed kopiowaniem:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

Grupa dociera tylko do własnych kont i tylko z elementami, które współdzieli, więc podsumowanie wymienia konta grupy. Gdy grupa nie współdzieli któregoś ze skopiowanych elementów albo niektóre konta do niej nie należą, ccam wypisuje polecenia `ccam group set <group> --add <items>` i `ccam group join <group> <accounts>`, które to zmieniają.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `[folder]` | Folder Claude Code, z którego kopiować | `~/.claude` |
| `--settings` | Kopiuje `settings.json` | wyłączone |
| `--memory` | Kopiuje globalny `CLAUDE.md` | wyłączone |
| `--extensions` | Kopiuje `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | wyłączone |
| `--mcp` | Kopiuje serwery MCP użytkownika (`mcpServers` w `.claude.json`) | wyłączone |
| `--history` | Kopiuje historię rozmów (`projects/`, `history.jsonl`) | wyłączone |
| `-g`, `--group string` | Grupa współdzielenia, do której kopiować | `default` |
| `-y`, `--yes` | Kopiuje bez pytania; przy konflikcie zachowywana jest kopia grupy | wyłączone |
| `--dry-run` (globalna) | Wypisuje plan i zatrzymuje się przed jakimkolwiek pytaniem lub kopiowaniem | wyłączone |

Bez opcji elementów w terminalu decyduje lista wyboru. Poza terminalem albo z `--yes` ccam kopiuje wszystko oprócz historii rozmów. Historia jest domyślnie pomijana, bo `projects/` może być duże.

`--group on` oznacza `default`. `--group off` i `--group local` są odrzucane komunikatem `--group takes a group name, not "off"`. `ccam migrate` respektuje `--dry-run`, a lista wyboru pojawia się także przy próbnym uruchomieniu, więc możesz podejrzeć wybór.

## Jak to działa

`internal/cli/migrate.go` wybiera elementy, a potem `migrate.Build` odczytuje źródło i grupę i zwraca plan. Nic nie jest zapisywane przed `Plan.Apply`, czyli przed twoim potwierdzeniem.

1. **Znalezienie źródła.** `checkSource` odrzuca ścieżkę, która nie istnieje, nie jest folderem albo nie zawiera żadnego znanego wpisu Claude Code. Dla samych serwerów MCP pusty lub brakujący `~/.claude` wystarcza, o ile `~/.claude.json` zawiera serwery. `migrate.ClaudeJSON` znajduje `.claude.json` należący do folderu: w jego wnętrzu albo, dla `~/.claude` bez takiego pliku w środku, `~/.claude.json` obok niego. Tam Claude Code trzyma ten plik, gdy `CLAUDE_CONFIG_DIR` nie jest ustawione.
2. **Przypisanie elementów do ścieżek.** Każda ścieżka pochodzi z `share.Layout`, tej samej listy dozwolonych ścieżek, której używają grupy. Serwery MCP są odczytywane z `.claude.json` przez `mcp.LocalServers`, które używa skanera `profile.TopLevelValue` zamiast parsować cały plik.
3. **Kopiowanie, nigdy przenoszenie.** Pliki trafiają do `~/.local/share/ccam/groups/<group>/` przez `fsutil.AtomicWrite` lub `fsutil.CopyFile`, które zapisuje plik tymczasowy obok celu i zmienia jego nazwę, więc przerwane kopiowanie nigdy nie zostawia uciętego pliku. Serwery MCP trafiają do `mcp.json` grupy pod jego blokadą: serwer dodany w międzyczasie przez inny ccam zostaje, a kopie zrobione dla [Remote Control](/guides/remote-control) są potem aktualizowane.

Każdy element przypisuje jedną ścieżkę źródłową do jednej ścieżki grupy:

<!-- motion:migrate-items -->

Dwa wiersze historii są przerywane, bo są kopiowane tylko z `--history`. `.credentials.json` i `.claude.json` znajdują się w `share.Never` i nigdy nie trafiają do grupy: `Build` zwraca błąd wewnętrzny, jeśli któryś z nich kiedykolwiek pojawi się w układzie, a przeglądanie folderu pomija je po nazwie.

`buildEntry` nadaje każdej ścieżce stan, który opisuje odpowiedni wiersz planu:

<!-- motion:migrate-states -->

- **Missing**: nie ma jej w źródle albo jest plikiem tam, gdzie oczekiwano folderu. Wymieniana raz na element jako `not found, skipped`.
- **Linked**: ścieżka źródłowa już prowadzi do kopii grupy, na przykład gdy źródłem jest folder profilu ccam.
- **Same**: grupa ma już identyczną zawartość.
- **New**: grupa jej nie ma, więc jest kopiowana w całości.
- **Merge**: folder po obu stronach. `copyMissing` dodaje tylko pliki, których grupa nie ma; pliki, które się różnią, zostają w wersji grupy, a plan podaje `group copies kept`. Dla `plugins/` rejestry `installed_plugins.json` i `known_marketplaces.json` są scalane, a bezwzględne ścieżki wtyczek przepisywane tak, żeby wskazywały folder grupy.
- **Conflict**: plik po obu stronach z różną zawartością. Wybierasz dla każdego pliku: `keep group`, `use mine`, a do tego `merge JSON` dla plików `.json` (głębokie scalenie, w którym przy sprzecznych wartościach wygrywa grupa) lub `merge by line` dla `history.jsonl`. Serwer MCP o innej definicji też jest konfliktem, z wyborem `keep group` lub `use mine`. Pierwszy wymieniony wybór jest domyślny po naciśnięciu Enter; `--yes` zawsze zachowuje kopię grupy.

Sama lista wyboru to `tui.MigratePicker`, zbudowany na komponencie `checklist`, którego używa `ccam add`; `migrate.Found` oznacza brakujące elementy na podstawie wywołań lstat i listy serwerów MCP.

### Logowania MCP

Serwery MCP logujące się przez OAuth przechowują swoje tokeny w pliku logowania, `.credentials.json` pod kluczem `mcpOAuth`, a nie na liście serwerów. ccam ich nie kopiuje. Skopiowane logowanie zawiera ten sam token odświeżania co oryginał, a serwer wymienia ten token przy pierwszym odnowieniu: każdy inny posiadacz, w tym źródłowy `~/.claude`, zostaje wtedy wylogowany z tego serwera. Dlatego wiersz planu `sign-ins` wymienia dla każdego kopiowanego serwera, do którego zalogowane jest źródło, konta grupy, które same logują się przez `/mcp`, a podsumowanie to powtarza. Wpis bez tokenu nie jest logowaniem i nie jest wymieniany, podobnie jak konto, które ma już własne logowanie, czy serwer, którego definicję zachowuje grupa. Żaden plik logowania nie jest zapisywany.

## Bezpieczeństwo i przypadki brzegowe

- Źródło jest tylko odczytywane. Żaden plik w nim nie jest przenoszony, linkowany ani zmieniany, a login źródła nigdy nie jest zapisywany, także na macOS, gdzie może być odczytywany z pęku kluczy.
- Zanim plik grupy zostanie zastąpiony, przez `use mine`, scalenie albo nowe serwery MCP w istniejącym `mcp.json`, kopia grupy jest zapisywana w `<group>/.ccam-backup/<time>/`. Nazwa folderu jest unikalna dla każdego uruchomienia, więc dwa uruchomienia w tej samej sekundzie nie nadpisują sobie kopii zapasowych. Podsumowanie wypisuje ścieżkę.
- Z `--yes` każdy konflikt zachowuje kopię grupy, a ccam podpowiada, jak zamiast tego wybrać: uruchom to samo polecenie bez `--yes`.
- Jeśli nic by się nie zmieniło, ccam wypisuje `Nothing to copy.`; jeśli każdy konflikt zachowuje kopię grupy, `Nothing to copy: the group keeps its own copies.`
- Jeśli kopiowanie przerwie się w połowie, ccam zatrzymuje się z komunikatem `The source folder was not changed. Group files replaced so far are saved in`, po którym następuje folder kopii zapasowej. Ponowne uruchomienie kopiuje pliki, które nie dotarły w całości.
- Gdy MCP współdzieli kilka grup, polecenia podpowiadane przez `ccam doctor`, `ccam ext` i `ccam mcp list` mają postać `ccam migrate --mcp --group <a|b>`.
- Brak źródła kończy się komunikatem `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Powiązane

- [Współdzielenie serwerów MCP, ustawień i historii](/docs/sharing)
- [Sprawdzanie i naprawa](/features/doctor)
- [Kilka kont Claude Code](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
