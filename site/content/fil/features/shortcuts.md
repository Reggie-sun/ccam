---
title: Mga shortcut command
description: Bigyan ang isang account ng sariling command gaya ng ccwork gamit ang ccam shortcut, at gawin muli ang mga link sa ~/.local/bin gamit ang ccam shims install.
updated: 2026-10-26
order: 3
translation: machine
sourceHash: "066e7e622517"
---

Ang shortcut ay isang command na ikaw mismo ang nagpangalan at nagsisimula ng Claude Code bilang isang account: ginagawa ng `ccwork --resume` ang ginagawa ng `ccam work --resume`. Isa itong link sa ccam sa `~/.local/bin`, kaya gumagana ito sa bash, zsh, mga script at anumang ibang nagpapatakbo ng program mula sa `PATH`. Opsyonal ang mga shortcut; gumawa ng isa para sa mga account na madalas mong sinisimulan.

## Paggamit

```sh
ccam shortcut work ccwork     # create the shortcut ccwork for work
ccam shortcut 2 ccside        # the account with number 2
ccam shortcut work            # show it
ccam shortcut work ccw        # rename it; the old link is removed
ccam shortcut work --remove   # delete it
ccam add side --shortcut ccside
ccam shims install            # recreate missing or outdated links of all accounts
ccwork -p "summarise" --output-format json
```

Kinukumpirma ng ccam ang bawat pagbabago:

```console
$ ccam shortcut work ccwork
✓ shortcut for work: ccwork
$ ccam shortcut work ccw
✓ shortcut for work: ccw (was ccwork)
$ ccam shortcut work --remove
✓ ccw removed; work has no shortcut now
$ ccam shortcut work
work: none (add one with `ccam shortcut work ccwork`)
```

Hindi hihigit sa isa ang shortcut ng isang account, na naka-store bilang `shortcut` sa `~/.config/ccam/config.toml`. Nagtatanong ang wizard ng `ccam add` ng `Create a shortcut command for <name>? [y/N]` at pagkatapos ay `Shortcut name [cc<name>]:`. Sa labas ng wizard, gumagawa lang ng isa ang `ccam add` kapag may `--shortcut`. Ipinapakita ng `ccam list` ang shortcut ng bawat account, o `-`.

Ipinapasa ng shortcut ang bawat argument sa Claude Code, sa parehong paraan ng `ccam <name>` pagkatapos ng pangalan. Kasama rito ang `--dry-run` at `--help`: para ma-preview ang isang pagsisimula, gamitin ang `ccam --dry-run work`.

## Mga option

| Option | Ano ang ginagawa | Default |
| --- | --- | --- |
| `--remove` | Burahin ang shortcut ng account at ang link nito | off |
| `-h`, `--help` | Help para sa command | |
| `--dry-run` | Sinusuportahan ito ng `ccam shortcut`: sinusuri nito ang pangalan at ipi-print ang `Would set the shortcut of work to ccwork.` o `Would remove ccwork, the shortcut of work.` nang hindi nagse-save o nagli-link | off |

