---
title: Nutzungslimits
description: ccam usage zeigt die 5-Stunden- und Wochenlimits aller Konten auf einmal, und ccam statusline bringt Kontoname und Nutzung in Claude Code.
updated: 2026-10-25
order: 7
translation: machine
sourceHash: "c98993fcdc7d"
---

`ccam usage` zeigt, wie viel jedes Konto von seinem 5-Stunden-Limit, seinem Wochenlimit und dem eigenen Wochenlimit jedes Modells (etwa Fable) verbraucht hat, als Balken pro Limit mit der Zeit bis zum Zurücksetzen. Verwende es, bevor du eine lange Aufgabe beginnst oder wenn Claude Code meldet, dass ein Limit nahe ist, um zu sehen, welches Konto noch Spielraum hat. `ccam statusline` ist der Statuszeilenbefehl für Claude Code: Er zeigt Kontoname und Nutzung in Claude Code an und speichert die Zahlen für `ccam usage`. Für die Tokens und geschätzten Kosten hinter der Nutzung führt `ccam ccusage` [ccusage](https://ccusage.com) über die Gespräche aller Konten aus ([unten](#tokens-and-costs-ccam-ccusage)).

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
ccam usage --source claude  # let Claude Code look it up: ccam sends no token
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
 3 side      access token expired; `ccam side` renews it when it starts
→ Most headroom: personal   (open it with ccam next --auto)
```

Jedes Konto bekommt eine Zeile mit seiner Nummer, seinem Namen und dem Alter seiner Zahlen; `(cached)` bedeutet, dass sie aus einer gespeicherten Datei stammen, nicht vom Endpunkt. Darunter zeigt eine Zeile pro Limit einen Balken aus 20 Zellen, den Prozentwert und die Zeit bis zum Zurücksetzen des Zeitfensters: `5-hour`, `Weekly` (alle Modelle), dann das eigene Wochenfenster jedes Modells, etwa `Fable`. Ein Zeitfenster, das das Konto nicht hat, zeigt `-`. Ein Prozentwert wird ab 80 % gelb und bei 100 % rot. Ein Zeitfenster, das Claude Code als voll gemeldet hat, als es eine Anfrage an einem Nutzungslimit abgelehnt hat, zeigt bis zu seinem Zurücksetzen 100 %, egal was die gespeicherten Zahlen sagen (siehe [den Limit-Eintrag](#the-limit-record)). In einem Terminal erscheint jedes Konto sofort als `fetching…`, und seine Zeilen werden gezeichnet, sobald es antwortet. Die letzte Zeile nennt das Konto, das [`ccam next --auto`](/features/next) nach denselben Regeln wählen würde (`auto_max_usage`, `auto_strategy`, die laufenden Sitzungen, ohne gesetztes `cross_org` nur Konten aus der Organisation des zuletzt verwendeten Kontos, jedes Zeitfenster bei 100 % wird übersprungen, sodass ein Konto mit voller Woche nie genannt wird), und erscheint nur, wenn du mehr als ein Konto hast. Mit der Standardstrategie `headroom` lautet sie `→ Most headroom: <name>`; mit `expiring` oder `balance` nennt sie den Grund, etwa `→ next --auto picks personal: its 5-hour window resets first, in 1h 12m (40% of its 5 hours used)`.

**Team- und Enterprise-Plätze.** Die automatische Wahl (`ccam next`, `--fallback`, `ccam pick`, `n` im Dashboard und [`auto_switch`](/features/next#auto-switch), wenn du es einschaltest) wechselt auch zwischen Plätzen derselben Claude Team- oder Enterprise-Organisation. Im Hilfe-Center von Anthropic ([Team-Tarife](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); ebenso für [Enterprise-Tarife](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)) heißt es: „Each seat on a Team plan is assigned to one named individual, and only that person may use it“ (jeder Platz in einem Team-Tarif ist einer namentlich benannten Person zugewiesen, und nur diese Person darf ihn nutzen) und „Users can't route their work through another user's account to use that user's remaining usage“ (Benutzer dürfen ihre Arbeit nicht über das Konto eines anderen Benutzers leiten, um dessen verbleibendes Nutzungskontingent zu verwenden). Wenn du die automatische Wahl zwischen Plätzen nutzt, kann das daher gegen die Bedingungen von Anthropic verstoßen. Ob du sie dort nutzt, ist deine Entscheidung: Kläre das mit deiner Organisation und mit Anthropic ab, bevor du dich darauf verlässt. Für private Tarife (Pro, Max) gilt der [Hinweis zu den Bedingungen in den FAQ](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

Ein nicht angemeldetes Konto zeigt `not logged in → ccam login <name>`, und ein Konto mit API-Schlüssel zeigt `no limits (API)`. Ein abgelaufenes Zugriffstoken zeigt ``access token expired; `ccam <name>` renews it when it starts``, über den gespeicherten Zahlen, wenn es welche gibt: Claude Code erneuert es beim nächsten Start. Eine Anmeldung, die als Ganzes abgelaufen ist, weil ihr Refresh-Token abgelaufen ist oder Claude Code es gelöscht hat, nachdem der Server es abgelehnt hat, zeigt `` login expired → `ccam login <name>` ``.

Nach `ccam statusline install` zeigt Claude Code eine Zeile wie `[work] · Opus · 5h 24% · 7d 41%`.

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `--json` | Die Berichte als JSON ausgeben: `name`, `slot`, `status`, `limits` (mit `credits`, wenn das Konto Nutzungsguthaben aktiviert hat), `source`, `at`, `error` sowie `limited` und `limited_windows`, wenn ein Limit-Eintrag Zeitfenster auf 100 % setzt. Ohne Konten gibt es `[]` aus | aus |
| `--watch[=N]` | Alle `N` Sekunden neu zeichnen, bis du Ctrl+C drückst. `--watch` allein oder ein Wert von 0 oder weniger bedeutet 60. Der Wert steht nach `=`: `--watch 10` wird mit einem Hinweis abgelehnt | aus |
| `-v`, `--verbose` | Unter jedem Konto eine Zeile `source:` hinzufügen: woher seine Zahlen stammen (der Nutzungs-Endpunkt, Claude Code, die Statuszeile, eine gespeicherte Datei oder das von Claude Code gemeldete Nutzungslimit). Liegt ein Limit-Eintrag über den Zahlen, folgt eine Zeile wie `Weekly full: Claude Code reported the usage limit 1h 0m ago, so 100% until the reset` | aus |
| `--no-cache` | Den 60-Sekunden-Cache ignorieren und den Endpunkt erneut abfragen | aus |
| `--source live\|claude\|cache` | `live` fragt den Endpunkt mit dem Token des Kontos ab; `claude` lässt Claude Code nachsehen (siehe unten); `cache` zeigt nur gespeicherte Zahlen | `usage_source`, sonst `live` |
| `--refresh` | Abgelaufene Zugriffstokens erneuern, nur für Konten, die nicht laufen. Eine Erneuerung kann bis zu 30 Sekunden dauern | aus |
| `--dry-run` | Wird akzeptiert: `ccam usage` ändert nichts außer seinem eigenen Cache. Zusammen mit `--refresh` abgelehnt: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | aus |

Jeder andere Wert für `--source` schlägt mit `--source must be live, claude or cache` fehl.

### Nutzung über Claude Code selbst

Mit `usage_source = "claude"` am Anfang von `config.toml` (vor dem ersten `[[group]]`, `[[profile]]` oder `[[folder]]`) oder `--source claude` liest ccam das Anmeldetoken überhaupt nicht. Für jedes Konto startet es `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` im Ordner des Kontos und sendet die Anfrage `get_usage` von Claude Code über die Standardeingabe (Claude Code 2.1.270 oder neuer); Claude Code sieht die Nutzung mit seiner eigenen Anmeldung nach und antwortet mit dem 5-Stunden-, dem Wochen- und den Modellfenstern. `--safe-mode` verhindert, dass Hooks, MCP-Server und Plugins des Kontos starten, und es wird keine Sitzung gespeichert. Jedes Konto kostet einen Start von Claude Code, etwa eine Sekunde, daher ist diese Quelle langsamer als `live`; die Konten werden gleichzeitig abgefragt, jedes innerhalb von 20 Sekunden, und der 60-Sekunden-Cache gilt wie gewohnt. Eine Anmeldung, deren Zugriffstoken abgelaufen ist oder innerhalb von 5 Minuten abläuft, bekommt 50 Sekunden, da Claude Code es beim Start erneuert und eine abgebrochene Erneuerung die Anmeldung kosten kann. Eine Abfrage, die aufgegeben wird, bittet Claude Code mit SIGTERM, sich zu beenden, sodass es eine begonnene Erneuerung abschließt (Windows hat kein solches Signal; dort wird es gestoppt). Eine Anmeldung, die als Ganzes abgelaufen ist, zeigt `` login expired → `ccam login <name>` ``, ohne Claude Code zu starten. Claude Code kennzeichnet die Anfrage als experimentell; antwortet sie nicht mehr, zeigt ccam die gespeicherten Zahlen, und `ccam doctor` meldet es. Antwortet Claude Code, dass seine eigene Abfrage fehlgeschlagen ist, zeigt ccam die gespeicherten Zahlen mit `the usage endpoint gave Claude Code no usage (it may be refusing lookups for now)` und keinen Plan ohne Limits. Ein `--dry-run` startet Claude Code nicht (Claude Code schreibt in den Ordner des Kontos) und zeigt die gespeicherten Zahlen. `usage_source` nimmt `live` oder `claude` an; `cache` gibt es nur als Wahl bei `--source`. `ccam usage` nimmt keine Argumente.

`ccam statusline` liest das Statuszeilen-JSON von Claude Code von der Standardeingabe. Seine einzige Option ist `--wrap <command>`: ein vorhandener Statuszeilenbefehl, dessen erster Ausgabezeile der Kontoname vorangestellt wird. `ccam statusline install` nimmt keine Optionen. Beide lehnen `--dry-run` mit `` --dry-run is not supported by `ccam statusline`; nothing was changed `` ab.

## Tokens und Kosten: `ccam ccusage`

`ccam usage` zeigt, wie viel ein Konto von jedem Limit verbraucht hat. [ccusage](https://ccusage.com) zählt dagegen die Tokens und geschätzten Kosten in den Gesprächsprotokollen von Claude Code. Für sich allein sucht es die Protokolle von Claude Code nur in `~/.claude` und `~/.config/claude` (oder in den Ordnern, die `CLAUDE_CONFIG_DIR` nennt) und sieht daher von ccams Konten höchstens eines, das an Ort und Stelle unter `~/.claude` importiert wurde. `ccam ccusage` startet es und setzt dabei `CLAUDE_CONFIG_DIR` auf die Ordner aller Konten, die Gespräche haben, sodass ein einziger Bericht sie alle abdeckt:

```sh
ccam ccusage                           # every account, ccusage's daily report
ccam ccusage daily --breakdown         # cost per model
ccam ccusage claude monthly --json     # Claude Code alone, as JSON
ccam ccusage -a work blocks --active   # one account's current 5-hour block
ccam ccusage -a work,side weekly       # these accounts, and any that share their history
```

ccam liest zuerst seine eigenen Optionen: `-a` (`--account`) nimmt Kontonamen oder -nummern, kommagetrennt oder mehrfach angegeben, und `--dry-run` gibt den Befehl aus, statt ihn auszuführen. Ab dem ersten Argument, das keine dieser Optionen ist, geht alles unverändert an ccusage; `ccam ccusage -- --help` zeigt die eigene Hilfe von ccusage. ccam startet ein installiertes `ccusage`, sonst `bun x ccusage`, sonst `npx ccusage@latest`, wobei die letzten beiden es von npm holen; ist keines davon vorhanden, fordert ccam dich auf, Node.js, Bun oder ccusage zu installieren. Unter Linux und macOS tritt ccusage an die Stelle von ccam; unter Windows wartet ccam auf ccusage und endet mit dessen Exit-Code. Unter WSL überspringt ccam ein `ccusage`, `bun` oder `npx` auf einem Windows-Laufwerk: Es liefe als Windows-Programm und bekäme die Kontoordner nicht.

- Der Bericht umfasst die Konten, deren Ordner Gespräche enthalten (einen Ordner `projects`). Ein Gespräch, das Konten über eine Gruppe teilen, zählt nur einmal.
- Mit `-a` bringt ein Konto, das seinen Verlauf über eine Gruppe teilt, die Gespräche der ganzen Gruppe mit, und ccam sagt, wessen Gespräche das sind: `ccam: work shares its conversation history with personal, so the report counts their conversations too`.
- `blocks` zeigt 5-Stunden-Fenster, und jedes Konto hat seine eigenen, also nenne dafür ein einziges Konto, wie oben.
- ccusage erfasst auch die anderen Coding-Tools, die es findet, etwa Codex, sofern du nicht eines davon nennst: `ccam ccusage claude daily` zeigt nur Claude Code.
- Während ccam es ausführt, liest ccusage seine Einstellungen (`ccusage.json`) nur aus `./.ccusage` und den Kontoordnern. `--config` nennt eine andere Datei: `ccam ccusage daily --config ~/.config/claude/ccusage.json`.
- ccusage würde einen Ordner, dessen Pfad ein Komma enthält, als zwei Ordner lesen, daher lehnt ccam ihn ab und nennt das Konto; lass dieses Konto mit `-a` weg.
- `ccusage` ist ein ccam-Befehl, daher wird ein Konto namens ccusage mit `ccam run ccusage` gestartet; `ccam doctor` weist darauf hin.

## So funktioniert es

Die Zahlen von `ccam usage` kommen aus zwei Quellen: dem nicht dokumentierten Nutzungs-Endpunkt von Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, aufgerufen mit der eigenen Anmeldung jedes Kontos, und den Werten, die Claude Code an seinen Statuszeilenbefehl übergibt. All das liegt in `internal/usage`.

1. `cli.usageCmd` lädt `config.toml` und baut pro Konto einen `usage.Account`: seinen Ordner `~/.local/share/ccam/profiles/<name>/`, seine Sitzungsmarkierungen (zusammen mit den eigenen Einträgen von Claude Code im Kontoordner zeigen sie, ob es läuft) und ob es an Ort und Stelle importiert wurde.
2. `Fetcher.FetchAll` startet eine Goroutine pro Konto. Alle teilen sich einen `http.Client`, sodass Verbindungen zum Endpunkt wiederverwendet werden. Der Verbindungsaufbau darf 2 s dauern, und jedes Konto bekommt insgesamt 3 s. Der Client folgt keinen Weiterleitungen, da eine Weiterleitung das Token an eine andere Stelle senden würde, und liest höchstens 1 MB der Antwort.
3. Sobald ein Konto fertig ist, wird seine Zeile neu gezeichnet. Ein Konto, das fehlschlägt oder nicht rechtzeitig antwortet, greift auf gespeicherte Zahlen zurück, sodass ein langsames Konto nie die anderen aufhält.

Jedes Konto läuft in seiner eigenen Goroutine gegen den Endpunkt; das Konto mit Zeitüberschreitung zeigt stattdessen seine gespeicherten Zahlen.

<!-- motion:usage-fetch -->

`personal` und `work` antworten rechtzeitig und zeigen Live-Zahlen. `side` braucht länger als 3 s, daher zeigt `fallback()` stattdessen seine gespeicherten Zahlen. `usage.Table()` zeichnet einen Block pro Konto, und `usage.Best()` wählt das Konto, das in der letzten Zeile genannt wird.

Für ein Konto führt `Fetcher.Fetch` diese Prüfungen der Reihe nach aus:

<!-- motion:usage-cache -->

1. Mit `--source cache` geht es direkt zu `fallback()`.
2. Sofern `--no-cache` nicht gesetzt ist, wird ein Live-Ergebnis in `~/.local/state/ccam/usage/<name>.json`, das vor weniger als 60 Sekunden gemessen wurde (`CacheTTL`), verwendet, bewertet nach dem jetzigen Zeitpunkt: Ein Zeitfenster, das seitdem zurückgesetzt wurde, zeigt 0 %, und eine Datei mit einem Datum nach der aktuellen Uhrzeit gilt nicht als frisch. Zweimal hintereinander `ccam usage` kostet eine Anfrage pro Konto.
3. `readCredentials()` liest `.credentials.json` (oder den macOS-Schlüsselbund), ohne etwas zu ändern, und dekodiert nur den Schlüssel `claudeAiOauth`. Ohne Abo-Anmeldung lautet das Ergebnis `not logged in`, oder `no limits (API)`, wenn `.claude.json` einen API-Schlüssel enthält. Eine Anmeldung, deren Refresh-Token abgelaufen ist oder die Claude Code gelöscht hat, ist `` login expired → `ccam login <name>` `` und wird nicht abgefragt.
4. Ein abgelaufenes Zugriffstoken wird nicht verwendet. Ohne `--refresh` zeigt das Konto seine gespeicherten Zahlen mit ``access token expired; `ccam <name>` renews it when it starts``.
5. Andernfalls ruft ccam den Endpunkt auf. Ein 401 oder 403 zählt als abgelaufene Anmeldung (`login expired: authentication failed (401)`), jeder andere Status oder eine nicht lesbare Antwort zählt als Fehler, und beides fällt auf gespeicherte Zahlen zurück. Eine mit Ctrl+C beendete Abfrage meldet `cancelled`. Ein 429 bedeutet, dass der Endpunkt zu oft gefragt wird: ccam schreibt `<name>.pause` mit dem Zeitpunkt für die nächste Anfrage (`Retry-After`, sonst 5 Minuten, höchstens eine Stunde), zeigt gespeicherte Zahlen mit `rate limited by the usage endpoint; asking again after 15:10` und ruft den Endpunkt für dieses Konto bis dahin nicht auf, auch nicht mit `--no-cache`. Ein sofortiger neuer Versuch würde das Limit nur verlängern.
6. Eine gültige Antwort wird mit dem Messzeitpunkt nach `<name>.json` geschrieben.

`fallback()` nimmt die neueste gespeicherte Abfrage (`<name>.json` oder `<name>.claude.json` aus den Abfragen über Claude Code) und legt, wenn `<name>.statusline.json` danach geschrieben wurde, das 5-Stunden- und das Wochenfenster der Statuszeile darüber; die Modellfenster, die die Statuszeile nicht hat, bleiben, wie die Abfrage sie gespeichert hat. Bevor es gespeicherte Zahlen zeigt, prüft es den Zurücksetzungszeitpunkt jedes Zeitfensters: Ein Fenster, das seit dem Speichern der Zahlen zurückgesetzt wurde, zeigt 0 %, weil der alte Prozentsatz nicht mehr gilt. Deshalb zeigt `personal` oben `0%` ohne Zurücksetzungszeit.

Der Endpunkt hat zwei Antwortformate. Das neuere hat ein Array `limits` mit den Einträgen `session` (5 Stunden), `weekly_all` und `weekly_scoped` (pro Modell). Das ältere hat die Objekte `five_hour`, `seven_day`, `seven_day_opus` und `seven_day_sonnet`. `ParseLive` liest beide: Das Array `limits` hat Vorrang, und die alten Objekte ergänzen, was es nicht abdeckt. Eine Antwort ganz ohne Limits wird als geändertes Format gemeldet, worauf `ccam doctor` dich hinweist.

### Der Limit-Eintrag

Wenn Claude Code einen Durchgang an einem Nutzungslimit stoppt, liest der Hook des [ccam-Plugins](/docs/install#claude-code-plugin) die abgelehnte Anfrage aus der Unterhaltung: welches Zeitfenster voll war (5 Stunden, Woche, Opus oder Sonnet) und wann es zurückgesetzt wird. Das schreibt er nach `~/.local/state/ccam/usage/<name>.limit.json`. `Fetch` legt den Eintrag über alles, was es sonst weiß, aus jeder Quelle: Jedes darin genannte Zeitfenster zählt bis zu seinem Zurücksetzen als 100 %, egal wie alt oder niedrig die gespeicherten Zahlen sind, sodass `ccam next --auto`, `ccam pick`, `--fallback` und [`auto_switch`](/features/next#auto-switch) das Konto bis dahin nicht wieder wählen. Ist sonst nichts bekannt, zeigt das Konto allein den Eintrag (`source: the usage limit Claude Code reported` mit `-v`). Sitzungen, an denen niemand sitzt (`-p`, Hintergrundsitzungen), schreiben keinen Eintrag, ebenso wenig der Stopp einer Organisation (Nutzung aufgebraucht, eine auf $0 gesetzte Gruppe), den ihre Admins jederzeit aufheben können. `ccam login` in ein anderes Claude-Konto verwirft den Eintrag; die Datei zu löschen, entfernt ihn ebenfalls.

### Tokens erneuern

Refresh-Tokens rotieren: Wer eines erneuert, bekommt ein neues Refresh-Token, und das alte wird ungültig. Würde ccam die Anmeldung erneuern, während Claude Code als dasselbe Konto läuft, würde diese Sitzung abgemeldet. Deshalb liest `ccam usage` Tokens und schreibt sie nie, es sei denn, du gibst `--refresh` an. Dann geht `Fetcher.refresh` so vor:

1. Es lehnt Konten ab, die mit `ccam import` an Ort und Stelle registriert wurden, da einfaches `claude` diese Anmeldung verwenden könnte, ohne dass ccam davon weiß;
2. es lehnt Anmeldungen ab, die im macOS-Schlüsselbund gespeichert sind;
3. es nimmt eine Sperrdatei, `~/.local/state/ccam/locks/profile-<name>.refresh.lock`, damit zwei ccam-Prozesse nie dasselbe Token zweimal erneuern; solange sie gehalten wird, zeigt das Konto `token not renewed: another ccam is using this login right now (a refresh, rename, remove or logout)`. `ccam rename`, `remove` und `logout` nehmen sie ebenfalls und lehnen ab, solange eine Erneuerung läuft;
4. unter der Sperre lehnt es ab, wenn eine Sitzung des Kontos läuft, egal wie sie gestartet wurde (über ccam, mit einfachem `claude`, aus einer IDE oder als Hintergrundsitzung samt dem Daemon dahinter, der die Anmeldung selbst erneuert), und überspringt die Anfrage, wenn ein anderer Prozess das Token bereits erneuert hat;
5. es sendet nichts, solange Claude Code seine eigenen Erneuerungssperren im Kontoordner hält (`.oauth_refresh.lock`, `<folder>.lock`): Claude Code erneuert diese Anmeldung gerade;
6. es sendet das Refresh-Token an den Token-Endpunkt und wartet bis zu 30 Sekunden auf die Antwort, so lange, wie Claude Code seinen eigenen Erneuerungen gibt. Eine einmal gesendete Anfrage wird nie abgebrochen, auch nicht mit Ctrl+C, da der Server das Refresh-Token vielleicht schon rotiert hat und nur noch seine Antwort die Anmeldung enthält: Ctrl+C oder ein Auflegen beendet die anderen Abfragen, und ccam beendet sich (nach Ctrl+C mit 130), sobald die Antwort gespeichert ist;
7. es speichert unter der Schreibsperre von Claude Code (`.storage-write.lock`) und ändert nur die Token-Felder von `.credentials.json` in ihrem dann aktuellen Stand: Was Claude Code inzwischen geschrieben hat (eine MCP-Anmeldung), bleibt, eine Anmeldung, die Claude Code inzwischen erneuert hat, wird behalten, und eine Abmeldung wird nie rückgängig gemacht. Nennt die Antwort die Lebensdauer des neuen Refresh-Tokens, wird sie wie bei Claude Code als `refreshTokenExpiresAt` gespeichert.

Eine Erneuerung, die nicht gelingt, nennt den Grund, in der Kontozeile oder über den gespeicherten Zahlen: zum Beispiel `token not renewed: a session is running`, oder `login may be lost: …`, wenn eine Anfrage gesendet wurde, ihre Antwort aber nicht verwendet oder gespeichert werden konnte, mit dem Rat, `ccam login <name>` auszuführen, falls Claude Code dich dann zur Anmeldung auffordert.

### Der Statuszeilen-Sammler

Claude Code führt den Statuszeilenbefehl aus und leitet ihm ein JSON-Dokument zu. ccam liest daraus den Modellnamen und `rate_limits`:

<!-- motion:usage-statusline -->

1. `runStatusline` findet das Konto über `CCAM_PROFILE`, das ccam beim Start von Claude Code setzt, oder sonst durch Abgleich von `CLAUDE_CONFIG_DIR` mit den Kontoordnern.
2. `ParseStatusLine` liest `rate_limits.five_hour` und `rate_limits.seven_day`. Bei Konten mit API-Schlüssel oder vor der ersten Antwort der Sitzung enthält das JSON keine Ratenlimits; dann wird nichts gespeichert.
3. `SaveStatusLine` schreibt `~/.local/state/ccam/usage/<name>.statusline.json`. `ccam usage` verwendet diese Datei, wenn der Endpunkt nicht antwortet.
4. ccam gibt `[name] · model · 5h N% · 7d N%` aus. Mit `--wrap` führt es deinen Befehl über `sh -c` mit derselben Eingabe aus und gibt seine Ausgabe so aus, wie Claude Code eine Statuszeile zeigt: jede nicht leere Zeile, mit `[name]` vor der ersten. Ein fehlschlagender Befehl zeigt nur den Kontonamen. Unter Windows führt es den Befehl so aus wie Claude Code: mit Git Bash, wenn es installiert ist, sonst mit PowerShell.
5. Wenn die 5-Stunden- oder Wochennutzung bei 90 % oder mehr liegt und es ein weiteres Konto gibt, endet die Zeile mit `→ ccam next --auto --continue`, zum Beispiel `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

Der umschlossene Befehl hat kein eigenes Zeitlimit: Claude Code gibt einem Statuszeilenbefehl 10 Minuten und stoppt einen Lauf, wenn eine neuere Aktualisierung beginnt, sodass ein langsamer Befehl sich so zeigt wie ohne ccam. Wenn Claude Code ccam beendet, endet auch der Befehl (unter Linux und macOS mit jedem Prozess, den er gestartet hat).

`ccam statusline install` schreibt `statusLine` in die `settings.json` jedes Kontos und behält die anderen Schlüssel des Objekts, etwa `padding`, `refreshInterval` und `hideVimModeIndicator`. Hat ein Konto bereits einen Statuszeilenbefehl, behält ccam ihn als `--wrap '<old command>'` bei; unter Windows als `--wrap-base64 <base64>`, was in keiner der beiden Shells Anführungszeichen braucht. Eine Statuszeile, die bereits ein ccam ausführt, von einem anderen Pfad oder so in Anführungszeichen gesetzt, wie ein älteres ccam sie geschrieben hat, wird ersetzt statt umschlossen, wobei erhalten bleibt, was sie umschließt. Unter Windows wird der Pfad zu ccam mit Schrägstrichen und `~` für den Benutzerordner geschrieben (`~/.local/bin/ccam.exe statusline`), was sowohl Git Bash als auch PowerShell verstehen. Ein Pfad, der Anführungszeichen braucht, wird vollständig für die Shell geschrieben, die Claude Code verwendet: `'…'` für Git Bash, `& '…'` für PowerShell. Eine über eine Gruppe geteilte `settings.json` wird einmal über ihren Link geschrieben, sodass jedes Konto der Gruppe die Einstellung bekommt. Ein erneuter Aufruf gibt `(already set)` für Dateien aus, die ccam bereits verwenden.

## Sicherheit und Sonderfälle

- `ccam usage` ändert nie eine Token-Datei, es sei denn, du gibst `--refresh` an. Ist ein Zugriffstoken abgelaufen, starte das Konto einmal mit `ccam <name>`, dann erneuert Claude Code es. Eine Anmeldung, die als Ganzes abgelaufen ist, braucht `ccam login <name>`.
- `--refresh` fasst nie laufende Konten (auch wenn Claude Code ohne ccam gestartet wurde), importierte Konten oder Anmeldungen im Schlüsselbund an. Eine abgelehnte Erneuerung zeigt die gespeicherten Zahlen und den Grund (`token not renewed: …`).
- Der Nutzungs-Endpunkt ist nicht öffentlich dokumentiert und kann sich ändern. Wenn das passiert, siehst du weiterhin gespeicherte Zahlen, und `ccam doctor` meldet die Formatänderung.
- Ohne gespeicherte Zahlen und ohne Antwort zeigt die Kontozeile `lookup failed:` und den Grund, zum Beispiel `no saved usage` mit `--source cache`.
- Die Cache-Dateien enthalten nur Prozentsätze und Zurücksetzungszeiten, nie Tokens.

## Verwandte Themen

- [Nutzungslimits und Kontowechsel](/docs/usage-limits)
- [Zum nächsten Konto wechseln](/features/next)
- [Umgang mit Nutzungslimits](/guides/usage-limit-workflow)
- [Befehle](/docs/commands)
