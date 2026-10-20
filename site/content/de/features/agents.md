---
title: Hintergrundsitzungen aller Konten
description: ccam agents listet die Hintergrund-Agents und offenen Sitzungen von Claude Code aller Konten in einer Tabelle, und ccam attach öffnet eine anhand ihrer ID.
updated: 2026-10-20
order: 16
translation: machine
sourceHash: "1f459764d66d"
---

Claude Code führt Hintergrundsitzungen (`claude --bg`, Agent-Ansicht) pro Konto: `claude agents` zeigt nur die Sitzungen des Kontos, als das es läuft. `ccam agents` fragt alle Konten gleichzeitig ab und listet sie zusammen auf, und `ccam attach <id>` öffnet eine Hintergrundsitzung mit dem Konto, das sie gestartet hat.

## Verwendung

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique start of the id is enough
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` ist der Zustand der Hintergrundsitzung (`working`, `blocked`, `done`, `failed`, `stopped`) und das, was ihr Prozess gerade tut (`busy`, `waiting`, `idle`); eine offene interaktive Sitzung zeigt `open`.

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `--all` | Auch Hintergrundsitzungen auflisten, die fertig, fehlgeschlagen oder angehalten sind | aus |
| `--json` | Die Sitzungen als JSON ausgeben, jede mit dem `account`, zu dem sie gehört | aus |

`ccam attach <id> [claude args]` übergibt die Argumente nach der ID an `claude attach`. Mit `--dry-run` gibt es den Befehl aus.

## So funktioniert es

1. Für jedes Konto führt ccam `claude agents --json` (mit `--all`, wenn angegeben) mit dem `CLAUDE_CONFIG_DIR` des Kontos aus, alle gleichzeitig, jedes innerhalb von 15 Sekunden. Das liest die eigene Sitzungsliste des Kontos und startet nicht den Hintergrunddienst von Claude Code.
2. Jeder Eintrag erhält den Namen des Kontos, und die Liste wird nach Startzeit sortiert, die neueste zuerst. Ein Konto, das nicht antwortet, wird auf der Standardfehlerausgabe genannt und ausgelassen, und ccam beendet sich dann mit Code 1, damit ein Skript einen Fehler nicht für eine leere Liste hält.
3. `ccam attach` sucht im Ordner jedes Kontos nach `jobs/<id>`, wo Claude Code eine Hintergrundsitzung ablegt, und führt `claude attach <id>` als dieses Konto aus, so wie `ccam <name>` Claude Code startet.

## Sicherheit und Sonderfälle

- Die IDs von Hintergrundsitzungen gehören zu einem Konto: Ein `claude attach`, das als anderes Konto läuft, findet sie nicht, und genau das löst `ccam attach`.
- Zwei Konten könnten Sitzungen haben, deren IDs gleich beginnen; `ccam attach` nennt dann beide und bittet um mehr von der ID.
- ccam übergibt seine geteilten MCP-Server mit `--mcp-config=<file>` in einem Argument an Claude Code, was `claude agents` und `claude attach` akzeptieren; Hintergrundsitzungen verwenden diese Datei weiter, wenn sie neu starten.
- Der Hintergrunddienst eines API-Schlüssel-Kontos nimmt seinen Schlüssel aus der Shell, die ihn gestartet hat. Starte die erste Hintergrundsitzung eines solchen Kontos über ccam (`ccam gw --bg …`), damit der Dienst den Schlüssel hat.

## Verwandte Themen

- [Claude Code starten](/features/launch)
- [Ein Konto für Skripte wählen](/features/pick)
- [Das Dashboard](/features/dashboard)
