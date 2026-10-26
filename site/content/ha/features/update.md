---
title: Sabunta ccam
description: ccam update yana maye gurbin ccam da sabon sakin GitHub bayan duba SHA-256; --check yana bayar da rahoto kawai, --version yana girka takamaiman tag na saki.
updated: 2026-10-26
order: 10
translation: machine
sourceHash: 598931e4801f
---

`ccam update` yana sauke sakin ccam daga GitHub, yana duba shi da `checksums.txt` na sakin kuma yana maye gurbin binary da ke gudana da rename ɗaya. Yi amfani da shi idan ka girka ccam da rubutun girkawa (install script) ko daga archive. [Dashboard](/features/dashboard) ma yana tambayar GitHub sabon saki duk lokacin da ya buɗe kuma yana ba da shi da `U`, wanda ke bin matakai iri ɗaya bayan ya tambaya. Sauran umarnai ba sa taɓa tuntuɓar GitHub, don haka suna aiki iri ɗaya ba tare da intanet ba.

## Amfani

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

Idan babu sabon saki, ccam yana buga `ccam 0.4.0 is up to date.` kuma ba ya canza komai. Idan `--version` ya ambaci sigar da kake da ita tuni, yana buga `ccam 0.4.0 is already installed.`

## Zaɓuɓɓuka

| Zaɓi | Abin da yake yi | Tsoho |
| --- | --- | --- |
| `--check` | Bayar da rahoto kawai ko akwai sabon saki | a kashe |
| `--version <tag>` | Girka wannan tag na saki, kamar `v0.3.0`, sabo ko tsoho fiye da na yanzu | na baya-bayan nan |
| `--dry-run` | Ana tallafa masa: yana neman tag na baya-bayan nan (sai dai idan `--version` ya ambaci ɗaya), yana gudanar da dubawa iri ɗaya a kan binary, kuma yana buga hanyoyin saukewa, checksum da maye gurbi ba tare da saukewa ba | a kashe |

Ba za a iya amfani da `--check` da `--version` tare ba. Dole tag ɗin ya yi kama da `vMAJOR.MINOR.PATCH`, tare da ƙarin pre-release ko build idan ana so; ana iya barin `v` na farko. Duk wani abu dabam yana gazawa, misali `--version takes a release tag such as v0.3.0, not "1.2"`.

Mai canjin muhalli `CCAM_RELEASES_URL` yana maye gurbin `https://github.com/jhyunwoo/ccam/releases`, don madubi (mirror) ko don gwaje-gwaje. Dole ya kasance da tsari iri ɗaya: `<url>/latest` yana turawa zuwa `<url>/tag/<tag>`, kuma assets suna a `<url>/download/<tag>/<asset>`.

## Yadda yake aiki

ccam yana duba saukewar kafin ya taɓa binary da aka girka, kuma yana maye gurbin binary da rename ɗaya. Lambar tana cikin `internal/update`, wanda `ccam update` ne kaɗai ke shigo da shi. Hanyar farawa ba ta taɓa kiran sa.

1. ccam yana nemo binary da ke gudana da `os.Executable` kuma yana warware symlinks da `filepath.EvalSymlinks` (`cli.updateCmd`). Wannan ainihin hanyar ita ce ake maye gurbinta.
2. Ba tare da `--version` ba, ccam yana neman `<releases>/latest` ba tare da bin turawar ba (`Updater.Latest`), kuma yana karanta tag ɗin daga header na `Location` (`…/tag/v0.5.0`). Wannan yana amfani da shafin saki maimakon GitHub API, don haka iyakar API ta buƙatu 60 a awa ba ta shafa ba. Neman yana ƙarewa bayan daƙiƙa 10.
3. ccam yana duba tag ɗin da tsari mai tsauri (`parseVersion`). Tag yana ƙarewa a URLs na saukewa da a umarnin shell da aka buga, don haka babu abin da ke da `/`, `%`, sarari ko haruffan shell da ke wucewa. Sannan ccam yana kwatanta shi da sigar da ke gudana (`compareVersions`); pre-release yana zuwa kafin sakinsa.
4. ccam yana yanke shawara ko zai iya maye gurbin binary kwata-kwata (`managed()`, `writable()`; duba ƙasa).
5. ccam yana sauke `ccam_<os>_<arch>.tar.gz` (a Windows `ccam_windows_<arch>.zip`) da `checksums.txt` a lokaci ɗaya, da mintuna 2 ga duka biyun (`Updater.Fetch`). Kowane saukewa an iyakance shi zuwa 256 MB.
6. Dole SHA-256 na archive ɗin ya yi daidai da layinsa a `checksums.txt` kafin wani abu ya faru. Layin da ya ɓace ko checksum daban yana dakatar da sabuntawar da `nothing was changed`.
7. ccam yana karanta tar ɗin kuma yana ɗaukar fayil na yau da kullum na farko mai suna `ccam` (`extract()`; `extractZip` yana ɗaukar `ccam.exe` daga zip ɗin). Ana watsi da foldoji, mahaɗai da sauran shigarwa, kuma binary ɗin ma an iyakance shi zuwa 256 MB.
8. ccam yana rubuta sabon binary da `fsutil.AtomicWrite` (`Updater.Install`): fayil na wucin gadi `.ccam.tmp*` kusa da binary, da bits na izinin tsohon fayil, ana tura shi zuwa faifai kuma ana rename shi a kan tsohon.

