---
title: Panel
description: ccam w terminalu otwiera pełnoekranowy panel kont, grup, reguł folderów, rozszerzeń i sprawdzeń; każda zmiana najpierw pokazuje swój plan.
updated: 2026-10-21
order: 15
translation: machine
sourceHash: "27bd153b48b2"
---

`ccam` bez argumentów otwiera pełnoekranowy panel, gdy standardowe wejście i wyjście są terminalem; w przeciwnym razie wypisuje `ccam list`. Panel pokazuje każde konto z jego użyciem i grupami i robi to samo co polecenia: uruchamia, dodaje, loguje, zmienia grupy, ustawia reguły folderów, sprawdza i naprawia. Uruchomienie konta go zamyka; dodawanie konta lub logowanie oddaje na chwilę terminal i potem do niego wraca.

```console
 ccam v0.10.0 · 3 accounts · 2 groups     ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health
 Your Claude accounts and their usage limits; enter launches the selected one.

    # Account      5-hour       Weekly       Fable        Last  Groups
   1. personal ✔   ████▍░  72%  ██▌░░░  41%  █████▎  88%  -     default
 ❯ 2. work      ▸• █▌░░░░  25%  ██████ 100%  ▋░░░░░  10%  -     default team
   3. side         ✘ not logged in                        -     no group

 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

Wiersz pod nazwami kart mówi jednym zdaniem, do czego służy bieżąca karta. W szerszym terminalu karta Accounts pokazuje też adres e-mail każdego konta.

Pierwsza klatka jest rysowana wyłącznie z `config.toml`, w znacznie mniej niż 100 ms; użycie pojawia się, gdy tylko dane konto odpowie. Każdy limit ma własny pasek: okno 5-godzinne, tygodniowe okno dla wszystkich modeli i własne tygodniowe okno Fable; `-` oznacza, że konto nie ma takiego okna. Procent zmienia kolor na żółty od 80% i na czerwony przy 100%, a wartości z pamięci podręcznej po nieudanym zapytaniu są rysowane bladym kolorem.

## Karty

| Karta | Co pokazuje | Klawisze |
| --- | --- | --- |
| Accounts | numer, nazwa, adres e-mail, paski 5-godzinny, tygodniowy i Fable, ostatnie użycie, grupy; `✔` konto domyślne, `▸` konto tego folderu, `•` uruchomione | poniżej |
| Groups | każda grupa, co współdzieli i jej konta, oraz konta, które nie należą do żadnej grupy | `enter` edycja, `n` nowa, `e` zmiana nazwy, `x` usunięcie |
| Folders | reguły [`ccam folder`](/features/folder) i ta, która obowiązuje tutaj | `n` użyj tu konta, `enter` zmiana, `x` usunięcie |
| Extensions | to, co pokazuje [`ccam ext`](/features/ext), każde rozszerzenie raz, z grupami, z których pochodzi | `r` ponowne wczytanie |
| Health | to, co zgłasza [`ccam doctor`](/features/doctor), jako ramka dla każdego konta (poniżej) | `f` naprawa zaznaczonej ramki, `F` naprawa wszystkiego (`doctor --fix`), `↑`/`↓` przejście między ramkami |

`←`/`→` lub `Tab` przełączają karty, `?` wyświetla klawisze bieżącej karty, `r` sprawdza użycie ponownie z pominięciem krótkiej pamięci podręcznej, a `q` lub `Esc` zamyka panel.

## Klawisze na karcie Accounts

| Klawisz | Działanie |
| --- | --- |
| `Enter` | Uruchamia Claude Code jako zaznaczone konto |
| `1`–`9` | Przechodzi do konta o tym numerze; wpisz szybko `1` `2`, żeby wybrać 12. Uruchamia tylko `Enter` |
| `↑`/`↓`, `k`/`j` | Przesuwa zaznaczenie |
| `n` | Uruchamia konto, które wybrałby [`ccam next --auto`](/features/next) |
| `o` | Uruchamia z opcjami: `--continue`, `--resume`, inne argumenty albo próbne uruchomienie |
| `u` | Szczegóły użycia: okna 5-godzinne, tygodniowe i modeli z ich resetami oraz źródło liczb |
| `a` | Dodaje konto (kreator opisany niżej) |
| `l` | Loguje ponownie |
| `s` | Wybiera grupy konta |
| `d` | Ustawia je jako konto domyślne |
| `f` | Używa go w bieżącym folderze |
| `e` | Zmienia jego nazwę |
| `c` | Ustawia lub usuwa jego polecenie skrótu |
| `x` | Usuwa je (najpierw pyta, z zaznaczonym No) |

## Karta Health

Karta Health po otwarciu wykonuje sprawdzenia [`ccam doctor`](/features/doctor) i rysuje ten sam widok: pierwszy wiersz ze znacznikiem dla każdego konta, a pod nim ramkę dla każdego konta. Każdy wiersz w ramce to jeden sprawdzony obszar, najpierw `Folder`, `Login` i `Share links`, oznaczony `✔` w porządku, `ℹ` informacja, `⚠` ostrzeżenie lub `✘` błąd, więc od razu widać, które konto ma problem i gdzie. Ustalenia, które nie należą do żadnego konta, trafiają do ramki `General`.

```console
 ✘ personal   ⚠ work

 ╭─ ✘ personal · 3 problems ──────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                 │
 │ ℹ              unknown entries (Claude Code layout may have changed):      │
 │                newthing                                                    │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)      │
 │ ✔              login valid (me@example.com)                                │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json: │
 │                set to shared but is a real file, not a link                │
 ╰────────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ───────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                          │
 │ ⚠ Login        not logged in → `ccam login work`                           │
 │ ✔ Share links  share links OK                                              │
 ╰────────────────────────────────────────────────────────────────────────────╯

 4 problems; f repairs the selected box, F all of them (ccam doctor --fix).
