---
title: Kusasisha ccam
description: ccam update hubadilisha ccam kwa toleo jipya zaidi la GitHub baada ya ukaguzi wa SHA-256; --check huripoti tu na --version husakinisha tagi mahususi ya toleo.
updated: 2026-10-20
order: 10
translation: machine
sourceHash: d91499b95d3d
---

`ccam update` hupakua toleo la ccam kutoka GitHub, hulikagua dhidi ya `checksums.txt` ya toleo hilo na hubadilisha binary inayoendeshwa kwa kubadilisha jina mara moja. Itumie ikiwa ulisakinisha ccam kwa script ya kusakinisha au kutoka kwenye archive. ccam kamwe haikagui masasisho yenyewe: amri hii pekee ndiyo huwasiliana na GitHub, kwa hiyo kila amri nyingine hufanya kazi vivyo hivyo bila mtandao.

## Matumizi

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

Ikiwa hakuna toleo jipya zaidi, ccam huchapisha `ccam 0.4.0 is up to date.` na haibadilishi chochote. Kwa `--version` na toleo unaloendesha tayari, huchapisha `ccam 0.4.0 is already installed.`

## Chaguo

| Chaguo | Inachofanya | Chaguo-msingi |
| --- | --- | --- |
| `--check` | Huripoti tu kama toleo jipya zaidi lipo | imezimwa |
| `--version <tag>` | Husakinisha tagi hii ya toleo, kama `v0.3.0`, jipya au la zamani kuliko la sasa | jipya zaidi |
| `--dry-run` | Inaheshimiwa: hutafuta tagi ya karibuni (isipokuwa `--version` imetaja moja), hufanya ukaguzi ule ule kwenye binary, na huchapisha njia za kupakua, checksum na kubadilisha bila kupakua | imezimwa |

`--check` na `--version` haziwezi kutumiwa pamoja. Tagi lazima ionekane kama `vMAJOR.MINOR.PATCH`, ikiwa na kiambishi cha hiari cha pre-release au build; `v` ya mwanzo inaweza kuachwa. Kitu kingine chochote hushindwa, kwa mfano `--version takes a release tag such as v0.3.0, not "1.2"`.

Kigezo cha mazingira `CCAM_RELEASES_URL` huchukua nafasi ya `https://github.com/jhyunwoo/ccam/releases`, kwa ajili ya kioo (mirror) au kwa majaribio. Lazima kipangwe kwa njia ile ile: `<url>/latest` huelekeza kwenda `<url>/tag/<tag>`, na faili za toleo ziko katika `<url>/download/<tag>/<asset>`.

## Jinsi inavyofanya kazi

Msimbo uko katika `internal/update`, ambayo ni `ccam update` pekee inayoiingiza. Njia ya kuanzisha haiiiti kamwe.

1. `cli.updateCmd` hupata binary inayoendeshwa kwa `os.Executable` na hufuata symlinks kwa `filepath.EvalSymlinks`. Njia hiyo halisi ndiyo hubadilishwa.
2. Bila `--version`, `Updater.Latest` huomba `<releases>/latest` bila kufuata uelekezaji, na husoma tagi kutoka kichwa cha `Location` (`…/tag/v0.5.0`). Hii hutumia ukurasa wa matoleo badala ya GitHub API, kwa hiyo kikomo cha API cha maombi 60 kwa saa hakihusiki. Utafutaji una sekunde 10.
3. `parseVersion` hukagua tagi dhidi ya mchoro mkali. Tagi huishia kwenye URL za kupakua na kwenye amri ya shell inayochapishwa, kwa hiyo hakuna kitu chenye `/`, `%`, nafasi au herufi za shell kinachopita. `compareVersions` huilinganisha na toleo linaloendeshwa; pre-release hupangwa kabla ya toleo lake.
4. `managed()` na `writable()` huamua kama ccam inaruhusiwa kubadilisha binary hata kidogo (tazama hapa chini).
5. `Updater.Fetch` hupakua `ccam_<os>_<arch>.tar.gz` (kwenye Windows `ccam_windows_<arch>.zip`) na `checksums.txt` kwa wakati mmoja, kwa dakika 2 kwa zote mbili. Kila upakuaji una kikomo cha MB 256.
6. SHA-256 ya archive lazima ilingane na mstari wake katika `checksums.txt` kabla ya jambo lingine lolote kutokea. Mstari unaokosekana au checksum tofauti husimamisha sasisho kwa `nothing was changed`.
7. `extract()` husoma tar na huchukua faili ya kawaida ya kwanza yenye jina `ccam` (`extractZip` huchukua `ccam.exe` kutoka kwenye zip). Folda, viungo na vitu vingine hupuuzwa, na binary nayo ina kikomo cha MB 256.
8. `Updater.Install` huandika binary mpya kwa `fsutil.AtomicWrite`: faili ya muda `.ccam.tmp*` kando ya binary, yenye ruhusa za faili ya zamani, huandikwa kwenye diski na hubadilishwa jina juu ya ile ya zamani.

