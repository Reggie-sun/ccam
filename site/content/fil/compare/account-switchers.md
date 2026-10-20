---
title: Ang ccam at iba pang Claude Code account switcher
description: Paano ikinukumpara ang ccam sa claude-swap, CCS, clauth, CAAM, cux at iba pang tool sa pag-switch ng Claude Code account, at aling paraan ang bagay saan.
updated: 2026-10-20
order: 2
translation: machine
sourceHash: d058da937218
---

May dalawang paraan ang mga tool sa pag-switch ng Claude Code account. Karamihan ay nagpapalit ng login sa loob ng iisang `~/.claude` folder, kaya isang account lang ang aktibo sa bawat pagkakataon. Ang iba, kasama ang ccam, ay nagbibigay sa bawat account ng sariling folder gamit ang `CLAUDE_CONFIG_DIR`, kaya nananatiling naka-log in ang bawat account at puwedeng tumakbo nang sabay ang ilan. Nagdadagdag ang ccam ng shared na settings at MCP server, at itinutuloy nito ang isang usapan sa folder ng ibang account.

## Pagpapalit ng iisang login o magkakahiwalay na folder?

**Ang pagpapalit ng login** ay nagpapanatili ng iisang set ng settings, MCP server at kasaysayan, dahil iisa lang ang folder. Mabilis ang paglipat, at puwedeng gamitin ng tumatakbong session ang bagong login. Pero isang account lang ang aktibo sa bawat pagkakataon, para sa lahat ng terminal at editor nang sabay.

Para makapagpalit, nag-iingat ang tool ng kopya ng login ng bawat ibang account at isinusulat ito sa lugar nito. Pinapalitan ng Claude Code ang refresh token sa isang login tuwing nire-renew nito ito, kaya kailangang manatiling tugma ang bawat naka-save na kopya sa mga tumatakbong session. Kabilang sa mga problemang pinakamadalas ayusin ng mga project na ito ang mga account na nala-log out o kailangang mag-`/login` ulit pagkatapos lumipat.

