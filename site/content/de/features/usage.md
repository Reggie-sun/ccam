---
title: Nutzungslimits
description: ccam usage zeigt die 5-Stunden- und Wochenlimits aller Konten auf einmal, und ccam statusline bringt Kontoname und Nutzung in Claude Code.
updated: 2026-10-21
order: 7
translation: machine
sourceHash: "b26800c6fa54"
---

`ccam usage` zeigt, wie viel jedes Konto von seinem 5-Stunden-Limit, seinem Wochenlimit und dem eigenen Wochenlimit jedes Modells (etwa Fable) verbraucht hat, als Balken pro Limit mit der Zeit bis zum Zurücksetzen. Verwende es, bevor du eine lange Aufgabe beginnst oder wenn Claude Code meldet, dass ein Limit nahe ist, um zu sehen, welches Konto noch Spielraum hat. `ccam statusline` ist der Statuszeilenbefehl für Claude Code: Er zeigt Kontoname und Nutzung in Claude Code an und speichert die Zahlen für `ccam usage`.

## Verwendung

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam usage --source claude  # let Claude Code look it up: ccam never touches a token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

Jedes Konto bekommt eine Zeile mit seiner Nummer, seinem Namen und dem Alter seiner Zahlen; `(cached)` bedeutet, dass sie aus einer gespeicherten Datei stammen, nicht vom Endpunkt. Darunter zeigt eine Zeile pro Limit einen Balken aus 20 Zellen, den Prozentwert und die Zeit bis zum Zurücksetzen des Zeitfensters: `5-hour`, `Weekly` (alle Modelle), dann das eigene Wochenfenster jedes Modells, etwa `Fable`. Ein Zeitfenster, das das Konto nicht hat, zeigt `-`. Ein Prozentwert wird ab 80 % gelb und bei 100 % rot. In einem Terminal erscheint jedes Konto sofort als `fetching…`, und seine Zeilen werden gezeichnet, sobald es antwortet. Die letzte Zeile nennt das Konto, das [`ccam next --auto`](/features/next) nach denselben Regeln wählen würde (`auto_max_usage`, `auto_strategy`, die laufenden Sitzungen, ohne gesetztes `cross_org` nur Konten aus der Organisation des zuletzt verwendeten Kontos, jedes Zeitfenster bei 100 % wird übersprungen, sodass ein Konto mit voller Woche nie genannt wird), und erscheint nur, wenn du mehr als ein Konto hast. Ein nicht angemeldetes Konto zeigt `not logged in → ccam login <name>`, und ein Konto mit API-Schlüssel zeigt `no limits (API)`.

Nach `ccam statusline install` zeigt Claude Code eine Zeile wie `[work] · Opus · 5h 24% · 7d 41%`.

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `--json` | Die Berichte als JSON ausgeben: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | aus |
| `--watch[=N]` | Alle `N` Sekunden neu zeichnen, bis du Ctrl+C drückst. `--watch` allein oder ein Wert von 0 oder weniger bedeutet 60 | aus |
| `-v`, `--verbose` | Unter jedem Konto eine Zeile `source:` hinzufügen: woher seine Zahlen stammen (der Nutzungs-Endpunkt, Claude Code, die Statuszeile oder eine gespeicherte Datei) | aus |
| `--no-cache` | Den 60-Sekunden-Cache ignorieren und den Endpunkt erneut abfragen | aus |
| `--source live\|claude\|cache` | `live` fragt den Endpunkt mit dem Token des Kontos ab; `claude` lässt Claude Code nachsehen (siehe unten); `cache` zeigt nur gespeicherte Zahlen | `usage_source`, sonst `live` |
| `--refresh` | Abgelaufene Zugriffstokens erneuern, nur für Konten, die nicht laufen | aus |
| `--dry-run` | Wird akzeptiert: `ccam usage` ändert nichts außer seinem eigenen Cache. Zusammen mit `--refresh` abgelehnt: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | aus |

Jeder andere Wert für `--source` schlägt mit `--source must be live, claude or cache` fehl.

