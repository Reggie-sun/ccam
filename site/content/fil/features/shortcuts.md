---
title: Mga shortcut command
description: Bigyan ang isang account ng sariling command gaya ng ccwork gamit ang ccam shortcut, at gawin muli ang mga link sa ~/.local/bin gamit ang ccam shims install.
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 5f0cb5a57042
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
ccam shims install            # recreate missing links of all accounts
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
| `--dry-run` | Sinusunod ito ng `ccam shortcut`: sinusuri nito ang pangalan at ipi-print ang `Would set the shortcut of work to ccwork.` o `Would remove ccwork, the shortcut of work.` nang hindi nagse-save o nagli-link | off |

Walang tinatanggap na argument ang `ccam shims install` at tinatanggihan nito ang `--dry-run` na may `` ccam: --dry-run is not supported by `ccam shims install`; nothing was changed ``.

Mga patakaran sa pangalan, na sinusuri sa ganitong pagkakasunod-sunod:

- 1 hanggang 64 na character ng `a-z`, `0-9`, `.`, `-` at `_`, na nagsisimula sa titik o numero;
- hindi `ccam`, `claude` o ccam command gaya ng `run` o `usage`;
- hindi shortcut ng ibang account;
- hindi file sa `~/.local/bin` na hindi pa link sa ccam;
- hindi pangalan ng ibang executable sa iyong `PATH`.

Kapag sabay na ibinigay ang command at `--remove`, humihinto ito na may `give a command or --remove, not both`. Kapag itinakda ang pangalang mayroon na ang account, ipi-print ang `work: ccwork (unchanged)` at gagawin muli ang link nito kung nawawala.

## Paano ito gumagana

Ang pagtatakda ng pangalan ay nagpapatakbo ng `checkShortcut`, na tumatawag sa `shim.Check` kasama ang bin folder at ang totoong path ng ccam binary:

<!-- motion:shortcuts-check -->

Bumababa ang pangalan sa mga pagsusuri nang sunod-sunod. Katabi ng bawat pagsusuri ang isang pangalang tinatanggihan nito: pumapalya ang `Work` sa pattern, ccam command ang `run`, pag-aari ng ibang account ang `ccside`, umiiral na file sa `~/.local/bin` ang `ccfile`, at itatago ng `git` ang `/usr/bin/git`. Ang pangalan lang na pumasa sa bawat pagsusuri ang sine-save at nili-link.

1. Itinutugma ng `shim.Check` ang pangalan sa `^[a-z0-9][a-z0-9._-]{0,63}$`, saka sa `config.Reserved`, saka sa `cfg.FindShortcut`. Tinatawag nito ang `os.Lstat` sa `~/.local/bin/<name>`, at para sa bawat ibang folder sa `PATH` ay tinatawag nito ang `os.Stat` at tinatanggihan ang executable file na hindi link sa ccam.
2. Sine-save ng ccam ang `config.toml` gamit ang `fsutil.AtomicWrite` (pansamantalang file, `fsync`, rename).
3. Binubura ng `shim.Remove` ang link ng lumang shortcut, pero kapag link lang ito sa ccam.
4. Ginagawa ng `shim.Install` ang link ng bawat naka-configure na shortcut na nawawala, gamit ang `os.Symlink` papunta sa ccam binary, na hinahanap gamit ang `os.Executable` at `filepath.EvalSymlinks`. Hindi ito kailanman nag-aalis ng anuman.

Itinuturing na sariling link ng ccam ang isang link kapag ibinabalik ng `os.Readlink` ang path ng ccam, o kapag tumuturo sa parehong file ang link at ang path ng ccam. Dumadaan sa pagsusuring ito ang bawat pag-alis at bawat pagsusuri bago mag-overwrite, kaya hindi ginagalaw ng ccam ang mga file na hindi nito ginawa.

Kapag nagpatakbo ka ng shortcut, ang pangalang tina-type mo lang ang kailangan ng ccam:

<!-- motion:shortcuts-link -->

Hinahanap ng shell ang `ccwork` sa `PATH`, sinusundan ng kernel ang link at sinisimulan ang ccam binary, at `ccwork` pa rin ang `argv[0]`. Sa `app.Main`, nilo-load nang isang beses ng anumang pangalan maliban sa `ccam` ang `config.toml` at tinatawag ang `shim.Lookup`, na hinahanap ang account na ang `shortcut` ay katumbas ng `filepath.Base(argv[0])`. Pagkatapos ay ginagawa ng `launcher.LaunchWith` ang parehong mga hakbang ng [`ccam <name>`](/features/launch) na naka-load na ang config, at ipinapasa ang `argv[1:]` sa `claude` nang hindi ginagalaw.

Bakit symlink at hindi shell alias o wrapper script: gumagana ang symlink sa bawat shell at sa mga script, hindi ito nagdadagdag ng karagdagang pagsisimula ng shell (mga 1 hanggang 2 ms), at iisang file pa rin ang binabasa ng pagsisimula. Sinusukat ng `scripts/perf-check.sh` ang isang shortcut gamit ang `hyperfine` laban sa parehong 10 ms (p95) na limit ng `ccam <name>`.

## Kaligtasan at mga espesyal na kaso

- **Link na walang account na umaangkin.** Kung may link sa ccam sa `~/.local/bin` na may pangalang hindi ginagamit ng anumang account, halimbawa pagkatapos mong mano-manong baguhin ang `shortcut` sa `config.toml`, hindi babalik ang ccam sa default na account. Humihinto ito na may `` ccam: ccold is not the shortcut of any account (see `ccam list`, `ccam shortcut`) ``. Itakda ulit ang pangalan gamit ang `ccam shortcut`, o burahin ang link.
- **Mga umiiral na file.** Hindi kailanman pinapalitan ng ccam ang file na hindi nito ginawa. Humihinto ang `ccam shims install` na may `<path> already exists and is not a ccam link, so it was left alone`.
- **Inilipat na ccam binary.** Tumuturo na ang mga lumang link sa path na wala na, at hindi itinuturing ng ccam na sarili nito ang mga ito. Burahin ang mga ito, saka patakbuhin ang `ccam shims install`.
- **Rename at remove.** Pinapanatili ng `ccam rename` ang pangalan ng shortcut, at gumagana pa rin ang link, dahil tumuturo ito sa ccam at hinahanap ng ccam ang pangalan sa `config.toml`. Binubura ng `ccam remove` ang link ng account. Binubura ng `ccam shortcut <account> <new>` ang lumang link bago gawin ang bago.
- **PATH.** Kapag gumawa ang ccam ng link at wala sa iyong `PATH` ang `~/.local/bin`, ipi-print nito ang `Note: /home/me/.local/bin is not in PATH. Add it to ~/.bashrc or ~/.zshrc:` kasunod ang linyang `export PATH=...` na idadagdag.
- **Pagsusuri sa ibang pagkakataon.** Iniuulat ng `ccam doctor` ang `missing shortcuts in /home/me/.local/bin: ccwork` at nagbababala kapag wala sa `PATH` ang folder. Ginagawa ng `ccam doctor --fix` ang mga nawawalang link.

## Kaugnay

- [Pagsisimula ng Claude Code](/features/launch)
- [Pamamahala ng mga account](/features/accounts)
- [Sanggunian ng mga command](/docs/commands#shortcuts)
- [I-install ang ccam](/docs/install)
