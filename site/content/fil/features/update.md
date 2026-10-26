---
title: Pag-update ng ccam
description: Pinapalitan ng ccam update ang ccam ng pinakabagong GitHub release matapos ang SHA-256 check; nag-uulat lang ang --check; nag-i-install ng tag ang --version.
updated: 2026-10-26
order: 10
translation: machine
sourceHash: "598931e4801f"
---

Nagda-download ang `ccam update` ng ccam release mula sa GitHub, sinusuri ito laban sa `checksums.txt` ng release at pinapalitan ang tumatakbong binary sa isang rename. Gamitin ito kapag in-install mo ang ccam gamit ang install script o mula sa isang archive. Nagtatanong din ang [dashboard](/features/dashboard) sa GitHub tungkol sa pinakabagong release tuwing nagbubukas ito at iniaalok ito gamit ang `U`, na nagpapatakbo ng parehong mga hakbang pagkatapos magtanong. Hindi kailanman kumokontak sa GitHub ang ibang mga command, kaya gumagana sila nang pareho offline.

## Paggamit

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, even an older one
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

Kapag walang mas bagong release, ipi-print ng ccam ang `ccam 0.4.0 is up to date.` at walang babaguhin. Kung pinangalanan ng `--version` ang bersyong mayroon ka na, ipi-print nito ang `ccam 0.4.0 is already installed.`

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--check` | Iulat lang kung may mas bagong release | off |
| `--version <tag>` | I-install ang release tag na ito, gaya ng `v0.3.0`, mas bago man o mas luma kaysa sa kasalukuyan | pinakabago |
| `--dry-run` | Sinusuportahan: hinahanap ang pinakabagong tag (maliban kung may pinangalanan ang `--version`), pinapatakbo ang parehong mga pagsusuri sa binary, at ipi-print ang mga path ng download, checksum at papalitan nang hindi nagda-download | off |

Hindi puwedeng pagsabayin ang `--check` at `--version`. Dapat mukhang `vMAJOR.MINOR.PATCH` ang tag, na may opsyonal na pre-release o build suffix; puwedeng wala ang `v` sa unahan. Pumapalya ang anumang iba, halimbawa `--version takes a release tag such as v0.3.0, not "1.2"`.

Pinapalitan ng environment variable na `CCAM_RELEASES_URL` ang `https://github.com/jhyunwoo/ccam/releases`, para sa mirror o para sa mga test. Dapat pareho ang ayos nito: nagre-redirect ang `<url>/latest` sa `<url>/tag/<tag>`, at nasa `<url>/download/<tag>/<asset>` ang mga asset.

## Paano ito gumagana

Sinusuri ng ccam ang download bago nito galawin ang naka-install na binary, at pinapalitan ang binary sa isang rename. Nasa `internal/update` ang code, na ang `ccam update` lang ang nag-i-import. Hindi kailanman tinatawag ito ng launch path.

1. Hinahanap ng ccam ang tumatakbong binary gamit ang `os.Executable` at nire-resolve ang mga symlink gamit ang `filepath.EvalSymlinks` (`cli.updateCmd`). Ang totoong path na iyon ang pinapalitan.
2. Kapag walang `--version`, hinihiling ng ccam ang `<releases>/latest` nang hindi sinusundan ang redirect (`Updater.Latest`), at binabasa ang tag mula sa header na `Location` (`…/tag/v0.5.0`). Ginagamit nito ang releases page sa halip na ang GitHub API, kaya hindi nalalapat ang limit ng API na 60 request kada oras. Nagta-time out ang paghahanap pagkalipas ng 10 segundo.
3. Sinusuri ng ccam ang tag laban sa isang mahigpit na pattern (`parseVersion`). Napupunta ang tag sa mga download URL at sa isang ipini-print na shell command, kaya walang nakakalusot na may `/`, `%`, espasyo o shell character. Pagkatapos, ikinukumpara ito ng ccam sa tumatakbong bersyon (`compareVersions`); nauuna ang pre-release bago ang release nito.
4. Pinagpapasyahan ng ccam kung puwede ba nitong palitan ang binary (`managed()`, `writable()`; tingnan sa ibaba).
5. Sabay na dina-download ng ccam ang `ccam_<os>_<arch>.tar.gz` (sa Windows ang `ccam_windows_<arch>.zip`) at `checksums.txt`, na may 2 minuto para sa dalawa (`Updater.Fetch`). Hanggang 256 MB lang ang bawat download.
6. Dapat tumugma ang SHA-256 ng archive sa linya nito sa `checksums.txt` bago mangyari ang anumang iba. Pinapahinto ng nawawalang linya o ng naiibang checksum ang update na may `nothing was changed`.
7. Binabasa ng ccam ang tar at kinukuha ang unang regular na file na may pangalang `ccam` (`extract()`; kinukuha ng `extractZip` ang `ccam.exe` mula sa zip). Hindi pinapansin ang mga directory, link at ibang entry, at hanggang 256 MB lang din ang binary.
8. Isinusulat ng ccam ang bagong binary gamit ang `fsutil.AtomicWrite` (`Updater.Install`): isang pansamantalang file na `.ccam.tmp*` katabi ng binary, na may permission bit ng lumang file, na fina-flush sa disk at nire-rename sa ibabaw ng luma.