### Nutzung über Claude Code selbst

Mit `usage_source = "claude"` in `config.toml` (oder `--source claude`) liest ccam das Anmeldetoken überhaupt nicht. Für jedes Konto startet es `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` im Ordner des Kontos und sendet die Anfrage `get_usage` von Claude Code über die Standardeingabe (Claude Code 2.1.270 oder neuer); Claude Code sieht die Nutzung mit seiner eigenen Anmeldung nach und antwortet mit dem 5-Stunden-, dem Wochen- und den Modellfenstern. `--safe-mode` verhindert, dass Hooks, MCP-Server und Plugins des Kontos starten, und es wird keine Sitzung gespeichert. Jedes Konto kostet einen Start von Claude Code, etwa eine Sekunde, daher ist diese Quelle langsamer als `live`; die Konten werden gleichzeitig abgefragt, jedes innerhalb von 20 Sekunden, und der 60-Sekunden-Cache gilt wie gewohnt. Claude Code kennzeichnet die Anfrage als experimentell; antwortet sie nicht mehr, zeigt ccam die gespeicherten Zahlen, und `ccam doctor` meldet es. Ein `--dry-run` startet Claude Code nicht (Claude Code schreibt in den Ordner des Kontos) und zeigt die gespeicherten Zahlen. `usage_source` nimmt `live` oder `claude` an; `cache` gibt es nur als Wahl bei `--source`. `ccam usage` nimmt keine Argumente.

`ccam statusline` liest das Statuszeilen-JSON von Claude Code von der Standardeingabe. Seine einzige Option ist `--wrap <command>`: ein vorhandener Statuszeilenbefehl, dessen erster Ausgabezeile der Kontoname vorangestellt wird. `ccam statusline install` nimmt keine Optionen. Beide lehnen `--dry-run` mit `` --dry-run is not supported by `ccam statusline`; nothing was changed `` ab.

## So funktioniert es

Die Zahlen kommen aus zwei Quellen: dem nicht dokumentierten Nutzungs-Endpunkt von Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, aufgerufen mit der eigenen Anmeldung jedes Kontos, und den Werten, die Claude Code an seinen Statuszeilenbefehl übergibt. All das liegt in `internal/usage`.

1. `cli.usageCmd` lädt `config.toml` und baut pro Konto einen `usage.Account`: seinen Ordner `~/.local/share/ccam/profiles/<name>/`, seinen Sitzungsordner (um zu sehen, ob es läuft) und ob es an Ort und Stelle importiert wurde.
2. `Fetcher.FetchAll` startet eine Goroutine pro Konto. Alle teilen sich einen `http.Client`, sodass Verbindungen zum Endpunkt wiederverwendet werden. Der Verbindungsaufbau darf 2 s dauern, und jedes Konto bekommt insgesamt 3 s. Der Client folgt keinen Weiterleitungen, da eine Weiterleitung das Token an eine andere Stelle senden würde, und liest höchstens 1 MB der Antwort.
3. Sobald ein Konto fertig ist, wird seine Zeile neu gezeichnet. Ein Konto, das fehlschlägt oder nicht rechtzeitig antwortet, greift auf gespeicherte Zahlen zurück, sodass ein langsames Konto nie die anderen aufhält.

Jedes Konto läuft in seiner eigenen Goroutine gegen den Endpunkt; das Konto mit Zeitüberschreitung zeigt stattdessen seine gespeicherten Zahlen.

<!-- motion:usage-fetch -->

`personal` und `work` antworten rechtzeitig und zeigen Live-Zahlen. `side` braucht länger als 3 s, daher liest `fallback()` seine beiden Cache-Dateien und verwendet die zuletzt geschriebene. `usage.Table()` zeichnet einen Block pro Konto, und `usage.Best()` wählt das Konto, das in der letzten Zeile genannt wird.

Für ein Konto führt `Fetcher.Fetch` diese Prüfungen der Reihe nach aus:

<!-- motion:usage-cache -->