**Ang magkakahiwalay na folder** ang paraang [inilalarawan sa dokumentasyon ng Claude Code](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Nananatiling naka-log in ang bawat account, at puwede mong patakbuhin nang magkatabi ang session para sa trabaho at ang personal na session. Nananatili ang bawat login sa sarili nitong folder, kung saan nire-renew ito ng Claude Code gaya ng dati. Ang kapalit nito: may sariling settings, MCP server at kasaysayan ang bawat folder sa simula, maliban kung may tool na nagshe-share ng mga ito.

## Paghahambing ng mga tool

Kinuha mula sa sariling dokumentasyon at source code ng bawat project noong 2026-10-20. Ang "Hindi nakadokumento" ay nangangahulugang hindi ito binabanggit ng project, hindi na wala ang feature.

| Tool | Paraan | Sabay-sabay na account | 5-oras at lingguhang usage | Kapag naabot ang limit | Platform |
| --- | --- | --- | --- | --- | --- |
| ccam | Magkakahiwalay na folder | Oo | Oo, `ccam usage` | Ikaw ang nagpapatakbo ng `ccam next --auto --continue` para ituloy ang usapan sa account na may pinakamaraming natitira; kusang inuulit ng `ccam --fallback` doon ang mga `-p` job | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | Nagpapalit ng login; magkakahiwalay na folder gamit ang `cswap run` | Oo, gamit ang `cswap run` | Oo | Nagpapalit ng login ang `cswap auto` sa isang threshold; tuloy ang tumatakbong session | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Magkakahiwalay na folder | Oo | Bahagya | Ipinapasa ng mga account pool ang mga request sa susunod na account sa pamamagitan ng lokal na proxy | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Nagpapalit ng login; magkakahiwalay na folder gamit ang `clauth start` | Oo, gamit ang `clauth start` | Oo | Lumilipat pababa sa isang fallback chain | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Nagpapalit ng login, o mga isolated na profile | Oo, gamit ang mga isolated na profile | Oo | Inuulit ng `caam run` ang mga non-interactive na command sa ibang account | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Nagpapalit ng login at binabalot ang `claude` | Hindi, isang aktibong account | Oo | Nagpapalit ng login at itinutuloy ang parehong usapan | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Nagpapalit ng login | Hindi | Oo | Opsyonal na paglipat sa isang threshold | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Nagpapalit ng login | Hindi | Oo | Opsyonal na paglipat sa 100% | macOS |

Ang ccswap ay fork ng claude-swap na namamahala rin ng mga Codex login. Nagbibigay rin ng sariling folder sa bawat account ang [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) at [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles). Ilan sa mga tool na ito ay humahawak din ng ibang coding agent, gaya ng Codex o Gemini.

## Paano ginawa ang mga ito

| Tool | Isinulat sa | Habang tumatakbo ang Claude Code | Kumokopya ng login sa pagitan ng mga lugar | Ano ang puwedeng i-share ng mga account |
| --- | --- | --- | --- | --- |
| ccam | Go, iisang static binary | Wala: pinapalitan ng ccam ang sarili nito ng `claude` (sa Windows, na hindi ito kaya, hinihintay nito ito) | Hindi | Bawat item nang hiwalay: lokal o isang may-pangalang grupo, bawat account |
| claude-swap | Python 3.12+, ini-install gamit ang uv o pipx | Wala kapag `cswap run`; patuloy na tumatakbo ang `cswap auto` at ang menu bar app | Oo, kapag nagpapalit ito ng login | Sa `cswap run`: settings, `CLAUDE.md` at extension nang magkakasama, opsyonal ang kasaysayan; kinokopya ang mga MCP server sa bawat launch |
| CCS | Node.js 18+ o Bun | Nananatili ang CCS bilang parent process; nagdadagdag ang mga pool ng lokal na proxy | Hindi para sa mga account folder nito; may sariling login ang mga pool | Shared o lokal bawat account, kasaysayan sa mga context group; kinokopya ang mga MCP server |
| clauth | Rust | Nananatili ang clauth bilang parent process | Oo, kapag nagpapalit ito ng login | Ang buong `~/.claude`, o wala |
| CAAM | Go | Depende sa mode; may opsyonal na daemon | Oo, sa vault mode nito | Itinakda ng mode; kinokopya ang mga MCP server sa bawat launch |
| cux | Go | May wrapper na nananatili sa paligid ng `claude` | Oo | Lahat: iisang folder |
| CCSwitcher, Claude Switcher | Mga macOS app (Swift, Python) | Ang menu bar app | Oo | Lahat: iisang folder |

## Kung saan magaling ang ccam

### Nananatili ang mga login kung saan inilagay ng Claude Code

Hindi kailanman kinokopya ng ccam ang isang Claude login mula sa isang lugar papunta sa iba. Nananatili ang login ng bawat account sa sarili nitong folder, o sa macOS, sa sariling Keychain item ng folder na iyon, at ang Claude Code ang nagre-renew nito. Walang naka-save na kopya na puwedeng mawala sa pagkakatugma. Ang dati nang `~/.claude` na inirehistro gamit ang `ccam import` ay nananatili sa Keychain item na ginagamit ng karaniwang `claude`, kaya naka-sign in pa rin ito alinmang paraan.

Nire-renew lang mismo ng ccam ang expired na token kapag pinatakbo mo ang `ccam usage --refresh`. Kumukuha muna ito ng lock at tumatanggi habang tumatakbo ang account na iyon, kaya hindi nito kailanman pinapalitan ang token habang may live na session.

### Walang tumatakbo sa pagitan mo at ng Claude Code

Nagbabasa ang ccam ng isang file, sinusuri ang mga link nito at pinapalitan ang sarili nito ng `claude` (`exec`). Habang nagtatrabaho ka, walang ccam process, wrapper, proxy o daemon: sa Claude Code ang standard input, mga signal at ang exit code, at diretso sa Anthropic ang mga request nito. Hindi kayang palitan ng Windows ang isang process, kaya doon sinisimulan ng ccam ang Claude Code sa parehong console at hinihintay ang exit code nito, nang walang ibang ginagawa.

Nagdadagdag lang ng mga 4 ms ang pagsisimula ng account sa ganitong paraan, mas mababa sa 6 ms sa 95th percentile, na may 8 MB na peak memory ([`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) gamit ang hyperfine, 300 run). Sa parehong laptop, 13 ms na ang pagsisimula pa lang ng Python interpreter at 31 ms ang Node.js, bago pa may gawin ang isang program.

### Mag-share bawat item, bawat account

Limang magkakahiwalay na item ang mga MCP server, `settings.json`, `CLAUDE.md`, mga extension (agent, command, skill, plugin) at ang kasaysayan ng usapan. Itinatakda ng bawat account ang bawat item sa `local` o sa isang grupo, kaya puwedeng mag-share ng kasaysayan ang dalawang account sa trabaho sa isang `work` na grupo habang may sariling kasaysayan ang personal na account, at iisang MCP server ang shine-share ng tatlo.

- Ipinapasa ang mga MCP server gamit ang `--mcp-config` sa launch, kaya bilang default ay hindi kailanman ine-edit ng ccam ang `.claude.json`, ang file na isinusulat muli ng Claude Code habang tumatakbo ito. Nagdadagdag ang `ccam mcp add` ng server nang isang beses para sa buong grupo.
- Ang ibang item ay mga link papunta sa folder ng grupo, kaya lumalabas agad sa iba ang pagbabagong ginawa sa isang account. Walang kinokopya sa launch.
- Hindi kailanman shine-share ang mga login, ang `.claude.json` at ang tumatakbong state.

### Ipinapakita muna at bina-backup ang bawat pagbabago

Nagpi-print ng plano ang `ccam share` at `ccam migrate` at nagtatanong muna bago maglipat ng anuman. Kapag magkaiba ang isang file sa account at sa grupo nito, ikaw ang pipili: panatilihin ang sa grupo, gamitin ang sa iyo, o i-merge ang JSON. Napupunta ang mga orihinal sa `.ccam-backup/<time>/`, napupunta sa isang trash folder ang folder ng inalis na account, at pinapalitan ang mga file sa isang rename, kaya naiiwan ang lumang file kapag naantala ang pagsusulat. Ipinapakita ng `--dry-run` kung ano ang gagawin ng launch, pag-share, migration, shortcut, update at pag-aayos. Inaayos ng `ccam doctor --fix` ang mga sirang link at permission.

### Ang parehong usapan sa susunod na account

Kapag shared ang kasaysayan, sabay-sabay na tinatanong ng `ccam next --auto --continue` ang bawat account tungkol sa usage nito, binubuksan ang may pinakamaraming natitira at itinutuloy ang usapan. Ikaw ang nagpapasya kung kailan lilipat; hindi binabantayan ng ccam ang usage mo sa background. Iminumungkahi ng status line ang command kapag umabot na sa 90% ang isang account.

Hindi ka kayang hintayin ng mga script, kaya ang `ccam --fallback work -p "…"` ang lumilipat para sa kanila: kapag naabot ng account ang limit nito, pinapatakbo ulit ang parehong command bilang account na may pinakamaraming natitira, at ang output lang ng run na natapos ang ipi-print. Ang pagkabigong hindi dahil sa usage limit ay ipinapasa nang gaya ng dati.

### Ang tamang account sa bawat folder

Ginagawa ng `ccam folder work ~/code/company` na gamitin ang `work` ng bawat launch na walang pangalan ng account sa loob ng folder na iyon, habang ginagamit pa rin ng iba mong project ang default na account. Nasa parehong `config.toml` na binabasa ng bawat launch ang setting na ito, kaya wala itong dagdag na oras.

### Iisang binary, sinusuring download

Ang ccam ay static binary na mga 10 MB at walang runtime na kailangang i-install, para sa Linux, macOS at Windows. Sinusuri ng mga install script (`install.sh`, `install.ps1`) at ng `ccam update` ang bawat download laban sa mga SHA-256 checksum ng release, at hindi kusang naghahanap ng update ang ccam. Wala itong telemetry: ang tanging mga request nito ay ang mga usage lookup sa Anthropic para sa mga account mo, nang sabay-sabay na may 3-segundong limit at 60-segundong cache, at ang mga download na hinihingi mo.

## Kung saan mas marami ang nagagawa ng ibang tool

- Kayang kusang ilipat ng claude-swap, clauth, cux, CAAM at CCSwitcher ang isang tumatakbong interactive session kapag lumampas ang isang account sa usage threshold. Inililipat ng ccam ang mga `-p` run gamit ang `--fallback` at iniiwan sa iyo ang mga interactive session.
- Sa Windows nang walang Developer Mode, nagshe-share ang ccam ng mga folder pero hindi ng mga iisang file gaya ng `settings.json`; kinokopya naman ang mga ito roon ng claude-swap at CCS.
- Ipinapakita ng CCSwitcher, Claude Switcher at claude-swap ang usage sa macOS menu bar.
- Namamahala rin ang CCS, CAAM, clauth, ccswap at Claude Switcher ng Codex at iba pang coding agent.
- Nai-install ang CAAM, CCSwitcher at Claude Switcher gamit ang Homebrew, at pinipirmahan ng clauth at CAAM ang mga release checksum nila. Nai-install ang ccam gamit ang script nito, isang `.deb` o `go install`, at sinusuri nito ang mga checksum nang walang pirma.

## Alin ang bagay sa iyo?

Pumili ng tool na nagpapalit ng login kung isang account lang ang ginagamit mo sa bawat pagkakataon o gusto mong awtomatikong lumipat ang mga tumatakbong session. Piliin ang ccam kung gusto mo ng mga account na magkatabi, iisang shared na setup, ang parehong usapan sa susunod na account at walang tumatakbo sa tabi ng Claude Code, sa Linux, macOS o Windows. Kung wala kang kailangang i-share, sapat na ang [simpleng mga shell alias](/compare/aliases).

Nakadepende sa terms ng bawat account kung puwede kang gumamit ng maraming account. Tingnan muna ang mga ito bago mo gawin.
