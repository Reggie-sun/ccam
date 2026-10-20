---
title: Prüfen und reparieren
description: Mit ccam doctor Freigabelinks, Dateiberechtigungen, Anmeldungen und die Claude Code-Version prüfen und mit ccam doctor --fix gefahrlos Reparierbares beheben.
updated: 2026-10-20
order: 9
translation: machine
sourceHash: a45ca25bb722
---

`ccam doctor` prüft deine ccam-Installation: die Claude Code-Binärdatei, Ordner, Anmeldung und Freigabelinks jedes Kontos, die geteilten MCP-Server, den Nutzungs-Endpunkt und deine Kurzbefehle. `ccam doctor --fix` repariert, was ohne Datenverlust möglich ist, und verschiebt jede Datei, die es ersetzt, vorher in einen Sicherungsordner. Führe es aus, wenn `ccam <name>` mit einem Problem bei den Freigabelinks anhält, nachdem du Dateien von Hand wiederhergestellt hast oder nach einem Update von Claude Code.

## Verwendung

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

Jeder Befund ist eine Zeile mit einem Zeichen: `✓` in Ordnung, `·` Information, `!` eine Warnung, `✗` ein Fehler. Befunde zu Konten beginnen mit dem Kontonamen.

```console
$ ccam doctor
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700)
! personal: .credentials.json permissions are 644 (should be 0600)
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp

`ccam doctor --fix` repairs what it can.
```

Mit `--fix` endet jeder reparierte Befund mit `→ fixed`:

```console
$ ccam doctor --fix
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700) → fixed
! personal: .credentials.json permissions are 644 (should be 0600) → fixed
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link → fixed
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp → fixed
```

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `--fix` | Reparieren, was sich gefahrlos reparieren lässt (geänderte Dateien werden gesichert) | aus |
| `--dry-run` (global) | Alle Prüfungen ausführen, aber nichts reparieren, auch nicht mit `--fix` | aus |

`ccam doctor` nimmt keine Argumente. Es beachtet `--dry-run`: `ccam --dry-run doctor --fix` gibt dieselben Befunde aus wie `ccam doctor` und ändert nichts.

Der Exit-Code ist `1`, wenn mindestens ein `✗`-Befund unrepariert bleibt, sonst `0`. Warnungen und Informationen lassen den Lauf nie fehlschlagen. Schlägt der Lauf ohne `--fix` fehl, endet ccam mit einer Leerzeile und `` `ccam doctor --fix` repairs what it can. ``. Eine Konfiguration, die sich nicht laden lässt, endet ebenfalls mit `1` und einem `ccam:`-Fehler.

## So funktioniert es

`internal/cli/doctor.go` lädt `~/.config/ccam/config.toml` und ruft `doctor.Run` in `internal/doctor/doctor.go` auf, das die Prüfungen in fester Reihenfolge ausführt und Befunde zurückgibt. `--fix` wird als `Options.Fix` durchgereicht und bei einem Probelauf auf false gesetzt.

<!-- motion:doctor-checks -->

Der mittlere Block läuft einmal pro Konto. Die Prüfungen und was `--fix` jeweils tut:

1. **Claude Code-Version** (`checkClaude`). Ermittelt die Binärdatei `claude` wie der Starter und führt `claude --version` aus. Beginnt die Version nicht mit `2.1.302`, der Version, gegen die ccam die Ordnerstruktur geprüft hat (`doctor.LayoutVersion`), fügt es `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed` hinzu. Eine fehlende Binärdatei oder ein fehlschlagendes `--version` ist ein Fehler. Keine Reparatur.
2. **Zugangsdaten in der Umgebung** (`checkAuthEnv`). Warnt, wenn `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` oder `ANTHROPIC_AUTH_TOKEN` gesetzt ist: Jedes aus dieser Shell gestartete Konto erbt die Variable, und Claude Code zieht sie der eigenen Anmeldung des Kontos vor. Keine Reparatur; entferne die Variable selbst.
3. **Ordner** (`checkProfile`). Ein fehlender Kontoordner ist ein Fehler und beendet die Prüfungen dieses Kontos. Berechtigungen mit gesetzten Gruppen- oder Fremd-Bits sind eine Warnung; `--fix` führt `chmod 0700` aus. Windows hat keine Berechtigungsbits, deshalb werden diese Prüfung und die `0600`-Prüfung unten dort übersprungen; den Ordner in deinem Benutzerprofil schützt seine Zugriffsliste.
4. **Anmeldedatei** (`checkCredentials`). `.credentials.json` wird mit `lstat` geprüft. Ein symbolischer Link ist ein Fehler, weil jedes Konto ein eigenes Token haben muss, und wird nicht repariert. Ein anderer Modus als `0600` ist eine Warnung; `--fix` führt `chmod 0600` aus. Dann wird die Anmeldung gelesen (unter macOS zuerst aus dem Schlüsselbund). Keine Anmeldung ist eine Warnung, die auf `ccam login <name>` verweist. Ein abgelaufenes Refresh-Token ist der Fehler `login expired`; ein abgelaufenes Zugriffstoken ist nur eine Information, da Claude Code es beim nächsten Start erneuert. Nichts davon wird repariert.
5. **Freigabelinks** (`share.Check`). Vergleicht jeden Pfad in `share.Layout` mit den Freigabeeinstellungen des Kontos und verwendet dabei nur `lstat` und `readlink`. Es meldet vier Arten: einen fehlenden Link, einen Link, der woandershin zeigt, eine echte Datei, wo ein Link hingehört, und einen übrig gebliebenen Link für ein Element, das inzwischen lokal ist. Ein fehlender Link ist eine Warnung, die anderen sind Fehler.
6. **Unbekannte Einträge** (`checkUnknownEntries`). Listet Namen auf der obersten Ebene des Kontoordners auf, die ccam nicht kennt (`share.Known`), als Hinweis darauf, dass sich die Ordnerstruktur von Claude Code geändert hat. Nur Information.
7. **Plugin-Pfade** (`checkPlugins`). Plugin-Register (`installed_plugins.json`, `known_marketplaces.json`) speichern absolute Pfade. Im `plugins/` jeder Gruppe gehen Pfade, die durch den Ordner eines Kontos führen, kaputt, wenn dieses Konto entfernt oder umbenannt wird. `--fix` schreibt sie mit `profile.RewritePluginPaths` auf den Gruppenordner um, als atomaren Schreibvorgang der Registerdatei.
8. **MCP-Server** (`checkMCP`). Lädt die `mcp.json` jeder Gruppe einmal. Eine Datei, die sich nicht lesen lässt, ist ein Fehler; ein Server, der `mcp.Validate` nicht besteht, ist eine Warnung. Keine Reparatur.
9. **Antwortformat der Nutzung** (`checkUsageFormat`). Fragt den Nutzungs-Endpunkt ohne Cache und mit dem üblichen Limit von 3 Sekunden für ein Konto nach dem anderen ab, bis eines antwortet. Lässt sich die Antwort nicht mehr auswerten, warnt es, dass `ccam usage` auf den Cache der Statuszeile zurückgreift. Das ist die einzige Prüfung, die das Netzwerk nutzt.
10. **Kurzbefehle** (`checkShims`). Nur wenn ein Konto einen Kurzbefehl hat. Fehlende Links in `~/.local/bin` sind eine Warnung; `--fix` ruft `shim.Install` auf, das nur fehlende Links anlegt und nie andere Dateien anfasst. Ein bin-Ordner, der nicht im `PATH` liegt, ist eine Warnung und wird nicht repariert.
11. **Ordner** (`checkFolders`). Nur wenn es Einträge von [`ccam folder`](/features/folder) gibt. Ein Eintrag, dessen Konto nicht existiert, ist ein Fehler, und ein Eintrag, dessen Ordner verschwunden ist, ist eine Warnung. Keines von beiden wird repariert, weil das Entfernen eines Eintrags ändern würde, welches Konto dort startet.