1. Mit `--source cache` geht es direkt zu `fallback()`.
2. Sofern `--no-cache` nicht gesetzt ist, wird ein Live-Ergebnis in `~/.local/state/ccam/usage/<name>.json`, das jünger als 60 Sekunden ist (`CacheTTL`), unverändert zurückgegeben. Zweimal hintereinander `ccam usage` kostet eine Anfrage pro Konto.
3. `readCredentials()` liest `.credentials.json` (oder den macOS-Schlüsselbund), ohne etwas zu ändern, und dekodiert nur den Schlüssel `claudeAiOauth`. Ohne Abo-Anmeldung lautet das Ergebnis `not logged in`, oder `no limits (API)`, wenn `.claude.json` einen API-Schlüssel enthält.
4. Ein abgelaufenes Zugriffstoken wird nicht verwendet. Ohne `--refresh` zeigt das Konto seine gespeicherten Zahlen, als abgelaufen markiert.
5. Andernfalls ruft ccam den Endpunkt auf. Ein 401 oder 403 zählt als abgelaufene Anmeldung, jeder andere Status oder eine nicht lesbare Antwort zählt als Fehler, und beides fällt auf gespeicherte Zahlen zurück. Ein 429 bedeutet, dass der Endpunkt zu oft gefragt wird: ccam schreibt `<name>.pause` mit dem Zeitpunkt für die nächste Anfrage (`Retry-After`, sonst 5 Minuten, höchstens eine Stunde), zeigt gespeicherte Zahlen mit `rate limited by the usage endpoint; asking again after 15:10` und ruft den Endpunkt für dieses Konto bis dahin nicht auf, auch nicht mit `--no-cache`. Ein sofortiger neuer Versuch würde das Limit nur verlängern.
6. Eine gültige Antwort wird mit dem Messzeitpunkt nach `<name>.json` geschrieben.

`fallback()` liest sowohl `<name>.json` als auch `<name>.statusline.json` und behält die neuere. Bevor es gespeicherte Zahlen zeigt, prüft es den Zurücksetzungszeitpunkt jedes Zeitfensters: Ein Fenster, das seit dem Speichern der Zahlen zurückgesetzt wurde, zeigt 0 %, weil der alte Prozentsatz nicht mehr gilt. Deshalb zeigt `personal` oben `0%` ohne Zurücksetzungszeit.

Der Endpunkt hat zwei Antwortformate. Das neuere hat ein Array `limits` mit den Einträgen `session` (5 Stunden), `weekly_all` und `weekly_scoped` (pro Modell). Das ältere hat die Objekte `five_hour`, `seven_day`, `seven_day_opus` und `seven_day_sonnet`. `ParseLive` liest beide: Das Array `limits` hat Vorrang, und die alten Objekte ergänzen, was es nicht abdeckt. Eine Antwort ganz ohne Limits wird als geändertes Format gemeldet, worauf `ccam doctor` dich hinweist.

### Tokens erneuern

Refresh-Tokens rotieren: Wer eines erneuert, bekommt ein neues Refresh-Token, und das alte wird ungültig. Würde ccam die Anmeldung erneuern, während Claude Code als dasselbe Konto läuft, würde diese Sitzung abgemeldet. Deshalb liest `ccam usage` Tokens und schreibt sie nie, es sei denn, du gibst `--refresh` an. Dann geht `Fetcher.refresh` so vor:

1. Es lehnt Konten ab, die mit `ccam import` an Ort und Stelle registriert wurden, da einfaches `claude` diese Anmeldung verwenden könnte, ohne dass ccam davon weiß;
2. es lehnt Anmeldungen ab, die im macOS-Schlüsselbund gespeichert sind;
3. es nimmt eine Sperrdatei, `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock` (unter `~/.local/state/ccam/locks/`, wenn diese Variable nicht gesetzt ist), damit zwei ccam-Prozesse nie dasselbe Token zweimal erneuern;
4. unter der Sperre lehnt es ab, wenn eine Sitzung des Kontos läuft, und überspringt die Anfrage, wenn ein anderer Prozess das Token bereits erneuert hat;
5. es sendet das Refresh-Token an den Token-Endpunkt und schreibt die neuen Tokens mit `fsutil.AtomicWrite` zurück nach `.credentials.json`, wobei jeder andere Schlüssel der Datei erhalten bleibt.