Ang buong pipeline, mula sa paghahanap ng tag hanggang sa rename:

<!-- motion:update-pipeline -->

Sabay na dina-download ang archive at `checksums.txt`, at walang isinusulat hangga't hindi tumutugma ang checksum. Atomic ang rename, kaya ang crash o punong disk ay mag-iiwan ng lumang binary o ng bago, hindi kailanman kalahati ng isa. Mga symlink sa ccam binary ang mga shortcut command gaya ng `~/.local/bin/ccwork`, kaya pinapatakbo nila ang bagong bersyon nang walang anumang pagbabago.

Sa Windows, hindi mapapalitan o mabubura ang file ng tumatakbong program, mapapalitan lang ang pangalan nito. Isinusulat ng ccam ang bagong binary sa `ccam.exe.new`, pinapalitan ang pangalan ng tumatakbong `ccam.exe` ng `ccam.exe.old`, at pinapalitan ang pangalan ng bagong file ng `ccam.exe` (`Updater.Install`). Kapag pumalya ang huling hakbang, ibinabalik ang pangalan ng lumang binary. Binubura ng susunod na update ang mga `ccam.exe.old` file na kaya nitong burahin; habang may session pang tumatakbo mula sa mas naunang `ccam.exe.old`, nagiging `ccam.exe.old-<time>` na lang ang tumatakbong binary. Nire-renew agad pagkatapos ang mga shortcut na hard link sa halip na symbolic link ([Mga shortcut command](/features/shortcuts)), kaya pinapatakbo rin nila ang bagong bersyon; saklaw nito ang mga hard link ng anumang mas lumang bersyon ng ccam na ginawa ng ccam, hindi lang ng binary na kapapalit lang.

### Mga installation na pinamamahalaan ng package

Pag-aari ng apt ang kopya ng ccam na in-install mula sa `.deb`, at kapag pinalitan ito nang hindi alam ng apt, hindi na magtutugma ang package database at ang file. Bago mag-download, sinusuri ng ccam kung sino ang may-ari ng binary:

<!-- motion:update-managed -->

1. Itinuturing na pinamamahalaan ng package ang binary sa `/usr/bin`, `/bin`, `/usr/sbin` o `/sbin`.
2. Sa ibang lugar, kung naka-install ang `dpkg` at nagtagumpay ang `dpkg -S <path>`, itinuturing din itong pinamamahalaan ng package.
3. Kung hindi, gumagawa at nag-aalis ang ccam ng pansamantalang file na `.ccam-update-*` katabi ng binary para matiyak na kaya nitong magsulat doon.

Para sa binary na pinamamahalaan ng package, humihinto ang ccam at ipi-print ang mga hakbang para sa apt:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://ccam.dev/install.sh | sh
```

Sa mga system na hindi Linux, ang linya ng install script lang ang ipi-print: `irm https://ccam.dev/install.ps1 | iex` sa Windows. Kapag hindi masulatan ang folder, imumungkahi ng ccam na patakbuhin ulit ang install script, na nag-i-install ng ccam sa `~/.local/bin`.

## Kaligtasan at mga espesyal na kaso

- Iniuulat ng development build ang bersyon nito bilang `dev` at walang release na maikukumpara, kaya tinatanggihan ito ng `ccam update` at `ccam update --check`. Nag-i-install ang `ccam update --version vX.Y.Z` ng release sa ibabaw nito. Para sa mga build na ginawa gamit ang `go install`, tingnan ang [I-install ang ccam](/docs/install#other-ways-to-install).
- Iniiwan nang gaya ng dati ang naka-install na binary kapag pumalya ang download, nawawala ang asset para sa iyong platform, o hindi tumugma ang checksum.
- Pumapalya ang release na hindi umiiral na may `there is no release <tag>`. Kapag tumanggi ang GitHub na may HTTP 403 o 429, sasabihin ng ccam na malamang rate limit ito at hihilingin na subukan mo ulit mamaya.
- Nagbabasa lang ang `--check`. Hindi nito sinusuri kung pinamamahalaan ng package ang binary o kung masusulatan ito.
- Ang `--version` na may mas lumang tag ang paraan para mag-roll back.

## Kaugnay

- [I-install ang ccam](/docs/install)
- [Sanggunian ng mga command](/docs/commands)
- [Mga madalas itanong](/docs/faq)
