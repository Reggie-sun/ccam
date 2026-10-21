---
title: ccam aktualisieren
description: ccam update ersetzt ccam nach einer SHA-256-Prüfung durch das neueste GitHub-Release; --check meldet nur, und --version installiert ein bestimmtes Release-Tag.
updated: 2026-10-21
order: 10
translation: machine
sourceHash: a92dbf37399f
---

`ccam update` lädt ein ccam-Release von GitHub herunter, prüft es gegen die `checksums.txt` des Releases und ersetzt die laufende Binärdatei mit einem einzigen Umbenennen. Verwende es, wenn du ccam mit dem Installationsskript oder aus einem Archiv installiert hast. Auch das [Dashboard](/features/dashboard) fragt bei jedem Öffnen GitHub nach dem neuesten Release und bietet es mit `U` an, das nach einer Rückfrage dieselben Schritte ausführt. Kein anderer Befehl kontaktiert GitHub, daher funktionieren sie offline genauso.

## Verwendung

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, also an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

Gibt es kein neueres Release, gibt ccam `ccam 0.4.0 is up to date.` aus und ändert nichts. Mit `--version` und der Version, die du bereits verwendest, gibt es `ccam 0.4.0 is already installed.` aus.

## Optionen

| Option | Wirkung | Standard |
| --- | --- | --- |
| `--check` | Nur melden, ob ein neueres Release existiert | aus |
| `--version <tag>` | Dieses Release-Tag installieren, etwa `v0.3.0`, neuer oder älter als das aktuelle | neuestes |
| `--dry-run` | Wird beachtet: ermittelt das neueste Tag (sofern `--version` keines nennt), führt dieselben Prüfungen an der Binärdatei aus und gibt Download-, Prüfsummen- und Ersetzungspfad aus, ohne herunterzuladen | aus |

`--check` und `--version` lassen sich nicht zusammen verwenden. Das Tag muss die Form `vMAJOR.MINOR.PATCH` haben, mit optionaler Vorabversions- oder Build-Endung; das führende `v` darf fehlen. Alles andere schlägt fehl, zum Beispiel `--version takes a release tag such as v0.3.0, not "1.2"`.

Die Umgebungsvariable `CCAM_RELEASES_URL` ersetzt `https://github.com/jhyunwoo/ccam/releases`, für einen Mirror oder für Tests. Sie muss gleich aufgebaut sein: `<url>/latest` leitet auf `<url>/tag/<tag>` weiter, und die Dateien liegen unter `<url>/download/<tag>/<asset>`.

## So funktioniert es

Der Code liegt in `internal/update`, das nur `ccam update` importiert. Der Startpfad ruft ihn nie auf.

1. `cli.updateCmd` findet die laufende Binärdatei mit `os.Executable` und löst symbolische Links mit `filepath.EvalSymlinks` auf. Dieser echte Pfad wird ersetzt.
2. Ohne `--version` ruft `Updater.Latest` die Adresse `<releases>/latest` ab, ohne der Weiterleitung zu folgen, und liest das Tag aus dem Header `Location` (`…/tag/v0.5.0`). Dabei wird die Release-Seite statt der GitHub-API verwendet, daher gilt das Limit der API von 60 Anfragen pro Stunde nicht. Die Abfrage hat 10 Sekunden Zeit.
3. `parseVersion` prüft das Tag gegen ein strenges Muster. Ein Tag landet in Download-URLs und in einem ausgegebenen Shell-Befehl, daher kommt nichts mit `/`, `%`, Leerzeichen oder Shell-Zeichen durch. `compareVersions` vergleicht es mit der laufenden Version; eine Vorabversion wird vor ihrem Release einsortiert.
4. `managed()` und `writable()` entscheiden, ob ccam die Binärdatei überhaupt ersetzen darf (siehe unten).
5. `Updater.Fetch` lädt `ccam_<os>_<arch>.tar.gz` (unter Windows `ccam_windows_<arch>.zip`) und `checksums.txt` gleichzeitig herunter, mit 2 Minuten für beide. Jeder Download ist auf 256 MB begrenzt.
6. Die SHA-256-Prüfsumme des Archivs muss mit seiner Zeile in `checksums.txt` übereinstimmen, bevor irgendetwas anderes geschieht. Eine fehlende Zeile oder eine abweichende Prüfsumme bricht das Update mit `nothing was changed` ab.
7. `extract()` liest das tar-Archiv und nimmt die erste reguläre Datei namens `ccam` (`extractZip` nimmt `ccam.exe` aus dem zip). Ordner, Links und andere Einträge werden ignoriert, und auch die Binärdatei ist auf 256 MB begrenzt.
8. `Updater.Install` schreibt die neue Binärdatei mit `fsutil.AtomicWrite`: eine temporäre Datei `.ccam.tmp*` neben der Binärdatei, mit den Berechtigungsbits der alten Datei, auf die Festplatte geschrieben und über die alte umbenannt.

