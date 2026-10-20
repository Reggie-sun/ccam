---
title: Konten nach Ordner
description: ccam folder lässt einen Ordner und alles darunter ein Konto nutzen, wenn du Claude Code ohne Kontonamen startest, etwa work für den Code deiner Firma.
updated: 2026-10-20
order: 11
translation: machine
sourceHash: 5071b04a194d
---

`ccam folder work ~/code/company` sorgt dafür, dass jeder Start ohne Kontonamen in `~/code/company` und in jedem Ordner darunter Claude Code als `work` startet. Außerhalb davon startet wie bisher das Standardkonto. Wenn du ein Konto nennst, mit `ccam personal` oder einem Kurzbefehl wie `ccpersonal`, startet immer dieses Konto.

## Verwendung

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

Diese Starts verwenden das Konto des Ordners:

- `ccam -p "…"` und jeder andere Start, dessen erstes Argument eine Option für Claude Code ist
- `ccam run` ohne Namen und `ccam --dry-run`
- `ccam --fallback -p "…"` für seinen ersten Versuch
- die Auswahl, `ccam` allein, die mit dem Cursor auf dem Konto des Ordners beginnt

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `<account>` | Name oder Nummer des zu verwendenden Kontos | |
| `[folder]` | Der Ordner, relativ oder absolut | der aktuelle Ordner |
| `--remove` | Den Eintrag für genau diesen Ordner entfernen | aus |
| `--dry-run` | Wird beachtet: die Änderung ausgeben, ohne sie zu speichern | aus |

## So funktioniert es

Jeder Ordner ist ein `[[folder]]`-Eintrag in `~/.config/ccam/config.toml`, der Datei, die ohnehin jeder Start liest:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder` wandelt den Ordner in einen absoluten Pfad mit aufgelösten symbolischen Links um und speichert ihn. Wird ein Ordner erneut festgelegt, ersetzt das sein Konto.
2. Ein Start ohne Namen ruft `launcher.HereProfile` auf. Ohne Einträge nimmt er sofort das Standardkonto.
3. Andernfalls vergleicht er das Arbeitsverzeichnis als Zeichenkette mit jedem `path`: der Ordner selbst oder einer darunter, nie ein Ordner, der nur mit denselben Buchstaben beginnt (`/code/company` umfasst nicht `/code/companyx`). Der tiefste Eintrag gewinnt, sodass `~/code` `personal` sein kann, während `~/code/company` `work` ist.
4. Nur wenn das Arbeitsverzeichnis so, wie es ist, auf nichts passt, löst ccam seine symbolischen Links auf und vergleicht noch einmal.
5. Unter Windows ignoriert der Vergleich Groß- und Kleinschreibung, wie Windows selbst.

Außer `config.toml` wird keine Datei gelesen, daher kostet ein Start in einem Ordner keine messbare Zeit (4,0 ms mit oder ohne Einträge, `hyperfine`, 300 Durchläufe).

## Sicherheit und Sonderfälle

- `ccam rename` überträgt die Ordnereinträge auf den neuen Namen, und `ccam remove` löscht die Einträge des entfernten Kontos.
- Ein Eintrag, dessen Konto nicht mehr existiert, hält einen Start ohne Namen in diesem Ordner an: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` meldet solche Einträge und Einträge, deren Ordner verschwunden ist.
- `ccam folder --remove` in einem Ordner, der sein Konto von einem darüberliegenden Ordner bekommt, nennt diesen Ordner, statt etwas zu entfernen.
- `ccam next` wählt nach Rotation oder Nutzung und berücksichtigt keine Ordner.

## Verwandte Themen

- [Claude Code starten](/features/launch)
- [Kurzbefehle](/features/shortcuts)
- [Konten verwalten](/features/accounts)
- [Befehlsreferenz](/docs/commands)