### Freigabelinks reparieren

Der Startpfad legt einen fehlenden Link sofort an, hält aber bei allem an, was Daten kosten könnte, und verweist hierher:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Bei jedem Problem mit Freigabelinks ruft `--fix` die Funktion `share.Build` mit den Freigabeeinstellungen des Kontos als aktuellem und als Zielzustand auf, dann `Plan.Apply`. `Build` leitet jeden Schritt neu aus dem ab, was auf der Festplatte liegt, derselbe Plan, den `ccam share` verwendet, daher folgt die Reparatur denselben Regeln wie das Einschalten des Teilens.

<!-- motion:doctor-repair -->

Hier ist `settings.json` eine echte Datei. Weil die Gruppe bereits eine andere `settings.json` hat, nimmt der Konflikt seine Standardwahl `KeepGroup`, sodass die Kopie der Gruppe bleibt. Die Datei des Kontos wandert nach `.ccam-backup/<time>/settings.json` im Kontoordner, und `os.Symlink` verlinkt den Pfad mit der Gruppe. Andere Fälle:

- Hat die Gruppe noch keine Kopie, wandert die Datei des Kontos in die Gruppe und wird ihr erster Inhalt.
- Bei einem Ordner wandern Dateien, die der Gruppe fehlen, in die Gruppe; der Rest kommt in die Sicherung.
- Bei `history.jsonl` ist die Standardwahl `merge by line`: die Zeilen der Gruppe, dann deine neuen. Die vorherige Kopie der Gruppe wird in `.ccam-backup/<time>/group/` gesichert.
- Ein Link auf eine andere Gruppe wird ersetzt. Ein Link auf einen Ordner außerhalb von ccam kommt in die Sicherung, und sein Inhalt wird zum Startinhalt der Gruppe, wenn diese leer ist.
- Ein übrig gebliebener Link für ein Element, das inzwischen lokal ist, wird entfernt, und das Konto erhält eine Kopie des Inhalts der Gruppe.

Abgesehen vom Entfernen von Links verschiebt jeder Schritt Daten in die Gruppe oder in den Sicherungsordner. Nichts wird gelöscht.

## Sicherheit und Sonderfälle

- Ein einfaches `ccam doctor` ändert nie etwas. Verwende in Skripten `--dry-run` mit `--fix`, um sicherzugehen.
- `doctor --fix` prüft nicht, ob ein Konto läuft. Beende Claude Code für ein Konto, bevor du seine Freigabelinks reparierst, weil die Sitzung weiter in diese Dateien schreibt.
- Schlägt eine Reparatur der Freigaben fehl, meldet ccam für dieses Konto `could not repair share links: <error>` und überspringt seine Prüfung auf unbekannte Einträge. Die anderen Konten werden trotzdem geprüft.
- Wenn `ccam share` mittendrin fehlschlägt, speichert es zuerst die Zieleinstellungen und rät dir, mit `ccam doctor --fix` abzuschließen.
- Eine per symbolischem Link verbundene `.credentials.json`, eine abgelaufene Anmeldung und Zugangsdaten in der Umgebung werden gemeldet, aber nie geändert.
- Umgeschriebene Plugin-Register werden atomar geschrieben, aber nicht nach `.ccam-backup/` kopiert.

## Verwandte Themen

- [MCP-Server, Einstellungen und Verlauf teilen](/docs/sharing)
- [Eine vorhandene Einrichtung kopieren](/features/migrate)
- [Nutzungslimits und Kontowechsel](/docs/usage-limits)
- [Befehle](/docs/commands)