### Der Statuszeilen-Sammler

Claude Code führt den Statuszeilenbefehl aus und leitet ihm ein JSON-Dokument zu. ccam liest daraus den Modellnamen und `rate_limits`:

<!-- motion:usage-statusline -->

1. `runStatusline` findet das Konto über `CCAM_PROFILE`, das ccam beim Start von Claude Code setzt, oder sonst durch Abgleich von `CLAUDE_CONFIG_DIR` mit den Kontoordnern.
2. `ParseStatusLine` liest `rate_limits.five_hour` und `rate_limits.seven_day`. Bei Konten mit API-Schlüssel oder vor der ersten Antwort der Sitzung enthält das JSON keine Ratenlimits; dann wird nichts gespeichert.
3. `SaveStatusLine` schreibt `~/.local/state/ccam/usage/<name>.statusline.json`. `ccam usage` verwendet diese Datei, wenn der Endpunkt nicht antwortet.
4. ccam gibt `[name] · model · 5h N% · 7d N%` aus. Mit `--wrap` führt es deinen Befehl über `sh -c` mit derselben Eingabe aus und gibt stattdessen `[name]` gefolgt von der ersten Zeile des Befehls aus. Unter Windows führt es den Befehl so aus wie Claude Code: mit Git Bash, wenn es installiert ist, sonst mit PowerShell.
5. Wenn die 5-Stunden- oder Wochennutzung bei 90 % oder mehr liegt und es ein weiteres Konto gibt, endet die Zeile mit `→ ccam next --auto --continue`, zum Beispiel `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

Der umschlossene Befehl bekommt 2 Sekunden. Danach wird er beendet, und die Zeile zeigt nur den Kontonamen, sodass ein langsamer Befehl die Statuszeile von Claude Code nicht blockieren kann.

`ccam statusline install` schreibt `statusLine` in die `settings.json` jedes Kontos. Hat ein Konto bereits einen Statuszeilenbefehl, behält ccam ihn als `--wrap '<old command>'` bei. Unter Windows wird der Pfad zu ccam mit Schrägstrichen und `~` für den Benutzerordner geschrieben (`~/.local/bin/ccam.exe statusline`), was sowohl Git Bash als auch PowerShell verstehen. Eine über eine Gruppe geteilte `settings.json` wird einmal über ihren Link geschrieben, sodass jedes Konto der Gruppe die Einstellung bekommt. Ein erneuter Aufruf gibt `(already set)` für Dateien aus, die ccam bereits verwenden.

## Sicherheit und Sonderfälle

- `ccam usage` ändert nie eine Token-Datei, es sei denn, du gibst `--refresh` an. Ist eine Anmeldung abgelaufen, starte das Konto einmal mit `ccam <name>`, dann erneuert Claude Code sie.
- `--refresh` fasst nie laufende Konten, importierte Konten oder Anmeldungen im Schlüsselbund an. Eine abgelehnte Erneuerung zeigt die gespeicherten Zahlen und den Grund.
- Der Nutzungs-Endpunkt ist nicht öffentlich dokumentiert und kann sich ändern. Wenn das passiert, siehst du weiterhin gespeicherte Zahlen, und `ccam doctor` meldet die Formatänderung.
- Ohne gespeicherte Zahlen und ohne Antwort zeigt die Kontozeile `lookup failed:` und den Grund, zum Beispiel `no saved usage` mit `--source cache`.
- Die Cache-Dateien enthalten nur Prozentsätze und Zurücksetzungszeiten, nie Tokens.

## Verwandte Themen

- [Nutzungslimits und Kontowechsel](/docs/usage-limits)
- [Zum nächsten Konto wechseln](/features/next)
- [Umgang mit Nutzungslimits](/guides/usage-limit-workflow)
- [Befehle](/docs/commands)
