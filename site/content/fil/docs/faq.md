---
title: Mga madalas itanong
metaTitle: "ccam FAQ: maraming Claude Code account"
description: Mga sagot tungkol sa ccam, mga Claude Code account, usage limit, sharing, privacy at mga suportadong system.
updated: 2026-10-19
order: 6
translation: machine
sourceHash: 8876e4f00bd2
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

I-share ang kasaysayan ng usapan sa pagitan ng dalawang account gamit ang `ccam share <name> --history on`, saka patakbuhin ang `ccam next --continue`. Bubuksan ng Claude Code ang pinakahuling usapan ng kasalukuyang folder sa susunod na account.

## Nakikita o ipinapadala ba ng ccam ang login ko?

Hindi kailanman kumokopya o nagshe-share ng login ang ccam. Binabasa lang nito ang access token ng isang account para itanong sa usage endpoint ng Anthropic ang usage ng account na iyon, kapag pinatakbo mo ang `ccam usage`, `ccam next --auto` o `ccam doctor`. Walang telemetry ang ccam at wala itong ipinapadala saanman bukod doon.

## Saan naka-store ang mga account ko?

Sa `~/.local/share/ccam/profiles/<name>/`, isang folder bawat account. Sa macOS, iniimbak mismo ng Claude Code ang login sa Keychain, sa ilalim ng pangalang hango sa folder na iyon.

## Ano ang mangyayari sa dati kong ~/.claude?

Wala, maliban kung patakbuhin mo ang `ccam import` o `ccam migrate`. Nirerehistro ng import ang folder kung nasaan ito, kaya ginagamit pa rin ng karaniwang `claude` at ng `ccam <name>` ang parehong account at settings. Kinokopya ng `ccam migrate` ang settings, `CLAUDE.md`, mga extension at MCP server nito sa shared na settings, at binabasa lang ang folder.

## Aling mga system ang suportado ng ccam?

Linux (Ubuntu, Debian, Fedora, Arch, Alpine at iba pa) at macOS. Sa Windows, gamitin ito sa loob ng WSL.

## Bakit sa iisang account lang lumalabas ang isang MCP server?

Ang mga server na idinagdag sa loob ng session gamit ang `claude mcp add` ay naka-store para sa account na iyon lang. Sasabihin ito sa iyo ng ccam sa susunod na simulan mo ang account; ilipat ang server sa shared group gamit ang `ccam mcp import <account> <server>`.

## Pinababagal ba ng ccam ang Claude Code?

Nagdadagdag ng mga 4 millisecond ang pagsisimula ng Claude Code sa pamamagitan ng ccam, na sinukat gamit ang `scripts/perf-check.sh` sa isang laptop. Nagbabasa ang ccam ng isang maliit na file, sinusuri ang mga shared na link at saka pinapalitan ang sarili nito ng `claude`, kaya hindi ito tumatakbo habang nagtatrabaho ka.

## Paano i-update ang ccam?

Patakbuhin ang `ccam update`. Dina-download nito ang pinakabagong release, sinusuri ang checksum nito at pinapalitan ang `ccam` file; sinasabi lang ng `ccam update --check` kung may mas bagong release. Hindi kailanman kusang nag-a-update ang ccam. Kung in-install mo ang `.deb` package, i-install na lang ang bagong `.deb` gamit ang apt. Tingnan ang [pag-update](/docs/install#update).

## Paano i-uninstall ang ccam?

Alisin ang shortcut ng bawat account gamit ang `ccam shortcut <account> --remove`, burahin ang `ccam` mula sa `~/.local/bin`, saka alisin ang `~/.config/ccam`, `~/.local/share/ccam` at `~/.local/state/ccam` kung hindi mo na kailangan ang mga account. Tingnan ang [install](/docs/install#uninstall).