Die gesamte Abfolge, von der Tag-Abfrage bis zum Umbenennen:

<!-- motion:update-pipeline -->

Das Archiv und `checksums.txt` werden nebeneinander heruntergeladen, und nichts wird geschrieben, bis die Prüfsumme stimmt. Das Umbenennen ist atomar, daher hinterlässt ein Absturz oder eine volle Festplatte entweder die alte oder die neue Binärdatei, nie eine halbe. Kurzbefehle wie `~/.local/bin/ccwork` sind symbolische Links auf die ccam-Binärdatei, daher führen sie ohne jede Änderung die neue Version aus.

Unter Windows lässt sich die Datei eines laufenden Programms nicht ersetzen oder löschen, nur umbenennen. `Updater.Install` schreibt die neue Binärdatei nach `ccam.exe.new`, benennt das laufende `ccam.exe` in `ccam.exe.old` um und die neue Datei in `ccam.exe`. Schlägt der letzte Schritt fehl, wird die alte Binärdatei zurückbenannt. Das nächste Update löscht `ccam.exe.old`. Kurzbefehle, die Hardlinks statt symbolischer Links sind ([Kurzbefehle](/features/shortcuts)), werden direkt danach erneuert, sodass auch sie die neue Version ausführen.

### Über Paketverwaltung installiert

Ein aus dem `.deb` installiertes ccam gehört apt, und es an apt vorbei zu ersetzen, würde die Paketdatenbank verfälschen. Vor dem Download prüft ccam, wem die Binärdatei gehört:

<!-- motion:update-managed -->

1. Eine Binärdatei in `/usr/bin`, `/bin`, `/usr/sbin` oder `/sbin` gilt als paketverwaltet.
2. Anderswo gilt sie ebenfalls als paketverwaltet, wenn `dpkg` installiert ist und `dpkg -S <path>` erfolgreich ist.
3. Andernfalls legt ccam eine temporäre Datei `.ccam-update-*` neben der Binärdatei an und entfernt sie wieder, um sicherzugehen, dass es dort schreiben kann.

Bei einer paketverwalteten Binärdatei hält ccam an und gibt die Schritte für apt aus:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

Auf anderen Systemen als Linux wird nur die Zeile mit dem Installationsskript ausgegeben: unter Windows `irm https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.ps1 | iex`. Ist der Ordner nicht beschreibbar, schlägt ccam vor, das Installationsskript erneut auszuführen, das ccam in `~/.local/bin` installiert.

## Sicherheit und Sonderfälle

- Ein Entwicklungs-Build meldet seine Version als `dev` und hat kein Release, mit dem er sich vergleichen ließe, daher lehnen `ccam update` und `ccam update --check` ihn ab. `ccam update --version vX.Y.Z` installiert ein Release darüber.
- Ein fehlgeschlagener Download, eine fehlende Datei für deine Plattform oder eine abweichende Prüfsumme lassen die installierte Binärdatei unverändert.
- Ein Release, das es nicht gibt, schlägt mit `there is no release <tag>` fehl. Lehnt GitHub mit HTTP 403 oder 429 ab, sagt ccam, dass es sich wahrscheinlich um ein Rate-Limit handelt, und bittet dich, es später erneut zu versuchen.
- `--check` liest nur. Es prüft nicht, ob die Binärdatei paketverwaltet oder beschreibbar ist.
- `--version` mit einem älteren Tag ist der Weg zurück zu einer früheren Version.

## Verwandte Themen

- [Installieren](/docs/install)
- [Befehle](/docs/commands)
- [FAQ](/docs/faq)
