---
title: Konten verwalten
description: "Claude Code-Konten anlegen, anmelden, auflisten, umbenennen, entfernen und importieren: mit ccam add, login, list, rename, remove und import."
updated: 2026-10-20
order: 2
translation: machine
sourceHash: 9eb992214132
---

Jedes ccam-Konto ist ein eigener Claude Code-Konfigurationsordner mit eigener Anmeldung. Diese Befehle legen Konten an, melden sie an, listen sie auf, benennen sie um, entfernen sie und registrieren Ordner, die du bereits hast.

## Verwendung

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Wo ein Befehl `<name|number>` erwartet, kannst du den Kontonamen oder seine Slot-Nummer angeben, also die Spalte `#` von `ccam list`.

`ccam add` ohne Optionen öffnet in einem Terminal einen Assistenten. Er fragt nach dem Namen, zeigt dann eine Checkliste der geteilten Elemente (alles außer dem Gesprächsverlauf ist angehakt), fragt anschließend `Create a shortcut command for work? [y/N]` und bei `y` dann `Shortcut name [ccwork]:`. Jede der Optionen `--share`, `--shortcut` oder `--no-login` oder eine Pipe statt eines Terminals überspringt den Assistenten; dann ist ein Name Pflicht, und ein Kurzbefehl wird nur angelegt, wenn `--shortcut` einen nennt.

`ccam list` zeigt Slot, Name, Kurzbefehl, E-Mail-Adresse, geteilte Elemente und letzte Verwendung, mit `*` vor dem Standardkonto:

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

Mit `--json` ist jedes Konto ein Objekt mit `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (fehlt, wenn das Konto nie gestartet wurde) und `share`. `share` ordnet jedem Element (`mcp`, `settings`, `memory`, `extensions`, `history`) seinen Gruppennamen oder `"local"` zu.

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `ccam add --no-login` | Das Konto ohne Anmeldung anlegen; später mit `ccam login <name>` anmelden | aus |
| `ccam add --share <item=group>` | Ein Element mit einer Gruppe teilen, oder mit `item=local` privat halten. Mehrfach verwendbar | `mcp`, `settings`, `memory`, `extensions` in `default`; `history` lokal |
| `ccam add --shortcut <command>` | Zusätzlich diesen Kurzbefehl in `~/.local/bin` anlegen | keiner |
| `ccam add --api --key-from <source>` | Ein Konto, das einen API-Schlüssel statt einer Anmeldung verwendet; siehe [API-Schlüssel-Konten](/features/api) | aus |
| `ccam list --json` | Die Konten als JSON ausgeben | aus |
| `ccam remove -y`, `--yes` | Ohne Nachfrage entfernen | aus |
| `ccam import --name <name>` | Kontoname für den angegebenen Ordner | aus dem Ordnernamen abgeleitet |
| `ccam import -y`, `--yes` | Jeden gefundenen Ordner ohne Nachfrage importieren | aus |
| `--dry-run` | Wird von `ccam list` und `ccam completion` beachtet, die nichts ändern. Abgelehnt von `add`, `login`, `logout`, `use`, `rename`, `remove` und `import`, mit `` --dry-run is not supported by `ccam add`; nothing was changed `` (mit dem jeweiligen Befehl) | aus |

`ccam login`, `ccam logout`, `ccam use` und `ccam rename` haben keine eigenen Optionen. `ccam completion` nimmt `bash`, `zsh` oder `powershell`; alles andere schlägt mit `only bash, zsh and powershell are supported` fehl.

Kontonamen bestehen aus 1 bis 32 Zeichen aus `a-z`, `0-9`, `-` und `_` und beginnen mit einem Buchstaben oder einer Ziffer. Ein Name darf nicht nur aus Ziffern bestehen, weil `ccam 2` Slot 2 bedeutet. Ein Name darf auch kein ccam-Befehl sein (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `folder`, `pick`, `ext`), damit `ccam <name>` nie mehrdeutig ist. Die Liste ist `config.Reserved` in `internal/config`.

## So funktioniert es

`ccam add work --shortcut ccwork` führt diese Schritte in `addProfile` (`internal/cli/profiles.go`) aus:

<!-- motion:accounts-add -->

1. `config.ValidateName` prüft den Namen, und `shim.Check` prüft den Kurzbefehl, bevor irgendetwas angelegt wird.
2. `profile.Add` legt `~/.local/share/ccam/profiles/work/` mit Modus `0700` an und gibt dem Konto die niedrigste freie Slot-Nummer (`Config.NextSlot`), sodass die Nummer eines entfernten Kontos wiederverwendet wird. Das erste Konto wird zum Standardkonto.
3. `applyShare` verlinkt den neuen, leeren Ordner mit `share.Build` und anschließend `Apply` in seine Gruppen, zum Beispiel `settings.json` auf `groups/default/settings.json`. Schlägt das fehl, wird der neue Ordner entfernt, weil ein übrig gebliebener Ordner jeden weiteren Versuch mit `folder already exists` blockieren würde.
4. Das Konto wird über `fsutil.AtomicWrite` in `~/.config/ccam/config.toml` gespeichert: eine temporäre Datei, `fsync`, dann `rename`, sodass ein Absturz nie eine halbe Konfiguration hinterlässt.
5. `runLogin` startet `claude auth login` mit `CLAUDE_CONFIG_DIR` auf den Kontoordner gesetzt. Schlägt die Anmeldung fehl, bleibt das Konto bestehen, und ccam gibt den Fehler aus und rät dir, später `ccam login work` auszuführen.
6. `profile.ReadEmail` liest `oauthAccount.emailAddress` aus `profiles/work/.claude.json`. Diese Datei kann mehrere Megabyte groß sein, daher sucht `TopLevelValue` nur nach diesem einen Schlüssel, statt die ganze Datei zu parsen.
7. `shim.Install` verlinkt `~/.local/bin/ccwork` mit der ccam-Binärdatei, die beim Start `argv[0]` in `config.toml` nachschlägt.

Eine Anmeldung kann Minuten dauern, und ein anderes ccam kann die Konfiguration in der Zwischenzeit speichern. Deshalb speichern Befehle, die warten, über `env.update`: Es lädt `config.toml` neu, wendet die eine Änderung an (hier die E-Mail-Adresse) und speichert. `ccam remove` lädt aus demselben Grund nach seiner Nachfrage neu.

`ccam login` führt Schritt 5 und 6 für ein vorhandenes Konto aus. `ccam logout` führt `claude auth logout` im Konto aus und löscht die gespeicherte E-Mail-Adresse. `ccam use` schreibt nur `default` in `config.toml` neu.

### Anmeldung auf einem entfernten Rechner

Über SSH kann der Browserteil von `claude auth login` scheitern. Claude Code wartet auf die Anmeldung über einen Callback-Server unter `http://localhost:<port>/callback` auf dem Server, aber der Browser, der die Anmeldung abschließt, läuft auf deinem Computer. Dort führt diese Adresse ins Leere, und der Browser zeigt einen Fehler wie `ERR_CONNECTION_REFUSED`. Die Adresse enthält trotzdem den Einmalcode, also kann ccam ihn zustellen:

