---
title: Mga account ayon sa folder
description: Sa ccam folder, iisang account ang gagamitin ng folder at laman nito kapag sinimulan ang Claude Code nang walang pangalan ng account, gaya ng work.
updated: 2026-10-26
order: 11
translation: machine
sourceHash: "25f82da64580"
---

Pagkatapos ng `ccam folder work ~/code/company`, sinisimulan ng anumang launch na walang pangalan ng account sa `~/code/company` o sa mga subfolder nito ang Claude Code bilang `work`. Sa labas nito, sinisimulan ang default na account gaya ng dati. Kapag pinangalanan mo ang isang account, gamit ang `ccam personal` o isang shortcut gaya ng `ccpersonal`, laging sinisimulan ang account na iyon.

## Paggamit

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

Ginagamit ng mga launch na ito ang account ng folder:

- `ccam -p "…"` at anumang ibang launch na ang unang argument ay isang option para sa Claude Code
- `ccam run` na walang pangalan, at `ccam --dry-run`
- `ccam --fallback -p "…"` para sa unang pagtatangka nito
- ang dashboard, ang `ccam` na mag-isa, na nagsisimulang nakatapat ang cursor sa account ng folder, kahit narating mo ang folder sa pamamagitan ng symbolic link

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `<account>` | Pangalan o numero ng account na gagamitin | |
| `[folder]` | Ang folder, relative o absolute | ang kasalukuyang folder |
| `--remove` | Alisin ang entry para sa eksaktong folder na iyon | off |
| `--dry-run` | Sinusuportahan: ipi-print ang pagbabago nang hindi ito sine-save | off |

## Paano ito gumagana

Isang `[[folder]]` entry ang bawat folder sa `~/.config/ccam/config.toml`, ang file na binabasa na rin ng bawat launch:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. Ginagawang absolute path ng `ccam folder` ang folder, na naresolba ang mga symbolic link, at sine-save ito. Kapag itinakda ulit ang isang folder, pinapalitan ang account nito.
2. Tinatawag ng launch na walang pangalan ang `launcher.HereProfile`. Kapag walang entry, kinukuha nito agad ang default na account.
3. Kung hindi, ikinukumpara nito ang working directory sa bawat `path` bilang mga string: ang folder mismo o isang nasa ilalim nito, hindi kailanman folder na nagsisimula lang sa parehong mga titik (hindi saklaw ng `/code/company` ang `/code/companyx`). Ang pinakamalalim na entry ang nananalo, kaya puwedeng `personal` ang `~/code` habang `work` ang `~/code/company`.
4. Kapag walang tumugma sa working directory gaya ng pagkakabigay nito, saka lang nireresolba ng ccam ang mga symbolic link nito at nagkukumpara ulit.
5. Sa Windows at macOS, hindi pinapansin ng paghahambing ang malaki at maliit na titik, gaya ng karaniwang ginagawa ng kanilang mga file system.

Walang ibang file na binabasa bukod sa `config.toml`, kaya walang nasusukat na dagdag na oras ang launch sa loob ng isang folder (4.0 ms may entry man o wala, `hyperfine`, 300 run).

## Kaligtasan at mga espesyal na kaso

- Dinadala ng `ccam rename` ang mga folder entry sa bagong pangalan, at binubura ng `ccam remove` ang mga entry ng inalis na account.
- Pinapahinto ng entry na wala nang umiiral na account ang launch na walang pangalan sa folder na iyon: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. Iniuulat ng `ccam doctor` ang mga ganitong entry, at ang mga entry na wala na ang folder.
- Kung patatakbuhin mo ang `ccam folder --remove` sa isang subfolder na nagmamana ng account nito mula sa isang parent folder, pinapangalanan ng ccam ang parent folder at walang inaalis.
- Pumipili ang `ccam next` ayon sa rotation o usage at hindi tumitingin sa mga folder.

## Kaugnay

- [Pagsisimula ng Claude Code](/features/launch)
- [Mga shortcut command](/features/shortcuts)
- [Pamamahala ng mga account](/features/accounts)
- [Sanggunian ng mga command](/docs/commands)
