---
title: Befehlsreferenz
metaTitle: "ccam-Befehlsreferenz"
description: Alle ccam-Befehle und -Optionen, von ccam add und ccam <name> bis zu ccam usage, ccam next und ccam doctor.
updated: 2026-10-25
order: 3
translation: machine
sourceHash: "abba525cae48"
---

Setze `--dry-run` vor einen Befehl, um zu sehen, was er tun würde, ohne es zu tun, zum Beispiel `ccam --dry-run work`. Das funktioniert beim Starten von Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`) sowie mit `group`, `migrate`, `shortcut`, `update` und `doctor --fix`. Andere Befehle, die etwas ändern, lehnen die Option ab, statt zu laufen.

<!-- commands -->

## Claude Code starten

`ccam <name|number> [arguments]` startet Claude Code mit diesem Konto. ccam liest die Argumente nach dem Namen nicht; es gibt sie in derselben Reihenfolge an `claude` weiter, und Standardeingabe, Exit-Code und Signale verhalten sich genau wie bei `claude` selbst. Ein `--` nach dem Namen ist optional.

`ccam` mit Argumenten, aber ohne Kontonamen, etwa `ccam -p "summarise"`, verwendet das Konto, das für den aktuellen Ordner festgelegt ist, sonst das Standardkonto. Lege sie mit `ccam folder <name> [folder]` und `ccam use <name>` fest.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` vor dem Kontonamen ist für `-p`/`--print`-Durchläufe gedacht. Wenn das Konto sein Nutzungslimit erreicht, führt ccam denselben Befehl erneut mit dem angemeldeten Konto mit dem meisten Spielraum aus und gibt nur die Ausgabe des Durchlaufs aus, der fertig wurde. Siehe [Konten nach Ordner](/features/folder) und [Am Limit erneut ausführen](/features/next#run-again-at-the-limit).

**Team- und Enterprise-Plätze.** Die automatische Wahl (`ccam next`, `--fallback`, `ccam pick`, `n` im Dashboard und [`auto_switch`](/features/next#auto-switch), wenn du es einschaltest) wechselt auch zwischen Plätzen derselben Claude Team- oder Enterprise-Organisation. Im Hilfe-Center von Anthropic ([Team-Tarife](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); ebenso für [Enterprise-Tarife](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)) heißt es: „Each seat on a Team plan is assigned to one named individual, and only that person may use it“ (jeder Platz in einem Team-Tarif ist einer namentlich benannten Person zugewiesen, und nur diese Person darf ihn nutzen) und „Users can't route their work through another user's account to use that user's remaining usage“ (Benutzer dürfen ihre Arbeit nicht über das Konto eines anderen Benutzers leiten, um dessen verbleibendes Nutzungskontingent zu verwenden). Wenn du die automatische Wahl zwischen Plätzen nutzt, kann das daher gegen die Bedingungen von Anthropic verstoßen. Ob du sie dort nutzt, ist deine Entscheidung: Kläre das mit deiner Organisation und mit Anthropic ab, bevor du dich darauf verlässt. Für private Tarife (Pro, Max) gilt der [Hinweis zu den Bedingungen in den FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

Ein Kurzbefehl wie `ccwork` startet sein Konto auf dieselbe Weise: `ccwork --resume` macht dasselbe wie `ccam work --resume`. Siehe [Kurzbefehle](/docs/commands#shortcuts).

## Konten hinzufügen und entfernen

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam logout side              # asks first when MCP sign-ins go with the login
ccam rename work job
ccam remove side              # logs out, then the folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` ohne Ordner sucht nach `~/.claude` und `~/.claude-*` und fragt bei jedem nach; ein Ordner, bei dem es fehlschlägt, wird als `Skipped <path>: <reason>` gemeldet, und die übrigen werden registriert, und einer, dessen Name bereits einem Konto gehört, braucht `--name`, das zusammen mit dem Ordner angegeben wird. Importierte Ordner werden dort registriert, wo sie liegen; ccam verschiebt oder kopiert nie eine Anmeldung.

`ccam login --expired` meldet auch Konten an, deren Anmeldung Claude Code gelöscht hat. `ccam remove` meldet ein verwaltetes Konto ab, bevor sein Ordner in den Papierkorb kommt (nicht, wenn ein anderes Konto dieselbe Anmeldung hat), und `ccam remove` eines importierten Kontos verlässt seine Gruppen mit einer Kopie dessen, was sie teilen. Unter macOS behält `ccam rename` den Namen des Ordners bei, da der Schlüsselbund die Anmeldung unter dem Pfad des Ordners speichert.

`ccam add --api` fügt ein Konto hinzu, das einen API-Schlüssel statt einer Claude-Anmeldung verwendet; ccam liest den Schlüssel bei jedem Start und speichert ihn nie. Siehe [API-Schlüssel-Konten](/features/api).

`ccam rename`, `ccam remove` und die `ccam group`-Befehle ändern ein Konto nicht, solange Claude Code mit diesem Konto läuft, egal wie es gestartet wurde (auch Hintergrundsitzungen). Beende zuerst die Sitzung. Sie lehnen außerdem ab, solange eine Änderung von Gruppen nicht abgeschlossen ist: ``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` gibt jedes Konto mit Slot-Nummer, Kurzbefehl, E-Mail-Adresse und Ordner aus. Das Array `groups` listet die Gruppen des Kontos auf, und das Objekt `share` ordnet jedem Element (`mcp`, `settings`, `memory`, `extensions`, `history`) die Gruppe zu, aus der es kommt, oder `local`.

## Kurzbefehle

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Ein Kurzbefehl ist ein Befehl, den du selbst benennst und der Claude Code mit einem bestimmten Konto startet. Wie `ccam work` gibt er jedes Argument unverändert an Claude Code weiter, sodass `ccwork --resume` und `ccwork -p "hi"` wie erwartet funktionieren. Er ist ein Link auf ccam in `~/.local/bin`.

Kurzbefehle sind optional. Der Assistent von `ccam add` fragt, ob einer angelegt werden soll, `ccam add <name> --shortcut <command>` legt einen ohne Nachfrage an, und `ccam list` zeigt den Kurzbefehl jedes Kontos oder `-`, wenn es keinen hat.

Namen bestehen aus den Zeichen a-z, 0-9, `.`, `-` und `_`. Ein Name darf kein ccam-Befehl sein, nicht der Kurzbefehl eines anderen Kontos, nicht der Name eines anderen Programms in deinem `PATH` und keine Datei in `~/.local/bin`, die ccam nicht angelegt hat, etwa deine eigene Kopie von ccam oder ein Link darauf. Ein Kurzbefehl, der von Hand in `config.toml` bearbeitet wurde und gegen diese Regeln verstößt, `claude` oder `ccam` heißt oder zwei Konten zugewiesen ist, lässt das Laden von `config.toml` fehlschlagen, wobei das Konto genannt wird.

`ccam rename` behält den Namen des Kurzbefehls bei, und `ccam remove` löscht den Kurzbefehl zusammen mit dem Konto. `ccam shims install` erstellt die Links aller Konten neu, zum Beispiel nachdem du ccam verschoben hast; eine Datei, die einem Kurzbefehl im Weg steht, hält die anderen nicht auf.

## Teilen

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Eine Gruppe teilt Elemente zwischen ihren Konten; ein Konto kann in mehreren Gruppen sein, aber jedes Element kommt nur aus einer Gruppe. Siehe [Teilen mit Gruppen](/features/share).

## Eine vorhandene Einrichtung kopieren

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` kopiert `settings.json`, `CLAUDE.md`, Erweiterungen und Benutzer-MCP-Server aus einem Claude Code-Ordner in eine Gruppe, sodass jedes Konto, das diese Elemente teilt, sie bekommt. Der Ordner wird nur gelesen. Die Optionen `--settings`, `--memory`, `--extensions`, `--mcp` und `--history` kopieren nur die genannten Elemente. Ohne sie kannst du in einem Terminal in einer Checkliste auswählen, in der alles außer dem Gesprächsverlauf angehakt ist; anderswo wird alles außer dem Gesprächsverlauf kopiert. `-g` oder `--group` wählt die Gruppe, und `-y` oder `--yes` kopiert ohne Nachfrage und behält die Kopie der Gruppe, wenn eine Datei abweicht. MCP-Anmeldungen werden nicht kopiert, da eine kopierte Anmeldung die anderen Inhaber abmeldet; der Plan nennt die Konten, die sich mit `/mcp` anmelden. `--mcp` funktioniert auch, wenn `~/.claude` leer ist oder fehlt, `~/.claude.json` aber Server enthält. Siehe [Eine vorhandene Einrichtung kopieren](/docs/sharing#copy-an-existing-setup).

## MCP-Server

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` verschiebt die Server, die nur dieses Konto hat, in seine Gruppe, auch die, die es nur für einen Ordner hat (einfaches `claude mcp add`); diese werden dann in jedem Ordner geladen. Ein Server, dessen Namen die Gruppe bereits für eine andere Definition verwendet, wird übersprungen; mit `--force` ersetzt du die Kopie der Gruppe, die ccam vorher sichert.

## Nutzung und Wechseln

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam sends no token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --continue --brief    # a new conversation from a summary of the last one
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it (or by its name)
```

`next --auto`, `--fallback` und `ccam pick` wählen auf dieselbe Weise: zuerst unter `auto_max_usage` (`config.toml`, Standard 100), nach `auto_strategy` (`headroom`, `expiring` oder `balance`), unter Berücksichtigung des Wochenfensters von `--model`, nie ein Konto bei 100 % (ein Fenster, das Claude Code als voll gemeldet hat, zählt bis zu seinem Zurücksetzen als 100 %), nie eine Anmeldung per API-Schlüssel und innerhalb der Organisation des Kontos, von dem sie ausgehen (`--cross-org` oder `cross_org = true`, um weiter zu suchen). `--wait` nimmt seinen Wert nach `=` (`--wait=2h`); mit Leerzeichen geschrieben wird es mit einem Hinweis abgelehnt. Wenn keines gewählt werden kann, beenden sie sich mit Code 75 oder warten mit `--wait`, bis das erste Konto zurückgesetzt wird. `ccam next --continue` nimmt ein Gespräch auch zwischen Konten mit, die keinen Verlauf teilen: Claude Code kopiert es dann in das neue Konto. Mit `--brief` beginnt das neue Konto stattdessen ein neues Gespräch, ausgehend von einer Zusammenfassung des letzten, die ccam ohne Modell schreibt ([`--brief`](/features/next#brief)). Anders als `ccam <name>` liest `ccam next` seine eigenen Optionen, daher gehören die Optionen von Claude Code hinter `--`; ohne es bricht ccam mit `unknown flag` ab. In Claude Code übergibt `/ccam:switch [account]` aus dem Plugin ein Gespräch, und `ccam hook` ist das, was seine Hooks ausführen. Mit `auto_switch = true` in `config.toml` (Linux und macOS; standardmäßig aus) beendet der Limit-Hook außerdem eine interaktive Sitzung, die ccam gestartet hat, und führt sie als ein anderes Konto mit freiem Kontingent weiter, und zwar nur als ein gleichartiges: einen anderen Platz derselben Organisation oder einen anderen privaten Tarif ([`auto_switch`](/features/next#auto-switch)). Siehe [Zum nächsten Konto wechseln](/features/next), [Ein Konto für Skripte wählen](/features/pick) und [Hintergrundsitzungen aller Konten](/features/agents).

## Wartung

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` prüft die SHA-256-Prüfsumme des Releases, bevor es ccam ersetzt, und `--version v0.2.0` installiert ein bestimmtes Release. Ein aus dem `.deb`-Paket installiertes ccam aktualisierst du stattdessen mit apt. Siehe [Aktualisieren](/docs/install#update).

## Dateien

| Pfad | Inhalt |
| --- | --- |
| `~/.config/ccam/config.toml` | Konten, Gruppen und Einstellungen (`theme`, `usage_source`, `relay`, `auto_switch`, `cross_org` …). Einstellungen stehen am Anfang, über dem ersten `[[group]]`, `[[profile]]` oder `[[folder]]`: Eine Einstellung unter einer Tabelle gehört zu dieser Tabelle, daher nennt ccam sie bei jedem Start und in `ccam doctor`, und Befehle, die die Konfiguration speichern, lehnen ab, bis sie verschoben ist. Kommentare bleiben nicht erhalten, wenn ccam die Datei speichert |
| `~/.local/share/ccam/profiles/<name>/` | Das `CLAUDE_CONFIG_DIR` jedes Kontos |
| `~/.local/share/ccam/groups/<group>/` | Geteilte Einstellungen, MCP-Server, Erweiterungen und Verlauf |
| `~/.local/state/ccam/` | Nutzungscache und Limit-Einträge (`usage/<name>.json`, `usage/<name>.limit.json`), Zeitpunkte der letzten Verwendung, Sitzungsmarker und Sperren (`locks/`), die Kurzbefehle, die ccam angelegt hat (`shortcuts.json`), eine nicht abgeschlossene Änderung von Gruppen (`group-change.json`), die Notizen, die `/ccam:switch` und der Limit-Hook hinterlassen (`handoff/<account>/<session>.json`), wann `ccam pick` oder `auto_switch` jedes Konto zuletzt gewählt hat (`picks/<name>`; eine Wahl zählt 2 Minuten lang), und die Zusammenfassungen von `--brief`, 7 Tage lang aufbewahrt (`brief/<session>.md`) |
| `~/.local/bin/<shortcut>` | Kurzbefehle wie `ccwork`, Links auf ccam |

Mit den Variablen `XDG_CONFIG_HOME`, `XDG_DATA_HOME` und `XDG_STATE_HOME` lassen sich diese Ordner verlegen. Unter Windows ist `~` gleich `%USERPROFILE%`, und Kurzbefehle sind `<shortcut>.exe`.
