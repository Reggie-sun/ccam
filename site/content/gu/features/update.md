---
title: ccam અપડેટ કરવું
description: ccam update SHA-256 તપાસ પછી ccam ને GitHub ની નવીનતમ રિલીઝથી બદલે છે; --check ફક્ત જણાવે છે અને --version આપેલો રિલીઝ tag ઇન્સ્ટૉલ કરે છે.
updated: 2026-10-26
order: 10
translation: machine
sourceHash: 598931e4801f
---

`ccam update` GitHub પરથી ccam ની રિલીઝ ડાઉનલોડ કરે છે, તેને રિલીઝની `checksums.txt` સામે તપાસે છે અને ચાલતી binary ને એક rename માં બદલે છે. જો તમે ccam ને install સ્ક્રિપ્ટ વડે કે archive માંથી ઇન્સ્ટૉલ કર્યું હોય તો તેને વાપરો. [ડેશબોર્ડ](/features/dashboard) પણ દરેક વખતે ખૂલે ત્યારે GitHub ને નવીનતમ રિલીઝ વિશે પૂછે છે અને `U` વડે તે ઑફર કરે છે, જે પૂછ્યા પછી આ જ પગલાં ચલાવે છે. બીજા કમાન્ડ ક્યારેય GitHub નો સંપર્ક કરતા નથી, તેથી તે offline પણ એ જ રીતે કામ કરે છે.

## ઉપયોગ

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

જ્યારે કોઈ નવી રિલીઝ ન હોય, ત્યારે ccam `ccam 0.4.0 is up to date.` છાપે છે અને કંઈ બદલતું નથી. જો `--version` તમારી પાસે પહેલેથી હોય તે વર્ઝન આપે, તો તે `ccam 0.4.0 is already installed.` છાપે છે.

## ઓપ્શન

| ઓપ્શન | શું કરે છે | ડિફૉલ્ટ |
| --- | --- | --- |
| `--check` | ફક્ત જણાવે છે કે નવી રિલીઝ છે કે નહીં | off |
| `--version <tag>` | આ રિલીઝ tag ઇન્સ્ટૉલ કરે છે, જેમ કે `v0.3.0`, હાલના કરતાં નવો કે જૂનો | નવીનતમ |
| `--dry-run` | સમર્થિત: નવીનતમ tag શોધે છે (સિવાય કે `--version` કોઈ tag આપે), binary પર એ જ તપાસો ચલાવે છે, અને ડાઉનલોડ કર્યા વગર download, checksum અને replace ના પાથ છાપે છે | off |

`--check` અને `--version` સાથે વાપરી શકાતા નથી. tag `vMAJOR.MINOR.PATCH` જેવો હોવો જોઈએ, વૈકલ્પિક pre-release કે build પ્રત્યય સાથે; શરૂઆતનો `v` છોડી શકાય છે. બીજું કંઈ પણ નિષ્ફળ જાય છે, ઉદાહરણ તરીકે `--version takes a release tag such as v0.3.0, not "1.2"`.

environment variable `CCAM_RELEASES_URL` mirror માટે કે tests માટે `https://github.com/jhyunwoo/ccam/releases` ને બદલે છે. તેની રચના એ જ હોવી જોઈએ: `<url>/latest` `<url>/tag/<tag>` પર redirect કરે, અને assets `<url>/download/<tag>/<asset>` પર હોય.

## તે કેવી રીતે કામ કરે છે

ccam ઇન્સ્ટૉલ થયેલી binary ને અડતાં પહેલાં download તપાસે છે, અને binary ને એક rename માં બદલે છે. કોડ `internal/update` માં છે, જે ફક્ત `ccam update` import કરે છે. શરૂઆતનો પાથ તેને ક્યારેય બોલાવતો નથી.