<!-- motion:accounts-relay -->

`runLogin` übergibt dein Terminal nicht an Claude Code. Es verbindet die Standardeingabe von Claude Code mit einer Pipe, und `relayLogin` kopiert jede Zeile, die du tippst, hinein, sodass ein Code, den du bei der Eingabeaufforderung `Paste code here` von Claude Code einfügst, trotzdem ankommt. Ist eine Zeile eine Callback-Adresse, akzeptiert `callbackURL` sie nur, wenn der Host `localhost`, `127.0.0.1` oder `[::1]` ist, ein Port angegeben ist, der Pfad `/callback` lautet und sowohl `code` als auch `state` vorhanden sind; `http://` darf fehlen. `deliverCallback` ruft diese Adresse dann auf dem Server mit einem Limit von 10 Sekunden und ohne Weiterleitungen zu folgen ab, so wie es der Browser getan hätte. Die Anfrage erreicht den Callback-Server des wartenden `claude auth login` (die gestrichelte Linie im Diagramm). Claude Code erhält den Code und speichert die Anmeldung, unter Linux in `.credentials.json` im Kontoordner. Eine eingefügte Zeile kann ccam nie dazu bringen, einen anderen Host zu kontaktieren.

ccam gibt dies aus, bevor die Anmeldung beginnt:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

Bei Erfolg gibt es `→ Passed the sign-in to Claude Code on this machine.` aus.

Während die Anmeldung in einem Terminal auf eine Eingabe wartet, schaltet ccam das Echo und den Zeilenmodus des Terminals aus und zeigt für jedes getippte Zeichen ein `*`, sodass ein eingefügter Code oder eine Callback-Adresse nicht auf dem Bildschirm erscheint. `readLine` übernimmt die Zeilenbearbeitung selbst: `Enter` sendet die Zeile, `Backspace` löscht ein Zeichen, und `Ctrl+D` in einer leeren Zeile beendet die Eingabe. Das Terminal wird wiederhergestellt, wenn die Anmeldung endet, und auch bei `Ctrl+C`, `SIGTERM` oder `SIGHUP`, die ccam weiterhin beenden. Ohne Terminal an Standardeingabe und Standardfehlerausgabe werden Zeilen so kopiert, wie sie ankommen.

