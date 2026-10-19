---
title: Pag-update ng ccam
description: Pinapalitan ng ccam update ang ccam ng pinakabagong GitHub release pagkatapos ng SHA-256 check; nag-uulat lang ang --check at nag-i-install ang --version ng partikular na release tag.
updated: 2026-10-19
order: 10
translation: machine
sourceHash: 83fc20191328
---

Nagda-download ang `ccam update` ng ccam release mula sa GitHub, sinusuri ito laban sa `checksums.txt` ng release at pinapalitan ang tumatakbong binary sa isang rename. Gamitin ito kapag in-install mo ang ccam gamit ang install script o mula sa isang archive. Hindi kailanman kusang tumitingin ang ccam ng mga update: ang command na ito lang ang kumokontak sa GitHub, kaya gumagana nang pareho offline ang bawat ibang command.

## Paggamit

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

Kapag walang mas bagong release, ipi-print ng ccam ang `ccam 0.4.0 is up to date.` at walang babaguhin. Kapag may `--version` at ang bersyong tumatakbo na sa iyo, ipi-print nito ang `ccam 0.4.0 is already installed.`

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--check` | Iulat lang kung may mas bagong release | off |
| `--version <tag>` | I-install ang release tag na ito, gaya ng `v0.3.0`, mas bago man o mas luma kaysa sa kasalukuyan | pinakabago |
| `--dry-run` | Sinusunod: hinahanap ang pinakabagong tag (maliban kung may pinangalanan ang `--version`), pinapatakbo ang parehong mga pagsusuri sa binary, at ipi-print ang mga path ng download, checksum at papalitan nang hindi nagda-download | off |

Hindi puwedeng pagsabayin ang `--check` at `--version`. Dapat mukhang `vMAJOR.MINOR.PATCH` ang tag, na may opsyonal na pre-release o build suffix; puwedeng wala ang `v` sa unahan. Pumapalya ang anumang iba, halimbawa `--version takes a release tag such as v0.3.0, not "1.2"`.

Pinapalitan ng environment variable na `CCAM_RELEASES_URL` ang `https://github.com/jhyunwoo/ccam/releases`, para sa mirror o para sa mga test. Dapat pareho ang ayos nito: nagre-redirect ang `<url>/latest` sa `<url>/tag/<tag>`, at nasa `<url>/download/<tag>/<asset>` ang mga asset.

## Paano ito gumagana

Nasa `internal/update` ang code, na ang `ccam update` lang ang nag-i-import. Hindi kailanman tinatawag ito ng launch path.

1. Hinahanap ng `cli.updateCmd` ang tumatakbong binary gamit ang `os.Executable` at nire-resolve ang mga symlink gamit ang `filepath.EvalSymlinks`. Ang totoong path na iyon ang pinapalitan.
2. Kapag walang `--version`, hinihiling ng `Updater.Latest` ang `<releases>/latest` nang hindi sinusundan ang redirect, at binabasa ang tag mula sa header na `Location` (`…/tag/v0.5.0`). Ginagamit nito ang releases page sa halip na ang GitHub API, kaya hindi nalalapat ang limit ng API na 60 request kada oras. May 10 segundo ang paghahanap.
3. Sinusuri ng `parseVersion` ang tag laban sa isang mahigpit na pattern. Napupunta ang tag sa mga download URL at sa isang ipini-print na shell command, kaya walang nakakalusot na may `/`, `%`, espasyo o shell character. Ikinukumpara ito ng `compareVersions` sa tumatakbong bersyon; nauuna ang pre-release bago ang release nito.
4. Pinagpapasyahan ng `managed()` at `writable()` kung puwede bang palitan ng ccam ang binary (tingnan sa ibaba).
5. Sabay na dina-download ng `Updater.Fetch` ang `ccam_<os>_<arch>.tar.gz` at `checksums.txt`, na may 2 minuto para sa dalawa. Hanggang 256 MB lang ang bawat download.
6. Dapat tumugma ang SHA-256 ng archive sa linya nito sa `checksums.txt` bago mangyari ang anumang iba. Pinapahinto ng nawawalang linya o ng naiibang checksum ang update na may `nothing was changed`.
7. Binabasa ng `extract()` ang tar at kinukuha ang unang regular na file na may pangalang `ccam`. Hindi pinapansin ang mga directory, link at ibang entry, at hanggang 256 MB lang din ang binary.
8. Isinusulat ng `Updater.Install` ang bagong binary gamit ang `fsutil.AtomicWrite`: isang pansamantalang file na `.ccam.tmp*` katabi ng binary, na may permission bit ng lumang file, na fina-flush sa disk at nire-rename sa ibabaw ng luma.

Ang buong pipeline, mula sa paghahanap ng tag hanggang sa rename:

<!-- motion:update-pipeline -->

Magkatabing dina-download ang archive at `checksums.txt`, at walang isinusulat hangga't hindi tumutugma ang checksum. Atomic ang rename, kaya ang crash o punong disk ay mag-iiwan ng lumang binary o ng bago, hindi kailanman kalahati ng isa. Mga symlink sa ccam binary ang mga shortcut command gaya ng `~/.local/bin/ccwork`, kaya pinapatakbo nila ang bagong bersyon nang walang anumang pagbabago.

### Mga installation na pinamamahalaan ng package

Pag-aari ng apt ang ccam na in-install mula sa `.deb`, at kapag pinalitan ito nang hindi alam ng apt, magiging mali ang package database. Bago mag-download, sinusuri ng ccam kung sino ang may-ari ng binary:

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
  curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

Sa mga system na hindi Linux, ang linya ng install script lang ang ipi-print. Kapag hindi masulatan ang folder, imumungkahi ng ccam na patakbuhin ulit ang install script, na nag-i-install ng ccam sa `~/.local/bin`.

## Kaligtasan at mga espesyal na kaso

- Iniuulat ng development build ang bersyon nito bilang `dev` at walang release na maikukumpara, kaya tinatanggihan ito ng `ccam update` at `ccam update --check`. Nag-i-install ang `ccam update --version vX.Y.Z` ng release sa ibabaw nito.
- Iniiwan nang gaya ng dati ang naka-install na binary kapag pumalya ang download, nawawala ang asset para sa iyong platform, o hindi tumugma ang checksum.
- Pumapalya ang release na hindi umiiral na may `there is no release <tag>`. Kapag tumanggi ang GitHub na may HTTP 403 o 429, sasabihin ng ccam na malamang rate limit ito at hihilingin na subukan mo ulit mamaya.
- Nagbabasa lang ang `--check`. Hindi nito sinusuri kung pinamamahalaan ng package ang binary o kung masusulatan ito.
- Ang `--version` na may mas lumang tag ang paraan para mag-roll back.

## Kaugnay

- [I-install](/docs/install)
- [Mga command](/docs/commands)
- [FAQ](/docs/faq)
