---
title: Mga madalas itanong
metaTitle: "ccam FAQ: maraming Claude Code account"
description: Mga sagot tungkol sa ccam, mga Claude Code account, usage limit, sharing, privacy at mga suportadong system.
updated: 2026-10-20
order: 6
translation: machine
sourceHash: e7d4a2c877ac
---

## Opisyal bang tool ng Anthropic ang ccam?

Hindi. Ang ccam ay hindi opisyal na open-source na proyekto at hindi kaanib o inendorso ng Anthropic. Ginagamit nito ang setting na `CLAUDE_CONFIG_DIR` na [idinodokumento ng Claude Code](https://code.claude.com/docs/en/env-vars) para sa pagpapatakbo ng maraming account nang magkatabi.

## Pinapayagan ba akong gumamit ng maraming Claude Code account?

Nakadepende iyan sa terms ng bawat account at, para sa mga work account, sa patakaran ng iyong organisasyon. Suriin muna ang mga ito bago ka gumamit ng maraming account. Pinapadali lang ng ccam ang paglipat sa pagitan ng mga account na mayroon ka na.

## Nalulusutan ba ng ccam ang mga usage limit?

Hindi. May sariling limit ang bawat account. Ipinapakita ng ccam ang mga ito at hinahayaan kang magpatuloy sa ibang account na pag-aari mo kapag naubos na ang isa.

## Paano gumamit ng dalawang Claude Code account sa iisang computer?

Bigyan ang bawat account ng sariling configuration folder gamit ang `CLAUDE_CONFIG_DIR`. Ginagawa ito ng ccam para sa iyo: patakbuhin ang `ccam add` para sa bawat account, saka simulan ang mga ito gamit ang `ccam personal` at `ccam work`, o gamit ang mga shortcut na ikaw mismo ang nagpangalan, gaya ng `ccwork`. Ipinapakita rin ng [gabay](/guides/multiple-claude-code-accounts) kung paano ito gawin nang mano-mano.

## Paano ko itutuloy ang isang usapan sa ibang account?

I-share ang kasaysayan ng usapan sa pagitan ng mga account gamit ang `ccam group set default --add history` (o sa isang group na para lang sa kanila: `ccam group add pair --items history --members work,side`), saka patakbuhin ang `ccam next --continue`. Bubuksan ng Claude Code ang pinakahuling usapan ng kasalukuyang folder sa susunod na account.

## Nakikita o ipinapadala ba ng ccam ang login ko?

Hindi kailanman ishe-share ng ccam ang iyong Claude login sa pagitan ng mga account. Binabasa lang nito ang access token ng isang account para itanong sa usage endpoint ng Anthropic ang usage ng account na iyon, kapag pinatakbo mo ang `ccam usage`, `ccam next --auto` o `ccam doctor`. Kinokopya ng `ccam migrate` ang mga naka-save na sign-in ng mga MCP server mula sa `~/.claude` papunta sa sariling login file ng iyong mga account sa parehong computer. Walang telemetry ang ccam at wala itong ipinapadala saanman bukod doon.

## Paano mag-sign in sa isang remote server?

Patakbuhin ang `ccam add` o `ccam login <name>` sa SSH session gaya ng dati. Naghihintay ang Claude Code ng sign-in sa server, kaya mapupunta ang browser sa computer mo sa isang `http://localhost:<port>/callback?code=…` na page na hindi naglo-load. Kopyahin ang address na iyon mula sa browser at i-paste sa terminal: ipapasa ito ng ccam sa Claude Code sa server, na siyang tatapos sa sign-in. May one-time code ang address, kaya huwag itong i-share.

## Saan naka-store ang mga account ko?

Sa `~/.local/share/ccam/profiles/<name>/`, isang folder bawat account. Sa macOS, iniimbak mismo ng Claude Code ang login sa Keychain, sa ilalim ng pangalang hango sa folder na iyon.

## Ano ang mangyayari sa dati kong ~/.claude?

Wala, maliban kung patakbuhin mo ang `ccam import` o `ccam migrate`. Nirerehistro ng import ang folder kung nasaan ito, kaya ginagamit pa rin ng karaniwang `claude` at ng `ccam <name>` ang parehong account at settings. Kinokopya ng `ccam migrate` ang settings, `CLAUDE.md`, mga extension at MCP server nito sa shared na settings, at binabasa lang ang folder.

## Aling mga system ang suportado ng ccam?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine at iba pa), macOS at Windows 10 at 11, direkta o sa loob ng WSL. Sa Windows, kailangan ng Developer Mode para mag-share ng mga file gaya ng `settings.json`; gumagana ang mga shared na folder kahit wala ito ([sharing sa Windows](/features/share#on-windows)).

## Bakit sa iisang account lang lumalabas ang isang MCP server?

Ang mga server na idinagdag sa loob ng session gamit ang `claude mcp add` ay naka-store para sa account na iyon lang. Sasabihin ito sa iyo ng ccam sa susunod na simulan mo ang account; ilipat ang server sa shared group gamit ang `ccam mcp import <account> <server>`.

## Pinababagal ba ng ccam ang Claude Code?

Nagdadagdag ng mga 4 millisecond ang pagsisimula ng Claude Code sa pamamagitan ng ccam, na sinukat gamit ang `scripts/perf-check.sh` sa isang laptop. Nagbabasa ang ccam ng isang maliit na file, sinusuri ang mga shared na link at saka pinapalitan ang sarili nito ng `claude`, kaya hindi ito tumatakbo habang nagtatrabaho ka. Hindi kayang palitan ng Windows ang isang process, kaya doon sinisimulan ng ccam ang Claude Code at hinihintay ito, nang walang ibang ginagawa.

## Paano i-update ang ccam?

Patakbuhin ang `ccam update`. Dina-download nito ang pinakabagong release, sinusuri ang checksum nito at pinapalitan ang `ccam` file; sinasabi lang ng `ccam update --check` kung may mas bagong release. Hindi kailanman kusang nag-a-update ang ccam. Kung in-install mo ang `.deb` package, i-install na lang ang bagong `.deb` gamit ang apt. Tingnan ang [pag-update](/docs/install#update).

## Paano i-uninstall ang ccam?

Alisin ang shortcut ng bawat account gamit ang `ccam shortcut <account> --remove`, burahin ang `ccam` mula sa `~/.local/bin`, saka alisin ang `~/.config/ccam`, `~/.local/share/ccam` at `~/.local/state/ccam` kung hindi mo na kailangan ang mga account. Tingnan ang [install](/docs/install#uninstall).