Walang tinatanggap na argument ang `ccam shims install` at tinatanggihan nito ang `--dry-run` na may `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``.

Mga patakaran sa pangalan, na sinusuri sa ganitong pagkakasunod-sunod:

- 1 hanggang 64 na character ng `a-z`, `0-9`, `.`, `-` at `_`, na nagsisimula sa titik o numero;
- hindi `ccam`, `claude` o ccam command gaya ng `run` o `usage`;
- hindi shortcut ng ibang account;
- hindi umiiral na file sa `~/.local/bin`, maliban kung link ito na ginawa ng ccam (gaya ng sariling shortcut ng account na ito);
- hindi pangalan ng ibang executable sa iyong `PATH`.

Sinusuri ang `config.toml` laban sa parehong mga patakaran kapag nilo-load ito. Kung mano-mano mong binago ang isang shortcut kaya lumalabag ito sa pattern, `claude` o `ccam` ito, o ginagamit ito ng dalawang account, papalya ang bawat command na may mensaheng nagpapangalan sa account hanggang ayusin mo ito.

Kapag sabay na ibinigay ang command at `--remove`, humihinto ito na may `give a command or --remove, not both`. Kapag itinakda ang pangalang mayroon na ang account, ipi-print ang `work: ccwork (unchanged)` at gagawin muli ang link nito kung nawawala.

## Paano ito gumagana

Ang pagtatakda ng pangalan ay nagpapatakbo ng `checkShortcut`, na tumatawag sa `shim.Check` kasama ang bin folder at ang totoong path ng ccam binary:

<!-- motion:shortcuts-check -->

Bumababa ang pangalan sa mga pagsusuri nang sunod-sunod. Katabi ng bawat pagsusuri ang isang pangalang tinatanggihan nito: pumapalya ang `Work` sa pattern, ccam command ang `run`, pag-aari ng ibang account ang `ccside`, umiiral na file sa `~/.local/bin` ang `ccfile`, at itatago ng `git` ang `/usr/bin/git`. Ang pangalan lang na pumasa sa bawat pagsusuri ang sine-save at nili-link.

1. Itinutugma ng ccam ang pangalan sa `^[a-z0-9][a-z0-9._-]{0,63}$`, saka sa `config.Reserved`, saka sa `cfg.FindShortcut` (`shim.Check`). Tinatawag nito ang `os.Lstat` sa `~/.local/bin/<name>`, at para sa bawat ibang folder sa `PATH` ay tinatawag nito ang `os.Stat` at tinatanggihan ang executable file na hindi link sa ccam.
2. Sine-save ng ccam ang `config.toml` gamit ang `fsutil.AtomicWrite` (pansamantalang file, `fsync`, rename).
3. Binubura ng ccam ang link ng lumang shortcut, pero kapag link lang ito sa ccam (`shim.Remove`).
4. Dinadaanan ng ccam ang bawat naka-configure na shortcut ayon sa pagkakasunod-sunod ng pangalan at ginagawa ang link ng bawat isang nawawala (`shim.Install`), gamit ang `os.Symlink` papunta sa ccam binary, na hinahanap gamit ang `os.Executable` at `filepath.EvalSymlinks`. Hindi napapahinto ang iba ng file na nakaharang sa isang shortcut; pinapangalanan ng error ang bawat isa.

Itinatala ng ccam ang bawat link na ginagawa nito, at kung aling binary ang pinapatakbo ng link na iyon, sa `~/.local/state/ccam/shortcuts.json`. Itinuturing na sariling link ng ccam ang isang link kapag ibinabalik ng `os.Readlink` ang path ng ccam, o kapag tumuturo sa parehong file ang link at ang path ng ccam. Luma na ang naitalang link na ang binary ay inilipat o pinalitan mula noon (ang ccam binary na inilipat sa ibang folder, muling pag-install), at itinuturo ito ng ccam sa ccam binary na ito (`shim.Install`). Hindi ginagalaw ang naitalang link sa ibang naka-install na kopya ng ccam. Dumadaan sa mga pagsusuring ito ang bawat pag-alis at bawat pagsusuri bago mag-overwrite, kaya hindi ginagalaw ng ccam ang mga file na hindi nito ginawa, gaya ng sarili mong kopya ng ccam o link na ikaw mismo ang gumawa.

Kapag nagpatakbo ka ng shortcut, ang pangalang tina-type mo lang ang kailangan ng ccam:

<!-- motion:shortcuts-link -->

Hinahanap ng shell ang `ccwork` sa `PATH`, sinusundan ng kernel ang link at sinisimulan ang ccam binary, at `ccwork` pa rin ang `argv[0]`. Sa `app.Main`, nilo-load nang isang beses ng anumang pangalan maliban sa `ccam` ang `config.toml` at tinatawag ang `shim.Lookup`, na hinahanap ang account na ang `shortcut` ay katumbas ng `filepath.Base(argv[0])`. Pagkatapos ay ginagawa ng ccam ang parehong mga hakbang ng [`ccam <name>`](/features/launch) na naka-load na ang config, at ipinapasa ang `argv[1:]` sa `claude` nang hindi ginagalaw (`launcher.LaunchWith`).

Bakit symlink at hindi shell alias o wrapper script: gumagana ang symlink sa bawat shell at sa mga script, naiiwasan nito ang pagsisimula ng karagdagang shell, na gagastos ng mga 1 hanggang 2 ms, at iisang file pa rin ang binabasa ng pagsisimula. Sinusukat ng `scripts/perf-check.sh` ang isang shortcut gamit ang `hyperfine` laban sa parehong 10 ms (p95) na limit ng `ccam <name>`.

### Sa Windows

Ang shortcut ay `ccwork.exe` sa `%USERPROFILE%\.local\bin`, at inaalis ng ccam ang `.exe` (anuman ang laki ng titik) bago ang paghahanap (`shim.Name`), kaya sinisimulan ng `ccwork`, `ccwork.exe` at `CCWORK.EXE` lahat ang `work`. Gumagawa muna ang ccam ng symbolic link, na pinapayagan ng Windows sa Developer Mode o para sa isang administrator. Kapag wala ang karapatang iyon, gumagawa ito ng hard link: pangalawang pangalan para sa parehong `ccam.exe` file, na hindi nangangailangan ng pribilehiyo. Itinuturing na sariling link ng ccam ang hard link kapag ito ang parehong file ng `ccam.exe`.

Laging pinapatakbo ng hard link ang file na pinagmulan nito, kaya kapag naglagay ang `ccam update` ng bagong file kapalit ng `ccam.exe`, nire-renew nito ang bawat shortcut na hard link gamit ang `shim.Relink` at ipi-print ang `✓ shortcut ccwork renewed`. Ganoon din ang ginagawa ng muling pagpapatakbo ng `install.ps1` sa pamamagitan ng `ccam shims install` (`OK 2 shortcuts now run the new ccam`), at gayundin ang `ccam shims install` at `ccam doctor --fix` para sa anumang shortcut na kopya ng mas lumang bersyon ng ccam. Hinahanap ng pagsusuri ng pangalan ang `<name>` kasama ang bawat extension sa `PATHEXT` (`.exe`, `.cmd` at iba pa) sa ibang mga folder ng `PATH`.

I-type ang shortcut sa isang shell (PowerShell, cmd.exe, Git Bash). Ang symbolic link na binuksan mula sa Explorer, sa Run dialog o sa `Start-Process` ay nagsisimula bilang karaniwang `ccam`, na nagbubukas ng dashboard sa halip na ng account.

## Kaligtasan at mga espesyal na kaso

- **Link na walang account na umaangkin.** Kung may link sa ccam sa `~/.local/bin` na may pangalang hindi ginagamit ng anumang account, halimbawa pagkatapos mong mano-manong baguhin ang `shortcut` sa `config.toml`, hindi babalik ang ccam sa default na account. Humihinto ito na may `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) ``. Itakda ulit ang pangalan gamit ang `ccam shortcut`, o burahin ang link.
- **Mga umiiral na file.** Hindi kailanman pinapalitan ng ccam ang file na hindi nito ginawa. Iniuulat ng `ccam shims install` ang `<path> already exists and is not a ccam link, so it was left alone` at ginagawa pa rin ang ibang mga link. Ipi-print naman ng `ccam add` at `ccam import` ang `shortcut: ccwork (not created; see above)`.
- **Inilipat na ccam binary.** Kapag inilipat mo ang ccam binary, tumuturo ang mga link na ginawa ng ccam sa path na wala na. Itinuturo ng `ccam shims install` (o ng `ccam doctor --fix`) ang mga ito sa ccam binary na ginamit mo sa pagpapatakbo nito.
- **Rename at remove.** Pinapanatili ng `ccam rename` ang pangalan ng shortcut, at gumagana pa rin ang link, dahil tumuturo ito sa ccam at hinahanap ng ccam ang pangalan sa `config.toml`. Binubura ng `ccam remove` ang link ng account. Binubura ng `ccam shortcut <account> <new>` ang lumang link bago gawin ang bago.
- **PATH.** Kapag gumawa ang ccam ng link at wala sa iyong `PATH` ang `~/.local/bin`, ipi-print nito ang `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` kasunod ang linyang `export PATH=...` na idadagdag. Sa Windows, ipi-print na lang nito ang PowerShell na linya na nagdadagdag ng folder sa user `PATH` mo.
- **Pagsusuri sa ibang pagkakataon.** Iniuulat ng `ccam doctor` ang `missing shortcuts in /home/me/.local/bin: ccwork`, ang mga shortcut na nagpapatakbo ng mas lumang bersyon ng ccam, at nagbababala kapag wala sa `PATH` ang folder. Ginagawa ng `ccam doctor --fix` ang mga nawawalang link at nire-renew ang mga luma na, at sinasabi lang nito ang `→ fixed` para sa mga ginawa nito. Gumagana ang shortcut na nagpapatakbo ng ibang ccam na naka-install pa, kaya inililista lang ito ng `ccam doctor` bilang impormasyon; itinuturo ito ng `ccam doctor --fix` sa ccam binary na ginamit mo sa pagpapatakbo nito, at ipinapakita ito ng `ccam add` at `ccam import` bilang ``shortcut: ccwork (runs another ccam, <path>; `ccam doctor --fix` points it at this one)``.

## Kaugnay

- [Pagsisimula ng Claude Code](/features/launch)
- [Pamamahala ng mga account](/features/accounts)
- [Sanggunian ng mga command](/docs/commands#shortcuts)
- [I-install ang ccam](/docs/install)
