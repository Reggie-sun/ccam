---
title: Prüfen und reparieren
description: Mit ccam doctor Freigabelinks, Dateiberechtigungen, Anmeldungen und die Claude Code-Version prüfen und mit ccam doctor --fix gefahrlos Reparierbares beheben.
updated: 2026-10-25
order: 9
translation: machine
sourceHash: "0f5ae4ab9a54"
---

`ccam doctor` prüft deine ccam-Installation: die Claude Code-Binärdatei, Ordner, Anmeldung und Freigabelinks jedes Kontos, die geteilten MCP-Server, den Nutzungs-Endpunkt und deine Kurzbefehle. `ccam doctor --fix` repariert, was ohne Datenverlust möglich ist, und verschiebt jede Datei, die es ersetzt, vorher in einen Sicherungsordner. Führe es aus, wenn `ccam <name>` mit einem Problem bei den Freigabelinks anhält, nachdem du Dateien von Hand wiederhergestellt hast oder nach einem Update von Claude Code.

## Verwendung

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

Die erste Zeile markiert alle Konten auf einmal, sodass du siehst, wo die Probleme liegen, bevor du irgendetwas anderes liest. Darunter bekommt jedes Konto einen Kasten, und Befunde, die zu keinem Konto gehören (Claude Code, Umgebungsvariablen, Plugins, MCP-Server, Kurzbefehle, Ordnerregeln, Gruppen und Einstellungen, `auto_switch`), kommen in einen Kasten `General`. Der Titel des Kastens zeigt das schlimmste Zeichen des Kontos und wie viele Probleme es hat. Jede Zeile darin ist ein geprüfter Bereich, immer `Folder`, `Login` und `Share links` in dieser Reihenfolge, dann alle weiteren, mit einem Zeichen: `✔` in Ordnung, `ℹ` Information, `⚠` eine Warnung, `✘` ein Fehler. Ein Ordner ohne Probleme zeigt `✔ Folder  OK`.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ───────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp              │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

Mit `--fix` endet jeder reparierte Befund mit `→ fixed` und zählt nicht mehr als Problem:

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ───────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed      │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ──────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed        │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →   │
 │                fixed                                                      │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link → fixed     │
 │ ℹ              settings.json differed from the group's copy, which it now │
 │                uses; this account's copy is in                            │
 │                /home/me/.local/share/ccam/profiles/personal/.ccam-        │
 │                backup/20261023-212553/settings.json                       │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯
```

In einem Terminal passen sich die Kästen seiner Breite an, lange Meldungen werden umbrochen, und es gelten die Farben des Terminals (`NO_COLOR` schaltet sie aus). Ein Kasten ist nie schmaler als 29 Spalten, sodass in einem schmalen Terminal kein Wort abgeschnitten wird. In eine Pipe oder Datei geleitet, hat die Ausgabe keine Farben, und keine Meldung wird umbrochen oder abgeschnitten, sodass `ccam doctor | grep settings.json` die ganze Zeile findet; jeder Kasten ist dann so breit wie die längste Zeile. Der Tab Health des Dashboards zeigt dieselben Kästen.

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `--fix` | Reparieren, was sich gefahrlos reparieren lässt (geänderte Dateien werden gesichert) | aus |
| `--dry-run` (global) | Alle Prüfungen ausführen, aber nichts reparieren, auch nicht mit `--fix` | aus |

`ccam doctor` nimmt keine Argumente. Es beachtet `--dry-run`: `ccam --dry-run doctor --fix` gibt die Befunde von `ccam doctor` aus, ergänzt, was `--fix` beiseitelegen würde, und ändert nichts:

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

Der Exit-Code ist `1`, wenn mindestens ein `✘`-Befund unrepariert bleibt, sonst `0`. Warnungen und Informationen lassen den Lauf nie fehlschlagen. Schlägt der Lauf ohne `--fix` fehl, endet ccam mit einer Leerzeile und `` `ccam doctor --fix` repairs what it can. ``. Eine Konfiguration, die sich nicht laden lässt, endet ebenfalls mit `1` und einem `ccam:`-Fehler.

## So funktioniert es

`internal/cli/doctor.go` lädt `~/.config/ccam/config.toml` und ruft `doctor.Run` in `internal/doctor/doctor.go` auf, das die Prüfungen in fester Reihenfolge ausführt und Befunde zurückgibt, jeweils mit dem Konto, zu dem sie gehören, und dem geprüften Bereich. `tui.HealthReport` ordnet sie in die Kästen ein. `--fix` wird als `Options.Fix` durchgereicht und bei einem Probelauf auf false gesetzt.

<!-- motion:doctor-checks -->

Der mittlere Block läuft einmal pro Konto. Die Prüfungen und was `--fix` jeweils tut:

1. **Claude Code-Version** (`checkClaude`). Ermittelt die Binärdatei `claude` wie der Starter und führt `claude --version` aus. Eine fehlende Binärdatei, ein `claude_bin`, das keine ausführbare Datei ist, oder ein fehlschlagendes `--version` ist ein Fehler. Keine Reparatur.
2. **Zugangsdaten in der Umgebung** (`checkAuthEnv`). Warnt, wenn `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` oder `ANTHROPIC_AUTH_TOKEN` gesetzt ist: Jedes aus dieser Shell gestartete Konto erbt die Variable, und Claude Code zieht sie der eigenen Anmeldung des Kontos vor. Keine Reparatur; entferne die Variable selbst.
3. **Ordner** (`checkProfile`). Ein fehlender Kontoordner ist ein Fehler und beendet die Prüfungen dieses Kontos. Berechtigungen mit gesetzten Gruppen- oder Fremd-Bits sind eine Warnung; `--fix` führt `chmod 0700` aus. Ein importiertes `~/.claude` mit einer eigenen `.claude.json` neben `~/.claude.json` ist eine Information: `ccam main` und ein einfaches `claude` führen dann getrennte Anmeldedaten und MCP-Server. Windows hat keine Berechtigungsbits, deshalb werden diese Prüfung und die `0600`-Prüfung unten dort übersprungen; den Ordner in deinem Benutzerprofil schützt seine Zugriffsliste.
4. **Anmeldedatei** (`checkCredentials`). `.credentials.json` wird mit `lstat` geprüft. Ein symbolischer Link ist ein Fehler, weil jedes Konto ein eigenes Token haben muss, und wird nicht repariert. Ein anderer Modus als `0600` ist eine Warnung; `--fix` führt `chmod 0600` aus. Dann wird die Anmeldung gelesen (unter macOS zuerst aus dem Schlüsselbund). Keine Anmeldung oder eine Anmeldedatei ohne Anmeldung ist eine Warnung, die auf `ccam login <name>` verweist. Ein abgelaufenes Refresh-Token oder eine Anmeldung, die Claude Code gelöscht hat, nachdem der Server sie abgelehnt hat, ist der Fehler `` login expired → `ccam login <name>` ``; ein abgelaufenes Zugriffstoken ist nur eine Information, da Claude Code es beim nächsten Start erneuert. Eine Anmeldung in der Windows-Anmeldeinformationsverwaltung ist eine Information, da ccam sie nicht lesen kann. Eine gültige Anmeldung zeigt die E-Mail-Adresse, mit der der Ordner gerade angemeldet ist, und den Tarif (`login valid (dev@company.com, Claude Team · Acme)`). Nichts davon wird repariert; `ccam login --expired` meldet jedes Konto an, das es braucht.
5. **Freigabelinks** (`share.Check`). Vergleicht jeden Pfad in `share.Layout` mit den Gruppen des Kontos und verwendet dabei nur `lstat` und `readlink`. Es meldet fünf Arten: einen fehlenden Link, einen Link, der woandershin zeigt, eine echte Datei, wo ein Link hingehört, einen übrig gebliebenen Link für ein Element, das inzwischen lokal ist, und einen richtigen Link, dessen Gruppenordner fehlt (`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)`; ein Start prüft das nicht). Ein fehlender Link und ein fehlender Gruppenordner sind Warnungen, die anderen sind Fehler. Auch ein Zurückkopieren aus der Gruppe, das mittendrin abgebrochen ist, wird gemeldet: `--fix` setzt eine vollständige vorbereitete Kopie an ihren Platz, doch wenn der Pfad seitdem neu angelegt wurde (durch einen Start), nennt es den Ordner mit der vollständigen Kopie, damit du das selbst klärst. Pfade, die in einer späteren Version von ccam zu einem Element hinzugekommen sind (`file-history/`, `tasks/` und `paste-cache/` für den Verlauf, `keybindings.json` für die Einstellungen), sind nur eine Warnung, solange das Konto noch eine eigene Kopie hat: `is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`. Ein Start hält dafür nie an.
6. **Unbekannte Einträge** (`checkUnknownEntries`). Listet Namen auf der obersten Ebene des Kontoordners auf, die ccam nicht kennt (`share.Known`), als Hinweis darauf, dass sich die Ordnerstruktur von Claude Code geändert hat. Nur Information.
7. **Plugin-Pfade** (`checkPlugins`). Plugin-Register (`installed_plugins.json`, `known_marketplaces.json`) speichern absolute Pfade. Im `plugins/` jeder Gruppe gehen Pfade, die durch den Ordner eines Kontos führen, kaputt, wenn dieses Konto entfernt oder umbenannt wird. `--fix` schreibt sie mit `profile.RewritePluginPaths` auf den Gruppenordner um, als atomaren Schreibvorgang der Registerdatei.
8. **MCP-Server** (`checkMCP`). Lädt die `mcp.json` jeder Gruppe einmal. Eine Datei, die sich nicht lesen lässt, ist ein Fehler; ein Server, der `mcp.Validate` nicht besteht, ist eine Warnung. Server in der `~/.claude.json` des einfachen `claude`, die kein Konto lädt, sind eine Warnung, die den Befehl `ccam migrate --mcp` nennt, der sie kopiert (mit `--group <a|b>`, wenn mehrere Gruppen MCP teilen). Server, die nur ein Konto einer Gruppe hat, sind eine Information mit dem Befehl `ccam mcp import`, der sie teilt, und ein Konto, das bei MCP-Servern nicht angemeldet ist, bei denen andere Konten angemeldet sind, bekommt eine Warnung, sich mit `/mcp` anzumelden (ccam kopiert nie Anmeldungen). Keine Reparatur.
9. **Antwortformat der Nutzung** (`checkUsageFormat`). Fragt den Nutzungs-Endpunkt ohne Cache für ein Konto nach dem anderen ab, bis eines antwortet, alles innerhalb eines gemeinsamen Limits von 4 Sekunden. Mit `usage_source = "claude"` fragt es stattdessen Claude Code für ein Konto. Lässt sich die Antwort nicht mehr auswerten, warnt es, dass `ccam usage` auf den Cache der Statuszeile zurückgreift; hat kein Konto geantwortet, sagt es das als Information (`could not check the usage response format: …`). Das ist die einzige Prüfung, die das Netzwerk nutzt.
10. **Kurzbefehle** (`checkShims`). Nur wenn ein Konto einen Kurzbefehl hat. Fehlende Links in `~/.local/bin` sind eine Warnung, ebenso Kurzbefehle, die ccam angelegt hat und die noch ein älteres ccam starten, eines, das ersetzt oder verschoben wurde (`shortcuts in <bin> that run another ccam (one replaced or moved): …`); `--fix` ruft `shim.Install` auf, das nur die Links anlegt und erneuert, die ccam angelegt und in `~/.local/state/ccam/shortcuts.json` festgehalten hat, und nie andere Dateien anfasst. Die Warnung endet nur dann mit `→ fixed`, wenn jeder dieser Links angelegt wurde; konnte einer nicht angelegt werden, wird der Grund angefügt. Ein Kurzbefehl, der ein anderes, noch installiertes ccam startet, funktioniert und ist eine Information; `--fix` richtet ihn auf dieses ccam aus. Ein bin-Ordner, der nicht im `PATH` liegt, ist eine Warnung und wird nicht repariert.
11. **Ordner** (`checkFolders`). Nur wenn es Einträge von [`ccam folder`](/features/folder) gibt. Ein Eintrag, dessen Konto nicht existiert, ist ein Fehler, und ein Eintrag, dessen Ordner verschwunden ist, ist eine Warnung. Keines von beiden wird repariert, weil das Entfernen eines Eintrags ändern würde, welches Konto dort startet.
12. **Gruppen und Einstellungen** (`checkGroups`). Ein Gruppenmitglied, das kein Konto ist, ein Konto, dessen Name inzwischen ein ccam-Befehl geworden ist (starte es mit `ccam run <name>` oder benenne es um), ein unbekanntes `theme` oder `usage_source` und Zeilen der `config.toml`, die ccam nicht liest, sind Warnungen. Letztere sind meist eine Einstellung, die unter einer Tabelle `[[group]]`, `[[profile]]` oder `[[folder]]` steht und damit zu dieser Tabelle gehört (`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`); Befehle, die die Konfiguration ändern, verweigern den Dienst, bis das behoben ist. Keine Reparatur.
13. **Anmeldungen** (`checkIdentities`). Liest aus der `.claude.json` jedes Ordners, als wer er angemeldet ist. Ein Ordner, der inzwischen bei einem anderen Konto angemeldet ist, als `ccam login` festgehalten hat (zum Beispiel nach `/login` in einer Sitzung), ist eine Warnung; ein Konto ohne Anmeldeeintrag, etwa ein importiertes oder eines von einem älteren ccam, wird anhand seiner festgehaltenen E-Mail-Adresse verglichen. Ein Ordner, der beim selben Claude-Konto angemeldet ist wie ein anderer, ist ebenfalls eine Warnung: Die beiden teilen sich ihre Nutzungslimits. Ein Ordner, der dasselbe Anmeldetoken enthält wie ein anderer, also ein kopierter Ordner (auch das `~/.claude` des einfachen `claude`), ist ein Fehler, denn wer es zuerst erneuert, meldet den anderen ab; führe `ccam login <name>` aus, um ihm ein eigenes zu geben. Keine Reparatur.
14. **`auto_switch`** (`checkAutoSwitch`). Nur bei `auto_switch = true`. Eine Information im Kasten `General` sagt, was `auto_switch` tut, und verlinkt auf [seine Beschreibung](/features/next#auto-switch). Für jeden Tarif mit Plätzen unter den Konten (Claude Team, Claude Enterprise) nennt eine Zeile diese Konten, sagt, dass `auto_switch` eine Sitzung zwischen Plätzen einer Organisation verschiebt, zitiert die Platzregel von Anthropic mit dem Link zum Hilfe-Center und überlässt dir die Entscheidung. Ein Claude Code älter als 2.1.214, das Hooks seinen Prozess (`CLAUDE_PID`) nicht mitteilt, ist dort eine Warnung, da sich dann keine Sitzung verschieben lässt. Ein Konto ohne das ccam-Plugin oder mit ausgeschaltetem Plugin bekommt in seinem eigenen Kasten eine Warnung mit dem Befehl, der das in Claude Code behebt: `auto_switch is on, but work does not have the ccam plugin, so its sessions are not moved. In Claude Code as work: /plugin install ccam --marketplace jhyunwoo/ccam` (oder `/plugin enable ccam@ccam`). Konten mit API-Schlüssel werden ausgelassen. Unter Windows ist die einzige Zeile `auto_switch is ignored on Windows; with relay = true, /ccam:switch and then /exit carry a conversation on.` Keine Reparatur.

Vor den Konten ist eine **Änderung der Gruppen**, die mittendrin abgebrochen ist, ein Fehler im Kasten `General`: ``the last change of groups stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``. `--fix` schließt sie in der geplanten Reihenfolge ab, mit den Antworten, die auf ihre Konflikte gegeben wurden; bis dahin beginnt keine andere Änderung der Gruppen, kein Umbenennen und kein Entfernen. Die beteiligten Konten bekommen nur ihre fehlenden Links angelegt, bis sie abgeschlossen ist. Lässt sie sich nicht abschließen, sagt der Fehler, dass das Entfernen dieser Datei den Rest der Änderung aufgibt. `ccam --dry-run doctor --fix` zeigt die Reihenfolge (`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`) und was der Teil jedes Kontos beiseitelegen würde. Eine Änderung, die ein anderes ccam gerade ausführt, ist eine Information, und solange sie läuft, repariert `--fix` keine Freigabelinks (`its share links are not repaired while another ccam changes the groups; …`). `--fix` hält `~/.local/state/ccam/group-change.lock`, während es repariert, und liest `config.toml` nach dem Abschluss einer Änderung neu ein.

### Freigabelinks reparieren

Der Startpfad legt einen fehlenden Link sofort an, hält aber bei allem an, was Daten kosten könnte, und verweist hierher:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Bei jedem Problem mit Freigabelinks ruft `--fix` die Funktion `share.Build` mit den Freigabeeinstellungen des Kontos als aktuellem und als Zielzustand auf, dann `Plan.Apply`. `Build` leitet jeden Schritt neu aus dem ab, was auf der Festplatte liegt, derselbe Plan, den eine Änderung mit `ccam group` verwendet, daher folgt die Reparatur denselben Regeln wie der Beitritt zu einer Gruppe.

<!-- motion:doctor-repair -->

Hier ist `settings.json` eine echte Datei. Weil die Gruppe bereits eine andere `settings.json` hat, nimmt der Konflikt seine Standardwahl `KeepGroup`, sodass die Kopie der Gruppe bleibt. Die Datei des Kontos wandert nach `.ccam-backup/<time>/settings.json` im Kontoordner, und `os.Symlink` verlinkt den Pfad mit der Gruppe. `--fix` sagt das: `settings.json differed from the group's copy, which it now uses; this account's copy is in …`. Andere Fälle:

- Hat die Gruppe noch keine Kopie, wandert die Datei des Kontos in die Gruppe und wird ihr erster Inhalt.
- Bei einem Ordner wandern Dateien, die der Gruppe fehlen, in die Gruppe. Dateien, die die Gruppe Byte für Byte gleich hat, werden verworfen, und die abweichenden kommen in die Sicherung und werden in einem Hinweis genannt (`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`).
- Ein fehlender Gruppenordner wird zuerst angelegt.
- Bei `history.jsonl` ist die Standardwahl `merge by line`: die Zeilen der Gruppe, dann deine neuen. Die vorherige Kopie der Gruppe wird in `.ccam-backup/<time>/group/` gesichert.
- Ein Link auf eine andere Gruppe wird ersetzt. Ein Link auf einen Ordner außerhalb von ccam kommt in die Sicherung, und sein Inhalt wird zum Startinhalt der Gruppe, wenn diese leer ist.
- Ein übrig gebliebener Link für ein Element, das inzwischen lokal ist, wird entfernt, und das Konto erhält eine Kopie des Inhalts der Gruppe.

Abgesehen vom Entfernen von Links und vom Verwerfen von Dateien, die die Gruppe bereits Byte für Byte gleich hat, verschiebt jeder Schritt Daten in die Gruppe oder in den Sicherungsordner. Nichts geht verloren.

## Sicherheit und Sonderfälle

- Ein einfaches `ccam doctor` ändert nie etwas. Verwende in Skripten `--dry-run` mit `--fix`, um sicherzugehen.
- Für ein laufendes Konto legt `doctor --fix`, egal wie Claude Code gestartet wurde, nur fehlende Links und fehlende Gruppenordner an, was nichts kostet. Bleibt noch etwas übrig, fügt es `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest `` hinzu, weil die Sitzung weiter in die Dateien schreibt, die eine Reparatur verschieben würde.
- Schlägt eine Reparatur der Freigaben fehl, meldet ccam für dieses Konto `could not repair share links: <error>` und überspringt seine Prüfung auf unbekannte Einträge. Die anderen Konten werden trotzdem geprüft.
- Wenn eine Änderung mit `ccam group` mittendrin fehlschlägt, sind die neuen Gruppen bereits gespeichert, und die Änderung ist in `~/.local/state/ccam/group-change.json` festgehalten; ccam rät dir, mit `ccam doctor --fix` abzuschließen, das den Rest davon ausführt, MCP-Schritte eingeschlossen.
- Eine per symbolischem Link verbundene `.credentials.json`, eine abgelaufene Anmeldung und Zugangsdaten in der Umgebung werden gemeldet, aber nie geändert.
- Umgeschriebene Plugin-Register werden atomar geschrieben, aber nicht nach `.ccam-backup/` kopiert.

## Verwandte Themen

- [MCP-Server, Einstellungen und Verlauf teilen](/docs/sharing)
- [Eine vorhandene Einrichtung kopieren](/features/migrate)
- [Nutzungslimits und Kontowechsel](/docs/usage-limits)
- [Befehle](/docs/commands)