Mchakato wote, kuanzia kutafuta tagi hadi kubadilisha jina:

<!-- motion:update-pipeline -->

Archive na `checksums.txt` hupakuliwa sambamba, na hakuna kinachoandikwa hadi checksum ilingane. Kubadilisha jina ni atomiki, kwa hiyo hitilafu ya ghafla au diski iliyojaa huacha binary ya zamani au mpya, kamwe si nusu ya moja. Amri za njia fupi kama `~/.local/bin/ccwork` ni symlinks kuelekea binary ya ccam, kwa hiyo huendesha toleo jipya bila badiliko lolote.

Kwenye Windows faili la programu inayoendeshwa haliwezi kubadilishwa wala kufutwa, ila kupewa jina jipya tu. `Updater.Install` huandika binary mpya kwenye `ccam.exe.new`, hubadilisha jina la `ccam.exe` inayoendeshwa kuwa `ccam.exe.old`, na hubadilisha jina la faili jipya kuwa `ccam.exe`. Hatua ya mwisho ikishindwa, binary ya zamani hurudishiwa jina lake. Sasisho linalofuata hufuta `ccam.exe.old`. Njia fupi ambazo ni hard links badala ya viungo vya ishara ([Amri za njia fupi](/features/shortcuts)) hufanywa upya mara moja baadaye, kwa hiyo nazo huendesha toleo jipya.

### Usakinishaji unaosimamiwa na kifurushi

ccam iliyosakinishwa kutoka `.deb` ni mali ya apt, na kuibadilisha bila apt kujua kungeacha hifadhidata ya vifurushi ikiwa si sahihi. Kabla ya kupakua, ccam hukagua nani anamiliki binary:

<!-- motion:update-managed -->

1. Binary iliyo katika `/usr/bin`, `/bin`, `/usr/sbin` au `/sbin` huhesabiwa kuwa inasimamiwa na kifurushi.
2. Mahali pengine, ikiwa `dpkg` imesakinishwa na `dpkg -S <path>` inafaulu, nayo huhesabiwa kuwa inasimamiwa na kifurushi.
3. Vinginevyo ccam huunda na kuondoa faili ya muda `.ccam-update-*` kando ya binary ili kuhakikisha inaweza kuandika hapo.

Kwa binary inayosimamiwa na kifurushi ccam husimama na huchapisha hatua za apt:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

Kwenye mifumo isiyo Linux, mstari wa script ya kusakinisha pekee ndio huchapishwa: `irm https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.ps1 | iex` kwenye Windows. Folda isipoweza kuandikwa, ccam hupendekeza kuendesha tena script ya kusakinisha, ambayo husakinisha ccam katika `~/.local/bin`.

## Usalama na hali maalum

- Build ya maendeleo huripoti toleo lake kama `dev` na haina toleo la kulinganisha nalo, kwa hiyo `ccam update` na `ccam update --check` huikataa. `ccam update --version vX.Y.Z` husakinisha toleo juu yake.
- Upakuaji ulioshindwa, faili ya toleo inayokosekana kwa jukwaa lako, au checksum isiyolingana huacha binary iliyosakinishwa kama ilivyokuwa.
- Toleo lisilokuwepo hushindwa na `there is no release <tag>`. GitHub inapokataa kwa HTTP 403 au 429, ccam husema huenda ni kikomo cha maombi na hukuomba ujaribu tena baadaye.
- `--check` husoma tu. Haikagui kama binary inasimamiwa na kifurushi au inaweza kuandikwa.
- `--version` yenye tagi ya zamani ndiyo njia ya kurudi kwenye toleo la awali.

## Zinazohusiana

- [Kusakinisha](/docs/install)
- [Amri](/docs/commands)
- [Maswali yanayoulizwa mara kwa mara](/docs/faq)