1. ccam `os.Executable` વડે ચાલતી binary શોધે છે અને `filepath.EvalSymlinks` વડે symlink ઉકેલે છે (`cli.updateCmd`). એ સાચો પાથ જ બદલાય છે.
2. `--version` વગર, ccam redirect અનુસર્યા વગર `<releases>/latest` ને request કરે છે (`Updater.Latest`), અને `Location` header (`…/tag/v0.5.0`) માંથી tag વાંચે છે. આ GitHub API ને બદલે releases પેજ વાપરે છે, તેથી API ની કલાકના 60 request ની મર્યાદા લાગુ પડતી નથી. lookup 10 સેકન્ડ પછી timeout થાય છે.
3. ccam tag ને કડક pattern સામે તપાસે છે (`parseVersion`). tag download URL માં અને છપાયેલા shell કમાન્ડમાં જાય છે, તેથી `/`, `%`, space કે shell અક્ષરોવાળું કંઈ પણ પસાર થતું નથી. પછી ccam તેની ચાલતા વર્ઝન સાથે સરખામણી કરે છે (`compareVersions`); pre-release તેની રિલીઝ પહેલાં આવે છે.
4. ccam નક્કી કરે છે કે તે binary બદલી શકે કે નહીં (`managed()`, `writable()`; નીચે જુઓ).
5. ccam `ccam_<os>_<arch>.tar.gz` (Windows પર `ccam_windows_<arch>.zip`) અને `checksums.txt` એકસાથે ડાઉનલોડ કરે છે, બંને માટે 2 મિનિટ સાથે (`Updater.Fetch`). દરેક download 256 MB સુધી મર્યાદિત છે.
6. બીજું કંઈ પણ થાય તે પહેલાં archive નો SHA-256 `checksums.txt` માંની તેની લાઇન સાથે મેળ ખાવો જોઈએ. લાઇન ન હોય કે checksum જુદો હોય તો અપડેટ `nothing was changed` સાથે અટકે છે.
7. ccam tar વાંચે છે અને `ccam` નામની પહેલી સામાન્ય ફાઇલ લે છે (`extract()`; `extractZip` zip માંથી `ccam.exe` લે છે). ફોલ્ડર, લિંક અને બીજી એન્ટ્રીઓ અવગણાય છે, અને binary પણ 256 MB સુધી મર્યાદિત છે.
8. ccam નવી binary `fsutil.AtomicWrite` વડે લખે છે (`Updater.Install`): binary ની બાજુમાં કામચલાઉ ફાઇલ `.ccam.tmp*`, જૂની ફાઇલના permission bits સાથે, જે ડિસ્ક પર flush થાય છે અને જૂની ફાઇલ પર rename થાય છે.

tag lookup થી rename સુધીની આખી pipeline:

<!-- motion:update-pipeline -->

archive અને `checksums.txt` સમાંતર રીતે ડાઉનલોડ થાય છે, અને checksum મેળ ન ખાય ત્યાં સુધી કંઈ લખાતું નથી. rename atomic છે, તેથી crash કે ભરાયેલી ડિસ્ક જૂની કે નવી binary છોડે છે, ક્યારેય અડધી નહીં. `~/.local/bin/ccwork` જેવા શૉર્ટકટ કમાન્ડ ccam binary ની symlink છે, તેથી તેઓ કોઈ ફેરફાર વગર નવું વર્ઝન ચલાવે છે.

Windows પર ચાલતા પ્રોગ્રામની ફાઇલ બદલી કે કાઢી શકાતી નથી, ફક્ત તેનું નામ બદલી શકાય છે. ccam નવી binary `ccam.exe.new` માં લખે છે, ચાલતી `ccam.exe` નું નામ `ccam.exe.old` કરે છે, અને નવી ફાઇલનું નામ `ccam.exe` કરે છે (`Updater.Install`). છેલ્લું પગલું નિષ્ફળ જાય તો જૂની binary નું નામ પાછું બદલાય છે. આગલું અપડેટ જે `ccam.exe.old` ફાઇલો કાઢી શકે તે કાઢી નાખે છે; અગાઉની `ccam.exe.old` માંથી કોઈ સેશન હજી ચાલતું હોય ત્યારે, ચાલતી binary તેના બદલે `ccam.exe.old-<time>` બને છે. symbolic links ને બદલે hard links હોય તેવા શૉર્ટકટ ([શૉર્ટકટ કમાન્ડ](/features/shortcuts)) તરત પછી નવા કરાય છે, જેથી તેઓ પણ નવું વર્ઝન ચલાવે; તેમાં ccam એ બનાવેલી, ccam ના કોઈ પણ જૂના વર્ઝનની hard links આવી જાય છે, ફક્ત હમણાં બદલાયેલી binary ની જ નહીં.

