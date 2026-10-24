---
title: Kilka kont Claude Code na jednym komputerze
description: Uruchamiaj kilka kont Claude Code obok siebie dzięki CLAUDE_CONFIG_DIR, ręcznie z aliasami powłoki albo automatycznie z ccam, i przełączaj się między nimi.
updated: 2026-10-24
order: 1
translation: machine
sourceHash: "5463dec2f6eb"
---

Daj każdemu kontu własny folder konfiguracji. Claude Code odczytuje folder ze zmiennej środowiskowej `CLAUDE_CONFIG_DIR`, a każdy folder przechowuje własny login, ustawienia i historię, więc każde konto pozostaje zalogowane i możesz używać ich jednocześnie. Możesz to skonfigurować ręcznie za pomocą aliasu powłoki albo zostawić to ccam.

## Jak skonfigurować to ręcznie?

To podejście [opisuje dokumentacja Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Dodaj do `~/.bashrc` lub `~/.zshrc` jeden wiersz na każde dodatkowe konto:

```sh
alias claude-work='CLAUDE_CONFIG_DIR=~/.claude-work claude'
```

Otwórz nowy terminal i uruchom `claude-work`. Claude Code poprosi o zalogowanie, bo folder jest nowy. Od tej chwili `claude` to twoje pierwsze konto, a `claude-work` drugie. Dodaj po jednym aliasie na konto, każdy z własnym folderem.

W Windows bez WSL to samo robi funkcja PowerShell w twoim profilu (`notepad $PROFILE`):

```powershell
function claude-work { $env:CLAUDE_CONFIG_DIR = "$HOME\.claude-work"; try { claude @args } finally { Remove-Item Env:CLAUDE_CONFIG_DIR } }
```

## Czego brakuje w podejściu z aliasami?

Przy dwóch kontach to dobrze działa. Warto wiedzieć:

- Każdy folder jest na początku pusty: bez ustawień, bez `CLAUDE.md`, bez serwerów MCP, bez agentów i umiejętności. Kopiujesz lub linkujesz to, czego potrzebujesz, i sam dbasz o zgodność kopii.
- Aliasy istnieją tylko w interaktywnych powłokach. Skrypty, pliki Makefile i inne programy ich nie widzą.
- Serwery MCP są w `.claude.json`, obok loginu, więc nie da się ich po prostu podlinkować między folderami.
- Rozmowy zostają na koncie, na którym się zaczęły, więc nie możesz kontynuować ich na innym koncie.
- Nic nie pokazuje użycia wszystkich kont naraz.

## Czy zamiast tego można przełączać konto przez /login?

Tak, jeśli przełączasz się tylko od czasu do czasu. Wpisanie `/login` w Claude Code [zmienia konto](https://code.claude.com/docs/en/quickstart) folderu, w którym jesteś. Naraz zalogowane jest tylko jedno konto, więc przy każdej zmianie logujesz się ponownie i nie możesz uruchomić dwóch kont obok siebie. Osobne foldery rozwiązują oba problemy.

## Jak skonfigurować to z ccam?

ccam automatyzuje ten sam pomysł i dodaje brakujące elementy:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
ccam add personal
ccam add work --shortcut ccwork
```

Teraz `ccam personal` i `ccam work` uruchamiają oba konta, podobnie jak `ccwork`, skrót nadany kontu służbowemu. Ponieważ są to prawdziwe programy, a nie aliasy, działają też w skryptach. Uruchom samo `ccam`, żeby wybrać konto z listy, a `ccam use work`, żeby ustawić konto domyślne.

Ustawienia, `CLAUDE.md`, serwery MCP i rozszerzenia są współdzielone od początku. `ccam next --continue` przenosi rozmowę na drugie konto, a współdzielona historia rozmów, którą możesz włączyć, sprawia, że pozostaje ona jedną rozmową. `ccam usage` pokazuje 5-godzinne i tygodniowe użycie wszystkich kont naraz.

Jeśli już używasz `~/.claude`, uruchom `ccam import` zamiast dodawać to konto ponownie albo `ccam migrate`, aby skopiować jego ustawienia, serwery MCP i rozszerzenia do współdzielonych ustawień.

## Jak sprawdzić, które konto działa?

W Claude Code `/status` pokazuje konto, na które jesteś zalogowany. Z ccam uruchom raz `ccam statusline install`, a pasek stanu będzie stale pokazywał nazwę konta i jego użycie, na przykład `[work] · Opus · 5h 24% · 7d 41%`.

## Czy to działa w Windows i WSL?

Claude Code działa [natywnie w Windows i w WSL](https://code.claude.com/docs/en/setup). `CLAUDE_CONFIG_DIR` działa w obu, więc w natywnym Windows wystarczy alias lub funkcja PowerShell opisane wyżej. ccam też działa natywnie w Windows: zainstaluj go w PowerShell poleceniem `irm https://ccam.dev/install.ps1 | iex` albo w WSL poleceniem dla Linux. Skróty stają się plikami `ccwork.exe`, a współdzielenie plików takich jak `settings.json` wymaga Trybu dewelopera ([Instalacja ccam](/docs/install)).

## Rozwiązywanie problemów

**Claude Code za każdym razem prosi o zalogowanie.** Pierwsze uruchomienie w nowym folderze zawsze o to prosi. Jeśli prośba się powtarza, login wygasł lub został usunięty: uruchom `ccam login work` (albo raz uruchom alias i zaloguj się).

**Serwer MCP pojawia się tylko na jednym koncie.** Serwery dodane w sesji przez `claude mcp add` należą do tego konta, a bez `-s user` tylko do tego folderu. Przenieś serwer na wszystkie konta przez `ccam mcp import work <server>`.

**Ustawienia różnią się między kontami.** `ccam share work` pokazuje, co konto współdzieli i przez którą grupę; `ccam group join default work` sprawia, że współdzieli to, co współdzieli `default`, po pokazaniu, co się zmieni. `ccam doctor --fix` naprawia zepsute linki i zachowuje kopię zapasową każdego zmienianego pliku.

**Dwa konta otwierają tę samą rozmowę.** Nie pracuj nad jedną rozmową na dwóch kontach jednocześnie. ccam ostrzeże cię, jeśli spróbujesz.

## Co wybrać?

Dla dwóch kont bez niczego współdzielonego wystarczy alias. Gdy chcesz współdzielić ustawienia i serwery MCP, mieć więcej kont, widzieć użycie na pierwszy rzut oka albo kontynuować rozmowy, ccam oszczędza ci ręcznego pilnowania wszystkiego. [Porównanie z aliasami](/compare/aliases) zestawia różnice obok siebie, a [ccam i inne narzędzia do przełączania kont](/compare/account-switchers) omawia narzędzia, które zamiast tego podmieniają jeden login.

To, czy wolno ci używać kilku kont, zależy od warunków każdego konta. Sprawdź je, zanim zaczniesz.

**Miejsca w planach Team i Enterprise.** Automatyczny wybór konta (`ccam next`, `--fallback`, `ccam pick`, klawisz `n` w panelu) przełącza też między miejscami w tej samej organizacji Claude Team lub Enterprise. Centrum pomocy Anthropic podaje: „Each seat on a Team plan is assigned to one named individual, and only that person may use it” (każde miejsce w planie Team jest przypisane do jednej imiennie wskazanej osoby i tylko ta osoba może z niego korzystać) oraz „Users can't route their work through another user's account to use that user's remaining usage” (użytkownicy nie mogą kierować swojej pracy przez konto innego użytkownika, aby wykorzystać pozostały limit użycia tego użytkownika) ([plany Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); podobnie w przypadku [planów Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Dlatego korzystanie z automatycznego wyboru między miejscami może naruszać warunki Anthropic. Zanim zaczniesz na nim polegać, skonsultuj to ze swoją organizacją i z Anthropic. Planów prywatnych (Pro, Max) dotyczy uwaga powyżej.
