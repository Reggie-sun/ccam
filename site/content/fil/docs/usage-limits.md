---
title: Mga usage limit at paglipat ng account
description: Tingnan ang 5-oras at lingguhang usage ng bawat Claude Code account gamit ang ccam usage, at buksan ang may pinakamaraming natitira gamit ang ccam next --auto.
updated: 2026-10-19
order: 5
translation: machine
sourceHash: d288177682d7
---

Ipinapakita ng `ccam usage` kung gaano kalaki na ang nagamit ng bawat account sa 5-oras at lingguhang limit. Binubuksan ng `ccam next --auto` ang account na may pinakamaraming natitira, at dinadala ng `--continue` ang kasalukuyang usapan.

```console
$ ccam usage
 Account     5-hour        Resets    Weekly        Resets    As of
 1 personal  ██████░░  72% 1h 40m    ███░░░░░  41% 3d 6h     just now
 2 work      █░░░░░░░   8% 4h 12m    ██░░░░░░  19% 5d 2h     just now
 3 side      ████████ 100% 22m       █████░░░  68% 1d 9h     12m ago (cached)
→ Most headroom: work   (open it with ccam next --auto)
```

Ang mga column ay: account, 5-oras na usage, oras bago ito mag-reset, lingguhang usage, oras bago ito mag-reset, at kung gaano na katagal ang mga numero.

## Saan galing ang mga numero

Sabay-sabay na nagtatanong ang ccam sa usage endpoint ng Anthropic para sa bawat account, gamit ang sariling login ng bawat account. Naghihintay ito nang hanggang 3 segundo bawat account at ipinapakita ang bawat account sa sandaling sumagot ito.

Hindi pampublikong dokumentado ang endpoint na ito at puwede itong magbago. Kapag pumalya ang isang request, ipinapakita ng ccam ang huling alam na mga numero at minamarkahan ang mga ito bilang cached. Ang window na nag-reset na mula nang i-save ang mga numerong iyon ay nagpapakita ng 0%. Sasabihin sa iyo ng `ccam doctor` kung nagbago ang format ng response.

Muling ginagamit ang mga numerong kinuha wala pang 60 segundo ang nakalipas. Mga option:

| Option | Epekto |
| --- | --- |
| `--watch [seconds]` | Tuloy-tuloy na i-refresh, tuwing 60 segundo bilang default |
| `--json` | Output na nababasa ng makina |
| `-v` | Ipakita rin ang mga lingguhang limit ng bawat model at kung saan galing ang bawat numero |
| `--no-cache` | Laging magtanong ulit |
| `--source cache` | Huwag kailanman magtanong; ipakita lang ang mga naka-save na numero |

## Mga expired na login

Kapag nag-expire na ang access token ng isang account, hindi ito nire-renew ng ccam, dahil kapag ni-renew ito habang ginagamit ng isang Claude Code session ang parehong login, puwedeng ma-sign out ang session na iyon. Simulan nang isang beses ang account gamit ang `ccam <name>` at ire-renew ito ng Claude Code.

Nire-renew mismo ng `ccam usage --refresh` ang mga expired na token, pero para lang sa mga account na hindi tumatakbo. Hindi nito kailanman nire-renew ang mga login na naka-store sa macOS Keychain o sa isang in-import na `~/.claude`.

## Usage sa status line

```sh
ccam statusline install
```

Ipapakita na ng Claude Code ang pangalan ng account at ang usage nito sa status line, halimbawa `[work] · Opus · 5h 24% · 7d 41%`. Kung may status line command ka na, pinapanatili ito ng ccam at inilalagay sa unahan ang pangalan ng account. Sine-save din ng status line ang mga numero nito, na ginagamit ng `ccam usage` kapag hindi sumasagot ang endpoint.

## Paglipat

```sh
ccam next             # the account after the one you used last
ccam next --auto      # the account with the lowest 5-hour usage
ccam next --continue  # the next account that shares your conversation history
```

Nilalaktawan ng `ccam next` ang mga account na hindi naka-log in, at nilalaktawan ng `--auto` ang mga account na nasa 100%. Para maituloy ang isang usapan, dapat parehong nagshe-share ng kasaysayan ng usapan ang dalawang account; tingnan ang [sharing](/docs/sharing).

Ang [gabay sa workflow kapag naabot ang usage limit](/guides/usage-limit-workflow) ay nagpapakita ng isang buong halimbawa.
