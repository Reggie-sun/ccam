---
title: Konta według folderu
description: ccam folder sprawia, że folder i wszystko w nim używa jednego konta, gdy uruchamiasz Claude Code bez nazwy konta, na przykład work dla kodu firmy.
updated: 2026-10-20
order: 11
translation: machine
sourceHash: 5071b04a194d
---

`ccam folder work ~/code/company` sprawia, że każde uruchomienie bez nazwy konta w `~/code/company` i w każdym folderze poniżej uruchamia Claude Code jako `work`. Poza nim konto domyślne uruchamia się jak wcześniej. Podanie nazwy konta, przez `ccam personal` albo skrót taki jak `ccpersonal`, zawsze uruchamia to konto.

## Użycie

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

Te uruchomienia używają konta folderu:

- `ccam -p "…"` i każde inne uruchomienie, którego pierwszy argument jest opcją dla Claude Code
- `ccam run` bez nazwy oraz `ccam --dry-run`
- `ccam --fallback -p "…"` przy pierwszej próbie
- selektor, czyli samo `ccam`, który startuje z kursorem na koncie folderu

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `<account>` | Nazwa lub numer konta, którego używać | |
| `[folder]` | Folder, ścieżka względna lub bezwzględna | bieżący folder |
| `--remove` | Usuwa wpis dokładnie dla tego folderu | wyłączone |
| `--dry-run` | Respektowana: wypisuje zmianę bez zapisywania jej | wyłączone |

## Jak to działa

Każdy folder to jeden wpis `[[folder]]` w `~/.config/ccam/config.toml`, czyli w pliku, który każde uruchomienie i tak czyta:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` zamienia folder na ścieżkę bezwzględną z rozwiązanymi dowiązaniami symbolicznymi i ją zapisuje. Ponowne ustawienie folderu zastępuje jego konto.
2. Uruchomienie bez nazwy wywołuje `launcher.HereProfile`. Gdy nie ma wpisów, od razu bierze konto domyślne.
3. W przeciwnym razie porównuje katalog roboczy z każdym `path` jako ciągi znaków: sam folder albo folder poniżej, nigdy folder, który tylko zaczyna się od tych samych liter (`/code/company` nie obejmuje `/code/companyx`). Wygrywa najgłębszy wpis, więc `~/code` może być `personal`, a `~/code/company` `work`.
4. Dopiero gdy katalog roboczy w obecnej postaci do niczego nie pasuje, ccam rozwiązuje jego dowiązania symboliczne i porównuje jeszcze raz.
5. W Windows porównanie ignoruje wielkość liter, tak jak sam Windows.

Poza `config.toml` nie jest czytany żaden plik, więc uruchomienie w folderze nie kosztuje mierzalnego czasu (4,0 ms z wpisami i bez, `hyperfine`, 300 przebiegów).

## Bezpieczeństwo i przypadki brzegowe

- `ccam rename` przenosi wpisy folderów na nową nazwę, a `ccam remove` usuwa wpisy usuniętego konta.
- Wpis, którego konto już nie istnieje, zatrzymuje uruchomienie bez nazwy w tym folderze: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` zgłasza takie wpisy oraz wpisy, których folderu już nie ma.
- `ccam folder --remove` w folderze, który dostaje konto od folderu wyżej, wskazuje ten folder zamiast cokolwiek usuwać.
- `ccam next` wybiera według rotacji lub użycia i nie patrzy na foldery.

## Powiązane

- [Uruchamianie Claude Code](/features/launch)
- [Polecenia skrótów](/features/shortcuts)
- [Zarządzanie kontami](/features/accounts)
- [Opis poleceń](/docs/commands)