### package થી ઇન્સ્ટૉલ થયેલું

`.deb` માંથી ઇન્સ્ટૉલ થયેલી ccam ની કૉપિ apt ની છે, અને apt ની જાણ બહાર તેને બદલવાથી package database ફાઇલ સાથે મેળ ગુમાવી દે. ડાઉનલોડ કરતાં પહેલાં, ccam તપાસે છે કે binary કોની છે:

<!-- motion:update-managed -->

1. `/usr/bin`, `/bin`, `/usr/sbin` કે `/sbin` માંની binary package-managed ગણાય છે.
2. બીજે ક્યાંક, જો `dpkg` ઇન્સ્ટૉલ હોય અને `dpkg -S <path>` સફળ થાય, તો તે પણ package-managed ગણાય છે.
3. નહીં તો ccam ત્યાં લખી શકે છે તેની ખાતરી કરવા binary ની બાજુમાં કામચલાઉ `.ccam-update-*` ફાઇલ બનાવે છે અને દૂર કરે છે.

package-managed binary માટે ccam અટકે છે અને apt નાં પગલાં છાપે છે:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://ccam.dev/install.sh | sh
```

Linux સિવાયની સિસ્ટમ પર ફક્ત install સ્ક્રિપ્ટની લાઇન છપાય છે: Windows પર `irm https://ccam.dev/install.ps1 | iex`. જ્યારે ફોલ્ડરમાં લખી ન શકાય, ત્યારે ccam install સ્ક્રિપ્ટ ફરી ચલાવવાનું સૂચવે છે, જે ccam ને `~/.local/bin` માં ઇન્સ્ટૉલ કરે છે.

## સલામતી અને ખાસ કિસ્સા

- development build પોતાનું વર્ઝન `dev` બતાવે છે અને તેની પાસે સરખામણી માટે કોઈ રિલીઝ નથી, તેથી `ccam update` અને `ccam update --check` તેનો ઇનકાર કરે છે. `ccam update --version vX.Y.Z` તેના પર રિલીઝ ઇન્સ્ટૉલ કરે છે. `go install` વડે બનાવેલી binary માટે, [ccam ઇન્સ્ટૉલ કરો](/docs/install#other-ways-to-install) જુઓ.
- નિષ્ફળ download, તમારા platform માટે ખૂટતું asset, કે checksum નો મેળ ન ખાવો ઇન્સ્ટૉલ થયેલી binary ને જેમ હતી તેમ રહેવા દે છે.
- અસ્તિત્વમાં ન હોય તેવી રિલીઝ `there is no release <tag>` સાથે નિષ્ફળ જાય છે. જ્યારે GitHub HTTP 403 કે 429 સાથે ઇનકાર કરે, ત્યારે ccam કહે છે કે તે કદાચ rate limit છે અને પછીથી ફરી પ્રયાસ કરવા કહે છે.
- `--check` ફક્ત વાંચે છે. તે binary package-managed છે કે લખી શકાય તેવી છે તે તપાસતું નથી.
- જૂના tag સાથેનું `--version` એ પાછા જવાની રીત છે.

## સંબંધિત

- [ccam ઇન્સ્ટૉલ કરો](/docs/install)
- [કમાન્ડ સંદર્ભ](/docs/commands)
- [વારંવાર પૂછાતા પ્રશ્નો](/docs/faq)
