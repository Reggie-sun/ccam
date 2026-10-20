---
title: Aktualizacja ccam
description: ccam update zastępuje ccam najnowszym wydaniem z GitHub po sprawdzeniu SHA-256; --check tylko informuje, a --version instaluje wskazany tag wydania.
updated: 2026-10-20
order: 10
translation: machine
sourceHash: d91499b95d3d
---

`ccam update` pobiera wydanie ccam z GitHub, sprawdza je względem `checksums.txt` wydania i zastępuje działający plik binarny jedną zmianą nazwy. Używaj go, jeśli zainstalowałeś ccam skryptem instalacyjnym lub z archiwum. ccam nigdy sam nie sprawdza aktualizacji: tylko to polecenie łączy się z GitHub, więc każde inne polecenie działa tak samo bez sieci.

## Użycie

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, also an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

Gdy nie ma nowszego wydania, ccam wypisuje `ccam 0.4.0 is up to date.` i niczego nie zmienia. Z `--version` i wersją, której już używasz, wypisuje `ccam 0.4.0 is already installed.`

## Opcje

| Opcja | Działanie | Domyślnie |
| --- | --- | --- |
| `--check` | Tylko informuje, czy istnieje nowsze wydanie | wyłączone |
| `--version <tag>` | Instaluje ten tag wydania, na przykład `v0.3.0`, nowszy lub starszy od bieżącego | najnowsze |
| `--dry-run` | Respektowana: sprawdza najnowszy tag (chyba że `--version` podaje konkretny), wykonuje te same sprawdzenia pliku binarnego i wypisuje adresy pobrania i sumy kontrolnej oraz ścieżkę do zastąpienia, niczego nie pobierając | wyłączone |

`--check` i `--version` nie mogą być użyte razem. Tag musi mieć postać `vMAJOR.MINOR.PATCH`, z opcjonalnym przyrostkiem wersji przedpremierowej lub kompilacji; początkowe `v` można pominąć. Wszystko inne kończy się błędem, na przykład `--version takes a release tag such as v0.3.0, not "1.2"`.

Zmienna środowiskowa `CCAM_RELEASES_URL` zastępuje `https://github.com/jhyunwoo/ccam/releases`, na potrzeby serwera lustrzanego lub testów. Musi mieć ten sam układ: `<url>/latest` przekierowuje do `<url>/tag/<tag>`, a zasoby są pod `<url>/download/<tag>/<asset>`.

## Jak to działa

Kod znajduje się w `internal/update`, który importuje tylko `ccam update`. Ścieżka uruchamiania nigdy go nie wywołuje.

1. `cli.updateCmd` znajduje działający plik binarny przez `os.Executable` i rozwiązuje dowiązania symboliczne przez `filepath.EvalSymlinks`. Zastępowana jest właśnie ta rzeczywista ścieżka.
2. Bez `--version` `Updater.Latest` wysyła żądanie `<releases>/latest` bez podążania za przekierowaniem i odczytuje tag z nagłówka `Location` (`…/tag/v0.5.0`). Używa strony wydań zamiast API GitHub, więc nie obowiązuje limit API wynoszący 60 żądań na godzinę. Sprawdzenie ma 10 sekund.
3. `parseVersion` sprawdza tag ścisłym wzorcem. Tag trafia do adresów pobierania i do wypisywanego polecenia powłoki, więc nie przechodzi nic, co zawiera `/`, `%`, spacje lub znaki powłoki. `compareVersions` porównuje go z działającą wersją; wersja przedpremierowa jest sortowana przed swoim wydaniem.
4. `managed()` i `writable()` decydują, czy ccam w ogóle może zastąpić plik binarny (zobacz niżej).
5. `Updater.Fetch` pobiera jednocześnie `ccam_<os>_<arch>.tar.gz` (w Windows `ccam_windows_<arch>.zip`) i `checksums.txt`, z limitem 2 minut na oba. Każde pobranie jest ograniczone do 256 MB.
6. SHA-256 archiwum musi zgadzać się z jego wierszem w `checksums.txt`, zanim stanie się cokolwiek innego. Brakujący wiersz lub inna suma kontrolna zatrzymuje aktualizację z komunikatem `nothing was changed`.
7. `extract()` odczytuje archiwum tar i bierze pierwszy zwykły plik o nazwie `ccam` (`extractZip` bierze `ccam.exe` z archiwum zip). Foldery, linki i inne wpisy są ignorowane, a plik binarny też jest ograniczony do 256 MB.
8. `Updater.Install` zapisuje nowy plik binarny przez `fsutil.AtomicWrite`: plik tymczasowy `.ccam.tmp*` obok pliku binarnego, z bitami uprawnień starego pliku, zrzucony na dysk i przemianowany w miejsce starego.