Dukkan hanyar, daga neman tag zuwa rename:

<!-- motion:update-pipeline -->

Ana sauke archive da `checksums.txt` a lokaci guda, kuma ba a rubuta komai har sai checksum ya yi daidai. Rename ɗin atomic ne, don haka faɗuwar shiri ko cikakken faifai yana barin ko dai tsohon binary ko sabon, ba taɓa rabin ɗaya ba. Gajerun umarni kamar `~/.local/bin/ccwork` symlinks ne zuwa binary na ccam, don haka suna gudanar da sabuwar sigar ba tare da wani canji ba.

A Windows ba za a iya maye gurbin ko goge fayil na shirin da ke gudana ba, sai dai a sauya sunansa. ccam yana rubuta sabon binary zuwa `ccam.exe.new`, yana sauya sunan `ccam.exe` da ke gudana zuwa `ccam.exe.old`, kuma yana sauya sunan sabon fayil zuwa `ccam.exe` (`Updater.Install`). Idan mataki na ƙarshe ya gaza, ana mayar da sunan tsohon binary. Sabuntawa ta gaba tana goge fayilolin `ccam.exe.old` da za ta iya; yayin da wani zama ke gudana har yanzu daga `ccam.exe.old` na baya, binary da ke gudana yana zama `ccam.exe.old-<time>` maimakon haka. Gajerun umarni da ke hard links maimakon symbolic links ([Gajerun umarni](/features/shortcuts)) ana sabunta su nan take bayan haka, don haka su ma suna gudanar da sabuwar sigar; hakan ya haɗa da hard links na kowace tsohuwar sigar ccam da ccam ya ƙirƙira, ba na binary da aka maye gurbinsa yanzu kawai ba.

### Girkin da mai sarrafa kunshi ke sarrafawa

Kwafin ccam da aka girka daga `.deb` na apt ne, kuma maye gurbinsa a bayan apt zai bar bayanan kunshin ba su yi daidai da fayil ɗin ba. Kafin saukewa, ccam yana duba wanene ke da binary ɗin:

<!-- motion:update-managed -->

1. Binary a `/usr/bin`, `/bin`, `/usr/sbin` ko `/sbin` ana ɗaukar sa a matsayin wanda kunshi ke sarrafawa.
2. A wani wuri, idan an girka `dpkg` kuma `dpkg -S <path>` ya yi nasara, shi ma ana ɗaukar sa a matsayin wanda kunshi ke sarrafawa.
3. In ba haka ba ccam yana ƙirƙira kuma yana cire fayil na wucin gadi `.ccam-update-*` kusa da binary don tabbatar da cewa zai iya rubutu a can.

Ga binary da kunshi ke sarrafawa ccam yana tsayawa kuma yana buga matakan apt:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://ccam.dev/install.sh | sh
```

A tsarin da ba Linux ba, layin rubutun girkawa ne kawai ake bugawa: `irm https://ccam.dev/install.ps1 | iex` a Windows. Idan ba za a iya rubutu a foldar ba, ccam yana ba da shawarar sake gudanar da rubutun girkawa, wanda ke girka ccam a `~/.local/bin`.

## Tsaro da yanayi na musamman

- Build na haɓakawa (development) yana bayar da sigarsa a matsayin `dev` kuma ba shi da sakin da za a kwatanta da shi, don haka `ccam update` da `ccam update --check` suna ƙin sa. `ccam update --version vX.Y.Z` yana girka saki a kansa. Don builds da aka yi da `go install`, duba [Girka ccam](/docs/install#other-ways-to-install).
- Saukewa da ya gaza, asset da ya ɓace ga dandalinka, ko checksum da bai yi daidai ba yana barin binary da aka girka yadda yake.
- Sakin da babu shi yana gazawa da `there is no release <tag>`. Idan GitHub ya ƙi da HTTP 403 ko 429, ccam yana cewa wataƙila iyakar buƙatu ce kuma yana roƙon ka ka sake gwadawa daga baya.
- `--check` yana karantawa ne kawai. Ba ya duba ko kunshi ke sarrafa binary ɗin ko ana iya rubutu a kansa.
- `--version` tare da tsohon tag shi ne yadda ake komawa baya (roll back).

## Masu alaƙa

- [Girka ccam](/docs/install)
- [Bayanin umarnai](/docs/commands)
- [Tambayoyin da ake yawan yi](/docs/faq)