### Umbenennen, entfernen und importieren

`ccam rename` und `ccam remove` wirken auf den Ordner des Kontos und auf das, was ccam an anderer Stelle unter dem Namen des Kontos aufbewahrt. `nameState` listet diese Pfade auf: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` und die Nutzungscaches `~/.local/state/ccam/usage/<name>.json` und `<name>.statusline.json`. `ccam remove work` behandelt sie in dieser Reihenfolge:

<!-- motion:accounts-remove -->

`ccam remove` lehnt zuerst ab, wenn das Konto läuft, und fragt dann nach. Teilt das Konto Erweiterungen, schreibt `RewritePluginPaths` die Dateien `installed_plugins.json` und `known_marketplaces.json` der Gruppe neu, weil Claude Code absolute Plugin-Pfade speichert, die durch den Ordner dieses Kontos führen können. `profile.Trash` verschiebt den Ordner nach `~/.local/share/ccam/trash/<name>-<time>/`, das Konto wird aus `config.toml` entfernt, die `nameState`-Pfade werden gelöscht, damit ein späteres Konto mit demselben Namen sie nicht erbt, und der Link des Kurzbefehls wird entfernt. Zum Wiederherstellen verschiebst du den Ordner aus dem Papierkorb zurück.

`ccam rename` benennt den Ordner um, schreibt die Plugin-Pfade in seinem `plugins/` vom alten auf den neuen Ordner um, übernimmt die Standardeinstellung und verschiebt die `nameState`-Pfade auf den neuen Namen. Der Kurzbefehl behält seinen Namen. Schlägt das Umschreiben der Plugins fehl, wird der Ordner zurückbenannt, da die Konfiguration nicht gespeichert wird.

`ccam import` registriert einen Ordner dort, wo er liegt, als `dir` in `config.toml`. Es verschiebt oder kopiert ihn nie: Ein kopiertes Refresh-Token würde dazu führen, dass sich die Kopien gegenseitig abmelden. Ohne Ordnerargument sucht es nach nicht registrierten Ordnern `~/.claude` und `~/.claude-*` und fragt bei jedem nach; die Namen ergeben sich aus dem Ordner (`~/.claude` wird `main`, `~/.claude-work` wird `work`, `~/.claude-2` wird `acct2`), und ein Ordner, aus dem sich kein Name ableiten lässt, wird mit dem passenden `--name`-Befehl aufgelistet. Ohne `CLAUDE_CONFIG_DIR` hält Claude Code `.claude.json` in deinem Home-Ordner neben `~/.claude`. Wenn du also `~/.claude` importierst und es keine eigene `.claude.json` hat, legt ccam den Link `~/.claude/.claude.json` auf `~/.claude.json` an. Einfaches `claude` und `ccam main` sehen dann dasselbe Konto. Importierte Konten teilen nichts, bis du `ccam share` ausführst, und es wird kein Kurzbefehl angelegt.

`ccam completion bash|zsh` gibt ein Vervollständigungsskript aus, das auch Kontonamen vervollständigt.

## Sicherheit und Sonderfälle

- `ccam rename` und `ccam remove` lehnen ab, solange das Konto läuft (`refuseRunning`), mit `work is running; quit its Claude Code session first`, weil die Sitzung weiter in den Ordner schreibt.
- `ccam remove` hebt bei einem importierten Konto nur die Registrierung auf: `✓ main unregistered (folder ... left in place)`. Dein `~/.claude` wird nicht in den Papierkorb verschoben.
- `ccam add` lehnt einen Namen ab, dessen Ordner in `profiles/` bereits existiert, und verweist auf `ccam import`.
- `ccam import` lehnt einen bereits registrierten Ordner, einen Pfad, der kein Ordner ist, und einen vergebenen Namen ab.
- `ccam logout` und `ccam login` prüfen nicht, ob das Konto läuft.
- Wenn die Weiterleitung Claude Code nicht erreicht, zum Beispiel nachdem die Anmeldung beendet wurde, gibt sie `ccam: could not pass the sign-in to Claude Code:` mit dem Grund aus. Ein HTTP-Fehler bedeutet meist, dass der Code abgelaufen ist; starte die Anmeldung neu.

## Verwandte Themen

- [Claude Code starten](/features/launch)
- [Kurzbefehle](/features/shortcuts)
- [MCP-Server, Einstellungen und Verlauf teilen](/docs/sharing)
- [Mehrere Claude Code-Konten](/guides/multiple-claude-code-accounts)