Cały proces, od sprawdzenia tagu do zmiany nazwy:

<!-- motion:update-pipeline -->

Archiwum i `checksums.txt` są pobierane równolegle i nic nie jest zapisywane, dopóki suma kontrolna się nie zgadza. Zmiana nazwy jest atomowa, więc awaria lub pełny dysk zostawia albo stary plik binarny, albo nowy, nigdy połowę. Polecenia skrótów, takie jak `~/.local/bin/ccwork`, są dowiązaniami symbolicznymi do pliku binarnego ccam, więc uruchamiają nową wersję bez żadnych zmian.

W Windows pliku działającego programu nie można podmienić ani usunąć, można tylko zmienić jego nazwę. `Updater.Install` zapisuje nowy plik binarny jako `ccam.exe.new`, zmienia nazwę działającego `ccam.exe` na `ccam.exe.old`, a nowego pliku na `ccam.exe`. Jeśli ostatni krok się nie powiedzie, stary plik binarny wraca do dawnej nazwy. Następna aktualizacja usuwa `ccam.exe.old`. Skróty będące dowiązaniami twardymi zamiast symbolicznych ([Polecenia skrótów](/features/shortcuts)) są odnawiane zaraz potem, więc też uruchamiają nową wersję.

### Instalacje zarządzane przez menedżer pakietów

ccam zainstalowany z `.deb` należy do apt, a zastąpienie go za plecami apt zostawiłoby błędną bazę pakietów. Przed pobraniem ccam sprawdza, do kogo należy plik binarny:

<!-- motion:update-managed -->

1. Plik binarny w `/usr/bin`, `/bin`, `/usr/sbin` lub `/sbin` jest uznawany za zarządzany przez menedżer pakietów.
2. W innych miejscach, jeśli zainstalowano `dpkg` i `dpkg -S <path>` kończy się powodzeniem, plik też jest uznawany za zarządzany przez menedżer pakietów.
3. W przeciwnym razie ccam tworzy i usuwa tymczasowy plik `.ccam-update-*` obok pliku binarnego, żeby upewnić się, że może tam zapisywać.

Dla pliku binarnego zarządzanego przez menedżer pakietów ccam zatrzymuje się i wypisuje kroki dla apt:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

W systemach innych niż Linux wypisywany jest tylko wiersz ze skryptem instalacyjnym, w Windows: `irm https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.ps1 | iex`. Gdy do folderu nie da się zapisywać, ccam proponuje ponowne uruchomienie skryptu instalacyjnego, który instaluje ccam w `~/.local/bin`.

## Bezpieczeństwo i przypadki brzegowe

- Wersja deweloperska zgłasza swoją wersję jako `dev` i nie ma wydania do porównania, więc `ccam update` i `ccam update --check` ją odrzucają. `ccam update --version vX.Y.Z` instaluje na niej wydanie.
- Nieudane pobranie, brak zasobu dla twojej platformy lub niezgodna suma kontrolna zostawiają zainstalowany plik binarny bez zmian.
- Nieistniejące wydanie kończy się błędem `there is no release <tag>`. Gdy GitHub odmawia z kodem HTTP 403 lub 429, ccam informuje, że to prawdopodobnie limit zapytań, i prosi o ponowną próbę później.
- `--check` tylko odczytuje. Nie sprawdza, czy plik binarny jest zarządzany przez menedżer pakietów ani czy można do niego zapisywać.
- `--version` ze starszym tagiem to sposób na powrót do poprzedniej wersji.

## Powiązane

- [Instalacja](/docs/install)
- [Polecenia](/docs/commands)
- [FAQ](/docs/faq)
