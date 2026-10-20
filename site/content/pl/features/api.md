---
title: Konta z kluczem API
description: ccam add --api dodaje konto z kluczem API do API Anthropic lub bramy; klucz jest czytany przy starcie ze zmiennej lub polecenia i nigdy nie jest zapisywany.
updated: 2026-10-20
order: 13
translation: machine
sourceHash: de2c5cde2c48
---

Konto z kluczem API uruchamia Claude Code z kluczem API zamiast loginu subskrypcji Claude: z API Anthropic albo z endpointem obsługującym ten sam format Messages, takim jak brama LiteLLM lub OpenRouter. ccam nie przechowuje klucza. Odczytuje go przy każdym uruchomieniu konta, ze zmiennej środowiskowej albo z polecenia, np. menedżera haseł, i przekazuje go tylko temu procesowi Claude Code.

## Użycie

```sh
ccam add api --api --key-from env:ANTHROPIC_KEY_WORK
ccam add gw --api --base-url https://llm.example.com --auth bearer \
  --model claude-sonnet-5-5 --key-from 'command:pass show llm/gateway-key'
ccam add mac --api --key-from 'command:security find-generic-password -s anthropic-api -w'
ccam gw -p "summarise the open issues"
ccam list
```

```console
$ ccam list
   #  Name      Shortcut  Email                    Shared               Last used
*  1  personal  -         me@example.com           mcp settings md ext  2h ago
   2  gw        -         (API: llm.example.com)   mcp settings md ext  -
```

Konto z kluczem API współdzieli ustawienia, `CLAUDE.md`, rozszerzenia i serwery MCP jak każde inne konto i działa ze skrótami, `ccam folder` i `ccam ext`.

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--api` | Dodaje konto z kluczem API zamiast logowania | wyłączone |
| `--key-from env:NAME` | Odczytuje klucz ze zmiennej środowiskowej `NAME` | |
| `--key-from command:COMMAND` | Uruchamia `COMMAND` (przez `sh -c`, a w Windows przez `cmd /C`) i używa pierwszego wypisanego wiersza | |
| `--base-url <url>` | Endpoint (`ANTHROPIC_BASE_URL`) | API Anthropic |
| `--auth api-key\|bearer` | Wysyła klucz jako `ANTHROPIC_API_KEY` (`x-api-key`) albo jako `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, czego zwykle wymagają bramy) | `api-key` |
| `--model <model>` | Model (`ANTHROPIC_MODEL`) | domyślny model Claude Code |

Ustawienia są zapisywane w `config.toml`:

```toml
[[profile]]
name = 'gw'
slot = 2

[profile.api]
base_url = 'https://llm.example.com'
key_from = 'command:pass show llm/gateway-key'
auth = 'bearer'
model = 'claude-sonnet-5-5'
```

## Jak to działa

Przy każdym uruchomieniu konta, po zwykłych sprawdzeniach i przed startem Claude Code:

1. ccam odczytuje klucz: `env:` ze środowiska, `command:` przez uruchomienie polecenia z limitem 30 sekund. Polecenie może poprosić o hasło w terminalu; jego wyjście błędów też trafia do terminala. Pusty klucz przerywa uruchomienie z komunikatem, gdzie klucz był oczekiwany.
2. Usuwa każdą odziedziczoną zmienną, która decyduje, jak Claude Code się loguje lub jakiego dostawcy używa: `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_BASE_URL`, `ANTHROPIC_MODEL`, `CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX` i `CLAUDE_CODE_USE_FOUNDRY`. Klucz ustawiony w twojej powłoce do czegoś innego nie trafi do tego konta, a login subskrypcji nie może wygrać z jego kluczem.
3. Ustawia `ANTHROPIC_BASE_URL` i `ANTHROPIC_MODEL`, jeśli je podano, oraz klucz jako `ANTHROPIC_API_KEY` lub `ANTHROPIC_AUTH_TOKEN`.

Konta z loginem subskrypcji zachowują odziedziczone środowisko, jak wcześniej. W sesji interaktywnej Claude Code może raz zapytać, czy użyć klucza znalezionego w środowisku.

## Bezpieczeństwo i przypadki brzegowe

- Klucz nigdy nie jest zapisywany w `config.toml`, w pamięci podręcznej ani w logu, a `ccam --dry-run gw` go nie odczytuje: próbne uruchomienie nie wykonuje źródła `command:`.
- `ccam add --api` od razu sprawdza źródło `env:` i ostrzega, jeśli zmienna jest pusta; źródło `command:` jest uruchamiane po raz pierwszy przy starcie konta.
- `ccam login` i `ccam logout` odrzucają konto z kluczem API, bo nie ma ono loginu. `ccam doctor` pokazuje jego endpoint i źródło klucza oraz ostrzega, gdy zmienna `env:` jest pusta.
- `ccam usage` pokazuje dla niego `API`. `ccam next --auto`, `--fallback` i `ccam pick` wybierają tylko spośród kont z subskrypcją, bo klucz API nie ma 5-godzinnego ani tygodniowego limitu do porównania.
- Współdzielony `settings.json` z blokiem `env` lub `apiKeyHelper` dotyczy każdego konta w grupie, także kont z kluczem API. Trzymaj ustawienia dostawcy poza współdzielonymi ustawieniami.

## Powiązane

- [Zarządzanie kontami](/features/accounts)
- [Uruchamianie Claude Code](/features/launch)
- [Opis poleceń](/docs/commands)