```

`↑`/`↓` przechodzą od ramki do ramki, a zaznaczona ramka ma obramowanie w kolorze akcentu. `f` naprawia tylko zaznaczoną ramkę, tak jak zrobiłby to `ccam doctor --fix`: w ramce konta naprawia tylko to konto, a w ramce `General` ustalenia, które nie należą do żadnego konta (ścieżki wtyczek i skróty); pozostałe konta są tylko sprawdzane ponownie. Najpierw pyta, a przy ramce, w której nie ma nic do naprawy, wyświetla tylko `nothing to repair in <name>`. `F` naprawia wszystko naraz, jak `ccam doctor --fix`, a `r` sprawdza ponownie.

## Zmiana grup

Na karcie Groups `Enter` otwiera edytor: elementy po lewej, konta po prawej, `Tab` przełącza stronę, `Space` zaznacza lub odznacza. Element, który któreś z kont dostaje już z innej grupy, jest wyszarzony z podanym powodem, podobnie jak konto, które dostałoby jakiś element dwa razy; element może pochodzić tylko z jednej grupy. `Enter` pokazuje plan:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Najpierw padają pytania o pliki, które różnią się od kopii grupy (zachować wersję grupy, użyć mojej albo scalić). Potem zmiana wykonuje się w tle tym samym kodem co [`ccam group`](/features/share): najpierw zapisywana jest konfiguracja, dane każdego konta są przenoszone po kolei, a oryginały trafiają do `.ccam-backup/`. Zmiana, która przeniosłaby dane uruchomionego konta, jest odrzucana z podaniem nazwy tego konta.

## Dodawanie konta

`a` otwiera kreator, którego używa też `ccam add` w terminalu:

1. Nazwa konta.
2. Sposób logowania: subskrypcja Claude teraz (w przeglądarce), [klucz API](/features/api) albo subskrypcja Claude później.
3. Dla klucza API: skąd ccam odczytuje klucz przy każdym uruchomieniu (`env:NAME` lub `command:COMMAND`), punkt końcowy, sposób wysyłania klucza i model.
4. Grupy do dołączenia, z zaznaczonym `default`; zanim powstaną jakiekolwiek grupy: co będzie współdzielić grupa `default`.
5. Polecenie skrótu, jeśli go chcesz.
6. Podsumowanie do potwierdzenia.

`Esc` cofa o jeden krok, a w pierwszym kroku zamyka kreator.

## Kolory

Kolory odpowiadają rolom z terminalowego interfejsu Claude Code, więc ccam dobrze wygląda obok niego; nazwy grup dostają kolor wyliczony z nazwy. ccam pyta terminal o kolor tła i wybiera jasny lub ciemny motyw. Możesz ustawić go sam w `config.toml`:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

Motywy `-daltonized` są przeznaczone dla osób z zaburzeniem rozróżniania barwy czerwonej i zielonej, a motywy `-ansi` używają 16 kolorów samego terminala. `NO_COLOR` wyłącza kolory.

## Bezpieczeństwo i przypadki brzegowe

- Panel wywołuje te same funkcje co polecenia, więc obowiązuje każda zasada: plany przed przeniesieniem danych, odmowa dla uruchomionych kont, kopie zapasowe i tokeny logowania, które nigdy nie opuszczają konta.
- Panel nie ma trybu próbnego: `ccam --dry-run` bez nazwy wypisuje, jak uruchomiłoby się konto tego folderu, a `o` na koncie proponuje próbne uruchomienie tylko tego jednego startu.
- Terminal jest przywracany, zanim wystartuje Claude Code albo otworzy się logowanie.
- Gdy ekran jest za mały na okno dialogowe, najpierw znikają wiersze karty; wiersze są obcinane znakiem `…` zamiast zawijane.

## Powiązane

- [Współdzielenie przez grupy](/features/share)
- [Uruchamianie Claude Code](/features/launch)
- [Zarządzanie kontami](/features/accounts)
- [Limity użycia](/features/usage)
